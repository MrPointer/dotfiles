package shell

import (
	"context"
	"fmt"
	"path/filepath"

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
	escalator privilege.Escalator,
) *DefaultShellChanger {
	return &DefaultShellChanger{
		shellName:  shellName,
		brewPath:   brewPath,
		logger:     logger,
		osManager:  osManager,
		fileSystem: fileSystem,
		escalator:  escalator,
	}
}

// GetShellPath returns the full path to the shell binary.
// Checks system PATH first (which includes brew if PATH is configured),
// then falls back to checking brew's bin directory directly.
func (c *DefaultShellChanger) GetShellPath() (string, error) {
	// First check system PATH (includes brew if PATH is configured)
	shellPath, err := c.osManager.GetProgramPath(c.shellName)
	if err == nil {
		if err := c.validateShellPath(shellPath); err == nil {
			c.logger.Debug("Found %s in PATH at: %s", c.shellName, shellPath)
			return shellPath, nil
		}
		c.logger.Debug("Shell found in PATH at %s but validation failed", shellPath)
	}

	// If brew is available, check brew's bin directory directly as fallback
	// This handles cases where brew installed the shell but PATH isn't updated yet
	if c.brewPath != "" {
		brewBinDir := filepath.Dir(c.brewPath)
		brewShellPath := filepath.Join(brewBinDir, c.shellName)
		c.logger.Debug("Checking brew bin directory for %s at: %s", c.shellName, brewShellPath)

		if err := c.validateShellPath(brewShellPath); err == nil {
			c.logger.Debug("Found %s in brew bin at: %s", c.shellName, brewShellPath)
			return brewShellPath, nil
		}
	}

	return "", fmt.Errorf("could not find %s: not in PATH and not in brew bin directory", c.shellName)
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
		if err := c.osManager.EnsureShellInEtcShells(shellPath); err != nil {
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
