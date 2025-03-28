package compatibility

import (
	"path/filepath"
	"testing"

	"github.com/MrPointer/dotfiles/installer/utils/files/current"
	"github.com/spf13/viper"
)

func TestCompatibilityConfigCanBeLoadedFromEmbeddedSource(t *testing.T) {
	// Create a new Viper instance
	v := viper.New()

	// Load the embedded compatibility config
	compatibilityConfig, err := LoadCompatibilityConfig(v, "")
	if err != nil {
		t.Fatalf("Expected no error when loading embedded compatibility config, got: %v", err)
	}

	if compatibilityConfig == nil {
		t.Fatal("Expected compatibility config to be non-nil, got nil")
	}
}

func TestCompatibilityConfigCanBeLoadedFromFile(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}

	rootDir, err := current.RootDirectory()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// Test loading from a file
	compatibilityConfigFile := filepath.Join(rootDir, "internal", "config", "compatibility.yaml")
	v := viper.New()
	v.SetConfigFile(compatibilityConfigFile)

	compatibilityConfig, err := LoadCompatibilityConfig(v, compatibilityConfigFile)
	if err != nil {
		t.Fatalf("Expected no error when loading compatibility config from file, got: %v", err)
	}

	if compatibilityConfig == nil {
		t.Fatal("Expected compatibility config to be non-nil, got nil")
	}
}
