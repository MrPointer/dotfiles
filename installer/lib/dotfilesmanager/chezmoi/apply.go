package chezmoi

import (
	"fmt"

	"github.com/MrPointer/dotfiles/installer/utils"
)

func (c *ChezmoiManager) Apply() error {
	c.logger.Debug("Applying chezmoi")

	// Always remove existing chezmoi clone first, just in case
	err := c.filesystem.RemovePath(c.chezmoiConfig.chezmoiCloneDir)
	if err != nil {
		return err
	}

	c.logger.Trace("Building chezmoi apply command")
	chezmoiApplyCmdArgs := []string{"init", "--apply"}
	if c.chezmoiConfig.chezmoiCloneDir != "" {
		chezmoiApplyCmdArgs = append(chezmoiApplyCmdArgs, "--source", c.chezmoiConfig.chezmoiCloneDir)
	}
	if c.chezmoiConfig.cloneViaSSH {
		chezmoiApplyCmdArgs = append(chezmoiApplyCmdArgs, "--ssh")
	}
	chezmoiApplyCmdArgs = append(chezmoiApplyCmdArgs, c.chezmoiConfig.githubUsername)

	var discardOutputOption utils.Option = utils.EmptyOption()
	if c.displayMode != utils.DisplayModePassthrough {
		discardOutputOption = utils.WithDiscardOutput()
	}

	result, err := c.commander.RunCommand("chezmoi", chezmoiApplyCmdArgs, discardOutputOption)
	if err != nil {
		return err
	}
	if result.ExitCode != 0 {
		return fmt.Errorf("chezmoi init failed with exit code %d: %s", result.ExitCode, result.StderrString())
	}

	c.logger.Debug("Chezmoi has been applied successfully")
	return nil
}
