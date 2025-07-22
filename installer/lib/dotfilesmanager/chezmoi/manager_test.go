package chezmoi_test

import (
	"errors"
	"testing"

	"github.com/MrPointer/dotfiles/installer/lib/dotfilesmanager/chezmoi"
	"github.com/MrPointer/dotfiles/installer/utils"
	"github.com/MrPointer/dotfiles/installer/utils/osmanager"
	"github.com/stretchr/testify/require"
)

func Test_NewChezmoiManager_ReturnsValidInstance(t *testing.T) {
	mockFileSystem := &utils.MoqFileSystem{}
	configFilePath := "/home/user/.config/chezmoi.toml"

	initializer := chezmoi.NewChezmoiManager(configFilePath, "", mockFileSystem)

	require.NotNil(t, initializer)
}

func Test_TryNewDefaultChezmoiManager_ReturnsValidInstance_WhenUserConfigDirAndHomeDirAreAvailable(t *testing.T) {
	mockFileSystem := &utils.MoqFileSystem{}

	userManager := &osmanager.MoqUserManager{}
	userManager.GetConfigDirFunc = func() (string, error) {
		return "/home/user/.config", nil
	}
	userManager.GetHomeDirFunc = func() (string, error) {
		return "/home/user", nil
	}

	initializer, err := chezmoi.TryNewDefaultChezmoiManager(mockFileSystem, userManager)

	require.NoError(t, err)
	require.NotNil(t, initializer)
}

func Test_TryNewDefaultChezmoiManager_ReturnsError_WhenUserConfigDirIsUnavailable(t *testing.T) {
	mockFileSystem := &utils.MoqFileSystem{}

	userManager := &osmanager.MoqUserManager{}
	userManager.GetConfigDirFunc = func() (string, error) {
		return "", errors.New("failed to get user config directory")
	}
	userManager.GetHomeDirFunc = func() (string, error) {
		return "/home/user", nil
	}

	initializer, err := chezmoi.TryNewDefaultChezmoiManager(mockFileSystem, userManager)

	require.Error(t, err)
	require.Nil(t, initializer)
}

func Test_TryNewDefaultChezmoiManager_ReturnsError_WhenUserHomeDirIsUnavailable(t *testing.T) {
	mockFileSystem := &utils.MoqFileSystem{}

	userManager := &osmanager.MoqUserManager{}
	userManager.GetConfigDirFunc = func() (string, error) {
		return "/home/user/.config", nil
	}
	userManager.GetHomeDirFunc = func() (string, error) {
		return "", errors.New("failed to get user home directory")
	}

	initializer, err := chezmoi.TryNewDefaultChezmoiManager(mockFileSystem, userManager)

	require.Error(t, err)
	require.Nil(t, initializer)
}
