package cmd

import (
	"fmt"
	"testing"

	"github.com/MrPointer/dotfiles/installer/lib/dnf"
	"github.com/MrPointer/dotfiles/installer/lib/packageresolver"
	"github.com/MrPointer/dotfiles/installer/lib/pkgmanager"
	"github.com/MrPointer/dotfiles/installer/utils"
	"github.com/MrPointer/dotfiles/installer/utils/logger"
	"github.com/MrPointer/dotfiles/installer/utils/osmanager"
	"github.com/MrPointer/dotfiles/installer/utils/privilege"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

// Test removed due to dependency on uninitialized global variables

func Test_PrerequisiteTypePreservation_EndToEnd(t *testing.T) {
	// This test verifies that the type information flows correctly
	// from package resolution through to package manager installation

	// Create mock components
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(command string, args []string, options ...utils.Option) (*utils.Result, error) {
			// Capture the actual command being run
			if command == "sudo" && len(args) >= 4 {
				if args[1] == "group" && args[2] == "install" {
					// This is a group installation command - exactly what we want!
					t.Logf("✅ Correct group install command: %s %v", command, args)
					require.Contains(t, args, "Development Tools", "Should install Development Tools group")
					return &utils.Result{}, nil
				} else if args[1] == "install" {
					// This is a regular package installation
					t.Logf("✅ Correct package install command: %s %v", command, args)
					require.Contains(t, args, "procps-ng", "Should install procps-ng package")
					return &utils.Result{}, nil
				}
			}
			return &utils.Result{}, nil
		},
	}

	mockProgramQuery := &osmanager.MoqProgramQuery{
		GetProgramVersionFunc: func(program string, versionExtractor osmanager.VersionExtractor, queryArgs ...string) (string, error) {
			return "4.14.0", nil
		},
	}

	mockEscalator := &privilege.MoqEscalator{
		EscalateCommandFunc: func(baseCmd string, baseArgs []string) (privilege.EscalationResult, error) {
			return privilege.EscalationResult{
				Command: "sudo",
				Args:    append([]string{baseCmd}, baseArgs...),
			}, nil
		},
	}

	// Create DNF package manager
	dnfManager := dnf.NewDnfPackageManager(
		logger.DefaultLogger,
		mockCommander,
		mockProgramQuery,
		mockEscalator,
		utils.DisplayModeProgress,
	)

	// Load package mappings
	v := viper.New()
	mappings, err := packageresolver.LoadPackageMappings(v, "")
	require.NoError(t, err, "Should load package mappings")

	// Create resolver
	resolver, err := packageresolver.NewResolver(mappings, dnfManager)
	require.NoError(t, err, "Should create resolver")

	// Test the critical flow: prerequisite name → resolved package → installation
	prerequisiteNames := []string{"development-tools", "procps-ng"}

	for _, prereqName := range prerequisiteNames {
		t.Run(fmt.Sprintf("InstallPrerequisite_%s", prereqName), func(t *testing.T) {
			// Step 1: Resolve prerequisite to get proper package info
			resolvedPackage, err := resolver.Resolve(prereqName, "")
			require.NoError(t, err, "Should resolve prerequisite %s", prereqName)

			// Step 2: Create RequestedPackageInfo with preserved type
			packageInfo := pkgmanager.NewRequestedPackageInfoWithType(
				resolvedPackage.Name,
				resolvedPackage.Type,
				resolvedPackage.VersionConstraints,
			)

			// Verify type preservation
			if prereqName == "development-tools" {
				require.Equal(t, "Development Tools", packageInfo.Name, "Should have resolved group name")
				require.Equal(t, "group", packageInfo.Type, "Should preserve group type")
			} else if prereqName == "procps-ng" {
				require.Equal(t, "procps-ng", packageInfo.Name, "Should have package name")
				require.Equal(t, "", packageInfo.Type, "Should have empty type for regular package")
			}

			// Step 3: Install package (this will call our mock commander)
			err = dnfManager.InstallPackage(packageInfo)
			require.NoError(t, err, "Should install package successfully")

			t.Logf("✅ Successfully installed %s with correct type handling", prereqName)
		})
	}
}

func Test_PrerequisiteInstallation_UsesCorrectDNFCommands(t *testing.T) {
	// This test specifically verifies that the right DNF commands are generated

	var capturedCommands [][]string

	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(command string, args []string, options ...utils.Option) (*utils.Result, error) {
			capturedCommands = append(capturedCommands, append([]string{command}, args...))
			return &utils.Result{}, nil
		},
	}

	mockProgramQuery := &osmanager.MoqProgramQuery{
		GetProgramVersionFunc: func(program string, versionExtractor osmanager.VersionExtractor, queryArgs ...string) (string, error) {
			return "4.14.0", nil
		},
	}

	mockEscalator := &privilege.MoqEscalator{
		EscalateCommandFunc: func(baseCmd string, baseArgs []string) (privilege.EscalationResult, error) {
			return privilege.EscalationResult{
				Command: "sudo",
				Args:    append([]string{baseCmd}, baseArgs...),
			}, nil
		},
	}

	dnfManager := dnf.NewDnfPackageManager(
		logger.DefaultLogger,
		mockCommander,
		mockProgramQuery,
		mockEscalator,
		utils.DisplayModeProgress,
	)

	// Test group package installation
	groupPackageInfo := pkgmanager.NewRequestedPackageInfoWithType("Development Tools", "group", nil)
	err := dnfManager.InstallPackage(groupPackageInfo)
	require.NoError(t, err)

	// Test regular package installation
	regularPackageInfo := pkgmanager.NewRequestedPackageInfoWithType("procps-ng", "", nil)
	err = dnfManager.InstallPackage(regularPackageInfo)
	require.NoError(t, err)

	// Verify the correct commands were generated
	require.Len(t, capturedCommands, 2, "Should have captured 2 commands")

	// Check group install command
	groupCmd := capturedCommands[0]
	require.Equal(t, []string{"sudo", "dnf", "group", "install", "-y", "Development Tools"}, groupCmd)
	t.Logf("✅ Group install command: %v", groupCmd)

	// Check regular install command
	regularCmd := capturedCommands[1]
	require.Equal(t, []string{"sudo", "dnf", "install", "-y", "procps-ng"}, regularCmd)
	t.Logf("✅ Regular install command: %v", regularCmd)
}
