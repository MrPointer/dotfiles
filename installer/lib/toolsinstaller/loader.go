package toolsinstaller

import (
	"bytes"
	"fmt"

	"github.com/MrPointer/dotfiles/installer/internal/config"
	"github.com/spf13/viper"
)

// LoadToolsConfig loads the tools configuration.
// It first tries to load from the specified toolsConfigFile. If toolsConfigFile is empty,
// it loads from the embedded default configuration.
// Callers must pass viper.New() to avoid state pollution from other config loaders.
func LoadToolsConfig(v *viper.Viper, toolsConfigFile string) (*ToolsConfig, error) {
	if toolsConfigFile != "" {
		v.SetConfigFile(toolsConfigFile)
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("error reading tools config file '%s': %w", toolsConfigFile, err)
		}
	} else {
		v.SetConfigType("yaml")

		embeddedData, err := config.GetRawEmbeddedToolsConfig()
		if err != nil {
			return nil, fmt.Errorf("error loading embedded tools config: %w", err)
		}
		if err := v.ReadConfig(bytes.NewBuffer(embeddedData)); err != nil {
			return nil, fmt.Errorf("error reading embedded tools config: %w", err)
		}
	}

	var toolsCfg ToolsConfig
	if err := v.Unmarshal(&toolsCfg); err != nil {
		return nil, fmt.Errorf("error parsing tools configuration: %w", err)
	}

	return &toolsCfg, nil
}
