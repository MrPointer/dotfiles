package cli

import (
	"errors"

	"github.com/charmbracelet/huh"
)

// MultiSelectItem represents an item that can be selected in a multi-select interface.
type MultiSelectItem[T comparable] struct {
	Value       T      // The actual value to be returned.
	Label       string // The display label for the item.
	Description string // Optional description for additional context.
}

// MultiSelectSelector defines the interface for selecting multiple items from a list.
type MultiSelectSelector[T comparable] interface {
	// SelectMultiple prompts the user to select multiple items from the provided list.
	SelectMultiple(title string, items []MultiSelectItem[T]) ([]T, error)
}

var _ MultiSelectSelector[string] = (*HuhMultiSelectSelector[string])(nil)

// HuhMultiSelectSelector implements MultiSelectSelector using the huh library.
type HuhMultiSelectSelector[T comparable] struct{}

// NewHuhMultiSelectSelector constructs a HuhMultiSelectSelector.
func NewHuhMultiSelectSelector[T comparable]() *HuhMultiSelectSelector[T] {
	return &HuhMultiSelectSelector[T]{}
}

// SelectMultiple implements MultiSelectSelector.
func (s *HuhMultiSelectSelector[T]) SelectMultiple(title string, items []MultiSelectItem[T]) ([]T, error) {
	if len(items) == 0 {
		return nil, errors.New("no items available for selection")
	}

	// Create options for the multi-select widget
	options := make([]huh.Option[T], len(items))
	for i, item := range items {
		// Combine label and description for display
		displayText := item.Label
		if item.Description != "" {
			displayText = item.Label + " - " + item.Description
		}
		options[i] = huh.NewOption(displayText, item.Value)
	}

	var selectedItems []T
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[T]().
				Title(title).
				Description("Use space to select/deselect, enter to confirm").
				Options(options...).
				Value(&selectedItems),
		),
	)

	err := form.Run()
	if err != nil {
		return nil, err
	}

	return selectedItems, nil
}
