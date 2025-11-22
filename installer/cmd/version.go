package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	versionInfo VersionInfo
)

// VersionInfo holds version information for the application
type VersionInfo struct {
	Version string
	Commit  string
	Date    string
	BuiltBy string
}

// SetVersionInfo sets the version information for the application
func SetVersionInfo(version, commit, date, builtBy string) {
	versionInfo = VersionInfo{
		Version: version,
		Commit:  commit,
		Date:    date,
		BuiltBy: builtBy,
	}
	// Set the version on the root command
	rootCmd.Version = version
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long: `Print detailed version information including build details.

This command displays the version number, Git commit hash, build date,
and other build information for the dotfiles installer.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("dotfiles-installer version %s\n", versionInfo.Version)
		fmt.Printf("  commit: %s\n", versionInfo.Commit)
		fmt.Printf("  built: %s\n", versionInfo.Date)
		fmt.Printf("  built by: %s\n", versionInfo.BuiltBy)
		fmt.Printf("  go version: %s\n", runtime.Version())
		fmt.Printf("  platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	},
}

//nolint:gochecknoinits // Cobra requires init function to register commands
func init() {
	rootCmd.AddCommand(versionCmd)
}
