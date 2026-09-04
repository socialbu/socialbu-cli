# socialbu CLI

Use [SocialBu](https://socialbu.com) from your terminal, scripts, or AI agents. Manage connected accounts, drafts and scheduled posts, analytics, notifications, curated content, media, teams, and SocialBu AI.

## Install

### macOS and Linux

```bash
curl -fsSL https://raw.githubusercontent.com/socialbu/socialbu-cli/main/scripts/install.sh | sh
```

### Windows

```powershell
irm https://raw.githubusercontent.com/socialbu/socialbu-cli/main/scripts/install.ps1 | iex
```

Both installers select the right binary for your system and verify its checksum.

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

### Chocolatey

```powershell
choco install socialbu
```

### Go

```bash
go install github.com/socialbu/socialbu-cli/cmd/socialbu@latest
```

You can also download a binary for macOS, Linux, or Windows from [GitHub Releases](https://github.com/socialbu/socialbu-cli/releases/latest).

Check the installation:

```bash
socialbu version
```

## Connect your SocialBu account

Copy your token from **Settings > API for Developers** in SocialBu, then save it locally:

```bash
socialbu config set-key YOUR_API_TOKEN
socialbu whoami
```

See the [SocialBu API documentation](https://socialbu.com/developers/docs) if you need help finding or creating a token.

Treat the token like a password. Do not commit it, paste it into an issue, or include it in logs. `socialbu config show` confirms whether a token is configured without printing it.

## Start using the CLI

```bash
# See connected social accounts and their IDs
socialbu account list

# See scheduled posts
socialbu post list --type scheduled

# See account statistics
socialbu analytics stats

# See unread notifications
socialbu notifications unread
```

Add `--json` when another program or agent needs structured output:

```bash
socialbu account list --json
socialbu analytics stats --json
```

## Create a draft or schedule a post

Use an account ID from `socialbu account list`. This example creates a draft and does not publish it:

```bash
socialbu post create \
  --accounts 123 \
  --content "Review this before publishing" \
  --publish-at "2099-01-01 10:00:00" \
  --draft
```

Remove `--draft` only when you want SocialBu to schedule the post. `--publish-at` uses UTC and the format `YYYY-MM-DD HH:MM:SS`.

```bash
socialbu post create \
  --accounts 123,456 \
  --content "This post is ready" \
  --publish-at "2099-01-01 10:00:00"
```

## Use with AI agents and automation

For agents, CI, and scripts, pass the token through the `SOCIALBU_API_KEY` environment variable instead of storing it in a repository.

macOS and Linux:

```bash
export SOCIALBU_API_KEY="YOUR_API_TOKEN"
socialbu whoami --json
```

Windows PowerShell:

```powershell
$env:SOCIALBU_API_KEY = "YOUR_API_TOKEN"
socialbu whoami --json
```

Environment variables override local CLI configuration for the current process. Store the value in your agent or CI platform's secret manager.

You can give an agent this instruction:

```text
Use the socialbu CLI with --json. Start with `socialbu whoami` and
`socialbu account list`. Do not create, publish, upload, delete, or change
anything unless I explicitly approve it. When I ask for a post draft, use
`socialbu post create` with --draft and show me the returned post details.
Never print or repeat SOCIALBU_API_KEY.
```

## Commands

| Command | What it does |
| --- | --- |
| `socialbu account` | List connected accounts and inspect account details |
| `socialbu post` | List, inspect, create, and schedule posts |
| `socialbu analytics` | Read post, account, follower, engagement, automation, and team metrics |
| `socialbu ai` | Generate or autocomplete content with SocialBu AI |
| `socialbu notifications` | Read notifications and update their read state |
| `socialbu curation` | Browse curated topics and content |
| `socialbu media` | Upload media and check processing status |
| `socialbu team` | List, create, and delete teams |
| `socialbu config` | Manage local CLI configuration |
| `socialbu completion` | Generate shell completions for Bash, Zsh, Fish, or PowerShell |

Run `socialbu --help` to see every command, or ask for help with a specific command:

```bash
socialbu post create --help
socialbu analytics --help
```

## Support

- [SocialBu help center](https://help.socialbu.com/)
- [Report a CLI problem](https://github.com/socialbu/socialbu-cli/issues)
- [View releases and release notes](https://github.com/socialbu/socialbu-cli/releases)
