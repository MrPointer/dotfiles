package chezmoi_test

import (
	"errors"
	"testing"

	"github.com/MrPointer/dotfiles/installer/lib/dotfilesmanager/chezmoi"
	"github.com/MrPointer/dotfiles/installer/lib/pkgmanager"
	"github.com/MrPointer/dotfiles/installer/utils"
	"github.com/MrPointer/dotfiles/installer/utils/httpclient"
	"github.com/MrPointer/dotfiles/installer/utils/osmanager"
	"github.com/stretchr/testify/require"
)

func Test_Apply_RemovesExistingCloneDirectory(t *testing.T) {
	mockFileSystem := &utils.MoqFileSystem{}
	mockUserManager := &osmanager.MoqUserManager{}
	mockCommander := &utils.MoqCommander{}
	mockPackageManager := &pkgmanager.MoqPackageManager{}
	mockHTTPClient := &httpclient.MoqHTTPClient{}

	chezmoiConfig := chezmoi.DefaultChezmoiConfig("/home/user/.config/chezmoi/chezmoi.toml", "/home/user/.local/share/chezmoi")
	manager := chezmoi.NewChezmoiManager(mockFileSystem, mockUserManager, mockCommander, mockPackageManager, mockHTTPClient, chezmoiConfig)

	mockFileSystem.RemovePathFunc = func(path string) error {
		require.Equal(t, "/home/user/.local/share/chezmoi", path)
		return nil
	}

	mockCommander.RunCommandFunc = func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
		return &utils.Result{ExitCode: 0}, nil
	}

	err := manager.Apply()

	require.NoError(t, err)
	require.Len(t, mockFileSystem.RemovePathCalls(), 1)
}

func Test_Apply_ReturnsError_WhenRemovePathFails(t *testing.T) {
	mockFileSystem := &utils.MoqFileSystem{}
	mockUserManager := &osmanager.MoqUserManager{}
	mockCommander := &utils.MoqCommander{}
	mockPackageManager := &pkgmanager.MoqPackageManager{}
	mockHTTPClient := &httpclient.MoqHTTPClient{}

	chezmoiConfig := chezmoi.DefaultChezmoiConfig("/home/user/.config/chezmoi/chezmoi.toml", "/home/user/.local/share/chezmoi")
	manager := chezmoi.NewChezmoiManager(mockFileSystem, mockUserManager, mockCommander, mockPackageManager, mockHTTPClient, chezmoiConfig)

	mockFileSystem.RemovePathFunc = func(path string) error {
		return errors.New("permission denied")
	}

	err := manager.Apply()

	require.Error(t, err)
	require.Contains(t, err.Error(), "permission denied")
	require.Len(t, mockFileSystem.RemovePathCalls(), 1)
}

func Test_Apply_RunsChezmoiInitApplyCommand_WithBasicArgs(t *testing.T) {
	mockFileSystem := &utils.MoqFileSystem{}
	mockUserManager := &osmanager.MoqUserManager{}
	mockCommander := &utils.MoqCommander{}
	mockPackageManager := &pkgmanager.MoqPackageManager{}
	mockHTTPClient := &httpclient.MoqHTTPClient{}

	chezmoiConfig := chezmoi.DefaultChezmoiConfig("/home/user/.config/chezmoi/chezmoi.toml", "")
	manager := chezmoi.NewChezmoiManager(mockFileSystem, mockUserManager, mockCommander, mockPackageManager, mockHTTPClient, chezmoiConfig)

	mockFileSystem.RemovePathFunc = func(path string) error {
		return nil
	}

	mockCommander.RunCommandFunc = func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
		require.Equal(t, "chezmoi", name)
		require.Equal(t, []string{"init", "--apply", "MrPointer"}, args)
		return &utils.Result{ExitCode: 0}, nil
	}

	err := manager.Apply()

	require.NoError(t, err)
	require.Len(t, mockCommander.RunCommandCalls(), 1)
}

