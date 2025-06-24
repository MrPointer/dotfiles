package gpg

import (
	"context"
	"errors"
	"strings"

	"github.com/Masterminds/semver"

	"github.com/MrPointer/dotfiles/installer/lib/pkgmanager"
	"github.com/MrPointer/dotfiles/installer/utils"
	"github.com/MrPointer/dotfiles/installer/utils/logger"
	"github.com/MrPointer/dotfiles/installer/utils/osmanager"
)

const supportedGpgVersionConstraintString = ">=2.2.0"

type GpgClientInstaller interface {
	IsAvailable() (bool, error)
	Install(ctx context.Context) error
}

type gpgInstaller struct {
	logger         logger.Logger
	commander      utils.Commander
	osManager      osmanager.OsManager
	packageManager pkgmanager.PackageManager
}

func NewGpgInstaller(
	logger logger.Logger,
	commander utils.Commander,
	osManager osmanager.OsManager,
	packageManager pkgmanager.PackageManager,
) GpgClientInstaller {
	return &gpgInstaller{
		logger:         logger,
		commander:      commander,
		osManager:      osManager,
		packageManager: packageManager,
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

	constraints, err := semver.NewConstraint(supportedGpgVersionConstraintString)
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
	versionConstraints, err := semver.NewConstraint(supportedGpgVersionConstraintString)
	if err != nil {
		return errors.New("failed to create version constraints: " + err.Error())
	}

	err = g.packageManager.InstallPackage(pkgmanager.NewRequestedPackageInfo("gpg", versionConstraints))
	if err != nil {
		return errors.New("failed to install GPG client: " + err.Error())
	}

	return nil
}
