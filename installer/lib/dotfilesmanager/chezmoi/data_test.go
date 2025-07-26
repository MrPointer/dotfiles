package chezmoi_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/samber/mo"
	"github.com/stretchr/testify/require"

	"github.com/MrPointer/dotfiles/installer/lib/dotfilesmanager"
	"github.com/MrPointer/dotfiles/installer/lib/dotfilesmanager/chezmoi"
	"github.com/MrPointer/dotfiles/installer/lib/pkgmanager"
	"github.com/MrPointer/dotfiles/installer/utils"
	"github.com/MrPointer/dotfiles/installer/utils/httpclient"
	"github.com/MrPointer/dotfiles/installer/utils/osmanager"
)

func Test_Initialize_CreatesConfigDirectory_WhenDirectoryDoesNotExist(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "config")
	configFilePath := filepath.Join(configDir, "chezmoi.toml")
	cloneDir := filepath.Join(tempDir, "clone")

	fileSystem := utils.NewDefaultFileSystem()
	mockUserManager := &osmanager.MoqUserManager{}
	mockPackageManager := &pkgmanager.MoqPackageManager{}
	mockHTTPClient := &httpclient.MoqHTTPClient{}
	mockCommander := &utils.MoqCommander{}

	config := chezmoi.DefaultChezmoiConfig(configFilePath, cloneDir)
	manager := chezmoi.NewChezmoiManager(fileSystem, mockUserManager, mockCommander, mockPackageManager, mockHTTPClient, config)

	data := dotfilesmanager.DotfilesData{
		Email:         "test@example.com",
		FirstName:     "John",
		LastName:      "Doe",
		GpgSigningKey: mo.None[string](),
		WorkEnv:       mo.None[dotfilesmanager.DotfilesWorkEnvData](),
		SystemData:    mo.None[dotfilesmanager.DotfilesSystemData](),
	}

	err := manager.Initialize(data)

	require.NoError(t, err)
	require.DirExists(t, configDir)
	require.FileExists(t, configFilePath)
}

func Test_Initialize_DoesNotCreateDirectory_WhenDirectoryExists(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "config")
	configFilePath := filepath.Join(configDir, "chezmoi.toml")
	cloneDir := filepath.Join(tempDir, "clone")

	err := os.MkdirAll(configDir, 0755)
	require.NoError(t, err)

	fileSystem := utils.NewDefaultFileSystem()
	mockUserManager := &osmanager.MoqUserManager{}
	mockPackageManager := &pkgmanager.MoqPackageManager{}
	mockHTTPClient := &httpclient.MoqHTTPClient{}
	mockCommander := &utils.MoqCommander{}

	config := chezmoi.DefaultChezmoiConfig(configFilePath, cloneDir)
	manager := chezmoi.NewChezmoiManager(fileSystem, mockUserManager, mockCommander, mockPackageManager, mockHTTPClient, config)

	data := dotfilesmanager.DotfilesData{
		Email:         "test@example.com",
		FirstName:     "John",
		LastName:      "Doe",
		GpgSigningKey: mo.None[string](),
		WorkEnv:       mo.None[dotfilesmanager.DotfilesWorkEnvData](),
		SystemData:    mo.None[dotfilesmanager.DotfilesSystemData](),
	}

	err = manager.Initialize(data)

	require.NoError(t, err)
	require.FileExists(t, configFilePath)
}

func Test_Initialize_ReturnsError_WhenDirectoryCreationFails(t *testing.T) {
	expectedError := errors.New("permission denied")
	mockFileSystem := &utils.MoqFileSystem{
		PathExistsFunc: func(path string) (bool, error) {
			return false, os.ErrNotExist
		},
		CreateDirectoryFunc: func(path string) error {
			return expectedError
		},
	}
	mockUserManager := &osmanager.MoqUserManager{}
	mockPackageManager := &pkgmanager.MoqPackageManager{}
	mockHTTPClient := &httpclient.MoqHTTPClient{}
	mockCommander := &utils.MoqCommander{}

	configFilePath := "/home/user/.config/chezmoi.toml"
	cloneDir := "/home/user/.local/share/chezmoi"
	config := chezmoi.DefaultChezmoiConfig(configFilePath, cloneDir)
	manager := chezmoi.NewChezmoiManager(mockFileSystem, mockUserManager, mockCommander, mockPackageManager, mockHTTPClient, config)

	data := dotfilesmanager.DotfilesData{
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}

	err := manager.Initialize(data)

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to create chezmoi config directory")
	require.Contains(t, err.Error(), expectedError.Error())
}

