package brew_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/MrPointer/dotfiles/installer/lib/brew"
	"github.com/MrPointer/dotfiles/installer/lib/pkgmanager"
	"github.com/MrPointer/dotfiles/installer/utils"
	"github.com/MrPointer/dotfiles/installer/utils/logger"
	"github.com/MrPointer/dotfiles/installer/utils/osmanager"
)

func Test_NewBrewPackageManager_ReturnsValidInstance(t *testing.T) {
	mockLogger := &logger.MoqLogger{}
	mockCommander := &utils.MoqCommander{}
	mockProgramQuery := &osmanager.MoqProgramQuery{}

	packageManager := brew.NewBrewPackageManager(mockLogger, mockCommander, mockProgramQuery)

	require.NotNil(t, packageManager)
}

func Test_GetInfo_ReturnsCorrectInfo_WhenBrewVersionIsRetrieved(t *testing.T) {
	mockLogger := &logger.MoqLogger{}
	mockCommander := &utils.MoqCommander{}
	mockProgramQuery := &osmanager.MoqProgramQuery{
		GetProgramVersionFunc: func(program string, versionExtractor osmanager.VersionExtractor, queryArgs ...string) (string, error) {
			if program == "brew" {
				return "3.6.5", nil
			}
			return "", errors.New("unexpected program")
		},
	}

	packageManager := brew.NewBrewPackageManager(mockLogger, mockCommander, mockProgramQuery)

	info, err := packageManager.GetInfo()

	require.NoError(t, err)
	require.Equal(t, "brew", info.Name)
	require.Equal(t, "3.6.5", info.Version)
}

func Test_GetInfo_HandlesBrewVersionExtractionCorrectly(t *testing.T) {
	tests := []struct {
		name            string
		rawVersion      string
		expectedVersion string
		expectedError   bool
	}{
		{
			name:            "Standard brew version format",
			rawVersion:      "Homebrew 3.6.5",
			expectedVersion: "3.6.5",
			expectedError:   false,
		},
		{
			name:            "Empty version string",
			rawVersion:      "",
			expectedVersion: "",
			expectedError:   false,
		},
		{
			name:            "Invalid version format",
			rawVersion:      "InvalidFormat",
			expectedVersion: "",
			expectedError:   true,
		},
		{
			name:            "Version with additional components",
			rawVersion:      "Homebrew 3.6.5-beta (extra info)",
			expectedVersion: "3.6.5-beta",
			expectedError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger := &logger.MoqLogger{}
			mockCommander := &utils.MoqCommander{}
			mockProgramQuery := &osmanager.MoqProgramQuery{
				GetProgramVersionFunc: func(program string, versionExtractor osmanager.VersionExtractor, queryArgs ...string) (string, error) {
					return versionExtractor(tt.rawVersion)
				},
			}

			packageManager := brew.NewBrewPackageManager(mockLogger, mockCommander, mockProgramQuery)

			info, err := packageManager.GetInfo()

			if tt.expectedError {
				require.Error(t, err)
				require.Contains(t, err.Error(), "Homebrew version")
			} else {
				require.NoError(t, err)
				require.Equal(t, "brew", info.Name)
				require.Equal(t, tt.expectedVersion, info.Version)
			}
		})
	}
}

func Test_GetInfo_ReturnsError_WhenProgramQueryFails(t *testing.T) {
	mockLogger := &logger.MoqLogger{}
	mockCommander := &utils.MoqCommander{}
	mockProgramQuery := &osmanager.MoqProgramQuery{
		GetProgramVersionFunc: func(program string, versionExtractor osmanager.VersionExtractor, queryArgs ...string) (string, error) {
			return "", errors.New("program not found")
		},
	}

	packageManager := brew.NewBrewPackageManager(mockLogger, mockCommander, mockProgramQuery)

	info, err := packageManager.GetInfo()

	require.Error(t, err)
	require.Contains(t, err.Error(), "Homebrew version")
	require.Equal(t, pkgmanager.DefaultPackageManagerInfo(), info)
}

