package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/socialbu/socialbu-cli/internal/config"
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

func TestAccountListRendersProductionResponseShape(t *testing.T) {
	withTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode([]map[string]any{{
			"id":     190345,
			"name":   "Main Page",
			"type":   "facebook.page",
			"_type":  "Facebook Page",
			"active": true,
		}})
	})

	out := captureStdout(t, func() {
		cmd := newRootCmd()
		cmd.SetArgs([]string{"account", "list"})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("execute: %v", err)
		}
	})

	for _, want := range []string{"Facebook Page", "facebook.page", "active"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestAccountListJSONPreservesRawAPIResponse(t *testing.T) {
	want := []map[string]any{{
		"id":               190345,
		"name":             "Main Page",
		"type":             "facebook.page",
		"attachment_types": []string{"jpg", "png"},
	}}
	withTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(want)
	})

	out := captureStdout(t, func() {
		cmd := newRootCmd()
		cmd.SetArgs([]string{"--json", "account", "list"})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("execute: %v", err)
		}
	})

	var got []map[string]any
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("decode output: %v\n%s", err, out)
	}
	if len(got) != 1 || got[0]["name"] != "Main Page" {
		t.Fatalf("unexpected output: %#v", got)
	}
	if _, ok := got[0]["attachment_types"]; !ok {
		t.Fatalf("raw account fields were discarded: %#v", got[0])
	}
}

func TestPostListUsesItemsAndPreservesRawJSON(t *testing.T) {
	withTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"items": []map[string]any{{
				"id":         73,
				"content":    "Scheduled launch post",
				"status":     "scheduled",
				"publish_at": "2026-09-01 12:00:00",
			}},
			"currentPage": 1,
			"total":       1,
		})
	})

	out := captureStdout(t, func() {
		cmd := newRootCmd()
		cmd.SetArgs([]string{"--json", "post", "list"})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("execute: %v", err)
		}
	})

	var got map[string]any
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("decode output: %v\n%s", err, out)
	}
	items, ok := got["items"].([]any)
	if !ok || len(items) != 1 {
		t.Fatalf("items were not preserved: %#v", got)
	}
	if got["total"] != float64(1) {
		t.Fatalf("pagination metadata was not preserved: %#v", got)
	}
}

func TestPostListRendersItemsInTableOutput(t *testing.T) {
	withTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"items": []map[string]any{{
				"id":         73,
				"content":    "Scheduled launch post",
				"status":     "scheduled",
				"publish_at": "2026-09-01 12:00:00",
			}},
		})
	})

	out, err := executeRoot(t, "post", "list")
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	for _, want := range []string{"73", "scheduled", "2026-09-01 12:00:00", "Scheduled launch post"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestTeamListJSONPreservesRawArray(t *testing.T) {
	withTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode([]map[string]any{{"id": 3, "name": "Growth", "owner_id": 13}})
	})

	out, err := executeRoot(t, "--json", "team", "list")
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	var got []map[string]any
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("decode output: %v\n%s", err, out)
	}
	if len(got) != 1 || got[0]["owner_id"] != float64(13) {
		t.Fatalf("raw team fields were discarded: %#v", got)
	}
}

func TestAIGenerateFromPostUsesDocumentedEndpointAndPayload(t *testing.T) {
	var seenPath string
	var seenBody map[string]any
	withTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		seenPath = r.URL.Path
		if err := json.NewDecoder(r.Body).Decode(&seenBody); err != nil {
			t.Errorf("decode body: %v", err)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"content": "Generated content"})
	})

	captureStdout(t, func() {
		cmd := newRootCmd()
		cmd.SetArgs([]string{"ai", "from-post", "--post", "42"})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("execute: %v", err)
		}
	})

	if seenPath != "/generated_content/generate_by_post" {
		t.Fatalf("path = %q", seenPath)
	}
	if seenBody["post_id"] != float64(42) {
		t.Fatalf("body = %#v", seenBody)
	}
}

func TestTeamCreateSendsRequiredAccounts(t *testing.T) {
	var seenBody map[string]any
	withTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&seenBody); err != nil {
			t.Errorf("decode body: %v", err)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"success": true})
	})

	captureStdout(t, func() {
		cmd := newRootCmd()
		cmd.SetArgs([]string{"team", "create", "Core", "--accounts", "7,9", "--require-approval"})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("execute: %v", err)
		}
	})

	accounts, ok := seenBody["accounts"].([]any)
	if !ok || len(accounts) != 2 {
		t.Fatalf("accounts = %#v", seenBody["accounts"])
	}
	if seenBody["requires_content_approval"] != true {
		t.Fatalf("requires_content_approval = %#v", seenBody["requires_content_approval"])
	}
}

func TestTruncatePreservesUTF8(t *testing.T) {
	got := truncate("🙂🙂🙂", 3)
	if !utf8.ValidString(got) {
		t.Fatalf("truncate returned invalid UTF-8: %q", got)
	}
	if got != "🙂🙂🙂" {
		t.Fatalf("truncate = %q", got)
	}

	got = truncate("🙂🙂🙂🙂", 3)
	if got != "🙂🙂…" {
		t.Fatalf("truncated value = %q", got)
	}
}
