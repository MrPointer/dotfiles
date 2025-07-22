package chezmoi

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Masterminds/semver"
	"github.com/MrPointer/dotfiles/installer/lib/pkgmanager"
)

func (c *ChezmoiManager) Install() error {
	chezmoiInstalled, err := c.pkgManager.IsPackageInstalled(pkgmanager.NewPackageInfo("chezmoi", ""))
	if err != nil {
		return err
	}
	if chezmoiInstalled {
		return nil
	}

	err = c.tryPackageManagerInstall()
	if err == nil {
		return nil
	}

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
	resp, err := c.httpClient.Get("get.chezmoi.io")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download chezmoi binary: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	userHomeDir, err := c.usermanager.GetHomeDir()
	if err != nil {
		return err
	}

	manualInstallDir := fmt.Sprintf("%s/.local/bin", userHomeDir)

	result, err := c.commander.RunCommand("sh", []string{"-c", string(body), "--", "-b", manualInstallDir})
	if err != nil {
		return err
	}
	if result.ExitCode != 0 {
		return fmt.Errorf("failed to install chezmoi manually: %s", result.Stderr)
	}

	return nil
}
