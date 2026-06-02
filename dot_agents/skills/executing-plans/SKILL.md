---
name: executing-plans
description: Use when executing implementation plans. Orchestrates task execution with testability-gated TDD, progress checkpointing, and dispute batching. Works with any plan structure that has identifiable tasks and acceptance criteria. Handles resumption after interruptions.
---

# Executing Plans

Orchestrate the execution of implementation plans. Read the plan, identify tasks and their ordering, maintain progress checkpoints, and handle interruptions gracefully.

When the code surface supports it, this skill structurally separates test authoring from implementation — a separate agent writes tests from acceptance criteria alone, without seeing design decisions. The implementer then receives the full task plus the tests. This prevents the common failure mode where agents write tests that rationalize their implementation rather than verify requirements. When the code isn't ready for TDD (missing test infrastructure, no testable seams, tight coupling), the executor skips structural separation rather than fighting the codebase.

## Runtime Binding

This skill has one canonical workflow. Runtime files only map that workflow to the active agent runtime's mechanics.

Before doing any work, determine the active runtime and read exactly one adapter:

- **OpenCode runtime** → [references/runtime-opencode.md](references/runtime-opencode.md)
- **Codex runtime** → [references/runtime-codex.md](references/runtime-codex.md)
- **Claude runtime** → [references/runtime-claude.md](references/runtime-claude.md)

**Determining the active runtime**: Check the system prompt and environment banner for identifying markers (e.g., "OpenCode", "Claude Code", "Codex CLI"). If the signal is ambiguous, ask the user rather than guessing — reading the wrong adapter silently breaks assumptions downstream.

Do not load or mix instructions from the other runtime adapter in the same turn. If a runtime adapter conflicts with this file, this file is authoritative.

**Terminology bridge**: This skill uses runtime-neutral terms. Claude's runtime calls execution bindings "worker agent definitions"; Codex's runtime calls them "dispatch recipes"; OpenCode's runtime uses custom subagent definitions. Use whichever term is native to the active runtime when writing or reading concrete artifacts; the canonical workflow terms are used only in this file.

## Core Principles

1. **Tests and Implementation Are Structurally Separated — When Feasible**: The test author sees only acceptance criteria. The implementer sees the full task context plus the tests. Neither can influence the other's context. However, this separation requires a testable code surface. When the target code lacks the infrastructure or seams to support isolated testing, the executor skips structural TDD rather than producing tests that fight the codebase.
2. **Progress Is Always Persisted**: After every meaningful step, update the progress file. If the session drops, the executor resumes from the last checkpoint — not from scratch.
3. **Tests Are Immutable to the Implementer**: The implementer cannot modify tests. If it believes a test is wrong, it reports this and moves on. Disputes are batched for human resolution.
4. **Independent Work Continues**: When a task is blocked (test dispute, failure), the executor continues with independent tasks.
5. **Parallel Implementation Is Workspace-Isolated**: DAG independence means tasks do not need each other's outputs; it does not make a shared dirty workspace safe. When implementation tasks run concurrently, dispatch each task in its own isolated worktree and integrate the result deliberately. For build-heavy projects, isolated worktrees must be seeded with ignored build/cache artifacts by hard link before dispatch so workers do not rebuild the world from an empty cache. If isolated dispatch or required cache seeding cannot be verified, serialize the tasks or ask the user rather than running concurrent implementers in the same workspace.
6. **Prerequisite State Is Materialized Through Checkpoint Commits**: Dependent worktrees must not build on dirty coordinator state they cannot see. Execute multi-sub-plan plans on a local integration branch, commit integrated sub-plan results there as temporary checkpoint commits, and launch dependent worktrees from a checkpoint that contains their prerequisites. These commits are execution plumbing, not review units.
7. **Final Review Is An Aggregate Dirty Diff**: The user reviews the complete feature at the end, not each checkpoint. After final verification, reset the integration branch back to the execution base with mixed-reset semantics so the final tree appears as dirty changes for normal IDE review.
8. **Execution Bindings Are Correctness, Not Optimization**: When a plan assigns worker bindings or model tiers, executing through those bindings is mandatory. The coordinator does not silently self-execute implementation work just because it can edit files.

## Workflow

### Step 1: Read the Plan

Read the plan file(s) at the given path. Understand the structure — plans vary in format:

- **Single file with numbered tasks**: Extract task list in order.
- **Directory with a master plan and task files**: Read the master plan for ordering, dependencies, and coordination details. Each task file is a self-contained unit.
- **Other structures**: Adapt. The executor needs to identify: what are the tasks, what order do they run in, and which can run independently.

Extract what's available:
- Task list and execution order
- Dependencies between tasks (if specified — otherwise assume sequential)
- Which tasks can run independently / in parallel (if specified)
- Execution binding assignments (if specified)
- File ownership per task (if specified — helps prevent conflicts during parallel execution)
- Whether each parallel group is concurrency-ready: non-overlapping file ownership, no unsequenced shared artifact, and an isolated implementer worktree path available through the active runtime adapter
- Build/cache artifact directories that must be available in isolated worktrees before compilation or test execution (if specified)

If the plan or user references an active feature anchor, read it for feature-level context before executing. If the plan names a feature but not an anchor path, look for one using project conventions, then `docs/context/<topic>-anchor.md`. The anchor explains intent, constraints, rationale, rejected alternatives, and handoff context; it is not the execution checkpoint. Do not infer an anchor from unrelated old feature docs.

### Step 2: Initialize or Resume Progress

Check for an existing progress file alongside the plan (`progress.md` in the plan directory, or next to a single-file plan).

**If progress exists**: Read it, determine current state, and resume from the last checkpoint. Verify that completed tasks' artifacts still exist (files, tests). If something is missing, mark that task for re-execution.

