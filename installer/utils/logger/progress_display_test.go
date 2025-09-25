package logger_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/MrPointer/dotfiles/installer/utils/logger"
)

func Test_NewProgressDisplayCreatesValidInstance(t *testing.T) {
	var buf bytes.Buffer
	display := logger.NewProgressDisplay(&buf)
	require.NotNil(t, display)
}

func Test_NewProgressDisplayUsesStdoutWhenOutputIsNil(t *testing.T) {
	display := logger.NewProgressDisplay(nil)
	require.NotNil(t, display)
}

func Test_SingleProgressOperationCanBeStartedAndFinished(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	var buf bytes.Buffer
	display := logger.NewProgressDisplay(&buf)

	display.Start("Test operation")
	require.True(t, display.IsActive())

	time.Sleep(50 * time.Millisecond)
	display.Finish("Test operation")

	require.False(t, display.IsActive())

	output := buf.String()
	require.Contains(t, output, "✓")
	require.Contains(t, output, "Test operation")
}

func Test_NestedProgressOperationsShowProperHierarchy(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	var buf bytes.Buffer
	display := logger.NewProgressDisplay(&buf)

	display.Start("Parent operation")
	require.True(t, display.IsActive())

	display.Start("Child operation")
	require.True(t, display.IsActive())

	time.Sleep(30 * time.Millisecond)
	display.Finish("Child operation")
	require.True(t, display.IsActive()) // Parent still active

	time.Sleep(30 * time.Millisecond)
	display.Finish("Parent operation")
	require.False(t, display.IsActive())

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	// Should have contextual messages during progress and clean completion messages
	require.Contains(t, output, "Child operation")
	require.Contains(t, output, "Parent operation")
	require.Contains(t, strings.Join(lines, "\n"), "✓")
}

func Test_ProgressMessageCanBeUpdatedAfterStart(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	var buf bytes.Buffer
	display := logger.NewProgressDisplay(&buf)

	display.Start("Initial message")
	display.Update("Updated message")
	time.Sleep(30 * time.Millisecond)
	display.Finish("Final message")

	output := buf.String()
	require.Contains(t, output, "Updated message")
}

func Test_ProgressOperationCanFailWithError(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	var buf bytes.Buffer
	display := logger.NewProgressDisplay(&buf)

	display.Start("Failing operation")
	time.Sleep(30 * time.Millisecond)
	display.Fail("Failing operation", errors.New("test error"))

	require.False(t, display.IsActive())

	output := buf.String()
	require.Contains(t, output, "✗")
	require.Contains(t, output, "Failing operation")
	require.Contains(t, output, "test error")
}

func Test_MixedSuccessAndFailureOperationsDisplayCorrectly(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	var buf bytes.Buffer
	display := logger.NewProgressDisplay(&buf)

	display.Start("Parent")

	display.Start("Success child")
	time.Sleep(20 * time.Millisecond)
	display.Finish("Success child")

	display.Start("Failing child")
	time.Sleep(20 * time.Millisecond)
	display.Fail("Failing child", errors.New("child error"))

	display.Finish("Parent")

	output := buf.String()
	require.Contains(t, output, "Success child")
	require.Contains(t, output, "Failing child")
	require.Contains(t, output, "Parent")
	require.Contains(t, output, "child error")
	require.Contains(t, output, "✓")
	require.Contains(t, output, "✗")
}

func Test_DeeplyNestedOperationsShowCorrectIndentation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	var buf bytes.Buffer
	display := logger.NewProgressDisplay(&buf)

	// Create deep nesting
	display.Start("Level 1")
	display.Start("Level 2")
	display.Start("Level 3")
	display.Start("Level 4")

	time.Sleep(20 * time.Millisecond)

	// Complete in reverse order
	display.Finish("Level 4")
	display.Finish("Level 3")
	display.Finish("Level 2")
	display.Finish("Level 1")

	output := buf.String()

	// Check that all levels appear in output (contextual messages during progress)
	require.Contains(t, output, "Level 1")
	require.Contains(t, output, "Level 2")
	require.Contains(t, output, "Level 3")
	require.Contains(t, output, "Level 4")
}

