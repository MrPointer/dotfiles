---
description: Generate an actionable, dependency-ordered tasks.md for the feature based on available design artifacts.
handoffs:
  - label: Analyze For Consistency
    agent: speckit.analyze
    prompt: Run a project analysis for consistency
    send: true
  - label: Implement Project
    agent: speckit.implement
    prompt: Start the implementation in phases
    send: true
scripts:
  sh: scripts/bash/setup-tasks.sh --json
  ps: scripts/powershell/setup-tasks.ps1 -Json
---

## Preset Compatibility Preflight

Before extension hooks, prerequisite scripts, or any write, read
`.specify/presets/timors-agentic-workflow/preset.yml` and obtain the active
Spec Kit version with `specify --version`. Parse the manifest's
`requires.speckit_version` range and fail closed unless the active version is
within `>=0.12.11,<0.13.0`. If the manifest, version, or version range cannot
be read and parsed, stop and report the compatibility failure. Do not invoke a
hook, setup script, or write a project artifact after a failed preflight.

## User Input

```text
$ARGUMENTS
```

You **MUST** consider the user input before proceeding (if not empty).

## Pre-Execution Checks

**Check for extension hooks (before tasks generation)**:
- Check if `.specify/extensions.yml` exists in the project root.
- If it exists, read it and look for entries under the `hooks.before_tasks` key
- If the YAML cannot be parsed or is invalid, skip hook checking silently and continue normally
- Filter out hooks where `enabled` is explicitly `false`. Treat hooks without an `enabled` field as enabled by default.
- For each remaining hook, do **not** attempt to interpret or evaluate hook `condition` expressions:
  - If the hook has no `condition` field, or it is null/empty, treat the hook as executable
  - If the hook defines a non-empty `condition`, skip the hook and leave condition evaluation to the HookExecutor implementation
- For each executable hook, output the following based on its `optional` flag:
  - **Optional hook** (`optional: true`):
    ```
    ## Extension Hooks

    **Optional Pre-Hook**: {extension}
    Command: `/{command}`
    Description: {description}

    Prompt: {prompt}
    To execute: `/{command}`
    ```
  - **Mandatory hook** (`optional: false`):
    ```
    ## Extension Hooks

    **Automatic Pre-Hook**: {extension}
    Executing: `/{command}`
    EXECUTE_COMMAND: {command}

    Wait for the result of the hook command before proceeding to the Outline.
    ```
    After emitting the block above you MUST actually invoke the hook and wait for it to finish before continuing. Run it the same way you would run the command yourself in this agent/session (the invocation may differ from the literal `{command}` id shown above, e.g. a skills-mode agent runs it as `/skill:speckit-...` or `$speckit-...`). Emitting the block alone does not run the hook.
- If no hooks are registered or `.specify/extensions.yml` does not exist, skip silently

## Outline

1. **Setup**: Run `{SCRIPT}` from repo root and parse FEATURE_DIR, TASKS_TEMPLATE, and AVAILABLE_DOCS list. `FEATURE_DIR` and `TASKS_TEMPLATE` must be absolute paths when provided. `AVAILABLE_DOCS` is a list of document names/relative paths available under `FEATURE_DIR` (for example `research.md` or `contracts/`). For single quotes in args like "I'm Groot", use escape syntax: e.g 'I'\''m Groot' (or double-quote if possible: "I'm Groot"). Do not write an artifact until setup succeeds.

2. **Load inputs progressively**: Read required `plan.md` and `spec.md`, then
   available `data-model.md`, `contracts/`, `research.md`, and `quickstart.md`
   under FEATURE_DIR. If `/memory/constitution.md` exists, load it for project
   principles and governance constraints. Use only available artifacts.

3. **Load preset policy progressively** from the installed preset:
   - Read `.specify/presets/timors-agentic-workflow/references/decomposition.md`
     to derive tasks, groups, contracts, and data flows.
   - Read `.specify/presets/timors-agentic-workflow/references/testable-work.md`
     when assigning test expectations.
   - Read `.specify/presets/timors-agentic-workflow/references/concurrency-policy.md`
     when deriving the graph and order.
   - Read `.specify/presets/timors-agentic-workflow/references/model-and-worker-selection.md`
     when selecting semantic model tiers, capabilities, skills, or project
     review roles.
   - Read `.specify/presets/timors-agentic-workflow/references/documentation-planning.md`
     when planning documentation.
   - Read `.specify/presets/timors-agentic-workflow/references/artifact-validation.md`
     before rendering or validating the artifact pair. This is the
     authoritative protocol contract.

