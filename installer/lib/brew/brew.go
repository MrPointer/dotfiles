package brew

import (
	"errors"
	"fmt"
	"strings"

	"github.com/MrPointer/dotfiles/installer/lib/pkgmanager"
	"github.com/MrPointer/dotfiles/installer/utils"
	"github.com/MrPointer/dotfiles/installer/utils/logger"
	"github.com/MrPointer/dotfiles/installer/utils/osmanager"
)

type BrewPackageManager struct {
	logger       logger.Logger
	commander    utils.Commander
	programQuery osmanager.ProgramQuery
	brewPath     string
	displayMode  utils.DisplayMode
}

var _ pkgmanager.PackageManager = &BrewPackageManager{}

// NewBrewPackageManager creates a new BrewPackageManager instance.
func NewBrewPackageManager(logger logger.Logger, commander utils.Commander, programQuery osmanager.ProgramQuery, brewPath string, displayMode utils.DisplayMode) *BrewPackageManager {
	return &BrewPackageManager{
		logger:       logger,
		commander:    commander,
		programQuery: programQuery,
		brewPath:     brewPath,
		displayMode:  displayMode,
	}
}

func (b *BrewPackageManager) GetInfo() (pkgmanager.PackageManagerInfo, error) {
	b.logger.Debug("Getting info about Homebrew")

	brewVersion, err := b.programQuery.GetProgramVersion(b.brewPath, func(version string) (string, error) {
		if version == "" {
			return "", nil
		}

		// Brew's version is typically in the format "Homebrew 3.4.0", so we extract the version number.
		parts := strings.Split(version, " ")
		if len(parts) < 2 {
			return "", errors.New("unexpected version format: " + version)
		}

		return parts[1], nil
	})
	if err != nil {
		return pkgmanager.DefaultPackageManagerInfo(), errors.New("failed to get Homebrew version: " + err.Error())
	}

	return pkgmanager.PackageManagerInfo{
		Name:    "brew",
		Version: brewVersion,
	}, nil
}

// GetPackageVersion implements pkgmanager.PackageManager.
func (b *BrewPackageManager) GetPackageVersion(packageName string) (string, error) {
	b.logger.Debug("Getting version of package %s with Homebrew", packageName)

	// Get list of installed packages with versions, then find the requested package.
	packages, err := b.ListInstalledPackages()
	if err != nil {
		return "", errors.New("failed to list installed packages with Homebrew: " + err.Error())
	}

	for _, pkg := range packages {
		if pkg.Name == packageName {
			b.logger.Debug("Found package %s with version %s", packageName, pkg.Version)
			return pkg.Version, nil
		}
	}

	b.logger.Debug("Package %s not found with Homebrew", packageName)
	return "", fmt.Errorf("package %s is not installed with Homebrew", packageName)
}

// InstallPackage implements pkgmanager.PackageManager.
func (b *BrewPackageManager) InstallPackage(requestedPackageInfo pkgmanager.RequestedPackageInfo) error {
	b.logger.Debug("Installing package %s with Homebrew", requestedPackageInfo.Name)

	if requestedPackageInfo.VersionConstraints != nil {
		b.logger.Warning("Homebrew doesn't support version constraints, installing the latest version of package %s", requestedPackageInfo.Name)
	}

	var discardOutputOption utils.Option = utils.EmptyOption()
	if b.displayMode.ShouldDiscardOutput() {
		discardOutputOption = utils.WithDiscardOutput()
	}

	_, err := b.commander.RunCommand(b.brewPath, []string{"install", requestedPackageInfo.Name}, discardOutputOption)
	if err != nil {
		return fmt.Errorf("failed to install package %s with Homebrew: %v", requestedPackageInfo.Name, err)
	}

	b.logger.Debug("Package %s installed successfully with Homebrew", requestedPackageInfo.Name)
	return nil
}

// IsPackageInstalled implements pkgmanager.PackageManager.
func (b *BrewPackageManager) IsPackageInstalled(packageInfo pkgmanager.PackageInfo) (bool, error) {
	b.logger.Debug("Checking if package %s is installed with Homebrew", packageInfo.Name)

	// Check if the package is installed by listing all installed packages and checking for the package name.
	packages, err := b.ListInstalledPackages()
	if err != nil {
		return false, errors.New("failed to list installed packages with Homebrew: " + err.Error())
	}

	for _, pkg := range packages {
		if pkg.Name == packageInfo.Name {
			b.logger.Debug("Package %s is installed with Homebrew", packageInfo.Name)
			return true, nil
		}
	}

	b.logger.Debug("Package %s is not installed with Homebrew", packageInfo.Name)
	return false, nil
}

// ListInstalledPackages implements pkgmanager.PackageManager.
func (b *BrewPackageManager) ListInstalledPackages() ([]pkgmanager.PackageInfo, error) {
	b.logger.Debug("Listing packages installed by Homebrew")

	// Run `brew list` to get the list of installed packages.
	output, err := b.commander.RunCommand(b.brewPath, []string{"list", "--versions"}, utils.WithCaptureOutput())
	if err != nil {
		return nil, errors.New("failed to list installed packages with Homebrew: " + err.Error())
	}

	// Split the output into lines and create PackageInfo for each package.
	trimmedOutput := strings.TrimSpace(output.String())
	if trimmedOutput == "" {
		return []pkgmanager.PackageInfo{}, nil
	}

	b.logger.Trace("Raw list of installed packages: %s", trimmedOutput)

	rawPackages := strings.Split(trimmedOutput, "\n")
	var packages []pkgmanager.PackageInfo
	for _, pkg := range rawPackages {
		name, version, found := strings.Cut(pkg, " ")
		if !found {
			return nil, errors.New("failed to parse package line: " + pkg)
		}
		packages = append(packages, pkgmanager.NewPackageInfo(name, version))
	}

	return packages, nil
}

// UninstallPackage implements pkgmanager.PackageManager.
func (b *BrewPackageManager) UninstallPackage(packageInfo pkgmanager.PackageInfo) error {
	b.logger.Debug("Uninstalling package %s with Homebrew", packageInfo.Name)

	var discardOutputOption utils.Option = utils.EmptyOption()
	if b.displayMode.ShouldDiscardOutput() {
		discardOutputOption = utils.WithDiscardOutput()
	}

	_, err := b.commander.RunCommand(b.brewPath, []string{"uninstall", packageInfo.Name}, discardOutputOption)
	if err != nil {
		return fmt.Errorf("failed to uninstall package %s with Homebrew: %v", packageInfo.Name, err)
	}

	b.logger.Debug("Package %s uninstalled successfully", packageInfo.Name)
	return nil
}
