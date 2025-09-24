package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/MrPointer/dotfiles/installer/utils/logger"
)

// Commander defines an interface for running system commands
// This allows us to proxy commands into a container for testing
// or use the real exec.Command in production.
type Commander interface {
	// RunCommand executes a command with flexible options
	// It returns the result containing output, error information, and exit code
	RunCommand(name string, args []string, opts ...Option) (*Result, error)
}

// Result contains the output and metadata from a command execution
type Result struct {
	// Stdout contains the standard output
	Stdout []byte
	// Stderr contains the standard error output
	Stderr []byte
	// ExitCode is the exit code of the command
	ExitCode int
	// Duration is how long the command took to execute
	Duration time.Duration
}

// String returns the stdout as a string
func (r *Result) String() string {
	return string(r.Stdout)
}

// StderrString returns the stderr as a string
func (r *Result) StderrString() string {
	return string(r.Stderr)
}

// Options contains all configurable options for command execution
type Options struct {
	// Env contains environment variables to set for the command
	Env map[string]string
	// Dir is the working directory for the command
	Dir string
	// Input is data to send to the command's stdin
	Input []byte
	// CaptureOutput determines whether to capture stdout/stderr or pipe to current process
	CaptureOutput bool
	// Interactive determines whether this is an interactive command that needs direct terminal access
	Interactive bool
	// Timeout specifies a timeout for the command execution
	Timeout time.Duration
	// Stdout specifies where to write stdout (only used when CaptureOutput is false)
	Stdout io.Writer
	// Stderr specifies where to write stderr (only used when CaptureOutput is false)
	Stderr io.Writer
}

// Option is a function that modifies Options
type Option func(*Options)

// WithEnv sets environment variables for the command
func WithEnv(env map[string]string) Option {
	return func(o *Options) {
		if o.Env == nil {
			o.Env = make(map[string]string)
		}
		for k, v := range env {
			o.Env[k] = v
		}
	}
}

// WithEnvVar sets a single environment variable
func WithEnvVar(key, value string) Option {
	return func(o *Options) {
		if o.Env == nil {
			o.Env = make(map[string]string)
		}
		o.Env[key] = value
	}
}

// WithDir sets the working directory for the command
func WithDir(dir string) Option {
	return func(o *Options) {
		o.Dir = dir
	}
}

// WithInput provides input to send to the command's stdin
func WithInput(input []byte) Option {
	return func(o *Options) {
		o.Input = input
	}
}

// WithInputString provides string input to send to the command's stdin
func WithInputString(input string) Option {
	return func(o *Options) {
		o.Input = []byte(input)
	}
}

// WithCaptureOutput enables capturing stdout and stderr in the result
func WithCaptureOutput() Option {
	return func(o *Options) {
		o.CaptureOutput = true
	}
}

// WithInteractive enables interactive mode for commands that need user input
// This ensures stdin/stdout/stderr are connected to the terminal and not captured
func WithInteractive() Option {
	return func(o *Options) {
		o.Interactive = true
		o.CaptureOutput = false // Interactive commands should not capture output
	}
}

// WithInteractiveCapture enables interactive mode while still capturing output
// This allows user interaction but also captures output for parsing
func WithInteractiveCapture() Option {
	return func(o *Options) {
		o.Interactive = true
		o.CaptureOutput = true // Capture output for parsing while allowing interaction
	}
}

// WithTimeout sets a timeout for command execution
func WithTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.Timeout = timeout
	}
}

// WithStdout sets where stdout should be written (when not capturing)
func WithStdout(w io.Writer) Option {
	return func(o *Options) {
		o.Stdout = w
	}
}

// WithStderr sets where stderr should be written (when not capturing)
func WithStderr(w io.Writer) Option {
	return func(o *Options) {
		o.Stderr = w
	}
}

// DefaultCommander is the production implementation using os/exec
type DefaultCommander struct {
	logger logger.Logger
}

func NewDefaultCommander(logger logger.Logger) *DefaultCommander {
	return &DefaultCommander{
		logger: logger,
	}
}

var _ Commander = (*DefaultCommander)(nil)

func (c *DefaultCommander) RunCommand(name string, args []string, opts ...Option) (*Result, error) {
	c.logger.Trace("Running command: %s %s", name, strings.Join(args, " "))

	// Apply default options
	options := &Options{
		CaptureOutput: false,
		Interactive:   false,
		Stdout:        os.Stdout,
		Stderr:        os.Stderr,
	}

	// Apply provided options
	for _, opt := range opts {
		opt(options)
	}

	// Create the command
	cmd := exec.Command(name, args...)

	// Set working directory
	if options.Dir != "" {
		cmd.Dir = options.Dir
	}

	// Set environment variables
	if len(options.Env) > 0 {
		cmd.Env = os.Environ() // Start with current environment
		for key, value := range options.Env {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
		}
	}

	// Set up input
	if options.Interactive {
		// For interactive commands, connect stdin directly to terminal
		cmd.Stdin = os.Stdin
	} else if len(options.Input) > 0 {
		cmd.Stdin = bytes.NewReader(options.Input)
	}

	var stdout, stderr bytes.Buffer
	var result *Result

	start := time.Now()

	if options.Interactive && !options.CaptureOutput {
		// For pure interactive commands, connect directly to terminal
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else if options.Interactive && options.CaptureOutput {
		// For interactive commands that also need output capture,
		// use io.MultiWriter to both display and capture
		cmd.Stdout = io.MultiWriter(os.Stdout, &stdout)
		cmd.Stderr = io.MultiWriter(os.Stderr, &stderr)
	} else if options.CaptureOutput {
		// Capture output in buffers
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
	} else {
		// Pipe to specified writers
		cmd.Stdout = options.Stdout
		cmd.Stderr = options.Stderr
	}

	// Handle timeout
	var err error
	if options.Timeout > 0 {
		done := make(chan error, 1)
		go func() {
			done <- cmd.Run()
		}()

		select {
		case err = <-done:
			// Command completed normally
		case <-time.After(options.Timeout):
			// Timeout occurred
			if cmd.Process != nil {
				cmd.Process.Kill()
			}
			err = fmt.Errorf("command timed out after %v", options.Timeout)
		}
	} else {
		err = cmd.Run()
	}

	duration := time.Since(start)

	// Create result
	result = &Result{
		Duration: duration,
	}

	if options.CaptureOutput {
		result.Stdout = stdout.Bytes()
		result.Stderr = stderr.Bytes()
	}

	// Get exit code
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitError.ExitCode()
		} else {
			// Non-exit error (e.g., command not found, timeout)
			result.ExitCode = -1
		}
	} else {
		result.ExitCode = 0
	}

	return result, err
}
