package gpg_test

import (
	"testing"

	"github.com/MrPointer/dotfiles/installer/lib/compatibility"
	"github.com/MrPointer/dotfiles/installer/lib/gpg"
	"github.com/MrPointer/dotfiles/installer/utils"
	"github.com/MrPointer/dotfiles/installer/utils/logger"
	"github.com/MrPointer/dotfiles/installer/utils/osmanager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_GpgIsReportedAsUnavailable_WhenGpgIsNotInstalled(t *testing.T) {
	// Arrange
	systemInfo := &compatibility.SystemInfo{OSName: "linux", Arch: "amd64"}

	loggerMock := &logger.MoqLogger{
		WarningFunc: func(format string, args ...any) {
			assert.Equal(t, "GPG is not available. Required for GPG operations.", format)
		},
	}

	commanderMock := &utils.MoqCommander{}

	osManagerMock := &osmanager.MoqOsManager{
		ProgramExistsFunc: func(program string) (bool, error) {
			if program == "gpg" {
				return false, nil
			}
			return true, nil
		},
	}

	installer := gpg.NewGpgInstaller(systemInfo, loggerMock, commanderMock, osManagerMock)

	// Act
	available, err := installer.IsAvailable()

	// Assert
	require.NoError(t, err)
	assert.False(t, available)
	assert.Len(t, osManagerMock.ProgramExistsCalls(), 1)
	assert.Equal(t, "gpg", osManagerMock.ProgramExistsCalls()[0].Program)
	assert.Len(t, loggerMock.WarningCalls(), 1)
}

func Test_GpgAvailabilityCheckFails_WhenGpgProgramExistsFails(t *testing.T) {
	// Arrange
	systemInfo := &compatibility.SystemInfo{OSName: "linux", Arch: "amd64"}

	loggerMock := &logger.MoqLogger{}
	commanderMock := &utils.MoqCommander{}

	osManagerMock := &osmanager.MoqOsManager{
		ProgramExistsFunc: func(program string) (bool, error) {
			return false, assert.AnError
		},
	}

	installer := gpg.NewGpgInstaller(systemInfo, loggerMock, commanderMock, osManagerMock)

	// Act
	available, err := installer.IsAvailable()

	// Assert
	require.Error(t, err)
	assert.False(t, available)
	assert.Len(t, osManagerMock.ProgramExistsCalls(), 1)
	assert.Equal(t, "gpg", osManagerMock.ProgramExistsCalls()[0].Program)
}

func Test_GpgIsReportedAsUnavailable_WhenGpgVersionIsIncompatible(t *testing.T) {
	// Arrange
	systemInfo := &compatibility.SystemInfo{OSName: "linux", Arch: "amd64"}

	loggerMock := &logger.MoqLogger{
		WarningFunc: func(format string, args ...any) {
			assert.Equal(t, "GPG version is not compatible. Required version is >=2.2.0", format)
		},
	}

	commanderMock := &utils.MoqCommander{}

	osManagerMock := &osmanager.MoqOsManager{
		ProgramExistsFunc: func(program string) (bool, error) {
			return true, nil
		},
		GetProgramVersionFunc: func(program string, versionExtractor osmanager.VersionExtractor, queryArgs ...string) (string, error) {
			return "2.1.0", nil
		},
	}

	installer := gpg.NewGpgInstaller(systemInfo, loggerMock, commanderMock, osManagerMock)

	// Act
	available, err := installer.IsAvailable()

	// Assert
	require.NoError(t, err)
	assert.False(t, available)
	assert.Len(t, osManagerMock.ProgramExistsCalls(), 1)
	assert.Len(t, osManagerMock.GetProgramVersionCalls(), 1)
	assert.Equal(t, "gpg", osManagerMock.GetProgramVersionCalls()[0].Program)
	assert.Len(t, loggerMock.WarningCalls(), 1)
}

func Test_GpgAvailabilityCheckFails_WhenGetProgramVersionFails(t *testing.T) {
	// Arrange
	systemInfo := &compatibility.SystemInfo{OSName: "linux", Arch: "amd64"}

	loggerMock := &logger.MoqLogger{}
	commanderMock := &utils.MoqCommander{}

	osManagerMock := &osmanager.MoqOsManager{
		ProgramExistsFunc: func(program string) (bool, error) {
			return true, nil
		},
		GetProgramVersionFunc: func(program string, versionExtractor osmanager.VersionExtractor, queryArgs ...string) (string, error) {
			return "", assert.AnError
		},
	}

	installer := gpg.NewGpgInstaller(systemInfo, loggerMock, commanderMock, osManagerMock)

	// Act
	available, err := installer.IsAvailable()

	// Assert
	require.Error(t, err)
	assert.False(t, available)
	assert.Len(t, osManagerMock.ProgramExistsCalls(), 1)
	assert.Len(t, osManagerMock.GetProgramVersionCalls(), 1)
}

