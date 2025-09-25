package apt_test

import (
	"testing"

	"github.com/MrPointer/dotfiles/installer/lib/apt"
	"github.com/MrPointer/dotfiles/installer/lib/pkgmanager"
	"github.com/MrPointer/dotfiles/installer/utils"
	"github.com/MrPointer/dotfiles/installer/utils/logger"
	"github.com/MrPointer/dotfiles/installer/utils/osmanager"
	"github.com/MrPointer/dotfiles/installer/utils/privilege"
	"github.com/stretchr/testify/require"
)

func Test_AptPackageManager_CanCheckIfPackageExists_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	mockLogger := &logger.MoqLogger{
		InfoFunc:    func(format string, args ...interface{}) {},
		WarningFunc: func(format string, args ...interface{}) {},
		ErrorFunc:   func(format string, args ...interface{}) {},
		SuccessFunc: func(format string, args ...interface{}) {},
	}

	// This test only runs on systems where apt is available
	defaultCommander := utils.NewDefaultCommander(mockLogger)
	defaultOsManager := osmanager.NewUnixOsManager(mockLogger, defaultCommander, false)

	// Check if apt is available on this system
	aptExists, err := defaultOsManager.ProgramExists("apt")
	if err != nil || !aptExists {
		t.Skip("apt not available on this system")
	}

	escalator := privilege.NewDefaultEscalator(mockLogger, defaultCommander, defaultOsManager)
	aptManager := apt.NewAptPackageManager(mockLogger, defaultCommander, defaultOsManager, escalator, utils.DisplayModeProgress)

	// Test with a commonly available package that should exist
	packageInfo := pkgmanager.NewPackageInfo("libc6", "")
	isInstalled, err := aptManager.IsPackageInstalled(packageInfo)

	require.NoError(t, err)
	// libc6 should be installed on any Debian/Ubuntu system
	require.True(t, isInstalled)
}

func Test_AptPackageManager_CanListInstalledPackages_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	mockLogger := &logger.MoqLogger{
		InfoFunc:    func(format string, args ...interface{}) {},
		WarningFunc: func(format string, args ...interface{}) {},
		ErrorFunc:   func(format string, args ...interface{}) {},
		SuccessFunc: func(format string, args ...interface{}) {},
	}

	// This test only runs on systems where apt is available
	defaultCommander := utils.NewDefaultCommander(mockLogger)
	defaultOsManager := osmanager.NewUnixOsManager(mockLogger, defaultCommander, false)

	// Check if apt is available on this system
	aptExists, err := defaultOsManager.ProgramExists("apt")
	if err != nil || !aptExists {
		t.Skip("apt not available on this system")
	}

	escalator := privilege.NewDefaultEscalator(mockLogger, defaultCommander, defaultOsManager)
	aptManager := apt.NewAptPackageManager(mockLogger, defaultCommander, defaultOsManager, escalator, utils.DisplayModeProgress)

	packages, err := aptManager.ListInstalledPackages()

	require.NoError(t, err)
	require.NotEmpty(t, packages)

	// Verify structure of returned packages
	for _, pkg := range packages {
		require.NotEmpty(t, pkg.Name)
		require.NotEmpty(t, pkg.Version)
	}
}

func Test_AptPackageManager_CanGetManagerInfo_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	mockLogger := &logger.MoqLogger{
		InfoFunc:    func(format string, args ...interface{}) {},
		WarningFunc: func(format string, args ...interface{}) {},
		ErrorFunc:   func(format string, args ...interface{}) {},
		SuccessFunc: func(format string, args ...interface{}) {},
	}

	// This test only runs on systems where apt is available
	defaultCommander := utils.NewDefaultCommander(mockLogger)
	defaultOsManager := osmanager.NewUnixOsManager(mockLogger, defaultCommander, false)

	// Check if apt is available on this system
	aptExists, err := defaultOsManager.ProgramExists("apt")
	if err != nil || !aptExists {
		t.Skip("apt not available on this system")
	}

	escalator := privilege.NewDefaultEscalator(mockLogger, defaultCommander, defaultOsManager)
	aptManager := apt.NewAptPackageManager(mockLogger, defaultCommander, defaultOsManager, escalator, utils.DisplayModeProgress)

	info, err := aptManager.GetInfo()

	require.NoError(t, err)
	require.Equal(t, "apt", info.Name)
	require.NotEmpty(t, info.Version)
}

func Test_AptPackageManager_PrerequisiteInstallationWorkflow_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	mockLogger := &logger.MoqLogger{
		InfoFunc:    func(format string, args ...interface{}) {},
		WarningFunc: func(format string, args ...interface{}) {},
		ErrorFunc:   func(format string, args ...interface{}) {},
		SuccessFunc: func(format string, args ...interface{}) {},
	}

	// This test only runs on systems where apt is available
	defaultCommander := utils.NewDefaultCommander(mockLogger)
	defaultOsManager := osmanager.NewUnixOsManager(mockLogger, defaultCommander, false)

	// Check if apt is available on this system
	aptExists, err := defaultOsManager.ProgramExists("apt")
	if err != nil || !aptExists {
		t.Skip("apt not available on this system")
	}

	// Check if we're running as root or can use sudo
	_, err = defaultCommander.RunCommand("sudo", []string{"-n", "true"}, utils.WithCaptureOutput())
	if err != nil {
		t.Skip("sudo not available or requires password - skipping package installation test")
	}

	escalator := privilege.NewDefaultEscalator(mockLogger, defaultCommander, defaultOsManager)
	aptManager := apt.NewAptPackageManager(mockLogger, defaultCommander, defaultOsManager, escalator, utils.DisplayModeProgress)

	// Test with a lightweight package that's commonly available but might not be installed
	testPackage := "file"
	packageInfo := pkgmanager.NewRequestedPackageInfo(testPackage, nil)

	// Check initial state
	initiallyInstalled, err := aptManager.IsPackageInstalled(pkgmanager.NewPackageInfo(testPackage, ""))
	require.NoError(t, err)

	if !initiallyInstalled {
		// Install the package
		err = aptManager.InstallPackage(packageInfo)
		require.NoError(t, err)

		// Verify it's now installed
		isInstalled, err := aptManager.IsPackageInstalled(pkgmanager.NewPackageInfo(testPackage, ""))
		require.NoError(t, err)
		require.True(t, isInstalled)

		// Clean up by removing it
		err = aptManager.UninstallPackage(pkgmanager.NewPackageInfo(testPackage, ""))
		require.NoError(t, err)
	} else {
		t.Log("Package already installed, skipping installation test")
	}
}
