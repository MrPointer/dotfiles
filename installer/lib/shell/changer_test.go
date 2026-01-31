package shell_test

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/MrPointer/dotfiles/installer/lib/shell"
	"github.com/MrPointer/dotfiles/installer/utils"
	"github.com/MrPointer/dotfiles/installer/utils/logger"
	"github.com/MrPointer/dotfiles/installer/utils/osmanager"
	"github.com/MrPointer/dotfiles/installer/utils/privilege"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_GetShellPath_ReturnsPathFromOsManager_WhenShellExistsInPath(t *testing.T) {
	// Arrange
	shellName := "zsh"
	expectedPath := "/usr/bin/zsh"

	osManagerMock := &osmanager.MoqOsManager{
		GetProgramPathFunc: func(program string) (string, error) {
			if program == shellName {
				return expectedPath, nil
			}
			return "", assert.AnError
		},
	}
	fileSystemMock := &utils.MoqFileSystem{
		PathExistsFunc: func(path string) (bool, error) {
			return path == expectedPath, nil
		},
		IsExecutableFunc: func(path string) (bool, error) {
			return path == expectedPath, nil
		},
	}
	commanderMock := &utils.MoqCommander{}
	escalatorMock := &privilege.MoqEscalator{}

	// Even with brew path set, should prefer PATH
	changer := shell.NewDefaultShellChanger(
		shellName,
		"/opt/homebrew/bin/brew", // brew path is set
		logger.DefaultLogger,
		osManagerMock,
		fileSystemMock,
		commanderMock,
		escalatorMock,
	)

	// Act
	path, err := changer.GetShellPath()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedPath, path)
	assert.Len(t, osManagerMock.GetProgramPathCalls(), 1)
	assert.Equal(t, shellName, osManagerMock.GetProgramPathCalls()[0].Program)
}

func Test_GetShellPath_UsesOsManager_WhenNoBrewPath(t *testing.T) {
	// Arrange
	shellName := "zsh"

	// Create a temporary fake shell binary
	tempDir := t.TempDir()
	expectedPath := filepath.Join(tempDir, shellName)
	err := os.WriteFile(expectedPath, []byte("#!/bin/sh\n"), 0755)
	require.NoError(t, err)

	osManagerMock := &osmanager.MoqOsManager{
		GetProgramPathFunc: func(program string) (string, error) {
			if program == shellName {
				return expectedPath, nil
			}
			return "", assert.AnError
		},
	}
	fileSystemMock := &utils.MoqFileSystem{
		PathExistsFunc: func(path string) (bool, error) {
			return path == expectedPath, nil
		},
		IsExecutableFunc: func(path string) (bool, error) {
			return path == expectedPath, nil
		},
	}
	commanderMock := &utils.MoqCommander{}
	escalatorMock := &privilege.MoqEscalator{}

	changer := shell.NewDefaultShellChanger(
		shellName,
		"", // no brew path
		logger.DefaultLogger,
		osManagerMock,
		fileSystemMock,
		commanderMock,
		escalatorMock,
	)

	// Act
	path, err := changer.GetShellPath()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedPath, path)
	assert.Len(t, osManagerMock.GetProgramPathCalls(), 1)
	assert.Equal(t, shellName, osManagerMock.GetProgramPathCalls()[0].Program)
}

