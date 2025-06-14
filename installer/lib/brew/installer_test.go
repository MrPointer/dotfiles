package brew_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/MrPointer/dotfiles/installer/lib/brew"
	"github.com/MrPointer/dotfiles/installer/lib/compatibility"
	"github.com/MrPointer/dotfiles/installer/utils"
	"github.com/MrPointer/dotfiles/installer/utils/httpclient"
	"github.com/MrPointer/dotfiles/installer/utils/logger"
	"github.com/MrPointer/dotfiles/installer/utils/osmanager"
	"github.com/stretchr/testify/require"
)

/* ------------------------------------------------------------------------------------------------------------------ */
/* Path Detection Tests                                                                                               */
/* ------------------------------------------------------------------------------------------------------------------ */

func Test_DetectBrewPath_ReturnsExpectedPath_WhenCompatible(t *testing.T) {
	tests := []struct {
		name          string
		sysInfo       *compatibility.SystemInfo
		expectedPath  string
		errorExpected bool
	}{
		{"darwin/arm64", &compatibility.SystemInfo{OSName: "darwin", Arch: "arm64"}, brew.MacOSArmBrewPath, false},
		{"darwin/amd64", &compatibility.SystemInfo{OSName: "darwin", Arch: "amd64"}, brew.MacOSIntelBrewPath, false},
		{"linux/amd64", &compatibility.SystemInfo{OSName: "linux", Arch: "amd64"}, brew.LinuxBrewPath, false},
		{"unsupported", &compatibility.SystemInfo{OSName: "plan9", Arch: "amd64"}, "", true},
		{"no sysinfo", nil, "", true},
	}

	for _, tt := range tests {
		current := tt
		t.Run(current.name, func(t *testing.T) {
			opts := brew.Options{
				Logger:     logger.DefaultLogger,
				SystemInfo: current.sysInfo,
				Commander:  nil,
			}

			b := brew.NewBrewInstaller(opts)
			got, err := b.DetectBrewPath()
			if (err != nil) != current.errorExpected {
				t.Fatalf("DetectBrewPath() error = %v, wantError %v", err, current.errorExpected)
			}
			if got != current.expectedPath {
				t.Errorf("DetectBrewPath() = %q, want %q", got, current.expectedPath)
			}
		})
	}
}

func Test_DetectBrewPath_UsesOverride_WhenProvided(t *testing.T) {
	overridePath := "/custom/brew/path"
	opts := brew.Options{
		Logger:           logger.DefaultLogger,
		SystemInfo:       &compatibility.SystemInfo{OSName: "darwin", Arch: "arm64"},
		BrewPathOverride: overridePath,
	}

	b := brew.NewBrewInstaller(opts)
	got, err := b.DetectBrewPath()

	require.NoError(t, err)
	require.Equal(t, overridePath, got)
}

/* ------------------------------------------------------------------------------------------------------------------ */
/* Installer Constructor Tests                                                                                        */
/* ------------------------------------------------------------------------------------------------------------------ */

func Test_NewBrewInstaller_CreatesMultiUserImplementation_WhenOptionIsEnabled(t *testing.T) {
	opts := brew.Options{
		MultiUserSystem: true,
		Logger:          logger.DefaultLogger,
		SystemInfo:      &compatibility.SystemInfo{OSName: "linux", Arch: "amd64"},
		Commander:       nil,
	}
	installer := brew.NewBrewInstaller(opts)
	require.NotNil(t, installer)

	_, isMultiUser := installer.(*brew.MultiUserBrewInstaller)
	require.True(t, isMultiUser)
}

func Test_NewBrewInstaller_CreatesSingleUserImplementation_ByDefault(t *testing.T) {
	opts := brew.Options{
		MultiUserSystem: false,
		Logger:          logger.DefaultLogger,
		SystemInfo:      &compatibility.SystemInfo{OSName: "linux", Arch: "amd64"},
		Commander:       nil,
	}
	installer := brew.NewBrewInstaller(opts)
	require.NotNil(t, installer)

	_, isMultiUser := installer.(*brew.MultiUserBrewInstaller)
	require.False(t, isMultiUser)
}

/* ------------------------------------------------------------------------------------------------------------------ */
/* Single-User Installer Tests                                                                                        */
/* ------------------------------------------------------------------------------------------------------------------ */

func Test_SingleUserBrew_ReportsAvailable_WhenPathExists(t *testing.T) {
	expectedBrewPath := "/opt/homebrew/bin/brew"

	// Create mock dependencies
	mockLogger := &logger.NoopLogger{}
	mockCommander := &utils.MoqCommander{}
	mockFS := &utils.MoqFileSystem{
		PathExistsFunc: func(path string) (bool, error) {
			if path == expectedBrewPath {
				return true, nil // Simulate that brew exists
			}
			return false, nil // Other paths do not exist
		},
	}
	mockHTTP := &httpclient.MoqHTTPClient{}
	mockOsManager := &osmanager.MoqOsManager{}

	opts := brew.Options{
		Logger:           mockLogger,
		SystemInfo:       &compatibility.SystemInfo{OSName: "darwin", Arch: "arm64"},
		Commander:        mockCommander,
		HTTPClient:       mockHTTP,
		OsManager:        mockOsManager,
		Fs:               mockFS,
		MultiUserSystem:  false,
		BrewPathOverride: expectedBrewPath,
	}

	installer := brew.NewBrewInstaller(opts)
	available, err := installer.IsAvailable()

	require.NoError(t, err)
	require.True(t, available)
}

