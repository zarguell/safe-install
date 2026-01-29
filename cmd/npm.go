package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var npmCmd = &cobra.Command{
	Use:   "npm",
	Short: "Wraps Node npm with security policies",
	Long: `Wraps Node npm with security policies, enforcing registry URLs and blocking 
interactive prompts. (Coming Soon)`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("NPM support coming soon")
	},
}

func init() {
	rootCmd.AddCommand(npmCmd)
}
