package chezmoi

func (c *ChezmoiManager) Apply() error {
	// Always remove existing chezmoi clone first, just in case
	err := c.filesystem.RemovePath(c.chezmoiConfigDir)
	if err != nil {
		return err
	}

	// Implementation
	return nil
}
