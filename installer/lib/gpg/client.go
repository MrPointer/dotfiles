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

	// Run the command interactively while capturing output for parsing
	args := []string{"--gen-key", "--default-new-key-algo", "nistp256"}
	result, err := c.commander.RunCommand("gpg", args, utils.WithInteractiveCapture(), utils.WithEnv(map[string]string{"GPG_TTY": activeTerminal}))
	if err != nil {
		return "", err
	} else if result.ExitCode != 0 {
		return "", errors.New("failed to create GPG key pair: " + result.StderrString())
	}

	// Parse the output to extract the key ID using multiple robust methods
	keyID, err := c.extractKeyIDFromGPGOutput(result.String())
	if err != nil {
		return "", errors.New("failed to extract GPG key ID from output: " + err.Error())
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
	c.logger.Trace("Detecting TTY for GPG operations")

	// Method 1: Check if GPG_TTY is already set in the environment
	if tty := c.osMgr.Getenv("GPG_TTY"); tty != "" {
		c.logger.Trace("Using GPG_TTY from environment: " + tty)
		return strings.TrimSpace(tty), nil
	}

	// Method 2: Try the tty command
	ttyOutput, err := c.commander.RunCommand("tty", []string{}, utils.WithCaptureOutput())
	if err == nil && ttyOutput.ExitCode == 0 && len(ttyOutput.Stdout) > 0 {
		tty := strings.TrimSpace(ttyOutput.String())
		if tty != "" && tty != "not a tty" {
			c.logger.Trace("Detected TTY using tty command: " + tty)
			return tty, nil
		}
	}

	// Method 3: Try /dev/tty directly
	if exists, err := c.fs.PathExists("/dev/tty"); err == nil && exists {
		c.logger.Trace("Using /dev/tty as fallback")
		return "/dev/tty", nil
	}

	// Method 4: Check common TTY environment variables
	for _, envVar := range []string{"TTY", "TERM_TTY"} {
		if tty := c.osMgr.Getenv(envVar); tty != "" {
			c.logger.Trace("Using TTY from " + envVar + ": " + tty)
			return strings.TrimSpace(tty), nil
		}
	}

	return "", errors.New("unable to detect TTY for GPG operations - ensure you're running in an interactive terminal or set GPG_TTY environment variable")
}

// extractKeyIDFromGPGOutput parses GPG output to extract the key ID using multiple robust methods.
func (c *DefaultGpgClient) extractKeyIDFromGPGOutput(output string) (string, error) {
	c.logger.Trace("Extracting key ID from GPG output")

	// Method 1: Look for "key <KEYID> marked as ultimately trusted" pattern
	// This is the most reliable pattern across distributions
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "marked as ultimately trusted") {
			// Pattern: "gpg: key ABC123DEF456 marked as ultimately trusted"
			parts := strings.Fields(line)
			for i, part := range parts {
				if part == "key" && i+1 < len(parts) {
					keyID := parts[i+1]
					c.logger.Trace("Found key ID using 'ultimately trusted' pattern: " + keyID)
					return keyID, nil
				}
			}
		}
	}

	// Method 2: Look for "gpg: <KEYID>: public key" pattern
	for _, line := range lines {
		if strings.Contains(line, ": public key") {
			// Pattern: "gpg: ABC123DEF456: public key"
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				keyID := strings.TrimSpace(parts[1])
				if keyID != "" && keyID != "gpg" {
					c.logger.Trace("Found key ID using 'public key' pattern: " + keyID)
					return keyID, nil
				}
			}
		}
	}

	// Method 3: Look for fingerprint patterns in pub lines
	// Parse lines that start with "pub" and extract the key ID
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "pub") {
			// The key ID might be on the same line or the next line
			// Format: "pub   nistp256 2024-01-01 [SC]" followed by "      ABC123DEF456"
			if i+1 < len(lines) {
				nextLine := strings.TrimSpace(lines[i+1])
				// Check if next line looks like a key ID (alphanumeric, reasonable length)
				if len(nextLine) >= 8 && len(nextLine) <= 40 && isAlphanumeric(nextLine) {
					c.logger.Trace("Found key ID using pub line pattern: " + nextLine)
					return nextLine, nil
				}
			}

			// Also check if key ID is on the same line after the algorithm
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				// Sometimes format is "pub   nistp256/ABC123DEF456 2024-01-01 [SC]"
				for _, part := range parts {
					if strings.Contains(part, "/") {
						keyParts := strings.Split(part, "/")
						if len(keyParts) == 2 {
							keyID := keyParts[1]
							if len(keyID) >= 8 && isAlphanumeric(keyID) {
								c.logger.Trace("Found key ID using pub/keyid pattern: " + keyID)
								return keyID, nil
							}
						}
					}
				}
			}
		}
	}

	// Method 4: Look for revocation certificate patterns
	for _, line := range lines {
		if strings.Contains(line, "revocation certificate stored") {
			// Pattern: "gpg: revocation certificate stored as '/path/ABC123DEF456.rev'"
			// Extract the filename and get the key ID from it
			if strings.Contains(line, ".rev") {
				parts := strings.Split(line, "/")
				if len(parts) > 0 {
					filename := parts[len(parts)-1]
					if strings.HasSuffix(filename, ".rev'") || strings.HasSuffix(filename, ".rev\"") {
						keyID := strings.TrimSuffix(strings.TrimSuffix(filename, ".rev'"), ".rev\"")
						if len(keyID) >= 8 && isAlphanumeric(keyID) {
							c.logger.Trace("Found key ID using revocation certificate pattern: " + keyID)
							return keyID, nil
						}
					}
				}
			}
		}
	}

	return "", errors.New("could not find key ID in GPG output using any known pattern")
}

// isAlphanumeric checks if a string contains only alphanumeric characters.
func isAlphanumeric(s string) bool {
	for _, r := range s {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')) {
			return false
		}
	}
	return true
}
