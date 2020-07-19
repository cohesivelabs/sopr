package git

import (
	"errors"
	"fmt"
	"os"
	"path"
	"sopr/helpers"
	"sopr/types"

	"github.com/spf13/viper"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func refFromProject(project types.ProjectConfig) (*git.Repository, error) {
	var fullPath string
	rootPath := helpers.ProjectRoot()
	projectsDir := viper.GetString("projectDirectory")

	if project.Path != nil && len(project.Remotes) > 0 {
		fullPath = path.Join(rootPath, projectsDir, *project.Path)
	} else {
		return nil, errors.New("No remotes or repo path configured for " + project.Name)
	}

	ref, err := git.PlainOpen(fullPath)
	if err != nil {
		return nil, err
	}

	return ref, nil
}

func branch(project types.ProjectConfig) (string, error) {
	ref, err := refFromProject(project)
	if err != nil {
		return "", err
	}

	head, err := ref.Head()
	if err != nil {
		return "", err
	}

	return head.Name().Short(), nil
}

func isClean(project types.ProjectConfig) (bool, error) {
	ref, err := refFromProject(project)
	if err != nil {
		return false, err
	}

	tree, err := ref.Worktree()
	if err != nil {
		return false, err
	}

	status, err := tree.Status()
	if err != nil {
		return false, err
	}

	return status.IsClean(), nil
}

func pull(project types.ProjectConfig) error {
	ref, err := refFromProject(project)
	if err != nil {
		return err
	}

	tree, err := ref.Worktree()
	if err != nil {
		return err
	}

	err = tree.Pull(&git.PullOptions{
		RemoteName: "origin",
		Progress:   os.Stdout,
	})

	return err
}

func checkout(project types.ProjectConfig, branchName string, create bool) error {
	var branchRef *plumbing.Reference

	ref, err := refFromProject(project)
	if err != nil {
		return err
	}

	if create == false {
		branch, err := ref.Branch(branchName)
		if err != nil {
			return err
		}

		branchRef, err = ref.Reference(branch.Merge, false)
		if err != nil {
			return err
		}
	} else {
		head, err := ref.Head()
		if err != nil {
			return err
		}

		branchRef = plumbing.NewReferenceFromStrings(fmt.Sprintf("refs/heads/%s", branchName), head.Hash().String())
	}

	tree, err := ref.Worktree()
	if err != nil {
		return err
	}

	err = tree.Checkout(&git.CheckoutOptions{
		Branch: branchRef.Name(),
		Create: create,
	})

	return err
}
