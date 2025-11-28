package packageresolver_test

import (
	"errors"
	"testing"

	"github.com/Masterminds/semver"
	"github.com/MrPointer/dotfiles/installer/lib/compatibility"
	"github.com/MrPointer/dotfiles/installer/lib/packageresolver"
	"github.com/MrPointer/dotfiles/installer/lib/pkgmanager"
	"github.com/stretchr/testify/require"
)

// createTestSystemInfo creates a basic SystemInfo for testing.
func createTestSystemInfo() *compatibility.SystemInfo {
	return &compatibility.SystemInfo{
		OSName:     "linux",
		DistroName: "fedora",
		Arch:       "amd64",
	}
}

func Test_NewResolver_CanCreateResolver_WithValidInputs(t *testing.T) {
	mappings := &packageresolver.PackageMappingCollection{
		Packages: make(map[string]packageresolver.PackageMapping),
	}
	mockPM := &pkgmanager.MoqPackageManager{
		GetInfoFunc: func() (pkgmanager.PackageManagerInfo, error) {
			return pkgmanager.PackageManagerInfo{Name: "apt"}, nil
		},
	}

	resolver, err := packageresolver.NewResolver(mappings, mockPM, createTestSystemInfo())

	require.NoError(t, err)
	require.NotNil(t, resolver)
}

func Test_NewResolver_ReturnsError_WhenMappingsIsNil(t *testing.T) {
	mockPM := &pkgmanager.MoqPackageManager{
		GetInfoFunc: func() (pkgmanager.PackageManagerInfo, error) {
			return pkgmanager.PackageManagerInfo{Name: "apt"}, nil
		},
	}

	resolver, err := packageresolver.NewResolver(nil, mockPM, createTestSystemInfo())

	require.Error(t, err)
	require.Nil(t, resolver)
	require.Contains(t, err.Error(), "package mappings cannot be nil")
}

func Test_NewResolver_ReturnsError_WhenPackageManagerIsNil(t *testing.T) {
	mappings := &packageresolver.PackageMappingCollection{
		Packages: make(map[string]packageresolver.PackageMapping),
	}

	resolver, err := packageresolver.NewResolver(mappings, nil, createTestSystemInfo())

	require.Error(t, err)
	require.Nil(t, resolver)
	require.Contains(t, err.Error(), "package manager cannot be nil")
}

func Test_NewResolver_ReturnsError_WhenSystemInfoIsNil(t *testing.T) {
	mappings := &packageresolver.PackageMappingCollection{
		Packages: make(map[string]packageresolver.PackageMapping),
	}
	mockPM := &pkgmanager.MoqPackageManager{
		GetInfoFunc: func() (pkgmanager.PackageManagerInfo, error) {
			return pkgmanager.PackageManagerInfo{Name: "apt"}, nil
		},
	}

	resolver, err := packageresolver.NewResolver(mappings, mockPM, nil)

	require.Error(t, err)
	require.Nil(t, resolver)
	require.Contains(t, err.Error(), "system info cannot be nil")
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

	resolver, err := packageresolver.NewResolver(mappings, mockPM, createTestSystemInfo())

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

	resolver, err := packageresolver.NewResolver(mappings, mockPM, createTestSystemInfo())

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
			name:   "uses zypper name directly",
			pmName: "zypper",
		},
		{
			name:   "uses pacman name directly",
			pmName: "pacman",
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

			resolver, err := packageresolver.NewResolver(mappings, mockPM, createTestSystemInfo())

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

	resolver, err := packageresolver.NewResolver(mappings, mockPM, createTestSystemInfo())
	require.NoError(t, err)

	result, err := resolver.Resolve("", "")

	require.Error(t, err)
	require.Equal(t, pkgmanager.RequestedPackageInfo{}, result)
	require.Contains(t, err.Error(), "generic package code cannot be empty")
}

