# AGENTS.md — homelab-dashboard

## What this is

A minimal Go web server that reads a YAML config file and renders an HTML table of service links (LAN/VPN). Designed to run as a container on a homelab, showing quick-access links to internal services.

## Commands

| Command | What it does |
|---------|-------------|
| `task run` | `go run ./cmd/` (defined in Taskfile.yml) |
| `go build ./cmd/` | Compile the binary |
| `DASHBOARD_CONFIG=./myconfig.yml task run` | Run with a custom config path |

No test or lint commands exist. No Makefile.

## Architecture

Single-file app in `cmd/main.go` (~115 lines). Three components:

1. **Config loading** (`getConfig()`) — reads `config.yml` (or `DASHBOARD_CONFIG` env var), unmarshals YAML into `Config{Entries []Entry}` where each `Entry` has `Name`, `Lan`, `Vpn` strings. Called once at startup; `log.Fatal` on errors.
2. **HTTP handlers**:
   - `GET /` — renders an HTML table with inline CSS (dark theme, chartreuse-on-black, auto font size via `calc(20px + 1.5vw)`). Each config entry becomes a row with service name and clickable LAN/VPN links.
   - `GET /reload` — writes "restarting" to response, then calls `exitAfterWhile()` (1ms sleep + `os.Exit(0)`). Relies on container orchestration (Docker/K8s restart policy) to bring the process back with the new config.
3. **Server** — `http.Server` on `:9091` with 5s read/write timeouts. Uses `log.Fatal(s.ListenAndServe())`.

## Key gotchas

- **Config reload = crash**: There is no graceful SIGHUP or config watch. The `/reload` endpoint exits the process. The container must have a restart policy. This is intentional for a minimal homelab setup.
- **No tests at all**: Any changes must be manually verified. No linting or CI tests.
- **All HTML/CSS is inline**: The HTML response is built via `fmt.Sprintf` concatenation inside the handler. CSS is in a `<style>` tag embedded in the Go string literal. No templates, no external assets.
- **Single dependency**: `github.com/goccy/go-yaml` — not the standard `gopkg.in/yaml.v3`. This matters if you're adding new config fields or serialization.
- **Config only read at startup**: Changing `config.yml` while the server is running has no effect until the process restarts (via `/reload` or container restart).
- **No TLS**: Plain HTTP only. Suitable for an internal homelab network behind a reverse proxy.

## CI/CD

- **Dockerfile**: Multi-stage build (`golang:1.26.2-alpine3.23` → `scratch`). Binary is built in `/dist/dash`, then copied to scratch. Exposes port 9091.
- **GitHub Actions** (`.github/workflows/docker-publish.yml`): On push to `main`, builds and pushes to `ghcr.io` with cosign signing. Uses Docker Buildx with GitHub Actions cache.

## Config format

```yaml
entries:
  - name: Service Name
    lan: http://192.168.1.100:8080
    vpn: https://vpn.example.com/service
  - name: Another Service
    lan: http://192.168.1.101:3000
    vpn: https://vpn.example.com/another
```

## Code style

- Standard Go layout (`package main` in `cmd/main.go` — no `internal/` or `pkg/` dirs yet)
- No error return from handlers — errors are `log.Println(err)` (best-effort write)
- No exported symbols outside `main` package
- Simple struct types with no JSON/tags beyond what goccy/go-yaml needs