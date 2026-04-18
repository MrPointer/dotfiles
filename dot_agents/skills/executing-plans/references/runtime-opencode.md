# Runtime Adapter: OpenCode

Use this adapter only when the active runtime is OpenCode.

This adapter maps the canonical execution workflow in `../SKILL.md` to OpenCode-native mechanics. It must not redefine task ordering, dispute policy, progress rules, or the test-author/implementer separation.

## Exploration and Dispatch

- Use OpenCode subagents for exploration, test authoring, and implementation work when subagents materially help.
- Prefer the built-in `explore` subagent for cheap, read-only codebase exploration such as the testability gate.
- Prefer a custom implementer or test-author subagent when the plan assigns a specific model tier or required skills; the built-in `general` subagent is acceptable only when no explicit binding is required.
- Keep prompts narrow for test authors and complete for implementers, matching the canonical workflow.

## Execution Bindings

OpenCode execution bindings should be **markdown-defined custom subagents** under `.opencode/agents/` (project-local) or `~/.config/opencode/agents/` (global).

- Reuse existing custom subagents when they match the task's model tier, permissions, and required skills.
- Prefer project-local bindings when the task depends on project-local skills or conventions.
- Set `mode: subagent` and set `model` explicitly. In OpenCode, subagents inherit the caller's model when `model` is omitted, so omission is not a reliable execution binding.
- Configure access with `permission`, not the deprecated `tools` field, for new or updated bindings.
- If the binding must use named skills, allow them via `permission.skill` and instruct the subagent in its prompt to load them with the `skill` tool at startup. OpenCode does not provide a declarative skill-preload field in agent frontmatter.
- Keep bindings thin: the task file provides task truth; the binding provides runtime-native identity, model, permissions, and skill-loading instructions.
- If the binding should not delegate further, deny or tightly scope `permission.task` so the worker does not spawn unrelated child agents.

## Test Author Isolation

**Workspace isolation**:

1. Prefer `wt` if it is installed and already part of the project's workflow.
2. Otherwise use `git worktree` via Bash to create a temporary isolated workspace.

When `wt` is used, create or enter the isolated workspace with `wt switch`. Per Worktrunk's documented behavior, `wt switch` is the command that switches to a worktree and creates one if needed; use `wt switch --create <branch>` when the isolation branch does not exist yet. Its `--execute` mode can also be useful when you need to launch the agent directly inside the isolated worktree.

OpenCode's public docs describe subagents, custom agent files, and permissions, but they do not currently document a first-class mechanism for dispatching a subagent into an arbitrary alternate workspace. Therefore:

- Do not assume structural TDD is possible just because subagents exist.
- Structural TDD is allowed only when you can verify that the test-author subagent will actually run against the isolated worktree rather than the main workspace.
- If you cannot verify that workspace routing, skip structural TDD and record `runtime cannot provide isolated test-author workspace` in progress.

**Physical isolation**: The isolated worktree must contain only the tracked code surface the test author needs. Do not copy plan files into it.

**Contextual isolation**: Even with a separate worktree, do not reveal the plan path, task file path, feature name, or design rationale to the test author. Pass only acceptance criteria and the relevant code surface.

**Bringing test files back**: If structural TDD is used successfully, return the resulting test files to the main workspace using the mechanism that matches the isolation method:

- If the isolated workspace was created with `wt`, prefer `wt merge` so the branch is merged back and the worktree is removed in the same documented workflow.
- If the isolated workspace was created with plain `git worktree`, use normal git/worktree mechanics that preserve the authored tests without exposing the planning artifacts. Prefer a small commit plus cherry-pick, or another explicit patch-transfer mechanism you can verify.

In both cases, verify that the test files now exist in the main execution workspace and that the isolated workspace has been removed unless there is an explicit reason to keep it.

## Implementer Dispatch

- Dispatch a separate implementer subagent with the full task, test file paths, prerequisite outputs, and required skills.
- Tell the implementer explicitly that tests are immutable.
- If the implementer reports a dispute, record it in progress and continue with independent tasks per the canonical workflow.

## Progress and Artifacts

- The parent executor owns `progress.md` and updates it after each meaningful step.
- Even if a subagent writes files directly, the parent remains responsible for checkpointing and verifying that expected artifacts exist.

## Model Assignment

- Use the custom subagent's explicit `model` field as the source of truth.
- Treat plan-assigned tiers as binding. Do not silently upgrade or downgrade.
- Because built-in subagents and custom subagents without `model` inherit from the caller in OpenCode, inherited model selection is not sufficient when the plan made an explicit model decision.
- If the required model is unavailable, stop and ask the user how to proceed.
