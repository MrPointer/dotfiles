# Chezmoi Data Contract

## Overview

This document is the canonical contract between the Go installer and the chezmoi source. The installer produces `~/.config/chezmoi/chezmoi.toml`; chezmoi exposes the entries below to templates without the TOML `data` prefix. The templates and `.chezmoiignore` consume those values while rendering the target state.

The installer owns production of this file. Dotfile templates are read-only consumers. Running installation initializes the file from current installer input rather than reconciling arbitrary manual changes.

## Contract Boundary

| Side | Implementation | Responsibility |
|------|----------------|----------------|
| Producer | [`installer/lib/dotfilesmanager.DotfilesData`][dotfiles-data] and [`chezmoi.ChezmoiManager.Initialize`][data-initializer] | Assemble installer state and map it to `data.*` keys in the TOML config |
| Transport | `~/.config/chezmoi/chezmoi.toml` | Persist the data used for one or more chezmoi renders |
| Consumers | `.chezmoiignore`, `*.tmpl`, and templated source files | Read `.personal.*`, `.system.*`, and other produced namespaces while deciding what to render |

Chezmoi's built-in `.chezmoi.*` values are not part of this installer-owned contract.

## Data Schema

### `personal`

| TOML key | Template access | Presence | Producer source | Current consumers |
|----------|-----------------|----------|-----------------|-------------------|
| `data.personal.email` | `.personal.email` | Always | `DotfilesData.Email` | Common Git identity template |
| `data.personal.full_name` | `.personal.full_name` | Always | `FirstName` and `LastName`, joined with a space | Common Git identity template |
| `data.personal.work_env` | `.personal.work_env` | Always | Whether `DotfilesData.WorkEnv` is present | `.chezmoiignore`, shell templates, and work-dependent agent/Git templates |
| `data.personal.work_name` | `.personal.work_name` | Work environments only | `DotfilesWorkEnvData.WorkName` | Work profile extensions and work-dependent agent/Git templates |
| `data.personal.work_email` | `.personal.work_email` | Work environments only | `DotfilesWorkEnvData.WorkEmail` | No current source-template consumer |

### `system`

| TOML key | Template access | Presence | Producer source | Current consumers |
|----------|-----------------|----------|-----------------|-------------------|
| `data.system.shell` | `.system.shell` | Always in the install flow | `DotfilesSystemData.Shell` | No current source-template consumer |
| `data.system.work_generic_dotfiles_profile` | `.system.work_generic_dotfiles_profile` | Work environments only | `$HOME/.work/profile` | Main `.zshenv` and the work VPN setup script |
| `data.system.work_specific_dotfiles_profile` | `.system.work_specific_dotfiles_profile` | Work environments only | `$HOME/.work/{work_name}/profile` | Generic work profile template |

### `gpg`

| TOML key | Template access | Presence | Producer source | Current consumers |
|----------|-----------------|----------|-----------------|-------------------|
| `data.gpg.signing_key` | `.gpg.signing_key` | Only when interactive GPG setup selects or creates a key | `DotfilesData.GpgSigningKey` | No current source-template consumer |

## Work-Environment Invariants

When `personal.work_env` is `false`, the installer omits `work_name`, `work_email`, and both work-profile paths. `.chezmoiignore` excludes the managed work tree and work-dependent agent configuration in that mode.

When `personal.work_env` is `true`, the installer writes the work identity and both profile paths. Templates may therefore treat those conditional values as a group. The dotfile meaning of the two profile paths is documented in the [dotfiles domain][dotfiles-domain], and their runtime use is documented by [work environment loading][work-environment-loading].

## Current Contract Drift

The installer writes the selected GPG key as `data.gpg.signing_key`, but [`dot_gitconfig.tmpl`][gitconfig-template] and [`dot_zshrc.tmpl`][zshrc-template] currently test and read `personal.signing_key`. No template reads `.gpg.signing_key`, and the installer does not write `data.personal.signing_key`. Consequently, installer-selected GPG data does not currently activate those signing-key template blocks.

This document records the implemented mismatch; it does not redefine either side. A code change must align the producer and consumers before documentation can describe GPG signing as connected end to end.

## Change Rules

Changes to this boundary are contract changes. Keep the producer and every consumer aligned in one change:

1. Update the installer data type when new input or optionality is required.
2. Update `ChezmoiManager.Initialize` with the exact `data.*` key and presence rule.
3. Update all source consumers to use the corresponding template path.
4. Update this schema, including current consumers and work-environment invariants.
5. Validate both work and personal renders so conditional keys are neither missing nor read outside their valid mode.

Installer-only concepts such as package resolution, display modes, shell source strategy, and optional-tool selection are outside this contract because templates do not consume them. See the [installer domain][installer-domain].

[dotfiles-data]: ../../installer/lib/dotfilesmanager/data.go
[data-initializer]: ../../installer/lib/dotfilesmanager/chezmoi/data.go
[dotfiles-domain]: ../dotfiles/domain.md
[work-environment-loading]: ../dotfiles/processes/work-environment-loading.md
[gitconfig-template]: ../../dot_gitconfig.tmpl
[zshrc-template]: ../../dot_zshrc.tmpl
[installer-domain]: ../../installer/docs/domain.md