4. **Generate the canonical pair together**:
   - Read TASKS_TEMPLATE and use it for `tasks.md`; if TASKS_TEMPLATE is empty,
     use `.specify/templates/tasks-template.md`.
   - Use `.specify/presets/timors-agentic-workflow/templates/execution-plan-template.md`
     for `execution-plan.md`.
   - Organize tasks by setup, isolated foundation, priority-ordered user
     stories, and final cross-cutting work. Preserve story goals, independent
     tests, checkpoints, strict checklist records, and exact file paths.
   - Generate tests by default for testable behavior. When no new test task is
     practical, record concrete existing-coverage evidence or a concrete
     not-applicable reason as required by the templates and policy.
   - Keep task descriptions and task-local `[P]` judgments exclusively in
     `tasks.md`. Put group ownership, the global DAG, models, capabilities,
     contracts, data flow, and execution policy exclusively in
     `execution-plan.md`.
   - Ask the user before selecting or omitting a project review role when its
     applicability remains ambiguous after inspecting repository packets and
     change triggers.

5. **Validate before writing**: Validate the complete in-memory pair against
   `.specify/presets/timors-agentic-workflow/references/artifact-validation.md`,
   including exact identity, headings, tables, vocabularies, identifiers, task
   coverage, graph/order agreement, file ownership, role packets, contract
   joins, and identical producer and consumer shapes. Unknown fields are
   invalid. Write both files only after the pair passes; do not leave one newly
   generated artifact without the other.

6. **Adopt or regenerate**: Existing upstream `tasks.md` without a conforming
   `execution-plan.md`, and any malformed or unsupported pair, must be
   regenerated by this command. Do not migrate, repair, normalize, or infer
   missing protocol data.

## Mandatory Post-Execution Hooks

**You MUST complete this section before reporting completion to the user.**

Check if `.specify/extensions.yml` exists in the project root.
- If it does not exist, or no hooks are registered under `hooks.after_tasks`, skip to the Completion Report.
- If it exists, read it and look for entries under the `hooks.after_tasks` key.
- If the YAML cannot be parsed or is invalid, skip hook checking silently and continue to the Completion Report.
- Filter out hooks where `enabled` is explicitly `false`. Treat hooks without an `enabled` field as enabled by default.
- For each remaining hook, do **not** attempt to interpret or evaluate hook `condition` expressions:
  - If the hook has no `condition` field, or it is null/empty, treat the hook as executable
  - If the hook defines a non-empty `condition`, skip the hook and leave condition evaluation to the HookExecutor implementation
- For each executable hook, output the following based on its `optional` flag:
  - **Mandatory hook** (`optional: false`) — **You MUST emit `EXECUTE_COMMAND:` for each mandatory hook**:
    ```
    ## Extension Hooks

    **Automatic Hook**: {extension}
    Executing: `/{command}`
    EXECUTE_COMMAND: {command}
    ```
    After emitting the block above you MUST actually invoke the hook and wait for it to finish before continuing. Run it the same way you would run the command yourself in this agent/session (the invocation may differ from the literal `{command}` id shown above, e.g. a skills-mode agent runs it as `/skill:speckit-...` or `$speckit-...`). Emitting the block alone does not run the hook.
  - **Optional hook** (`optional: true`):
    ```
    ## Extension Hooks

    **Optional Hook**: {extension}
    Command: `/{command}`
    Description: {description}

    Prompt: {prompt}
    To execute: `/{command}`
    ```

## Completion Report

Report both generated paths, total tasks, task count and independent test per
story, checkpoints, task-local parallel opportunities, execution groups and
order, concurrency decision, model-tier rationale, selected project review
roles, suggested MVP scope, and pair-validation result. Include task-local
parallel execution examples per story without treating them as group
concurrency. Confirm every task has a checkbox, sequential `TNNN`, applicable
story label, and exact path.

Context for task generation: {ARGS}

The artifact pair must be immediately usable by cold readers without duplicated
task descriptions or undeclared orchestration.

## Task Generation Rules

- Organize work by user story. Isolate shared foundation and final
  cross-cutting work only where they cannot belong to one independently
  testable story.
- Use `- [ ] TNNN [P?] [Story?] Description with exact path`. IDs begin at
  `T001` and increase without gaps. Story labels are required only in story
  phases. `[P]` means only that the individual task uses different files and
  has no dependency on an incomplete task; it never declares group concurrency.
- Keep tests and the behavior they prove in the same execution group. Tests are
  the default for testable behavior; otherwise record the concrete exception.
- `tasks.md` contains no execution-group IDs, global dependency graph, model or
  worker assignment, group concurrency, or cross-group contracts.

## Done When

- [ ] `tasks.md` and `execution-plan.md` were generated together and pass the
  installed artifact-validation contract.
- [ ] Extension hooks were dispatched or skipped according to Mandatory
  Post-Execution Hooks.
- [ ] Completion reports task coverage, story tests, execution decomposition,
  and protocol validation.
