package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func New(baseURL, apiKey string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  strings.TrimSpace(apiKey),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func (c *Client) Request(ctx context.Context, method, endpoint string, query url.Values, body any, out any) error {
	base, err := url.Parse(c.baseURL)
	if err != nil {
		return fmt.Errorf("parse base url: %w", err)
	}
	base.Path = path.Join(base.Path, endpoint)
	if query != nil {
		base.RawQuery = query.Encode()
	}

	var reader io.Reader
	if body != nil {
		buf := &bytes.Buffer{}
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return fmt.Errorf("encode request body: %w", err)
		}
		reader = buf
	}

	req, err := http.NewRequestWithContext(ctx, method, base.String(), reader)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("perform request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		payload, _ := io.ReadAll(resp.Body)
		apiErr := ErrorResponse{}
		if json.Unmarshal(payload, &apiErr) == nil && (apiErr.Message != "" || apiErr.Error != "") {
			msg := apiErr.Message
			if msg == "" {
				msg = apiErr.Error
			}
			return fmt.Errorf("api error (%s): %s", resp.Status, msg)
		}
		return fmt.Errorf("api error (%s): %s", resp.Status, strings.TrimSpace(string(payload)))
	}

	if out == nil {
		io.Copy(io.Discard, resp.Body)
		return nil
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}
