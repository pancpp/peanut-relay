package app

import (
	"context"
	"log"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pancpp/peanut-relay/conf"
)

func Init(ctx context.Context) error {
	// get discovery server address info
	discoveryAddrInfo, err := getDiscoveryAddrs()
	if err != nil {
		return err
	}
	log.Println("app: discovery addr info:", discoveryAddrInfo)

	// init allowlist
	allowlist, err := newAllowlist()
	if err != nil {
		return err
	}

	// init connection gater
	connGater, err := newConnGater(allowlist, discoveryAddrInfo)
	if err != nil {
		return err
	}

	// create p2p host
	p2pHost, err := newHost(connGater, discoveryAddrInfo)
	if err != nil {
		return err
	}

	log.Println("PeerID:", p2pHost.ID())
	log.Println("Listen Addrs:", p2pHost.Addrs())

	return nil
}

func getDiscoveryAddrs() ([]peer.AddrInfo, error) {
	var discoveryAddrInfo []peer.AddrInfo
	for _, addrStr := range conf.GetStringSlice("p2p.discovery_multiaddrs") {
		addrInfo, err := peer.AddrInfoFromString(addrStr)
		if err != nil {
			return nil, err
		}
		discoveryAddrInfo = append(discoveryAddrInfo, *addrInfo)
	}

	return discoveryAddrInfo, nil
}
