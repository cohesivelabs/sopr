package cmd

import (
	"github.com/spf13/cobra"
	"sopr/lib/actions"
)

var AllReposDeps bool
var DryRun bool

// depsCmd represents the init command
var depsCmd = &cobra.Command{
	Use:   "deps",
	Short: "install dependencies",
	Args:  cobra.MinimumNArgs(1),
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "install dependencies for repo",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		actions.DepsInstall(AllReposDeps, DryRun)
	},
}

var removeCmd = &cobra.Command{
	Use:   "rm",
	Short: "remove dependencies for repo",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		actions.DepsRemove(AllReposDeps, DryRun)
	},
}

func init() {
	installCmd.Flags().BoolVarP(&AllReposDeps, "all", "a", false, "use all repos")
	removeCmd.Flags().BoolVarP(&AllReposDeps, "all", "a", false, "use all repos")
	installCmd.Flags().BoolVarP(&DryRun, "dryrun", "d", false, "show what command executes but don't execute them")
	removeCmd.Flags().BoolVarP(&DryRun, "dryrun", "d", false, "show what command executes but don't execute them")
	rootCmd.AddCommand(depsCmd)
	depsCmd.AddCommand(installCmd)
	depsCmd.AddCommand(removeCmd)
}
