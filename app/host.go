package app

import (
	"encoding/base64"
	"log"
	"os"
	"strings"
	"time"

	"github.com/libp2p/go-libp2p"
	coreconnmgr "github.com/libp2p/go-libp2p/core/connmgr"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/core/pnet"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"
	quic "github.com/libp2p/go-libp2p/p2p/transport/quic"
	"github.com/pancpp/peanut-relay/conf"
)

func newHost(connGater coreconnmgr.ConnectionGater, discoveryAddrInfo []peer.AddrInfo) (host.Host, error) {
	// p2p opts
	var opts []libp2p.Option

	// private key
	privateKeyPath := conf.GetString("p2p.private_key_path")
	privateKeyB64, err := os.ReadFile(privateKeyPath)
	if err != nil {
		log.Printf("reading private key err: %v, path: %s", err, privateKeyPath)
		return nil, err
	}
	privateKeyBytes, err := base64.StdEncoding.DecodeString(strings.TrimSpace(string(privateKeyB64)))
	if err != nil {
		log.Printf("base64 unmarshal err: %v, string: %s", err, string(privateKeyB64))
		return nil, err
	}
	privateKey, err := crypto.UnmarshalPrivateKey(privateKeyBytes)
	if err != nil {
		log.Printf("invalid private key, err: %v, string: %s", err, string(privateKeyBytes))
		return nil, err
	}
	opts = append(opts, libp2p.Identity(privateKey))

	// listen addresses
	listenAddrs := conf.GetStringSlice("p2p.listen_multiaddrs")
	opts = append(opts,
		libp2p.Transport((quic.NewTransport)),
		libp2p.ListenAddrStrings(listenAddrs...))

	// relay
	relayRes := relay.DefaultResources()
	relayRes.Limit = nil
	relayRes.ReservationTTL = time.Duration(conf.GetInt("p2p.relay_reservation_ttl")) * time.Second
	relayOpts := []relay.Option{relay.WithResources(relayRes)}

	// add connection gater
	if connGater != nil {
		opts = append(opts, libp2p.ConnectionGater(connGater))
	}

	opts = append(opts,
		libp2p.ForceReachabilityPublic(),
		libp2p.EnableRelayService(relayOpts...),
		libp2p.DisableRelay(),
	)

	// NAT service
	opts = append(opts, libp2p.EnableNATService())

	// pnet psk
	pskPath := conf.GetString("p2p.pnet_psk_path")
	if pskPath != "" {
		pskFile, err := os.Open(pskPath)
		if err != nil {
			return nil, err
		}
		defer pskFile.Close()

		psk, err := pnet.DecodeV1PSK(pskFile)
		if err != nil {
			return nil, err
		}

		opts = append(opts, libp2p.PrivateNetwork(psk))
		log.Println("private network is enabled")
	}

	// connection manager
	connMgr, err := connmgr.NewConnManager(
		conf.GetInt("p2p.relay_conn_lo"),
		conf.GetInt("p2p.relay_conn_hi"),
		connmgr.WithGracePeriod(time.Duration(conf.GetInt("p2p.relay_conn_grace"))*time.Second),
	)
	if err != nil {
		return nil, err
	}
	opts = append(opts, libp2p.ConnectionManager(connMgr))

	// create libp2p host
	h, err := libp2p.New(opts...)
	if err != nil {
		return nil, err
	}

	// add discovery server address info
	for _, addrInfo := range discoveryAddrInfo {
		h.Peerstore().AddAddrs(addrInfo.ID, addrInfo.Addrs, peerstore.PermanentAddrTTL)
	}

	return h, nil
}
