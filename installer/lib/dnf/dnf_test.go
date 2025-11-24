package dnf_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/MrPointer/dotfiles/installer/lib/dnf"
	"github.com/MrPointer/dotfiles/installer/lib/pkgmanager"
	"github.com/MrPointer/dotfiles/installer/utils"
	"github.com/MrPointer/dotfiles/installer/utils/logger"
	"github.com/MrPointer/dotfiles/installer/utils/osmanager"
	"github.com/MrPointer/dotfiles/installer/utils/privilege"
)

func Test_DnfPackageManager_ImplementsPackageManagerInterface(t *testing.T) {
	mockCommander := &utils.MoqCommander{}
	mockProgramQuery := &osmanager.MoqProgramQuery{}
	mockEscalator := &privilege.MoqEscalator{}

	dnfManager := dnf.NewDnfPackageManager(logger.DefaultLogger, mockCommander, mockProgramQuery, mockEscalator, utils.DisplayModeProgress)

	require.Implements(t, (*pkgmanager.PackageManager)(nil), dnfManager)
}

func Test_NewDnfPackageManager_ReturnsValidInstance(t *testing.T) {
	mockCommander := &utils.MoqCommander{}
	mockProgramQuery := &osmanager.MoqProgramQuery{}
	mockEscalator := &privilege.MoqEscalator{}

	dnfManager := dnf.NewDnfPackageManager(logger.DefaultLogger, mockCommander, mockProgramQuery, mockEscalator, utils.DisplayModeProgress)

	require.NotNil(t, dnfManager)
}

func Test_GetInfo_ReturnsValidDnfManagerInfo(t *testing.T) {
	mockCommander := &utils.MoqCommander{}
	mockProgramQuery := &osmanager.MoqProgramQuery{
		GetProgramVersionFunc: func(program string, versionExtractor osmanager.VersionExtractor, queryArgs ...string) (string, error) {
			if program == "dnf" {
				return versionExtractor("dnf 4.14.0")
			}
			return "", nil
		},
	}
	mockEscalator := &privilege.MoqEscalator{}

	dnfManager := dnf.NewDnfPackageManager(logger.DefaultLogger, mockCommander, mockProgramQuery, mockEscalator, utils.DisplayModeProgress)

	info, err := dnfManager.GetInfo()

	require.NoError(t, err)
	require.Equal(t, "dnf", info.Name)
	require.Equal(t, "4.14.0", info.Version)
}

func Test_GetInfo_ReturnsError_WhenVersionQueryFails(t *testing.T) {
	mockCommander := &utils.MoqCommander{}
	mockProgramQuery := &osmanager.MoqProgramQuery{
		GetProgramVersionFunc: func(program string, versionExtractor osmanager.VersionExtractor, queryArgs ...string) (string, error) {
			return "", errors.New("command not found")
		},
	}
	mockEscalator := &privilege.MoqEscalator{}

	dnfManager := dnf.NewDnfPackageManager(logger.DefaultLogger, mockCommander, mockProgramQuery, mockEscalator, utils.DisplayModeProgress)

	info, err := dnfManager.GetInfo()

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to get DNF version")
	require.Equal(t, pkgmanager.DefaultPackageManagerInfo(), info)
}

