package packageresolver

import (
	"bytes"
	"fmt"

	"github.com/MrPointer/dotfiles/installer/internal/config"
	"github.com/spf13/viper"
)

// LoadPackageMappings loads the package mapping configuration.
// It first tries to load from the specified `packageMapFile`. If `packageMapFile` is empty
// or loading fails and fallback is implicitly desired, it loads from the embedded default configuration.
func LoadPackageMappings(v *viper.Viper, packageMapFile string) (*PackageMappingCollection, error) {
	var mappingsCfg PackageMappingCollection

	if packageMapFile != "" {
		v.SetConfigFile(packageMapFile)
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("error reading package map file '%s': %w", packageMapFile, err)
		}
		// Consider adding logging for which config file is used, e.g.:
		// fmt.Println("Using package map file:", v.ConfigFileUsed())
	} else {
		// Use embedded configuration if no file is specified
		v.SetConfigType("yaml")
		embeddedData, err := config.GetRawEmbeddedPackageMapConfig()
		if err != nil {
			return nil, fmt.Errorf("error loading embedded package map config: %w", err)
		}
		if err := v.ReadConfig(bytes.NewBuffer(embeddedData)); err != nil {
			return nil, fmt.Errorf("error reading embedded package map config: %w", err)
		}
	}

	if err := v.Unmarshal(&mappingsCfg); err != nil {
		return nil, fmt.Errorf("error parsing package map configuration: %w", err)
	}

	// Ensure Packages map is initialized to prevent nil pointer dereference later.
	// This is important if the "packages" key is missing or empty in the config.
	if mappingsCfg.Packages == nil {
		mappingsCfg.Packages = make(map[string]PackageMapping)
	}

	return &mappingsCfg, nil
}
