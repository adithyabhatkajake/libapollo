package msg

import (
	chain "github.com/adithyabhatkajake/libapollo/chain"
	crypto "github.com/adithyabhatkajake/libchatter/crypto"
	"google.golang.org/protobuf/proto"
)

// Implement Proposal interface

// Author returns the node that created the block
func (p *ProtoProp) Author() uint64 {
	return p.Blk.Author()
}

// GetBlock returns the block interface
func (p *ProtoProp) GetBlock() chain.Block {
	return p.Blk
}

// IsValid checks if the proposal is correctly formed, i.e.,
// correctly signed
func (p *ProtoProp) IsValid(pk crypto.PubKey) bool {
	internalBlk := p.Blk
	data, err := proto.Marshal(internalBlk.Header)
	if err != nil {
		return false
	}

	isCorrect, err := pk.Verify(data, internalBlk.Proof)
	if err != nil {
		return false
	}
	return isCorrect
}

// ToProto gives the underlying implementation of the proposal
func (p *ProtoProp) ToProto() *ProtoProp {
	return p
}

// GetProof returns the valid signature on the block
func (p *ProtoProp) GetProof() []byte {
	return p.Blk.Proof
}
