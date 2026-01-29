package cmd

import (
	"fmt"
	"os"

	"safe-install/internal/config"
	"safe-install/internal/policy"

	"github.com/spf13/cobra"
)

var pipCmd = &cobra.Command{
	Use:   "pip",
	Short: "Wraps Python pip with security policies",
	Long: `Wraps Python pip with security policies, enforcing virtual environment 
requirements, registry URLs, and blocking interactive prompts.`,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}

		if !cfg.Managers.Pip.Enabled {
			fmt.Fprintf(os.Stderr, "Pip is disabled in configuration\n")
			os.Exit(1)
		}

		if err := policy.ValidatePipEnv(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Validation error: %v\n", err)
			os.Exit(1)
		}

		buildArgs, err := policy.BuildPipArgs(cfg, args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error building pip arguments: %v\n", err)
			os.Exit(1)
		}

		if err := policy.ExecutePip(cfg, buildArgs); err != nil {
			fmt.Fprintf(os.Stderr, "Error executing pip: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(pipCmd)
}
