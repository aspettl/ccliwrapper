package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(writeConfigCmd)
}

var writeConfigCmd = &cobra.Command{
	Use:   "write-config",
	Short: "Write config file",
	Long: `Writes the current configuration to the specified config file or
the default config file if nothing else is specified. Useful to generate the
initial template or to add all keys and their defaults to the current config.`,
	Run: func(cmd *cobra.Command, args []string) {
		fileName := viper.ConfigFileUsed()
		if fileName == "" {
			home, err := os.UserHomeDir()
			cobra.CheckErr(err)
			fileName = path.Join(home, cfgFileDefault)
		}
		err := viper.WriteConfigAs(fileName)
		cobra.CheckErr(err)
		fmt.Fprintln(os.Stderr, "Written to config file:", fileName)
	},
}
