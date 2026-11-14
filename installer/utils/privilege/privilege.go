package privilege

import (
	"fmt"
	"strings"

	"github.com/MrPointer/dotfiles/installer/utils"
	"github.com/MrPointer/dotfiles/installer/utils/logger"
)

// EscalationMethod represents the type of privilege escalation being used.
type EscalationMethod string

const (
	// EscalationNone indicates no privilege escalation is needed (running as root).
	EscalationNone EscalationMethod = "none"
	// EscalationSudo indicates using sudo for privilege escalation.
	EscalationSudo EscalationMethod = "sudo"
	// EscalationDoas indicates using doas for privilege escalation.
	EscalationDoas EscalationMethod = "doas"
	// EscalationDirect is the legacy unprivileged fallback, no longer returned by DefaultEscalator.
	EscalationDirect EscalationMethod = "direct"
)

// EscalationResult contains information about how a command should be escalated.
type EscalationResult struct {
	// Method indicates which escalation method was chosen.
	Method EscalationMethod
	// Command is the final command to execute.
	Command string
	// Args are the final arguments to pass to the command.
	Args []string
	// NeedsEscalation indicates whether privilege escalation is being used.
	NeedsEscalation bool
}

// Escalator provides smart privilege escalation for system commands.
type Escalator interface {
	// EscalateCommand takes a base command and arguments and returns the appropriate
	// escalated command based on the current system state and installed tools.
	// Authentication is handled when the returned command is executed.
	EscalateCommand(baseCmd string, baseArgs []string) (EscalationResult, error)

	// IsRunningAsRoot checks if the current process has root privileges.
	IsRunningAsRoot() (bool, error)

	// GetAvailableEscalationMethods returns the escalation methods available on this system.
	GetAvailableEscalationMethods() ([]EscalationMethod, error)
}

// DefaultEscalator implements smart privilege escalation using various system tools.
type DefaultEscalator struct {
	logger       logger.Logger
	commander    utils.Commander
	programQuery utils.ProgramQuery
}

var _ Escalator = (*DefaultEscalator)(nil)

// NewDefaultEscalator creates a new DefaultEscalator instance.
func NewDefaultEscalator(logger logger.Logger, commander utils.Commander, programQuery utils.ProgramQuery) *DefaultEscalator {
	return &DefaultEscalator{
		logger:       logger,
		commander:    commander,
		programQuery: programQuery,
	}
}

// EscalateCommand implements smart privilege escalation strategy.
// It determines the best method for privilege escalation based on:
//
// 1. Root containers/users: Run commands directly (no privilege escalation needed).
// 2. Non-root with sudo installed: Use sudo for privilege escalation.
// 3. Non-root with only doas installed: Use doas for privilege escalation.
// 4. Non-root without either tool: Return an actionable error.
//
// Tool selection does not authenticate; execution handles authentication and authorization.
func (e *DefaultEscalator) EscalateCommand(baseCmd string, baseArgs []string) (EscalationResult, error) {
	e.logger.Debug("Escalating command to run with privileges: %s %s", baseCmd, strings.Join(baseArgs, " "))

	if baseCmd == "" {
		return EscalationResult{}, fmt.Errorf("base command cannot be empty")
	}

	// Check if we're already running as root
	isRoot, err := e.IsRunningAsRoot()
	if err != nil {
		e.logger.Warning("Failed to check if running as root, assuming non-root: %v", err)
		isRoot = false
	}

	if isRoot {
		e.logger.Trace("Already running as root")
		return EscalationResult{
			Method:          EscalationNone,
			Command:         baseCmd,
			Args:            baseArgs,
			NeedsEscalation: false,
		}, nil
	}

	e.logger.Trace("Not running as root")

	// Check if sudo is available
	if e.isSudoAvailable() {
		e.logger.Trace("Sudo is available")
		args := make([]string, 0, len(baseArgs)+1)
		args = append(args, baseCmd)
		args = append(args, baseArgs...)
		return EscalationResult{
			Method:          EscalationSudo,
			Command:         "sudo",
			Args:            args,
			NeedsEscalation: true,
		}, nil
	}

	e.logger.Trace("Sudo is not available")

	// Check if doas is available (alternative to sudo on some systems)
	if e.isDoasAvailable() {
		e.logger.Trace("Doas is available")
		args := make([]string, 0, len(baseArgs)+1)
		args = append(args, baseCmd)
		args = append(args, baseArgs...)
		return EscalationResult{
			Method:          EscalationDoas,
			Command:         "doas",
			Args:            args,
			NeedsEscalation: true,
		}, nil
	}

	return EscalationResult{}, fmt.Errorf("privilege escalation required: install sudo or doas and configure access")
}

// IsRunningAsRoot checks if the current process is running as root.
// Uses `id -u` command to get the user ID, where 0 indicates root.
// This is more reliable than checking environment variables or using
// os.Geteuid() which might not work correctly in all container environments.
func (e *DefaultEscalator) IsRunningAsRoot() (bool, error) {
	e.logger.Trace("Checking if running as root")

	exists, err := e.programQuery.ProgramExists("id")
	if err != nil {
		return false, fmt.Errorf("failed to check if 'id' command exists: %w", err)
	}
	if !exists {
		return false, fmt.Errorf("'id' command not available on this system")
	}

	result, err := e.commander.RunCommand("id", []string{"-u"}, utils.WithCaptureOutput())
	if err != nil {
		return false, fmt.Errorf("failed to execute 'id -u': %w", err)
	}

	uid := strings.TrimSpace(string(result.Stdout))
	return uid == "0", nil
}

// GetAvailableEscalationMethods returns installed escalation tools in preference order,
// or EscalationNone when running as root. It does not check authentication or authorization.
func (e *DefaultEscalator) GetAvailableEscalationMethods() ([]EscalationMethod, error) {
	methods := []EscalationMethod{}

	// Check if running as root
	isRoot, err := e.IsRunningAsRoot()
	if err != nil {
		e.logger.Warning("Failed to check root status: %v", err)
	} else if isRoot {
		methods = append(methods, EscalationNone)
		return methods, nil // If root, no other methods needed
	}

	// Check available privilege escalation tools
	if e.isSudoAvailable() {
		methods = append(methods, EscalationSudo)
	}

	if e.isDoasAvailable() {
		methods = append(methods, EscalationDoas)
	}

	if len(methods) == 0 {
		return nil, fmt.Errorf("privilege escalation required: install sudo or doas and configure access")
	}

	return methods, nil
}

// isSudoAvailable checks if sudo is installed, regardless of authentication state.
func (e *DefaultEscalator) isSudoAvailable() bool {
	e.logger.Trace("Checking if sudo is available")

	exists, err := e.programQuery.ProgramExists("sudo")
	return err == nil && exists
}

// isDoasAvailable checks if doas is installed, regardless of authentication state.
func (e *DefaultEscalator) isDoasAvailable() bool {
	e.logger.Trace("Checking if doas is available")

	exists, err := e.programQuery.ProgramExists("doas")
	return err == nil && exists
}
