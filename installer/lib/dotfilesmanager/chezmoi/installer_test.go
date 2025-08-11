package chezmoi_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/MrPointer/dotfiles/installer/lib/dotfilesmanager/chezmoi"
	"github.com/MrPointer/dotfiles/installer/lib/pkgmanager"
	"github.com/MrPointer/dotfiles/installer/utils"
	"github.com/MrPointer/dotfiles/installer/utils/httpclient"
	"github.com/MrPointer/dotfiles/installer/utils/logger"
	"github.com/MrPointer/dotfiles/installer/utils/osmanager"
)

func Test_Install_ReturnsEarly_WhenChezmoiAlreadyInstalled(t *testing.T) {
	mockFileSystem := &utils.MoqFileSystem{}
	mockUserManager := &osmanager.MoqUserManager{}
	mockCommander := &utils.MoqCommander{}
	mockPackageManager := &pkgmanager.MoqPackageManager{}
	mockHTTPClient := &httpclient.MoqHTTPClient{}

	chezmoiConfig := chezmoi.DefaultChezmoiConfig("/home/user/.config/chezmoi/chezmoi.toml", "/home/user/.local/share/chezmoi")
	manager := chezmoi.NewChezmoiManager(logger.DefaultLogger, mockFileSystem, mockUserManager, mockCommander, mockPackageManager, mockHTTPClient, chezmoiConfig)

	mockPackageManager.IsPackageInstalledFunc = func(packageInfo pkgmanager.PackageInfo) (bool, error) {
		require.Equal(t, "chezmoi", packageInfo.Name)
		return true, nil
	}

	err := manager.Install()

	require.NoError(t, err)
	require.Len(t, mockPackageManager.IsPackageInstalledCalls(), 1)
	require.Len(t, mockPackageManager.InstallPackageCalls(), 0)
	require.Len(t, mockHTTPClient.GetCalls(), 0)
	require.Len(t, mockCommander.RunCommandCalls(), 0)
}

func Test_Install_ReturnsError_WhenPackageInstalledCheckFails(t *testing.T) {
	mockFileSystem := &utils.MoqFileSystem{}
	mockUserManager := &osmanager.MoqUserManager{}
	mockCommander := &utils.MoqCommander{}
	mockPackageManager := &pkgmanager.MoqPackageManager{}
	mockHTTPClient := &httpclient.MoqHTTPClient{}

	chezmoiConfig := chezmoi.DefaultChezmoiConfig("/home/user/.config/chezmoi/chezmoi.toml", "/home/user/.local/share/chezmoi")
	manager := chezmoi.NewChezmoiManager(logger.DefaultLogger, mockFileSystem, mockUserManager, mockCommander, mockPackageManager, mockHTTPClient, chezmoiConfig)

	mockPackageManager.IsPackageInstalledFunc = func(packageInfo pkgmanager.PackageInfo) (bool, error) {
		return false, errors.New("package manager unavailable")
	}

	err := manager.Install()

	require.Error(t, err)
	require.Contains(t, err.Error(), "package manager unavailable")
	require.Len(t, mockPackageManager.IsPackageInstalledCalls(), 1)
}

func Test_Install_SucceedsWithPackageManager_WhenChezmoiNotInstalledAndPackageManagerWorks(t *testing.T) {
	mockFileSystem := &utils.MoqFileSystem{}
	mockUserManager := &osmanager.MoqUserManager{}
	mockCommander := &utils.MoqCommander{}
	mockPackageManager := &pkgmanager.MoqPackageManager{}
	mockHTTPClient := &httpclient.MoqHTTPClient{}

	chezmoiConfig := chezmoi.DefaultChezmoiConfig("/home/user/.config/chezmoi/chezmoi.toml", "/home/user/.local/share/chezmoi")
	manager := chezmoi.NewChezmoiManager(logger.DefaultLogger, mockFileSystem, mockUserManager, mockCommander, mockPackageManager, mockHTTPClient, chezmoiConfig)

	mockPackageManager.IsPackageInstalledFunc = func(packageInfo pkgmanager.PackageInfo) (bool, error) {
		return false, nil
	}

	mockPackageManager.InstallPackageFunc = func(packageInfo pkgmanager.RequestedPackageInfo) error {
		require.Equal(t, "chezmoi", packageInfo.Name)
		require.NotNil(t, packageInfo.VersionConstraints)
		return nil
	}

	err := manager.Install()

	require.NoError(t, err)
	require.Len(t, mockPackageManager.IsPackageInstalledCalls(), 1)
	require.Len(t, mockPackageManager.InstallPackageCalls(), 1)
	require.Len(t, mockHTTPClient.GetCalls(), 0)
	require.Len(t, mockCommander.RunCommandCalls(), 0)
}

