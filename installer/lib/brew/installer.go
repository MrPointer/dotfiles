package brew

import (
	"fmt"
	"net/http" // Keep for http.StatusOK and potentially other http constants if needed by other functions
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"slices"

	"github.com/MrPointer/dotfiles/installer/lib/compatibility"
	"github.com/MrPointer/dotfiles/installer/utils"
	"github.com/MrPointer/dotfiles/installer/utils/httpclient"
	"github.com/MrPointer/dotfiles/installer/utils/logger"
	"github.com/MrPointer/dotfiles/installer/utils/osmanager"
	"github.com/MrPointer/dotfiles/installer/utils/privilege"
)

const (
	// LinuxBrewPath is the default installation path for Homebrew on Linux.
	LinuxBrewPath = "/home/linuxbrew/.linuxbrew/bin/brew"

	// MacOSIntelBrewPath is the default installation path for Homebrew on Intel macOS.
	MacOSIntelBrewPath = "/usr/local/bin/brew"

	// MacOSArmBrewPath is the default installation path for Homebrew on ARM macOS.
	MacOSArmBrewPath = "/opt/homebrew/bin/brew"
)

// DetectBrewPath returns the appropriate brew binary path based on the system information.
func DetectBrewPath(systemInfo *compatibility.SystemInfo, pathOverride string) (string, error) {
	if pathOverride != "" {
		return pathOverride, nil
	}

	if systemInfo != nil {
		switch systemInfo.OSName {
		case "darwin":
			if systemInfo.Arch == "arm64" {
				return MacOSArmBrewPath, nil
			}

			return MacOSIntelBrewPath, nil

		case "linux":
			return LinuxBrewPath, nil

		default:
			return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
		}
	}

	return "", fmt.Errorf("system information is not provided")
}

// UpdatePathWithBrewBinaries updates the current process's PATH to include brew's binary directories.
// This function is primarily used internally by BrewInstaller but can be called directly if needed.
func UpdatePathWithBrewBinaries(brewPath string) error {
	// Extract the bin directory from the brew path
	// e.g., /opt/homebrew/bin/brew -> /opt/homebrew/bin
	brewBinDir := filepath.Dir(brewPath)

	// Get current PATH
	currentPath := os.Getenv("PATH")

	// Check if already in PATH to avoid duplicates
	pathDirs := strings.Split(currentPath, string(os.PathListSeparator))
	if slices.Contains(pathDirs, brewBinDir) {
		return nil // Already in PATH
	}

	// Prepend brew bin directory to PATH
	newPath := brewBinDir + string(os.PathListSeparator) + currentPath
	return os.Setenv("PATH", newPath)
}

// BrewInstaller defines the interface for Homebrew operations.
// (moq is used for generating mocks in tests.)
type BrewInstaller interface {
	IsAvailable() (bool, error)
	Install() error
}

// brewInstaller implements BrewInstaller for single-user systems.
// Holds configuration options for Homebrew operations.
type brewInstaller struct {
	logger           logger.Logger
	systemInfo       *compatibility.SystemInfo
	commander        utils.Commander
	httpClient       httpclient.HTTPClient
	osManager        osmanager.OsManager
	fs               utils.FileSystem
	brewPathOverride string // for testing only
	displayMode      utils.DisplayMode
}

var _ BrewInstaller = (*brewInstaller)(nil)

// NewBrewInstaller creates a new BrewInstaller with the given options.
func NewBrewInstaller(opts Options) BrewInstaller {
	return &brewInstaller{
		logger:           opts.Logger,
		systemInfo:       opts.SystemInfo,
		commander:        opts.Commander,
		httpClient:       opts.HTTPClient,
		osManager:        opts.OsManager,
		fs:               opts.Fs,
		brewPathOverride: opts.BrewPathOverride,
		displayMode:      opts.DisplayMode,
	}
}

// DetectBrewPath returns the appropriate brew binary path based on the system information.
func (b *brewInstaller) DetectBrewPath() (string, error) {
	return DetectBrewPath(b.systemInfo, b.brewPathOverride)
}