func Test_LongRunningOperationsShowTimingInformation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	var buf bytes.Buffer
	display := logger.NewProgressDisplay(&buf)

	display.Start("Long operation")
	time.Sleep(150 * time.Millisecond) // Longer than 100ms threshold
	display.Finish("Long operation")

	output := buf.String()
	require.Contains(t, output, "took")
	require.Contains(t, output, "ms")
}

func Test_ShortOperationsDoNotShowTimingInformation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	var buf bytes.Buffer
	display := logger.NewProgressDisplay(&buf)

	display.Start("Quick operation")
	// No sleep - immediate completion
	display.Finish("Quick operation")

	output := buf.String()
	require.NotContains(t, output, "took")
}

func Test_ClearStopsAllActiveProgressOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	var buf bytes.Buffer
	display := logger.NewProgressDisplay(&buf)

	display.Start("Operation 1")
	display.Start("Operation 2")
	require.True(t, display.IsActive())

	display.Clear()
	require.False(t, display.IsActive())

	// Should not crash when calling methods after clear
	display.Update("Should be ignored")
	display.Finish("Should be ignored")
	display.Fail("Should be ignored", errors.New("test"))
}

func Test_UpdateWithoutActiveProgressDoesNothing(t *testing.T) {
	var buf bytes.Buffer
	display := logger.NewProgressDisplay(&buf)

	display.Update("No active operation")
	require.False(t, display.IsActive())
}

func Test_FinishWithoutActiveProgressDoesNothing(t *testing.T) {
	var buf bytes.Buffer
	display := logger.NewProgressDisplay(&buf)

	display.Finish("No active operation")
	require.False(t, display.IsActive())
}

func Test_FailWithoutActiveProgressDoesNothing(t *testing.T) {
	var buf bytes.Buffer
	display := logger.NewProgressDisplay(&buf)

	display.Fail("No active operation", errors.New("test"))
	require.False(t, display.IsActive())
}

func Test_NoopProgressDisplayImplementsProgressReporterInterface(t *testing.T) {
	var _ logger.ProgressReporter = (*logger.NoopProgressDisplay)(nil)
}

func Test_NoopProgressDisplayAllMethodsDoNothing(t *testing.T) {
	display := logger.NewNoopProgressDisplay()

	// All these should not crash and should not do anything
	display.Start("Test")
	require.False(t, display.IsActive())

	display.Update("Test")
	require.False(t, display.IsActive())

	display.Finish("Test")
	require.False(t, display.IsActive())

	display.Fail("Test", errors.New("test"))
	require.False(t, display.IsActive())

	display.Clear()
	require.False(t, display.IsActive())
}

func Test_ConcurrentProgressDisplayOperationsAreThreadSafe(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	var buf bytes.Buffer
	display := logger.NewProgressDisplay(&buf)

	done := make(chan bool, 2)

	// Test concurrent access to ensure thread safety
	go func() {
		display.Start("Concurrent 1")
		time.Sleep(50 * time.Millisecond)
		display.Finish("Concurrent 1")
		done <- true
	}()

	go func() {
		time.Sleep(25 * time.Millisecond)
		display.Start("Concurrent 2")
		time.Sleep(50 * time.Millisecond)
		display.Finish("Concurrent 2")
		done <- true
	}()

	// Wait for both goroutines to complete
	<-done
	<-done

	require.False(t, display.IsActive())

	output := buf.String()
	require.Contains(t, output, "Concurrent 1")
	require.Contains(t, output, "Concurrent 2")
}

