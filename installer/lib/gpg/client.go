package gpg

import (
	"errors"
	"strings"

	"github.com/MrPointer/dotfiles/installer/utils"
	"github.com/MrPointer/dotfiles/installer/utils/logger"
	"github.com/MrPointer/dotfiles/installer/utils/osmanager"
)

// GpgClient defines the interface for interacting with GPG.
type GpgClient interface {
	// CreateKeyPair creates a GPG key pair interactively.
	CreateKeyPair() (string, error)
	// KeysAvailable returns true if there are secret keys, false otherwise.
	KeysAvailable() (bool, error)
	// ListAvailableKeys lists all available GPG keys.
	ListAvailableKeys() ([]string, error)
}

var _ GpgClient = (*DefaultGpgClient)(nil)

// DefaultGpgClient implements GpgClient using OsManager and Commander.
type DefaultGpgClient struct {
	osMgr     osmanager.OsManager
	fs        utils.FileSystem
	commander utils.Commander
	logger    logger.Logger
}

// NewDefaultGpgClient constructs a DefaultGpgClient with the given OsManager, Commander, and Logger.
func NewDefaultGpgClient(
	osMgr osmanager.OsManager,
	fs utils.FileSystem,
	commander utils.Commander,
	logger logger.Logger,
) *DefaultGpgClient {
	return &DefaultGpgClient{
		osMgr:     osMgr,
		fs:        fs,
		commander: commander,
		logger:    logger,
	}
}

// CreateKeyPair implements GpgClient.
func (c *DefaultGpgClient) CreateKeyPair() (string, error) {
	c.logger.Debug("Creating GPG key pair")

	activeTerminal, err := c.detectTTY()
	if err != nil {
		return "", err
	}

	// Run the command interactively using the working variant with ECC NIST P-256
	args := []string{"--gen-key", "--pinentry-mode", "loopback", "--default-new-key-algo", "nistp256"}
	result, err := c.commander.RunCommand("gpg", args, utils.WithCaptureOutput(), utils.WithEnv(map[string]string{"GPG_TTY": activeTerminal}))
	if err != nil {
		return "", err
	} else if result.ExitCode != 0 {
		return "", errors.New("failed to create GPG key pair: " + result.StderrString())
	}

	trimmedOutput := strings.TrimSpace(result.String())
	outputLines := strings.Split(trimmedOutput, "\n")

	// Take last 3 lines of output containing the summary of the key creation.
	if len(outputLines) < 3 {
		return "", errors.New("unexpected output format from GPG key creation")
	}
	keySummary := outputLines[len(outputLines)-3:]

	// From that summary we extract the full key ID, which is expected to be the 2nd line after the "pub" field.
	var keyID string
	parseCurrentLine := false
	for _, line := range keySummary {
		if parseCurrentLine {
			keyID = strings.TrimSpace(line)
			break
		} else if strings.HasPrefix(line, "pub") {
			parseCurrentLine = true
		}
	}

	if keyID == "" {
		return "", errors.New("failed to extract GPG key ID from output")
	}
	return keyID, nil
}

// ListAvailableKeys implements GpgClient.
func (c *DefaultGpgClient) ListAvailableKeys() ([]string, error) {
	c.logger.Debug("Listing available GPG keys")

	args := []string{"--list-secret-keys", "--keyid-format", "LONG"}
	result, err := c.commander.RunCommand("gpg", args, utils.WithCaptureOutput())
	if err != nil {
		return nil, err
	}

	// Search for "sec" in output.
	trimmedOutput := strings.TrimSpace(result.String())
	outputLines := strings.Split(trimmedOutput, "\n")

	keys := make([]string, 0, len(outputLines))
	for _, line := range outputLines {
		if strings.HasPrefix(line, "sec") {
			// Extract the key ID from the line.
			parts := strings.Fields(line)
			if len(parts) > 1 {
				fullKeyID := parts[1]
				// The key ID is everything after the last slash.
				slashIndex := strings.LastIndex(fullKeyID, "/")
				if slashIndex != -1 {
					keyID := fullKeyID[slashIndex+1:]
					keys = append(keys, keyID)
				}
			}
		}
	}

	if len(keys) == 0 {
		return nil, nil // No keys found
	}

	return keys, nil
}

// KeysAvailable returns true if there are secret keys, false otherwise.
func (c *DefaultGpgClient) KeysAvailable() (bool, error) {
	c.logger.Debug("Checking for available GPG keys")

	availableKeys, err := c.ListAvailableKeys()
	if err != nil {
		return false, err
	}

	return len(availableKeys) > 0, nil
}

// detectTTY attempts to detect the current TTY using multiple fallback methods.
// This is essential for GPG operations that require interactive input, as GPG needs
// to know which terminal to use for user prompts.
//
// The method tries the following approaches in order:
// 1. Check if GPG_TTY environment variable is already set
// 2. Try running the 'tty' command to get the current terminal
// 3. Use /dev/tty as a fallback if it exists
// 4. Check other common TTY environment variables (TTY, TERM_TTY)
//
// Returns the detected TTY path or an error if no TTY can be determined.
func (c *DefaultGpgClient) detectTTY() (string, error) {
	c.logger.Debug("Detecting TTY for GPG operations")

	// Method 1: Check if GPG_TTY is already set in the environment
	if tty := c.osMgr.Getenv("GPG_TTY"); tty != "" {
		c.logger.Debug("Using GPG_TTY from environment: " + tty)
		return strings.TrimSpace(tty), nil
	}

	// Method 2: Try the tty command
	ttyOutput, err := c.commander.RunCommand("tty", []string{}, utils.WithCaptureOutput())
	if err == nil && ttyOutput.ExitCode == 0 && len(ttyOutput.Stdout) > 0 {
		tty := strings.TrimSpace(ttyOutput.String())
		if tty != "" && tty != "not a tty" {
			c.logger.Debug("Detected TTY using tty command: " + tty)
			return tty, nil
		}
	}

	// Method 3: Try /dev/tty directly
	if exists, err := c.fs.PathExists("/dev/tty"); err == nil && exists {
		c.logger.Debug("Using /dev/tty as fallback")
		return "/dev/tty", nil
	}

	// Method 4: Check common TTY environment variables
	for _, envVar := range []string{"TTY", "TERM_TTY"} {
		if tty := c.osMgr.Getenv(envVar); tty != "" {
			c.logger.Debug("Using TTY from " + envVar + ": " + tty)
			return strings.TrimSpace(tty), nil
		}
	}

	return "", errors.New("unable to detect TTY for GPG operations - ensure you're running in an interactive terminal or set GPG_TTY environment variable")
}
