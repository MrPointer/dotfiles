package compatibility_test

import (
	"errors"
	"testing"

	"github.com/MrPointer/dotfiles/installer/lib/compatibility"
	"github.com/MrPointer/dotfiles/installer/utils/osmanager"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
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
		OSName:        m.osName,
		DistroName:    m.distroName,
		Arch:          "amd64",                            // Mock architecture
		Prerequisites: compatibility.PrerequisiteStatus{}, // Will be populated by compatibility check
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
						Prerequisites: []compatibility.PrerequisiteConfig{
							{
								Name:        "git",
								Command:     "git",
								Description: "Git version control system",
								InstallHint: "sudo apt-get install git",
							},
							{
								Name:        "curl",
								Command:     "curl",
								Description: "Command line tool for transferring data",
								InstallHint: "sudo apt-get install curl",
							},
						},
					},
					"debian": {
						Supported:         true,
						VersionConstraint: ">= 10",
						Notes:             "Debian is supported",
						Prerequisites: []compatibility.PrerequisiteConfig{
							{
								Name:        "git",
								Command:     "git",
								Description: "Git version control system",
								InstallHint: "sudo apt-get install git",
							},
						},
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
				Prerequisites: []compatibility.PrerequisiteConfig{
					{
						Name:        "git",
						Command:     "git",
						Description: "Git version control system",
						InstallHint: "Install Command Line Tools: xcode-select --install",
					},
				},
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

			// Create a mock program query that returns true for all programs
			mockProgramQuery := &osmanager.MoqProgramQuery{
				ProgramExistsFunc: func(program string) (bool, error) {
					return true, nil
				},
			}
			sysInfo, err := compatibility.CheckCompatibilityWithDetector(config, detector, mockProgramQuery)

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

	// Create a mock program query that returns true for all programs
	mockProgramQuery := &osmanager.MoqProgramQuery{
		ProgramExistsFunc: func(program string) (bool, error) {
			return true, nil
		},
	}
	sysInfo, err := compatibility.CheckCompatibilityWithDetector(compatibilityConfig, detector, mockProgramQuery)
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
	// Create a mock program query
	mockProgramQuery := &osmanager.MoqProgramQuery{
		ProgramExistsFunc: func(program string) (bool, error) {
			return true, nil
		},
	}
	sysInfo, err := compatibility.CheckCompatibilityWithDetector(nil, detector, mockProgramQuery)

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

func Test_DefaultPrerequisiteChecker_CheckPrerequisites_WhenAllPrerequisitesAvailable(t *testing.T) {
	// Arrange
	mockProgramQuery := &osmanager.MoqProgramQuery{
		ProgramExistsFunc: func(program string) (bool, error) {
			return true, nil
		},
	}
	checker := compatibility.NewDefaultPrerequisiteChecker(mockProgramQuery)

	prereqConfig := map[string]compatibility.PrerequisiteConfig{
		"git": {
			Name:        "git",
			Command:     "git",
			Description: "Git version control system",
			InstallHint: "Install git",
		},
		"curl": {
			Name:        "curl",
			Command:     "curl",
			Description: "HTTP client",
			InstallHint: "Install curl",
		},
	}

	// Act
	status, err := checker.CheckPrerequisites(prereqConfig)

	// Assert
	require.NoError(t, err)
	require.Len(t, status.Available, 2)
	require.Len(t, status.Missing, 0)
	require.Contains(t, status.Available, "git")
	require.Contains(t, status.Available, "curl")
	require.Len(t, status.Details, 2)

	gitDetail := status.Details["git"]
	require.True(t, gitDetail.Available)
	require.Equal(t, "git", gitDetail.Name)
	require.Equal(t, "Git version control system", gitDetail.Description)
}

func Test_DefaultPrerequisiteChecker_CheckPrerequisites_WhenSomePrerequisitesMissing(t *testing.T) {
	// Arrange
	mockProgramQuery := &osmanager.MoqProgramQuery{
		ProgramExistsFunc: func(program string) (bool, error) {
			// Only git is available, curl is missing
			return program == "git", nil
		},
	}
	checker := compatibility.NewDefaultPrerequisiteChecker(mockProgramQuery)

	prereqConfig := map[string]compatibility.PrerequisiteConfig{
		"git": {
			Name:        "git",
			Command:     "git",
			Description: "Git version control system",
			InstallHint: "Install git",
		},
		"curl": {
			Name:        "curl",
			Command:     "curl",
			Description: "HTTP client",
			InstallHint: "Install curl",
		},
	}

	// Act
	status, err := checker.CheckPrerequisites(prereqConfig)

	// Assert
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing prerequisites")
	require.Len(t, status.Available, 1)
	require.Len(t, status.Missing, 1)
	require.Contains(t, status.Available, "git")
	require.Contains(t, status.Missing, "curl")

	gitDetail := status.Details["git"]
	require.True(t, gitDetail.Available)

	curlDetail := status.Details["curl"]
	require.False(t, curlDetail.Available)
}

func Test_DefaultPrerequisiteChecker_CheckPrerequisites_WhenAllPrerequisitesMissing(t *testing.T) {
	// Arrange
	mockProgramQuery := &osmanager.MoqProgramQuery{
		ProgramExistsFunc: func(program string) (bool, error) {
			return false, nil
		},
	}
	checker := compatibility.NewDefaultPrerequisiteChecker(mockProgramQuery)

	prereqConfig := map[string]compatibility.PrerequisiteConfig{
		"git": {
			Name:        "git",
			Command:     "git",
			Description: "Git version control system",
			InstallHint: "Install git",
		},
		"make": {
			Name:        "make",
			Command:     "make",
			Description: "Build tool",
			InstallHint: "Install make",
		},
	}

	// Act
	status, err := checker.CheckPrerequisites(prereqConfig)

	// Assert
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing prerequisites")
	require.Len(t, status.Available, 0)
	require.Len(t, status.Missing, 2)
	require.Contains(t, status.Missing, "git")
	require.Contains(t, status.Missing, "make")
}

func Test_DefaultPrerequisiteChecker_CheckPrerequisites_WhenProgramQueryReturnsError(t *testing.T) {
	// Arrange
	mockProgramQuery := &osmanager.MoqProgramQuery{
		ProgramExistsFunc: func(program string) (bool, error) {
			return false, errors.New("program query failed")
		},
	}
	checker := compatibility.NewDefaultPrerequisiteChecker(mockProgramQuery)

	prereqConfig := map[string]compatibility.PrerequisiteConfig{
		"git": {
			Name:        "git",
			Command:     "git",
			Description: "Git version control system",
			InstallHint: "Install git",
		},
	}

	// Act
	status, err := checker.CheckPrerequisites(prereqConfig)

	// Assert
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing prerequisites")
	require.Len(t, status.Available, 0)
	require.Len(t, status.Missing, 1)
	require.Contains(t, status.Missing, "git")

	gitDetail := status.Details["git"]
	require.False(t, gitDetail.Available)
}

func Test_DefaultPrerequisiteChecker_CheckPrerequisites_WhenNoPrerequisitesProvided(t *testing.T) {
	// Arrange
	mockProgramQuery := &osmanager.MoqProgramQuery{
		ProgramExistsFunc: func(program string) (bool, error) {
			return true, nil
		},
	}
	checker := compatibility.NewDefaultPrerequisiteChecker(mockProgramQuery)

	prereqConfig := map[string]compatibility.PrerequisiteConfig{}

	// Act
	status, err := checker.CheckPrerequisites(prereqConfig)

	// Assert
	require.NoError(t, err)
	require.Len(t, status.Available, 0)
	require.Len(t, status.Missing, 0)
	require.Len(t, status.Details, 0)
}
