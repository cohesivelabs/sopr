package prompts

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"os"
	"sopr/lib/git"
	"sort"
	"strings"
)

func BranchNamePrompt() string {
	prompt := promptui.Prompt{
		Label: "Branch Name",
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	return result
}

func RepoSelectPrompt(repoList []git.Repo) []git.Repo {
	var options []string
	var results []git.Repo

	for _, repo := range repoList {
		branch, _ := repo.Branch()
		options = append(options, fmt.Sprintf("%s (%s)", repo.Config.Name, branch))
	}

	sort.Strings(options)
	options = append(options, "Done")

	index := -1
	for index < len(options)-1 {
		var err error
		var result string

		if index >= 0 {
			options = append(options[:index], options[index+1:]...)
		}

		repoSelectPrompt := promptui.Select{
			Label: "Select Repositories",
			Items: options,
		}

		index, result, err = repoSelectPrompt.Run()
		if err != nil {
			fmt.Printf("Error: failed to open repo select list: %s", err)
			os.Exit(1)
		}

		if index < len(options)-1 {
			name := strings.Split(result, " ")[0]
			for _, repo := range repoList {
				if repo.Config.Name == name {
					results = append(results, repo)
				}
			}
		}

	}

	return results
}
