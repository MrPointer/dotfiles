# RFC-0001: Switch To uv For Python Tooling

Status: Accepted
Revision: R4
Last Updated: 2026-05-08

## Review Record

| Reviewer | Scope | Status | Notes |
|----------|-------|--------|-------|
| `rfc-architect-reviewer` | Architecture, boundaries, contracts, current-state fit, planning readiness | Passed | `docs/rfcs/reviews/RFC-0001.rfc-architect-reviewer.md`; R4 passed with no concerns. |
| `rfc-risk-reviewer` | Technical risks, migration, compatibility, rollback, hidden complexity | Passed | `docs/rfcs/reviews/RFC-0001.rfc-risk-reviewer.md`; R4 passed with no concerns. |
| `rfc-clarity-reviewer` | Clarity and actionability | Passed | `docs/rfcs/reviews/RFC-0001.rfc-clarity-reviewer.md`; R4 passed with no concerns. |

## Summary

This RFC changes the dotfiles' Python tooling posture from pyenv-centered shell
integration to uv-compatible, on-demand Python tooling. The repo does not
contain Python project dependencies or a Python-based installer path today, so
the chosen design is intentionally small: offer `uv` as a normal optional
installer tool, remove pyenv-specific shell startup behavior, and keep low-cost
completions for Python tools that are still useful in existing projects.

The design does not introduce a dedicated Python setup subsystem. It does not
preinstall Python versions, migrate existing `pyenv` or `pipx` state, convert
Poetry projects, or add standalone uv installation fallback for non-Homebrew
systems. The immediate goal is to make uv available where the current optional
tools mechanism supports it, simplify shell startup, and leave future Python
needs to uv's on-demand workflows when uv is installed. No generated shell
behavior may require uv to be present.

## Problem

The current shell configuration contains pyenv-specific behavior even though
pyenv's main benefit in this setup, automatic directory-sensitive switching for
bare `python` and `pip`, is not part of the user's actual workflow. This creates
shell complexity without corresponding value:

- `.zshenv` adds pyenv shims and pyenv's binary directory when `~/.pyenv`
  exists.
- `.zshrc` lazily initializes pyenv and pyenv-virtualenv based on directory
  markers or Python-related commands.
- A local `pyenv-shell` shim exists to work around pyenv shell-integration
  behavior.

At the same time, uv has become the preferred tool the user wants available for
Python environments, Python-based tools, and ad hoc Python work. The dotfiles
should shift toward uv-compatible shell behavior without over-engineering a full
Python lifecycle manager inside the installer or making shell startup depend on
uv availability.

## Goals

- Make `uv` available through the existing installer optional-tools flow where
  the active package manager supports it.
- Remove pyenv-specific shell startup behavior from the dotfiles.
- Use Homebrew-provided uv/uvx completions when uv is installed by Homebrew.
- Preserve `pip`, `poetry`, and guarded `pipx` completions because they remain
  useful for core Python workflows, existing projects, and machines where pipx
  is still installed.
- Keep shell startup uv-compatible rather than uv-dependent.
- Keep the migration small, reversible, and aligned with the existing optional
  tools architecture.

## Non-Goals

- Do not add a dedicated Python setup subsystem to the installer.
- Do not preinstall Python versions with `uv python install`.
- Do not migrate, delete, uninstall, or convert existing `pyenv`, `pipx`, or
  Poetry state.
- Do not convert existing Poetry projects to uv projects.
- Do not preserve pyenv-style automatic directory-sensitive switching for bare
  `python` or `pip`.
- Do not add optional-tool preselection, recommendation tiers, or special
  prompt behavior for `uv`.
- Do not add standalone uv installer fallback for non-Homebrew active package
  managers in this RFC.

## Constraints

- This repository is the chezmoi source directory; shell behavior changes belong
  in source templates such as `dot_zshenv.tmpl` and `dot_zshrc.tmpl`.
- The installer writes chezmoi data, but Python tooling is not part of the
  current chezmoi data contract.
- Shell startup should stay fast and progressively skip unavailable tools.
- The existing optional-tools flow is package-manager-backed and non-fatal.
- Existing brew-only optional tools are silently filtered out when the active
  package manager is not Homebrew.
