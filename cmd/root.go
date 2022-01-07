package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ccliwrapper",
	Short: "Simplify use of containerized CLI tools",
	Long: `ccliwrapper produces wrapper scripts for CLI tools running in containers
via Docker or Podman. This helps to avoid local installation of many tools
and limits their access to the host system.`,
}

const cfgFileDefault = ".ccliwrapper.yaml"
const outputDirDefault = "~/.local/bin"

// ToolConfig represents the configuration for a CLI tool in the config file
type ToolConfig struct {
	name string
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("config file (default is $HOME/%v)", cfgFileDefault))

	viper.SetDefault("output-dir", outputDirDefault)
	viper.SetDefault("tools", []ToolConfig{})
	viper.BindPFlag("output-dir", generateCmd.Flags().Lookup("output-dir"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Use yaml as file format when file name does not have a file extension
	viper.SetConfigType("yaml")

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".ccliwrapper.yaml".
		viper.AddConfigPath(home)
		viper.SetConfigName(cfgFileDefault)
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
