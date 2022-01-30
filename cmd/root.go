package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/aspettl/ccliwrapper/cfg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const cfgFileDefault = ".ccliwrapper.yaml"
const outputDirDefault = "~/.local/bin"

var homeDir string
var cfgFile string
var config cfg.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ccliwrapper",
	Short: "Simplify use of containerized CLI tools",
	Long: `ccliwrapper produces wrapper scripts for CLI tools running in containers
via Docker or Podman. This helps to avoid local installation of many tools
and limits their access to the host system.`,
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
	// Find home directory.
	var err error
	homeDir, err = os.UserHomeDir()
	cobra.CheckErr(err)

	// Use Podman by default (if installed, otherwise fall back to Docker).
	engineDefault := "docker"
	if fileExists("/usr/bin/podman") {
		engineDefault = "podman"
	}

	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("config file (default is ~/%v)", cfgFileDefault))

	viper.SetDefault("OutputDir", outputDirDefault)
	viper.SetDefault("Engine", engineDefault)
	viper.SetDefault("Tools", map[string]cfg.ToolConfig{})
}

// initConfig sets up the config file in viper, but does not actually try to read.
func initConfig() {
	// Use yaml as file format when file name does not have a file extension.
	viper.SetConfigType("yaml")

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".ccliwrapper.yaml".
		viper.AddConfigPath(homeDir)
		viper.SetConfigName(cfgFileDefault)
	}
}

// readConfig reads the config file, this method needs to be explicitly called by all the commands needing the config.
func readConfig() error {
	// Read config file.
	err := viper.ReadInConfig()
	if viper.ConfigFileUsed() != "" {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
	if err != nil {
		return err
	}

	// Parse the data into our config struct.
	if err := viper.Unmarshal(&config); err != nil {
		return err
	}

	// Apply some default values for configured tools.
	config.ApplyToolDefaults()

	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}
