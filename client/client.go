package main

/*
 * A client does the following:
 * Read the config to get public key and IP maps
 * Let B be the number of commands.
 * Send B commands to the nodes and wait for f+1 acknowledgements for every acknowledgement
 */

import (
	"bufio"
	"context"
	"encoding/binary"
	"sync"
	"time"

	"github.com/adithyabhatkajake/libapollo/config"
	"github.com/adithyabhatkajake/libchatter/log"

	"github.com/adithyabhatkajake/libapollo/consensus"
	"github.com/adithyabhatkajake/libapollo/msg"
	"github.com/adithyabhatkajake/libchatter/crypto"
	"github.com/adithyabhatkajake/libchatter/io"

	pb "github.com/golang/protobuf/proto"
	"github.com/libp2p/go-libp2p"
	p2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
)

var (
	// BufferCommands defines how many commands to wait for
	// acknowledgement in a batch
	BufferCommands = uint64(10)
	// PendingCommands tells how many commands we are waiting
	// for acknowledgements from replicas
	PendingCommands = uint64(0)
	cmdMutex        = &sync.Mutex{}
	streamMutex     = &sync.Mutex{}
	voteMutex       = &sync.Mutex{}
	condLock        = &sync.RWMutex{}
	voteChannel     chan *msg.CommitAck
	idMap           = make(map[string]uint64)
	votes           = make(map[crypto.Hash]uint64)
	timeMap         = make(map[crypto.Hash]time.Time)
	f               uint64
	rwMap           = make(map[uint64]*bufio.Writer)
)

func sendCommandToServer(cmd *msg.ApolloMsg) {
	log.Trace("Sending a Command to the server")
	cmdHash := crypto.DoHash(cmd.GetTx())
	data, err := pb.Marshal(cmd)
	if err != nil {
		log.Error("Marshaling error", err)
		return
	}
	condLock.Lock()
	timeMap[cmdHash] = time.Now()
	condLock.Unlock()

	// Ship command off to the nodes
	streamMutex.Lock()
	for idx, rw := range rwMap {
		rw.Write(data)
		rw.Flush()
		log.Trace("Sending command to node", idx)
	}
	streamMutex.Unlock()
}

func ackMsgHandler(s network.Stream, from uint64) {
	reader := bufio.NewReader(s)
	// Prepare a buffer to receive an acknowledgement
	msgBuf := make([]byte, msg.MaxMsgSize)
	log.Trace("Started acknowledgement message handler")
	for {
		len, err := reader.Read(msgBuf)
		if err != nil {
			log.Error("bufio read error", err)
			return
		}
		log.Trace("Received a message from the server ", from)
		msg := &msg.ApolloMsg{}
		err = pb.Unmarshal(msgBuf[0:len], msg)
		if err != nil {
			log.Error("Unmarshalling error", err)
			continue
		}
		voteChannel <- msg.GetAck()
	}
}

func handleVotes(cmdChannel chan *msg.ApolloMsg) {
	voteMap := make(map[crypto.Hash]uint64)
	commitMap := make(map[crypto.Hash]bool)
	for {
		// Get Acknowledgements from nodes after consensus
		ack, ok := <-voteChannel
		timeReceived := time.Now()
		log.Trace("Received an acknowledgement")
		if !ok {
			log.Error("vote channel closed")
			return
		}
		if ack == nil {
			continue
		}
		bhash := ack.GetBlock().GetBlockHash()
		_, exists := voteMap[bhash]
		if !exists {
			voteMap[bhash] = 1       // 1 means we have seen one vote so far.
			commitMap[bhash] = false // To say that we have not yet committed this value
		} else {
			voteMap[bhash]++ // Add another vote
		}
		// To ensure this is executed only once, check old committed state
		old := commitMap[bhash]
		if voteMap[bhash] <= f {
			// Not enough votes for this block
			// So this is not yet committed
			// Deal with it later
			continue
		}
		commitMap[bhash] = true
		new := commitMap[bhash]
		log.Trace("Committed block. Processing block",
			ack.GetBlock().GetHeader().GetHeight())
		sendNewCommands := old != new
		txs := ack.GetBlock().GetBody().TxHashes
		for _, tx := range txs {
			cmdHash := crypto.ToHash(tx)
			condLock.Lock()
			commitTimeMetric[cmdHash] = timeReceived.Sub(timeMap[cmdHash])
			// Time from sending to getting back in some block
			condLock.Unlock()
			// If we commit the block for the first time, then ship off a new command to the server
			if sendNewCommands { // will be triggered once when commitMap value changes
				cmd := <-cmdChannel
				// log.Info("Sending command ", cmd, " to the servers")
				go sendCommandToServer(cmd)
			}
		}
	}
}

