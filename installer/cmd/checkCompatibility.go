/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/MrPointer/dotfiles/installer/lib/compatibility"
	"github.com/MrPointer/dotfiles/installer/utils/logger"
	"github.com/spf13/cobra"
)

// checkCompatibilityCmd represents the check-compatibility command.
var checkCompatibilityCmd = &cobra.Command{
	Use:   "check-compatibility",
	Short: "Check compatibility of your dotfiles with the current system",
	Long: `Checks whether the current system is compatible with the dotfiles,
as they have some distribution-specific configurations. This command will
provide a report on the compatibility status.

It's recommended to run this command before attempting to install the dotfiles.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Create a CLI logger for output.
		log := logger.NewCliLogger()

		// Get the globally loaded compatibility config.
		config := GetCompatibilityConfig()

		// Check system compatibility.
		sysInfo, err := compatibility.CheckCompatibility(config)
		if err != nil {
			// Print the error symbol and message.
			fmt.Fprint(os.Stderr, "✘ ")
			log.Error("Your system isn't compatible with these dotfiles: %v", err)
			os.Exit(1)
		}

		// Print the success symbol and message.
		fmt.Print("✔︎ ")
		log.Success("Your system is compatible with these dotfiles!")

		// Print detected system information if verbose flag is set.
		if verbose {
			fmt.Println() // Add an empty line for better spacing.
			log.Info("Detected system information:")
			fmt.Printf("OS: %s\n", sysInfo.OSName)
			fmt.Printf("Distribution: %s\n", sysInfo.DistroName)
			fmt.Printf("Architecture: %s\n", sysInfo.Arch)
		}
	},
}

//nolint:gochecknoinits // Cobra requires an init function to set up the command structure.
func init() {
	rootCmd.AddCommand(checkCompatibilityCmd)
	// No need for additional flags here, as we use the global compatibility config.
}