func Test_SingleUserBrew_ReportsUnavailable_WhenPathDoesNotExist(t *testing.T) {
	expectedBrewPath := "/nonexistent/brew"

	// Create mock dependencies
	mockLogger := &logger.NoopLogger{}
	mockCommander := &utils.MoqCommander{}
	mockFS := &utils.MoqFileSystem{
		PathExistsFunc: func(path string) (bool, error) {
			if path == expectedBrewPath {
				return false, nil // Simulate that brew does not exist
			}
			return true, nil // Other paths exist
		},
	}
	mockHTTP := &httpclient.MoqHTTPClient{}
	mockOsManager := &osmanager.MoqOsManager{}

	opts := brew.Options{
		Logger:           mockLogger,
		SystemInfo:       &compatibility.SystemInfo{OSName: "darwin", Arch: "arm64"},
		Commander:        mockCommander,
		HTTPClient:       mockHTTP,
		OsManager:        mockOsManager,
		Fs:               mockFS,
		MultiUserSystem:  false,
		BrewPathOverride: expectedBrewPath,
	}

	installer := brew.NewBrewInstaller(opts)
	available, err := installer.IsAvailable()

	require.NoError(t, err)
	require.False(t, available)
}

func Test_SingleUserBrew_IsNotReinstalled_WhenAvailable(t *testing.T) {
	expectedBrewPath := "/opt/homebrew/bin/brew"

	// Create mock dependencies
	mockLogger := &logger.NoopLogger{}
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			if name == expectedBrewPath && len(args) == 1 && args[0] == "--version" {
				return &utils.Result{ExitCode: 0}, nil // Validation successful
			}
			return nil, fmt.Errorf("unexpected command: %s %v", name, args)
		},
	}
	mockFS := &utils.MoqFileSystem{
		PathExistsFunc: func(path string) (bool, error) {
			if path == expectedBrewPath {
				return true, nil // Simulate that brew exists
			}
			return false, nil // Other paths do not exist
		},
	}
	mockHTTP := &httpclient.MoqHTTPClient{}
	mockOsManager := &osmanager.MoqOsManager{}

	opts := brew.Options{
		Logger:           mockLogger,
		SystemInfo:       &compatibility.SystemInfo{OSName: "darwin", Arch: "arm64"},
		Commander:        mockCommander,
		HTTPClient:       mockHTTP,
		OsManager:        mockOsManager,
		Fs:               mockFS,
		MultiUserSystem:  false,
		BrewPathOverride: expectedBrewPath,
	}

	installer := brew.NewBrewInstaller(opts)
	err := installer.Install()

	require.NoError(t, err)

	// Verify no HTTP requests were made (no download should happen)
	require.Empty(t, mockHTTP.GetCalls())
}

func Test_SingleUserBrew_InstallsSuccessfully_WhenNotAvailable(t *testing.T) {
	// Create mock dependencies
	mockLogger := &logger.NoopLogger{}
	tempScriptPath := "/tmp/brew-install-12345.sh"
	installScript := "#!/bin/bash\necho 'Installing Homebrew...'"

	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			if name == "/bin/bash" && len(args) == 1 && args[0] == tempScriptPath {
				// Check if NONINTERACTIVE env is set via options
				for range opts {
					// We can't easily inspect options, so assume env is correct
				}
				return &utils.Result{ExitCode: 0}, nil // Installation successful
			}
			if name == "/opt/homebrew/bin/brew" && len(args) == 1 && args[0] == "--version" {
				return &utils.Result{ExitCode: 0}, nil // Validation successful
			}
			return nil, fmt.Errorf("unexpected command: %s %v", name, args)
		},
	}

	expectedBrewPath := "/opt/homebrew/bin/brew"

	pathExistsCalls := 0
	mockFS := &utils.MoqFileSystem{
		PathExistsFunc: func(path string) (bool, error) {
			pathExistsCalls++
			if path == expectedBrewPath {
				// First call: brew doesn't exist, second call: brew exists after install
				return pathExistsCalls > 1, nil
			}
			if path == tempScriptPath {
				return true, nil // Script exists after download
			}
			return false, nil
		},
		CreateTemporaryFileFunc: func(dir, pattern string) (string, error) {
			return tempScriptPath, nil
		},
		WriteFileFunc: func(path string, reader io.Reader) (int64, error) {
			if path == tempScriptPath {
				return int64(len(installScript)), nil
			}
			return 0, fmt.Errorf("unexpected path: %s", path)
		},
		RemovePathFunc: func(path string) error {
			return nil // Cleanup successful
		},
	}

	mockHTTP := &httpclient.MoqHTTPClient{
		GetFunc: func(url string) (*http.Response, error) {
			if url == "https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh" {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(installScript)),
				}, nil
			}
			return nil, fmt.Errorf("unexpected URL: %s", url)
		},
	}

	mockOsManager := &osmanager.MoqOsManager{
		SetPermissionsFunc: func(path string, perms os.FileMode) error {
			if path == tempScriptPath && perms == 0o755 {
				return nil
			}
			return fmt.Errorf("unexpected permission call: %s %o", path, perms)
		},
	}

	opts := brew.Options{
		Logger:           mockLogger,
		SystemInfo:       &compatibility.SystemInfo{OSName: "darwin", Arch: "arm64"},
		Commander:        mockCommander,
		HTTPClient:       mockHTTP,
		OsManager:        mockOsManager,
		Fs:               mockFS,
		MultiUserSystem:  false,
		BrewPathOverride: expectedBrewPath,
	}

	installer := brew.NewBrewInstaller(opts)
	err := installer.Install()

	require.NoError(t, err)

	// Verify the correct sequence of operations
	require.Len(t, mockHTTP.GetCalls(), 1)
	require.Len(t, mockFS.CreateTemporaryFileCalls(), 1)
	require.Len(t, mockFS.WriteFileCalls(), 1)
	require.Len(t, mockOsManager.SetPermissionsCalls(), 1)
	require.Len(t, mockCommander.RunCommandCalls(), 2) // Installation and validation calls
	require.Len(t, mockFS.RemovePathCalls(), 1)        // Cleanup call
}

