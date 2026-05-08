# RFC Risk Review: RFC-0001 Switch To uv For Python Tooling

## Verdict

PASS

## Critical Findings

Empty.

## Concerns

Empty.

## Observations

- R4 resolves the prior stale `~/.zcompdump` concern at the design level. The RFC now explicitly states that `compinit -C` may reuse an existing completion dump, that installer-driven Homebrew installs do not gain a new invalidation path, and that immediate uv/uvx completion discovery is outside the MVP success boundary.
- R4 resolves the prior pipx completion guard concern. The completion contract, control flow, data/state section, compatibility section, and success criteria now require both `pipx` and `register-python-argcomplete` before refreshing `_pipx`.
- The uv/uvx completion-generation scope remains appropriately narrow: rely on Homebrew's uv formula and existing Homebrew `fpath` loading rather than adding dotfiles-managed uv/uvx cache files.
- The pyenv removal boundary remains clear: no shim PATH setup, no `~/.pyenv/bin` PATH setup, no lazy initialization for `.python-version`/`.envrc`, and no `pyenv-shell` workaround.
- The no-migration and residual-state behavior is explicit enough for planning: existing pyenv, pipx, Poetry, uv-managed state, and old completion-cache files may remain and are not treated as authoritative current dotfile state.
- Brew-only, optional uv availability is acknowledged as an accepted MVP tradeoff; the future non-Homebrew uv path is correctly treated as a separate design extension with its own completion-source and security decisions.
- Removing pyenv shell integration is a real compatibility break for projects that depended on bare `python` switching, but the RFC states this as an accepted user decision and preserves rollback at the source-template level.