- The installer is used for initial bootstrapping, not ongoing state keeping.
- The user has explicitly accepted Homebrew-only uv availability for the MVP
  because most target setups use Homebrew.
- The user has explicitly confirmed that directory-sensitive bare `python`
  switching is not important in day-to-day use.

## Current State

The project has three relevant boundaries:

- The Go installer bootstraps prerequisites, applies dotfiles, and optionally
  installs selected CLI tools.
- Chezmoi templates define shell runtime behavior through files such as
  `dot_zshenv.tmpl` and `dot_zshrc.tmpl`.
- The generated shell files configure runtime PATH, completions, Homebrew
  loading, work environment loading, and tool integrations.

Python support is currently shell-centric:

- `dot_zshenv.tmpl` adds `~/.pyenv/shims` and `~/.pyenv/bin` to `PATH` when
  `~/.pyenv` exists.
- `dot_zshrc.tmpl` lazily initializes pyenv and pyenv-virtualenv when a Python
  marker or Python-related command is detected.
- `dot_zshrc.tmpl` caches completions for `poetry`, `pip`, and `pipx`.
- `private_dot_local/bin/executable_pyenv-shell` exists only to make `pyenv
  shell` behave when full shell integration is not loaded.

The optional tools system is package-manager-centered:

- `installer/internal/config/tools.yaml` lists optional tool names and
  descriptions.
- `installer/internal/config/packagemap.yaml` maps abstract tool names to
  package-manager-specific names.
- During installation, `installOptionalTools` pre-filters tools by resolving
  them for the active package manager.
- Tools without a valid mapping for the active package manager are dropped
  before the interactive prompt and before `--install-tools` installation.
- Existing tools such as `sheldon`, `eza`, and `difftastic` are brew-only and do
  not appear when the active package manager is not Homebrew.

The repository currently has no Python project payload to preserve. Searches
found no `pyproject.toml`, `requirements*.txt`, `poetry.lock`, or `uv.lock` in
the repo. Python references are limited to shell integration, documentation,
examples, and package-resolver test fixtures.

## Chosen Approach

Adopt uv as the preferred Python tooling entrypoint while keeping the change
inside existing boundaries:

- Add `uv` to `tools.yaml` as a normal optional tool.
- Add a Homebrew mapping for `uv` in `packagemap.yaml`.
- Let non-Homebrew active package managers continue using existing unsupported-
  tool behavior: `uv` is filtered out when no mapping exists.
- Remove pyenv-specific PATH setup, lazy initialization, and local shim support.
- Keep cached completions for `pip`, `poetry`, and `pipx`.
- Rely on Homebrew-provided completions for `uv` and `uvx` in the MVP.
- Document the shell startup process as uv-compatible rather than pyenv-oriented.

This design treats uv like other optional tools for installation while treating
pyenv as no longer part of the dotfiles' preferred shell runtime model.

## Decision Summary

| Decision | Rationale | Consequence |
|----------|-----------|-------------|
| Use uv as the preferred Python tooling entrypoint when available. | uv covers on-demand Python versions, project environments, and Python CLI tools without pyenv shell shims. | Python work should route through uv workflows when project-specific behavior is needed, but shell startup must not require uv. |
| Do not preserve automatic bare `python` switching. | Settled user decision: the user does not rely on this behavior. | pyenv shell integration can be removed without replacing its directory-sensitive shim model. |
| Install uv through optional tools. | The repo has no Python payload requiring a stronger bootstrap dependency. | uv remains user-selected and non-fatal like other optional tools. |
| Support Homebrew first. | The existing package-map flow already supports brew-only tools, and settled user decision accepts brew-only support for the MVP. | uv is hidden for non-Homebrew active package managers until explicit mappings or fallback install support are added. |
| Keep `pip`, `poetry`, and `pipx` completion caching. | They are low-cost and useful for existing Python projects or machines where pipx remains installed. | Removing pyenv and adding uv does not remove useful legacy completions. |
| Rely on Homebrew for uv/uvx completions. | The Homebrew uv formula generates uv and uvx completions, and `.zshrc` already adds Homebrew's zsh completion directory to `fpath`. | No explicit uv/uvx completion cache generation is needed for the brew-only MVP. |
| Do not migrate user state. | The installer is not an ongoing state keeper. | Existing pyenv/pipx/Poetry files remain untouched. |

