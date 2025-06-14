package packageresolver

import (
	"fmt"

	"github.com/Masterminds/semver"
	"github.com/MrPointer/dotfiles/installer/lib/pkgmanager"
)

// Resolver translates generic package codes and version constraints into manager-specific package information.
type Resolver struct {
	mappings           *PackageMappingCollection
	packageManagerName string // Normalized name like "apt", "brew"
}

var _ PackageManagerResolver = (*Resolver)(nil)

// NewResolver creates a new package resolver.
// It requires the loaded package mappings and an active PackageManager to determine the current manager.
func NewResolver(
	mappings *PackageMappingCollection,
	pm pkgmanager.PackageManager,
) (*Resolver, error) {
	if mappings == nil {
		return nil, fmt.Errorf("package mappings cannot be nil")
	}
	if pm == nil {
		return nil, fmt.Errorf("package manager cannot be nil")
	}

	pmInfo, err := pm.GetInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get package manager info: %w", err)
	}

	return &Resolver{
		mappings:           mappings,
		packageManagerName: pmInfo.Name,
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
		// If no direct mapping, try to use the generic code as the package name itself.
		// This allows users to specify packages not in the map, assuming the name is consistent.
		// No version constraints can be applied in this case unless the package manager handles it.
		var constraints *semver.Constraints
		var err error
		if versionConstraintString != "" {
			constraints, err = semver.NewConstraint(versionConstraintString)
			if err != nil {
				return pkgmanager.RequestedPackageInfo{}, fmt.Errorf("invalid version constraint string '%s' for package '%s': %w", versionConstraintString, genericPackageCode, err)
			}
		}
		// Consider logging a warning here that a direct mapping was not found.
		return pkgmanager.RequestedPackageInfo{
			Name:               genericPackageCode, // Use the code as the name
			VersionConstraints: constraints,
		}, nil
	}

	var specificPackageName string
	managerSpecificCfg, managerFound := packageMapping[r.packageManagerName]

	if managerFound && managerSpecificCfg.Name != "" {
		specificPackageName = managerSpecificCfg.Name
	} else {
		// No specific name for this manager, fall back to generic code
		specificPackageName = genericPackageCode
		// Consider logging a warning here that a specific mapping was not found,
		// and the generic code is being used as the package name.
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
		VersionConstraints: constraints, // This will be nil if versionConstraintString was empty
	}, nil
}

// PackageManagerResolver defines the interface for resolving package information.
// This is added here to allow for the var _ PackageManagerResolver = (*Resolver)(nil) check.
type PackageManagerResolver interface {
	Resolve(genericPackageCode string, versionConstraintString string) (pkgmanager.RequestedPackageInfo, error)
}