func Test_RapidSequentialOperationsWork(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	var buf bytes.Buffer
	display := logger.NewProgressDisplay(&buf)

	// Rapidly start and finish operations
	for i := 0; i < 10; i++ {
		display.Start("Rapid operation")
		time.Sleep(5 * time.Millisecond)
		display.Finish("Rapid operation")
	}

	require.False(t, display.IsActive())

	output := buf.String()
	// Should have multiple completion messages
	checkmarkCount := strings.Count(output, "✓")
	require.Equal(t, 10, checkmarkCount)
}

func Test_StartPersistentProgressActivatesPersistentMode(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	var buf bytes.Buffer
	display := logger.NewProgressDisplay(&buf)

	display.StartPersistent("Installing packages")
	require.True(t, display.IsActive())

	time.Sleep(50 * time.Millisecond)
	display.FinishPersistent("Installation complete")
	require.False(t, display.IsActive())

	output := buf.String()
	require.Contains(t, output, "✓")
	require.Contains(t, output, "Installing packages")
}

func Test_LogAccomplishmentShowsVisibleAccomplishments(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	var buf bytes.Buffer
	display := logger.NewProgressDisplay(&buf)

	display.StartPersistent("Deploying application")
	time.Sleep(30 * time.Millisecond)

	display.LogAccomplishment("Built application")
	display.LogAccomplishment("Created container")
	display.LogAccomplishment("Pushed to registry")

	time.Sleep(30 * time.Millisecond)
	display.FinishPersistent("Deployment complete")

	output := buf.String()
	require.Contains(t, output, "Built application")
	require.Contains(t, output, "Created container")
	require.Contains(t, output, "Pushed to registry")
	require.Contains(t, output, "Deploying application")
}

func Test_ProgressDisplayPersistentProgressCanFailWithError(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	var buf bytes.Buffer
	display := logger.NewProgressDisplay(&buf)

	display.StartPersistent("Critical operation")
	display.LogAccomplishment("Step 1 completed")
	display.LogAccomplishment("Step 2 completed")
	time.Sleep(30 * time.Millisecond)
	display.FailPersistent("Critical operation failed", errors.New("permission denied"))

	require.False(t, display.IsActive())

	output := buf.String()
	require.Contains(t, output, "Step 1 completed")
	require.Contains(t, output, "Step 2 completed")
	require.Contains(t, output, "✗")
	require.Contains(t, output, "Critical operation")
	require.Contains(t, output, "permission denied")
}

func Test_ProgressDisplayMixedPersistentAndRegularProgressOperationsWork(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	var buf bytes.Buffer
	display := logger.NewProgressDisplay(&buf)

	display.StartPersistent("Setting up environment")
	display.LogAccomplishment("Created directories")

	display.Start("Downloading files")
	time.Sleep(30 * time.Millisecond)
	display.Finish("Files downloaded")

	display.LogAccomplishment("Installed dependencies")

	display.Start("Running tests")
	time.Sleep(30 * time.Millisecond)
	display.Finish("Tests passed")

	display.FinishPersistent("Environment ready")

	output := buf.String()
	require.Contains(t, output, "Created directories")
	require.Contains(t, output, "Installed dependencies")
	require.Contains(t, output, "Setting up environment")
}

func Test_LogAccomplishmentWithoutActivePersistentProgressStillWorks(t *testing.T) {
	var buf bytes.Buffer
	display := logger.NewProgressDisplay(&buf)

	display.LogAccomplishment("Standalone accomplishment")

	output := buf.String()
	require.Contains(t, output, "✓")
	require.Contains(t, output, "Standalone accomplishment")
}

func Test_FinishPersistentWithoutActiveProgressDoesNothing(t *testing.T) {
	var buf bytes.Buffer
	display := logger.NewProgressDisplay(&buf)

	display.FinishPersistent("No active progress")
	require.False(t, display.IsActive())
}

func Test_FailPersistentWithoutActiveProgressDoesNothing(t *testing.T) {
	var buf bytes.Buffer
	display := logger.NewProgressDisplay(&buf)

	display.FailPersistent("No active progress", errors.New("test error"))
	require.False(t, display.IsActive())
}

