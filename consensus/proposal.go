package consensus

import (
	"container/list"

	"github.com/adithyabhatkajake/libapollo/chain"
	msg "github.com/adithyabhatkajake/libapollo/msg"
	"github.com/adithyabhatkajake/libchatter/crypto"
	"github.com/adithyabhatkajake/libchatter/log"
	"github.com/adithyabhatkajake/libchatter/util"
	pb "github.com/golang/protobuf/proto"
)

// Monitor pending commands, if there is any change and the current node is the
// leader, then start proposing blocks
func (n *Apollo) propose() {
	log.Debug("Starting a propose step")

	// Hey, give me a block to propose please
	n.TxExtractCh <- TxCleave{}
	// Thanks for the block
	cand := <-n.ExtractBlk
	if cand == nil || len(cand) == 0 {
		return // we do not have enough blocks
	}
	// Propose this block
	head := n.bc.GetHead()
	headht := head.GetHeight()
	prop := n.NewCandidateProposal(cand, head, headht+1, nil)
	nblk := prop.GetBlock()
	bhash := nblk.GetBlockHash()
	n.bc.Storage[bhash] = nblk
	n.bc.DeliveredBlocks[bhash] = nblk

	// Ship proposal to processing
	relayMsg := &msg.ApolloMsg{
		Msg: &msg.ApolloMsg_Prop{
			Prop: prop.ToProto(),
		},
	}

	log.Debug("Proposing block: ", nblk.GetHeight())
	log.Debug("Hashes: ", util.HashToString(
		nblk.GetBlockHash()))

	// Leader sends new block to all the other nodes
	n.Broadcast(relayMsg)

	// Self Process block
	n.OnRecvPropose(prop)
}

// IsValid checks
// 1. The block has correct hashes, signatures
func (n *Apollo) IsValid(blk chain.Block) bool {
	author := blk.Author()
	if author == n.GetID() {
		return true
	}
	return blk.IsValid(n.GetPubKeyFromID(author))
}

// OnRecvPropose reacts to the proposal
func (n *Apollo) OnRecvPropose(prop msg.Proposal) {
	blk := prop.GetBlock()
	ht := blk.GetHeight()
	bhash := blk.GetBlockHash()

	log.Trace("Handling proposal ", ht)

	if !n.IsValid(blk) {
		log.Warn("Invalid proposal received")
		return
	}

	nblk := n.bc.Chain[ht]
	if nblk != nil && nblk.GetBlockHash() != bhash {
		log.Warn("Equivocation detected.")
		return
	}
	// Add this block to our storage
	n.bc.AddToStorage(blk)

	proposer := prop.Author()
	// Ensure that the block and all its parents are delivered
	n.ensureDelivered(bhash, proposer,
		func(_ chain.Block) {
			n.OnDeliveredProposal(prop)
		})
	// Stop blame timer for author, since we got a valid proposal
	// TODO
	// Start blame timer for nextLeader
	// TODO
}