func Test_Resolve_UsesManagerSpecificName_WhenAvailable(t *testing.T) {
	mappings := &packageresolver.PackageMappingCollection{
		Packages: map[string]packageresolver.PackageMapping{
			"neovim": packageresolver.PackageMapping{
				"apt": packageresolver.ManagerSpecificMapping{
					Name: "neovim",
				},
			},
		},
	}
	mockPM := &pkgmanager.MoqPackageManager{
		GetInfoFunc: func() (pkgmanager.PackageManagerInfo, error) {
			return pkgmanager.PackageManagerInfo{Name: "apt"}, nil
		},
	}

	resolver, err := packageresolver.NewResolver(mappings, mockPM, createTestSystemInfo())
	require.NoError(t, err)

	result, err := resolver.Resolve("neovim", "")

	require.NoError(t, err)
	require.Equal(t, "neovim", result.Name)
	require.Equal(t, "", result.Type)
	require.Nil(t, result.VersionConstraints)
}

func Test_Resolve_ReturnsError_WhenManagerSpecificNameNotFound(t *testing.T) {
	mappings := &packageresolver.PackageMappingCollection{
		Packages: map[string]packageresolver.PackageMapping{
			"nodejs": packageresolver.PackageMapping{
				"brew": packageresolver.ManagerSpecificMapping{
					Name: "node",
				},
			},
		},
	}
	mockPM := &pkgmanager.MoqPackageManager{
		GetInfoFunc: func() (pkgmanager.PackageManagerInfo, error) {
			return pkgmanager.PackageManagerInfo{Name: "apt"}, nil
		},
	}

	resolver, err := packageresolver.NewResolver(mappings, mockPM, createTestSystemInfo())
	require.NoError(t, err)

	result, err := resolver.Resolve("nodejs", "")

	require.Error(t, err)
	require.Equal(t, pkgmanager.RequestedPackageInfo{}, result)
	require.Contains(t, err.Error(), "no package mapping found")
	require.Contains(t, err.Error(), "nodejs")
	require.Contains(t, err.Error(), "apt")
}

func Test_Resolve_ReturnsError_WhenNoMappingFound(t *testing.T) {
	mappings := &packageresolver.PackageMappingCollection{
		Packages: make(map[string]packageresolver.PackageMapping),
	}
	mockPM := &pkgmanager.MoqPackageManager{
		GetInfoFunc: func() (pkgmanager.PackageManagerInfo, error) {
			return pkgmanager.PackageManagerInfo{Name: "apt"}, nil
		},
	}

	resolver, err := packageresolver.NewResolver(mappings, mockPM, createTestSystemInfo())
	require.NoError(t, err)

	result, err := resolver.Resolve("unknown-package", "")

	require.Error(t, err)
	require.Equal(t, pkgmanager.RequestedPackageInfo{}, result)
	require.Contains(t, err.Error(), "no package mapping found")
	require.Contains(t, err.Error(), "unknown-package")
}

func Test_Resolve_ReturnsError_WhenManagerSpecificNameIsEmpty(t *testing.T) {
	mappings := &packageresolver.PackageMappingCollection{
		Packages: map[string]packageresolver.PackageMapping{
			"test-package": packageresolver.PackageMapping{
				"apt": packageresolver.ManagerSpecificMapping{
					Name: "",
				},
			},
		},
	}
	mockPM := &pkgmanager.MoqPackageManager{
		GetInfoFunc: func() (pkgmanager.PackageManagerInfo, error) {
			return pkgmanager.PackageManagerInfo{Name: "apt"}, nil
		},
	}

	resolver, err := packageresolver.NewResolver(mappings, mockPM, createTestSystemInfo())
	require.NoError(t, err)

	result, err := resolver.Resolve("test-package", "")

	require.Error(t, err)
	require.Equal(t, pkgmanager.RequestedPackageInfo{}, result)
	require.Contains(t, err.Error(), "no package mapping found")
	require.Contains(t, err.Error(), "test-package")
	require.Contains(t, err.Error(), "apt")
}

