# Timor's Agentic Workflow

`timors-agentic-workflow` is a project-installed Spec Kit preset that adds
grounded planning and validated execution artifacts while preserving upstream
requirements and clarification commands.

- **Preset version**: `0.1.0`
- **Spec Kit compatibility**: `>=0.12.11,<0.13.0`
- **Protocol version**: `0.1.0`

The installed manifest is the compatibility authority. Every composed or
replaced command checks the installed manifest and `specify --version` before
hooks, scripts, dispatch, or writes, and fails closed outside that range.

## Provided Surfaces

| Surface | Strategy | Installed source |
|---------|----------|------------------|
| `speckit.plan` | Composed (`wrap`) | `commands/speckit.plan.md` |
| `plan-template` | Composed (`wrap`) | `templates/plan-template.md` |
| `speckit.tasks` | Replaced | `commands/speckit.tasks.md` |
| `tasks-template` | Replaced | `templates/tasks-template.md` |
| `speckit.analyze` | Replaced | `commands/speckit.analyze.md` |
| `speckit.implement` | Replaced | `commands/speckit.implement.md` |
| `execution-plan-template` | Replaced | `templates/execution-plan-template.md` |
| `analysis-template` | Replaced | `templates/analysis-template.md` |
| `review-report-template` | Replaced | `templates/review-report-template.md` |
| `progress-template` | Replaced | `templates/progress-template.md` |

`speckit.specify` and `speckit.clarify` remain upstream. The composed plan
retains upstream planning behavior and adds grounding; the task, analysis, and
implementation replacements operate on the preset protocol.

## Install In A Project

Chezmoi manages the global source package, but it does not activate the preset
for a project. From the project root, initialize Spec Kit if needed, then
install the preset explicitly:

```sh
specify init --here --integration opencode --ignore-agent-tools
specify preset add --dev "$HOME/.config/specify/presets/timors-agentic-workflow"
```

The direct `specify preset add --dev` command copies the source into the
project's `.specify/presets/timors-agentic-workflow/` and generates the active
integration's command files. The copy is independent of the global source;
global changes take effect in a project only after an explicit refresh.

Spec Kit `0.12.11` has no preset-update command. Refresh a development install
by removing and adding it again, never by assuming an in-place update:

```sh
specify preset remove timors-agentic-workflow
specify preset add --dev "$HOME/.config/specify/presets/timors-agentic-workflow"
```

## Diagnose Installation And Resolution

Run these commands from the project root:

```sh
specify preset list
specify preset info timors-agentic-workflow
specify preset resolve plan-template
specify preset resolve tasks-template
specify preset resolve execution-plan-template
specify preset resolve analysis-template
specify preset resolve review-report-template
specify preset resolve progress-template
```

All six listed templates should resolve through `timors-agentic-workflow`.
Spec Kit `0.12.11` `preset resolve` is template-only, so command names report
as not found rather than identifying their source. Verify the four installed
command surfaces by checking the generated OpenCode files for project-relative
`.specify/presets/timors-agentic-workflow/` references; those files must never
contain the global source path. `speckit.specify` and `speckit.clarify` remain
upstream and outside this preset.

## Artifact And State Ownership

The project tracks the installed preset copy, generated integration commands,
normal Spec Kit design artifacts, `plan.md`, `tasks.md`,
`execution-plan.md`, `reviews/*.md`, and `analysis.md`. Only
`specs/**/progress.md` is ignored: it is local runtime state and can contain
binding, workspace, checkpoint, failure, and test evidence.

`tasks.md` owns atomic task records and checkboxes. `execution-plan.md` owns
execution groups and orchestration. Role reports own cumulative findings, and
`analysis.md` owns aggregate review state and human approval. `progress.md`
owns local execution state.

The chezmoi source package is the editable distribution source. The installed
project copy is an informational, refresh-replaceable copy used by command
resolution. Commands read installed preset files but never mutate the installed
copy or its README.

## Adoption And Lifecycle

- A project with only `spec.md` follows the normal enhanced flow: run plan,
  then tasks, analyze, and implement as applicable.
- A feature with `plan.md` and no `tasks.md` reruns `/speckit.plan` before
  task generation.
- A feature with upstream `tasks.md` and no conforming `execution-plan.md`
  reruns `/speckit.tasks`.
- A malformed, unsupported, missing, or nonconforming task/execution-plan pair
  reruns `/speckit.tasks`; this release defines no migration, repair, or
  inferred protocol data.
- A conforming task/execution-plan pair with missing or nonconforming role
  reports or `analysis.md` reruns `/speckit.analyze`. Replace an unparseable
  report only after explicit human confirmation so cumulative findings and an
  existing approval decision are not silently discarded.
- A feature already executing through upstream behavior or another preset must
  finish or be explicitly abandoned before adoption. Mixed execution semantics
  are rejected.

Removing the preset is safe before enhanced artifacts exist. For planned work
that has not started execution, revert the enhanced artifacts or retain them as
inert history; do not treat them as core implementation plans. For in-flight
work, finish with the installed preset, pause with the pinned preset and local
state retained, or explicitly abandon back to the recorded execution base.
Never automatically delete branches, worktrees, checkpoint refs, dirty output,
progress, or feature artifacts.

For a CLI upgrade, use the preset-first sequence: review upstream
command/template changes, release a preset whose manifest and preflight cover
the new compatibility range, refresh the project-installed preset while using
a compatible CLI, verify list/info/resolve results, and only then upgrade or use
the new CLI. If the CLI is upgraded first, installed enhanced commands fail
closed until a compatible preset is installed.

## No-Probe Semantics And Limitations

The preset defines data contracts, not worker provisioning, provider/model
bindings, or runtime-specific behavior. Binding discovery and checks use runtime
metadata and configuration; no suitability probe, calibration task, canary,
sacrificial invocation, or model call occurs before real assigned work. The
first dispatch is the real review or execution assignment.

Installation and resolution validation prove package compatibility, not runtime
operation. Operational support requires a normal feature to complete the
required real-work handshake. The preset does not migrate existing artifacts,
does not automatically refresh project copies, and cannot safely be removed
during active execution.

## Operational Status

> This is a non-normative installation observation, not a runtime-support or
> binding guarantee. It contains no provider, model, or sensitive binding data.

- **Installation**: Verified by temporary source-to-target materialization,
  project installation, resolution, refresh, and tracking checks.
- **Runtime**: Not yet exercised.
- **Preset version**: `0.1.0`
- **Spec Kit version**: `0.12.11`
- **OpenCode version**: `1.17.18`
- **Binding configuration class**: Not exercised.
- **Validation date**: `2026-07-12`

Installed README copies are informational and refresh-replaceable. Commands
never mutate them.
