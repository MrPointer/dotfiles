package compatibility

import (
	"testing"

	"github.com/spf13/viper"
)

// MockOSDetector implements OSDetector for testing
type MockOSDetector struct {
	osName     string
	distroName string
}

// GetOSName returns the mocked OS name
func (m *MockOSDetector) GetOSName() string {
	return m.osName
}

// GetDistroName returns the mocked distro name
func (m *MockOSDetector) GetDistroName() string {
	return m.distroName
}

// createMockConfig creates a mock compatibility configuration for testing
func createMockConfig() *CompatibilityConfig {
	return &CompatibilityConfig{
		OperatingSystems: map[string]OSConfig{
			"linux": {
				Supported: true,
				Notes:     "Supported Linux",
				Distributions: map[string]DistroConfig{
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

			err := CheckCompatibilityWithDetector(config, detector)

			if tc.expectError {
				if err == nil {
					t.Fatalf("Expected error but got nil")
				}
				if err.Error() != tc.errorMsg {
					t.Fatalf("Expected error message '%s', got '%s'", tc.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("Expected no error, got %v", err)
				}
			}
		})
	}
}

func TestCompatibilityCanBeCheckedWithMockDetectorAndEmbeddedConfig(t *testing.T) {
	detector := &MockOSDetector{osName: "linux", distroName: "ubuntu"}

	compatibilityConfig, err := LoadCompatibilityConfig(viper.New(), "")
	if err != nil {
		t.Fatalf("Expected no error when loading embedded compatibility config, got: %v", err)
	}

	err = CheckCompatibilityWithDetector(compatibilityConfig, detector)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestCompatibilityCheckRejectsNilConfig(t *testing.T) {
	detector := &MockOSDetector{osName: "linux", distroName: "ubuntu"}
	err := CheckCompatibilityWithDetector(nil, detector)

	if err == nil {
		t.Fatal("Expected error with nil config, got nil")
	}

	expectedMsg := "compatibility configuration is nil"
	if err.Error() != expectedMsg {
		t.Fatalf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}
