package cmd

import (
	"errors"
	"fmt"
	"io/fs"
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
		fmt.Println("Generating wrapper scripts in output folder:", outputDir)

		if len(config.Tools) == 0 {
			fmt.Println("Warning: no tools are configured, nothing is generated.")
		}

		for toolName, toolConfig := range config.Tools {
			switch toolConfig.Type {
			case cfg.WrapperScript:
				fmt.Println("Generating script:", toolName)

				for i, mount := range toolConfig.Mounts {
					mount.Source = expandPath(mount.Source)
					toolConfig.Mounts[i] = mount

					_, err := os.Stat(mount.Source)
					if err == nil {
						continue
					}
					if !errors.Is(err, fs.ErrNotExist) {
						fmt.Fprintln(os.Stderr, "Error:", err)
						continue
					}
					fmt.Println("Warning: mount source path does not exist, creating folder:", mount.Source)
					if err := os.MkdirAll(mount.Source, 0755); err != nil {
						fmt.Fprintln(os.Stderr, "Error:", err)
					}
				}

				err := gen.Generate(outputDir, toolName, toolConfig)
				if err != nil {
					fmt.Fprintln(os.Stderr, "Failed:", err)
				}
			case cfg.Alias:
				fmt.Println("Generating alias:", toolName)
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
