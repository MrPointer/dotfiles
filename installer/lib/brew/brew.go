package brew

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strconv"
	"syscall"

	"github.com/MrPointer/dotfiles/installer/lib/compatibility"
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

// Options holds configuration options for Homebrew operations
type Options struct {
	MultiUserSystem bool
	Logger          logger.Logger
	SystemInfo      *compatibility.SystemInfo
}

// DefaultOptions returns the default options
func DefaultOptions() Options {
	return Options{
		MultiUserSystem: false,
		Logger:          logger.DefaultLogger,
		SystemInfo:      nil,
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

// DetectBrewPath returns the appropriate brew binary path based on the system information
func DetectBrewPath(opts Options) string {
	// If system info is provided, use it to determine the brew path
	if opts.SystemInfo != nil {
		switch opts.SystemInfo.OSName {
		case "darwin":
			if opts.SystemInfo.Arch == "arm64" {
				return MacOSARMBrewPath
			}
			return MacOSIntelBrewPath
		case "linux":
			return LinuxBrewPath
		default:
			// Fallback to runtime detection
		}
	}

	// Fallback to runtime-based detection if no system info or unsupported OS
	switch runtime.GOOS {
	case "darwin":
		if runtime.GOARCH == "arm64" {
			return MacOSARMBrewPath
		}
		return MacOSIntelBrewPath
	default:
		return LinuxBrewPath
	}
}

// IsAvailable checks if Homebrew is already installed and available
func IsAvailable(opts Options) bool {
	brewPath := DetectBrewPath(opts)

	if opts.MultiUserSystem {
		// For multi-user systems, check if the brew binary exists and is owned by the correct user
		fileInfo, err := os.Stat(brewPath)
		if err != nil {
			return false
		}

		// On Linux, check the ownership of the brew binary
		if opts.SystemInfo != nil && opts.SystemInfo.OSName == "linux" ||
			opts.SystemInfo == nil && runtime.GOOS == "linux" {
			stat, ok := fileInfo.Sys().(*syscall.Stat_t)
			if !ok {
				return false
			}

			brewUser, err := user.LookupId(strconv.FormatUint(uint64(stat.Uid), 10))
			if err != nil {
				return false
			}

			return brewUser.Username == BrewUserOnMultiUserSystem
		}

		// We don't officially support multi-user on macOS, but at least check if brew exists
		return true
	}

	// For single-user systems, just check if the file exists
	_, err := os.Stat(brewPath)
	return err == nil
}

// RunBrewCommand runs a brew command with appropriate user permissions
func RunBrewCommand(opts Options, args ...string) error {
	brewPath := DetectBrewPath(opts)

	if opts.MultiUserSystem && (opts.SystemInfo == nil || opts.SystemInfo.OSName == "linux") {
		// On multi-user Linux systems, run brew as the brew user
		sudoArgs := []string{"-Hu", BrewUserOnMultiUserSystem, brewPath}
		sudoArgs = append(sudoArgs, args...)

		cmd := exec.Command("sudo", sudoArgs...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		return cmd.Run()
	}

	// Regular brew command for single-user systems
	cmd := exec.Command(brewPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// Install installs Homebrew if not already installed
func Install(opts Options) error {
	if IsAvailable(opts) {
		opts.Logger.Success("Homebrew is already installed")
		return nil
	}

	opts.Logger.Info("Installing Homebrew")

	if opts.MultiUserSystem {
		if runtime.GOOS == "darwin" {
			return fmt.Errorf("multi-user Homebrew installation is not supported on macOS, please install manually")
		}

		return installMultiUserLinux(opts)
	}

	// Regular installation for single-user systems
	return installHomebrew(opts, false, "")
}

// installMultiUserLinux installs Homebrew in a multi-user configuration on Linux
func installMultiUserLinux(opts Options) error {
	// Check if the brew user already exists
	_, err := user.Lookup(BrewUserOnMultiUserSystem)
	if err != nil {
		// Create the brew user
		opts.Logger.Info("Creating user '%s' for Homebrew", BrewUserOnMultiUserSystem)

		cmd := exec.Command("sudo", "useradd", "-m", "-p", "", BrewUserOnMultiUserSystem)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed creating user '%s' for Homebrew: %w", BrewUserOnMultiUserSystem, err)
		}

		// Add user to sudo group
		cmd = exec.Command("sudo", "usermod", "-aG", "sudo", BrewUserOnMultiUserSystem)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed adding user '%s' to sudo group: %w", BrewUserOnMultiUserSystem, err)
		}

		// Add user to passwordless sudo
		sudoersLine := fmt.Sprintf("%s ALL=(ALL) NOPASSWD:ALL", BrewUserOnMultiUserSystem)
		cmd = exec.Command("sudo", "bash", "-c", fmt.Sprintf("echo '%s' | sudo tee -a /etc/sudoers", sudoersLine))
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed adding user '%s' to passwordless-sudoers: %w", BrewUserOnMultiUserSystem, err)
		}

		opts.Logger.Success("Successfully created Homebrew user '%s'", BrewUserOnMultiUserSystem)
	}

	// Set proper ownership of the brew user home directory
	opts.Logger.Info("Setting ownership of Homebrew directories")
	brewUserHomeDir := "/home/linuxbrew"
	cmd := exec.Command("sudo", "chown", "-R",
		fmt.Sprintf("%s:%s", BrewUserOnMultiUserSystem, BrewUserOnMultiUserSystem),
		brewUserHomeDir)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed changing ownership of %s to %s: %w",
			brewUserHomeDir, BrewUserOnMultiUserSystem, err)
	}

	// Install Homebrew as the brew user
	return installHomebrew(opts, true, BrewUserOnMultiUserSystem)
}

