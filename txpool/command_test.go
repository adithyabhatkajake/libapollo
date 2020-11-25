package txpool_test

import (
	"testing"

	"github.com/adithyabhatkajake/libapollo/txpool"
	"github.com/adithyabhatkajake/libchatter/crypto"
)

func TestPool(t *testing.T) {
	pool := txpool.NewTxPool()
	if pool == nil {
		t.Error("unable to create the transaction pool")
	}
	cmds := make([]crypto.Hash, 100)
	for i := 0; i < 100; i++ {
		cmd := []byte{uint8(i)}
		cmds[i] = crypto.DoHash(cmd)
		pool.AddCommand(cmd)
	}
	gotHash := pool.GetBlock(1)[0]
	expectedHash := crypto.DoHash([]byte{uint8(0)})
	if gotHash != expectedHash {
		t.Error("Order not preserved")
		t.Log(gotHash)
		t.Log(expectedHash)
	}
}