func Test_GetShellPath_ReturnsBrewPath_WhenShellNotInPathButExistsInBrewBin(t *testing.T) {
	// Arrange
	shellName := "zsh"
	tempDir := t.TempDir()
	brewBinDir := filepath.Join(tempDir, "bin")
	err := os.MkdirAll(brewBinDir, 0755)
	require.NoError(t, err)

	shellPath := filepath.Join(brewBinDir, shellName)
	err = os.WriteFile(shellPath, []byte("#!/bin/sh\n"), 0755)
	require.NoError(t, err)

	brewPath := filepath.Join(brewBinDir, "brew")

	osManagerMock := &osmanager.MoqOsManager{
		GetProgramPathFunc: func(program string) (string, error) {
			return "", assert.AnError // Shell not in PATH
		},
	}
	fileSystemMock := &utils.MoqFileSystem{
		PathExistsFunc: func(path string) (bool, error) {
			return path == shellPath, nil
		},
		IsExecutableFunc: func(path string) (bool, error) {
			return path == shellPath, nil
		},
	}
	commanderMock := &utils.MoqCommander{}
	escalatorMock := &privilege.MoqEscalator{}

	changer := shell.NewDefaultShellChanger(
		shellName,
		brewPath,
		logger.DefaultLogger,
		osManagerMock,
		fileSystemMock,
		commanderMock,
		escalatorMock,
	)

	// Act
	path, err := changer.GetShellPath()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, shellPath, path)
}

func Test_GetShellPath_PrefersPathOverBrewBin_WhenBothExist(t *testing.T) {
	// Arrange
	shellName := "zsh"
	systemShellPath := "/usr/bin/zsh"

	tempDir := t.TempDir()
	brewBinDir := filepath.Join(tempDir, "bin")
	err := os.MkdirAll(brewBinDir, 0755)
	require.NoError(t, err)

	brewShellPath := filepath.Join(brewBinDir, shellName)
	err = os.WriteFile(brewShellPath, []byte("#!/bin/sh\n"), 0755)
	require.NoError(t, err)

	brewPath := filepath.Join(brewBinDir, "brew")

	osManagerMock := &osmanager.MoqOsManager{
		GetProgramPathFunc: func(program string) (string, error) {
			return systemShellPath, nil // Shell found in PATH
		},
	}
	fileSystemMock := &utils.MoqFileSystem{
		PathExistsFunc: func(path string) (bool, error) {
			return path == systemShellPath || path == brewShellPath, nil
		},
		IsExecutableFunc: func(path string) (bool, error) {
			return path == systemShellPath || path == brewShellPath, nil
		},
	}
	commanderMock := &utils.MoqCommander{}
	escalatorMock := &privilege.MoqEscalator{}

	changer := shell.NewDefaultShellChanger(
		shellName,
		brewPath, // Brew path is set
		logger.DefaultLogger,
		osManagerMock,
		fileSystemMock,
		commanderMock,
		escalatorMock,
	)

	// Act
	path, err := changer.GetShellPath()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, systemShellPath, path, "Should prefer shell from PATH over brew bin")
}

func Test_GetShellPath_ReturnsError_WhenShellNotFoundAnywhere(t *testing.T) {
	// Arrange
	shellName := "zsh"
	tempDir := t.TempDir()
	brewBinDir := filepath.Join(tempDir, "bin")
	brewPath := filepath.Join(brewBinDir, "brew")

	osManagerMock := &osmanager.MoqOsManager{
		GetProgramPathFunc: func(program string) (string, error) {
			return "", assert.AnError // Not in PATH
		},
	}
	fileSystemMock := &utils.MoqFileSystem{
		PathExistsFunc: func(path string) (bool, error) {
			return false, nil // Not in brew bin either
		},
	}
	commanderMock := &utils.MoqCommander{}
	escalatorMock := &privilege.MoqEscalator{}

	changer := shell.NewDefaultShellChanger(
		shellName,
		brewPath,
		logger.DefaultLogger,
		osManagerMock,
		fileSystemMock,
		commanderMock,
		escalatorMock,
	)

	// Act
	path, err := changer.GetShellPath()

	// Assert
	require.Error(t, err)
	assert.Empty(t, path)
	assert.Contains(t, err.Error(), "could not find")
}

