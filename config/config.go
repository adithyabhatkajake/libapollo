package config

import (
	"fmt"
	"time"

	"github.com/adithyabhatkajake/libchatter/crypto"
	peerstore "github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
)

// Implement all the interfaces, i.e.,
// 1. net
// 2. crypto
// 3. config

// GetID returns the Id of this instance
func (apl *NodeConfig) GetID() uint64 {
	return apl.GetId()
}

// GetP2PAddrFromID gets the P2P address of the node rid
func (apl *NodeConfig) GetP2PAddrFromID(rid uint64) string {
	address := apl.GetNodeAddressMap()[rid]
	addr := fmt.Sprintf("/ip4/%s/tcp/%s", address.IP, address.Port)
	return addr
}

// GetMyKey returns the private key of this instance
func (apl *NodeConfig) GetMyKey() crypto.PrivKey {
	return apl.pvtKey
}

// GetPubKeyFromID returns the Public key of node whose ID is nid
func (apl *NodeConfig) GetPubKeyFromID(nid uint64) crypto.PubKey {
	return apl.nodeKeyMap[nid]
}

// GetPeerFromID returns libp2p peerInfo from the config
func (apl *NodeConfig) GetPeerFromID(nid uint64) peerstore.AddrInfo {
	pID, err := peerstore.IDFromPublicKey(apl.GetPubKeyFromID(nid))
	if err != nil {
		panic(err)
	}
	addr, err := ma.NewMultiaddr(apl.GetP2PAddrFromID(nid))
	if err != nil {
		panic(err)
	}
	pInfo := peerstore.AddrInfo{
		ID:    pID,
		Addrs: []ma.Multiaddr{addr},
	}
	return pInfo
}

// GetNumNodes returns the protocol size
func (apl *NodeConfig) GetNumNodes() uint64 {
	return apl.GetInfo().GetNodeSize()
}

// GetClientListenAddr returns the address where to talk to/from clients
func (apl *NodeConfig) GetClientListenAddr() string {
	id := apl.GetID()
	address := apl.GetNodeAddressMap()[id]
	addr := fmt.Sprintf("/ip4/%s/tcp/%s", address.IP, apl.ClientPort)
	return addr
}

// GetBlockSize returns the number of commands that can be inserted in one block
func (apl *NodeConfig) GetBlockSize() uint64 {
	return apl.GetInfo().GetBlockSize()
}

// GetDelta returns the synchronous wait time
func (apl *NodeConfig) GetDelta() time.Duration {
	timeInSeconds := apl.ProtocolConfig.GetDelta()
	return time.Duration(int(timeInSeconds*1000)) * time.Millisecond
}

// GetCommitWaitTime returns how long to wait before committing a block
func (apl *NodeConfig) GetCommitWaitTime() time.Duration {
	return apl.GetDelta() * 2
}

// GetNPBlameWaitTime returns how long to wait before sending the NP Blame
func (apl *NodeConfig) GetNPBlameWaitTime() time.Duration {
	return apl.GetDelta() * 3
}

// GetNumberOfFaultyNodes returns f
// We do this because f can be less than floor(n-1)/2
func (apl *NodeConfig) GetNumberOfFaultyNodes() uint64 {
	return apl.GetInfo().GetFaults()
}
