package packageresolver

import (
	"fmt"

	"github.com/Masterminds/semver"
	"github.com/MrPointer/dotfiles/installer/lib/compatibility"
	"github.com/MrPointer/dotfiles/installer/lib/pkgmanager"
)

// Resolver translates generic package codes and version constraints into manager-specific package information.
type Resolver struct {
	mappings           *PackageMappingCollection
	packageManagerName string // Normalized name like "apt", "brew"
	systemInfo         *compatibility.SystemInfo
}

var _ PackageManagerResolver = (*Resolver)(nil)

// NewResolver creates a new package resolver.
// It requires the loaded package mappings, an active PackageManager to determine the current manager,
// and system information for distro-specific mappings.
func NewResolver(
	mappings *PackageMappingCollection,
	pm pkgmanager.PackageManager,
	sysInfo *compatibility.SystemInfo,
) (*Resolver, error) {
	if mappings == nil {
		return nil, fmt.Errorf("package mappings cannot be nil")
	}
	if pm == nil {
		return nil, fmt.Errorf("package manager cannot be nil")
	}
	if sysInfo == nil {
		return nil, fmt.Errorf("system info cannot be nil")
	}

	pmInfo, err := pm.GetInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get package manager info: %w", err)
	}

	return &Resolver{
		mappings:           mappings,
		packageManagerName: pmInfo.Name,
		systemInfo:         sysInfo,
	}, nil
}

// Resolve takes a generic package code (e.g., "neovim") and a version constraint string
// (e.g., ">=0.5, <0.7 || >0.8.0"). It returns a RequestedPackageInfo struct containing the
// manager-specific package name and the parsed semver constraints.
// If versionConstraintString is empty, VersionConstraints in the result will be nil.
func (r *Resolver) Resolve(
	genericPackageCode string,
	versionConstraintString string,
) (pkgmanager.RequestedPackageInfo, error) {
	if genericPackageCode == "" {
		return pkgmanager.RequestedPackageInfo{}, fmt.Errorf("generic package code cannot be empty")
	}

	packageMapping, ok := r.mappings.Packages[genericPackageCode]
	if !ok {
		// No mapping found for this package
		return pkgmanager.RequestedPackageInfo{}, fmt.Errorf("no package mapping found for package '%s'", genericPackageCode)
	}

	var specificPackageName string
	var packageType string
	managerSpecificCfg, managerFound := packageMapping[r.packageManagerName]

	if managerFound {
		resolvedName, found := managerSpecificCfg.ResolvePackageName(r.systemInfo.DistroName)
		if found && resolvedName != "" {
			specificPackageName = resolvedName
			packageType = managerSpecificCfg.Type
		} else {
			// Check if this package has distro-specific mappings
			if r.hasDistroSpecificMappings(managerSpecificCfg) {
				// Package requires specific distro handling but current distro is not mapped
				return pkgmanager.RequestedPackageInfo{}, fmt.Errorf("package '%s' requires distro-specific mapping for '%s' distribution, but no mapping is defined", genericPackageCode, r.systemInfo.DistroName)
			}
			// No distro-specific mappings exist, but still no mapping for this package manager
			return pkgmanager.RequestedPackageInfo{}, fmt.Errorf("no package mapping found for package '%s' on package manager '%s'", genericPackageCode, r.packageManagerName)
		}
	} else {
		// No mapping for this package manager
		return pkgmanager.RequestedPackageInfo{}, fmt.Errorf("no package mapping found for package '%s' on package manager '%s'", genericPackageCode, r.packageManagerName)
	}

	var constraints *semver.Constraints
	var err error
	if versionConstraintString != "" {
		constraints, err = semver.NewConstraint(versionConstraintString)
		if err != nil {
			return pkgmanager.RequestedPackageInfo{}, fmt.Errorf("invalid version constraint string '%s' for package '%s': %w", versionConstraintString, genericPackageCode, err)
		}
	}

	return pkgmanager.RequestedPackageInfo{
		Name:               specificPackageName,
		Type:               packageType,
		VersionConstraints: constraints, // This will be nil if versionConstraintString was empty
	}, nil
}

// hasDistroSpecificMappings checks if the ManagerSpecificMapping uses distro-specific name mappings.
func (r *Resolver) hasDistroSpecificMappings(cfg ManagerSpecificMapping) bool {
	switch nameValue := cfg.Name.(type) {
	case string:
		// Simple string case - no distro-specific mappings
		return false
	case map[string]interface{}:
		// Map case - has distro-specific mappings
		return len(nameValue) > 0
	case NameMapping:
		// Direct NameMapping case - has distro-specific mappings
		return len(nameValue) > 0
	default:
		// Unsupported type - assume no distro-specific mappings
		return false
	}
}

// PackageManagerResolver defines the interface for resolving package information.
// This is added here to allow for the var _ PackageManagerResolver = (*Resolver)(nil) check.
type PackageManagerResolver interface {
	Resolve(genericPackageCode string, versionConstraintString string) (pkgmanager.RequestedPackageInfo, error)
}