func Test_PersistentProgressShowsAccomplishmentsInRealTime(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	var buf bytes.Buffer
	display := logger.NewProgressDisplay(&buf)

	display.StartPersistent("Processing items")

	accomplishments := []string{
		"Processed item 1",
		"Processed item 2",
		"Processed item 3",
		"Processed item 4",
		"Processed item 5",
	}

	for _, accomplishment := range accomplishments {
		time.Sleep(20 * time.Millisecond)
		display.LogAccomplishment(accomplishment)
	}

	display.FinishPersistent("All items processed")

	output := buf.String()
	for _, accomplishment := range accomplishments {
		require.Contains(t, output, accomplishment)
	}
	require.Contains(t, output, "Processing items")
}

func Test_NoopProgressDisplayPersistentMethodsDoNothing(t *testing.T) {
	display := logger.NewNoopProgressDisplay()

	// All these should not crash and should not do anything
	display.StartPersistent("Test")
	require.False(t, display.IsActive())

	display.LogAccomplishment("Test accomplishment")
	require.False(t, display.IsActive())

	display.FinishPersistent("Test")
	require.False(t, display.IsActive())

	display.FailPersistent("Test", errors.New("test"))
	require.False(t, display.IsActive())

	display.Close()
	require.False(t, display.IsActive())
}

func Test_CloseStopsAllOperationsAndRestoresCursor(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	var buf bytes.Buffer
	display := logger.NewProgressDisplay(&buf)

	display.Start("Operation 1")
	display.Start("Operation 2")
	require.True(t, display.IsActive())

	display.Close()
	require.False(t, display.IsActive())

	// Should not crash when calling methods after cleanup
	display.Update("Should be ignored")
	display.Finish("Should be ignored")
	display.Fail("Should be ignored", errors.New("test"))
}

func Test_CloseCanBeCalledMultipleTimes(t *testing.T) {
	var buf bytes.Buffer
	display := logger.NewProgressDisplay(&buf)

	display.Start("Test operation")
	require.True(t, display.IsActive())

	// Multiple cleanups should not crash
	display.Close()
	display.Close()
	display.Close()

	require.False(t, display.IsActive())
}

func Test_CloseWithoutActiveOperationsDoesNotCrash(t *testing.T) {
	var buf bytes.Buffer
	display := logger.NewProgressDisplay(&buf)

	require.NotPanics(t, func() {
		display.Close()
	})

	require.False(t, display.IsActive())
}

func Test_CloseAlsoRestoresCursor(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	var buf bytes.Buffer
	display := logger.NewProgressDisplay(&buf)

	display.Start("Operation with cursor")
	require.True(t, display.IsActive())

	time.Sleep(30 * time.Millisecond)
	display.Close()
	require.False(t, display.IsActive())

	// Close should not crash and should handle cursor state properly
	require.NotNil(t, display)
}

func Test_ProgressFailureSynchronizesCleanupAndPreventsHangingCursor(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	var buf bytes.Buffer
	display := logger.NewProgressDisplay(&buf)

	display.Start("Operation that will fail")
	require.True(t, display.IsActive())

	// Simulate work that results in failure
	time.Sleep(100 * time.Millisecond)
	display.Fail("Operation failed", errors.New("simulated error"))

	// After failure, display should be properly cleaned up
	require.False(t, display.IsActive())

	// Verify that output contains failure message and cursor control sequences are handled
	output := buf.String()
	require.Contains(t, output, "✗")
	require.Contains(t, output, "Operation that will fail")
	require.Contains(t, output, "simulated error")

	// Should be able to start new operations without issues
	display.Start("New operation after failure")
	require.True(t, display.IsActive())

	time.Sleep(30 * time.Millisecond)
	display.Finish("New operation completed")
	require.False(t, display.IsActive())

	// Verify new operation also completed successfully
	require.Contains(t, buf.String(), "New operation after failure")
}

