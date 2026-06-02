---
name: executing-plans
description: Use when executing implementation plans. Orchestrates task execution with testability-gated TDD, progress checkpointing, and dispute batching. Works with any plan structure that has identifiable tasks and acceptance criteria. Handles resumption after interruptions.
---

# Executing Plans

Execute implementation plans with persisted progress, assigned workers, optional structural TDD, and checkpointed integration when needed.

This file is the execution spine. Load referenced files only when their trigger applies.

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
- **Structural TDD is strict when used**: the test author sees only acceptance criteria and code surface; the implementer sees the full task plus tests. If isolation or testability cannot be verified, skip structural TDD only according to [structural-tdd.md](references/structural-tdd.md).
- **Tests are immutable to implementers**: implementers report disputes instead of editing test-author tests.
- **Independent work continues**: blocked tasks do not stop unrelated tasks.
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
| Any task has testable acceptance criteria and structural TDD may apply | [references/structural-tdd.md](references/structural-tdd.md) |
| Any task uses an isolated worktree, parallel implementation, or build/cache reuse | [references/workspace-isolation.md](references/workspace-isolation.md) |
| Plan has multiple sub-plans, dependencies, checkpoint commits, or task worktree integration | [references/checkpoint-integration.md](references/checkpoint-integration.md) |

## Workflow

### 1. Read the Plan

Read the plan file or directory at the given path. Plans vary in shape, so extract what exists:

- tasks and execution order
- dependencies, or sequential order when dependencies are not specified
- tasks that can run independently or in parallel
- execution binding and model assignments
- file ownership and conflict boundaries
- required verification per task
- build/cache reuse notes for isolated workspaces

If the plan or user references an active feature anchor, read it for intent, constraints, rationale, rejected alternatives, unresolved questions, and handoff context. If the plan names a feature but not an anchor path, look for one using project conventions, then `docs/context/<topic>-anchor.md`. Do not infer an anchor from unrelated old feature docs.

### 2. Initialize Or Resume Progress

The progress file lives alongside the plan: `progress.md` in the plan directory, or next to a single-file plan.

- If progress exists, read it, resume from the last checkpoint, and verify completed tasks' artifacts still exist. If artifacts are missing, mark the task for re-execution.
- If progress does not exist, create it from [assets/progress-template.md](assets/progress-template.md).

`progress.md` is the source of truth for task status, dispatch audit, test artifacts, completed artifacts, blockers, disputes, regressions, checkpoints, and final review state.

If an anchor is active, update it only for durable feature context: new decisions, changed constraints, spec or plan deviations, unresolved questions, or handoff state. Do not duplicate progress tables in the anchor.

### 3. Run Execution Preflight

Before editing implementation files, verify execution mechanics and record the audit in progress.

Check these items:

- every planned implementer worker, test author worker, model tier, and runtime dispatch mechanism
- whether assigned bindings are discoverable and invokable; diagnose and retry once for mechanical failures, then stop and ask
- whether coordinator self-execution is allowed; it is allowed only for progress/coordination artifacts or a single trivial sub-plan with no explicit binding requirement
- whether checkpoint integration is required; if yes, load [checkpoint-integration.md](references/checkpoint-integration.md) and establish or resume the integration branch before implementation
- whether structural TDD is applicable; if yes, load [structural-tdd.md](references/structural-tdd.md)
- whether isolated workspaces or cache reuse is required; if yes, load [workspace-isolation.md](references/workspace-isolation.md) and run its dirty-state preflight before creating worktrees
- planned worker, actual worker, model/effort, dispatch evidence, implementation workspace, dirty-state preflight, cache reuse, checkpoint commit, integration status, and TDD gate for each task

Proceed only after the audit is complete or the user explicitly authorizes a deviation.

### 4. Execute Tasks

Execute tasks in dependency order. If a task blocks, continue with independent tasks and pause only when remaining tasks depend on blocked work.

#### 4.1 Choose The Workspace

- Sequential tasks may run in the main execution workspace unless the plan or user requires isolation.
- Concurrent implementation tasks require task-scoped worktrees. Use [workspace-isolation.md](references/workspace-isolation.md).
- Dependent task worktrees must start from an integration-branch checkpoint that contains all prerequisite outputs. If prerequisite outputs are only dirty files, checkpoint them first or serialize the work.

#### 4.2 Gate Structural TDD

Skip test authoring for tasks with no testable work: documentation, pure file moves, configuration-only changes, or tasks without acceptance criteria.

For tasks with testable acceptance criteria, use [structural-tdd.md](references/structural-tdd.md) to decide whether structural TDD is feasible. Record the result in progress:

- `used isolated workspace` when structural TDD was enforced
- `skipped: <reason>` when no testable surface or no enforceable isolation exists
- `blocked: <reason>` when an available isolation path or required scaffold fails and user input is needed

When structural TDD is used, the test author writes failing tests from acceptance criteria only and returns test file paths plus a summary of what each test verifies.

#### 4.3 Dispatch Implementation

Dispatch the implementer through the assigned runtime binding.

Without structural TDD, pass the full task packet inline and have the implementer work directly from the acceptance criteria. Existing tests still must pass.

With structural TDD, pass the full task packet inline, test file paths, prerequisite outputs, and required skills. Tell the implementer that test-author tests are immutable. It may refactor implementation after tests pass, but it must not edit those tests. If it believes a test is wrong, it reports a dispute with the test file, test name, and explanation.

The implementer returns status, files created or modified, test results, and any disputes.

#### 4.4 Handle Results

Record every result in progress.

Use these task states:

- `ready for integration/checkpoint`: implementation is complete and required tests pass
- `blocked: implementation failure`: tests fail and the implementer does not claim the tests are wrong
- `blocked: test dispute`: implementer disputes one or more test-author tests
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
| 02-api-layer | test dispute | test_returns_404: "AC says 'return error' but status code is unspecified" |
| 04-integration | implementation failure | Cannot satisfy timeout recovery without async support outside task scope |
```

The user resolves each item by adjusting the acceptance criteria, fixing a test, changing the task, accepting the implementation, or authorizing a deviation. Re-run affected tasks from the appropriate step: test authoring if acceptance criteria changed, implementation if tests or implementation context changed.

### 6. Complete Execution

When all tasks are `done`:

1. Run the required full verification on the final integration state.
2. Record final verification results in progress.
3. If checkpoint integration was used, record the final checkpoint range and materialize the aggregate dirty review diff according to [checkpoint-integration.md](references/checkpoint-integration.md).
4. Update progress to `complete` and note the final review state.
5. If an anchor is active, update its final active-work state and note durable decisions that should graduate to ADRs or permanent docs.
6. Report what was built, issues encountered, final verification results, checkpoint range when applicable, final review branch/state, and any ADR or docs follow-up.

## Final Rules

- Do not let the test author see design decisions, plan paths, task paths, feature names, or other breadcrumbs to the plan.
- Do not claim structural TDD without verified physical isolation and prompt isolation.
- Do not let implementers modify test-author tests.
- Do not run concurrent implementers in one workspace.
- Do not launch dependent work from uncheckpointed dirty state.
- Do not create isolated worktrees from a dirty workspace unless the user explicitly authorizes it after being warned.
- Do not push checkpoint commits.
- Do not copy plan/progress artifacts into worker worktrees.
- Do not assume fresh worktrees contain ignored build artifacts.
- Do not silently upgrade, downgrade, or replace assigned workers or models.
- Do not analyze testability improvements when TDD is skipped; record the skip reason briefly and continue.
