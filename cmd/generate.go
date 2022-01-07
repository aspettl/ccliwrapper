package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(generateCmd)
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate wrapper scripts",
	Long:  `Generate wrapper scripts for all configured CLI tools.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO")
	},
}
