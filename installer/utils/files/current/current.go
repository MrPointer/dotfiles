package current

import (
	"errors"
	"path/filepath"
	"runtime"
)

// Filename returns the name of the current file.
// It uses [runtime.Caller] to get the file name of the caller.
func Filename() (string, error) {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return "", errors.New("unable to get the current filename")
	}
	return filename, nil
}

// Dirname returns the directory name of the current file.
// It uses the [Filename] function
// to get the file name and then extracts the directory part.
func Dirname() (string, error) {
	filename, err := Filename()
	if err != nil {
		return "", err
	}
	return filepath.Dir(filename), nil
}
