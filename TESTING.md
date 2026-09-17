# Testing

## Prerequisites

- Go 1.22+ (see `go.mod` for the exact version)
- [golangci-lint](https://golangci-lint.run/welcome/install/) on `PATH`
- Optional: [pre-commit](https://pre-commit.com/#install) for local git hooks

## Quick checks

```bash
# Format
gofmt -w -s .

# Lint
golangci-lint run

# Unit tests (internal packages)
go test ./internal/...

# All packages
go test ./...
```

### Race detector and short mode

CI runs tests with the race detector and `-short`. Locally:

```bash
go test -race -short -count=1 ./...
```

Coverage (atomic mode, matching CI):

```bash
go test -race -short -count=1 \
  -coverprofile=coverage.out -covermode=atomic ./...
go tool cover -func=coverage.out
```

## Package focus

| Path | What it covers |
|------|----------------|
| `./internal/proxy/...` | CONNECT tunnel, HTTPS/// TLS rewrite, cache rules, host-prefix routing |
| `./internal/cli/...` | Daemon lifecycle, Fiber app, shutdown, production proxy routes |
| `./internal/distro/...` | Registry, cache rules, URL patterns |
| `./internal/api/...` | Auth middleware, cache stats API |
| `./internal/mirrors/...` | Mirror selection / benchmarks |

Example:

```bash
go test ./internal/proxy/ -run 'TLS|Connect|ParseTLS' -v
go test ./internal/cli/ -run 'ProductionProxy|Shutdown' -v
```

## Pre-commit hooks

This repo ships a minimal [`.pre-commit-config.yaml`](.pre-commit-config.yaml) that runs `gofmt` and `golangci-lint` before each commit.

### Install

```bash
# Once per machine (pip, brew, or https://pre-commit.com/#install)
pip install pre-commit
# or: brew install pre-commit
```

Ensure `gofmt` and `golangci-lint` are on your `PATH`.

### Enable in this clone

```bash
pre-commit install
```

### Run manually

```bash
# Against staged files (same as the git hook)
pre-commit run

# Against the whole tree
pre-commit run --all-files
```

To skip hooks for a single commit (not recommended for normal work):

```bash
git commit --no-verify
```

## Tips

- Flaky MISS→HIT cache assertions under `-race` on macOS are mitigated with short retries in the relevant daemon tests; prefer `go test -race -short` when reproducing CI.
- Prefer package-scoped runs while iterating (`go test ./internal/proxy/...`) for faster feedback.
- After changing handler or CONNECT code, run both `./internal/proxy/...` and `./internal/cli/...`.
