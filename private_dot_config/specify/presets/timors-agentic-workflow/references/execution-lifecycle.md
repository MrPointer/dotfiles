# Delegated Implementation Lifecycle

## Authority And Boundaries

This reference governs `speckit.implement` for preset protocol `0.1.0`. The command executes the validated decomposition; it does not redesign it. Never invent, remove, regroup, or change planned tasks, dependencies, concurrency policy, file ownership, model tier, skills, capabilities, contracts, tests, acceptance criteria, or verification. A needed change is a planning correction and blocks affected execution.

The command owns coordination directly. It MUST NOT invoke `executing-plans`, provision workers through the preset, or let workers dispatch other workers. Installed preset references are the only preset policy source.

## Command-Wide Preflight

Before extension hooks, `{SCRIPT}`, progress creation or update, workspace or branch mutation, project-file mutation, or dispatch:

1. Run the compatibility preflight in `protocol-compatibility.md`.
2. Resolve the current feature and its installed feature directory read-only, without a setup or prerequisite script. Require readable `spec.md`, `plan.md`, `tasks.md`, and `execution-plan.md`, plus every protocol artifact selected by those files.
3. Read `.specify/presets/timors-agentic-workflow/references/artifact-validation.md` and rerun that complete deterministic validation against the current `tasks.md`, `execution-plan.md`, and every selected reviewer packet. This is the same installed validation source used by `speckit.analyze`; do not copy, abbreviate, or define a second invariant set here.
4. Require the latest analysis report for this exact artifact set to have final status `Complete`. Match the exact selected core and project roles, each role's `Role Verdict`, and every applicable report round. Reject missing, extra, stale, duplicated, or mismatched role/report identities.
5. Verify structural consistency among the final aggregate verdict, findings, conflict dispositions, deviations, and artifact identities. Require an authorizing Approval Record for the exact validated artifact set and report round.

Fail closed on any missing, malformed, stale, contradictory, or unauthorized evidence. This gate checks structural consumption safety only: it does not repeat semantic review, reassess findings, or revoke an approval because the implementer disagrees with it. Regeneration or correction remains owned by the command that owns the invalid artifact.

Only after the entire gate passes may before-hooks and `{SCRIPT}` run. If either changes a gated artifact, rerun the entire gate before any prohibited mutation or dispatch.

## Progress Initialization And Resume

After the preflight and upstream checklist gate, initialize or resume the ignored `specs/<feature>/progress.md` from `templates/progress-template.md`. Ensure the project ignore mechanism excludes this file before writing it. Preserve every exact section, column, and controlled state in the template, and update it after each attributable transition.

Use this resume matrix before dispatch:

| Observed state | Required action |
|---|---|
| Tasks checked, group `done`, and valid checkpoint plus Completed Artifacts evidence | Trust the completed coverage. |
| Valid result/checkpoint/artifacts with an interrupted progress transition | Finish only the interrupted progress and task-checkbox transition. |
| Task checked but completion evidence missing or invalid | Block; reconstruct evidence or explicitly return it to pending. |
| Task unchecked but group is `done` or has a completion checkpoint | Block the contradiction; do not infer which side is correct. |
| Local progress, workspace, branch, checkpoint, or attributable dispatch/result state is missing | Block and recover explicitly; never redispatch automatically. |

## Ready-Group Gate

Bind a group only when it is ready. Future groups need not bind workers or workspaces early. Immediately before a ready group mutates a workspace, edits, or dispatches, verify all of the following for that group:

- every logical and policy prerequisite and required checkpoint is present;
- every prerequisite output is recorded and matches its producer contract;
- the planned tier resolves to a concrete worker and concrete model;
- every exact required skill is available and invokable at its recorded path;
- required capabilities map to effective permissions, with no undeclared capability;
- worker discovery, selection, binding, and dispatch are available;
- the absolute target workspace and its evidence are recorded;
- dirty-state policy is satisfied;
- concurrent worker capacity and output/build/cache safety are proven;
- the planned Test Expectation and verification commands remain enforceable.

A parallel set starts only after every member passes. If any member fails, none of that set may mutate or dispatch. `Parallel allowed` groups may be serialized when capacity or safety requires it, but the fallback and reason must be recorded. Never parallelize a `Linear DAG` or remove a dependency.

## Canonical Worker Packet

Construct each worker packet without interpretation by joining:

1. the canonical Execution Groups row;
2. that group's complete Execution Group Details section;
3. the current canonical `tasks.md` lines for exactly its `Covers` task IDs; and
4. prerequisite outputs from Completed Artifacts, joined through the matching `DFNN` flow and `CTNN` producer/consumer contracts.

Pass absolute workspace targeting and require the Result evidence below. Do not add or change policy, ownership, implementation direction, contracts, tests, acceptance criteria, or verification. A producer result that does not match its planned contract blocks source-group completion and requires planning correction; it is not normalized in progress or in a consumer packet.

## Provider-Neutral Handshake Evidence

Record these stages for each binding:

1. **Discovery**: candidate identity, scope, purpose, planned tier source, concrete model source, exact skill paths, capability-to-permission mapping, and invokability.
2. **Selection**: search project-local candidates before global candidates, then apply explicit configuration. Select only a unique preferred candidate, the sole eligible candidate, or a human decision; ambiguity blocks.
3. **Binding**: selected identity, concrete model, exact skills, and effective permissions.
4. **Workspace**: absolute target and evidence that invocation will execute there.
5. **Dispatch**: attributable invocation ID and start evidence.
6. **Result**: the same attributable identity, terminal status, changed files, tests, verification, and blockers.

Do not run suitability probes, calibration tasks, canaries, or model-confirmation calls. Feature-local planned tier classification is the only tier evidence; do not infer tier from general reputation or observed output. Reuse only recorded metadata or an explicit operator decision.

Retry is allowed only with proof that no invocation started. If start is ambiguous or an invocation may have started, retain all state and mark `Recovery required`; never retry automatically.

## Tests, Results, And Completion

Tests remain with the implementation worker that owns the group. Enforce the exact planned Test Expectation. Missing evidence is pushed back once to the same binding and workspace; if it remains missing, block rather than substituting a different worker or explanation.

A successful worker result moves the group to `ready for integration/checkpoint`, not `done`. Complete a group only in this order:

1. integrate its exact result;
2. run required post-integration verification;
3. create the signed checkpoint;
4. persist checkpoint, handshake, test, verification, and Completed Artifacts evidence;
5. update every covered canonical task checkbox;
6. mark the group `done`.

Completed Artifacts records the group's `Flow ID`, `Contract ID`, `Files / Outputs`, `Result Summary`, and `Checkpoint`. Use `None` for Flow ID and Contract ID only when output does not cross a group boundary. Dependent groups consume only these checkpointed records.

## Adoption, Removal, And Unsupported Transitions

Adoption requires upstream `tasks.md` plus a protocol-valid execution plan. Reject projects with only upstream tasks, no execution plan, or mixed protocol/non-protocol execution state.

For planned work removed before execution, either revert the planning artifacts or retain inert history; do not execute it. For in-flight removal, finish the pinned work, pause it with all state retained, or explicitly abandon it back to the recorded execution base. Never automatically delete branches, worktrees, checkpoint refs, dirty files, progress, or artifacts.

Removing this preset while execution is active is unsupported. Pause or complete/abandon execution first.
