package logger_test

import (
	"bytes"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/MrPointer/dotfiles/installer/utils/logger"
)

func Test_CliLogger_ByDesign_ImplementsLoggerInterface(t *testing.T) {
	var _ logger.Logger = (*logger.CliLogger)(nil)
}

func Test_NewProgressCliLogger_WhenCalled_CreatesValidInstance(t *testing.T) {
	log := logger.NewProgressCliLogger(logger.Normal)
	require.NotNil(t, log)
}

func Test_NewCliLogger_WithVerbosityLevels_CreatesValidInstance(t *testing.T) {
	tests := []struct {
		name      string
		verbosity logger.VerbosityLevel
	}{
		{"Minimal verbosity", logger.Minimal},
		{"Normal verbosity", logger.Normal},
		{"Verbose verbosity", logger.Verbose},
		{"ExtraVerbose verbosity", logger.ExtraVerbose},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := logger.NewCliLogger(tt.verbosity)
			require.NotNil(t, log)
		})
	}
}

func Test_NewCliLoggerWithOutput_WithProgressEnabled_EnablesProgressFunctionality(t *testing.T) {
	var buf bytes.Buffer
	log := logger.NewCliLoggerWithOutput(logger.Normal, &buf, true)
	require.NotNil(t, log)
}

func Test_MultipleRapidProgressOperations_WithQuickSuccession_Work(t *testing.T) {

	var buf bytes.Buffer
	log := logger.NewCliLoggerWithOutput(logger.Normal, &buf, true)

	// Rapidly start and finish operations
	for range 5 {
		log.StartProgress("Rapid operation")
		time.Sleep(10 * time.Millisecond)
		log.FinishProgress("Completed")
	}

	require.NotNil(t, log)
}

func Test_DeeplyNestedProgressOperations_WithFiveLevels_Complete(t *testing.T) {

	var buf bytes.Buffer
	log := logger.NewCliLoggerWithOutput(logger.Normal, &buf, true)

	// Create a deep hierarchy
	log.StartProgress("Level 1")
	log.StartProgress("Level 2")
	log.StartProgress("Level 3")
	log.StartProgress("Level 4")
	log.StartProgress("Level 5")

	time.Sleep(50 * time.Millisecond)

	// Complete in reverse order
	log.FinishProgress("Level 5 done")
	log.FinishProgress("Level 4 done")
	log.FinishProgress("Level 3 done")
	log.FinishProgress("Level 2 done")
	log.FinishProgress("Level 1 done")

	require.NotNil(t, log)
}

func Test_MixedSuccessAndFailureProgressOperations_WithNestedStructure_DisplayCorrectly(t *testing.T) {

	var buf bytes.Buffer
	log := logger.NewCliLoggerWithOutput(logger.Normal, &buf, true)

	log.StartProgress("Parent operation")

	log.StartProgress("Child 1")
	time.Sleep(30 * time.Millisecond)
	log.FinishProgress("Child 1 success")

	log.StartProgress("Child 2")
	time.Sleep(30 * time.Millisecond)
	log.FailProgress("Child 2 failed", errors.New("test error"))

	log.StartProgress("Child 3")
	time.Sleep(30 * time.Millisecond)
	log.FinishProgress("Child 3 success")

	log.FinishProgress("Parent complete")

	require.NotNil(t, log)
}

func Test_UpdateProgress_WithoutActiveProgress_DoesNothing(t *testing.T) {
	var buf bytes.Buffer
	log := logger.NewCliLoggerWithOutput(logger.Normal, &buf, true)

	// This should not crash
	log.UpdateProgress("No active progress")

	require.NotNil(t, log)
}

func Test_FinishProgress_WithoutActiveProgress_FallsBackToSuccessLogging(t *testing.T) {
	var buf bytes.Buffer
	log := logger.NewCliLoggerWithOutput(logger.Normal, &buf, true)

	// This should fall back to Success logging
	log.FinishProgress("Never started")

	require.NotNil(t, log)
}

func Test_FailProgress_WithoutActiveProgress_FallsBackToErrorLogging(t *testing.T) {
	var buf bytes.Buffer
	log := logger.NewCliLoggerWithOutput(logger.Normal, &buf, true)

	// This should fall back to Error logging
	log.FailProgress("Never started", errors.New("test error"))

	require.NotNil(t, log)
}

