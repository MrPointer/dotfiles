package shell_test

import (
	"context"
	"testing"

	"github.com/MrPointer/dotfiles/installer/lib/pkgmanager"
	"github.com/MrPointer/dotfiles/installer/lib/shell"
	"github.com/MrPointer/dotfiles/installer/utils/logger"
	"github.com/MrPointer/dotfiles/installer/utils/osmanager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ShellIsReportedAsAvailable_WhenShellProgramExists(t *testing.T) {
	// Arrange
	shellName := "zsh"

	programQueryMock := &osmanager.MoqProgramQuery{
		ProgramExistsFunc: func(program string) (bool, error) {
			return true, nil
		},
	}

	pkgManagerMock := &pkgmanager.MoqPackageManager{}

	installer := shell.NewDefaultShellInstaller(shellName, programQueryMock, pkgManagerMock, logger.DefaultLogger)

	// Act
	available, err := installer.IsAvailable()

	// Assert
	require.NoError(t, err)
	assert.True(t, available)
	assert.Len(t, programQueryMock.ProgramExistsCalls(), 1)
	assert.Equal(t, shellName, programQueryMock.ProgramExistsCalls()[0].Program)
}

func Test_ShellIsReportedAsUnavailable_WhenShellProgramDoesNotExist(t *testing.T) {
	// Arrange
	shellName := "fish"

	programQueryMock := &osmanager.MoqProgramQuery{
		ProgramExistsFunc: func(program string) (bool, error) {
			return false, nil
		},
	}

	pkgManagerMock := &pkgmanager.MoqPackageManager{}

	installer := shell.NewDefaultShellInstaller(shellName, programQueryMock, pkgManagerMock, logger.DefaultLogger)

	// Act
	available, err := installer.IsAvailable()

	// Assert
	require.NoError(t, err)
	assert.False(t, available)
	assert.Len(t, programQueryMock.ProgramExistsCalls(), 1)
	assert.Equal(t, shellName, programQueryMock.ProgramExistsCalls()[0].Program)
}

func Test_ShellAvailabilityCheckFails_WhenProgramExistsFails(t *testing.T) {
	// Arrange
	shellName := "bash"

	programQueryMock := &osmanager.MoqProgramQuery{
		ProgramExistsFunc: func(program string) (bool, error) {
			return false, assert.AnError
		},
	}

	pkgManagerMock := &pkgmanager.MoqPackageManager{}

	installer := shell.NewDefaultShellInstaller(shellName, programQueryMock, pkgManagerMock, logger.DefaultLogger)

	// Act
	available, err := installer.IsAvailable()

	// Assert
	require.Error(t, err)
	assert.False(t, available)
	assert.Len(t, programQueryMock.ProgramExistsCalls(), 1)
	assert.Equal(t, shellName, programQueryMock.ProgramExistsCalls()[0].Program)
}

func Test_ShellInstallationSucceeds_WhenPackageManagerSucceeds(t *testing.T) {
	// Arrange
	shellName := "zsh"

	programQueryMock := &osmanager.MoqProgramQuery{}

	pkgManagerMock := &pkgmanager.MoqPackageManager{
		InstallPackageFunc: func(pkg pkgmanager.RequestedPackageInfo) error {
			return nil
		},
	}

	installer := shell.NewDefaultShellInstaller(shellName, programQueryMock, pkgManagerMock, logger.DefaultLogger)

	// Act
	err := installer.Install(context.Background())

	// Assert
	require.NoError(t, err)
	assert.Len(t, pkgManagerMock.InstallPackageCalls(), 1)
	assert.Equal(t, shellName, pkgManagerMock.InstallPackageCalls()[0].RequestedPackageInfo.Name)
	assert.Nil(t, pkgManagerMock.InstallPackageCalls()[0].RequestedPackageInfo.VersionConstraints)
}

func Test_ShellInstallationFails_WhenPackageManagerFails(t *testing.T) {
	// Arrange
	shellName := "fish"

	programQueryMock := &osmanager.MoqProgramQuery{}

	pkgManagerMock := &pkgmanager.MoqPackageManager{
		InstallPackageFunc: func(pkg pkgmanager.RequestedPackageInfo) error {
			return assert.AnError
		},
	}

	installer := shell.NewDefaultShellInstaller(shellName, programQueryMock, pkgManagerMock, logger.DefaultLogger)

	// Act
	err := installer.Install(context.Background())

	// Assert
	require.Error(t, err)
	assert.Len(t, pkgManagerMock.InstallPackageCalls(), 1)
	assert.Equal(t, shellName, pkgManagerMock.InstallPackageCalls()[0].RequestedPackageInfo.Name)
	assert.Nil(t, pkgManagerMock.InstallPackageCalls()[0].RequestedPackageInfo.VersionConstraints)
}
