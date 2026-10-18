# Agent Commit Signing

## Overview

Signs Git commits from nono-launched OpenCode, Claude Code, and Codex with the host's Ed25519 key through `ssh-agent`. The [shared agent Git configuration][agent-git] selects the signing key, while the [common configuration][common-git] supplies the shared name and personal-email fallback. See [agent Git signing][domain-signing] for identity selection and configuration ownership.

## Trigger

A coding agent starts through a shell wrapper and later runs a Git commit with signing enabled.

## Actors

- **Host user**: Registers the public signing key with the Git hosting account and makes the private key available to the host `ssh-agent`.
- **[Shared launcher][shared-launcher]**: Attempts key preloading, sets the agent Git config, and passes the agent socket to nono.
- **Nono profiles**: Permit public metadata and socket access without private-key reads.
- **Git and `ssh-agent`**: Request and produce the SSH signature.
- **Local Git or Git host**: Verify the signature using separate trust configuration.

## Flow

### Host Setup

1. Register `~/.ssh/id_ed25519.pub` as a signing key on each Git hosting account that needs to verify its signatures. Where the host supports multiple signing keys, register each machine's public key with the account.
2. Ensure the host `ssh-agent` is available through `SSH_AUTH_SOCK` and can use the matching private key. Key loading happens on the host, before the sandbox starts; private keys are not granted to the sandbox.
3. Configure repository authentication separately. Registering a signing key does not by itself grant clone, fetch, or push access, even when the same public key is used for authentication.

### Agent Launch and Signing

1. The [OpenCode][opencode-wrapper], [Claude Code][claude-wrapper], or [Codex][codex-wrapper] wrapper passes its profile, command, and preload key paths to `_nono_agent_run`. The defaults are `~/.ssh/id_ed25519.pub` and `~/.ssh/id_ed25519`.
2. `_nono_agent_ensure_ssh_signing_key` checks the socket, `ssh-add`, and public key, then tests agent signing with `ssh-add -T`. If needed and readable on the host, it attempts to load the private key. On macOS it tries `--apple-use-keychain` before ordinary `ssh-add`.
3. `_nono_agent_run` launches nono with `GIT_CONFIG_GLOBAL="$HOME/.config/git/agent.gitconfig"`. If `SSH_AUTH_SOCK` is set, it adds `--allow-unix-socket "$SSH_AUTH_SOCK"`.
4. Git reads the shared agent config, which includes `common.gitconfig`, enables SSH signature format, and selects `~/.ssh/id_ed25519.pub`. The common config enables commit signing. Git asks the host agent to sign using the corresponding private key.
5. For Git SSH transport, `core.sshCommand` selects [`nono-git-ssh`][git-ssh], which uses `ncat` for proxy-aware SSH connections. This transport setting is separate from commit signing.

Preloading is best-effort: missing prerequisites or failed key loading do not stop agent launch. The `OPENCODE_SSH_SIGNING_PUBLIC_KEY` / `OPENCODE_SSH_SIGNING_KEY`, `CLAUDE_SSH_SIGNING_PUBLIC_KEY` / `CLAUDE_SSH_SIGNING_KEY`, and `CODEX_SSH_SIGNING_PUBLIC_KEY` / `CODEX_SSH_SIGNING_KEY` overrides change only preload inputs. They do not change Git's configured `~/.ssh/id_ed25519.pub` signing key.

### Identity Selection

The common configuration sets the personal email and `user.useConfigOnly=true`. On work machines, conditional includes override it with the work email for GitLab `gitlab.solaredge.com`, Azure SSH `ssh.dev.azure.com`, and Azure HTTP(S) `dev.azure.com` URLs. Any remote can match; Git does not limit the check to `origin`.

Repositories that do not match, including those without remotes, use the personal email. Repository-local settings and environment overrides still take precedence. URL spelling and SSH aliases affect matching, and Git 2.36 or later is required. These identity rules do not change GPG or SSH signing methods, signing keys, or repository authentication. Conditional configuration files select identity only; they must not define remote URLs.

### Sandbox Access

The [OpenCode][opencode-profile], [Claude Code][claude-profile], and [Codex][codex-profile] profiles grant these exact files in both `filesystem.read_file` and `filesystem.bypass_protection`:

- `$HOME/.ssh/id_ed25519.pub`
- `$HOME/.ssh/known_hosts`
- `$HOME/.ssh/config`

They also allow reading the agent Git configuration and SSH transport helper, with common configuration and local signer trust available through exact-file or Git-directory grants. The launch-time socket grant permits requests to the host agent. None of these grants requires reading private key files or granting the whole `~/.ssh` directory.

### Verification

- **Local verification**: `common.gitconfig` sets `gpg.ssh.allowedSignersFile` to `~/.config/git/allowed_signers`. The [signer template][allowed-signers] includes available local `id_rsa.pub` and `id_ed25519.pub` keys for the personal email and, on work machines, the work email.
- **Git host verification**: The host checks its registered signing keys and account identity rules. A successful local verification does not register a key with the host or guarantee its verified status there.

## Failure Scenarios

| Condition | Result |
|-----------|--------|
| Socket absent, key not loaded, or preload fails | The agent still launches; a later signed commit can fail because the configured key is unavailable |
| Preload override selects a different key | Git still requests the key matching `~/.ssh/id_ed25519.pub` |
| Public metadata or socket access missing | Signing or SSH transport can fail; inspect the exact profile grants rather than exposing private keys |
| Local signer entry missing | Local verification cannot establish trust even if the host recognizes the signature |
| Signing key not registered with the host | Local verification can succeed without host verification |
| Authentication not configured | Signing can work while repository access fails |

## State Changes

- The host `ssh-agent` may gain a loaded key.
- The launched process and its children inherit the shared agent Git config selection.
- A successful signed commit stores an SSH signature in the repository.

Manual Git outside these wrappers reads [`~/.gitconfig`][manual-git], which selects the installer-configured GPG key when present.

## Dependencies

- Rendered agent/common Git configs and nono profiles, plus a host key pair and usable `ssh-agent`.
- Host signing-key registration for host verification; local `allowed_signers` for local verification.
- SSH transport additionally requires `ncat` and separate authentication access.

[domain-signing]: ../domain.md#agent-git-signing
[agent-git]: ../../../private_dot_config/git/agent.gitconfig.tmpl
[common-git]: ../../../private_dot_config/git/common.gitconfig.tmpl
[shared-launcher]: ../../../dot_agents/shell/nono.zsh.tmpl
[opencode-wrapper]: ../../../dot_agents/shell/opencode.zsh
[claude-wrapper]: ../../../dot_agents/shell/claude.zsh
[codex-wrapper]: ../../../dot_agents/shell/codex.zsh
[opencode-profile]: ../../../private_dot_config/nono/profiles/my-opencode.jsonc
[claude-profile]: ../../../private_dot_config/nono/profiles/my-claude-code.jsonc.tmpl
[codex-profile]: ../../../private_dot_config/nono/profiles/my-codex.jsonc
[git-ssh]: ../../../private_dot_local/bin/executable_nono-git-ssh
[allowed-signers]: ../../../private_dot_config/git/allowed_signers.tmpl
[manual-git]: ../../../dot_gitconfig.tmpl
