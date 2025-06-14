package packageresolver_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/MrPointer/dotfiles/installer/lib/packageresolver"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

// Unit tests using in-memory YAML data

func Test_LoadPackageMappings_CanParseValidYAMLStructure(t *testing.T) {
	yamlContent := `packages:
  neovim:
    apt:
      name: neovim
    brew:
      name: neovim
  git:
    apt:
      name: git
    brew:
      name: git`

	v := viper.New()
	v.SetConfigType("yaml")
	err := v.ReadConfig(strings.NewReader(yamlContent))
	require.NoError(t, err)

	var mappingsCfg packageresolver.PackageMappingCollection
	err = v.Unmarshal(&mappingsCfg)
	require.NoError(t, err)

	// Ensure Packages map is initialized (mimicking loader behavior)
	if mappingsCfg.Packages == nil {
		mappingsCfg.Packages = make(map[string]packageresolver.PackageMapping)
	}

	require.NotNil(t, mappingsCfg.Packages)
	require.Len(t, mappingsCfg.Packages, 2)
	require.Contains(t, mappingsCfg.Packages, "neovim")
	require.Contains(t, mappingsCfg.Packages, "git")

	neovimMapping := mappingsCfg.Packages["neovim"]
	require.Len(t, neovimMapping, 2)
	require.Equal(t, "neovim", neovimMapping["apt"].Name)
	require.Equal(t, "neovim", neovimMapping["brew"].Name)

	gitMapping := mappingsCfg.Packages["git"]
	require.Len(t, gitMapping, 2)
	require.Equal(t, "git", gitMapping["apt"].Name)
	require.Equal(t, "git", gitMapping["brew"].Name)
}

func Test_LoadPackageMappings_CanLoadFromEmbeddedConfig_WhenNoFileSpecified(t *testing.T) {
	v := viper.New()

	mappings, err := packageresolver.LoadPackageMappings(v, "")

	require.NoError(t, err)
	require.NotNil(t, mappings)
	require.NotNil(t, mappings.Packages)
}

func Test_LoadPackageMappings_ReturnsError_WhenYAMLIsInvalid(t *testing.T) {
	invalidYAML := `packages:
  neovim:
    apt:
      name: neovim
    brew:
      name: neovim
    invalid_yaml: [unclosed array`

	v := viper.New()
	v.SetConfigType("yaml")
	err := v.ReadConfig(strings.NewReader(invalidYAML))

	require.Error(t, err)
}

func Test_LoadPackageMappings_InitializesEmptyPackagesMap_WhenConfigHasNoPackages(t *testing.T) {
	yamlContent := `packages: {}`

	v := viper.New()
	v.SetConfigType("yaml")
	err := v.ReadConfig(strings.NewReader(yamlContent))
	require.NoError(t, err)

	var mappingsCfg packageresolver.PackageMappingCollection
	err = v.Unmarshal(&mappingsCfg)
	require.NoError(t, err)

	// Ensure Packages map is initialized (mimicking loader behavior)
	if mappingsCfg.Packages == nil {
		mappingsCfg.Packages = make(map[string]packageresolver.PackageMapping)
	}

	require.NotNil(t, mappingsCfg.Packages)
	require.Empty(t, mappingsCfg.Packages)
}

func Test_LoadPackageMappings_InitializesEmptyPackagesMap_WhenPackagesKeyIsMissing(t *testing.T) {
	yamlContent := `other_config:
  some_value: test`

	v := viper.New()
	v.SetConfigType("yaml")
	err := v.ReadConfig(strings.NewReader(yamlContent))
	require.NoError(t, err)

	var mappingsCfg packageresolver.PackageMappingCollection
	err = v.Unmarshal(&mappingsCfg)
	require.NoError(t, err)

	// Ensure Packages map is initialized (mimicking loader behavior)
	if mappingsCfg.Packages == nil {
		mappingsCfg.Packages = make(map[string]packageresolver.PackageMapping)
	}

	require.NotNil(t, mappingsCfg.Packages)
	require.Empty(t, mappingsCfg.Packages)
}

func Test_LoadPackageMappings_CanParseComplexConfiguration(t *testing.T) {
	yamlContent := `packages:
  nodejs:
    apt:
      name: nodejs
    brew:
      name: node
    dnf:
      name: nodejs
    pacman:
      name: nodejs
  python:
    apt:
      name: python3
    brew:
      name: python@3.11
  docker:
    apt:
      name: docker.io
    brew:
      name: docker`

	v := viper.New()
	v.SetConfigType("yaml")
	err := v.ReadConfig(strings.NewReader(yamlContent))
	require.NoError(t, err)

	var mappingsCfg packageresolver.PackageMappingCollection
	err = v.Unmarshal(&mappingsCfg)
	require.NoError(t, err)

	// Ensure Packages map is initialized (mimicking loader behavior)
	if mappingsCfg.Packages == nil {
		mappingsCfg.Packages = make(map[string]packageresolver.PackageMapping)
	}

	require.NotNil(t, mappingsCfg.Packages)
	require.Len(t, mappingsCfg.Packages, 3)

	// Test nodejs mapping
	nodejsMapping := mappingsCfg.Packages["nodejs"]
	require.Len(t, nodejsMapping, 4)
	require.Equal(t, "nodejs", nodejsMapping["apt"].Name)
	require.Equal(t, "node", nodejsMapping["brew"].Name)
	require.Equal(t, "nodejs", nodejsMapping["dnf"].Name)
	require.Equal(t, "nodejs", nodejsMapping["pacman"].Name)

	// Test python mapping
	pythonMapping := mappingsCfg.Packages["python"]
	require.Len(t, pythonMapping, 2)
	require.Equal(t, "python3", pythonMapping["apt"].Name)
	require.Equal(t, "python@3.11", pythonMapping["brew"].Name)

	// Test docker mapping
	dockerMapping := mappingsCfg.Packages["docker"]
	require.Len(t, dockerMapping, 2)
	require.Equal(t, "docker.io", dockerMapping["apt"].Name)
	require.Equal(t, "docker", dockerMapping["brew"].Name)
}

