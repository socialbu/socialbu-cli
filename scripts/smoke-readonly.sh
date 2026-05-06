#!/usr/bin/env bash
set -euo pipefail

if [[ -z "${SOCIALBU_API_KEY:-}" ]]; then
  echo "SOCIALBU_API_KEY is required for live smoke checks" >&2
  exit 1
fi

run() {
  echo "+ go run . $*"
  go run . "$@" >/tmp/socialbu-cli-smoke.out
}

# Non-mutating endpoints only. Keep create/update/delete/media/AI-generation out of CI smoke.
run whoami
run --json account list
run --json post list
run --json team list
run --json notifications list
run --json notifications unread
run --json curation topics
run --json curation items --per-page 5
run --json analytics stats
run --json analytics team-activity --limit 5
