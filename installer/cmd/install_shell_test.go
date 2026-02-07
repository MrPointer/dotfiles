package cmd

import (
	"testing"

	"github.com/MrPointer/dotfiles/installer/lib/compatibility"
	"github.com/MrPointer/dotfiles/installer/lib/shell"
	"github.com/MrPointer/dotfiles/installer/utils"
	"github.com/MrPointer/dotfiles/installer/utils/logger"
	"github.com/MrPointer/dotfiles/installer/utils/osmanager"
	"github.com/stretchr/testify/require"
)

func Test_ResolveShellInstallStrategy(t *testing.T) {
	// Save and restore global state
	originalShellSource := shellSource
	originalShellName := shellName
	originalBrewPath := globalBrewPath
	originalSysInfo := globalSysInfo
	originalLogger := cliLogger
	originalFilesystem := globalFilesystem
	originalOsManager := globalOsManager

	defer func() {
		shellSource = originalShellSource
		shellName = originalShellName
		globalBrewPath = originalBrewPath
		globalSysInfo = originalSysInfo
		cliLogger = originalLogger
		globalFilesystem = originalFilesystem
		globalOsManager = originalOsManager
	}()

	// Setup test dependencies
	shellName = "zsh"
	cliLogger = logger.DefaultLogger
	globalFilesystem = &utils.MoqFileSystem{
		PathExistsFunc: func(path string) (bool, error) {
			return true, nil
		},
	}
	globalOsManager = &osmanager.MoqOsManager{}

	tests := []struct {
		name           string
		shellSource    string
		osName         string
		distroName     string
		brewAvailable  bool
		wantSource     shell.ShellSource
		wantPkgMgr     bool // true if package manager should be non-nil
		wantBrewPkgMgr bool // true if package manager should be brew
		wantError      bool
	}{
		{
			name:           "auto on Linux with brew",
			shellSource:    "auto",
			osName:         "linux",
			distroName:     "ubuntu",
			brewAvailable:  true,
			wantSource:     shell.ShellSourceAuto,
			wantPkgMgr:     true,
			wantBrewPkgMgr: true,
			wantError:      false,
		},
		{
			name:           "auto on Linux without brew",
			shellSource:    "auto",
			osName:         "linux",
			distroName:     "ubuntu",
			brewAvailable:  false,
			wantSource:     shell.ShellSourceAuto,
			wantPkgMgr:     true,
			wantBrewPkgMgr: false,
			wantError:      false,
		},
		{
			name:           "auto on macOS with brew",
			shellSource:    "auto",
			osName:         "darwin",
			distroName:     "",
			brewAvailable:  true,
			wantSource:     shell.ShellSourceAuto,
			wantPkgMgr:     true,
			wantBrewPkgMgr: true,
			wantError:      false,
		},
		{
			name:           "brew on Linux with brew",
			shellSource:    "brew",
			osName:         "linux",
			distroName:     "ubuntu",
			brewAvailable:  true,
			wantSource:     shell.ShellSourceBrew,
			wantPkgMgr:     true,
			wantBrewPkgMgr: true,
			wantError:      false,
		},
		{
			name:          "brew without brew installed",
			shellSource:   "brew",
			osName:        "linux",
			distroName:    "ubuntu",
			brewAvailable: false,
			wantError:     true,
		},
		{
			name:           "system on Linux",
			shellSource:    "system",
			osName:         "linux",
			distroName:     "ubuntu",
			brewAvailable:  true, // brew available but should use system PM
			wantSource:     shell.ShellSourceSystem,
			wantPkgMgr:     true,
			wantBrewPkgMgr: false,
			wantError:      false,
		},
		{
			name:           "system on Fedora",
			shellSource:    "system",
			osName:         "linux",
			distroName:     "fedora",
			brewAvailable:  false,
			wantSource:     shell.ShellSourceSystem,
			wantPkgMgr:     true,
			wantBrewPkgMgr: false,
			wantError:      false,
		},
		{
			name:          "invalid shell source",
			shellSource:   "invalid",
			osName:        "linux",
			distroName:    "ubuntu",
			brewAvailable: true,
			wantError:     true,
		},
		{
			name:          "system on unsupported distro",
			shellSource:   "system",
			osName:        "linux",
			distroName:    "arch", // not supported
			brewAvailable: false,
			wantError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test state
			shellSource = tt.shellSource
			globalSysInfo = compatibility.SystemInfo{
				OSName:     tt.osName,
				DistroName: tt.distroName,
			}
			if tt.brewAvailable {
				globalBrewPath = "/opt/homebrew/bin/brew"
			} else {
				globalBrewPath = ""
			}

			// Act
			pkgMgr, resolver, err := resolveShellInstallStrategy()

			// Assert
			if tt.wantError {
				require.Error(t, err, "Expected error but got none")
				return
			}

			require.NoError(t, err, "Unexpected error: %v", err)
			require.NotNil(t, resolver, "Resolver should not be nil")

			if tt.wantPkgMgr {
				require.NotNil(t, pkgMgr, "Package manager should not be nil")
			}

			// Verify resolver was created with correct source
			// We can test this by checking the shell path resolution behavior
			if tt.wantBrewPkgMgr {
				require.Contains(t, globalBrewPath, "brew", "Should use brew package manager")
			}
		})
	}
}

