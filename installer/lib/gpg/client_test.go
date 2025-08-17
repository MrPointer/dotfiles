package gpg_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/MrPointer/dotfiles/installer/lib/gpg"
	"github.com/MrPointer/dotfiles/installer/utils"
	"github.com/MrPointer/dotfiles/installer/utils/logger"
	"github.com/MrPointer/dotfiles/installer/utils/osmanager"
)

func Test_NewDefaultGpgClient_ReturnsValidInstance(t *testing.T) {
	mockOsManager := &osmanager.MoqOsManager{}
	mockFilesystem := &utils.MoqFileSystem{}
	mockCommander := &utils.MoqCommander{}

	client := gpg.NewDefaultGpgClient(mockOsManager, mockFilesystem, mockCommander, logger.DefaultLogger)

	require.NotNil(t, client)
}

func Test_CreateKeyPair_ReturnsKeyID_WhenCommandSucceeds(t *testing.T) {
	testCases := []struct {
		name          string
		commandOutput string
		expectedKeyID string
		gpgTTY        string
	}{
		{
			name: "RSA 3072-bit key with GPG_TTY set",
			commandOutput: `gpg: key ABC123DEF456 marked as ultimately trusted
gpg: revocation certificate stored as '/home/user/.gnupg/openpgp-revocs.d/ABC123DEF456789.rev'
pub   rsa3072 2024-01-01 [SC]
      ABC123DEF456789
uid                      Test User <test@example.com>`,
			expectedKeyID: "ABC123DEF456",
			gpgTTY:        "/dev/pts/0",
		},
		{
			name: "RSA 4096-bit key with tty command fallback",
			commandOutput: `gpg: key XYZ789ABC123 marked as ultimately trusted
gpg: revocation certificate stored as '/home/user/.gnupg/openpgp-revocs.d/XYZ789ABC123456.rev'
pub   rsa4096 2024-01-02 [SC]
      XYZ789ABC123456
uid                      Another User <another@example.com>`,
			expectedKeyID: "XYZ789ABC123",
			gpgTTY:        "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockOsManager := &osmanager.MoqOsManager{
				GetenvFunc: func(key string) string {
					if key == "GPG_TTY" {
						return tc.gpgTTY
					}
					return ""
				},
			}
			mockFilesystem := &utils.MoqFileSystem{
				PathExistsFunc: func(path string) (bool, error) {
					return path == "/dev/tty", nil
				},
			}
			mockCommander := &utils.MoqCommander{
				RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
					if name == "tty" && tc.gpgTTY == "" {
						return &utils.Result{
							Stdout:   []byte("/dev/pts/1\n"),
							ExitCode: 0,
						}, nil
					}
					if name == "gpg" {
						return &utils.Result{
							Stdout:   []byte(tc.commandOutput),
							ExitCode: 0,
						}, nil
					}
					return &utils.Result{ExitCode: 0}, nil
				},
			}

			client := gpg.NewDefaultGpgClient(mockOsManager, mockFilesystem, mockCommander, logger.DefaultLogger)

			keyID, err := client.CreateKeyPair()

			require.NoError(t, err)
			require.Equal(t, tc.expectedKeyID, keyID)

			calls := mockCommander.RunCommandCalls()
			// Find GPG call
			var gpgCall bool
			for _, call := range calls {
				if call.Name == "gpg" && len(call.Args) > 0 && call.Args[0] == "--gen-key" {
					gpgCall = true
					require.Equal(t, []string{"--gen-key", "--pinentry-mode", "loopback", "--default-new-key-algo", "nistp256"}, call.Args)
					break
				}
			}
			require.True(t, gpgCall, "Expected GPG command call")
		})
	}
}

func Test_CreateKeyPair_ReturnsError_WhenTTYDetectionFails(t *testing.T) {
	mockOsManager := &osmanager.MoqOsManager{
		GetenvFunc: func(key string) string {
			return ""
		},
	}
	mockFilesystem := &utils.MoqFileSystem{
		PathExistsFunc: func(path string) (bool, error) {
			return path == "", errors.New("path not found")
		},
	}
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			return &utils.Result{ExitCode: 1}, errors.New("command failed")
		},
	}

	client := gpg.NewDefaultGpgClient(mockOsManager, mockFilesystem, mockCommander, logger.DefaultLogger)

	keyID, err := client.CreateKeyPair()
	require.Error(t, err)
	require.Empty(t, keyID)
	require.Contains(t, err.Error(), "unable to detect TTY")
}

