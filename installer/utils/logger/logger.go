package logger

// Logger defines a minimal logging interface that our installer utilities need
type Logger interface {
	// Info logs an informational message
	Info(format string, args ...any)
	// Success logs a success message
	Success(format string, args ...any)
	// Warning logs a warning message
	Warning(format string, args ...any)
	// Error logs an error message
	Error(format string, args ...any)
}

// NoopLogger implements Logger but does nothing
type NoopLogger struct{}

func (l NoopLogger) Info(format string, args ...any)    {}
func (l NoopLogger) Success(format string, args ...any) {}
func (l NoopLogger) Warning(format string, args ...any) {}
func (l NoopLogger) Error(format string, args ...any)   {}

// DefaultLogger is the default logger used if none is provided
var DefaultLogger Logger = NoopLogger{}
