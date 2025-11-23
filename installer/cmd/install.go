package cmd

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/MrPointer/dotfiles/installer/cli"
	"github.com/MrPointer/dotfiles/installer/lib/apt"
	"github.com/MrPointer/dotfiles/installer/lib/brew"
	"github.com/MrPointer/dotfiles/installer/lib/compatibility"
	"github.com/MrPointer/dotfiles/installer/lib/dotfilesmanager"
	"github.com/MrPointer/dotfiles/installer/lib/dotfilesmanager/chezmoi"
	"github.com/MrPointer/dotfiles/installer/lib/gpg"
	"github.com/MrPointer/dotfiles/installer/lib/pkgmanager"
	"github.com/MrPointer/dotfiles/installer/lib/shell"
	"github.com/MrPointer/dotfiles/installer/utils/logger"
	"github.com/MrPointer/dotfiles/installer/utils/privilege"
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

	gitCloneProtocol     string
	verbose              bool
	installPrerequisites bool
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
		// Use the global logger (already configured with proper progress/verbosity settings)
		installLogger := cliLogger

		// Get basic system info first to determine if we need to install Homebrew early
		osDetector := compatibility.NewDefaultOSDetector()
		basicSysInfo, err := osDetector.DetectSystem()
		if err != nil {
			installLogger.Error("Failed to get basic system information: %v", err)
			os.Exit(1)
		}

		// For macOS, install Homebrew BEFORE checking prerequisites if requested
		// This solves the chicken-and-egg problem where we need Homebrew to install prerequisites
		if basicSysInfo.OSName == "darwin" && installBrew {
			brewPath, err := installHomebrew(basicSysInfo, installLogger)
			if err != nil {
				installLogger.Error("Failed to install Homebrew: %v", err)
				os.Exit(1)
			}
			globalPackageManager = brew.NewBrewPackageManager(installLogger, globalCommander, globalOsManager, brewPath, GetDisplayMode())
		}

		// Check system compatibility and get system info.
		config := GetCompatibilityConfig()

		// Now check system compatibility and get full system info including prerequisites
		sysInfo, err := compatibility.CheckCompatibilityWithDetector(config, osDetector, globalOsManager)
		if err != nil {
			if handlePrerequisiteInstallation(sysInfo, installLogger) {
				// Prerequisites were installed, re-check compatibility
				sysInfo, err = compatibility.CheckCompatibility(config, globalOsManager)
				if err != nil {
					HandleCompatibilityError(err, sysInfo, installLogger)
				}
			} else {
				HandleCompatibilityError(err, sysInfo, installLogger)
			}
		}
		installLogger.Success("System compatibility check passed")

		// Install Homebrew for non-macOS systems or if not already installed
		if installBrew && (basicSysInfo.OSName != "darwin" || globalPackageManager == nil) {
			brewPath, err := installHomebrew(sysInfo, installLogger)
			if err != nil {
				installLogger.Error("Failed to install Homebrew: %v", err)
				os.Exit(1)
			}
			globalPackageManager = brew.NewBrewPackageManager(installLogger, globalCommander, globalOsManager, brewPath, GetDisplayMode())
		}

		if err := installShell(installLogger); err != nil {
			installLogger.Error("Failed to install shell: %v", err)
			os.Exit(1)
		}

		if err := setupGpgKeys(installLogger); err != nil {
			installLogger.Error("Failed to setup GPG keys: %v", err)
			os.Exit(1)
		}

		if err := setupDotfilesManager(installLogger); err != nil {
			installLogger.Error("Failed to setup dotfiles manager: %v", err)
			os.Exit(1)
		}

		installLogger.Success("Installation completed successfully")
	},
}

// createPackageManagerForSystem creates the appropriate package manager for the current system.
func createPackageManagerForSystem(sysInfo *compatibility.SystemInfo) pkgmanager.PackageManager {
	// If we already have a global package manager set up (e.g., Homebrew installed early for macOS), use it
	if globalPackageManager != nil {
		return globalPackageManager
	}

	switch sysInfo.OSName {
	case "linux":
		switch sysInfo.DistroName {
		case "ubuntu", "debian":
			return apt.NewAptPackageManager(cliLogger, globalCommander, globalOsManager, privilege.NewDefaultEscalator(cliLogger, globalCommander, globalOsManager), GetDisplayMode())
		default:
			cliLogger.Warning("Unsupported Linux distribution for automatic package installation: %s", sysInfo.DistroName)
			return nil
		}
	case "darwin":
		brewPath, err := brew.DetectBrewPath(sysInfo, "")
		if err != nil {
			cliLogger.Error("Failed to detect Homebrew path: %v", err)
			return nil
		}
		return brew.NewBrewPackageManager(cliLogger, globalCommander, globalOsManager, brewPath, GetDisplayMode())
	default:
		cliLogger.Warning("Unsupported operating system for automatic package installation: %s", sysInfo.OSName)
		return nil
	}
}