// OnDeliveredProposal runs after ensuring that all the parents of a block are delivered
func (n *Apollo) OnDeliveredProposal(prop msg.Proposal) {
	blk := prop.GetBlock()
	proposer := prop.Author()
	ht := blk.GetHeight()
	n.bc.UpdateHead(blk)
	unHandledProps := list.New()
	log.Trace("Handling delivered proposal ", ht)

	lastSeen := n.pmaker.lastBlock
	log.Debug("Last Seen Block: ", lastSeen.GetHeight())
	head := blk

	if lastSeen.GetHeight() >= blk.GetHeight() {
		log.Warn("Already processed this block")
		return
	}
	// Go back until the head meets lastSeen
	for head.GetHeight() > lastSeen.GetHeight() {
		unHandledProps.PushFront(head)
		log.Debug("Added ", head.GetHeight(), " block to the list")
		log.Debug("LastSeen ", lastSeen.GetHeight())
		log.Debug("Block ", blk.GetHeight())
		// update head backwards
		headHash := head.GetParentHash()
		head = n.bc.GetFromStorage(headHash)
	}
	log.Trace("Finished adding blocks to the queue")
	for e := unHandledProps.Front(); e != nil; {
		tempBlk := e.Value.(chain.Block)
		log.Debug("For proposal ", ht, " Processing block ", tempBlk.GetHeight())
		if tempBlk.Author() != n.pmaker.GetProposer() {
			log.Info("Proposal received from invalid node ", proposer,
				", not leader ", n.pmaker.GetProposer())
			return
		}
		if tempBlk.GetHeight() != lastSeen.GetHeight()+1 {
			log.Warn("Invalid height parent")
			log.Warn("Last seen height: ", lastSeen.GetHeight())
			log.Warn("Have height: ", tempBlk.GetHeight())
			return
		}
		// Add this block to our chain
		n.bc.AddBlock(tempBlk)
		n.NewBlockNotify <- tempBlk // Remove the proposals from this block

		log.Trace("Updating leaders")
		n.pmaker.OnFinishPropose(tempBlk) // Update leaders

		unHandledProps.Remove(e)
		e = unHandledProps.Front()
	}
	if ht > n.GetNumberOfFaultyNodes() {
		commitHeight := ht - n.GetNumberOfFaultyNodes()
		cblk, _ := n.bc.CheckExists(commitHeight)
		n.CommitNotifyCh <- cblk
		// n.OnCommit(cblk)
	}
	nextLeader := n.pmaker.GetProposer()
	log.Trace("New Leader", nextLeader)
	if n.GetID() == nextLeader {
		n.propose()
	} else {
		msg := &msg.ApolloMsg{
			Msg: &msg.ApolloMsg_Prop{
				Prop: prop.ToProto(),
			},
		}
		n.SendTo(nextLeader, msg)
	}
}

// This function fetches the block with hash bHash
// Ensures the new block's parents are delivered
// And finally runs the callback function with the block requested
func (n *Apollo) ensureDelivered(bHash crypto.Hash, peer uint64,
	cb func(chain.Block)) {
	log.Trace("Ensuring the existence of parent ", util.HashToString(bHash))
	if blk := n.bc.GetDeliveredBlock(bHash); blk != nil {
		cb(blk)
		return
	}
	blk := n.bc.GetFromStorage(bHash)
	// No such block is known
	// Get the block with bHash
	// Recursively ensure that its parents are also delivered
	n.GetBlock(blk.GetParentHash(), peer,
		func(fetchedBlk chain.Block) {
			log.Info("Fetched the parent block from peers")
			// We got some block with hash bHash
			author := fetchedBlk.Author()
			if !fetchedBlk.IsValid(n.GetPubKeyFromID(author)) {
				return
			}
			// The fetched block is semantically correct
			n.ensureDelivered(fetchedBlk.GetParentHash(), peer,
				func(_ chain.Block) {
					n.bc.AddDeliveredBlock(fetchedBlk)
					cb(blk)
				})
		})

	log.Trace("All parents are delivered after fetching")
}

// NewCandidateProposal returns a proposal message built using commands
func (n *Apollo) NewCandidateProposal(cand []crypto.Hash, parent chain.Block,
	newHeight uint64, extra []byte) msg.Proposal {
	// Start setting block fields
	pbody := &chain.ProtoBody{
		TxHashes:  make([][]byte, len(cand)),
		Responses: nil,
	}
	for i := 0; i < len(cand); i++ {
		pbody.TxHashes[i] = cand[i].GetBytes()
	}
	data, _ := pb.Marshal(pbody)
	pheader := &chain.ProtoHeader{
		Author:     n.GetID(),
		BodyHash:   crypto.DoHash(data).GetBytes(),
		Extra:      extra,
		Height:     newHeight,
		ParentHash: parent.GetBlockHash().GetBytes(),
		View:       n.view,
	}

	// Sign
	data, _ = pb.Marshal(pheader)
	newBlockHash := crypto.DoHash(data)
	sig, err := n.GetMyKey().Sign(data)
	if err != nil {
		log.Error("Error in signing a block during proposal")
		panic(err)
	}
	prop := &msg.ProtoProp{
		Blk: &chain.ProtoBlock{
			Header: pheader,
			Body:   pbody,
			Proof:  sig,
			Hash:   newBlockHash.GetBytes(),
		},
	}
	return prop
}
