package consensus

import (
	"github.com/adithyabhatkajake/libapollo/msg"
	"github.com/adithyabhatkajake/libchatter/log"
)

// OnNpBlame reacts to a blame message
func (n *Apollo) OnNpBlame(bl msg.NPBlame) {
	// TODO
	log.Warn("Np blame: UNIMPLEMENTED")
}

// OnEqBlame reacts to Equivocation blames
func (n *Apollo) OnEqBlame(bl msg.EqBlame) {
	// TODO
	log.Warn("Eq Blame: UNIMPLEMENTED")
}
