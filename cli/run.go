package cli

import (
	"github.com/alecthomas/colour"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"path"
	"sopr/helpers"
	action "sopr/run"
	"sopr/types"
	"strings"
)

// runCmd represents the init command
var runCmd = &cobra.Command{
	Use:   "run [flags] [script name]",
	Short: "run user defined scripts",
	Args:  cobra.MinimumNArgs(0),
}

func cmdFactory(name string, desc string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   name,
		Short: desc,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			if helpers.IsTopLevelScript(name) {
				var scripts *[]types.Script
				if err := viper.UnmarshalKey("scripts", &scripts); err != nil {
					log.Print(err)
					log.Fatal("Error: could not parse scripts from config file")
				}

				p := path.Join(helpers.ProjectRoot(), viper.GetString("ProjectDirectory"))

				action.RunScript(name, &p, *scripts, func(_ types.Script) {})

			} else {
				projects := helpers.ProjectsForScript(name)

				if len(projects) <= 0 {
					colour.Print("^1Error^R - no project configured for script")
					return
				}

				selectedProjects := getSelectedRepos(AllRepos, projects)

				action.RunScriptForProjects(name, selectedProjects)
			}
		},
	}

	return cmd
}

func init() {
	viper.AddConfigPath(helpers.ProjectRoot())
	viper.SetConfigFile("sopr.yaml")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Can not find or parse config file")
	}
	descriptions := map[string]string{}

	if err := viper.UnmarshalKey("descriptions", &descriptions); err != nil {
		log.Print(err)
		log.Fatal("Error: could not parse descriptions from config file")
	}

	scripts := action.ListScripts()

	for _, script := range scripts {
		key := strings.ToLower(script)

		cmd := cmdFactory(script, descriptions[key])
		cmd.Flags().BoolVarP(&AllRepos, "all", "a", false, "use all repos")

		runCmd.AddCommand(cmd)
	}

	rootCmd.AddCommand(runCmd)

}
