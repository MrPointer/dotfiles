package logger

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
)

// ANSI escape codes for terminal control
const (
	clearLine = "\033[K" // Clear line from cursor to end
)

// ProgressOperation represents an active progress operation.
type ProgressOperation struct {
	Message    string
	StartTime  time.Time
	Level      int
	done       int32 // atomic
	Success    bool
	Error      error
	CancelFunc context.CancelFunc
}

// IsDone returns whether this operation is completed.
func (op *ProgressOperation) IsDone() bool {
	return atomic.LoadInt32(&op.done) == 1
}

// SetDone marks this operation as completed.
func (op *ProgressOperation) SetDone() {
	atomic.StoreInt32(&op.done, 1)
}

// ProgressDisplay provides hierarchical progress reporting with npm-style output.
type ProgressDisplay struct {
	output              io.Writer
	progressStack       []*ProgressOperation
	activeSpinner       *ProgressOperation
	operationInProgress int32          // atomic counter
	persistentMode      bool           // whether we're in persistent mode
	paused              int32          // atomic flag for paused state
	pauseMutex          sync.Mutex     // protects pause/resume operations
	spinnerWaitGroup    sync.WaitGroup // tracks active spinner goroutines
}

var _ ProgressReporter = (*ProgressDisplay)(nil)

// NewProgressDisplay creates a new hierarchical progress display.
func NewProgressDisplay(output io.Writer) *ProgressDisplay {
	if output == nil {
		output = os.Stdout
	}
	return &ProgressDisplay{
		output: output,
	}
}

// Start begins a new progress operation with the given message.
func (p *ProgressDisplay) Start(message string) {
	level := len(p.progressStack)

	// Stop any currently active spinner
	if p.activeSpinner != nil && p.activeSpinner.CancelFunc != nil {
		p.activeSpinner.CancelFunc()
	}

	// Create context for this operation
	ctx, cancel := context.WithCancel(context.Background())

	operation := &ProgressOperation{
		Message:    message,
		StartTime:  time.Now(),
		Level:      level,
		CancelFunc: cancel,
	}
	p.progressStack = append(p.progressStack, operation)
	p.activeSpinner = operation

	// Increment operation counter
	atomic.AddInt32(&p.operationInProgress, 1)

	// Create contextual message showing hierarchy
	displayMessage := p.buildContextualMessage()

	// Start spinner in background
	p.spinnerWaitGroup.Add(1)
	go p.runSpinner(ctx, operation, displayMessage)
}

// Update modifies the message of the current progress operation.
func (p *ProgressDisplay) Update(message string) {
	if len(p.progressStack) == 0 {
		return
	}

	// Update the most recent progress operation
	currentIndex := len(p.progressStack) - 1
	p.progressStack[currentIndex].Message = message
}

// Finish completes the current progress operation successfully.
func (p *ProgressDisplay) Finish(message string) {
	if len(p.progressStack) == 0 {
		return
	}

	// Pop from progress stack
	currentIndex := len(p.progressStack) - 1
	operation := p.progressStack[currentIndex]
	p.progressStack = p.progressStack[:currentIndex]

	// Stop the spinner for this operation
	if operation.CancelFunc != nil {
		operation.CancelFunc()
	}
	operation.SetDone()
	operation.Success = true

	// Decrement operation counter
	atomic.AddInt32(&p.operationInProgress, -1)

	// Resume parent operation if exists
	p.resumeParentOperation()

	// Display completion message
	p.displayCompletion(operation, true, nil)
}

// Fail completes the current progress operation with an error.
func (p *ProgressDisplay) Fail(message string, err error) {
	if len(p.progressStack) == 0 {
		return
	}

	// Pop from progress stack
	currentIndex := len(p.progressStack) - 1
	operation := p.progressStack[currentIndex]
	p.progressStack = p.progressStack[:currentIndex]

	// Stop the spinner for this operation
	if operation.CancelFunc != nil {
		operation.CancelFunc()
	}
	operation.SetDone()
	operation.Success = false
	operation.Error = err

	// Decrement operation counter
	atomic.AddInt32(&p.operationInProgress, -1)

	// Resume parent operation if exists
	p.resumeParentOperation()

	// Display failure message
	p.displayCompletion(operation, false, err)
}

// IsActive returns true if there are any active progress operations.
func (p *ProgressDisplay) IsActive() bool {
	return atomic.LoadInt32(&p.operationInProgress) > 0
}