## Proposed Architecture

The proposed architecture keeps Python tooling split across the same lifecycle
boundaries as the rest of the dotfiles:

- Installer lifecycle: `uv` is offered during optional tools installation when
  the active package manager can resolve it.
- Template lifecycle: shell templates stop adding pyenv runtime behavior and
  retain guarded completion cache generation for pip/Poetry/pipx. uv/uvx
  completions come from Homebrew's zsh completion directory when uv is installed
  by Homebrew.
- Runtime lifecycle: shells no longer initialize pyenv, add pyenv shims, or
  depend on a pyenv workaround shim. Shell startup remains valid without uv;
  users invoke uv directly when it is installed and Python project or tool
  behavior is needed.

No new chezmoi data keys, installer state files, or persistent migration records
are introduced.

## Components And Responsibilities

| Component / Boundary | Responsibility | Current / New / Modified | Notes |
|----------------------|----------------|--------------------------|-------|
| `installer/internal/config/tools.yaml` | Lists optional tools shown by the installer after package-manager filtering. | Modified | Add `uv` with a user-facing description. |
| `installer/internal/config/packagemap.yaml` | Maps abstract optional tool names to package-manager package names. | Modified | Add `uv` for Homebrew only in the MVP. |
| `dot_zshenv.tmpl` | Defines fast environment setup for all Zsh sessions. | Modified | Remove pyenv PATH/shim setup. |
| `dot_zshrc.tmpl` | Defines interactive shell behavior and completion caching. | Modified | Remove lazy pyenv initialization; keep pip/Poetry/pipx completion caching; keep Homebrew completion loading for brew-installed uv/uvx. |
| `private_dot_local/bin/executable_pyenv-shell` | Provides a pyenv-specific workaround command. | Removed | No longer needed once pyenv shell integration is removed. |
| `docs/processes/shell-startup.md` | Documents shell startup behavior. | Modified | Update pyenv-oriented current-state documentation to the uv-compatible runtime model. |
| `docs/domain.md` | Describes optional tools and domain concepts. | Modified | Keep optional-tool documentation accurate after adding uv to the configured optional tools. |
| `docs/architecture.md` | Describes shell runtime dependencies. | Modified | Stop naming pyenv as a shell runtime dependency; describe Python tooling as optional and guarded. |

## Contracts And Interfaces

The optional tools contract remains unchanged:

- Tool definitions remain entries in `tools.yaml`.
- Package resolution remains driven by `packagemap.yaml`.
- Unsupported tools remain filtered out before selection.
- `--install-tools` continues to install all currently resolvable tools.
- Interactive selection continues to show all currently resolvable tools
  unselected by default.

The shell runtime contract changes by removal rather than replacement:

- The dotfiles no longer guarantee that pyenv shims are placed before other
  Python executables.
- The dotfiles no longer guarantee that `~/.pyenv/bin` is placed on `PATH`, so
  explicit `pyenv` commands are available only if the user manages pyenv outside
  these dotfiles.
- The dotfiles no longer initialize pyenv or pyenv-virtualenv in interactive
  shells, including for `.python-version`, `.envrc`, or explicit `pyenv*`
  command detection.
- The dotfiles no longer provide a `pyenv-shell` workaround command.
- The dotfiles retain guarded completion cache hooks for `pip`, `poetry`, and
  `pipx`.
- The dotfiles do not generate cached uv/uvx completions in the MVP. Brew-
  installed uv/uvx completions are provided through Homebrew's zsh completion
  directory, which is already added to `fpath` before `compinit`.
- Newly installed Homebrew completions may not be visible immediately if
  `compinit -C` reuses an existing `~/.zcompdump`; the MVP does not add installer-
  level zcompdump invalidation.
- pipx remains completion-supported when installed, but it is no longer the
  preferred Python tool manager in the dotfiles' documentation or direction.

The completion contract is intentionally concrete:

- `pip` and `poetry` completion generation remains guarded by command existence
  checks and stale cache checks, matching the existing completion-cache pattern.
- `pipx` completion generation is guarded by both `pipx` and
  `register-python-argcomplete` availability, plus the same stale cache check.
