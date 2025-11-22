package compatibility

import (
	"bytes"
	"fmt"

	"github.com/MrPointer/dotfiles/installer/internal/config"
	"github.com/spf13/viper"
)

// LoadCompatibilityConfig loads compatibility config from file or embedded source.
func LoadCompatibilityConfig(v *viper.Viper, compatibilityConfigFile string) (*CompatibilityConfig, error) {
	// Create a separate viper instance for compatibility config
	var compatibilityConfig CompatibilityConfig

	// If compatibility config file is specified, load from there
	if compatibilityConfigFile != "" {
		v.SetConfigFile(compatibilityConfigFile)

		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("error reading compatibility config file: %v", err)
		}

		fmt.Println("Using compatibility config file:", v.ConfigFileUsed())

		if err := v.Unmarshal(&compatibilityConfig); err != nil {
			return nil, fmt.Errorf("error parsing compatibility config: %v", err)
		}
	} else {
		// If no file is specified, use the embedded compatibility config
		v.SetConfigType("yaml")

		embedded_config, err := config.GetRawEmbeddedCompatibilityConfig()
		if err != nil {
			return nil, fmt.Errorf("error loading embedded compatibility config: %v", err)
		}

		if err := v.ReadConfig(bytes.NewBuffer(embedded_config)); err != nil {
			return nil, fmt.Errorf("error reading embedded compatibility config: %v", err)
		}

		if err := v.Unmarshal(&compatibilityConfig); err != nil {
			return nil, fmt.Errorf("error parsing embedded compatibility config: %v", err)
		}
	}

	return &compatibilityConfig, nil
}