func Test_GetPackageVersion_ReturnsVersion_WhenPackageIsInstalled(t *testing.T) {
	mockLogger := &logger.MoqLogger{}
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			if name == "brew" && len(args) == 2 && args[0] == "list" && args[1] == "--versions" {
				output := "git 2.39.0\nnode 18.12.1\nvim 9.0.0500"
				return &utils.Result{
					Stdout: []byte(output),
				}, nil
			}
			return nil, errors.New("unexpected command")
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{}

	packageManager := brew.NewBrewPackageManager(mockLogger, mockCommander, mockProgramQuery)

	version, err := packageManager.GetPackageVersion("node")

	require.NoError(t, err)
	require.Equal(t, "18.12.1", version)
}

func Test_GetPackageVersion_ReturnsError_WhenPackageIsNotInstalled(t *testing.T) {
	mockLogger := &logger.MoqLogger{}
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			if name == "brew" && len(args) == 2 && args[0] == "list" && args[1] == "--versions" {
				output := "git 2.39.0\nvim 9.0.0500"
				return &utils.Result{
					Stdout: []byte(output),
				}, nil
			}
			return nil, errors.New("unexpected command")
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{}

	packageManager := brew.NewBrewPackageManager(mockLogger, mockCommander, mockProgramQuery)

	version, err := packageManager.GetPackageVersion("nonexistent")

	require.Error(t, err)
	require.Contains(t, err.Error(), "not installed")
	require.Empty(t, version)
}

func Test_GetPackageVersion_ReturnsError_WhenListInstalledPackagesFails(t *testing.T) {
	mockLogger := &logger.MoqLogger{}
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			return nil, errors.New("command failed")
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{}

	packageManager := brew.NewBrewPackageManager(mockLogger, mockCommander, mockProgramQuery)

	version, err := packageManager.GetPackageVersion("git")

	require.Error(t, err)
	require.Contains(t, err.Error(), "list installed packages")
	require.Empty(t, version)
}

func Test_InstallPackage_InstallsPackageSuccessfully(t *testing.T) {
	expectedWarning := ""
	mockLogger := &logger.MoqLogger{
		WarningFunc: func(format string, args ...any) {
			expectedWarning = format
		},
	}
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			if name == "brew" && len(args) == 2 && args[0] == "install" && args[1] == "git" {
				return &utils.Result{}, nil
			}
			return nil, errors.New("unexpected command")
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{}

	packageManager := brew.NewBrewPackageManager(mockLogger, mockCommander, mockProgramQuery)
	requestedPackage := pkgmanager.RequestedPackageInfo{Name: "git"}

	err := packageManager.InstallPackage(requestedPackage)

	require.NoError(t, err)
	require.Contains(t, expectedWarning, "Homebrew doesn't support version constraints")
}

func Test_InstallPackage_ReturnsError_WhenInstallationFails(t *testing.T) {
	mockLogger := &logger.MoqLogger{
		WarningFunc: func(format string, args ...any) {},
	}
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			return nil, errors.New("installation failed")
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{}

	packageManager := brew.NewBrewPackageManager(mockLogger, mockCommander, mockProgramQuery)
	requestedPackage := pkgmanager.RequestedPackageInfo{Name: "git"}

	err := packageManager.InstallPackage(requestedPackage)

	require.Error(t, err)
	require.Contains(t, err.Error(), "install package")
}

func Test_IsPackageInstalled_ReturnsTrue_WhenPackageIsInstalled(t *testing.T) {
	mockLogger := &logger.MoqLogger{}
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			if name == "brew" && len(args) == 2 && args[0] == "list" && args[1] == "--versions" {
				output := "git 2.39.0\nnode 18.12.1\nvim 9.0.0500"
				return &utils.Result{
					Stdout: []byte(output),
				}, nil
			}
			return nil, errors.New("unexpected command")
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{}

	packageManager := brew.NewBrewPackageManager(mockLogger, mockCommander, mockProgramQuery)
	packageInfo := pkgmanager.NewPackageInfo("git", "2.39.0")

	isInstalled, err := packageManager.IsPackageInstalled(packageInfo)

	require.NoError(t, err)
	require.True(t, isInstalled)
}

