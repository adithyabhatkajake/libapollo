package consensus

import (
	"bufio"
	"sync"

	"github.com/adithyabhatkajake/libapollo/txpool"

	"github.com/adithyabhatkajake/libapollo/chain"
	"github.com/adithyabhatkajake/libapollo/msg"
	"github.com/adithyabhatkajake/libchatter/crypto"
	"github.com/adithyabhatkajake/libchatter/log"
	pb "github.com/golang/protobuf/proto"
	"github.com/libp2p/go-libp2p-core/network"
)

// Implement how to talk to clients
const (
	ClientProtocolID = "apollo/client/0.0.1"
)

// ClientManager manages client side protocol
type ClientManager struct {
	cliMap  map[*bufio.ReadWriter]bool
	blkSize uint64

	newTxCh        chan chain.Command
	commitNotifyCh chan chain.Block
	errCh          chan error
	poolFull       chan ProposeReady
	responseCh     chan chain.Block
	cleave         chan TxCleave
	giveBlk        chan []crypto.Hash
	newIncomingBlk chan chain.Block

	cliMutex sync.RWMutex
}

// NewClientManager returns a new initialized client manager
func NewClientManager(n *Apollo) *ClientManager {
	clim := &ClientManager{
		cliMap:         make(map[*bufio.ReadWriter]bool),
		blkSize:        n.GetBlockSize(),
		newTxCh:        n.NewTxCh,
		commitNotifyCh: n.CommitNotifyCh,
		errCh:          n.errCh,
		poolFull:       n.PoolFull,
		cleave:         n.TxExtractCh,
		giveBlk:        n.ExtractBlk,
		responseCh:     make(chan chain.Block, len(n.CommitNotifyCh)),
		newIncomingBlk: n.NewBlockNotify,
	}
	log.Debug("Using block size:", clim.blkSize)
	return clim
}

// AddClient adds a new client to cliMap
func (clim *ClientManager) AddClient(rw *bufio.ReadWriter) {
	clim.cliMap[rw] = true
}

// RemoveClient removes rw from cliMap after disconnection
func (clim *ClientManager) RemoveClient(rw *bufio.ReadWriter) {
	delete(clim.cliMap, rw)
}

// ClientHandler manages the clients
func (clim *ClientManager) ClientHandler() {
	// Handle transactions here
	txm := txpool.NewTxManager()
	poolFull := false
	// TODO: Only notify once when the pool becomes full
	for {
		select {
		case tx := <-clim.newTxCh:
			if tx == nil {
				continue
			}
			log.Debug("Received a transaction:", tx)
			// Process new transaction
			txm.AddCommand(tx) // Add command
			if numCmds := txm.Size(); numCmds >= clim.blkSize && !poolFull {
				// Let the proposer know that you can propose
				clim.poolFull <- ProposeReady{}
				poolFull = true // Prevent telling again that the pool is full
			}
		case <-clim.cleave:
			// Someone wants to propose, and is requesting a block
			if numCmds := txm.Size(); numCmds >= clim.blkSize {
				newCandBlk := txm.GetBlock(clim.blkSize)
				poolFull = txm.Size() >= clim.blkSize
				clim.giveBlk <- newCandBlk
			} else {
				poolFull = false
				clim.giveBlk <- nil
			}
		case blk := <-clim.newIncomingBlk:
			log.Debug("We have a new block. Update the tx pool")
			txm.Clear(blk.GetTxs())
			poolFull = txm.Size() >= clim.blkSize
		case blk := <-clim.commitNotifyCh:
			log.Debug("Committing block:", blk)
			// Send committed block to the clients
			clim.responseCh <- blk
		case err, ok := <-clim.errCh:
			log.Warn("Received an error.", err)
			if !ok {
				return
			}
		}
	}
}

// ClientResponder lets clients know that we have committed blocks
func (clim *ClientManager) ClientResponder() {
	for {
		select {
		case blk := <-clim.responseCh:
			log.Debug("Letting the clients know that we have a new block", blk.GetHeight())
			m := &msg.ApolloMsg{
				Msg: &msg.ApolloMsg_Ack{
					Ack: &msg.CommitAck{
						Block: blk.ToProto(),
					},
				},
			}
			clim.clientBroadcast(m)
		case err, ok := <-clim.errCh:
			log.Warn("Received an error.", err)
			if !ok {
				return
			}
		}
	}
}

// ClientMsgHandler defines how to talk to client messages
func (clim *ClientManager) ClientMsgHandler(s network.Stream) {
	// A buffer to collect messages
	buf := make([]byte, msg.MaxMsgSize)
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	clim.AddClient(rw) // Add client for later contact
	// Event Handler
	for {
		// Receive a message from a client and process them
		len, err := rw.Read(buf)
		if err != nil {
			log.Error("Error receiving a message from the client-", err)
			clim.RemoveClient(rw)
			return
		}
		// Send a copy for reacting
		inMsg := &msg.ApolloMsg{}
		err = pb.Unmarshal(buf[0:len], inMsg)
		if err != nil {
			log.Error("Error unmarshalling cmd from client")
			log.Error(err)
			continue
		}
		var cmd chain.Command
		if cmd = inMsg.GetTx(); cmd == nil {
			log.Error("Invalid command received from client")
			continue
		}
		clim.newTxCh <- cmd
	}
}

// ClientBroadcast sends a protocol message to all the clients known to this instance
func (clim *ClientManager) clientBroadcast(m *msg.ApolloMsg) {
	data, err := pb.Marshal(m)
	if err != nil {
		log.Error("Failed to send message", m, "to client")
		return
	}

	for cliBuf := range clim.cliMap {
		log.Trace("Sending to", cliBuf)
		cliBuf.Write(data)
		cliBuf.Flush()
	}
	log.Trace("Finish client broadcast")
}