**If no progress exists**: Create the initial progress file from the plan. See [Progress File Format](#progress-file-format).

The progress file is the source of truth for task-by-task execution status, tests, blockers, disputes, and completion. If an anchor is active, update it only for feature-level context: new durable decisions, changed constraints, spec/plan deviations, unresolved questions, or handoff state. Do not duplicate progress tables or per-task status in the anchor.

### Step 2.5: Establish Execution Integration Branch

For plans with multiple sub-plans, task dependencies, or isolated task worktrees, create or resume a local execution integration branch before implementation begins. This branch exists so sub-plan outputs become real Git commits that later worktrees can use as their base without forcing the user to review each sub-plan separately.

1. **Record the review starting point**: Capture the current review branch or detached-HEAD state, the execution base commit (`HEAD` before execution starts), and the intended integration branch name in `progress.md`. Use a local branch name that is clearly temporary, such as `agent/<plan-id>` or another project convention.
2. **Require a safe implementation state**: The coordinator workspace may contain local plan files, review files, `progress.md`, anchors, or other explicitly coordinator-owned artifacts. These are expected and must not be removed, stashed, or committed as execution checkpoints. If there are dirty changes in implementation source, tests, config, generated code, or any file that could affect build/test behavior or final review diff, stop and ask the user how to isolate them before checkpoint execution starts. Do not create checkpoint/reset workflows on top of unknown implementation changes.
3. **Create or resume the integration branch**: Create the integration branch from the execution base, or resume the recorded branch if progress already exists. Verify it still descends from the recorded execution base before continuing.
4. **Checkpoint integrated results**: After a sub-plan result is integrated and its required verification passes, create a local checkpoint commit on the integration branch. Use the project's normal commit-signing policy; if signing is required and fails, stop rather than silently creating an unsigned checkpoint. Do not push checkpoint commits.
5. **Launch dependents from checkpoints**: A dependent task may start only from an integration-branch commit that contains all of its prerequisite outputs. If a prerequisite exists only as dirty files, commit it as a checkpoint first or serialize/block the dependent task.
6. **Prepare final review state**: After all tasks pass final verification, materialize the aggregate result for review by resetting the integration branch to the execution base with mixed-reset semantics (`git reset --mixed <execution-base>`). This moves checkpoint commits out of the branch tip while leaving the final feature changes as dirty files in the coordinator workspace. Switch back to the original review branch only when it is still at the execution base and checkout is safe; otherwise leave the dirty aggregate on the integration branch and report the branch/base clearly.

### Step 3: Execution Preflight

Before editing implementation files, verify the execution mechanics and record them in progress.

1. **Worker binding audit**: For every sub-plan, identify the planned implementer worker, test author worker, model tier, and runtime dispatch mechanism. If the plan has two or more sub-plans, or if any sub-plan has an explicit model tier, the coordinator must dispatch through the assigned execution binding. Coordinator self-execution is allowed only for progress files, coordination artifacts, or a single trivial sub-plan with no explicit binding requirement.
2. **Binding availability check**: Use the active runtime adapter to verify every assigned worker is discoverable and invokable before executing any sub-plan. If a binding is missing or cannot be invoked, diagnose and retry once when the cause is mechanical (for example, stale discovery or permission shape). If it still cannot be invoked, stop and ask the user. Do not fall back to the coordinator or a more expensive model.
3. **Integration branch check**: Verify the execution integration branch exists, starts from the recorded execution base, and is the branch used for checkpoint commits. Verify that final review will be materialized as a mixed-reset dirty diff rather than leaving checkpoint commits as the review surface.
4. **TDD isolation check**: For each task with testable acceptance criteria, verify whether the active runtime can route the test author into a Worktrunk (`wt`) or plain `git worktree` workspace. If a matching isolation path exists, attempt or otherwise concretely verify it before skipping structural TDD. If verification fails, record the attempted mechanism and stop to ask whether to fix isolation or explicitly skip structural TDD. Do not record a generic "runtime cannot provide isolated workspace" reason while an untried Worktrunk or git isolation path exists.
5. **Build/cache seeding check**: For every task that will run in an isolated workspace, identify ignored in-repository build/cache artifact directories that the package manager or compiler depends on for fast incremental builds. Examples include Rust `target/`, Bazel output trees, or other project-local build directories named by docs or config. If such a directory exists in the coordinator workspace, verify the isolated workspace can hard-link seed it from the coordinator workspace before dispatch. `cp -al <source-dir>/. <worktree-dir>/<relative-dir>/` is acceptable when the platform supports it and both paths are on the same filesystem; otherwise use an equivalent hard-link-preserving copy. Seed only ignored build/cache artifacts, never source files, plan files, review files, or `progress.md`. If required hard-link seeding cannot be verified for a build-heavy project, record the failure and ask before dispatching a worker that would rebuild from scratch.
6. **Implementer workspace check**: For every task that may run concurrently with another implementation task, verify a task-scoped implementer worktree can be created and that the active runtime can dispatch the assigned worker inside it. Prefer Worktrunk (`wt`) when it is installed and suitable. Otherwise fall back to plain `git worktree`. Do not use runtime-native worktree switching for this skill; runtime adapters are used to dispatch workers into the selected worktree and verify results. If no isolated implementation path can be verified, serialize that parallel group or ask the user for approval before weakening isolation.
7. **Progress audit update**: Add or update the progress file's execution audit with planned worker, actual worker, model/effort, dispatch evidence, implementation workspace, build/cache seeding status, checkpoint commit, integration status, and TDD gate status for each task.

Only proceed once the audit is complete or the user explicitly authorizes a deviation.

### Step 4: Execute Tasks

For each task in execution order (respecting dependencies):

#### 4.0 Workspace Selection

Before launching a task, choose its implementation workspace:

- **Sequential task**: The main execution workspace is acceptable unless the plan or user requests isolation.
- **Concurrent task**: Create a task-scoped implementer worktree using the active runtime adapter's workspace isolation priority: Worktrunk (`wt`) first when available, then plain `git worktree`.

Create the task worktree from the current integration branch checkpoint that already includes all prerequisite outputs, using the runtime adapter's branch, merge, or patch-transfer mechanism as needed. If prerequisite outputs exist only as local changes that cannot be represented in the task worktree, do not create a stale worktree; integrate and checkpoint the prerequisites first or serialize the dependent task. Do not copy plan files, review files, or `progress.md` into the task worktree. The coordinator keeps those artifacts in the original workspace and passes the sub-plan to the implementer as an inline task packet, not as a file path.

If the task worktree will compile or test build-heavy code, seed required ignored build/cache artifact directories before dispatch. Prefer hard-link seeding from the coordinator workspace so the initial artifact state is shared without duplicating the cache. Verify the destination contains the expected files and record the seeded relative paths in progress. If hard-link seeding fails because the platform, filesystem, or copy tool cannot create links, ask before dispatching a worker that would rebuild from an empty worktree.

If the runtime cannot verify isolated dispatch for a concurrent task, do not run that task concurrently in the shared workspace. Serialize it or stop and ask the user.

#### 4a. Test Authoring

**When to skip (no testable work)**: Some tasks don't have testable work — documentation updates, file moves, configuration changes, or tasks with no acceptance criteria. Skip test authoring for these and proceed directly to implementation.

**When to skip (code surface not testable)**: For tasks that have testable acceptance criteria, check whether the code surface supports TDD before spawning the test author. This is a gate, not an analysis — the goal is a binary decision, not a report on what needs fixing.

1. **Check project docs first**: Look for explicit testability statements in `AGENTS.md`, project documentation, or task-level notes. If the project or area is declared untestable, skip immediately — no exploration needed.

2. **Lightweight code surface check**: If docs don't cover it, use the active runtime adapter's cheapest suitable exploration mechanism to answer one question: can the components this task touches be tested in isolation? The agent looks for:
   - Interfaces or injectable dependencies (substitution points for test doubles)
   - An existing test framework and test patterns in the area
   - Whether tests could target the code without restructuring it first

   The agent returns a yes/no — not an analysis of what's wrong or recommendations for making it testable.

If the code surface isn't ready for TDD, skip test authoring and proceed directly to implementation without structural separation. Record the decision briefly in the progress file (Tests column: `skipped`, Notes: e.g., "declared untestable in AGENTS.md" or "no testable seams"). The implementer still receives the full task with acceptance criteria — it just doesn't receive pre-written tests and isn't bound by the test immutability constraint.

**Isolation capability gate**: Before structural TDD, use the active runtime adapter to verify that the test author can be dispatched inside a Worktrunk (`wt`) or plain `git worktree` workspace that does not contain the planning artifacts. If no enforceable isolation path exists, skip structural TDD and proceed directly to implementation. If the adapter names an isolation path but verification fails, stop and ask whether to fix isolation or explicitly skip structural TDD. Record the decision briefly in the progress file (Tests column: `skipped`, Notes: e.g., "no isolated worktree dispatch mechanism" or "user approved TDD skip after failed `<mechanism>` dispatch").

**Build/cache readiness gate**: If the test author will run in a temporary isolated workspace for build-heavy code, hard-link seed the same ignored build/cache artifact directories required for implementation worktrees before running compile or test checks. The test author can see build artifacts but must not see planning artifacts. If required seeding cannot be verified, ask before proceeding; do not silently make the test author rebuild the whole project from an empty worktree.

**Compilation readiness gate**: Before spawning the test author, verify that the target package compiles at the current execution state and that the required test infrastructure already exists. Earlier sub-plans may have changed interfaces or compile-time assertions in ways that leave the package temporarily uncompilable even though the task is still meant to follow structural TDD. Check for missing or stale mocks, fixtures, and test helpers as well. If the project uses generated mocks, confirm that the mock for the dependency under test exists and is up to date.

If either check fails, resolve the blocker within the TDD framework instead of skipping TDD or inverting the order. Acceptable scaffolding includes adding method stubs that panic or return zero values to satisfy interface assertions, adding the package to the mock generator configuration and regenerating mocks, or creating missing test helpers. These are temporary unblockers for the test author, not the task implementation itself. Record every scaffold created in the progress file so the implementer knows what must be replaced in step 4b.

**Isolation via isolated workspace**: The test author must NOT have access to plan files. Use the selected Worktrunk (`wt`) or plain `git worktree` workspace plus the active runtime adapter's dispatch mechanism to provide an isolated workspace containing the relevant code surface plus any explicitly seeded ignored build/cache artifacts, but not the planning artifacts.

If the task already has a task-scoped implementer worktree, run the test author in that worktree and keep it for implementation. If the task will implement in the main workspace, create a temporary test-author workspace and remove it after the test files are brought back.

**Spawn a test author sub-agent inside the isolated workspace** with a deliberately limited context:

**What the test author receives:**
- The acceptance criteria from the task (extracted and relayed as inline text by the executor — NOT as a file path)
- The relevant code surface (existing files, interfaces, types the tests will interact with — these are in the isolated workspace)

The test author's skills (`test-driven-development` + project-specific testing skills) must be loaded using the active runtime adapter's execution-binding mechanism. If the plan specifies a test author binding, use it. If not, the test author runs with the runtime's default sub-agent setup and the executor compensates by providing the minimum required task-local context.

**What the test author does NOT receive:**
- Design decisions, implementation context, integration contracts, or any other section of the task
- Plan file paths, plan directory paths, feature names, or any reference that could lead to the plan files
- The task file itself

**Two layers of isolation protect this separation:**
1. **Physical**: The isolated workspace contains only the code surface and explicitly seeded ignored build/cache artifacts the test author needs, not the planning artifacts.
2. **Contextual**: The test author's prompt contains no plan paths, directory names, feature names, or breadcrumbs that would lead it to look for the planning artifacts.

The test author writes tests grounded in the acceptance criteria and confirms they fail (RED). It returns the test file paths and a summary of what each test verifies.

**After the test author finishes**: If the test author ran in the task's implementer worktree, leave the tests there for the implementer. Otherwise, bring the test files back to the main execution workspace using the active runtime adapter's mechanism, then remove the temporary isolated workspace.

**Update progress**: Record test authoring as complete, note test file paths.

#### 4b. Implementation

**Without structural TDD** (test authoring was skipped due to testability): Spawn an implementer sub-agent in the chosen implementation workspace with the complete task as an inline task packet. No pre-written tests exist, so the immutability constraint doesn't apply. The implementer implements against the acceptance criteria directly. It returns implementation status and files created/modified. Existing tests must still pass — regressions are still caught in step 4c.

**With structural TDD** (tests were written in 4a): Spawn an **implementer** sub-agent in the chosen implementation workspace with full context:

**What the implementer receives:**
- The complete task as an inline task packet (all sections — design decisions, context, contracts, acceptance criteria), not a plan file path
- The tests written in step 4a (file paths)
- The project's required skills (if the plan specifies them for this task, or via the implementer's execution binding)
- Any outputs from prerequisite tasks (relayed by the executor)

