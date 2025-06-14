package packageresolver_test

import (
	"errors"
	"testing"

	"github.com/Masterminds/semver"
	"github.com/MrPointer/dotfiles/installer/lib/packageresolver"
	"github.com/MrPointer/dotfiles/installer/lib/pkgmanager"
	"github.com/stretchr/testify/require"
)

func Test_NewResolver_CanCreateResolver_WithValidInputs(t *testing.T) {
	mappings := &packageresolver.PackageMappingCollection{
		Packages: make(map[string]packageresolver.PackageMapping),
	}
	mockPM := &pkgmanager.MoqPackageManager{
		GetInfoFunc: func() (pkgmanager.PackageManagerInfo, error) {
			return pkgmanager.PackageManagerInfo{Name: "apt"}, nil
		},
	}

	resolver, err := packageresolver.NewResolver(mappings, mockPM)

	require.NoError(t, err)
	require.NotNil(t, resolver)
}

func Test_NewResolver_ReturnsError_WhenMappingsIsNil(t *testing.T) {
	mockPM := &pkgmanager.MoqPackageManager{
		GetInfoFunc: func() (pkgmanager.PackageManagerInfo, error) {
			return pkgmanager.PackageManagerInfo{Name: "apt"}, nil
		},
	}

	resolver, err := packageresolver.NewResolver(nil, mockPM)

	require.Error(t, err)
	require.Nil(t, resolver)
	require.Contains(t, err.Error(), "package mappings cannot be nil")
}

func Test_NewResolver_ReturnsError_WhenPackageManagerIsNil(t *testing.T) {
	mappings := &packageresolver.PackageMappingCollection{
		Packages: make(map[string]packageresolver.PackageMapping),
	}

	resolver, err := packageresolver.NewResolver(mappings, nil)

	require.Error(t, err)
	require.Nil(t, resolver)
	require.Contains(t, err.Error(), "package manager cannot be nil")
}

func Test_NewResolver_ReturnsError_WhenPackageManagerGetInfoFails(t *testing.T) {
	mappings := &packageresolver.PackageMappingCollection{
		Packages: make(map[string]packageresolver.PackageMapping),
	}
	mockPM := &pkgmanager.MoqPackageManager{
		GetInfoFunc: func() (pkgmanager.PackageManagerInfo, error) {
			return pkgmanager.PackageManagerInfo{}, errors.New("failed to get info")
		},
	}

	resolver, err := packageresolver.NewResolver(mappings, mockPM)

	require.Error(t, err)
	require.Nil(t, resolver)
	require.Contains(t, err.Error(), "failed to get package manager info")
}

func Test_NewResolver_AcceptsAnyPackageManagerName(t *testing.T) {
	mappings := &packageresolver.PackageMappingCollection{
		Packages: make(map[string]packageresolver.PackageMapping),
	}
	mockPM := &pkgmanager.MoqPackageManager{
		GetInfoFunc: func() (pkgmanager.PackageManagerInfo, error) {
			return pkgmanager.PackageManagerInfo{Name: "custom-manager"}, nil
		},
	}

	resolver, err := packageresolver.NewResolver(mappings, mockPM)

	require.NoError(t, err)
	require.NotNil(t, resolver)
}

func Test_NewResolver_UsesPackageManagerNameDirectly(t *testing.T) {
	testCases := []struct {
		name   string
		pmName string
	}{
		{
			name:   "uses apt name directly",
			pmName: "apt",
		},
		{
			name:   "uses brew name directly",
			pmName: "brew",
		},
		{
			name:   "uses dnf name directly",
			pmName: "dnf",
		},
		{
			name:   "uses pacman name directly",
			pmName: "pacman",
		},
		{
			name:   "uses custom name directly",
			pmName: "custom-pm",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mappings := &packageresolver.PackageMappingCollection{
				Packages: make(map[string]packageresolver.PackageMapping),
			}
			mockPM := &pkgmanager.MoqPackageManager{
				GetInfoFunc: func() (pkgmanager.PackageManagerInfo, error) {
					return pkgmanager.PackageManagerInfo{Name: tc.pmName}, nil
				},
			}

			resolver, err := packageresolver.NewResolver(mappings, mockPM)

			require.NoError(t, err)
			require.NotNil(t, resolver)
		})
	}
}

