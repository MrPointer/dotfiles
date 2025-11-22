package current_test

import (
	"strings"
	"testing"

	"github.com/MrPointer/dotfiles/installer/utils/collections"
	"github.com/MrPointer/dotfiles/installer/utils/files/current"
)

func TestCurrentFileIsSelf(t *testing.T) {
	// Get the current file path
	currentFile, err := current.Filename()
	if err != nil {
		t.Fatalf("Failed to get current file: %v", err)
	}

	// Get the expected last path element
	expectedLastElement := "current_test.go"

	// Check if the current file is the expected file by comparing the last path element
	lastPathElement, err := collections.Last(strings.Split(currentFile, "/"))
	if err != nil {
		t.Fatalf("Failed to get last path element of returned current file: %v", err)
	}

	if lastPathElement != expectedLastElement {
		t.Errorf("Expected current file to be %s, got %s", expectedLastElement, lastPathElement)
	}
}

func TestCurrentDirIsSelf(t *testing.T) {
	// Get the current directory path
	currentDir, err := current.Dirname()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// Get the expected last path element
	expectedLastElement := "current"

	// Check if the current directory is the expected directory by comparing the last path element
	lastPathElement, err := collections.Last(strings.Split(currentDir, "/"))
	if err != nil {
		t.Fatalf("Failed to get last path element of returned current directory: %v", err)
	}

	if lastPathElement != expectedLastElement {
		t.Errorf("Expected current directory to be %s, got %s", expectedLastElement, lastPathElement)
	}
}
