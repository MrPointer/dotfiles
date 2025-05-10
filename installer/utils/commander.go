package utils

import (
	"fmt"
	"os"
	"os/exec"
)

// Commander defines an interface for running system commands
// This allows us to proxy commands into a container for testing
// or use the real exec.Command in production.
//
//go:generate moq -out commander_mock.go . Commander
type Commander interface {
	// Run executes a command with the given name and arguments
	// It returns an error if the command fails.
	Run(name string, args ...string) error

	// RunWithEnv executes a command with the given environment variables
	// It returns an error if the command fails.
	RunWithEnv(env map[string]string, name string, args ...string) error
}

// DefaultCommander is the production implementation using os/exec
// It implements the Commander interface.
type DefaultCommander struct{}

func NewDefaultCommander() *DefaultCommander {
	return &DefaultCommander{}
}

var _ Commander = (*DefaultCommander)(nil)

func (c *DefaultCommander) Run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (c *DefaultCommander) RunWithEnv(env map[string]string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Set the environment variables for the command
	for key, value := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}

	return cmd.Run()
}
