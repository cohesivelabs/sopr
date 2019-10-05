package actions

import (
	"bytes"
	"context"
	"fmt"
	"github.com/alecthomas/colour"
	"io"
	"log"
	"os"
	"os/exec"
	"sopr/lib/git"
	"sopr/lib/prompts"
	"strings"
	"time"
)

func execute(rawCommand string) {
	var stdoutBuffer, stderrBuffer bytes.Buffer

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	command := strings.Split(rawCommand, " ")
	cmd := exec.CommandContext(ctx, command[0], command[1:]...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	cmd.Start()

	go writeToConsole(ctx, &stdoutBuffer, os.Stdout, stdout)
	go writeToConsole(ctx, &stderrBuffer, os.Stderr, stderr)

	cmd.Wait()
}

func writeToConsole(ctx context.Context, buffer *bytes.Buffer, target io.Writer, pipe io.ReadCloser) {
	writer := io.MultiWriter(target, buffer)
	io.Copy(writer, pipe)
	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}

func DepsInstall(allRepos, dryrun bool) {
	var selectedRepos []git.Repo

	repos, err := git.RepoList(true)
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
		if repo.Config.InstallDeps != nil {
			colour.Printf("^running command (^2%s^R) for ^2%s^R. \n", *repo.Config.InstallDeps, repo.Config.Name)

			if dryrun == false {
				execute(*repo.Config.InstallDeps)
			}
		}
	}
}

func DepsRemove(allRepos, dryrun bool) {
	var selectedRepos []git.Repo

	repos, err := git.RepoList(true)
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
		if repo.Config.InstallDeps != nil {
			colour.Printf("^running command (^2%s^R) for ^2%s^R. \n", repo.Config.RemoveDeps, repo.Config.Name)

			if dryrun == false {
				execute(*repo.Config.RemoveDeps)
			}
		}
	}
}
