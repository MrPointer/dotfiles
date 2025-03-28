package current

import "path/filepath"

// Gets the root directory of the project by navigating three levels up from the current directory.
// This is useful for locating files or directories that are relative to the root of the project.
//
// Note: Should be used mostly in tests.
func RootDirectory() (string, error) {
	currentDir, err := Dirname()
	if err != nil {
		return "", err
	}

	return filepath.Join(currentDir, "..", "..", ".."), nil
}