func Test_Install_FallsBackToManualInstall_WhenPackageManagerFails(t *testing.T) {
	mockFileSystem := &utils.MoqFileSystem{}
	mockUserManager := &osmanager.MoqUserManager{}
	mockCommander := &utils.MoqCommander{}
	mockPackageManager := &pkgmanager.MoqPackageManager{}
	mockHTTPClient := &httpclient.MoqHTTPClient{}

	chezmoiConfig := chezmoi.DefaultChezmoiConfig("/home/user/.config/chezmoi/chezmoi.toml", "/home/user/.local/share/chezmoi")
	manager := chezmoi.NewChezmoiManager(logger.DefaultLogger, mockFileSystem, mockUserManager, mockCommander, mockPackageManager, mockHTTPClient, chezmoiConfig)

	mockPackageManager.IsPackageInstalledFunc = func(packageInfo pkgmanager.PackageInfo) (bool, error) {
		return false, nil
	}

	mockPackageManager.InstallPackageFunc = func(packageInfo pkgmanager.RequestedPackageInfo) error {
		return errors.New("package manager failed")
	}

	mockHTTPClient.GetFunc = func(url string) (*http.Response, error) {
		require.Equal(t, "get.chezmoi.io", url)
		response := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString("#!/bin/sh\necho 'Installing chezmoi...'")),
		}
		return response, nil
	}

	mockUserManager.GetHomeDirFunc = func() (string, error) {
		return "/home/user", nil
	}

	mockCommander.RunCommandFunc = func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
		require.Equal(t, "sh", name)
		require.Equal(t, []string{"-c", "#!/bin/sh\necho 'Installing chezmoi...'", "--", "-b", "/home/user/.local/bin"}, args)
		return &utils.Result{ExitCode: 0}, nil
	}

	err := manager.Install()

	require.NoError(t, err)
	require.Len(t, mockPackageManager.IsPackageInstalledCalls(), 1)
	require.Len(t, mockPackageManager.InstallPackageCalls(), 1)
	require.Len(t, mockHTTPClient.GetCalls(), 1)
	require.Len(t, mockUserManager.GetHomeDirCalls(), 1)
	require.Len(t, mockCommander.RunCommandCalls(), 1)
}

func Test_Install_ReturnsError_WhenHTTPRequestFails(t *testing.T) {
	mockFileSystem := &utils.MoqFileSystem{}
	mockUserManager := &osmanager.MoqUserManager{}
	mockCommander := &utils.MoqCommander{}
	mockPackageManager := &pkgmanager.MoqPackageManager{}
	mockHTTPClient := &httpclient.MoqHTTPClient{}

	chezmoiConfig := chezmoi.DefaultChezmoiConfig("/home/user/.config/chezmoi/chezmoi.toml", "/home/user/.local/share/chezmoi")
	manager := chezmoi.NewChezmoiManager(logger.DefaultLogger, mockFileSystem, mockUserManager, mockCommander, mockPackageManager, mockHTTPClient, chezmoiConfig)

	mockPackageManager.IsPackageInstalledFunc = func(packageInfo pkgmanager.PackageInfo) (bool, error) {
		return false, nil
	}

	mockPackageManager.InstallPackageFunc = func(packageInfo pkgmanager.RequestedPackageInfo) error {
		return errors.New("package manager failed")
	}

	mockHTTPClient.GetFunc = func(url string) (*http.Response, error) {
		return nil, errors.New("network unavailable")
	}

	err := manager.Install()

	require.Error(t, err)
	require.Contains(t, err.Error(), "network unavailable")
	require.Len(t, mockPackageManager.IsPackageInstalledCalls(), 1)
	require.Len(t, mockPackageManager.InstallPackageCalls(), 1)
	require.Len(t, mockHTTPClient.GetCalls(), 1)
}

