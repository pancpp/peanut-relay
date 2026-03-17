package app

import "context"

func Init(ctx context.Context) error {
	// create whitelist ACL
	w, err := newWhitelistACL()
	if err != nil {
		return err
	}

	// init relay
	if err := initRelay(ctx, w); err != nil {
		return err
	}

	return nil
}
