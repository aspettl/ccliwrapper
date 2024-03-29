package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"sort"
	"strings"

	"github.com/aspettl/ccliwrapper/cfg"
	"github.com/aspettl/ccliwrapper/gen"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	generateCmd.Flags().StringP("output-dir", "o", "", "Output directory for wrapper scripts")
	generateCmd.Flags().StringP("template-file", "t", "", "Path to custom wrapper script template file")

	viper.BindPFlag("OutputDir", generateCmd.Flags().Lookup("output-dir"))
	viper.BindPFlag("TemplateFile", generateCmd.Flags().Lookup("template-file"))

	rootCmd.AddCommand(generateCmd)
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate wrapper scripts",
	Long:  `Generate wrapper scripts for all configured CLI tools.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := readConfig(); err != nil {
			cobra.CheckErr(err)
		}

		outputDir := expandPath(config.OutputDir)
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			cobra.CheckErr(err)
		}
		fmt.Println("Generating wrapper scripts in output folder:", outputDir)

		templateFile := expandPath(config.TemplateFile)
		if templateFile != "" {
			fmt.Println("Using custom wrapper script template:", templateFile)
		}

		if len(config.Tools) == 0 {
			fmt.Fprintln(os.Stderr, "Warning: no tools are configured, nothing is generated.")
		}

		errorCount := 0
		for _, toolName := range sortedToolNames() {
			toolConfig := config.Tools[toolName]
			switch toolConfig.Type {
			case cfg.WrapperScript:
				fmt.Println("Generating script:", toolName)
				err := generateWrapperScript(outputDir, templateFile, toolName, toolConfig)
				if err != nil {
					fmt.Fprintln(os.Stderr, "Error:", err)
					errorCount++
				}
			case cfg.Alias:
				fmt.Println("Generating alias:", toolName)
				var err error
				if toolConfig.ForceTemplate {
					err = generateWrapperScript(outputDir, templateFile, toolName, config.Tools[toolConfig.AliasFor])
				} else {
					err = gen.GenerateAlias(outputDir, toolConfig.AliasFor, toolName)
				}
				if err != nil {
					fmt.Fprintln(os.Stderr, "Error:", err)
					errorCount++
				}
			}
		}

		if errorCount > 0 {
			fmt.Fprintln(os.Stderr, errorCount, "errors occured!")
			os.Exit(1)
		}
	},
}

// expandPath makes sure environment variables and "~" are expanded in a path
func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") || strings.HasPrefix(path, "~\\") {
		path = homeDir + path[1:]
	}
	return os.ExpandEnv(path)
}

// sortedToolNames returns all configured tool names in a deterministic order: first all wrapper scripts, then all aliases, in both cases alphabetical
func sortedToolNames() []string {
	toolNames := make([]string, 0, len(config.Tools))
	for toolName := range config.Tools {
		toolNames = append(toolNames, toolName)
	}

	sort.Strings(toolNames)
	sort.SliceStable(toolNames, func(p, q int) bool { return config.Tools[toolNames[p]].Type > config.Tools[toolNames[q]].Type })

	return toolNames
}

// generateWrapperScript renders the template with the correct configuration, but also makes sure beforehand that
// local folders for mount points are expanded and exist
func generateWrapperScript(outputDir, templateFile, toolName string, toolConfig cfg.ToolConfig) error {
	for i, mount := range toolConfig.Mounts {
		mount.Source = expandPath(mount.Source)
		toolConfig.Mounts[i] = mount
	}

	createMountFolders(toolConfig)

	toolParams := gen.ToolParams{
		Engine:       config.Engine,
		Name:         toolName,
		ImageName:    toolConfig.ImageName,
		ImageTag:     toolConfig.ImageTag,
		WorkDir:      toolConfig.WorkDir,
		HomeDir:      toolConfig.HomeDir,
		Command:      toolConfig.Command,
		Mounts:       toolConfig.Mounts,
		Env:          toolConfig.Env,
		NetworkMode:  toolConfig.NetworkMode,
		CustomScript: toolConfig.CustomScript,
	}
	return gen.Generate(outputDir, templateFile, toolParams)
}

// createMountFolders tries to create all local folders for mount points if they do not yet exist - errors
// are considered noncritical and are thus only written to stderr
func createMountFolders(toolConfig cfg.ToolConfig) {
	for _, mount := range toolConfig.Mounts {
		_, err := os.Stat(mount.Source)
		if err == nil {
			continue
		}
		if !errors.Is(err, fs.ErrNotExist) {
			fmt.Fprintln(os.Stderr, "Warning:", err)
			continue
		}
		fmt.Println("Information: mount source path does not exist, creating folder:", mount.Source)
		if err := os.MkdirAll(mount.Source, 0755); err != nil {
			fmt.Fprintln(os.Stderr, "Warning:", err)
		}
	}
}
