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
type ManagerSpecificMapping struct {
	// Name is the package name as recognized by the specific package manager.
	Name string `mapstructure:"name"`
}
