package msg

import (
	chain "github.com/adithyabhatkajake/libapollo/chain"
	crypto "github.com/adithyabhatkajake/libchatter/crypto"
)

const (
	// MaxMsgSize defines the biggest message to be ever recived in the system
	MaxMsgSize = 1024 * 1024 * 1024 // 500 kB
)

// Proposal is an interface for new blocks proposed
type Proposal interface {
	Author() uint64
	GetBlock() chain.Block
	IsValid(crypto.PubKey) bool
	GetProof() []byte
	ToProto() *ProtoProp
}

// PartCert is a partial certificate consisting of multiple signatures on a fixed data
type PartCert interface {
	GetNumSigners() uint64            // Return the number of signatures contained in the certificate
	AddSignature(uint64, []byte)      // Add signer's signature to the certificate
	SetData([]byte)                   // Set this as the data for signing
	GetData() []byte                  // Returns the Data on which we have signatures
	GetSignatureFromID(uint64) []byte // Return the signature for signer
	GetSigners() []uint64
}

// NPBlame is a no-progress blame
type NPBlame interface {
	ToProto() *NoProgressBlame
}

// EqBlame is an equivocating blame
type EqBlame interface {
	ToProto() *EquivocationBlame
}