func Test_Install_ReturnsError_WhenHTTPResponseIsNotOK(t *testing.T) {
	mockFileSystem := &utils.MoqFileSystem{}
	mockUserManager := &osmanager.MoqUserManager{}
	mockCommander := &utils.MoqCommander{}
	mockPackageManager := &pkgmanager.MoqPackageManager{}
	mockHTTPClient := &httpclient.MoqHTTPClient{}

	chezmoiConfig := chezmoi.DefaultChezmoiConfig("/home/user/.config/chezmoi/chezmoi.toml", "/home/user/.local/share/chezmoi")
	manager := chezmoi.NewChezmoiManager(logger.DefaultLogger, mockFileSystem, mockUserManager, mockCommander, mockPackageManager, mockHTTPClient, chezmoiConfig)

	mockPackageManager.IsPackageInstalledFunc = func(packageInfo pkgmanager.PackageInfo) (bool, error) {
		return false, nil
	}

	mockPackageManager.InstallPackageFunc = func(packageInfo pkgmanager.RequestedPackageInfo) error {
		return errors.New("package manager failed")
	}

	mockHTTPClient.GetFunc = func(url string) (*http.Response, error) {
		response := &http.Response{
			StatusCode: http.StatusNotFound,
			Status:     "404 Not Found",
			Body:       io.NopCloser(bytes.NewBufferString("")),
		}
		return response, nil
	}

	err := manager.Install()

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to download chezmoi binary")
	require.Contains(t, err.Error(), "404 Not Found")
}

func Test_Install_ReturnsError_WhenResponseBodyReadFails(t *testing.T) {
	mockFileSystem := &utils.MoqFileSystem{}
	mockUserManager := &osmanager.MoqUserManager{}
	mockCommander := &utils.MoqCommander{}
	mockPackageManager := &pkgmanager.MoqPackageManager{}
	mockHTTPClient := &httpclient.MoqHTTPClient{}

	chezmoiConfig := chezmoi.DefaultChezmoiConfig("/home/user/.config/chezmoi/chezmoi.toml", "/home/user/.local/share/chezmoi")
	manager := chezmoi.NewChezmoiManager(logger.DefaultLogger, mockFileSystem, mockUserManager, mockCommander, mockPackageManager, mockHTTPClient, chezmoiConfig)

	mockPackageManager.IsPackageInstalledFunc = func(packageInfo pkgmanager.PackageInfo) (bool, error) {
		return false, nil
	}

	mockPackageManager.InstallPackageFunc = func(packageInfo pkgmanager.RequestedPackageInfo) error {
		return errors.New("package manager failed")
	}

	failingReader := &FailingReader{}
	mockHTTPClient.GetFunc = func(url string) (*http.Response, error) {
		response := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(failingReader),
		}
		return response, nil
	}

	err := manager.Install()

	require.Error(t, err)
	require.Contains(t, err.Error(), "read failed")
}

func Test_Install_ReturnsError_WhenGetHomeDirFails(t *testing.T) {
	mockFileSystem := &utils.MoqFileSystem{}
	mockUserManager := &osmanager.MoqUserManager{}
	mockCommander := &utils.MoqCommander{}
	mockPackageManager := &pkgmanager.MoqPackageManager{}
	mockHTTPClient := &httpclient.MoqHTTPClient{}

	chezmoiConfig := chezmoi.DefaultChezmoiConfig("/home/user/.config/chezmoi/chezmoi.toml", "/home/user/.local/share/chezmoi")
	manager := chezmoi.NewChezmoiManager(logger.DefaultLogger, mockFileSystem, mockUserManager, mockCommander, mockPackageManager, mockHTTPClient, chezmoiConfig)

	mockPackageManager.IsPackageInstalledFunc = func(packageInfo pkgmanager.PackageInfo) (bool, error) {
		return false, nil
	}

	mockPackageManager.InstallPackageFunc = func(packageInfo pkgmanager.RequestedPackageInfo) error {
		return errors.New("package manager failed")
	}

	mockHTTPClient.GetFunc = func(url string) (*http.Response, error) {
		response := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString("#!/bin/sh\necho 'Installing chezmoi...'")),
		}
		return response, nil
	}

	mockUserManager.GetHomeDirFunc = func() (string, error) {
		return "", errors.New("home directory not available")
	}

	err := manager.Install()

	require.Error(t, err)
	require.Contains(t, err.Error(), "home directory not available")
}

