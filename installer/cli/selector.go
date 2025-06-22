package cli

import (
	"errors"

	"github.com/charmbracelet/huh"
)

// Selector defines the interface for selecting items from a list.
type Selector[T comparable] interface {
	// Select prompts the user to select an item from the provided list.
	Select(title string, items []T) (T, error)
	// SelectWithLabels prompts the user to select an item with custom labels.
	SelectWithLabels(title string, items []T, labels []string) (T, error)
}

var _ Selector[string] = (*HuhSelector[string])(nil)

// HuhSelector implements Selector using the huh library.
type HuhSelector[T comparable] struct{}

// NewHuhSelector constructs a HuhSelector.
func NewHuhSelector[T comparable]() *HuhSelector[T] {
	return &HuhSelector[T]{}
}

// Select implements Selector.
func (s *HuhSelector[T]) Select(title string, items []T) (T, error) {
	var zero T
	if len(items) == 0 {
		return zero, errors.New("no items available for selection")
	}

	if len(items) == 1 {
		return items[0], nil
	}

	var selectedItem T
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[T]().
				Title(title).
				Options(huh.NewOptions(items...)...).
				Value(&selectedItem),
		),
	)

	err := form.Run()
	if err != nil {
		return zero, err
	}

	return selectedItem, nil
}

// SelectWithLabels implements Selector.
func (s *HuhSelector[T]) SelectWithLabels(title string, items []T, labels []string) (T, error) {
	var zero T
	if len(items) == 0 {
		return zero, errors.New("no items available for selection")
	}

	if len(items) != len(labels) {
		return zero, errors.New("items and labels must have the same length")
	}

	if len(items) == 1 {
		return items[0], nil
	}

	options := make([]huh.Option[T], len(items))
	for i, item := range items {
		options[i] = huh.NewOption(labels[i], item)
	}

	var selectedItem T
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[T]().
				Title(title).
				Options(options...).
				Value(&selectedItem),
		),
	)

	err := form.Run()
	if err != nil {
		return zero, err
	}

	return selectedItem, nil
}