func Test_CreateKeyPair_ExtractsKeyIDFromDifferentOutputFormats(t *testing.T) {
	testCases := []struct {
		name          string
		gpgOutput     string
		expectedKeyID string
	}{
		{
			name: "ultimately trusted pattern",
			gpgOutput: `gpg: directory '/home/user/.gnupg' created
gpg: keybox '/home/user/.gnupg/pubring.kbx' created
gpg: key ABC123DEF456 marked as ultimately trusted
gpg: directory '/home/user/.gnupg/openpgp-revocs.d' created`,
			expectedKeyID: "ABC123DEF456",
		},
		{
			name: "public key pattern",
			gpgOutput: `gpg: directory '/home/user/.gnupg' created
gpg: XYZ789ABC123: public key "Test User <test@example.com>" imported
gpg: Total number processed: 1`,
			expectedKeyID: "XYZ789ABC123",
		},
		{
			name: "pub line pattern with next line key",
			gpgOutput: `gpg: directory '/home/user/.gnupg' created
pub   nistp256 2024-01-01 [SC]
      DEF456GHI789
uid                      Test User <test@example.com>`,
			expectedKeyID: "DEF456GHI789",
		},
		{
			name: "revocation certificate pattern",
			gpgOutput: `gpg: directory '/home/user/.gnupg' created
gpg: keybox '/home/user/.gnupg/pubring.kbx' created
gpg: revocation certificate stored as '/home/user/.gnupg/openpgp-revocs.d/JKL012MNO345.rev'`,
			expectedKeyID: "JKL012MNO345",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockOsManager := &osmanager.MoqOsManager{
				GetenvFunc: func(key string) string {
					if key == "GPG_TTY" {
						return "/dev/pts/0"
					}
					return ""
				},
			}
			mockFileSystem := &utils.MoqFileSystem{}
			mockCommander := &utils.MoqCommander{
				RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
					if name == "gpg" {
						return &utils.Result{
							Stdout:   []byte(tc.gpgOutput),
							ExitCode: 0,
						}, nil
					}
					return &utils.Result{ExitCode: 0}, nil
				},
			}
			client := gpg.NewDefaultGpgClient(mockOsManager, mockFileSystem, mockCommander, logger.DefaultLogger)

			keyID, err := client.CreateKeyPair()
			require.NoError(t, err)
			require.Equal(t, tc.expectedKeyID, keyID)
		})
	}
}

func Test_CreateKeyPair_ReturnsError_WhenKeyIDCannotBeExtractedFromOutput(t *testing.T) {
	mockOsManager := &osmanager.MoqOsManager{
		GetenvFunc: func(key string) string {
			if key == "GPG_TTY" {
				return "/dev/pts/0"
			}
			return ""
		},
	}
	mockFileSystem := &utils.MoqFileSystem{}
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			if name == "gpg" {
				return &utils.Result{
					Stdout: []byte(`gpg: directory '/home/user/.gnupg' created
gpg: keybox '/home/user/.gnupg/pubring.kbx' created
Some random output without key information`),
					ExitCode: 0,
				}, nil
			}
			return &utils.Result{ExitCode: 0}, nil
		},
	}
	client := gpg.NewDefaultGpgClient(mockOsManager, mockFileSystem, mockCommander, logger.DefaultLogger)

	keyID, err := client.CreateKeyPair()
	require.Error(t, err)
	require.Empty(t, keyID)
	require.Contains(t, err.Error(), "could not find key ID in GPG output")
}

func Test_CreateKeyPair_ReturnsError_WhenGpgCommandExitsWithNonZeroCode(t *testing.T) {
	mockOsManager := &osmanager.MoqOsManager{
		GetenvFunc: func(key string) string {
			if key == "GPG_TTY" {
				return "/dev/pts/0"
			}
			return ""
		},
	}
	mockFilesystem := &utils.MoqFileSystem{}
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			if name == "gpg" {
				return &utils.Result{
					Stderr:   []byte("GPG error occurred"),
					ExitCode: 1,
				}, nil
			}
			return &utils.Result{ExitCode: 0}, nil
		},
	}

	client := gpg.NewDefaultGpgClient(mockOsManager, mockFilesystem, mockCommander, logger.DefaultLogger)

	keyID, err := client.CreateKeyPair()

	require.Error(t, err)
	require.Empty(t, keyID)
	require.Contains(t, err.Error(), "failed to create GPG key pair")
}

