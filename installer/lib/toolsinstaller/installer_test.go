package toolsinstaller

import (
	"errors"
	"testing"

	"github.com/MrPointer/dotfiles/installer/lib/packageresolver"
	"github.com/MrPointer/dotfiles/installer/lib/pkgmanager"
	"github.com/MrPointer/dotfiles/installer/utils/logger"
	"github.com/stretchr/testify/require"
)

func Test_InstallToolsSucceedsForAllTools(t *testing.T) {
	resolver := &packageresolver.MoqPackageManagerResolver{
		ResolveFunc: func(genericPackageCode, versionConstraintString string) (pkgmanager.RequestedPackageInfo, error) {
			return pkgmanager.RequestedPackageInfo{Name: genericPackageCode + "-resolved"}, nil
		},
	}

	pm := &pkgmanager.MoqPackageManager{
		InstallPackageFunc: func(requestedPackageInfo pkgmanager.RequestedPackageInfo) error {
			return nil
		},
	}

	installer := NewToolsInstaller(resolver, pm, &logger.NoopLogger{})

	results := installer.InstallTools([]string{"git", "fzf"})

	require.Len(t, results, 2)
	require.True(t, results[0].Success)
	require.Equal(t, "git", results[0].Name)
	require.NoError(t, results[0].Error)
	require.True(t, results[1].Success)
	require.Equal(t, "fzf", results[1].Name)
	require.NoError(t, results[1].Error)
}

func Test_InstallToolsContinuesWhenResolverFails(t *testing.T) {
	resolver := &packageresolver.MoqPackageManagerResolver{
		ResolveFunc: func(genericPackageCode, versionConstraintString string) (pkgmanager.RequestedPackageInfo, error) {
			if genericPackageCode == "bad-tool" {
				return pkgmanager.RequestedPackageInfo{}, errors.New("unknown tool")
			}

			return pkgmanager.RequestedPackageInfo{Name: genericPackageCode + "-resolved"}, nil
		},
	}

	pm := &pkgmanager.MoqPackageManager{
		InstallPackageFunc: func(requestedPackageInfo pkgmanager.RequestedPackageInfo) error {
			return nil
		},
	}

	installer := NewToolsInstaller(resolver, pm, &logger.NoopLogger{})

	results := installer.InstallTools([]string{"git", "bad-tool", "fzf"})

	require.Len(t, results, 3)
	require.True(t, results[0].Success)
	require.False(t, results[1].Success)
	require.Contains(t, results[1].Error.Error(), "unknown tool")
	require.True(t, results[2].Success)
}

func Test_InstallToolsContinuesWhenInstallFails(t *testing.T) {
	callCount := 0
	resolver := &packageresolver.MoqPackageManagerResolver{
		ResolveFunc: func(genericPackageCode, versionConstraintString string) (pkgmanager.RequestedPackageInfo, error) {
			return pkgmanager.RequestedPackageInfo{Name: genericPackageCode + "-resolved"}, nil
		},
	}

	pm := &pkgmanager.MoqPackageManager{
		InstallPackageFunc: func(requestedPackageInfo pkgmanager.RequestedPackageInfo) error {
			callCount++
			if callCount == 2 {
				return errors.New("installation failed")
			}

			return nil
		},
	}

	installer := NewToolsInstaller(resolver, pm, &logger.NoopLogger{})

	results := installer.InstallTools([]string{"git", "bad-pkg", "fzf"})

	require.Len(t, results, 3)
	require.True(t, results[0].Success)
	require.False(t, results[1].Success)
	require.Contains(t, results[1].Error.Error(), "installation failed")
	require.True(t, results[2].Success)
}

func Test_InstallToolsReturnsEmptyResultsForEmptyList(t *testing.T) {
	resolver := &packageresolver.MoqPackageManagerResolver{}
	pm := &pkgmanager.MoqPackageManager{}
	installer := NewToolsInstaller(resolver, pm, &logger.NoopLogger{})

	results := installer.InstallTools([]string{})

	require.Len(t, results, 0)
}

func Test_InstallToolsTracksSuccessAndFailureAccurately(t *testing.T) {
	toolResults := map[string]bool{
		"git":     true,
		"fzf":     false,
		"bat":     true,
		"neovim":  false,
		"ripgrep": true,
	}

	resolver := &packageresolver.MoqPackageManagerResolver{
		ResolveFunc: func(genericPackageCode, versionConstraintString string) (pkgmanager.RequestedPackageInfo, error) {
			if !toolResults[genericPackageCode] {
				return pkgmanager.RequestedPackageInfo{}, errors.New("failed to resolve")
			}

			return pkgmanager.RequestedPackageInfo{Name: genericPackageCode + "-resolved"}, nil
		},
	}

	pm := &pkgmanager.MoqPackageManager{
		InstallPackageFunc: func(requestedPackageInfo pkgmanager.RequestedPackageInfo) error {
			return nil
		},
	}

	installer := NewToolsInstaller(resolver, pm, &logger.NoopLogger{})

	tools := []string{"git", "fzf", "bat", "neovim", "ripgrep"}
	results := installer.InstallTools(tools)

	require.Len(t, results, 5)
	require.True(t, results[0].Success)
	require.False(t, results[1].Success)
	require.True(t, results[2].Success)
	require.False(t, results[3].Success)
	require.True(t, results[4].Success)
}
