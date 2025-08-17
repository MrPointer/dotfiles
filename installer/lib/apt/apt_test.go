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

func Test_AptPackageManager_ImplementsPackageManagerInterface(t *testing.T) {
	mockCommander := &utils.MoqCommander{}
	mockProgramQuery := &osmanager.MoqProgramQuery{}
	mockEscalator := &privilege.MoqEscalator{}

	aptManager := apt.NewAptPackageManager(logger.DefaultLogger, mockCommander, mockProgramQuery, mockEscalator)

	require.Implements(t, (*pkgmanager.PackageManager)(nil), aptManager)
}

func Test_NewAptPackageManager_ReturnsValidInstance(t *testing.T) {
	mockCommander := &utils.MoqCommander{}
	mockProgramQuery := &osmanager.MoqProgramQuery{}
	mockEscalator := &privilege.MoqEscalator{}

	aptManager := apt.NewAptPackageManager(logger.DefaultLogger, mockCommander, mockProgramQuery, mockEscalator)

	require.NotNil(t, aptManager)
}

func Test_GetInfo_ReturnsAptManagerInfo(t *testing.T) {
	mockCommander := &utils.MoqCommander{}
	mockProgramQuery := &osmanager.MoqProgramQuery{
		GetProgramVersionFunc: func(program string, versionExtractor osmanager.VersionExtractor, queryArgs ...string) (string, error) {
			if program == "apt" {
				return versionExtractor("apt 2.4.8 (amd64)")
			}
			return "", nil
		},
	}
	mockEscalator := &privilege.MoqEscalator{}

	aptManager := apt.NewAptPackageManager(logger.DefaultLogger, mockCommander, mockProgramQuery, mockEscalator)

	info, err := aptManager.GetInfo()

	require.NoError(t, err)
	require.Equal(t, "apt", info.Name)
	require.Equal(t, "2.4.8", info.Version)
}

func Test_InstallPackage_CallsAptInstallCommand_AsRoot(t *testing.T) {
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			return &utils.Result{}, nil
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{}
	mockEscalator := &privilege.MoqEscalator{
		EscalateCommandFunc: func(baseCmd string, baseArgs []string) (privilege.EscalationResult, error) {
			return privilege.EscalationResult{
				Method:          privilege.EscalationNone,
				Command:         baseCmd,
				Args:            baseArgs,
				NeedsEscalation: false,
			}, nil
		},
	}

	aptManager := apt.NewAptPackageManager(logger.DefaultLogger, mockCommander, mockProgramQuery, mockEscalator)
	packageInfo := pkgmanager.NewRequestedPackageInfo("git", nil)

	err := aptManager.InstallPackage(packageInfo)

	require.NoError(t, err)

	calls := mockCommander.RunCommandCalls()
	require.GreaterOrEqual(t, len(calls), 2)

	// Find apt calls - should be direct (no sudo) since running as root
	var updateCall, installCall *struct {
		Name string
		Args []string
		Opts []utils.Option
	}

	for _, call := range calls {
		if call.Name == "apt" && len(call.Args) >= 1 && call.Args[0] == "update" {
			updateCall = &call
		}
		if call.Name == "apt" && len(call.Args) >= 3 && call.Args[0] == "install" && call.Args[1] == "-y" {
			installCall = &call
		}
	}

	require.NotNil(t, updateCall, "apt update call not found")
	require.NotNil(t, installCall, "apt install call not found")
}