**What the implementer is told:**
- Tests are immutable — do not modify test files
- Make the tests pass by implementing correctly
- Once tests are green, refactor the implementation for quality — clean up structure, naming, duplication. Tests are the safety net; refactor freely as long as they keep passing
- If a test cannot be satisfied, report it as a dispute with a clear explanation — do not modify the test, do not loop, do not skip it

The implementer implements the task and runs tests. It returns:
- Implementation status (complete, blocked, or disputes)
- Files created/modified
- Test results
- Any disputes (test file, test name, explanation of why it believes the test is wrong)

**Update progress**: Record implementation status, note created files, record any disputes.

#### 4c. Handle Results

**With structural TDD:**

- **All tests pass**: Mark task as `ready for integration/checkpoint`. Proceed according to the task worktree integration or sequential checkpointing rule below.
- **Tests fail, no disputes**: The implementer couldn't make tests pass but doesn't claim they're wrong. Mark as `blocked: implementation failure`. Continue with independent tasks.
- **Disputes reported**: Record each dispute in the progress file with the implementer's explanation. Mark the task as `blocked: test dispute`. Continue with independent tasks.
- **Existing tests regress**: Mark as `blocked: regression`. This takes priority — it means the implementation broke something outside its scope. Continue with independent tasks.

**Without structural TDD:**

