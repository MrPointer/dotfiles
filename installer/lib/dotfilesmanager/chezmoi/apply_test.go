package chezmoi_test

import (
	"errors"
	"io"
	"os"
	"testing"

	"github.com/MrPointer/dotfiles/installer/lib/dotfilesmanager/chezmoi"
	"github.com/MrPointer/dotfiles/installer/lib/pkgmanager"
	"github.com/MrPointer/dotfiles/installer/utils"
	"github.com/MrPointer/dotfiles/installer/utils/httpclient"
	"github.com/MrPointer/dotfiles/installer/utils/logger"
	"github.com/MrPointer/dotfiles/installer/utils/osmanager"
	"github.com/stretchr/testify/require"
)

func Test_Apply_RemovesExistingCloneDirectory(t *testing.T) {
	mockFileSystem := &utils.MoqFileSystem{}
	mockFileSystem.RemovePathFunc = func(path string) error {
		require.Equal(t, "/home/user/.local/share/chezmoi", path)
		return nil
	}

	mockUserManager := &osmanager.MoqUserManager{}

	mockCommander := &utils.MoqCommander{}
	mockCommander.RunCommandFunc = func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
		return &utils.Result{ExitCode: 0}, nil
	}
	mockPackageManager := &pkgmanager.MoqPackageManager{}
	mockHTTPClient := &httpclient.MoqHTTPClient{}

	chezmoiConfig := chezmoi.DefaultChezmoiConfig("/home/user/.config/chezmoi/chezmoi.toml", "/home/user/.local/share/chezmoi")

	manager := chezmoi.NewChezmoiManager(logger.DefaultLogger, mockFileSystem, mockUserManager, mockCommander, mockPackageManager, mockHTTPClient, utils.DisplayModeProgress, chezmoiConfig)

	err := manager.Apply()

	require.NoError(t, err)
	require.Len(t, mockFileSystem.RemovePathCalls(), 1)
}

func Test_Apply_ReturnsError_WhenRemovePathFails(t *testing.T) {
	mockFileSystem := &utils.MoqFileSystem{}
	mockFileSystem.RemovePathFunc = func(path string) error {
		return errors.New("permission denied")
	}

	mockUserManager := &osmanager.MoqUserManager{}
	mockCommander := &utils.MoqCommander{}
	mockPackageManager := &pkgmanager.MoqPackageManager{}
	mockHTTPClient := &httpclient.MoqHTTPClient{}

	chezmoiConfig := chezmoi.DefaultChezmoiConfig("/home/user/.config/chezmoi/chezmoi.toml", "/home/user/.local/share/chezmoi")

	manager := chezmoi.NewChezmoiManager(logger.DefaultLogger, mockFileSystem, mockUserManager, mockCommander, mockPackageManager, mockHTTPClient, utils.DisplayModeProgress, chezmoiConfig)

	err := manager.Apply()

	require.Error(t, err)
	require.Contains(t, err.Error(), "permission denied")
	require.Len(t, mockFileSystem.RemovePathCalls(), 1)
}

func Test_Apply_RunsChezmoiInitApplyCommand_WithBasicArgs(t *testing.T) {
	mockFileSystem := &utils.MoqFileSystem{}
	mockFileSystem.RemovePathFunc = func(path string) error {
		return nil
	}

	mockUserManager := &osmanager.MoqUserManager{}

	mockCommander := &utils.MoqCommander{}
	mockCommander.RunCommandFunc = func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
		require.Equal(t, "chezmoi", name)
		require.Equal(t, []string{"init", "--apply", "MrPointer", "--config", "/home/user/.config/chezmoi/chezmoi.toml"}, args)
		return &utils.Result{ExitCode: 0}, nil
	}

	mockPackageManager := &pkgmanager.MoqPackageManager{}
	mockHTTPClient := &httpclient.MoqHTTPClient{}

	chezmoiConfig := chezmoi.DefaultChezmoiConfig("/home/user/.config/chezmoi/chezmoi.toml", "")

	manager := chezmoi.NewChezmoiManager(logger.DefaultLogger, mockFileSystem, mockUserManager, mockCommander, mockPackageManager, mockHTTPClient, utils.DisplayModeProgress, chezmoiConfig)

	err := manager.Apply()

	require.NoError(t, err)
	require.Len(t, mockCommander.RunCommandCalls(), 1)
}

