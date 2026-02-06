package osmanager_test

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/MrPointer/dotfiles/installer/utils"
	"github.com/MrPointer/dotfiles/installer/utils/logger"
	"github.com/MrPointer/dotfiles/installer/utils/osmanager"
	"github.com/MrPointer/dotfiles/installer/utils/privilege"
	"github.com/stretchr/testify/require"
)

func Test_UnixOsManager_SetUserShell_UsesEscalatorCommand(t *testing.T) {
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			require.Equal(t, "doas", name)
			if runtime.GOOS == "darwin" {
				require.Equal(t, []string{"dscl", ".", "-create", "/Users/alice", "UserShell", "/bin/zsh"}, args)
				return &utils.Result{}, nil
			}
			require.Equal(t, []string{"usermod", "-s", "/bin/zsh", "alice"}, args)
			return &utils.Result{}, nil
		},
	}

	mockEscalator := &privilege.MoqEscalator{
		EscalateCommandFunc: func(baseCmd string, baseArgs []string) (privilege.EscalationResult, error) {
			if runtime.GOOS == "darwin" {
				require.Equal(t, "dscl", baseCmd)
				require.Equal(t, []string{".", "-create", "/Users/alice", "UserShell", "/bin/zsh"}, baseArgs)
				return privilege.EscalationResult{Command: "doas", Args: append([]string{baseCmd}, baseArgs...)}, nil
			}
			require.Equal(t, "usermod", baseCmd)
			require.Equal(t, []string{"-s", "/bin/zsh", "alice"}, baseArgs)
			return privilege.EscalationResult{Command: "doas", Args: append([]string{baseCmd}, baseArgs...)}, nil
		},
		IsRunningAsRootFunc: func() (bool, error) { return false, nil },
		GetAvailableEscalationMethodsFunc: func() ([]privilege.EscalationMethod, error) {
			return []privilege.EscalationMethod{privilege.EscalationDoas, privilege.EscalationDirect}, nil
		},
	}

	m := osmanager.NewUnixOsManagerWithEscalator(logger.DefaultLogger, mockCommander, mockEscalator)
	err := m.SetUserShell("alice", "/bin/zsh")
	require.NoError(t, err)
}

func Test_UnixOsManager_AddSudoAccess_UsesEscalatorAndTeeWithInput(t *testing.T) {
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			require.Equal(t, "sudo", name)
			require.Equal(t, []string{"tee", "/etc/sudoers.d/alice"}, args)
			return &utils.Result{}, nil
		},
	}

	mockEscalator := &privilege.MoqEscalator{
		EscalateCommandFunc: func(baseCmd string, baseArgs []string) (privilege.EscalationResult, error) {
			require.Equal(t, "tee", baseCmd)
			require.Equal(t, []string{"/etc/sudoers.d/alice"}, baseArgs)
			return privilege.EscalationResult{Command: "sudo", Args: append([]string{baseCmd}, baseArgs...)}, nil
		},
		IsRunningAsRootFunc: func() (bool, error) { return false, nil },
		GetAvailableEscalationMethodsFunc: func() ([]privilege.EscalationMethod, error) {
			return []privilege.EscalationMethod{privilege.EscalationSudo, privilege.EscalationDirect}, nil
		},
	}

	m := osmanager.NewUnixOsManagerWithEscalator(logger.DefaultLogger, mockCommander, mockEscalator)
	err := m.AddSudoAccess("alice")
	require.NoError(t, err)

	calls := mockCommander.RunCommandCalls()
	require.Len(t, calls, 1)

	// Validate we pass stdin via WithInputString (not a shell pipeline).
	var hasInput bool
	for _, opt := range calls[0].Opts {
		cfg := &utils.Options{}
		opt(cfg)
		if string(cfg.Input) == "alice ALL=(ALL) NOPASSWD:ALL\n" {
			hasInput = true
			break
		}
	}
	require.True(t, hasInput, "expected sudoers line to be passed via stdin")
}

func Test_UnixOsManager_runPrivileged_ReturnsEscalatorError(t *testing.T) {
	mockCommander := &utils.MoqCommander{}
	mockEscalator := &privilege.MoqEscalator{
		EscalateCommandFunc: func(baseCmd string, baseArgs []string) (privilege.EscalationResult, error) {
			return privilege.EscalationResult{}, fmt.Errorf("boom")
		},
		IsRunningAsRootFunc: func() (bool, error) { return false, nil },
		GetAvailableEscalationMethodsFunc: func() ([]privilege.EscalationMethod, error) {
			return []privilege.EscalationMethod{privilege.EscalationDirect}, nil
		},
	}

	m := osmanager.NewUnixOsManagerWithEscalator(logger.DefaultLogger, mockCommander, mockEscalator)
	err := m.SetOwnership("/tmp", "alice")
	require.Error(t, err)
	require.Contains(t, err.Error(), "boom")
}
