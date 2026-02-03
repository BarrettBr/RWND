# Usage

This doc covers basic usage for recording and replaying traffic.

## Install

```bash
go install github.com/BarrettBr/RWND/cmd/rwnd@latest
```

## Record Traffic

Start the proxy and point it at your backend:

```bash
rwnd proxy --listen :8080 --target http://localhost:3000
```

Logs are written to `.rwnd/logs/` by default, one file per proxy run.

## Replay Traffic

Replay is interactive by default and uses the latest log file:

```bash
rwnd replay
```

Controls:

- `Enter`: Step to the next request
- `r`: Replay the current request and show old/new responses
- `q`: Quit

To replay a specific log file:

```bash
rwnd replay --log path/to/file.jsonl
```

## Flags

Proxy:

- `--listen`: Address to listen on (default `:8080`)
- `--target`: Upstream service to forward traffic to (required)
- `--log`: Path to write recorded traffic (default `.rwnd/logs/`)

Replay:

- `--log`: Path to a recorded traffic log or log directory (default `.rwnd/logs/`)
