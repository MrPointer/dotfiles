package shell

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/MrPointer/dotfiles/installer/utils"
	"github.com/MrPointer/dotfiles/installer/utils/logger"
	"github.com/MrPointer/dotfiles/installer/utils/osmanager"
	"github.com/MrPointer/dotfiles/installer/utils/privilege"
)

// ShellChanger provides functionality to change the user's default login shell.
type ShellChanger interface {
	// GetShellPath returns the full path to the shell binary.
	GetShellPath() (string, error)

	// IsCurrentDefault checks if the shell is already the user's default.
	IsCurrentDefault() (bool, error)

	// SetAsDefault sets the shell as the user's default login shell.
	SetAsDefault(ctx context.Context) error
}

// DefaultShellChanger implements ShellChanger with platform-specific logic.
type DefaultShellChanger struct {
	shellName  string
	brewPath   string // path to brew binary, empty if not installed via brew
	logger     logger.Logger
	osManager  osmanager.OsManager
	fileSystem utils.FileSystem
	commander  utils.Commander
	escalator  privilege.Escalator
}

var _ ShellChanger = (*DefaultShellChanger)(nil)

// NewDefaultShellChanger creates a new DefaultShellChanger instance.
// brewPath should be the path to the brew binary if the shell was installed via Homebrew,
// or empty string if installed via system package manager.
func NewDefaultShellChanger(
	shellName string,
	brewPath string,
	logger logger.Logger,
	osManager osmanager.OsManager,
	fileSystem utils.FileSystem,
	commander utils.Commander,
	escalator privilege.Escalator,
) *DefaultShellChanger {
	return &DefaultShellChanger{
		shellName:  shellName,
		brewPath:   brewPath,
		logger:     logger,
		osManager:  osManager,
		fileSystem: fileSystem,
		commander:  commander,
		escalator:  escalator,
	}
}

// GetShellPath returns the full path to the shell binary.
// If installed via brew, returns the brew-installed shell path.
// Otherwise, uses GetProgramPath to find the shell in PATH.
func (c *DefaultShellChanger) GetShellPath() (string, error) {
	if c.brewPath != "" {
		// Shell was installed via Homebrew - use brew prefix to find it
		brewBinDir := filepath.Dir(c.brewPath)
		shellPath := filepath.Join(brewBinDir, c.shellName)

		c.logger.Debug("Looking for brew-installed shell at: %s", shellPath)

		if err := c.validateShellPath(shellPath); err != nil {
			return "", fmt.Errorf("brew-installed shell not found at %s: %w", shellPath, err)
		}

		return shellPath, nil
	}

	// Shell was installed via system package manager - use GetProgramPath to find it
	c.logger.Debug("Looking for system-installed shell using GetProgramPath(%s)", c.shellName)

	shellPath, err := c.osManager.GetProgramPath(c.shellName)
	if err != nil {
		return "", fmt.Errorf("failed to find %s in PATH: %w", c.shellName, err)
	}

	if err := c.validateShellPath(shellPath); err != nil {
		return "", err
	}

	return shellPath, nil
}

// IsCurrentDefault checks if the shell is already the user's default.
func (c *DefaultShellChanger) IsCurrentDefault() (bool, error) {
	shellPath, err := c.GetShellPath()
	if err != nil {
		return false, err
	}

	username, err := c.osManager.GetCurrentUsername()
	if err != nil {
		return false, fmt.Errorf("failed to get current username: %w", err)
	}

	currentShell, err := c.osManager.GetUserShell(username)
	if err != nil {
		return false, fmt.Errorf("failed to get current shell: %w", err)
	}

	c.logger.Debug("Current shell: %s, target shell: %s", currentShell, shellPath)

	return currentShell == shellPath, nil
}

// SetAsDefault sets the shell as the user's default login shell.
func (c *DefaultShellChanger) SetAsDefault(ctx context.Context) error {
	shellPath, err := c.GetShellPath()
	if err != nil {
		return fmt.Errorf("failed to get shell path: %w", err)
	}

	// Check if already default
	isDefault, err := c.IsCurrentDefault()
	if err != nil {
		c.logger.Warning("Failed to check current default shell: %v", err)
		// Continue anyway - we'll try to set it
	} else if isDefault {
		c.logger.Debug("Shell %s is already the default", shellPath)
		return nil
	}

	// Log warning if running as root (common in CI containers)
	isRoot, err := c.escalator.IsRunningAsRoot()
	if err == nil && isRoot {
		c.logger.Warning("Running as root - shell change will affect root user's default shell")
	}

	// Ensure shell is in /etc/shells (only needed for Homebrew installations)
	if c.brewPath != "" {
		if err := c.ensureShellInEtcShells(shellPath); err != nil {
			return fmt.Errorf("failed to add shell to /etc/shells: %w", err)
		}
	}

	// Set default shell using OsManager
	username, err := c.osManager.GetCurrentUsername()
	if err != nil {
		return fmt.Errorf("failed to get current username: %w", err)
	}

	c.logger.Debug("Setting default shell to %s for user %s", shellPath, username)

	if err := c.osManager.SetUserShell(username, shellPath); err != nil {
		return fmt.Errorf("failed to set default shell: %w", err)
	}

	return nil
}

// validateShellPath verifies that the shell binary exists and is executable.
func (c *DefaultShellChanger) validateShellPath(shellPath string) error {
	exists, err := c.fileSystem.PathExists(shellPath)
	if err != nil {
		return fmt.Errorf("failed to check shell binary at %s: %w", shellPath, err)
	}
	if !exists {
		return fmt.Errorf("shell binary not found at %s", shellPath)
	}

	executable, err := c.fileSystem.IsExecutable(shellPath)
	if err != nil {
		return fmt.Errorf("failed to check if shell binary is executable at %s: %w", shellPath, err)
	}
	if !executable {
		return fmt.Errorf("shell binary at %s is not executable", shellPath)
	}

	return nil
}

// ensureShellInEtcShells adds the shell to /etc/shells if not already present.
// This is required for Homebrew-installed shells before changing the default.
func (c *DefaultShellChanger) ensureShellInEtcShells(shellPath string) error {
	c.logger.Debug("Checking if %s is in /etc/shells", shellPath)

	// Read current contents using FileSystem interface
	content, err := c.fileSystem.ReadFileContents("/etc/shells")
	if err != nil {
		return fmt.Errorf("failed to read /etc/shells: %w", err)
	}

	// Check if already present
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == shellPath {
			c.logger.Debug("Shell %s already in /etc/shells", shellPath)
			return nil
		}
	}

	c.logger.Info("Adding %s to /etc/shells", shellPath)

	// Append using tee with privilege escalation
	escalated, err := c.escalator.EscalateCommand("tee", []string{"-a", "/etc/shells"})
	if err != nil {
		return fmt.Errorf("failed to escalate command for /etc/shells modification: %w", err)
	}

	_, err = c.commander.RunCommand(
		escalated.Command,
		escalated.Args,
		utils.WithCaptureOutput(),
		utils.WithInputString(shellPath+"\n"),
	)
	if err != nil {
		return fmt.Errorf("failed to append shell to /etc/shells: %w", err)
	}

	return nil
}
