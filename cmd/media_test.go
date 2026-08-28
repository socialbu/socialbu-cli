package cmd

import (
	"strings"
	"testing"
)

func TestDetectMimeType(t *testing.T) {
	tests := map[string]string{
		"photo.JPG":  "image/jpeg",
		"photo.jpeg": "image/jpeg",
		"image.png":  "image/png",
		"image.gif":  "image/gif",
		"image.webp": "image/webp",
		"video.mp4":  "video/mp4",
		"video.mov":  "video/quicktime",
		"file.pdf":   "application/pdf",
		"file.bin":   "application/octet-stream",
	}
	for file, want := range tests {
		if got := detectMimeType(file); got != want {
			t.Errorf("detectMimeType(%q) = %q, want %q", file, got, want)
		}
	}
}

func TestRenderMediaUploadText(t *testing.T) {
	jsonOutput = false
	out := captureStdout(t, func() {
		if err := renderMediaUpload(map[string]any{
			"name":       "image.png",
			"mime_type":  "image/png",
			"key":        "key-1",
			"secure_key": "secure-1",
			"status_result": map[string]any{
				"status":       "processed",
				"upload_token": "token-1",
			},
		}); err != nil {
			t.Fatalf("render: %v", err)
		}
	})
	for _, want := range []string{"Upload complete", "image.png", "key-1", "Upload status", "processed", "token-1"} {
		if !strings.Contains(out, want) {
			t.Fatalf("output missing %q:\n%s", want, out)
		}
	}
}
