package brew_test

import (
	"os"
	"runtime"
	"testing"

	"github.com/MrPointer/dotfiles/installer/lib/brew"
	"github.com/MrPointer/dotfiles/installer/lib/compatibility"
	"github.com/MrPointer/dotfiles/installer/utils/logger"
)

func Test_BrewReportedAsUnavailable_WhenNotInstalled(t *testing.T) {
	opts := brew.DefaultOptions().
		WithLogger(logger.DefaultLogger).
		WithSystemInfo(&compatibility.SystemInfo{OSName: runtime.GOOS, Arch: runtime.GOARCH}).
		WithBrewPathOverride("/tmp/nonexistent-brew-binary")

	installer := brew.NewBrewInstaller(*opts)

	available, err := installer.IsAvailable()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if available {
		t.Error("expected not available, got available")
	}
}

func Test_BrewReportedAsAvailable_WhenInstalled(t *testing.T) {
	tempFile, err := os.CreateTemp("", "brew-binary-*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	opts := brew.DefaultOptions().
		WithLogger(logger.DefaultLogger).
		WithSystemInfo(&compatibility.SystemInfo{OSName: runtime.GOOS, Arch: runtime.GOARCH}).
		WithBrewPathOverride(tempFile.Name())

	installer := brew.NewBrewInstaller(*opts)

	available, err := installer.IsAvailable()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !available {
		t.Error("expected available, got not available")
	}
}
