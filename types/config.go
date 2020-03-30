package types

//A ScriptOptions struct
type ScriptOptions struct {
	Path *string
}

//A Script struct
type Script struct {
	Name    string
	Command string
	Options *ScriptOptions
}

//A Remote struct
type Remote struct {
	Name string
	URL  string
}

//A ProjectConfig struct
type ProjectConfig struct {
	Name     string
	Path     *string
	Remotes  []Remote
	Scripts  *[]Script
	FullPath *string
}

//A Config struct
type Config struct {
	ProjectDirectory string `yaml:"projectDirectory"`
	Projects         []ProjectConfig
	Scripts          *[]Script
	Descriptions     *map[string]string
}
