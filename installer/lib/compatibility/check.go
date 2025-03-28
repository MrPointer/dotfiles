package compatibility

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// OSDetector provides operating system detection capabilities
type OSDetector interface {
	GetOSName() string
	GetDistroName() string
}

// DefaultOSDetector uses runtime and file system to detect OS information
type DefaultOSDetector struct{}

// GetOSName returns the current operating system name
func (d *DefaultOSDetector) GetOSName() string {
	return runtime.GOOS
}

// GetDistroName returns the current Linux distribution name
func (d *DefaultOSDetector) GetDistroName() string {
	return getLinuxDistro()
}

// CompatibilityConfig represents the structure of the compatibility.yaml file
type CompatibilityConfig struct {
	OperatingSystems map[string]OSConfig `yaml:"operatingSystems"`
}

// OSConfig represents configuration for an operating system
type OSConfig struct {
	Supported     bool                    `yaml:"supported"`
	Notes         string                  `yaml:"notes,omitempty"`
	Distributions map[string]DistroConfig `yaml:"distributions,omitempty"`
}

// DistroConfig represents configuration for a Linux distribution
type DistroConfig struct {
	Supported         bool   `yaml:"supported"`
	VersionConstraint string `yaml:"version_constraint,omitempty"`
	Notes             string `yaml:"notes,omitempty"`
}

// CheckCompatibility checks if the current system is compatible based on provided config
func CheckCompatibility(config *CompatibilityConfig) error {
	if config == nil {
		return fmt.Errorf("compatibility configuration is nil")
	}

	detector := &DefaultOSDetector{}
	return CheckCompatibilityWithDetector(config, detector)
}

// CheckCompatibilityWithDetector checks compatibility using the provided detector
func CheckCompatibilityWithDetector(config *CompatibilityConfig, detector OSDetector) error {
	if config == nil {
		return fmt.Errorf("compatibility configuration is nil")
	}

	osName := detector.GetOSName()

	osConfig, exists := config.OperatingSystems[osName]
	if !exists {
		return fmt.Errorf("unsupported operating system: %s", osName)
	}

	if !osConfig.Supported {
		return fmt.Errorf("unsupported operating system: %s - %s", osName, osConfig.Notes)
	}

	// If Linux, check distribution compatibility
	if osName == "linux" {
		distroName := detector.GetDistroName()
		distroConfig, exists := osConfig.Distributions[distroName]

		if !exists {
			return fmt.Errorf("unsupported Linux distribution: %s", distroName)
		}

		if !distroConfig.Supported {
			return fmt.Errorf("unsupported Linux distribution: %s - %s", distroName, distroConfig.Notes)
		}

		// TODO: If needed, add version constraint checking here
	}

	return nil
}

func getLinuxDistro() string {
	// Check for /etc/os-release (freedesktop.org and systemd)
	if _, err := os.Stat("/etc/os-release"); err == nil {
		content, err := os.ReadFile("/etc/os-release")
		if err == nil {
			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "ID=") {
					id := strings.TrimPrefix(line, "ID=")
					// Remove quotes if present
					id = strings.Trim(id, "\"")
					return strings.ToLower(id)
				}
			}
		}
	}

	// Check for /etc/lsb-release
	if _, err := os.Stat("/etc/lsb-release"); err == nil {
		content, err := os.ReadFile("/etc/lsb-release")
		if err == nil {
			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "DISTRIB_ID=") {
					id := strings.TrimPrefix(line, "DISTRIB_ID=")
					// Remove quotes if present
					id = strings.Trim(id, "\"")
					return strings.ToLower(id)
				}
			}
		}
	}

	// Check for /etc/debian_version
	if _, err := os.Stat("/etc/debian_version"); err == nil {
		return "debian"
	}

	// Check for /etc/SuSe-release
	if _, err := os.Stat("/etc/SuSe-release"); err == nil {
		return "suse"
	}

	// Check for /etc/redhat-release
	if _, err := os.Stat("/etc/redhat-release"); err == nil {
		return "redhat"
	}

	// Fallback: use uname
	cmd := exec.Command("uname", "-s")
	output, err := cmd.Output()
	if err == nil {
		return strings.ToLower(strings.TrimSpace(string(output)))
	}

	return "unknown"
}