// handlePrerequisiteInstallation handles automatic installation of missing prerequisites.
// Returns true if prerequisites were installed and compatibility should be re-checked.
func handlePrerequisiteInstallation(sysInfo compatibility.SystemInfo, log logger.Logger) bool {
	// Only attempt installation if we have missing prerequisites and the flag is set
	if len(sysInfo.Prerequisites.Missing) == 0 {
		return false
	}

	// Create package manager for this system
	packageManager := createPackageManagerForSystem(&sysInfo)
	if packageManager == nil {
		log.Warning("Cannot install prerequisites automatically on this system")
		return false
	}

	var prerequisitesToInstall []string

	// In non-interactive mode, or if explicitly requested, install all missing prerequisites automatically
	if IsNonInteractive() || installPrerequisites {
		prerequisitesToInstall = sysInfo.Prerequisites.Missing
		log.StartProgress("Installing missing prerequisites automatically")
	} else {
		// In interactive mode, let user select which prerequisites to install
		prerequisiteSelector := cli.NewDefaultPrerequisiteSelector()

		// Convert compatibility.PrerequisiteDetail to cli.PrerequisiteDetail
		cliDetails := make(map[string]cli.PrerequisiteDetail)
		for name, detail := range sysInfo.Prerequisites.Details {
			cliDetails[name] = cli.PrerequisiteDetail{
				Name:        detail.Name,
				Available:   detail.Available,
				Command:     detail.Command,
				Description: detail.Description,
				InstallHint: detail.InstallHint,
			}
		}

		selectedPrerequisites, err := prerequisiteSelector.SelectPrerequisites(
			sysInfo.Prerequisites.Missing,
			cliDetails,
		)
		if err != nil {
			log.Error("Failed to select prerequisites: %v", err)
			return false
		}

		if len(selectedPrerequisites) == 0 {
			log.Info("No prerequisites selected for installation")
			return false
		}

		prerequisitesToInstall = selectedPrerequisites
		log.StartProgress(fmt.Sprintf("Installing %d selected prerequisites", len(selectedPrerequisites)))
	}

	// Install each selected prerequisite
	installed := false
	for _, name := range prerequisitesToInstall {
		if detail, exists := sysInfo.Prerequisites.Details[name]; exists {
			log.StartProgress(fmt.Sprintf("Installing %s", detail.Description))

			// Use the prerequisite name directly as the package name
			packageInfo := pkgmanager.NewRequestedPackageInfo(name, nil)

			err := packageManager.InstallPackage(packageInfo)
			if err != nil {
				log.FailProgress(fmt.Sprintf("Failed to install %s", detail.Description), err)
				return false
			}

			log.FinishProgress(fmt.Sprintf("%s installed successfully", detail.Description))
			installed = true
		}
	}

	if installed {
		log.FinishProgress("Prerequisites installation completed")
		return true
	}

	return false
}

// installHomebrew installs Homebrew if not already installed.
func installHomebrew(sysInfo compatibility.SystemInfo, log logger.Logger) (string, error) {
	log.StartProgress("Setting up Homebrew")

	// Create BrewInstaller using the new API.
	installer := brew.NewBrewInstaller(brew.Options{
		SystemInfo:  &sysInfo,
		Logger:      cliLogger,
		Commander:   globalCommander,
		HTTPClient:  globalHttpClient,
		OsManager:   globalOsManager,
		Fs:          globalFilesystem,
		DisplayMode: GetDisplayMode(),
	})

	log.StartProgress("Checking Homebrew availability")
	isAvailable, err := installer.IsAvailable()
	if err != nil {
		log.FailProgress("Failed to check Homebrew availability", err)
		return "", err
	}

	if isAvailable {
		log.FinishProgress("Homebrew is already available")

		log.Debug("Detecting Homebrew path")
		brewPath, err := brew.DetectBrewPath(&sysInfo, "")
		if err != nil {
			log.FailProgress("Failed to detect Homebrew path", err)
			return "", err
		}
		log.Debug("Homebrew path detected: %s", brewPath)

		log.Debug("Updating PATH environment variable with Homebrew binaries")
		// Although Homebrew is already installed, we still need to update the PATH environment variable,
		// because it may not be set correctly.
		err = brew.UpdatePathWithBrewBinaries(brewPath)
		if err != nil {
			log.FailProgress("Failed to update PATH with Homebrew binaries", err)
			return "", err
		}
		log.Debug("PATH updated with Homebrew binaries")

		log.FinishProgress("Homebrew is ready")
		return brewPath, nil
	}
	log.FinishProgress("Homebrew not found")

	log.StartProgress("Installing Homebrew")
	if err := installer.Install(); err != nil {
		log.FailProgress("Failed to install Homebrew", err)
		return "", err
	}
	log.FinishProgress("Homebrew installation completed")

	log.Debug("Detecting Homebrew path after installation")
	brewPath, err := brew.DetectBrewPath(&sysInfo, "")
	if err != nil {
		log.FailProgress("Failed to detect Homebrew path after installation", err)
		return "", err
	}
	log.Debug("Homebrew path detected: %s", brewPath)

	log.FinishProgress("Homebrew setup completed")
	return brewPath, nil
}