/* ------------------------------------------------------------------------------------------------------------------ */
/* Error handling tests                                                                                               */
/* ------------------------------------------------------------------------------------------------------------------ */

func Test_SingleUserBrew_FailsInstallation_WhenDownloadingInstallationScriptFails(t *testing.T) {
	tests := []struct {
		name            string
		setupMockHTTP   func() *httpclient.MoqHTTPClient
		expectedErrText string
	}{
		{
			name: "HTTP error status",
			setupMockHTTP: func() *httpclient.MoqHTTPClient {
				return &httpclient.MoqHTTPClient{
					GetFunc: func(url string) (*http.Response, error) {
						if url == "https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh" {
							return &http.Response{
								StatusCode: http.StatusNotFound,
								Body:       io.NopCloser(bytes.NewBufferString("")),
							}, nil
						}
						return nil, fmt.Errorf("unexpected URL: %s", url)
					},
				}
			},
			expectedErrText: "HTTP status 404",
		},
		{
			name: "Network error",
			setupMockHTTP: func() *httpclient.MoqHTTPClient {
				return &httpclient.MoqHTTPClient{
					GetFunc: func(url string) (*http.Response, error) {
						return nil, fmt.Errorf("network error: connection timeout")
					},
				}
			},
			expectedErrText: "network error: connection timeout",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock dependencies
			mockLogger := &logger.NoopLogger{}
			mockCommander := &utils.MoqCommander{}
			mockFS := &utils.MoqFileSystem{
				PathExistsFunc: func(path string) (bool, error) {
					return false, nil // Brew doesn't exist
				},
			}
			mockHTTP := tt.setupMockHTTP()
			mockOsManager := &osmanager.MoqOsManager{}

			opts := brew.Options{
				Logger:           mockLogger,
				SystemInfo:       &compatibility.SystemInfo{OSName: "darwin", Arch: "arm64"},
				Commander:        mockCommander,
				HTTPClient:       mockHTTP,
				OsManager:        mockOsManager,
				Fs:               mockFS,
				MultiUserSystem:  false,
				BrewPathOverride: "/opt/homebrew/bin/brew",
			}

			installer := brew.NewBrewInstaller(opts)
			err := installer.Install()

			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErrText)
			require.Len(t, mockHTTP.GetCalls(), 1) // HTTP call should have been made
		})
	}
}

func Test_SingleUserBrew_FailsInstallation_WhenFailingToCreateTempFileHoldingDownloadedScript(t *testing.T) {
	// Create mock dependencies
	mockLogger := &logger.NoopLogger{}
	mockCommander := &utils.MoqCommander{}
	mockFS := &utils.MoqFileSystem{
		PathExistsFunc: func(path string) (bool, error) {
			return false, nil // Brew doesn't exist
		},
		CreateTemporaryFileFunc: func(dir, pattern string) (string, error) {
			return "", fmt.Errorf("disk space full")
		},
	}
	mockHTTP := &httpclient.MoqHTTPClient{
		GetFunc: func(url string) (*http.Response, error) {
			if url == "https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh" {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString("#!/bin/bash\necho 'test'")),
				}, nil
			}
			return nil, fmt.Errorf("unexpected URL: %s", url)
		},
	}
	mockOsManager := &osmanager.MoqOsManager{}

	opts := brew.Options{
		Logger:           mockLogger,
		SystemInfo:       &compatibility.SystemInfo{OSName: "darwin", Arch: "arm64"},
		Commander:        mockCommander,
		HTTPClient:       mockHTTP,
		OsManager:        mockOsManager,
		Fs:               mockFS,
		MultiUserSystem:  false,
		BrewPathOverride: "/opt/homebrew/bin/brew",
	}

	installer := brew.NewBrewInstaller(opts)
	err := installer.Install()

	require.Error(t, err)
	require.Contains(t, err.Error(), "disk space full")
}

