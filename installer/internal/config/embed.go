package config

import (
	"embed"
	"fmt"
)

//go:embed compatibility.yaml
var configFS embed.FS

// GetRawEmbeddedCompatibilityConfig returns the raw content of an embedded file
// This is useful for testing purposes.
func GetRawEmbeddedCompatibilityConfig() ([]byte, error) {
	// Read the embedded compatibility config
	data, err := configFS.ReadFile("compatibility.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded compatibility config: %w", err)
	}

	return data, nil
}
