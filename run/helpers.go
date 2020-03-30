package scripts

import (
	"bytes"
	"context"
	"github.com/spf13/viper"
	"io"
	"log"
	"os"
	"os/exec"
	"sopr/types"
	"sort"
	"strings"
)

func writeToConsole(buffer *bytes.Buffer, target io.Writer, pipe io.ReadCloser, errCh chan<- error) {
	writer := io.MultiWriter(target, buffer)

	_, err := io.Copy(writer, pipe)
	if err != nil {
		errCh <- err
		return
	}
}

func scriptsFromConfig() ([]types.ProjectConfig, *[]types.Script) {
	var projects []types.ProjectConfig
	var topLevelScripts *[]types.Script

	if err := viper.UnmarshalKey("Projects", &projects); err != nil {
		log.Print(err)
		log.Fatal("Error: could not parse projects from config file")
	}

	if err := viper.UnmarshalKey("scripts", &topLevelScripts); err != nil {
		log.Print(err)
		log.Fatal("Error: could not parse scripts from config file")
	}

	return projects, topLevelScripts
}

func buildScriptList(repos []types.ProjectConfig, topLevelScripts *[]types.Script) []string {
	scripts := make([]string, 0)
	scriptsMap := make(map[string]bool)

	if topLevelScripts != nil {
		for _, script := range *topLevelScripts {
			if scriptsMap[script.Name] {
				continue
			}

			scripts = append(scripts, script.Name)
			scriptsMap[script.Name] = true
		}
	}

	for _, repo := range repos {
		if repo.Scripts == nil {
			continue
		}

		for _, script := range *repo.Scripts {
			if scriptsMap[script.Name] {
				continue
			}

			scripts = append(scripts, script.Name)
			scriptsMap[script.Name] = true
		}
	}

	sort.Strings(scripts)

	return scripts
}

// Execute shell command
func execute(ctx context.Context, rawCommand string, path *string, doneCh chan<- bool, errCh chan<- error) {
	go func() {
		var stdoutBuffer, stderrBuffer bytes.Buffer

		command := strings.Split(rawCommand, " ")
		cmd := exec.CommandContext(ctx, command[0], command[1:]...)

		if path != nil {
			cmd.Dir = *path
		}

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			errCh <- err
			return
		}

		defer stdout.Close()

		stderr, err := cmd.StderrPipe()
		if err != nil {
			errCh <- err
			return
		}

		defer stderr.Close()

		err = cmd.Start()
		if err != nil {
			errCh <- err
			return
		}

		go writeToConsole(&stdoutBuffer, os.Stdout, stdout, errCh)
		go writeToConsole(&stderrBuffer, os.Stderr, stderr, errCh)

		err = cmd.Wait()
		if err != nil {
			errCh <- err
			return
		}

		doneCh <- true
	}()
}
