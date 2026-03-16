package app

import "context"

func Init(ctx context.Context) error {
	if err := initRelay(ctx); err != nil {
		return err
	}

	return nil
}
