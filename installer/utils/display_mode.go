package utils

// DisplayMode represents different output display modes for external tool execution.
type DisplayMode int

const (
	// DisplayModeProgress shows progress indicators and hides command output (default interactive mode).
	DisplayModeProgress DisplayMode = iota
	// DisplayModePlain shows simple progress messages without spinners, hides command output.
	DisplayModePlain
	// DisplayModePassthrough shows all command output directly to stdout.
	DisplayModePassthrough
)

// String returns the string representation of the display mode.
func (d DisplayMode) String() string {
	switch d {
	case DisplayModeProgress:
		return "progress"
	case DisplayModePlain:
		return "plain"
	case DisplayModePassthrough:
		return "passthrough"
	default:
		return "unknown"
	}
}

// ShouldDiscardOutput returns true if command output should be discarded/hidden.
func (d DisplayMode) ShouldDiscardOutput() bool {
	switch d {
	case DisplayModeProgress, DisplayModePlain:
		return true
	case DisplayModePassthrough:
		return false
	default:
		return true
	}
}
