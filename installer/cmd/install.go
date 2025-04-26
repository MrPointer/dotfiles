/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	workEnvironment      bool
	workName             string
	workEmail            string
	shell                string
	installBrew          bool
	installShellWithBrew bool
	multiUserSystem      bool
	gitCloneProtocol     string
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install dotfiles",
	Long: `Install dotfiles on the current system.
This command will set up the necessary configurations and
install essential packages and tools that I use on a daily basis.
It automates the process of setting up the dotfiles,
making it easier to get started with a new system.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get the globally loaded compatibility config
		checkCompatibilityCmd.Run(cmd, args)

		fmt.Println("Installing dotfiles...")
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Flags().BoolVar(&workEnvironment, "work-env", false,
		"Treat this installation as a work environment (affects some dotfiles)")
	installCmd.Flags().StringVar(&workName, "work-name", "",
		"Use the given name as the work's name")
	installCmd.Flags().StringVar(&workEmail, "work-email", "",
		"Use the given email address as work's email address")
	installCmd.Flags().StringVar(&shell, "shell", "zsh",
		"Install given shell if required and set it as user's default")
	installCmd.Flags().BoolVar(&installBrew, "install-brew", true,
		"Install brew if not already installed")
	installCmd.Flags().BoolVar(&installShellWithBrew, "install-shell-with-brew", true,
		"Install shell with brew if not already installed")
	installCmd.Flags().BoolVar(&multiUserSystem, "multi-user-system", false,
		"Treat this system as a multi-user system (affects some dotfiles)")
	installCmd.Flags().StringVar(&gitCloneProtocol, "git-clone-protocol", "ssh",
		"Use the given git clone protocol (ssh or https) for git operations")

	viper.BindPFlag("work-env", installCmd.Flags().Lookup("work-env"))
	viper.BindPFlag("work-name", installCmd.Flags().Lookup("work-name"))
	viper.BindPFlag("work-email", installCmd.Flags().Lookup("work-email"))
	viper.BindPFlag("shell", installCmd.Flags().Lookup("shell"))
	viper.BindPFlag("install-brew", installCmd.Flags().Lookup("install-brew"))
	viper.BindPFlag("install-shell-with-brew", installCmd.Flags().Lookup("install-shell-with-brew"))
	viper.BindPFlag("multi-user-system", installCmd.Flags().Lookup("multi-user-system"))
	viper.BindPFlag("git-clone-protocol", installCmd.Flags().Lookup("git-clone-protocol"))
}
