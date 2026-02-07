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
	// EscalationDirect indicates falling back to direct execution (may fail).
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
	// escalated command based on the current system state and available tools.
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
// 1. Root containers/users: Run commands directly (no privilege escalation needed)
// 2. Non-root with sudo: Use sudo for privilege escalation
// 3. Non-root with doas: Use doas as alternative to sudo (OpenBSD, some Linux)
// 4. Bare containers: Fall back to direct execution with warning
//
// This approach handles common scenarios like:
// - Docker containers running as root (no sudo needed)
// - Minimal Ubuntu containers without sudo package
// - Systems with alternative privilege escalation tools
// - Regular user systems with proper sudo setup
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

	// Fall back to direct execution (might fail, but let's try)
	e.logger.Warning("Running as non-root without sudo/doas - command may fail due to insufficient privileges")
	return EscalationResult{
		Method:          EscalationDirect,
		Command:         baseCmd,
		Args:            baseArgs,
		NeedsEscalation: false,
	}, nil
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

// GetAvailableEscalationMethods returns all escalation methods available on this system.
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

	// Direct execution is always available as fallback
	methods = append(methods, EscalationDirect)

	return methods, nil
}

// isSudoAvailable checks if sudo is available and usable.
// First checks if the sudo command exists, then tests if it can be used
// without a password prompt using the -n (non-interactive) flag.
// This handles cases where sudo exists but requires password authentication.
func (e *DefaultEscalator) isSudoAvailable() bool {
	e.logger.Trace("Checking if sudo is available")

	// First check if sudo command exists
	exists, err := e.programQuery.ProgramExists("sudo")
	if err != nil || !exists {
		return false
	}

	// Test if sudo can be used (with -n flag for non-interactive)
	_, err = e.commander.RunCommand("sudo", []string{"-n", "true"}, utils.WithCaptureOutput())
	return err == nil
}

// isDoasAvailable checks if doas is available and usable.
// doas is an alternative to sudo commonly found on OpenBSD and some Linux systems.
// Like sudo, we test both existence and usability with non-interactive flag.
func (e *DefaultEscalator) isDoasAvailable() bool {
	e.logger.Trace("Checking if doas is available")

	// First check if doas command exists
	exists, err := e.programQuery.ProgramExists("doas")
	if err != nil || !exists {
		return false
	}

	// Test if doas can be used (with -n flag for non-interactive)
	_, err = e.commander.RunCommand("doas", []string{"-n", "true"}, utils.WithCaptureOutput())
	return err == nil
}
