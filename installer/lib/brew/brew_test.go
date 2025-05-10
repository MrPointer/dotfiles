package brew

import (
	"os"
	"runtime"
	"testing"

	"github.com/MrPointer/dotfiles/installer/lib/compatibility"
	"github.com/MrPointer/dotfiles/installer/utils/logger"
)

func TestDetectBrewPathReturnsCorrectPath(t *testing.T) {
	tests := []struct {
		name      string
		sysInfo   *compatibility.SystemInfo
		wantPath  string
		wantError bool
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
			if (err != nil) != current.wantError {
				t.Fatalf("DetectBrewPath() error = %v, wantError %v", err, current.wantError)
			}
			if got != current.wantPath {
				t.Errorf("DetectBrewPath() = %q, want %q", got, current.wantPath)
			}
		})
	}
}

func TestIsAvailableReturnsFalseWhenBrewNotInstalled(t *testing.T) {
	b := &brewInstaller{
		logger:           logger.DefaultLogger,
		systemInfo:       &compatibility.SystemInfo{OSName: runtime.GOOS, Arch: runtime.GOARCH},
		commander:        nil,
		brewPathOverride: "/tmp/nonexistent-brew-binary",
	}
	avail, err := b.IsAvailable()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if avail {
		t.Error("expected not available, got available")
	}
}

func TestIsAvailableReturnsTrueWhenBrewInstalled(t *testing.T) {
	tempFile, err := os.CreateTemp("", "brew-binary-*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	b := &brewInstaller{
		logger:           logger.DefaultLogger,
		systemInfo:       &compatibility.SystemInfo{OSName: runtime.GOOS, Arch: runtime.GOARCH},
		commander:        nil,
		brewPathOverride: tempFile.Name(),
	}
	avail, err := b.IsAvailable()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !avail {
		t.Error("expected available, got not available")
	}
}

func TestNewBrewInstallerReturnsSingleUserImpl(t *testing.T) {
	opts := Options{
		MultiUserSystem: false,
		Logger:          logger.DefaultLogger,
		SystemInfo:      &compatibility.SystemInfo{OSName: "linux", Arch: "amd64"},
		Commander:       nil,
	}
	installer := NewBrewInstaller(opts)
	if installer == nil {
		t.Fatal("expected non-nil installer")
	}
	if _, ok := installer.(*brewInstaller); !ok {
		t.Errorf("expected *brewInstaller, got %T", installer)
	}
}

func TestNewBrewInstallerReturnsMultiUserImpl(t *testing.T) {
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