func Test_Apply_RunsChezmoiInitApplyCommand_WithSourceDir(t *testing.T) {
	mockFileSystem := &utils.MoqFileSystem{}
	mockUserManager := &osmanager.MoqUserManager{}
	mockCommander := &utils.MoqCommander{}
	mockPackageManager := &pkgmanager.MoqPackageManager{}
	mockHTTPClient := &httpclient.MoqHTTPClient{}

	chezmoiConfig := chezmoi.DefaultChezmoiConfig("/home/user/.config/chezmoi/chezmoi.toml", "/home/user/.local/share/chezmoi")
	manager := chezmoi.NewChezmoiManager(mockFileSystem, mockUserManager, mockCommander, mockPackageManager, mockHTTPClient, chezmoiConfig)

	mockFileSystem.RemovePathFunc = func(path string) error {
		return nil
	}

	mockCommander.RunCommandFunc = func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
		require.Equal(t, "chezmoi", name)
		expectedArgs := []string{"init", "--apply", "--source", "/home/user/.local/share/chezmoi", "MrPointer"}
		require.Equal(t, expectedArgs, args)
		return &utils.Result{ExitCode: 0}, nil
	}

	err := manager.Apply()

	require.NoError(t, err)
	require.Len(t, mockCommander.RunCommandCalls(), 1)
}

func Test_Apply_RunsChezmoiInitApplyCommand_WithSSHCloningPreference(t *testing.T) {
	mockFileSystem := &utils.MoqFileSystem{}
	mockUserManager := &osmanager.MoqUserManager{}
	mockCommander := &utils.MoqCommander{}
	mockPackageManager := &pkgmanager.MoqPackageManager{}
	mockHTTPClient := &httpclient.MoqHTTPClient{}

	chezmoiConfig := chezmoi.NewChezmoiConfig(
		"/home/user/.config/chezmoi",
		"/home/user/.config/chezmoi/chezmoi.toml",
		"/home/user/.local/share/chezmoi",
		"testuser",
		true,
	)
	manager := chezmoi.NewChezmoiManager(mockFileSystem, mockUserManager, mockCommander, mockPackageManager, mockHTTPClient, chezmoiConfig)

	mockFileSystem.RemovePathFunc = func(path string) error {
		return nil
	}

	mockCommander.RunCommandFunc = func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
		require.Equal(t, "chezmoi", name)
		expectedArgs := []string{"init", "--apply", "--source", "/home/user/.local/share/chezmoi", "--ssh", "testuser"}
		require.Equal(t, expectedArgs, args)
		return &utils.Result{ExitCode: 0}, nil
	}

	err := manager.Apply()

	require.NoError(t, err)
	require.Len(t, mockCommander.RunCommandCalls(), 1)
}

func Test_Apply_RunsChezmoiInitApplyCommand_WithCustomUsername(t *testing.T) {
	mockFileSystem := &utils.MoqFileSystem{}
	mockUserManager := &osmanager.MoqUserManager{}
	mockCommander := &utils.MoqCommander{}
	mockPackageManager := &pkgmanager.MoqPackageManager{}
	mockHTTPClient := &httpclient.MoqHTTPClient{}

	chezmoiConfig := chezmoi.NewChezmoiConfig(
		"/home/user/.config/chezmoi",
		"/home/user/.config/chezmoi/chezmoi.toml",
		"",
		"customuser",
		false,
	)
	manager := chezmoi.NewChezmoiManager(mockFileSystem, mockUserManager, mockCommander, mockPackageManager, mockHTTPClient, chezmoiConfig)

	mockFileSystem.RemovePathFunc = func(path string) error {
		return nil
	}

	mockCommander.RunCommandFunc = func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
		require.Equal(t, "chezmoi", name)
		require.Equal(t, []string{"init", "--apply", "customuser"}, args)
		return &utils.Result{ExitCode: 0}, nil
	}

	err := manager.Apply()

	require.NoError(t, err)
	require.Len(t, mockCommander.RunCommandCalls(), 1)
}

func Test_Apply_ReturnsError_WhenCommandExecutionFails(t *testing.T) {
	mockFileSystem := &utils.MoqFileSystem{}
	mockUserManager := &osmanager.MoqUserManager{}
	mockCommander := &utils.MoqCommander{}
	mockPackageManager := &pkgmanager.MoqPackageManager{}
	mockHTTPClient := &httpclient.MoqHTTPClient{}

	chezmoiConfig := chezmoi.DefaultChezmoiConfig("/home/user/.config/chezmoi/chezmoi.toml", "/home/user/.local/share/chezmoi")
	manager := chezmoi.NewChezmoiManager(mockFileSystem, mockUserManager, mockCommander, mockPackageManager, mockHTTPClient, chezmoiConfig)

	mockFileSystem.RemovePathFunc = func(path string) error {
		return nil
	}

	mockCommander.RunCommandFunc = func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
		return nil, errors.New("command not found")
	}

	err := manager.Apply()

	require.Error(t, err)
	require.Contains(t, err.Error(), "command not found")
}

