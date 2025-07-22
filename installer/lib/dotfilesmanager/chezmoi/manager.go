package chezmoi

import (
	"fmt"
	"path/filepath"

	"github.com/MrPointer/dotfiles/installer/lib/dotfilesmanager"
	"github.com/MrPointer/dotfiles/installer/utils"
	"github.com/MrPointer/dotfiles/installer/utils/osmanager"
)

type ChezmoiManager struct {
	chezmoiConfigDir      string
	chezmoiConfigFilePath string
	chezmoiCloneDir       string
	filesystem            utils.FileSystem
}

var _ dotfilesmanager.DotfilesDataInitializer = (*ChezmoiManager)(nil)

func NewChezmoiManager(chezmoiConfigFilePath string, chezmoiCloneDir string, filesystem utils.FileSystem) *ChezmoiManager {
	configFileBasePath := filepath.Dir(chezmoiConfigFilePath)

	return &ChezmoiManager{
		chezmoiConfigDir:      configFileBasePath,
		chezmoiConfigFilePath: chezmoiConfigFilePath,
		chezmoiCloneDir:       chezmoiCloneDir,
		filesystem:            filesystem,
	}
}

func TryNewDefaultChezmoiManager(filesystem utils.FileSystem, userManager osmanager.UserManager) (*ChezmoiManager, error) {
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

	return NewChezmoiManager(chezmoiConfigFilePath, chezmoiCloneDir, filesystem), nil
}