func Test_Apply_RunsChezmoiInitApplyCommand_WithSourceDir(t *testing.T) {
	mockFileSystem := &utils.MoqFileSystem{}
	mockFileSystem.RemovePathFunc = func(path string) error {
		return nil
	}

	mockUserManager := &osmanager.MoqUserManager{}

	mockCommander := &utils.MoqCommander{}
	mockCommander.RunCommandFunc = func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
		require.Equal(t, "chezmoi", name)
		expectedArgs := []string{"init", "--apply", "--source", "/home/user/.local/share/chezmoi", "MrPointer", "--config", "/home/user/.config/chezmoi/chezmoi.toml"}
		require.Equal(t, expectedArgs, args)
		return &utils.Result{ExitCode: 0}, nil
	}

	mockPackageManager := &pkgmanager.MoqPackageManager{}
	mockHTTPClient := &httpclient.MoqHTTPClient{}

	chezmoiConfig := chezmoi.DefaultChezmoiConfig("/home/user/.config/chezmoi/chezmoi.toml", "/home/user/.local/share/chezmoi")

	manager := chezmoi.NewChezmoiManager(logger.DefaultLogger, mockFileSystem, mockUserManager, mockCommander, mockPackageManager, mockHTTPClient, utils.DisplayModeProgress, chezmoiConfig)

	err := manager.Apply()

	require.NoError(t, err)
	require.Len(t, mockCommander.RunCommandCalls(), 1)
}

func Test_Apply_RunsChezmoiInitApplyCommand_WithSSHCloningPreference(t *testing.T) {
	mockFileSystem := &utils.MoqFileSystem{}
	mockFileSystem.RemovePathFunc = func(path string) error {
		return nil
	}

	mockUserManager := &osmanager.MoqUserManager{}

	mockCommander := &utils.MoqCommander{}
	mockCommander.RunCommandFunc = func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
		require.Equal(t, "chezmoi", name)
		expectedArgs := []string{"init", "--apply", "--source", "/home/user/.local/share/chezmoi", "--ssh", "testuser", "--config", "/home/user/.config/chezmoi/chezmoi.toml"}
		require.Equal(t, expectedArgs, args)
		return &utils.Result{ExitCode: 0}, nil
	}

	mockPackageManager := &pkgmanager.MoqPackageManager{}
	mockHTTPClient := &httpclient.MoqHTTPClient{}

	chezmoiConfig := chezmoi.NewChezmoiConfig(
		"/home/user/.config/chezmoi",
		"/home/user/.config/chezmoi/chezmoi.toml",
		"/home/user/.local/share/chezmoi",
		"testuser",
		true,
	)

	manager := chezmoi.NewChezmoiManager(logger.DefaultLogger, mockFileSystem, mockUserManager, mockCommander, mockPackageManager, mockHTTPClient, utils.DisplayModeProgress, chezmoiConfig)

	err := manager.Apply()

	require.NoError(t, err)
	require.Len(t, mockCommander.RunCommandCalls(), 1)
}

