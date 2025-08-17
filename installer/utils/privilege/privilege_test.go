package privilege_test

import (
	"fmt"
	"testing"

	"github.com/MrPointer/dotfiles/installer/utils"
	"github.com/MrPointer/dotfiles/installer/utils/logger"
	"github.com/MrPointer/dotfiles/installer/utils/osmanager"
	"github.com/MrPointer/dotfiles/installer/utils/privilege"
	"github.com/stretchr/testify/require"
)

func Test_DefaultEscalator_ImplementsEscalatorInterface(t *testing.T) {
	mockCommander := &utils.MoqCommander{}
	mockProgramQuery := &osmanager.MoqProgramQuery{}

	escalator := privilege.NewDefaultEscalator(logger.DefaultLogger, mockCommander, mockProgramQuery)

	require.Implements(t, (*privilege.Escalator)(nil), escalator)
}

func Test_NewDefaultEscalator_ReturnsValidInstance(t *testing.T) {
	mockCommander := &utils.MoqCommander{}
	mockProgramQuery := &osmanager.MoqProgramQuery{}

	escalator := privilege.NewDefaultEscalator(logger.DefaultLogger, mockCommander, mockProgramQuery)

	require.NotNil(t, escalator)
}

func Test_EscalateCommand_RunningAsRoot_ReturnsDirectCommand(t *testing.T) {
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			if name == "id" && len(args) == 1 && args[0] == "-u" {
				return &utils.Result{Stdout: []byte("0")}, nil
			}
			return &utils.Result{}, nil
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{
		ProgramExistsFunc: func(program string) (bool, error) {
			return program == "id", nil
		},
	}

	escalator := privilege.NewDefaultEscalator(logger.DefaultLogger, mockCommander, mockProgramQuery)

	result, err := escalator.EscalateCommand("apt", []string{"install", "git"})

	require.NoError(t, err)
	require.Equal(t, privilege.EscalationNone, result.Method)
	require.Equal(t, "apt", result.Command)
	require.Equal(t, []string{"install", "git"}, result.Args)
	require.False(t, result.NeedsEscalation)
}

func Test_EscalateCommand_NonRootWithSudo_ReturnsSudoCommand(t *testing.T) {
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			if name == "id" && len(args) == 1 && args[0] == "-u" {
				return &utils.Result{Stdout: []byte("1000")}, nil
			}
			if name == "sudo" && len(args) == 2 && args[0] == "-n" && args[1] == "true" {
				return &utils.Result{}, nil
			}
			return &utils.Result{}, nil
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{
		ProgramExistsFunc: func(program string) (bool, error) {
			return program == "id" || program == "sudo", nil
		},
	}

	escalator := privilege.NewDefaultEscalator(logger.DefaultLogger, mockCommander, mockProgramQuery)

	result, err := escalator.EscalateCommand("apt", []string{"install", "git"})

	require.NoError(t, err)
	require.Equal(t, privilege.EscalationSudo, result.Method)
	require.Equal(t, "sudo", result.Command)
	require.Equal(t, []string{"apt", "install", "git"}, result.Args)
	require.True(t, result.NeedsEscalation)
}

func Test_EscalateCommand_NonRootWithDoas_ReturnsDoasCommand(t *testing.T) {
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			if name == "id" && len(args) == 1 && args[0] == "-u" {
				return &utils.Result{Stdout: []byte("1000")}, nil
			}
			if name == "sudo" && len(args) == 2 && args[0] == "-n" && args[1] == "true" {
				return &utils.Result{ExitCode: 1}, fmt.Errorf("sudo not available")
			}
			if name == "doas" && len(args) == 2 && args[0] == "-n" && args[1] == "true" {
				return &utils.Result{}, nil
			}
			return &utils.Result{}, nil
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{
		ProgramExistsFunc: func(program string) (bool, error) {
			if program == "sudo" {
				return false, nil
			}
			return program == "id" || program == "doas", nil
		},
	}

	escalator := privilege.NewDefaultEscalator(logger.DefaultLogger, mockCommander, mockProgramQuery)

	result, err := escalator.EscalateCommand("dnf", []string{"install", "vim"})

	require.NoError(t, err)
	require.Equal(t, privilege.EscalationDoas, result.Method)
	require.Equal(t, "doas", result.Command)
	require.Equal(t, []string{"dnf", "install", "vim"}, result.Args)
	require.True(t, result.NeedsEscalation)
}

