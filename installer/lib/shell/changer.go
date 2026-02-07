package shell

import (
	"context"
	"fmt"

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
	shellName string
	resolver  ShellResolver
	logger    logger.Logger
	osManager osmanager.OsManager
	escalator privilege.Escalator
}

var _ ShellChanger = (*DefaultShellChanger)(nil)

// NewDefaultShellChanger creates a new DefaultShellChanger instance.
func NewDefaultShellChanger(
	shellName string,
	resolver ShellResolver,
	logger logger.Logger,
	osManager osmanager.OsManager,
	escalator privilege.Escalator,
) *DefaultShellChanger {
	return &DefaultShellChanger{
		shellName: shellName,
		resolver:  resolver,
		logger:    logger,
		osManager: osManager,
		escalator: escalator,
	}
}

// GetShellPath returns the full path to the shell binary.
func (c *DefaultShellChanger) GetShellPath() (string, error) {
	return c.resolver.GetShellPath()
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

	// Ensure shell is in /etc/shells
	if err := c.osManager.EnsureShellInEtcShells(shellPath); err != nil {
		return fmt.Errorf("failed to add shell to /etc/shells: %w", err)
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
