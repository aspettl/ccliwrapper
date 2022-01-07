package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	generateCmd.Flags().StringP("output-dir", "o", "", "Output directory for wrapper scripts")

	rootCmd.AddCommand(generateCmd)
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate wrapper scripts",
	Long:  `Generate wrapper scripts for all configured CLI tools.`,
	Run: func(cmd *cobra.Command, args []string) {
		outputDir := viper.GetString("output-dir")
		fmt.Println("TODO, output dir: ", outputDir)
	},
}