func Test_EscalateCommand_NonRootWithoutPrivilegeEscalation_ReturnsDirectCommand(t *testing.T) {
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			if name == "id" && len(args) == 1 && args[0] == "-u" {
				return &utils.Result{Stdout: []byte("1000")}, nil
			}
			// Both sudo and doas fail
			if (name == "sudo" || name == "doas") && len(args) == 2 && args[0] == "-n" && args[1] == "true" {
				return &utils.Result{ExitCode: 1}, fmt.Errorf("privilege escalation not available")
			}
			return &utils.Result{}, nil
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{
		ProgramExistsFunc: func(program string) (bool, error) {
			if program == "sudo" || program == "doas" {
				return false, nil
			}
			return program == "id", nil
		},
	}

	escalator := privilege.NewDefaultEscalator(logger.DefaultLogger, mockCommander, mockProgramQuery)

	result, err := escalator.EscalateCommand("apt", []string{"install", "git"})

	require.NoError(t, err)
	require.Equal(t, privilege.EscalationDirect, result.Method)
	require.Equal(t, "apt", result.Command)
	require.Equal(t, []string{"install", "git"}, result.Args)
	require.False(t, result.NeedsEscalation)
}

func Test_EscalateCommand_EmptyCommand_ReturnsError(t *testing.T) {
	mockCommander := &utils.MoqCommander{}
	mockProgramQuery := &osmanager.MoqProgramQuery{}

	escalator := privilege.NewDefaultEscalator(logger.DefaultLogger, mockCommander, mockProgramQuery)

	result, err := escalator.EscalateCommand("", []string{"install", "git"})

	require.Error(t, err)
	require.Contains(t, err.Error(), "base command cannot be empty")
	require.Equal(t, privilege.EscalationResult{}, result)
}

func Test_IsRunningAsRoot_ReturnsTrue_WhenUidIsZero(t *testing.T) {
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			if name == "id" && len(args) == 1 && args[0] == "-u" {
				return &utils.Result{Stdout: []byte("0")}, nil
			}
			return &utils.Result{}, nil
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{
		ProgramExistsFunc: func(program string) (bool, error) {
			return program == "id", nil
		},
	}

	escalator := privilege.NewDefaultEscalator(logger.DefaultLogger, mockCommander, mockProgramQuery)

	isRoot, err := escalator.IsRunningAsRoot()

	require.NoError(t, err)
	require.True(t, isRoot)
}

func Test_IsRunningAsRoot_ReturnsFalse_WhenUidIsNotZero(t *testing.T) {
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			if name == "id" && len(args) == 1 && args[0] == "-u" {
				return &utils.Result{Stdout: []byte("1000")}, nil
			}
			return &utils.Result{}, nil
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{
		ProgramExistsFunc: func(program string) (bool, error) {
			return program == "id", nil
		},
	}

	escalator := privilege.NewDefaultEscalator(logger.DefaultLogger, mockCommander, mockProgramQuery)

	isRoot, err := escalator.IsRunningAsRoot()

	require.NoError(t, err)
	require.False(t, isRoot)
}

func Test_IsRunningAsRoot_ReturnsError_WhenIdCommandNotAvailable(t *testing.T) {
	mockCommander := &utils.MoqCommander{}
	mockProgramQuery := &osmanager.MoqProgramQuery{
		ProgramExistsFunc: func(program string) (bool, error) {
			return false, nil
		},
	}

	escalator := privilege.NewDefaultEscalator(logger.DefaultLogger, mockCommander, mockProgramQuery)

	isRoot, err := escalator.IsRunningAsRoot()

	require.Error(t, err)
	require.Contains(t, err.Error(), "'id' command not available")
	require.False(t, isRoot)
}

func Test_IsRunningAsRoot_ReturnsError_WhenIdCommandFails(t *testing.T) {
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			if name == "id" && len(args) == 1 && args[0] == "-u" {
				return &utils.Result{}, fmt.Errorf("command failed")
			}
			return &utils.Result{}, nil
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{
		ProgramExistsFunc: func(program string) (bool, error) {
			return program == "id", nil
		},
	}

	escalator := privilege.NewDefaultEscalator(logger.DefaultLogger, mockCommander, mockProgramQuery)

	isRoot, err := escalator.IsRunningAsRoot()

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to execute 'id -u'")
	require.False(t, isRoot)
}

