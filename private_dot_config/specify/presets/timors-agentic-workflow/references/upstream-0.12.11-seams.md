# Spec Kit 0.12.11 Seam Inventory

## Pinned Source Record

This inventory was inspected from the locally installed `specify-cli==0.12.11`
package and pinned against the Spec Kit `v0.12.11` release through `gh`.

- Release tag: `v0.12.11`
- Release commit: `e802a7dd52a6eceba9403cbbf40e60dced043238`
- Release tree: `fbbb26e1579064e55aea13759fa5cf080ea89843`
- Installed package metadata locator:
  `specify_cli-0.12.11.dist-info/METADATA` (Version: `0.12.11`)
- Installed core-pack root locator:
  `specify_cli/core_pack/`
- Upstream source root locator at the pinned tag: `templates/`

The source locators below pair the version-pinned upstream path and Git blob
with the package-accessible core-pack path used for the inspection. The plan
command content SHA-256 was verified equal between the local package and the
pinned upstream blob:
`074f8c273868ef69881708f89d764a1fb05fca272fe96ffa770d0626af3714e0`.

## Shared Extension-Hook Seam

Every listed command has a pre-execution hook section and all except analyze
place their post-execution hook section before the completion report. The exact
event key changes by command, but the hook behavior is the same:

1. Read `.specify/extensions.yml`; silently skip invalid YAML or absent event.
2. Exclude only hooks with `enabled: false`; unspecified `enabled` means true.
3. Treat absent, null, or empty `condition` as executable; skip a nonempty
   condition for HookExecutor evaluation instead of interpreting it.
4. Optional hooks emit the `Optional Pre-Hook` or `Optional Hook` block with
   `{extension}`, `{command}`, `{description}`, and `{prompt}`.
5. Mandatory hooks emit the `Automatic Pre-Hook` or `Automatic Hook` block with
   `EXECUTE_COMMAND: {command}`, invoke the hook, and wait for completion.

The composed `speckit.plan` adds its compatibility preflight before this shared
pre-hook seam; it does not remove or alter the upstream hook behavior.

## `speckit.plan`

- Upstream locator:
  `templates/commands/plan.md` at `v0.12.11`, blob
  `312c5ab1812913eecd4dce85110825a144edb0e8`.
- Package locator: `specify_cli/core_pack/commands/plan.md`.
- Frontmatter fields:
  - `description: Execute the implementation planning workflow using the plan template to generate design artifacts.`
  - `handoffs[0]`: `label: Create Tasks`, `agent: speckit.tasks`,
    `prompt: Break the plan into tasks`, `send: true`.
  - `handoffs[1]`: `label: Create Checklist`, `agent: speckit.checklist`,
    `prompt: Create a checklist for the following domain...`.
  - `scripts.sh: scripts/bash/setup-plan.sh --json`.
  - `scripts.ps: scripts/powershell/setup-plan.ps1 -Json`.
- Integration placeholders: `$ARGUMENTS` and `{SCRIPT}`. There are no
  `__SPECKIT_COMMAND_*__` placeholders in this command.
- Prerequisite script seam: `{SCRIPT}` runs the listed setup-plan script from
  repository root and consumes JSON `FEATURE_SPEC`, `IMPL_PLAN`, `SPECS_DIR`,
  and `BRANCH`.
- Hook keys: `hooks.before_plan` and `hooks.after_plan`.
- Preserved workflow seams: setup; feature-spec and constitution loading;
  Technical Context; Constitution Check before and after design; Phase 0
  `research.md`; Phase 1 `data-model.md`, `contracts/`, and `quickstart.md`;
  completion report; and both extension-hook sections.

## `speckit.tasks`

- Upstream locator:
  `templates/commands/tasks.md` at `v0.12.11`, blob
  `ae7192c3d3528ab0a4685a05110973b45a8e20fe`.
- Package locator: `specify_cli/core_pack/commands/tasks.md`.
- Frontmatter fields:
  - `description: Generate an actionable, dependency-ordered tasks.md for the feature based on available design artifacts.`
  - `handoffs[0]`: `label: Analyze For Consistency`, `agent: speckit.analyze`,
    `prompt: Run a project analysis for consistency`, `send: true`.
  - `handoffs[1]`: `label: Implement Project`, `agent: speckit.implement`,
    `prompt: Start the implementation in phases`, `send: true`.
  - `scripts.sh: scripts/bash/setup-tasks.sh --json`.
  - `scripts.ps: scripts/powershell/setup-tasks.ps1 -Json`.
- Integration placeholders: `$ARGUMENTS`, `{SCRIPT}`, and `{ARGS}`. There are
  no `__SPECKIT_COMMAND_*__` placeholders in this command.
