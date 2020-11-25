package consensus

import (
	"bufio"
	"context"
	"sync"

	chain "github.com/adithyabhatkajake/libapollo/chain"
	config "github.com/adithyabhatkajake/libapollo/config"
	"github.com/adithyabhatkajake/libchatter/crypto"

	"github.com/libp2p/go-libp2p-core/host"
	peerstore "github.com/libp2p/go-libp2p-core/peer"
)

// Apollo implements the consensus protocol
type Apollo struct {
	// Network data structures
	host    host.Host
	cliHost host.Host
	ctx     context.Context

	// Maps
	// Mapping between ID and libp2p-peer
	pMap map[uint64]*peerstore.AddrInfo
	// A map of node ID to its corresponding RW stream
	streamMap map[uint64]*bufio.ReadWriter

	/* Locks - We separate all the locks, so that acquiring
	one lock does not make other goroutines stop */
	netMutex    sync.RWMutex // The lock to modify streamMap: Use mutex when using network streams to talk to other nodes
	blTimerLock sync.RWMutex // The lock to modify blTimer

	// Channels
	NewTxCh        chan chain.Command
	CommitNotifyCh chan chain.Block
	PoolFull       chan ProposeReady
	TxExtractCh    chan TxCleave
	ExtractBlk     chan []crypto.Hash
	NewBlockNotify chan chain.Block

	msgChannel chan internalMsg // All messages come here first
	errCh      chan error       // All errors are sent here

	// ClientManager
	clim *ClientManager

	// BlockChain
	bc     *chain.BlockChain
	view   uint64
	pmaker *RRPaceMaker

	// Embed the config
	*config.NodeConfig
}

// Signals used internally
type (
	// ProposeReady tells that we are ready to propose
	ProposeReady struct{}
	// TxCleave lets the tx handler know that we want to extract a block
	TxCleave struct{}
)
