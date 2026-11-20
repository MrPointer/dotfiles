package cmd

import (
	"errors"
	"testing"

	"github.com/MrPointer/dotfiles/installer/lib/dotfilesmanager"
	"github.com/MrPointer/dotfiles/installer/utils/logger"
	"github.com/stretchr/testify/require"
)

const (
	callInstall    = "Install"
	callInitialize = "Initialize"
	callApply      = "Apply"
)

func newRecordingManager(calls *[]string) *dotfilesmanager.MoqDotfilesManager {
	return &dotfilesmanager.MoqDotfilesManager{
		InstallFunc: func() error {
			*calls = append(*calls, callInstall)
			return nil
		},
		InitializeFunc: func(data dotfilesmanager.DotfilesData) error {
			*calls = append(*calls, callInitialize)
			return nil
		},
		ApplyFunc: func() error {
			*calls = append(*calls, callApply)
			return nil
		},
	}
}

func Test_PrepareOnlyStopsAfterInitializeAndSkipsApply(t *testing.T) {
	var calls []string

	dm := newRecordingManager(&calls)

	err := runDotfilesManagerSetup(&logger.NoopLogger{}, dm, true)

	require.NoError(t, err)
	require.Equal(t, []string{callInstall, callInitialize}, calls)
	require.False(t, shouldInstallOptionalTools(true, false))
}

func Test_PrepareOnlyWithInstallToolsRunsOptionalToolsAndSkipsApply(t *testing.T) {
	var calls []string

	dm := newRecordingManager(&calls)

	err := runDotfilesManagerSetup(&logger.NoopLogger{}, dm, true)

	require.NoError(t, err)
	require.Equal(t, []string{callInstall, callInitialize}, calls)
	require.True(t, shouldInstallOptionalTools(true, true))
}

func Test_PrepareOnlyAbsentRunsInstallInitializeAndApplyInOrder(t *testing.T) {
	var calls []string

	dm := newRecordingManager(&calls)

	err := runDotfilesManagerSetup(&logger.NoopLogger{}, dm, false)

	require.NoError(t, err)
	require.Equal(t, []string{callInstall, callInitialize, callApply}, calls)
	require.True(t, shouldInstallOptionalTools(false, false))
	require.True(t, shouldInstallOptionalTools(false, true))
}

func Test_PrepareOnlyFailingInitializeReturnsError(t *testing.T) {
	var calls []string

	dm := newRecordingManager(&calls)
	dm.InitializeFunc = func(data dotfilesmanager.DotfilesData) error {
		return errors.New("initialize failed")
	}

	err := runDotfilesManagerSetup(&logger.NoopLogger{}, dm, true)

	require.Error(t, err)
	require.Contains(t, err.Error(), "initialize failed")
	require.Equal(t, []string{callInstall}, calls)
}
