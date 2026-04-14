---
name: executing-plans
description: Use when executing implementation plans. Orchestrates task execution with testability-gated TDD, progress checkpointing, and dispute batching. Works with any plan structure that has identifiable tasks and acceptance criteria. Handles resumption after interruptions.
---

# Executing Plans

Orchestrate the execution of implementation plans. Read the plan, identify tasks and their ordering, maintain progress checkpoints, and handle interruptions gracefully.

When the code surface supports it, this skill structurally separates test authoring from implementation — a separate agent writes tests from acceptance criteria alone, without seeing design decisions. The implementer then receives the full task plus the tests. This prevents the common failure mode where agents write tests that rationalize their implementation rather than verify requirements. When the code isn't ready for TDD (missing test infrastructure, no testable seams, tight coupling), the executor skips structural separation rather than fighting the codebase.

## Core Principles

1. **Tests and Implementation Are Structurally Separated — When Feasible**: The test author sees only acceptance criteria. The implementer sees the full task context plus the tests. Neither can influence the other's context. However, this separation requires a testable code surface. When the target code lacks the infrastructure or seams to support isolated testing, the executor skips structural TDD rather than producing tests that fight the codebase.
2. **Progress Is Always Persisted**: After every meaningful step, update the progress file. If the session drops, the executor resumes from the last checkpoint — not from scratch.
3. **Tests Are Immutable to the Implementer**: The implementer cannot modify tests. If it believes a test is wrong, it reports this and moves on. Disputes are batched for human resolution.
4. **Independent Work Continues**: When a task is blocked (test dispute, failure), the executor continues with independent tasks.

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
- Worker agent assignments (if specified)
- File ownership per task (if specified — helps prevent conflicts during parallel execution)

### Step 2: Initialize or Resume Progress

Check for an existing progress file alongside the plan (`progress.md` in the plan directory, or next to a single-file plan).

**If progress exists**: Read it, determine current state, and resume from the last checkpoint. Verify that completed tasks' artifacts still exist (files, tests). If something is missing, mark that task for re-execution.

**If no progress exists**: Create the initial progress file from the plan. See [Progress File Format](#progress-file-format).

### Step 3: Execute Tasks

For each task in execution order (respecting dependencies):

#### 3a. Test Authoring

**When to skip (no testable work)**: Some tasks don't have testable work — documentation updates, file moves, configuration changes, or tasks with no acceptance criteria. Skip test authoring for these and proceed directly to implementation.

**When to skip (code surface not testable)**: For tasks that have testable acceptance criteria, check whether the code surface supports TDD before spawning the test author. This is a gate, not an analysis — the goal is a binary decision, not a report on what needs fixing.

1. **Check project docs first**: Look for explicit testability statements in `AGENTS.md`, project documentation, or task-level notes. If the project or area is declared untestable, skip immediately — no exploration needed.

2. **Lightweight code surface check**: If docs don't cover it, spawn an exploration agent at the cheapest available model tier to answer one question: can the components this task touches be tested in isolation? The agent looks for:
   - Interfaces or injectable dependencies (substitution points for test doubles)
   - An existing test framework and test patterns in the area
   - Whether tests could target the code without restructuring it first

   The agent returns a yes/no — not an analysis of what's wrong or recommendations for making it testable.

If the code surface isn't ready for TDD, skip test authoring and proceed directly to implementation without structural separation. Record the decision briefly in the progress file (Tests column: `skipped`, Notes: e.g., "declared untestable in AGENTS.md" or "no testable seams"). The implementer still receives the full task with acceptance criteria — it just doesn't receive pre-written tests and isn't bound by the test immutability constraint.

**Isolation via worktree**: The test author must NOT have access to plan files. Since plans live outside the tracked codebase (under `plans/`, typically gitignored), spawning the test author in a worktree achieves physical isolation — the worktree contains the source code but not the plan directory.

Create the worktree using whichever mechanism is available, in priority order:
1. **Worktrunk** (`wt switch <branch-name>`) — if the `worktrunk` skill or `wt` CLI is available. Preferred because it handles setup hooks and cleanup.
2. **Agent-native worktrees** (`EnterWorktree`) — if the agent framework provides worktree tooling.
3. **Git CLI** (`git worktree add`) — always available as a fallback.

The worktree is temporary — remove it after the test author finishes. The test files it writes exist on a branch that gets merged or cherry-picked back.

**Spawn a test author sub-agent inside the worktree** with a deliberately limited context:

**What the test author receives:**
- The acceptance criteria from the task (extracted and relayed as inline text by the executor — NOT as a file path)
- The relevant code surface (existing files, interfaces, types the tests will interact with — these are in the worktree)

The test author's skills (`test-driven-development` + project-specific testing skills) must be preloaded via a worker agent definition — they cannot be passed dynamically at spawn time. If the plan specifies a test author worker agent, use it. If not, the test author runs without preloaded skills and relies on whatever the default agent has access to.

**What the test author does NOT receive:**
- Design decisions, implementation context, integration contracts, or any other section of the task
- Plan file paths, plan directory paths, feature names, or any reference that could lead to the plan files
- The task file itself

**Two layers of isolation protect this separation:**
1. **Physical**: The worktree contains only tracked source code. If plans are gitignored (the default), they don't exist in the worktree.
2. **Contextual**: Even if plans are tracked and present in the worktree, the test author's prompt contains no plan paths, directory names, or breadcrumbs that would lead it to look for them. An agent won't search for files it doesn't know exist.

The test author writes tests grounded in the acceptance criteria and confirms they fail (RED). It returns the test file paths and a summary of what each test verifies.

**After the test author finishes**: Bring the test files back to the main working tree (merge the worktree branch or cherry-pick the commit), then remove the worktree.

**Update progress**: Record test authoring as complete, note test file paths.

#### 3b. Implementation

**Without structural TDD** (test authoring was skipped due to testability): Spawn an implementer sub-agent with the complete task. No pre-written tests exist, so the immutability constraint doesn't apply. The implementer implements against the acceptance criteria directly. It returns implementation status and files created/modified. Existing tests must still pass — regressions are still caught in step 3c.

**With structural TDD** (tests were written in 3a): Spawn an **implementer** sub-agent with full context:

**What the implementer receives:**
- The complete task (all sections — design decisions, context, contracts, acceptance criteria)
- The tests written in step 3a (file paths)
- The project's required skills (if the plan specifies them for this task, or via the implementer's worker agent definition)
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

