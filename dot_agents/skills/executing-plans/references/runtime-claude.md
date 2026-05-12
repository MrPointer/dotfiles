# Runtime Adapter: Claude

Use this adapter only when the active runtime is Claude.

This adapter maps the canonical execution workflow in `../SKILL.md` to Claude-native mechanics. It must not redefine task ordering, dispute policy, progress rules, or the test-author/implementer separation.

## Exploration and Dispatch

- Use Claude sub-agents (launched via the `Agent` tool) for exploration, test authoring, and implementation work.
- Respect the model tier assigned by the plan — use the worker agent's `model` field, not prompt wording, to control model selection.
- For testability gate exploration, prefer `subagent_type: "Explore"` at a cheap model tier.
- Keep prompts narrow for test authors and complete for implementers, matching the canonical workflow.
- Before executing a multi-sub-plan plan, verify every assigned worker agent is discoverable to the current session. If a worker is missing or cannot be launched, diagnose and retry once when the cause is mechanical. If it still fails, stop and ask the user; do not perform assigned implementation in the coordinator context.

## Execution Bindings

Claude execution bindings are **file-defined worker agents** under `.claude/agents/` (project-local) or `~/.claude/agents/` (global). The `model` and `skills` frontmatter fields are the reliable mechanism for controlling model selection and skill preload.

- Reuse existing worker agents when they match the sub-plan's needs.
- If a worker is missing, the planning skill should have created it. If it was missed, create one in the project-local agent directory following the frontmatter rules documented in the planning skill's Claude adapter.
- Warn the user if new workers were just created — a session restart is required before they become discoverable.

## Workspace Isolation Strategy

Use one ordered fallback chain for both structural TDD workspaces and task-scoped implementer worktrees:

1. **`EnterWorktree` / `ExitWorktree`** — Claude Code's native worktree tools. Prefer these when available because they are runtime-aware and provide verifiable entry, exit, result collection, and cleanup.
2. **Worktrunk (`wt`)** — if the `worktrunk:worktrunk` skill is installed, load it and use `wt switch <branch-name>` when native Claude worktree tooling is unavailable or unsuitable.
3. **Git CLI** (`git worktree add`) — always available as a fallback.

## Test Author Isolation

**Physical isolation**: The plan directory lives outside the tracked codebase (under `plans/`, typically gitignored). A fresh worktree therefore contains the source code but not the plan files — an agent cannot read what isn't present. This is the foundation of physical isolation; do not undermine it by copying plan files, review files, or `progress.md` into the worktree.

**Contextual isolation**: Even with physical isolation, the test author's prompt must not reveal the plan path, task file path, feature name, or design rationale. Pass only acceptance criteria (inline as text, not as a file path) and the code surface the tests interact with.

**Structural TDD gate**: If Claude cannot create the isolated worktree (e.g., the repository has uncommitted changes that block worktree creation and the user declines to resolve them), skip structural TDD only after recording the attempted mechanism and failure. If an isolation mechanism exists but dispatch into it fails, stop and ask whether to fix isolation or explicitly skip structural TDD.

**Bringing test files back**: If the test author ran in a task-scoped implementer worktree, leave the test files there for the implementer. Otherwise, after the test author finishes, return the test files to the main execution workspace:

- If the worktree was created via `EnterWorktree`, use `ExitWorktree` to merge the changes back.
- If created via Worktrunk (`wt`) or `git worktree add`, merge the worktree's branch into the main branch or cherry-pick the commit containing the tests.
- Remove the worktree once the test files are back.

## Implementer Dispatch

- For implementation tasks that run concurrently with any other implementation task, create or enter a task-scoped worktree using the Workspace Isolation Strategy.
- Dispatch a separate implementer worker in the chosen workspace with the full task packet, test file paths, prerequisite outputs, and required skills. Pass the task content inline; do not rely on plan file paths inside worker worktrees.
- Tell the implementer explicitly that tests are immutable.
- If the implementer reports a dispute, record it in progress and continue with independent tasks per the canonical workflow.
- After a task-scoped worktree worker finishes, inspect and integrate the result using `ExitWorktree` when the native tool created it, `wt merge` from Worktrunk (`wt`) when that tool created it, or git merge/cherry-pick/patch transfer when plain git created it. Remove the worktree only after the result is integrated or intentionally abandoned.

## Progress and Artifacts

- The parent executor owns `progress.md` and updates it after each meaningful step.
- Even if a worker writes files directly, the parent remains responsible for checkpointing and verifying that expected artifacts exist.
- Record dispatch evidence in progress: planned worker, actual worker, model/effort, runtime dispatch mechanism, implementation workspace path, integration status, and TDD isolation outcome.

## Model Assignment

- Use the worker agent's explicit `model` field as the source of truth.
- Treat plan-assigned tiers as binding. Do not silently upgrade or downgrade.
- If a worker fails at its assigned model, diagnose and fix (permission mode, tool access, model availability). If it cannot be resolved, stop and ask the user — never fall back to a more expensive model silently.
