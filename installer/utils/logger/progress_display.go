package logger

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
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
	operationInProgress int32 // atomic counter
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

	// Create indented message for hierarchical display
	indent := strings.Repeat("  ", level)
	displayMessage := indent + message

	// Start spinner in background
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
			go p.runSpinner(ctx, prevOperation, displayMessage)
		}
	} else {
		p.activeSpinner = nil
	}
}

// displayCompletion shows the completion message for an operation.
func (p *ProgressDisplay) displayCompletion(operation *ProgressOperation, success bool, err error) {
	duration := time.Since(operation.StartTime)
	indent := strings.Repeat("  ", operation.Level)

	var displayMessage string
	if success {
		if duration > 100*time.Millisecond {
			displayMessage = fmt.Sprintf("%s%s (took %v)", indent, operation.Message, duration.Round(10*time.Millisecond))
		} else {
			displayMessage = fmt.Sprintf("%s%s", indent, operation.Message)
		}

		// Print success message
		checkmark := lipgloss.NewStyle().Foreground(lipgloss.Color("#2ecc71")).Render("✓")
		fmt.Fprintf(p.output, "\r%s %s\n", checkmark, displayMessage)
	} else {
		if duration > 100*time.Millisecond {
			displayMessage = fmt.Sprintf("%s%s (failed after %v)", indent, operation.Message, duration.Round(10*time.Millisecond))
		} else {
			displayMessage = fmt.Sprintf("%s%s", indent, operation.Message)
		}

		// Print failure message
		cross := lipgloss.NewStyle().Foreground(lipgloss.Color("#e74c3c")).Render("✗")
		errorMsg := fmt.Sprintf("\r%s %s", cross, displayMessage)
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
