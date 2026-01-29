package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_Success(t *testing.T) {
	tests := []struct {
		name       string
		configJSON string
		wantErr    bool
		validate   func(*testing.T, *Config)
	}{
		{
			name: "valid config with all fields",
			configJSON: `{
				"common": {
					"block_interactive": true
				},
				"managers": {
					"pip": {
						"enabled": true,
						"require_venv": true,
						"registry_url": "https://artifactory.local/api/pypi/simple",
						"binary_path": "python3"
					},
					"npm": {
						"enabled": false,
						"registry_url": "https://artifactory.local/api/npm/npm-repo"
					}
				}
			}`,
			wantErr: false,
			validate: func(t *testing.T, cfg *Config) {
				if cfg.Common.BlockInteractive != true {
					t.Errorf("BlockInteractive = %v, want true", cfg.Common.BlockInteractive)
				}
				if cfg.Managers.Pip.Enabled != true {
					t.Errorf("Pip.Enabled = %v, want true", cfg.Managers.Pip.Enabled)
				}
				if cfg.Managers.Pip.RequireVenv != true {
					t.Errorf("Pip.RequireVenv = %v, want true", cfg.Managers.Pip.RequireVenv)
				}
				if cfg.Managers.Pip.RegistryURL != "https://artifactory.local/api/pypi/simple" {
					t.Errorf("Pip.RegistryURL = %v, want 'https://artifactory.local/api/pypi/simple'", cfg.Managers.Pip.RegistryURL)
				}
				if cfg.Managers.Pip.BinaryPath != "python3" {
					t.Errorf("Pip.BinaryPath = %v, want 'python3'", cfg.Managers.Pip.BinaryPath)
				}
				if cfg.Managers.NPM.Enabled != false {
					t.Errorf("NPM.Enabled = %v, want false", cfg.Managers.NPM.Enabled)
				}
				if cfg.Managers.NPM.RegistryURL != "https://artifactory.local/api/npm/npm-repo" {
					t.Errorf("NPM.RegistryURL = %v, want 'https://artifactory.local/api/npm/npm-repo'", cfg.Managers.NPM.RegistryURL)
				}
			},
		},
		{
			name: "valid config with minimal fields",
			configJSON: `{
				"common": {
					"block_interactive": false
				},
				"managers": {
					"pip": {
						"enabled": true
					}
				}
			}`,
			wantErr: false,
			validate: func(t *testing.T, cfg *Config) {
				if cfg.Common.BlockInteractive != false {
					t.Errorf("BlockInteractive = %v, want false", cfg.Common.BlockInteractive)
				}
				if cfg.Managers.Pip.Enabled != true {
					t.Errorf("Pip.Enabled = %v, want true", cfg.Managers.Pip.Enabled)
				}
			},
		},
		{
			name: "valid config with empty strings",
			configJSON: `{
				"common": {},
				"managers": {
					"pip": {
						"enabled": true,
						"require_venv": false,
						"registry_url": "",
						"binary_path": ""
					},
					"npm": {
						"enabled": false
					}
				}
			}`,
			wantErr: false,
			validate: func(t *testing.T, cfg *Config) {
				if cfg.Managers.Pip.RegistryURL != "" {
					t.Errorf("Pip.RegistryURL = %v, want empty string", cfg.Managers.Pip.RegistryURL)
				}
				if cfg.Managers.Pip.BinaryPath != "" {
					t.Errorf("Pip.BinaryPath = %v, want empty string", cfg.Managers.Pip.BinaryPath)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			configPath := filepath.Join(tmpDir, "config.json")

			if err := os.WriteFile(configPath, []byte(tt.configJSON), 0644); err != nil {
				t.Fatalf("Failed to write test config file: %v", err)
			}

			cfg, err := loadWithPath(configPath)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Load() expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Load() unexpected error = %v", err)
				}
				if tt.validate != nil {
					tt.validate(t, cfg)
				}
			}
		})
	}
}

func TestLoad_Error(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(t *testing.T) string
		wantErr   bool
		errorMsg  string
	}{
		{
			name: "config file does not exist",
			setupFunc: func(t *testing.T) string {
				tmpDir := t.TempDir()
				return filepath.Join(tmpDir, "nonexistent.json")
			},
			wantErr:  true,
			errorMsg: "no such file or directory",
		},
		{
			name: "malformed JSON",
			setupFunc: func(t *testing.T) string {
				tmpDir := t.TempDir()
				configPath := filepath.Join(tmpDir, "config.json")

				malformedJSON := `{invalid json`
				if err := os.WriteFile(configPath, []byte(malformedJSON), 0644); err != nil {
					t.Fatalf("Failed to write malformed config: %v", err)
				}

				return configPath
			},
			wantErr:  true,
			errorMsg: "invalid character",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configPath := tt.setupFunc(t)

			cfg, err := loadWithPath(configPath)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Load() expected error but got nil")
				}
				if tt.errorMsg != "" && err != nil {
					if !contains(err.Error(), tt.errorMsg) {
						t.Errorf("Load() error = %v, want error containing %v", err.Error(), tt.errorMsg)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Load() unexpected error = %v", err)
				}
				if cfg != nil {
					t.Logf("Config loaded successfully: %+v", cfg)
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func loadWithPath(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
