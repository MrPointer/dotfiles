package brew

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"runtime"
	"strconv"
	"syscall"

	"github.com/MrPointer/dotfiles/installer/lib/compatibility"
	"github.com/MrPointer/dotfiles/installer/utils"
	"github.com/MrPointer/dotfiles/installer/utils/logger"
)

const (
	// BrewUserOnMultiUserSystem is the username used for Homebrew on multi-user systems
	BrewUserOnMultiUserSystem = "linuxbrew-manager"

	// LinuxBrewPath is the default installation path for Homebrew on Linux
	LinuxBrewPath = "/home/linuxbrew/.linuxbrew/bin/brew"

	// MacOSIntelBrewPath is the default installation path for Homebrew on Intel macOS
	MacOSIntelBrewPath = "/usr/local/bin/brew"

	// MacOSARMBrewPath is the default installation path for Homebrew on ARM macOS
	MacOSARMBrewPath = "/opt/homebrew/bin/brew"
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
		logger:     opts.Logger,
		systemInfo: opts.SystemInfo,
		commander:  opts.Commander,
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
				return MacOSARMBrewPath, nil
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
	return b.installHomebrew("")
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
	if runtime.GOOS == "darwin" {
		return fmt.Errorf("multi-user Homebrew installation is not supported on macOS, please install manually")
	}

	return m.installMultiUserLinux()
}

// installMultiUserLinux installs Homebrew in a multi-user configuration on Linux
func (m *MultiUserBrewInstaller) installMultiUserLinux() error {
	// Check if the brew user already exists
	_, err := user.Lookup(BrewUserOnMultiUserSystem)
	if err != nil {
		// Create the brew user
		m.logger.Info("Creating user '%s' for Homebrew", BrewUserOnMultiUserSystem)

		err := m.commander.Run("sudo", "useradd", "-m", "-p", "", BrewUserOnMultiUserSystem)
		if err != nil {
			return fmt.Errorf("failed creating user '%s' for Homebrew: %w", BrewUserOnMultiUserSystem, err)
		}

		// Add user to sudo group
		err = m.commander.Run("sudo", "usermod", "-aG", "sudo", BrewUserOnMultiUserSystem)
		if err != nil {
			return fmt.Errorf("failed adding user '%s' to sudo group: %w", BrewUserOnMultiUserSystem, err)
		}

		// Add user to passwordless sudo
		sudoersLine := fmt.Sprintf("%s ALL=(ALL) NOPASSWD:ALL", BrewUserOnMultiUserSystem)
		err = m.commander.Run("sudo", "bash", "-c", fmt.Sprintf("echo '%s' | sudo tee -a /etc/sudoers", sudoersLine))
		if err != nil {
			return fmt.Errorf("failed adding user '%s' to passwordless-sudoers: %w", BrewUserOnMultiUserSystem, err)
		}

		m.logger.Success("Successfully created Homebrew user '%s'", BrewUserOnMultiUserSystem)
	}

	// Set proper ownership of the brew user home directory
	m.logger.Info("Setting ownership of Homebrew directories")
	brewUserHomeDir := "/home/linuxbrew"

	err = m.commander.Run("sudo", "chown", "-R",
		fmt.Sprintf("%s:%s", BrewUserOnMultiUserSystem, BrewUserOnMultiUserSystem),
		brewUserHomeDir)

	if err != nil {
		return fmt.Errorf("failed changing ownership of %s to %s: %w",
			brewUserHomeDir, BrewUserOnMultiUserSystem, err)
	}

	// Install Homebrew as the brew user
	return m.installHomebrew(BrewUserOnMultiUserSystem)
}

// installHomebrew handles both regular and multi-user Homebrew installations
func (b *brewInstaller) installHomebrew(asUser string) error {
	// Download and prepare the installation script
	installScriptPath, cleanup, err := b.downloadAndPrepareInstallScript()
	if err != nil {
		return err
	}
	defer cleanup() // Ensure the temporary script is removed after execution

	if _, err := os.Stat(installScriptPath); err == nil {
		b.logger.Debug("Homebrew install script downloaded to %s", installScriptPath)
	} else {
		return fmt.Errorf("failed to download Homebrew install script: %w", err)
	}

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

	resp, err := http.Get(installScriptURL)
	if err != nil {
		return "", nil, fmt.Errorf("failed to download Homebrew install script: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", nil, fmt.Errorf("failed to download Homebrew install script: HTTP status %d", resp.StatusCode)
	}

	// Create a temporary file for the install script
	b.logger.Debug("Creating temporary file for Homebrew install script")
	tempFile, err := os.CreateTemp("", "brew-install-*.sh")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temporary file for Homebrew install script: %w", err)
	}

	// Create a cleanup function to remove the temp file
	cleanup := func() {
		os.Remove(tempFile.Name())
	}

	// Copy the script to the temp file
	b.logger.Debug("Writing Homebrew install script to temporary file")
	bytesWritten, err := io.Copy(tempFile, resp.Body)
	if err != nil {
		cleanup()
		return "", nil, fmt.Errorf("failed to write Homebrew install script to temporary file: %w", err)
	}
	if bytesWritten == 0 {
		cleanup()
		return "", nil, fmt.Errorf("failed to write Homebrew install script: no bytes written")
	}
	b.logger.Debug("Homebrew install script downloaded successfully")
	b.logger.Debug("First line of script: %s", b.readFirstLine(tempFile.Name()))

	// Close the file to ensure all data is written
	b.logger.Debug("Closing temporary file for Homebrew install script")
	if err = tempFile.Close(); err != nil {
		cleanup()
		return "", nil, fmt.Errorf("failed to close temporary file: %w", err)
	}

	// Make the script executable
	b.logger.Debug("Making Homebrew install script executable")
	if err = os.Chmod(tempFile.Name(), 0777); err != nil {
		cleanup()
		return "", nil, fmt.Errorf("failed to make Homebrew install script executable: %w", err)
	}

	return tempFile.Name(), cleanup, nil
}

// readFirstLine reads the first line of a file
func (b *brewInstaller) readFirstLine(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		return ""
	}
	defer file.Close()

	var firstLine string
	_, err = fmt.Fscanln(file, &firstLine)
	if err != nil {
		return ""
	}

	return firstLine
}

// Options holds configuration options for Homebrew operations
type Options struct {
	MultiUserSystem bool
	Logger          logger.Logger
	SystemInfo      *compatibility.SystemInfo
	Commander       utils.Commander
}

// DefaultOptions returns the default options
func DefaultOptions() Options {
	return Options{
		MultiUserSystem: false,
		Logger:          logger.DefaultLogger,
		SystemInfo:      nil,
		Commander:       utils.NewDefaultCommander(),
	}
}

// WithLogger sets a custom logger for the brew operations
func (o Options) WithLogger(log logger.Logger) Options {
	o.Logger = log
	return o
}

// WithMultiUserSystem configures for multi-user system operation
func (o Options) WithMultiUserSystem(multiUser bool) Options {
	o.MultiUserSystem = multiUser
	return o
}

// WithSystemInfo sets system information for brew operations
func (o Options) WithSystemInfo(sysInfo *compatibility.SystemInfo) Options {
	o.SystemInfo = sysInfo
	return o
}

// WithCommander sets a custom Commander for Homebrew operations
func (o Options) WithCommander(cmdr utils.Commander) Options {
	o.Commander = cmdr
	return o
}