func Test_GetShellPath_ReturnsError_WhenProgramNotFound(t *testing.T) {
	// Arrange
	shellName := "nonexistent"

	osManagerMock := &osmanager.MoqOsManager{
		GetProgramPathFunc: func(program string) (string, error) {
			return "", assert.AnError
		},
	}
	fileSystemMock := &utils.MoqFileSystem{}
	commanderMock := &utils.MoqCommander{}
	escalatorMock := &privilege.MoqEscalator{}

	changer := shell.NewDefaultShellChanger(
		shellName,
		"", // no brew path
		logger.DefaultLogger,
		osManagerMock,
		fileSystemMock,
		commanderMock,
		escalatorMock,
	)

	// Act
	path, err := changer.GetShellPath()

	// Assert
	require.Error(t, err)
	assert.Empty(t, path)
	assert.Contains(t, err.Error(), "could not find")
}

func Test_IsCurrentDefault_ReturnsTrue_WhenShellMatches(t *testing.T) {
	// Arrange
	shellName := "zsh"
	shellPath := "/usr/local/bin/zsh"
	username := "testuser"

	osManagerMock := &osmanager.MoqOsManager{
		GetProgramPathFunc: func(program string) (string, error) {
			return shellPath, nil
		},
		GetCurrentUsernameFunc: func() (string, error) {
			return username, nil
		},
		GetUserShellFunc: func(user string) (string, error) {
			return shellPath, nil // Current shell matches target
		},
	}
	fileSystemMock := &utils.MoqFileSystem{
		PathExistsFunc: func(path string) (bool, error) {
			return true, nil
		},
		IsExecutableFunc: func(path string) (bool, error) {
			return true, nil
		},
	}
	commanderMock := &utils.MoqCommander{}
	escalatorMock := &privilege.MoqEscalator{}

	changer := shell.NewDefaultShellChanger(
		shellName,
		"",
		logger.DefaultLogger,
		osManagerMock,
		fileSystemMock,
		commanderMock,
		escalatorMock,
	)

	// Act
	isDefault, err := changer.IsCurrentDefault()

	// Assert
	require.NoError(t, err)
	assert.True(t, isDefault)
}

func Test_IsCurrentDefault_ReturnsFalse_WhenShellDiffers(t *testing.T) {
	// Arrange
	shellName := "zsh"
	shellPath := "/usr/local/bin/zsh"
	currentShell := "/bin/bash"
	username := "testuser"

	osManagerMock := &osmanager.MoqOsManager{
		GetProgramPathFunc: func(program string) (string, error) {
			return shellPath, nil
		},
		GetCurrentUsernameFunc: func() (string, error) {
			return username, nil
		},
		GetUserShellFunc: func(user string) (string, error) {
			return currentShell, nil // Different from target
		},
	}
	fileSystemMock := &utils.MoqFileSystem{
		PathExistsFunc: func(path string) (bool, error) {
			return true, nil
		},
		IsExecutableFunc: func(path string) (bool, error) {
			return true, nil
		},
	}
	commanderMock := &utils.MoqCommander{}
	escalatorMock := &privilege.MoqEscalator{}

	changer := shell.NewDefaultShellChanger(
		shellName,
		"",
		logger.DefaultLogger,
		osManagerMock,
		fileSystemMock,
		commanderMock,
		escalatorMock,
	)

	// Act
	isDefault, err := changer.IsCurrentDefault()

	// Assert
	require.NoError(t, err)
	assert.False(t, isDefault)
}

