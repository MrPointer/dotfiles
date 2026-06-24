# Runtime Adapter: OpenCode

Use this adapter only when the active runtime is OpenCode.

This adapter maps the canonical execution workflow in `../SKILL.md` to OpenCode-native mechanics. It must not redefine task ordering, progress rules, or implementer-owned test discipline.

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

- Use OpenCode subagents for exploration and implementation work when subagents materially help. OpenCode subagents can be invoked through native subagent dispatch, including `@<subagent-name>` mentions in an interactive session, and by primary agents through the Task tool when permitted.
- Prefer the project-provided `semble-search` subagent for cheap, read-only codebase exploration when available; otherwise use the built-in `explore` subagent.
- Prefer a custom implementer subagent when the plan assigns a specific model tier or required skills; the built-in `general` subagent is acceptable only when no explicit binding is required.
- Give implementers complete task packets, including testing expectations when the task is testable.

For same-workspace assigned workers, prefer native OpenCode subagent dispatch when it can invoke the named custom subagent and preserve the assigned model/permissions. In interactive sessions, mentioning `@<worker-name>` is a valid native invocation path. If the current runtime surface cannot invoke project-local custom subagents directly, or when explicit workspace routing is required, the CLI dispatch path is an acceptable fallback:

```bash
opencode run --agent <worker-name> --dir <workspace-path> "<task prompt>"
```

Before executing a multi-sub-plan plan, verify that every assigned worker name appears in `opencode agent list` or the runtime's native subagent picker. If native dispatch of an assigned worker fails, diagnose and retry once. If native dispatch is unavailable in the current surface, use `opencode run --agent ...` as a fallback. If no dispatch path works, stop and ask the user; do not perform the assigned implementation in the coordinator context.

## Execution Bindings

OpenCode execution bindings should be **markdown-defined custom subagents** under `.opencode/agents/` (project-local) or `~/.config/opencode/agents/` (global).

- Reuse existing custom subagents when they match the task's model tier, permissions, and required skills.
- Prefer project-local bindings when the task depends on project-local skills or conventions.
- Set `mode: subagent` and set `model` explicitly. In OpenCode, subagents inherit the caller's model when `model` is omitted, so omission is not a reliable execution binding.
- Configure access with `permission`, not the deprecated `tools` field, for new or updated bindings.
- If the binding must use named skills, allow them via `permission.skill` and instruct the subagent in its prompt to load them with the `skill` tool at startup. OpenCode does not provide a declarative skill-preload field in agent frontmatter.
- Keep bindings thin: the task packet provides task truth; the binding provides runtime-native identity, model, permissions, and skill-loading instructions.
- If the binding should not delegate further, deny or tightly scope `permission.task` so the worker does not spawn unrelated child agents.

## Workspace Isolation Strategy

Use one ordered fallback chain for task-scoped implementer worktrees:

1. Use Worktrunk (`wt`) if it is installed and suitable.
2. Otherwise use `git worktree` via Bash to create an isolated workspace.

Do not use runtime-native worktree switching for isolated workspace creation or switching. Use OpenCode-native mechanics to dispatch workers into the selected Worktrunk or git worktree and to verify the worker actually ran there.

Before creating or entering an isolated workspace, apply the dirty-state preflight and build/cache reuse rules in [workspace-isolation.md](workspace-isolation.md).

After creating or entering an isolated workspace, verify its initial status according to [workspace-isolation.md](workspace-isolation.md) before dispatching a worker.

When Worktrunk (`wt`) is used, create or enter the isolated workspace with `wt switch`. Per Worktrunk (`wt`) documentation, `wt switch` is the command that switches to a worktree and creates one if needed; use `wt switch --create <branch>` when the isolation branch does not exist yet. Its `--execute` mode can also be useful when you need to launch the agent directly inside the isolated worktree.

Do not assume implementer isolation is possible just because subagents exist. The executor must verify that the worker actually runs in the isolated worktree rather than the main workspace.

## Integration Branch And Checkpoints

- Create or resume the canonical execution integration branch from the recorded execution base before dispatching implementation workers.
- Create task worktrees from the integration branch checkpoint that contains the task's prerequisites.
- After integrating a task result into the integration branch and passing verification, create a local checkpoint commit using the project's normal signing policy.
- Do not push checkpoint commits. At completion, leave the aggregate result for review by mixed-resetting the integration branch to the execution base.

## Implementer Worktree Isolation

- For implementation tasks that run concurrently with any other implementation task, create or enter a task-scoped worktree using the Workspace Isolation Strategy.
- Before dispatch, apply only the project-documented build/cache reuse strategy allowed by [workspace-isolation.md](workspace-isolation.md).
- Dispatch the implementer with `opencode run --agent <implementer-worker> --dir <task-worktree-path> "<full task packet>"` unless a native OpenCode dispatch path can target that worktree and preserve the worker's configured model and permissions.
- Do not copy plan files, review files, or `progress.md` into the task worktree. The parent executor passes the full sub-plan content, prerequisite outputs, and testing expectations as inline task context.
- After the implementer finishes, inspect the task worktree diff and integrate it into the execution integration branch using the mechanism that matches how the worktree was created: `wt merge` from Worktrunk (`wt`), or explicit git merge/cherry-pick/patch transfer for plain `git worktree`.
- If integration conflicts, checkpoint commit creation fails, or verification fails after integration, record the task as blocked and keep enough worktree state for diagnosis. Remove the worktree only after the result is integrated and checkpointed or intentionally abandoned.

## Implementer Dispatch

- Dispatch a separate implementer subagent with the full task packet, prerequisite outputs, and required skills.
- Use native OpenCode subagent dispatch for assigned implementer workers when available. Use `opencode run --agent <implementer-worker> --dir <workspace-path> "<full task packet>"` when the current runtime surface cannot invoke the custom worker directly or when explicit workspace routing is needed.
- For testable behavior changes, tell the implementer to follow the required testing/TDD skills, make a test-first attempt when practical, and report tests added or updated plus verification results. If test-first work or new tests are not practical, it must explain why.
- If the implementer reports a blocker, record it in progress and continue with independent tasks per the canonical workflow.

## Progress and Artifacts

- The parent executor owns `progress.md` and updates it after each meaningful step.
- Even if a subagent writes files directly, the parent remains responsible for checkpointing and verifying that expected artifacts exist.
- Record dispatch evidence in progress: planned worker, actual worker, model/effort from the worker binding, command or runtime mechanism used, implementation workspace path, dirty-state preflight result, build/cache reuse status, checkpoint commit, integration status, test evidence, and verification status.

## Model Assignment

- Use the custom subagent's explicit `model` field as the source of truth.
- Treat plan-assigned tiers as binding. Do not silently upgrade or downgrade.
- Because built-in subagents and custom subagents without `model` inherit from the caller in OpenCode, inherited model selection is not sufficient when the plan made an explicit model decision.
- If the required model is unavailable, stop and ask the user how to proceed.
