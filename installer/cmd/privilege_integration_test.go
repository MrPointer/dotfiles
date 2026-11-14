package cmd

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/MrPointer/dotfiles/installer/lib/compatibility"
	"github.com/MrPointer/dotfiles/installer/lib/pkgmanager"
	"github.com/stretchr/testify/require"
)

// Test_InstallingSystemChanges_Should_RequireAuthorization_When_RunAsNormalUser
// is driven by tests/privilege/test.exp on a disposable Linux machine.
func Test_InstallingSystemChanges_Should_RequireAuthorization_When_RunAsNormalUser(t *testing.T) {
	work := os.Getenv("INSTALLER_PRIVILEGE_TEST_DIR")
	if testing.Short() || work == "" {
		t.Skip("requires the disposable installer privilege test environment")
	}

	require.Equal(t, "linux", runtime.GOOS)
	marker, err := os.ReadFile(filepath.Join(work, "ready"))
	require.NoError(t, err)
	require.Equal(t, "installer-privilege-tests\n", string(marker))
	require.NotZero(t, os.Geteuid(), "the installer must run as a normal user")
	terminal, terminalErr := os.Open("/dev/tty")
	require.NoError(t, terminalErr)
	require.NoError(t, terminal.Close())

	// Use the CLI's real display policy, logger, commander, and OS manager.
	// Do not pause progress or attach stdin on behalf of production callers.
	plainFlag = false
	nonInteractive = false
	initLogger()
	t.Cleanup(func() { require.NoError(t, cliLogger.Close()) })
	initCommander()
	initOsManager()

	switch os.Getenv("INSTALLER_PRIVILEGE_TEST_CALLER") {
	case "shell":
		shellName = filepath.Base(os.Getenv("INSTALLER_PRIVILEGE_TEST_SHELL"))
		shellSource = "system"
		globalSysInfo = compatibility.SystemInfo{OSName: "linux", DistroName: "debian"}
		err = installShell(cliLogger)
	case "apt":
		require.NoError(t, cliLogger.StartProgress("Installing privilege regression package"))
		manager := createNativePackageManager("linux", "debian")
		err = manager.InstallPackage(pkgmanager.RequestedPackageInfo{Name: filepath.Join(work, "protected.deb")})
		if err != nil {
			require.NoError(t, cliLogger.FailProgress("Package installation failed", err))
		} else {
			require.NoError(t, cliLogger.FinishProgress("Package installation completed"))
		}
	default:
		t.Fatal("unknown privilege test caller")
	}

	if os.Getenv("INSTALLER_PRIVILEGE_TEST_SCENARIO") == "reject" {
		require.Error(t, err)
		if os.Getenv("INSTALLER_PRIVILEGE_TEST_CALLER") == "apt" {
			require.Contains(t, err.Error(), "failed to update package list")
		} else {
			require.Contains(t, err.Error(), "failed to add shell to /etc/shells")
		}
		t.Log("AUTHORIZATION_REJECTED")
		return
	}

	require.NoError(t, err)
	t.Log("AUTHORIZED_OPERATION_COMPLETED")
}