// Clear stops all progress operations without displaying completion messages.
func (p *ProgressDisplay) Clear() {
	// Stop all active spinners
	for _, operation := range p.progressStack {
		if operation.CancelFunc != nil {
			operation.CancelFunc()
		}
	}

	// Clear the stack and reset counter
	p.progressStack = nil
	p.activeSpinner = nil
	atomic.StoreInt32(&p.operationInProgress, 0)
	atomic.StoreInt32(&p.paused, 0)
}

// resumeParentOperation resumes the spinner for the parent operation if one exists.
func (p *ProgressDisplay) resumeParentOperation() {
	if len(p.progressStack) > 0 {
		prevOperation := p.progressStack[len(p.progressStack)-1]
		if !prevOperation.IsDone() {
			p.activeSpinner = prevOperation

			// Create new context for resumed operation
			ctx, cancel := context.WithCancel(context.Background())
			prevOperation.CancelFunc = cancel

			// Resume spinner for previous operation
			indent := strings.Repeat("  ", prevOperation.Level)
			displayMessage := indent + prevOperation.Message
			p.spinnerWaitGroup.Add(1)
			go p.runSpinner(ctx, prevOperation, displayMessage)
		}
	} else {
		p.activeSpinner = nil
	}
}

// displayCompletion shows the completion message for an operation.
func (p *ProgressDisplay) displayCompletion(operation *ProgressOperation, success bool, err error) {
	duration := time.Since(operation.StartTime)

	// In persistent mode, don't show individual completion messages
	// unless it's the top-level operation
	if p.persistentMode && operation.Level > 0 {
		return
	}

	var displayMessage string
	if success {
		if duration > 100*time.Millisecond {
			displayMessage = fmt.Sprintf("%s (took %v)", operation.Message, duration.Round(10*time.Millisecond))
		} else {
			displayMessage = operation.Message
		}

		// Print success message without indentation for minimal output
		checkmark := lipgloss.NewStyle().Foreground(lipgloss.Color("#2ecc71")).Render("✓")
		fmt.Fprintf(p.output, "\r%s%s %s\n", clearLine, checkmark, displayMessage)
	} else {
		if duration > 100*time.Millisecond {
			displayMessage = fmt.Sprintf("%s (failed after %v)", operation.Message, duration.Round(10*time.Millisecond))
		} else {
			displayMessage = operation.Message
		}

		// Print failure message without indentation for minimal output
		cross := lipgloss.NewStyle().Foreground(lipgloss.Color("#e74c3c")).Render("✗")
		errorMsg := fmt.Sprintf("\r%s%s %s", clearLine, cross, displayMessage)
		if err != nil {
			errorMsg += fmt.Sprintf("\n  Error: %v", err)
		}

		// Write to stderr for errors, but use the configured output writer
		if p.output == os.Stdout {
			fmt.Fprintf(os.Stderr, "%s\n", errorMsg)
		} else {
			fmt.Fprintf(p.output, "%s\n", errorMsg)
		}
	}
}

// runSpinner runs a spinner for the given operation in the background.
func (p *ProgressDisplay) runSpinner(ctx context.Context, operation *ProgressOperation, displayMessage string) {
	// Signal completion when function exits
	defer p.spinnerWaitGroup.Done()

	// Don't start spinner if paused
	if p.IsPaused() {
		return
	}

	// Create spinner with huh
	s := spinner.New().
		Title(displayMessage).
		Type(spinner.Dots).
		Output(p.output).
		Accessible(false).
		Context(ctx)

	// Run spinner with a simple action that waits for completion
	s.ActionWithErr(func(spinnerCtx context.Context) error {
		ticker := time.NewTicker(50 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-spinnerCtx.Done():
				return spinnerCtx.Err()
			case <-ticker.C:
				// Stop spinner if paused
				if p.IsPaused() {
					return nil
				}
				if operation == nil {
					return errors.New("no operation in progress")
				}
				if operation.IsDone() {
					return nil
				}
			}
		}
	})

	// Run the spinner
	s.Run()
}

// Pause temporarily stops all spinner operations for interactive commands
func (p *ProgressDisplay) Pause() {
	p.pauseMutex.Lock()
	defer p.pauseMutex.Unlock()

	// Set paused state first
	atomic.StoreInt32(&p.paused, 1)

	// Cancel all active spinners
	for _, operation := range p.progressStack {
		if operation.CancelFunc != nil {
			operation.CancelFunc()
		}
	}

	// Wait for all spinner goroutines to finish
	p.spinnerWaitGroup.Wait()
}

