package consensus

import (
	"bufio"
	"container/list"

	"github.com/adithyabhatkajake/libapollo/msg"
	"github.com/adithyabhatkajake/libchatter/crypto"
	pb "github.com/golang/protobuf/proto"

	"github.com/adithyabhatkajake/libapollo/chain"
)

var (
	pendingFetches = list.New()
)

type pair struct {
	first  interface{}
	second interface{}
}

// GetBlock fetches the block from the hash from the neighbors
func (n *Apollo) GetBlock(h crypto.Hash, peer uint64, callback func(chain.Block)) {
	req := &msg.ApolloMsg{
		Msg: &msg.ApolloMsg_ReqBlk{
			ReqBlk: &msg.RequestBlock{
				Hash: h.GetBytes(),
			},
		},
	}
	pendingFetches.PushBack(pair{
		first:  h,
		second: callback,
	})
	// Use sync stream to request a block from the peer who asked for it.
	n.SendTo(peer, req)
}

// OnFetchBlock is invoked when a block is fetched via synchronization
// It adds the block to storage
func (n *Apollo) OnFetchBlock(blk chain.Block) {
	if blk == nil {
		return
	}
	var cb func(chain.Block)
	for e := pendingFetches.Front(); e != pendingFetches.Back(); e = e.Next() {
		p := e.Value.(pair)
		if p.first.(crypto.Hash) == blk.GetBlockHash() {
			pendingFetches.Remove(e)
			cb = p.second.(func(chain.Block))
			break
		}
	}
	n.bc.AddToStorage(blk)
	cb(blk)
}

// OnReqBlock is invoked when some node asks for a block
func (n *Apollo) OnReqBlock(bhash crypto.Hash, by *bufio.ReadWriter) {
	blk := n.bc.GetFromStorage(bhash)
	resp := &msg.ApolloMsg{
		Msg: &msg.ApolloMsg_RespBlk{
			RespBlk: &msg.ResponseBlock{
				Block: blk.ToProto(),
			},
		},
	}

	data, err := pb.Marshal(resp)
	if err != nil {
		return
	}
	by.Write(data)
	by.Flush()
}