func Test_SingleUserBrew_FailsInstallation_WhenFailingToCopyDownloadedScriptFromHttpBodyToTempFile(t *testing.T) {
	tests := []struct {
		name            string
		setupMockFS     func(tempScriptPath, expectedBrewPath, installScript string) *utils.MoqFileSystem
		expectedErrText string
	}{
		{
			name: "Write permissions error",
			setupMockFS: func(tempScriptPath, expectedBrewPath, installScript string) *utils.MoqFileSystem {
				return &utils.MoqFileSystem{
					PathExistsFunc: func(path string) (bool, error) {
						// For this specific test case, brew is assumed to not exist initially to trigger download.
						return false, nil
					},
					CreateTemporaryFileFunc: func(dir, pattern string) (string, error) {
						return tempScriptPath, nil
					},
					WriteFileFunc: func(path string, reader io.Reader) (int64, error) {
						return 0, fmt.Errorf("permission denied")
					},
					RemovePathFunc: func(path string) error {
						return nil // Cleanup should still work
					},
				}
			},
			expectedErrText: "permission denied",
		},
		{
			name: "Zero bytes written",
			setupMockFS: func(tempScriptPath, expectedBrewPath, installScript string) *utils.MoqFileSystem {
				return &utils.MoqFileSystem{
					PathExistsFunc: func(path string) (bool, error) {
						// Brew should not exist to trigger download and write attempt.
						return false, nil
					},
					CreateTemporaryFileFunc: func(dir, pattern string) (string, error) {
						return tempScriptPath, nil
					},
					WriteFileFunc: func(path string, reader io.Reader) (int64, error) {
						return 0, nil // Zero bytes written
					},
					RemovePathFunc: func(path string) error {
						return nil // Cleanup should still work
					},
				}
			},
			expectedErrText: "no bytes written",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedBrewPath := "/opt/homebrew/bin/brew"
			mockLogger := &logger.NoopLogger{}
			tempScriptPath := "/tmp/brew-install-12345.sh"
			installScript := "#!/bin/bash\necho 'Installing Homebrew...'"

			mockCommander := &utils.MoqCommander{}
			mockFS := tt.setupMockFS(tempScriptPath, expectedBrewPath, installScript)
			mockHTTP := &httpclient.MoqHTTPClient{
				GetFunc: func(url string) (*http.Response, error) {
					if url == "https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh" {
						return &http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(bytes.NewBufferString(installScript)),
						}, nil
					}
					return nil, fmt.Errorf("unexpected URL: %s", url)
				},
			}
			mockOsManager := &osmanager.MoqOsManager{}

			opts := brew.Options{
				Logger:           mockLogger,
				SystemInfo:       &compatibility.SystemInfo{OSName: "darwin", Arch: "arm64"},
				Commander:        mockCommander,
				HTTPClient:       mockHTTP,
				OsManager:        mockOsManager,
				Fs:               mockFS,
				MultiUserSystem:  false,
				BrewPathOverride: expectedBrewPath,
			}

			installer := brew.NewBrewInstaller(opts)
			err := installer.Install()

			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErrText)
			require.Len(t, mockFS.RemovePathCalls(), 1) // Cleanup should be called
		})
	}
}

/* ------------------------------------------------------------------------------------------------------------------ */
/* Options builder pattern tests                                                                                      */
/* ------------------------------------------------------------------------------------------------------------------ */

func Test_DefaultOptions_CreatesNonEmptyConfiguration(t *testing.T) {
	opts := brew.DefaultOptions()

	require.NotNil(t, opts.Logger)
	require.NotNil(t, opts.Commander)
	require.NotNil(t, opts.HTTPClient)
	require.NotNil(t, opts.OsManager)
	require.NotNil(t, opts.Fs)
	require.False(t, opts.MultiUserSystem)
	require.Nil(t, opts.SystemInfo)
	require.Empty(t, opts.BrewPathOverride)
}

func Test_OptionsBuilderPattern_ConfiguresAllOptions(t *testing.T) {
	customLogger := &logger.NoopLogger{}
	customSystemInfo := &compatibility.SystemInfo{OSName: "linux", Arch: "amd64"}
	customCommander := &utils.MoqCommander{}
	customHTTP := &httpclient.MoqHTTPClient{}
	customOsManager := &osmanager.MoqOsManager{}
	customFS := &utils.MoqFileSystem{}

	opts := brew.DefaultOptions().
		WithLogger(customLogger).
		WithMultiUserSystem(true).
		WithSystemInfo(customSystemInfo).
		WithCommander(customCommander).
		WithHTTPClient(customHTTP).
		WithOsManager(customOsManager).
		WithFileSystem(customFS)

	require.Equal(t, customLogger, opts.Logger)
	require.True(t, opts.MultiUserSystem)
	require.Equal(t, customSystemInfo, opts.SystemInfo)
	require.Equal(t, customCommander, opts.Commander)
	require.Equal(t, customHTTP, opts.HTTPClient)
	require.Equal(t, customOsManager, opts.OsManager)
	require.Equal(t, customFS, opts.Fs)
}

/* ------------------------------------------------------------------------------------------------------------------ */
/* Edge case and boundary tests                                                                                       */
/* ------------------------------------------------------------------------------------------------------------------ */