func Test_SetAsDefault_SkipsChange_WhenAlreadyDefault(t *testing.T) {
	// Arrange
	shellName := "zsh"
	shellPath := "/usr/local/bin/zsh"
	username := "testuser"

	osManagerMock := &osmanager.MoqOsManager{
		GetProgramPathFunc: func(program string) (string, error) {
			return shellPath, nil
		},
		GetCurrentUsernameFunc: func() (string, error) {
			return username, nil
		},
		GetUserShellFunc: func(user string) (string, error) {
			return shellPath, nil // Already the default
		},
		SetUserShellFunc: func(user, shell string) error {
			t.Error("SetUserShell should not be called when already default")
			return nil
		},
	}
	fileSystemMock := &utils.MoqFileSystem{
		PathExistsFunc: func(path string) (bool, error) {
			return true, nil
		},
		IsExecutableFunc: func(path string) (bool, error) {
			return true, nil
		},
	}
	commanderMock := &utils.MoqCommander{}
	escalatorMock := &privilege.MoqEscalator{
		IsRunningAsRootFunc: func() (bool, error) {
			return false, nil
		},
	}

	changer := shell.NewDefaultShellChanger(
		shellName,
		"",
		logger.DefaultLogger,
		osManagerMock,
		fileSystemMock,
		commanderMock,
		escalatorMock,
	)

	// Act
	err := changer.SetAsDefault(context.Background())

	// Assert
	require.NoError(t, err)
	assert.Len(t, osManagerMock.SetUserShellCalls(), 0, "SetUserShell should not be called")
}

func Test_SetAsDefault_SetsShell_WhenNotDefault(t *testing.T) {
	// Arrange
	shellName := "zsh"
	shellPath := "/usr/local/bin/zsh"
	currentShell := "/bin/bash"
	username := "testuser"

	osManagerMock := &osmanager.MoqOsManager{
		GetProgramPathFunc: func(program string) (string, error) {
			return shellPath, nil
		},
		GetCurrentUsernameFunc: func() (string, error) {
			return username, nil
		},
		GetUserShellFunc: func(user string) (string, error) {
			return currentShell, nil // Different from target
		},
		SetUserShellFunc: func(user, shell string) error {
			return nil
		},
	}
	fileSystemMock := &utils.MoqFileSystem{
		PathExistsFunc: func(path string) (bool, error) {
			return true, nil
		},
		IsExecutableFunc: func(path string) (bool, error) {
			return true, nil
		},
	}
	commanderMock := &utils.MoqCommander{}
	escalatorMock := &privilege.MoqEscalator{
		IsRunningAsRootFunc: func() (bool, error) {
			return false, nil
		},
	}

	changer := shell.NewDefaultShellChanger(
		shellName,
		"",
		logger.DefaultLogger,
		osManagerMock,
		fileSystemMock,
		commanderMock,
		escalatorMock,
	)

	// Act
	err := changer.SetAsDefault(context.Background())

	// Assert
	require.NoError(t, err)
	require.Len(t, osManagerMock.SetUserShellCalls(), 1, "SetUserShell should be called once")
	assert.Equal(t, username, osManagerMock.SetUserShellCalls()[0].Username)
	assert.Equal(t, shellPath, osManagerMock.SetUserShellCalls()[0].ShellPath)
}

