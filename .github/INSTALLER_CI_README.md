# Installer CI Documentation

This directory contains GitHub Actions workflows for the dotfiles installer project located in the `installer/` subdirectory.

## Workflow Overview

### 🔧 installer-ci.yml
**Purpose**: Build and test the installer  
**Trigger**: Push to main or PRs affecting `installer/` directory  

**Jobs (in order)**:
1. **Build**: Uses GoReleaser to create cross-platform binaries
2. **Test**: Runs Go test suite with race detection  
3. **E2E Tests**: Tests installer on multiple platforms in real environments

**Platforms Tested**:
- Ubuntu (latest)
- Debian (bookworm container)
- macOS (latest)

## Build Process

The CI uses GoReleaser in **snapshot mode** to:
- Build cross-platform binaries (Linux/macOS, AMD64/ARM64)
- Generate consistent build artifacts
- **No releases** - just builds for testing

## E2E Testing

The pipeline includes end-to-end testing that:

1. **Downloads** the built installer binary
2. **Tests** on multiple OS distributions  
3. **Verifies** compatibility detection works
4. **Runs** installer in non-interactive mode
5. **Validates** graceful behavior in CI environments

### Test Configuration

E2E tests use these flags for CI compatibility:
- `--non-interactive`: Skips all user prompts
- `--plain`: Disables progress indicators for cleaner logs  
- `--install-brew=false`: Skips Homebrew for faster testing
- `--install-prerequisites=false`: Skips prerequisite installation
- `--git-clone-protocol=https`: Uses HTTPS instead of SSH

### Test Environment

- **Isolated**: Uses `/tmp/test-home` as HOME directory
- **Timeout**: 300 seconds (5 minutes) to prevent hangs
- **Graceful Failures**: Expected in CI since we don't have full system setup

## Local Testing

To test the workflow locally:

```bash
cd installer/
task build
./bin/dotfiles-installer install --non-interactive --plain --install-brew=false
```

## Workflow Structure

```
.github/
└── workflows/
    └── installer-ci.yml    # Main CI pipeline
```

## Troubleshooting

### Common Issues

1. **Build Fails**: Check Go version in `installer/go.mod`
2. **E2E Test Fails**: Review platform-specific requirements
3. **Test Timeout**: E2E tests timeout after 5 minutes

### Debugging E2E Tests

E2E tests are designed to handle CI environment limitations:
- Allow expected failures (exit code 1) in restricted environments
- Create isolated test directories
- Skip complex system modifications

To debug locally:
```bash
export HOME="/tmp/test-home"
mkdir -p "$HOME"
./installer/bin/dotfiles-installer install --non-interactive --plain
```

## Future Enhancements

- **Linting**: Will integrate golangci-lint later
- **Security Scanning**: Can add Trivy scans if needed
- **Release Automation**: Will add when ready for releases