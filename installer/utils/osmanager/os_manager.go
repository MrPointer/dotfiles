package osmanager

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"syscall"

	"github.com/MrPointer/dotfiles/installer/utils"
	"github.com/MrPointer/dotfiles/installer/utils/logger"
)

// UserManager defines operations for managing system users
type UserManager interface {
	// AddUser creates a new user in the system
	AddUser(username string) error

	// AddUserToGroup adds a user to a specified group
	AddUserToGroup(username, group string) error

	// UserExists checks if a user exists in the system
	UserExists(username string) (bool, error)
}

// SudoManager defines operations for managing sudo permissions
type SudoManager interface {
	// AddSudoAccess grants password-less sudo access to a user
	AddSudoAccess(username string) error
}

// FilePermissionManager defines operations for managing filesystem permissions
type FilePermissionManager interface {
	// SetOwnership sets ownership of a directory to a user
	SetOwnership(path, username string) error

	// SetPermissions sets permissions for a file or directory
	SetPermissions(path string, mode os.FileMode) error

	// GetFileOwner returns the username of the file owner
	GetFileOwner(path string) (string, error)
}

// OsManager combines all system operation interfaces
type OsManager interface {
	UserManager
	SudoManager
	FilePermissionManager
}

// UnixOsManager implements OsManager for Unix-like systems
type UnixOsManager struct {
	logger    logger.Logger
	commander utils.Commander
	isRoot    bool
}

var _ OsManager = (*UnixOsManager)(nil)

// NewUnixOsManager creates a new UnixOsManager
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
	u.logger.Info("User '%s' does not exist, creating...", username)

	// Try useradd, fallback to adduser
	useraddCmd := []string{"useradd", "-m", "-s", "/bin/bash", username}
	if !u.isRoot {
		useraddCmd = append([]string{"sudo"}, useraddCmd...)
	}

	err := u.commander.Run(useraddCmd[0], useraddCmd[1:]...)
	if err != nil {
		// Try adduser as fallback
		adduserCmd := []string{"adduser", "--disabled-password", "--gecos", "''", username}
		if !u.isRoot {
			adduserCmd = append([]string{"sudo"}, adduserCmd...)
		}

		err = u.commander.Run(adduserCmd[0], adduserCmd[1:]...)
		if err != nil {
			return fmt.Errorf("failed to create user '%s' with useradd/adduser: %w", username, err)
		}
	}

	return nil
}

func (u *UnixOsManager) AddUserToGroup(username, group string) error {
	u.logger.Info("Adding '%s' to %s group", username, group)
	usermodCmd := []string{"usermod", "-aG", group, username}
	if !u.isRoot {
		usermodCmd = append([]string{"sudo"}, usermodCmd...)
	}

	err := u.commander.Run(usermodCmd[0], usermodCmd[1:]...)
	// Often we don't care if the user is already in the group
	if err != nil {
		u.logger.Debug("Note: User might already be in the %s group", group)
	}

	return nil
}

func (u *UnixOsManager) AddSudoAccess(username string) error {
	sudoersFile := fmt.Sprintf("/etc/sudoers.d/%s", username)
	sudoersLine := fmt.Sprintf("%s ALL=(ALL) NOPASSWD:ALL", username)

	var sudoPrefix string
	if !u.isRoot {
		sudoPrefix = "sudo "
	}

	// Use shell to echo and tee the line into the sudoers file
	shCmd := []string{"sh", "-c", fmt.Sprintf("echo '%s' | %stee %s", sudoersLine, sudoPrefix, sudoersFile)}
	err := u.commander.Run(shCmd[0], shCmd[1:]...)
	if err != nil {
		return fmt.Errorf("failed to add passwordless sudo for '%s': %w", username, err)
	}

	return nil
}

func (u *UnixOsManager) SetOwnership(path, username string) error {
	u.logger.Info("Setting ownership of %s to %s", path, username)
	chownCmd := []string{"chown", "-R", fmt.Sprintf("%s:%s", username, username), path}
	if !u.isRoot {
		chownCmd = append([]string{"sudo"}, chownCmd...)
	}

	err := u.commander.Run(chownCmd[0], chownCmd[1:]...)
	if err != nil {
		return fmt.Errorf("failed to chown %s: %w", path, err)
	}

	return nil
}

func (u *UnixOsManager) SetPermissions(path string, mode os.FileMode) error {
	u.logger.Info("Setting permissions of %s to %o", path, mode)
	chmodCmd := []string{"chmod", fmt.Sprintf("%o", mode), path}
	if !u.isRoot {
		chmodCmd = append([]string{"sudo"}, chmodCmd...)
	}

	err := u.commander.Run(chmodCmd[0], chmodCmd[1:]...)
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
