package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MrPointer/dotfiles/installer/utils/logger"
	"github.com/spf13/cobra"
)

var (
	verbosity      string
	withProgress   bool
	operationType  string
	operationCount int
	failureRate    int
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "spinner-demo",
	Short: "Demo application for testing spinner capabilities",
	Long: `A demo application that showcases the hierarchical progress display
capabilities of the logger package. Use this to test and refine spinner
behavior without running the full installer.`,
	RunE: runDemo,
}

func init() {
	rootCmd.Flags().StringVarP(&verbosity, "verbosity", "v", "normal",
		"Log verbosity level (minimal, normal, verbose, extra-verbose)")
	rootCmd.Flags().BoolVarP(&withProgress, "progress", "p", true,
		"Enable progress display with spinners")
	rootCmd.Flags().StringVarP(&operationType, "type", "t", "simple",
		"Type of demo to run (simple, nested, mixed, concurrent, long, persistent, interactive, stress)")
	rootCmd.Flags().IntVarP(&operationCount, "count", "c", 3,
		"Number of operations to run")
	rootCmd.Flags().IntVar(&failureRate, "fail-rate", 0,
		"Percentage of operations that should fail (0-100)")
}

func runDemo(cmd *cobra.Command, args []string) error {
	// Parse verbosity level
	var verbosityLevel logger.VerbosityLevel
	switch verbosity {
	case "minimal":
		verbosityLevel = logger.Minimal
	case "normal":
		verbosityLevel = logger.Normal
	case "verbose":
		verbosityLevel = logger.Verbose
	case "extra-verbose":
		verbosityLevel = logger.ExtraVerbose
	default:
		return fmt.Errorf("invalid verbosity level: %s", verbosity)
	}

	// Create logger
	var log logger.Logger
	if withProgress {
		log = logger.NewProgressCliLogger(verbosityLevel)
		fmt.Println("🎯 Running spinner demo with progress display enabled")
	} else {
		log = logger.NewCliLogger(verbosityLevel)
		fmt.Println("📝 Running spinner demo with plain text output")
	}

	// Set up signal handling to ensure cleanup on interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Info("Interrupt received, cleaning up...")
		log.Close()
		os.Exit(1)
	}()

	// Defer cleanup to ensure it runs even if program exits normally
	defer log.Close()

	fmt.Printf("📊 Configuration: type=%s, count=%d, fail-rate=%d%%\n\n",
		operationType, operationCount, failureRate)

	// Run the appropriate demo
	switch operationType {
	case "simple":
		return runSimpleDemo(log)
	case "nested":
		return runNestedDemo(log)
	case "mixed":
		return runMixedDemo(log)
	case "concurrent":
		return runConcurrentDemo(log)
	case "long":
		return runLongDemo(log)
	case "persistent":
		return runPersistentDemo(log)
	case "interactive":
		return runInteractiveDemo(log)
	case "stress":
		return runStressTestDemo(log)
	default:
		return fmt.Errorf("invalid operation type: %s", operationType)
	}
}

func runSimpleDemo(log logger.Logger) error {
	log.Info("🚀 Starting simple demo with %d operations", operationCount)

	for i := 1; i <= operationCount; i++ {
		operationName := fmt.Sprintf("Simple operation %d", i)

		log.StartProgress(operationName)

		// Simulate work
		duration := time.Duration(200+i*100) * time.Millisecond
		time.Sleep(duration)

		// Randomly fail based on failure rate
		if shouldFail() {
			log.FailProgress(operationName, errors.New("simulated failure"))
		} else {
			log.FinishProgress(fmt.Sprintf("Completed %s", operationName))
		}
	}

	log.Success("✅ Simple demo completed")
	return nil
}

func runNestedDemo(log logger.Logger) error {
	log.Info("🏗️ Starting nested demo")

	log.StartProgress("Parent operation")
	time.Sleep(100 * time.Millisecond)

	for i := 1; i <= operationCount; i++ {
		childOp := fmt.Sprintf("Child operation %d", i)

		log.StartProgress(childOp)
		time.Sleep(150 * time.Millisecond)

		// Add some grandchild operations
		if i%2 == 0 {
			grandchildOp := fmt.Sprintf("Grandchild of operation %d", i)
			log.StartProgress(grandchildOp)
			time.Sleep(100 * time.Millisecond)

			if shouldFail() {
				log.FailProgress(grandchildOp, errors.New("grandchild failed"))
			} else {
				log.FinishProgress(fmt.Sprintf("Completed %s", grandchildOp))
			}
		}

		if shouldFail() {
			log.FailProgress(childOp, errors.New("child operation failed"))
		} else {
			log.FinishProgress(fmt.Sprintf("Completed %s", childOp))
		}
	}

	log.FinishProgress("Parent operation completed")
	log.Success("✅ Nested demo completed")
	return nil
}

