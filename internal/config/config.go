package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type CommonConfig struct {
	BlockInteractive bool `json:"block_interactive"`
}

type PipConfig struct {
	Enabled     bool   `json:"enabled"`
	RequireVenv bool   `json:"require_venv"`
	RegistryURL string `json:"registry_url"`
	BinaryPath  string `json:"binary_path"`
}

type NPMConfig struct {
	Enabled     bool   `json:"enabled"`
	RegistryURL string `json:"registry_url"`
}

type ManagersConfig struct {
	Pip PipConfig `json:"pip"`
	NPM NPMConfig `json:"npm"`
}

type Config struct {
	Common   CommonConfig   `json:"common"`
	Managers ManagersConfig `json:"managers"`
}

func getConfigPath() (string, error) {
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(os.Getenv("ProgramData"), "SafeInstall", "config.json"), nil
	case "linux":
		return "/etc/safe-install/config.json", nil
	case "darwin":
		return "/Library/Application Support/SafeInstall/config.json", nil
	default:
		return "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func Load() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, fmt.Errorf("failed to determine config path: %w", err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file at %s: %w", configPath, err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}
