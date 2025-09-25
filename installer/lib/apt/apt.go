package apt

import (
	"fmt"
	"strings"

	"github.com/MrPointer/dotfiles/installer/lib/pkgmanager"
	"github.com/MrPointer/dotfiles/installer/utils"
	"github.com/MrPointer/dotfiles/installer/utils/logger"
	"github.com/MrPointer/dotfiles/installer/utils/osmanager"
	"github.com/MrPointer/dotfiles/installer/utils/privilege"
)

type AptPackageManager struct {
	logger       logger.Logger
	commander    utils.Commander
	programQuery osmanager.ProgramQuery
	escalator    privilege.Escalator
	displayMode  utils.DisplayMode
}

var _ pkgmanager.PackageManager = (*AptPackageManager)(nil)

// NewAptPackageManager creates a new AptPackageManager instance.
func NewAptPackageManager(logger logger.Logger, commander utils.Commander, programQuery osmanager.ProgramQuery, escalator privilege.Escalator, displayMode utils.DisplayMode) *AptPackageManager {
	return &AptPackageManager{
		logger:       logger,
		commander:    commander,
		programQuery: programQuery,
		escalator:    escalator,
		displayMode:  displayMode,
	}
}

// GetInfo retrieves information about the APT package manager.
func (a *AptPackageManager) GetInfo() (pkgmanager.PackageManagerInfo, error) {
	a.logger.Debug("Getting info about apt")

	aptVersion, err := a.programQuery.GetProgramVersion("apt", func(version string) (string, error) {
		if version == "" {
			return "", nil
		}

		// APT version output typically contains "apt 2.4.8 (amd64)" format.
		// Extract just the version number.
		parts := strings.Fields(version)
		if len(parts) >= 2 {
			return parts[1], nil
		}

		return version, nil
	})
	if err != nil {
		return pkgmanager.DefaultPackageManagerInfo(), fmt.Errorf("failed to get APT version: %w", err)
	}

	return pkgmanager.PackageManagerInfo{
		Name:    "apt",
		Version: aptVersion,
	}, nil
}

// GetPackageVersion retrieves the version of an installed package.
func (a *AptPackageManager) GetPackageVersion(packageName string) (string, error) {
	packages, err := a.ListInstalledPackages()
	if err != nil {
		return "", fmt.Errorf("failed to list installed packages: %w", err)
	}

	for _, pkg := range packages {
		if pkg.Name == packageName {
			return pkg.Version, nil
		}
	}

	return "", fmt.Errorf("package %s is not installed", packageName)
}

// InstallPackage installs a package using APT.
func (a *AptPackageManager) InstallPackage(requestedPackageInfo pkgmanager.RequestedPackageInfo) error {
	a.logger.Debug("Installing package %s with apt", requestedPackageInfo.Name)

	if requestedPackageInfo.VersionConstraints != nil {
		a.logger.Warning("APT doesn't support version constraints, installing the latest version of package %s", requestedPackageInfo.Name)
	}

	// Update package list first to ensure we have the latest package information.
	escalatedUpdate, err := a.escalator.EscalateCommand("apt", []string{"update"})
	if err != nil {
		return fmt.Errorf("failed to determine privilege escalation for apt update: %w", err)
	}

	var discardOutputOption utils.Option = utils.EmptyOption()
	if a.displayMode.ShouldDiscardOutput() {
		discardOutputOption = utils.WithDiscardOutput()
	}

	_, err = a.commander.RunCommand(escalatedUpdate.Command, escalatedUpdate.Args, discardOutputOption)
	if err != nil {
		return fmt.Errorf("failed to update package list: %w", err)
	}

	// Install the package with automatic yes confirmation.
	escalatedInstall, err := a.escalator.EscalateCommand("apt", []string{"install", "-y", requestedPackageInfo.Name})
	if err != nil {
		return fmt.Errorf("failed to determine privilege escalation for apt install: %w", err)
	}

	discardOutputOption = utils.EmptyOption()
	if a.displayMode.ShouldDiscardOutput() {
		discardOutputOption = utils.WithDiscardOutput()
	}

	_, err = a.commander.RunCommand(escalatedInstall.Command, escalatedInstall.Args, discardOutputOption)
	if err != nil {
		return fmt.Errorf("failed to install package %s: %w", requestedPackageInfo.Name, err)
	}

	a.logger.Debug("Package %s installed successfully with apt", requestedPackageInfo.Name)
	return nil
}

// IsPackageInstalled checks if a package is installed.
func (a *AptPackageManager) IsPackageInstalled(packageInfo pkgmanager.PackageInfo) (bool, error) {
	a.logger.Debug("Checking if package %s is installed with apt", packageInfo.Name)

	packages, err := a.ListInstalledPackages()
	if err != nil {
		return false, fmt.Errorf("failed to list installed packages: %w", err)
	}

	for _, pkg := range packages {
		if pkg.Name == packageInfo.Name {
			a.logger.Debug("Package %s is installed with apt", packageInfo.Name)
			return true, nil
		}
	}

	a.logger.Debug("Package %s is not installed with apt", packageInfo.Name)
	return false, nil
}

// ListInstalledPackages returns a list of all installed packages.
func (a *AptPackageManager) ListInstalledPackages() ([]pkgmanager.PackageInfo, error) {
	a.logger.Debug("Listing packages installed with apt")

	// Use dpkg-query to get installed packages with versions.
	output, err := a.commander.RunCommand("dpkg-query", []string{"-W", "-f=${Package} ${Version}\n"}, utils.WithCaptureOutput())
	if err != nil {
		return nil, fmt.Errorf("failed to list installed packages: %w", err)
	}

	trimmedOutput := strings.TrimSpace(string(output.Stdout))
	if trimmedOutput == "" {
		return []pkgmanager.PackageInfo{}, nil
	}

	a.logger.Trace("Raw output from dpkg-query: %s", trimmedOutput)

	lines := strings.Split(trimmedOutput, "\n")
	packages := make([]pkgmanager.PackageInfo, 0, len(lines))

	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) >= 2 {
			name := parts[0]
			version := parts[1]
			packages = append(packages, pkgmanager.NewPackageInfo(name, version))
		}
	}

	return packages, nil
}

// UninstallPackage uninstalls a package using APT.
func (a *AptPackageManager) UninstallPackage(packageInfo pkgmanager.PackageInfo) error {
	a.logger.Debug("Uninstalling package %s with apt", packageInfo.Name)

	removeResult, err := a.escalator.EscalateCommand("apt", []string{"remove", "-y", packageInfo.Name})
	if err != nil {
		return fmt.Errorf("failed to determine privilege escalation for apt remove: %w", err)
	}

	var discardOutputOption utils.Option = utils.EmptyOption()
	if a.displayMode.ShouldDiscardOutput() {
		discardOutputOption = utils.WithDiscardOutput()
	}

	_, err = a.commander.RunCommand(removeResult.Command, removeResult.Args, discardOutputOption)
	if err != nil {
		return fmt.Errorf("failed to uninstall package %s: %w", packageInfo.Name, err)
	}

	a.logger.Debug("Package %s uninstalled successfully with apt", packageInfo.Name)
	return nil
}
