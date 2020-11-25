package consensus

import (
	"bufio"
	"context"
	"sync"
	"time"

	"github.com/adithyabhatkajake/libapollo/chain"
	"github.com/adithyabhatkajake/libchatter/crypto"
	"github.com/adithyabhatkajake/libchatter/log"

	"github.com/libp2p/go-libp2p"

	config "github.com/adithyabhatkajake/libapollo/config"
	"github.com/adithyabhatkajake/libchatter/net"

	msg "github.com/adithyabhatkajake/libapollo/msg"

	"github.com/libp2p/go-libp2p-core/network"
	peerstore "github.com/libp2p/go-libp2p-core/peer"
)

const (
	// ProtocolID is the ID for E2C Protocol
	ProtocolID = "apollo/proto/0.0.1"
	// ProtocolMsgBuffer defines how many protocol messages can be buffered
	ProtocolMsgBuffer = 100
)

// Init implements the Protocol interface
func (n *Apollo) Init(c *config.NodeConfig) {
	n.NodeConfig = c

	// Setup maps
	n.streamMap = make(map[uint64]*bufio.ReadWriter)

	// Setup channels
	n.NewTxCh = make(chan chain.Command, 5*n.GetBlockSize())
	n.CommitNotifyCh = make(chan chain.Block, 1000)
	n.PoolFull = make(chan ProposeReady, 1000)
	n.errCh = make(chan error, 1)
	n.msgChannel = make(chan internalMsg, ProtocolMsgBuffer)
	n.TxExtractCh = make(chan TxCleave, 1000)
	n.ExtractBlk = make(chan []crypto.Hash, 1000)
	n.NewBlockNotify = make(chan chain.Block, 1000)

	n.clim = NewClientManager(n)
	// Setup genesis
	n.bc = chain.NewChain()
	n.pmaker = &RRPaceMaker{
		currentLeader: DefaultLeaderID,
		lastBlock:     n.bc.Head,
		numNodes:      n.GetNumNodes(),
	}
	n.view = 1
}

// Setup sets up the network components
func (apl *Apollo) Setup(n *net.Network) error {
	apl.host = n.H
	host, err := libp2p.New(
		context.Background(),
		libp2p.ListenAddrStrings(apl.GetClientListenAddr()),
		libp2p.Identity(apl.GetMyKey()),
	)
	if err != nil {
		panic(err)
	}
	apl.pMap = n.PeerMap
	apl.cliHost = host
	apl.ctx = n.Ctx

	// Obtain a new chain
	apl.bc = chain.NewChain()
	// TODO: create a new chain only if no chain is present in the data directory

	// How to react to Protocol Messages
	apl.host.SetStreamHandler(ProtocolID, apl.ProtoMsgHandler)

	// How to react to Client Messages
	apl.cliHost.SetStreamHandler(ClientProtocolID, apl.clim.ClientMsgHandler)

	// Connect to all the other nodes talking E2C protocol
	wg := &sync.WaitGroup{} // For faster setup
	for idx, p := range apl.pMap {
		wg.Add(1)
		go func(idx uint64, p *peerstore.AddrInfo) {
			log.Trace("Attempting to open a stream with", p, "using protocol", ProtocolID)
			retries := 300
			for i := retries; i > 0; i-- {
				s, err := apl.host.NewStream(apl.ctx, p.ID, ProtocolID)
				if err != nil {
					log.Error("Error connecting to peers:", err)
					log.Info("Retry attempt ", retries-i+1, " to connect to node ", idx, " in a second")
					<-time.After(10 * time.Millisecond)
					continue
				}
				apl.netMutex.Lock()
				apl.streamMap[idx] = bufio.NewReadWriter(
					bufio.NewReader(s), bufio.NewWriter(s))
				apl.netMutex.Unlock()
				log.Info("Connected to Node ", idx)
				break
			}
			wg.Done()
		}(idx, p)
	}
	wg.Wait()
	log.Info("Setup Finished. Ready to do SMR:)")

	return nil
}

// Start implements the Protocol Interface
func (n *Apollo) Start() {
	// Start handling the clients
	go n.clim.ClientResponder()
	go n.clim.ClientHandler()
	// Start E2C Protocol - Start message handler
	n.protocol()
}

// ProtoMsgHandler reacts to all protocol messages in the network
func (n *Apollo) ProtoMsgHandler(s network.Stream) {
	// A global buffer to collect messages
	buf := make([]byte, msg.MaxMsgSize)
	// Event Handler
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	for {
		// Receive a message from anyone and process them
		len, err := rw.Read(buf)
		if err != nil {
			return
		}
		// Use a copy of the message and send it to off for processing
		msgBuf := make([]byte, len)
		copy(msgBuf, buf[0:len])
		// React to the message in parallel and continue
		n.react(msgBuf, rw)
	}
}
