# uv Python Tooling Anchor

## Intent

Decide whether and how this dotfiles repo should move Python tooling from
`pip`, `pipx`, `poetry`, and possibly `pyenv` toward `uv`, with the goal of
reducing shell complexity while preserving portable machine setup behavior.

## Quick Summary

| Settled Direction | What It Means |
|-------------------|---------------|
| RFC accepted. | `docs/rfcs/RFC-switch-to-uv-python-tooling.md` R4 is the accepted design baseline. |
| Implementation complete. | `uv` is configured as a brew-only optional tool, pyenv shell runtime behavior has been removed from the dotfiles source, and docs describe the uv-compatible optional Python tooling model. |
| Directory-sensitive bare `python` selection is not required. | This removes the strongest reason to keep pyenv shims and shell integration. |
| `uv` should be a normal optional installer tool. | The installer should offer `uv`, but not preselect or special-case it; adding recommendation tiers is not worth the complexity. |
| Current repo has no Python project payload. | Searches found no `pyproject.toml`, `requirements*.txt`, `poetry.lock`, or `uv.lock`; Python references are shell integration, docs/examples, and packageresolver test fixtures. |
| Avoid a full Python setup subsystem. | `uv` should provide enough Python capability on demand; preinstalling Python versions or migrating old Python state is over-engineering for current repo priorities. |
| Keep useful Python completions. | Preserve `pip`, `poetry`, and guarded `pipx` completion caching; rely on Homebrew-provided `uv`/`uvx` completions in the MVP, with no installer-level zcompdump invalidation. |
| MVP accepts brew-only `uv` installation. | Add `uv` through the existing optional-tool/package-map flow; apt/dnf fallback can be revisited later. |
| No user-state migration. | Existing pyenv/pipx/Poetry state should be left alone; the installer is for initial bootstrapping, not ongoing state keeping. |

## Current State

- [x] Project docs and shell startup files were read for current architecture.
- [x] Current Python shell integration was identified: pyenv shims in `.zshenv`, lazy pyenv initialization in `.zshrc`, cached completions for Poetry, pip, and pipx.
- [x] `uv` documentation was checked for managed Python versions, `.python-version`, tools, auto-downloads, and executable behavior.
- [x] User confirmed directory-sensitive bare `python` selection is not important in day-to-day use.
- [x] User confirmed that if `uv` manages Python environments, the project installer should offer it.
- [x] Official `uv` installation docs were checked: package-manager installs are supported, and the standalone installer can install into `~/.local/bin` with `UV_NO_MODIFY_PATH=1` to avoid editing shell profiles.
- [x] Current repo usage was re-checked: there is no Python project dependency graph or Python-based installer path to preserve.
- [x] User rejected a complete Python setup subsystem as unnecessary right now.
- [x] User wants to keep `pip`, `poetry`, and guarded `pipx` completions, while adding `uv`/`uvx` completions and removing pyenv-specific shell behavior.
- [x] User decided `uv` should not be preselected or treated as more important than other optional tools.
- [x] Installer support was checked: standalone script installation exists as a pattern for chezmoi/Homebrew, but optional tools are currently package-manager-only.
- [x] User accepted package-manager-only `uv` availability for now, with brew-only support matching current setup reality and existing brew-only optional tools.
- [x] User confirmed no migration of existing pyenv, pipx, or Poetry user state is desired.
- [x] RFC R2 was written and passed architecture, risk, and clarity re-review after preserving guarded pipx completions.
- [x] User identified that brew-installed `uv`/`uvx` already provide completions through Homebrew; explicit uv/uvx completion caching is unnecessary for the brew-only MVP.
- [x] R3 review found completion-cache ambiguity; R4 clarifies stale zcompdump behavior and requires a safer pipx argcomplete guard.
- [x] RFC R4 passed architecture, risk, and clarity review with no concerns.
- [x] User accepted RFC R4.
- [x] Implementation plan `plans/features/switch-to-uv-python-tooling/` executed completely; final installer verification passed with `rtk go test ./...` from `installer/`.

## Constraints

- The repository is a chezmoi source directory; shell behavior changes should be made in source templates, not generated files.
- Shell startup performance matters; avoid reintroducing eager Python tool initialization or mandatory uv startup work.
- The installer is the only writer of chezmoi data, but Python tooling is currently not part of the installer data contract.
- Existing design favors progressive enhancement: missing tools should not break shell startup.

## Decisions

### Python Version Selection Model

**Decision:** The migration does not need to preserve pyenv-style automatic
directory-sensitive switching for bare `python` or `pip` commands.

**Reason:** The feature was added recently but has not been useful in practice;
removing it simplifies shell startup and reduces pyenv-specific hacks.

**Rejected:** Keep pyenv primarily to preserve automatic `cd project && python`
behavior.

