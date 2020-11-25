package txpool

import (
	"sync"

	"github.com/adithyabhatkajake/libchatter/crypto"
	"github.com/elliotchance/orderedmap"
)

// TxPool is a map of hash to pending commands
type TxPool struct {
	unproposedCommands map[crypto.Hash]struct{}
	newCommands        *orderedmap.OrderedMap
	sync.RWMutex
}

// TxManager is an unsafe transaction pool manager
type TxManager struct {
	// These are commands that are not yet committed
	unCommitted map[crypto.Hash]struct{}
	// These are commands that are not yet proposed
	unproposed *orderedmap.OrderedMap
}