func main() {
	log.Info("I am the client")
	ctx := context.Background()

	// Set values based on command line arguments
	ParseOptions()

	// Get client config
	confData := &config.ClientConfig{}
	io.ReadFromFile(confData, *ConfFile)

	f = confData.GetFaults()
	// Start networking stack
	node, err := p2p.New(ctx,
		libp2p.Identity(confData.GetMyKey()),
	)
	if err != nil {
		panic(err)
	}

	// Print self information
	log.Info("Client at", node.Addrs())

	// Handle all messages received using ackMsgHandler
	// node.SetStreamHandler(e2cconsensus.ClientProtocolID, ackMsgHandler)
	// Setting stream handler is useless :/

	pMap := make(map[uint64]peer.AddrInfo)
	streamMap := make(map[uint64]network.Stream)
	connectedNodes := uint64(0)
	wg := &sync.WaitGroup{}
	updateLock := &sync.Mutex{}

	for i := uint64(0); i < confData.GetNumNodes(); i++ {
		wg.Add(1)
		go func(i uint64, peer peer.AddrInfo) {
			defer wg.Done()
			// Connect to node i
			log.Trace("Attempting connection to node ", peer)
			err := node.Connect(ctx, peer)
			if err != nil {
				log.Error("Connection Error ", err)
				return
			}
			for {
				stream, err := node.NewStream(ctx, peer.ID,
					consensus.ClientProtocolID)
				if err != nil {
					log.Trace("Stream opening Error-", err)
					<-time.After(10 * time.Millisecond)
					continue
				}
				updateLock.Lock()
				defer updateLock.Unlock()
				streamMap[i] = stream
				pMap[i] = peer
				idMap[stream.ID()] = i
				connectedNodes++
				rwMap[i] = bufio.NewWriter(stream)
				go ackMsgHandler(stream, i)
				break
			}
			log.Debug("Successfully connected to node ", i)
		}(i, confData.GetPeerFromID(i))
	}
	wg.Wait()

	// Ensure we are connected to sufficient nodes
	if connectedNodes <= f {
		log.Warn("Insufficient connections to replicas")
		return
	}

	blksize := confData.GetBlockSize()
	initialSendCmds := 10 * (f + 1) * blksize

	cmdChannel := make(chan *msg.ApolloMsg, BufferCommands)
	voteChannel = make(chan *msg.CommitAck, blksize)

	// First, spawn a thread that handles acknowledgement received for the
	// various requests
	go handleVotes(cmdChannel)

	idx := uint64(0)

	log.Info("Sending initial batch of ", initialSendCmds, " commands")
	// Then, run a goroutine that sends the first Blocksize requests to the nodes
	for ; idx < initialSendCmds; idx++ {
		// Build a command
		cmd := make([]byte, 8+*Payload)
		binary.LittleEndian.PutUint64(cmd, idx)

		// Build a protocol message
		cmdMsg := &msg.ApolloMsg{}
		cmdMsg.Msg = &msg.ApolloMsg_Tx{
			Tx: cmd,
		}

		// log.Info("Sending command ", idx, " to the servers")
		sendCommandToServer(cmdMsg)
	}
	log.Info("Finished sending initial batch of commands", idx)

	go printMetrics()

	// Make sure we always fill the channel with commands
	for {
		// Build a command
		cmd := make([]byte, 8+*Payload)

		// Make every command unique so that the hashes are unique
		binary.LittleEndian.PutUint64(cmd, idx)

		// Build a protocol message
		cmdMsg := &msg.ApolloMsg{}
		cmdMsg.Msg = &msg.ApolloMsg_Tx{
			Tx: cmd,
		}

		// Dispatch message for processing
		// This will block until some of the commands are committed
		cmdChannel <- cmdMsg
		// Increment command number
		idx++
	}
}
