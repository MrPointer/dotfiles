package config

import (
	"testing"
)

func TestEmbeddedCompatibilityConfigCanBeLoaded(t *testing.T) {
	// Test basic loading functionality
	config, err := GetRawEmbeddedCompatibilityConfig()
	if err != nil {
		t.Fatalf("Expected no error when loading embedded config, got: %v", err)
	}

	if config == nil {
		t.Fatal("Expected config to be non-nil, got nil")
	}
}