func Test_IsPackageInstalled_ReturnsFalse_WhenPackageIsNotInstalled(t *testing.T) {
	mockLogger := &logger.MoqLogger{}
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			if name == "brew" && len(args) == 2 && args[0] == "list" && args[1] == "--versions" {
				output := "git 2.39.0\nvim 9.0.0500"
				return &utils.Result{
					Stdout: []byte(output),
				}, nil
			}
			return nil, errors.New("unexpected command")
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{}

	packageManager := brew.NewBrewPackageManager(mockLogger, mockCommander, mockProgramQuery)
	packageInfo := pkgmanager.NewPackageInfo("nonexistent", "1.0.0")

	isInstalled, err := packageManager.IsPackageInstalled(packageInfo)

	require.NoError(t, err)
	require.False(t, isInstalled)
}

func Test_IsPackageInstalled_ReturnsError_WhenListInstalledPackagesFails(t *testing.T) {
	mockLogger := &logger.MoqLogger{}
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			return nil, errors.New("command failed")
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{}

	packageManager := brew.NewBrewPackageManager(mockLogger, mockCommander, mockProgramQuery)
	packageInfo := pkgmanager.NewPackageInfo("git", "2.39.0")

	isInstalled, err := packageManager.IsPackageInstalled(packageInfo)

	require.Error(t, err)
	require.Contains(t, err.Error(), "list installed packages")
	require.False(t, isInstalled)
}

func Test_ListInstalledPackages_ReturnsPackageList_WhenCommandSucceeds(t *testing.T) {
	mockLogger := &logger.MoqLogger{}
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			if name == "brew" && len(args) == 2 && args[0] == "list" && args[1] == "--versions" {
				output := "git 2.39.0\nnode 18.12.1\nvim 9.0.0500"
				return &utils.Result{
					Stdout: []byte(output),
				}, nil
			}
			return nil, errors.New("unexpected command")
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{}

	packageManager := brew.NewBrewPackageManager(mockLogger, mockCommander, mockProgramQuery)

	packages, err := packageManager.ListInstalledPackages()

	require.NoError(t, err)
	require.Len(t, packages, 3)
	require.Equal(t, "git", packages[0].Name)
	require.Equal(t, "2.39.0", packages[0].Version)
	require.Equal(t, "node", packages[1].Name)
	require.Equal(t, "18.12.1", packages[1].Version)
	require.Equal(t, "vim", packages[2].Name)
	require.Equal(t, "9.0.0500", packages[2].Version)
}

func Test_ListInstalledPackages_ReturnsError_WhenCommandFails(t *testing.T) {
	mockLogger := &logger.MoqLogger{}
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			return nil, errors.New("command failed")
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{}

	packageManager := brew.NewBrewPackageManager(mockLogger, mockCommander, mockProgramQuery)

	packages, err := packageManager.ListInstalledPackages()

	require.Error(t, err)
	require.Contains(t, err.Error(), "list installed packages")
	require.Nil(t, packages)
}

func Test_ListInstalledPackages_ReturnsError_WhenOutputFormatIsInvalid(t *testing.T) {
	mockLogger := &logger.MoqLogger{}
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			if name == "brew" && len(args) == 2 && args[0] == "list" && args[1] == "--versions" {
				output := "git\ninvalid-line-without-version\nvim 9.0.0500"
				return &utils.Result{
					Stdout: []byte(output),
				}, nil
			}
			return nil, errors.New("unexpected command")
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{}

	packageManager := brew.NewBrewPackageManager(mockLogger, mockCommander, mockProgramQuery)

	packages, err := packageManager.ListInstalledPackages()

	require.Error(t, err)
	require.Contains(t, err.Error(), "parse package")
	require.Nil(t, packages)
}