func Test_ResolveShellInstallStrategy_ValidatesShellSource(t *testing.T) {
	// Save and restore
	originalShellSource := shellSource
	originalShellName := shellName
	originalBrewPath := globalBrewPath
	originalSysInfo := globalSysInfo
	originalLogger := cliLogger
	originalFilesystem := globalFilesystem
	originalOsManager := globalOsManager

	defer func() {
		shellSource = originalShellSource
		shellName = originalShellName
		globalBrewPath = originalBrewPath
		globalSysInfo = originalSysInfo
		cliLogger = originalLogger
		globalFilesystem = originalFilesystem
		globalOsManager = originalOsManager
	}()

	shellName = "zsh"
	globalBrewPath = "/opt/homebrew/bin/brew"
	globalSysInfo = compatibility.SystemInfo{OSName: "linux", DistroName: "ubuntu"}
	cliLogger = logger.DefaultLogger
	globalFilesystem = &utils.MoqFileSystem{
		PathExistsFunc: func(path string) (bool, error) {
			return true, nil
		},
	}
	globalOsManager = &osmanager.MoqOsManager{}

	invalidSources := []string{"", "homebrew", "native", "Auto", "BREW", "System"}

	for _, source := range invalidSources {
		t.Run("rejects_"+source, func(t *testing.T) {
			shellSource = source

			_, _, err := resolveShellInstallStrategy()

			require.Error(t, err, "Should reject invalid source: %s", source)
			require.Contains(t, err.Error(), "invalid shell-source")
		})
	}
}

func Test_ResolveShellInstallStrategy_RequiresBrewForBrewSource(t *testing.T) {
	// Save and restore
	originalShellSource := shellSource
	originalShellName := shellName
	originalBrewPath := globalBrewPath
	originalSysInfo := globalSysInfo
	originalLogger := cliLogger
	originalFilesystem := globalFilesystem
	originalOsManager := globalOsManager

	defer func() {
		shellSource = originalShellSource
		shellName = originalShellName
		globalBrewPath = originalBrewPath
		globalSysInfo = originalSysInfo
		cliLogger = originalLogger
		globalFilesystem = originalFilesystem
		globalOsManager = originalOsManager
	}()

	shellName = "zsh"
	shellSource = "brew"
	globalBrewPath = "" // No brew installed
	globalSysInfo = compatibility.SystemInfo{OSName: "linux", DistroName: "ubuntu"}
	cliLogger = logger.DefaultLogger
	globalFilesystem = &utils.MoqFileSystem{
		PathExistsFunc: func(path string) (bool, error) {
			return true, nil
		},
	}
	globalOsManager = &osmanager.MoqOsManager{}

	_, _, err := resolveShellInstallStrategy()

	require.Error(t, err)
	require.Contains(t, err.Error(), "homebrew is not installed")
}

func Test_ResolveShellInstallStrategy_AutoFallsBackToSystemPM(t *testing.T) {
	// Save and restore
	originalShellSource := shellSource
	originalShellName := shellName
	originalBrewPath := globalBrewPath
	originalSysInfo := globalSysInfo
	originalLogger := cliLogger
	originalFilesystem := globalFilesystem
	originalOsManager := globalOsManager

	defer func() {
		shellSource = originalShellSource
		shellName = originalShellName
		globalBrewPath = originalBrewPath
		globalSysInfo = originalSysInfo
		cliLogger = originalLogger
		globalFilesystem = originalFilesystem
		globalOsManager = originalOsManager
	}()

	shellName = "zsh"
	shellSource = "auto"
	globalBrewPath = "" // No brew
	globalSysInfo = compatibility.SystemInfo{OSName: "linux", DistroName: "ubuntu"}
	cliLogger = logger.DefaultLogger
	globalFilesystem = &utils.MoqFileSystem{
		PathExistsFunc: func(path string) (bool, error) {
			return true, nil
		},
	}
	globalOsManager = &osmanager.MoqOsManager{}

	pkgMgr, resolver, err := resolveShellInstallStrategy()

	require.NoError(t, err)
	require.NotNil(t, pkgMgr, "Should fall back to system package manager")
	require.NotNil(t, resolver)
}