func Test_CreateKeyPair_ReturnsError_WhenOutputHasInsufficientLines(t *testing.T) {
	mockOsManager := &osmanager.MoqOsManager{
		GetenvFunc: func(key string) string {
			if key == "GPG_TTY" {
				return "/dev/pts/0"
			}
			return ""
		},
	}
	mockFilesystem := &utils.MoqFileSystem{}
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			if name == "gpg" {
				return &utils.Result{
					Stdout:   []byte("line1\nline2"),
					ExitCode: 0,
				}, nil
			}
			return &utils.Result{ExitCode: 0}, nil
		},
	}

	client := gpg.NewDefaultGpgClient(mockOsManager, mockFilesystem, mockCommander, logger.DefaultLogger)

	keyID, err := client.CreateKeyPair()

	require.Error(t, err)
	require.Empty(t, keyID)
	require.Contains(t, err.Error(), "failed to extract GPG key ID")
}

func Test_CreateKeyPair_ReturnsError_WhenKeyIDCannotBeExtracted(t *testing.T) {
	mockOsManager := &osmanager.MoqOsManager{
		GetenvFunc: func(key string) string {
			if key == "GPG_TTY" {
				return "/dev/pts/0"
			}
			return ""
		},
	}
	mockFilesystem := &utils.MoqFileSystem{}
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			if name == "gpg" {
				return &utils.Result{
					Stdout: []byte(`line1
line2
line3 without pub prefix`),
					ExitCode: 0,
				}, nil
			}
			return &utils.Result{ExitCode: 0}, nil
		},
	}

	client := gpg.NewDefaultGpgClient(mockOsManager, mockFilesystem, mockCommander, logger.DefaultLogger)

	keyID, err := client.CreateKeyPair()

	require.Error(t, err)
	require.Empty(t, keyID)
	require.Contains(t, err.Error(), "could not find key ID in GPG output")
}

func Test_ListAvailableKeys_ReturnsKeys_WhenKeysExist(t *testing.T) {
	testCases := []struct {
		name          string
		commandOutput string
		expectedKeys  []string
	}{
		{
			name: "single key",
			commandOutput: `sec   rsa3072/ABC123DEF456 2024-01-01 [SC]
uid                 [ultimate] Test User <test@example.com>
ssb   rsa3072/789GHI012JKL 2024-01-01 [E]`,
			expectedKeys: []string{"ABC123DEF456"},
		},
		{
			name: "multiple keys",
			commandOutput: `sec   rsa3072/ABC123DEF456 2024-01-01 [SC]
uid                 [ultimate] Test User <test@example.com>
ssb   rsa3072/789GHI012JKL 2024-01-01 [E]
sec   rsa4096/XYZ789ABC123 2024-01-02 [SC]
uid                 [ultimate] Another User <another@example.com>
ssb   rsa4096/456DEF789GHI 2024-01-02 [E]`,
			expectedKeys: []string{"ABC123DEF456", "XYZ789ABC123"},
		},
		{
			name: "key with different format",
			commandOutput: `sec   ed25519/FEDCBA987654 2024-01-01 [SC]
uid                   [ultimate] Ed User <ed@example.com>
ssb   cv25519/321FED654CBA 2024-01-01 [E]`,
			expectedKeys: []string{"FEDCBA987654"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockOsManager := &osmanager.MoqOsManager{}
			mockFilesystem := &utils.MoqFileSystem{}
			mockCommander := &utils.MoqCommander{
				RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
					return &utils.Result{
						Stdout:   []byte(tc.commandOutput),
						ExitCode: 0,
					}, nil
				},
			}

			client := gpg.NewDefaultGpgClient(mockOsManager, mockFilesystem, mockCommander, logger.DefaultLogger)

			keys, err := client.ListAvailableKeys()

			require.NoError(t, err)
			require.Equal(t, tc.expectedKeys, keys)

			calls := mockCommander.RunCommandCalls()
			require.Len(t, calls, 1)
			require.Equal(t, "gpg", calls[0].Name)
			require.Equal(t, []string{"--list-secret-keys", "--keyid-format", "LONG"}, calls[0].Args)
		})
	}
}

