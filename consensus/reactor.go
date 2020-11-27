package consensus

import (
	"bufio"

	msg "github.com/adithyabhatkajake/libapollo/msg"
	"github.com/adithyabhatkajake/libchatter/crypto"
	"github.com/adithyabhatkajake/libchatter/log"
	pb "github.com/golang/protobuf/proto"
)

type internalMsg struct {
	sender *bufio.ReadWriter
	msg    *msg.ApolloMsg
}

func (n *Apollo) react(m []byte, rw *bufio.ReadWriter) {
	log.Trace("Received a message of size", len(m))
	inMessage := &msg.ApolloMsg{}
	err := pb.Unmarshal(m, inMessage)
	if err != nil {
		log.Error("Received an invalid protocol message from the node", err)
		return
	}
	n.msgChannel <- internalMsg{
		sender: rw,
		msg:    inMessage,
	}
}

func (n *Apollo) handleProtocolMsg(intMsg internalMsg) {
	msgIn := intMsg.msg
	// log.Trace("Received msg", msgIn.String())
	switch x := msgIn.Msg.(type) {
	case *msg.ApolloMsg_Prop:
		prop := msgIn.GetProp()
		log.Debug("Received a proposal from ", prop.Author())
		// Send proposal to propose handler
		n.OnRecvPropose(prop)
	case *msg.ApolloMsg_Npblame:
		blMsg := msgIn.GetNpblame()
		log.Debug("Received an NP blame from ", blMsg.Blame.BlOrigin)
		// Process a No-progress Blame
		n.OnNpBlame(blMsg)
	case *msg.ApolloMsg_Eqblame:
		blMsg := msgIn.GetEqblame()
		log.Debug("Received an equivocating blame from ",
			blMsg.Blame.BlOrigin)
		// Handle an equivocating blame message
		n.OnEqBlame(blMsg)
	case *msg.ApolloMsg_ReqBlk:
		// Someone has asked for a block
		req := msgIn.GetReqBlk()
		n.OnReqBlock(crypto.ToHash(req.Hash), intMsg.sender)
	case *msg.ApolloMsg_RespBlk:
		// Someone has responded to our request
		resp := msgIn.GetRespBlk()
		n.OnFetchBlock(resp.Block)
	case nil:
		log.Warn("Unspecified msg type", x)
	default:
		log.Warn("Unknown msg type", x)
	}
}

// Process protocol messages
func (n *Apollo) protocol() {
	myID := n.GetId()
	for {
		select {
		case intMsg, ok := <-n.msgChannel:
			if !ok {
				log.Error("Msg channel error")
				return
			}
			n.handleProtocolMsg(intMsg)
		case <-n.PoolFull:
			// Try to propose
			if n.pmaker.GetProposer() != myID {
				continue
			}
			log.Trace("Pool is full. Trying to propose")
			n.propose()
		case err, ok := <-n.errCh:
			log.Warn("Received an error.", err)
			if !ok {
				return
			}
		}
	}
}