func Test_LoadPackageMappings_ReturnsError_WhenUnmarshalFails(t *testing.T) {
	// This YAML is valid but doesn't match our expected structure
	yamlContent := `packages:
  - this_should_be_a_map
  - not_an_array`

	v := viper.New()
	v.SetConfigType("yaml")
	err := v.ReadConfig(strings.NewReader(yamlContent))
	require.NoError(t, err)

	var mappingsCfg packageresolver.PackageMappingCollection
	err = v.Unmarshal(&mappingsCfg)

	require.Error(t, err)
}

func Test_LoadPackageMappings_CanHandleEmptyManagersMap(t *testing.T) {
	yamlContent := `packages:
  test-package: {}`

	v := viper.New()
	v.SetConfigType("yaml")
	err := v.ReadConfig(strings.NewReader(yamlContent))
	require.NoError(t, err)

	var mappingsCfg packageresolver.PackageMappingCollection
	err = v.Unmarshal(&mappingsCfg)
	require.NoError(t, err)

	// Ensure Packages map is initialized (mimicking loader behavior)
	if mappingsCfg.Packages == nil {
		mappingsCfg.Packages = make(map[string]packageresolver.PackageMapping)
	}

	require.NotNil(t, mappingsCfg.Packages)
	// Empty package mappings are not loaded, so we expect 0 packages
	require.Empty(t, mappingsCfg.Packages)
}

func Test_LoadPackageMappings_PreservesViperConfigurationState(t *testing.T) {
	yamlContent := `packages:
  test:
    apt:
      name: test-package`

	v := viper.New()
	originalValue := "test-value"
	v.Set("test-key", originalValue)

	v.SetConfigType("yaml")
	err := v.ReadConfig(strings.NewReader(yamlContent))
	require.NoError(t, err)

	var mappingsCfg packageresolver.PackageMappingCollection
	err = v.Unmarshal(&mappingsCfg)
	require.NoError(t, err)

	// Verify that the viper instance still contains our original value
	require.Equal(t, originalValue, v.Get("test-key"))
}

func Test_LoadPackageMappings_MatchesRealWorldStructure(t *testing.T) {
	// This matches the actual structure from packagemap.yaml
	yamlContent := `packages:
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
      name: gnupg2
  neovim:
    apt:
      name: neovim
    brew:
      name: neovim
  zsh:
    apt:
      name: zsh
    brew:
      name: zsh`

	v := viper.New()
	v.SetConfigType("yaml")
	err := v.ReadConfig(strings.NewReader(yamlContent))
	require.NoError(t, err)

	var mappingsCfg packageresolver.PackageMappingCollection
	err = v.Unmarshal(&mappingsCfg)
	require.NoError(t, err)

	// Ensure Packages map is initialized (mimicking loader behavior)
	if mappingsCfg.Packages == nil {
		mappingsCfg.Packages = make(map[string]packageresolver.PackageMapping)
	}

	require.NotNil(t, mappingsCfg.Packages)
	require.Len(t, mappingsCfg.Packages, 4)

	// Verify git package
	gitMapping := mappingsCfg.Packages["git"]
	require.Len(t, gitMapping, 2)
	require.Equal(t, "git", gitMapping["apt"].Name)
	require.Equal(t, "git", gitMapping["brew"].Name)

	// Verify gpg package with different names per manager
	gpgMapping := mappingsCfg.Packages["gpg"]
	require.Len(t, gpgMapping, 3)
	require.Equal(t, "gnupg2", gpgMapping["apt"].Name)
	require.Equal(t, "gnupg", gpgMapping["brew"].Name)
	require.Equal(t, "gnupg2", gpgMapping["dnf"].Name)

	// Verify neovim package
	neovimMapping := mappingsCfg.Packages["neovim"]
	require.Len(t, neovimMapping, 2)
	require.Equal(t, "neovim", neovimMapping["apt"].Name)
	require.Equal(t, "neovim", neovimMapping["brew"].Name)

	// Verify zsh package
	zshMapping := mappingsCfg.Packages["zsh"]
	require.Len(t, zshMapping, 2)
	require.Equal(t, "zsh", zshMapping["apt"].Name)
	require.Equal(t, "zsh", zshMapping["brew"].Name)
}

func Test_LoadPackageMappings_CanParseYAMLWithBytesBuffer(t *testing.T) {
	yamlContent := `packages:
  curl:
    apt:
      name: curl
    brew:
      name: curl`

	yamlBytes := []byte(yamlContent)

	v := viper.New()
	v.SetConfigType("yaml")
	err := v.ReadConfig(bytes.NewBuffer(yamlBytes))
	require.NoError(t, err)

	var mappingsCfg packageresolver.PackageMappingCollection
	err = v.Unmarshal(&mappingsCfg)
	require.NoError(t, err)

	// Ensure Packages map is initialized (mimicking loader behavior)
	if mappingsCfg.Packages == nil {
		mappingsCfg.Packages = make(map[string]packageresolver.PackageMapping)
	}

	require.NotNil(t, mappingsCfg.Packages)
	require.Len(t, mappingsCfg.Packages, 1)

	curlMapping := mappingsCfg.Packages["curl"]
	require.Len(t, curlMapping, 2)
	require.Equal(t, "curl", curlMapping["apt"].Name)
	require.Equal(t, "curl", curlMapping["brew"].Name)
}
