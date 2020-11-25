package txpool

import (
	"github.com/adithyabhatkajake/libapollo/chain"
	"github.com/adithyabhatkajake/libchatter/crypto"
	"github.com/elliotchance/orderedmap"
)

// NewTxPool returns the txpool object after initializing all the internal trackers
func NewTxPool() *TxPool {
	return &TxPool{
		unproposedCommands: make(map[crypto.Hash]struct{}),
		newCommands:        orderedmap.NewOrderedMap(),
	}
}

// AddCommand adds the command to the internal data structures
func (txp *TxPool) AddCommand(cmd chain.Command) {
	txp.Lock()
	defer txp.Unlock()

	hash := cmd.Hash()
	txp.unproposedCommands[hash] = struct{}{}
	txp.newCommands.Set(hash, struct{}{})
}

// GetBlock returns an array if sufficient blocks are available
func (txp *TxPool) GetBlock(blkSize uint64) []crypto.Hash {
	// Now check if the queue is full
	txp.RLock()
	numCmds := uint64(txp.newCommands.Len())
	txp.RUnlock()

	if numCmds < blkSize {
		// log.Trace("Insufficient commands ", numCmds)
		return nil
	}

	txp.Lock()
	defer txp.Unlock()

	candidateCmds := make([]crypto.Hash, blkSize)
	i := uint64(0)
	for e := txp.newCommands.Front(); e != nil; e = e.Next() {
		if i == blkSize {
			break
		}
		cmd := e.Key.(crypto.Hash)
		candidateCmds[i] = cmd
		txp.newCommands.Delete(cmd)
		i++
	}
	return candidateCmds
}

// Clear removes the commands from the pool
// So that future proposals do not show these commands
func (txp *TxPool) Clear(cmd [][]byte) {
	txp.Lock()
	defer txp.Unlock()

	for _, tx := range cmd {
		hash := crypto.ToHash(tx)
		delete(txp.unproposedCommands, hash)
		txp.newCommands.Delete(hash)
	}
}
