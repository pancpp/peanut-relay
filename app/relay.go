package app

import (
	"context"
	"encoding/base64"
	"log"
	"os"
	"strings"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/pnet"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"
	quic "github.com/libp2p/go-libp2p/p2p/transport/quic"
	"github.com/pancpp/peanut-relay/conf"
)

var (
	gHost host.Host
)

func initRelay(ctx context.Context) error {
	// p2p opts
	var opts []libp2p.Option

	// private key
	privateKeyPath := conf.GetString("p2p.private_key_path")
	privateKeyB64, err := os.ReadFile(privateKeyPath)
	if err != nil {
		log.Printf("reading private key err: %v, path: %s", err, privateKeyPath)
		return err
	}
	privateKeyBytes, err := base64.StdEncoding.DecodeString(strings.TrimSpace(string(privateKeyB64)))
	if err != nil {
		log.Printf("base64 unmarshal err: %v, string: %s", err, string(privateKeyB64))
		return err
	}
	privateKey, err := crypto.UnmarshalPrivateKey(privateKeyBytes)
	if err != nil {
		log.Printf("invalid private key, err: %v, string: %s", err, string(privateKeyBytes))
		return err
	}
	opts = append(opts, libp2p.Identity(privateKey))

	// listen addresses
	opts = append(opts, libp2p.Transport((quic.NewTransport)))
	listenAddrs := conf.GetStringSlice("p2p.listen_multiaddrs")
	if len(listenAddrs) > 0 {
		opts = append(opts, libp2p.ListenAddrStrings(listenAddrs...))
	}

	// relay
	relayRes := relay.DefaultResources()
	relayRes.Limit = nil
	relayRes.ReservationTTL = time.Duration(conf.GetInt("relay.reservation_ttl")) * time.Second
	relayOpts := []relay.Option{relay.WithResources(relayRes)}

	// peer whitelist ACL
	whitelistPeers, err := loadWhitelist()
	if err != nil {
		return err
	}
	if len(whitelistPeers) > 0 {
		acl, err := newPeerWhitelistACL(whitelistPeers)
		if err != nil {
			return err
		}
		relayOpts = append(relayOpts, relay.WithACL(acl))
		log.Printf("relay peer whitelist enabled: %d peers", len(whitelistPeers))
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
			return err
		}
		defer pskFile.Close()

		psk, err := pnet.DecodeV1PSK(pskFile)
		if err != nil {
			return err
		}

		opts = append(opts, libp2p.PrivateNetwork(psk))
		log.Println("private network is enabled")
	}

	// connection manager
	connMgr, err := connmgr.NewConnManager(
		conf.GetInt("relay.conn_lo"),
		conf.GetInt("relay.conn_hi"),
		connmgr.WithGracePeriod(time.Duration(conf.GetInt("relay.conn_grace"))*time.Second),
	)
	if err != nil {
		return err
	}
	opts = append(opts, libp2p.ConnectionManager(connMgr))

	// create libp2p host
	h, err := libp2p.New(opts...)
	if err != nil {
		return err
	}

	// save variables to global
	gHost = h

	log.Println("PeerID:", h.ID())
	log.Println("Listen Addrs:", h.Addrs())

	return nil
}
