# SocialBu CLI

A fast Go CLI for SocialBu.

## What it does

The current CLI supports:
- `config`
- `whoami`
- `account list|get`
- `post list|get|create`
- `team list|create|delete`
- `analytics posts-count|posts-metrics|top-posts|accounts-metrics|followers|followers-growth|engagement-rate|engagement-trend|inbox-unread-count|automation-logs|team-metrics|team-activity|stats`
- `ai generate|from-post|autocomplete`
- `notifications list|unread|get|mark-read|mark-unread|mark-all-read`
- `curation topics|items|get`
- `media upload|status`
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

If local GitHub auth is unavailable, run the `release` workflow manually and provide:
- `tag`: the semver tag to create, for example `v0.1.2`
- `target`: the commit SHA or branch to tag, for example `94398c6c6080269f7d38104f69c7da1617d4d196`

The manual workflow creates the annotated tag in GitHub Actions, pushes it, runs tests, and publishes the GoReleaser binaries.

## Smoke workflow

`.github/workflows/smoke.yml` runs `go run . whoami` on pull requests and manual dispatch only when the `SOCIALBU_TEST_KEY` GitHub secret is configured. The workflow passes that secret through `SOCIALBU_API_KEY` without printing it.

## Fixture capture workflow

When a valid API key is available, capture real endpoint responses before changing renderer assumptions:

```bash
go run . fixtures capture > /tmp/socialbu-capture.sh
bash /tmp/socialbu-capture.sh
```

The generated script resolves the repo root automatically, writes the current fixture set into `artifacts/samples/`, and reuses `~/.socialbu/config.json` when `SOCIALBU_API_KEY` is not exported.

The smoke workflow runs `./scripts/smoke-readonly.sh` when the `SOCIALBU_TEST_KEY` repository secret is configured. The suite covers non-mutating identity, accounts, posts, teams, notifications, curation, analytics, inbox-count, and automation-log commands.