- uv/uvx completion behavior follows the Homebrew formula for the brew-only MVP.
- If future work adds non-Homebrew uv installation, that work must decide whether
  to generate uv/uvx completions explicitly or rely on that install method.

## Data And State

No new durable data is introduced.

Existing user-owned state is intentionally left untouched:

- `~/.pyenv` may continue to exist but is no longer managed or placed on PATH by
  the dotfiles.
- Existing pipx tool environments may continue to exist and can still receive
  guarded completion refreshes when both the `pipx` command and
  `register-python-argcomplete` helper are available.
- Existing completion cache files for tools that are no longer installed may
  remain; they are residual user state and are not authoritative for the current
  dotfiles source.
- Existing Poetry projects and Poetry user configuration are not modified.
- uv-managed Python versions, caches, and tools are created later by uv itself
  when the user invokes uv workflows.

## Control And Data Flow

During installer optional tools setup:

1. The installer loads `tools.yaml`.
2. The package resolver tries to resolve each tool for the active package
   manager.
3. `uv` resolves on Homebrew systems through the new `packagemap.yaml` entry.
4. `uv` does not resolve when the active package manager is not Homebrew in the
   MVP and is filtered out like existing brew-only tools.
5. Interactive users may select `uv` if it is visible; `--install-tools` installs
   it automatically only when it is resolvable.

During shell startup:

1. `.zshenv` performs base PATH, Homebrew, Rust, and work-environment setup
   without pyenv-specific PATH mutation.
2. `.zshrc` initializes completions and interactive shell integrations without
   pyenv lazy hooks.
3. Completion cache generation creates or refreshes `pip` and `poetry`
   completion files only when the corresponding command is available and the
   cache is missing or stale.
4. `pipx` completion cache generation also requires
   `register-python-argcomplete` to be available before refreshing `_pipx`.
5. Brew-installed uv/uvx completions are discovered through Homebrew's zsh
   completion directory on `fpath`; `.zshrc` does not generate uv/uvx completion
   cache files in the MVP.
6. If `~/.zcompdump` is stale, Homebrew-provided uv/uvx completions may not be
   discovered until the completion cache is refreshed. The existing interactive
   `brew` shell wrapper handles future interactive brew changes, but installer-
   driven package installs do not gain a new cache invalidation mechanism in this
   RFC.
7. Python project behavior is handled by explicit uv commands or by whatever
   Python environment the user activates outside the dotfiles.

## Failure Modes And Recovery

| Failure Mode | Expected Behavior | Recovery / User Impact |
|--------------|-------------------|------------------------|
| `uv` has no mapping for the active package manager. | It is filtered out before prompt/installation. | User can install uv manually or wait for future standalone fallback support. |
| Homebrew installation of `uv` fails. | Optional tool installation records the failure and continues with other tools. | Main installation remains successful; user can retry or install manually. |
| `uv` is not installed at shell startup. | No uv/uvx completions are available from Homebrew. | Shell starts normally without uv completions. |
| Homebrew installs uv but zsh completion cache is stale. | `_uv` and `_uvx` files may exist in Homebrew's completion directory but not be discovered by `compinit -C`. | Completions appear after the zsh completion cache is refreshed; this RFC does not add installer-level cache invalidation. |
| Existing project expects pyenv shims, `pyenv` on PATH, or pyenv initialization from `.envrc`/`.python-version`. | Dotfiles no longer provide those shims, PATH entries, initialization hooks, or `pyenv-shell` workaround behavior. | User can invoke uv workflows, activate a project environment, or manage pyenv PATH/init outside these dotfiles. |
| Existing pipx-installed command remains on PATH from prior user state. | Dotfiles do not remove it and continue guarded completion refreshes when pipx and `register-python-argcomplete` are available. | Command may still work if user state provides it; pipx remains completion-supported but is not the preferred tool manager. |

## Security, Privacy, And Permissions

This RFC does not add a new download script, credential surface, or privilege
model. uv installation uses the same package-manager path as other optional
tools in the MVP. Future standalone installer support would require a separate
security review because it would download and execute an external installer
script.

## Operations And Observability

The installer observability model remains unchanged:

- Unsupported optional tools are silently filtered before prompt/installation.
- Installation failures for selected tools are reported in the optional-tools
  summary.