func installShell(log logger.Logger) error {
	log.StartProgress(fmt.Sprintf("Setting up %s shell", shellName))

	shellInstaller := shell.NewDefaultShellInstaller(shellName, globalOsManager, globalPackageManager, log)

	log.StartProgress(fmt.Sprintf("Checking %s shell availability", shellName))
	isAvailable, err := shellInstaller.IsAvailable()
	if err != nil {
		log.FailProgress(fmt.Sprintf("Failed to check %s shell availability", shellName), err)
		return err
	}

	if isAvailable {
		log.FinishProgress(fmt.Sprintf("%s shell is already available", shellName))
		log.FinishProgress(fmt.Sprintf("%s shell is ready", shellName))
		return nil
	}
	log.FinishProgress(fmt.Sprintf("%s shell not found", shellName))

	log.StartProgress(fmt.Sprintf("Installing %s shell", shellName))
	if err := shellInstaller.Install(context.TODO()); err != nil {
		log.FailProgress(fmt.Sprintf("Failed to install %s shell", shellName), err)
		return err
	}
	log.FinishProgress(fmt.Sprintf("%s shell installed successfully", shellName))

	log.FinishProgress(fmt.Sprintf("%s shell setup completed", shellName))
	return nil
}

func setupGpgKeys(log logger.Logger) error {
	err := installGpgClient(log)
	if err != nil {
		return err
	}

	if IsNonInteractive() {
		log.Warning("Skipping GPG key setup in non-interactive mode - You will need to set them up manually")
		return nil
	}

	log.StartProgress("Setting up GPG keys")

	gpgClient := gpg.NewDefaultGpgClient(
		globalOsManager,
		globalFilesystem,
		globalCommander,
		cliLogger,
	)

	log.StartProgress("Checking for existing GPG keys")
	existingKeys, err := gpgClient.ListAvailableKeys()
	if err != nil {
		log.FailProgress("Failed to list available GPG keys", err)
		return err
	}
	log.FinishProgress("GPG keys check completed")

	if len(existingKeys) == 0 {
		log.StartInteractiveProgress("Creating new GPG key pair")
		keyId, err := gpgClient.CreateKeyPair()
		if err != nil {
			log.FailInteractiveProgress("Failed to create GPG key pair", err)
			return err
		}
		selectedGpgKey = keyId
		log.FinishInteractiveProgress("GPG key pair created successfully")
	} else {
		log.StartInteractiveProgress("Selecting GPG key from existing keys")
		gpgSelector := cli.NewDefaultGpgKeySelector()
		selectedKey, err := gpgSelector.SelectKey(existingKeys)
		if err != nil {
			log.FailInteractiveProgress("Failed to select GPG key", err)
			return err
		}
		selectedGpgKey = selectedKey
		log.FinishInteractiveProgress("GPG key selected successfully")
	}

	log.FinishProgress("GPG keys set up successfully")
	return nil
}

// installGpgClient installs the GPG client if not already available.
func installGpgClient(log logger.Logger) error {
	log.StartProgress("Setting up GPG client")

	// Create GpgClientInstaller using the new API.
	installer := gpg.NewGpgInstaller(
		log,
		globalCommander,
		globalOsManager,
		globalPackageManager,
	)

	log.StartProgress("Checking GPG client availability")
	isAvailable, err := installer.IsAvailable()
	if err != nil {
		log.FailProgress("Failed to check GPG client availability", err)
		return err
	}

	if isAvailable {
		log.FinishProgress("GPG client is already available")
		log.FinishProgress("GPG client is ready")
		return nil
	}
	log.FinishProgress("GPG client not found")

	log.StartProgress("Installing GPG client")
	if err := installer.Install(context.TODO()); err != nil {
		log.FailProgress("Failed to install GPG client", err)
		return err
	}
	log.FinishProgress("GPG client installed successfully")

	log.FinishProgress("GPG client setup completed")
	return nil
}

