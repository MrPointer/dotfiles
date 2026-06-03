---
name: managing-nono-profiles
description: Manage nono sandbox profiles and their shell launch wrappers in this dotfiles repo. Use when editing private_dot_config/nono/profiles, nono profile schemas, dot_agents/shell wrappers, SSH signing access for sandboxed agents, or adding/renaming agent-specific nono profiles.
---

# Managing Nono Profiles

Maintain nono sandbox profiles and the matching shell functions that launch coding agents with those profiles.

For chezmoi source/target mechanics, templates, `chezmoi apply`, or drift checks, also load the `managing-chezmoi` skill. This skill covers nono-specific profile and wrapper conventions only.

## Source Locations

Edit chezmoi source files, not applied targets:

| Purpose | Source path | Target path |
|---------|-------------|-------------|
| User profiles | `private_dot_config/nono/profiles/` | `~/.config/nono/profiles/` |
| Profile schema | `private_dot_config/nono/schemas/nono-profile.schema.json` | `~/.config/nono/schemas/nono-profile.schema.json` |
| Agent shell wrappers | `dot_agents/shell/` | `~/.agents/shell/` |
| Wrapper sourcing | `dot_zshrc.tmpl` | `~/.zshrc` |

`.zshrc` sources `~/.agents/shell/nono.zsh` first, then per-agent shell files, only when `nono` is installed.

## Core Rule

Profiles intended for interactive use need matching shell functions.

When adding or renaming a profile:

1. Add or update the profile in `private_dot_config/nono/profiles/`.
2. Keep the profile filename, `meta.name`, and wrapper profile argument aligned.
3. Add or update the matching function in the relevant file under `dot_agents/shell/`.
4. If the profile is for a new agent shell file, add that file to the source list in `dot_zshrc.tmpl` after `nono.zsh`.

Example: profile `claude-code-go` should have a profile file named `claude-code-go.jsonc`, `meta.name: "claude-code-go"`, and a wrapper that calls `nono run --profile claude-code-go`.

## Path Conventions

Prefer portable nono variables over `~` and absolute machine paths:

| Prefer | Avoid |
|--------|-------|
| `$HOME/.ssh/config` | `~/.ssh/config` |
| `$XDG_CONFIG_HOME/git/common.gitconfig` | `~/.config/git/common.gitconfig` |
| `$XDG_CACHE_HOME/uv` | `~/.cache/uv` |
| `$XDG_DATA_HOME/uv` | `~/.local/share/uv` |
| `$TMPDIR` | `/tmp` or macOS temp paths |

Use `{{ .personal.work_name }}` inline inside JSON strings for work-specific template paths. Do not put template variable declarations before the opening `{` in `.jsonc.tmpl` profile files; editors should still treat them as JSONC.

## SSH Signing Pattern

Do not grant sandboxed agents read access to private SSH keys.

Use this pattern instead:

- Grant exact public SSH metadata with `filesystem.read_file`.
- Add the same exact public/metadata files to `filesystem.bypass_protection` when they live under denied paths like `$HOME/.ssh`.
- Never add `$HOME/.ssh` as a directory grant or bypass unless explicitly requested and security-reviewed.
- Let shell wrappers preload keys into the host `ssh-agent` before entering nono.
- Let wrappers pass `--allow-unix-socket "$SSH_AUTH_SOCK"` so the sandbox can ask the host agent to sign without reading private key files.

Public SSH metadata commonly needed by Git:

- `$HOME/.ssh/<key>.pub`
- `$HOME/.ssh/known_hosts`
- `$HOME/.ssh/config`

Validate private keys remain blocked with:

```sh
nono why --profile <profile-or-file> --path "$HOME/.ssh/<private-key>" --op read
```

Expected result: `DENIED` by `deny_credentials`.

## Shell Wrapper Pattern

Shared nono wrapper helpers live in `dot_agents/shell/nono.zsh.tmpl`.

Per-agent files should call `_nono_agent_run` instead of duplicating nono setup:

- `dot_agents/shell/opencode.zsh`
- `dot_agents/shell/claude.zsh.tmpl`
- `dot_agents/shell/codex.zsh`

Use wrapper functions, not aliases, when SSH signing support or profile selection is needed.

## Validation

For plain profiles:

```sh
nono profile validate private_dot_config/nono/profiles/<profile>.jsonc
```

For templated profiles:

```sh
chezmoi execute-template < private_dot_config/nono/profiles/<profile>.jsonc.tmpl > /tmp/<profile>.jsonc
nono profile validate /tmp/<profile>.jsonc
```

For shell wrappers:

```sh
chezmoi execute-template < dot_agents/shell/nono.zsh.tmpl > /tmp/nono.zsh
zsh -n /tmp/nono.zsh
zsh -n dot_agents/shell/opencode.zsh
```

For `.zshrc` integration:

```sh
chezmoi execute-template < dot_zshrc.tmpl > /tmp/zshrc
zsh -n /tmp/zshrc
```

Clean up temporary rendered files after validation.

## Applying Changes

After editing source files, use chezmoi from an appropriate shell profile:

```sh
chezmoi apply --dry-run -v
chezmoi apply
```

If current opencode is already running inside nono, some target-side checks may be blocked by the active profile. Validate source files directly when target inspection is not available.