- Shell completion generation remains silent and guarded by command existence
  checks.

No new metrics, logs, or state files are introduced.

## Compatibility And Migration

This is a forward-looking shell behavior migration, not a user-state migration.

Compatibility expectations:

- Existing machines with pyenv installed are not modified on disk.
- Existing shells generated after this change will stop adding pyenv behavior.
- Existing projects or scripts that call `pyenv` directly must manage pyenv PATH
  and initialization outside these dotfiles.
- Existing Poetry and pip workflows remain supported at the completion level.
- Existing pipx workflows are not actively broken; pipx remains
  completion-supported when installed and when `register-python-argcomplete` is
  available, but is no longer the preferred Python tool manager.
- Existing cached completions for removed tools or pyenv-related target files
  may remain as residual state; this RFC does not require cleanup.
- Non-Homebrew systems continue to behave like they do for existing brew-only
  optional tools: unsupported tools are omitted from selection.

Rollback is source-level: reintroducing pyenv shell blocks and
`private_dot_local/bin/executable_pyenv-shell` in the chezmoi source, then
applying dotfiles, restores the previous shell integration model assuming user
pyenv state still exists.

## Alternatives Considered

| Alternative | Strengths | Why Rejected |
|-------------|-----------|--------------|
| Keep pyenv and add uv alongside it. | Maximum compatibility for bare `python` directory switching. | Preserves shell complexity for behavior the user does not use. |
| Create a dedicated Python setup subsystem. | Could install uv, preinstall Python versions, and define a complete Python baseline. | Over-engineered for a repo with no Python project payload and no installer-time Python dependency. |
| Add standalone uv installer fallback now. | Would support non-Homebrew active package managers even without native package mappings. | More installer design work than the immediate need justifies; most target setups use Homebrew. |
| Add default-selected/recommended optional tools. | Could nudge users toward uv during setup. | Adds UX and installer complexity for a tool that should remain optional. |
| Remove all Python completions. | Simplest shell cleanup. | `pip` remains core, Poetry remains common in existing projects, and installed pipx still benefits from accurate completions. |
| Migrate existing pyenv/pipx/Poetry state. | Could produce a cleaner local machine after migration. | Outside installer responsibility and risks changing user-owned state unexpectedly. |

## Risks And Tradeoffs

| Risk / Tradeoff | Impact | Mitigation Or Acceptance |
|-----------------|--------|--------------------------|
| `uv` is initially Homebrew-only in optional tools. | Non-Homebrew active package managers will not offer uv through the installer. | Accepted for MVP; future standalone fallback can be designed if non-Homebrew setups become important. |
| Removing pyenv shell integration may surprise projects that relied on bare `python` switching, `pyenv` on PATH, or pyenv activation from `.envrc`/`.python-version`. | Such projects may resolve a different Python or fail to call `pyenv` unless the user invokes uv, activates an environment, or manages pyenv separately. | Accepted because the user does not rely on this behavior; pyenv state is not deleted. |
| Keeping Poetry/pip/pipx completion caching preserves some legacy Python surface. | Shell config remains slightly broader than a pure uv setup. | Accepted because completion generation is guarded and useful for existing projects or machines where pipx remains installed. |
| Relying on Homebrew for uv/uvx completions ties completion availability to the brew formula. | If uv is installed outside Homebrew, uv/uvx completions may not be available from the dotfiles. | Accepted because the MVP installs uv through Homebrew only; non-Homebrew uv installation is future work. |
| Installer-driven Homebrew installs can leave `~/.zcompdump` stale. | uv/uvx completion files may not be discovered immediately after installer-run `brew install uv`. | Accepted for MVP; the RFC narrows success to Homebrew-provided completion files being available on `fpath`, not immediate compinit cache refresh. |
| No user-state cleanup leaves old pyenv/pipx files on disk. | Machines may contain unused legacy state. | Accepted because the installer is not a state keeper; cleanup can be manual or future explicit work. |

## Success Criteria

- `uv` is offered and resolvable as an optional tool on Homebrew-backed installs.
- pyenv-specific shell startup behavior is removed from generated Zsh files.
- Shell startup documentation no longer describes pyenv as part of the runtime
  model.
