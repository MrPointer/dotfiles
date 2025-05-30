package brew

import (
	"fmt"
	"net/http" // Keep for http.StatusOK and potentially other http constants if needed by other functions
	"os"
	"os/user"
	"runtime"
	"strconv"
	"syscall"

	"github.com/MrPointer/dotfiles/installer/lib/compatibility"
	"github.com/MrPointer/dotfiles/installer/utils"
	"github.com/MrPointer/dotfiles/installer/utils/httpclient"
	"github.com/MrPointer/dotfiles/installer/utils/logger"
	"github.com/MrPointer/dotfiles/installer/utils/osmanager"
)

const (
	// BrewUserOnMultiUserSystem is the username used for Homebrew on multi-user systems
	BrewUserOnMultiUserSystem = "linuxbrew-manager"

	// LinuxBrewPath is the default installation path for Homebrew on Linux
	LinuxBrewPath = "/home/linuxbrew/.linuxbrew/bin/brew"

	// MacOSIntelBrewPath is the default installation path for Homebrew on Intel macOS
	MacOSIntelBrewPath = "/usr/local/bin/brew"

	// MacOSArmBrewPath is the default installation path for Homebrew on ARM macOS
	MacOSArmBrewPath = "/opt/homebrew/bin/brew"
)

// BrewInstaller defines the interface for Homebrew operations
// (moq is used for generating mocks in tests)
//
//go:generate moq -out brew_installer_moq.go . BrewInstaller
type BrewInstaller interface {
	DetectBrewPath() (string, error)
	IsAvailable() (bool, error)
	Install() error
}

// brewInstaller implements BrewInstaller for single-user systems
// Holds configuration options for Homebrew operations
type brewInstaller struct {
	logger           logger.Logger
	systemInfo       *compatibility.SystemInfo
	commander        utils.Commander
	httpClient       httpclient.HTTPClient
	osManager        osmanager.OsManager
	fs               utils.FileSystem
	brewPathOverride string // for testing only
}

var _ BrewInstaller = (*brewInstaller)(nil)

// MultiUserBrewInstaller implements BrewInstaller for multi-user systems
// It composes a regular brewInstaller and adds multi-user logic
type MultiUserBrewInstaller struct {
	*brewInstaller
	brewUser string
}

var _ BrewInstaller = (*MultiUserBrewInstaller)(nil)

// NewBrewInstaller creates a new BrewInstaller with the given options
func NewBrewInstaller(opts Options) BrewInstaller {
	base := &brewInstaller{
		logger:           opts.Logger,
		systemInfo:       opts.SystemInfo,
		commander:        opts.Commander,
		httpClient:       opts.HTTPClient,
		osManager:        opts.OsManager,
		fs:               opts.Fs,
		brewPathOverride: opts.BrewPathOverride,
	}

	if opts.MultiUserSystem {
		return &MultiUserBrewInstaller{
			brewInstaller: base,
			brewUser:      BrewUserOnMultiUserSystem,
		}
	}

	return base
}

