package brew

import (
	"testing"

	"github.com/MrPointer/dotfiles/installer/lib/compatibility"
	"github.com/MrPointer/dotfiles/installer/utils/logger"
)

func TestDetectBrewPathReturnsExpectedPath(t *testing.T) {
	tests := []struct {
		name          string
		sysInfo       *compatibility.SystemInfo
		expectedPath  string
		errorExpected bool
	}{
		{"darwin/arm64", &compatibility.SystemInfo{OSName: "darwin", Arch: "arm64"}, MacOSARMBrewPath, false},
		{"darwin/amd64", &compatibility.SystemInfo{OSName: "darwin", Arch: "amd64"}, MacOSIntelBrewPath, false},
		{"linux/amd64", &compatibility.SystemInfo{OSName: "linux", Arch: "amd64"}, LinuxBrewPath, false},
		{"unsupported", &compatibility.SystemInfo{OSName: "plan9", Arch: "amd64"}, "", true},
		{"no sysinfo", nil, "", true},
	}
	for _, tt := range tests {
		current := tt
		t.Run(current.name, func(t *testing.T) {
			b := &brewInstaller{
				logger:     logger.DefaultLogger,
				systemInfo: current.sysInfo,
				commander:  nil,
			}
			got, err := b.DetectBrewPath()
			if (err != nil) != current.errorExpected {
				t.Fatalf("DetectBrewPath() error = %v, wantError %v", err, current.errorExpected)
			}
			if got != current.expectedPath {
				t.Errorf("DetectBrewPath() = %q, want %q", got, current.expectedPath)
			}
		})
	}
}

func TestNewBrewInstallerCreatesMultiUserImplementationWhenOptionIsEnabled(t *testing.T) {
	opts := Options{
		MultiUserSystem: true,
		Logger:          logger.DefaultLogger,
		SystemInfo:      &compatibility.SystemInfo{OSName: "linux", Arch: "amd64"},
		Commander:       nil,
	}
	installer := NewBrewInstaller(opts)
	if installer == nil {
		t.Fatal("expected non-nil installer")
	}
	if _, ok := installer.(*MultiUserBrewInstaller); !ok {
		t.Errorf("expected *MultiUserBrewInstaller, got %T", installer)
	}
}
