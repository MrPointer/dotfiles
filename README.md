# MrPointer's dotfiles

## Motivation

Like any other dotfiles project, I'm looking to create myself a templated solution that will help me
apply it on new environments, mostly Unix ones.  
I'm using a dotfiles manager alongside a custom installer binary to achieve this,
managing both home and office/work environments.

## Installation

The dotfiles are installed using a dedicated Go binary that handles system setup and configuration.

### Quick Setup (Recommended)

Use our get script to download the installer binary:

```bash
curl -fsSL https://raw.githubusercontent.com/MrPointer/dotfiles/main/get.sh | bash
```

**One-command download and install dotfiles:**
```bash
curl -fsSL https://raw.githubusercontent.com/MrPointer/dotfiles/main/get.sh | bash -s -- --run
```

**Download and install dotfiles with custom options:**
```bash
curl -fsSL https://raw.githubusercontent.com/MrPointer/dotfiles/main/get.sh | bash -s -- --run -- --work-env --install-prerequisites
```

**Download to custom directory:**
```bash
curl -fsSL https://raw.githubusercontent.com/MrPointer/dotfiles/main/get.sh | bash -s -- --dir /usr/local/bin
```

### Manual Download

Download pre-built binaries from [GitHub Releases](https://github.com/MrPointer/dotfiles/releases):

**macOS (Apple Silicon):**
```bash
curl -L -o dotfiles-installer.tar.gz https://github.com/MrPointer/dotfiles/releases/latest/download/dotfiles-installer-*-darwin-arm64.tar.gz
tar -xzf dotfiles-installer.tar.gz
chmod +x dotfiles-installer
```

**Linux (x86_64):**
```bash
curl -L -o dotfiles-installer.tar.gz https://github.com/MrPointer/dotfiles/releases/latest/download/dotfiles-installer-*-linux-x86_64.tar.gz
tar -xzf dotfiles-installer.tar.gz
chmod +x dotfiles-installer
```

**Linux (ARM64):**
```bash
curl -L -o dotfiles-installer.tar.gz https://github.com/MrPointer/dotfiles/releases/latest/download/dotfiles-installer-*-linux-arm64.tar.gz
tar -xzf dotfiles-installer.tar.gz
chmod +x dotfiles-installer
```

### Build from Source

Clone the repository and build the installer:

```bash
git clone https://github.com/MrPointer/dotfiles.git
cd dotfiles/installer
go build -o dotfiles-installer .
./dotfiles-installer install
```

### Using Go Install

If you have Go installed:

```bash
go install github.com/MrPointer/dotfiles/installer@latest
dotfiles-installer install
```

## Usage

The installer provides several commands and options:

### Basic Commands

- `dotfiles-installer install` - Install dotfiles on the current system
- `dotfiles-installer check-compatibility` - Check system compatibility
- `dotfiles-installer --help` - Show all available commands and options

### Installation Options

The following options can be passed to the `install` command:

| Option                      | Description                                   | Default                      |
| --------------------------- | --------------------------------------------- | ---------------------------- |
| `--work-env`                | Treat this installation as a work environment | `false`                      |
| `--work-name`               | Work environment name                         | `sedg`                       |
| `--work-email`              | Work email address                            | `timor.gruber@solaredge.com` |
| `--shell`                   | Shell to install and set as default           | `zsh`                        |
| `--install-brew`            | Install Homebrew if not present               | `true`                       |
| `--install-shell-with-brew` | Install shell using Homebrew                  | `true`                       |
| `--multi-user-system`       | Configure for multi-user system               | `false`                      |
| `--git-clone-protocol`      | Git protocol for operations                   | `https`                      |
| `--install-prerequisites`   | Automatically install missing prerequisites   | `false`                      |

### Global Options

These options work with any command:

| Option              | Description                                         |
| ------------------- | --------------------------------------------------- |
| `-v, --verbose`     | Enable verbose output (use `-vv` for extra verbose) |
| `--plain`           | Show plain text instead of progress indicators      |
| `--non-interactive` | Disable interactive prompts                         |
| `--extra-verbose`   | Enable maximum verbosity                            |

### Example Usage

**Basic setup:**

```bash
./dotfiles-installer install
```

**Work environment installation:**

```bash
./dotfiles-installer install --work-env --work-email your.email@company.com
```

**Non-interactive dotfiles installation with prerequisites:**

```bash
./dotfiles-installer install --non-interactive --install-prerequisites --git-clone-protocol=https
```

**Check system compatibility first:**

```bash
./dotfiles-installer check-compatibility
./dotfiles-installer install
```

## Overview

### Dotfiles Manager

I'm using [chezmoi] as the dotfiles manager, which provides templating abilities,
per-machine differences, and much more. The installer sets up chezmoi and populates it with
the necessary configuration.

### Installation Process

The Go installer handles the complete setup process:

1. **System Compatibility Check** - Verifies the system can run the dotfiles
2. **Prerequisites Installation** - Installs required tools and dependencies
3. **Homebrew Setup** - Installs Homebrew on macOS (optional on Linux)
4. **Shell Installation** - Installs and configures the specified shell
5. **GPG Setup** - Configures GPG keys for secure operations
6. **Dotfiles Manager Setup** - Installs and configures chezmoi
7. **Template Application** - Applies the dotfiles with user-specific configuration

The installer provides real-time progress indicators and detailed logging to track the installation process.

[chezmoi]: https://www.chezmoi.io/
