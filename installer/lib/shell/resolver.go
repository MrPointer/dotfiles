package shell

import (
	"fmt"
	"path/filepath"

	"github.com/MrPointer/dotfiles/installer/utils"
	"github.com/MrPointer/dotfiles/installer/utils/logger"
	"github.com/MrPointer/dotfiles/installer/utils/osmanager"
)

// ShellResolver resolves shell paths based on source and availability.
type ShellResolver interface {
	GetShellPath() (string, error)
	IsAvailable() (bool, error)
}

// DefaultShellResolver implements ShellResolver.
type DefaultShellResolver struct {
	shellName   string
	shellSource ShellSource
	brewPath    string
	osName      string
	osManager   osmanager.OsManager
	fileSystem  utils.FileSystem
	logger      logger.Logger
}

// NewDefaultShellResolver creates a new DefaultShellResolver.
func NewDefaultShellResolver(
	shellName string,
	shellSource ShellSource,
	brewPath string,
	osName string,
	osManager osmanager.OsManager,
	fileSystem utils.FileSystem,
	logger logger.Logger,
) *DefaultShellResolver {
	return &DefaultShellResolver{
		shellName:   shellName,
		shellSource: shellSource,
		brewPath:    brewPath,
		osName:      osName,
		osManager:   osManager,
		fileSystem:  fileSystem,
		logger:      logger,
	}
}

func (r *DefaultShellResolver) GetShellPath() (string, error) {
	switch r.shellSource {
	case ShellSourceBrew:
		return r.getBrewShellPath()
	case ShellSourceSystem:
		return r.getSystemShellPath()
	case ShellSourceAuto:
		return r.getAutoShellPath()
	default:
		return "", fmt.Errorf("unknown shell source: %s", r.shellSource)
	}
}

func (r *DefaultShellResolver) getBrewShellPath() (string, error) {
	if r.brewPath == "" {
		return "", fmt.Errorf("homebrew is not installed")
	}
	brewBinDir := filepath.Join(string(r.brewPath), "bin")
	shellPath := filepath.Join(brewBinDir, r.shellName)
	return shellPath, nil
}

func (r *DefaultShellResolver) getSystemShellPath() (string, error) {
	for _, dir := range knownSystemShellDirs(r.osName) {
		path := filepath.Join(dir, r.shellName)
		exists, err := r.fileSystem.PathExists(path)
		if err != nil {
			return "", fmt.Errorf("failed to check if %s exists: %w", path, err)
		}
		if exists {
			return path, nil
		}
	}
	return "", fmt.Errorf("shell %s not found in system directories", r.shellName)
}

func (r *DefaultShellResolver) getAutoShellPath() (string, error) {
	// Try brew first (if available)
	if r.brewPath != "" {
		brewPath, err := r.getBrewShellPath()
		if err == nil {
			exists, checkErr := r.fileSystem.PathExists(brewPath)
			if checkErr == nil && exists {
				r.logger.Debug("Found shell at brew path: %s", brewPath)
				return brewPath, nil
			}
		}
	}

	// Fall back to system
	r.logger.Debug("Brew shell not found, falling back to system shell")
	return r.getSystemShellPath()
}

func knownSystemShellDirs(osName string) []string {
	if osName == "darwin" {
		return []string{"/bin", "/usr/bin"}
	}
	// Linux and other Unix-like systems
	return []string{"/bin", "/usr/bin", "/usr/local/bin"}
}

func (r *DefaultShellResolver) IsAvailable() (bool, error) {
	shellPath, err := r.GetShellPath()
	if err != nil {
		// Shell not found in expected locations - this is not an error, just means it's not available
		return false, nil
	}
	exists, err := r.fileSystem.PathExists(shellPath)
	if err != nil {
		return false, fmt.Errorf("failed to check if shell exists at %s: %w", shellPath, err)
	}
	return exists, nil
}
