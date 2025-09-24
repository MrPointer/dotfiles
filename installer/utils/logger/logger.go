package logger

import "io"

// Logger defines a minimal logging interface that our installer utilities need.
type Logger interface {
	io.Closer

	// Trace logs a trace message
	Trace(format string, args ...any)
	// Debug logs a debug message
	Debug(format string, args ...any)
	// Info logs an informational message
	Info(format string, args ...any)
	// Success logs a success message
	Success(format string, args ...any)
	// Warning logs a warning message
	Warning(format string, args ...any)
	// Error logs an error message
	Error(format string, args ...any)

	// StartProgress starts a progress indicator with the given message
	StartProgress(message string) error
	// UpdateProgress updates the current progress message
	UpdateProgress(message string) error
	// FinishProgress completes the progress with success
	FinishProgress(message string) error
	// FailProgress completes the progress with failure and shows error
	FailProgress(message string, err error) error

	// StartPersistentProgress starts a progress indicator that shows accomplishments
	StartPersistentProgress(message string) error
	// LogAccomplishment logs an accomplishment that stays visible
	LogAccomplishment(message string) error
	// FinishPersistentProgress completes persistent progress with success
	FinishPersistentProgress(message string) error
	// FailPersistentProgress completes persistent progress with failure
	FailPersistentProgress(message string, err error) error

	// StartInteractiveProgress starts progress and pauses spinners for interactive commands
	StartInteractiveProgress(message string) error
	// FinishInteractiveProgress completes interactive progress and resumes spinners
	FinishInteractiveProgress(message string) error
	// FailInteractiveProgress completes interactive progress with error and resumes spinners
	FailInteractiveProgress(message string, err error) error
}

// NoopLogger implements Logger but does nothing.
type NoopLogger struct{}

var _ Logger = (*NoopLogger)(nil)

func (l NoopLogger) Trace(format string, args ...any)   {}
func (l NoopLogger) Debug(format string, args ...any)   {}
func (l NoopLogger) Info(format string, args ...any)    {}
func (l NoopLogger) Success(format string, args ...any) {}
func (l NoopLogger) Warning(format string, args ...any) {}
func (l NoopLogger) Error(format string, args ...any)   {}

// Progress methods - no-op implementations
func (l NoopLogger) StartProgress(message string) error                      { return nil }
func (l NoopLogger) UpdateProgress(message string) error                     { return nil }
func (l NoopLogger) FinishProgress(message string) error                     { return nil }
func (l NoopLogger) FailProgress(message string, err error) error            { return nil }
func (l NoopLogger) StartPersistentProgress(message string) error            { return nil }
func (l NoopLogger) LogAccomplishment(message string) error                  { return nil }
func (l NoopLogger) FinishPersistentProgress(message string) error           { return nil }
func (l NoopLogger) FailPersistentProgress(message string, err error) error  { return nil }
func (l NoopLogger) StartInteractiveProgress(message string) error           { return nil }
func (l NoopLogger) FinishInteractiveProgress(message string) error          { return nil }
func (l NoopLogger) FailInteractiveProgress(message string, err error) error { return nil }
func (l NoopLogger) Close() error                                            { return nil }

// DefaultLogger is the default logger used if none is provided.
var DefaultLogger Logger = NoopLogger{}
