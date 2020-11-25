package txpool_test

import (
	"testing"

	"github.com/adithyabhatkajake/libapollo/txpool"
	"github.com/adithyabhatkajake/libchatter/crypto"
)

func TestManager(t *testing.T) {
	txman := txpool.NewTxManager()
	if txman == nil {
		t.Error("unable to create the transaction pool")
	}
	cmds := make([]crypto.Hash, 100)
	for i := 0; i < 100; i++ {
		cmd := []byte{uint8(i)}
		cmds[i] = crypto.DoHash(cmd)
		txman.AddCommand(cmd)
	}
	if txman.Size() != 100 {
		t.Error("Ordered map not added enough elements")
	}
	blk := txman.GetBlock(100)
	for i := 0; i < 100; i++ {
		gotHash := blk[i]
		expectedHash := crypto.DoHash([]byte{uint8(i)})
		if gotHash != expectedHash {
			t.Error("Order not preserved for ", i)
			t.Log(gotHash)
			t.Log(expectedHash)
		}
	}
}
