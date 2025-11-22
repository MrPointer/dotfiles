# Installation Guide

This guide covers different methods to download the dotfiles installer and install your dotfiles.

## Quick Start

The fastest way to get started is using our get script:

```bash
curl -fsSL https://raw.githubusercontent.com/MrPointer/dotfiles/main/get.sh | bash
```

This will download the latest installer binary to `~/.local/bin/dotfiles-installer`.

## Installation Methods

### 1. Using the Get Script (Recommended)

The get script automatically detects your platform and downloads the appropriate binary:

```bash
# Download to default location (~/.local/bin)
curl -fsSL https://raw.githubusercontent.com/MrPointer/dotfiles/main/get.sh | bash

# Download and install dotfiles immediately (one-command setup)
curl -fsSL https://raw.githubusercontent.com/MrPointer/dotfiles/main/get.sh | bash -s -- --run

# Download and install dotfiles with custom options
curl -fsSL https://raw.githubusercontent.com/MrPointer/dotfiles/main/get.sh | bash -s -- --run -- --work-env --install-prerequisites

# Download and install dotfiles non-interactively
curl -fsSL https://raw.githubusercontent.com/MrPointer/dotfiles/main/get.sh | bash -s -- --run -- install --non-interactive --plain

# Download to custom directory
curl -fsSL https://raw.githubusercontent.com/MrPointer/dotfiles/main/get.sh | bash -s -- --dir /usr/local/bin

# Download specific version
curl -fsSL https://raw.githubusercontent.com/MrPointer/dotfiles/main/get.sh | bash -s -- --version v1.0.0
```

#### Get Script Options

- `--dir DIR`: Download directory (default: `$HOME/.local/bin`)
- `--version VER`: Specific version to download (default: latest)
- `--run`: Run the installer after download (this installs your dotfiles)
- `--help`: Show help message

#### Passing Arguments to the Installer

When using `--run`, you can pass arguments to the installer by placing them after `--`:

```bash
# Basic usage - runs 'install' command to install your dotfiles
curl -fsSL https://raw.githubusercontent.com/MrPointer/dotfiles/main/get.sh | bash -s -- --run

# Pass specific command and options for dotfiles installation
curl -fsSL https://raw.githubusercontent.com/MrPointer/dotfiles/main/get.sh | bash -s -- --run -- install --work-env --non-interactive

# Run compatibility check instead
curl -fsSL https://raw.githubusercontent.com/MrPointer/dotfiles/main/get.sh | bash -s -- --run -- check-compatibility
```

### 2. Manual Download

Download pre-built binaries from [GitHub Releases](https://github.com/MrPointer/dotfiles/releases).

#### macOS
```bash
curl -L -o dotfiles-installer.tar.gz https://github.com/MrPointer/dotfiles/releases/latest/download/dotfiles-installer-*-darwin-arm64.tar.gz
tar -xzf dotfiles-installer.tar.gz
chmod +x dotfiles-installer
sudo mv dotfiles-installer /usr/local/bin/
```

#### Linux (x86_64)
```bash
curl -L -o dotfiles-installer.tar.gz https://github.com/MrPointer/dotfiles/releases/latest/download/dotfiles-installer-*-linux-x86_64.tar.gz
tar -xzf dotfiles-installer.tar.gz
chmod +x dotfiles-installer
sudo mv dotfiles-installer /usr/local/bin/
```

#### Linux (ARM64)
```bash
curl -L -o dotfiles-installer.tar.gz https://github.com/MrPointer/dotfiles/releases/latest/download/dotfiles-installer-*-linux-arm64.tar.gz
tar -xzf dotfiles-installer.tar.gz
chmod +x dotfiles-installer
sudo mv dotfiles-installer /usr/local/bin/
```

### 3. Using Go Install

If you have Go installed and want the latest development version:

```bash
go install github.com/MrPointer/dotfiles/installer@latest
```

### 4. Build from Source

For development or if you want to modify the installer:

```bash
git clone https://github.com/MrPointer/dotfiles.git
cd dotfiles/installer
go build -o dotfiles-installer .
```

## Verification

### Check Installation
```bash
dotfiles-installer version
```

### Verify Binary Signature (Optional)

All release binaries are signed with cosign. To verify:

```bash
# Install cosign first
# macOS: brew install cosign
# Linux: See https://docs.sigstore.dev/cosign/installation

# Verify the binary
cosign verify-blob \
  --certificate-identity-regexp 'https://github.com/MrPointer/dotfiles' \
  --certificate-oidc-issuer 'https://token.actions.githubusercontent.com' \
  --bundle dotfiles-installer.bundle \
  dotfiles-installer
```

## Usage

Once downloaded, you can use the dotfiles installer:

```bash
# Check system compatibility
dotfiles-installer check-compatibility

# Install dotfiles (interactive)
dotfiles-installer install

# Install dotfiles (non-interactive)
dotfiles-installer install --non-interactive --install-prerequisites

# Show help
dotfiles-installer --help
```

## Supported Platforms

The installer supports the following platforms:

- **macOS**
  - ARM64 (Apple Silicon - M1, M2, M3, etc.)
- **Linux**
  - x86_64 (AMD64)
  - ARM64 (AArch64)

## Troubleshooting

### Binary Not Found After Installation

If the binary is not found after installation, ensure the install directory is in your PATH:

```bash
# Check if directory is in PATH
echo $PATH | grep -o "$HOME/.local/bin"

# Add to PATH if missing (add to your shell's RC file)
export PATH="$HOME/.local/bin:$PATH"
```

### Permission Denied

If you get permission denied errors:

```bash
# Make sure the binary is executable
chmod +x /path/to/dotfiles-installer

# For system-wide setup, use sudo
sudo mv dotfiles-installer /usr/local/bin/
```

### Download Issues

If downloads fail:

1. Check your internet connection
2. Try using `wget` instead of `curl`:
   ```bash
   wget https://github.com/MrPointer/dotfiles/releases/latest/download/dotfiles-installer-*-linux-x86_64.tar.gz
   ```
3. Download from the [releases page](https://github.com/MrPointer/dotfiles/releases) manually

### Version Mismatch

To ensure you have the latest version:

```bash
# Check current version
dotfiles-installer version

# Download latest version again
curl -fsSL https://raw.githubusercontent.com/MrPointer/dotfiles/main/get.sh | bash
```

## Uninstalling

To remove the dotfiles installer:

```bash
# Remove the binary
rm $(which dotfiles-installer)

# Or if installed in ~/.local/bin
rm ~/.local/bin/dotfiles-installer
```

Note: This only removes the installer binary, not your installed dotfiles themselves. To remove dotfiles, see the main documentation.

## Getting Help

- **Documentation**: Check the main [README](../README.md)
- **Issues**: Report bugs on [GitHub Issues](https://github.com/MrPointer/dotfiles/issues)
- **Discussions**: Ask questions in [GitHub Discussions](https://github.com/MrPointer/dotfiles/discussions)

## Security

- All binaries are built in GitHub Actions with reproducible builds
- Binaries are signed with cosign using GitHub's OIDC provider
- Checksums are provided for all releases
- Source code is available for audit

For security concerns, please see our [security policy](../SECURITY.md).