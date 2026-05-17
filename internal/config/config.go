package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

const (
	defaultBaseURL = "https://socialbu.com/api/v1"
	configDirName  = ".socialbu"
	configName     = "config"
)

type Config struct {
	APIKey  string `mapstructure:"api_key"`
	BaseURL string `mapstructure:"base_url"`
}

func Init() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("resolve home dir: %w", err)
	}
	configDir := filepath.Join(home, configDirName)
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	viper.SetConfigName(configName)
	viper.SetConfigType("json")
	viper.AddConfigPath(configDir)
	viper.SetEnvPrefix("SOCIALBU")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.SetDefault("base_url", defaultBaseURL)
	viper.SetDefault("api_key", "")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("read config: %w", err)
		}
	}
	return nil
}

func Current() Config {
	baseURL := strings.TrimRight(strings.TrimSpace(viper.GetString("base_url")), "/")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	return Config{
		APIKey:  strings.TrimSpace(viper.GetString("api_key")),
		BaseURL: baseURL,
	}
}

func Save(cfg Config) error {
	baseURL := strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	viper.Set("api_key", strings.TrimSpace(cfg.APIKey))
	viper.Set("base_url", baseURL)
	if viper.ConfigFileUsed() == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("resolve home dir: %w", err)
		}
		configDir := filepath.Join(home, configDirName)
		if err := os.MkdirAll(configDir, 0o755); err != nil {
			return fmt.Errorf("create config dir: %w", err)
		}
		viper.SetConfigFile(filepath.Join(configDir, configName+".json"))
	}
	return viper.WriteConfigAs(viper.ConfigFileUsed())
}

func Reset() error {
	return Save(Config{BaseURL: defaultBaseURL})
}
