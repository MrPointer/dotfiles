package cli

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

// MockMultiSelectSelector is an inline mock for MultiSelectSelector[string].
type MockMultiSelectSelector struct {
	SelectMultipleFunc func(title string, items []MultiSelectItem[string]) ([]string, error)
}

var _ MultiSelectSelector[string] = (*MockMultiSelectSelector)(nil)

func (m *MockMultiSelectSelector) SelectMultiple(title string, items []MultiSelectItem[string]) ([]string, error) {
	if m.SelectMultipleFunc != nil {
		return m.SelectMultipleFunc(title, items)
	}
	return nil, nil
}

func Test_ItemsCreatedCorrectlyFromToolDetails(t *testing.T) {
	availableTools := []string{"tool1", "tool2", "tool3"}
	toolDetails := map[string]ToolDetail{
		"tool1": {Name: "tool1", Description: "First tool"},
		"tool2": {Name: "tool2", Description: "Second tool"},
		"tool3": {Name: "tool3", Description: "Third tool"},
	}

	var capturedItems []MultiSelectItem[string]
	mockSelector := &MockMultiSelectSelector{
		SelectMultipleFunc: func(title string, items []MultiSelectItem[string]) ([]string, error) {
			capturedItems = items
			return []string{}, nil
		},
	}

	selector := NewToolSelector(mockSelector)
	_, err := selector.SelectTools(availableTools, toolDetails)

	require.NoError(t, err)
	require.Len(t, capturedItems, 3)

	require.Equal(t, "tool1", capturedItems[0].Value)
	require.Equal(t, "tool1", capturedItems[0].Label)
	require.Equal(t, "First tool", capturedItems[0].Description)

	require.Equal(t, "tool2", capturedItems[1].Value)
	require.Equal(t, "tool2", capturedItems[1].Label)
	require.Equal(t, "Second tool", capturedItems[1].Description)

	require.Equal(t, "tool3", capturedItems[2].Value)
	require.Equal(t, "tool3", capturedItems[2].Label)
	require.Equal(t, "Third tool", capturedItems[2].Description)
}

func Test_EmptyAvailableToolsReturnsError(t *testing.T) {
	mockSelector := &MockMultiSelectSelector{
		SelectMultipleFunc: func(title string, items []MultiSelectItem[string]) ([]string, error) {
			return nil, nil
		},
	}

	selector := NewToolSelector(mockSelector)
	_, err := selector.SelectTools([]string{}, map[string]ToolDetail{})

	require.Error(t, err)
	require.Equal(t, "no tools available for selection", err.Error())
}

func Test_EmptySelectionReturnsEmptySlice(t *testing.T) {
	availableTools := []string{"tool1", "tool2"}
	toolDetails := map[string]ToolDetail{
		"tool1": {Name: "tool1", Description: "First tool"},
		"tool2": {Name: "tool2", Description: "Second tool"},
	}

	mockSelector := &MockMultiSelectSelector{
		SelectMultipleFunc: func(title string, items []MultiSelectItem[string]) ([]string, error) {
			return []string{}, nil
		},
	}

	selector := NewToolSelector(mockSelector)
	result, err := selector.SelectTools(availableTools, toolDetails)

	require.NoError(t, err)
	require.Empty(t, result)
	require.Equal(t, []string{}, result)
}

func Test_SelectorErrorReturnsError(t *testing.T) {
	availableTools := []string{"tool1"}
	toolDetails := map[string]ToolDetail{
		"tool1": {Name: "tool1", Description: "First tool"},
	}

	expectedErr := errors.New("user cancelled selection")
	mockSelector := &MockMultiSelectSelector{
		SelectMultipleFunc: func(title string, items []MultiSelectItem[string]) ([]string, error) {
			return nil, expectedErr
		},
	}

	selector := NewToolSelector(mockSelector)
	_, err := selector.SelectTools(availableTools, toolDetails)

	require.Error(t, err)
	require.Equal(t, expectedErr, err)
}

func Test_SelectionReturnsCorrectToolNames(t *testing.T) {
	availableTools := []string{"tool1", "tool2", "tool3"}
	toolDetails := map[string]ToolDetail{
		"tool1": {Name: "tool1", Description: "First tool"},
		"tool2": {Name: "tool2", Description: "Second tool"},
		"tool3": {Name: "tool3", Description: "Third tool"},
	}

	expectedSelection := []string{"tool1", "tool3"}
	mockSelector := &MockMultiSelectSelector{
		SelectMultipleFunc: func(title string, items []MultiSelectItem[string]) ([]string, error) {
			return expectedSelection, nil
		},
	}

	selector := NewToolSelector(mockSelector)
	result, err := selector.SelectTools(availableTools, toolDetails)

	require.NoError(t, err)
	require.Equal(t, expectedSelection, result)
}

func Test_ItemsWithoutDetailsHaveEmptyDescription(t *testing.T) {
	availableTools := []string{"tool1", "tool2"}
	toolDetails := map[string]ToolDetail{
		"tool1": {Name: "tool1", Description: "First tool"},
	}

	var capturedItems []MultiSelectItem[string]
	mockSelector := &MockMultiSelectSelector{
		SelectMultipleFunc: func(title string, items []MultiSelectItem[string]) ([]string, error) {
			capturedItems = items
			return []string{}, nil
		},
	}

	selector := NewToolSelector(mockSelector)
	_, err := selector.SelectTools(availableTools, toolDetails)

	require.NoError(t, err)
	require.Len(t, capturedItems, 2)

	// tool1 has description
	require.Equal(t, "First tool", capturedItems[0].Description)

	// tool2 has no description
	require.Equal(t, "", capturedItems[1].Description)
}