- **Implementation complete, no regressions**: Mark task as `ready for integration/checkpoint`. Proceed according to the task worktree integration or sequential checkpointing rule below.
- **Existing tests regress**: Mark as `blocked: regression`. Same treatment as above — the implementation broke something outside its scope.

**Task worktree integration**: If the task ran in a task-scoped implementer worktree, inspect and integrate its result into the execution integration branch using the mechanism that matches the worktree tool before marking the task `done`: `wt merge` for Worktrunk worktrees, or explicit git merge/cherry-pick/patch transfer for plain `git worktree` worktrees. After successful integration and required verification, create a local checkpoint commit on the integration branch, record the commit in progress, then mark the task `done`. If integration conflicts, checkpoint commit creation fails, or post-merge verification fails, mark the task `blocked: integration` or `blocked: regression`, keep enough workspace state for diagnosis, and continue with independent tasks where possible. Remove the task worktree only after its result has been integrated and checkpointed or intentionally abandoned.

**Sequential task checkpointing**: If a sequential task runs directly in the integration branch workspace rather than a task worktree, create the same checkpoint commit after verification, record it in progress, then mark the task `done`. A task is not a completed prerequisite until its checkpoint commit exists.

### Step 5: Resolve Blocks

When the executor reaches a natural breakpoint — all independent work is done, or all remaining tasks depend on a blocked one — present all blocked items to the user at once:

```markdown
## Blocked Items

| Task | Status | Details |
|------|--------|---------|
| 02-api-layer | test dispute | test_returns_404: "AC says 'return error' — unclear if 404 or 500" |
| 04-integration | impl failure | Cannot satisfy test_timeout_recovery — needs async pattern not in scope |
```

The user resolves each item (fix the test, fix the AC, adjust the task, or accept as-is). The executor then re-runs the affected tasks from the appropriate step (test authoring if AC changed, implementation if tests changed).