// IsAvailable checks if Homebrew is already installed and available (single-user).
func (b *brewInstaller) IsAvailable() (bool, error) {
	b.logger.Debug("Checking if Homebrew is available")

	brewPath, err := b.DetectBrewPath()
	if err != nil {
		return false, err
	}

	exists, err := b.fs.PathExists(brewPath)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// Install installs Homebrew if not already installed.
func (b *brewInstaller) Install() error {
	isAvailable, err := b.IsAvailable()
	if err != nil {
		return fmt.Errorf("failed checking Homebrew availability: %w", err)
	}

	if isAvailable {
		b.logger.Debug("Homebrew is already installed")

		// Update PATH to include brew binaries so installed tools can be found
		brewPath, err := b.DetectBrewPath()
		if err != nil {
			return fmt.Errorf("failed to detect brew path for PATH update: %w", err)
		}
		if err := UpdatePathWithBrewBinaries(brewPath); err != nil {
			b.logger.Warning("Failed to update PATH with brew binaries: %v", err)
		}

		return nil
	}

	b.logger.Debug("Installing Homebrew")
	err = b.installHomebrew()
	if err != nil {
		return err
	}

	// Self-validation: check that brew is available and works
	if err := b.validateInstall(); err != nil {
		return fmt.Errorf("brew self-validation failed: %w", err)
	}

	// Update PATH to include brew binaries so installed tools can be found
	brewPath, err := b.DetectBrewPath()
	if err != nil {
		return fmt.Errorf("failed to detect brew path for PATH update: %w", err)
	}
	if err := UpdatePathWithBrewBinaries(brewPath); err != nil {
		b.logger.Warning("Failed to update PATH with brew binaries: %v", err)
	}

	return nil
}

// validateInstall checks that the brew binary exists and is functional.
func (b *brewInstaller) validateInstall() error {
	brewPath, err := b.DetectBrewPath()
	if err != nil {
		return fmt.Errorf("could not detect brew path: %w", err)
	}

	exists, err := b.fs.PathExists(brewPath)
	if err != nil {
		return fmt.Errorf("failed to check if brew binary exists: %w", err)
	}
	if !exists {
		return fmt.Errorf("brew binary not found at %s", brewPath)
	}

	// Try running 'brew --version' to verify it works
	if b.commander != nil {
		var discardOutputOption utils.Option = utils.EmptyOption()
		if b.displayMode != utils.DisplayModePassthrough {
			discardOutputOption = utils.WithDiscardOutput()
		}

		_, err = b.commander.RunCommand(brewPath, []string{"--version"}, discardOutputOption)
		if err != nil {
			return fmt.Errorf("brew --version failed: %w", err)
		}
	}

	return nil
}

// installHomebrew handles Homebrew installation.
func (b *brewInstaller) installHomebrew() error {
	// Download and prepare the installation script
	installScriptPath, cleanup, err := b.downloadAndPrepareInstallScript()
	if err != nil {
		return err
	}
	defer cleanup() // Ensure the temporary script is removed after execution

	b.logger.Debug("Downloading Homebrew install script")
	exists, err := b.fs.PathExists(installScriptPath)
	if err != nil {
		return fmt.Errorf("failed checking if install script exists: %w", err)
	}
	if !exists {
		return fmt.Errorf("install script does not exist at %s", installScriptPath)
	}
	b.logger.Debug("Successfully downloaded Homebrew install script")
	b.logger.Trace("Homebrew install script downloaded to %s", installScriptPath)

	// Execute the downloaded install script
	b.logger.Debug("Running Homebrew install script")

	var discardOutputOption utils.Option = utils.EmptyOption()
	if b.displayMode != utils.DisplayModePassthrough {
		discardOutputOption = utils.WithDiscardOutput()
	}

	_, err = b.commander.RunCommand("/bin/bash", []string{installScriptPath}, discardOutputOption, utils.WithEnv(map[string]string{"NONINTERACTIVE": "1"}))
	if err != nil {
		return err
	}

	b.logger.Debug("Homebrew installed successfully")

	return nil
}

// downloadAndPrepareInstallScript downloads the Homebrew installation script and prepares it for execution.
func (b *brewInstaller) downloadAndPrepareInstallScript() (string, func(), error) {
	installScriptURL := "https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh"

	resp, err := b.httpClient.Get(installScriptURL) // Changed to use httpClient
	if err != nil {
		return "", nil, fmt.Errorf("failed to download Homebrew install script: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", nil, fmt.Errorf("failed to download Homebrew install script: HTTP status %d", resp.StatusCode)
	}

	// Create a temporary file for the install script
	b.logger.Trace("Creating temporary file for Homebrew install script")
	tempFilePath, err := b.fs.CreateTemporaryFile("", "brew-install-*.sh")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temporary file for Homebrew install script: %w", err)
	}

	// Create a cleanup function to remove the temp file
	cleanup := func() {
		b.logger.Trace("Cleaning up temporary file for Homebrew install script")
		err := b.fs.RemovePath(tempFilePath)
		if err != nil {
			b.logger.Trace("Failed to remove temporary file: %w", err)
		}
	}

	// Copy the script to the temp file
	b.logger.Trace("Writing Homebrew install script to temporary file")
	bytesWritten, err := b.fs.WriteFile(tempFilePath, resp.Body)
	if err != nil {
		cleanup()
		return "", nil, fmt.Errorf("failed to write Homebrew install script to temporary file: %w", err)
	}
	if bytesWritten == 0 {
		cleanup()
		return "", nil, fmt.Errorf("failed to write Homebrew install script: no bytes written")
	}

	// Make the script executable.
	const permissions = 0o755 // Standard executable permissions
	b.logger.Trace("Making Homebrew install script executable")
	if err = b.osManager.SetPermissions(tempFilePath, permissions); err != nil {
		cleanup()
		return "", nil, fmt.Errorf("failed to make Homebrew install script executable: %w", err)
	}

	return tempFilePath, cleanup, nil
}

// Options holds configuration options for Homebrew operations.
type Options struct {
	Logger           logger.Logger
	SystemInfo       *compatibility.SystemInfo
	Commander        utils.Commander
	HTTPClient       httpclient.HTTPClient
	OsManager        osmanager.OsManager
	Fs               utils.FileSystem
	BrewPathOverride string
	DisplayMode      utils.DisplayMode
}

// DefaultOptions returns the default options.
func DefaultOptions() *Options {
	commander := utils.NewDefaultCommander(logger.DefaultLogger)
	escalator := privilege.NewDefaultEscalator(logger.DefaultLogger, commander, utils.NewGoNativeProgramQuery())
	fileSystem := utils.NewDefaultFileSystem()

	return &Options{
		Logger:           logger.DefaultLogger,
		SystemInfo:       nil,
		Commander:        utils.NewDefaultCommander(logger.DefaultLogger),
		HTTPClient:       httpclient.NewDefaultHTTPClient(),
		OsManager:        osmanager.NewUnixOsManager(logger.DefaultLogger, commander, escalator, fileSystem),
		Fs:               fileSystem,
		BrewPathOverride: "",
		DisplayMode:      utils.DisplayModeProgress,
	}
}

func (o *Options) WithBrewPathOverride(path string) *Options {
	o.BrewPathOverride = path
	return o
}

// WithLogger sets a custom logger for the brew operations.
func (o *Options) WithLogger(log logger.Logger) *Options {
	o.Logger = log
	return o
}

// WithSystemInfo sets system information for brew operations.
func (o *Options) WithSystemInfo(sysInfo *compatibility.SystemInfo) *Options {
	o.SystemInfo = sysInfo
	return o
}

// WithCommander sets a custom Commander for Homebrew operations.
func (o *Options) WithCommander(cmdr utils.Commander) *Options {
	o.Commander = cmdr
	return o
}

// WithHTTPClient sets a custom HTTP client for Homebrew operations.
func (o *Options) WithHTTPClient(client httpclient.HTTPClient) *Options {
	o.HTTPClient = client
	return o
}

// WithOsManager sets a custom OS manager for Homebrew operations.
func (o *Options) WithOsManager(osMgr osmanager.OsManager) *Options {
	o.OsManager = osMgr
	return o
}

// WithFileSystem sets a custom FileSystem for Homebrew operations.
func (o *Options) WithFileSystem(fs utils.FileSystem) *Options {
	o.Fs = fs
	return o
}

// WithDisplayMode sets the display mode for external tool output.
func (o *Options) WithDisplayMode(mode utils.DisplayMode) *Options {
	o.DisplayMode = mode
	return o
}
