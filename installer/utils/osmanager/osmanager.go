package osmanager

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strconv"
	"strings"
	"syscall"

	"path/filepath"

	"github.com/MrPointer/dotfiles/installer/utils"
	"github.com/MrPointer/dotfiles/installer/utils/logger"
)

// UserManager defines operations for managing system users.
type UserManager interface {
	// AddUser creates a new user in the system.
	AddUser(username string) error

	// AddUserToGroup adds a user to a specified group.
	AddUserToGroup(username, group string) error

	// UserExists checks if a user exists in the system.
	UserExists(username string) (bool, error)

	// GetHomeDirectory returns the home directory of the current user.
	GetHomeDir() (string, error)

	// GetConfigDir returns the configuration directory of the current user.
	GetConfigDir() (string, error)

	// GetChezmoiConfigHome returns the configuration directory where chezmoi actually looks for its config.
	// This is always ~/.config regardless of XDG specification on different platforms.
	GetChezmoiConfigHome() (string, error)

	// GetCurrentUsername returns the current user's username.
	GetCurrentUsername() (string, error)

	// GetUserShell returns the default login shell for the specified user.
	GetUserShell(username string) (string, error)

	// SetUserShell sets the default login shell for the specified user.
	// On Linux, uses usermod -s. On macOS, uses dscl.
	SetUserShell(username, shellPath string) error
}

// SudoManager defines operations for managing sudo permissions.
type SudoManager interface {
	// AddSudoAccess grants password-less sudo access to a user.
	AddSudoAccess(username string) error
}

// FilePermissionManager defines operations for managing filesystem permissions.
type FilePermissionManager interface {
	// SetOwnership sets ownership of a directory to a user.
	SetOwnership(path, username string) error

	// SetPermissions sets permissions for a file or directory.
	SetPermissions(path string, mode os.FileMode) error

	// GetFileOwner returns the username of the file owner.
	GetFileOwner(path string) (string, error)
}

type VersionExtractor func(string) (string, error)

type ProgramQuery interface {
	// GetProgramPath retrieves the full path of a program if it's available in one of the system's PATH directories.
	// If the program is not found, it returns an error.
	GetProgramPath(program string) (string, error)

	// ProgramExists checks if a program exists in the system's PATH directories.
	// It returns true if the program is found, false if not, and an error if there was an issue checking.
	ProgramExists(program string) (bool, error)

	// GetProgramVersion retrieves the version of a program by executing it with the provided query arguments.
	GetProgramVersion(program string, versionExtractor VersionExtractor, queryArgs ...string) (string, error)
}

// EnvironmentManager defines operations for managing environment variables.
type EnvironmentManager interface {
	// Getenv retrieves the value of the environment variable named by the key.
	Getenv(key string) string
}

// OsManager combines all system operation interfaces.
type OsManager interface {
	UserManager
	SudoManager
	FilePermissionManager
	ProgramQuery
	EnvironmentManager
}

// UnixOsManager implements OsManager for Unix-like systems.
type UnixOsManager struct {
	logger    logger.Logger
	commander utils.Commander
	isRoot    bool
}

var _ OsManager = (*UnixOsManager)(nil)

// NewUnixOsManager creates a new UnixOsManager.
func NewUnixOsManager(logger logger.Logger, commander utils.Commander, isRoot bool) *UnixOsManager {
	return &UnixOsManager{
		logger:    logger,
		commander: commander,
		isRoot:    isRoot,
	}
}

