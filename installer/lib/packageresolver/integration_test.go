package packageresolver_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/MrPointer/dotfiles/installer/lib/packageresolver"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func Test_LoadPackageMappings_CanLoadFromActualFile(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "test-config.yaml")

	// Use the exact structure from embedded packagemap.yaml
	configContent := `packages:
  git:
    apt:
      name: git
    brew:
      name: git
  gpg:
    apt:
      name: gnupg2
    brew:
      name: gnupg
    dnf:
      name: gnupg2`

	err := os.WriteFile(configFile, []byte(configContent), 0644)
	require.NoError(t, err)

	v := viper.New()
	mappings, err := packageresolver.LoadPackageMappings(v, configFile)

	require.NoError(t, err)
	require.NotNil(t, mappings)
	require.NotNil(t, mappings.Packages)
	require.Len(t, mappings.Packages, 2)

	// Verify packages were loaded correctly from file
	gitMapping := mappings.Packages["git"]
	require.Equal(t, "git", gitMapping["apt"].Name)
	require.Equal(t, "git", gitMapping["brew"].Name)

	gpgMapping := mappings.Packages["gpg"]
	require.Equal(t, "gnupg2", gpgMapping["apt"].Name)
	require.Equal(t, "gnupg", gpgMapping["brew"].Name)
	require.Equal(t, "gnupg2", gpgMapping["dnf"].Name)
}

func Test_LoadPackageMappings_ReturnsError_WhenFileDoesNotExist(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	v := viper.New()
	nonExistentFile := "/non/existent/file.yaml"

	mappings, err := packageresolver.LoadPackageMappings(v, nonExistentFile)

	require.Error(t, err)
	require.Nil(t, mappings)
	require.Contains(t, err.Error(), "error reading package map file")
	require.Contains(t, err.Error(), nonExistentFile)
}

func Test_LoadPackageMappings_ReturnsError_WhenFileHasInvalidYAML(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "invalid.yaml")

	invalidYAML := `packages:
  git:
    apt:
      name: git
    invalid_yaml: [unclosed array`

	err := os.WriteFile(configFile, []byte(invalidYAML), 0644)
	require.NoError(t, err)

	v := viper.New()
	mappings, err := packageresolver.LoadPackageMappings(v, configFile)

	require.Error(t, err)
	require.Nil(t, mappings)
	require.Contains(t, err.Error(), "error reading package map file")
}

func Test_LoadPackageMappings_CanLoadFromEmbeddedConfig(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	v := viper.New()
	mappings, err := packageresolver.LoadPackageMappings(v, "")

	require.NoError(t, err)
	require.NotNil(t, mappings)
	require.NotNil(t, mappings.Packages)

	// The embedded config should contain the packages from packagemap.yaml
	// We don't assert specific content since it might change, but verify it loads
}

func Test_LoadPackageMappings_HandlesLargeConfigFile(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "large-config.yaml")

	// Generate a large config with many packages
	configContent := "packages:\n"
	packageManagers := []string{"apt", "brew", "dnf", "pacman"}

	// Create 50 packages (reasonable size for performance testing)
	for i := 0; i < 50; i++ {
		pkgName := fmt.Sprintf("package%02d", i)
		configContent += "  " + pkgName + ":\n"

		for _, pm := range packageManagers {
			configContent += "    " + pm + ":\n"
			configContent += "      name: " + pm + "-pkg-" + fmt.Sprintf("%02d", i) + "\n"
		}
	}

	err := os.WriteFile(configFile, []byte(configContent), 0644)
	require.NoError(t, err)

	v := viper.New()
	mappings, err := packageresolver.LoadPackageMappings(v, configFile)

	require.NoError(t, err)
	require.NotNil(t, mappings)

	require.Len(t, mappings.Packages, 50)

	// Verify a few random packages to ensure structure is correct
	pkg00 := mappings.Packages["package00"]
	require.Len(t, pkg00, len(packageManagers))
	require.Equal(t, "apt-pkg-00", pkg00["apt"].Name)
	require.Equal(t, "brew-pkg-00", pkg00["brew"].Name)

	pkg49 := mappings.Packages["package49"]
	require.Len(t, pkg49, len(packageManagers))
	require.Equal(t, "apt-pkg-49", pkg49["apt"].Name)
}