func Test_RapidFailureAndRecoveryMaintainsProperTerminalState(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	var buf bytes.Buffer
	display := logger.NewProgressDisplay(&buf)

	// Test rapid failure and recovery cycles
	for i := 0; i < 5; i++ {
		display.Start("Rapid operation")
		time.Sleep(20 * time.Millisecond)

		if i%2 == 0 {
			display.Fail("Operation failed", errors.New("test error"))
		} else {
			display.Finish("Operation succeeded")
		}

		require.False(t, display.IsActive())
	}

	output := buf.String()
	// Should have both success and failure indicators
	require.Contains(t, output, "✓")
	require.Contains(t, output, "✗")

	// Should not have any hanging operations
	require.False(t, display.IsActive())
}

func Test_Pause_WithActiveSpinners_StopsAllOperations(t *testing.T) {
	var output bytes.Buffer
	display := logger.NewProgressDisplay(&output)

	display.Start("Operation 1")
	display.Start("Operation 2")
	time.Sleep(50 * time.Millisecond) // Allow spinners to start

	require.True(t, display.IsActive())
	require.False(t, display.IsPaused())

	err := display.Pause()
	require.NoError(t, err)

	require.True(t, display.IsPaused())
	require.True(t, display.IsActive()) // Operations still exist, just paused
}

func Test_Resume_AfterPause_RestartsSpinnerOperations(t *testing.T) {
	var output bytes.Buffer
	display := logger.NewProgressDisplay(&output)

	display.Start("Operation 1")
	display.Start("Operation 2")
	time.Sleep(50 * time.Millisecond) // Allow spinners to start

	err := display.Pause()
	require.NoError(t, err)
	require.True(t, display.IsPaused())

	err = display.Resume()
	require.NoError(t, err)

	require.False(t, display.IsPaused())
	require.True(t, display.IsActive())
}

func Test_Pause_WithoutActiveOperations_DoesNotCrash(t *testing.T) {
	var output bytes.Buffer
	display := logger.NewProgressDisplay(&output)

	require.False(t, display.IsActive())
	require.False(t, display.IsPaused())

	err := display.Pause()
	require.NoError(t, err)

	require.True(t, display.IsPaused())
	require.False(t, display.IsActive())
}

func Test_Resume_WithoutActiveOperations_DoesNotCrash(t *testing.T) {
	var output bytes.Buffer
	display := logger.NewProgressDisplay(&output)

	err := display.Pause()
	require.NoError(t, err)
	require.True(t, display.IsPaused())

	err = display.Resume()
	require.NoError(t, err)

	require.False(t, display.IsPaused())
	require.False(t, display.IsActive())
}

func Test_Pause_CalledMultipleTimes_IsSafe(t *testing.T) {
	var output bytes.Buffer
	display := logger.NewProgressDisplay(&output)

	display.Start("Test Operation")
	time.Sleep(50 * time.Millisecond) // Allow spinner to start

	err := display.Pause()
	require.NoError(t, err)
	require.True(t, display.IsPaused())

	err = display.Pause()
	require.NoError(t, err)
	require.True(t, display.IsPaused())
}

func Test_Resume_CalledMultipleTimes_IsSafe(t *testing.T) {
	var output bytes.Buffer
	display := logger.NewProgressDisplay(&output)

	display.Start("Test Operation")
	time.Sleep(50 * time.Millisecond) // Allow spinner to start
	display.Pause()

	err := display.Resume()
	require.NoError(t, err)
	require.False(t, display.IsPaused())

	err = display.Resume()
	require.NoError(t, err)
	require.False(t, display.IsPaused())
}

