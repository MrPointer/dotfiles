package main

import "github.com/MrPointer/dotfiles/installer/cmd"

// Version information - populated by goreleaser
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

func main() {
	// Set version information for the CLI
	cmd.SetVersionInfo(version, commit, date, builtBy)
	cmd.Execute()
}
