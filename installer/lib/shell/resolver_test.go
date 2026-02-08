package shell_test

import (
	"testing"

	"github.com/MrPointer/dotfiles/installer/lib/shell"
	"github.com/MrPointer/dotfiles/installer/utils"
	"github.com/MrPointer/dotfiles/installer/utils/logger"
	"github.com/MrPointer/dotfiles/installer/utils/osmanager"
	"github.com/stretchr/testify/require"
)

func Test_ResolverReturnsBrewPath_WhenSourceIsBrew(t *testing.T) {
	// Arrange
	shellName := "zsh"
	brewPath := "/opt/homebrew"
	expectedPath := "/opt/homebrew/bin/zsh"

	fileSystemMock := &utils.MoqFileSystem{
		PathExistsFunc: func(path string) (bool, error) {
			return path == expectedPath, nil
		},
	}

	resolver := shell.NewDefaultShellResolver(
		shellName,
		shell.ShellSourceBrew,
		brewPath,
		"darwin",
		&osmanager.MoqOsManager{},
		fileSystemMock,
		logger.DefaultLogger,
	)

	// Act
	path, err := resolver.GetShellPath()

	// Assert
	require.NoError(t, err)
	require.Equal(t, expectedPath, path)
}

func Test_ResolverReturnsSystemPath_WhenSourceIsSystem(t *testing.T) {
	// Arrange
	shellName := "zsh"
	expectedPath := "/usr/bin/zsh"

	fileSystemMock := &utils.MoqFileSystem{
		PathExistsFunc: func(path string) (bool, error) {
			return path == expectedPath, nil
		},
	}

	resolver := shell.NewDefaultShellResolver(
		shellName,
		shell.ShellSourceSystem,
		"", // no brew path
		"linux",
		&osmanager.MoqOsManager{},
		fileSystemMock,
		logger.DefaultLogger,
	)

	// Act
	path, err := resolver.GetShellPath()

	// Assert
	require.NoError(t, err)
	require.Equal(t, expectedPath, path)
}

func Test_ResolverPrefersBrewInAutoMode_WhenBrewPathAvailable(t *testing.T) {
	// Arrange
	shellName := "zsh"
	brewPath := "/opt/homebrew"
	brewShellPath := "/opt/homebrew/bin/zsh"

	fileSystemMock := &utils.MoqFileSystem{
		PathExistsFunc: func(path string) (bool, error) {
			return path == brewShellPath, nil
		},
	}

	resolver := shell.NewDefaultShellResolver(
		shellName,
		shell.ShellSourceAuto,
		brewPath,
		"darwin",
		&osmanager.MoqOsManager{},
		fileSystemMock,
		logger.DefaultLogger,
	)

	// Act
	path, err := resolver.GetShellPath()

	// Assert
	require.NoError(t, err)
	require.Equal(t, brewShellPath, path)
}

func Test_ResolverFallsBackToSystem_WhenBrewShellNotFound(t *testing.T) {
	// Arrange
	shellName := "zsh"
	brewPath := "/opt/homebrew"
	systemPath := "/usr/bin/zsh"

	fileSystemMock := &utils.MoqFileSystem{
		PathExistsFunc: func(path string) (bool, error) {
			return path == systemPath, nil // Only system path exists
		},
	}

	resolver := shell.NewDefaultShellResolver(
		shellName,
		shell.ShellSourceAuto,
		brewPath,
		"linux",
		&osmanager.MoqOsManager{},
		fileSystemMock,
		logger.DefaultLogger,
	)

	// Act
	path, err := resolver.GetShellPath()

	// Assert
	require.NoError(t, err)
	require.Equal(t, systemPath, path)
}

func Test_ResolverReturnsError_WhenShellNotFound(t *testing.T) {
	// Arrange
	shellName := "nonexistent"

	fileSystemMock := &utils.MoqFileSystem{
		PathExistsFunc: func(path string) (bool, error) {
			return false, nil // Nothing exists
		},
	}

	resolver := shell.NewDefaultShellResolver(
		shellName,
		shell.ShellSourceSystem,
		"",
		"linux",
		&osmanager.MoqOsManager{},
		fileSystemMock,
		logger.DefaultLogger,
	)

	// Act
	_, err := resolver.GetShellPath()

	// Assert
	require.Error(t, err)
	require.Contains(t, err.Error(), "not found")
}

func Test_IsAvailableReturnsTrue_WhenShellExists(t *testing.T) {
	// Arrange
	shellName := "zsh"
	shellPath := "/usr/bin/zsh"

	fileSystemMock := &utils.MoqFileSystem{
		PathExistsFunc: func(path string) (bool, error) {
			return path == shellPath, nil
		},
	}

	resolver := shell.NewDefaultShellResolver(
		shellName,
		shell.ShellSourceSystem,
		"",
		"linux",
		&osmanager.MoqOsManager{},
		fileSystemMock,
		logger.DefaultLogger,
	)

	// Act
	available, err := resolver.IsAvailable()

	// Assert
	require.NoError(t, err)
	require.True(t, available)
}

func Test_IsAvailableReturnsFalse_WhenShellNotFound(t *testing.T) {
	// Arrange
	shellName := "nonexistent"

	fileSystemMock := &utils.MoqFileSystem{
		PathExistsFunc: func(path string) (bool, error) {
			return false, nil
		},
	}

	resolver := shell.NewDefaultShellResolver(
		shellName,
		shell.ShellSourceSystem,
		"",
		"linux",
		&osmanager.MoqOsManager{},
		fileSystemMock,
		logger.DefaultLogger,
	)

	// Act
	available, err := resolver.IsAvailable()

	// Assert
	require.NoError(t, err) // IsAvailable should not error when shell not found, just return false
	require.False(t, available)
}
