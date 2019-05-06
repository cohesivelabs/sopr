package config

import (
    "github.com/go-yaml/yaml"
    "sopr/lib"
	"io/ioutil"
	"log"
	"path"
)

func ParseConfig() (Config, error) {
	basePath := lib.ProjectRoot()

	data, err := ioutil.ReadFile(path.Join(basePath, "sopr.yaml"))
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	config := Config{}

	err = yaml.Unmarshal([]byte(data), &config)

	return config, err
}