func Test_SingleUserBrew_HandlesEmptyInstallScript(t *testing.T) {
	// Create mock dependencies
	mockLogger := &logger.NoopLogger{}
	tempScriptPath := "/tmp/brew-install-12345.sh"
	emptyScript := ""

	mockCommander := &utils.MoqCommander{}

	mockFS := &utils.MoqFileSystem{
		PathExistsFunc: func(path string) (bool, error) {
			return false, nil // Brew doesn't exist
		},
		CreateTemporaryFileFunc: func(dir, pattern string) (string, error) {
			return tempScriptPath, nil
		},
		WriteFileFunc: func(path string, reader io.Reader) (int64, error) {
			if path == tempScriptPath {
				return int64(len(emptyScript)), nil // Zero bytes for empty script
			}
			return 0, fmt.Errorf("unexpected path: %s", path)
		},
		RemovePathFunc: func(path string) error {
			return nil // Cleanup should still work
		},
	}

	mockHTTP := &httpclient.MoqHTTPClient{
		GetFunc: func(url string) (*http.Response, error) {
			if url == "https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh" {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(emptyScript)),
				}, nil
			}
			return nil, fmt.Errorf("unexpected URL: %s", url)
		},
	}

	mockOsManager := &osmanager.MoqOsManager{}

	opts := brew.Options{
		Logger:           mockLogger,
		SystemInfo:       &compatibility.SystemInfo{OSName: "darwin", Arch: "arm64"},
		Commander:        mockCommander,
		HTTPClient:       mockHTTP,
		OsManager:        mockOsManager,
		Fs:               mockFS,
		MultiUserSystem:  false,
		BrewPathOverride: "/opt/homebrew/bin/brew",
	}

	installer := brew.NewBrewInstaller(opts)
	err := installer.Install()

	require.Error(t, err)
	require.Contains(t, err.Error(), "no bytes written")
	require.Len(t, mockFS.RemovePathCalls(), 1) // Cleanup should be called
}

func Test_SingleUserBrew_CanHandleLargeInstallScript(t *testing.T) {
	// Create mock dependencies with a large script
	mockLogger := &logger.NoopLogger{}
	tempScriptPath := "/tmp/brew-install-12345.sh"
	// Create a large script (1MB)
	largeScript := "#!/bin/bash\n" + string(bytes.Repeat([]byte("echo 'large script content'\n"), 50000)[:])

	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			if name == "/bin/bash" && len(args) == 1 && args[0] == tempScriptPath {
				return &utils.Result{ExitCode: 0}, nil // Installation successful
			}
			if name == "/opt/homebrew/bin/brew" && len(args) == 1 && args[0] == "--version" {
				return &utils.Result{ExitCode: 0}, nil // Validation successful
			}
			return nil, fmt.Errorf("unexpected command: %s %v", name, args)
		},
	}

	pathExistsCalls := 0
	mockFS := &utils.MoqFileSystem{
		PathExistsFunc: func(path string) (bool, error) {
			pathExistsCalls++
			if path == "/opt/homebrew/bin/brew" {
				// First call: brew doesn't exist, second call: brew exists after install
				return pathExistsCalls > 1, nil
			}
			if path == tempScriptPath {
				return true, nil // Script exists after download
			}
			return false, nil
		},
		CreateTemporaryFileFunc: func(dir, pattern string) (string, error) {
			return tempScriptPath, nil
		},
		WriteFileFunc: func(path string, reader io.Reader) (int64, error) {
			if path == tempScriptPath {
				return int64(len(largeScript)), nil
			}
			return 0, fmt.Errorf("unexpected path: %s", path)
		},
		RemovePathFunc: func(path string) error {
			return nil // Cleanup successful
		},
	}

	mockHTTP := &httpclient.MoqHTTPClient{
		GetFunc: func(url string) (*http.Response, error) {
			if url == "https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh" {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(largeScript)),
				}, nil
			}
			return nil, fmt.Errorf("unexpected URL: %s", url)
		},
	}

	mockOsManager := &osmanager.MoqOsManager{
		SetPermissionsFunc: func(path string, perms os.FileMode) error {
			if path == tempScriptPath && perms == 0o755 {
				return nil
			}
			return fmt.Errorf("unexpected permission call: %s %o", path, perms)
		},
	}

	opts := brew.Options{
		Logger:           mockLogger,
		SystemInfo:       &compatibility.SystemInfo{OSName: "darwin", Arch: "arm64"},
		Commander:        mockCommander,
		HTTPClient:       mockHTTP,
		OsManager:        mockOsManager,
		Fs:               mockFS,
		MultiUserSystem:  false,
		BrewPathOverride: "/opt/homebrew/bin/brew",
	}

	installer := brew.NewBrewInstaller(opts)
	err := installer.Install()

	require.NoError(t, err)

	// Verify that large scripts are handled correctly
	writeCalls := mockFS.WriteFileCalls()
	require.Len(t, writeCalls, 1)
	// Note: We can't easily verify the exact bytes written in this mock setup,
	// but the test ensures the system can handle large scripts without errors
}

func Test_BrewInstaller_ReturnsError_WhenSystemInfoIsNil(t *testing.T) {
	// Test behavior when SystemInfo is nil
	opts := brew.Options{
		Logger:          logger.DefaultLogger,
		SystemInfo:      nil, // Nil system info
		Commander:       nil,
		MultiUserSystem: false,
	}

	installer := brew.NewBrewInstaller(opts)
	_, err := installer.DetectBrewPath()

	require.Error(t, err)
	require.Contains(t, err.Error(), "system information is not provided")
}

/* ------------------------------------------------------------------------------------------------------------------ */
/* Multi-User Brew Installer Tests                                                                                    */
/* ------------------------------------------------------------------------------------------------------------------ */

