# MrPointer's dotfiles

Personal dotfiles managed with [chezmoi], applied via a custom Go installer.
Supports both personal and work environments on macOS and Linux.

## Table of Contents

- [Overview](#overview)
- [Quick Start](#quick-start)
- [Installation Options](#installation-options)
- [Alternative Installation Methods](#alternative-installation-methods)
- [Installation Process](#installation-process)

## Overview

Like any other dotfiles project, this is a templated solution for applying a consistent
environment across new Unix machines. A dotfiles manager ([chezmoi]) handles templating
and per-machine differences, while a dedicated Go installer binary automates the full
setup - from prerequisites to shell configuration.

## Quick Start

Download and run the installer in one command:

```bash
curl -fsSL https://raw.githubusercontent.com/MrPointer/dotfiles/main/get.sh | bash -s -- --run
```

Or download first, then run manually:

```bash
curl -fsSL https://raw.githubusercontent.com/MrPointer/dotfiles/main/get.sh | bash
dotfiles-installer install
```

## Installation Options

### Install Command Options

| Option | Description | Default |
| --- | --- | --- |
| `--work-env` | Treat this installation as a work environment | `false` |
| `--work-name` | Work environment name | `sedg` |
| `--work-email` | Work email address | `timor.gruber@solaredge.com` |
| `--shell` | Shell to install and set as default | `zsh` |
| `--install-brew` | Install Homebrew if not present | `true` |
| `--shell-source` | Where to find the shell: `auto`, `brew`, `system` | `auto` |
| `--multi-user-system` | Configure for multi-user system | `false` |
| `--git-clone-protocol` | Git protocol for operations | `https` |
| `--install-prerequisites` | Automatically install missing prerequisites | `false` |

### Global Options

These work with any command:

| Option | Description |
| --- | --- |
| `-v, --verbose` | Enable verbose output (use `-vv` for extra verbose) |
| `--plain` | Show plain text instead of progress indicators |
| `--non-interactive` | Disable interactive prompts |
| `--extra-verbose` | Enable maximum verbosity |

### Examples

```bash
# Work environment
dotfiles-installer install --work-env --work-email your.email@company.com

# Non-interactive with prerequisites
dotfiles-installer install --non-interactive --install-prerequisites

# Check compatibility before installing
dotfiles-installer check-compatibility
dotfiles-installer install
```

### Get Script Options

The `get.sh` script itself accepts these flags:

| Option | Description | Default |
| --- | --- | --- |
| `-d, --dir` | Download directory | `$HOME/.local/bin` |
| `-v, --version` | Specific version to download | latest |
| `-r, --run` | Run the installer after download | `false` |

Pass installer flags after `--`:

```bash
curl -fsSL .../get.sh | bash -s -- --run -- --work-env --install-prerequisites
```

## Alternative Installation Methods

<details>
<summary>Manual download from GitHub Releases</summary>

Download pre-built binaries from [GitHub Releases](https://github.com/MrPointer/dotfiles/releases).

**macOS (Apple Silicon):**
```bash
curl -L -o dotfiles-installer.tar.gz https://github.com/MrPointer/dotfiles/releases/latest/download/dotfiles-installer-*-darwin-arm64.tar.gz
tar -xzf dotfiles-installer.tar.gz && chmod +x dotfiles-installer
```

**Linux (x86_64):**
```bash
curl -L -o dotfiles-installer.tar.gz https://github.com/MrPointer/dotfiles/releases/latest/download/dotfiles-installer-*-linux-x86_64.tar.gz
tar -xzf dotfiles-installer.tar.gz && chmod +x dotfiles-installer
```

**Linux (ARM64):**
```bash
curl -L -o dotfiles-installer.tar.gz https://github.com/MrPointer/dotfiles/releases/latest/download/dotfiles-installer-*-linux-arm64.tar.gz
tar -xzf dotfiles-installer.tar.gz && chmod +x dotfiles-installer
```
</details>

<details>
<summary>Build from source</summary>

```bash
git clone https://github.com/MrPointer/dotfiles.git
cd dotfiles/installer
go build -o dotfiles-installer .
./dotfiles-installer install
```
</details>

<details>
<summary>Using go install</summary>

```bash
go install github.com/MrPointer/dotfiles/installer@latest
dotfiles-installer install
```
</details>

## Installation Process

The installer walks through these steps:

1. **System Compatibility Check** - Verifies the system can run the dotfiles
2. **Prerequisites Installation** - Installs required tools and dependencies
3. **Homebrew Setup** - Installs Homebrew on macOS (optional on Linux)
4. **Shell Installation** - Installs and configures the specified shell
5. **GPG Setup** - Configures GPG keys for secure operations
6. **Dotfiles Manager Setup** - Installs and configures chezmoi
7. **Template Application** - Applies the dotfiles with user-specific configuration

Real-time progress indicators and detailed logging track each step.

[chezmoi]: https://www.chezmoi.io/
