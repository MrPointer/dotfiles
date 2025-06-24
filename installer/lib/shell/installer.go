package shell

import (
	"context"

	"github.com/MrPointer/dotfiles/installer/lib/pkgmanager"
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
}

var _ ShellInstaller = (*DefaultShellInstaller)(nil)

func NewDefaultShellInstaller(
	shellName string,
	programQuery osmanager.ProgramQuery,
	pkgManager pkgmanager.PackageManager,
) *DefaultShellInstaller {
	return &DefaultShellInstaller{
		shellName:    shellName,
		programQuery: programQuery,
		pkgManager:   pkgManager,
	}
}

func (d *DefaultShellInstaller) IsAvailable() (bool, error) {
	shellAvailable, err := d.programQuery.ProgramExists(d.shellName)
	if err != nil {
		return false, err
	}
	return shellAvailable, nil
}

func (d *DefaultShellInstaller) Install(ctx context.Context) error {
	return d.pkgManager.InstallPackage(pkgmanager.NewRequestedPackageInfo(d.shellName, nil))
}
