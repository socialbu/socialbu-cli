package config

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	defaultBaseURL = "https://socialbu.com/api/v1"
	configDirName  = ".socialbu"
	configName     = "config.json"
)

type Config struct {
	APIKey  string `json:"api_key"`
	BaseURL string `json:"base_url"`
}

var (
	mu     sync.RWMutex
	stored = Config{BaseURL: defaultBaseURL}
)

func Init() error {
	path, err := filePath()
	if err != nil {
		return err
	}

	cfg := Config{BaseURL: defaultBaseURL}
	data, err := os.ReadFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("read config: %w", err)
		}
	} else {
		if err := json.Unmarshal(data, &cfg); err != nil {
			return fmt.Errorf("decode config: %w", err)
		}
		if err := os.Chmod(path, 0o600); err != nil {
			return fmt.Errorf("secure config file: %w", err)
		}
	}

	cfg.APIKey = strings.TrimSpace(cfg.APIKey)
	cfg.BaseURL = normalizeBaseURL(cfg.BaseURL)
	if err := validateBaseURL(cfg.BaseURL); err != nil {
		return fmt.Errorf("invalid base URL in config: %w", err)
	}

	mu.Lock()
	stored = cfg
	mu.Unlock()
	return nil
}

func Current() Config {
	mu.RLock()
	cfg := stored
	mu.RUnlock()

	if value := strings.TrimSpace(os.Getenv("SOCIALBU_API_KEY")); value != "" {
		cfg.APIKey = value
	}
	if value := strings.TrimSpace(os.Getenv("SOCIALBU_BASE_URL")); value != "" {
		cfg.BaseURL = value
	}
	cfg.BaseURL = normalizeBaseURL(cfg.BaseURL)
	return cfg
}

func Save(cfg Config) error {
	cfg.APIKey = strings.TrimSpace(cfg.APIKey)
	cfg.BaseURL = normalizeBaseURL(cfg.BaseURL)
	if err := validateBaseURL(cfg.BaseURL); err != nil {
		return err
	}

	path, err := filePath()
	if err != nil {
		return err
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}
	if err := os.Chmod(dir, 0o700); err != nil {
		return fmt.Errorf("secure config directory: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("encode config: %w", err)
	}
	data = append(data, '\n')

	temp, err := os.CreateTemp(dir, "config-*.tmp")
	if err != nil {
		return fmt.Errorf("create temporary config: %w", err)
	}
	tempPath := temp.Name()
	defer os.Remove(tempPath)

	if err := temp.Chmod(0o600); err != nil {
		temp.Close()
		return fmt.Errorf("secure temporary config: %w", err)
	}
	if _, err := temp.Write(data); err != nil {
		temp.Close()
		return fmt.Errorf("write temporary config: %w", err)
	}
	if err := temp.Sync(); err != nil {
		temp.Close()
		return fmt.Errorf("sync temporary config: %w", err)
	}
	if err := temp.Close(); err != nil {
		return fmt.Errorf("close temporary config: %w", err)
	}
	if err := os.Rename(tempPath, path); err != nil {
		return fmt.Errorf("replace config: %w", err)
	}

	mu.Lock()
	stored = cfg
	mu.Unlock()
	return nil
}

func SetAPIKey(apiKey string) error {
	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return fmt.Errorf("API key cannot be empty")
	}
	mu.RLock()
	cfg := stored
	mu.RUnlock()
	cfg.APIKey = apiKey
	return Save(cfg)
}

func SetBaseURL(baseURL string) error {
	baseURL = normalizeBaseURL(baseURL)
	if err := validateBaseURL(baseURL); err != nil {
		return err
	}
	mu.RLock()
	cfg := stored
	mu.RUnlock()
	cfg.BaseURL = baseURL
	return Save(cfg)
}

func Reset() error {
	return Save(Config{BaseURL: defaultBaseURL})
}

func filePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home directory: %w", err)
	}
	return filepath.Join(home, configDirName, configName), nil
}

func normalizeBaseURL(value string) string {
	value = strings.TrimRight(strings.TrimSpace(value), "/")
	if value == "" {
		return defaultBaseURL
	}
	return value
}

func validateBaseURL(value string) error {
	parsed, err := url.Parse(value)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "https" && parsed.Scheme != "http") {
		return fmt.Errorf("base URL must be an absolute http:// or https:// URL")
	}
	if parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" {
		return fmt.Errorf("base URL cannot contain credentials, a query, or a fragment")
	}
	return nil
}
