# Installer Domain

## Overview

The installer domain covers the concepts used to bootstrap a supported machine before chezmoi renders the dotfiles. It owns platform compatibility, package selection, shell installation strategy, optional tools, and command-output behavior. The values handed from the installer to chezmoi are defined separately by the shared [chezmoi data contract][chezmoi-data-contract].

## Key Concepts

### Package Resolution

Package resolution translates a platform-independent package key into a package that the active package manager can install.

- **Abstract package key**: A stable name such as `build-essential`, `gpg`, or `uv` used by installer configuration.
- **Package mapping**: An entry in [`packagemap.yaml`][packagemap-yaml] that associates an abstract key with package-manager-specific configuration.
- **Manager-specific mapping**: The concrete name for `apt`, `dnf`, or `brew`, optionally classified by package type or specialized by distribution.

| Package type | Meaning |
|--------------|---------|
| Empty | A regular package installed by the manager's normal install operation |
| `group` | A package-manager group, used by DNF group installation |
| `pattern` | A package-manager installation pattern |

The [package resolution process][package-resolution] describes the lookup and failure behavior.

### Prerequisite

A prerequisite is a command required for compatibility or installation. [`compatibility.yaml`][compatibility-yaml] associates prerequisites with supported operating systems and distributions. Missing prerequisites can be resolved and installed before the compatibility check is repeated.

### Optional Tool

An optional tool is a daily-use CLI tool that may enhance the resulting environment but is not required to apply the dotfiles. [`tools.yaml`][tools-yaml] supplies each tool's abstract package key and user-facing description.

Optional tools use the same package resolver as prerequisites. A tool is available to the selection step only when its key resolves for the active package manager. Selections are not written to chezmoi data, and individual installation failures do not make the main installation fail.

`uv` follows this ordinary optional-tool model. Its current package mapping is Homebrew-only; it does not imply installer-managed Python versions or migration of existing Python state.

See [optional tools installation][tools-installation] for the end-to-end flow.

### Shell Source Strategy

The shell source strategy chooses both the package manager used to install the requested shell and the location used to resolve its executable.

| Strategy | Flag value | Meaning |
|----------|------------|---------|
| Auto | `auto` | Prefer Homebrew when the installer has a Homebrew path; otherwise use the supported native package manager |
| Homebrew | `brew` | Require Homebrew and resolve the shell under its prefix |
| System | `system` | Require the supported native package manager and resolve the shell from system locations |

The [shell setup process][shell-setup] documents path resolution and default-shell changes.

### Display Mode

A display mode controls how the installer presents progress and external command output.

| Mode | Behavior |
|------|----------|
| Progress | Show interactive progress indicators and hide command output |
| Plain | Show simple progress messages without spinners and hide command output |
| Passthrough | Send command output directly to the terminal |

Progress and Plain discard external command output; Passthrough preserves it.

## Domain Rules

- **Resolution uses the active manager**: A package key must have a mapping for the package manager selected by the installer.
- **Distribution mappings are exact**: When a manager mapping contains distribution-specific names, the detected distribution must have an explicit entry; the resolver does not guess a fallback.
- **Optional tools remain optional**: Tool configuration and tool-install failures do not alter the installer-to-template data contract.
- **Shell source is explicit**: Values other than `auto`, `brew`, and `system` are invalid.

## Glossary

| Term | Definition |
|------|------------|
| Abstract package key | Platform-independent identifier resolved before package installation |
| Active package manager | The package-manager implementation selected for the current installer step |
| Display mode | Installer policy for progress presentation and external command output |
| Optional tool | Non-required CLI tool offered after dotfiles setup when it resolves for the active manager |
| Package mapping | Configuration that maps an abstract key to manager-specific package names and types |
| Prerequisite | Required command checked as part of platform compatibility |
| Shell source strategy | Policy that chooses where the requested shell is installed and resolved |

[chezmoi-data-contract]: ../../docs/contracts/chezmoi-data.md
[compatibility-yaml]: ../internal/config/compatibility.yaml
[packagemap-yaml]: ../internal/config/packagemap.yaml
[tools-yaml]: ../internal/config/tools.yaml
[package-resolution]: processes/package-resolution.md
[shell-setup]: processes/shell-setup.md
[tools-installation]: processes/tools-installation.md
