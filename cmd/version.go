package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

// Version is set via compiler flags, see .goreleaser.yaml
var Version = "0.0.0-dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version string",
	Long:  `Print the exact version number of ccliwrapper.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ccliwrapper", Version)
	},
}
