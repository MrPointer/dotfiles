# Chezmoi Data Contract

## Overview

This document is the canonical contract between the Go installer and the chezmoi source. The installer initially produces `~/.config/chezmoi/chezmoi.toml`; chezmoi exposes the entries below to templates without the TOML `data` prefix. The templates and `.chezmoiignore` consume those values while rendering the target state.

The installer supplies the machine-specific `data.*` values during installation. On a later `chezmoi init`, the repository-root `.chezmoi.toml.tmpl` regenerates the config from the data already exposed to templates. Normal dotfile templates remain read-only consumers of the config.

## Contract Boundary

| Side | Implementation | Responsibility |
|------|----------------|----------------|
| Initial producer | [`installer/lib/dotfilesmanager.DotfilesData`][dotfiles-data] and [`chezmoi.ChezmoiManager.Initialize`][data-initializer] | Assemble installer state and map it to `data.*` keys in the TOML config |
| Init-time producer | [`.chezmoi.toml.tmpl`][config-template] | Regenerate the config during `chezmoi init` while carrying forward the available template data |
| Transport | `~/.config/chezmoi/chezmoi.toml` | Persist the data used for one or more chezmoi renders |
| Consumers | `.chezmoiignore` and normal dotfile templates | Read `.personal.*`, `.system.*`, and other namespaces while deciding what to render |

Chezmoi's built-in `.chezmoi.*` values are not part of this shared data contract.

### Init-Time Preservation

For an existing installation, run `chezmoi init` or `chezmoi update --init` once to adopt the repository config template. The template builds a `data` table from `omit . "chezmoi"` and serializes it with `toToml`, so personal, work, GPG, and unknown custom data keys survive repeated initialization. A fresh init with no existing data writes no installer-specific data; it does not invent identity, work-environment, system, or GPG values.

This preservation applies only to template data. Manually customized top-level config outside `data` is not carried forward; the template currently writes the preserved `data` plus `diff.exclude = ["scripts"]`. A normal `chezmoi apply` does not regenerate the config.

## Data Schema

### `personal`

| TOML key | Template access | Presence | Producer source | Current consumers |
|----------|-----------------|----------|-----------------|-------------------|
| `data.personal.email` | `.personal.email` | Always | `DotfilesData.Email` | Common Git personal-email fallback and signer templates |
| `data.personal.full_name` | `.personal.full_name` | Always | `FirstName` and `LastName`, joined with a space | Common Git identity template |
| `data.personal.work_env` | `.personal.work_env` | Always | Whether `DotfilesData.WorkEnv` is present | `.chezmoiignore`, work shell templates, work Git identity, and allowed signer identities |
| `data.personal.work_name` | `.personal.work_name` | Work environments only | `DotfilesWorkEnvData.WorkName` | Work shell profile extensions and employer nono profile filtering |
| `data.personal.work_email` | `.personal.work_email` | Work environments only | `DotfilesWorkEnvData.WorkEmail` | Work Git identity and signer templates (work environments) |

### `system`

| TOML key | Template access | Presence | Producer source | Current consumers |
|----------|-----------------|----------|-----------------|-------------------|
| `data.system.shell` | `.system.shell` | Always in the install flow | `DotfilesSystemData.Shell` | No current source-template consumer |
| `data.system.work_generic_dotfiles_profile` | `.system.work_generic_dotfiles_profile` | Work environments only | `$HOME/.work/profile` | Main `.zshenv` and the work VPN setup script |
| `data.system.work_specific_dotfiles_profile` | `.system.work_specific_dotfiles_profile` | Work environments only | `$HOME/.work/{work_name}/profile` | Generic work profile template |

### `gpg`

| TOML key | Template access | Presence | Producer source | Current consumers |
|----------|-----------------|----------|-----------------|-------------------|
| `data.gpg.signing_key` | `.gpg.signing_key` | Only when interactive GPG setup selects or creates a key | `DotfilesData.GpgSigningKey` | `dot_gitconfig.tmpl` (`user.signingkey`) and `dot_zshrc.tmpl` (`gpg-unlock` alias) |

## Work-Environment Invariants

When `personal.work_env` is `false`, the installer omits `work_name`, `work_email`, and both work-profile paths. `.chezmoiignore` excludes the managed work tree, employer-specific nono profiles, and the work Git configuration in that mode.

When `personal.work_env` is `true`, the installer writes the work identity and both profile paths. The work Git configuration overrides the common personal-email fallback only for matching configured GitLab and Azure DevOps repository URLs. Templates may therefore treat the conditional values as a group. The dotfile meaning of the two profile paths is documented in the [dotfiles domain][dotfiles-domain], and their runtime use is documented by [work environment loading][work-environment-loading].

## GPG Signing Key

The installer writes the selected GPG key as `data.gpg.signing_key`. [`dot_gitconfig.tmpl`][gitconfig-template] reads `.gpg.signing_key` to set `user.signingkey`, and [`dot_zshrc.tmpl`][zshrc-template] reads it to define the `gpg-unlock` alias. Both templates guard on the presence of the `gpg` namespace, so a personal or work render without a selected key omits both blocks.

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
[config-template]: ../../.chezmoi.toml.tmpl
[dotfiles-domain]: ../dotfiles/domain.md
[work-environment-loading]: ../dotfiles/processes/work-environment-loading.md
[gitconfig-template]: ../../dot_gitconfig.tmpl
[zshrc-template]: ../../dot_zshrc.tmpl
[installer-domain]: ../../installer/docs/domain.md