func Test_Resolve_ParsesVersionConstraints_WhenProvided(t *testing.T) {
	mappings := &packageresolver.PackageMappingCollection{
		Packages: map[string]packageresolver.PackageMapping{
			"nodejs": packageresolver.PackageMapping{
				"apt": packageresolver.ManagerSpecificMapping{
					Name: "nodejs",
				},
			},
		},
	}
	mockPM := &pkgmanager.MoqPackageManager{
		GetInfoFunc: func() (pkgmanager.PackageManagerInfo, error) {
			return pkgmanager.PackageManagerInfo{Name: "apt"}, nil
		},
	}

	resolver, err := packageresolver.NewResolver(mappings, mockPM, createTestSystemInfo())
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
			"nodejs": packageresolver.PackageMapping{
				"apt": packageresolver.ManagerSpecificMapping{
					Name: "nodejs",
				},
			},
		},
	}
	mockPM := &pkgmanager.MoqPackageManager{
		GetInfoFunc: func() (pkgmanager.PackageManagerInfo, error) {
			return pkgmanager.PackageManagerInfo{Name: "apt"}, nil
		},
	}

	resolver, err := packageresolver.NewResolver(mappings, mockPM, createTestSystemInfo())
	require.NoError(t, err)

	result, err := resolver.Resolve("nodejs", "invalid-version-constraint")

	require.Error(t, err)
	require.Equal(t, pkgmanager.RequestedPackageInfo{}, result)
	require.Contains(t, err.Error(), "invalid version constraint string")
}

func Test_Resolve_ParsesComplexVersionConstraints(t *testing.T) {
	testCases := []struct {
		name        string
		constraint  string
		description string
	}{
		{
			name:        "range constraint",
			constraint:  ">=16.0.0, <20.0.0",
			description: "should parse range constraints",
		},
		{
			name:        "or constraint",
			constraint:  "^16.0.0 || ^18.0.0",
			description: "should parse OR constraints",
		},
		{
			name:        "exact version",
			constraint:  "18.17.1",
			description: "should parse exact version constraints",
		},
		{
			name:        "tilde constraint",
			constraint:  "~16.14.0",
			description: "should parse tilde constraints",
		},
		{
			name:        "caret constraint",
			constraint:  "^16.0.0",
			description: "should parse caret constraints",
		},
		{
			name:        "complex mixed constraint",
			constraint:  ">=14.0.0, <16.0.0 || >=16.14.0, <18.0.0 || >=18.12.0",
			description: "should parse complex mixed constraints",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mappings := &packageresolver.PackageMappingCollection{
				Packages: map[string]packageresolver.PackageMapping{
					"test-package": packageresolver.PackageMapping{
						"apt": packageresolver.ManagerSpecificMapping{
							Name: "test-package",
						},
					},
				},
			}
			mockPM := &pkgmanager.MoqPackageManager{
				GetInfoFunc: func() (pkgmanager.PackageManagerInfo, error) {
					return pkgmanager.PackageManagerInfo{Name: "apt"}, nil
				},
			}

			resolver, err := packageresolver.NewResolver(mappings, mockPM, createTestSystemInfo())
			require.NoError(t, err)

			result, err := resolver.Resolve("test-package", tc.constraint)

			require.NoError(t, err, tc.description)
			require.Equal(t, "test-package", result.Name)
			require.NotNil(t, result.VersionConstraints, tc.description)

			// Verify the constraint was parsed correctly by creating the same constraint
			expectedConstraints, err := semver.NewConstraint(tc.constraint)
			require.NoError(t, err)

			// Test with a sample version to ensure constraints work the same way
			testVersion, _ := semver.NewVersion("16.0.0")
			require.Equal(t, expectedConstraints.Check(testVersion), result.VersionConstraints.Check(testVersion), tc.description)
		})
	}
}

