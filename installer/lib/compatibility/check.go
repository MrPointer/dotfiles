package compatibility

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/MrPointer/dotfiles/installer/utils/osmanager"
)

// SystemInfo contains information about the detected system.
type SystemInfo struct {
	OSName        string             // Operating system name (e.g., "linux", "darwin").
	DistroName    string             // Linux distribution name (e.g., "ubuntu", "debian").
	Arch          string             // Architecture (e.g., "amd64", "arm64").
	Prerequisites PrerequisiteStatus // Status of system prerequisites.
}

// PrerequisiteStatus contains the status of system prerequisites.
type PrerequisiteStatus struct {
	Available []string                      // List of available prerequisites.
	Missing   []string                      // List of missing prerequisites.
	Details   map[string]PrerequisiteDetail // Detailed status for each prerequisite.
}

// PrerequisiteDetail contains detailed information about a prerequisite.
type PrerequisiteDetail struct {
	Name        string // Name of the prerequisite.
	Available   bool   // Whether the prerequisite is available.
	Command     string // Command used to check availability.
	Description string // Human-readable description.
	InstallHint string // Hint for installing the prerequisite.
}

// OSDetector provides operating system detection capabilities.
type OSDetector interface {
	GetOSName() string
	GetDistroName() string
	DetectSystem() (SystemInfo, error)
}

// DefaultOSDetector uses runtime and file system to detect OS information.
type DefaultOSDetector struct{}

// Ensure that DefaultOSDetector implements OSDetector.
var _ OSDetector = (*DefaultOSDetector)(nil)

// NewDefaultOSDetector creates a new DefaultOSDetector.
func NewDefaultOSDetector() *DefaultOSDetector {
	return &DefaultOSDetector{}
}

// GetOSName returns the current operating system name.
func (d *DefaultOSDetector) GetOSName() string {
	return runtime.GOOS
}

// GetDistroName returns the current Linux distribution name.
func (d *DefaultOSDetector) GetDistroName() string {
	return getLinuxDistro()
}

// DetectSystem detects the current system information.
func (d *DefaultOSDetector) DetectSystem() (SystemInfo, error) {
	osName := d.GetOSName()
	var distroName string
	if osName == "linux" {
		distroName = d.GetDistroName()
	} else if osName == "darwin" {
		distroName = "mac"
	}

	return SystemInfo{
		OSName:        osName,
		DistroName:    distroName,
		Arch:          runtime.GOARCH,
		Prerequisites: PrerequisiteStatus{}, // Will be populated by compatibility check
	}, nil
}

// PrerequisiteChecker provides prerequisite checking capabilities.
type PrerequisiteChecker interface {
	CheckPrerequisites(config map[string]PrerequisiteConfig) (PrerequisiteStatus, error)
}

// DefaultPrerequisiteChecker uses ProgramQuery to check for prerequisites.
type DefaultPrerequisiteChecker struct {
	programQuery osmanager.ProgramQuery
}

// Ensure that DefaultPrerequisiteChecker implements PrerequisiteChecker.
var _ PrerequisiteChecker = (*DefaultPrerequisiteChecker)(nil)

// NewDefaultPrerequisiteChecker creates a new DefaultPrerequisiteChecker with the provided ProgramQuery.
func NewDefaultPrerequisiteChecker(programQuery osmanager.ProgramQuery) *DefaultPrerequisiteChecker {
	return &DefaultPrerequisiteChecker{
		programQuery: programQuery,
	}
}

// CheckPrerequisites checks if required prerequisites are available on the system.
func (d *DefaultPrerequisiteChecker) CheckPrerequisites(config map[string]PrerequisiteConfig) (PrerequisiteStatus, error) {
	status := PrerequisiteStatus{
		Available: make([]string, 0),
		Missing:   make([]string, 0),
		Details:   make(map[string]PrerequisiteDetail),
	}

	for name, prereqConfig := range config {
		detail := PrerequisiteDetail{
			Name:        name,
			Command:     prereqConfig.Command,
			Description: prereqConfig.Description,
			InstallHint: prereqConfig.InstallHint,
		}

		// Check if the command is available
		exists, err := d.programQuery.ProgramExists(prereqConfig.Command)
		if err != nil {
			detail.Available = false
			status.Missing = append(status.Missing, name)
		} else if exists {
			detail.Available = true
			status.Available = append(status.Available, name)
		} else {
			detail.Available = false
			status.Missing = append(status.Missing, name)
		}

		status.Details[name] = detail
	}

	// Return error if prerequisites are missing
	if len(status.Missing) > 0 {
		return status, fmt.Errorf("missing prerequisites: %v", status.Missing)
	}

	return status, nil
}

// CompatibilityConfig represents the structure of the compatibility.yaml file.
type CompatibilityConfig struct {
	OperatingSystems map[string]OSConfig `mapstructure:"operatingSystems"`
}

// PrerequisiteConfig represents configuration for a system prerequisite.
type PrerequisiteConfig struct {
	Name        string `mapstructure:"name"`
	Command     string `mapstructure:"command"`
	Description string `mapstructure:"description"`
	InstallHint string `mapstructure:"install_hint"`
}

// OSConfig represents configuration for an operating system.
type OSConfig struct {
	Supported     bool                    `mapstructure:"supported"`
	Notes         string                  `mapstructure:"notes,omitempty"`
	Prerequisites []PrerequisiteConfig    `mapstructure:"prerequisites,omitempty"`
	Distributions map[string]DistroConfig `mapstructure:"distributions,omitempty"`
}

