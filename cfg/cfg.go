package cfg

import (
	"strings"
)

// Config represents the whole config file
type Config struct {
	OutputDir string
	Engine    string
	Tools     map[string]ToolConfig
}

// ToolType determines which kind of tool should be generated
type ToolType string

// Values for ToolType
const (
	WrapperScript = "WrapperScript"
	Alias         = "Alias"
)

// ImageTagType determines how the container image tag is determined
type ImageTagType string

// Values for ImageTagType
const (
	Fixed    = "Fixed"
	FromFile = "FromFile"
)

// CommandType determines how to infer the executable name for the container
type CommandType string

// Values for CommandType
const (
	DoNotSpecify = "DoNotSpecify"
	ReuseName    = "ReuseName"
)

// ToolConfig represents the configuration for a CLI tool in the config file
type ToolConfig struct {
	Type ToolType

	// Following keys only when Type=WrapperScript
	ImageName    string
	ImageTag     ImageTagConfig
	WorkDir      string
	HomeDir      string
	Command      CommandConfig
	Mounts       []MountConfig
	Env          []EnvConfig
	CustomScript string

	// Following keys only when Type=Alias
	AliasFor string
}

type ImageTagConfig struct {
	Type ImageTagType
	// Following keys only when Type=Fixed
	Value string
	// Following keys only when Type=FromFile
	File     string
	Sed      []string
	Fallback string
}

type CommandConfig struct {
	Type CommandType
	// Following keys only when Type=ReuseName
	Folder string
}

type MountConfig struct {
	Source string
	Target string
}

type EnvConfig struct {
	Name  string
	Value string
}

func (toolType ToolType) IsWrapperScript() bool {
	return strings.EqualFold(string(toolType), WrapperScript)
}

func (toolType ToolType) IsAlias() bool {
	return strings.EqualFold(string(toolType), Alias)
}

func (imageTagType ImageTagType) IsFixed() bool {
	return strings.EqualFold(string(imageTagType), Fixed)
}

func (imageTagType ImageTagType) IsFromFile() bool {
	return strings.EqualFold(string(imageTagType), FromFile)
}

func (commandType CommandType) IsDoNotSpecify() bool {
	return strings.EqualFold(string(commandType), DoNotSpecify)
}

func (commandType CommandType) IsReuseName() bool {
	return strings.EqualFold(string(commandType), ReuseName)
}