func Test_Resolve_ReturnsError_WhenGenericPackageCodeIsEmpty(t *testing.T) {
	mappings := &packageresolver.PackageMappingCollection{
		Packages: make(map[string]packageresolver.PackageMapping),
	}
	mockPM := &pkgmanager.MoqPackageManager{
		GetInfoFunc: func() (pkgmanager.PackageManagerInfo, error) {
			return pkgmanager.PackageManagerInfo{Name: "apt"}, nil
		},
	}

	resolver, err := packageresolver.NewResolver(mappings, mockPM)
	require.NoError(t, err)

	result, err := resolver.Resolve("", "")

	require.Error(t, err)
	require.Equal(t, pkgmanager.RequestedPackageInfo{}, result)
	require.Contains(t, err.Error(), "generic package code cannot be empty")
}

func Test_Resolve_UsesManagerSpecificName_WhenAvailable(t *testing.T) {
	mappings := &packageresolver.PackageMappingCollection{
		Packages: map[string]packageresolver.PackageMapping{
			"neovim": {
				"apt":  {Name: "neovim"},
				"brew": {Name: "neovim"},
			},
		},
	}
	mockPM := &pkgmanager.MoqPackageManager{
		GetInfoFunc: func() (pkgmanager.PackageManagerInfo, error) {
			return pkgmanager.PackageManagerInfo{Name: "apt"}, nil
		},
	}

	resolver, err := packageresolver.NewResolver(mappings, mockPM)
	require.NoError(t, err)

	result, err := resolver.Resolve("neovim", "")

	require.NoError(t, err)
	require.Equal(t, "neovim", result.Name)
	require.Nil(t, result.VersionConstraints)
}

func Test_Resolve_FallsBackToGenericCode_WhenManagerSpecificNameNotFound(t *testing.T) {
	mappings := &packageresolver.PackageMappingCollection{
		Packages: map[string]packageresolver.PackageMapping{
			"nodejs": {
				"brew": {Name: "node"},
			},
		},
	}
	mockPM := &pkgmanager.MoqPackageManager{
		GetInfoFunc: func() (pkgmanager.PackageManagerInfo, error) {
			return pkgmanager.PackageManagerInfo{Name: "apt"}, nil
		},
	}

	resolver, err := packageresolver.NewResolver(mappings, mockPM)
	require.NoError(t, err)

	result, err := resolver.Resolve("nodejs", "")

	require.NoError(t, err)
	require.Equal(t, "nodejs", result.Name) // Falls back to generic code
	require.Nil(t, result.VersionConstraints)
}

func Test_Resolve_FallsBackToGenericCode_WhenNoMappingFound(t *testing.T) {
	mappings := &packageresolver.PackageMappingCollection{
		Packages: make(map[string]packageresolver.PackageMapping),
	}
	mockPM := &pkgmanager.MoqPackageManager{
		GetInfoFunc: func() (pkgmanager.PackageManagerInfo, error) {
			return pkgmanager.PackageManagerInfo{Name: "apt"}, nil
		},
	}

	resolver, err := packageresolver.NewResolver(mappings, mockPM)
	require.NoError(t, err)

	result, err := resolver.Resolve("unknown-package", "")

	require.NoError(t, err)
	require.Equal(t, "unknown-package", result.Name)
	require.Nil(t, result.VersionConstraints)
}

func Test_Resolve_FallsBackToGenericCode_WhenManagerSpecificNameIsEmpty(t *testing.T) {
	mappings := &packageresolver.PackageMappingCollection{
		Packages: map[string]packageresolver.PackageMapping{
			"test-package": {
				"apt": {Name: ""}, // Empty name
			},
		},
	}
	mockPM := &pkgmanager.MoqPackageManager{
		GetInfoFunc: func() (pkgmanager.PackageManagerInfo, error) {
			return pkgmanager.PackageManagerInfo{Name: "apt"}, nil
		},
	}

	resolver, err := packageresolver.NewResolver(mappings, mockPM)
	require.NoError(t, err)

	result, err := resolver.Resolve("test-package", "")

	require.NoError(t, err)
	require.Equal(t, "test-package", result.Name)
	require.Nil(t, result.VersionConstraints)
}

