package app

import (
	"log"
	"os"

	"github.com/libp2p/go-libp2p/core/peer"
	multiaddr "github.com/multiformats/go-multiaddr"
	"github.com/pancpp/peanut-relay/conf"
	"go.yaml.in/yaml/v2"
)

type whiteListACL struct {
	allowed map[peer.ID]struct{}
}

func newWhitelistACL() (*whiteListACL, error) {
	var peerIdList []peer.ID

	// load peer IDs from discovery server
	discMultiAddrs := conf.GetStringSlice("disc.multiaddrs")
	for _, addr := range discMultiAddrs {
		maddr, err := multiaddr.NewMultiaddr(addr)
		if err != nil {
			log.Printf("discovery server multi-addr parsing err: %v, %v", err, addr)
			return nil, err
		}

		info, err := peer.AddrInfoFromP2pAddr(maddr)
		if err != nil {
			return nil, err
		}

		peerIdList = append(peerIdList, info.ID)
	}

	// load peer IDs from whitelist file
	type whitelist struct {
		PeerIDs []string `yaml:"peer_ids"`
	}

	path := conf.GetString("p2p.acl_whitelist_path")
	if path == "" {
		return nil, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("reading whitelist file err: %v, path: %s", err, path)
		return nil, err
	}

	var wl whitelist
	if err := yaml.Unmarshal(data, &wl); err != nil {
		log.Printf("parsing whitelist file err: %v", err)
		return nil, err
	}

	for _, peerID := range wl.PeerIDs {
		id, err := peer.Decode(peerID)
		if err != nil {
			return nil, err
		}
		peerIdList = append(peerIdList, id)
	}

	allowed := make(map[peer.ID]struct{}, len(peerIdList))
	for _, id := range peerIdList {
		allowed[id] = struct{}{}
	}
	return &whiteListACL{allowed: allowed}, nil
}

func (a *whiteListACL) AllowReserve(p peer.ID, addr multiaddr.Multiaddr) bool {
	_, ok := a.allowed[p]
	if !ok {
		log.Printf("relay reserve denied for peer %s from %s", p, addr)
	}
	return ok
}

func (a *whiteListACL) AllowConnect(src peer.ID, srcAddr multiaddr.Multiaddr, dest peer.ID) bool {
	_, ok := a.allowed[src]
	if !ok {
		log.Printf("relay connect denied for peer %s from %s to %s", src, srcAddr, dest)
	}
	return ok
}
