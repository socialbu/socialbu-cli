package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestCurrentFallsBackToDefaultBaseURLWhenConfigValueIsBlank(t *testing.T) {
	t.Cleanup(viper.Reset)
	viper.Reset()

	home := t.TempDir()
	t.Setenv("HOME", home)
	configDir := filepath.Join(home, configDirName)
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("create config dir: %v", err)
	}
	configPath := filepath.Join(configDir, configName+".json")
	if err := os.WriteFile(configPath, []byte(`{"api_key":"abc123","base_url":""}`), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	if err := Init(); err != nil {
		t.Fatalf("init config: %v", err)
	}

	got := Current()
	if got.BaseURL != defaultBaseURL {
		t.Fatalf("BaseURL = %q, want %q", got.BaseURL, defaultBaseURL)
	}
}

func TestSaveFallsBackToDefaultBaseURLWhenConfigValueIsBlank(t *testing.T) {
	t.Cleanup(viper.Reset)
	viper.Reset()

	home := t.TempDir()
	t.Setenv("HOME", home)
	if err := Init(); err != nil {
		t.Fatalf("init config: %v", err)
	}

	if err := Save(Config{APIKey: "abc123"}); err != nil {
		t.Fatalf("save config: %v", err)
	}

	got := Current()
	if got.BaseURL != defaultBaseURL {
		t.Fatalf("BaseURL = %q, want %q", got.BaseURL, defaultBaseURL)
	}
}
