package current

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRootDirectoryContainsMain(t *testing.T) {
	// Test that the root directory contains the main.go file
	rootDir, err := RootDirectory()
	if err != nil {
		t.Fatalf("Expected no error when getting root directory, got: %v", err)
	}

	mainFilePath := filepath.Join(rootDir, "main.go")
	if _, err := os.Stat(mainFilePath); os.IsNotExist(err) {
		t.Fatalf("Expected main.go to exist in root directory, but it does not: %v", err)
	}
}
