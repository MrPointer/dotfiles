package cli

import (
	"fmt"
)

// GpgKeySelector provides GPG-specific key selection functionality.
type GpgKeySelector struct {
	selector Selector[string]
}

// NewGpgKeySelector constructs a GpgKeySelector with the given generic selector.
func NewGpgKeySelector(selector Selector[string]) *GpgKeySelector {
	return &GpgKeySelector{
		selector: selector,
	}
}

// NewDefaultGpgKeySelector constructs a GpgKeySelector with the default HuhSelector.
func NewDefaultGpgKeySelector() *GpgKeySelector {
	return &GpgKeySelector{
		selector: NewHuhSelector[string](),
	}
}

// SelectKey prompts the user to select a GPG key from the available keys.
// It provides GPG-specific context and formatting for better user experience.
func (s *GpgKeySelector) SelectKey(availableKeys []string) (string, error) {
	if len(availableKeys) == 0 {
		return "", fmt.Errorf("no GPG keys available for selection")
	}

	if len(availableKeys) == 1 {
		return availableKeys[0], nil
	}

	title := fmt.Sprintf("Select a GPG key to use (%d available):", len(availableKeys))

	// Format keys with additional context for better display
	labels := make([]string, len(availableKeys))
	for i, key := range availableKeys {
		labels[i] = fmt.Sprintf("Key ID: %s", key)
	}

	return s.selector.SelectWithLabels(title, availableKeys, labels)
}
