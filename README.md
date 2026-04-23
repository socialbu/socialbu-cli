# SocialBu CLI

A fast Go CLI for SocialBu.

## What it does

The current CLI supports:
- `config`
- `whoami`
- `account list|get`
- `post list|get|create`
- `team list|create|delete`
- `fixtures capture`

It uses the same API key as your SocialBu account and stores config in `~/.socialbu/config.json`.

## Install

Download the right binary from the GitHub Releases page and make it executable on macOS or Linux.

```bash
chmod +x ./socialbu_<VERSION>_linux_amd64
sudo mv ./socialbu_<VERSION>_linux_amd64 /usr/local/bin/socialbu
```

Windows releases ship as `.exe` binaries.

## Quick start

```bash
socialbu config set-key <your-api-key>
socialbu whoami
socialbu account list
socialbu post list --type scheduled
socialbu post create --accounts 123 --content "Hello" --publish-at "2026-04-21 10:00:00"
```

Supported environment variables:

```bash
SOCIALBU_API_KEY
SOCIALBU_BASE_URL
```

## Build from source

```bash
go build ./...
```

## Releases

GitHub Actions + GoReleaser publish standalone binaries for:
- macOS: amd64, arm64
- Linux: amd64, arm64
- Windows: amd64, arm64

Cut a release by pushing a semver tag:

```bash
git tag vX.Y.Z
git push origin vX.Y.Z
```

## Fixture capture workflow

When a valid API key is available, capture real endpoint responses before changing renderer assumptions:

```bash
go run . fixtures capture > /tmp/socialbu-capture.sh
bash /tmp/socialbu-capture.sh
```

The generated script resolves the repo root automatically, writes the current fixture set into `artifacts/samples/`, and reuses `~/.socialbu/config.json` when `SOCIALBU_API_KEY` is not exported.
