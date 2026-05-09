# RFC Architecture Review: RFC-0001 Switch To uv For Python Tooling

## Verdict

PASS

## Critical Findings

None.

## Concerns

None.

## Observations

- The main architectural boundaries remain sound: uv is a normal optional tool; no Python setup subsystem, new chezmoi data contract, standalone installer fallback, or user-state migration is introduced.
- R4 resolves the prior blocking finding about Homebrew-owned completions and `compinit -C` cache semantics. The RFC now explicitly states that stale `~/.zcompdump` may delay discovery of Homebrew-provided `_uv` and `_uvx`, and that installer-level cache invalidation is out of scope for this MVP.
- R4's completion ownership split is architecturally coherent: pip, Poetry, and pipx remain dotfiles-managed guarded cache targets, while uv/uvx completions are owned by Homebrew for the brew-only MVP.
- The pipx completion contract now covers both required executables: `pipx` and `register-python-argcomplete`.
- Current-state claims checked against source remain broadly consistent: pyenv PATH setup exists in `dot_zshenv.tmpl`, lazy pyenv initialization and pip/Poetry/pipx completion caching exist in `dot_zshrc.tmpl`, `pyenv-shell` exists as a local shim, optional tools are pre-filtered through package resolution, and `uv` is not currently in `tools.yaml` or `packagemap.yaml`.
- The active package-manager wording concern from the prior review is resolved; R4 consistently describes unsupported uv availability as non-Homebrew active package-manager behavior rather than OS-level Linux behavior.
