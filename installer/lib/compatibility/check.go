package compatibility

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// SystemInfo contains information about the detected system
type SystemInfo struct {
	OSName     string // Operating system name (e.g., "linux", "darwin")
	DistroName string // Linux distribution name (e.g., "ubuntu", "debian")
	Arch       string // Architecture (e.g., "amd64", "arm64")
}

// OSDetector provides operating system detection capabilities
type OSDetector interface {
	GetOSName() string
	GetDistroName() string
	DetectSystem() (SystemInfo, error)
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

// DetectSystem detects the current system information
func (d *DefaultOSDetector) DetectSystem() (SystemInfo, error) {
	osName := d.GetOSName()
	var distroName string
	if osName == "linux" {
		distroName = d.GetDistroName()
	} else if osName == "darwin" {
		distroName = "mac"
	}

	return SystemInfo{
		OSName:     osName,
		DistroName: distroName,
		Arch:       runtime.GOARCH,
	}, nil
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

// CheckCompatibility checks if the current system is compatible
func CheckCompatibility(config *CompatibilityConfig) (SystemInfo, error) {
	if config == nil {
		return SystemInfo{}, fmt.Errorf("compatibility configuration is nil")
	}

	detector := &DefaultOSDetector{}
	return CheckCompatibilityWithDetector(config, detector)
}

// CheckCompatibilityWithDetector checks compatibility using the provided detector
func CheckCompatibilityWithDetector(config *CompatibilityConfig, detector OSDetector) (SystemInfo, error) {
	if config == nil {
		return SystemInfo{}, fmt.Errorf("compatibility configuration is nil")
	}

	// Detect system information
	sysInfo, err := detector.DetectSystem()
	if err != nil {
		return SystemInfo{}, fmt.Errorf("failed to detect system: %w", err)
	}

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
			return sysInfo, fmt.Errorf("unsupported Linux distribution: %s - %s", sysInfo.DistroName, distroConfig.Notes)
		}

		// TODO: If needed, add version constraint checking here
	}

	// System is compatible
	return sysInfo, nil
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
