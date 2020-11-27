package consensus

import (
	msg "github.com/adithyabhatkajake/libapollo/msg"
	"github.com/adithyabhatkajake/libchatter/log"
	pb "github.com/golang/protobuf/proto"
)

// Broadcast broadcasts a protocol message to all the nodes
func (n *Apollo) Broadcast(m *msg.ApolloMsg) error {
	n.netMutex.Lock()
	defer n.netMutex.Unlock()

	data, err := pb.Marshal(m)
	if err != nil {
		return err
	}
	log.Trace("Broadcasting a message of size ", len(data))

	// If we fail to send a message to someone, continue
	for idx, s := range n.streamMap {
		writeLen, err := s.Writer.Write(data)
		log.Trace("Wrote ", writeLen, " bytes into the stream")
		if err != nil {
			log.Error("Error while sending to node", idx)
			log.Error("Error:", err)
			continue
		}
		err = s.Writer.Flush()
		if err != nil {
			log.Error("Error while sending to node", idx)
			log.Error("Error:", err)
		}
	}
	return nil
}

// SendTo sends a message to a particular node
func (n *Apollo) SendTo(peer uint64, m *msg.ApolloMsg) {
	n.netMutex.Lock()
	defer n.netMutex.Unlock()

	data, err := pb.Marshal(m)
	log.Trace("Sending to ", peer, " a message of size ", len(data))
	if err != nil {
		return
	}

	n.streamMap[peer].Write(data)
	n.streamMap[peer].Flush()
}