func Test_ListInstalledPackages_HandlesEmptyOutput(t *testing.T) {
	mockLogger := &logger.MoqLogger{}
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			if name == "brew" && len(args) == 2 && args[0] == "list" && args[1] == "--versions" {
				return &utils.Result{
					Stdout: []byte(""),
				}, nil
			}
			return nil, errors.New("unexpected command")
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{}

	packageManager := brew.NewBrewPackageManager(mockLogger, mockCommander, mockProgramQuery)

	packages, err := packageManager.ListInstalledPackages()

	require.NoError(t, err)
	require.Empty(t, packages)
}

func Test_UninstallPackage_UninstallsPackageSuccessfully(t *testing.T) {
	mockLogger := &logger.MoqLogger{}
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			if name == "brew" && len(args) == 2 && args[0] == "uninstall" && args[1] == "git" {
				return &utils.Result{}, nil
			}
			return nil, errors.New("unexpected command")
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{}

	packageManager := brew.NewBrewPackageManager(mockLogger, mockCommander, mockProgramQuery)
	packageInfo := pkgmanager.NewPackageInfo("git", "2.39.0")

	err := packageManager.UninstallPackage(packageInfo)

	require.NoError(t, err)
}

func Test_UninstallPackage_ReturnsError_WhenUninstallationFails(t *testing.T) {
	mockLogger := &logger.MoqLogger{}
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			return nil, errors.New("uninstallation failed")
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{}

	packageManager := brew.NewBrewPackageManager(mockLogger, mockCommander, mockProgramQuery)
	packageInfo := pkgmanager.NewPackageInfo("git", "2.39.0")

	err := packageManager.UninstallPackage(packageInfo)

	require.Error(t, err)
	require.Contains(t, err.Error(), "uninstall package")
}

func Test_ListInstalledPackages_HandlesPackageNamesWithVersionsInMultipleFormats(t *testing.T) {
	tests := []struct {
		name          string
		brewOutput    string
		expectedCount int
		expectedFirst pkgmanager.PackageInfo
		expectError   bool
	}{
		{
			name:          "Standard package format",
			brewOutput:    "git 2.39.0\nnode 18.12.1",
			expectedCount: 2,
			expectedFirst: pkgmanager.NewPackageInfo("git", "2.39.0"),
			expectError:   false,
		},
		{
			name:          "Package with complex version",
			brewOutput:    "postgresql@14 14.6_1\nredis 7.0.5",
			expectedCount: 2,
			expectedFirst: pkgmanager.NewPackageInfo("postgresql@14", "14.6_1"),
			expectError:   false,
		},
		{
			name:          "Single package",
			brewOutput:    "vim 9.0.0500",
			expectedCount: 1,
			expectedFirst: pkgmanager.NewPackageInfo("vim", "9.0.0500"),
			expectError:   false,
		},
		{
			name:          "Package with whitespace-containing version",
			brewOutput:    "package-name 1.0.0 2.0.0",
			expectedCount: 1,
			expectedFirst: pkgmanager.NewPackageInfo("package-name", "1.0.0 2.0.0"),
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger := &logger.MoqLogger{}
			mockCommander := &utils.MoqCommander{
				RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
					if name == "brew" && len(args) == 2 && args[0] == "list" && args[1] == "--versions" {
						return &utils.Result{
							Stdout: []byte(tt.brewOutput),
						}, nil
					}
					return nil, errors.New("unexpected command")
				},
			}
			mockProgramQuery := &osmanager.MoqProgramQuery{}

			packageManager := brew.NewBrewPackageManager(mockLogger, mockCommander, mockProgramQuery)

			packages, err := packageManager.ListInstalledPackages()

			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Len(t, packages, tt.expectedCount)
				if tt.expectedCount > 0 {
					require.Equal(t, tt.expectedFirst.Name, packages[0].Name)
					require.Equal(t, tt.expectedFirst.Version, packages[0].Version)
				}
			}
		})
	}
}
