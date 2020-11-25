package txpool

import (
	"github.com/adithyabhatkajake/libapollo/chain"
	"github.com/adithyabhatkajake/libchatter/crypto"
	"github.com/elliotchance/orderedmap"
)

// NewTxManager returns the txpool object after initializing all the internal
// trackers
func NewTxManager() *TxManager {
	return &TxManager{
		unCommitted: make(map[crypto.Hash]struct{}),
		unproposed:  orderedmap.NewOrderedMap(),
	}
}

// AddCommand adds the command to the internal data structures
func (txm TxManager) AddCommand(cmd chain.Command) {
	hash := cmd.Hash()
	txm.unCommitted[hash] = struct{}{}
	txm.unproposed.Set(hash, struct{}{})
}

// GetBlock returns an array if sufficient blocks are available
func (txm TxManager) GetBlock(blkSize uint64) []crypto.Hash {
	// Now check if the queue is full
	numCmds := uint64(txm.unproposed.Len())
	if numCmds < blkSize {
		// log.Trace("Insufficient commands ", numCmds)
		return nil
	}

	candidateCmds := make([]crypto.Hash, blkSize)
	i := uint64(0)
	for e := txm.unproposed.Front(); e != nil; {
		if i == blkSize {
			break
		}
		cmd := e.Key.(crypto.Hash)
		candidateCmds[i] = cmd
		txm.unproposed.Delete(cmd)
		e = txm.unproposed.Front()
		i++
	}
	return candidateCmds
}

// Size returns the size of the transaction pool
func (txm TxManager) Size() uint64 {
	return uint64(txm.unproposed.Len())
}

// Clear removes the commands from the pool
// So that future proposals do not show these commands
func (txm TxManager) Clear(cmd [][]byte) {
	for _, tx := range cmd {
		hash := crypto.ToHash(tx)
		delete(txm.unCommitted, hash)
		txm.unproposed.Delete(hash)
	}
}