- `pip` and `poetry` completions are generated only when the corresponding
  command exists and the cache is missing or stale.
- `pipx` completions are generated only when both `pipx` and
  `register-python-argcomplete` exist and the cache is missing or stale.
- Brew-installed uv/uvx completion files are provided in Homebrew's zsh
  completion directory and made discoverable through `fpath`, without
  dotfiles-managed uv/uvx completion generation or installer-level zcompdump
  invalidation.
- Existing user state for pyenv, pipx, and Poetry is not modified by the
  dotfiles or installer.

## Future Work

- Non-Homebrew uv installation can be revisited later, either through a
  generalized optional-tool alternate-install mechanism or a uv-specific
  standalone fallback. That decision is intentionally outside this RFC and must
  include a completion-source decision for uv/uvx.
- Installer-level zcompdump invalidation after optional tool installation can be
  considered separately if completion freshness after package-manager installs
  becomes important.

## Planning Handoff

Planning should preserve the MVP boundary. The design is not a request for a
general Python management system, optional-tool recommendation UX, state
migration, or standalone installer support. Any implementation plan should keep
changes within the existing optional-tools mapping model, shell templates,
pyenv-specific source files, and documentation that currently describes shell
startup behavior.

If future planning proposes non-Homebrew uv installation, it should be treated
as a separate design extension because it changes installer installation methods
and security considerations.

## Source References

| Source | What It Confirms |
|--------|------------------|
| `docs/domain.md` | Optional tools are daily-use CLI tools, are not persisted in chezmoi data, and may have platform-dependent availability. |
| `docs/architecture.md` | The project separates installer bootstrap, chezmoi templates, and shell runtime behavior. |
| `docs/processes/tools-installation.md` | Optional tools are pre-filtered through package resolution, selected interactively or via `--install-tools`, installed through the active package manager, and non-fatal. |
| `docs/processes/shell-startup.md` | Current documented shell startup includes pyenv setup, lazy pyenv initialization, and completion caching for Poetry, pip, and pipx. |
| `dot_zshenv.tmpl` | Current all-shell startup adds pyenv shims and pyenv binary path when `~/.pyenv` exists. |
| `dot_zshrc.tmpl` | Current interactive startup lazily initializes pyenv and caches Poetry, pip, and pipx completions. |
| `private_dot_local/bin/executable_pyenv-shell` | Current repo includes a pyenv-specific workaround shim. |
| `installer/internal/config/tools.yaml` | Current optional tool list does not include uv. |
| `installer/internal/config/packagemap.yaml` | Current package map has existing brew-only tools and no uv mapping. |
| `installer/cmd/install.go` | Current optional-tools flow filters resolvable tools before prompt/install and installs through `ToolsInstaller`. |
| `installer/lib/toolsinstaller/installer.go` | Current `ToolsInstaller` resolves each tool and installs through the package manager, continuing after failures. |
| uv documentation: `https://docs.astral.sh/uv/concepts/python-versions/` | uv supports managed Python versions, `.python-version`, automatic downloads, and system Python discovery. |
| uv documentation: `https://docs.astral.sh/uv/concepts/tools/` | uv provides `uv tool install` and `uvx` for isolated Python CLI tools. |
| uv documentation: `https://docs.astral.sh/uv/getting-started/installation/` | uv is available through Homebrew and other install paths, including a standalone installer. |
| Homebrew uv formula fetched via `gh api repos/Homebrew/homebrew-core/contents/Formula/u/uv.rb` | Homebrew generates completions from both `uv generate-shell-completion` and `uvx --generate-shell-completion` during formula installation. |
| `dot_zshrc.tmpl` | Homebrew's zsh `site-functions` directory is added to `fpath` before `compinit`, allowing brew-installed completions to be discovered. |
| Settled user decision | Directory-sensitive bare `python` switching is not important in day-to-day use. |
| Settled user decision | A full Python setup subsystem, preinstalled Python versions, and user-state migration are out of scope. |
| Settled user decision | `pip`, `poetry`, and guarded `pipx` completions should remain. |
| Settled user decision | `uv` should not be preselected or specially recommended in the optional-tools prompt. |
| Settled user decision | Brew-only optional-tool support is acceptable for the MVP because most target setups use Homebrew. |