func Test_Resolve_HandlesMultiplePackageManagers(t *testing.T) {
	mappings := &packageresolver.PackageMappingCollection{
		Packages: map[string]packageresolver.PackageMapping{
			"nodejs": packageresolver.PackageMapping{
				"apt": packageresolver.ManagerSpecificMapping{
					Name: "nodejs",
				},
				"brew": packageresolver.ManagerSpecificMapping{
					Name: "node",
				},
				"dnf": packageresolver.ManagerSpecificMapping{
					Name: "nodejs",
				},
			},
		},
	}

	testCases := []struct {
		name                string
		packageManagerName  string
		expectedPackageName string
	}{
		{
			name:                "apt manager uses nodejs",
			packageManagerName:  "apt",
			expectedPackageName: "nodejs",
		},
		{
			name:                "brew manager uses node",
			packageManagerName:  "brew",
			expectedPackageName: "node",
		},
		{
			name:                "dnf manager uses nodejs",
			packageManagerName:  "dnf",
			expectedPackageName: "nodejs",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockPM := &pkgmanager.MoqPackageManager{
				GetInfoFunc: func() (pkgmanager.PackageManagerInfo, error) {
					return pkgmanager.PackageManagerInfo{Name: tc.packageManagerName}, nil
				},
			}

			resolver, err := packageresolver.NewResolver(mappings, mockPM, createTestSystemInfo())
			require.NoError(t, err)

			result, err := resolver.Resolve("nodejs", "")

			require.NoError(t, err)
			require.Equal(t, tc.expectedPackageName, result.Name)
			require.Equal(t, "", result.Type)
			require.Nil(t, result.VersionConstraints)
		})
	}
}

func Test_Resolve_WorksWithRealWorldStructure(t *testing.T) {
	// This test uses a structure that closely resembles the actual packagemap.yaml
	mappings := &packageresolver.PackageMappingCollection{
		Packages: map[string]packageresolver.PackageMapping{
			"git": packageresolver.PackageMapping{
				"apt": packageresolver.ManagerSpecificMapping{
					Name: "git",
				},
				"brew": packageresolver.ManagerSpecificMapping{
					Name: "git",
				},
				"dnf": packageresolver.ManagerSpecificMapping{
					Name: "git",
				},
			},
			"neovim": packageresolver.PackageMapping{
				"apt": packageresolver.ManagerSpecificMapping{
					Name: "neovim",
				},
				"brew": packageresolver.ManagerSpecificMapping{
					Name: "neovim",
				},
				"dnf": packageresolver.ManagerSpecificMapping{
					Name: "neovim",
				},
			},
		},
	}

	testCases := []struct {
		packageCode        string
		packageManagerName string
		expectedName       string
	}{
		{"git", "apt", "git"},
		{"git", "brew", "git"},
		{"git", "dnf", "git"},
		{"neovim", "apt", "neovim"},
		{"neovim", "brew", "neovim"},
		{"neovim", "dnf", "neovim"},
	}

	for _, tc := range testCases {
		t.Run("resolves "+tc.packageCode+" for "+tc.packageManagerName, func(t *testing.T) {
			mockPM := &pkgmanager.MoqPackageManager{
				GetInfoFunc: func() (pkgmanager.PackageManagerInfo, error) {
					return pkgmanager.PackageManagerInfo{Name: tc.packageManagerName}, nil
				},
			}

			resolver, err := packageresolver.NewResolver(mappings, mockPM, createTestSystemInfo())
			require.NoError(t, err)

			result, err := resolver.Resolve(tc.packageCode, "")

			require.NoError(t, err)
			require.Equal(t, tc.expectedName, result.Name)
			require.Equal(t, "", result.Type)
			require.Nil(t, result.VersionConstraints)
		})
	}
}

func Test_Resolve_ParsesVersionConstraintsCorrectlyForMultiplePackages(t *testing.T) {
	mappings := &packageresolver.PackageMappingCollection{
		Packages: map[string]packageresolver.PackageMapping{
			"git": packageresolver.PackageMapping{
				"apt": packageresolver.ManagerSpecificMapping{
					Name: "git",
				},
			},
		},
	}
	mockPM := &pkgmanager.MoqPackageManager{
		GetInfoFunc: func() (pkgmanager.PackageManagerInfo, error) {
			return pkgmanager.PackageManagerInfo{Name: "apt"}, nil
		},
	}

	resolver, err := packageresolver.NewResolver(mappings, mockPM, createTestSystemInfo())
	require.NoError(t, err)

	result, err := resolver.Resolve("git", ">=2.0.0")

	require.NoError(t, err)
	require.Equal(t, "git", result.Name)
	require.Equal(t, "", result.Type)
	require.NotNil(t, result.VersionConstraints)

	// Test that the constraint works as expected
	version200, _ := semver.NewVersion("2.0.0")
	version100, _ := semver.NewVersion("1.0.0")
	require.True(t, result.VersionConstraints.Check(version200))
	require.False(t, result.VersionConstraints.Check(version100))
}

