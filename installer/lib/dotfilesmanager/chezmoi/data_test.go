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
	"github.com/MrPointer/dotfiles/installer/utils"
	"github.com/MrPointer/dotfiles/installer/utils/osmanager"
)

func Test_NewChezmoiDataInitializer_ReturnsValidInstance(t *testing.T) {
	mockFileSystem := &utils.MoqFileSystem{}
	configFilePath := "/home/user/.config/chezmoi.toml"

	initializer := chezmoi.NewChezmoiDataInitializer(configFilePath, mockFileSystem)

	require.NotNil(t, initializer)
}

func Test_TryNewDefaultChezmoiDataInitializer_ReturnsValidInstance_WhenUserConfigDirIsAvailable(t *testing.T) {
	mockFileSystem := &utils.MoqFileSystem{}

	userManager := &osmanager.MoqUserManager{}
	userManager.GetConfigDirFunc = func() (string, error) {
		return "/home/user/.config", nil
	}

	initializer, err := chezmoi.TryNewDefaultChezmoiDataInitializer(mockFileSystem, userManager)

	require.NoError(t, err)
	require.NotNil(t, initializer)
}

func Test_TryNewDefaultChezmoiDataInitializer_ReturnsError_WhenUserConfigDirIsUnavailable(t *testing.T) {
	mockFileSystem := &utils.MoqFileSystem{}

	userManager := &osmanager.MoqUserManager{}
	userManager.GetConfigDirFunc = func() (string, error) {
		return "", errors.New("failed to get user config directory")
	}

	initializer, err := chezmoi.TryNewDefaultChezmoiDataInitializer(mockFileSystem, userManager)

	require.Error(t, err)
	require.Nil(t, initializer)
}

func Test_Initialize_CreatesConfigDirectory_WhenDirectoryDoesNotExist(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "config")
	configFilePath := filepath.Join(configDir, "chezmoi.toml")

	mockFileSystem := &utils.MoqFileSystem{
		PathExistsFunc: func(path string) (bool, error) {
			if path == configDir {
				return false, os.ErrNotExist
			}
			return false, nil
		},
		CreateDirectoryFunc: func(path string) error {
			return os.MkdirAll(path, 0755)
		},
	}

	initializer := chezmoi.NewChezmoiDataInitializer(configFilePath, mockFileSystem)

	data := dotfilesmanager.DotfilesData{
		Email:         "test@example.com",
		FirstName:     "John",
		LastName:      "Doe",
		GpgSigningKey: mo.None[string](),
		WorkEnv:       mo.None[dotfilesmanager.DotfilesWorkEnvData](),
		SystemData:    mo.None[dotfilesmanager.DotfilesSystemData](),
	}

	err := initializer.Initialize(data)

	require.NoError(t, err)
	require.Len(t, mockFileSystem.PathExistsCalls(), 1)
	require.Equal(t, configDir, mockFileSystem.PathExistsCalls()[0].Path)
	require.Len(t, mockFileSystem.CreateDirectoryCalls(), 1)
	require.Equal(t, configDir, mockFileSystem.CreateDirectoryCalls()[0].Path)

	require.FileExists(t, configFilePath)
}

