package helpers

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"sopr/types"

	"github.com/spf13/viper"
)

//ProjectRoot returns the root
func ProjectRoot() string {
	searchPath, _ := os.Getwd()

	for {
		newSearchPath := filepath.Dir(searchPath)

		if newSearchPath == searchPath {
			break
		}

		if _, err := os.Stat(path.Join(searchPath, "sopr.yaml")); !os.IsNotExist(err) {
			break
		}

		searchPath = newSearchPath
	}

	return searchPath
}

//IsTopLevelScript used to determine if the provided script is global
func IsTopLevelScript(scriptName string) bool {
	var scripts *[]types.Script
	exists := false

	if err := viper.UnmarshalKey("scripts", &scripts); err != nil {
		log.Print(err)
		log.Fatal("Error: could not parse scripts from config file")
	}

	if scripts != nil {
		for _, script := range *scripts {
			if script.Name == scriptName {
				exists = true
				continue
			}
		}
	}

	return exists
}

//ProjectsForScript return all projects that have the required script configured.
func ProjectsForScript(scriptName string) []types.ProjectConfig {
	projects := Projects()

	projectScripts := make([]types.ProjectConfig, 0)

	for _, project := range projects {
		if project.Scripts == nil {
			continue
		}

		for _, script := range *project.Scripts {
			if script.Name == scriptName {
				projectScripts = append(projectScripts, project)
				break
			}
		}
	}

	return projectScripts
}

//ProjectsWithRepos returns a list of all configured projects that have git repositories configured
func ProjectsWithRepos() []types.ProjectConfig {
	projects := Projects()

	projectRepos := make([]types.ProjectConfig, 0)

	for _, project := range projects {
		if len(project.Remotes) <= 0 {
			continue
		}

		projectRepos = append(projectRepos, project)
	}

	return projectRepos
}