### Step 6: Completion

When all tasks are `done`:

1. Run the full test suite on the integration branch HEAD to catch cross-task regressions
2. Record the final checkpoint range in progress (`<execution-base>..<integration-branch-head>`)
3. Convert the checkpointed result into an aggregate dirty review diff with `git reset --mixed <execution-base>` in the coordinator workspace
4. Update progress file to reflect completion and the final dirty review state
5. If an anchor is active, update its final active-work state without duplicating progress details, and note any durable decisions that should graduate to ADRs or permanent docs
6. Report summary to the user: what was built, any issues encountered, final test results, the checkpoint range that was reset for review, and whether any ADR/doc follow-up was identified

## Progress File Format

The progress file lives alongside the plan files: `<plan-directory>/progress.md`

A reference template is in this skill's `assets/` directory:

- **[Progress file template][progress-template]**

The executor updates this file after every meaningful step. It must be readable enough that a human (or a resumed executor) can understand exactly where execution stands.

## Integration with Execution Bindings

If the plan specifies runtime-specific execution bindings, the executor uses them:

- **Test author binding**: Should preload testing and code-writing skills. The executor dispatches this agent in an isolated workspace and passes only acceptance criteria as inline text.
- **Implementer binding**: Should preload the task's required skills (coding conventions, operational skills). The executor dispatches this agent with the full task and test file paths.

For tasks without testable AC (docs, config, file moves), only the implementer binding is needed.

If no execution bindings are specified, the executor may use the runtime's default sub-agent setup only for a single-sub-plan plan with no explicit model tier. For multi-sub-plan plans or plans with explicit model assignments, missing bindings are a blocker that must be resolved before implementation starts.

## Rules

- **Never let the test author see design decisions** — the structural separation is the entire point. If the test author's prompt accidentally includes plan context, the separation is broken.
- **Never claim structural TDD without enforced isolation** — if the active runtime cannot dispatch the test author inside the required Worktrunk (`wt`) or plain `git worktree` workspace, skip structural TDD only after the isolation path is verified unavailable or the user explicitly approves the skip. Prompt hygiene alone does not satisfy the physical isolation requirement.
- **Never let the implementer modify tests** — disputes are recorded and batched, not resolved by the implementer.
- **Never run concurrent implementers in the same workspace** — DAG independence is not workspace isolation. Use task-scoped implementer worktrees for concurrent work, or serialize the tasks.
- **Never launch dependent work from uncheckpointed dirty state** — if a task depends on prior sub-plan output, that output must exist in an integration-branch checkpoint commit before the dependent worktree starts.
- **Never push checkpoint commits** — checkpoint commits are local execution plumbing. The user's review surface is the final aggregate dirty diff after mixed reset, not the checkpoint history.
- **Never copy plan/progress artifacts into worker worktrees** — plan files, review files, and `progress.md` stay in the coordinator workspace. Workers receive task packets and prerequisite outputs through the executor.
- **Never assume fresh worktrees contain ignored build artifacts** — build-heavy projects need explicit hard-link seeding of ignored build/cache directories before isolated TDD or parallel implementation dispatch. If required seeding cannot be verified, ask before triggering expensive rebuilds.
- **Never self-execute assigned worker tasks** — if a task has an assigned worker or model tier, dispatch it through the runtime binding. If dispatch fails, diagnose and retry once, then stop and ask the user rather than doing the task in the coordinator context.
- **Always update progress after each step** — this is the checkpoint mechanism. If you skip an update and the session drops, work is lost.
- **Do not use anchors as progress files** — anchors preserve feature-level rationale and handoff context. Progress files preserve mechanical execution state.
- **Continue independent work when blocked** — don't stop the entire execution because one task has a dispute. If the plan specifies dependencies, use them to determine what can proceed. If no dependencies are specified, treat remaining tasks as sequential and pause at the blocked one.
- **Relay prerequisite outputs between dependent tasks** — sub-agents cannot communicate with each other. The executor passes results between tasks when needed.
- **Respect model assignments** — if the plan specifies a model tier for a task, use it. Don't silently upgrade or downgrade.
- **Record testability gate decisions but don't analyze** — when TDD is skipped, note the reason briefly in progress (e.g., "declared untestable in AGENTS.md", "no testable seams"). Don't spend tokens analyzing what would need to change — that's a separate effort if the user chooses to pursue it.

[progress-template]: assets/progress-template.md
