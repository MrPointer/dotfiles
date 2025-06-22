package gpg

import (
	"errors"
	"strings"

	"github.com/MrPointer/dotfiles/installer/utils"
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
	commander utils.Commander
}

// NewDefaultGpgClient constructs a DefaultGpgClient with the given OsManager and Commander.
func NewDefaultGpgClient(
	osMgr osmanager.OsManager,
	commander utils.Commander,
) *DefaultGpgClient {
	return &DefaultGpgClient{
		osMgr:     osMgr,
		commander: commander,
	}
}

// CreateKeyPair implements GpgClient.
func (c *DefaultGpgClient) CreateKeyPair() (string, error) {
	args := []string{"--expert", "--full-generate-key"}
	// Run the command interactively.
	result, err := c.commander.RunCommand("gpg", args, utils.WithCaptureOutput())
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
	availableKeys, err := c.ListAvailableKeys()
	if err != nil {
		return false, err
	}

	return availableKeys != nil && len(availableKeys) > 0, nil
}
