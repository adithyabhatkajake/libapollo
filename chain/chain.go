package chain

import (
	"github.com/adithyabhatkajake/libchatter/crypto"
)

// NewChain returns an empty chain
func NewChain() *BlockChain {
	c := &BlockChain{}
	// genesis := GetGenesis()

	c.Storage = make(map[crypto.Hash]Block)
	c.DeliveredBlocks = make(map[crypto.Hash]Block)
	c.Chain = make(map[uint64]Block)

	// Set genesis block as the first block
	c.Chain[0] = GenesisBlock
	c.Storage[GenesisBlock.GetBlockHash()] = GenesisBlock
	c.DeliveredBlocks[GenesisBlock.GetBlockHash()] = GenesisBlock

	// Set the head
	c.Head = GenesisBlock

	return c
}

// CheckExists returns a block if it exists, otherwise returns false
func (bc *BlockChain) CheckExists(ht uint64) (Block, bool) {
	blk, exists := bc.Chain[ht]
	return blk, exists
}

// AddBlock adds a block to the blockchain
// and updates the head if this is the highest block
func (bc *BlockChain) AddBlock(b Block) {
	ht := b.GetHeight()
	bc.Chain[ht] = b
	if ht > bc.Head.GetHeight() {
		bc.Head = b
	}
	if bc.DeliveredBlocks[b.GetParentHash()] != nil {
		bc.DeliveredBlocks[b.GetBlockHash()] = b
	}
}

// AddToStorage adds a block to the storage but not the blockchain
func (bc *BlockChain) AddToStorage(b Block) {
	// Set hash and height accessor
	bc.Storage[b.GetBlockHash()] = b
	if bc.DeliveredBlocks[b.GetParentHash()] != nil {
		bc.DeliveredBlocks[b.GetBlockHash()] = b
	}
}

// GetFromStorage returns a block from the storage
func (bc *BlockChain) GetFromStorage(bhash crypto.Hash) Block {
	return bc.Storage[bhash]
}

// GetHead returns the head block
func (bc *BlockChain) GetHead() Block {
	return bc.Head
}

// UpdateHead updates the head of the chain
func (bc *BlockChain) UpdateHead(b Block) {
	ht := b.GetHeight()
	if ht > bc.Head.GetHeight() {
		bc.Head = b
	}
}

// GetDeliveredBlock returns the delivered block with hash h
func (bc *BlockChain) GetDeliveredBlock(h crypto.Hash) Block {
	return bc.DeliveredBlocks[h]
}

// AddDeliveredBlock adds the block to the delivered block cache
// to prevent future delivery requests
func (bc *BlockChain) AddDeliveredBlock(b Block) {
	bc.DeliveredBlocks[b.GetBlockHash()] = b
}
