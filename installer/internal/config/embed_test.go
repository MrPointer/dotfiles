package config_test

import (
	"testing"

	"github.com/MrPointer/dotfiles/installer/internal/config"
)

func TestEmbeddedCompatibilityConfigCanBeLoaded(t *testing.T) {
	// Test basic loading functionality
	config, err := config.GetRawEmbeddedCompatibilityConfig()
	if err != nil {
		t.Fatalf("Expected no error when loading embedded config, got: %v", err)
	}

	if config == nil {
		t.Fatal("Expected config to be non-nil, got nil")
	}
}

func TestEmbeddedPackageMapConfigCanBeLoaded(t *testing.T) {
	// Test basic loading functionality
	config, err := config.GetRawEmbeddedPackageMapConfig()
	if err != nil {
		t.Fatalf("Expected no error when loading embedded package map config, got: %v", err)
	}

	if config == nil {
		t.Fatal("Expected package map config to be non-nil, got nil")
	}
}
