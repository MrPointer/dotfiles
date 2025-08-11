package shell

import (
	"context"

	"github.com/MrPointer/dotfiles/installer/lib/pkgmanager"
	"github.com/MrPointer/dotfiles/installer/utils/logger"
	"github.com/MrPointer/dotfiles/installer/utils/osmanager"
)

type ShellInstaller interface {
	IsAvailable() (bool, error)
	Install(ctx context.Context) error
}

type DefaultShellInstaller struct {
	shellName    string
	programQuery osmanager.ProgramQuery
	pkgManager   pkgmanager.PackageManager
	logger       logger.Logger
}

var _ ShellInstaller = (*DefaultShellInstaller)(nil)

func NewDefaultShellInstaller(
	shellName string,
	programQuery osmanager.ProgramQuery,
	pkgManager pkgmanager.PackageManager,
	logger logger.Logger,
) *DefaultShellInstaller {
	return &DefaultShellInstaller{
		shellName:    shellName,
		programQuery: programQuery,
		pkgManager:   pkgManager,
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
