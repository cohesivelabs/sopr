package git

import (
	"fmt"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"log"
	"os"
	"path"
	"sopr/lib"
	"sopr/lib/config"
)

func RepoList(returnMissing bool) ([]Repo, error) {
	var repos []Repo

	basePath := lib.ProjectRoot()

	config, err := config.ParseConfig()

	if err != nil {
		log.Fatalf("error: unable to parse config %v", err)
	}

	for _, repoConfig := range config.Repos {
		repo := Repo{}

		repo.Config = repoConfig
		repo.FullPath = path.Join(basePath, config.RepoDirectory, repoConfig.Path)

		if ref, err := git.PlainOpen(repo.FullPath); err == nil {
			repo.Ref = ref
		} else if !returnMissing {
			continue
		}

		repos = append(repos, repo)
	}

	return repos, err
}

func (repo Repo) Branch() (string, error) {
	head, err := repo.Ref.Head()
	if err != nil {
		return "", err
	}

	return head.Name().Short(), nil
}

func (repo Repo) IsClean() bool {
	tree, err := repo.Ref.Worktree()
	if err != nil {
		fmt.Printf("Error: getting working tree for %s - %s", repo.Config.Name, err)
		os.Exit(1)
	}

	status, err := tree.Status()
	if err != nil {
		fmt.Printf("Error: getting working tree statue for %s - %s", repo.Config.Name, err)
		os.Exit(1)
	}

	return status.IsClean()
}

func (repo Repo) Pull() error {
	tree, err := repo.Ref.Worktree()
	if err != nil {
		fmt.Printf("Error: getting working tree for %s - %s", repo.Config.Name, err)
		os.Exit(1)
	}

	err = tree.Pull(&git.PullOptions{
		RemoteName: "origin",
		Progress:   os.Stdout,
	})

	return err
}

func (repo Repo) Checkout(branchName string, create bool) error {
	var branchRef *plumbing.Reference

	if create == false {
		branch, err := repo.Ref.Branch(branchName)
		if err != nil {
			fmt.Printf("Error: resolving branch %s in %s - %s", branchName, repo.Config.Name, err)
			os.Exit(1)
		}

		branchRef, err = repo.Ref.Reference(branch.Merge, false)
		if err != nil {
			fmt.Printf("Error: resolving branch %s in %s - %s", branchName, repo.Config.Name, err)
			os.Exit(1)
		}
	} else {
		head, err := repo.Ref.Head()
		if err != nil {
			fmt.Printf("Error: getting current HEAD %s - %s", repo.Config.Name, err)
			os.Exit(1)
		}

		branchRef = plumbing.NewReferenceFromStrings(fmt.Sprintf("refs/heads/%s", branchName), head.Hash().String())
	}

	tree, err := repo.Ref.Worktree()
	if err != nil {
		fmt.Printf("Error: getting working tree for %s - %s", repo.Config.Name, err)
		os.Exit(1)
	}

	err = tree.Checkout(&git.CheckoutOptions{
		Branch: branchRef.Name(),
		Create: create,
	})

	return err
}