func Test_Initialize_WritesBasicPersonalData_WhenOnlyRequiredFieldsProvided(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "config")
	configFilePath := filepath.Join(configDir, "chezmoi.toml")
	cloneDir := filepath.Join(tempDir, "clone")

	err := os.MkdirAll(configDir, 0755)
	require.NoError(t, err)

	fileSystem := utils.NewDefaultFileSystem()
	mockUserManager := &osmanager.MoqUserManager{}
	mockPackageManager := &pkgmanager.MoqPackageManager{}
	mockHTTPClient := &httpclient.MoqHTTPClient{}
	mockCommander := &utils.MoqCommander{}

	config := chezmoi.DefaultChezmoiConfig(configFilePath, cloneDir)
	manager := chezmoi.NewChezmoiManager(fileSystem, mockUserManager, mockCommander, mockPackageManager, mockHTTPClient, config)

	data := dotfilesmanager.DotfilesData{
		Email:         "test@example.com",
		FirstName:     "John",
		LastName:      "Doe",
		GpgSigningKey: mo.None[string](),
		WorkEnv:       mo.None[dotfilesmanager.DotfilesWorkEnvData](),
		SystemData:    mo.None[dotfilesmanager.DotfilesSystemData](),
	}

	err = manager.Initialize(data)

	require.NoError(t, err)
	require.FileExists(t, configFilePath)

	configContent, err := os.ReadFile(configFilePath)
	require.NoError(t, err)
	configStr := string(configContent)
	require.Contains(t, configStr, "test@example.com")
	require.Contains(t, configStr, "John Doe")
	require.Contains(t, configStr, "work_env = false")
}

func Test_Initialize_WritesGpgSigningKey_WhenProvided(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "config")
	configFilePath := filepath.Join(configDir, "chezmoi.toml")
	cloneDir := filepath.Join(tempDir, "clone")

	err := os.MkdirAll(configDir, 0755)
	require.NoError(t, err)

	fileSystem := utils.NewDefaultFileSystem()
	mockUserManager := &osmanager.MoqUserManager{}
	mockPackageManager := &pkgmanager.MoqPackageManager{}
	mockHTTPClient := &httpclient.MoqHTTPClient{}
	mockCommander := &utils.MoqCommander{}

	config := chezmoi.DefaultChezmoiConfig(configFilePath, cloneDir)
	manager := chezmoi.NewChezmoiManager(fileSystem, mockUserManager, mockCommander, mockPackageManager, mockHTTPClient, config)

	data := dotfilesmanager.DotfilesData{
		Email:         "test@example.com",
		FirstName:     "John",
		LastName:      "Doe",
		GpgSigningKey: mo.Some("ABC123DEF456"),
		WorkEnv:       mo.None[dotfilesmanager.DotfilesWorkEnvData](),
		SystemData:    mo.None[dotfilesmanager.DotfilesSystemData](),
	}

	err = manager.Initialize(data)

	require.NoError(t, err)
	require.FileExists(t, configFilePath)

	configContent, err := os.ReadFile(configFilePath)
	require.NoError(t, err)
	configStr := string(configContent)
	require.Contains(t, configStr, "ABC123DEF456")
}

func Test_Initialize_WritesWorkEnvironmentData_WhenProvided(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "config")
	configFilePath := filepath.Join(configDir, "chezmoi.toml")
	cloneDir := filepath.Join(tempDir, "clone")

	err := os.MkdirAll(configDir, 0755)
	require.NoError(t, err)

	fileSystem := utils.NewDefaultFileSystem()
	mockUserManager := &osmanager.MoqUserManager{}
	mockPackageManager := &pkgmanager.MoqPackageManager{}
	mockHTTPClient := &httpclient.MoqHTTPClient{}
	mockCommander := &utils.MoqCommander{}

	config := chezmoi.DefaultChezmoiConfig(configFilePath, cloneDir)
	manager := chezmoi.NewChezmoiManager(fileSystem, mockUserManager, mockCommander, mockPackageManager, mockHTTPClient, config)

	workEnvData := dotfilesmanager.DotfilesWorkEnvData{
		WorkName:  "Acme Corp",
		WorkEmail: "john.doe@acme.com",
	}
	data := dotfilesmanager.DotfilesData{
		Email:         "test@example.com",
		FirstName:     "John",
		LastName:      "Doe",
		GpgSigningKey: mo.None[string](),
		WorkEnv:       mo.Some(workEnvData),
		SystemData:    mo.None[dotfilesmanager.DotfilesSystemData](),
	}

	err = manager.Initialize(data)

	require.NoError(t, err)
	require.FileExists(t, configFilePath)

	configContent, err := os.ReadFile(configFilePath)
	require.NoError(t, err)
	configStr := string(configContent)
	require.Contains(t, configStr, "work_env = true")
	require.Contains(t, configStr, "Acme Corp")
	require.Contains(t, configStr, "john.doe@acme.com")
}

