package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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
		return fmt.Errorf("parse base URL: %w", err)
	}
	base.Path = path.Join(base.Path, strings.TrimLeft(endpoint, "/"))
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
		const maxErrorBody = 1 << 20
		payload, _ := io.ReadAll(io.LimitReader(resp.Body, maxErrorBody+1))
		truncated := len(payload) > maxErrorBody
		if truncated {
			payload = payload[:maxErrorBody]
		}
		apiErr := ErrorResponse{}
		if json.Unmarshal(payload, &apiErr) == nil && (apiErr.Message != "" || apiErr.Error != "") {
			msg := apiErr.Message
			if msg == "" {
				msg = apiErr.Error
			}
			return fmt.Errorf("api error (%s): %s", resp.Status, msg)
		}
		message := strings.TrimSpace(string(payload))
		if message == "" {
			return fmt.Errorf("api error (%s)", resp.Status)
		}
		if truncated {
			message += "…"
		}
		return fmt.Errorf("api error (%s): %s", resp.Status, message)
	}

	if out == nil {
		io.Copy(io.Discard, resp.Body)
		return nil
	}
	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()
	if err := decoder.Decode(out); err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}