func Test_InstallPackage_CallsAptInstallCommand_WithSudo(t *testing.T) {
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			return &utils.Result{}, nil
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{}
	mockEscalator := &privilege.MoqEscalator{
		EscalateCommandFunc: func(baseCmd string, baseArgs []string) (privilege.EscalationResult, error) {
			return privilege.EscalationResult{
				Method:          privilege.EscalationSudo,
				Command:         "sudo",
				Args:            append([]string{baseCmd}, baseArgs...),
				NeedsEscalation: true,
			}, nil
		},
	}

	aptManager := apt.NewAptPackageManager(logger.DefaultLogger, mockCommander, mockProgramQuery, mockEscalator)
	packageInfo := pkgmanager.NewRequestedPackageInfo("git", nil)

	err := aptManager.InstallPackage(packageInfo)

	require.NoError(t, err)

	calls := mockCommander.RunCommandCalls()
	require.GreaterOrEqual(t, len(calls), 2)

	// Find sudo apt calls
	var updateCall, installCall *struct {
		Name string
		Args []string
		Opts []utils.Option
	}

	for _, call := range calls {
		if call.Name == "sudo" && len(call.Args) >= 2 && call.Args[0] == "apt" && call.Args[1] == "update" {
			updateCall = &call
		}
		if call.Name == "sudo" && len(call.Args) >= 4 && call.Args[0] == "apt" && call.Args[1] == "install" && call.Args[2] == "-y" {
			installCall = &call
		}
	}

	require.NotNil(t, updateCall, "sudo apt update call not found")
	require.NotNil(t, installCall, "sudo apt install call not found")
}

func Test_InstallPackage_CallsAptInstallCommand_WithoutPrivilegeEscalation(t *testing.T) {
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			return &utils.Result{}, nil
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{}
	mockEscalator := &privilege.MoqEscalator{
		EscalateCommandFunc: func(baseCmd string, baseArgs []string) (privilege.EscalationResult, error) {
			return privilege.EscalationResult{
				Method:          privilege.EscalationDirect,
				Command:         baseCmd,
				Args:            baseArgs,
				NeedsEscalation: false,
			}, nil
		},
	}

	aptManager := apt.NewAptPackageManager(logger.DefaultLogger, mockCommander, mockProgramQuery, mockEscalator)
	packageInfo := pkgmanager.NewRequestedPackageInfo("git", nil)

	err := aptManager.InstallPackage(packageInfo)

	require.NoError(t, err)

	calls := mockCommander.RunCommandCalls()
	require.GreaterOrEqual(t, len(calls), 2)

	// Find direct apt calls (no privilege escalation)
	var updateCall, installCall *struct {
		Name string
		Args []string
		Opts []utils.Option
	}

	for _, call := range calls {
		if call.Name == "apt" && len(call.Args) >= 1 && call.Args[0] == "update" {
			updateCall = &call
		}
		if call.Name == "apt" && len(call.Args) >= 3 && call.Args[0] == "install" && call.Args[1] == "-y" {
			installCall = &call
		}
	}

	require.NotNil(t, updateCall, "apt update call not found")
	require.NotNil(t, installCall, "apt install call not found")
}

func Test_IsPackageInstalled_ChecksInstalledPackages(t *testing.T) {
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			if name == "dpkg-query" {
				return &utils.Result{
					Stdout: []byte("git 1:2.34.1-1ubuntu1.9\ncurl 7.81.0-1ubuntu1.10\n"),
				}, nil
			}
			return &utils.Result{}, nil
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{}
	mockEscalator := &privilege.MoqEscalator{}

	aptManager := apt.NewAptPackageManager(logger.DefaultLogger, mockCommander, mockProgramQuery, mockEscalator)
	packageInfo := pkgmanager.NewPackageInfo("git", "")

	isInstalled, err := aptManager.IsPackageInstalled(packageInfo)

	require.NoError(t, err)
	require.True(t, isInstalled)
}

func Test_IsPackageInstalled_ReturnsFalse_WhenPackageNotInstalled(t *testing.T) {
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			if name == "dpkg-query" {
				return &utils.Result{
					Stdout: []byte("curl 7.81.0-1ubuntu1.10\n"),
				}, nil
			}
			return &utils.Result{}, nil
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{}
	mockEscalator := &privilege.MoqEscalator{}

	aptManager := apt.NewAptPackageManager(logger.DefaultLogger, mockCommander, mockProgramQuery, mockEscalator)
	packageInfo := pkgmanager.NewPackageInfo("git", "")

	isInstalled, err := aptManager.IsPackageInstalled(packageInfo)

	require.NoError(t, err)
	require.False(t, isInstalled)
}

