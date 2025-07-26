package chezmoi

import (
	"fmt"
	"path/filepath"

	"github.com/MrPointer/dotfiles/installer/lib/dotfilesmanager"
	"github.com/MrPointer/dotfiles/installer/lib/pkgmanager"
	"github.com/MrPointer/dotfiles/installer/utils"
	"github.com/MrPointer/dotfiles/installer/utils/httpclient"
	"github.com/MrPointer/dotfiles/installer/utils/osmanager"
)

const DefaultGitHubUsername = "MrPointer"

type ChezmoiConfig struct {
	chezmoiConfigDir      string
	chezmoiConfigFilePath string
	chezmoiCloneDir       string
	githubUsername        string
	cloneViaSSH           bool
}

func NewChezmoiConfig(configDir, configFilePath, cloneDir, githubUsername string, cloneViaSSH bool) ChezmoiConfig {
	return ChezmoiConfig{
		chezmoiConfigDir:      configDir,
		chezmoiConfigFilePath: configFilePath,
		chezmoiCloneDir:       cloneDir,
		githubUsername:        githubUsername,
		cloneViaSSH:           cloneViaSSH,
	}
}

func DefaultChezmoiConfig(chezmoiConfigFilePath string, chezmoiCloneDir string) ChezmoiConfig {
	return ChezmoiConfig{
		chezmoiConfigDir:      filepath.Dir(chezmoiConfigFilePath),
		chezmoiConfigFilePath: chezmoiConfigFilePath,
		chezmoiCloneDir:       chezmoiCloneDir,
		githubUsername:        "MrPointer",
		cloneViaSSH:           false,
	}
}

type ChezmoiManager struct {
	chezmoiConfig ChezmoiConfig
	filesystem    utils.FileSystem
	usermanager   osmanager.UserManager
	commander     utils.Commander
	pkgManager    pkgmanager.PackageManager
	httpClient    httpclient.HTTPClient
}

var _ dotfilesmanager.DotfilesDataInitializer = (*ChezmoiManager)(nil)

func NewChezmoiManager(filesystem utils.FileSystem, userManager osmanager.UserManager, commander utils.Commander, pkgManager pkgmanager.PackageManager, httpClient httpclient.HTTPClient, chezmoiConfig ChezmoiConfig) *ChezmoiManager {
	return &ChezmoiManager{
		chezmoiConfig: chezmoiConfig,
		filesystem:    filesystem,
		usermanager:   userManager,
		commander:     commander,
		pkgManager:    pkgManager,
		httpClient:    httpClient,
	}
}

func TryStandardChezmoiManager(filesystem utils.FileSystem, userManager osmanager.UserManager, commander utils.Commander, pkgManager pkgmanager.PackageManager, httpClient httpclient.HTTPClient, githubUsername string, cloneViaSSH bool) (*ChezmoiManager, error) {
	userConfigDir, err := userManager.GetConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user config directory: %w", err)
	}

	chezmoiConfigDir := fmt.Sprintf("%s/chezmoi", userConfigDir)
	chezmoiConfigFilePath := fmt.Sprintf("%s/chezmoi.toml", chezmoiConfigDir)

	userHomeDir, err := userManager.GetHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	chezmoiCloneDir := fmt.Sprintf("%s/.local/share/chezmoi", userHomeDir)

	return NewChezmoiManager(
		filesystem,
		userManager,
		commander,
		pkgManager,
		httpClient,
		NewChezmoiConfig(
			chezmoiConfigDir,
			chezmoiConfigFilePath,
			chezmoiCloneDir,
			githubUsername,
			cloneViaSSH,
		),
	), nil
}

func TryStandardChezmoiManagerWithDefaults(filesystem utils.FileSystem, userManager osmanager.UserManager, commander utils.Commander, pkgManager pkgmanager.PackageManager, httpClient httpclient.HTTPClient) (*ChezmoiManager, error) {
	return TryStandardChezmoiManager(filesystem, userManager, commander, pkgManager, httpClient, DefaultGitHubUsername, false)
}
