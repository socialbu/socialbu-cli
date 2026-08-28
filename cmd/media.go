package cmd

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/socialbu/socialbu-cli/internal/output"
	"github.com/spf13/cobra"
)

const mediaUploadTimeout = 2 * time.Minute

type mediaUploadInitResponse struct {
	Name      string `json:"name"`
	MimeType  string `json:"mime_type"`
	SignedURL string `json:"signed_url"`
	Key       string `json:"key"`
	SecureKey string `json:"secure_key"`
}

func newMediaCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "media", Short: "Upload media and check upload status"}
	cmd.AddCommand(newMediaUploadCmd())
	cmd.AddCommand(newMediaStatusCmd())
	return cmd
}

func newMediaUploadCmd() *cobra.Command {
	var filePath, mimeType string
	cmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload a local file through SocialBu's signed URL flow",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(filePath) == "" {
				return fmt.Errorf("--file is required")
			}
			info, err := os.Stat(filePath)
			if err != nil {
				return err
			}
			if !info.Mode().IsRegular() {
				return fmt.Errorf("--file must point to a regular file")
			}
			if mimeType == "" {
				mimeType = detectMimeType(filePath)
			}
			cli, err := apiClient()
			if err != nil {
				return err
			}
			payload := map[string]any{"name": filepath.Base(filePath), "mime_type": mimeType}
			var initResp mediaUploadInitResponse
			if err := cli.Request(cmd.Context(), "POST", "/upload_media", nil, payload, &initResp); err != nil {
				return err
			}
			if strings.TrimSpace(initResp.SignedURL) == "" || strings.TrimSpace(initResp.Key) == "" {
				return fmt.Errorf("media upload initialization response is missing signed_url or key")
			}
			if err := putSignedMedia(cmd.Context(), initResp.SignedURL, filePath, mimeType, info.Size()); err != nil {
				return err
			}
			var statusResp map[string]any
			if err := cli.Request(cmd.Context(), "GET", "/upload_media/status", addOptionalString(url.Values{}, "key", initResp.Key), nil, &statusResp); err != nil {
				return err
			}
			result := map[string]any{
				"name":          initResp.Name,
				"mime_type":     initResp.MimeType,
				"key":           initResp.Key,
				"secure_key":    initResp.SecureKey,
				"status_result": statusResp,
			}
			return renderMediaUpload(result)
		},
	}
	cmd.Flags().StringVar(&filePath, "file", "", "local file path")
	cmd.Flags().StringVar(&mimeType, "mime-type", "", "explicit MIME type")
	return cmd
}

func newMediaStatusCmd() *cobra.Command {
	var key string
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Check uploaded media status and retrieve upload token",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if key == "" {
				return fmt.Errorf("--key is required")
			}
			cli, err := apiClient()
			if err != nil {
				return err
			}
			var resp map[string]any
			if err := cli.Request(cmd.Context(), "GET", "/upload_media/status", addOptionalString(url.Values{}, "key", key), nil, &resp); err != nil {
				return err
			}
			return renderMediaStatus(resp)
		},
	}
	cmd.Flags().StringVar(&key, "key", "", "media key from init step")
	return cmd
}

func putSignedMedia(ctx context.Context, signedURL, filePath, mimeType string, size int64) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, signedURL, file)
	if err != nil {
		return fmt.Errorf("build upload request: %w", err)
	}
	req.Header.Set("Content-Type", mimeType)
	req.Header.Set("x-amz-acl", "private")
	req.ContentLength = size

	resp, err := (&http.Client{Timeout: mediaUploadTimeout}).Do(req)
	if err != nil {
		return fmt.Errorf("upload media: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload media failed (%s): %s", resp.Status, strings.TrimSpace(string(body)))
	}
	return nil
}

func detectMimeType(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".mp4":
		return "video/mp4"
	case ".mov":
		return "video/quicktime"
	case ".pdf":
		return "application/pdf"
	default:
		return "application/octet-stream"
	}
}

func renderMediaUpload(resp map[string]any) error {
	if jsonOutput {
		return output.JSON(resp)
	}
	output.PrintSection("Upload complete")
	fmt.Printf("Name: %s\n", output.StringFromMap(resp, "name"))
	fmt.Printf("MIME type: %s\n", output.StringFromMap(resp, "mime_type"))
	fmt.Printf("Key: %s\n", output.StringFromMap(resp, "key"))
	fmt.Printf("Secure key: %s\n", output.StringFromMap(resp, "secure_key"))
	if status, ok := resp["status_result"].(map[string]any); ok {
		fmt.Println()
		return renderMediaStatus(status)
	}
	return nil
}

func renderMediaStatus(resp map[string]any) error {
	if jsonOutput {
		return output.JSON(resp)
	}
	data := resp
	if nested := output.MapFromMap(resp, "data", "result", "upload"); nested != nil {
		data = nested
	}
	values := map[string]string{}
	for _, key := range []string{"key", "secure_key", "status", "mime_type", "name", "url", "upload_token", "message"} {
		if value := output.StringFromMap(data, key); value != "" {
			values[key] = value
		}
	}
	if meta := output.MapFromMap(data, "meta", "media", "file"); meta != nil {
		for _, key := range []string{"width", "height", "duration", "size"} {
			if value := output.StringFromMap(meta, key); value != "" {
				values[key] = value
			}
		}
	}
	if len(values) == 0 {
		return output.JSON(resp)
	}
	output.KeyValue("Upload status", values)
	return nil
}
