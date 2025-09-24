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
	StartProgress(message string)
	// UpdateProgress updates the current progress message
	UpdateProgress(message string)
	// FinishProgress completes the progress with success
	FinishProgress(message string)
	// FailProgress completes the progress with failure and shows error
	FailProgress(message string, err error)
	// StartPersistentProgress starts a progress indicator that shows accomplishments
	StartPersistentProgress(message string)
	// LogAccomplishment logs an accomplishment that stays visible
	LogAccomplishment(message string)
	// FinishPersistentProgress completes persistent progress with success
	FinishPersistentProgress(message string)
	// FailPersistentProgress completes persistent progress with failure
	FailPersistentProgress(message string, err error)
	// StartInteractiveProgress starts progress and pauses spinners for interactive commands
	StartInteractiveProgress(message string)
	// FinishInteractiveProgress completes interactive progress and resumes spinners
	FinishInteractiveProgress(message string)
	// FailInteractiveProgress completes interactive progress with error and resumes spinners
	FailInteractiveProgress(message string, err error)
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
func (l NoopLogger) StartProgress(message string)                      {}
func (l NoopLogger) UpdateProgress(message string)                     {}
func (l NoopLogger) FinishProgress(message string)                     {}
func (l NoopLogger) FailProgress(message string, err error)            {}
func (l NoopLogger) StartPersistentProgress(message string)            {}
func (l NoopLogger) LogAccomplishment(message string)                  {}
func (l NoopLogger) FinishPersistentProgress(message string)           {}
func (l NoopLogger) FailPersistentProgress(message string, err error)  {}
func (l NoopLogger) StartInteractiveProgress(message string)           {}
func (l NoopLogger) FinishInteractiveProgress(message string)          {}
func (l NoopLogger) FailInteractiveProgress(message string, err error) {}
func (l NoopLogger) Close() error                                      { return nil }

// DefaultLogger is the default logger used if none is provided.
var DefaultLogger Logger = NoopLogger{}
