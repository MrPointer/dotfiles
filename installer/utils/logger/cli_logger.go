package logger

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
)

// Styles for different types of messages using lipgloss
var (
	infoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#3498db")).Bold(true) // Blue
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#2ecc71")).Bold(true) // Green
	warningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#f39c12")).Bold(true) // Yellow/Orange
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#e74c3c")).Bold(true) // Red
)

// CliLogger implements the Logger interface using lipgloss styling
type CliLogger struct{}

// NewCliLogger creates a new CLI logger that uses lipgloss styling
func NewCliLogger() *CliLogger {
	return &CliLogger{}
}

// Info logs an informational message with blue styling
func (l *CliLogger) Info(format string, args ...any) {
	printStyled(os.Stdout, infoStyle, format, args...)
}

// Success logs a success message with green styling
func (l *CliLogger) Success(format string, args ...any) {
	printStyled(os.Stdout, successStyle, format, args...)
}

// Warning logs a warning message with yellow styling
func (l *CliLogger) Warning(format string, args ...any) {
	printStyled(os.Stdout, warningStyle, format, args...)
}

// Error logs an error message with red styling
func (l *CliLogger) Error(format string, args ...any) {
	printStyled(os.Stderr, errorStyle, format, args...)
}

// Helper function to print styled text to the specified writer
func printStyled(writer *os.File, style lipgloss.Style, format string, args ...any) {
	fmt.Fprintln(writer, style.Render(fmt.Sprintf(format, args...)))
}
