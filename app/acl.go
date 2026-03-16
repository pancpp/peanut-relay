package app

import (
	"log"

	"github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
)

type peerWhitelistACL struct {
	allowed map[peer.ID]struct{}
}

func newPeerWhitelistACL(peerIDs []string) (*peerWhitelistACL, error) {
	allowed := make(map[peer.ID]struct{}, len(peerIDs))
	for _, s := range peerIDs {
		id, err := peer.Decode(s)
		if err != nil {
			return nil, err
		}
		allowed[id] = struct{}{}
	}
	return &peerWhitelistACL{allowed: allowed}, nil
}

func (a *peerWhitelistACL) AllowReserve(p peer.ID, addr ma.Multiaddr) bool {
	_, ok := a.allowed[p]
	if !ok {
		log.Printf("relay reserve denied for peer %s from %s", p, addr)
	}
	return ok
}

func (a *peerWhitelistACL) AllowConnect(src peer.ID, srcAddr ma.Multiaddr, dest peer.ID) bool {
	_, ok := a.allowed[src]
	if !ok {
		log.Printf("relay connect denied for peer %s from %s to %s", src, srcAddr, dest)
	}
	return ok
}
