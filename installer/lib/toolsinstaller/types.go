package toolsinstaller

// ToolDefinition represents a single tool entry from tools.yaml.
type ToolDefinition struct {
	Name        string `mapstructure:"name"`
	Description string `mapstructure:"description"`
}

// ToolsConfig is the top-level structure for tools.yaml.
type ToolsConfig struct {
	Tools []ToolDefinition `mapstructure:"tools"`
}