func Test_Install_ReturnsError_WhenManualInstallCommandFails(t *testing.T) {
	mockFileSystem := &utils.MoqFileSystem{}
	mockUserManager := &osmanager.MoqUserManager{}
	mockCommander := &utils.MoqCommander{}
	mockPackageManager := &pkgmanager.MoqPackageManager{}
	mockHTTPClient := &httpclient.MoqHTTPClient{}

	chezmoiConfig := chezmoi.DefaultChezmoiConfig("/home/user/.config/chezmoi/chezmoi.toml", "/home/user/.local/share/chezmoi")
	manager := chezmoi.NewChezmoiManager(logger.DefaultLogger, mockFileSystem, mockUserManager, mockCommander, mockPackageManager, mockHTTPClient, chezmoiConfig)

	mockPackageManager.IsPackageInstalledFunc = func(packageInfo pkgmanager.PackageInfo) (bool, error) {
		return false, nil
	}

	mockPackageManager.InstallPackageFunc = func(packageInfo pkgmanager.RequestedPackageInfo) error {
		return errors.New("package manager failed")
	}

	mockHTTPClient.GetFunc = func(url string) (*http.Response, error) {
		response := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString("#!/bin/sh\necho 'Installing chezmoi...'")),
		}
		return response, nil
	}

	mockUserManager.GetHomeDirFunc = func() (string, error) {
		return "/home/user", nil
	}

	mockCommander.RunCommandFunc = func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
		return nil, errors.New("command execution failed")
	}

	err := manager.Install()

	require.Error(t, err)
	require.Contains(t, err.Error(), "command execution failed")
}

func Test_Install_ReturnsError_WhenManualInstallCommandExitsWithNonZeroCode(t *testing.T) {
	mockFileSystem := &utils.MoqFileSystem{}
	mockUserManager := &osmanager.MoqUserManager{}
	mockCommander := &utils.MoqCommander{}
	mockPackageManager := &pkgmanager.MoqPackageManager{}
	mockHTTPClient := &httpclient.MoqHTTPClient{}

	chezmoiConfig := chezmoi.DefaultChezmoiConfig("/home/user/.config/chezmoi/chezmoi.toml", "/home/user/.local/share/chezmoi")
	manager := chezmoi.NewChezmoiManager(logger.DefaultLogger, mockFileSystem, mockUserManager, mockCommander, mockPackageManager, mockHTTPClient, chezmoiConfig)

	mockPackageManager.IsPackageInstalledFunc = func(packageInfo pkgmanager.PackageInfo) (bool, error) {
		return false, nil
	}

	mockPackageManager.InstallPackageFunc = func(packageInfo pkgmanager.RequestedPackageInfo) error {
		return errors.New("package manager failed")
	}

	mockHTTPClient.GetFunc = func(url string) (*http.Response, error) {
		response := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString("#!/bin/sh\necho 'Installing chezmoi...'")),
		}
		return response, nil
	}

	mockUserManager.GetHomeDirFunc = func() (string, error) {
		return "/home/user", nil
	}

	mockCommander.RunCommandFunc = func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
		return &utils.Result{
			ExitCode: 1,
			Stderr:   []byte("installation script failed"),
		}, nil
	}

	err := manager.Install()

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to install chezmoi manually")
	require.Contains(t, err.Error(), "installation script failed")
}

func Test_Install_ClosesResponseBody(t *testing.T) {
	mockFileSystem := &utils.MoqFileSystem{}
	mockUserManager := &osmanager.MoqUserManager{}
	mockCommander := &utils.MoqCommander{}
	mockPackageManager := &pkgmanager.MoqPackageManager{}
	mockHTTPClient := &httpclient.MoqHTTPClient{}

	chezmoiConfig := chezmoi.DefaultChezmoiConfig("/home/user/.config/chezmoi/chezmoi.toml", "/home/user/.local/share/chezmoi")
	manager := chezmoi.NewChezmoiManager(logger.DefaultLogger, mockFileSystem, mockUserManager, mockCommander, mockPackageManager, mockHTTPClient, chezmoiConfig)

	mockPackageManager.IsPackageInstalledFunc = func(packageInfo pkgmanager.PackageInfo) (bool, error) {
		return false, nil
	}

	mockPackageManager.InstallPackageFunc = func(packageInfo pkgmanager.RequestedPackageInfo) error {
		return errors.New("package manager failed")
	}

	mockBody := &MockReadCloser{
		Reader: bytes.NewBufferString("#!/bin/sh\necho 'Installing chezmoi...'"),
	}

	mockHTTPClient.GetFunc = func(url string) (*http.Response, error) {
		response := &http.Response{
			StatusCode: http.StatusOK,
			Body:       mockBody,
		}
		return response, nil
	}

	mockUserManager.GetHomeDirFunc = func() (string, error) {
		return "/home/user", nil
	}

	mockCommander.RunCommandFunc = func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
		return &utils.Result{ExitCode: 0}, nil
	}

	err := manager.Install()

	require.NoError(t, err)
	require.True(t, mockBody.CloseCalled, "Response body should be closed")
}

// Helper types for testing
type FailingReader struct{}

func (fr *FailingReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("read failed")
}

type MockReadCloser struct {
	io.Reader
	CloseCalled bool
}

func (mrc *MockReadCloser) Close() error {
	mrc.CloseCalled = true
	return nil
}