func Test_Resolve_ParsesVersionConstraints_WhenProvided(t *testing.T) {
	mappings := &packageresolver.PackageMappingCollection{
		Packages: map[string]packageresolver.PackageMapping{
			"nodejs": {
				"apt": {Name: "nodejs"},
			},
		},
	}
	mockPM := &pkgmanager.MoqPackageManager{
		GetInfoFunc: func() (pkgmanager.PackageManagerInfo, error) {
			return pkgmanager.PackageManagerInfo{Name: "apt"}, nil
		},
	}

	resolver, err := packageresolver.NewResolver(mappings, mockPM)
	require.NoError(t, err)

	result, err := resolver.Resolve("nodejs", ">=16.0.0")

	require.NoError(t, err)
	require.Equal(t, "nodejs", result.Name)
	require.NotNil(t, result.VersionConstraints)

	// Test that the constraint works as expected
	version160, _ := semver.NewVersion("16.0.0")
	version140, _ := semver.NewVersion("14.0.0")
	require.True(t, result.VersionConstraints.Check(version160))
	require.False(t, result.VersionConstraints.Check(version140))
}

func Test_Resolve_ReturnsError_WhenVersionConstraintIsInvalid(t *testing.T) {
	mappings := &packageresolver.PackageMappingCollection{
		Packages: map[string]packageresolver.PackageMapping{
			"nodejs": {
				"apt": {Name: "nodejs"},
			},
		},
	}
	mockPM := &pkgmanager.MoqPackageManager{
		GetInfoFunc: func() (pkgmanager.PackageManagerInfo, error) {
			return pkgmanager.PackageManagerInfo{Name: "apt"}, nil
		},
	}

	resolver, err := packageresolver.NewResolver(mappings, mockPM)
	require.NoError(t, err)

	result, err := resolver.Resolve("nodejs", "invalid-version-constraint")

	require.Error(t, err)
	require.Equal(t, pkgmanager.RequestedPackageInfo{}, result)
	require.Contains(t, err.Error(), "invalid version constraint string")
	require.Contains(t, err.Error(), "invalid-version-constraint")
	require.Contains(t, err.Error(), "nodejs")
}

func Test_Resolve_ParsesComplexVersionConstraints(t *testing.T) {
	testCases := []struct {
		name              string
		constraint        string
		versionToTest     string
		expectedSatisfied bool
	}{
		{
			name:              "simple greater than constraint",
			constraint:        ">1.0.0",
			versionToTest:     "1.1.0",
			expectedSatisfied: true,
		},
		{
			name:              "simple greater than constraint not satisfied",
			constraint:        ">1.0.0",
			versionToTest:     "0.9.0",
			expectedSatisfied: false,
		},
		{
			name:              "range constraint satisfied",
			constraint:        ">=1.0.0, <2.0.0",
			versionToTest:     "1.5.0",
			expectedSatisfied: true,
		},
		{
			name:              "range constraint not satisfied",
			constraint:        ">=1.0.0, <2.0.0",
			versionToTest:     "2.1.0",
			expectedSatisfied: false,
		},
		{
			name:              "OR constraint satisfied by first part",
			constraint:        "<1.0.0 || >2.0.0",
			versionToTest:     "0.5.0",
			expectedSatisfied: true,
		},
		{
			name:              "OR constraint satisfied by second part",
			constraint:        "<1.0.0 || >2.0.0",
			versionToTest:     "3.0.0",
			expectedSatisfied: true,
		},
		{
			name:              "OR constraint not satisfied",
			constraint:        "<1.0.0 || >2.0.0",
			versionToTest:     "1.5.0",
			expectedSatisfied: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mappings := &packageresolver.PackageMappingCollection{
				Packages: map[string]packageresolver.PackageMapping{
					"test-package": {
						"apt": {Name: "test-pkg"},
					},
				},
			}
			mockPM := &pkgmanager.MoqPackageManager{
				GetInfoFunc: func() (pkgmanager.PackageManagerInfo, error) {
					return pkgmanager.PackageManagerInfo{Name: "apt"}, nil
				},
			}

			resolver, err := packageresolver.NewResolver(mappings, mockPM)
			require.NoError(t, err)

			result, err := resolver.Resolve("test-package", tc.constraint)

			require.NoError(t, err)
			require.NotNil(t, result.VersionConstraints)

			testVersion, err := semver.NewVersion(tc.versionToTest)
			require.NoError(t, err)

			satisfied := result.VersionConstraints.Check(testVersion)
			require.Equal(t, tc.expectedSatisfied, satisfied)
		})
	}
}