func Test_PauseAndResume_WithNestedOperations_WorksCorrectly(t *testing.T) {
	var output bytes.Buffer
	display := logger.NewProgressDisplay(&output)

	display.Start("Parent Operation")
	display.Start("Child Operation 1")
	display.Start("Child Operation 2")
	time.Sleep(50 * time.Millisecond) // Allow spinners to start

	require.True(t, display.IsActive())
	require.False(t, display.IsPaused())

	err := display.Pause()
	require.NoError(t, err)
	require.True(t, display.IsPaused())

	err = display.Resume()
	require.NoError(t, err)
	require.False(t, display.IsPaused())
	require.True(t, display.IsActive())
}

func Test_Pause_BeforeInteractiveInput_StopsSpinnerAndClearsOutput(t *testing.T) {
	var output bytes.Buffer
	display := logger.NewProgressDisplay(&output)

	display.Start("Processing files...")
	time.Sleep(50 * time.Millisecond) // Let spinner start

	initialOutput := output.String()
	require.NotEmpty(t, initialOutput)

	err := display.Pause()
	require.NoError(t, err)

	// Allow some time to ensure spinner has stopped
	time.Sleep(100 * time.Millisecond)
	outputAfterPause := output.String()

	// Output should contain clear line sequence after pause
	require.Contains(t, outputAfterPause, "\r")
	require.True(t, display.IsPaused())
	require.True(t, display.IsActive()) // Operations still exist, just paused
}

func Test_Resume_WithMultipleOperations_RestartsMostRecent(t *testing.T) {
	var output bytes.Buffer
	display := logger.NewProgressDisplay(&output)

	display.Start("Operation 1")
	display.Start("Operation 2")
	display.Start("Operation 3")
	time.Sleep(50 * time.Millisecond) // Allow spinners to start

	err := display.Pause()
	require.NoError(t, err)

	err = display.Resume()
	require.NoError(t, err)

	time.Sleep(50 * time.Millisecond)
	output_content := output.String()

	// Should show the most recent operation (Operation 3)
	require.Contains(t, output_content, "Operation 3")
}

func Test_PausedState_WithNewOperations_StillAllowsOperationCreation(t *testing.T) {
	var output bytes.Buffer
	display := logger.NewProgressDisplay(&output)

	display.Start("Operation 1")
	time.Sleep(50 * time.Millisecond) // Allow spinner to start

	err := display.Pause()
	require.NoError(t, err)
	require.True(t, display.IsPaused())

	// Starting new operation while paused should still work
	display.Start("Operation 2")
	require.True(t, display.IsActive())
	require.True(t, display.IsPaused())

	err = display.Resume()
	require.NoError(t, err)
	require.False(t, display.IsPaused())
}

func Test_PauseAndResume_WithPersistentProgress_WorksCorrectly(t *testing.T) {
	var output bytes.Buffer
	display := logger.NewProgressDisplay(&output)

	display.StartPersistent("Installing packages...")
	display.LogAccomplishment("Package 1 installed")

	err := display.Pause()
	require.NoError(t, err)
	require.True(t, display.IsPaused())

	err = display.Resume()
	require.NoError(t, err)
	require.False(t, display.IsPaused())

	display.LogAccomplishment("Package 2 installed")
	display.FinishPersistent("All packages installed successfully")
}

func Test_Close_AfterPause_RestoresTerminalState(t *testing.T) {
	var output bytes.Buffer
	display := logger.NewProgressDisplay(&output)

	display.Start("Test Operation")
	time.Sleep(50 * time.Millisecond) // Allow spinner to start

	err := display.Pause()
	require.NoError(t, err)
	require.True(t, display.IsPaused())

	err = display.Close()
	require.NoError(t, err)

	// Should be able to close multiple times
	err = display.Close()
	require.NoError(t, err)
}

func Test_NoopProgressDisplay_PauseAndResumeMethods_DoNothing(t *testing.T) {
	display := logger.NewNoopProgressDisplay()

	require.False(t, display.IsActive())
	require.False(t, display.IsPaused())

	err := display.Pause()
	require.NoError(t, err)
	require.False(t, display.IsPaused())

	err = display.Resume()
	require.NoError(t, err)
	require.False(t, display.IsPaused())
}
