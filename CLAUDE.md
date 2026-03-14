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

Peanut-relay is a P2P relay node built in Go. The startup sequence in `main.go` is: config → logger → app, with graceful shutdown on SIGTERM/SIGINT/SIGHUP.

### Packages

- **conf** — Configuration via Viper with YAML file support (default: `/etc/peanut/relay.yaml`). CLI flags parsed with pflag. Access config values through package-level getters (`conf.GetString`, `conf.GetInt`, etc.). Build metadata (`gVersion`, `gBuildTime`, `gGitHash`, `gBuildNumber`) is injected via ldflags.
- **logger** — Log initialization using stdlib `log` with lumberjack for file rotation. Writes to file by default; set `enable_console_log: true` in config to also write to stderr.
- **app** — Application logic entry point (currently a stub).

### Configuration keys

P2P settings under `p2p.*`: `private_key`, `pnet_psk`, `fqdn`, `rendezvous`, `conn_lo`/`conn_hi`/`conn_grace` (connection manager), `reservation_ttl`, `bootstraps` (string slice), `listen_port` (default 19881).

Log settings under `log.*`: `path`, `max_size` (MB), `max_backups`, `local_time`, `compress`.
