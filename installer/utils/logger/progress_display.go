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
	clearLine  = "\033[K"    // Clear line from cursor to end
	showCursor = "\033[?25h" // Show cursor
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

// ProgressReporter defines the interface for hierarchical progress reporting.
type ProgressReporter interface {
	io.Closer
	// Start begins a new progress operation with the given message
	Start(message string) error
	// Update modifies the message of the current progress operation
	Update(message string) error
	// Finish completes the current progress operation successfully
	Finish(message string) error
	// Fail completes the current progress operation with an error
	Fail(message string, err error) error

	// StartPersistent begins a persistent progress operation that shows accomplishments
	StartPersistent(message string) error
	// LogAccomplishment logs an accomplishment that stays visible
	LogAccomplishment(message string) error
	// FinishPersistent completes persistent progress with success
	FinishPersistent(message string) error
	// FailPersistent completes persistent progress with failure
	FailPersistent(message string, err error) error

	// Clear stops all progress operations without displaying completion messages
	Clear() error
	// Pause temporarily stops all spinner operations for interactive commands
	Pause() error
	// Resume restarts spinner operations after interactive commands complete
	Resume() error

	// IsActive returns true if there are any active progress operations
	IsActive() bool
	// IsPaused returns whether the progress display is currently paused
	IsPaused() bool
}

// ProgressDisplay provides hierarchical progress reporting with npm-style output.
type ProgressDisplay struct {
	output              io.Writer
	progressStack       []*ProgressOperation
	activeSpinner       *ProgressOperation
	operationInProgress int32          // atomic counter
	persistentMode      bool           // whether we're in persistent mode
	cursorHidden        int32          // atomic flag for cursor state
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
	pd := &ProgressDisplay{
		output: output,
	}

	// Ensure cursor is restored on program exit
	pd.setupCleanup()

	return pd
}

// Start begins a new progress operation with the given message.
func (p *ProgressDisplay) Start(message string) error {
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

	return nil
}

// Update modifies the message of the current progress operation.
func (p *ProgressDisplay) Update(message string) error {
	if len(p.progressStack) == 0 {
		// Updating an inexistent operation is not an error
		return nil
	}

	// Update the most recent progress operation
	currentIndex := len(p.progressStack) - 1
	p.progressStack[currentIndex].Message = message

	return nil
}

