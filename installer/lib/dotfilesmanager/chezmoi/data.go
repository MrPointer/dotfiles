package chezmoi

import (
	"errors"
	"fmt"
	"os"

	"github.com/MrPointer/dotfiles/installer/lib/dotfilesmanager"
	"github.com/spf13/viper"
)

func (c *ChezmoiManager) Initialize(data dotfilesmanager.DotfilesData) error {
	c.logger.Debug("Initializing chezmoi data")

	c.logger.Trace("Creating chezmoi config directory")
	configDirExists, err := c.filesystem.PathExists(c.chezmoiConfig.chezmoiConfigDir)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to check if chezmoi config directory exists: %w", err)
	}
	if !configDirExists {
		if err := c.filesystem.CreateDirectory(c.chezmoiConfig.chezmoiConfigDir); err != nil {
			return fmt.Errorf("failed to create chezmoi config directory: %w", err)
		}
	}
	c.logger.Trace("Chezmoi config directory created")

	c.logger.Trace("Creating chezmoi config file")
	_, err = c.filesystem.CreateFile(c.chezmoiConfig.chezmoiConfigFilePath)
	if err != nil {
		return fmt.Errorf("failed to create chezmoi config file: %w", err)
	}
	c.logger.Trace("Chezmoi config file created")

	c.logger.Trace("Building viper object to contain chezmoi data")
	viperObject := viper.New()
	viperObject.SetConfigFile(c.chezmoiConfig.chezmoiConfigFilePath)

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
		viperObject.Set("data.system.multi_user_system", value.MultiUserSystem)
		viperObject.Set("data.system.brew_multi_user", value.BrewMultiUser)

		if value.GenericWorkProfile.IsPresent() {
			viperObject.Set("data.system.work_generic_dotfiles_profile", value.GenericWorkProfile)
			if value.SpecificWorkProfile.IsPresent() {
				viperObject.Set("data.system.work_specific_dotfiles_profile", value.SpecificWorkProfile)
			}
		}

		return value
	})

	c.logger.Trace("Writing viper object to chezmoi config file")
	return viperObject.WriteConfig()
}