func Test_ListAvailableKeys_ReturnsNil_WhenNoKeysExist(t *testing.T) {
	mockOsManager := &osmanager.MoqOsManager{}
	mockFilesystem := &utils.MoqFileSystem{}
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			return &utils.Result{
				Stdout:   []byte("gpg: no secret keys found"),
				ExitCode: 0,
			}, nil
		},
	}

	client := gpg.NewDefaultGpgClient(mockOsManager, mockFilesystem, mockCommander, logger.DefaultLogger)

	keys, err := client.ListAvailableKeys()

	require.NoError(t, err)
	require.Nil(t, keys)
}

func Test_ListAvailableKeys_ReturnsError_WhenCommandFails(t *testing.T) {
	mockOsManager := &osmanager.MoqOsManager{}
	mockFilesystem := &utils.MoqFileSystem{}
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			return nil, errors.New("gpg command failed")
		},
	}

	client := gpg.NewDefaultGpgClient(mockOsManager, mockFilesystem, mockCommander, logger.DefaultLogger)

	keys, err := client.ListAvailableKeys()

	require.Error(t, err)
	require.Nil(t, keys)
	require.Contains(t, err.Error(), "gpg command failed")
}

func Test_ListAvailableKeys_HandlesIncompleteSecLines(t *testing.T) {
	mockOsManager := &osmanager.MoqOsManager{}
	mockFilesystem := &utils.MoqFileSystem{}
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			return &utils.Result{
				Stdout: []byte(`sec
sec   incomplete
sec   rsa3072/ABC123DEF456 2024-01-01 [SC]`),
				ExitCode: 0,
			}, nil
		},
	}

	client := gpg.NewDefaultGpgClient(mockOsManager, mockFilesystem, mockCommander, logger.DefaultLogger)

	keys, err := client.ListAvailableKeys()

	require.NoError(t, err)
	require.Equal(t, []string{"ABC123DEF456"}, keys)
}

func Test_KeysAvailable_ReturnsTrue_WhenKeysExist(t *testing.T) {
	mockOsManager := &osmanager.MoqOsManager{}
	mockFilesystem := &utils.MoqFileSystem{}
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			return &utils.Result{
				Stdout: []byte(`sec   rsa3072/ABC123DEF456 2024-01-01 [SC]
uid                 [ultimate] Test User <test@example.com>`),
				ExitCode: 0,
			}, nil
		},
	}

	client := gpg.NewDefaultGpgClient(mockOsManager, mockFilesystem, mockCommander, logger.DefaultLogger)

	available, err := client.KeysAvailable()

	require.NoError(t, err)
	require.True(t, available)
}

func Test_KeysAvailable_ReturnsFalse_WhenNoKeysExist(t *testing.T) {
	mockOsManager := &osmanager.MoqOsManager{}
	mockFilesystem := &utils.MoqFileSystem{}
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			return &utils.Result{
				Stdout:   []byte("gpg: no secret keys found"),
				ExitCode: 0,
			}, nil
		},
	}

	client := gpg.NewDefaultGpgClient(mockOsManager, mockFilesystem, mockCommander, logger.DefaultLogger)

	available, err := client.KeysAvailable()

	require.NoError(t, err)
	require.False(t, available)
}

func Test_KeysAvailable_ReturnsError_WhenListAvailableKeysFails(t *testing.T) {
	mockOsManager := &osmanager.MoqOsManager{}
	mockFilesystem := &utils.MoqFileSystem{}
	mockCommander := &utils.MoqCommander{
		RunCommandFunc: func(name string, args []string, opts ...utils.Option) (*utils.Result, error) {
			return nil, errors.New("gpg list command failed")
		},
	}

	client := gpg.NewDefaultGpgClient(mockOsManager, mockFilesystem, mockCommander, logger.DefaultLogger)

	available, err := client.KeysAvailable()

	require.Error(t, err)
	require.False(t, available)
	require.Contains(t, err.Error(), "gpg list command failed")
}