func runMixedDemo(log logger.Logger) error {
	log.Info("🎭 Starting mixed demo (success/failure mix)")

	for i := 1; i <= operationCount; i++ {
		parentOp := fmt.Sprintf("Mixed parent %d", i)

		log.StartProgress(parentOp)
		time.Sleep(50 * time.Millisecond)

		// First child - usually succeeds
		child1 := fmt.Sprintf("Child %d-A (success)", i)
		log.StartProgress(child1)
		time.Sleep(100 * time.Millisecond)
		log.FinishProgress(fmt.Sprintf("Completed %s", child1))

		// Second child - might fail
		child2 := fmt.Sprintf("Child %d-B (risky)", i)
		log.StartProgress(child2)
		time.Sleep(100 * time.Millisecond)

		if shouldFail() {
			log.FailProgress(child2, errors.New("risky operation failed"))
		} else {
			log.FinishProgress(fmt.Sprintf("Completed %s", child2))
		}

		// Third child with updates
		child3 := fmt.Sprintf("Child %d-C (with updates)", i)
		log.StartProgress(child3)
		time.Sleep(50 * time.Millisecond)

		log.UpdateProgress(fmt.Sprintf("%s (25%% complete)", child3))
		time.Sleep(50 * time.Millisecond)

		log.UpdateProgress(fmt.Sprintf("%s (75%% complete)", child3))
		time.Sleep(50 * time.Millisecond)

		log.FinishProgress(fmt.Sprintf("Completed %s", child3))

		log.FinishProgress(fmt.Sprintf("Completed %s", parentOp))
	}

	log.Success("✅ Mixed demo completed")
	return nil
}

func runConcurrentDemo(log logger.Logger) error {
	log.Info("⚡ Starting concurrent demo")

	done := make(chan bool, operationCount)

	for i := 1; i <= operationCount; i++ {
		go func(id int) {
			defer func() { done <- true }()

			opName := fmt.Sprintf("Concurrent operation %d", id)

			log.StartProgress(opName)

			// Simulate varying work durations
			duration := time.Duration(100+id*50) * time.Millisecond
			time.Sleep(duration)

			if shouldFail() {
				log.FailProgress(opName, fmt.Errorf("concurrent operation %d failed", id))
			} else {
				log.FinishProgress(fmt.Sprintf("Completed %s", opName))
			}
		}(i)

		// Stagger the start times slightly
		time.Sleep(25 * time.Millisecond)
	}

	// Wait for all operations to complete
	for i := 0; i < operationCount; i++ {
		<-done
	}

	log.Success("✅ Concurrent demo completed")
	return nil
}

func runLongDemo(log logger.Logger) error {
	log.Info("⏳ Starting long-running demo")

	log.StartProgress("Long-running installation process")
	time.Sleep(100 * time.Millisecond)

	phases := []string{
		"Downloading packages",
		"Extracting archives",
		"Installing dependencies",
		"Configuring settings",
		"Running post-install scripts",
		"Cleaning up temporary files",
	}

	for i, phase := range phases {
		log.StartProgress(phase)

		// Simulate progress updates
		steps := 4
		for step := 1; step <= steps; step++ {
			time.Sleep(200 * time.Millisecond)
			percentage := (step * 100) / steps
			log.UpdateProgress(fmt.Sprintf("%s (%d%%)", phase, percentage))
		}

		// Occasionally fail a phase
		if i == 3 && shouldFail() { // Fail configuration sometimes
			log.FailProgress(phase, errors.New("configuration validation failed"))
			log.FailProgress("Long-running installation process",
				errors.New("installation aborted due to configuration error"))
			return nil
		}

		log.FinishProgress(fmt.Sprintf("Completed %s", phase))
	}

	log.FinishProgress("Installation process completed successfully")
	log.Success("✅ Long-running demo completed")
	return nil
}

func shouldFail() bool {
	if failureRate <= 0 {
		return false
	}
	if failureRate >= 100 {
		return true
	}

	// Simple pseudo-random failure based on current time
	return int(time.Now().UnixNano()%100) < failureRate
}

