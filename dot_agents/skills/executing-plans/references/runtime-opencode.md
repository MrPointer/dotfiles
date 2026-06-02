# Runtime Adapter: OpenCode

Use this adapter only when the active runtime is OpenCode.

This adapter maps the canonical execution workflow in `../SKILL.md` to OpenCode-native mechanics. It must not redefine task ordering, dispute policy, progress rules, or the test-author/implementer separation.

## Exploration and Dispatch

- Use OpenCode subagents for exploration, test authoring, and implementation work when subagents materially help. OpenCode subagents can be invoked through native subagent dispatch, including `@<subagent-name>` mentions in an interactive session, and by primary agents through the Task tool when permitted.
- Prefer the built-in `explore` subagent for cheap, read-only codebase exploration such as the testability gate.
- Prefer a custom implementer or test-author subagent when the plan assigns a specific model tier or required skills; the built-in `general` subagent is acceptable only when no explicit binding is required.
- Keep prompts narrow for test authors and complete for implementers, matching the canonical workflow.

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

Use one ordered fallback chain for both structural TDD workspaces and task-scoped implementer worktrees:

1. Use Worktrunk (`wt`) if it is installed and suitable.
2. Otherwise use `git worktree` via Bash to create an isolated workspace.

Do not use runtime-native worktree switching for isolated workspace creation or switching. Use OpenCode-native mechanics to dispatch workers into the selected Worktrunk or git worktree and to verify the worker actually ran there.

When Worktrunk (`wt`) is used, create or enter the isolated workspace with `wt switch`. Per Worktrunk (`wt`) documentation, `wt switch` is the command that switches to a worktree and creates one if needed; use `wt switch --create <branch>` when the isolation branch does not exist yet. Its `--execute` mode can also be useful when you need to launch the agent directly inside the isolated worktree.

After creating or entering an isolated workspace, seed ignored build/cache artifacts when the project needs them for practical compile or test performance. Identify the relative cache directories from project docs or config, such as Rust `target/`. When the source cache exists in the coordinator workspace and the paths are on the same filesystem, create the destination directory and hard-link copy the contents, for example with `cp -al <source-dir>/. <worktree-dir>/<relative-dir>/` when supported. Verify the expected files exist in the worktree. If hard-link seeding is required for a build-heavy project but cannot be verified, record the failure and ask before dispatching the worker.

Do not assume structural TDD or implementer isolation is possible just because subagents exist. The executor must verify that the worker actually runs in the isolated worktree rather than the main workspace.

## Integration Branch And Checkpoints

- Create or resume the canonical execution integration branch from the recorded execution base before dispatching implementation workers.
- Create task worktrees from the integration branch checkpoint that contains the task's prerequisites.
- After integrating a task result into the integration branch and passing verification, create a local checkpoint commit using the project's normal signing policy.
- Do not push checkpoint commits. At completion, leave the aggregate result for review by mixed-resetting the integration branch to the execution base.

## Test Author Isolation

For OpenCode structural TDD, the required verification is that the test-author process actually runs in the isolated workspace. Same-workspace `@<subagent-name>` invocation is not sufficient for this gate unless the runtime can prove it routes that subagent into the isolated worktree. Acceptable isolation-routing mechanisms include:

- Worktrunk (`wt`) command: `wt switch --create <isolation-branch> --execute 'opencode run --agent <test-author-worker> --dir "$PWD" "<AC-only prompt>"'`
- Creating the isolated worktree with Worktrunk (`wt`), determining its path, then running `opencode run --agent <test-author-worker> --dir <isolated-worktree-path> "<AC-only prompt>"`
- Creating the isolated worktree with `git worktree`, determining its path, then running `opencode run --agent <test-author-worker> --dir <isolated-worktree-path> "<AC-only prompt>"`

OpenCode's public docs describe subagents, custom agent files, and permissions, but they do not currently document a first-class mechanism for dispatching a subagent into an arbitrary alternate workspace. Therefore:

