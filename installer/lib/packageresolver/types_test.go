package packageresolver_test

import (
	"testing"

	"github.com/MrPointer/dotfiles/installer/lib/packageresolver"
	"github.com/stretchr/testify/require"
)

func Test_PackageMappingCollection_CanBeInstantiated(t *testing.T) {
	collection := packageresolver.PackageMappingCollection{
		Packages: make(map[string]packageresolver.PackageMapping),
	}

	require.NotNil(t, collection.Packages)
	require.Empty(t, collection.Packages)
}

func Test_PackageMappingCollection_CanStorePackageMappings(t *testing.T) {
	collection := packageresolver.PackageMappingCollection{
		Packages: map[string]packageresolver.PackageMapping{
			"test-package": {
				"apt":  {Name: "apt-test-package"},
				"brew": {Name: "brew-test-package"},
			},
		},
	}

	require.Len(t, collection.Packages, 1)
	require.Contains(t, collection.Packages, "test-package")

	packageMapping := collection.Packages["test-package"]
	require.Len(t, packageMapping, 2)
	require.Contains(t, packageMapping, "apt")
	require.Contains(t, packageMapping, "brew")
	require.Equal(t, "apt-test-package", packageMapping["apt"].Name)
	require.Equal(t, "brew-test-package", packageMapping["brew"].Name)
}

func Test_PackageMapping_CanBeInstantiated(t *testing.T) {
	mapping := packageresolver.PackageMapping{
		"apt": {Name: "test-package"},
	}

	require.Len(t, mapping, 1)
	require.Contains(t, mapping, "apt")
	require.Equal(t, "test-package", mapping["apt"].Name)
}

func Test_PackageMapping_CanBeInstantiated_AsEmptyMap(t *testing.T) {
	mapping := packageresolver.PackageMapping{}

	require.Empty(t, mapping)
	require.NotNil(t, mapping)
}

func Test_PackageMapping_CanStoreMultipleManagers(t *testing.T) {
	mapping := packageresolver.PackageMapping{
		"apt":    {Name: "apt-package"},
		"brew":   {Name: "homebrew-package"},
		"dnf":    {Name: "dnf-package"},
		"pacman": {Name: "arch-package"},
	}

	require.Len(t, mapping, 4)

	require.Equal(t, "apt-package", mapping["apt"].Name)
	require.Equal(t, "homebrew-package", mapping["brew"].Name)
	require.Equal(t, "dnf-package", mapping["dnf"].Name)
	require.Equal(t, "arch-package", mapping["pacman"].Name)
}

func Test_PackageMapping_CanAccessNonExistentManager(t *testing.T) {
	mapping := packageresolver.PackageMapping{
		"apt": {Name: "test-package"},
	}

	// Accessing non-existent key should return zero value
	nonExistent := mapping["non-existent"]
	require.Empty(t, nonExistent.Name)

	// Check with ok pattern
	value, exists := mapping["non-existent"]
	require.False(t, exists)
	require.Empty(t, value.Name)
}

func Test_ManagerSpecificMapping_CanBeInstantiated(t *testing.T) {
	mapping := packageresolver.ManagerSpecificMapping{
		Name: "specific-package-name",
	}

	require.Equal(t, "specific-package-name", mapping.Name)
}

func Test_ManagerSpecificMapping_CanBeInstantiated_WithEmptyName(t *testing.T) {
	mapping := packageresolver.ManagerSpecificMapping{}

	require.Empty(t, mapping.Name)
}

func Test_PackageMapping_SupportsRealWorldStructure(t *testing.T) {
	// This test mimics the actual structure from packagemap.yaml
	mapping := packageresolver.PackageMapping{
		"apt":  {Name: "git"},
		"brew": {Name: "git"},
	}

	require.Len(t, mapping, 2)
	require.Equal(t, "git", mapping["apt"].Name)
	require.Equal(t, "git", mapping["brew"].Name)
}

func Test_PackageMappingCollection_SupportsRealWorldStructure(t *testing.T) {
	// This test mimics the actual structure from packagemap.yaml
	collection := packageresolver.PackageMappingCollection{
		Packages: map[string]packageresolver.PackageMapping{
			"git": {
				"apt":  {Name: "git"},
				"brew": {Name: "git"},
			},
			"gpg": {
				"apt":  {Name: "gnupg2"},
				"brew": {Name: "gnupg"},
				"dnf":  {Name: "gnupg2"},
			},
			"neovim": {
				"apt":  {Name: "neovim"},
				"brew": {Name: "neovim"},
			},
		},
	}

	require.Len(t, collection.Packages, 3)

	// Test git package
	gitMapping := collection.Packages["git"]
	require.Len(t, gitMapping, 2)
	require.Equal(t, "git", gitMapping["apt"].Name)
	require.Equal(t, "git", gitMapping["brew"].Name)

	// Test gpg package with different manager names
	gpgMapping := collection.Packages["gpg"]
	require.Len(t, gpgMapping, 3)
	require.Equal(t, "gnupg2", gpgMapping["apt"].Name)
	require.Equal(t, "gnupg", gpgMapping["brew"].Name)
	require.Equal(t, "gnupg2", gpgMapping["dnf"].Name)

	// Test neovim package
	neovimMapping := collection.Packages["neovim"]
	require.Len(t, neovimMapping, 2)
	require.Equal(t, "neovim", neovimMapping["apt"].Name)
	require.Equal(t, "neovim", neovimMapping["brew"].Name)
}

