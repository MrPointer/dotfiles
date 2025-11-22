# MrPointer's dotfiles

## Motivation

Like any other dotfiles project, I'm looking to create myself a templated solution that will help me
apply it on new environments, mostly Unix ones.  
I'm using a dotfiles manager alongside a custom installer binary to achieve this,
managing both home and office/work environments.

## Installation

The dotfiles are installed using a dedicated Go binary that handles system setup and configuration.

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
go install github.com/MrPointer/dotfiles/installer@main
dotfiles-installer install
```

### Pre-built Releases

Pre-built binaries will be available in GitHub releases once the first version is tagged.
Until then, use the build-from-source method above.

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

**Basic installation:**

```bash
./dotfiles-installer install
```

**Work environment setup:**

```bash
./dotfiles-installer install --work-env --work-email your.email@company.com
```

**Non-interactive installation with prerequisites:**

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