func Test_Apply_RunsChezmoiInitApplyCommand_WithCustomUsername(t *testing.T) {
	mockFileSystem := &utils.MoqFileSystem{}
	mockFileSystem.RemovePathFunc = func(path string) error {
		return nil
	}

	mockUserManager := &osmanager.MoqUserManager{}

	mockCommander := &utils.MoqCommander{}
	mockCommander.RunCommandFunc = func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
		require.Equal(t, "chezmoi", name)
		require.Equal(t, []string{"init", "--apply", "customuser", "--config", "/home/user/.config/chezmoi/chezmoi.toml"}, args)
		return &utils.Result{ExitCode: 0}, nil
	}

	mockPackageManager := &pkgmanager.MoqPackageManager{}
	mockHTTPClient := &httpclient.MoqHTTPClient{}

	chezmoiConfig := chezmoi.NewChezmoiConfig(
		"/home/user/.config/chezmoi",
		"/home/user/.config/chezmoi/chezmoi.toml",
		"",
		"customuser",
		false,
	)

	manager := chezmoi.NewChezmoiManager(logger.DefaultLogger, mockFileSystem, mockUserManager, mockCommander, mockPackageManager, mockHTTPClient, utils.DisplayModeProgress, chezmoiConfig)

	err := manager.Apply()

	require.NoError(t, err)
	require.Len(t, mockCommander.RunCommandCalls(), 1)
}

func Test_Apply_ReturnsError_WhenCommandExecutionFails(t *testing.T) {
	mockFileSystem := &utils.MoqFileSystem{}
	mockFileSystem.RemovePathFunc = func(path string) error {
		return nil
	}

	mockUserManager := &osmanager.MoqUserManager{}

	mockCommander := &utils.MoqCommander{}
	mockCommander.RunCommandFunc = func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
		return nil, errors.New("command not found")
	}

	mockPackageManager := &pkgmanager.MoqPackageManager{}
	mockHTTPClient := &httpclient.MoqHTTPClient{}

	chezmoiConfig := chezmoi.DefaultChezmoiConfig("/home/user/.config/chezmoi/chezmoi.toml", "/home/user/.local/share/chezmoi")

	manager := chezmoi.NewChezmoiManager(logger.DefaultLogger, mockFileSystem, mockUserManager, mockCommander, mockPackageManager, mockHTTPClient, utils.DisplayModeProgress, chezmoiConfig)

	err := manager.Apply()

	require.Error(t, err)
	require.Contains(t, err.Error(), "command not found")
}

func Test_Apply_ReturnsError_WhenCommandExitsWithNonZeroCode(t *testing.T) {
	mockFileSystem := &utils.MoqFileSystem{}
	mockFileSystem.RemovePathFunc = func(path string) error {
		return nil
	}

	mockUserManager := &osmanager.MoqUserManager{}

	mockCommander := &utils.MoqCommander{}
	mockCommander.RunCommandFunc = func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
		return &utils.Result{
			ExitCode: 1,
			Stderr:   []byte("initialization failed"),
		}, nil
	}

	mockPackageManager := &pkgmanager.MoqPackageManager{}
	mockHTTPClient := &httpclient.MoqHTTPClient{}

	chezmoiConfig := chezmoi.DefaultChezmoiConfig("/home/user/.config/chezmoi/chezmoi.toml", "/home/user/.local/share/chezmoi")

	manager := chezmoi.NewChezmoiManager(logger.DefaultLogger, mockFileSystem, mockUserManager, mockCommander, mockPackageManager, mockHTTPClient, utils.DisplayModeProgress, chezmoiConfig)

	err := manager.Apply()

	require.Error(t, err)
	require.Contains(t, err.Error(), "chezmoi init failed with exit code 1")
	require.Contains(t, err.Error(), "initialization failed")
}

