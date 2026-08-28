package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestInitFallsBackToDefaultBaseURLWhenConfigValueIsBlank(t *testing.T) {
	home := useTempHome(t)
	writeTestConfig(t, home, Config{APIKey: "abc123"}, 0o600)

	if err := Init(); err != nil {
		t.Fatalf("init config: %v", err)
	}

	got := Current()
	if got.BaseURL != defaultBaseURL {
		t.Fatalf("BaseURL = %q, want %q", got.BaseURL, defaultBaseURL)
	}
	if got.APIKey != "abc123" {
		t.Fatalf("APIKey = %q", got.APIKey)
	}
}

func TestInitDoesNotCreateConfigDirectory(t *testing.T) {
	home := useTempHome(t)
	if err := Init(); err != nil {
		t.Fatalf("init config: %v", err)
	}
	if _, err := os.Stat(filepath.Join(home, configDirName)); !os.IsNotExist(err) {
		t.Fatalf("read-only init created config directory: %v", err)
	}
}

func TestSaveUsesSecurePermissionsAndDefaultBaseURL(t *testing.T) {
	home := useTempHome(t)
	if err := Init(); err != nil {
		t.Fatalf("init config: %v", err)
	}
	if err := Save(Config{APIKey: "abc123"}); err != nil {
		t.Fatalf("save config: %v", err)
	}

	path := filepath.Join(home, configDirName, configName)
	got := readTestConfig(t, path)
	if got.BaseURL != defaultBaseURL || got.APIKey != "abc123" {
		t.Fatalf("saved config = %#v", got)
	}
	if runtime.GOOS != "windows" {
		if info, err := os.Stat(path); err != nil {
			t.Fatalf("stat config: %v", err)
		} else if perm := info.Mode().Perm(); perm != 0o600 {
			t.Fatalf("config permissions = %o, want 600", perm)
		}
		if info, err := os.Stat(filepath.Dir(path)); err != nil {
			t.Fatalf("stat config dir: %v", err)
		} else if perm := info.Mode().Perm(); perm != 0o700 {
			t.Fatalf("config directory permissions = %o, want 700", perm)
		}
	}
}

func TestEnvironmentOverridesAreNotPersistedByPartialUpdates(t *testing.T) {
	home := useTempHome(t)
	writeTestConfig(t, home, Config{APIKey: "disk-key", BaseURL: "https://disk.example/api"}, 0o600)
	t.Setenv("SOCIALBU_API_KEY", "env-key")
	t.Setenv("SOCIALBU_BASE_URL", "https://env.example/api")
	if err := Init(); err != nil {
		t.Fatalf("init config: %v", err)
	}

	if got := Current(); got.APIKey != "env-key" || got.BaseURL != "https://env.example/api" {
		t.Fatalf("effective config = %#v", got)
	}
	if err := SetBaseURL("https://new.example/api/"); err != nil {
		t.Fatalf("set base URL: %v", err)
	}

	path := filepath.Join(home, configDirName, configName)
	saved := readTestConfig(t, path)
	if saved.APIKey != "disk-key" {
		t.Fatalf("environment API key was persisted: %#v", saved)
	}
	if saved.BaseURL != "https://new.example/api" {
		t.Fatalf("BaseURL = %q", saved.BaseURL)
	}

	if err := SetAPIKey("new-disk-key"); err != nil {
		t.Fatalf("set API key: %v", err)
	}
	saved = readTestConfig(t, path)
	if saved.APIKey != "new-disk-key" || saved.BaseURL != "https://new.example/api" {
		t.Fatalf("saved config = %#v", saved)
	}
}

func TestSetBaseURLRejectsUnsafeOrInvalidValues(t *testing.T) {
	useTempHome(t)
	if err := Init(); err != nil {
		t.Fatalf("init config: %v", err)
	}
	for _, value := range []string{
		"socialbu.com/api/v1",
		"ftp://socialbu.com/api/v1",
		"https://user:pass@socialbu.com/api/v1",
		"https://socialbu.com/api/v1?token=secret",
		"https://socialbu.com/api/v1#fragment",
	} {
		if err := SetBaseURL(value); err == nil {
			t.Errorf("SetBaseURL(%q) succeeded", value)
		}
	}
}

func TestSetAPIKeyRejectsBlankValue(t *testing.T) {
	useTempHome(t)
	if err := Init(); err != nil {
		t.Fatalf("init config: %v", err)
	}
	if err := SetAPIKey("  "); err == nil {
		t.Fatal("blank API key was accepted")
	}
}

func TestResetClearsStoredAPIKey(t *testing.T) {
	home := useTempHome(t)
	writeTestConfig(t, home, Config{APIKey: "abc123", BaseURL: "https://example.com/api"}, 0o600)
	if err := Init(); err != nil {
		t.Fatalf("init config: %v", err)
	}
	if err := Reset(); err != nil {
		t.Fatalf("reset config: %v", err)
	}

	got := readTestConfig(t, filepath.Join(home, configDirName, configName))
	if got.APIKey != "" || got.BaseURL != defaultBaseURL {
		t.Fatalf("reset config = %#v", got)
	}
}

func TestInitRejectsMalformedConfig(t *testing.T) {
	home := useTempHome(t)
	dir := filepath.Join(home, configDirName)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatalf("create config dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, configName), []byte("{"), 0o600); err != nil {
		t.Fatalf("write malformed config: %v", err)
	}
	if err := Init(); err == nil || !strings.Contains(err.Error(), "decode config") {
		t.Fatalf("Init error = %v", err)
	}
}

func useTempHome(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("SOCIALBU_API_KEY", "")
	t.Setenv("SOCIALBU_BASE_URL", "")
	return home
}

func writeTestConfig(t *testing.T, home string, cfg Config, mode os.FileMode) {
	t.Helper()
	dir := filepath.Join(home, configDirName)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatalf("create config dir: %v", err)
	}
	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("marshal config: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, configName), data, mode); err != nil {
		t.Fatalf("write config: %v", err)
	}
}

func readTestConfig(t *testing.T, path string) Config {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("decode config: %v", err)
	}
	return cfg
}
