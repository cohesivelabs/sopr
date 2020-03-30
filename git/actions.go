package git

import (
	"context"
	"errors"
	"fmt"
	"sopr/helpers"
	"sopr/types"
	"sort"
	"time"

	"github.com/alecthomas/colour"
	gogit "gopkg.in/src-d/go-git.v4"
	gogitConfig "gopkg.in/src-d/go-git.v4/config"
)

type cloneResult struct {
	Error    error
	Path     string
	RepoName string
}

//Initialize runs clones all configured git repos.
func Initialize() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	projects := helpers.Projects()

	var resultChannels []<-chan cloneResult

	for _, project := range projects {
		resultChannels = append(resultChannels, clone(ctx, project))
	}

monitor:
	for _, ch := range resultChannels {
		select {
		case <-ctx.Done():
			break monitor
		case result := <-ch:
			if result.Error == nil {
				colour.Println(colour.Sprintf("^2%s^R cloned to ^2%s^R.", result.RepoName, result.Path))
			} else if result.Error.Error() != "repository already exists" {
				colour.Println(colour.Sprintf("^1ERROR^R - failed cloning ^2%s^R: %s", result.RepoName, result.Error.Error()))
			}
		}
	}

	select {
	case <-ctx.Done():
		colour.Println("^1ERROR^R - initialization timed out")
	default:
		return
	}
}

func clone(ctx context.Context, project types.ProjectConfig) <-chan cloneResult {
	result := make(chan cloneResult)
	repoName := project.Name

	finalizeResult := func(ctx context.Context, r cloneResult, ch chan<- cloneResult) {
		select {
		case <-ctx.Done():
			fmt.Println("cancelled.....")
			return
		case result <- r:
			return
		}
	}

	go func(result chan<- cloneResult) {
		select {
		case <-ctx.Done():
			return
		default:
			break
		}

		if project.Remotes == nil || len(project.Remotes) == 0 {
			r := cloneResult{
				RepoName: repoName,
				Error:    errors.New("No remotes configured"),
			}

			finalizeResult(ctx, r, result)
			return
		}

		// determine which remote should be the main remote
		// first attempt to see if an "origin" is provided
		// if no "origin" is found, use the first remote
		var mainRemote *types.Remote

		for _, remote := range project.Remotes {
			if remote.Name == "origin" {
				mainRemote = &remote
				break
			}
		}

		if mainRemote == nil {
			mainRemote = &project.Remotes[0]
		}

		ref, err := gogit.PlainCloneContext(ctx, *project.FullPath, false, &gogit.CloneOptions{
			URL:        mainRemote.URL,
			RemoteName: mainRemote.Name,
			Progress:   nil,
		})

		if err != nil {
			r := cloneResult{
				RepoName: repoName,
				Error:    err,
			}

			finalizeResult(ctx, r, result)
			return
		}

		for _, remote := range project.Remotes {
			if mainRemote.Name == remote.Name {
				continue
			}

			remoteConfig := &gogitConfig.RemoteConfig{
				Name: remote.Name,
				URLs: []string{remote.URL},
			}

			if _, err := ref.CreateRemote(remoteConfig); err != nil {
				colour.Println(colour.Sprintf("^3WARNING^R: Could not configure remote ^2%s^R for repo ^2%s^R.", remote.URL, repoName))
			}
		}

		r := cloneResult{
			RepoName: repoName,
			Path:     *project.FullPath,
		}

		finalizeResult(ctx, r, result)
		return
	}(result)

	return result
}

// CheckoutBranch used to create or switch branches in a repo
func CheckoutBranch(selectedProjects []types.ProjectConfig, branchName string, allRepos bool, create bool) {
	pristineRepos := getPristineRepos(selectedProjects)

	for _, repo := range pristineRepos {
		colour.Println(colour.Sprintf("Checking out Branch ^4%s^R in: ^2%s^R.", branchName, repo.Name))

		err := checkout(repo, branchName, create)
		if err != nil {
			fmt.Println(fmt.Sprintf("Error: %s - %s", repo.Name, err))
		}
	}
}

//ListRepos lists the available repos
func ListRepos() {
	var output []map[string]string

	projects := helpers.Projects()

	for _, project := range projects {
		b, err := branch(project)
		if err != nil {
			colour.Println(colour.Sprintf(colour.Sprintf("^1Error^R: could not get branch for ^2%s^R", project.Name)))
			continue
		}

		output = append(output, map[string]string{
			"Name":   project.Name,
			"Branch": b,
		})
	}

	sort.Slice(output, func(i, j int) bool { return output[i]["Name"] < output[j]["Name"] })

	for _, repo := range output {
		colour.Println(colour.Sprintf("^2%s^R (^4%s^R)", repo["Name"], repo["Branch"]))
	}
}

func getPristineRepos(selectedProjects []types.ProjectConfig) []types.ProjectConfig {
	var pristineRepos []types.ProjectConfig

	fmt.Println("Checking local working tree status")
	for _, project := range selectedProjects {
		if clean, _ := isClean(project); !clean {
			colour.Println(colour.Sprintf("Working tree for ^2%s^R is not clean, skipping.", project.Name))
			continue
		}

		pristineRepos = append(pristineRepos, project)
	}

	return pristineRepos
}

//Update updates a repo
func Update(selectedProjects []types.ProjectConfig) {
	pristineRepos := getPristineRepos(selectedProjects)

	for _, repo := range pristineRepos {
		err := pull(repo)
		if err != nil && err.Error() == "already up-to-date" {
			colour.Println(colour.Sprintf("Skipping ^2%s^R because its already up to date.", repo.Name))
			continue
		} else if err != nil {
			fmt.Println(fmt.Sprintf("Error: %s - %s", repo.Name, err))
			continue
		}

		colour.Println(colour.Sprintf("Updating ^2%s^R.", repo.Name))
		fmt.Println(fmt.Sprintf("Updating %s", repo.Name))
	}
}