func Test_ListInstalledPackages_ParsesDpkgQueryOutput(t *testing.T) {
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			if name == "dpkg-query" {
				return &utils.Result{
					Stdout: []byte("git 1:2.34.1-1ubuntu1.9\ncurl 7.81.0-1ubuntu1.10\nbuild-essential 12.9ubuntu3\n"),
				}, nil
			}
			return &utils.Result{}, nil
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{}
	mockEscalator := &privilege.MoqEscalator{}

	aptManager := apt.NewAptPackageManager(logger.DefaultLogger, mockCommander, mockProgramQuery, mockEscalator)

	packages, err := aptManager.ListInstalledPackages()

	require.NoError(t, err)
	require.Len(t, packages, 3)
	require.Equal(t, "git", packages[0].Name)
	require.Equal(t, "1:2.34.1-1ubuntu1.9", packages[0].Version)
	require.Equal(t, "curl", packages[1].Name)
	require.Equal(t, "7.81.0-1ubuntu1.10", packages[1].Version)
	require.Equal(t, "build-essential", packages[2].Name)
	require.Equal(t, "12.9ubuntu3", packages[2].Version)
}

func Test_UninstallPackage_CallsAptRemoveCommand_AsRoot(t *testing.T) {
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			return &utils.Result{}, nil
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{}
	mockEscalator := &privilege.MoqEscalator{
		EscalateCommandFunc: func(baseCmd string, baseArgs []string) (privilege.EscalationResult, error) {
			return privilege.EscalationResult{
				Method:          privilege.EscalationNone,
				Command:         baseCmd,
				Args:            baseArgs,
				NeedsEscalation: false,
			}, nil
		},
	}

	aptManager := apt.NewAptPackageManager(logger.DefaultLogger, mockCommander, mockProgramQuery, mockEscalator)
	packageInfo := pkgmanager.NewPackageInfo("git", "1.0.0")

	err := aptManager.UninstallPackage(packageInfo)

	require.NoError(t, err)

	calls := mockCommander.RunCommandCalls()
	require.NotEmpty(t, calls)

	// Find the direct apt remove call (no sudo since running as root)
	var removeCall *struct {
		Name string
		Args []string
		Opts []utils.Option
	}

	for _, call := range calls {
		if call.Name == "apt" && len(call.Args) >= 3 && call.Args[0] == "remove" && call.Args[1] == "-y" {
			removeCall = &call
		}
	}

	require.NotNil(t, removeCall, "apt remove call not found")
}

func Test_GetPackageVersion_ReturnsVersion_WhenPackageIsInstalled(t *testing.T) {
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			if name == "dpkg-query" {
				return &utils.Result{
					Stdout: []byte("git 1:2.34.1-1ubuntu1.9\ncurl 7.81.0-1ubuntu1.10\n"),
				}, nil
			}
			return &utils.Result{}, nil
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{}
	mockEscalator := &privilege.MoqEscalator{}

	aptManager := apt.NewAptPackageManager(logger.DefaultLogger, mockCommander, mockProgramQuery, mockEscalator)

	version, err := aptManager.GetPackageVersion("git")

	require.NoError(t, err)
	require.Equal(t, "1:2.34.1-1ubuntu1.9", version)
}

func Test_GetPackageVersion_ReturnsError_WhenPackageNotInstalled(t *testing.T) {
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			if name == "dpkg-query" {
				return &utils.Result{
					Stdout: []byte("curl 7.81.0-1ubuntu1.10\n"),
				}, nil
			}
			return &utils.Result{}, nil
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{}
	mockEscalator := &privilege.MoqEscalator{}

	aptManager := apt.NewAptPackageManager(logger.DefaultLogger, mockCommander, mockProgramQuery, mockEscalator)

	version, err := aptManager.GetPackageVersion("git")

	require.Error(t, err)
	require.Contains(t, err.Error(), "package git is not installed")
	require.Empty(t, version)
}