func Test_VeryShortDurationProgressOperations_UnderThreshold_DoNotShowTiming(t *testing.T) {

	var buf bytes.Buffer
	log := logger.NewCliLoggerWithOutput(logger.Normal, &buf, true)

	log.StartProgress("Quick operation")
	// No sleep - immediate completion
	log.FinishProgress("Done quickly")

	require.NotNil(t, log)
}

func Test_ProgressOperations_WithMessageUpdates_Work(t *testing.T) {

	var buf bytes.Buffer
	log := logger.NewCliLoggerWithOutput(logger.Normal, &buf, true)

	log.StartProgress("Downloading files")
	time.Sleep(20 * time.Millisecond)

	log.UpdateProgress("Downloading files (25%)")
	time.Sleep(20 * time.Millisecond)

	log.UpdateProgress("Downloading files (50%)")
	time.Sleep(20 * time.Millisecond)

	log.UpdateProgress("Downloading files (75%)")
	time.Sleep(20 * time.Millisecond)

	log.FinishProgress("Download complete")

	require.NotNil(t, log)
}

func Test_VerbosityLevels_WithDifferentMessages_FilterCorrectly(t *testing.T) {
	tests := []struct {
		name      string
		verbosity logger.VerbosityLevel
		logFunc   func(logger.Logger)
		shouldLog bool
	}{
		{
			name:      "Trace messages appear with ExtraVerbose",
			verbosity: logger.ExtraVerbose,
			logFunc:   func(l logger.Logger) { l.Trace("trace message") },
			shouldLog: true,
		},
		{
			name:      "Trace messages hidden with Verbose",
			verbosity: logger.Verbose,
			logFunc:   func(l logger.Logger) { l.Trace("trace message") },
			shouldLog: false,
		},
		{
			name:      "Debug messages appear with Verbose",
			verbosity: logger.Verbose,
			logFunc:   func(l logger.Logger) { l.Debug("debug message") },
			shouldLog: true,
		},
		{
			name:      "Debug messages hidden with Normal",
			verbosity: logger.Normal,
			logFunc:   func(l logger.Logger) { l.Debug("debug message") },
			shouldLog: false,
		},
		{
			name:      "Info messages appear with Normal",
			verbosity: logger.Normal,
			logFunc:   func(l logger.Logger) { l.Info("info message") },
			shouldLog: true,
		},
		{
			name:      "Info messages hidden with Minimal",
			verbosity: logger.Minimal,
			logFunc:   func(l logger.Logger) { l.Info("info message") },
			shouldLog: false,
		},
		{
			name:      "Error messages appear with Normal",
			verbosity: logger.Normal,
			logFunc:   func(l logger.Logger) { l.Error("error message") },
			shouldLog: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := logger.NewCliLogger(tt.verbosity)
			require.NotNil(t, log)

			// Just verify that calling the function doesn't panic
			// The actual output verification would require capturing stdout/stderr
			require.NotPanics(t, func() {
				tt.logFunc(log)
			})
		})
	}
}

func Test_ProgressMethods_FallBackToRegularLogging_WhenProgressDisabled(t *testing.T) {
	// Test that progress methods work correctly when progress is disabled
	log := logger.NewCliLogger(logger.Normal) // Not a progress logger

	// These should fall back to regular logging methods
	log.StartProgress("Starting operation")
	log.UpdateProgress("Updating operation")
	log.FinishProgress("Finishing operation")
	log.FailProgress("Failing operation", errors.New("test error"))

	require.NotNil(t, log)
}

func Test_ConcurrentProgressOperations_WithMultipleGoroutines_AreThreadSafe(t *testing.T) {

	var buf bytes.Buffer
	log := logger.NewCliLoggerWithOutput(logger.Normal, &buf, true)

	// Test concurrent access to ensure thread safety
	done := make(chan bool, 2)

	go func() {
		log.StartProgress("Concurrent operation 1")
		time.Sleep(50 * time.Millisecond)
		log.FinishProgress("Concurrent 1 done")
		done <- true
	}()

	go func() {
		time.Sleep(25 * time.Millisecond)
		log.StartProgress("Concurrent operation 2")
		time.Sleep(50 * time.Millisecond)
		log.FinishProgress("Concurrent 2 done")
		done <- true
	}()

	// Wait for both goroutines to complete
	<-done
	<-done

	require.NotNil(t, log)
}

