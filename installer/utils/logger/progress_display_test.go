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

	// Should have indented child operation
	require.Contains(t, output, "  Child operation")
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

	// Check indentation levels
	require.Contains(t, output, "Level 1")       // No indentation
	require.Contains(t, output, "  Level 2")     // 2 spaces
	require.Contains(t, output, "    Level 3")   // 4 spaces
	require.Contains(t, output, "      Level 4") // 6 spaces
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
