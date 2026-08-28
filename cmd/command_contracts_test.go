package cmd

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestReadOnlyAPICommandContracts(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		path     string
		query    url.Values
		response string
	}{
		{name: "whoami", args: []string{"whoami"}, path: "/user", response: `{}`},
		{name: "account list", args: []string{"account", "list", "--type", "shared"}, path: "/accounts", query: url.Values{"type": {"shared"}}, response: `[]`},
		{name: "account get", args: []string{"account", "get", "7"}, path: "/accounts/7", response: `{}`},
		{name: "post list", args: []string{"post", "list", "--type", "published"}, path: "/posts", query: url.Values{"type": {"published"}}, response: `{"items":[]}`},
		{name: "post get", args: []string{"post", "get", "8"}, path: "/posts/8", response: `{}`},
		{name: "team list", args: []string{"team", "list", "--type", "joined"}, path: "/teams", query: url.Values{"type": {"joined"}}, response: `[]`},
		{name: "notification list", args: []string{"notifications", "list"}, path: "/notifications", response: `{"items":[]}`},
		{name: "notification unread", args: []string{"notifications", "unread"}, path: "/notifications/unread", response: `{"items":[]}`},
		{name: "notification get", args: []string{"notifications", "get", "9"}, path: "/notifications/9", response: `{}`},
		{name: "curation topics", args: []string{"curation", "topics", "--query", "ai"}, path: "/curation/topics", query: url.Values{"q": {"ai"}}, response: `[]`},
		{name: "curation items", args: []string{"curation", "items", "--page", "2", "--per-page", "5", "--feed-id", "3", "--search", "launch", "--from", "2026-08-01", "--to", "2026-08-31", "--authors", "A,B", "--topics", "x,y"}, path: "/curation/items", query: url.Values{"page": {"2"}, "per_page": {"5"}, "feed_id": {"3"}, "search": {"launch"}, "from": {"2026-08-01"}, "to": {"2026-08-31"}, "authors[]": {"A", "B"}, "topics[]": {"x", "y"}}, response: `{"items":[]}`},
		{name: "curation get", args: []string{"curation", "get", "10"}, path: "/curation/items/10", response: `{}`},
		{name: "analytics posts count", args: []string{"analytics", "posts-count", "--start", "2026-08-01", "--end", "2026-08-31", "--post-type", "image", "--accounts", "7,9", "--team", "3"}, path: "/insights/posts/counts", query: url.Values{"start": {"2026-08-01"}, "end": {"2026-08-31"}, "post_type": {"image"}, "accounts[]": {"7", "9"}, "team": {"3"}}, response: `{"data":[]}`},
		{name: "analytics posts metrics", args: []string{"analytics", "posts-metrics", "--start", "2026-08-01", "--end", "2026-08-31", "--post-type", "video", "--metrics", "likes,comments", "--accounts", "7", "--team", "3"}, path: "/insights/posts/metrics", query: url.Values{"start": {"2026-08-01"}, "end": {"2026-08-31"}, "post_type": {"video"}, "metrics": {"likes,comments"}, "accounts[]": {"7"}, "team": {"3"}}, response: `{}`},
		{name: "analytics top posts", args: []string{"analytics", "top-posts", "--start", "2026-08-01", "--end", "2026-08-31", "--metrics", "likes"}, path: "/insights/posts/top_posts", query: url.Values{"start": {"2026-08-01"}, "end": {"2026-08-31"}, "metrics": {"likes"}}, response: `{}`},
		{name: "analytics account metrics", args: []string{"analytics", "accounts-metrics", "--start", "2026-08-01", "--end", "2026-08-31", "--metrics", "followers", "--calculate-growth=false"}, path: "/insights/accounts/metrics", query: url.Values{"start": {"2026-08-01"}, "end": {"2026-08-31"}, "metrics": {"followers"}, "calculate_growth": {"false"}}, response: `{}`},
		{name: "analytics followers", args: []string{"analytics", "followers"}, path: "/insights/accounts/followers", response: `{}`},
		{name: "analytics followers growth", args: []string{"analytics", "followers-growth", "--start", "2026-08-01", "--end", "2026-08-31"}, path: "/insights/accounts/followers/growth", query: url.Values{"start": {"2026-08-01"}, "end": {"2026-08-31"}}, response: `{"data":[]}`},
		{name: "analytics engagement rate", args: []string{"analytics", "engagement-rate"}, path: "/insights/accounts/engagement/rate", response: `{}`},
		{name: "analytics engagement trend", args: []string{"analytics", "engagement-trend", "--start", "2026-08-01", "--end", "2026-08-31"}, path: "/insights/accounts/engagement/trend", query: url.Values{"start": {"2026-08-01"}, "end": {"2026-08-31"}}, response: `{"data":[]}`},
		{name: "analytics inbox unread", args: []string{"analytics", "inbox-unread-count"}, path: "/insights/inbox/unread-count", response: `{}`},
		{name: "analytics automation logs", args: []string{"analytics", "automation-logs", "--limit", "3"}, path: "/insights/automations/logs", query: url.Values{"limit": {"3"}}, response: `{}`},
		{name: "analytics team metrics", args: []string{"analytics", "team-metrics", "--start", "2026-08-01", "--end", "2026-08-31", "--metrics", "posts", "--accounts", "7", "--team", "3"}, path: "/insights/teams/metrics", query: url.Values{"start": {"2026-08-01"}, "end": {"2026-08-31"}, "metrics": {"posts"}, "accounts[]": {"7"}, "team": {"3"}}, response: `{"data":[]}`},
		{name: "analytics team activity", args: []string{"analytics", "team-activity", "--limit", "4"}, path: "/insights/teams/activity", query: url.Values{"limit": {"4"}}, response: `{"data":[]}`},
		{name: "analytics stats", args: []string{"analytics", "stats"}, path: "/insights/stats", response: `{}`},
		{name: "media status", args: []string{"media", "status", "--key", "media-key"}, path: "/upload_media/status", query: url.Values{"key": {"media-key"}}, response: `{}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			withTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
				called = true
				if r.Method != http.MethodGet {
					t.Errorf("method = %q", r.Method)
				}
				if r.URL.Path != tt.path {
					t.Errorf("path = %q, want %q", r.URL.Path, tt.path)
				}
				gotQuery := r.URL.Query()
				if (len(gotQuery) != 0 || len(tt.query) != 0) && !reflect.DeepEqual(gotQuery, tt.query) {
					t.Errorf("query = %#v, want %#v", gotQuery, tt.query)
				}
				_, _ = io.WriteString(w, tt.response)
			})

			args := append([]string{"--json"}, tt.args...)
			if _, err := executeRoot(t, args...); err != nil {
				t.Fatalf("execute: %v", err)
			}
			if !called {
				t.Fatal("API was not called")
			}
		})
	}
}

func TestMutatingAPICommandContracts(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		method     string
		path       string
		wantBody   map[string]any
		statusCode int
	}{
		{name: "post create", args: []string{"post", "create", "--accounts", "7,9", "--content", "Launch", "--publish-at", "2099-09-01 12:00:00", "--draft"}, method: http.MethodPost, path: "/posts", wantBody: map[string]any{"accounts": []any{float64(7), float64(9)}, "content": "Launch", "publish_at": "2099-09-01 12:00:00", "draft": true}},
		{name: "team create", args: []string{"team", "create", "Core", "--accounts", "7,9", "--require-approval"}, method: http.MethodPost, path: "/teams", wantBody: map[string]any{"name": "Core", "accounts": []any{map[string]any{"id": float64(7)}, map[string]any{"id": float64(9)}}, "requires_content_approval": true}},
		{name: "team delete", args: []string{"team", "delete", "3"}, method: http.MethodDelete, path: "/teams/3", statusCode: http.StatusNoContent},
		{name: "AI generate", args: []string{"ai", "generate", "--type", "tweet", "--topic", "launch", "--account", "7", "--team", "3"}, method: http.MethodPost, path: "/generated_content", wantBody: map[string]any{"type": "tweet", "topic": "launch", "account_id": float64(7), "team_id": float64(3)}},
		{name: "AI from post", args: []string{"ai", "from-post", "--post", "8"}, method: http.MethodPost, path: "/generated_content/generate_by_post", wantBody: map[string]any{"post_id": float64(8)}},
		{name: "AI autocomplete", args: []string{"ai", "autocomplete", "--content", "Draft", "--account", "7", "--team", "3"}, method: http.MethodPost, path: "/generated_content/autocomplete_post", wantBody: map[string]any{"content": "Draft", "account_id": float64(7), "team_id": float64(3)}},
		{name: "notification mark read", args: []string{"notifications", "mark-read", "9"}, method: http.MethodPost, path: "/notifications/9/mark_read", wantBody: map[string]any{}},
		{name: "notification mark unread", args: []string{"notifications", "mark-unread", "9"}, method: http.MethodPost, path: "/notifications/9/mark_unread", wantBody: map[string]any{}},
		{name: "notification mark all read", args: []string{"notifications", "mark-all-read"}, method: http.MethodPost, path: "/notifications/mark_all_read", wantBody: map[string]any{}, statusCode: http.StatusNoContent},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			withTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
				called = true
				if r.Method != tt.method || r.URL.Path != tt.path {
					t.Errorf("request = %s %s, want %s %s", r.Method, r.URL.Path, tt.method, tt.path)
				}
				if tt.wantBody != nil {
					var got map[string]any
					if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
						t.Errorf("decode body: %v", err)
					}
					if !reflect.DeepEqual(got, tt.wantBody) {
						t.Errorf("body = %#v, want %#v", got, tt.wantBody)
					}
				}
				if tt.statusCode != 0 {
					w.WriteHeader(tt.statusCode)
					return
				}
				_, _ = io.WriteString(w, `{"content":"generated","success":true}`)
			})

			args := append([]string{"--json"}, tt.args...)
			if _, err := executeRoot(t, args...); err != nil {
				t.Fatalf("execute: %v", err)
			}
			if !called {
				t.Fatal("API was not called")
			}
		})
	}
}

func TestMediaUploadFlow(t *testing.T) {
	var uploaded string
	storage := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("upload method = %q", r.Method)
		}
		if r.Header.Get("Content-Type") != "image/png" || r.Header.Get("x-amz-acl") != "private" {
			t.Errorf("upload headers = %#v", r.Header)
		}
		if r.ContentLength != int64(len("png-data")) {
			t.Errorf("content length = %d", r.ContentLength)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("read upload: %v", err)
		}
		uploaded = string(body)
		w.WriteHeader(http.StatusOK)
	}))
	defer storage.Close()

	withTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/upload_media":
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Errorf("decode init body: %v", err)
			}
			if body["name"] != "image.png" || body["mime_type"] != "image/png" {
				t.Errorf("init body = %#v", body)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"name": "image.png", "mime_type": "image/png", "signed_url": storage.URL, "key": "key-1", "secure_key": "secure-1"})
		case r.Method == http.MethodGet && r.URL.Path == "/upload_media/status":
			if r.URL.Query().Get("key") != "key-1" {
				t.Errorf("status query = %q", r.URL.RawQuery)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"status": "processed", "upload_token": "token-1"})
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.String())
			w.WriteHeader(http.StatusNotFound)
		}
	})

	file := t.TempDir() + string(os.PathSeparator) + "image.png"
	if err := os.WriteFile(file, []byte("png-data"), 0o600); err != nil {
		t.Fatalf("write media: %v", err)
	}
	out, err := executeRoot(t, "--json", "media", "upload", "--file", file)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if uploaded != "png-data" {
		t.Fatalf("uploaded body = %q", uploaded)
	}
	if strings.Contains(out, storage.URL) {
		t.Fatalf("output exposed the temporary signed URL: %s", out)
	}
}

func executeRoot(t *testing.T, args ...string) (string, error) {
	t.Helper()
	var executeErr error
	out := captureStdout(t, func() {
		cmd := newRootCmd()
		cmd.SetArgs(args)
		executeErr = cmd.Execute()
	})
	return out, executeErr
}
