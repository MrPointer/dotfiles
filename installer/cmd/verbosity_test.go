package cmd

import (
	"testing"

	"github.com/MrPointer/dotfiles/installer/utils"
	"github.com/MrPointer/dotfiles/installer/utils/logger"
	"github.com/stretchr/testify/require"
)

func Test_VerbosityLevelDeterminationLogic(t *testing.T) {
	tests := []struct {
		name         string
		verboseCount int
		extraVerbose bool
		expected     logger.VerbosityLevel
	}{
		{
			name:         "Normal_WhenNoVerbosityFlags",
			verboseCount: 0,
			extraVerbose: false,
			expected:     logger.Normal,
		},
		{
			name:         "Verbose_WhenSingleVerboseFlag",
			verboseCount: 1,
			extraVerbose: false,
			expected:     logger.Verbose,
		},
		{
			name:         "ExtraVerbose_WhenDoubleVerboseFlag",
			verboseCount: 2,
			extraVerbose: false,
			expected:     logger.ExtraVerbose,
		},
		{
			name:         "ExtraVerbose_WhenTripleVerboseFlag",
			verboseCount: 3,
			extraVerbose: false,
			expected:     logger.ExtraVerbose,
		},
		{
			name:         "ExtraVerbose_WhenExtraVerboseFlagSet",
			verboseCount: 0,
			extraVerbose: true,
			expected:     logger.ExtraVerbose,
		},
		{
			name:         "ExtraVerbose_WhenBothExtraVerboseAndVerboseFlags",
			verboseCount: 1,
			extraVerbose: true,
			expected:     logger.ExtraVerbose,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset global state
			verboseCount = tt.verboseCount
			extraVerbose = tt.extraVerbose
			plainFlag = false
			globalVerbosity = logger.Normal

			// Call the verbosity determination logic
			initLogger()

			// Verify the result
			actual := GetVerbosity()
			require.Equal(t, tt.expected, actual)
		})
	}
}

func Test_CliLoggerCreationWithDifferentVerbosityLevels(t *testing.T) {
	tests := []struct {
		name      string
		verbosity logger.VerbosityLevel
	}{
		{
			name:      "Normal_VerbosityLevel",
			verbosity: logger.Normal,
		},
		{
			name:      "Verbose_VerbosityLevel",
			verbosity: logger.Verbose,
		},
		{
			name:      "ExtraVerbose_VerbosityLevel",
			verbosity: logger.ExtraVerbose,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create logger with specific verbosity level
			testLogger := logger.NewCliLogger(tt.verbosity)

			// Verify logger was created successfully
			require.NotNil(t, testLogger)
		})
	}
}

func Test_ShouldShowProgress_Logic(t *testing.T) {
	tests := []struct {
		name           string
		plainFlag      bool
		nonInteractive bool
		expected       bool
	}{
		{
			name:           "Progress_WhenNoFlags",
			plainFlag:      false,
			nonInteractive: false,
			expected:       true,
		},
		{
			name:           "NoProgress_WhenPlainFlagSet",
			plainFlag:      true,
			nonInteractive: false,
			expected:       false,
		},
		{
			name:           "NoProgress_WhenNonInteractiveSet",
			plainFlag:      false,
			nonInteractive: true,
			expected:       false,
		},
		{
			name:           "NoProgress_WhenBothPlainAndNonInteractiveSet",
			plainFlag:      true,
			nonInteractive: true,
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset global state
			plainFlag = tt.plainFlag
			nonInteractive = tt.nonInteractive

			// Verify the result
			actual := ShouldShowProgress()
			require.Equal(t, tt.expected, actual)
		})
	}
}

func Test_GetDisplayMode_Logic(t *testing.T) {
	tests := []struct {
		name                string
		nonInteractiveFlag  bool
		plainFlag           bool
		verbosity           logger.VerbosityLevel
		expectedDisplayMode utils.DisplayMode
	}{
		{
			name:                "Progress_WhenDefaultSettings",
			nonInteractiveFlag:  false,
			plainFlag:           false,
			verbosity:           logger.Normal,
			expectedDisplayMode: utils.DisplayModeProgress,
		},
		{
			name:                "Plain_WhenPlainFlagSet",
			nonInteractiveFlag:  false,
			plainFlag:           true,
			verbosity:           logger.Normal,
			expectedDisplayMode: utils.DisplayModePlain,
		},
		{
			name:                "Passthrough_WhenNonInteractiveSet",
			nonInteractiveFlag:  true,
			plainFlag:           false,
			verbosity:           logger.Normal,
			expectedDisplayMode: utils.DisplayModePassthrough,
		},
		{
			name:                "Passthrough_WhenVerboseSet",
			nonInteractiveFlag:  false,
			plainFlag:           false,
			verbosity:           logger.Verbose,
			expectedDisplayMode: utils.DisplayModePassthrough,
		},
		{
			name:                "Passthrough_WhenExtraVerboseSet",
			nonInteractiveFlag:  false,
			plainFlag:           false,
			verbosity:           logger.ExtraVerbose,
			expectedDisplayMode: utils.DisplayModePassthrough,
		},
		{
			name:                "Passthrough_WhenNonInteractiveAndPlainBothSet",
			nonInteractiveFlag:  true,
			plainFlag:           true,
			verbosity:           logger.Normal,
			expectedDisplayMode: utils.DisplayModePassthrough,
		},
		{
			name:                "Passthrough_WhenNonInteractiveAndVerboseSet",
			nonInteractiveFlag:  true,
			plainFlag:           false,
			verbosity:           logger.Verbose,
			expectedDisplayMode: utils.DisplayModePassthrough,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Store original values
			origNonInteractive := nonInteractive
			origPlainFlag := plainFlag
			origGlobalVerbosity := globalVerbosity

			// Set test values
			nonInteractive = tt.nonInteractiveFlag
			plainFlag = tt.plainFlag
			globalVerbosity = tt.verbosity

			// Test the function
			result := GetDisplayMode()
			require.Equal(t, tt.expectedDisplayMode, result)

			// Restore original values
			nonInteractive = origNonInteractive
			plainFlag = origPlainFlag
			globalVerbosity = origGlobalVerbosity
		})
	}
}