func (u *UnixOsManager) UserExists(username string) (bool, error) {
	_, err := user.Lookup(username)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (u *UnixOsManager) AddUser(username string) error {
	u.logger.Debug("User '%s' does not exist, creating...", username)

	// Try useradd, fallback to adduser.
	useraddCmd := []string{"useradd", "-m", "-s", "/bin/bash", username}
	if !u.isRoot {
		useraddCmd = append([]string{"sudo"}, useraddCmd...)
	}

	_, err := u.commander.RunCommand(useraddCmd[0], useraddCmd[1:])
	if err != nil {
		// Try adduser as fallback.
		adduserCmd := []string{"adduser", "--disabled-password", "--gecos", "''", username}
		if !u.isRoot {
			adduserCmd = append([]string{"sudo"}, adduserCmd...)
		}

		_, err = u.commander.RunCommand(adduserCmd[0], adduserCmd[1:])
		if err != nil {
			return fmt.Errorf("failed to create user '%s' with useradd/adduser: %w", username, err)
		}
	}

	return nil
}

func (u *UnixOsManager) AddUserToGroup(username, group string) error {
	u.logger.Debug("Adding '%s' to %s group", username, group)
	usermodCmd := []string{"usermod", "-aG", group, username}
	if !u.isRoot {
		usermodCmd = append([]string{"sudo"}, usermodCmd...)
	}

	_, err := u.commander.RunCommand(usermodCmd[0], usermodCmd[1:])
	// Often we don't care if the user is already in the group.
	if err != nil {
		u.logger.Debug("Note: User might already be in the %s group", group)
	}

	return nil
}

func (u *UnixOsManager) GetHomeDir() (string, error) {
	return os.UserHomeDir()
}

func (u *UnixOsManager) GetConfigDir() (string, error) {
	return os.UserConfigDir()
}

func (u *UnixOsManager) GetChezmoiConfigHome() (string, error) {
	homeDir, err := u.GetHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, ".config"), nil
}

func (u *UnixOsManager) AddSudoAccess(username string) error {
	sudoersFile := fmt.Sprintf("/etc/sudoers.d/%s", username)
	sudoersLine := fmt.Sprintf("%s ALL=(ALL) NOPASSWD:ALL", username)

	var sudoPrefix string
	if !u.isRoot {
		sudoPrefix = "sudo "
	}

	// Use shell to echo and tee the line into the sudoers file.
	shCmd := []string{"sh", "-c", fmt.Sprintf("echo '%s' | %stee %s", sudoersLine, sudoPrefix, sudoersFile)}
	_, err := u.commander.RunCommand(shCmd[0], shCmd[1:])
	if err != nil {
		return fmt.Errorf("failed to add passwordless sudo for '%s': %w", username, err)
	}

	return nil
}

func (u *UnixOsManager) SetOwnership(path, username string) error {
	u.logger.Debug("Setting ownership of %s to %s", path, username)
	chownCmd := []string{"chown", "-R", fmt.Sprintf("%s:%s", username, username), path}
	if !u.isRoot {
		chownCmd = append([]string{"sudo"}, chownCmd...)
	}

	_, err := u.commander.RunCommand(chownCmd[0], chownCmd[1:])
	if err != nil {
		return fmt.Errorf("failed to chown %s: %w", path, err)
	}

	return nil
}

func (u *UnixOsManager) SetPermissions(path string, mode os.FileMode) error {
	u.logger.Debug("Setting permissions of %s to %o", path, mode)
	chmodCmd := []string{"chmod", fmt.Sprintf("%o", mode), path}
	if !u.isRoot {
		chmodCmd = append([]string{"sudo"}, chmodCmd...)
	}

	_, err := u.commander.RunCommand(chmodCmd[0], chmodCmd[1:])
	if err != nil {
		return fmt.Errorf("failed to chmod %s: %w", path, err)
	}

	return nil
}

func (u *UnixOsManager) GetFileOwner(path string) (string, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("failed to get file info for %s: %w", path, err)
	}

	stat, ok := fileInfo.Sys().(*syscall.Stat_t)
	if !ok {
		return "", fmt.Errorf("failed to get file info")
	}

	owner, err := user.LookupId(strconv.FormatUint(uint64(stat.Uid), 10))
	if err != nil {
		return "", fmt.Errorf("failed to lookup owner for %s: %w", path, err)
	}

	return owner.Username, nil
}

func (u *UnixOsManager) GetProgramPath(program string) (string, error) {
	return exec.LookPath(program)
}

func (u *UnixOsManager) ProgramExists(program string) (bool, error) {
	_, err := u.GetProgramPath(program)
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) || errors.Is(err, os.ErrNotExist) {
			return false, nil // Program not found.
		}
		return false, fmt.Errorf("error checking program existence: %w", err)
	}
	return true, nil // Program found.
}

