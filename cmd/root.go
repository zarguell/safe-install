package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "safe-install",
	Short: "A unified, cross-platform CLI wrapper that enforces security policies",
	Long: `Safe-Install is a unified, cross-platform CLI wrapper that enforces security 
policies (like internal registries and lockfiles) across multiple package 
managers (pip, npm, etc.).`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