func Test_Resolve_HandlesMultiplePackageManagers(t *testing.T) {
	testCases := []struct {
		name                string
		packageManagerName  string
		expectedPackageName string
	}{
		{
			name:                "resolves for apt",
			packageManagerName:  "apt",
			expectedPackageName: "nodejs",
		},
		{
			name:                "resolves for brew",
			packageManagerName:  "brew",
			expectedPackageName: "node",
		},
		{
			name:                "resolves for dnf",
			packageManagerName:  "dnf",
			expectedPackageName: "nodejs",
		},
		{
			name:                "resolves for pacman",
			packageManagerName:  "pacman",
			expectedPackageName: "nodejs",
		},
	}

	mappings := &packageresolver.PackageMappingCollection{
		Packages: map[string]packageresolver.PackageMapping{
			"nodejs": {
				"apt":    {Name: "nodejs"},
				"brew":   {Name: "node"},
				"dnf":    {Name: "nodejs"},
				"pacman": {Name: "nodejs"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockPM := &pkgmanager.MoqPackageManager{
				GetInfoFunc: func() (pkgmanager.PackageManagerInfo, error) {
					return pkgmanager.PackageManagerInfo{Name: tc.packageManagerName}, nil
				},
			}

			resolver, err := packageresolver.NewResolver(mappings, mockPM)
			require.NoError(t, err)

			result, err := resolver.Resolve("nodejs", "")

			require.NoError(t, err)
			require.Equal(t, tc.expectedPackageName, result.Name)
		})
	}
}

func Test_Resolve_WorksWithRealWorldStructure(t *testing.T) {
	// This test uses the actual structure from packagemap.yaml
	mappings := &packageresolver.PackageMappingCollection{
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
			"zsh": {
				"apt":  {Name: "zsh"},
				"brew": {Name: "zsh"},
			},
		},
	}

	testCases := []struct {
		packageManagerName  string
		packageCode         string
		expectedPackageName string
	}{
		{"apt", "git", "git"},
		{"brew", "git", "git"},
		{"apt", "gpg", "gnupg2"},
		{"brew", "gpg", "gnupg"},
		{"dnf", "gpg", "gnupg2"},
		{"apt", "neovim", "neovim"},
		{"brew", "neovim", "neovim"},
		{"apt", "zsh", "zsh"},
		{"brew", "zsh", "zsh"},
	}

	for _, tc := range testCases {
		t.Run("resolves "+tc.packageCode+" for "+tc.packageManagerName, func(t *testing.T) {
			mockPM := &pkgmanager.MoqPackageManager{
				GetInfoFunc: func() (pkgmanager.PackageManagerInfo, error) {
					return pkgmanager.PackageManagerInfo{Name: tc.packageManagerName}, nil
				},
			}

			resolver, err := packageresolver.NewResolver(mappings, mockPM)
			require.NoError(t, err)

			result, err := resolver.Resolve(tc.packageCode, "")

			require.NoError(t, err)
			require.Equal(t, tc.expectedPackageName, result.Name)
		})
	}
}

func Test_Resolve_HandlesPackageWithVersionConstraints_UsingRealWorldStructure(t *testing.T) {
	mappings := &packageresolver.PackageMappingCollection{
		Packages: map[string]packageresolver.PackageMapping{
			"git": {
				"apt":  {Name: "git"},
				"brew": {Name: "git"},
			},
		},
	}

	mockPM := &pkgmanager.MoqPackageManager{
		GetInfoFunc: func() (pkgmanager.PackageManagerInfo, error) {
			return pkgmanager.PackageManagerInfo{Name: "apt"}, nil
		},
	}

	resolver, err := packageresolver.NewResolver(mappings, mockPM)
	require.NoError(t, err)

	result, err := resolver.Resolve("git", ">=2.0.0")

	require.NoError(t, err)
	require.Equal(t, "git", result.Name)
	require.NotNil(t, result.VersionConstraints)

	// Verify constraint works
	version200, _ := semver.NewVersion("2.0.0")
	version100, _ := semver.NewVersion("1.0.0")
	require.True(t, result.VersionConstraints.Check(version200))
	require.False(t, result.VersionConstraints.Check(version100))
}
