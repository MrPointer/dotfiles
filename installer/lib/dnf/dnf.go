package dnf

import (
	"fmt"
	"strings"

	"github.com/MrPointer/dotfiles/installer/lib/pkgmanager"
	"github.com/MrPointer/dotfiles/installer/utils"
	"github.com/MrPointer/dotfiles/installer/utils/logger"
	"github.com/MrPointer/dotfiles/installer/utils/osmanager"
	"github.com/MrPointer/dotfiles/installer/utils/privilege"
)

type DnfPackageManager struct {
	logger       logger.Logger
	commander    utils.Commander
	programQuery osmanager.ProgramQuery
	escalator    privilege.Escalator
	displayMode  utils.DisplayMode
}

var _ pkgmanager.PackageManager = (*DnfPackageManager)(nil)

// NewDnfPackageManager creates a new DnfPackageManager instance.
func NewDnfPackageManager(logger logger.Logger, commander utils.Commander, programQuery osmanager.ProgramQuery, escalator privilege.Escalator, displayMode utils.DisplayMode) *DnfPackageManager {
	return &DnfPackageManager{
		logger:       logger,
		commander:    commander,
		programQuery: programQuery,
		escalator:    escalator,
		displayMode:  displayMode,
	}
}

// GetInfo retrieves information about the DNF package manager.
func (d *DnfPackageManager) GetInfo() (pkgmanager.PackageManagerInfo, error) {
	d.logger.Debug("Getting info about dnf")

	dnfVersion, err := d.programQuery.GetProgramVersion("dnf", func(version string) (string, error) {
		if version == "" {
			return "", nil
		}

		// DNF version output typically contains "dnf 4.14.0" format.
		// Extract just the version number.
		parts := strings.Fields(version)
		if len(parts) >= 2 {
			return parts[1], nil
		}

		return version, nil
	})
	if err != nil {
		return pkgmanager.DefaultPackageManagerInfo(), fmt.Errorf("failed to get DNF version: %w", err)
	}

	return pkgmanager.PackageManagerInfo{
		Name:    "dnf",
		Version: dnfVersion,
	}, nil
}

// GetPackageVersion retrieves the version of an installed package.
func (d *DnfPackageManager) GetPackageVersion(packageName string) (string, error) {
	packages, err := d.ListInstalledPackages()
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

// InstallPackage installs a package using DNF.
func (d *DnfPackageManager) InstallPackage(requestedPackageInfo pkgmanager.RequestedPackageInfo) error {
	d.logger.Debug("Installing package %s with dnf", requestedPackageInfo.Name)

	if requestedPackageInfo.VersionConstraints != nil {
		d.logger.Debug("DNF doesn't support version constraints, installing the latest version of package %s", requestedPackageInfo.Name)
	}

	var installArgs []string
	if requestedPackageInfo.Type == "group" {
		d.logger.Debug("Installing package group %s with dnf", requestedPackageInfo.Name)
		installArgs = []string{"group", "install", "-y", requestedPackageInfo.Name}
	} else {
		d.logger.Debug("Installing regular package %s with dnf", requestedPackageInfo.Name)
		installArgs = []string{"install", "-y", requestedPackageInfo.Name}
	}

	escalatedInstall, err := d.escalator.EscalateCommand("dnf", installArgs)
	if err != nil {
		return fmt.Errorf("failed to determine privilege escalation for dnf install: %w", err)
	}

	var discardOutputOption utils.Option = utils.EmptyOption()
	if d.displayMode.ShouldDiscardOutput() {
		discardOutputOption = utils.WithDiscardOutput()
	}

	_, err = d.commander.RunCommand(escalatedInstall.Command, escalatedInstall.Args, discardOutputOption)
	if err != nil {
		return fmt.Errorf("failed to install package %s: %w", requestedPackageInfo.Name, err)
	}

	d.logger.Debug("Package %s installed successfully with dnf", requestedPackageInfo.Name)
	return nil
}

// IsPackageInstalled checks if a package is installed.
func (d *DnfPackageManager) IsPackageInstalled(packageInfo pkgmanager.PackageInfo) (bool, error) {
	d.logger.Debug("Checking if package %s is installed with dnf", packageInfo.Name)

	if packageInfo.Type == "group" {
		return d.isGroupInstalled(packageInfo.Name)
	}

	packages, err := d.ListInstalledPackages()
	if err != nil {
		return false, fmt.Errorf("failed to list installed packages: %w", err)
	}

	for _, pkg := range packages {
		if pkg.Name == packageInfo.Name {
			d.logger.Debug("Package %s is installed with dnf", packageInfo.Name)
			return true, nil
		}
	}

	d.logger.Debug("Package %s is not installed with dnf", packageInfo.Name)
	return false, nil
}

// isGroupInstalled checks if a package group is installed.
func (d *DnfPackageManager) isGroupInstalled(groupName string) (bool, error) {
	d.logger.Debug("Checking if group %s is installed with dnf", groupName)

	output, err := d.commander.RunCommand("dnf", []string{"group", "list", "installed"}, utils.WithCaptureOutput())
	if err != nil {
		return false, fmt.Errorf("failed to list installed groups: %w", err)
	}

	trimmedOutput := strings.TrimSpace(string(output.Stdout))
	d.logger.Trace("Raw output from dnf group list installed: %s", trimmedOutput)

	// DNF group list output contains group names, check if our group is listed
	lines := strings.Split(trimmedOutput, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, groupName) {
			d.logger.Debug("Group %s is installed with dnf", groupName)
			return true, nil
		}
	}

	d.logger.Debug("Group %s is not installed with dnf", groupName)
	return false, nil
}

