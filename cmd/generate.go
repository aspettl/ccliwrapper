package cmd

import (
	"fmt"
	"os"
	"strings"

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
		outputDir := expandPath(config.OutputDir)
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			cobra.CheckErr(err)
		}
		fmt.Fprintln(os.Stderr, "Generating wrapper scripts in output folder:", outputDir)

		for toolName, toolConfig := range config.Tools {
			switch toolConfig.Type {
			case cfg.WrapperScript:
				fmt.Fprintln(os.Stderr, "Generating script:", toolName)
				err := gen.Generate(outputDir, toolName, toolConfig)
				if err != nil {
					fmt.Fprintln(os.Stderr, "Failed:", err)
				}
			case cfg.Alias:
				fmt.Fprintln(os.Stderr, "Generating alias:", toolName)
				err := gen.GenerateAlias(outputDir, toolConfig.AliasFor, toolName)
				if err != nil {
					fmt.Fprintln(os.Stderr, "Failed:", err)
				}
			}
		}
	},
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") || strings.HasPrefix(path, "~\\") {
		path = homeDir + path[1:]
	}
	return os.ExpandEnv(path)
}
