package config

import (
	"embed"
	"fmt"
)

//go:embed compatibility.yaml packagemap.yaml
var configFS embed.FS

// GetRawEmbeddedCompatibilityConfig returns the raw content of the embedded compatibility configuration file.
// This is useful for testing purposes.
func GetRawEmbeddedCompatibilityConfig() ([]byte, error) {
	data, err := configFS.ReadFile("compatibility.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded compatibility config: %w", err)
	}
	return data, nil
}

// GetRawEmbeddedPackageMapConfig returns the raw content of the embedded package map configuration file.
// This is useful for loading the default configuration or for testing purposes.
func GetRawEmbeddedPackageMapConfig() ([]byte, error) {
	data, err := configFS.ReadFile("packagemap.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded package map config: %w", err)
	}
	return data, nil
}