// Finish completes the current progress operation successfully.
func (p *ProgressDisplay) Finish(message string) error {
	if len(p.progressStack) == 0 {
		return nil
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

	// Wait for spinner goroutine to complete cleanup before proceeding
	p.spinnerWaitGroup.Wait()

	// Ensure cursor is restored before displaying completion message
	if err := p.restoreCursor(); err != nil {
		return err
	}

	// Decrement operation counter
	atomic.AddInt32(&p.operationInProgress, -1)

	// Display completion message
	if err := p.displayCompletion(operation, true, nil); err != nil {
		return err
	}

	// Resume parent operation if exists
	p.resumeParentOperation()

	return nil
}

// Fail completes the current progress operation with an error.
func (p *ProgressDisplay) Fail(message string, err error) error {
	if len(p.progressStack) == 0 {
		return nil
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

	// Wait for spinner goroutine to complete cleanup before proceeding
	p.spinnerWaitGroup.Wait()

	// Ensure cursor is restored before displaying completion message
	if err := p.restoreCursor(); err != nil {
		return err
	}

	// Decrement operation counter
	atomic.AddInt32(&p.operationInProgress, -1)

	// Display failure message
	if err := p.displayCompletion(operation, false, err); err != nil {
		return err
	}

	// Resume parent operation if exists
	p.resumeParentOperation()

	return nil
}

// IsActive returns true if there are any active progress operations.
func (p *ProgressDisplay) IsActive() bool {
	return atomic.LoadInt32(&p.operationInProgress) > 0
}

// Clear stops all progress operations without displaying completion messages.
func (p *ProgressDisplay) Clear() error {
	// Stop all active spinners
	for _, operation := range p.progressStack {
		if operation.CancelFunc != nil {
			operation.CancelFunc()
		}
		operation.SetDone()
	}

	// Wait for all spinner goroutines to complete
	p.spinnerWaitGroup.Wait()

	// Clear the stack and reset counter
	p.progressStack = nil
	p.activeSpinner = nil
	atomic.StoreInt32(&p.operationInProgress, 0)
	atomic.StoreInt32(&p.paused, 0)

	// Restore cursor if it was hidden
	return p.restoreCursor()
}

// Pause temporarily stops all spinner operations for interactive commands
func (p *ProgressDisplay) Pause() error {
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

	// Now it's safe to clean up terminal state
	if err := p.restoreCursor(); err != nil {
		return err
	}

	if file, ok := p.output.(*os.File); ok {
		_, err := file.WriteString("\r" + clearLine)
		if err != nil {
			return err
		}
	} else {
		_, err := fmt.Fprint(p.output, "\r"+clearLine)
		if err != nil {
			return err
		}
	}

	return nil
}

// Resume restarts spinner operations after interactive commands complete
func (p *ProgressDisplay) Resume() error {
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

	return nil
}

// IsPaused returns whether the progress display is currently paused
func (p *ProgressDisplay) IsPaused() bool {
	return atomic.LoadInt32(&p.paused) == 1
}

// StartPersistent begins a persistent progress operation that shows accomplishments.
func (p *ProgressDisplay) StartPersistent(message string) error {
	p.persistentMode = true
	return p.Start(message)
}

// LogAccomplishment logs an accomplishment that stays visible.
func (p *ProgressDisplay) LogAccomplishment(message string) error {
	checkmark := lipgloss.NewStyle().Foreground(lipgloss.Color("#2ecc71")).Render("✓")
	_, err := fmt.Fprintf(p.output, "\r%s   %s %s\n", clearLine, checkmark, message)
	return err
}

// FinishPersistent completes persistent progress with success.
func (p *ProgressDisplay) FinishPersistent(message string) error {
	p.persistentMode = false
	return p.Finish(message)
}

// FailPersistent completes persistent progress with failure.
func (p *ProgressDisplay) FailPersistent(message string, err error) error {
	p.persistentMode = false
	return p.Fail(message, err)
}

// Close ensures proper cleanup of terminal state.
func (p *ProgressDisplay) Close() error {
	return p.Clear()
}

// setupCleanup sets up signal handlers and cleanup mechanisms to ensure cursor is restored.
func (p *ProgressDisplay) setupCleanup() {
	// Note: We don't set up signal handlers here to avoid dependencies.
	// The cursor restoration will happen when Clear() is called or through defer.
}

// resumeParentOperation resumes the spinner for the parent operation if one exists.
func (p *ProgressDisplay) resumeParentOperation() {
	if len(p.progressStack) == 0 {
		p.activeSpinner = nil
		return
	}

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
}

// displayCompletion shows the completion message for an operation.
func (p *ProgressDisplay) displayCompletion(operation *ProgressOperation, success bool, err error) error {
	duration := time.Since(operation.StartTime)

	// In persistent mode, don't show individual completion messages
	// unless it's the top-level operation
	if p.persistentMode && operation.Level > 0 {
		return nil
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
			_, err := fmt.Fprintf(os.Stderr, "%s\n", errorMsg)
			if err != nil {
				return err
			}
		} else {
			_, err := fmt.Fprintf(p.output, "%s\n", errorMsg)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// runSpinner runs a spinner for the given operation in the background.
func (p *ProgressDisplay) runSpinner(ctx context.Context, operation *ProgressOperation, displayMessage string) {
	// Signal completion when function exits
	defer p.spinnerWaitGroup.Done()

	// Don't start spinner if paused
	if p.IsPaused() {
		return
	}

	// Mark cursor as hidden when spinner starts
	atomic.StoreInt32(&p.cursorHidden, 1)

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

// restoreCursor ensures the terminal cursor is visible.
func (p *ProgressDisplay) restoreCursor() error {
	if !atomic.CompareAndSwapInt32(&p.cursorHidden, 1, 0) {
		return errors.New("Failed swapping atomic value representing cursor visibility")
	}

	if file, ok := p.output.(*os.File); ok {
		file.WriteString(showCursor)
	} else {
		fmt.Fprint(p.output, showCursor)
	}

	return nil
}

// NoopProgressDisplay is a progress display that does nothing.
type NoopProgressDisplay struct{}

var _ ProgressReporter = (*NoopProgressDisplay)(nil)

// NewNoopProgressDisplay creates a progress display that does nothing.
func NewNoopProgressDisplay() *NoopProgressDisplay {
	return &NoopProgressDisplay{}
}

// Start does nothing.
func (n *NoopProgressDisplay) Start(message string) error {
	return nil
}

// Update does nothing.
func (n *NoopProgressDisplay) Update(message string) error {
	return nil
}

// Finish does nothing.
func (n *NoopProgressDisplay) Finish(message string) error {
	return nil
}

// Fail does nothing.
func (n *NoopProgressDisplay) Fail(message string, err error) error {
	return nil
}

// IsActive always returns false.
func (n *NoopProgressDisplay) IsActive() bool { return false }

// Clear does nothing.
func (n *NoopProgressDisplay) Clear() error { return nil }

// Pause does nothing.
func (n *NoopProgressDisplay) Pause() error {
	return nil
}

// Resume does nothing.
func (n *NoopProgressDisplay) Resume() error {
	return nil
}

// IsPaused always returns false.
func (n *NoopProgressDisplay) IsPaused() bool { return false }

// StartPersistent does nothing.
func (n *NoopProgressDisplay) StartPersistent(message string) error {
	return nil
}

// LogAccomplishment does nothing.
func (n *NoopProgressDisplay) LogAccomplishment(message string) error {
	return nil
}

// FinishPersistent does nothing.
func (n *NoopProgressDisplay) FinishPersistent(message string) error {
	return nil
}

// FailPersistent does nothing.
func (n *NoopProgressDisplay) FailPersistent(message string, err error) error {
	return nil
}

// Close does nothing.
func (n *NoopProgressDisplay) Close() error { return nil }