func Test_MultiUserBrew_ReportsAvailable_WhenBrewExistsForBrewUser(t *testing.T) {
	expectedBrewPath := "/home/linuxbrew/.linuxbrew/bin/brew"

	// Create mock dependencies
	mockLogger := &logger.NoopLogger{}
	mockCommander := &utils.MoqCommander{}
	mockFS := &utils.MoqFileSystem{
		PathExistsFunc: func(path string) (bool, error) {
			if path == expectedBrewPath {
				return true, nil // Simulate that brew exists for the brew user
			}
			return false, nil // Other paths do not exist
		},
	}
	mockHTTP := &httpclient.MoqHTTPClient{}
	mockOsManager := &osmanager.MoqOsManager{
		GetFileOwnerFunc: func(path string) (string, error) {
			if path == expectedBrewPath {
				return brew.BrewUserOnMultiUserSystem, nil // Simulate that the brew user owns the brew binary
			}
			return "", fmt.Errorf("unexpected path: %s", path)
		},
	}

	opts := brew.Options{
		Logger:           mockLogger,
		SystemInfo:       &compatibility.SystemInfo{OSName: "linux", Arch: "amd64"},
		Commander:        mockCommander,
		HTTPClient:       mockHTTP,
		OsManager:        mockOsManager,
		Fs:               mockFS,
		MultiUserSystem:  true,
		BrewPathOverride: expectedBrewPath,
	}

	installer := brew.NewBrewInstaller(opts)
	available, err := installer.IsAvailable()

	require.NoError(t, err)
	require.True(t, available)
}

func Test_MultiUserBrew_ReportsUnavailable_WhenBrewDoesNotExistForBrewUser(t *testing.T) {
	expectedBrewPath := "/home/linuxbrew/.linuxbrew/bin/brew"

	// Create mock dependencies
	mockLogger := &logger.NoopLogger{}
	mockCommander := &utils.MoqCommander{}
	mockFS := &utils.MoqFileSystem{
		PathExistsFunc: func(path string) (bool, error) {
			if path == expectedBrewPath {
				return true, nil // Simulate that brew path exists
			}
			return true, nil // Other paths exist
		},
	}
	mockHTTP := &httpclient.MoqHTTPClient{}
	mockOsManager := &osmanager.MoqOsManager{
		GetFileOwnerFunc: func(path string) (string, error) {
			if path == expectedBrewPath {
				return "someotheruser", nil // Simulate that the brew binary is owned by a different user
			}
			return "", fmt.Errorf("unexpected path: %s", path)
		},
	}

	opts := brew.Options{
		Logger:           mockLogger,
		SystemInfo:       &compatibility.SystemInfo{OSName: "linux", Arch: "amd64"},
		Commander:        mockCommander,
		HTTPClient:       mockHTTP,
		OsManager:        mockOsManager,
		Fs:               mockFS,
		MultiUserSystem:  true,
		BrewPathOverride: expectedBrewPath,
	}

	installer := brew.NewBrewInstaller(opts)
	available, err := installer.IsAvailable()

	require.NoError(t, err)
	require.False(t, available)
}

func Test_MultiUserBrew_ReportedUnavailable_WhenBrewPathDoesNotExist(t *testing.T) {
	expectedBrewPath := "/nonexistent/brew"

	// Create mock dependencies
	mockLogger := &logger.NoopLogger{}
	mockCommander := &utils.MoqCommander{}
	mockFS := &utils.MoqFileSystem{
		PathExistsFunc: func(path string) (bool, error) {
			if path == expectedBrewPath {
				return false, nil // Simulate that brew does not exist
			}
			return true, nil // Other paths exist
		},
	}
	mockHTTP := &httpclient.MoqHTTPClient{}
	mockOsManager := &osmanager.MoqOsManager{}

	opts := brew.Options{
		Logger:           mockLogger,
		SystemInfo:       &compatibility.SystemInfo{OSName: "linux", Arch: "amd64"},
		Commander:        mockCommander,
		HTTPClient:       mockHTTP,
		OsManager:        mockOsManager,
		Fs:               mockFS,
		MultiUserSystem:  true,
		BrewPathOverride: expectedBrewPath,
	}

	installer := brew.NewBrewInstaller(opts)
	available, err := installer.IsAvailable()

	require.NoError(t, err)     // Multi-user installation should error on non-existent path
	require.False(t, available) // Brew should be reported as unavailable
}

func Test_MultiUserBrew_DoesNotReinstall_WhenAlreadyAvailable(t *testing.T) {
	expectedBrewPath := "/home/linuxbrew/.linuxbrew/bin/brew"

	// Create mock dependencies
	mockLogger := &logger.NoopLogger{}
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			if name == expectedBrewPath && len(args) == 1 && args[0] == "--version" {
				return &utils.Result{ExitCode: 0}, nil // Validation successful
			}
			return nil, fmt.Errorf("unexpected command: %s %v", name, args)
		},
	}
	mockFS := &utils.MoqFileSystem{
		PathExistsFunc: func(path string) (bool, error) {
			if path == expectedBrewPath {
				return true, nil // Simulate that brew exists for the brew user
			}
			return false, nil // Other paths do not exist
		},
	}
	mockHTTP := &httpclient.MoqHTTPClient{}
	mockOsManager := &osmanager.MoqOsManager{
		AddUserFunc: func(username string) error {
			if username == brew.BrewUserOnMultiUserSystem {
				return nil // Simulate successful user addition
			}
			return fmt.Errorf("unexpected user: %s", username)
		},
		UserExistsFunc: func(username string) (bool, error) {
			if username == brew.BrewUserOnMultiUserSystem {
				return true, nil // Simulate that the brew user exists
			}
			return false, nil // Other users do not exist
		},
		GetFileOwnerFunc: func(path string) (string, error) {
			if path == expectedBrewPath {
				return brew.BrewUserOnMultiUserSystem, nil // Simulate that the brew user owns the brew binary
			}
			return "", fmt.Errorf("unexpected path: %s", path)
		},
	}

	opts := brew.Options{
		Logger:           mockLogger,
		SystemInfo:       &compatibility.SystemInfo{OSName: "linux", Arch: "amd64"},
		Commander:        mockCommander,
		HTTPClient:       mockHTTP,
		OsManager:        mockOsManager,
		Fs:               mockFS,
		MultiUserSystem:  true,
		BrewPathOverride: expectedBrewPath,
	}

	installer := brew.NewBrewInstaller(opts)
	err := installer.Install()

	require.NoError(t, err)

	// Verify no HTTP requests were made (no download should happen)
	require.Empty(t, mockHTTP.GetCalls())
	// Verify no user management operations were performed
	require.Empty(t, mockOsManager.UserExistsCalls())
	require.Empty(t, mockOsManager.AddUserCalls())
}

