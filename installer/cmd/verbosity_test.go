package cmd

import (
	"testing"

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