- If any priority-order isolation mechanism plus `opencode run --agent --dir` or verifiable native subagent routing is available, do not skip structural TDD without first attempting or otherwise concretely verifying one of those routing paths.
- If routing cannot be verified after an attempted dispatch, record the exact attempted mechanism and failure in progress, then stop and ask whether to fix isolation or explicitly skip structural TDD.

**Physical isolation**: The isolated worktree must contain only the tracked code surface and explicitly seeded ignored build/cache artifacts the test author needs. Do not copy plan files, review files, or `progress.md` into it.

**Contextual isolation**: Even with a separate worktree, do not reveal the plan path, task file path, feature name, or design rationale to the test author. Pass only acceptance criteria and the relevant code surface.

**Bringing test files back**: If the test author ran in a task-scoped implementer worktree, leave the test files there for the implementer. Otherwise, return the resulting test files to the main workspace using the mechanism that matches the isolation method:

- If the isolated workspace was created with Worktrunk (`wt`), prefer `wt merge` so the branch is merged back into the integration branch and the worktree is removed in the same documented workflow.
- If the isolated workspace was created with plain `git worktree`, use normal git/worktree mechanics that preserve the authored tests without exposing the planning artifacts. Prefer a small commit plus cherry-pick into the integration branch, or another explicit patch-transfer mechanism you can verify.

In both cases, verify that the test files now exist in the main execution workspace and that the isolated workspace has been removed unless there is an explicit reason to keep it.

## Implementer Worktree Isolation

- For implementation tasks that run concurrently with any other implementation task, create or enter a task-scoped worktree using the Workspace Isolation Strategy.
- Before dispatch, hard-link seed required ignored build/cache artifact directories into the task worktree when the project is build-heavy.
- Dispatch the implementer with `opencode run --agent <implementer-worker> --dir <task-worktree-path> "<full task packet>"` unless a native OpenCode dispatch path can target that worktree and preserve the worker's configured model and permissions.
- Do not copy plan files, review files, or `progress.md` into the task worktree. The parent executor passes the full sub-plan content, prerequisite outputs, and test file paths as inline task context.
- After the implementer finishes, inspect the task worktree diff and integrate it into the execution integration branch using the mechanism that matches how the worktree was created: `wt merge` from Worktrunk (`wt`), or explicit git merge/cherry-pick/patch transfer for plain `git worktree`.
- If integration conflicts, checkpoint commit creation fails, or verification fails after integration, record the task as blocked and keep enough worktree state for diagnosis. Remove the worktree only after the result is integrated and checkpointed or intentionally abandoned.

## Implementer Dispatch

- Dispatch a separate implementer subagent with the full task packet, test file paths, prerequisite outputs, and required skills.
- Use native OpenCode subagent dispatch for assigned implementer workers when available. Use `opencode run --agent <implementer-worker> --dir <workspace-path> "<full task packet>"` when the current runtime surface cannot invoke the custom worker directly or when explicit workspace routing is needed.
- Tell the implementer explicitly that tests are immutable.
- If the implementer reports a dispute, record it in progress and continue with independent tasks per the canonical workflow.

## Progress and Artifacts

- The parent executor owns `progress.md` and updates it after each meaningful step.
- Even if a subagent writes files directly, the parent remains responsible for checkpointing and verifying that expected artifacts exist.
- Record dispatch evidence in progress: planned worker, actual worker, model/effort from the worker binding, command or runtime mechanism used, implementation workspace path, build/cache seeding status, checkpoint commit, integration status, and TDD isolation outcome.

## Model Assignment

- Use the custom subagent's explicit `model` field as the source of truth.
- Treat plan-assigned tiers as binding. Do not silently upgrade or downgrade.
- Because built-in subagents and custom subagents without `model` inherit from the caller in OpenCode, inherited model selection is not sufficient when the plan made an explicit model decision.
- If the required model is unavailable, stop and ask the user how to proceed.
