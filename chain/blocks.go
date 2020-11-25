package chain

import (
	"github.com/adithyabhatkajake/libchatter/crypto"
	"github.com/adithyabhatkajake/libchatter/util"
	pb "github.com/golang/protobuf/proto"
	"github.com/prometheus/common/log"
)

// ComputeHash computes the hash for the block (i.e. the header)
func (b *ProtoBlock) ComputeHash() crypto.Hash {
	data, _ := pb.Marshal(b.Header)
	return crypto.DoHash(data)
}

// Author returns the proposer for the block
func (b *ProtoBlock) Author() uint64 {
	return b.Header.Author
}

// ======================================
// Implementing ProtoBlock as chain.Block
// ======================================

// GetSize returns the number of transactions present in the block
func (b *ProtoBlock) GetSize() uint64 {
	return uint64(len(b.Body.GetTxHashes()))
}

// GetBlockHash returns the hash of the block this header is referring to
func (b *ProtoBlock) GetBlockHash() crypto.Hash {
	return crypto.ToHash(b.Hash)
}

// IsValid returns whether the block has a valid hash and signatures
func (b *ProtoBlock) IsValid(pk crypto.PubKey) bool {
	data, err := pb.Marshal(b.Header)
	if err != nil {
		return false
	}
	h1 := crypto.DoHash(data)
	h2 := b.GetBlockHash()

	if h1 != h2 {
		log.Warn("Invalid block. Computed Hash and the obtained hash does not match")
		log.Warn("Block hash: ", util.HashToString(h1))
		log.Warn("Computed hash: ", util.HashToString(h2))
		return false
	}

	isCorrect, err := pk.Verify(data, b.Proof)
	if err != nil {
		return false
	}
	return isCorrect
}

// GetParentHash returns the hash of the parent block for this block
func (b *ProtoBlock) GetParentHash() crypto.Hash {
	return crypto.ToHash(b.Header.ParentHash)
}

// GetHeight returns the height of this block
func (b *ProtoBlock) GetHeight() uint64 {
	return b.Header.Height
}

// GetExtradata returns extra data from the block
func (b *ProtoBlock) GetExtradata() []byte {
	return b.Header.Extra
}

// ToProto returns the protocol buffer variant of the file
func (b *ProtoBlock) ToProto() *ProtoBlock {
	return b
}

// GetTxs returns the hashes of all the transactions in the block
func (b *ProtoBlock) GetTxs() [][]byte {
	return b.Body.TxHashes
}