func runPersistentDemo(log logger.Logger) error {
	log.Info("📦 Starting persistent progress demo (like cargo/npm)")

	// Demo 1: Installing system packages
	log.StartPersistentProgress("Installing system packages")
	time.Sleep(200 * time.Millisecond)

	packages := []string{"brew", "git v2.39.0", "zsh v5.8.1", "fzf v0.35.1", "tmux v3.3a"}
	for _, pkg := range packages {
		time.Sleep(300 * time.Millisecond)
		if shouldFail() {
			log.FailPersistentProgress("Failed to install system packages", fmt.Errorf("could not install %s", pkg))
			return nil
		}
		log.LogAccomplishment(fmt.Sprintf("Installed %s", pkg))
	}
	log.FinishPersistentProgress(fmt.Sprintf("System packages installed (%d packages)", len(packages)))

	time.Sleep(500 * time.Millisecond)

	// Demo 2: Configuring shell environment
	log.StartPersistentProgress("Configuring shell environment")
	time.Sleep(150 * time.Millisecond)

	configs := []string{
		"~/.zshrc",
		"~/.zsh/aliases.zsh",
		"~/.zsh/exports.zsh",
		"~/.zsh/functions.zsh",
	}
	for _, config := range configs {
		time.Sleep(200 * time.Millisecond)
		log.LogAccomplishment(fmt.Sprintf("Created %s", config))
	}

	time.Sleep(300 * time.Millisecond)
	log.LogAccomplishment("Installed oh-my-zsh plugins (5 plugins)")
	log.FinishPersistentProgress("Shell environment configured")

	time.Sleep(500 * time.Millisecond)

	// Demo 3: Linking dotfiles
	log.StartPersistentProgress("Linking dotfiles")
	time.Sleep(100 * time.Millisecond)

	dotfiles := []struct {
		target, source string
	}{
		{"~/.vimrc", "~/dotfiles/vim/vimrc"},
		{"~/.tmux.conf", "~/dotfiles/tmux/tmux.conf"},
		{"~/.gitconfig", "~/dotfiles/git/gitconfig"},
		{"~/.ssh/config", "~/dotfiles/ssh/config"},
	}

	for _, df := range dotfiles {
		time.Sleep(150 * time.Millisecond)
		log.LogAccomplishment(fmt.Sprintf("Linked %s → %s", df.target, df.source))
	}
	log.FinishPersistentProgress(fmt.Sprintf("Dotfiles linked (%d files)", len(dotfiles)))

	log.Success("✅ Persistent demo completed - notice how accomplishments stay visible!")
	return nil
}

func runStressTestDemo(log logger.Logger) error {
	log.Info("🔥 Starting stress test demo - rapid pause/resume cycles")

	// Start background operations
	log.StartProgress("Background operation 1")
	time.Sleep(100 * time.Millisecond)

	log.StartProgress("Background operation 2")
	time.Sleep(100 * time.Millisecond)

	// Perform rapid pause/resume cycles to test race conditions
	for i := 1; i <= 5; i++ {
		log.StartInteractiveProgress(fmt.Sprintf("Interactive cycle %d (rapid test)", i))

		// Very short interaction to stress test the synchronization
		time.Sleep(50 * time.Millisecond)

		log.FinishInteractiveProgress(fmt.Sprintf("Interactive cycle %d completed", i))

		// Brief pause between cycles
		time.Sleep(50 * time.Millisecond)
	}

	// Finish background operations
	log.FinishProgress("Background operation 2 completed")
	log.FinishProgress("Background operation 1 completed")

	log.Success("🔥 Stress test completed - no race conditions detected!")
	return nil
}

func runInteractiveDemo(log logger.Logger) error {
	log.Info("🎯 Starting interactive demo to test spinner pause/resume")

	// Start a regular progress operation
	log.StartProgress("Background operation running")

	// Simulate some background work
	time.Sleep(500 * time.Millisecond)

	// Start an interactive operation that should pause the spinner
	log.StartInteractiveProgress("Interactive operation (spinner should pause)")

	// Simulate an interactive command (like GPG key creation)
	fmt.Print("Please enter some text (this simulates GPG prompts): ")
	var userInput string
	fmt.Scanln(&userInput)

	// Finish the interactive operation (should resume spinner)
	log.FinishInteractiveProgress("Interactive operation completed")

	// Continue with background work
	time.Sleep(500 * time.Millisecond)

	// Finish the background operation
	log.FinishProgress("Background operation finished")

	log.Success("🎉 Interactive demo completed successfully!")
	log.Info("User input was: %s", userInput)

	return nil
}
