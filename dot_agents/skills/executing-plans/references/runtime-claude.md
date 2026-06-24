# Runtime Adapter: Claude

Use this adapter only when the active runtime is Claude.

This adapter maps the canonical execution workflow in `../SKILL.md` to Claude-native mechanics. It must not redefine task ordering, progress rules, or implementer-owned test discipline.

## Contents

- [Exploration and Dispatch](#exploration-and-dispatch)
- [Execution Bindings](#execution-bindings)
- [Workspace Isolation Strategy](#workspace-isolation-strategy)
- [Integration Branch And Checkpoints](#integration-branch-and-checkpoints)
- [Implementer Dispatch](#implementer-dispatch)
- [Progress and Artifacts](#progress-and-artifacts)
- [Model Assignment](#model-assignment)

## Exploration and Dispatch

- Use Claude sub-agents (launched via the `Agent` tool) for exploration and implementation work.
- Respect the model tier assigned by the plan — use the worker agent's `model` field, not prompt wording, to control model selection.
- Keep prompts complete for implementers, including testing expectations when the task is testable.
- Before executing a multi-sub-plan plan, verify every assigned worker agent is discoverable to the current session. If a worker is missing or cannot be launched, diagnose and retry once when the cause is mechanical. If it still fails, stop and ask the user; do not perform assigned implementation in the coordinator context.

## Execution Bindings

Claude execution bindings are **file-defined worker agents** under `.claude/agents/` (project-local) or `~/.claude/agents/` (global). The `model` and `skills` frontmatter fields are the reliable mechanism for controlling model selection and skill preload.

- Reuse existing worker agents when they match the sub-plan's needs.
- If a worker is missing, the planning skill should have created it. If it was missed, create one in the project-local agent directory following the frontmatter rules documented in the planning skill's Claude adapter.
- Warn the user if new workers were just created — a session restart is required before they become discoverable.

## Workspace Isolation Strategy

Use one ordered fallback chain for task-scoped implementer worktrees:

1. **Worktrunk (`wt`)** — if Worktrunk is installed, load the `worktrunk:worktrunk` skill when available and use `wt switch <branch-name>`.
2. **Git CLI** (`git worktree add`) — use plain git as the fallback.

Do not use `EnterWorktree` / `ExitWorktree` or other runtime-native worktree switching for isolated workspace creation or switching. Use Claude-native mechanics to dispatch workers into the selected Worktrunk or git worktree and to verify the worker actually ran there.

Before creating or entering an isolated workspace, apply the dirty-state preflight and build/cache reuse rules in [workspace-isolation.md](workspace-isolation.md).

After creating or entering an isolated workspace, verify its initial status according to [workspace-isolation.md](workspace-isolation.md) before dispatching a worker.

## Integration Branch And Checkpoints

- Create or resume the canonical execution integration branch from the recorded execution base before dispatching implementation workers.
- Create task worktrees from the integration branch checkpoint that contains the task's prerequisites.
- After integrating a task result into the integration branch and passing verification, create a local checkpoint commit using the project's normal signing policy.
- Do not push checkpoint commits. At completion, leave the aggregate result for review by mixed-resetting the integration branch to the execution base.

## Implementer Dispatch

- For implementation tasks that run concurrently with any other implementation task, create or enter a task-scoped worktree using the Workspace Isolation Strategy.
- Before dispatch, apply only the project-documented build/cache reuse strategy allowed by [workspace-isolation.md](workspace-isolation.md).
- Dispatch a separate implementer worker in the chosen workspace with the full task packet, prerequisite outputs, and required skills. Pass the task content inline; do not rely on plan file paths inside worker worktrees.
- For testable behavior changes, tell the implementer to follow the required testing/TDD skills, make a test-first attempt when practical, and report tests added or updated plus verification results. If test-first work or new tests are not practical, it must explain why.
- If the implementer reports a blocker, record it in progress and continue with independent tasks per the canonical workflow.
- After a task-scoped worktree worker finishes, inspect and integrate the result into the integration branch using `wt merge` from Worktrunk (`wt`) when that tool created it, or git merge/cherry-pick/patch transfer when plain git created it. Create the checkpoint commit after verification. Remove the worktree only after the result is integrated and checkpointed or intentionally abandoned.

## Progress and Artifacts

- The parent executor owns `progress.md` and updates it after each meaningful step.
- Even if a worker writes files directly, the parent remains responsible for checkpointing and verifying that expected artifacts exist.
- Record dispatch evidence in progress: planned worker, actual worker, model/effort, runtime dispatch mechanism, implementation workspace path, dirty-state preflight result, build/cache reuse status, checkpoint commit, integration status, test evidence, and verification status.

## Model Assignment

- Use the worker agent's explicit `model` field as the source of truth.
- Treat plan-assigned tiers as binding. Do not silently upgrade or downgrade.
- If a worker fails at its assigned model, diagnose and fix (permission mode, tool access, model availability). If it cannot be resolved, stop and ask the user — never fall back to a more expensive model silently.
