# SocialBu CLI PR follow-up pass — 2026-04-23 02:38 PKT

## What this run changed
- Removed eager config initialization from the root command so non-API commands like `version`, `help`, `completion`, and `fixtures capture` no longer depend on `$HOME` config state.
- Routed shell completion output through Cobra's configured writer instead of hardcoded stdout.
- Added `SilenceUsage` and `SilenceErrors` on the root command so `main.go` can exit cleanly without duplicate error printing.
- Made signed media uploads use `cmd.Context()` plus a bounded HTTP client timeout instead of `context.Background()` + `http.DefaultClient`.

## Files touched
- `cmd/root.go`
- `cmd/media.go`
- `main.go`

## Verification
- `gofmt -w cmd/root.go cmd/media.go main.go`
- `go test ./...`
- `go build ./...`
- `go run . version`
- `go run . fixtures capture`

## Remaining review items worth deciding later
- Output helpers still write directly to process stdout in `internal/output/*`; fixing that cleanly wants a broader writer-plumbing pass across renderers.
- If desired before merge, we can do that writer refactor in one contained follow-up instead of piecemeal.