func Test_InstallPackage_InstallsRegularPackage_Successfully(t *testing.T) {
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(command string, args []string, options ...utils.Option) (*utils.Result, error) {
			require.Equal(t, "sudo", command)
			require.Equal(t, []string{"dnf", "install", "-y", "git"}, args)
			return &utils.Result{}, nil
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{}
	mockEscalator := &privilege.MoqEscalator{
		EscalateCommandFunc: func(command string, args []string) (privilege.EscalationResult, error) {
			return privilege.EscalationResult{Command: "sudo", Args: append([]string{"dnf"}, args...)}, nil
		},
	}

	dnfManager := dnf.NewDnfPackageManager(logger.DefaultLogger, mockCommander, mockProgramQuery, mockEscalator, utils.DisplayModeProgress)

	err := dnfManager.InstallPackage(pkgmanager.NewRequestedPackageInfo("git", nil))

	require.NoError(t, err)
	require.Len(t, mockCommander.RunCommandCalls(), 1)
}

func Test_InstallPackage_InstallsGroupPackage_Successfully(t *testing.T) {
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(command string, args []string, options ...utils.Option) (*utils.Result, error) {
			require.Equal(t, "sudo", command)
			require.Equal(t, []string{"dnf", "group", "install", "-y", "Development Tools"}, args)
			return &utils.Result{}, nil
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{}
	mockEscalator := &privilege.MoqEscalator{
		EscalateCommandFunc: func(command string, args []string) (privilege.EscalationResult, error) {
			return privilege.EscalationResult{Command: "sudo", Args: append([]string{"dnf"}, args...)}, nil
		},
	}

	dnfManager := dnf.NewDnfPackageManager(logger.DefaultLogger, mockCommander, mockProgramQuery, mockEscalator, utils.DisplayModeProgress)

	err := dnfManager.InstallPackage(pkgmanager.NewRequestedPackageInfoWithType("Development Tools", "group", nil))

	require.NoError(t, err)
	require.Len(t, mockCommander.RunCommandCalls(), 1)
}

func Test_InstallPackage_ReturnsError_WhenEscalationFails(t *testing.T) {
	mockCommander := &utils.MoqCommander{}
	mockProgramQuery := &osmanager.MoqProgramQuery{}
	mockEscalator := &privilege.MoqEscalator{
		EscalateCommandFunc: func(command string, args []string) (privilege.EscalationResult, error) {
			return privilege.EscalationResult{}, errors.New("escalation failed")
		},
	}

	dnfManager := dnf.NewDnfPackageManager(logger.DefaultLogger, mockCommander, mockProgramQuery, mockEscalator, utils.DisplayModeProgress)

	err := dnfManager.InstallPackage(pkgmanager.NewRequestedPackageInfo("git", nil))

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to determine privilege escalation for dnf install")
}

func Test_InstallPackage_ReturnsError_WhenCommandFails(t *testing.T) {
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(command string, args []string, options ...utils.Option) (*utils.Result, error) {
			return nil, errors.New("package not found")
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{}
	mockEscalator := &privilege.MoqEscalator{
		EscalateCommandFunc: func(command string, args []string) (privilege.EscalationResult, error) {
			return privilege.EscalationResult{Command: "sudo", Args: append([]string{"dnf"}, args...)}, nil
		},
	}

	dnfManager := dnf.NewDnfPackageManager(logger.DefaultLogger, mockCommander, mockProgramQuery, mockEscalator, utils.DisplayModeProgress)

	err := dnfManager.InstallPackage(pkgmanager.NewRequestedPackageInfo("nonexistent", nil))

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to install package nonexistent")
}

func Test_UninstallPackage_UninstallsRegularPackage_Successfully(t *testing.T) {
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(command string, args []string, options ...utils.Option) (*utils.Result, error) {
			require.Equal(t, "sudo", command)
			require.Equal(t, []string{"dnf", "remove", "-y", "git"}, args)
			return &utils.Result{}, nil
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{}
	mockEscalator := &privilege.MoqEscalator{
		EscalateCommandFunc: func(command string, args []string) (privilege.EscalationResult, error) {
			return privilege.EscalationResult{Command: "sudo", Args: append([]string{"dnf"}, args...)}, nil
		},
	}

	dnfManager := dnf.NewDnfPackageManager(logger.DefaultLogger, mockCommander, mockProgramQuery, mockEscalator, utils.DisplayModeProgress)

	err := dnfManager.UninstallPackage(pkgmanager.NewPackageInfo("git", "2.39.0"))

	require.NoError(t, err)
	require.Len(t, mockCommander.RunCommandCalls(), 1)
}

func Test_UninstallPackage_UninstallsGroupPackage_Successfully(t *testing.T) {
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(command string, args []string, options ...utils.Option) (*utils.Result, error) {
			require.Equal(t, "sudo", command)
			require.Equal(t, []string{"dnf", "group", "remove", "-y", "Development Tools"}, args)
			return &utils.Result{}, nil
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{}
	mockEscalator := &privilege.MoqEscalator{
		EscalateCommandFunc: func(command string, args []string) (privilege.EscalationResult, error) {
			return privilege.EscalationResult{Command: "sudo", Args: append([]string{"dnf"}, args...)}, nil
		},
	}

	dnfManager := dnf.NewDnfPackageManager(logger.DefaultLogger, mockCommander, mockProgramQuery, mockEscalator, utils.DisplayModeProgress)

	err := dnfManager.UninstallPackage(pkgmanager.NewPackageInfoWithType("Development Tools", "1.0", "group"))

	require.NoError(t, err)
	require.Len(t, mockCommander.RunCommandCalls(), 1)
}

func Test_IsPackageInstalled_ReturnsTrueForInstalledRegularPackage(t *testing.T) {
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(command string, args []string, options ...utils.Option) (*utils.Result, error) {
			require.Equal(t, "dnf", command)
			require.Equal(t, []string{"list", "installed"}, args)
			return &utils.Result{
				Stdout: []byte("Installed Packages\ngit.x86_64                    2.39.0-1.fc38                    @fedora\nvim.x86_64                    9.0.1160-1.fc38                  @fedora"),
			}, nil
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{}
	mockEscalator := &privilege.MoqEscalator{}

	dnfManager := dnf.NewDnfPackageManager(logger.DefaultLogger, mockCommander, mockProgramQuery, mockEscalator, utils.DisplayModeProgress)

	isInstalled, err := dnfManager.IsPackageInstalled(pkgmanager.NewPackageInfo("git", ""))

	require.NoError(t, err)
	require.True(t, isInstalled)
}

func Test_IsPackageInstalled_ReturnsFalseForNotInstalledRegularPackage(t *testing.T) {
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(command string, args []string, options ...utils.Option) (*utils.Result, error) {
			require.Equal(t, "dnf", command)
			require.Equal(t, []string{"list", "installed"}, args)
			return &utils.Result{
				Stdout: []byte("Installed Packages\nvim.x86_64                    9.0.1160-1.fc38                  @fedora"),
			}, nil
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{}
	mockEscalator := &privilege.MoqEscalator{}

	dnfManager := dnf.NewDnfPackageManager(logger.DefaultLogger, mockCommander, mockProgramQuery, mockEscalator, utils.DisplayModeProgress)

	isInstalled, err := dnfManager.IsPackageInstalled(pkgmanager.NewPackageInfo("git", ""))

	require.NoError(t, err)
	require.False(t, isInstalled)
}

func Test_IsPackageInstalled_ReturnsTrueForInstalledGroup(t *testing.T) {
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(command string, args []string, options ...utils.Option) (*utils.Result, error) {
			require.Equal(t, "dnf", command)
			require.Equal(t, []string{"group", "list", "installed"}, args)
			return &utils.Result{
				Stdout: []byte("Available Environment Groups:\n   Development Tools\n   Server Tools"),
			}, nil
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{}
	mockEscalator := &privilege.MoqEscalator{}

	dnfManager := dnf.NewDnfPackageManager(logger.DefaultLogger, mockCommander, mockProgramQuery, mockEscalator, utils.DisplayModeProgress)

	isInstalled, err := dnfManager.IsPackageInstalled(pkgmanager.NewPackageInfoWithType("Development Tools", "", "group"))

	require.NoError(t, err)
	require.True(t, isInstalled)
}

func Test_IsPackageInstalled_ReturnsFalseForNotInstalledGroup(t *testing.T) {
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(command string, args []string, options ...utils.Option) (*utils.Result, error) {
			require.Equal(t, "dnf", command)
			require.Equal(t, []string{"group", "list", "installed"}, args)
			return &utils.Result{
				Stdout: []byte("Available Environment Groups:\n   Server Tools"),
			}, nil
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{}
	mockEscalator := &privilege.MoqEscalator{}

	dnfManager := dnf.NewDnfPackageManager(logger.DefaultLogger, mockCommander, mockProgramQuery, mockEscalator, utils.DisplayModeProgress)

	isInstalled, err := dnfManager.IsPackageInstalled(pkgmanager.NewPackageInfoWithType("Development Tools", "", "group"))

	require.NoError(t, err)
	require.False(t, isInstalled)
}

func Test_ListInstalledPackages_ReturnsValidPackageList(t *testing.T) {
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(command string, args []string, options ...utils.Option) (*utils.Result, error) {
			return &utils.Result{
				Stdout: []byte("Installed Packages\ngit.x86_64                    2.39.0-1.fc38                    @fedora\nvim.x86_64                    9.0.1160-1.fc38                  @fedora\nzsh.x86_64                    5.9-1.fc38                       @fedora"),
			}, nil
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{}
	mockEscalator := &privilege.MoqEscalator{}

	dnfManager := dnf.NewDnfPackageManager(logger.DefaultLogger, mockCommander, mockProgramQuery, mockEscalator, utils.DisplayModeProgress)

	packages, err := dnfManager.ListInstalledPackages()

	require.NoError(t, err)
	require.Len(t, packages, 3)

	// Check that packages are parsed correctly
	expectedPackages := []pkgmanager.PackageInfo{
		{Name: "git", Version: "2.39.0-1.fc38"},
		{Name: "vim", Version: "9.0.1160-1.fc38"},
		{Name: "zsh", Version: "5.9-1.fc38"},
	}

	for i, expected := range expectedPackages {
		require.Equal(t, expected.Name, packages[i].Name)
		require.Equal(t, expected.Version, packages[i].Version)
	}
}

func Test_ListInstalledPackages_ReturnsEmptyList_WhenNoPackagesInstalled(t *testing.T) {
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(command string, args []string, options ...utils.Option) (*utils.Result, error) {
			return &utils.Result{
				Stdout: []byte("Installed Packages\n"),
			}, nil
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{}
	mockEscalator := &privilege.MoqEscalator{}

	dnfManager := dnf.NewDnfPackageManager(logger.DefaultLogger, mockCommander, mockProgramQuery, mockEscalator, utils.DisplayModeProgress)

	packages, err := dnfManager.ListInstalledPackages()

	require.NoError(t, err)
	require.Empty(t, packages)
}

func Test_ListInstalledPackages_ReturnsError_WhenCommandFails(t *testing.T) {
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(command string, args []string, options ...utils.Option) (*utils.Result, error) {
			return nil, errors.New("dnf command failed")
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{}
	mockEscalator := &privilege.MoqEscalator{}

	dnfManager := dnf.NewDnfPackageManager(logger.DefaultLogger, mockCommander, mockProgramQuery, mockEscalator, utils.DisplayModeProgress)

	packages, err := dnfManager.ListInstalledPackages()

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to list installed packages")
	require.Nil(t, packages)
}

func Test_GetPackageVersion_ReturnsCorrectVersion_WhenPackageInstalled(t *testing.T) {
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(command string, args []string, options ...utils.Option) (*utils.Result, error) {
			return &utils.Result{
				Stdout: []byte("Installed Packages\ngit.x86_64                    2.39.0-1.fc38                    @fedora\nvim.x86_64                    9.0.1160-1.fc38                  @fedora"),
			}, nil
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{}
	mockEscalator := &privilege.MoqEscalator{}

	dnfManager := dnf.NewDnfPackageManager(logger.DefaultLogger, mockCommander, mockProgramQuery, mockEscalator, utils.DisplayModeProgress)

	version, err := dnfManager.GetPackageVersion("git")

	require.NoError(t, err)
	require.Equal(t, "2.39.0-1.fc38", version)
}

func Test_GetPackageVersion_ReturnsError_WhenPackageNotInstalled(t *testing.T) {
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(command string, args []string, options ...utils.Option) (*utils.Result, error) {
			return &utils.Result{
				Stdout: []byte("Installed Packages\nvim.x86_64                    9.0.1160-1.fc38                  @fedora"),
			}, nil
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{}
	mockEscalator := &privilege.MoqEscalator{}

	dnfManager := dnf.NewDnfPackageManager(logger.DefaultLogger, mockCommander, mockProgramQuery, mockEscalator, utils.DisplayModeProgress)

	version, err := dnfManager.GetPackageVersion("git")

	require.Error(t, err)
	require.Contains(t, err.Error(), "package git is not installed")
	require.Empty(t, version)
}

func Test_InstallPackage_DiscardsOutput_WhenDisplayModeRequiresIt(t *testing.T) {
	var capturedOptions []utils.Option
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(command string, args []string, options ...utils.Option) (*utils.Result, error) {
			capturedOptions = options
			return &utils.Result{}, nil
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{}
	mockEscalator := &privilege.MoqEscalator{
		EscalateCommandFunc: func(command string, args []string) (privilege.EscalationResult, error) {
			return privilege.EscalationResult{Command: "sudo", Args: append([]string{"dnf"}, args...)}, nil
		},
	}

	dnfManager := dnf.NewDnfPackageManager(logger.DefaultLogger, mockCommander, mockProgramQuery, mockEscalator, utils.DisplayModeProgress)

	err := dnfManager.InstallPackage(pkgmanager.NewRequestedPackageInfo("git", nil))

	require.NoError(t, err)
	require.Len(t, capturedOptions, 1)
	// We can't directly test the option type, but we can verify it was passed
}

func Test_UninstallPackage_DiscardsOutput_WhenDisplayModeRequiresIt(t *testing.T) {
	var capturedOptions []utils.Option
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(command string, args []string, options ...utils.Option) (*utils.Result, error) {
			capturedOptions = options
			return &utils.Result{}, nil
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{}
	mockEscalator := &privilege.MoqEscalator{
		EscalateCommandFunc: func(command string, args []string) (privilege.EscalationResult, error) {
			return privilege.EscalationResult{Command: "sudo", Args: append([]string{"dnf"}, args...)}, nil
		},
	}

	dnfManager := dnf.NewDnfPackageManager(logger.DefaultLogger, mockCommander, mockProgramQuery, mockEscalator, utils.DisplayModeProgress)

	err := dnfManager.UninstallPackage(pkgmanager.NewPackageInfo("git", "2.39.0"))

	require.NoError(t, err)
	require.Len(t, capturedOptions, 1)
	// We can't directly test the option type, but we can verify it was passed
}
