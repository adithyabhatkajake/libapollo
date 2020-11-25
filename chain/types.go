package chain

import (
	"sync"

	"github.com/adithyabhatkajake/libchatter/crypto"
)

// BlockChain is what we call a blockchain
type BlockChain struct {
	// Just a store of blocks
	Storage map[crypto.Hash]Block
	// A height block map - Final, i.e., a stable chain
	Chain map[uint64]Block
	// Delivered Blocks
	DeliveredBlocks map[crypto.Hash]Block

	// Chain head
	Head Block

	// A lock that we use to safely update the chain
	sync.RWMutex
}

// Block is an Ethereum Block
type Block interface {
	Author() uint64
	// Returns the number of transactions in the block
	GetSize() uint64
	// GetBlockHash returns the hash of the block (i.e header)
	GetBlockHash() crypto.Hash
	// IsValid checks if the hash provided in the block and the hash
	IsValid(pk crypto.PubKey) bool
	// GetParentHash returns the hash of the parent block of this block
	GetParentHash() crypto.Hash
	// GetHeight returns the height of this block
	GetHeight() uint64
	// GetExtradata returns extra data from the block
	GetExtradata() []byte
	// ToProto returns the Protocol Buffer variant
	ToProto() *ProtoBlock
	// GetTxs returns the hashes of all the transactions in the block
	GetTxs() [][]byte
}

// Command is a transaction in blockchain terms
type Command []byte
