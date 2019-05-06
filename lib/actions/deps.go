package actions

import (
    "fmt"
	"io"
	"github.com/alecthomas/colour"
	"log"
	"os"
	"os/exec"
	"sopr/lib/git"
	"sopr/lib/prompts"
	"strings"
	"sync"
	"bytes"
)

func execute(rawCommand string) {
	var stdoutBuffer bytes.Buffer
	var stderrBuffer bytes.Buffer
	var waitGroup sync.WaitGroup

	command := strings.Split(rawCommand, " ")
	cmd := exec.Command(command[0], command[1:]...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	stdoutWriter := io.MultiWriter(os.Stdout, &stdoutBuffer)
	stderrWriter := io.MultiWriter(os.Stderr, &stderrBuffer)

	cmd.Start()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		io.Copy(stdoutWriter, stdout)
	}()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		io.Copy(stderrWriter, stderr)
	}()

	cmd.Wait()
	waitGroup.Wait()
}

func DepsInstall(allRepos, dryrun bool) {
	var selectedRepos []git.Repo

	repos, err := git.RepoList(false)
	if err != nil {
		fmt.Printf("Error getting repo list %s", err)
		os.Exit(1)
	}

	if allRepos {
		selectedRepos = repos
	} else {
		selectedRepos = prompts.RepoSelectPrompt(repos)
	}

	if err != nil {
		log.Fatalf("Error: could not read repository list - %s", err)
	}

	for _, repo := range selectedRepos {
		if repo.Config.InstallDeps != "" {
			colour.Printf("^running command (^2%s^R) for ^2%s^R. \n", repo.Config.InstallDeps, repo.Config.Name)

			if dryrun == false {
				execute(repo.Config.InstallDeps)
			}
		}
	}
}

func DepsRemove(allRepos, dryrun bool) {
	var selectedRepos []git.Repo

	repos, err := git.RepoList(false)
	if err != nil {
		fmt.Printf("Error getting repo list %s", err)
		os.Exit(1)
	}

	if allRepos {
		selectedRepos = repos
	} else {
		selectedRepos = prompts.RepoSelectPrompt(repos)
	}

	if err != nil {
		log.Fatalf("Error: could not read repository list - %s", err)
	}

	for _, repo := range selectedRepos {
		if repo.Config.InstallDeps != "" {
			colour.Printf("^running command (^2%s^R) for ^2%s^R. \n", repo.Config.RemoveDeps, repo.Config.Name)

			if dryrun == false {
				execute(repo.Config.RemoveDeps)
			}
		}
	}
}
