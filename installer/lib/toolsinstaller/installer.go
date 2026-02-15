package toolsinstaller

import (
	"fmt"

	"github.com/MrPointer/dotfiles/installer/lib/packageresolver"
	"github.com/MrPointer/dotfiles/installer/lib/pkgmanager"
	"github.com/MrPointer/dotfiles/installer/utils/logger"
)

// ToolsInstaller manages the installation of optional tools.
type ToolsInstaller struct {
	resolver       packageresolver.PackageManagerResolver
	packageManager pkgmanager.PackageManager
	logger         logger.Logger
}

// InstallResult represents the outcome of installing a single tool.
type InstallResult struct {
	Name    string
	Success bool
	Error   error
}

var _ interface{} = (*ToolsInstaller)(nil)

// NewToolsInstaller creates a new ToolsInstaller with the provided dependencies.
func NewToolsInstaller(
	resolver packageresolver.PackageManagerResolver,
	pm pkgmanager.PackageManager,
	log logger.Logger,
) *ToolsInstaller {
	return &ToolsInstaller{
		resolver:       resolver,
		packageManager: pm,
		logger:         log,
	}
}

// InstallTools installs the provided list of tools and returns results for each.
// If any tool installation fails, it logs a warning and continues with the next tool.
// It never aborts on individual failures.
func (ti *ToolsInstaller) InstallTools(tools []string) []InstallResult {
	results := make([]InstallResult, len(tools))

	for i, tool := range tools {
		result := InstallResult{Name: tool}

		ti.logger.UpdateProgress(fmt.Sprintf("Installing %s (%d/%d)", tool, i+1, len(tools)))

		// Resolve the tool name using the resolver
		resolvedInfo, err := ti.resolver.Resolve(tool, "")
		if err != nil {
			result.Success = false
			result.Error = fmt.Errorf("failed to resolve tool '%s': %w", tool, err)
			results[i] = result
			continue
		}

		// Install the tool using the package manager
		err = ti.packageManager.InstallPackage(resolvedInfo)
		if err != nil {
			result.Success = false
			result.Error = fmt.Errorf("failed to install tool '%s': %w", tool, err)
			results[i] = result
			continue
		}

		ti.logger.LogAccomplishment(fmt.Sprintf("Installed %s", tool))
		result.Success = true
		results[i] = result
	}

	return results
}
