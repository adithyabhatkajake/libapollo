package consensus

import (
	"github.com/adithyabhatkajake/libapollo/chain"
)

const (
	// DefaultLeaderID is the ID of the Replica that the protocol starts with
	DefaultLeaderID uint64 = 1
)

// RRPaceMaker is a Round-Robin Pace Maker
type RRPaceMaker struct {
	currentLeader uint64
	numNodes      uint64
	lastBlock     chain.Block
}

// GetProposer returns the proposer for the next block
func (p *RRPaceMaker) GetProposer() uint64 {
	return p.currentLeader
}

// OnFinishPropose cleans up after a propose
// It updates the leader and the last seen block
func (p *RRPaceMaker) OnFinishPropose(blk chain.Block) {
	p.currentLeader = (p.currentLeader + 1) % p.numNodes
	if p.lastBlock.GetHeight() < blk.GetHeight() {
		p.lastBlock = blk
	}
}
