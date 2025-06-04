package compatibility_test

import (
	"testing"

	"github.com/MrPointer/dotfiles/installer/lib/compatibility"
	"github.com/spf13/viper"
)

// MockOSDetector implements OSDetector for testing.
type MockOSDetector struct {
	osName     string
	distroName string
}

// GetOSName returns the mocked OS name.
func (m *MockOSDetector) GetOSName() string {
	return m.osName
}

// GetDistroName returns the mocked distro name.
func (m *MockOSDetector) GetDistroName() string {
	return m.distroName
}

// DetectSystem implements the OSDetector interface for the mock.
func (m *MockOSDetector) DetectSystem() (compatibility.SystemInfo, error) {
	return compatibility.SystemInfo{
		OSName:     m.osName,
		DistroName: m.distroName,
		Arch:       "amd64", // Mock architecture
	}, nil
}

// createMockConfig creates a mock compatibility configuration for testing.
func createMockConfig() *compatibility.CompatibilityConfig {
	return &compatibility.CompatibilityConfig{
		OperatingSystems: map[string]compatibility.OSConfig{
			"linux": {
				Supported: true,
				Notes:     "Supported Linux",
				Distributions: map[string]compatibility.DistroConfig{
					"ubuntu": {
						Supported:         true,
						VersionConstraint: ">= 20.04",
						Notes:             "Ubuntu is supported",
					},
					"debian": {
						Supported:         true,
						VersionConstraint: ">= 10",
						Notes:             "Debian is supported",
					},
					"fedora": {
						Supported: false,
						Notes:     "Fedora is not supported",
					},
				},
			},
			"darwin": {
				Supported: true,
				Notes:     "macOS is supported",
			},
			"windows": {
				Supported: false,
				Notes:     "Windows is not supported",
			},
		},
	}
}

//gocognit:ignore
func TestCompatibilityCanBeCheckedWithMockDetectorAndMockConfig(t *testing.T) {
	tests := []struct {
		name        string
		osName      string
		distroName  string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Supported OS Darwin",
			osName:      "darwin",
			distroName:  "",
			expectError: false,
		},
		{
			name:        "Unsupported OS Windows",
			osName:      "windows",
			distroName:  "",
			expectError: true,
			errorMsg:    "unsupported operating system: windows - Windows is not supported",
		},
		{
			name:        "Supported Linux Ubuntu",
			osName:      "linux",
			distroName:  "ubuntu",
			expectError: false,
		},
		{
			name:        "Supported Linux Debian",
			osName:      "linux",
			distroName:  "debian",
			expectError: false,
		},
		{
			name:        "Unsupported Linux Fedora",
			osName:      "linux",
			distroName:  "fedora",
			expectError: true,
			errorMsg:    "unsupported Linux distribution: fedora - Fedora is not supported",
		},
		{
			name:        "Unknown Linux Distro",
			osName:      "linux",
			distroName:  "arch",
			expectError: true,
			errorMsg:    "unsupported Linux distribution: arch",
		},
		{
			name:        "Unknown OS",
			osName:      "solaris",
			distroName:  "",
			expectError: true,
			errorMsg:    "unsupported operating system: solaris",
		},
	}

	config := createMockConfig()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			detector := &MockOSDetector{
				osName:     tc.osName,
				distroName: tc.distroName,
			}

			sysInfo, err := compatibility.CheckCompatibilityWithDetector(config, detector)

			if tc.expectError {
				if err == nil {
					t.Fatalf("Expected error but got nil")
				}
				if err.Error() != tc.errorMsg {
					t.Fatalf("Expected error message '%s', got '%s'", tc.errorMsg, err.Error())
				}

				// Even when there's an error, sysInfo should contain the detected system information
				if tc.osName != sysInfo.OSName {
					t.Fatalf("Expected OS name '%s', got '%s'", tc.osName, sysInfo.OSName)
				}
				if tc.distroName != sysInfo.DistroName && tc.osName == "linux" {
					t.Fatalf("Expected distro name '%s', got '%s'", tc.distroName, sysInfo.DistroName)
				}
			} else {
				if err != nil {
					t.Fatalf("Expected no error, got %v", err)
				}

				// For successful compatibility checks, verify the system info is correct
				if tc.osName != sysInfo.OSName {
					t.Fatalf("Expected OS name '%s', got '%s'", tc.osName, sysInfo.OSName)
				}
				if tc.osName == "linux" && tc.distroName != sysInfo.DistroName {
					t.Fatalf("Expected distro name '%s', got '%s'", tc.distroName, sysInfo.DistroName)
				}
				if sysInfo.Arch == "" {
					t.Fatalf("Expected non-empty architecture")
				}
			}
		})
	}
}

func TestCompatibilityCanBeCheckedWithMockDetectorAndEmbeddedConfig(t *testing.T) {
	detector := &MockOSDetector{osName: "linux", distroName: "ubuntu"}

	compatibilityConfig, err := compatibility.LoadCompatibilityConfig(viper.New(), "")
	if err != nil {
		t.Fatalf("Expected no error when loading embedded compatibility config, got: %v", err)
	}

	sysInfo, err := compatibility.CheckCompatibilityWithDetector(compatibilityConfig, detector)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify the system info is correct
	if sysInfo.OSName != "linux" {
		t.Fatalf("Expected OS name 'linux', got '%s'", sysInfo.OSName)
	}
	if sysInfo.DistroName != "ubuntu" {
		t.Fatalf("Expected distro name 'ubuntu', got '%s'", sysInfo.DistroName)
	}
}

func TestCompatibilityCheckRejectsNilConfig(t *testing.T) {
	detector := &MockOSDetector{osName: "linux", distroName: "ubuntu"}
	sysInfo, err := compatibility.CheckCompatibilityWithDetector(nil, detector)

	if err == nil {
		t.Fatal("Expected error with nil config, got nil")
	}

	expectedMsg := "compatibility configuration is nil"
	if err.Error() != expectedMsg {
		t.Fatalf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}

	// When config is nil, sysInfo should be empty
	if sysInfo.OSName != "" || sysInfo.DistroName != "" || sysInfo.Arch != "" {
		t.Fatalf("Expected empty system info, got %+v", sysInfo)
	}
}
