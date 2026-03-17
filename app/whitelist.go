package app

import (
	"fmt"
	"os"

	"github.com/pancpp/peanut-relay/conf"
	"go.yaml.in/yaml/v3"
)

func loadWhitelist() ([]string, error) {
	type whitelist struct {
		PeerIDs []string `yaml:"peer_ids"`
	}

	path := conf.GetString("p2p.acl_whitelist_path")
	if path == "" {
		return nil, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading whitelist file: %w", err)
	}

	var wl whitelist
	if err := yaml.Unmarshal(data, &wl); err != nil {
		return nil, fmt.Errorf("parsing whitelist file: %w", err)
	}

	return wl.PeerIDs, nil
}
