package config

type Remote struct {
	Name string
	Url  string
}

type RepoConfig struct {
	Name        string
	Path        string
	Remotes     []Remote
	InstallDeps string `yaml:"installDeps"`
	RemoveDeps  string `yaml:"removeDeps"`
}

type Config struct {
	RepoDirectory string `yaml:"repoDirectory"`
	Repos         []RepoConfig
}