func Test_Resolve_RespectsTypeInfo_WhenProvided(t *testing.T) {
	mappings := &packageresolver.PackageMappingCollection{
		Packages: map[string]packageresolver.PackageMapping{
			"build-tools": packageresolver.PackageMapping{
				"dnf": packageresolver.ManagerSpecificMapping{
					Name: "Development Tools",
					Type: "group",
				},
			},
		},
	}
	mockPM := &pkgmanager.MoqPackageManager{
		GetInfoFunc: func() (pkgmanager.PackageManagerInfo, error) {
			return pkgmanager.PackageManagerInfo{Name: "dnf"}, nil
		},
	}

	resolver, err := packageresolver.NewResolver(mappings, mockPM, createTestSystemInfo())
	require.NoError(t, err)

	result, err := resolver.Resolve("build-tools", "")

	require.NoError(t, err)
	require.Equal(t, "Development Tools", result.Name)
	require.Equal(t, "group", result.Type)
	require.Nil(t, result.VersionConstraints)
}

func Test_Resolve_FallsBackToDirectMapping_WhenNoDistroSpecificFound(t *testing.T) {
	mappings := &packageresolver.PackageMappingCollection{
		Packages: map[string]packageresolver.PackageMapping{
			"git": packageresolver.PackageMapping{
				"dnf": packageresolver.ManagerSpecificMapping{
					Name: "git",
					Type: "",
				},
			},
		},
	}
	mockPM := &pkgmanager.MoqPackageManager{
		GetInfoFunc: func() (pkgmanager.PackageManagerInfo, error) {
			return pkgmanager.PackageManagerInfo{Name: "dnf"}, nil
		},
	}

	resolver, err := packageresolver.NewResolver(mappings, mockPM, createTestSystemInfo())
	require.NoError(t, err)

	result, err := resolver.Resolve("git", "")

	require.NoError(t, err)
	require.Equal(t, "git", result.Name)
	require.Equal(t, "", result.Type)
	require.Nil(t, result.VersionConstraints)
}

func Test_Resolve_ReturnsError_WhenSimpleStringPackageHasNoMapping(t *testing.T) {
	mappings := &packageresolver.PackageMappingCollection{
		Packages: map[string]packageresolver.PackageMapping{
			"simple-package": packageresolver.PackageMapping{
				"apt": packageresolver.ManagerSpecificMapping{
					Name: "apt-simple-package",
				},
				// No DNF mapping
			},
		},
	}

	// Test with DNF on any distro - should error since no mapping exists
	sysInfo := &compatibility.SystemInfo{
		OSName:     "linux",
		DistroName: "fedora",
		Arch:       "amd64",
	}
	mockPM := &pkgmanager.MoqPackageManager{
		GetInfoFunc: func() (pkgmanager.PackageManagerInfo, error) {
			return pkgmanager.PackageManagerInfo{Name: "dnf"}, nil
		},
	}

	resolver, err := packageresolver.NewResolver(mappings, mockPM, sysInfo)
	require.NoError(t, err)

	result, err := resolver.Resolve("simple-package", "")

	require.Error(t, err)
	require.Equal(t, pkgmanager.RequestedPackageInfo{}, result)
	require.Contains(t, err.Error(), "no package mapping found")
	require.Contains(t, err.Error(), "simple-package")
	require.Contains(t, err.Error(), "dnf")
}

