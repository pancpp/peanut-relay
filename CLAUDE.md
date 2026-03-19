# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run

```bash
go build -o peanut-relay .
./peanut-relay -c /path/to/relay.yaml
```

Build with version info via ldflags:
```bash
go build -ldflags "-X github.com/pancpp/peanut-relay/conf.gVersion=1.0.0 -X github.com/pancpp/peanut-relay/conf.gBuildTime=$(date -u +%Y%m%d%H%M%S) -X github.com/pancpp/peanut-relay/conf.gGitHash=$(git rev-parse --short HEAD)" -o peanut-relay .
```

Show version: `./peanut-relay -V`

## Architecture

Peanut-relay is a P2P relay node built in Go using libp2p. The startup sequence in `main.go` is: config → logger → app, with graceful shutdown on SIGTERM/SIGINT/SIGHUP.

### Packages

- **conf** — Configuration via Viper with YAML file support (default: `/etc/peanut-relay/relay.yaml`). CLI flags parsed with pflag. Access config values through package-level getters (`conf.GetString`, `conf.GetInt`, etc.). Build metadata (`gVersion`, `gBuildTime`, `gGitHash`, `gBuildNumber`) is injected via ldflags.
- **logger** — Log initialization using stdlib `log` with lumberjack for file rotation. Writes to file by default; set `enable_console_log: true` in config to also write to stderr.
- **app** — Application logic. Init sequence: resolve discovery addrs → load allowlist → create connection gater → create libp2p host.

### app package files

- `app.go` — `Init()` entry point; `getDiscoveryAddrs()` parses `p2p.discovery_multiaddrs` into `peer.AddrInfo` slices.
- `allowlist.go` — `Allowlist` struct: loads allowed peer IDs from a YAML file (`p2p.allowlist_path`). Provides bidirectional peer ID ↔ IP lookup with RW mutex protection.
- `conn_gater.go` — `ConnGater` implements libp2p `ConnectionGater`. Allows all outbound dials; on `InterceptSecured` only permits peers in the combined allowlist (allowlist peers + discovery servers).
- `host.go` — `newHost()` builds the libp2p host: loads private key from file (base64-encoded), QUIC-only transport, circuit relay v2 service (no per-connection limits), optional private network PSK, connection manager.

### Configuration keys

Default config path: `/etc/peanut-relay/relay.yaml`

P2P settings under `p2p.*`:

| Key | Default | Description |
|-----|---------|-------------|
| `private_key_path` | `/etc/peanut-relay/private-key.b64` | Path to base64-encoded private key file |
| `pnet_psk_path` | `""` | Path to private network PSK file (disabled if empty) |
| `listen_multiaddrs` | `[/ip4/0.0.0.0/udp/19881/quic-v1]` | Listen multiaddresses (QUIC transport) |
| `discovery_multiaddrs` | *(see conf.go)* | Discovery server multiaddresses (always allowed through gater) |
| `allowlist_path` | `/etc/peanut-relay/allowlist.yaml` | Path to YAML file with `peer_ids: [...]` list |
| `relay_conn_lo` | `4096` | Connection manager low watermark |
| `relay_conn_hi` | `8192` | Connection manager high watermark |
| `relay_conn_grace` | `60` | Connection manager grace period (seconds) |
| `relay_reservation_ttl` | `60` | Circuit relay reservation TTL (seconds) |

Log settings under `log.*`: `path`, `max_size` (MB), `max_backups`, `local_time`, `compress`, `enable_console_log`.

### Allowlist file format

```yaml
peer_ids:
  - 12D3KooW...
  - 12D3KooW...
```