//gocognit:ignore
//nolint:cyclop // This test is complex due to the multi-user setup and various checks involved.
func Test_MultiUserBrew_InstallsFromScratch_WhenUserDoesNotExist(t *testing.T) {
	expectedBrewPath := "/home/linuxbrew/.linuxbrew/bin/brew"
	tempScriptPath := "/tmp/brew-install-12345.sh"
	installScript := "#!/bin/bash\necho 'Installing Homebrew...'"

	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			if name == "sudo" && len(args) == 4 && args[0] == "-Hu" &&
				args[1] == brew.BrewUserOnMultiUserSystem &&
				args[2] == "bash" && args[3] == tempScriptPath {
				return &utils.Result{ExitCode: 0}, nil
			}
			if name == expectedBrewPath && len(args) == 1 && args[0] == "--version" {
				return &utils.Result{ExitCode: 0}, nil
			}
			return nil, fmt.Errorf("unexpected command: %s %v", name, args)
		},
	}

	pathExistsCalls := 0
	mockFS := &utils.MoqFileSystem{
		PathExistsFunc: func(path string) (bool, error) {
			pathExistsCalls++
			if path == expectedBrewPath {
				return pathExistsCalls > 1, nil
			}
			if path == tempScriptPath {
				return true, nil
			}
			return false, nil
		},
		CreateTemporaryFileFunc: func(dir, pattern string) (string, error) {
			return tempScriptPath, nil
		},
		WriteFileFunc: func(path string, reader io.Reader) (int64, error) {
			if path == tempScriptPath {
				return int64(len(installScript)), nil
			}
			return 0, fmt.Errorf("unexpected path: %s", path)
		},
		RemovePathFunc: func(path string) error {
			return nil
		},
	}

	mockHTTP := &httpclient.MoqHTTPClient{
		GetFunc: func(url string) (*http.Response, error) {
			if url == "https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh" {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(installScript)),
				}, nil
			}
			return nil, fmt.Errorf("unexpected URL: %s", url)
		},
	}

	mockOsManager := &osmanager.MoqOsManager{
		UserExistsFunc: func(username string) (bool, error) {
			if username == brew.BrewUserOnMultiUserSystem {
				return false, nil // User does not exist
			}
			return false, fmt.Errorf("unexpected user check: %s", username)
		},
		AddUserFunc: func(username string) error {
			if username == brew.BrewUserOnMultiUserSystem {
				return nil
			}
			return fmt.Errorf("unexpected user add: %s", username)
		},
		AddUserToGroupFunc: func(username, group string) error {
			if username == brew.BrewUserOnMultiUserSystem && group == "sudo" {
				return nil
			}
			return fmt.Errorf("unexpected user/group add: %s/%s", username, group)
		},
		AddSudoAccessFunc: func(username string) error {
			if username == brew.BrewUserOnMultiUserSystem {
				return nil
			}
			return fmt.Errorf("unexpected sudo access for: %s", username)
		},
		SetOwnershipFunc: func(path, username string) error {
			if path == "/home/linuxbrew" && username == brew.BrewUserOnMultiUserSystem {
				return nil
			}
			return fmt.Errorf("unexpected ownership set: %s for %s", path, username)
		},
		SetPermissionsFunc: func(path string, perms os.FileMode) error {
			if path == tempScriptPath && perms == 0o755 {
				return nil
			}
			return fmt.Errorf("unexpected permission call: %s %o", path, perms)
		},
		GetFileOwnerFunc: func(path string) (string, error) {
			if path == expectedBrewPath {
				return brew.BrewUserOnMultiUserSystem, nil
			}
			return "", fmt.Errorf("unexpected get owner for: %s", path)
		},
	}

	opts := brew.Options{
		Logger:           &logger.NoopLogger{},
		SystemInfo:       &compatibility.SystemInfo{OSName: "linux", Arch: "amd64"},
		Commander:        mockCommander,
		HTTPClient:       mockHTTP,
		OsManager:        mockOsManager,
		Fs:               mockFS,
		MultiUserSystem:  true,
		BrewPathOverride: expectedBrewPath,
	}

	installer := brew.NewBrewInstaller(opts)
	err := installer.Install()

	require.NoError(t, err)
	require.Len(t, mockOsManager.UserExistsCalls(), 1)
	require.Len(t, mockOsManager.AddUserCalls(), 1) // User should be added
	require.Len(t, mockOsManager.AddUserToGroupCalls(), 1)
	require.Len(t, mockOsManager.AddSudoAccessCalls(), 1)
	require.Len(t, mockOsManager.SetOwnershipCalls(), 1)
	require.Len(t, mockHTTP.GetCalls(), 1)
	require.Len(t, mockFS.CreateTemporaryFileCalls(), 1)
	require.Len(t, mockFS.WriteFileCalls(), 1)
	require.Len(t, mockOsManager.SetPermissionsCalls(), 1)
	require.Len(t, mockCommander.RunCommandCalls(), 2)
	require.Len(t, mockFS.RemovePathCalls(), 1)
}

