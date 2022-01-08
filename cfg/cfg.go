package cfg

// Config represents the whole config file
type Config struct {
	OutputDir string
	Tools     map[string]ToolConfig
}

// ToolConfig represents the configuration for a CLI tool in the config file
type ToolConfig struct {
	ImageName string
	AliasFor  string
}
