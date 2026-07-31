---
description: Perform a non-destructive cross-artifact consistency and quality analysis across spec.md, plan.md, and tasks.md after task generation.
scripts:
  sh: scripts/bash/check-prerequisites.sh --json --require-tasks --include-tasks
  ps: scripts/powershell/check-prerequisites.ps1 -Json -RequireTasks -IncludeTasks
  py: scripts/python/check_prerequisites.py --json --require-tasks --include-tasks
---

## Preset Compatibility Preflight

Before extension hooks, prerequisite scripts, dispatch, or any write, read
`.specify/presets/timors-agentic-workflow/preset.yml` and obtain the active
Spec Kit version with `specify --version`. Parse the manifest's
`requires.speckit_version` range and fail closed unless the active version is
within `>=0.12.11,<0.13.0`. If the manifest, version, or version range cannot
be read and parsed, stop and report the compatibility failure. Do not invoke a
hook, prerequisite script, reviewer, or write an artifact after a failed
preflight.

## User Input

```text
$ARGUMENTS
```

You **MUST** consider the user input before proceeding (if not empty).

## Pre-Execution Checks

**Check for extension hooks (before analysis)**:
- Check if `.specify/extensions.yml` exists in the project root.
- If it exists, read it and look for entries under the `hooks.before_analyze` key
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

    Wait for the result of the hook command before proceeding to the Goal.
    ```
    After emitting the block above you MUST actually invoke the hook and wait for it to finish before continuing. Run it the same way you would run the command yourself in this agent/session (the invocation may differ from the literal `{command}` id shown above, e.g. a skills-mode agent runs it as `/skill:speckit-...` or `$speckit-...`). Emitting the block alone does not run the hook.
- If no hooks are registered or `.specify/extensions.yml` does not exist, skip silently

## Goal

Coordinate independent, focused, read-only reviews after
`__SPECKIT_COMMAND_TASKS__` has produced `tasks.md` and `execution-plan.md`.
Persist cumulative source reports, a deterministic aggregate, and an explicit
human approval record. Follow
`.specify/presets/timors-agentic-workflow/references/analysis-and-approval.md`
exactly.

The coordinator performs no semantic review. It MUST NOT replace a missing or
failed fresh-context reviewer with its own judgment, add a broad review pass,
or repair reviewer conclusions.

## Operating Constraints

**READ-ONLY PLANNING REVIEW**: Reviewers receive read access only. Do not modify
`spec.md`, `plan.md`, `tasks.md`, `execution-plan.md`, design artifacts, the
constitution, or reviewer packets. The intended writes are only cumulative
`reviews/<role-id>.md` reports and aggregate `analysis.md`, and only after both
preflights succeed.

**Constitution Authority**: The project constitution (`/memory/constitution.md`)
is non-negotiable within review scope. A constitution change must occur through
an explicit owning workflow, never through analysis.

## Execution Steps

### 1. Initialize Analysis Context

Run `{SCRIPT}` once from the repository root and parse `FEATURE_DIR` and
`AVAILABLE_DOCS`. Derive `spec.md`, `plan.md`, `tasks.md`,
`execution-plan.md`, `analysis.md`, and `reviews/` paths under `FEATURE_DIR`.
Abort when any required artifact is missing and direct the user to its owning
command. For single quotes in args like "I'm Groot", use escape syntax such as
`'I'\''m Groot'`, or use double quotes when possible.

### 2. Run Deterministic Shared Validation

Read and execute the complete protocol in
`.specify/presets/timors-agentic-workflow/references/artifact-validation.md`
against the current artifacts and selected project-role packets. This installed
reference is the sole source for structural invariants; do not copy, infer,
weaken, or repair its rules. Any failure blocks before report writes or reviewer
dispatch and instructs the user to regenerate invalid execution artifacts with
`__SPECKIT_COMMAND_TASKS__`.

### 3. Establish Or Resume The Analysis Run

Apply the run, declaration, trigger-union, selected-role, packet, and recovery
rules in `analysis-and-approval.md`. The initial run selects every core and
project role. Later runs require an explicit human declaration. Resume an
incomplete `ARNN`; never start a newer run while it is incomplete.

Write `analysis.md` with the run marked `Incomplete` before dispatch and retain
the existing Approval Record unchanged. Retain valid applicable rounds. If an
invocation may have started without an attributable result, record
`Recovery required` and do not redispatch until an explicit human recovery
decision clears that state.

### 4. Resolve Reviewer Bindings Without Invocation

Resolve each role independently, project-local before global, using explicit
role preference, a unique preferred/default candidate, the sole eligible
candidate, or explicit human choice. Before invocation, verify the planned
tier, concrete model binding, read access to every required input, native fresh
context, invokability, and attributable-result support from runtime metadata or
configuration. Missing or malformed packets and unavailable or ambiguous
bindings block.

