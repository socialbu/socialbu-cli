package cmd

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/socialbu/socialbu-cli/internal/config"
)

func TestVersionCommandShowsBuildMetadata(t *testing.T) {
	old := buildInfo
	t.Cleanup(func() { buildInfo = old })
	SetBuildInfo("v1.2.3", "abc123", "2026-08-28T00:00:00Z")

	out, err := executeRoot(t, "version")
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	for _, want := range []string{"socialbu v1.2.3", "commit: abc123", "built: 2026-08-28T00:00:00Z"} {
		if !strings.Contains(out, want) {
			t.Fatalf("output missing %q:\n%s", want, out)
		}
	}
}

func TestCompletionRejectsUnsupportedShell(t *testing.T) {
	if _, err := executeRoot(t, "completion", "nushell"); err == nil || !strings.Contains(err.Error(), "unsupported shell") {
		t.Fatalf("error = %v", err)
	}
}

func TestAPICommandRequiresKey(t *testing.T) {
	old := cfg
	t.Cleanup(func() { cfg = old })
	cfg = config.Config{BaseURL: "https://socialbu.com/api/v1"}

	if _, err := executeRoot(t, "whoami"); err == nil || !strings.Contains(err.Error(), "missing API key") {
		t.Fatalf("error = %v", err)
	}
}

func TestNoArgumentCommandsRejectUnexpectedArguments(t *testing.T) {
	tests := [][]string{
		{"whoami", "extra"},
		{"account", "list", "extra"},
		{"post", "list", "extra"},
		{"post", "create", "extra"},
		{"team", "list", "extra"},
		{"ai", "generate", "extra"},
		{"notifications", "list", "extra"},
		{"notifications", "mark-all-read", "extra"},
		{"curation", "topics", "extra"},
		{"media", "status", "extra"},
		{"analytics", "stats", "extra"},
	}
	for _, args := range tests {
		if _, err := executeRoot(t, args...); err == nil {
			t.Errorf("%q accepted an unexpected argument", strings.Join(args, " "))
		}
	}
}

func TestConfigShowDoesNotExposeAPIKey(t *testing.T) {
	old := cfg
	t.Cleanup(func() { cfg = old })
	cfg = config.Config{APIKey: "super-secret", BaseURL: "https://socialbu.com/api/v1"}

	out, err := executeRoot(t, "config", "show")
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if strings.Contains(out, "super-secret") {
		t.Fatalf("config output exposed API key: %s", out)
	}
	if !strings.Contains(out, "api_key_set: true") {
		t.Fatalf("config output = %s", out)
	}
}

func TestRootHelpHidesUndeployedInboxCommand(t *testing.T) {
	out, err := executeRoot(t, "analytics", "--help")
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if strings.Contains(out, "inbox-unread-count") {
		t.Fatalf("help advertises undeployed endpoint:\n%s", out)
	}
}

func TestConfigCommandsPersistAndResetSettings(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("SOCIALBU_API_KEY", "")
	t.Setenv("SOCIALBU_BASE_URL", "")
	old := cfg
	t.Cleanup(func() { cfg = old })
	cfg = config.Config{}

	if _, err := executeRoot(t, "config", "set-key", "stored-key"); err != nil {
		t.Fatalf("set key: %v", err)
	}
	if _, err := executeRoot(t, "config", "set-base-url", "https://example.test/api/"); err != nil {
		t.Fatalf("set base URL: %v", err)
	}

	path := filepath.Join(home, ".socialbu", "config.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	var saved config.Config
	if err := json.Unmarshal(data, &saved); err != nil {
		t.Fatalf("decode config: %v", err)
	}
	if saved.APIKey != "stored-key" || saved.BaseURL != "https://example.test/api" {
		t.Fatalf("saved config = %#v", saved)
	}

	if _, err := executeRoot(t, "config", "reset"); err != nil {
		t.Fatalf("reset: %v", err)
	}
	data, err = os.ReadFile(path)
	if err != nil {
		t.Fatalf("read reset config: %v", err)
	}
	if err := json.Unmarshal(data, &saved); err != nil {
		t.Fatalf("decode reset config: %v", err)
	}
	if saved.APIKey != "" || saved.BaseURL != "https://socialbu.com/api/v1" {
		t.Fatalf("reset config = %#v", saved)
	}
}

func TestCommandValidationHappensBeforeAPICall(t *testing.T) {
	tests := [][]string{
		{"account", "list", "--type", "invalid"},
		{"account", "get", "abc"},
		{"post", "list", "--type", "drafts"},
		{"post", "create", "--accounts", "0", "--publish-at", "bad"},
		{"team", "list", "--type", "invalid"},
		{"team", "create", "Core"},
		{"team", "delete", "0"},
		{"ai", "generate", "--type", "invalid", "--topic", "x", "--account", "1"},
		{"ai", "autocomplete"},
		{"notifications", "get", "0"},
		{"curation", "items", "--from", "2026-08-01"},
		{"curation", "get", "bad"},
		{"media", "status"},
		{"analytics", "posts-count", "--start", "2026-08-31", "--end", "2026-08-01"},
		{"analytics", "posts-metrics", "--start", "2026-08-01", "--end", "2026-08-31", "--post-type", "invalid", "--metrics", "likes"},
		{"analytics", "automation-logs", "--limit", "-1"},
		{"analytics", "team-metrics", "--start", "2026-08-01", "--end", "2026-08-31", "--metrics", "posts"},
		{"analytics", "team-metrics", "--start", "2026-08-01", "--end", "2026-08-31", "--metrics", "posts", "--accounts", "0"},
	}

	for _, args := range tests {
		t.Run(strings.Join(args, "_"), func(t *testing.T) {
			called := false
			withTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
				called = true
			})
			if _, err := executeRoot(t, args...); err == nil {
				t.Fatalf("command accepted invalid input")
			}
			if called {
				t.Fatal("API was called for invalid input")
			}
		})
	}
}
