# NetStatus-API-Go

A small TCP ping API that tells you whether `ip:port` is reachable from the server it runs on. Run it
on a server inside mainland China and point NeXT-Panel's GFW detection at it.

## Usage

```text
netstatus-api-go [--config <file or directory>]
```

Without `--config`, `config.json` is searched in `/etc/netstatus-api-go/` and the working directory.
If neither has one, the defaults below are used.

### Docker

Images for `linux/amd64`, `linux/arm64` and `linux/riscv64` are published to GHCR. `latest` and `X.Y.Z`
follow releases, `edge` follows `main`.

```sh
docker run -d --name netstatus-api-go --restart unless-stopped -p 8080:8080 \
  -v /etc/netstatus-api-go:/etc/netstatus-api-go:ro \
  ghcr.io/the-next-project/netstatus-api-go:latest
```

To also push to Docker Hub, set the `DOCKERHUB_USERNAME` repository variable and the
`DOCKERHUB_TOKEN` secret. `DOCKERHUB_IMAGE` overrides the image name, which defaults to
`<DOCKERHUB_USERNAME>/netstatus-api-go`.

## API

`GET /v1/tcping?ip=<ip>&port=<port>[&token=<token>]`

| Response | Meaning |
| --- | --- |
| `200 {"status":"true","time":42,"message":"TCP connection successful"}` | Reachable, `time` is the connect latency in ms |
| `200 {"status":"false","message":"TCP connection failed"}` | Unreachable after every attempt |
| `4xx`/`5xx {"status":"error","message":"..."}` | No verdict: bad input (400), bad token (401), non-public target (403), timed out (504) |
| `429` | Rate limited |

`ip` must be a literal IPv4 or IPv6 address. Hostnames are not resolved.

## Configuration

| Key | Default | Description |
| --- | --- | --- |
| `port` | `8080` | Listen port |
| `api_timeout` | `3000` | Time limit for a whole request, in ms |
| `tcping_timeout` | `1000` | Time limit for one connection attempt, in ms |
| `tcping_attempts` | `1` | Attempts before reporting unreachable. `tcping_attempts × tcping_timeout` must not exceed `api_timeout` |
| `rate_limit` | `60` | Requests per minute per client IP, `0` disables it |
| `trust_proxy` | `false` | Take the client IP from `X-Forwarded-For` / `X-Real-IP`. Only enable this behind a reverse proxy |
| `allow_private` | `false` | Allow probing loopback, private, link-local and CGNAT addresses |
| `token` | `""` | If set, requests must pass it as the `token` query parameter |

## NeXT-Panel

Set the detection API URL in the panel's cron settings, adding the token if you configured one:

```text
https://probe.example.com/v1/tcping?ip={ip}&port={port}&token=<token>
```

- The panel probes nodes one after another from a single IP. If you have more nodes than `rate_limit`
  allows per minute, the rest get `429` and are skipped. Raise `rate_limit`, or set a `token` and use `0`.
- The panel only marks a node as blocked on a `200` with `"status":"false"`. Every other response,
  including timeouts, leaves the node's state unchanged.
- One lost SYN is enough to mark a node as blocked and notify users. Use `tcping_attempts` of `2` or `3`
  to reduce false alarms, and raise `api_timeout` to match. The panel waits up to 10 s per probe.
