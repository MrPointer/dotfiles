/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/MrPointer/dotfiles/installer/lib/compatibility"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// checkCompatibilityCmd represents the check-compatibility command
var checkCompatibilityCmd = &cobra.Command{
	Use:   "check-compatibility",
	Short: "Check compatibility of your dotfiles with the current system",
	Long: `Checks whether the current system is compatible with the dotfiles,
as they have some distribution-specific configurations. This command will
provide a report on the compatibility status.

It's recommended to run this command before attempting to install the dotfiles.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get the globally loaded compatibility config
		config := GetCompatibilityConfig()

		// Check system compatibility
		if err := compatibility.CheckCompatibility(config); err != nil {
			style := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#c71a3a"))

			fmt.Fprint(os.Stderr, style.Render("✘"))
			fmt.Fprintf(os.Stderr, " Your system isn't compatible with these dotfiles: %v\n", err)
			os.Exit(1)
		}

		style := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#2cb851"))

		fmt.Print(style.Render("✔︎"))
		fmt.Println(" Your system is compatible with these dotfiles!")
	},
}

func init() {
	rootCmd.AddCommand(checkCompatibilityCmd)
	// No need for additional flags here, as we use the global compatibility config
}