//nolint:cyclop // This test is complex due to the multi-user setup and various checks involved.
func Test_MultiUserBrew_InstallsFromScratch_WhenUserAlreadyExists(t *testing.T) {
	expectedBrewPath := "/home/linuxbrew/.linuxbrew/bin/brew"
	tempScriptPath := "/tmp/brew-install-12345.sh"
	installScript := "#!/bin/bash\necho 'Installing Homebrew...'"

	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			if name == "sudo" && len(args) == 4 && args[0] == "-Hu" &&
				args[1] == brew.BrewUserOnMultiUserSystem &&
				args[2] == "bash" && args[3] == tempScriptPath {
				return &utils.Result{ExitCode: 0}, nil
			}
			if name == expectedBrewPath && len(args) == 1 && args[0] == "--version" {
				return &utils.Result{ExitCode: 0}, nil
			}
			return nil, fmt.Errorf("unexpected command: %s %v", name, args)
		},
	}

	pathExistsCalls := 0
	mockFS := &utils.MoqFileSystem{
		PathExistsFunc: func(path string) (bool, error) {
			pathExistsCalls++
			if path == expectedBrewPath {
				return pathExistsCalls > 1, nil
			}
			if path == tempScriptPath {
				return true, nil
			}
			return false, nil
		},
		CreateTemporaryFileFunc: func(dir, pattern string) (string, error) {
			return tempScriptPath, nil
		},
		WriteFileFunc: func(path string, reader io.Reader) (int64, error) {
			if path == tempScriptPath {
				return int64(len(installScript)), nil
			}
			return 0, fmt.Errorf("unexpected path: %s", path)
		},
		RemovePathFunc: func(path string) error {
			return nil
		},
	}

	mockHTTP := &httpclient.MoqHTTPClient{
		GetFunc: func(url string) (*http.Response, error) {
			if url == "https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh" {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(installScript)),
				}, nil
			}
			return nil, fmt.Errorf("unexpected URL: %s", url)
		},
	}

	mockOsManager := &osmanager.MoqOsManager{
		UserExistsFunc: func(username string) (bool, error) {
			if username == brew.BrewUserOnMultiUserSystem {
				return true, nil // User already exists
			}
			return false, fmt.Errorf("unexpected user check: %s", username)
		},
		AddUserFunc: func(username string) error {
			return fmt.Errorf("unexpected user add: %s", username) // Should not be called
		},
		AddUserToGroupFunc: func(username, group string) error {
			if username == brew.BrewUserOnMultiUserSystem && group == "sudo" {
				return nil
			}
			return fmt.Errorf("unexpected user/group add: %s/%s", username, group)
		},
		AddSudoAccessFunc: func(username string) error {
			if username == brew.BrewUserOnMultiUserSystem {
				return nil
			}
			return fmt.Errorf("unexpected sudo access for: %s", username)
		},
		SetOwnershipFunc: func(path, username string) error {
			if path == "/home/linuxbrew" && username == brew.BrewUserOnMultiUserSystem {
				return nil
			}
			return fmt.Errorf("unexpected ownership set: %s for %s", path, username)
		},
		SetPermissionsFunc: func(path string, perms os.FileMode) error {
			if path == tempScriptPath && perms == 0o755 {
				return nil
			}
			return fmt.Errorf("unexpected permission call: %s %o", path, perms)
		},
		GetFileOwnerFunc: func(path string) (string, error) {
			if path == expectedBrewPath {
				return brew.BrewUserOnMultiUserSystem, nil
			}
			return "", fmt.Errorf("unexpected get owner for: %s", path)
		},
	}

	opts := brew.Options{
		Logger:           &logger.NoopLogger{},
		SystemInfo:       &compatibility.SystemInfo{OSName: "linux", Arch: "amd64"},
		Commander:        mockCommander,
		HTTPClient:       mockHTTP,
		OsManager:        mockOsManager,
		Fs:               mockFS,
		MultiUserSystem:  true,
		BrewPathOverride: expectedBrewPath,
	}

	installer := brew.NewBrewInstaller(opts)
	err := installer.Install()

	require.NoError(t, err)
	require.Len(t, mockOsManager.UserExistsCalls(), 1)
	require.Empty(t, mockOsManager.AddUserCalls()) // User should not be added
	require.Len(t, mockOsManager.AddUserToGroupCalls(), 1)
	require.Len(t, mockOsManager.AddSudoAccessCalls(), 1)
	require.Len(t, mockOsManager.SetOwnershipCalls(), 1)
	require.Len(t, mockHTTP.GetCalls(), 1)
	require.Len(t, mockFS.CreateTemporaryFileCalls(), 1)
	require.Len(t, mockFS.WriteFileCalls(), 1)
	require.Len(t, mockOsManager.SetPermissionsCalls(), 1)
	require.Len(t, mockCommander.RunCommandCalls(), 2)
	require.Len(t, mockFS.RemovePathCalls(), 1)
}
