# Domain

## Overview

This project automates personal machine setup: an installer bootstraps prerequisites and applies dotfiles via chezmoi. The dotfiles support both personal and work environments, with conditional configuration driven by data the installer collects at install time.

## Key Concepts

### Work Environment

The project implements a two-tier model for separating personal and work configuration.

- **Activation**: Controlled by the `personal.work_env` boolean in chezmoi data. When `false`, all work-related configuration is skipped.
- **Work Name**: A short identifier for the employer (e.g., `sedg`). Used to locate employer-specific files and derive environment variable prefixes.

**Two-tier profile system:**

| Layer | Path | Purpose |
|-------|------|---------|
| Generic work profile | `~/.work/profile` | Cross-employer shared config — sets `WORK_DIR`, `WORK_EXECUTABLES_DIR`, and shell extension paths (`WORK_ZSH_ENV_EXTENSION`, `WORK_ZSH_RC_EXTENSION`) |
| Specific work profile | `~/.work/{work_name}/profile` | Employer-specific config — sourced by the generic profile |

The generic profile delegates to the specific profile. Both layers produce shell extension files (`.zshenv` and `.zshrc` fragments) that are sourced during shell startup. See the [work environment loading process][work-env-loading] for the full loading flow.

**File structure in chezmoi source:**

```
private_dot_work/
├── profile.tmpl              # Generic work profile (sets vars, sources specific profile)
├── zsh/
│   ├── dot_zshenv.tmpl       # Generic work zshenv extension
│   └── dot_zshrc.tmpl        # Generic work zshrc extension
└── private_{work_name}/
    ├── profile               # Employer-specific profile
    └── zsh/
        ├── dot_zshenv.tmpl   # Employer-specific zshenv extension
        └── dot_zshrc.tmpl    # Employer-specific zshrc extension
```

### Chezmoi Data Schema

The installer collects user input and writes it to chezmoi's config file. Templates then consume these values. The data is organized into three namespaces:

| Namespace | Set By | Purpose |
|-----------|--------|---------|
| `personal` | Installer | User identity — `email`, `full_name`, `work_env`, `work_name`, `work_email` |
| `system` | Installer | Runtime paths — `shell`, `work_generic_dotfiles_profile`, `work_specific_dotfiles_profile` |
| `gpg` | Installer | Security — `signing_key` |

**Installer-side types** (in `dotfilesmanager.DotfilesData`):

- `WorkEnv` is `Option[DotfilesWorkEnvData]` — when present, work mode is enabled
- `SystemData` is `Option[DotfilesSystemData]` — contains shell name and optional work profile paths
- `GpgSigningKey` is `Option[string]` — set when the user selects or creates a GPG key

**Contract**: The installer writes `data.personal.*`, `data.system.*`, and `data.gpg.*` keys into chezmoi's YAML config. Templates read them via `.personal.*`, `.system.*`, and `.gpg.*`. Adding a new data key requires updating both `DotfilesData` (Go struct) and the `Initialize` method that maps struct fields to viper keys. See the [installation process][installation] for how these are populated.

### Package Resolution

The installer resolves abstract package names to concrete, installable names for the target system. This decouples prerequisite definitions from platform-specific package naming.

- **Abstract package key**: A platform-independent name (e.g., `build-essential`, `gpg`). Used in compatibility checks and prerequisite lists.
- **Package mapping**: A YAML entry in [`packagemap.yaml`][packagemap-yaml] that maps an abstract key to manager-specific configurations.
- **Manager-specific mapping**: The concrete package name for a given package manager (`apt`, `dnf`, `brew`), optionally with a type and distro-specific name overrides.

**Package types:**

| Type | Meaning | Example |
|------|---------|---------|
| *(empty)* | Regular package | `apt install git` |
| `group` | Package group (DNF concept) | `dnf groupinstall "Development Tools"` |
| `pattern` | Installation pattern | Used by zypper-style managers |

