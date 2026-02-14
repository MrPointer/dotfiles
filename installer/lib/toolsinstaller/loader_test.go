package toolsinstaller

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadToolsConfig_Embedded(t *testing.T) {
	cfg, err := LoadToolsConfig(viper.New(), "")
	require.NoError(t, err)
	require.NotNil(t, cfg)
	require.NotEmpty(t, cfg.Tools)

	for _, tool := range cfg.Tools {
		assert.NotEmpty(t, tool.Name, "tool name should not be empty")
		assert.NotEmpty(t, tool.Description, "tool description should not be empty")
	}
}

func TestLoadToolsConfig_CustomFile(t *testing.T) {
	content := `tools:
  - name: mytool
    description: "A custom tool"
  - name: anothertool
    description: "Another custom tool"
`
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "custom-tools.yaml")
	err := os.WriteFile(tmpFile, []byte(content), 0o644)
	require.NoError(t, err)

	cfg, err := LoadToolsConfig(viper.New(), tmpFile)
	require.NoError(t, err)
	require.NotNil(t, cfg)
	require.Len(t, cfg.Tools, 2)
	assert.Equal(t, "mytool", cfg.Tools[0].Name)
	assert.Equal(t, "A custom tool", cfg.Tools[0].Description)
	assert.Equal(t, "anothertool", cfg.Tools[1].Name)
	assert.Equal(t, "Another custom tool", cfg.Tools[1].Description)
}

func TestLoadToolsConfig_NonExistentFile(t *testing.T) {
	cfg, err := LoadToolsConfig(viper.New(), "/nonexistent/path/tools.yaml")
	require.Error(t, err)
	assert.Nil(t, cfg)
}

func TestLoadToolsConfig_EmptyConfig(t *testing.T) {
	content := `tools: []
`
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "empty-tools.yaml")
	err := os.WriteFile(tmpFile, []byte(content), 0o644)
	require.NoError(t, err)

	cfg, err := LoadToolsConfig(viper.New(), tmpFile)
	require.NoError(t, err)
	require.NotNil(t, cfg)
	assert.Empty(t, cfg.Tools)
}