func Test_GpgIsReportedAsUnavailable_WhenGpgAgentIsNotInstalled(t *testing.T) {
	// Arrange
	systemInfo := &compatibility.SystemInfo{OSName: "linux", Arch: "amd64"}

	loggerMock := &logger.MoqLogger{
		WarningFunc: func(format string, args ...any) {
			assert.Equal(t, "GPG agent is not available. Required for GPG operations.", format)
		},
	}

	commanderMock := &utils.MoqCommander{}

	osManagerMock := &osmanager.MoqOsManager{
		ProgramExistsFunc: func(program string) (bool, error) {
			if program == "gpg-agent" {
				return false, nil
			}
			return true, nil
		},
		GetProgramVersionFunc: func(program string, versionExtractor osmanager.VersionExtractor, queryArgs ...string) (string, error) {
			return "2.3.0", nil
		},
	}

	installer := gpg.NewGpgInstaller(systemInfo, loggerMock, commanderMock, osManagerMock)

	// Act
	available, err := installer.IsAvailable()

	// Assert
	require.NoError(t, err)
	assert.False(t, available)
	assert.Len(t, osManagerMock.ProgramExistsCalls(), 2)
	assert.Equal(t, "gpg-agent", osManagerMock.ProgramExistsCalls()[1].Program)
	assert.Len(t, loggerMock.WarningCalls(), 1)
}

func Test_GpgAvailabilityCheckFails_WhenGpgAgentProgramExistsFails(t *testing.T) {
	// Arrange
	systemInfo := &compatibility.SystemInfo{OSName: "linux", Arch: "amd64"}

	loggerMock := &logger.MoqLogger{}
	commanderMock := &utils.MoqCommander{}

	callCount := 0
	osManagerMock := &osmanager.MoqOsManager{
		ProgramExistsFunc: func(program string) (bool, error) {
			callCount++
			if callCount == 1 {
				return true, nil // First call for gpg
			}
			return false, assert.AnError // Second call for gpg-agent
		},
		GetProgramVersionFunc: func(program string, versionExtractor osmanager.VersionExtractor, queryArgs ...string) (string, error) {
			return "2.3.0", nil
		},
	}

	installer := gpg.NewGpgInstaller(systemInfo, loggerMock, commanderMock, osManagerMock)

	// Act
	available, err := installer.IsAvailable()

	// Assert
	require.Error(t, err)
	assert.False(t, available)
	assert.Len(t, osManagerMock.ProgramExistsCalls(), 2)
	assert.Equal(t, "gpg", osManagerMock.ProgramExistsCalls()[0].Program)
	assert.Equal(t, "gpg-agent", osManagerMock.ProgramExistsCalls()[1].Program)
}

func Test_GpgIsReportedAsAvailable_WhenAllRequirementsAreMet(t *testing.T) {
	// Arrange
	systemInfo := &compatibility.SystemInfo{OSName: "linux", Arch: "amd64"}

	loggerMock := &logger.MoqLogger{}
	commanderMock := &utils.MoqCommander{}

	osManagerMock := &osmanager.MoqOsManager{
		ProgramExistsFunc: func(program string) (bool, error) {
			return true, nil
		},
		GetProgramVersionFunc: func(program string, versionExtractor osmanager.VersionExtractor, queryArgs ...string) (string, error) {
			return "2.4.0", nil
		},
	}

	installer := gpg.NewGpgInstaller(systemInfo, loggerMock, commanderMock, osManagerMock)

	// Act
	available, err := installer.IsAvailable()

	// Assert
	require.NoError(t, err)
	assert.True(t, available)
	assert.Len(t, osManagerMock.ProgramExistsCalls(), 2)
	assert.Equal(t, "gpg", osManagerMock.ProgramExistsCalls()[0].Program)
	assert.Equal(t, "gpg-agent", osManagerMock.ProgramExistsCalls()[1].Program)
}