**Example mapping:**

```yaml
build-essential:
  apt:
    name: build-essential
  dnf:
    type: group
    name:
      fedora: development-tools
      centos: "Development Tools"
```

See the [package resolution process][pkg-resolution] for the resolution flow from abstract key to installable package.

### Shell Source Strategy

When installing the user's shell, the installer supports three strategies for locating and installing it:

| Strategy | Flag value | Behavior |
|----------|-----------|----------|
| Auto | `auto` (default) | Try Homebrew first if available, fall back to native package manager |
| Brew | `brew` | Use Homebrew exclusively — fails if brew is not installed |
| System | `system` | Use native package manager exclusively (apt, dnf) |

The strategy determines both the installation source and which binary path is registered as the user's default shell (brew-installed shells live under the brew prefix, system shells live in `/bin` or `/usr/bin`).

### Display Modes

The installer supports three output modes for controlling how external command output is presented:

| Mode | Behavior |
|------|----------|
| Progress | Interactive spinners, command output hidden (default) |
| Plain | Simple text messages, command output hidden (non-interactive/CI) |
| Passthrough | All command output shown directly (debugging) |

Progress and Plain discard command output; Passthrough does not.

### Deferred Homebrew Loading

A shell startup optimization pattern where Homebrew's expensive `shellenv` evaluation is postponed from `.zshenv` to `.zshrc` on macOS. This avoids the cost in non-interactive shells while keeping the environment available for interactive use.

**Key state:**

- `DEFER_BREW_LOAD` — flag set in `.zshenv` to signal that loading should happen later
- `BREW_LOADED` — guard that prevents double-evaluation across sourced files

**Platform behavior:**

| Platform | Behavior | Reason |
|----------|----------|--------|
| macOS | Deferred to `.zshrc` | Optimize non-interactive shell startup |
| Linux | Loaded immediately in `.zshenv` | Homebrew is less common; consistent PATH needed |
| Devbox | Simple PATH addition | `eval` is unnecessary in devbox environments |

See the [shell startup process][shell-startup] for the full shell startup flow including deferred loading.

## Domain Rules

- **Package resolution is exact-match**: Distro-specific name resolution has no fallback — if the exact distro name isn't in the mapping, resolution fails. This is intentional to prevent silent mismatches.
- **Work environment is all-or-nothing at the profile level**: If `work_env` is true, the entire generic profile is sourced. Individual work features cannot be toggled independently.
- **Installer populates, templates consume**: The installer is the sole writer of chezmoi data. Templates are read-only consumers. Manual edits to chezmoi's config file will be overwritten on next install.

## Glossary

| Term | Definition |
|------|-----------|
| Abstract package key | Platform-independent package name used in config files, resolved to a concrete name at install time |
| Chezmoi data | Key-value pairs written to chezmoi's config file by the installer, consumed by templates via `.personal.*`, `.system.*`, `.gpg.*` |
| Deferred brew loading | Pattern where Homebrew's shell environment setup is postponed from `.zshenv` to `.zshrc` on macOS |
| Display mode | Installer output verbosity level: progress (spinners), plain (text), or passthrough (raw output) |
| Generic work profile | Shared work configuration at `~/.work/profile`, loaded for any employer |
| Package mapping | YAML entry mapping an abstract package key to manager-specific names and types |
| Package type | Classification of a package: regular (default), group (DNF group install), or pattern |
| Shell source strategy | How the installer finds/installs the shell: auto, brew, or system |
| Specific work profile | Employer-specific configuration at `~/.work/{work_name}/profile` |
| Work name | Short employer identifier (e.g., `sedg`) used in paths and environment variable prefixes |

[installation]: processes/installation.md
[pkg-resolution]: processes/package-resolution.md
[shell-startup]: processes/shell-startup.md
[work-env-loading]: processes/work-environment-loading.md
[packagemap-yaml]: ../installer/internal/config/packagemap.yaml
