package chezmoi_test

import (
	"errors"
	"testing"

	"github.com/MrPointer/dotfiles/installer/lib/dotfilesmanager/chezmoi"
	"github.com/MrPointer/dotfiles/installer/lib/pkgmanager"
	"github.com/MrPointer/dotfiles/installer/utils"
	"github.com/MrPointer/dotfiles/installer/utils/httpclient"
	"github.com/MrPointer/dotfiles/installer/utils/logger"
	"github.com/MrPointer/dotfiles/installer/utils/osmanager"
	"github.com/stretchr/testify/require"
)

func Test_NewChezmoiManager_ReturnsValidInstance(t *testing.T) {
	mockFileSystem := &utils.MoqFileSystem{}
	mockUserManager := &osmanager.MoqUserManager{}
	mockCommander := &utils.MoqCommander{}
	mockPkgManager := &pkgmanager.MoqPackageManager{}
	mockHttpClient := &httpclient.MoqHTTPClient{}

	configFilePath := "/home/user/.config/chezmoi.toml"

	initializer := chezmoi.NewChezmoiManager(logger.DefaultLogger, mockFileSystem, mockUserManager, mockCommander, mockPkgManager, mockHttpClient, utils.DisplayModeProgress, chezmoi.DefaultChezmoiConfig(configFilePath, ""))

	require.NotNil(t, initializer)
}

func Test_TryNewDefaultChezmoiManager_ReturnsValidInstance_WhenUserConfigDirAndHomeDirAreAvailable(t *testing.T) {
	mockFileSystem := &utils.MoqFileSystem{}

	mockUserManager := &osmanager.MoqUserManager{
		GetChezmoiConfigHomeFunc: func() (string, error) {
			return "/home/user/.config", nil
		},
		GetHomeDirFunc: func() (string, error) {
			return "/home/user", nil
		},
	}

	mockCommander := &utils.MoqCommander{}
	mockPackageManager := &pkgmanager.MoqPackageManager{}
	mockHTTPClient := &httpclient.MoqHTTPClient{}

	initializer, err := chezmoi.TryStandardChezmoiManagerWithDefaults(logger.DefaultLogger, mockFileSystem, mockUserManager, mockCommander, mockPackageManager, mockHTTPClient, utils.DisplayModeProgress)

	require.NoError(t, err)
	require.NotNil(t, initializer)
}

func Test_TryNewDefaultChezmoiManager_ReturnsError_WhenUserConfigDirIsUnavailable(t *testing.T) {
	mockFileSystem := &utils.MoqFileSystem{}

	mockUserManager := &osmanager.MoqUserManager{
		GetChezmoiConfigHomeFunc: func() (string, error) {
			return "", errors.New("failed to get user config directory")
		},
		GetHomeDirFunc: func() (string, error) {
			return "/home/user", nil
		},
	}

	mockCommander := &utils.MoqCommander{}
	mockPackageManager := &pkgmanager.MoqPackageManager{}
	mockHTTPClient := &httpclient.MoqHTTPClient{}

	initializer, err := chezmoi.TryStandardChezmoiManagerWithDefaults(logger.DefaultLogger, mockFileSystem, mockUserManager, mockCommander, mockPackageManager, mockHTTPClient, utils.DisplayModeProgress)

	require.Error(t, err)
	require.Nil(t, initializer)
}

func Test_TryNewDefaultChezmoiManager_ReturnsError_WhenUserHomeDirIsUnavailable(t *testing.T) {
	mockFileSystem := &utils.MoqFileSystem{}

	mockUserManager := &osmanager.MoqUserManager{
		GetChezmoiConfigHomeFunc: func() (string, error) {
			return "/home/user/.config", nil
		},
		GetHomeDirFunc: func() (string, error) {
			return "", errors.New("failed to get user home directory")
		},
	}

	mockCommander := &utils.MoqCommander{}
	mockPackageManager := &pkgmanager.MoqPackageManager{}
	mockHTTPClient := &httpclient.MoqHTTPClient{}

	initializer, err := chezmoi.TryStandardChezmoiManagerWithDefaults(logger.DefaultLogger, mockFileSystem, mockUserManager, mockCommander, mockPackageManager, mockHTTPClient, utils.DisplayModeProgress)

	require.Error(t, err)
	require.Nil(t, initializer)
}