func Test_Initialize_WritesSystemData_WhenProvided(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "config")
	configFilePath := filepath.Join(configDir, "chezmoi.toml")
	cloneDir := filepath.Join(tempDir, "clone")

	err := os.MkdirAll(configDir, 0755)
	require.NoError(t, err)

	fileSystem := utils.NewDefaultFileSystem()
	mockUserManager := &osmanager.MoqUserManager{}
	mockPackageManager := &pkgmanager.MoqPackageManager{}
	mockHTTPClient := &httpclient.MoqHTTPClient{}
	mockCommander := &utils.MoqCommander{}

	config := chezmoi.DefaultChezmoiConfig(configFilePath, cloneDir)
	manager := chezmoi.NewChezmoiManager(fileSystem, mockUserManager, mockCommander, mockPackageManager, mockHTTPClient, config)

	systemData := dotfilesmanager.DotfilesSystemData{
		Shell:           "/bin/zsh",
		MultiUserSystem: true,
		BrewMultiUser:   "multifoo",
	}
	data := dotfilesmanager.DotfilesData{
		Email:         "test@example.com",
		FirstName:     "John",
		LastName:      "Doe",
		GpgSigningKey: mo.None[string](),
		WorkEnv:       mo.None[dotfilesmanager.DotfilesWorkEnvData](),
		SystemData:    mo.Some(systemData),
	}

	err = manager.Initialize(data)

	require.NoError(t, err)
	require.FileExists(t, configFilePath)

	configContent, err := os.ReadFile(configFilePath)
	require.NoError(t, err)
	configStr := string(configContent)
	require.Contains(t, configStr, "/bin/zsh")
	require.Contains(t, configStr, "multi_user_system = true")
	require.Contains(t, configStr, "brew_multi_user")
	require.Contains(t, configStr, "multifoo")
}

func Test_Initialize_WritesCompleteData_WhenAllFieldsProvided(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "config")
	configFilePath := filepath.Join(configDir, "chezmoi.toml")
	cloneDir := filepath.Join(tempDir, "clone")

	err := os.MkdirAll(configDir, 0755)
	require.NoError(t, err)

	fileSystem := utils.NewDefaultFileSystem()
	mockUserManager := &osmanager.MoqUserManager{}
	mockPackageManager := &pkgmanager.MoqPackageManager{}
	mockHTTPClient := &httpclient.MoqHTTPClient{}
	mockCommander := &utils.MoqCommander{}

	config := chezmoi.DefaultChezmoiConfig(configFilePath, cloneDir)
	manager := chezmoi.NewChezmoiManager(fileSystem, mockUserManager, mockCommander, mockPackageManager, mockHTTPClient, config)

	workEnvData := dotfilesmanager.DotfilesWorkEnvData{
		WorkName:  "Acme Corp",
		WorkEmail: "john.doe@acme.com",
	}
	systemData := dotfilesmanager.DotfilesSystemData{
		Shell:           "/bin/zsh",
		MultiUserSystem: true,
		BrewMultiUser:   "multifoo",
	}
	data := dotfilesmanager.DotfilesData{
		Email:         "test@example.com",
		FirstName:     "John",
		LastName:      "Doe",
		GpgSigningKey: mo.Some("ABC123DEF456"),
		WorkEnv:       mo.Some(workEnvData),
		SystemData:    mo.Some(systemData),
	}

	err = manager.Initialize(data)

	require.NoError(t, err)
	require.FileExists(t, configFilePath)

	configContent, err := os.ReadFile(configFilePath)
	require.NoError(t, err)
	configStr := string(configContent)
	require.Contains(t, configStr, "test@example.com")
	require.Contains(t, configStr, "John Doe")
	require.Contains(t, configStr, "ABC123DEF456")
	require.Contains(t, configStr, "work_env = true")
	require.Contains(t, configStr, "Acme Corp")
	require.Contains(t, configStr, "john.doe@acme.com")
	require.Contains(t, configStr, "/bin/zsh")
	require.Contains(t, configStr, "multi_user_system = true")
	require.Contains(t, configStr, "brew_multi_user")
	require.Contains(t, configStr, "multifoo")
}