#### 3c. Handle Results

**With structural TDD:**

- **All tests pass**: Mark task as `done`. Proceed to the next task.
- **Tests fail, no disputes**: The implementer couldn't make tests pass but doesn't claim they're wrong. Mark as `blocked: implementation failure`. Continue with independent tasks.
- **Disputes reported**: Record each dispute in the progress file with the implementer's explanation. Mark the task as `blocked: test dispute`. Continue with independent tasks.
- **Existing tests regress**: Mark as `blocked: regression`. This takes priority — it means the implementation broke something outside its scope. Continue with independent tasks.

**Without structural TDD:**

- **Implementation complete, no regressions**: Mark task as `done`. Proceed to the next task.
- **Existing tests regress**: Mark as `blocked: regression`. Same treatment as above — the implementation broke something outside its scope.

### Step 4: Resolve Blocks

When the executor reaches a natural breakpoint — all independent work is done, or all remaining tasks depend on a blocked one — present all blocked items to the user at once:

```markdown
## Blocked Items

| Task | Status | Details |
|------|--------|---------|
| 02-api-layer | test dispute | test_returns_404: "AC says 'return error' — unclear if 404 or 500" |
| 04-integration | impl failure | Cannot satisfy test_timeout_recovery — needs async pattern not in scope |
```

The user resolves each item (fix the test, fix the AC, adjust the task, or accept as-is). The executor then re-runs the affected tasks from the appropriate step (test authoring if AC changed, implementation if tests changed).

### Step 5: Completion

When all tasks are `done`:

1. Run the full test suite to catch cross-task regressions
2. Update progress file to reflect completion
3. Report summary to the user: what was built, any issues encountered, final test results

## Progress File Format

The progress file lives alongside the plan files: `<plan-directory>/progress.md`

A reference template is in this skill's `references/` directory:

- **[Progress file template][progress-template]**

The executor updates this file after every meaningful step. It must be readable enough that a human (or a resumed executor) can understand exactly where execution stands.

## Integration with Worker Agents

If the plan specifies worker agents, the executor uses them:

- **Test author worker**: Should preload testing and code-writing skills. The executor spawns this agent in a worktree and passes only acceptance criteria as inline text.
- **Implementer worker**: Should preload the task's required skills (coding conventions, operational skills). The executor spawns this agent with the full task and test file paths.

For tasks without testable AC (docs, config, file moves), only the implementer worker is needed.

If no worker agents are specified, the executor spawns sub-agents with the project's default model. The structural isolation (worktree + prompt hygiene) still applies regardless, but agents won't have skills preloaded.

## Rules

- **Never let the test author see design decisions** — the structural separation is the entire point. If the test author's prompt accidentally includes plan context, the separation is broken.
- **Never let the implementer modify tests** — disputes are recorded and batched, not resolved by the implementer.
- **Always update progress after each step** — this is the checkpoint mechanism. If you skip an update and the session drops, work is lost.
- **Continue independent work when blocked** — don't stop the entire execution because one task has a dispute. If the plan specifies dependencies, use them to determine what can proceed. If no dependencies are specified, treat remaining tasks as sequential and pause at the blocked one.
- **Relay prerequisite outputs between dependent tasks** — sub-agents cannot communicate with each other. The executor passes results between tasks when needed.
- **Respect model assignments** — if the plan specifies a model tier for a task, use it. Don't silently upgrade or downgrade.
- **Record testability gate decisions but don't analyze** — when TDD is skipped, note the reason briefly in progress (e.g., "declared untestable in AGENTS.md", "no testable seams"). Don't spend tokens analyzing what would need to change — that's a separate effort if the user chooses to pursue it.

[progress-template]: references/progress-template.md