func Test_Apply_SucceedsWithAllParametersCombined(t *testing.T) {
	mockFileSystem := &utils.MoqFileSystem{}
	mockFileSystem.RemovePathFunc = func(path string) error {
		require.Equal(t, "/custom/clone/dir", path)
		return nil
	}

	mockUserManager := &osmanager.MoqUserManager{}

	mockCommander := &utils.MoqCommander{}
	mockCommander.RunCommandFunc = func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
		require.Equal(t, "chezmoi", name)
		expectedArgs := []string{"init", "--apply", "--source", "/custom/clone/dir", "--ssh", "testuser123", "--config", "/home/user/.config/chezmoi/chezmoi.toml"}
		require.Equal(t, expectedArgs, args)
		return &utils.Result{ExitCode: 0}, nil
	}

	mockPackageManager := &pkgmanager.MoqPackageManager{}
	mockHTTPClient := &httpclient.MoqHTTPClient{}

	chezmoiConfig := chezmoi.NewChezmoiConfig(
		"/home/user/.config/chezmoi",
		"/home/user/.config/chezmoi/chezmoi.toml",
		"/custom/clone/dir",
		"testuser123",
		true,
	)

	manager := chezmoi.NewChezmoiManager(logger.DefaultLogger, mockFileSystem, mockUserManager, mockCommander, mockPackageManager, mockHTTPClient, utils.DisplayModeProgress, chezmoiConfig)

	err := manager.Apply()

	require.NoError(t, err)
	require.Len(t, mockFileSystem.RemovePathCalls(), 1)
	require.Len(t, mockCommander.RunCommandCalls(), 1)
}

func Test_Apply_DiscardsOutput_WhenDisplayModeIsNotPassthrough(t *testing.T) {
	mockFileSystem := &utils.MoqFileSystem{
		RemovePathFunc: func(path string) error {
			require.Equal(t, "/custom/clone/dir", path)
			return nil
		},
	}

	mockUserManager := &osmanager.MoqUserManager{}

	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(command string, args []string, options ...utils.Option) (*utils.Result, error) {
			cmdOptions := utils.Options{}

			// Apply all provided options
			for _, opt := range options {
				opt(&cmdOptions)
			}

			require.Equal(t, io.Discard, cmdOptions.Stdout)
			require.Equal(t, io.Discard, cmdOptions.Stderr)
			return &utils.Result{}, nil
		},
	}

	mockPackageManager := &pkgmanager.MoqPackageManager{}
	mockHTTPClient := &httpclient.MoqHTTPClient{}

	chezmoiConfig := chezmoi.NewChezmoiConfig(
		"/home/user/.config/chezmoi",
		"/home/user/.config/chezmoi/chezmoi.toml",
		"/custom/clone/dir",
		"testuser123",
		true,
	)

	manager := chezmoi.NewChezmoiManager(logger.DefaultLogger, mockFileSystem, mockUserManager, mockCommander, mockPackageManager, mockHTTPClient, utils.DisplayModeProgress, chezmoiConfig)

	err := manager.Apply()

	require.NoError(t, err)
}

func Test_Apply_DoesNotDiscardOutput_WhenDisplayModeIsPassthrough(t *testing.T) {
	mockFileSystem := &utils.MoqFileSystem{
		RemovePathFunc: func(path string) error {
			require.Equal(t, "/custom/clone/dir", path)
			return nil
		},
	}
	mockUserManager := &osmanager.MoqUserManager{}

	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(command string, args []string, options ...utils.Option) (*utils.Result, error) {
			cmdOptions := utils.Options{
				Stdout: os.Stdout,
				Stderr: os.Stderr,
			}

			// Apply all provided options
			for _, opt := range options {
				opt(&cmdOptions)
			}

			require.Equal(t, os.Stdout, cmdOptions.Stdout)
			require.Equal(t, os.Stderr, cmdOptions.Stderr)
			return &utils.Result{}, nil
		},
	}

	mockPackageManager := &pkgmanager.MoqPackageManager{}
	mockHTTPClient := &httpclient.MoqHTTPClient{}

	chezmoiConfig := chezmoi.NewChezmoiConfig(
		"/home/user/.config/chezmoi",
		"/home/user/.config/chezmoi/chezmoi.toml",
		"/custom/clone/dir",
		"testuser123",
		true,
	)

	manager := chezmoi.NewChezmoiManager(logger.DefaultLogger, mockFileSystem, mockUserManager, mockCommander, mockPackageManager, mockHTTPClient, utils.DisplayModePassthrough, chezmoiConfig)

	err := manager.Apply()

	require.NoError(t, err)
}
