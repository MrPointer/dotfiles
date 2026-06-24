# Runtime Adapter: Codex

Use this adapter only when the active runtime is Codex.

This adapter maps the canonical execution workflow in `../SKILL.md` to Codex-native mechanics. It must not redefine task ordering, progress rules, or implementer-owned test discipline.

## Contents

- [Exploration and Dispatch](#exploration-and-dispatch)
- [Execution Bindings](#execution-bindings)
- [Workspace Isolation Strategy](#workspace-isolation-strategy)
- [Integration Branch And Checkpoints](#integration-branch-and-checkpoints)
- [Implementer Worktree Isolation](#implementer-worktree-isolation)
- [Implementer Dispatch](#implementer-dispatch)
- [Progress and Artifacts](#progress-and-artifacts)
- [Model Assignment](#model-assignment)

## Exploration and Dispatch

- Use Codex sub-agent dispatch for exploration and implementation work when sub-agents materially help.
- Set the sub-agent `model` explicitly using Codex's actual dispatch mechanism when the plan assigns a tier.
- Give implementers complete task packets, including testing expectations when the task is testable.
- Before executing a multi-sub-plan plan, verify every assigned dispatch recipe or worker can be invoked. If a binding fails, diagnose and retry once. If it still fails, stop and ask the user; do not perform assigned implementation in the coordinator context.

## Execution Bindings

- In Codex, execution bindings should use the built-in worker path for implementation work unless the project explicitly defines a more specialized worker.
- In Codex, execution bindings are usually reusable dispatch recipes rather than checked-in worker files.
- Project-local Codex skills are expected under `.agents/skills/` in the project root. Global Codex skills are expected under `~/.agents/skills/`.
- Each binding should make explicit:
  - target `model`
  - required skills to attach explicitly as Codex `skill` items from the local or global Codex skills directories
  - whether the worker should write files directly or return results for the parent to persist
- Pass model explicitly in the worker dispatch and attach required skills explicitly as `skill` items.
- Do not rely on prompt text alone for skill loading when the runtime can attach the skills directly.
- If the binding is ephemeral rather than file-backed, record its parameters in the plan metadata or execution context so retries and resumed execution reuse the same model and skills.

## Workspace Isolation Strategy

Use one ordered fallback chain for task-scoped implementer worktrees:

1. Use Worktrunk (`wt`) if it is installed and suitable.
2. Otherwise use `git worktree` directly.

Do not use runtime-native worktree switching for isolated workspace creation or switching. Use Codex-native mechanics to dispatch workers into the selected Worktrunk or git worktree and to verify the worker actually ran there.

Before creating or entering an isolated workspace, apply the dirty-state preflight and build/cache reuse rules in [workspace-isolation.md](workspace-isolation.md).

After creating or entering an isolated workspace, verify its initial status according to [workspace-isolation.md](workspace-isolation.md) before dispatching a worker.

## Integration Branch And Checkpoints

- Create or resume the canonical execution integration branch from the recorded execution base before dispatching implementation workers.
- Create task worktrees from the integration branch checkpoint that contains the task's prerequisites.
- After integrating a task result into the integration branch and passing verification, create a local checkpoint commit using the project's normal signing policy.
- Do not push checkpoint commits. At completion, leave the aggregate result for review by mixed-resetting the integration branch to the execution base.

## Implementer Worktree Isolation

- For implementation tasks that run concurrently with any other implementation task, create or enter a task-scoped worktree using the Workspace Isolation Strategy.
- Before dispatch, apply only the project-documented build/cache reuse strategy allowed by [workspace-isolation.md](workspace-isolation.md).
- Dispatch the assigned implementer worker inside that worktree using Codex's actual worker dispatch mechanism and explicit model/skill binding.
- Do not copy plan files, review files, or `progress.md` into the worktree. The parent executor passes the full sub-plan content, prerequisite outputs, and testing expectations as inline task context.
- After the implementer finishes, inspect the task worktree diff and integrate it into the execution integration branch using the mechanism that matches how the worktree was created: `wt merge` from Worktrunk (`wt`), or explicit git merge/cherry-pick/patch transfer for plain `git worktree`.
- If integration conflicts, checkpoint commit creation fails, or verification fails after integration, record the task as blocked and keep enough worktree state for diagnosis. Remove the worktree only after the result is integrated and checkpointed or intentionally abandoned.

## Implementer Dispatch

- Spawn a separate implementer sub-agent with the full task packet, prerequisite outputs, and required skills/reference material.
- For testable behavior changes, tell the implementer to follow the required testing/TDD skills, make a test-first attempt when practical, and report tests added or updated plus verification results. If test-first work or new tests are not practical, it must explain why.
- If the implementer reports a blocker, record it and continue according to the canonical workflow.

## Progress and Artifacts

- The parent executor owns `progress.md` and should update it after each meaningful step.
- Even if a sub-agent writes files directly, the parent remains responsible for checkpointing and artifact verification.
- Record dispatch evidence in progress: planned worker, actual worker, model/effort, runtime dispatch mechanism, implementation workspace path, dirty-state preflight result, build/cache reuse status, checkpoint commit, integration status, test evidence, and verification status.

## Model Assignment

- Use Codex's explicit model-selection mechanism in the worker dispatch path; do not rely on prompt wording.
- Treat plan-assigned tiers as binding.
- If the required model is unavailable, stop and ask the user how to proceed.
