package cmd

import (
	"fmt"
	"os"

	"github.com/aspettl/ccliwrapper/gen"
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
		fmt.Fprintln(os.Stderr, "Generating wrapper scripts in output folder:", outputDir)

		for toolName := range viper.GetStringMap("tools") {
			props := viper.Sub("tools." + toolName)

			if aliasFor := props.GetString("alias_for"); aliasFor != "" {
				fmt.Fprintln(os.Stderr, "Generating alias:", toolName)
				err := gen.GenerateAlias(outputDir, aliasFor, toolName)
				if err != nil {
					fmt.Fprintln(os.Stderr, "Failed:", err)
				}
				continue
			}

			toolConfig := gen.ToolConfig{
				ImageName: props.GetString("image_name"),
			}
			fmt.Fprintln(os.Stderr, "Generating script:", toolName)
			err := gen.Generate(outputDir, toolName, toolConfig)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Failed:", err)
			}
		}
	},
}
