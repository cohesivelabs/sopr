package cli

import (
	"github.com/spf13/cobra"
	"sopr/git"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use: "init",
	Long: `Initializes the development environment
by cloning all of the repositories listed in repositories.yaml and placing them in the repo/ directory`,
	Run: func(cmd *cobra.Command, args []string) {
		git.Initialize()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
