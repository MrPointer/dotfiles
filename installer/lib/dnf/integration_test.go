package dnf_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/MrPointer/dotfiles/installer/lib/dnf"
	"github.com/MrPointer/dotfiles/installer/lib/pkgmanager"
	"github.com/MrPointer/dotfiles/installer/utils"
	"github.com/MrPointer/dotfiles/installer/utils/logger"
	"github.com/MrPointer/dotfiles/installer/utils/osmanager"
	"github.com/MrPointer/dotfiles/installer/utils/privilege"
)

func Test_DnfPackageManager_CanCheckIfPackageExists_Integration(t *testing.T) {
	if !isDnfAvailable() {
		t.Skip("DNF not available on this system")
	}

	defaultCommander := utils.NewDefaultCommander(logger.DefaultLogger)
	defaultOsManager := osmanager.NewUnixOsManager(logger.DefaultLogger, defaultCommander, false)

	escalator := privilege.NewDefaultEscalator(logger.DefaultLogger, defaultCommander, defaultOsManager)
	dnfManager := dnf.NewDnfPackageManager(logger.DefaultLogger, defaultCommander, defaultOsManager, escalator, utils.DisplayModeProgress)

	// Test with a commonly available package that should exist
	packageInfo := pkgmanager.NewPackageInfo("bash", "")
	isInstalled, err := dnfManager.IsPackageInstalled(packageInfo)

	require.NoError(t, err)
	// bash is typically installed on all systems, but we don't assert the result
	// since the test environment might vary
	t.Logf("bash package installed: %v", isInstalled)
}

func Test_DnfPackageManager_CanCheckIfGroupExists_Integration(t *testing.T) {
	if !isDnfAvailable() {
		t.Skip("DNF not available on this system")
	}

	defaultCommander := utils.NewDefaultCommander(logger.DefaultLogger)
	defaultOsManager := osmanager.NewUnixOsManager(logger.DefaultLogger, defaultCommander, false)

	escalator := privilege.NewDefaultEscalator(logger.DefaultLogger, defaultCommander, defaultOsManager)
	dnfManager := dnf.NewDnfPackageManager(logger.DefaultLogger, defaultCommander, defaultOsManager, escalator, utils.DisplayModeProgress)

	// Test with a commonly available group
	groupInfo := pkgmanager.NewPackageInfoWithType("Development Tools", "", "group")
	isInstalled, err := dnfManager.IsPackageInstalled(groupInfo)

	require.NoError(t, err)
	// Development Tools group existence varies by system, but the check should not fail
	t.Logf("Development Tools group installed: %v", isInstalled)
}

func Test_DnfPackageManager_CanListInstalledPackages_Integration(t *testing.T) {
	if !isDnfAvailable() {
		t.Skip("DNF not available on this system")
	}

	defaultCommander := utils.NewDefaultCommander(logger.DefaultLogger)
	defaultOsManager := osmanager.NewUnixOsManager(logger.DefaultLogger, defaultCommander, false)

	escalator := privilege.NewDefaultEscalator(logger.DefaultLogger, defaultCommander, defaultOsManager)
	dnfManager := dnf.NewDnfPackageManager(logger.DefaultLogger, defaultCommander, defaultOsManager, escalator, utils.DisplayModeProgress)

	packages, err := dnfManager.ListInstalledPackages()

	require.NoError(t, err)
	require.NotNil(t, packages)

	// On a typical system, there should be at least some packages installed
	if len(packages) > 0 {
		t.Logf("Found %d installed packages", len(packages))

		// Verify package structure
		firstPackage := packages[0]
		require.NotEmpty(t, firstPackage.Name)
		require.NotEmpty(t, firstPackage.Version)

		t.Logf("Example package: %s version %s", firstPackage.Name, firstPackage.Version)
	}
}

func Test_DnfPackageManager_CanGetManagerInfo_Integration(t *testing.T) {
	if !isDnfAvailable() {
		t.Skip("DNF not available on this system")
	}

	defaultCommander := utils.NewDefaultCommander(logger.DefaultLogger)
	defaultOsManager := osmanager.NewUnixOsManager(logger.DefaultLogger, defaultCommander, false)

	escalator := privilege.NewDefaultEscalator(logger.DefaultLogger, defaultCommander, defaultOsManager)
	dnfManager := dnf.NewDnfPackageManager(logger.DefaultLogger, defaultCommander, defaultOsManager, escalator, utils.DisplayModeProgress)

	info, err := dnfManager.GetInfo()

	require.NoError(t, err)
	require.Equal(t, "dnf", info.Name)
	require.NotEmpty(t, info.Version)

	t.Logf("DNF version: %s", info.Version)
}

func Test_DnfPackageManager_PrerequisiteInstallationWorkflow_Integration(t *testing.T) {
	if !isDnfAvailable() {
		t.Skip("DNF not available on this system")
	}

	// This test requires root privileges or sudo access
	if os.Getuid() != 0 && !hasSudoAccess() {
		t.Skip("This test requires root privileges or sudo access")
	}

	defaultCommander := utils.NewDefaultCommander(logger.DefaultLogger)
	defaultOsManager := osmanager.NewUnixOsManager(logger.DefaultLogger, defaultCommander, false)

	escalator := privilege.NewDefaultEscalator(logger.DefaultLogger, defaultCommander, defaultOsManager)
	dnfManager := dnf.NewDnfPackageManager(logger.DefaultLogger, defaultCommander, defaultOsManager, escalator, utils.DisplayModeProgress)

	// Test with a lightweight package that's commonly available but might not be installed
	testPackageName := "tree"
	packageInfo := pkgmanager.NewPackageInfo(testPackageName, "")
	requestedPackageInfo := pkgmanager.NewRequestedPackageInfo(testPackageName, nil)

	// Check if package is initially installed
	initiallyInstalled, err := dnfManager.IsPackageInstalled(packageInfo)
	require.NoError(t, err)

	// If not installed, install it
	if !initiallyInstalled {
		t.Logf("Installing test package: %s", testPackageName)
		err = dnfManager.InstallPackage(requestedPackageInfo)
		require.NoError(t, err)

		// Verify installation
		isInstalled, err := dnfManager.IsPackageInstalled(packageInfo)
		require.NoError(t, err)
		require.True(t, isInstalled, "Package should be installed after installation")

		// Get package version
		version, err := dnfManager.GetPackageVersion(testPackageName)
		require.NoError(t, err)
		require.NotEmpty(t, version)
		t.Logf("Installed package version: %s", version)

		// Clean up by removing the package
		t.Logf("Cleaning up test package: %s", testPackageName)
		err = dnfManager.UninstallPackage(packageInfo)
		require.NoError(t, err)

		// Verify removal
		isInstalled, err = dnfManager.IsPackageInstalled(packageInfo)
		require.NoError(t, err)
		require.False(t, isInstalled, "Package should be removed after uninstallation")
	} else {
		t.Logf("Package %s is already installed, skipping install/uninstall test", testPackageName)
	}
}

// isDnfAvailable checks if DNF is available on the system
func isDnfAvailable() bool {
	commander := utils.NewDefaultCommander(logger.DefaultLogger)
	_, err := commander.RunCommand("which", []string{"dnf"}, utils.WithCaptureOutput())
	return err == nil
}

// hasSudoAccess checks if the current user has sudo access
func hasSudoAccess() bool {
	commander := utils.NewDefaultCommander(logger.DefaultLogger)
	_, err := commander.RunCommand("sudo", []string{"-n", "true"}, utils.WithCaptureOutput())
	return err == nil
}