// Resume restarts spinner operations after interactive commands complete
func (p *ProgressDisplay) Resume() {
	p.pauseMutex.Lock()
	defer p.pauseMutex.Unlock()

	atomic.StoreInt32(&p.paused, 0)

	// Resume the most recent operation if there is one
	if len(p.progressStack) > 0 {
		currentOperation := p.progressStack[len(p.progressStack)-1]
		if !currentOperation.IsDone() {
			// Create new context for resumed operation
			ctx, cancel := context.WithCancel(context.Background())
			currentOperation.CancelFunc = cancel

			// Resume spinner for current operation
			displayMessage := p.buildContextualMessage()
			p.spinnerWaitGroup.Add(1)
			go p.runSpinner(ctx, currentOperation, displayMessage)
		}
	}
}

// IsPaused returns whether the progress display is currently paused
func (p *ProgressDisplay) IsPaused() bool {
	return atomic.LoadInt32(&p.paused) == 1
}

// buildContextualMessage creates a hierarchical message showing the full context
func (p *ProgressDisplay) buildContextualMessage() string {
	if len(p.progressStack) == 0 {
		return ""
	}

	// Build context from all operations in the stack
	var parts []string
	for _, op := range p.progressStack {
		parts = append(parts, op.Message)
	}

	// Join with separator to show hierarchy
	return strings.Join(parts, ": ")
}

// StartPersistent begins a persistent progress operation that shows accomplishments.
func (p *ProgressDisplay) StartPersistent(message string) {
	p.persistentMode = true
	p.Start(message)
}

// LogAccomplishment logs an accomplishment that stays visible.
func (p *ProgressDisplay) LogAccomplishment(message string) {
	checkmark := lipgloss.NewStyle().Foreground(lipgloss.Color("#2ecc71")).Render("✓")
	fmt.Fprintf(p.output, "\r%s   %s %s\n", clearLine, checkmark, message)
}

// FinishPersistent completes persistent progress with success.
func (p *ProgressDisplay) FinishPersistent(message string) {
	p.persistentMode = false
	p.Finish(message)
}

// FailPersistent completes persistent progress with failure.
func (p *ProgressDisplay) FailPersistent(message string, err error) {
	p.persistentMode = false
	p.Fail(message, err)
}

// ProgressReporter defines the interface for hierarchical progress reporting.
type ProgressReporter interface {
	// Start begins a new progress operation with the given message
	Start(message string)
	// Update modifies the message of the current progress operation
	Update(message string)
	// Finish completes the current progress operation successfully
	Finish(message string)
	// Fail completes the current progress operation with an error
	Fail(message string, err error)
	// IsActive returns true if there are any active progress operations
	IsActive() bool
	// Clear stops all progress operations without displaying completion messages
	Clear()
	// Pause temporarily stops all spinner operations for interactive commands
	Pause()
	// Resume restarts spinner operations after interactive commands complete
	Resume()
	// IsPaused returns whether the progress display is currently paused
	IsPaused() bool
	// StartPersistent begins a persistent progress operation that shows accomplishments
	StartPersistent(message string)
	// LogAccomplishment logs an accomplishment that stays visible
	LogAccomplishment(message string)
	// FinishPersistent completes persistent progress with success
	FinishPersistent(message string)
	// FailPersistent completes persistent progress with failure
	FailPersistent(message string, err error)
}

// NoopProgressDisplay is a progress display that does nothing.
type NoopProgressDisplay struct{}

var _ ProgressReporter = (*NoopProgressDisplay)(nil)

// NewNoopProgressDisplay creates a progress display that does nothing.
func NewNoopProgressDisplay() *NoopProgressDisplay {
	return &NoopProgressDisplay{}
}

// Start does nothing.
func (n *NoopProgressDisplay) Start(message string) {}

// Update does nothing.
func (n *NoopProgressDisplay) Update(message string) {}

// Finish does nothing.
func (n *NoopProgressDisplay) Finish(message string) {}

// Fail does nothing.
func (n *NoopProgressDisplay) Fail(message string, err error) {}

// IsActive always returns false.
func (n *NoopProgressDisplay) IsActive() bool { return false }

// Clear does nothing.
func (n *NoopProgressDisplay) Clear() {}

// Pause does nothing.
func (n *NoopProgressDisplay) Pause() {}

// Resume does nothing.
func (n *NoopProgressDisplay) Resume() {}

// IsPaused always returns false.
func (n *NoopProgressDisplay) IsPaused() bool { return false }

// StartPersistent does nothing.
func (n *NoopProgressDisplay) StartPersistent(message string) {}

// LogAccomplishment does nothing.
func (n *NoopProgressDisplay) LogAccomplishment(message string) {}

// FinishPersistent does nothing.
func (n *NoopProgressDisplay) FinishPersistent(message string) {}

// FailPersistent does nothing.
func (n *NoopProgressDisplay) FailPersistent(message string, err error) {}
