package pkgmanager

import "github.com/Masterminds/semver"

type PackageManagerInfo struct {
	// Name of the package manager.
	Name string `json:"name"`

	// Version of the package manager.
	Version string `json:"version"`
}

func NewPackageManagerInfo(name, version string) PackageManagerInfo {
	return PackageManagerInfo{
		Name:    name,
		Version: version,
	}
}

func DefaultPackageManagerInfo() PackageManagerInfo {
	return PackageManagerInfo{
		Name:    "Unknown",
		Version: "0.0.0",
	}
}

type PackageInfo struct {
	// Name of the package.
	Name string `json:"name"`

	// Version of the package.
	Version string `json:"version"`

	// Type of the package (e.g., "group", "pattern"). Empty means regular package.
	Type string `json:"type,omitempty"`
}

func NewPackageInfo(name, version string) PackageInfo {
	return PackageInfo{
		Name:    name,
		Version: version,
	}
}

func NewPackageInfoWithType(name, version, packageType string) PackageInfo {
	return PackageInfo{
		Name:    name,
		Version: version,
		Type:    packageType,
	}
}

type RequestedPackageInfo struct {
	// Name of the package.
	Name string `json:"name"`

	// Type of the package (e.g., "group", "pattern"). Empty means regular package.
	Type string `json:"type,omitempty"`

	// VersionConstraints defines the semver constraints for the requested package.
	// It's a pointer to allow for nil (no constraints).
	VersionConstraints *semver.Constraints `json:"version_constraint,omitempty"`
}

func NewRequestedPackageInfo(name string, versionConstraints *semver.Constraints) RequestedPackageInfo {
	return RequestedPackageInfo{
		Name:               name,
		VersionConstraints: versionConstraints,
	}
}

func NewRequestedPackageInfoWithType(name, packageType string, versionConstraints *semver.Constraints) RequestedPackageInfo {
	return RequestedPackageInfo{
		Name:               name,
		Type:               packageType,
		VersionConstraints: versionConstraints,
	}
}
