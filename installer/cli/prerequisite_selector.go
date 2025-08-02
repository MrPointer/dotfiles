package cli

import (
	"errors"
)

// PrerequisiteSelector provides prerequisite-specific selection functionality.
type PrerequisiteSelector struct {
	selector MultiSelectSelector[string]
}

// NewPrerequisiteSelector constructs a PrerequisiteSelector with the given multi-select selector.
func NewPrerequisiteSelector(selector MultiSelectSelector[string]) *PrerequisiteSelector {
	return &PrerequisiteSelector{
		selector: selector,
	}
}

// NewDefaultPrerequisiteSelector constructs a PrerequisiteSelector with the default HuhMultiSelectSelector.
func NewDefaultPrerequisiteSelector() *PrerequisiteSelector {
	return &PrerequisiteSelector{
		selector: NewHuhMultiSelectSelector[string](),
	}
}

// SelectPrerequisites prompts the user to select prerequisites to install from the available missing ones.
// It provides prerequisite-specific context and formatting for better user experience.
func (s *PrerequisiteSelector) SelectPrerequisites(missingPrerequisites []string,
	prerequisiteDetails map[string]PrerequisiteDetail) ([]string, error) {
	if len(missingPrerequisites) == 0 {
		return nil, errors.New("no missing prerequisites available for selection")
	}

	// Create items with descriptions for better user experience
	items := make([]MultiSelectItem[string], len(missingPrerequisites))
	for i, prerequisite := range missingPrerequisites {
		item := MultiSelectItem[string]{
			Value: prerequisite,
			Label: prerequisite,
		}

		if detail, exists := prerequisiteDetails[prerequisite]; exists {
			item.Description = detail.Description
			if detail.InstallHint != "" {
				item.Description += " (" + detail.InstallHint + ")"
			}
		}

		items[i] = item
	}

	title := "Select prerequisites to install:"
	if len(missingPrerequisites) == 1 {
		title = "Install the missing prerequisite?"
	}

	return s.selector.SelectMultiple(title, items)
}

// PrerequisiteDetail represents the details of a prerequisite (matching the compatibility package structure).
type PrerequisiteDetail struct {
	Name        string // Name of the prerequisite.
	Available   bool   // Whether the prerequisite is available.
	Command     string // Command used to check availability.
	Description string // Human-readable description.
	InstallHint string // Hint for installing the prerequisite.
}
