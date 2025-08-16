package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/MrPointer/dotfiles/installer/lib/compatibility"
	"github.com/MrPointer/dotfiles/installer/utils"
	"github.com/MrPointer/dotfiles/installer/utils/httpclient"
	"github.com/MrPointer/dotfiles/installer/utils/logger"
	"github.com/MrPointer/dotfiles/installer/utils/osmanager"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile                   string
	compatibilityConfigFile   string
	globalCompatibilityConfig *compatibility.CompatibilityConfig
	globalVerbosity           logger.VerbosityLevel
	verboseCount              int
	extraVerbose              bool
	progressFlag              bool
	nonInteractive            bool

	cliLogger        logger.Logger       = nil // Will be initialized before any command is executed
	globalCommander  utils.Commander     = nil // Will be initialized before any command is executed
	globalHttpClient                     = httpclient.NewDefaultHTTPClient()
	globalFilesystem                     = utils.NewDefaultFileSystem()
	globalOsManager  osmanager.OsManager = nil // Will be initialized before any command is executed
)

// HandleCompatibilityError displays compatibility error with install hints and exits.
func HandleCompatibilityError(err error, sysInfo compatibility.SystemInfo, log logger.Logger) {
	// Print the error symbol and message.
	fmt.Fprint(os.Stderr, "✘ ")
	log.Error("Your system isn't compatible with these dotfiles: %v", err)

	// Show install hints for missing prerequisites
	if len(sysInfo.Prerequisites.Missing) > 0 {
		fmt.Println()
		log.Info("Missing prerequisites and how to install them:")
		for _, name := range sysInfo.Prerequisites.Missing {
			if detail, exists := sysInfo.Prerequisites.Details[name]; exists {
				if detail.InstallHint != "" {
					fmt.Printf("  • %s: %s\n", detail.Description, detail.InstallHint)
				} else {
					fmt.Printf("  • %s\n", detail.Description)
				}
			} else {
				fmt.Printf("  • %s: (no install hint available)\n", name)
			}
		}
	}
	os.Exit(1)
}

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "dotfiles-installer",
	Short: "A tool to install (bootstrap) my dotfiles on any system",
	Long: `dotfiles-installer is a command-line tool that helps installing
my personal dotfiles on any system. It automates the process of setting up
necessary configurations, mostly for chezmoi to work properly.
It also installs essential packages and tools that I use on a daily basis.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

//nolint:gochecknoinits // Cobra requires an init function to set up the command structure.
func init() {
	cobra.OnInitialize(initConfig, initCompatibilityConfig, initLogger, initCommander, initOsManager)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().
		StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dotfiles-installer.yaml)")

	// Add compatibility config flag to root command so it's available globally.
	rootCmd.PersistentFlags().StringVar(&compatibilityConfigFile, "compat-config", "",
		"compatibility configuration file (uses embedded config by default)")

	// Verbosity flags: supports multiple levels
	// - No flags: Normal verbosity with progress indicators (default)
	// - -v or --verbose: Verbose level (adds Debug messages, progress disabled by default)
	// - -vv or --extra-verbose: Extra verbose level (adds Trace messages, progress disabled by default)
	// Note: Progress can be explicitly enabled with --progress even when using verbosity flags
	rootCmd.PersistentFlags().CountVarP(&verboseCount, "verbose", "v",
		"Enable verbose output (use -vv for extra verbose)")

	rootCmd.PersistentFlags().BoolVar(&extraVerbose, "extra-verbose", false,
		"Enable extra verbose output (equivalent to -vv)")

	// Progress flag: controls whether to show hierarchical progress indicators
	// - Default behavior: disabled (shows regular log messages instead)
	// - Explicit --progress: enables npm-style hierarchical display with spinners and timing information
	// - Always disabled in non-interactive mode regardless of other flags
	rootCmd.PersistentFlags().BoolVar(&progressFlag, "progress", false,
		"Show hierarchical progress indicators with spinners and timing")

	// Interactive flag: controls whether to allow user interaction
	// - When enabled, disables progress indicators and skips any prompts
	// - Affects all commands that might need user input
	rootCmd.PersistentFlags().BoolVar(&nonInteractive, "non-interactive", false,
		"Disable interactive mode (also disables progress indicators)")

	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("extra-verbose", rootCmd.PersistentFlags().Lookup("extra-verbose"))
	viper.BindPFlag("progress", rootCmd.PersistentFlags().Lookup("progress"))
	viper.BindPFlag("non-interactive", rootCmd.PersistentFlags().Lookup("non-interactive"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".dotfiles" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".dotfiles-installer")
	}

	viper.AutomaticEnv() // Read in environment variables that match.

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func initCompatibilityConfig() {
	// Initialize compatibility configuration.
	compatibilityConfig, err := compatibility.LoadCompatibilityConfig(viper.New(), compatibilityConfigFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading compatibility config: %v\n", err)
		os.Exit(1)
	}
	globalCompatibilityConfig = compatibilityConfig
}

// GetCompatibilityConfig returns the loaded compatibility configuration.
func GetCompatibilityConfig() *compatibility.CompatibilityConfig {
	return globalCompatibilityConfig
}

func initOsManager() {
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		globalOsManager = osmanager.NewUnixOsManager(cliLogger, globalCommander, osmanager.IsRoot())
	} else {
		cliLogger.Error("The system may be compatible, but we haven't implemented an OS manager for it yet. Please open an issue on GitHub to request support for this OS.")
		os.Exit(1)
	}
}

func initCommander() {
	globalCommander = utils.NewDefaultCommander(cliLogger)
}

func initLogger() {
	// Determine verbosity level based on flags
	if extraVerbose || verboseCount >= 2 {
		globalVerbosity = logger.ExtraVerbose
	} else if verboseCount >= 1 {
		globalVerbosity = logger.Verbose
	} else {
		globalVerbosity = logger.Normal
	}

	// Create logger with or without progress based on ShouldShowProgress()
	if ShouldShowProgress() {
		cliLogger = logger.NewProgressCliLogger(globalVerbosity)
	} else {
		cliLogger = logger.NewCliLogger(globalVerbosity)
	}
}

// ShouldShowProgress determines if hierarchical progress indicators should be shown.
// The logic is:
// 1. If --non-interactive is set, never show progress
// 2. If --progress was explicitly set, show progress (unless non-interactive)
// 3. Otherwise, default to no progress (regular logging)
func ShouldShowProgress() bool {
	// If non-interactive mode is enabled, never show progress
	if nonInteractive {
		return false
	}

	// Only show progress if explicitly requested
	return progressFlag
}

// GetVerbosity returns the current global verbosity level.
func GetVerbosity() logger.VerbosityLevel {
	return globalVerbosity
}

// IsNonInteractive returns whether non-interactive mode is enabled.
func IsNonInteractive() bool {
	return nonInteractive
}
