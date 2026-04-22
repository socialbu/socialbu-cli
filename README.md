# socialbu-cli

A Go-based SocialBu CLI aiming for `gh`-style ergonomics.

## Current bootstrap

This worktree now contains the first Go/Cobra rewrite scaffold with:
- Cobra root command and global `--json`
- Config management via `~/.socialbu/config.json`
- Environment variable support: `SOCIALBU_API_KEY`, `SOCIALBU_BASE_URL`
- HTTP client for SocialBu API bearer auth
- Working command groups for:
  - `config`
  - `whoami`
  - `account list|get`
  - `post list|get|create`
  - `team list|create|delete`

## Build

```bash
go build ./...
```

## Quick start

```bash
go run . config set-key <your-api-key>
go run . whoami
go run . account list
go run . post list --type scheduled
go run . post create --accounts 123 --content "Hello" --publish-at "2026-04-21 10:00:00"
```

## Config

Stored in:

```bash
~/.socialbu/config.json
```

Supported env vars:

```bash
SOCIALBU_API_KEY
SOCIALBU_BASE_URL
```

## Release automation

GitHub Actions + GoReleaser are configured for tagged releases across:
- macOS: amd64, arm64
- Linux: amd64, arm64
- Windows: amd64, arm64

Cut a release by pushing a semver tag:

```bash
git tag v0.1.0
git push origin v0.1.0
```

## Next targets

- Capture real `--json` fixtures into `artifacts/samples/` and replace remaining fallback field guesses with exact mappings
- Improve response shaping and table rendering per endpoint
- Add higher-level UX polish around analytics, AI, notifications, curation, and media flows
- Decide whether to remove the legacy Node source tree once the Go CLI fully replaces it

## Fixture capture workflow

When a valid API key is available, capture real endpoint responses before changing renderer assumptions:

```bash
go run . fixtures capture > /tmp/socialbu-capture.sh
bash /tmp/socialbu-capture.sh
```

Run those commands from within the repository checkout. The generated script resolves the repo root automatically, writes the current fixture set into `artifacts/samples/`, and reuses `~/.socialbu/config.json` when `SOCIALBU_API_KEY` is not exported.

Prioritize `whoami`, `account`, `post`, and `team` fixture coverage first, then extend the capture set for any endpoint whose renderer logic still relies on fallback field guesses.
