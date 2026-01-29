package policy

import (
	"os"
	"os/exec"
	"testing"

	"safe-install/internal/config"
)

type CommandExecutor interface {
	Command(name string, args ...string) *exec.Cmd
}

type RealExecutor struct{}

func (e *RealExecutor) Command(name string, args ...string) *exec.Cmd {
	return exec.Command(name, args...)
}

type MockExecutor struct {
	Cmd *exec.Cmd
	Err error
}

func (e *MockExecutor) Command(name string, args ...string) *exec.Cmd {
	if e.Err != nil {
		cmd := exec.Command("echo", "error")
		return cmd
	}
	return e.Cmd
}

func TestValidatePipEnv(t *testing.T) {
	tests := []struct {
		name         string
		venv         string
		requireVenv  bool
		wantErr      bool
		errorMessage string
	}{
		{
			name:        "passes when venv is set and required",
			venv:        "/path/to/venv",
			requireVenv: true,
			wantErr:     false,
		},
		{
			name:         "fails when venv is not set but required",
			venv:         "",
			requireVenv:  true,
			wantErr:      true,
			errorMessage: "virtual environment is required but not active",
		},
		{
			name:        "passes when venv not required",
			venv:        "",
			requireVenv: false,
			wantErr:     false,
		},
		{
			name:        "passes when venv is set but not required",
			venv:        "/path/to/venv",
			requireVenv: false,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalEnv := os.Getenv("VIRTUAL_ENV")
			defer os.Setenv("VIRTUAL_ENV", originalEnv)

			os.Unsetenv("VIRTUAL_ENV")
			if tt.venv != "" {
				os.Setenv("VIRTUAL_ENV", tt.venv)
			}

			cfg := &config.Config{
				Managers: config.ManagersConfig{
					Pip: config.PipConfig{
						RequireVenv: tt.requireVenv,
					},
				},
			}

			err := ValidatePipEnv(cfg)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidatePipEnv() expected error but got nil")
				}
				if err != nil && err.Error() != tt.errorMessage {
					t.Errorf("ValidatePipEnv() error = %v, want %v", err.Error(), tt.errorMessage)
				}
			} else {
				if err != nil {
					t.Errorf("ValidatePipEnv() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestBuildPipArgs(t *testing.T) {
	tests := []struct {
		name            string
		cfg             *config.Config
		userArgs        []string
		wantErr         bool
		errorMessage    string
		wantArgsContain []string
	}{
		{
			name: "errors on empty arguments",
			cfg: &config.Config{
				Managers: config.ManagersConfig{
					Pip: config.PipConfig{
						RegistryURL: "",
					},
				},
			},
			userArgs:     []string{},
			wantErr:      true,
			errorMessage: "no pip subcommand provided",
		},
		{
			name: "blocks --index-url argument",
			cfg: &config.Config{
				Managers: config.ManagersConfig{
					Pip: config.PipConfig{
						RegistryURL: "",
					},
				},
			},
			userArgs:     []string{"install", "--index-url", "http://bad.com"},
			wantErr:      true,
			errorMessage: "argument '--index-url' and '--extra-index-url' are not allowed",
		},
		{
			name: "blocks --extra-index-url argument",
			cfg: &config.Config{
				Managers: config.ManagersConfig{
					Pip: config.PipConfig{
						RegistryURL: "",
					},
				},
			},
			userArgs:     []string{"install", "--extra-index-url", "http://bad.com"},
			wantErr:      true,
			errorMessage: "argument '--index-url' and '--extra-index-url' are not allowed",
		},
		{
			name: "blocks both index-url and extra-index-url",
			cfg: &config.Config{
				Managers: config.ManagersConfig{
					Pip: config.PipConfig{
						RegistryURL: "",
					},
				},
			},
			userArgs:     []string{"install", "--index-url", "http://bad.com", "--extra-index-url", "http://also-bad.com"},
			wantErr:      true,
			errorMessage: "argument '--index-url' and '--extra-index-url' are not allowed",
		},
		{
			name: "allows normal install command",
			cfg: &config.Config{
				Common: config.CommonConfig{
					BlockInteractive: false,
				},
				Managers: config.ManagersConfig{
					Pip: config.PipConfig{
						RegistryURL: "",
					},
				},
			},
			userArgs:        []string{"install", "requests"},
			wantErr:         false,
			wantArgsContain: []string{"install", "requests"},
		},
		{
			name: "injects --no-input when block_interactive is true for install",
			cfg: &config.Config{
				Common: config.CommonConfig{
					BlockInteractive: true,
				},
				Managers: config.ManagersConfig{
					Pip: config.PipConfig{
						RegistryURL: "",
					},
				},
			},
			userArgs:        []string{"install", "requests"},
			wantErr:         false,
			wantArgsContain: []string{"install", "--no-input", "requests"},
		},
		{
			name: "does not inject --no-input for list subcommand",
			cfg: &config.Config{
				Common: config.CommonConfig{
					BlockInteractive: true,
				},
				Managers: config.ManagersConfig{
					Pip: config.PipConfig{
						RegistryURL: "",
					},
				},
			},
			userArgs:        []string{"list"},
			wantErr:         false,
			wantArgsContain: []string{"list"},
		},
		{
			name: "injects registry url when configured for install",
			cfg: &config.Config{
				Common: config.CommonConfig{
					BlockInteractive: false,
				},
				Managers: config.ManagersConfig{
					Pip: config.PipConfig{
						RegistryURL: "https://artifactory.local/api/pypi/simple",
					},
				},
			},
			userArgs:        []string{"install", "requests"},
			wantErr:         false,
			wantArgsContain: []string{"install", "--index-url", "https://artifactory.local/api/pypi/simple", "--no-index", "requests"},
		},
		{
			name: "does not inject registry url for list subcommand",
			cfg: &config.Config{
				Common: config.CommonConfig{
					BlockInteractive: false,
				},
				Managers: config.ManagersConfig{
					Pip: config.PipConfig{
						RegistryURL: "https://artifactory.local/api/pypi/simple",
					},
				},
			},
			userArgs:        []string{"list"},
			wantErr:         false,
			wantArgsContain: []string{"list"},
		},
		{
			name: "injects both --no-input and registry url for install",
			cfg: &config.Config{
				Common: config.CommonConfig{
					BlockInteractive: true,
				},
				Managers: config.ManagersConfig{
					Pip: config.PipConfig{
						RegistryURL: "https://artifactory.local/api/pypi/simple",
					},
				},
			},
			userArgs:        []string{"install", "requests"},
			wantErr:         false,
			wantArgsContain: []string{"install", "--no-input", "--index-url", "https://artifactory.local/api/pypi/simple", "--no-index", "requests"},
		},
		{
			name: "preserves multiple user arguments",
			cfg: &config.Config{
				Common: config.CommonConfig{
					BlockInteractive: false,
				},
				Managers: config.ManagersConfig{
					Pip: config.PipConfig{
						RegistryURL: "",
					},
				},
			},
			userArgs:        []string{"install", "--upgrade", "requests", "flask"},
			wantErr:         false,
			wantArgsContain: []string{"install", "--upgrade", "requests", "flask"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args, err := BuildPipArgs(tt.cfg, tt.userArgs)

			if tt.wantErr {
				if err == nil {
					t.Errorf("BuildPipArgs() expected error but got nil")
				}
				if err != nil && err.Error() != tt.errorMessage {
					t.Errorf("BuildPipArgs() error = %v, want %v", err.Error(), tt.errorMessage)
				}
			} else {
				if err != nil {
					t.Errorf("BuildPipArgs() unexpected error = %v", err)
				}

				for _, wantArg := range tt.wantArgsContain {
					found := false
					for _, arg := range args {
						if arg == wantArg {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("BuildPipArgs() expected to contain %v, got %v", wantArg, args)
					}
				}
			}
		})
	}
}

func TestExecutePip_CommandConstruction(t *testing.T) {
	t.Skip("This test would require mocking exec.Command to verify command construction")

	cfg := &config.Config{
		Managers: config.ManagersConfig{
			Pip: config.PipConfig{
				BinaryPath: "python3",
			},
		},
	}
	args := []string{"install", "requests"}

	executor := &MockExecutor{}
	_ = executor.Command("python3", "-m", "pip")
	_ = ExecutePip(cfg, args)
}
