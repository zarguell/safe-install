package policy

import (
	"fmt"
	"os"
	"os/exec"

	"safe-install/internal/config"
)

func ValidatePipEnv(cfg *config.Config) error {
	if cfg.Managers.Pip.RequireVenv {
		if os.Getenv("VIRTUAL_ENV") == "" {
			return fmt.Errorf("virtual environment is required but not active")
		}
	}
	return nil
}

func BuildPipArgs(cfg *config.Config, userArgs []string) ([]string, error) {
	if len(userArgs) == 0 {
		return nil, fmt.Errorf("no pip subcommand provided")
	}

	for _, arg := range userArgs {
		if arg == "--index-url" || arg == "--extra-index-url" {
			return nil, fmt.Errorf("argument '--index-url' and '--extra-index-url' are not allowed")
		}
	}

	subcommand := userArgs[0]
	packageArgs := userArgs[1:]

	var args []string
	args = append(args, subcommand)

	if cfg.Common.BlockInteractive && subcommand == "install" {
		args = append(args, "--no-input")
	}

	if cfg.Managers.Pip.RegistryURL != "" && subcommand == "install" {
		args = append(args, "--index-url", cfg.Managers.Pip.RegistryURL, "--no-index")
	}

	args = append(args, packageArgs...)

	return args, nil
}

func ExecutePip(cfg *config.Config, args []string) error {
	binaryPath := "python"
	if cfg.Managers.Pip.BinaryPath != "" {
		binaryPath = cfg.Managers.Pip.BinaryPath
	}

	cmdArgs := append([]string{"-m", "pip"}, args...)

	cmd := exec.Command(binaryPath, cmdArgs...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