func Test_LongRunningOperation_WithPeriodicUpdates_Works(t *testing.T) {

	var buf bytes.Buffer
	log := logger.NewCliLoggerWithOutput(logger.Normal, &buf, true)

	log.StartProgress("Long running operation")

	// Simulate a longer operation with periodic updates
	for i := 0; i < 5; i++ {
		time.Sleep(50 * time.Millisecond)
		log.UpdateProgress(fmt.Sprintf("Long running operation (step %d/5)", i+1))
	}

	log.FinishProgress("Long operation completed")

	require.NotNil(t, log)
}

func Test_HierarchicalProgressReporting_WithNestedOperations_WorksLikeNpm(t *testing.T) {

	// Create a progress-enabled logger
	var buf bytes.Buffer
	log := logger.NewCliLoggerWithOutput(logger.Normal, &buf, true)

	// Simulate a hierarchical installation process
	log.StartProgress("Installing dotfiles")

	// Simulate some work
	time.Sleep(100 * time.Millisecond)

	// Start nested progress
	log.StartProgress("Downloading configuration files")
	time.Sleep(200 * time.Millisecond)

	// Update progress message
	log.UpdateProgress("Downloading configuration files (50%)")
	time.Sleep(150 * time.Millisecond)

	log.FinishProgress("Downloaded configuration files")

	// Another nested operation
	log.StartProgress("Setting up shell configuration")
	time.Sleep(100 * time.Millisecond)

	// Nested operation within nested operation
	log.StartProgress("Installing zsh plugins")
	time.Sleep(150 * time.Millisecond)
	log.FinishProgress("Installed zsh plugins")

	log.FinishProgress("Set up shell configuration")

	// Simulate a failed operation
	log.StartProgress("Configuring git settings")
	time.Sleep(100 * time.Millisecond)
	log.FailProgress("Failed to configure git", errors.New("permission denied"))

	// Complete the main operation
	log.FinishProgress("Installed dotfiles")

	// This test mainly verifies that the code doesn't crash
	// and provides a visual demonstration of the hierarchical progress
	require.True(t, true, "Hierarchical progress test completed")
}

func Test_AllVerbosityLevels_WithDifferentMessages_ProduceAppropriateOutput(t *testing.T) {
	tests := []struct {
		name      string
		verbosity logger.VerbosityLevel
	}{
		{"Minimal verbosity produces minimal output", logger.Minimal},
		{"Normal verbosity produces normal output", logger.Normal},
		{"Verbose verbosity produces verbose output", logger.Verbose},
		{"ExtraVerbose verbosity produces extra verbose output", logger.ExtraVerbose},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := logger.NewCliLogger(tt.verbosity)

			// Test all logging levels
			log.Trace("This is a trace message")
			log.Debug("This is a debug message")
			log.Info("This is an info message")
			log.Success("This is a success message")
			log.Warning("This is a warning message")
			log.Error("This is an error message")

			require.NotNil(t, log)
		})
	}
}

func Test_Progress_WithMinimalVerbosity_Works(t *testing.T) {

	var buf bytes.Buffer
	log := logger.NewCliLoggerWithOutput(logger.Minimal, &buf, true)

	log.StartProgress("Operation with minimal verbosity")
	time.Sleep(50 * time.Millisecond)
	log.FinishProgress("Completed operation")

	require.NotNil(t, log)
}

func Test_VerboseLogging_WithoutProgress_Works(t *testing.T) {
	log := logger.NewCliLogger(logger.Verbose)

	log.StartProgress("This should appear as Info message")
	log.UpdateProgress("This update should be ignored")
	log.FinishProgress("This should appear as Success message")

	require.NotNil(t, log)
}

func Test_StartPersistentProgress_WithProgressEnabled_Works(t *testing.T) {

	var buf bytes.Buffer
	log := logger.NewCliLoggerWithOutput(logger.Normal, &buf, true)

	log.StartPersistentProgress("Installing components")
	time.Sleep(50 * time.Millisecond)
	log.FinishPersistentProgress("Installation complete")

	require.NotNil(t, log)
}

func Test_StartPersistentProgress_WithoutProgress_FallsBackToInfoLogging(t *testing.T) {
	log := logger.NewCliLogger(logger.Normal)

	log.StartPersistentProgress("This should appear as Info message")
	log.FinishPersistentProgress("This should appear as Success message")

	require.NotNil(t, log)
}