// DetectBrewPath returns the appropriate brew binary path based on the system information
func (b *brewInstaller) DetectBrewPath() (string, error) {
	if b.brewPathOverride != "" {
		return b.brewPathOverride, nil
	}

	if b.systemInfo != nil {
		switch b.systemInfo.OSName {
		case "darwin":
			if b.systemInfo.Arch == "arm64" {
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

// IsAvailable checks if Homebrew is already installed and available (single-user)
func (b *brewInstaller) IsAvailable() (bool, error) {
	brewPath, err := b.DetectBrewPath()
	if err != nil {
		return false, err
	}

	_, err = os.Stat(brewPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// Install installs Homebrew if not already installed (single-user)
func (b *brewInstaller) Install() error {
	isAvailable, err := b.IsAvailable()
	if err != nil {
		return fmt.Errorf("failed checking Homebrew availability: %w", err)
	}

	if isAvailable {
		b.logger.Success("Homebrew is already installed")
		return nil
	}

	b.logger.Info("Installing Homebrew")
	err = b.installHomebrew("")
	if err != nil {
		return err
	}

	// Self-validation: check that brew is available and works
	if err := b.validateInstall(); err != nil {
		return fmt.Errorf("brew self-validation failed: %w", err)
	}

	return nil
}

// validateInstall checks that the brew binary exists and is functional
func (b *brewInstaller) validateInstall() error {
	brewPath, err := b.DetectBrewPath()
	if err != nil {
		return fmt.Errorf("could not detect brew path: %w", err)
	}

	info, err := os.Stat(brewPath)
	if err != nil {
		return fmt.Errorf("brew binary not found at %s: %w", brewPath, err)
	}
	if info.Mode()&0111 == 0 {
		return fmt.Errorf("brew binary at %s is not executable", brewPath)
	}

	// Try running 'brew --version' to verify it works
	if b.commander != nil {
		err = b.commander.Run(brewPath, "--version")
		if err != nil {
			return fmt.Errorf("brew --version failed: %w", err)
		}
	}

	return nil
}

// Install installs Homebrew if not already installed (multi-user)
func (m *MultiUserBrewInstaller) Install() error {
	isAvailable, err := m.IsAvailable()
	if err != nil {
		return fmt.Errorf("failed checking Homebrew availability: %w", err)
	}

	if isAvailable {
		m.logger.Success("Homebrew is already installed (multi-user)")
		return nil
	}

	m.logger.Info("Installing Homebrew (multi-user)")
	if m.systemInfo.OSName == "darwin" {
		return fmt.Errorf("multi-user Homebrew installation is not supported on macOS, please install manually")
	}

	err = m.installMultiUserLinux()
	if err != nil {
		return err
	}

	// Self-validation: check that brew is available and works
	if err := m.validateInstall(); err != nil {
		return fmt.Errorf("brew self-validation failed: %w", err)
	}

	return nil
}

// Multi-user overrides
// IsAvailable checks if Homebrew is already installed and available (multi-user)
func (m *MultiUserBrewInstaller) IsAvailable() (bool, error) {
	brewPath, err := m.DetectBrewPath()
	if err != nil {
		return false, err
	}

	fileInfo, err := os.Stat(brewPath)
	if err != nil {
		return false, err
	}

	if m.systemInfo != nil && m.systemInfo.OSName == "linux" ||
		m.systemInfo == nil && runtime.GOOS == "linux" {
		stat, ok := fileInfo.Sys().(*syscall.Stat_t)
		if !ok {
			return false, fmt.Errorf("failed to get file info: %w", err)
		}

		brewUser, err := user.LookupId(strconv.FormatUint(uint64(stat.Uid), 10))
		if err != nil {
			return false, fmt.Errorf("failed to lookup user: %w", err)
		}
		return brewUser.Username == m.brewUser, nil
	}

	return true, nil
}

// installMultiUserLinux installs Homebrew in a multi-user configuration on Linux
func (m *MultiUserBrewInstaller) installMultiUserLinux() error {
	brewUser := BrewUserOnMultiUserSystem
	brewHome := "/home/linuxbrew"

	// 1. Check if user exists and create if needed
	exists, err := m.osManager.UserExists(brewUser)
	if err != nil {
		return fmt.Errorf("error checking if user '%s' exists: %w", brewUser, err)
	}

	if !exists {
		if err := m.osManager.AddUser(brewUser); err != nil {
			return fmt.Errorf("error creating user '%s': %w", brewUser, err)
		}
	}

	// 2. Add user to sudo group
	if err := m.osManager.AddUserToGroup(brewUser, "sudo"); err != nil {
		m.logger.Debug("Note: Failed to add user to sudo group, continuing anyway")
	}

	// 3. Add passwordless sudo for brew user
	if err := m.osManager.AddSudoAccess(brewUser); err != nil {
		return fmt.Errorf("failed to add sudo access for user '%s': %w", brewUser, err)
	}

	// 4. Set ownership of homebrew directory
	if err := m.osManager.SetOwnership(brewHome, brewUser); err != nil {
		return fmt.Errorf("failed to set ownership of '%s' to '%s': %w", brewHome, brewUser, err)
	}

	// 5. Install Homebrew as the brew user
	return m.installHomebrew(brewUser)
}

// isRoot returns true if the current user is root
func isRoot() bool {
	return os.Geteuid() == 0
}

// installHomebrew handles both regular and multi-user Homebrew installations
func (b *brewInstaller) installHomebrew(asUser string) error {
	// Download and prepare the installation script
	installScriptPath, cleanup, err := b.downloadAndPrepareInstallScript()
	if err != nil {
		return err
	}
	defer cleanup() // Ensure the temporary script is removed after execution

	exists, err := b.fs.PathExists(installScriptPath)
	if err != nil {
		return fmt.Errorf("failed checking if install script exists: %w", err)
	}
	if !exists {
		return fmt.Errorf("install script does not exist at %s", installScriptPath)
	}
	b.logger.Debug("Homebrew install script downloaded to %s", installScriptPath)

	// Execute the downloaded install script, optionally as a different user
	if asUser != "" {
		b.logger.Info("Running Homebrew install script as %s", asUser)
		err := b.commander.Run("sudo", "-Hu", asUser, "bash", installScriptPath)
		if err != nil {
			return fmt.Errorf("failed running Homebrew install script as %s: %w", asUser, err)
		}

		b.logger.Success("Successfully installed Homebrew for user %s", asUser)
	} else {
		b.logger.Info("Running Homebrew install script")
		err := b.commander.RunWithEnv(map[string]string{"NONINTERACTIVE": "1"}, "/bin/bash", installScriptPath)
		if err != nil {
			return err
		}

		b.logger.Success("Homebrew installed successfully")
	}

	return nil
}

// downloadAndPrepareInstallScript downloads the Homebrew installation script and prepares it for execution
func (b *brewInstaller) downloadAndPrepareInstallScript() (string, func(), error) {
	b.logger.Info("Downloading Homebrew install script")
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
	b.logger.Debug("Creating temporary file for Homebrew install script")
	tempFilePath, err := b.fs.CreateTemporaryFile("", "brew-install-*.sh")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temporary file for Homebrew install script: %w", err)
	}

	// Create a cleanup function to remove the temp file
	cleanup := func() {
		b.logger.Debug("Cleaning up temporary file for Homebrew install script")
		err := b.fs.RemovePath(tempFilePath)
		if err != nil {
			b.logger.Warning("Failed to remove temporary file: %w", err)
		}
	}

	// Copy the script to the temp file
	b.logger.Debug("Writing Homebrew install script to temporary file")
	bytesWritten, err := b.fs.WriteFile(tempFilePath, resp.Body)
	if err != nil {
		cleanup()
		return "", nil, fmt.Errorf("failed to write Homebrew install script to temporary file: %w", err)
	}
	if bytesWritten == 0 {
		cleanup()
		return "", nil, fmt.Errorf("failed to write Homebrew install script: no bytes written")
	}

	// Make the script executable
	b.logger.Debug("Making Homebrew install script executable")
	if err = os.Chmod(tempFilePath, 0777); err != nil {
		cleanup()
		return "", nil, fmt.Errorf("failed to make Homebrew install script executable: %w", err)
	}

	return tempFilePath, cleanup, nil
}

// Options holds configuration options for Homebrew operations
type Options struct {
	Logger           logger.Logger
	MultiUserSystem  bool
	SystemInfo       *compatibility.SystemInfo
	Commander        utils.Commander
	HTTPClient       httpclient.HTTPClient
	OsManager        osmanager.OsManager
	Fs               utils.FileSystem
	BrewPathOverride string
}

// DefaultOptions returns the default options
func DefaultOptions() Options {
	return Options{
		MultiUserSystem:  false,
		Logger:           logger.DefaultLogger,
		SystemInfo:       nil,
		Commander:        utils.NewDefaultCommander(),
		HTTPClient:       httpclient.NewDefaultHTTPClient(),
		OsManager:        osmanager.NewUnixOsManager(logger.DefaultLogger, utils.NewDefaultCommander(), isRoot()),
		Fs:               utils.NewDefaultFileSystem(),
		BrewPathOverride: "",
	}
}

// WithLogger sets a custom logger for the brew operations
func (o *Options) WithLogger(log logger.Logger) *Options {
	o.Logger = log
	return o
}

// WithMultiUserSystem configures for multi-user system operation
func (o *Options) WithMultiUserSystem(multiUser bool) *Options {
	o.MultiUserSystem = multiUser
	return o
}

// WithSystemInfo sets system information for brew operations
func (o *Options) WithSystemInfo(sysInfo *compatibility.SystemInfo) *Options {
	o.SystemInfo = sysInfo
	return o
}

// WithCommander sets a custom Commander for Homebrew operations
func (o *Options) WithCommander(cmdr utils.Commander) *Options {
	o.Commander = cmdr
	return o
}

// WithHTTPClient sets a custom HTTP client for Homebrew operations
func (o *Options) WithHTTPClient(client httpclient.HTTPClient) *Options {
	o.HTTPClient = client
	return o
}

// WithOsManager sets a custom OS manager for Homebrew operations
func (o *Options) WithOsManager(osMgr osmanager.OsManager) *Options {
	o.OsManager = osMgr
	return o
}

// WithFileSystem sets a custom FileSystem for Homebrew operations
func (o *Options) WithFileSystem(fs utils.FileSystem) *Options {
	o.Fs = fs
	return o
}