func Test_Initialize_DoesNotCreateDirectory_WhenDirectoryExists(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "config")
	configFilePath := filepath.Join(configDir, "chezmoi.toml")

	err := os.MkdirAll(configDir, 0755)
	require.NoError(t, err)

	mockFileSystem := &utils.MoqFileSystem{
		PathExistsFunc: func(path string) (bool, error) {
			if path == configDir {
				return true, nil
			}
			return false, nil
		},
	}

	initializer := chezmoi.NewChezmoiDataInitializer(configFilePath, mockFileSystem)

	data := dotfilesmanager.DotfilesData{
		Email:         "test@example.com",
		FirstName:     "John",
		LastName:      "Doe",
		GpgSigningKey: mo.None[string](),
		WorkEnv:       mo.None[dotfilesmanager.DotfilesWorkEnvData](),
		SystemData:    mo.None[dotfilesmanager.DotfilesSystemData](),
	}

	err = initializer.Initialize(data)

	require.NoError(t, err)
	require.Len(t, mockFileSystem.PathExistsCalls(), 1)
	require.Equal(t, configDir, mockFileSystem.PathExistsCalls()[0].Path)
	require.Len(t, mockFileSystem.CreateDirectoryCalls(), 0)

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
	configFilePath := "/home/user/.config/chezmoi.toml"

	initializer := chezmoi.NewChezmoiDataInitializer(configFilePath, mockFileSystem)

	data := dotfilesmanager.DotfilesData{
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}

	err := initializer.Initialize(data)

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

	err := os.MkdirAll(configDir, 0755)
	require.NoError(t, err)

	mockFileSystem := &utils.MoqFileSystem{
		PathExistsFunc: func(path string) (bool, error) {
			return true, nil
		},
	}

	initializer := chezmoi.NewChezmoiDataInitializer(configFilePath, mockFileSystem)

	data := dotfilesmanager.DotfilesData{
		Email:         "test@example.com",
		FirstName:     "John",
		LastName:      "Doe",
		GpgSigningKey: mo.None[string](),
		WorkEnv:       mo.None[dotfilesmanager.DotfilesWorkEnvData](),
		SystemData:    mo.None[dotfilesmanager.DotfilesSystemData](),
	}

	err = initializer.Initialize(data)

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

	err := os.MkdirAll(configDir, 0755)
	require.NoError(t, err)

	mockFileSystem := &utils.MoqFileSystem{
		PathExistsFunc: func(path string) (bool, error) {
			return true, nil
		},
	}

	initializer := chezmoi.NewChezmoiDataInitializer(configFilePath, mockFileSystem)

	data := dotfilesmanager.DotfilesData{
		Email:         "test@example.com",
		FirstName:     "John",
		LastName:      "Doe",
		GpgSigningKey: mo.Some("ABC123DEF456"),
		WorkEnv:       mo.None[dotfilesmanager.DotfilesWorkEnvData](),
		SystemData:    mo.None[dotfilesmanager.DotfilesSystemData](),
	}

	err = initializer.Initialize(data)

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

	err := os.MkdirAll(configDir, 0755)
	require.NoError(t, err)

	mockFileSystem := &utils.MoqFileSystem{
		PathExistsFunc: func(path string) (bool, error) {
			return true, nil
		},
	}

	initializer := chezmoi.NewChezmoiDataInitializer(configFilePath, mockFileSystem)

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

	err = initializer.Initialize(data)

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

	err := os.MkdirAll(configDir, 0755)
	require.NoError(t, err)

	mockFileSystem := &utils.MoqFileSystem{
		PathExistsFunc: func(path string) (bool, error) {
			return true, nil
		},
	}

	initializer := chezmoi.NewChezmoiDataInitializer(configFilePath, mockFileSystem)

	systemData := dotfilesmanager.DotfilesSystemData{
		Shell:           "/bin/zsh",
		User:            "johndoe",
		MultiUserSystem: true,
		BrewUser:        false,
	}
	data := dotfilesmanager.DotfilesData{
		Email:         "test@example.com",
		FirstName:     "John",
		LastName:      "Doe",
		GpgSigningKey: mo.None[string](),
		WorkEnv:       mo.None[dotfilesmanager.DotfilesWorkEnvData](),
		SystemData:    mo.Some(systemData),
	}

	err = initializer.Initialize(data)

	require.NoError(t, err)
	require.FileExists(t, configFilePath)

	configContent, err := os.ReadFile(configFilePath)
	require.NoError(t, err)
	configStr := string(configContent)
	require.Contains(t, configStr, "/bin/zsh")
	require.Contains(t, configStr, "johndoe")
	require.Contains(t, configStr, "multi_user_system = true")
	require.Contains(t, configStr, "brew_user = false")
}

func Test_Initialize_WritesCompleteData_WhenAllFieldsProvided(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "config")
	configFilePath := filepath.Join(configDir, "chezmoi.toml")

	err := os.MkdirAll(configDir, 0755)
	require.NoError(t, err)

	mockFileSystem := &utils.MoqFileSystem{
		PathExistsFunc: func(path string) (bool, error) {
			return true, nil
		},
	}

	initializer := chezmoi.NewChezmoiDataInitializer(configFilePath, mockFileSystem)

	workEnvData := dotfilesmanager.DotfilesWorkEnvData{
		WorkName:  "Acme Corp",
		WorkEmail: "john.doe@acme.com",
	}
	systemData := dotfilesmanager.DotfilesSystemData{
		Shell:           "/bin/zsh",
		User:            "johndoe",
		MultiUserSystem: true,
		BrewUser:        false,
	}
	data := dotfilesmanager.DotfilesData{
		Email:         "test@example.com",
		FirstName:     "John",
		LastName:      "Doe",
		GpgSigningKey: mo.Some("ABC123DEF456"),
		WorkEnv:       mo.Some(workEnvData),
		SystemData:    mo.Some(systemData),
	}

	err = initializer.Initialize(data)

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
	require.Contains(t, configStr, "johndoe")
	require.Contains(t, configStr, "multi_user_system = true")
	require.Contains(t, configStr, "brew_user = false")
}
