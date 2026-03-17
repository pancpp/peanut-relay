package app

import (
	"log"

	"github.com/libp2p/go-libp2p/core/peer"
	multiaddr "github.com/multiformats/go-multiaddr"
)

type peerWhitelistACL struct {
	allowed map[peer.ID]struct{}
}

func newPeerWhitelistACL(peerIDs []peer.ID) (*peerWhitelistACL, error) {
	allowed := make(map[peer.ID]struct{}, len(peerIDs))
	for _, id := range peerIDs {
		allowed[id] = struct{}{}
	}
	return &peerWhitelistACL{allowed: allowed}, nil
}

func (a *peerWhitelistACL) AllowReserve(p peer.ID, addr multiaddr.Multiaddr) bool {
	_, ok := a.allowed[p]
	if !ok {
		log.Printf("relay reserve denied for peer %s from %s", p, addr)
	}
	return ok
}

func (a *peerWhitelistACL) AllowConnect(src peer.ID, srcAddr multiaddr.Multiaddr, dest peer.ID) bool {
	_, ok := a.allowed[src]
	if !ok {
		log.Printf("relay connect denied for peer %s from %s to %s", src, srcAddr, dest)
	}
	return ok
}
