package cli

import (
	"errors"
)

// ToolSelector provides tool-specific selection functionality.
type ToolSelector struct {
	selector MultiSelectSelector[string]
}

// NewToolSelector constructs a ToolSelector with the given multi-select selector.
func NewToolSelector(selector MultiSelectSelector[string]) *ToolSelector {
	return &ToolSelector{
		selector: selector,
	}
}

// NewDefaultToolSelector constructs a ToolSelector with the default HuhMultiSelectSelector.
func NewDefaultToolSelector() *ToolSelector {
	return &ToolSelector{
		selector: NewHuhMultiSelectSelector[string](),
	}
}

// ToolDetail represents the details of a tool.
type ToolDetail struct {
	Name        string // Name of the tool.
	Description string // Human-readable description.
}

// SelectTools prompts the user to select tools to install from the available ones.
// All tools start unselected. Returns the list of selected tool names (generic package codes).
func (s *ToolSelector) SelectTools(availableTools []string,
	toolDetails map[string]ToolDetail,
) ([]string, error) {
	if len(availableTools) == 0 {
		return nil, errors.New("no tools available for selection")
	}

	items := make([]MultiSelectItem[string], len(availableTools))
	for i, tool := range availableTools {
		item := MultiSelectItem[string]{
			Value: tool,
			Label: tool,
		}

		if detail, exists := toolDetails[tool]; exists {
			item.Description = detail.Description
		}

		items[i] = item
	}

	return s.selector.SelectMultiple("Select optional tools to install:", items)
}
