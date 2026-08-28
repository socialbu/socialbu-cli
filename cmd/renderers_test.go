package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = w
	defer func() { os.Stdout = old }()

	fn()
	_ = w.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("copy stdout: %v", err)
	}
	_ = r.Close()
	return buf.String()
}

func TestRenderAnalyticsPostsCountUsesKnownKeys(t *testing.T) {
	jsonOutput = false
	resp := map[string]any{
		"data": []any{
			map[string]any{"date": "2026-04-20", "count": 12},
			map[string]any{"day": "2026-04-21", "posts_count": 7},
		},
	}

	out := captureStdout(t, func() {
		if err := renderAnalyticsPostsCount(resp); err != nil {
			t.Fatalf("render: %v", err)
		}
	})

	for _, want := range []string{"DATE", "POSTS", "2026-04-20", "12", "2026-04-21", "7"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestRenderNotificationDetailIncludesActorAndTarget(t *testing.T) {
	jsonOutput = false
	item := map[string]any{
		"id":         "42",
		"title":      "Post published",
		"message":    "Your queued post is live.",
		"created_at": "2026-04-21T03:00:00Z",
		"actor": map[string]any{
			"name": "Scheduler",
		},
		"target": map[string]any{
			"title": "Launch teaser",
		},
	}

	out := captureStdout(t, func() {
		if err := renderNotificationDetail(item); err != nil {
			t.Fatalf("render: %v", err)
		}
	})

	for _, want := range []string{"Notification", "actor: Scheduler", "target: Launch teaser", "title: Post published"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestRenderMediaStatusIncludesMetaFields(t *testing.T) {
	jsonOutput = false
	resp := map[string]any{
		"data": map[string]any{
			"key":          "media_123",
			"status":       "processed",
			"upload_token": "upl_tok_abc",
			"media": map[string]any{
				"width":  1080,
				"height": 1350,
				"size":   2048,
			},
		},
	}

	out := captureStdout(t, func() {
		if err := renderMediaStatus(resp); err != nil {
			t.Fatalf("render: %v", err)
		}
	})

	for _, want := range []string{"Upload status", "key: media_123", "status: processed", "upload_token: upl_tok_abc", "width: 1080", "height: 1350", "size: 2048"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestRenderCurationDetailIncludesSourceTopicsAndSummary(t *testing.T) {
	jsonOutput = false
	resp := map[string]any{
		"data": map[string]any{
			"id":      "item_1",
			"title":   "How to grow on LinkedIn",
			"summary": "A compact playbook.",
			"url":     "https://example.com/linkedin",
			"feed":    map[string]any{"name": "Growth Daily"},
			"topics":  []any{"linkedin", map[string]any{"name": "growth"}},
		},
	}

	out := captureStdout(t, func() {
		if err := renderCuration(resp); err != nil {
			t.Fatalf("render: %v", err)
		}
	})

	for _, want := range []string{"Curation item", "title: How to grow on LinkedIn", "source: Growth Daily", "topics: linkedin, growth", "summary: A compact playbook."} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestRenderAnalyticsEngagementTrendHandlesNestedSeries(t *testing.T) {
	jsonOutput = false
	resp := map[string]any{
		"data": map[string]any{
			"series": []any{
				map[string]any{"date": "2026-04-20", "engagement": 44},
				map[string]any{"day": "2026-04-21", "value": 51},
			},
		},
	}

	out := captureStdout(t, func() {
		if err := renderAnalyticsEngagementTrend(resp); err != nil {
			t.Fatalf("render: %v", err)
		}
	})

	for _, want := range []string{"DATE", "ENGAGEMENT", "2026-04-20", "44", "2026-04-21", "51"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestRenderAnalyticsFollowersGrowthUsesProductionTotal(t *testing.T) {
	jsonOutput = false
	resp := map[string]any{
		"data": []any{
			map[string]any{"date": "2026-08-28", "total_followers": 125},
		},
	}

	out := captureStdout(t, func() {
		if err := renderAnalyticsFollowersGrowth(resp); err != nil {
			t.Fatalf("render: %v", err)
		}
	})
	for _, want := range []string{"DATE", "FOLLOWERS", "2026-08-28", "125"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestRenderAnalyticsTeamActivityHandlesNestedActorName(t *testing.T) {
	jsonOutput = false
	resp := map[string]any{
		"data": map[string]any{
			"logs": []any{
				map[string]any{
					"created_at": "2026-04-21T03:00:00Z",
					"actor":      map[string]any{"name": "Usama"},
					"action":     map[string]any{"name": "approved post"},
				},
			},
		},
	}

	out := captureStdout(t, func() {
		if err := renderAnalyticsTeamActivity(resp); err != nil {
			t.Fatalf("render: %v", err)
		}
	})

	for _, want := range []string{"Recent team activity", "Usama", "approved post"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestRenderNotificationsListHandlesNestedDataResults(t *testing.T) {
	jsonOutput = false
	resp := map[string]any{
		"data": map[string]any{
			"results": []any{
				map[string]any{"id": "1", "title": "Comment", "message": "Reply posted", "created_at": "2026-04-21T05:00:00Z"},
				map[string]any{"id": "2", "subject": "Mention", "body": "You were tagged", "date": "2026-04-21T05:10:00Z"},
			},
		},
	}

	out := captureStdout(t, func() {
		if err := renderNotifications(resp); err != nil {
			t.Fatalf("render: %v", err)
		}
	})

	for _, want := range []string{"ID", "TITLE", "MESSAGE", "CREATED", "Comment", "Mention", "You were tagged"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestRenderMediaStatusHandlesNestedResultMetaFile(t *testing.T) {
	jsonOutput = false
	resp := map[string]any{
		"result": map[string]any{
			"key":    "media_789",
			"status": "processing",
			"file": map[string]any{
				"width":    1920,
				"height":   1080,
				"duration": 31,
			},
		},
	}

	out := captureStdout(t, func() {
		if err := renderMediaStatus(resp); err != nil {
			t.Fatalf("render: %v", err)
		}
	})

	for _, want := range []string{"Upload status", "key: media_789", "status: processing", "width: 1920", "height: 1080", "duration: 31"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestRenderCurationDetailHandlesNestedTopicObjects(t *testing.T) {
	jsonOutput = false
	item := map[string]any{
		"id":          "item_2",
		"title":       "Social media hooks",
		"description": "Examples of good hooks.",
		"source":      map[string]any{"name": "Creator Notes"},
		"topics": []any{
			map[string]any{"topic": map[string]any{"name": "copywriting"}},
			map[string]any{"title": "hooks"},
		},
	}

	out := captureStdout(t, func() {
		if err := renderCurationDetail(item); err != nil {
			t.Fatalf("render: %v", err)
		}
	})

	for _, want := range []string{"Curation item", "source: Creator Notes", "topics: copywriting, hooks", "description: Examples of good hooks."} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestRenderAnalyticsTeamMetricsHandlesMemberAndMetricMaps(t *testing.T) {
	jsonOutput = false
	resp := map[string]any{
		"data": map[string]any{
			"rows": []any{
				map[string]any{
					"member": map[string]any{"name": "Ayesha"},
					"metric": map[string]any{"name": "17 approvals"},
				},
			},
		},
	}

	out := captureStdout(t, func() {
		if err := renderAnalyticsTeamMetrics(resp); err != nil {
			t.Fatalf("render: %v", err)
		}
	})

	for _, want := range []string{"MEMBER", "METRIC", "Ayesha", "17 approvals"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestRenderAnalyticsTeamMetricsUsesProductionColumns(t *testing.T) {
	jsonOutput = false
	resp := map[string]any{
		"data": []any{
			map[string]any{
				"member_id":         13,
				"member_name":       "Usama",
				"total_engagements": 9,
				"scheduled_posts":   []any{map[string]any{"count": 2}, map[string]any{"count": 3}},
				"published_posts":   []any{map[string]any{"count": 4}},
				"rejected_posts":    []any{},
			},
		},
	}

	out := captureStdout(t, func() {
		if err := renderAnalyticsTeamMetrics(resp); err != nil {
			t.Fatalf("render: %v", err)
		}
	})
	for _, want := range []string{"MEMBER", "ENGAGEMENTS", "SCHEDULED", "PUBLISHED", "REJECTED", "Usama", "9", "5", "4", "0"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestRenderCurationDetailHandlesGraphNodeTopics(t *testing.T) {
	jsonOutput = false
	item := map[string]any{
		"id":      "item_3",
		"title":   "Repurpose a webinar",
		"summary": "Turn one talk into many clips.",
		"source":  map[string]any{"name": "Content Engine"},
		"topics": []any{
			map[string]any{"node": map[string]any{"title": "repurposing"}},
			map[string]any{"attributes": map[string]any{"name": "video"}},
		},
	}

	out := captureStdout(t, func() {
		if err := renderCurationDetail(item); err != nil {
			t.Fatalf("render: %v", err)
		}
	})

	for _, want := range []string{"Curation item", "source: Content Engine", "topics: repurposing, video", "summary: Turn one talk into many clips."} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestRenderNotificationsListHandlesNestedEventAndPayload(t *testing.T) {
	jsonOutput = false
	resp := map[string]any{
		"data": map[string]any{
			"notifications": []any{
				map[string]any{
					"id":         "3",
					"event":      map[string]any{"name": "approval_requested"},
					"payload":    map[string]any{"message": "Please review this post."},
					"timestamps": map[string]any{"created_at": "2026-04-21T05:20:00Z"},
				},
			},
		},
	}

	out := captureStdout(t, func() {
		if err := renderNotifications(resp); err != nil {
			t.Fatalf("render: %v", err)
		}
	})

	for _, want := range []string{"Notification", "event: approval_requested", "message: Please review this post.", "created_at: 2026-04-21T05:20:00Z"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestRenderAIResponseHandlesNestedResultContentAndMeta(t *testing.T) {
	jsonOutput = false
	resp := map[string]any{
		"data": map[string]any{
			"result": map[string]any{
				"content": "Try a stronger hook, then lead with the customer pain.",
				"model":   "gpt-x",
			},
			"meta": map[string]any{
				"tone":     "direct",
				"language": "en",
			},
		},
	}

	out := captureStdout(t, func() {
		if err := renderAIResponse(resp); err != nil {
			t.Fatalf("render: %v", err)
		}
	})

	for _, want := range []string{"AI output", "Try a stronger hook, then lead with the customer pain.", "tone: direct", "language: en", "model: gpt-x"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestRenderCurationHandlesTopLevelTopicArray(t *testing.T) {
	jsonOutput = false
	resp := []any{
		map[string]any{"id": 1, "name": "Marketing"},
		map[string]any{"slug": "ai", "title": "AI"},
	}

	out := captureStdout(t, func() {
		if err := renderCurationAny(resp); err != nil {
			t.Fatalf("render: %v", err)
		}
	})

	for _, want := range []string{"ID", "TITLE", "SOURCE", "PUBLISHED", "1", "Marketing", "ai", "AI"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestRenderCurationUsesProductionAuthorsField(t *testing.T) {
	jsonOutput = false
	resp := map[string]any{
		"items": []any{
			map[string]any{"id": 1, "title": "One", "authors": "Author One"},
			map[string]any{"id": 2, "title": "Two", "authors": "Author Two"},
		},
	}

	out := captureStdout(t, func() {
		if err := renderCuration(resp); err != nil {
			t.Fatalf("render: %v", err)
		}
	})
	for _, want := range []string{"Author One", "Author Two"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestEmptyCollectionsRenderHonestEmptyStates(t *testing.T) {
	jsonOutput = false
	notifications := captureStdout(t, func() {
		if err := renderNotifications(map[string]any{"items": []any{}}); err != nil {
			t.Fatalf("render notifications: %v", err)
		}
	})
	if !strings.Contains(notifications, "No notifications.") {
		t.Fatalf("notification output = %q", notifications)
	}

	curation := captureStdout(t, func() {
		if err := renderCuration(map[string]any{"items": []any{}}); err != nil {
			t.Fatalf("render curation: %v", err)
		}
	})
	if !strings.Contains(curation, "No curation items.") {
		t.Fatalf("curation output = %q", curation)
	}
}

func TestRenderAnalyticsStatsFlattensNestedMetricMaps(t *testing.T) {
	jsonOutput = false
	resp := map[string]any{
		"data": map[string]any{
			"accounts":        7,
			"scheduled_posts": map[string]any{"count": 31},
			"engagement":      map[string]any{"current": map[string]any{"value": 904}},
			"followers":       map[string]any{"totals": map[string]any{"count": 12043}},
		},
	}

	out := captureStdout(t, func() {
		if err := renderAnalyticsStats(resp); err != nil {
			t.Fatalf("render: %v", err)
		}
	})

	for _, want := range []string{"Stats", "accounts: 7", "scheduled_posts: 31", "engagement: 904", "followers: 12043"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}
