package lib

import (
	"os"
	"path"
	"path/filepath"
)

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