**Reason rejected:** That behavior is not part of the user's actual workflow.

**Reconsider if:** A future workflow depends on invoking bare `python` or `pip`
inside project directories without routing through `uv` or an activated virtual
environment.

### Installer Ownership

**Decision:** If `uv` becomes the preferred Python environment/tool manager, the
installer should offer it as a normal optional tool, but the migration should
not create a full Python provisioning subsystem or a special recommendation
tier.

**Reason:** The dotfiles should not depend on a manually installed external tool
for the primary Python workflow. At the same time, the repo has no Python
project payload today, so preinstalling Python versions, migrating existing
pyenv/pipx state, designing a dedicated Python setup flow, or adding optional
tool preselection is unnecessary.

**Rejected:** Leave `uv` entirely to manual Homebrew/system installation while
still making dotfiles assume `uv` workflows.

**Reason rejected:** That would create an implicit dependency and weaken the
portable bootstrap goal.

**Rejected:** Add a recommended/default-selected optional-tool tier for `uv`.

**Reason rejected:** It would be extra UX and installer complexity for a tool
that should remain optional like the rest of the daily-use CLI tools.

**Reconsider if:** The repo gains Python project dependencies, the installer
starts requiring Python-managed tools during setup, or machines need a
pre-provisioned Python runtime rather than on-demand `uv` behavior.

### Shell Completions

**Decision:** Keep low-cost completion caching for relevant legacy Python tools:
retain `pip`, `poetry`, and guarded `pipx` completions. For the brew-only uv MVP,
rely on Homebrew-provided `uv` and `uvx` completions instead of generating cached
uv/uvx completions in `.zshrc`.

**Reason:** Completion caching has minimal shell-startup cost and preserves
compatibility with projects that still use Poetry, core `pip` workflows, or an
installed `pipx` command. Homebrew's uv formula already generates uv/uvx
completions, and the shell already adds Homebrew's zsh completion directory to
`fpath`.

**Decision:** Do not add installer-level `~/.zcompdump` invalidation in this RFC.

**Reason:** The MVP relies on Homebrew completion files being present and on
normal shell completion cache refresh behavior. Immediate completion freshness
after installer-run `brew install uv` is not important enough to expand scope.

**Rejected:** Invalidate zsh completion cache after optional-tool installs.

**Reason rejected:** It is a broader installer behavior change that can be
considered separately if completion freshness becomes important.

**Rejected:** Remove legacy Python completions as part of the cleanup.

**Reason rejected:** `pip` remains core Python tooling, and Poetry remains common
enough that encountering existing projects is likely. `pipx` is no longer the
preferred manager, but if it is installed, accurate completions are cheap and
useful.

**Reconsider if:** Completion generation becomes a measurable startup cost or
the tools are no longer encountered in practice.

### Optional Tool Availability

**Decision:** Add `uv` through the existing optional tool/package mapping flow,
with Homebrew support as the initial MVP. Do not add standalone installer
fallback or unsupported-tool warnings for this RFC.

**Reason:** This matches existing brew-only optional tools (`sheldon`, `eza`,
`difftastic`) and is enough for the user's current setup pattern, where most
machines use Homebrew.

**Rejected:** Extend optional tools with alternate install methods before adding
`uv`.

**Reason rejected:** It is more work than the immediate need justifies.

**Reconsider if:** Non-Homebrew Linux setups become common enough that missing
`uv` materially hurts bootstrap completeness.

### User State Migration

**Decision:** Do not migrate, clean up, uninstall, or convert existing `pyenv`,
`pipx`, or Poetry state.

**Reason:** The installer is used for initial bootstrapping, not as an ongoing
state keeper. Existing state may remain on disk without being managed by the
dotfiles.

**Rejected:** Convert pyenv Python versions, pipx-installed tools, or Poetry
projects into uv-managed equivalents.

**Reason rejected:** It would add risk and complexity outside the dotfiles'
current responsibility.

**Reconsider if:** The installer later becomes a state reconciliation tool or a
dedicated cleanup/migration command is explicitly requested.

## Open Questions

No active open questions for the accepted RFC.

## References

- `dot_zshenv.tmpl`
- `dot_zshrc.tmpl`
- `docs/processes/shell-startup.md`
- `installer/internal/config/tools.yaml`
- `installer/internal/config/packagemap.yaml`
- `docs/rfcs/RFC-switch-to-uv-python-tooling.md`
- `plans/features/switch-to-uv-python-tooling/progress.md`
- https://docs.astral.sh/uv/concepts/python-versions/
- https://docs.astral.sh/uv/concepts/tools/
- https://docs.astral.sh/uv/getting-started/installation/
- https://docs.astral.sh/uv/reference/installer/