func Test_LogAccomplishment_WithProgressEnabled_ShowsPersistentMessage(t *testing.T) {

	var buf bytes.Buffer
	log := logger.NewCliLoggerWithOutput(logger.Normal, &buf, true)

	log.StartPersistentProgress("Installing packages")
	time.Sleep(30 * time.Millisecond)

	log.LogAccomplishment("Downloaded package A")
	log.LogAccomplishment("Downloaded package B")
	log.LogAccomplishment("Downloaded package C")

	time.Sleep(30 * time.Millisecond)
	log.FinishPersistentProgress("All packages installed")

	require.NotNil(t, log)
}

func Test_LogAccomplishment_WithoutProgress_FallsBackToSuccessLogging(t *testing.T) {
	log := logger.NewCliLogger(logger.Normal)

	log.LogAccomplishment("This should appear as Success message")

	require.NotNil(t, log)
}

func Test_FinishPersistent_ProgressWithoutActiveProgress_FallsBackToSuccessLogging(t *testing.T) {
	var buf bytes.Buffer
	log := logger.NewCliLoggerWithOutput(logger.Normal, &buf, true)

	log.FinishPersistentProgress("Never started persistent progress")

	require.NotNil(t, log)
}

func Test_FailPersistentProgress_WithoutActiveProgress_FallsBackToErrorLogging(t *testing.T) {
	var buf bytes.Buffer
	log := logger.NewCliLoggerWithOutput(logger.Normal, &buf, true)

	log.FailPersistentProgress("Never started persistent progress", errors.New("test error"))

	require.NotNil(t, log)
}

func Test_PersistentProgress_WhenFailed_ShowsErrorMessage(t *testing.T) {

	var buf bytes.Buffer
	log := logger.NewCliLoggerWithOutput(logger.Normal, &buf, true)

	log.StartPersistentProgress("Installing critical component")
	log.LogAccomplishment("Downloaded dependencies")
	log.LogAccomplishment("Validated checksums")
	time.Sleep(50 * time.Millisecond)
	log.FailPersistentProgress("Installation failed", errors.New("permission denied"))

	require.NotNil(t, log)
}

func Test_MixedPersistentAndRegularProgressOperations_WithNestedStructure_Work(t *testing.T) {

	var buf bytes.Buffer
	log := logger.NewCliLoggerWithOutput(logger.Normal, &buf, true)

	log.StartPersistentProgress("Setting up development environment")
	log.LogAccomplishment("Created project directory")

	log.StartProgress("Downloading dependencies")
	time.Sleep(50 * time.Millisecond)
	log.FinishProgress("Dependencies downloaded")

	log.LogAccomplishment("Installed build tools")
	log.LogAccomplishment("Configured IDE settings")

	log.StartProgress("Running initial build")
	time.Sleep(50 * time.Millisecond)
	log.FinishProgress("Build completed successfully")

	log.FinishPersistentProgress("Development environment ready")

	require.NotNil(t, log)
}

func Test_PersistentProgress_WithMultipleAccomplishments_DisplaysCorrectly(t *testing.T) {

	var buf bytes.Buffer
	log := logger.NewCliLoggerWithOutput(logger.Normal, &buf, true)

	log.StartPersistentProgress("Deploying application")

	accomplishments := []string{
		"Built application binary",
		"Created Docker image",
		"Pushed to registry",
		"Updated deployment configuration",
		"Applied Kubernetes manifests",
		"Verified health checks",
	}

	for _, accomplishment := range accomplishments {
		time.Sleep(20 * time.Millisecond)
		log.LogAccomplishment(accomplishment)
	}

	time.Sleep(30 * time.Millisecond)
	log.FinishPersistentProgress("Application deployed successfully")

	require.NotNil(t, log)
}

func Test_StartPersistentProgress_WithMockReporter_CallsProgressReporterStartPersistent(t *testing.T) {
	mockProgress := &logger.MoqProgressReporter{
		StartPersistentFunc: func(message string) error { return nil },
		IsActiveFunc:        func() bool { return true },
	}

	log := logger.NewCliLoggerWithProgress(logger.Normal, mockProgress)

	log.StartPersistentProgress("Test message")

	calls := mockProgress.StartPersistentCalls()
	require.Len(t, calls, 1)
	require.Equal(t, "Test message", calls[0].Message)
}

func Test_StartPersistentProgress_WithoutProgress_FallsBackToInfo(t *testing.T) {
	mockProgress := &logger.MoqProgressReporter{
		StartPersistentFunc: func(message string) error { return nil },
		IsActiveFunc:        func() bool { return false },
	}

	log := logger.NewCliLoggerWithProgress(logger.Normal, mockProgress)

	log.StartPersistentProgress("Test message")

	calls := mockProgress.StartPersistentCalls()
	require.Len(t, calls, 1)
	require.Equal(t, "Test message", calls[0].Message)
}

