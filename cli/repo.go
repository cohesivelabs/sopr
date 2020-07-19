package cli

import (
	"github.com/alecthomas/colour"
	"github.com/spf13/cobra"
	"sopr/git"
	"sopr/helpers"
)

// repoCmd represents the repo command
var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Commands for interactive with repositories",
	Args:  cobra.MinimumNArgs(1),
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Perform git pull",
	Args:  cobra.OnlyValidArgs,
	Run: func(cmd *cobra.Command, args []string) {
		repos := helpers.ProjectsWithRepos()

		if len(repos) <= 0 {
			colour.Print("^1Error^R - No projects configured with git remotes")
			return
		}

		git.Update(getSelectedRepos(AllRepos, repos))
	},
}

var branchCmd = &cobra.Command{
	Use:   "branch",
	Short: "Branch operations",
	Args:  cobra.MinimumNArgs(0),
}

var createBranchCmd = &cobra.Command{
	Use:   "create [branchName]",
	Short: "Create git branch",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		repos := helpers.ProjectsWithRepos()

		if len(repos) <= 0 {
			colour.Print("^1Error^R - No projects configured with git remotes")
			return
		}

		branchName := getBranchName(args)
		selectedProjects := getSelectedRepos(AllRepos, repos)

		git.CheckoutBranch(selectedProjects, branchName, AllRepos, true)
	},
}

var listRepoCmd = &cobra.Command{
	Use:   "list",
	Short: "List repos and their branches",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		repos := helpers.ProjectsWithRepos()

		if len(repos) <= 0 {
			colour.Print("^1Error^R - No projects configured with git remotes")
			return
		}

		git.ListRepos()
	},
}

var switchBranchCmd = &cobra.Command{
	Use:   "switch [branchName]",
	Short: "Change git branch",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		repos := helpers.ProjectsWithRepos()

		if len(repos) <= 0 {
			colour.Print("^1Error^R - No projects configured with git remotes")
			return
		}

		branchName := getBranchName(args)

		selectedProjects := getSelectedRepos(AllRepos, repos)

		git.CheckoutBranch(selectedProjects, branchName, AllRepos, false)
	},
}

func init() {
	updateCmd.Flags().BoolVarP(&AllRepos, "all", "a", false, "use all repos")
	switchBranchCmd.Flags().BoolVarP(&AllRepos, "all", "a", false, "use all repos")
	createBranchCmd.Flags().BoolVarP(&AllRepos, "all", "a", false, "use all repos")

	rootCmd.AddCommand(repoCmd)
	repoCmd.AddCommand(updateCmd)
	repoCmd.AddCommand(listRepoCmd)
	repoCmd.AddCommand(branchCmd)
	branchCmd.AddCommand(createBranchCmd)
	branchCmd.AddCommand(switchBranchCmd)
}
