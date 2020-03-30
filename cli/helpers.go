package cli

import (
	"sopr/cli/prompts"
	"sopr/types"
)

func getSelectedRepos(allProjects bool, projects []types.ProjectConfig) []types.ProjectConfig {
	var selectedProjects []types.ProjectConfig

	if AllRepos {
		selectedProjects = projects
	} else {
		selectedProjects = prompts.RepoSelectPrompt(projects)
	}

	return selectedProjects
}

func getBranchName(args []string) string {
	branchName := ""

	if len(args) > 0 {
		branchName = args[0]
	}

	if branchName == "" {
		branchName = prompts.BranchNamePrompt()
	}

	return branchName
}