// DistroConfig represents configuration for a Linux distribution.
type DistroConfig struct {
	Supported         bool                 `mapstructure:"supported"`
	VersionConstraint string               `mapstructure:"version_constraint,omitempty"`
	Notes             string               `mapstructure:"notes,omitempty"`
	Prerequisites     []PrerequisiteConfig `mapstructure:"prerequisites,omitempty"`
}

// CheckCompatibility checks if the current system is compatible.
func CheckCompatibility(config *CompatibilityConfig, programQuery osmanager.ProgramQuery) (SystemInfo, error) {
	if config == nil {
		return SystemInfo{}, fmt.Errorf("compatibility configuration is nil")
	}

	detector := NewDefaultOSDetector()
	prereqChecker := NewDefaultPrerequisiteChecker(programQuery)
	return CheckCompatibilityWithDetectors(config, detector, prereqChecker)
}

// CheckCompatibilityWithDetector checks compatibility using the provided detector.
// Deprecated: Use CheckCompatibilityWithDetectors instead.
func CheckCompatibilityWithDetector(config *CompatibilityConfig, detector OSDetector, programQuery osmanager.ProgramQuery) (SystemInfo, error) {
	prereqChecker := NewDefaultPrerequisiteChecker(programQuery)
	return CheckCompatibilityWithDetectors(config, detector, prereqChecker)
}

// CheckCompatibilityWithDetectors checks compatibility using the provided detectors.
func CheckCompatibilityWithDetectors(config *CompatibilityConfig,
	detector OSDetector,
	prereqChecker PrerequisiteChecker) (SystemInfo, error) {
	if config == nil {
		return SystemInfo{}, fmt.Errorf("compatibility configuration is nil")
	}

	// Detect system information
	sysInfo, err := detector.DetectSystem()
	if err != nil {
		return SystemInfo{}, fmt.Errorf("failed to detect system: %w", err)
	}

	// Check prerequisites based on OS and distribution
	var prerequisites []PrerequisiteConfig
	if osConfig, exists := config.OperatingSystems[sysInfo.OSName]; exists {
		// Add OS-level prerequisites
		prerequisites = append(prerequisites, osConfig.Prerequisites...)

		// Add distribution-specific prerequisites for Linux
		if sysInfo.OSName == "linux" {
			if distroConfig, exists := osConfig.Distributions[sysInfo.DistroName]; exists {
				prerequisites = append(prerequisites, distroConfig.Prerequisites...)
			}
		}
	}

	// Convert to map format for checker
	prereqMap := make(map[string]PrerequisiteConfig)
	for _, prereq := range prerequisites {
		prereqMap[prereq.Name] = prereq
	}

	prereqStatus, err := prereqChecker.CheckPrerequisites(prereqMap)
	if err != nil {
		sysInfo.Prerequisites = prereqStatus
		return sysInfo, fmt.Errorf("prerequisite check failed: %w", err)
	}
	sysInfo.Prerequisites = prereqStatus

	// Check if the operating system is supported
	osConfig, exists := config.OperatingSystems[sysInfo.OSName]
	if !exists {
		return sysInfo, fmt.Errorf("unsupported operating system: %s", sysInfo.OSName)
	}

	if !osConfig.Supported {
		return sysInfo, fmt.Errorf("unsupported operating system: %s - %s", sysInfo.OSName, osConfig.Notes)
	}

	// If Linux, check distribution compatibility
	if sysInfo.OSName == "linux" {
		distroConfig, exists := osConfig.Distributions[sysInfo.DistroName]

		if !exists {
			return sysInfo, fmt.Errorf("unsupported Linux distribution: %s", sysInfo.DistroName)
		}

		if !distroConfig.Supported {
			return sysInfo, fmt.Errorf(
				"unsupported Linux distribution: %s - %s",
				sysInfo.DistroName,
				distroConfig.Notes,
			)
		}

		// TODO: If needed, add version constraint checking here
	}

	// System is compatible
	return sysInfo, nil
}

func getLinuxDistro() string {
	// Check for /etc/os-release (freedesktop.org and systemd).
	if _, err := os.Stat("/etc/os-release"); err == nil {
		content, err := os.ReadFile("/etc/os-release")
		if err == nil {
			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "ID=") {
					id := strings.TrimPrefix(line, "ID=")
					// Remove quotes if present.
					id = strings.Trim(id, "\"")
					return strings.ToLower(id)
				}
			}
		}
	}

	// Check for /etc/lsb-release.
	if _, err := os.Stat("/etc/lsb-release"); err == nil {
		content, err := os.ReadFile("/etc/lsb-release")
		if err == nil {
			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "DISTRIB_ID=") {
					id := strings.TrimPrefix(line, "DISTRIB_ID=")
					// Remove quotes if present.
					id = strings.Trim(id, "\"")
					return strings.ToLower(id)
				}
			}
		}
	}

	// Check for /etc/debian_version.
	if _, err := os.Stat("/etc/debian_version"); err == nil {
		return "debian"
	}

	// Check for /etc/SuSe-release.
	if _, err := os.Stat("/etc/SuSe-release"); err == nil {
		return "suse"
	}

	// Check for /etc/redhat-release.
	if _, err := os.Stat("/etc/redhat-release"); err == nil {
		return "redhat"
	}

	// Fallback: use uname.
	cmd := exec.Command("uname", "-s")
	output, err := cmd.Output()
	if err == nil {
		return strings.ToLower(strings.TrimSpace(string(output)))
	}

	return "unknown"
}
