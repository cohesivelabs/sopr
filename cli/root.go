package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"sopr/helpers"
)

var cfgFile string

// AllRepos flag is used to run a action against all available repos
var AllRepos bool

var rootCmd = &cobra.Command{
	Use: "sopr",
}

// Execute is the main entrypoint for cobra
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.AddConfigPath(helpers.ProjectRoot())
	viper.SetConfigFile("sopr.yaml")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Can not find or parse config file")
	}
}
