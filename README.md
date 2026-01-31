# RWND

A local HTTP reverse proxy that logs request/response pairs and supports deterministic replay for debugging.

## Description

`RWND` is a lightweight CLI tool that sits in front of an HTTP service, records & logs incoming requests and outgoing responses, and let's you replay them later in a determinstic, step-by-step way so as to allow for easy debugging of your services.

It acts as a reverse proxy during recording, writing traffic to a structured log file, then replays that traffic back into a target service to help debug behavior, test changes, or reproduce tricky issues.

## Motivation

I often have run into bugs that showed up only specific request sequences or timing conditions. Logging helped but they didn't let you re-run the problem. Because of this manual reproduction was slow, unreliable and sometimes impossible.

Relying on normal test scripts and curl commands usually would fail to capture the shape of the traffic resulting in hard to reproduce bugs

So using `RWND` I can now:

- Capture real HTTP traffic with zero app changes
- Rewind and replay requests deterministically
- Step through interactions and observe behavior under identical conditions

By doing this I was able to turn this from guesswork or just harder generalized scripting to a cleaner quicker repeatable debugging process

## Quick Start

1. Install `rwnd`

```bash
go install github.com/yourusername/rwnd@latest
```

2. Record traffic

```bash
rwnd proxy --listen :8080 --target http://localhost:3000
```

3. Replay traffic

```bash
rwnd replay
```

## Usage

### Proxy Mode

Proxy mode is used for setting up the reverse proxy and recording the data

```bash
rwnd proxy [options]
```

Available Flags:

- `--listen`: Address to listen on (Default `:8080`)
- `--target`: Upstream service to forward traffic to
- `--log`: Path to write recorded traffic (Defaults to `.rwnd/logs/latest.jsonl`)
- `--verbose`: Enables verbose logging
- `--help / -h`: Shows help

### Replay Mode

Replay mode is used for replaying the recorded data and being able to step through it

```bash
rwnd replay [options]
```

Available Flags:

- `--log`: Path to a recorded traffic log
- `--step`: Step through each request interactively
- `--help / -h`: Shows help

### Docker

RWND can be built and run as a minimal container image using the provided Dockerfile.

Build the image:
docker build -t rwnd .

Run RWND Example:

```dockerfile
docker run --rm -it \
 -p 8080:8080 \
 -v $(pwd)/.rwnd:/app/.rwnd \
 rwnd --listen :8080 --target http://example.com
```

RWND runs as a non-root user inside the container and writes logs to
/app/.rwnd/logs so the command above would mount that to your working directories `./.rwnd/` folder.

## Contributing

Contributions are welcome.

Feel free to open an issue or submit a pull request if interested in:

- Improving replay controls
- Supporting additional protocols or formats
- Adding filtering or request selection
- Tightening performance or determinism
