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
