package collections

import "testing"

func TestLast(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected int
	}{
		{
			name:     "non-empty slice",
			input:    []int{1, 2, 3},
			expected: 3,
		},
		{
			name:     "empty slice",
			input:    []int{},
			expected: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := Last(test.input)
			if err != nil && len(test.input) > 0 {
				t.Errorf("Expected true, got false")
			}
			if result != test.expected {
				t.Errorf("Expected %d, got %d", test.expected, result)
			}
		})
	}
}
