package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/usamaejaz/socialbu-cli/internal/config"
)

func withTestAPI(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	oldCfg := cfg
	oldJSON := jsonOutput
	cfg = config.Config{APIKey: "test-key", BaseURL: server.URL}
	jsonOutput = false
	t.Cleanup(func() {
		cfg = oldCfg
		jsonOutput = oldJSON
	})

	return server
}

func TestWhoamiUsesUserEndpoint(t *testing.T) {
	seenPath := ""
	withTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		seenPath = r.URL.Path
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Fatalf("authorization header = %q", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":       42,
			"name":     "Usama",
			"email":    "usama@example.com",
			"verified": true,
		})
	})

	cmd := newRootCmd()
	cmd.SetArgs([]string{"whoami"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}

	if seenPath != "/user" {
		t.Fatalf("whoami path = %q, want /user", seenPath)
	}
}

func TestAccountListUsesPlatformColumn(t *testing.T) {
	withTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/accounts" {
			t.Fatalf("path = %q, want /accounts", r.URL.Path)
		}
		if got := r.URL.Query().Get("type"); got != "all" {
			t.Fatalf("type query = %q, want all", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"items": []map[string]any{{
				"id":       7,
				"name":     "Main Page",
				"type":     "facebook",
				"provider": "facebook_page",
				"status":   "active",
			}},
		})
	})

	out := captureStdout(t, func() {
		cmd := newRootCmd()
		cmd.SetArgs([]string{"account", "list"})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("execute: %v", err)
		}
	})

	for _, want := range []string{"PLATFORM", "facebook_page", "Main Page"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected output to contain %q, got:\n%s", want, out)
		}
	}
	if strings.Contains(out, "USERNAME") {
		t.Fatalf("output should not contain USERNAME column:\n%s", out)
	}
}
