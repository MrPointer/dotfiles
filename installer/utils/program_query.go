package utils

import (
	"errors"
	"os"
	"os/exec"
)

type VersionExtractor func(string) (string, error)

// ProgramExistsQuery is the minimal interface for checking whether a program exists.
//
// Keep this narrow so it can be used widely (e.g., privilege escalation) without pulling in extra concerns.
type ProgramExistsQuery interface {
	ProgramExists(program string) (bool, error)
}

// ProgramQuery provides utilities for checking program availability and retrieving versions.
type ProgramQuery interface {
	ProgramExistsQuery

	GetProgramPath(program string) (string, error)
	GetProgramVersion(program string, versionExtractor VersionExtractor, queryArgs ...string) (string, error)
}

// GoNativeProgramQuery is a standard library based ProgramQuery implementation.
type GoNativeProgramQuery struct{}

var _ ProgramQuery = (*GoNativeProgramQuery)(nil)

func NewGoNativeProgramQuery() *GoNativeProgramQuery {
	return &GoNativeProgramQuery{}
}

func (q *GoNativeProgramQuery) GetProgramPath(program string) (string, error) {
	return exec.LookPath(program)
}

func (q *GoNativeProgramQuery) ProgramExists(program string) (bool, error) {
	_, err := q.GetProgramPath(program)
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) || errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (q *GoNativeProgramQuery) GetProgramVersion(
	program string,
	versionExtractor VersionExtractor,
	queryArgs ...string,
) (string, error) {
	args := []string{"--version"}
	if len(queryArgs) > 0 {
		args = queryArgs
	}

	cmd := exec.Command(program, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return versionExtractor(string(output))
}