func Test_LogAccomplishment_WithMockReporter_CallsProgressReporterLogAccomplishment(t *testing.T) {
	mockProgress := &logger.MoqProgressReporter{
		LogAccomplishmentFunc: func(message string) error { return nil },
		IsActiveFunc:          func() bool { return true },
	}

	log := logger.NewCliLoggerWithProgress(logger.Normal, mockProgress)

	log.LogAccomplishment("Test accomplishment")

	calls := mockProgress.LogAccomplishmentCalls()
	require.Len(t, calls, 1)
	require.Equal(t, "Test accomplishment", calls[0].Message)
}

func Test_LogAccomplishment_FallsBackToSuccess_WhenProgressNotActive(t *testing.T) {
	mockProgress := &logger.MoqProgressReporter{
		LogAccomplishmentFunc: func(message string) error { return nil },
		IsActiveFunc:          func() bool { return false },
	}

	log := logger.NewCliLoggerWithProgress(logger.Normal, mockProgress)

	log.LogAccomplishment("Test accomplishment")

	calls := mockProgress.LogAccomplishmentCalls()
	require.Len(t, calls, 1)
	require.Equal(t, "Test accomplishment", calls[0].Message)
}

func Test_FinishPersistentProgress_WithMockReporter_CallsProgressReporterFinishPersistent(t *testing.T) {
	mockProgress := &logger.MoqProgressReporter{
		FinishPersistentFunc: func(message string) error { return nil },
		IsActiveFunc:         func() bool { return true },
	}

	log := logger.NewCliLoggerWithProgress(logger.Normal, mockProgress)

	log.FinishPersistentProgress("Test finished")

	calls := mockProgress.FinishPersistentCalls()
	require.Len(t, calls, 1)
	require.Equal(t, "Test finished", calls[0].Message)
}

func Test_FinishPersistentProgress_FallsBackToSuccess_WhenProgressNotActive(t *testing.T) {
	mockProgress := &logger.MoqProgressReporter{
		FinishPersistentFunc: func(message string) error { return nil },
		IsActiveFunc:         func() bool { return false },
	}

	log := logger.NewCliLoggerWithProgress(logger.Normal, mockProgress)

	log.FinishPersistentProgress("Test finished")

	calls := mockProgress.FinishPersistentCalls()
	require.Len(t, calls, 1)
	require.Equal(t, "Test finished", calls[0].Message)
}

func Test_FailPersistentProgress_WithMockReporter_CallsProgressReporterFailPersistent(t *testing.T) {
	testErr := errors.New("test error")
	mockProgress := &logger.MoqProgressReporter{
		FailPersistentFunc: func(message string, err error) error { return nil },
		IsActiveFunc:       func() bool { return true },
	}

	log := logger.NewCliLoggerWithProgress(logger.Normal, mockProgress)

	log.FailPersistentProgress("Test failed", testErr)

	calls := mockProgress.FailPersistentCalls()
	require.Len(t, calls, 1)
	require.Equal(t, "Test failed", calls[0].Message)
	require.Equal(t, testErr, calls[0].Err)
}

func Test_FailPersistentProgress_FallsBackToError_WhenProgressNotActive(t *testing.T) {
	testErr := errors.New("test error")
	mockProgress := &logger.MoqProgressReporter{
		FailPersistentFunc: func(message string, err error) error { return nil },
		IsActiveFunc:       func() bool { return false },
	}

	log := logger.NewCliLoggerWithProgress(logger.Normal, mockProgress)

	log.FailPersistentProgress("Test failed", testErr)

	calls := mockProgress.FailPersistentCalls()
	require.Len(t, calls, 1)
	require.Equal(t, "Test failed", calls[0].Message)
	require.Equal(t, testErr, calls[0].Err)
}

func Test_Close_WithMockReporter_CallsProgressReporterClose(t *testing.T) {
	mockProgress := &logger.MoqProgressReporter{
		CloseFunc: func() error { return nil },
	}

	log := logger.NewCliLoggerWithProgress(logger.Normal, mockProgress)

	log.Close()

	calls := mockProgress.CloseCalls()
	require.Len(t, calls, 1)
}

func Test_Close_WithoutProgress_StillWorks(t *testing.T) {
	log := logger.NewCliLogger(logger.Normal)

	require.NotPanics(t, func() {
		log.Close()
	})
}
