package cmd

import (
	"os"

	"github.com/MrPointer/dotfiles/installer/cli"
	"github.com/MrPointer/dotfiles/installer/lib/brew"
	"github.com/MrPointer/dotfiles/installer/lib/compatibility"
	"github.com/MrPointer/dotfiles/installer/lib/dotfilesmanager"
	"github.com/MrPointer/dotfiles/installer/lib/dotfilesmanager/chezmoi"
	"github.com/MrPointer/dotfiles/installer/lib/gpg"
	"github.com/MrPointer/dotfiles/installer/lib/pkgmanager"
	"github.com/MrPointer/dotfiles/installer/lib/shell"
	"github.com/MrPointer/dotfiles/installer/utils/osmanager"
	"github.com/samber/mo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Options and flags for the install command.
var (
	workEnvironment      bool
	workName             string
	workEmail            string
	shellName            string
	installBrew          bool
	installShellWithBrew bool
	multiUserSystem      bool
	gitCloneProtocol     string
	verbose              bool
	nonInteractive       bool
)

// global variables for the command execution context.
var (
	globalPackageManager pkgmanager.PackageManager = nil // set later based on passed flags
)

// output variables stored in global context
var (
	selectedGpgKey string
)

// installCmd represents the install command.
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install dotfiles",
	Long: `Install dotfiles on the current system.
This command will set up the necessary configurations and
install essential packages and tools that I use on a daily basis.
It automates the process of setting up the dotfiles,
making it easier to get started with a new system.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check system compatibility and get system info.
		config := GetCompatibilityConfig()
		sysInfo, err := compatibility.CheckCompatibility(config)
		if err != nil {
			cliLogger.Error("Your system isn't compatible with these dotfiles: %v", err)
			os.Exit(1)
		}
		if sysInfo.OSName == "linux" || sysInfo.OSName == "darwin" {
			globalOsManager = osmanager.NewUnixOsManager(cliLogger, globalCommander, osmanager.IsRoot())
		} else {
			cliLogger.Error("The system should be compatible, but we haven't implemented an OS manager for it yet. Please open an issue on GitHub to request support for this OS.")
			os.Exit(1)
		}
		cliLogger.Success("System compatibility check passed")

		cliLogger.Info("Installing dotfiles...")

		// Install Homebrew if specified and not already available.
		if installBrew {
			if err := installHomebrew(&sysInfo); err != nil {
				cliLogger.Error("Failed to install Homebrew: %v", err)
				os.Exit(1)
			}
			globalPackageManager = brew.NewBrewPackageManager(cliLogger, globalCommander, globalOsManager)
		}

		if err := installShell(); err != nil {
			cliLogger.Error("Failed to install shell: %v", err)
			os.Exit(1)
		}

		if err := setupGpgKeys(); err != nil {
			cliLogger.Error("Failed to setup GPG keys: %v", err)
			os.Exit(1)
		}

		if err := setupDotfilesManager(); err != nil {
			cliLogger.Error("Failed to setup dotfiles manager: %v", err)
			os.Exit(1)
		}

		cliLogger.Success("🪄 Installation completed 🎉")
	},
}

// installHomebrew installs Homebrew if not already installed.
func installHomebrew(sysInfo *compatibility.SystemInfo) error {
	// Create BrewInstaller using the new API.
	installer := brew.NewBrewInstaller(brew.Options{
		SystemInfo:      sysInfo,
		Logger:          cliLogger,
		Commander:       globalCommander,
		HTTPClient:      globalHttpClient,
		OsManager:       globalOsManager,
		Fs:              globalFilesystem,
		MultiUserSystem: multiUserSystem,
	})

	isAvailable, err := installer.IsAvailable()
	if err != nil {
		return err
	}
	if isAvailable {
		cliLogger.Success("Homebrew is already installed")
		return nil
	}

	if err := installer.Install(); err != nil {
		return err
	}

	cliLogger.Success("Homebrew installed successfully")
	return nil
}

func installShell() error {
	shellInstaller := shell.NewDefaultShellInstaller(shellName, globalOsManager, globalPackageManager)

	isAvailable, err := shellInstaller.IsAvailable()
	if err != nil {
		return err
	}
	if isAvailable {
		cliLogger.Success("%s shell is already installed", shellName)
		return nil
	}

	if err := shellInstaller.Install(nil); err != nil { // Pass context if needed.
		return err
	}

	cliLogger.Success("%s shell installed successfully", shellName)
	return nil
}

func setupGpgKeys() error {
	err := installGpgClient()
	if err != nil {
		return err
	}

	if nonInteractive {
		cliLogger.Warning("Skipping GPG key setup in non-interactive mode - You will need to set them up manually")
		return nil
	}

	gpgClient := gpg.NewDefaultGpgClient(
		globalOsManager,
		globalCommander,
	)

	existingKeys, err := gpgClient.ListAvailableKeys()
	if err != nil {
		return err
	}

	if len(existingKeys) == 0 {
		keyId, err := gpgClient.CreateKeyPair()
		if err != nil {
			return err
		}
		selectedGpgKey = keyId
	} else {
		gpgSelector := cli.NewDefaultGpgKeySelector()
		selectedKey, err := gpgSelector.SelectKey(existingKeys)
		if err != nil {
			return err
		}
		selectedGpgKey = selectedKey
	}

	cliLogger.Success("GPG keys set up successfully")
	return nil
}

// installGpgClient installs the GPG client if not already available.
func installGpgClient() error {
	// Create GpgClientInstaller using the new API.
	installer := gpg.NewGpgInstaller(
		cliLogger,
		globalCommander,
		globalOsManager,
		globalPackageManager,
	)

	isAvailable, err := installer.IsAvailable()
	if err != nil {
		return err
	}
	if isAvailable {
		cliLogger.Success("GPG client is already installed")
		return nil
	}

	if err := installer.Install(nil); err != nil { // Pass context if needed.
		return err
	}
	cliLogger.Success("✔️ GPG client installed successfully")

	return nil
}

func setupDotfilesManager() error {
	dm, err := chezmoi.TryStandardChezmoiManager(globalFilesystem, globalOsManager, globalCommander, globalPackageManager, globalHttpClient, chezmoi.DefaultGitHubUsername, gitCloneProtocol == "ssh")
	if err != nil {
		return err
	}

	cliLogger.Info("Installing dotfiles manager")
	err = dm.Install()
	if err != nil {
		return err
	}
	cliLogger.Success("✔️ Dotfiles manager installed successfully")

	cliLogger.Info("Initializing dotfiles manager data")
	if err := initDotfilesManagerData(dm); err != nil {
		return err
	}
	cliLogger.Success("✔️ Dotfiles manager data initialized successfully")

	cliLogger.Info("Applying dotfiles manager")
	if err := dm.Apply(); err != nil {
		return err
	}
	cliLogger.Success("✔️ Dotfiles manager data applied successfully")

	return nil
}

func initDotfilesManagerData(dm dotfilesmanager.DotfilesManager) error {
	dotfiles_data := dotfilesmanager.DotfilesData{
		FirstName: "Timor",
		LastName:  "Gruber",
		Email:     "timor.gruber@gmail.com",
	}

	if workEnvironment {
		work_data := dotfilesmanager.DotfilesWorkEnvData{
			WorkName:  workName,
			WorkEmail: workEmail,
		}
		dotfiles_data.WorkEnv = mo.Some(work_data)
	}

	if selectedGpgKey != "" {
		dotfiles_data.GpgSigningKey = mo.Some(selectedGpgKey)
	}

	return dm.Initialize(dotfiles_data)
}

//nolint:gochecknoinits // Cobra requires an init function to set up the command structure.
func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Flags().BoolVar(&workEnvironment, "work-env", false,
		"Treat this installation as a work environment (affects some dotfiles)")
	installCmd.Flags().StringVar(&workName, "work-name", "",
		"Use the given name as the work's name")
	installCmd.Flags().StringVar(&workEmail, "work-email", "",
		"Use the given email address as work's email address")
	installCmd.Flags().StringVar(&shellName, "shell", "zsh",
		"Install given shell if required and set it as user's default")
	installCmd.Flags().BoolVar(&installBrew, "install-brew", true,
		"Install brew if not already installed")
	installCmd.Flags().BoolVar(&installShellWithBrew, "install-shell-with-brew", true,
		"Install shell with brew if not already installed")
	installCmd.Flags().BoolVar(&multiUserSystem, "multi-user-system", false,
		"Treat this system as a multi-user system (affects some dotfiles)")
	installCmd.Flags().StringVar(&gitCloneProtocol, "git-clone-protocol", "ssh",
		"Use the given git clone protocol (ssh or https) for git operations")
	installCmd.Flags().BoolVarP(&verbose, "verbose", "v", false,
		"Enable verbose output")
	installCmd.Flags().BoolVar(&nonInteractive, "non-interactive", false,
		"Disable interactive mode")

	viper.BindPFlag("work-env", installCmd.Flags().Lookup("work-env"))
	viper.BindPFlag("work-name", installCmd.Flags().Lookup("work-name"))
	viper.BindPFlag("work-email", installCmd.Flags().Lookup("work-email"))
	viper.BindPFlag("shell", installCmd.Flags().Lookup("shell"))
	viper.BindPFlag("install-brew", installCmd.Flags().Lookup("install-brew"))
	viper.BindPFlag("install-shell-with-brew", installCmd.Flags().Lookup("install-shell-with-brew"))
	viper.BindPFlag("multi-user-system", installCmd.Flags().Lookup("multi-user-system"))
	viper.BindPFlag("git-clone-protocol", installCmd.Flags().Lookup("git-clone-protocol"))
	viper.BindPFlag("verbose", installCmd.Flags().Lookup("verbose"))
	viper.BindPFlag("interactive", installCmd.Flags().Lookup("interactive"))
}
