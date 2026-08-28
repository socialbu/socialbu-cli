package cmd

import (
	"strings"
	"testing"
)

func TestFixtureCaptureScriptIsRunnableWithAccountScopedAICommand(t *testing.T) {
	script := buildFixtureCaptureScript("/tmp/socialbu fixtures")
	for _, want := range []string{
		"set -euo pipefail",
		"account_id=$(python3",
		"ai autocomplete",
		"--account \"$account_id\"",
		"'/tmp/socialbu fixtures'/accounts-list.json",
		"trap 'rm -f \"$accounts_file\"' EXIT",
	} {
		if !strings.Contains(script, want) {
			t.Fatalf("script missing %q:\n%s", want, script)
		}
	}
}

func TestShellQuoteEscapesSingleQuotes(t *testing.T) {
	got := shellQuote("/tmp/usama's fixtures")
	if got != "'/tmp/usama'\\''s fixtures'" {
		t.Fatalf("shellQuote = %q", got)
	}
}

func TestFixtureCaptureCommandUsesBashCompatiblePath(t *testing.T) {
	out, err := executeRoot(t, "fixtures", "capture", "--out-dir", "test-output")
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if strings.Contains(out, "\\test-output") {
		t.Fatalf("script contains Windows path separators:\n%s", out)
	}
}