// ListInstalledPackages returns a list of all installed packages.
func (d *DnfPackageManager) ListInstalledPackages() ([]pkgmanager.PackageInfo, error) {
	d.logger.Debug("Listing packages installed with dnf")

	// Use dnf list installed to get installed packages with versions.
	output, err := d.commander.RunCommand("dnf", []string{"list", "installed"}, utils.WithCaptureOutput())
	if err != nil {
		return nil, fmt.Errorf("failed to list installed packages: %w", err)
	}

	trimmedOutput := strings.TrimSpace(string(output.Stdout))
	if trimmedOutput == "" {
		return []pkgmanager.PackageInfo{}, nil
	}

	d.logger.Trace("Raw output from dnf list installed: %s", trimmedOutput)

	lines := strings.Split(trimmedOutput, "\n")
	packages := make([]pkgmanager.PackageInfo, 0)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Installed Packages") {
			continue
		}

		// DNF output format: "package-name.arch version repo"
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			// Extract package name (remove architecture suffix if present)
			packageWithArch := parts[0]
			packageName := packageWithArch
			if dotIndex := strings.LastIndex(packageWithArch, "."); dotIndex != -1 {
				packageName = packageWithArch[:dotIndex]
			}

			version := parts[1]
			packages = append(packages, pkgmanager.NewPackageInfo(packageName, version))
		}
	}

	return packages, nil
}

// UninstallPackage uninstalls a package using DNF.
func (d *DnfPackageManager) UninstallPackage(packageInfo pkgmanager.PackageInfo) error {
	d.logger.Debug("Uninstalling package %s with dnf", packageInfo.Name)

	var removeArgs []string
	if packageInfo.Type == "group" {
		d.logger.Debug("Uninstalling package group %s with dnf", packageInfo.Name)
		removeArgs = []string{"group", "remove", "-y", packageInfo.Name}
	} else {
		d.logger.Debug("Uninstalling regular package %s with dnf", packageInfo.Name)
		removeArgs = []string{"remove", "-y", packageInfo.Name}
	}

	removeResult, err := d.escalator.EscalateCommand("dnf", removeArgs)
	if err != nil {
		return fmt.Errorf("failed to determine privilege escalation for dnf remove: %w", err)
	}

	var discardOutputOption utils.Option = utils.EmptyOption()
	if d.displayMode.ShouldDiscardOutput() {
		discardOutputOption = utils.WithDiscardOutput()
	}

	_, err = d.commander.RunCommand(removeResult.Command, removeResult.Args, discardOutputOption)
	if err != nil {
		return fmt.Errorf("failed to uninstall package %s: %w", packageInfo.Name, err)
	}

	d.logger.Debug("Package %s uninstalled successfully with dnf", packageInfo.Name)
	return nil
}
