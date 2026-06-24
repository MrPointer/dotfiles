---
name: executing-plans
description: Use whenever executing, continuing, resuming, or finishing an implementation plan or plan directory. Orchestrates assigned worker dispatch, integrated test-first discipline, progress checkpointing, isolated workspaces, checkpoint integration, and blocker batching.
---

# Executing Plans

Execute implementation plans with persisted progress, assigned workers, integrated test-first discipline, and checkpointed integration when needed.

This file is the execution spine. Load referenced files only when their trigger applies.

## Contents

- [Runtime Binding](#runtime-binding)
- [Core Invariants](#core-invariants)
- [Conditional References](#conditional-references)
- [Workflow](#workflow)
- [Final Rules](#final-rules)

## Runtime Binding

Before doing any work, determine the active runtime from the system prompt and environment banner, then read exactly one adapter:

- **OpenCode runtime**: [references/runtime-opencode.md](references/runtime-opencode.md)
- **Codex runtime**: [references/runtime-codex.md](references/runtime-codex.md)
- **Claude runtime**: [references/runtime-claude.md](references/runtime-claude.md)

If the runtime signal is ambiguous, ask the user rather than guessing. Do not load or mix other runtime adapters in the same turn. If a runtime adapter conflicts with this file, this file is authoritative.

Use runtime-native names when reading or writing concrete artifacts: worker agent definitions, dispatch recipes, custom subagents, or equivalent runtime bindings.

## Core Invariants

- **Progress is always persisted**: update `progress.md` after every meaningful step so execution can resume.
- **Assigned bindings are mandatory**: if a plan assigns a worker or model tier, dispatch through that binding. Do not self-execute assigned implementation work.
- **Testable work stays with the implementer**: the same worker owns tests and code, follows the required testing/TDD skills, makes a test-first attempt when practical, and explains when that is not practical.
- **Independent work continues**: blocked tasks do not stop unrelated tasks.
- **Plan concurrency policy is binding**: if the master plan says `Linear DAG`, run at most one implementation sub-plan at a time even when tasks appear logically independent.
- **Parallel implementation requires isolated workspaces**: never run concurrent implementers in the same dirty workspace.
- **Dependent work starts from checkpointed prerequisites**: a prerequisite is not complete until its output is in an integration-branch checkpoint commit.
- **Checkpoint commits are local plumbing**: never push them. Final review is the aggregate dirty diff after mixed reset, not checkpoint history.
- **Plan artifacts stay with the coordinator**: never copy plan files, review files, anchors, or `progress.md` into worker worktrees. Pass task packets inline.
- **Dirty worktree creation is user-authorized only**: before creating isolated worktrees, stop on tracked or untracked non-ignored dirty files unless the user explicitly accepts the risk.
- **Anchors are not progress files**: anchors preserve feature-level rationale and handoff context; `progress.md` preserves mechanical execution state.

## Conditional References

| Trigger | Load |
|---------|------|
| Creating or repairing progress | [assets/progress-template.md](assets/progress-template.md) |
| Any task has testable acceptance criteria | [references/testable-work.md](references/testable-work.md) |
| Any task uses an isolated worktree, policy-allowed parallel implementation, or build/cache reuse | [references/workspace-isolation.md](references/workspace-isolation.md) |
| Plan has multiple sub-plans, dependencies, checkpoint commits, or task worktree integration | [references/checkpoint-integration.md](references/checkpoint-integration.md) |

## Workflow

### 1. Read the Plan

Read the plan file or directory at the given path. Plans vary in shape, so extract what exists:

- tasks and execution order
- dependencies, or sequential order when dependencies are not specified
- the master plan's concurrency policy, including whether `Linear DAG` forbids parallel dispatch
- tasks that can run independently, and tasks that may run in parallel only when the concurrency policy allows it
- execution binding and model assignments
- file ownership and conflict boundaries
- required verification per task
- build/cache reuse notes for isolated workspaces

If the plan or user references an active feature anchor, read it for intent, constraints, rationale, rejected alternatives, unresolved questions, and handoff context. If the plan names a feature but not an anchor path, look for one using project conventions, then `docs/context/<topic>-anchor.md`. Do not infer an anchor from unrelated old feature docs.

### 2. Initialize Or Resume Progress

The progress file lives alongside the plan: `progress.md` in the plan directory, or next to a single-file plan.

- If progress exists, read it, resume from the last checkpoint, and verify completed tasks' artifacts still exist. If artifacts are missing, mark the task for re-execution.
- If progress does not exist, create it from [assets/progress-template.md](assets/progress-template.md).

`progress.md` is the source of truth for task status, dispatch audit, test evidence, completed artifacts, blockers, regressions, checkpoints, and final review state.

If an anchor is active, update it only for durable feature context: new decisions, changed constraints, spec or plan deviations, unresolved questions, or handoff state. Do not duplicate progress tables in the anchor.

### 3. Run Execution Preflight

Before editing implementation files, verify execution mechanics and record the audit in progress.

Check these items:

- every planned implementer worker, model tier, runtime dispatch mechanism, and required testing skill
- whether assigned bindings are discoverable and invokable; diagnose and retry once for mechanical failures, then stop and ask
- whether coordinator self-execution is allowed; it is allowed only for progress/coordination artifacts or a single trivial sub-plan with no explicit binding requirement
- whether checkpoint integration is required; if yes, load [checkpoint-integration.md](references/checkpoint-integration.md) and establish or resume the integration branch before implementation
- whether testable-work expectations apply; if yes, load [testable-work.md](references/testable-work.md)
- whether the plan's concurrency policy allows parallel implementation; if it says `Linear DAG`, record the serialized schedule and do not dispatch independent tasks concurrently
- whether isolated workspaces or cache reuse is required; if yes, load [workspace-isolation.md](references/workspace-isolation.md) and run its dirty-state preflight before creating worktrees
- planned worker, actual worker, model/effort, dispatch evidence, implementation workspace, dirty-state preflight, cache reuse, checkpoint commit, integration status, test evidence, and verification status for each task

Proceed only after the audit is complete or the user explicitly authorizes a deviation.

### 4. Execute Tasks

Execute tasks in dependency order and honor the plan's concurrency policy. If the policy says `Linear DAG`, dispatch at most one implementation task at a time. If a task blocks, continue with independent tasks and pause only when remaining tasks depend on blocked work.

#### 4.1 Choose The Workspace

- Sequential tasks may run in the main execution workspace unless the plan or user requires isolation.
- Concurrent implementation tasks require task-scoped worktrees. Use [workspace-isolation.md](references/workspace-isolation.md).
- A `Linear DAG` policy forbids concurrent implementation dispatch; do not create parallel task worktrees solely because tasks are logically independent.
- Dependent task worktrees must start from an integration-branch checkpoint that contains all prerequisite outputs. If prerequisite outputs are only dirty files, checkpoint them first or serialize the work.

#### 4.2 Set Testable-Work Expectations

Skip test-first expectations for tasks with no testable work: documentation, pure file moves, configuration-only changes, or tasks without acceptance criteria.

For tasks with testable acceptance criteria, use [testable-work.md](references/testable-work.md) to set expectations. Record the result in progress:

- `applies` when the implementer should make a test-first attempt or explain why not
- `skipped: <reason>` when the task has no testable behavior or project docs declare the area untestable
- `blocked: <reason>` when acceptance criteria are too ambiguous to test or verify

Do not create a separate test-writing workflow. The implementer owns tests and code.

#### 4.3 Dispatch Implementation

Dispatch the implementer through the assigned runtime binding.

Pass the full task packet inline with prerequisite outputs and required skills. For testable behavior changes, tell the implementer to follow the required testing/TDD skills, make a test-first attempt when practical, and report tests added or updated plus verification results. If test-first work or new tests are not practical, the implementer must explain why.

The implementer returns status, files created or modified, tests added or updated, verification results, and any blockers.

Before accepting a successful implementer result for testable behavior, check that the result includes test evidence or a concrete explanation for missing tests. If that evidence is absent, push back once through the same implementer binding before blocking the task.

#### 4.4 Handle Results

Record every result in progress.

Use these task states:

- `ready for integration/checkpoint`: implementation is complete and required tests pass
- `blocked: implementation failure`: required verification fails or the implementer cannot complete the task
- `blocked: missing test evidence`: testable behavior was implemented without tests or a concrete explanation
- `blocked: acceptance-criteria ambiguity`: the implementer cannot determine what behavior to test or verify
- `blocked: regression`: existing tests regressed; this takes priority over other states
- `blocked: integration`: task worktree integration, checkpoint commit creation, or post-integration verification failed
- `done`: result is integrated, verified, and checkpointed when checkpointing applies

If the task ran in a task-scoped worktree, integrate and checkpoint it before marking it `done`. If the task ran sequentially on an integration branch, create the checkpoint after verification before marking it `done`. Use [checkpoint-integration.md](references/checkpoint-integration.md) for the mechanics.

### 5. Resolve Blocks

At a natural breakpoint, present all blocked items together:

```markdown
## Blocked Items

| Task | Status | Details |
|------|--------|---------|
| 02-api-layer | acceptance-criteria ambiguity | AC says "return error" but does not specify status-code behavior |
| 04-integration | implementation failure | Cannot satisfy timeout recovery without async support outside task scope |
```

The user resolves each item by adjusting the acceptance criteria, authorizing a no-test exception, changing the task, accepting the implementation, or authorizing a deviation. Re-run affected tasks from the appropriate implementation step.

### 6. Complete Execution

When all tasks are `done`:

1. Run the required full verification on the final integration state.
2. Record final verification results in progress.
3. If checkpoint integration was used, record the final checkpoint range and materialize the aggregate dirty review diff according to [checkpoint-integration.md](references/checkpoint-integration.md).
4. Update progress to `complete` and note the final review state.
5. If an anchor is active, update its final active-work state and note durable decisions that should graduate to ADRs or permanent docs.
6. Report what was built, issues encountered, final verification results, checkpoint range when applicable, final review branch/state, and any ADR or docs follow-up.

## Final Rules

- Do not create a separate test-writing workflow unless the approved plan explicitly requires it.
- Do not accept testable implementation work with no test evidence and no concrete explanation; push back once, then block.
- Do not run concurrent implementers in one workspace.
- Do not override a master plan's `Linear DAG` policy with inferred independence, cache seeding, or isolated worktrees.
- Do not launch dependent work from uncheckpointed dirty state.
- Do not create isolated worktrees from a dirty workspace unless the user explicitly authorizes it after being warned.
- Do not push checkpoint commits.
- Do not copy plan/progress artifacts into worker worktrees.
- Do not assume fresh worktrees contain ignored build artifacts.
- Do not silently upgrade, downgrade, or replace assigned workers or models.
- Do not let test-process ceremony overtake execution. Use the required testing skills and acceptance criteria as the source of testing depth.
