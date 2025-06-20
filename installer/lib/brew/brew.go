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
}

var _ pkgmanager.PackageManager = &BrewPackageManager{}

// NewBrewPackageManager creates a new BrewPackageManager instance.
func NewBrewPackageManager(logger logger.Logger, commander utils.Commander, programQuery osmanager.ProgramQuery) *BrewPackageManager {
	return &BrewPackageManager{
		logger:       logger,
		commander:    commander,
		programQuery: programQuery,
	}
}

func (b *BrewPackageManager) GetInfo() (pkgmanager.PackageManagerInfo, error) {
	brewVersion, err := b.programQuery.GetProgramVersion("brew", func(version string) (string, error) {
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
	// Get list of installed packages with versions, then find the requested package.
	packages, err := b.ListInstalledPackages()
	if err != nil {
		return "", errors.New("failed to list installed packages with Homebrew: " + err.Error())
	}

	for _, pkg := range packages {
		if pkg.Name == packageName {
			return pkg.Version, nil
		}
	}

	return "", errors.New(fmt.Sprintf("package %s is not installed with Homebrew", packageName))
}

// InstallPackage implements pkgmanager.PackageManager.
func (b *BrewPackageManager) InstallPackage(requestedPackageInfo pkgmanager.RequestedPackageInfo) error {
	b.logger.Warning("Homebrew doesn't support version constraints, installing the latest version of package %s", requestedPackageInfo.Name)

	_, err := b.commander.RunCommand("brew", []string{"install", requestedPackageInfo.Name})
	if err != nil {
		return errors.New(fmt.Sprintf("failed to install package %s with Homebrew: %v", requestedPackageInfo.Name, err))
	}

	return nil
}

// IsPackageInstalled implements pkgmanager.PackageManager.
func (b *BrewPackageManager) IsPackageInstalled(packageInfo pkgmanager.PackageInfo) (bool, error) {
	// Check if the package is installed by listing all installed packages and checking for the package name.
	packages, err := b.ListInstalledPackages()
	if err != nil {
		return false, errors.New("failed to list installed packages with Homebrew: " + err.Error())
	}

	for _, pkg := range packages {
		if pkg.Name == packageInfo.Name {
			return true, nil
		}
	}

	return false, nil
}

// ListInstalledPackages implements pkgmanager.PackageManager.
func (b *BrewPackageManager) ListInstalledPackages() ([]pkgmanager.PackageInfo, error) {
	// Run `brew list` to get the list of installed packages.
	output, err := b.commander.RunCommand("brew", []string{"list", "--versions"})
	if err != nil {
		return nil, errors.New("failed to list installed packages with Homebrew: " + err.Error())
	}

	// Split the output into lines and create PackageInfo for each package.
	trimmedOutput := strings.TrimSpace(output.String())
	if trimmedOutput == "" {
		return []pkgmanager.PackageInfo{}, nil
	}

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
	_, err := b.commander.RunCommand("brew", []string{"uninstall", packageInfo.Name})
	if err != nil {
		return errors.New(fmt.Sprintf("failed to uninstall package %s with Homebrew: %v", packageInfo.Name, err))
	}
	return nil
}
