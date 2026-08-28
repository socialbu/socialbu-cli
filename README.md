# SocialBu CLI

A fast Go CLI for SocialBu.

## What it does

The current CLI supports:
- `config`
- `whoami`
- `account list|get`
- `post list|get|create`
- `team list|create|delete`
- `analytics posts-count|posts-metrics|top-posts|accounts-metrics|followers|followers-growth|engagement-rate|engagement-trend|automation-logs|team-metrics|team-activity|stats`
- `ai generate|from-post|autocomplete`
- `notifications list|unread|get|mark-read|mark-unread|mark-all-read`
- `curation topics|items|get`
- `media upload|status`
- `fixtures capture`

It uses the same API key as your SocialBu account and stores config in `~/.socialbu/config.json`.

## Install

Release downloads are public. The installer scripts verify the release checksum before installing the binary.

### macOS and Linux

The installer detects your operating system and architecture, verifies the release checksum, and installs `socialbu` into `/usr/local/bin` or `~/.local/bin`:

```bash
curl -fsSL https://raw.githubusercontent.com/socialbu/socialbu-cli/main/scripts/install.sh | sh
```

### Windows

The PowerShell installer verifies the release checksum, installs `socialbu.exe` under your local application data directory, and adds it to your user `PATH`:

```powershell
irm https://raw.githubusercontent.com/socialbu/socialbu-cli/main/scripts/install.ps1 | iex
```

### Go

If you already have Go installed:

```bash
go install github.com/socialbu/socialbu-cli/cmd/socialbu@latest
```

### Homebrew

```bash
brew tap socialbu/cli https://github.com/socialbu/socialbu-cli
brew install --cask socialbu/cli/socialbu
```

### Scoop

```powershell
scoop bucket add socialbu https://github.com/socialbu/socialbu-cli
scoop install socialbu/socialbu
```

Homebrew, Scoop, WinGet, and Chocolatey packages are generated from the same checksummed release archives. WinGet and Chocolatey commands will be listed here after their public catalog submissions are accepted.

### Manual download

Download the right binary or archive from the GitHub Releases page. Raw binaries and `.tar.gz` or `.zip` archives are published for each supported platform.

```bash
chmod +x ./socialbu_linux_amd64
sudo mv ./socialbu_linux_amd64 /usr/local/bin/socialbu
```

Windows releases ship as `.exe` binaries.

## Quick start

```bash
socialbu config set-key <your-api-key>
socialbu whoami
socialbu account list
socialbu post list --type scheduled
socialbu post create --accounts 123 --content "Hello" --publish-at "2030-01-01 10:00:00" --draft
```

`publish-at` is always UTC and must use `YYYY-MM-DD HH:MM:SS`. Keep `--draft` while testing post creation.

Supported environment variables:

```bash
SOCIALBU_API_KEY
SOCIALBU_BASE_URL
```

Environment variables override stored config for the current process and are never copied into `~/.socialbu/config.json`. On macOS and Linux, the CLI stores that file with mode `0600` inside a `0700` directory.

Common write commands:

```bash
socialbu team create "Marketing" --accounts 123,456
socialbu media upload --file ./image.png
socialbu ai autocomplete --account 123 --content "Draft caption"
```

## Build from source

```bash
go build ./...
```

## Test

```bash
go test -shuffle=on -count=1 ./...
go vet ./...
go build ./...
```

On a CGO-enabled Linux or macOS environment, also run:

```bash
go test -race ./...
```

CI runs tests and builds on Linux, macOS, and Windows. Linux also runs the race detector, enforces at least 80% statement coverage, runs vet, and checks module tidiness.

## Releases

GitHub Actions + GoReleaser publish checksummed binaries and archives for each platform. Public releases also receive GitHub artifact attestations.
- macOS: amd64, arm64
- Linux: amd64, arm64
- Windows: amd64, arm64

Cut a release by pushing a semver tag:

```bash
git tag vX.Y.Z
git push origin vX.Y.Z
```

If local GitHub auth is unavailable, run the `release` workflow manually and provide:
- `tag`: the `vMAJOR.MINOR.PATCH` tag to create, for example `v0.1.2`
- `target`: the commit SHA or branch to tag, for example `94398c6c6080269f7d38104f69c7da1617d4d196`

The manual workflow creates the annotated tag in GitHub Actions, pushes it, runs tests, and publishes the GoReleaser binaries. If a rerun finds that the tag already points to the selected target, it continues instead of failing.

## Smoke workflow

`.github/workflows/smoke.yml` runs `go run . whoami` on pull requests and manual dispatch only when the `SOCIALBU_TEST_KEY` GitHub secret is configured. The workflow passes that secret through `SOCIALBU_API_KEY` without printing it.

## Fixture capture workflow

When a valid API key is available, capture real endpoint responses before changing renderer assumptions:

```bash
go run . fixtures capture > /tmp/socialbu-capture.sh
bash /tmp/socialbu-capture.sh
```

The generated script resolves the repo root automatically, writes the current fixture set into `artifacts/samples/`, and reuses `~/.socialbu/config.json` when `SOCIALBU_API_KEY` is not exported.

The smoke workflow runs `./scripts/smoke-readonly.sh` when the `SOCIALBU_TEST_KEY` repository secret is configured. The suite covers non-mutating identity, accounts, posts, teams, notifications, curation, and deployed analytics endpoints. It never creates or publishes content.
