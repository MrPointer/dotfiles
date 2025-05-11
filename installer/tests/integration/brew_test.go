package brew_test

import (
	"os"
	"runtime"
	"testing"

	"github.com/MrPointer/dotfiles/installer/lib/brew"
	"github.com/MrPointer/dotfiles/installer/lib/compatibility"
	"github.com/MrPointer/dotfiles/installer/utils/logger"
)

func TestBrewReportedAsUnavailableWhenNotInstalled(t *testing.T) {
	opts := brew.Options{
		Logger:           logger.DefaultLogger,
		SystemInfo:       &compatibility.SystemInfo{OSName: runtime.GOOS, Arch: runtime.GOARCH},
		Commander:        nil,
		BrewPathOverride: "/tmp/nonexistent-brew-binary",
	}
	installer := brew.NewBrewInstaller(opts)

	available, err := installer.IsAvailable()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if available {
		t.Error("expected not available, got available")
	}
}

func TestBrewReportedAsAvailableWhenInstalled(t *testing.T) {
	tempFile, err := os.CreateTemp("", "brew-binary-*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	opts := brew.Options{
		Logger:           logger.DefaultLogger,
		SystemInfo:       &compatibility.SystemInfo{OSName: runtime.GOOS, Arch: runtime.GOARCH},
		Commander:        nil,
		BrewPathOverride: tempFile.Name(),
	}
	installer := brew.NewBrewInstaller(opts)

	available, err := installer.IsAvailable()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !available {
		t.Error("expected available, got not available")
	}
}
