package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_DisplayMode_String_ReturnsCorrectStringRepresentation(t *testing.T) {
	tests := []struct {
		name        string
		displayMode DisplayMode
		expectedStr string
	}{
		{
			name:        "Progress_Mode",
			displayMode: DisplayModeProgress,
			expectedStr: "progress",
		},
		{
			name:        "Plain_Mode",
			displayMode: DisplayModePlain,
			expectedStr: "plain",
		},
		{
			name:        "Passthrough_Mode",
			displayMode: DisplayModePassthrough,
			expectedStr: "passthrough",
		},
		{
			name:        "Unknown_Mode",
			displayMode: DisplayMode(999),
			expectedStr: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.displayMode.String()
			require.Equal(t, tt.expectedStr, result)
		})
	}
}

func Test_DisplayMode_ShouldDiscardOutput_ReturnsCorrectValue(t *testing.T) {
	tests := []struct {
		name          string
		displayMode   DisplayMode
		shouldDiscard bool
	}{
		{
			name:          "Progress_ShouldDiscardOutput",
			displayMode:   DisplayModeProgress,
			shouldDiscard: true,
		},
		{
			name:          "Plain_ShouldDiscardOutput",
			displayMode:   DisplayModePlain,
			shouldDiscard: true,
		},
		{
			name:          "Passthrough_ShouldNotDiscardOutput",
			displayMode:   DisplayModePassthrough,
			shouldDiscard: false,
		},
		{
			name:          "Unknown_ShouldDiscardOutput_ByDefault",
			displayMode:   DisplayMode(999),
			shouldDiscard: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.displayMode.ShouldDiscardOutput()
			require.Equal(t, tt.shouldDiscard, result)
		})
	}
}
