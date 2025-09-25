package chezmoi

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Masterminds/semver"
	"github.com/MrPointer/dotfiles/installer/lib/pkgmanager"
	"github.com/MrPointer/dotfiles/installer/utils"
)

func (c *ChezmoiManager) Install() error {
	chezmoiInstalled, err := c.pkgManager.IsPackageInstalled(pkgmanager.NewPackageInfo("chezmoi", ""))
	if err != nil {
		return err
	}
	if chezmoiInstalled {
		return nil
	}

	c.logger.Debug("Trying to install chezmoi using package manager")
	err = c.tryPackageManagerInstall()
	if err == nil {
		c.logger.Debug("chezmoi installed successfully using package manager")
		return nil
	}
	c.logger.Debug("Failed to install chezmoi using package manager")

	c.logger.Debug("Trying to install chezmoi manually")
	return c.tryManualInstall()
}

// tryPackageManagerInstall attempts to install chezmoi using the system package manager.
// Returns nil if successful, otherwise returns the package manager error.
func (c *ChezmoiManager) tryPackageManagerInstall() error {
	chezmoiVersionConstraint, err := semver.NewConstraint(">=2.60.0")
	if err != nil {
		return err
	}
	return c.pkgManager.InstallPackage(pkgmanager.NewRequestedPackageInfo("chezmoi", chezmoiVersionConstraint))
}

// tryManualInstall attempts to install chezmoi manually by downloading and executing
// the installation script from get.chezmoi.io.
// Returns nil if successful, otherwise returns an error describing the failure.
func (c *ChezmoiManager) tryManualInstall() error {
	c.logger.Debug("Downloading chezmoi binary from official website (get.chezmoi.io)")
	resp, err := c.httpClient.Get("get.chezmoi.io")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download chezmoi binary: %s", resp.Status)
	}

	c.logger.Trace("Reading HTTP response body")
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	userHomeDir, err := c.usermanager.GetHomeDir()
	if err != nil {
		return err
	}

	manualInstallDir := fmt.Sprintf("%s/.local/bin", userHomeDir)

	c.logger.Trace("Executing downloaded binary through the shell")

	var discardOutputOption utils.Option = utils.EmptyOption()
	if c.displayMode.ShouldDiscardOutput() {
		discardOutputOption = utils.WithDiscardOutput()
	}

	result, err := c.commander.RunCommand("sh", []string{"-c", string(body), "--", "-b", manualInstallDir}, discardOutputOption)
	if err != nil {
		return err
	}
	if result.ExitCode != 0 {
		return fmt.Errorf("failed to install chezmoi manually: %s", result.Stderr)
	}

	c.logger.Debug("Chezmoi installed manually successfully")
	return nil
}