func Test_NameMapping_GetNameForDistro_ReturnsExactMatch_WhenDistroExists(t *testing.T) {
	nameMapping := packageresolver.NameMapping{
		"fedora": "development-tools",
		"centos": "Development Tools",
		"ubuntu": "build-essential",
	}

	name, found := nameMapping.GetNameForDistro("fedora")
	require.True(t, found)
	require.Equal(t, "development-tools", name)

	name, found = nameMapping.GetNameForDistro("centos")
	require.True(t, found)
	require.Equal(t, "Development Tools", name)

	name, found = nameMapping.GetNameForDistro("ubuntu")
	require.True(t, found)
	require.Equal(t, "build-essential", name)
}

func Test_NameMapping_GetNameForDistro_ReturnsNotFound_WhenDistroNotMapped(t *testing.T) {
	nameMapping := packageresolver.NameMapping{
		"fedora": "development-tools",
		"centos": "Development Tools",
	}

	// Should return not found for unmapped distro
	name, found := nameMapping.GetNameForDistro("rhel")
	require.False(t, found)
	require.Equal(t, "", name)

	name, found = nameMapping.GetNameForDistro("unknown")
	require.False(t, found)
	require.Equal(t, "", name)
}

func Test_NameMapping_GetNameForDistro_ReturnsNotFound_WhenNoMatchOrDefault(t *testing.T) {
	nameMapping := packageresolver.NameMapping{
		"fedora": "development-tools",
		"centos": "Development Tools",
	}

	name, found := nameMapping.GetNameForDistro("unknown-distro")
	require.False(t, found)
	require.Equal(t, "", name)
}

func Test_NameMapping_GetNameForDistro_ReturnsExactMatchOnly(t *testing.T) {
	nameMapping := packageresolver.NameMapping{
		"fedora": "development-tools",
		"centos": "Development Tools",
	}

	// Should return exact matches only
	name, found := nameMapping.GetNameForDistro("fedora")
	require.True(t, found)
	require.Equal(t, "development-tools", name)

	name, found = nameMapping.GetNameForDistro("centos")
	require.True(t, found)
	require.Equal(t, "Development Tools", name)
}

func Test_ManagerSpecificMapping_ResolvePackageName_HandlesStringName(t *testing.T) {
	mapping := packageresolver.ManagerSpecificMapping{
		Name: "git",
		Type: "",
	}

	name, found := mapping.ResolvePackageName("fedora")
	require.True(t, found)
	require.Equal(t, "git", name)

	name, found = mapping.ResolvePackageName("centos")
	require.True(t, found)
	require.Equal(t, "git", name)
}

func Test_ManagerSpecificMapping_ResolvePackageName_HandlesMapName(t *testing.T) {
	mapping := packageresolver.ManagerSpecificMapping{
		Name: map[string]interface{}{
			"fedora": "development-tools",
			"centos": "Development Tools",
		},
		Type: "group",
	}

	name, found := mapping.ResolvePackageName("fedora")
	require.True(t, found)
	require.Equal(t, "development-tools", name)

	name, found = mapping.ResolvePackageName("centos")
	require.True(t, found)
	require.Equal(t, "Development Tools", name)
}

func Test_ManagerSpecificMapping_ResolvePackageName_HandlesMapNameWithoutFallback(t *testing.T) {
	mapping := packageresolver.ManagerSpecificMapping{
		Name: map[string]interface{}{
			"fedora": "development-tools",
			"centos": "Development Tools",
		},
		Type: "group",
	}

	name, found := mapping.ResolvePackageName("fedora")
	require.True(t, found)
	require.Equal(t, "development-tools", name)

	// Should not find unmapped distro
	name, found = mapping.ResolvePackageName("rhel")
	require.False(t, found)
	require.Equal(t, "", name)
}

func Test_ManagerSpecificMapping_ResolvePackageName_ReturnsNotFound_WhenNoMatch(t *testing.T) {
	mapping := packageresolver.ManagerSpecificMapping{
		Name: map[string]interface{}{
			"fedora": "development-tools",
			"centos": "Development Tools",
		},
		Type: "group",
	}

	name, found := mapping.ResolvePackageName("unknown-distro")
	require.False(t, found)
	require.Equal(t, "", name)
}

func Test_ManagerSpecificMapping_ResolvePackageName_HandlesEmptyString(t *testing.T) {
	mapping := packageresolver.ManagerSpecificMapping{
		Name: "",
		Type: "",
	}

	name, found := mapping.ResolvePackageName("fedora")
	require.False(t, found)
	require.Equal(t, "", name)
}

func Test_ManagerSpecificMapping_ResolvePackageName_HandlesNilName(t *testing.T) {
	mapping := packageresolver.ManagerSpecificMapping{
		Name: nil,
		Type: "",
	}

	name, found := mapping.ResolvePackageName("fedora")
	require.False(t, found)
	require.Equal(t, "", name)
}

func Test_ManagerSpecificMapping_ResolvePackageName_HandlesDirectNameMapping(t *testing.T) {
	nameMapping := packageresolver.NameMapping{
		"fedora": "development-tools",
		"centos": "Development Tools",
	}

	mapping := packageresolver.ManagerSpecificMapping{
		Name: nameMapping,
		Type: "group",
	}

	name, found := mapping.ResolvePackageName("fedora")
	require.True(t, found)
	require.Equal(t, "development-tools", name)

	name, found = mapping.ResolvePackageName("centos")
	require.True(t, found)
	require.Equal(t, "Development Tools", name)

	// Should not find unmapped distro
	name, found = mapping.ResolvePackageName("rhel")
	require.False(t, found)
	require.Equal(t, "", name)
}
