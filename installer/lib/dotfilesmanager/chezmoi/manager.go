package chezmoi

import (
	"fmt"
	"path/filepath"

	"github.com/MrPointer/dotfiles/installer/lib/dotfilesmanager"
	"github.com/MrPointer/dotfiles/installer/lib/pkgmanager"
	"github.com/MrPointer/dotfiles/installer/utils"
	"github.com/MrPointer/dotfiles/installer/utils/httpclient"
	"github.com/MrPointer/dotfiles/installer/utils/logger"
	"github.com/MrPointer/dotfiles/installer/utils/osmanager"
)

const DefaultGitHubUsername = "MrPointer"

type ChezmoiConfig struct {
	chezmoiConfigDir      string
	chezmoiConfigFilePath string
	chezmoiCloneDir       string
	githubUsername        string
	cloneViaSSH           bool
	branch                string
}

func NewChezmoiConfig(configDir, configFilePath, cloneDir, githubUsername string, cloneViaSSH bool, branch string) ChezmoiConfig {
	return ChezmoiConfig{
		chezmoiConfigDir:      configDir,
		chezmoiConfigFilePath: configFilePath,
		chezmoiCloneDir:       cloneDir,
		githubUsername:        githubUsername,
		cloneViaSSH:           cloneViaSSH,
		branch:                branch,
	}
}

func DefaultChezmoiConfig(chezmoiConfigFilePath string, chezmoiCloneDir string) ChezmoiConfig {
	return ChezmoiConfig{
		chezmoiConfigDir:      filepath.Dir(chezmoiConfigFilePath),
		chezmoiConfigFilePath: chezmoiConfigFilePath,
		chezmoiCloneDir:       chezmoiCloneDir,
		githubUsername:        "MrPointer",
		cloneViaSSH:           false,
		branch:                "",
	}
}

type ChezmoiManager struct {
	chezmoiConfig ChezmoiConfig
	logger        logger.Logger
	filesystem    utils.FileSystem
	usermanager   osmanager.UserManager
	commander     utils.Commander
	pkgManager    pkgmanager.PackageManager
	httpClient    httpclient.HTTPClient
	displayMode   utils.DisplayMode
}

var _ dotfilesmanager.DotfilesDataInitializer = (*ChezmoiManager)(nil)

func NewChezmoiManager(logger logger.Logger, filesystem utils.FileSystem, userManager osmanager.UserManager, commander utils.Commander, pkgManager pkgmanager.PackageManager, httpClient httpclient.HTTPClient, displayMode utils.DisplayMode, chezmoiConfig ChezmoiConfig) *ChezmoiManager {
	return &ChezmoiManager{
		chezmoiConfig: chezmoiConfig,
		logger:        logger,
		filesystem:    filesystem,
		usermanager:   userManager,
		commander:     commander,
		pkgManager:    pkgManager,
		httpClient:    httpClient,
		displayMode:   displayMode,
	}
}

func TryStandardChezmoiManager(logger logger.Logger, filesystem utils.FileSystem, userManager osmanager.UserManager, commander utils.Commander, pkgManager pkgmanager.PackageManager, httpClient httpclient.HTTPClient, displayMode utils.DisplayMode, githubUsername string, cloneViaSSH bool, branch string) (*ChezmoiManager, error) {
	chezmoiConfigHome, err := userManager.GetChezmoiConfigHome()
	if err != nil {
		return nil, fmt.Errorf("failed to get chezmoi config home directory: %w", err)
	}
	chezmoiConfigDir := fmt.Sprintf("%s/chezmoi", chezmoiConfigHome)
	chezmoiConfigFilePath := fmt.Sprintf("%s/chezmoi.toml", chezmoiConfigDir)

	userHomeDir, err := userManager.GetHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	chezmoiCloneDir := fmt.Sprintf("%s/.local/share/chezmoi", userHomeDir)

	return NewChezmoiManager(
		logger,
		filesystem,
		userManager,
		commander,
		pkgManager,
		httpClient,
		displayMode,
		NewChezmoiConfig(
			chezmoiConfigDir,
			chezmoiConfigFilePath,
			chezmoiCloneDir,
			githubUsername,
			cloneViaSSH,
			branch,
		),
	), nil
}

func TryStandardChezmoiManagerWithDefaults(logger logger.Logger, filesystem utils.FileSystem, userManager osmanager.UserManager, commander utils.Commander, pkgManager pkgmanager.PackageManager, httpClient httpclient.HTTPClient, displayMode utils.DisplayMode) (*ChezmoiManager, error) {
	return TryStandardChezmoiManager(logger, filesystem, userManager, commander, pkgManager, httpClient, displayMode, DefaultGitHubUsername, false, "")
}
