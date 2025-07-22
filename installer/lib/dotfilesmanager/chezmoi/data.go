package chezmoi

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/MrPointer/dotfiles/installer/lib/dotfilesmanager"
	"github.com/MrPointer/dotfiles/installer/utils"
	"github.com/MrPointer/dotfiles/installer/utils/osmanager"
	"github.com/spf13/viper"
)

type ChezmoiDataInitializer struct {
	chezmoiConfigDir      string
	chezmoiConfigFilePath string
	filesystem            utils.FileSystem
}

var _ dotfilesmanager.DotfilesDataInitializer = (*ChezmoiDataInitializer)(nil)

func NewChezmoiDataInitializer(chezmoiConfigFilePath string, filesystem utils.FileSystem) *ChezmoiDataInitializer {
	configFileBasePath := filepath.Dir(chezmoiConfigFilePath)

	return &ChezmoiDataInitializer{
		chezmoiConfigDir:      configFileBasePath,
		chezmoiConfigFilePath: chezmoiConfigFilePath,
		filesystem:            filesystem,
	}
}

func TryNewDefaultChezmoiDataInitializer(filesystem utils.FileSystem, userManager osmanager.UserManager) (*ChezmoiDataInitializer, error) {
	userConfigDir, err := userManager.GetConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user config directory: %w", err)
	}

	chezmoiConfigDir := fmt.Sprintf("%s/chezmoi", userConfigDir)
	chezmoiConfigFilePath := fmt.Sprintf("%s.toml", chezmoiConfigDir)

	return NewChezmoiDataInitializer(chezmoiConfigFilePath, filesystem), nil
}

func (c *ChezmoiDataInitializer) Initialize(data dotfilesmanager.DotfilesData) error {
	if _, err := c.filesystem.PathExists(c.chezmoiConfigDir); os.IsNotExist(err) {
		if err := c.filesystem.CreateDirectory(c.chezmoiConfigDir); err != nil {
			return fmt.Errorf("failed to create chezmoi config directory: %w", err)
		}
	}

	viperObject := viper.New()
	viperObject.SetConfigFile(c.chezmoiConfigFilePath)

	viperObject.Set("data.personal.email", data.Email)
	viperObject.Set("data.personal.full_name", fmt.Sprintf("%s %s", data.FirstName, data.LastName))

	data.GpgSigningKey.MapValue(func(value string) string {
		viperObject.Set("data.gpg.signing_key", value)
		return value
	})

	data.WorkEnv.Match(func(value dotfilesmanager.DotfilesWorkEnvData) (dotfilesmanager.DotfilesWorkEnvData, bool) {
		viperObject.Set("data.personal.work_env", true)
		viperObject.Set("data.personal.work_name", value.WorkName)
		viperObject.Set("data.personal.work_email", value.WorkEmail)
		return value, true
	}, func() (dotfilesmanager.DotfilesWorkEnvData, bool) {
		viperObject.Set("data.personal.work_env", false)
		return dotfilesmanager.DotfilesWorkEnvData{}, false
	})

	data.SystemData.MapValue(func(value dotfilesmanager.DotfilesSystemData) dotfilesmanager.DotfilesSystemData {
		viperObject.Set("data.system.shell", value.Shell)
		viperObject.Set("data.system.user", value.User)
		viperObject.Set("data.system.multi_user_system", value.MultiUserSystem)
		viperObject.Set("data.system.brew_user", value.BrewUser)
		return value
	})

	return viperObject.WriteConfig()
}