func (u *UnixOsManager) GetProgramVersion(
	program string,
	versionExtractor VersionExtractor,
	queryArgs ...string,
) (string, error) {
	args := []string{"--version"} // Default argument for version query.
	if len(queryArgs) > 0 {
		args = queryArgs
	}

	cmd := exec.Command(program, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get version for %s: %w", program, err)
	}

	version, err := versionExtractor(string(output))
	if err != nil {
		return "", fmt.Errorf("failed to extract version from output: %w", err)
	}

	return version, nil
}

func (u *UnixOsManager) Getenv(key string) string {
	return os.Getenv(key)
}

// GetCurrentUsername returns the current user's username.
func (u *UnixOsManager) GetCurrentUsername() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("failed to get current user: %w", err)
	}
	return currentUser.Username, nil
}

// GetUserShell returns the default login shell for the specified user.
// On macOS, uses dscl. On Linux, reads /etc/passwd.
func (u *UnixOsManager) GetUserShell(username string) (string, error) {
	// Try to get shell from user.Lookup first (works on both platforms).
	lookedUpUser, err := user.Lookup(username)
	if err != nil {
		return "", fmt.Errorf("failed to lookup user %s: %w", username, err)
	}

	// On Unix systems, we can try to get the shell from /etc/passwd via the user package.
	// However, the user.User struct doesn't expose the shell directly.
	// We need platform-specific approaches.

	// For macOS, use dscl.
	if isDarwin() {
		result, err := u.commander.RunCommand("dscl", []string{
			".", "-read",
			fmt.Sprintf("/Users/%s", username),
			"UserShell",
		}, utils.WithCaptureOutput())
		if err != nil {
			return "", fmt.Errorf("failed to read UserShell via dscl: %w", err)
		}

		// Output format: "UserShell: /path/to/shell"
		output := strings.TrimSpace(string(result.Stdout))
		parts := strings.SplitN(output, ":", 2)
		if len(parts) != 2 {
			return "", fmt.Errorf("unexpected dscl output format: %s", output)
		}

		return strings.TrimSpace(parts[1]), nil
	}

	// For Linux, parse /etc/passwd.
	// The user's home directory from Lookup gives us a hint that user exists.
	_ = lookedUpUser

	passwdFile, err := os.Open("/etc/passwd")
	if err != nil {
		return "", fmt.Errorf("failed to open /etc/passwd: %w", err)
	}
	defer passwdFile.Close()

	scanner := bufio.NewScanner(passwdFile)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, username+":") {
			fields := strings.Split(line, ":")
			if len(fields) >= 7 {
				return fields[6], nil
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading /etc/passwd: %w", err)
	}

	return "", fmt.Errorf("user %s not found in /etc/passwd", username)
}

// SetUserShell sets the default login shell for the specified user.
// On macOS, uses dscl. On Linux, uses usermod -s.
func (u *UnixOsManager) SetUserShell(username, shellPath string) error {
	if isDarwin() {
		// dscl . -create /Users/username UserShell /path/to/shell
		_, err := u.commander.RunCommand("dscl", []string{
			".", "-create",
			fmt.Sprintf("/Users/%s", username),
			"UserShell",
			shellPath,
		}, utils.WithCaptureOutput())
		if err != nil {
			return fmt.Errorf("failed to set shell via dscl: %w", err)
		}
		return nil
	}

	// Linux: usermod -s /path/to/shell username
	usermodCmd := []string{"usermod", "-s", shellPath, username}
	if !u.isRoot {
		usermodCmd = append([]string{"sudo"}, usermodCmd...)
	}

	_, err := u.commander.RunCommand(usermodCmd[0], usermodCmd[1:])
	if err != nil {
		return fmt.Errorf("failed to set shell via usermod: %w", err)
	}

	return nil
}

// isDarwin returns true if the current OS is macOS.
func isDarwin() bool {
	return runtime.GOOS == "darwin"
}

// IsRoot returns true if the current user is root.
func IsRoot() bool {
	return os.Geteuid() == 0
}
