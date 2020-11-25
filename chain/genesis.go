package chain

import "github.com/adithyabhatkajake/libchatter/crypto"

var (
	// EmptyHash is the hash used to fill empty values
	EmptyHash = crypto.Hash{0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0}
	// GenesisBlock is the root/genesis block
	GenesisBlock = &ProtoBlock{
		Header: &ProtoHeader{
			Author:     uint64(0),
			BodyHash:   EmptyHash.GetBytes(),
			Extra:      nil,
			Height:     uint64(0),
			ParentHash: EmptyHash.GetBytes(),
		},
		Body: &ProtoBody{
			Responses: nil,
			TxHashes:  nil,
		},
		Hash: EmptyHash.GetBytes(),
	}
)
