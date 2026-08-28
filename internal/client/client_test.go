package client

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestRequestBuildsAuthenticatedJSONRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q", r.Method)
		}
		if r.URL.Path != "/api/v1/posts" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("type") != "scheduled" {
			t.Errorf("query = %q", r.URL.RawQuery)
		}
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Errorf("authorization = %q", r.Header.Get("Authorization"))
		}
		if r.Header.Get("Accept") != "application/json" || r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("headers = %#v", r.Header)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode body: %v", err)
		}
		if body["content"] != "hello" {
			t.Errorf("body = %#v", body)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"id": 42})
	}))
	defer server.Close()

	var response map[string]any
	err := New(server.URL+"/api/v1/", " test-key ").Request(
		context.Background(),
		http.MethodPost,
		"/posts",
		url.Values{"type": {"scheduled"}},
		map[string]any{"content": "hello"},
		&response,
	)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	if response["id"] != json.Number("42") {
		t.Fatalf("response = %#v", response)
	}
}

func TestRequestAcceptsSuccessfulEmptyBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	var response map[string]any
	if err := New(server.URL, "key").Request(context.Background(), http.MethodPost, "/action", nil, map[string]any{}, &response); err != nil {
		t.Fatalf("request: %v", err)
	}
}

func TestRequestFormatsAPIErrorMessages(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		contains string
	}{
		{name: "message", body: `{"message":"invalid account"}`, contains: "invalid account"},
		{name: "error", body: `{"error":"forbidden"}`, contains: "forbidden"},
		{name: "plain", body: "service unavailable", contains: "service unavailable"},
		{name: "empty", body: "", contains: "api error (400 Bad Request)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = io.WriteString(w, tt.body)
			}))
			defer server.Close()

			err := New(server.URL, "key").Request(context.Background(), http.MethodGet, "/resource", nil, nil, nil)
			if err == nil || !strings.Contains(err.Error(), tt.contains) {
				t.Fatalf("error = %v", err)
			}
			if strings.Contains(err.Error(), "Bearer") {
				t.Fatalf("error leaked authorization data: %v", err)
			}
		})
	}
}

func TestRequestOmitsAuthorizationAndContentTypeWhenUnused(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "" {
			t.Errorf("unexpected authorization header")
		}
		if r.Header.Get("Content-Type") != "" {
			t.Errorf("unexpected content type %q", r.Header.Get("Content-Type"))
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	if err := New(server.URL, " ").Request(context.Background(), http.MethodGet, "/health", nil, nil, nil); err != nil {
		t.Fatalf("request: %v", err)
	}
}

func TestRequestRejectsInvalidBaseURL(t *testing.T) {
	err := New("://bad", "key").Request(context.Background(), http.MethodGet, "/resource", nil, nil, nil)
	if err == nil || !strings.Contains(err.Error(), "parse base URL") {
		t.Fatalf("error = %v", err)
	}
}

func TestRequestHonorsCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	client := New("https://example.test", "key")
	client.httpClient = &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return nil, r.Context().Err()
	})}
	err := client.Request(ctx, http.MethodGet, "/slow", nil, nil, nil)
	if err == nil || !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v", err)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}
