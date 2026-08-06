# Domain

## Overview

This document defines the dotfile and shell-runtime concepts used by the chezmoi source. The dotfiles support personal and work environments, with conditional rendering driven by the shared [chezmoi data contract][chezmoi-data-contract].

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

- **Work environment is all-or-nothing at the profile level**: If `work_env` is true, the entire generic profile is sourced. Individual work features cannot be toggled independently.
- **Templates consume shared configuration**: Dotfile templates read the values defined by the [chezmoi data contract][chezmoi-data-contract] and do not write the chezmoi config.

## Glossary

| Term | Definition |
|------|-----------|
| Deferred brew loading | Pattern where Homebrew's shell environment setup is postponed from `.zshenv` to `.zshrc` on macOS |
| Generic work profile | Shared work configuration at `~/.work/profile`, loaded for any employer |
| Specific work profile | Employer-specific configuration at `~/.work/{work_name}/profile` |
| Work environment | The conditionally rendered and loaded work-specific dotfile configuration |
| Work name | Short employer identifier (e.g., `sedg`) used in paths and environment variable prefixes |

[shell-startup]: processes/shell-startup.md
[work-env-loading]: processes/work-environment-loading.md
[chezmoi-data-contract]: ../contracts/chezmoi-data.md
