package gpg

import (
	"context"
	"errors"
	"strings"

	"github.com/Masterminds/semver"

	"github.com/MrPointer/dotfiles/installer/lib/compatibility"
	"github.com/MrPointer/dotfiles/installer/utils"
	"github.com/MrPointer/dotfiles/installer/utils/logger"
	"github.com/MrPointer/dotfiles/installer/utils/osmanager"
)

type GpgClientInstaller interface {
	IsAvailable() (bool, error)
	Install(ctx context.Context) error
}

type gpgInstaller struct {
	systemInfo *compatibility.SystemInfo
	logger     logger.Logger
	commander  utils.Commander
	osManager  osmanager.OsManager
}

func NewGpgInstaller(
	systemInfo *compatibility.SystemInfo,
	logger logger.Logger,
	commander utils.Commander,
	osManager osmanager.OsManager,
) GpgClientInstaller {
	return &gpgInstaller{
		systemInfo: systemInfo,
		logger:     logger,
		commander:  commander,
		osManager:  osManager,
	}
}

func (g *gpgInstaller) IsAvailable() (bool, error) {
	// Check if gpg is available.
	gpgExists, err := g.osManager.ProgramExists("gpg")
	if err != nil {
		return false, err
	}
	if !gpgExists {
		g.logger.Warning("GPG is not available. Required for GPG operations.")
		return false, nil
	}

	// gpg is available, now ensure its version is compatible (we required anything above 2.2).
	versionMatches, err := gpgVersionMatches(g)
	if err != nil {
		return false, err
	}
	if !versionMatches {
		g.logger.Warning("GPG version is not compatible. Required version is >=2.2.0")
		return false, nil
	}

	gpgAgentExists, err := g.osManager.ProgramExists("gpg-agent")
	if err != nil {
		return false, err
	}
	if !gpgAgentExists {
		g.logger.Warning("GPG agent is not available. Required for GPG operations.")
		return false, nil
	}

	return true, nil
}

func gpgVersionMatches(g *gpgInstaller) (bool, error) {
	gpgVersion, err := g.osManager.GetProgramVersion("gpg", extractGpgVersion)
	if err != nil {
		return false, err
	}
	constraints, err := semver.NewConstraint(">=2.2.0")
	if err != nil {
		return false, err
	}
	version, err := semver.NewVersion(gpgVersion)
	if err != nil {
		return false, err
	}
	if !constraints.Check(version) {
		return false, nil
	}

	return true, nil
}

func extractGpgVersion(rawVersion string) (string, error) {
	// Extract the version number from the raw version string.
	// Take the first row, split by space, and return the 3rd element (the version number).
	lines := strings.Split(rawVersion, "\n")
	if len(lines) == 1 {
		return "", errors.New("line count is 1, meaning there are no newlines in the version string")
	}

	const minimumRequiredElements = 3
	parts := strings.Split(lines[0], " ")
	if len(parts) < minimumRequiredElements {
		return "", errors.New("version string does not contain enough parts to extract version")
	}

	return parts[2], nil
}

func (g *gpgInstaller) Install(ctx context.Context) error {
	// Implementation to install GPG
	return nil
}
