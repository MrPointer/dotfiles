package brew_test

import (
	"os"
	"runtime"
	"testing"

	"github.com/MrPointer/dotfiles/installer/lib/brew"
	"github.com/MrPointer/dotfiles/installer/lib/compatibility"
	"github.com/MrPointer/dotfiles/installer/utils/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_BrewReportedAsUnavailable_WhenNotInstalled(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Arrange
	opts := brew.DefaultOptions().
		WithLogger(logger.DefaultLogger).
		WithSystemInfo(&compatibility.SystemInfo{OSName: runtime.GOOS, Arch: runtime.GOARCH}).
		WithBrewPathOverride("/tmp/nonexistent-brew-binary")

	installer := brew.NewBrewInstaller(*opts)

	// Act
	available, err := installer.IsAvailable()

	// Assert
	require.NoError(t, err)
	assert.False(t, available, "expected brew to be unavailable when not installed")
}

func Test_BrewReportedAsAvailable_WhenInstalled(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Arrange
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

	// Act
	available, err := installer.IsAvailable()

	// Assert
	require.NoError(t, err)
	assert.True(t, available, "expected brew to be available when installed")
}