- Prerequisite script seam: `{SCRIPT}` runs the listed setup-tasks script from
  repository root and consumes `FEATURE_DIR`, `TASKS_TEMPLATE`, and
  `AVAILABLE_DOCS`.
- Hook keys: `hooks.before_tasks` and `hooks.after_tasks`.
- Preserved workflow seams: required `plan.md` and `spec.md`; optional
  design-artifact loading; constitution loading; user-story organization;
  strict `TNNN` checklist task format; dependency graph; parallel examples;
  completion report; and both extension-hook sections.

## `speckit.analyze`

- Upstream locator:
  `templates/commands/analyze.md` at `v0.12.11`, blob
  `2cd83bd7c031e01af1f3e5745168982d9085a3aa`.
- Package locator: `specify_cli/core_pack/commands/analyze.md`.
- Frontmatter fields:
  - `description: Perform a non-destructive cross-artifact consistency and quality analysis across spec.md, plan.md, and tasks.md after task generation.`
  - `scripts.sh: scripts/bash/check-prerequisites.sh --json --require-tasks --include-tasks`.
  - `scripts.ps: scripts/powershell/check-prerequisites.ps1 -Json -RequireTasks -IncludeTasks`.
  - `scripts.py: scripts/python/check_prerequisites.py --json --require-tasks --include-tasks`.
- Integration placeholders: `$ARGUMENTS`, `{SCRIPT}`, `{ARGS}`,
  `__SPECKIT_COMMAND_TASKS__`, `__SPECKIT_COMMAND_IMPLEMENT__`,
  `__SPECKIT_COMMAND_SPECIFY__`, and `__SPECKIT_COMMAND_PLAN__`.
- Prerequisite script seam: `{SCRIPT}` runs the listed prerequisite script from
  repository root and consumes `FEATURE_DIR` and `AVAILABLE_DOCS`, requiring
  and including `tasks.md`.
- Hook keys: `hooks.before_analyze` and `hooks.after_analyze`.
- Preserved workflow seams: read-only operation; required-artifact abort;
  progressive artifact loading; constitution authority; semantic inventories;
  detection passes; severity assignment; compact report; next actions;
  optional remediation; and both extension-hook sections.

## `speckit.implement`

- Upstream locator:
  `templates/commands/implement.md` at `v0.12.11`, blob
  `1d312a1c37850226bed8ec1899664cc02f67f0a4`.
- Package locator: `specify_cli/core_pack/commands/implement.md`.
- Frontmatter fields:
  - `description: Execute the implementation plan by processing and executing all tasks defined in tasks.md`
  - `scripts.sh: scripts/bash/check-prerequisites.sh --json --require-tasks --include-tasks`.
  - `scripts.ps: scripts/powershell/check-prerequisites.ps1 -Json -RequireTasks -IncludeTasks`.
  - `scripts.py: scripts/python/check_prerequisites.py --json --require-tasks --include-tasks`.
- Integration placeholders: `$ARGUMENTS`, `{SCRIPT}`, and
  `__SPECKIT_COMMAND_TASKS__`.
- Prerequisite script seam: `{SCRIPT}` runs the listed prerequisite script from
  repository root and consumes `FEATURE_DIR` and `AVAILABLE_DOCS`, requiring
  and including `tasks.md`.
- Hook keys: `hooks.before_implement` and `hooks.after_implement`.
- Preserved workflow seams: checklist status gate; implementation-context
  loading; ignore-file verification; task parsing; phase/dependency execution;
  TDD ordering; progress and failure handling; task check-off; completion
  validation; completion report; and both extension-hook sections.

## `plan-template`

- Upstream locator: `templates/plan-template.md` at `v0.12.11`, blob
  `36f2eab16880bac670fe43cbe7ef2b9bc8c3aa2f`.
- Package locator: `specify_cli/core_pack/templates/plan-template.md`.
- Integration placeholders preserved by the composed template:
  `__SPECKIT_COMMAND_PLAN__` and `__SPECKIT_COMMAND_TASKS__`.
- Preserved template seams: title and branch/date/spec metadata; Input line;
  Technical Context; Constitution Check; Project Structure Documentation tree;
  `plan.md`, `research.md`, `data-model.md`, `quickstart.md`, `contracts/`, and
  `tasks.md` references; Source Code options; Structure Decision; and
  Complexity Tracking.

## Composition Verification Rule

For protocol `0.1.0`, the `speckit.plan` wrapper retains the upstream
frontmatter fields and places one compatibility preflight before
`{CORE_TEMPLATE}`. The `plan-template` has one `{CORE_TEMPLATE}` and adds only
the four headings specified in `planning-grounding.md`. No upstream seam may be
removed, renamed, or duplicated.
