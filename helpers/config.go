package helpers

import (
	"github.com/spf13/viper"
	"log"
	"path"
	"sopr/types"
)

// Projects gives a list of all projects from the current config
func Projects() []types.ProjectConfig {
	var projects []types.ProjectConfig
	var projectsWithPath []types.ProjectConfig
	rootPath := ProjectRoot()
	projectsDir := viper.GetString("projectDirectory")

	if err := viper.UnmarshalKey("projects", &projects); err != nil {
		log.Print(err)
		log.Fatal("Error: could not parse projects from config file")
	}

	// extrapolate the full repo path
	for _, project := range projects {
		var dir string

		// default to project name if path fragement is not configured
		if project.Path != nil {
			dir = *project.Path
		} else {
			dir = project.Name
		}

		fullPath := path.Join(rootPath, projectsDir, dir)
		project.FullPath = &fullPath
		projectsWithPath = append(projectsWithPath, project)
	}

	return projectsWithPath
}