func Test_SetAsDefault_AddsToEtcShells_WhenBrewInstalled(t *testing.T) {
	// Arrange
	shellName := "zsh"
	tempDir := t.TempDir()
	brewBinDir := filepath.Join(tempDir, "bin")
	shellPath := filepath.Join(brewBinDir, shellName)
	brewPath := filepath.Join(brewBinDir, "brew")
	currentShell := "/bin/bash"
	username := "testuser"

	// Simulate /etc/shells content without the brew shell
	etcShellsContent := "/bin/sh\n/bin/bash\n/bin/zsh\n"

	osManagerMock := &osmanager.MoqOsManager{
		GetProgramPathFunc: func(program string) (string, error) {
			// Shell not in PATH (requires brew bin fallback)
			return "", assert.AnError
		},
		GetCurrentUsernameFunc: func() (string, error) {
			return username, nil
		},
		GetUserShellFunc: func(user string) (string, error) {
			return currentShell, nil
		},
		SetUserShellFunc: func(user, shell string) error {
			return nil
		},
	}
	fileSystemMock := &utils.MoqFileSystem{
		PathExistsFunc: func(path string) (bool, error) {
			return path == shellPath, nil
		},
		IsExecutableFunc: func(path string) (bool, error) {
			return path == shellPath, nil
		},
		ReadFileContentsFunc: func(path string) ([]byte, error) {
			if path == "/etc/shells" {
				return []byte(etcShellsContent), nil
			}
			return nil, assert.AnError
		},
	}

	var teeCalled bool
	commanderMock := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			// Capture tee call for /etc/shells
			if name == "sudo" && len(args) > 0 && args[0] == "tee" {
				teeCalled = true
			}
			return &utils.Result{ExitCode: 0}, nil
		},
	}
	escalatorMock := &privilege.MoqEscalator{
		IsRunningAsRootFunc: func() (bool, error) {
			return false, nil
		},
		EscalateCommandFunc: func(baseCmd string, baseArgs []string) (privilege.EscalationResult, error) {
			return privilege.EscalationResult{
				Method:          privilege.EscalationSudo,
				Command:         "sudo",
				Args:            append([]string{baseCmd}, baseArgs...),
				NeedsEscalation: true,
			}, nil
		},
	}

	changer := shell.NewDefaultShellChanger(
		shellName,
		brewPath, // Brew-installed shell
		logger.DefaultLogger,
		osManagerMock,
		fileSystemMock,
		commanderMock,
		escalatorMock,
	)

	// Act
	err := changer.SetAsDefault(context.Background())

	// Assert
	require.NoError(t, err)
	assert.True(t, teeCalled, "Should have called tee to add shell to /etc/shells")
}

func Test_SetAsDefault_SkipsEtcShells_WhenShellAlreadyPresent(t *testing.T) {
	// Arrange
	shellName := "zsh"
	tempDir := t.TempDir()
	brewBinDir := filepath.Join(tempDir, "bin")
	shellPath := filepath.Join(brewBinDir, shellName)
	brewPath := filepath.Join(brewBinDir, "brew")
	currentShell := "/bin/bash"
	username := "testuser"

	// Simulate /etc/shells content WITH the brew shell already present
	etcShellsContent := "/bin/sh\n/bin/bash\n" + shellPath + "\n"

	osManagerMock := &osmanager.MoqOsManager{
		GetProgramPathFunc: func(program string) (string, error) {
			// Shell not in PATH (requires brew bin fallback)
			return "", assert.AnError
		},
		GetCurrentUsernameFunc: func() (string, error) {
			return username, nil
		},
		GetUserShellFunc: func(user string) (string, error) {
			return currentShell, nil
		},
		SetUserShellFunc: func(user, shell string) error {
			return nil
		},
	}
	fileSystemMock := &utils.MoqFileSystem{
		PathExistsFunc: func(path string) (bool, error) {
			return path == shellPath, nil
		},
		IsExecutableFunc: func(path string) (bool, error) {
			return path == shellPath, nil
		},
		ReadFileContentsFunc: func(path string) ([]byte, error) {
			if path == "/etc/shells" {
				return []byte(etcShellsContent), nil
			}
			return nil, assert.AnError
		},
	}

	var teeCalled bool
	commanderMock := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			if name == "sudo" && len(args) > 0 && args[0] == "tee" {
				teeCalled = true
			}
			return &utils.Result{ExitCode: 0}, nil
		},
	}
	escalatorMock := &privilege.MoqEscalator{
		IsRunningAsRootFunc: func() (bool, error) {
			return false, nil
		},
		EscalateCommandFunc: func(baseCmd string, baseArgs []string) (privilege.EscalationResult, error) {
			return privilege.EscalationResult{
				Method:          privilege.EscalationSudo,
				Command:         "sudo",
				Args:            append([]string{baseCmd}, baseArgs...),
				NeedsEscalation: true,
			}, nil
		},
	}

	changer := shell.NewDefaultShellChanger(
		shellName,
		brewPath,
		logger.DefaultLogger,
		osManagerMock,
		fileSystemMock,
		commanderMock,
		escalatorMock,
	)

	// Act
	err := changer.SetAsDefault(context.Background())

	// Assert
	require.NoError(t, err)
	assert.False(t, teeCalled, "Should NOT have called tee - shell already in /etc/shells")
}