func setupDotfilesManager(log logger.Logger) error {
	log.StartProgress("Setting up dotfiles manager")

	dm, err := chezmoi.TryStandardChezmoiManager(log, globalFilesystem, globalOsManager, globalCommander, globalPackageManager, globalHttpClient, GetDisplayMode(), chezmoi.DefaultGitHubUsername, gitCloneProtocol == "ssh")
	if err != nil {
		log.FailProgress("Failed to create dotfiles manager", err)
		return err
	}

	log.StartProgress("Installing dotfiles manager")
	err = dm.Install()
	if err != nil {
		log.FailProgress("Failed to install dotfiles manager", err)
		return err
	}
	log.FinishProgress("Dotfiles manager installed successfully")

	log.StartProgress("Initializing dotfiles manager data")
	if err := initDotfilesManagerData(dm); err != nil {
		log.FailProgress("Failed to initialize dotfiles manager data", err)
		return err
	}
	log.FinishProgress("Dotfiles manager data initialized successfully")

	log.StartProgress("Applying dotfiles configuration")
	if err := dm.Apply(); err != nil {
		log.FailProgress("Failed to apply dotfiles configuration", err)
		return err
	}
	log.FinishProgress("Dotfiles configuration applied successfully")

	log.FinishProgress("Dotfiles manager setup completed successfully")
	return nil
}

func initDotfilesManagerData(dm dotfilesmanager.DotfilesManager) error {
	dotfiles_data := dotfilesmanager.DotfilesData{
		FirstName: "Timor",
		LastName:  "Gruber",
		Email:     "timor.gruber@gmail.com",
		SystemData: mo.Some(dotfilesmanager.DotfilesSystemData{
			Shell: shellName,
		}),
	}

	if workEnvironment {
		work_data := dotfilesmanager.DotfilesWorkEnvData{
			WorkName:  workName,
			WorkEmail: workEmail,
		}
		dotfiles_data.WorkEnv = mo.Some(work_data)

		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		generic_work_dotfiles_dir := path.Join(userHomeDir, ".work")

		dotfiles_data.SystemData = dotfiles_data.SystemData.Map(func(value dotfilesmanager.DotfilesSystemData) (dotfilesmanager.DotfilesSystemData, bool) {
			value.GenericWorkProfile = mo.Some(path.Join(generic_work_dotfiles_dir, "profile"))
			value.SpecificWorkProfile = mo.Some(path.Join(generic_work_dotfiles_dir, workName, "profile"))
			return value, true
		})
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
	installCmd.Flags().StringVar(&workName, "work-name", "sedg",
		"Use the given name as the work's name")
	installCmd.Flags().StringVar(&workEmail, "work-email", "timor.gruber@solaredge.com",
		"Use the given email address as work's email address")
	installCmd.Flags().StringVar(&shellName, "shell", "zsh",
		"Install given shell if required and set it as user's default")
	installCmd.Flags().BoolVar(&installBrew, "install-brew", true,
		"Install brew if not already installed")
	installCmd.Flags().BoolVar(&installShellWithBrew, "install-shell-with-brew", true,
		"Install shell with brew if not already installed")
	installCmd.Flags().StringVar(&gitCloneProtocol, "git-clone-protocol", "https",
		"Use the given git clone protocol (ssh or https) for git operations")
	installCmd.Flags().BoolVar(&installPrerequisites, "install-prerequisites", false,
		"Automatically install missing prerequisites")

	viper.BindPFlag("work-env", installCmd.Flags().Lookup("work-env"))
	viper.BindPFlag("work-name", installCmd.Flags().Lookup("work-name"))
	viper.BindPFlag("work-email", installCmd.Flags().Lookup("work-email"))
	viper.BindPFlag("shell", installCmd.Flags().Lookup("shell"))
	viper.BindPFlag("install-brew", installCmd.Flags().Lookup("install-brew"))
	viper.BindPFlag("install-shell-with-brew", installCmd.Flags().Lookup("install-shell-with-brew"))
	viper.BindPFlag("git-clone-protocol", installCmd.Flags().Lookup("git-clone-protocol"))

	viper.BindPFlag("install-prerequisites", installCmd.Flags().Lookup("install-prerequisites"))
}
