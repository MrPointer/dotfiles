# Release Process

This document outlines the process for creating releases of the dotfiles installer.

## Prerequisites

- Push access to the main branch
- Ability to create tags and releases
- All changes merged to `main` branch
- Tests passing on the latest commit

## Release Steps

### 1. Prepare the Release

1. Ensure all changes for the release are merged into `main`
2. Update version references if needed (the build process handles most versioning automatically)
3. Run tests locally to ensure everything works:
   ```bash
   cd installer
   go test -race -v ./...
   ```

### 2. Create and Push the Tag

1. Create a new tag following semantic versioning (e.g., `v1.0.0`, `v1.1.0`, `v2.0.0-beta.1`):
   ```bash
   git tag v1.0.0
   ```

2. Push the tag to trigger the release workflow:
   ```bash
   git push origin v1.0.0
   ```

### 3. Monitor the Release Process

1. Go to the [Actions tab](https://github.com/MrPointer/dotfiles/actions) in GitHub
2. Monitor the "Release" workflow that gets triggered by the tag push
3. The workflow will:
   - Build binaries for all supported platforms
   - Sign the binaries with cosign
   - Create a GitHub release with the binaries attached
   - Generate checksums for verification

### 4. Verify the Release

1. Check the [Releases page](https://github.com/MrPointer/dotfiles/releases)
2. Verify that the release contains:
   - Binaries for all platforms (macOS arm64/x86_64, Linux arm64/x86_64)
   - Checksums file (`checksums.txt`)
   - Proper release notes
3. Test the get script with the new release:
   ```bash
   curl -fsSL https://raw.githubusercontent.com/MrPointer/dotfiles/main/get.sh | bash -s -- --version v1.0.0
   ```

## Supported Platforms

The automated release process builds binaries for:

- **macOS**
  - ARM64 (Apple Silicon)
  - x86_64 (Intel)
- **Linux**
  - ARM64
  - x86_64

## Versioning Strategy

We follow [Semantic Versioning (SemVer)](https://semver.org/):

- **MAJOR version** (`v2.0.0`): Incompatible API changes or major breaking changes
- **MINOR version** (`v1.1.0`): Backward-compatible functionality additions
- **PATCH version** (`v1.0.1`): Backward-compatible bug fixes

### Pre-release Versions

For beta or release candidate versions, append the pre-release identifier:
- `v1.0.0-beta.1`
- `v1.0.0-rc.1`

## Release Checklist

- [ ] All changes merged to `main`
- [ ] Tests passing locally and on CI
- [ ] Version tag created and follows SemVer
- [ ] Tag pushed to GitHub
- [ ] Release workflow completed successfully
- [ ] Release appears on GitHub with all expected assets
- [ ] Get script tested with new version
- [ ] Release announcement (if needed)

## Troubleshooting

### Release Workflow Failed

1. Check the workflow logs in the Actions tab
2. Common issues:
   - Go build failures (check for compilation errors)
   - Test failures (fix tests before releasing)
   - GoReleaser configuration issues (validate `.goreleaser.yaml`)
   - Permission issues (ensure GITHUB_TOKEN has proper permissions)

### Missing Binaries

If some platform binaries are missing:
1. Check the GoReleaser configuration in `installer/.goreleaser.yaml`
2. Verify the `builds` section includes all required platforms
3. Re-run the release by deleting and recreating the tag

### Binary Signing Issues

If cosign signing fails:
1. Check that the workflow has `id-token: write` permissions
2. Verify cosign installation in the workflow
3. The signing process uses GitHub's OIDC provider automatically

## Manual Release (Fallback)

If the automated process fails, you can create a release manually:

```bash
cd installer
goreleaser release --clean
```

This requires:
- GoReleaser installed locally
- `GITHUB_TOKEN` environment variable set
- cosign installed for signing (optional)

## Post-Release Tasks

1. Update any documentation that references version numbers
2. Consider updating the get script if there are breaking changes
3. Announce the release in relevant channels if it's a major release
4. Monitor for any issues reported by users

## Get Script Features

The `get.sh` script provides multiple ways for users to download the installer binary and optionally run it:

**Basic download:**
```bash
curl -fsSL https://raw.githubusercontent.com/MrPointer/dotfiles/main/get.sh | bash
```

**One-command download and install dotfiles:**
```bash
curl -fsSL https://raw.githubusercontent.com/MrPointer/dotfiles/main/get.sh | bash -s -- --run
```

**Download and install dotfiles with custom options:**
```bash
curl -fsSL https://raw.githubusercontent.com/MrPointer/dotfiles/main/get.sh | bash -s -- --run -- install --work-env --non-interactive
```

This solves the "pass custom options" problem by allowing users to specify installer arguments after `--run --`. Note that "install" here refers to installing the dotfiles, not installing the binary.

## Security Considerations

- All binaries are built in a secure GitHub Actions environment with reproducible builds
- SHA256 checksums are generated for all assets to enable verification
- The release process runs with minimal required permissions
- No secrets are required for the standard release process
- Builds are deterministic and traceable through GitHub Actions logs