func Test_GetAvailableEscalationMethods_AsRoot_ReturnsNoneOnly(t *testing.T) {
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			if name == "id" && len(args) == 1 && args[0] == "-u" {
				return &utils.Result{Stdout: []byte("0")}, nil
			}
			return &utils.Result{}, nil
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{
		ProgramExistsFunc: func(program string) (bool, error) {
			return program == "id", nil
		},
	}

	escalator := privilege.NewDefaultEscalator(logger.DefaultLogger, mockCommander, mockProgramQuery)

	methods, err := escalator.GetAvailableEscalationMethods()

	require.NoError(t, err)
	require.Equal(t, []privilege.EscalationMethod{privilege.EscalationNone}, methods)
}

func Test_GetAvailableEscalationMethods_AsNonRootWithSudo_ReturnsSudoAndDirect(t *testing.T) {
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			if name == "id" && len(args) == 1 && args[0] == "-u" {
				return &utils.Result{Stdout: []byte("1000")}, nil
			}
			if name == "sudo" && len(args) == 2 && args[0] == "-n" && args[1] == "true" {
				return &utils.Result{}, nil
			}
			return &utils.Result{ExitCode: 1}, fmt.Errorf("command failed")
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{
		ProgramExistsFunc: func(program string) (bool, error) {
			return program == "id" || program == "sudo", nil
		},
	}

	escalator := privilege.NewDefaultEscalator(logger.DefaultLogger, mockCommander, mockProgramQuery)

	methods, err := escalator.GetAvailableEscalationMethods()

	require.NoError(t, err)
	require.Contains(t, methods, privilege.EscalationSudo)
	require.Contains(t, methods, privilege.EscalationDirect)
	require.NotContains(t, methods, privilege.EscalationNone)
}

func Test_GetAvailableEscalationMethods_AsNonRootWithBothSudoAndDoas_ReturnsAll(t *testing.T) {
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			if name == "id" && len(args) == 1 && args[0] == "-u" {
				return &utils.Result{Stdout: []byte("1000")}, nil
			}
			if (name == "sudo" || name == "doas") && len(args) == 2 && args[0] == "-n" && args[1] == "true" {
				return &utils.Result{}, nil
			}
			return &utils.Result{}, nil
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{
		ProgramExistsFunc: func(program string) (bool, error) {
			return program == "id" || program == "sudo" || program == "doas", nil
		},
	}

	escalator := privilege.NewDefaultEscalator(logger.DefaultLogger, mockCommander, mockProgramQuery)

	methods, err := escalator.GetAvailableEscalationMethods()

	require.NoError(t, err)
	require.Contains(t, methods, privilege.EscalationSudo)
	require.Contains(t, methods, privilege.EscalationDoas)
	require.Contains(t, methods, privilege.EscalationDirect)
	require.Len(t, methods, 3)
}

func Test_GetAvailableEscalationMethods_AsNonRootWithoutPrivilegeEscalation_ReturnsDirectOnly(t *testing.T) {
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			if name == "id" && len(args) == 1 && args[0] == "-u" {
				return &utils.Result{Stdout: []byte("1000")}, nil
			}
			return &utils.Result{ExitCode: 1}, fmt.Errorf("command failed")
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{
		ProgramExistsFunc: func(program string) (bool, error) {
			return program == "id", nil
		},
	}

	escalator := privilege.NewDefaultEscalator(logger.DefaultLogger, mockCommander, mockProgramQuery)

	methods, err := escalator.GetAvailableEscalationMethods()

	require.NoError(t, err)
	require.Equal(t, []privilege.EscalationMethod{privilege.EscalationDirect}, methods)
}

func Test_EscalateCommand_HandlesRootCheckFailure_GracefullyFallsBackToNonRoot(t *testing.T) {
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			if name == "id" && len(args) == 1 && args[0] == "-u" {
				return &utils.Result{}, fmt.Errorf("id command failed")
			}
			if name == "sudo" && len(args) == 2 && args[0] == "-n" && args[1] == "true" {
				return &utils.Result{}, nil
			}
			return &utils.Result{}, nil
		},
	}
	mockProgramQuery := &osmanager.MoqProgramQuery{
		ProgramExistsFunc: func(program string) (bool, error) {
			return program == "id" || program == "sudo", nil
		},
	}

	escalator := privilege.NewDefaultEscalator(logger.DefaultLogger, mockCommander, mockProgramQuery)

	result, err := escalator.EscalateCommand("apt", []string{"install", "git"})

	require.NoError(t, err)
	require.Equal(t, privilege.EscalationSudo, result.Method)
	require.Equal(t, "sudo", result.Command)
	require.Equal(t, []string{"apt", "install", "git"}, result.Args)
	require.True(t, result.NeedsEscalation)
}
