# peanut-relay

A libp2p-based circuit relay node for the Peanut P2P network. It accepts relay reservations only from peers listed in an allowlist, using QUIC as the sole transport.

This is a simple example implementation of the relay server for small size private networks.

## Features

- Circuit relay v2 service (no per-connection resource limits)
- Peer allowlist enforced via libp2p `ConnectionGater`
- QUIC-only transport
- Optional private network (PSK)
- NAT service
- Configurable connection manager

## Build

```bash
go build -o peanut-relay .
```

With version info:
```bash
go build -ldflags "-X github.com/pancpp/peanut-relay/conf.gVersion=1.0.0 \
  -X github.com/pancpp/peanut-relay/conf.gBuildTime=$(date -u +%Y%m%d%H%M%S) \
  -X github.com/pancpp/peanut-relay/conf.gGitHash=$(git rev-parse --short HEAD)" \
  -o peanut-relay .
```

## Usage

```bash
./peanut-relay -c /etc/peanut-relay/relay.yaml
./peanut-relay -V   # show version
```

## Configuration

Default config path: `/etc/peanut-relay/relay.yaml`

```yaml
p2p:
  private_key_path: /etc/peanut-relay/private-key.b64   # base64-encoded libp2p private key
  pnet_psk_path: ""                                       # path to PSK file; empty disables private network
  listen_multiaddrs:
    - /ip4/0.0.0.0/udp/19881/quic-v1
  discovery_multiaddrs:
    - /dns4/discovery.example.com/udp/19880/quic-v1/p2p/<peer-id>
  allowlist_path: /etc/peanut-relay/allowlist.yaml
  relay_conn_lo: 4096        # connection manager low watermark
  relay_conn_hi: 8192        # connection manager high watermark
  relay_conn_grace: 60       # grace period in seconds
  relay_reservation_ttl: 60  # relay reservation TTL in seconds

log:
  path: /var/log/peanut/relay.log
  max_size: 500        # MB
  max_backups: 3
  local_time: true
  compress: true
  enable_console_log: false
```

### Private key

Generate a base64-encoded Ed25519 private key and save it to the path specified by `private_key_path`.

### Allowlist

`allowlist_path` points to a YAML file listing the peer IDs permitted to use this relay:

```yaml
peer_ids:
  - 12D3KooW...
  - 12D3KooW...
```

Discovery server peers are always allowed regardless of the allowlist.

## Systemd

Copy the binary and service file, then enable the service:

```bash
cp peanut-relay /srv/peanut-relay/peanut-relay
cp peanut-relay.service /etc/systemd/system/
systemctl daemon-reload
systemctl enable --now peanut-relay
```

The service file expects the binary at `/srv/peanut-relay/peanut-relay` and runs with `WorkingDirectory=/srv/peanut-relay`.