func Test_Resolve_UsesDistroSpecificMapping_WhenAvailable(t *testing.T) {
	mappings := &packageresolver.PackageMappingCollection{
		Packages: map[string]packageresolver.PackageMapping{
			"development-tools": packageresolver.PackageMapping{
				"dnf": packageresolver.ManagerSpecificMapping{
					Type: "group",
					Name: map[string]any{
						"fedora": "development-tools",
						"centos": "Development Tools",
					},
				},
			},
		},
	}

	testCases := []struct {
		name         string
		distroName   string
		expectedName string
		expectedType string
	}{
		{
			name:         "uses fedora specific mapping",
			distroName:   "fedora",
			expectedName: "development-tools",
			expectedType: "group",
		},
		{
			name:         "uses centos specific mapping",
			distroName:   "centos",
			expectedName: "Development Tools",
			expectedType: "group",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sysInfo := &compatibility.SystemInfo{
				OSName:     "linux",
				DistroName: tc.distroName,
				Arch:       "amd64",
			}
			mockPM := &pkgmanager.MoqPackageManager{
				GetInfoFunc: func() (pkgmanager.PackageManagerInfo, error) {
					return pkgmanager.PackageManagerInfo{Name: "dnf"}, nil
				},
			}

			resolver, err := packageresolver.NewResolver(mappings, mockPM, sysInfo)
			require.NoError(t, err)

			result, err := resolver.Resolve("development-tools", "")

			require.NoError(t, err)
			require.Equal(t, tc.expectedName, result.Name)
			require.Equal(t, tc.expectedType, result.Type)
			require.Nil(t, result.VersionConstraints)
		})
	}
}

func Test_Resolve_ReturnsError_WhenDistroNotMappedForDistroSpecificPackage(t *testing.T) {
	mappings := &packageresolver.PackageMappingCollection{
		Packages: map[string]packageresolver.PackageMapping{
			"development-tools": packageresolver.PackageMapping{
				"dnf": packageresolver.ManagerSpecificMapping{
					Type: "group",
					Name: map[string]any{
						"fedora": "development-tools",
						"centos": "Development Tools",
					},
				},
			},
		},
	}

	// Test with a distro that doesn't have specific mapping (should fail)
	sysInfo := &compatibility.SystemInfo{
		OSName:     "linux",
		DistroName: "ubuntu", // No specific mapping for Ubuntu
		Arch:       "amd64",
	}
	mockPM := &pkgmanager.MoqPackageManager{
		GetInfoFunc: func() (pkgmanager.PackageManagerInfo, error) {
			return pkgmanager.PackageManagerInfo{Name: "dnf"}, nil
		},
	}

	resolver, err := packageresolver.NewResolver(mappings, mockPM, sysInfo)
	require.NoError(t, err)

	result, err := resolver.Resolve("development-tools", "")

	require.Error(t, err)
	require.Equal(t, pkgmanager.RequestedPackageInfo{}, result)
	require.Contains(t, err.Error(), "requires distro-specific mapping")
	require.Contains(t, err.Error(), "ubuntu")
}

func Test_Resolve_HandlesAllSupportedDistros(t *testing.T) {
	mappings := &packageresolver.PackageMappingCollection{
		Packages: map[string]packageresolver.PackageMapping{
			"development-tools": packageresolver.PackageMapping{
				"dnf": packageresolver.ManagerSpecificMapping{
					Type: "group",
					Name: map[string]any{
						"fedora": "development-tools",
						"centos": "Development Tools",
						"rhel":   "Development Tools",
					},
				},
			},
		},
	}

	testCases := []struct {
		distroName   string
		expectedName string
	}{
		{"fedora", "development-tools"},
		{"centos", "Development Tools"},
		{"rhel", "Development Tools"},
	}

	for _, tc := range testCases {
		t.Run("handles "+tc.distroName+" distro", func(t *testing.T) {
			sysInfo := &compatibility.SystemInfo{
				OSName:     "linux",
				DistroName: tc.distroName,
				Arch:       "amd64",
			}
			mockPM := &pkgmanager.MoqPackageManager{
				GetInfoFunc: func() (pkgmanager.PackageManagerInfo, error) {
					return pkgmanager.PackageManagerInfo{Name: "dnf"}, nil
				},
			}

			resolver, err := packageresolver.NewResolver(mappings, mockPM, sysInfo)
			require.NoError(t, err)

			result, err := resolver.Resolve("development-tools", "")

			require.NoError(t, err)
			require.Equal(t, tc.expectedName, result.Name)
			require.Equal(t, "group", result.Type)
			require.Nil(t, result.VersionConstraints)
		})
	}
}