// installHomebrew handles both regular and multi-user Homebrew installations
func installHomebrew(opts Options, asUser bool, username string) error {
	// Download and prepare the installation script
	installScriptPath, cleanup, err := downloadAndPrepareInstallScript(opts.Logger)
	if err != nil {
		return err
	}
	defer cleanup() // Ensure the temporary script is removed after execution

	// Execute the downloaded install script, optionally as a different user
	if asUser {
		if username == "" {
			return fmt.Errorf("username must be provided when installing as a different user")
		}

		opts.Logger.Info("Running Homebrew install script as %s", username)
		cmd := exec.Command("sudo", "-Hu", username, "bash", installScriptPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed running Homebrew install script as %s: %w", username, err)
		}

		opts.Logger.Success("Successfully installed Homebrew for user %s", username)
	} else {
		opts.Logger.Info("Running Homebrew install script")
		cmd := exec.Command("bash", installScriptPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed installing Homebrew: %w", err)
		}

		opts.Logger.Success("Homebrew installed successfully")
	}

	return nil
}

// downloadAndPrepareInstallScript downloads the Homebrew installation script and prepares it for execution
func downloadAndPrepareInstallScript(logger logger.Logger) (string, func(), error) {
	logger.Info("Downloading Homebrew install script")
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
	tempFile, err := os.CreateTemp("", "brew-install-*.sh")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temporary file for Homebrew install script: %w", err)
	}

	// Create a cleanup function to remove the temp file
	cleanup := func() {
		os.Remove(tempFile.Name())
	}

	// Copy the script to the temp file
	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		cleanup()
		return "", nil, fmt.Errorf("failed to write Homebrew install script to temporary file: %w", err)
	}

	// Close the file to ensure all data is written
	if err = tempFile.Close(); err != nil {
		cleanup()
		return "", nil, fmt.Errorf("failed to close temporary file: %w", err)
	}

	// Make the script executable
	if err = os.Chmod(tempFile.Name(), 0755); err != nil {
		cleanup()
		return "", nil, fmt.Errorf("failed to make Homebrew install script executable: %w", err)
	}

	return tempFile.Name(), cleanup, nil
}
