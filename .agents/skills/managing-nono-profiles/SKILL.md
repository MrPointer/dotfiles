---
name: managing-nono-profiles
description: Manage nono sandbox profiles and their shell launch wrappers in this dotfiles repo. Use when editing private_dot_config/nono/profiles, nono profile schemas, dot_agents/shell wrappers, employer-specific work-* profiles and launchers under private_dot_work, SSH signing access for sandboxed agents, or adding/renaming agent-specific nono profiles.
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
| Employer profiles | `private_dot_config/nono/profiles/work-{work_name}-*.jsonc` | `~/.config/nono/profiles/work-{work_name}-*.jsonc` |
| Employer profile filtering | `.chezmoiignore` | — |
| Employer launchers | `private_dot_work/private_{work_name}/agents/shell/` | `~/.work/{work_name}/agents/shell/` |
| Employer launcher sourcing | `private_dot_work/private_{work_name}/zsh/dot_zshrc.tmpl` | `~/.work/{work_name}/zsh/.zshrc` |

`.zshrc` sources `~/.agents/shell/nono.zsh` first, then per-agent shell files, only when `nono` is installed. The employer zshrc extension sources the employer launchers, also only when `nono` is installed.

## Core Rule

Profiles intended for interactive use need matching shell functions.

When adding or renaming a profile:

1. Add or update the profile in `private_dot_config/nono/profiles/`.
2. Keep the profile filename, `meta.name`, and wrapper profile argument aligned.
3. Add or update the matching function in the relevant file under `dot_agents/shell/`.
4. If the profile is for a new agent shell file, add that file to the source list in `dot_zshrc.tmpl` after `nono.zsh`.

Example: profile `claude-code-go` should have a profile file named `claude-code-go.jsonc`, `meta.name: "claude-code-go"`, and a wrapper that calls `nono run --profile claude-code-go`.

## Employer-Specific Profiles

Profiles that serve a single employer or employer project are named `work-{work_name}-<profile>` and are installed only on that employer's machines. Their launchers belong to the employer's work tree, not to `dot_agents/shell/`.

1. Add the profile to `private_dot_config/nono/profiles/` as `work-{work_name}-<profile>.jsonc`, with the same `meta.name`. List `work-{work_name}-common` first in `extends`, followed by the personal profiles the project needs.
2. Make sure `.chezmoiignore` has a block that ignores `.config/nono/profiles/work-{work_name}-*` unless `personal.work_env` is true and `personal.work_name` matches. Each employer needs one block, not one entry per profile.
3. Add the launcher function to the matching agent file in `private_dot_work/private_{work_name}/agents/shell/`, calling the shared helper with the profile name (for example `_claude_nono work-sedg-claude-code-energy-runtime`).
4. For a new employer launcher file, source it from `private_dot_work/private_{work_name}/zsh/dot_zshrc.tmpl` inside the existing `command -v nono` block.

Each employer project gets its own profile and launcher; there are no employer variants of the general launchers such as `claude-go`. `work-{work_name}-common` is agent-neutral: it holds the access every project of that employer needs regardless of agent (employer domains, credentials, tool config), extends nothing, and is shared by employer profiles of any agent. Put agent-specific employer access in a separate `work-{work_name}-<agent>-common` profile. Keep employer access out of personal profiles such as `my-claude-code`.

Employer profiles may extend personal profiles by name (for example `claude-code-go`). Renaming a personal profile requires updating those `extends` entries.

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
- Use the shared `~/.config/git/agent.gitconfig` for agent Git configuration.
- Let `_nono_agent_run` attempt to preload the key into the host `ssh-agent` before entering nono. This is best-effort, not a fail-fast check; missing sockets, missing keys, or failed loads do not prevent launch.
- Let `_nono_agent_run` pass `--allow-unix-socket "$SSH_AUTH_SOCK"` when the socket variable is set, so the sandbox can ask the host agent to sign without reading private key files.

Public SSH metadata commonly needed by Git:

- `$HOME/.ssh/<key>.pub`
- `$HOME/.ssh/known_hosts`
- `$HOME/.ssh/config`

Profiles also need read access to the shared agent, common, and personal Git configs, `allowed_signers`, and `nono-git-ssh`. Grant the required public SSH metadata through both `read_file` and `bypass_protection`, while keeping private keys excluded.

Validate private keys remain blocked with:

```sh
nono why --profile <profile-or-file> --path "$HOME/.ssh/<private-key>" --op read
```

Expected result: `DENIED` by `deny_credentials`.

## Shell Wrapper Pattern

Shared nono wrapper helpers live in `dot_agents/shell/nono.zsh.tmpl`.

`_nono_agent_run` owns `GIT_CONFIG_GLOBAL="$HOME/.config/git/agent.gitconfig"` for all agents. Per-agent wrappers pass the profile, command, preload key paths, and command arguments. The `OPENCODE_SSH_SIGNING_*`, `CLAUDE_SSH_SIGNING_*`, and `CODEX_SSH_SIGNING_*` environment overrides affect preloading only; they do not change Git's configured signing key.

Per-agent files should call `_nono_agent_run` instead of duplicating nono setup:

- `dot_agents/shell/opencode.zsh`
- `dot_agents/shell/claude.zsh`
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
