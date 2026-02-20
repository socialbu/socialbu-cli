# socialbu-cli

Command-line interface for managing your [SocialBu](https://socialbu.com) accounts, posts, analytics, and more — right from the terminal.

## Installation

```bash
npm install -g socialbu-cli
```

## Quick Start

```bash
# Set your API key
socialbu config set-key <your-api-key>

# Check your identity
socialbu whoami

# List connected social accounts
socialbu account list

# Create a post
socialbu post create --content "Hello from the CLI!" --accounts 123 --publish-at "2025-04-14 15:30:00"

# List scheduled posts
socialbu post list --type scheduled

# Generate AI content
socialbu ai generate --topic "Mother's Day" --type tweet

# View analytics
socialbu analytics stats
socialbu analytics followers
```

## Commands

| Command | Description |
|---------|-------------|
| `socialbu config set-key <key>` | Store your API key |
| `socialbu config get-key` | Show stored API key |
| `socialbu config reset` | Remove stored config |
| `socialbu whoami` | Show current user info |
| `socialbu account list` | List social accounts |
| `socialbu account get <id>` | Get account details |
| `socialbu account delete <id>` | Delete an account |
| `socialbu post create` | Create a new post |
| `socialbu post list` | List posts |
| `socialbu post get <id>` | Get post details |
| `socialbu post delete <id>` | Delete a post |
| `socialbu ai generate` | Generate AI content |
| `socialbu ai autocomplete` | Autocomplete post content |
| `socialbu team list` | List teams |
| `socialbu team create` | Create a team |
| `socialbu team delete <id>` | Delete a team |
| `socialbu analytics stats` | User stats overview |
| `socialbu analytics followers` | Followers count |
| `socialbu analytics engagement` | Engagement rate |
| `socialbu notifications list` | List notifications |
| `socialbu notifications unread` | List unread notifications |
| `socialbu curate topics` | List curation topics |
| `socialbu curate items` | List curated items |

## Global Options

All commands support:
- `--json` — Output raw JSON instead of formatted tables
- `--help` — Show help for any command

## Authentication

Get your API token from SocialBu (Settings → API) or generate one:

```bash
socialbu config set-key YOUR_BEARER_TOKEN
```

The key is stored in `~/.socialbu/config.json`.

## API Reference

Based on the [SocialBu API v1.0.0](https://socialbu.com/developers/docs).

## License

MIT
