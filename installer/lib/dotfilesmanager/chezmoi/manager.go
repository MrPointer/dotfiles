package chezmoi

import (
	"fmt"
	"path/filepath"

	"github.com/MrPointer/dotfiles/installer/lib/dotfilesmanager"
	"github.com/MrPointer/dotfiles/installer/utils"
	"github.com/MrPointer/dotfiles/installer/utils/osmanager"
)

type ChezmoiConfig struct {
	chezmoiConfigDir      string
	chezmoiConfigFilePath string
	chezmoiCloneDir       string
	githubUsername        string
	cloneViaSSH           bool
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
	commander     utils.Commander
}

var _ dotfilesmanager.DotfilesDataInitializer = (*ChezmoiManager)(nil)

func NewChezmoiManager(filesystem utils.FileSystem, commander utils.Commander, chezmoiConfig ChezmoiConfig) *ChezmoiManager {
	return &ChezmoiManager{
		chezmoiConfig: chezmoiConfig,
		filesystem:    filesystem,
		commander:     commander,
	}
}

func TryStandardChezmoiManager(filesystem utils.FileSystem, userManager osmanager.UserManager, commander utils.Commander, githubUsername string, cloneViaSSH bool) (*ChezmoiManager, error) {
	userConfigDir, err := userManager.GetConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user config directory: %w", err)
	}

	chezmoiConfigDir := fmt.Sprintf("%s/chezmoi", userConfigDir)
	chezmoiConfigFilePath := fmt.Sprintf("%s.toml", chezmoiConfigDir)

	userHomeDir, err := userManager.GetHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	chezmoiCloneDir := fmt.Sprintf("%s/.local/share/chezmoi", userHomeDir)

	return NewChezmoiManager(
		filesystem,
		commander,
		ChezmoiConfig{
			chezmoiConfigDir:      chezmoiConfigDir,
			chezmoiConfigFilePath: chezmoiConfigFilePath,
			chezmoiCloneDir:       chezmoiCloneDir,
			githubUsername:        githubUsername,
			cloneViaSSH:           cloneViaSSH,
		}), nil
}

func TryStandardChezmoiManagerWithDefaults(filesystem utils.FileSystem, userManager osmanager.UserManager, commander utils.Commander) (*ChezmoiManager, error) {
	return TryStandardChezmoiManager(filesystem, userManager, commander, "MrPointer", false)
}
