package app

import (
	"context"

	coreconnmgr "github.com/libp2p/go-libp2p/core/connmgr"
	"github.com/pancpp/peanut-relay/conf"
)

func Init(ctx context.Context) error {
	// create connection gater
	var connGater coreconnmgr.ConnectionGater
	if conf.GetBool("p2p.enable_acl") {
		if g, err := newConnGater(); err != nil {
			return err
		} else {
			connGater = g
		}
	}

	// init relay
	if err := initRelay(ctx, connGater); err != nil {
		return err
	}

	return nil
}
