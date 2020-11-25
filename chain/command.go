package chain

import "github.com/adithyabhatkajake/libchatter/crypto"

// Hash returns the hash of the command
func (c Command) Hash() crypto.Hash {
	return crypto.DoHash(c)
}
