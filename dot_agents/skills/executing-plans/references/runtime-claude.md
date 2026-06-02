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

1. **Worktrunk (`wt`)** — if Worktrunk is installed, load the `worktrunk:worktrunk` skill when available and use `wt switch <branch-name>`.
2. **Git CLI** (`git worktree add`) — use plain git as the fallback.

Do not use `EnterWorktree` / `ExitWorktree` or other runtime-native worktree switching for isolated workspace creation or switching. Use Claude-native mechanics to dispatch workers into the selected Worktrunk or git worktree and to verify the worker actually ran there.

After creating or entering an isolated workspace, seed ignored build/cache artifacts when the project needs them for practical compile or test performance. Identify the relative cache directories from project docs or config, such as Rust `target/`. When the source cache exists in the coordinator workspace and the paths are on the same filesystem, create the destination directory and hard-link copy the contents, for example with `cp -al <source-dir>/. <worktree-dir>/<relative-dir>/` when supported. Verify the expected files exist in the worktree. If hard-link seeding is required for a build-heavy project but cannot be verified, record the failure and ask before dispatching the worker.

## Integration Branch And Checkpoints

- Create or resume the canonical execution integration branch from the recorded execution base before dispatching implementation workers.
- Create task worktrees from the integration branch checkpoint that contains the task's prerequisites.
- After integrating a task result into the integration branch and passing verification, create a local checkpoint commit using the project's normal signing policy.
- Do not push checkpoint commits. At completion, leave the aggregate result for review by mixed-resetting the integration branch to the execution base.

## Test Author Isolation

**Physical isolation**: The plan directory lives outside the tracked codebase (under `plans/`, typically gitignored). A fresh worktree therefore contains the source code but not the plan files; build-heavy projects may also hard-link seed ignored build/cache artifacts. This is the foundation of physical isolation; do not undermine it by copying plan files, review files, or `progress.md` into the worktree.

**Contextual isolation**: Even with physical isolation, the test author's prompt must not reveal the plan path, task file path, feature name, or design rationale. Pass only acceptance criteria (inline as text, not as a file path) and the code surface the tests interact with.

**Structural TDD gate**: If neither Worktrunk nor plain git can create the isolated worktree (e.g., the repository has uncommitted changes that block worktree creation and the user declines to resolve them), skip structural TDD only after recording the attempted mechanism and failure. If an isolation mechanism exists but dispatch into it fails, stop and ask whether to fix isolation or explicitly skip structural TDD.

**Bringing test files back**: If the test author ran in a task-scoped implementer worktree, leave the test files there for the implementer. Otherwise, after the test author finishes, return the test files to the main execution workspace:

- If the worktree was created via Worktrunk (`wt`), prefer `wt merge`.
- If created via `git worktree add`, merge the worktree's branch into the integration branch or cherry-pick the commit containing the tests.
- Remove the worktree once the test files are back.

## Implementer Dispatch

- For implementation tasks that run concurrently with any other implementation task, create or enter a task-scoped worktree using the Workspace Isolation Strategy.
- Before dispatch, hard-link seed required ignored build/cache artifact directories into the task worktree when the project is build-heavy.
- Dispatch a separate implementer worker in the chosen workspace with the full task packet, test file paths, prerequisite outputs, and required skills. Pass the task content inline; do not rely on plan file paths inside worker worktrees.
- Tell the implementer explicitly that tests are immutable.
- If the implementer reports a dispute, record it in progress and continue with independent tasks per the canonical workflow.
- After a task-scoped worktree worker finishes, inspect and integrate the result into the integration branch using `wt merge` from Worktrunk (`wt`) when that tool created it, or git merge/cherry-pick/patch transfer when plain git created it. Create the checkpoint commit after verification. Remove the worktree only after the result is integrated and checkpointed or intentionally abandoned.

## Progress and Artifacts

- The parent executor owns `progress.md` and updates it after each meaningful step.
- Even if a worker writes files directly, the parent remains responsible for checkpointing and verifying that expected artifacts exist.
- Record dispatch evidence in progress: planned worker, actual worker, model/effort, runtime dispatch mechanism, implementation workspace path, build/cache seeding status, checkpoint commit, integration status, and TDD isolation outcome.

## Model Assignment

- Use the worker agent's explicit `model` field as the source of truth.
- Treat plan-assigned tiers as binding. Do not silently upgrade or downgrade.
- If a worker fails at its assigned model, diagnose and fix (permission mode, tool access, model availability). If it cannot be resolved, stop and ask the user — never fall back to a more expensive model silently.