Do not probe, calibrate, canary, or invoke a candidate to test suitability. The
first invocation is the role's real review. Concrete provider, model, candidate,
and dispatch identifiers are transient and MUST NOT enter tracked reports.

### 5. Dispatch Focused Reviews

Dispatch only required unfinished roles, each in an independent fresh context
with its role packet, required read-only inputs, current `ARNN`, trigger, prior
report when present, and the report-template contract. A reviewer cannot edit
planning artifacts. Reject missing, malformed, semantically unattributable, or
schema-invalid results; do not reinterpret or complete them in coordinator
context. Retry only when runtime evidence proves no invocation started.

Each reviewer progressively loads only its packet's required artifacts, builds
the semantic inventory needed for its bounded question, performs that packet's
detection passes, and assigns the protocol severities. This preserves focused,
high-signal analysis without moving semantic inventories, detection, or
severity judgment into the coordinator.

Append each valid complete current-snapshot round to its cumulative report.
Record only `Yes` attestations for planned tier, fresh context, and attributable
result after observing the underlying transient evidence.

### 6. Aggregate Deterministically

Materialize `analysis.md` from the applicable source rounds without semantic
review. Preserve source IDs, group only findings with the same location,
defect, and correction, use the highest source severity, and expose all
disagreements as reviewer conflicts. Project the role-owned requirement
coverage and execution-plan validation rows without creating new semantic
conclusions. Keep the aggregate compact while retaining the complete current
source records in cumulative role reports. Preserve monotonic, non-reused
aggregate, conflict, and deviation IDs and preserve the Approval Record across
reruns.

Mark the run `Complete` only when exactly one valid applicable round exists for
every selected role. Compute the aggregate verdict and deviation statuses by
the installed reference. Never infer approval from a verdict and never
automatically revoke or reset a recorded human decision.

### 7. Provide Next Actions And Approval Gate

Present the aggregate and request an explicit human decision. Record only
`Pending`, `Approved`, `Changes requested`, or `Approved with deviations`, and
only a human may own deviation decisions and rationale. Enforce the exact
authorization matrix in `analysis-and-approval.md`; incomplete or uncovered
blocking states never authorize implementation.

Route remediation to the owning workflow. Use
`__SPECKIT_COMMAND_SPECIFY__` for requirements, `__SPECKIT_COMMAND_PLAN__` for
design, `__SPECKIT_COMMAND_TASKS__` for tasks or decomposition, and proceed to
`__SPECKIT_COMMAND_IMPLEMENT__` only when the stored aggregate and Approval
Record authorize it.

### 8. Offer Remediation

Ask: "Would you like me to suggest concrete remediation edits for the top N
issues?" Do not apply edits automatically. Any follow-up editing command
requires explicit user approval and runs separately through the owning command.

### 9. Check for extension hooks

After reporting, check if `.specify/extensions.yml` exists in the project root.
- If it exists, read it and look for entries under the `hooks.after_analyze` key
- If the YAML cannot be parsed or is invalid, skip hook checking silently and continue normally
- Filter out hooks where `enabled` is explicitly `false`. Treat hooks without an `enabled` field as enabled by default.
- For each remaining hook, do **not** attempt to interpret or evaluate hook `condition` expressions:
  - If the hook has no `condition` field, or it is null/empty, treat the hook as executable
  - If the hook defines a non-empty `condition`, skip the hook and leave condition evaluation to the HookExecutor implementation
- For each executable hook, output the following based on its `optional` flag:
  - **Optional hook** (`optional: true`):
    ```
    ## Extension Hooks

    **Optional Hook**: {extension}
    Command: `/{command}`
    Description: {description}

    Prompt: {prompt}
    To execute: `/{command}`
    ```
  - **Mandatory hook** (`optional: false`):
    ```
    ## Extension Hooks

    **Automatic Hook**: {extension}
    Executing: `/{command}`
    EXECUTE_COMMAND: {command}
    ```
    After emitting the block above you MUST actually invoke the hook and wait for it to finish before continuing. Run it the same way you would run the command yourself in this agent/session (the invocation may differ from the literal `{command}` id shown above, e.g. a skills-mode agent runs it as `/skill:speckit-...` or `$speckit-...`). Emitting the block alone does not run the hook.
- If no hooks are registered or `.specify/extensions.yml` does not exist, skip silently

## Operating Principles

- Preserve reviewer ownership: validate and aggregate source conclusions; do
  not create semantic findings in coordinator context.
- Fail closed on missing evidence, invalid artifacts, malformed reports,
  unavailable bindings, and incomplete or inconsistent state.
- Store no content hashes, artifact fingerprints, concrete model IDs, candidate
  IDs, or dispatch IDs.
- Report zero findings gracefully and still require explicit human approval.

## Context

{ARGS}
