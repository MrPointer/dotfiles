package packageresolver

// PackageMappingCollection holds all package mappings defined in the configuration.
// The top-level key in the configuration YAML (e.g., "packages") maps to this structure.
type PackageMappingCollection struct {
	Packages map[string]PackageMapping `mapstructure:"packages"`
}

// PackageMapping maps package manager names directly to their specific configurations.
// For example: "apt" -> ManagerSpecificMapping{Name: "git"}, "brew" -> ManagerSpecificMapping{Name: "git"}
type PackageMapping map[string]ManagerSpecificMapping

// ManagerSpecificMapping holds the actual package name for a specific package manager.
// The Name field can be either a string (single name for all distros) or a NameMapping (distro-specific names).
type ManagerSpecificMapping struct {
	// Name can be either:
	// - A string: single package name used for all distributions
	// - A NameMapping: map of distribution-specific names
	Name any `mapstructure:"name"`

	// Type is the package type (e.g., "group", "pattern"). Empty means regular package.
	Type string `mapstructure:"type,omitempty"`
}

// NameMapping holds distribution-specific package names.
// It supports only explicit distro mappings with no fallback behavior.
type NameMapping map[string]string

// GetNameForDistro resolves the package name for a specific distribution.
// It only returns exact matches - no fallback behavior.
func (nm NameMapping) GetNameForDistro(distroName string) (string, bool) {
	// Try exact distro match only
	if name, exists := nm[distroName]; exists {
		return name, true
	}

	// No match found
	return "", false
}

// ResolvePackageName resolves the package name from ManagerSpecificMapping for a given distribution.
// It handles both string and NameMapping types for the Name field.
func (msm *ManagerSpecificMapping) ResolvePackageName(distroName string) (string, bool) {
	switch nameValue := msm.Name.(type) {
	case string:
		// Simple string case - same name for all distros
		return nameValue, nameValue != ""
	case map[string]any:
		// Convert to NameMapping for processing
		nameMapping := make(NameMapping)
		for k, v := range nameValue {
			if str, ok := v.(string); ok {
				nameMapping[k] = str
			}
		}
		return nameMapping.GetNameForDistro(distroName)
	case NameMapping:
		// Direct NameMapping case
		return nameValue.GetNameForDistro(distroName)
	default:
		// Unsupported type
		return "", false
	}
}
