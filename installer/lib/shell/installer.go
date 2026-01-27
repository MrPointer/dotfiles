package shell

import (
	"context"

	"github.com/MrPointer/dotfiles/installer/lib/pkgmanager"
	"github.com/MrPointer/dotfiles/installer/utils/logger"
	"github.com/MrPointer/dotfiles/installer/utils/osmanager"
)

// ShellInstaller provides functionality to install a shell and set it as default.
type ShellInstaller interface {
	// IsAvailable checks if the shell is already installed and available in PATH.
	IsAvailable() (bool, error)

	// Install installs the shell using the configured package manager.
	Install(ctx context.Context) error

	// SetAsDefault sets the shell as the user's default login shell.
	SetAsDefault(ctx context.Context) error
}

// DefaultShellInstaller implements ShellInstaller.
type DefaultShellInstaller struct {
	shellName    string
	programQuery osmanager.ProgramQuery
	pkgManager   pkgmanager.PackageManager
	shellChanger ShellChanger
	logger       logger.Logger
}

var _ ShellInstaller = (*DefaultShellInstaller)(nil)

// NewDefaultShellInstaller creates a new DefaultShellInstaller instance.
func NewDefaultShellInstaller(
	shellName string,
	programQuery osmanager.ProgramQuery,
	pkgManager pkgmanager.PackageManager,
	shellChanger ShellChanger,
	logger logger.Logger,
) *DefaultShellInstaller {
	return &DefaultShellInstaller{
		shellName:    shellName,
		programQuery: programQuery,
		pkgManager:   pkgManager,
		shellChanger: shellChanger,
		logger:       logger,
	}
}

func (d *DefaultShellInstaller) IsAvailable() (bool, error) {
	d.logger.Debug("Checking if %s is available", d.shellName)

	shellAvailable, err := d.programQuery.ProgramExists(d.shellName)
	if err != nil {
		return false, err
	}

	if shellAvailable {
		d.logger.Debug("%s is available", d.shellName)
	} else {
		d.logger.Debug("%s is not available", d.shellName)
	}
	return shellAvailable, nil
}

func (d *DefaultShellInstaller) Install(ctx context.Context) error {
	d.logger.Debug("Installing %s via package manager", d.shellName)
	return d.pkgManager.InstallPackage(pkgmanager.NewRequestedPackageInfo(d.shellName, nil))
}

// SetAsDefault sets the shell as the user's default login shell.
func (d *DefaultShellInstaller) SetAsDefault(ctx context.Context) error {
	d.logger.Debug("Setting %s as default shell", d.shellName)
	return d.shellChanger.SetAsDefault(ctx)
}
