#!/usr/bin/env bash
set -euo pipefail

if [[ -z "${SOCIALBU_API_KEY:-}" ]]; then
  echo "SOCIALBU_API_KEY is required for live smoke checks" >&2
  exit 1
fi

run_required() {
  echo "+ go run . $*"
  local out
  out=$(mktemp)
  if ! go run . "$@" >"$out" 2>&1; then
    cat "$out" >&2
    rm -f "$out"
    return 1
  fi
  rm -f "$out"
}

run_optional() {
  echo "+ go run . $*"
  local out
  out=$(mktemp)
  if ! go run . "$@" >"$out" 2>&1; then
    echo "::warning::optional smoke command failed: go run . $*" >&2
    cat "$out" >&2
  fi
  rm -f "$out"
}

# Non-mutating endpoints only. Keep create/update/delete/media/AI-generation out of CI smoke.
# Core account/post identity checks must pass; broader resource checks are optional because
# SOCIALBU_TEST_KEY permissions and account fixtures can vary across live environments.
run_required whoami
run_required --json account list
run_required --json post list

run_optional --json team list
run_optional --json notifications list
run_optional --json notifications unread
run_optional --json curation topics
run_optional --json curation items --per-page 5
run_optional --json analytics stats
run_optional --json analytics followers
run_optional --json analytics engagement-rate
run_optional --json analytics inbox-unread-count
run_optional --json analytics automation-logs --limit 5
run_optional --json analytics team-activity --limit 5
