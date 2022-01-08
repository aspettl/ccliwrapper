package cmd

import (
	"fmt"
	"os"

	"github.com/aspettl/ccliwrapper/cfg"
	"github.com/aspettl/ccliwrapper/gen"
	"github.com/spf13/cobra"
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
		fmt.Fprintln(os.Stderr, "Generating wrapper scripts in output folder:", config.OutputDir)

		for toolName, toolConfig := range config.Tools {
			switch toolConfig.Type {
			case cfg.WrapperScript:
				fmt.Fprintln(os.Stderr, "Generating script:", toolName)
				err := gen.Generate(config.OutputDir, toolName, toolConfig)
				if err != nil {
					fmt.Fprintln(os.Stderr, "Failed:", err)
				}
			case cfg.Alias:
				fmt.Fprintln(os.Stderr, "Generating alias:", toolName)
				err := gen.GenerateAlias(config.OutputDir, toolConfig.AliasFor, toolName)
				if err != nil {
					fmt.Fprintln(os.Stderr, "Failed:", err)
				}
			}
		}
	},
}
