package scripts

import (
	"context"
	"fmt"
	"github.com/alecthomas/colour"
	"sopr/types"
	"time"
)

// ListScripts lists all possible scripts to execute
func ListScripts() []string {
	projects, topLevelScripts := scriptsFromConfig()

	return buildScriptList(projects, topLevelScripts)
}

//RunScript executes a given script
func RunScript(name string, path *string, scripts []types.Script, preRun func(types.Script)) error {
	var err error

	for _, script := range scripts {
		if script.Name == name {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
			defer cancel()

			preRun(script)

			doneCh := make(chan bool)
			errCh := make(chan error)
			execute(ctx, script.Command, path, doneCh, errCh)

		monitor:
			// block until command times out or finishes
			for {
				select {
				case <-ctx.Done():
					err = ctx.Err()
					break monitor
				case err = <-errCh:
					break monitor
				case <-doneCh:
					break monitor
				}
			}

			break
		}
	}

	return err
}

//RunScriptForProjects executes script with given name
func RunScriptForProjects(name string, selectedProjects []types.ProjectConfig) {
	for _, project := range selectedProjects {
		if project.Scripts == nil {
			continue
		}

		printScriptDetails := func(script types.Script) {
			fmt.Println("******************")
			fmt.Println(fmt.Sprintf("Executing %s for %s", script.Name, project.Name))
			fmt.Println(script.Command)
			fmt.Println("******************")
		}

		if err := RunScript(name, project.FullPath, *project.Scripts, printScriptDetails); err != nil {
			colour.Println("^1Error^R - " + err.Error())
			break
		}
	}

	fmt.Println("Script Completed")
}