func Test_Apply_ReturnsError_WhenCommandExitsWithNonZeroCode(t *testing.T) {
	mockFileSystem := &utils.MoqFileSystem{}
	mockUserManager := &osmanager.MoqUserManager{}
	mockCommander := &utils.MoqCommander{}
	mockPackageManager := &pkgmanager.MoqPackageManager{}
	mockHTTPClient := &httpclient.MoqHTTPClient{}

	chezmoiConfig := chezmoi.DefaultChezmoiConfig("/home/user/.config/chezmoi/chezmoi.toml", "/home/user/.local/share/chezmoi")
	manager := chezmoi.NewChezmoiManager(mockFileSystem, mockUserManager, mockCommander, mockPackageManager, mockHTTPClient, chezmoiConfig)

	mockFileSystem.RemovePathFunc = func(path string) error {
		return nil
	}

	mockCommander.RunCommandFunc = func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
		return &utils.Result{
			ExitCode: 1,
			Stderr:   []byte("initialization failed"),
		}, nil
	}

	err := manager.Apply()

	require.Error(t, err)
	require.Contains(t, err.Error(), "chezmoi init failed with exit code 1")
	require.Contains(t, err.Error(), "initialization failed")
}

func Test_Apply_PassesStdoutAndStderrOptions(t *testing.T) {
	mockFileSystem := &utils.MoqFileSystem{}
	mockUserManager := &osmanager.MoqUserManager{}
	mockCommander := &utils.MoqCommander{}
	mockPackageManager := &pkgmanager.MoqPackageManager{}
	mockHTTPClient := &httpclient.MoqHTTPClient{}

	chezmoiConfig := chezmoi.DefaultChezmoiConfig("/home/user/.config/chezmoi/chezmoi.toml", "/home/user/.local/share/chezmoi")
	manager := chezmoi.NewChezmoiManager(mockFileSystem, mockUserManager, mockCommander, mockPackageManager, mockHTTPClient, chezmoiConfig)

	mockFileSystem.RemovePathFunc = func(path string) error {
		return nil
	}

	mockCommander.RunCommandFunc = func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
		require.Equal(t, "chezmoi", name)
		require.Len(t, opts, 2)
		return &utils.Result{ExitCode: 0}, nil
	}

	err := manager.Apply()

	require.NoError(t, err)
	require.Len(t, mockCommander.RunCommandCalls(), 1)
}

func Test_Apply_SucceedsWithAllParametersCombined(t *testing.T) {
	mockFileSystem := &utils.MoqFileSystem{}
	mockUserManager := &osmanager.MoqUserManager{}
	mockCommander := &utils.MoqCommander{}
	mockPackageManager := &pkgmanager.MoqPackageManager{}
	mockHTTPClient := &httpclient.MoqHTTPClient{}

	chezmoiConfig := chezmoi.NewChezmoiConfig(
		"/home/user/.config/chezmoi",
		"/home/user/.config/chezmoi/chezmoi.toml",
		"/custom/clone/dir",
		"testuser123",
		true,
	)
	manager := chezmoi.NewChezmoiManager(mockFileSystem, mockUserManager, mockCommander, mockPackageManager, mockHTTPClient, chezmoiConfig)

	mockFileSystem.RemovePathFunc = func(path string) error {
		require.Equal(t, "/custom/clone/dir", path)
		return nil
	}

	mockCommander.RunCommandFunc = func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
		require.Equal(t, "chezmoi", name)
		expectedArgs := []string{"init", "--apply", "--source", "/custom/clone/dir", "--ssh", "testuser123"}
		require.Equal(t, expectedArgs, args)
		require.Len(t, opts, 2)
		return &utils.Result{ExitCode: 0}, nil
	}

	err := manager.Apply()

	require.NoError(t, err)
	require.Len(t, mockFileSystem.RemovePathCalls(), 1)
	require.Len(t, mockCommander.RunCommandCalls(), 1)
}
