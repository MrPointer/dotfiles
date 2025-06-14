package pkgmanager

type PackageManager interface {
	// InstallPackage installs a package by its name.
	InstallPackage(requestedPackageInfo RequestedPackageInfo) error

	// UninstallPackage uninstalls a package by its name.
	UninstallPackage(packageInfo PackageInfo) error

	// ListInstalledPackages returns a list of installed packages.
	ListInstalledPackages() ([]PackageInfo, error)

	// IsPackageInstalled checks if a package is installed by its name.
	IsPackageInstalled(packageInfo PackageInfo) (bool, error)

	// GetPackageVersion retrieves the version of a package by its name.
	GetPackageVersion(packageName string) (string, error)

	// GetInfo retrieves information about the package manager itself.
	GetInfo() (PackageManagerInfo, error)
}