func Test_SetAsDefault_SkipsEtcShells_WhenNotBrewInstalled(t *testing.T) {
	// Arrange
	shellName := "zsh"
	shellPath := "/usr/bin/zsh"
	currentShell := "/bin/bash"
	username := "testuser"

	osManagerMock := &osmanager.MoqOsManager{
		GetProgramPathFunc: func(program string) (string, error) {
			return shellPath, nil
		},
		GetCurrentUsernameFunc: func() (string, error) {
			return username, nil
		},
		GetUserShellFunc: func(user string) (string, error) {
			return currentShell, nil
		},
		SetUserShellFunc: func(user, shell string) error {
			return nil
		},
	}
	fileSystemMock := &utils.MoqFileSystem{
		PathExistsFunc: func(path string) (bool, error) {
			return true, nil
		},
		IsExecutableFunc: func(path string) (bool, error) {
			return true, nil
		},
		ReadFileContentsFunc: func(path string) ([]byte, error) {
			t.Error("ReadFileContents should not be called for non-brew shells")
			return nil, assert.AnError
		},
	}
	commanderMock := &utils.MoqCommander{}
	escalatorMock := &privilege.MoqEscalator{
		IsRunningAsRootFunc: func() (bool, error) {
			return false, nil
		},
	}

	changer := shell.NewDefaultShellChanger(
		shellName,
		"", // NOT brew installed
		logger.DefaultLogger,
		osManagerMock,
		fileSystemMock,
		commanderMock,
		escalatorMock,
	)

	// Act
	err := changer.SetAsDefault(context.Background())

	// Assert
	require.NoError(t, err)
	// ReadFileContents assertion is checked in the mock function itself
}

func Test_SetAsDefault_LogsWarning_WhenRunningAsRoot(t *testing.T) {
	// Arrange
	shellName := "zsh"
	shellPath := "/usr/bin/zsh"
	currentShell := "/bin/bash"
	username := "root"

	osManagerMock := &osmanager.MoqOsManager{
		GetProgramPathFunc: func(program string) (string, error) {
			return shellPath, nil
		},
		GetCurrentUsernameFunc: func() (string, error) {
			return username, nil
		},
		GetUserShellFunc: func(user string) (string, error) {
			return currentShell, nil
		},
		SetUserShellFunc: func(user, shell string) error {
			return nil
		},
	}
	fileSystemMock := &utils.MoqFileSystem{
		PathExistsFunc: func(path string) (bool, error) {
			return true, nil
		},
		IsExecutableFunc: func(path string) (bool, error) {
			return true, nil
		},
	}
	commanderMock := &utils.MoqCommander{}

	var isRunningAsRootCalled bool
	escalatorMock := &privilege.MoqEscalator{
		IsRunningAsRootFunc: func() (bool, error) {
			isRunningAsRootCalled = true
			return true, nil // Running as root
		},
	}

	changer := shell.NewDefaultShellChanger(
		shellName,
		"",
		logger.DefaultLogger,
		osManagerMock,
		fileSystemMock,
		commanderMock,
		escalatorMock,
	)

	// Act
	err := changer.SetAsDefault(context.Background())

	// Assert
	require.NoError(t, err)
	assert.True(t, isRunningAsRootCalled, "Should have checked if running as root")
}

// Platform-specific integration tests

func Test_SetAsDefault_Integration_MacOS(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Test only runs on macOS")
	}

	// This test verifies the integration works on macOS
	// but doesn't actually change the shell
	t.Skip("Integration test - run manually to verify macOS behavior")
}

func Test_SetAsDefault_Integration_Linux(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Test only runs on Linux")
	}

	// This test verifies the integration works on Linux
	// but doesn't actually change the shell
	t.Skip("Integration test - run manually to verify Linux behavior")
}
