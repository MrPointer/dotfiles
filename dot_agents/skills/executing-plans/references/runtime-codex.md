# Runtime Adapter: Codex

Use this adapter only when the active runtime is Codex.

This adapter maps the canonical execution workflow in `../SKILL.md` to Codex-native mechanics. It must not redefine task ordering, dispute policy, progress rules, or the test-author/implementer separation.

## Exploration and Dispatch

- Use Codex sub-agent dispatch for exploration, test authoring, and implementation work when sub-agents materially help.
- Set the sub-agent `model` explicitly using Codex's actual dispatch mechanism when the plan assigns a tier.
- Prefer narrow prompts with only the context required for the specific role.
- Before executing a multi-sub-plan plan, verify every assigned dispatch recipe or worker can be invoked. If a binding fails, diagnose and retry once. If it still fails, stop and ask the user; do not perform assigned implementation in the coordinator context.

## Execution Bindings

- In Codex, execution bindings should use the built-in worker path for implementation and test-author work unless the project explicitly defines a more specialized worker.
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

Use one ordered fallback chain for both structural TDD workspaces and task-scoped implementer worktrees:

1. Prefer Codex's native isolated-workspace or worktree mechanism when the current environment provides one and the executor can verify creation, worker dispatch, result collection, and cleanup.
2. Otherwise use Worktrunk (`wt`) to create an isolated workspace.
3. Otherwise use `git worktree` directly.

After creating or entering an isolated workspace, seed ignored build/cache artifacts when the project needs them for practical compile or test performance. Identify the relative cache directories from project docs or config, such as Rust `target/`. When the source cache exists in the coordinator workspace and the paths are on the same filesystem, create the destination directory and hard-link copy the contents, for example with `cp -al <source-dir>/. <worktree-dir>/<relative-dir>/` when supported. Verify the expected files exist in the worktree. If hard-link seeding is required for a build-heavy project but cannot be verified, record the failure and ask before dispatching the worker.

## Test Author Isolation

- Structural TDD in Codex is allowed only when the test author can be dispatched into that isolated workspace.
- For build-heavy projects, seed required ignored build/cache artifacts in the isolated workspace before dispatching the test author.
- If any priority-order worktree mechanism plus Codex worker dispatch is available, do not skip structural TDD without first attempting or concretely verifying the isolated dispatch path.
- If the active Codex environment cannot run the test-author worker inside the isolated workspace after an attempted dispatch, record the exact attempted mechanism and failure, then stop and ask whether to fix isolation or explicitly skip structural TDD.
- Do not reveal the plan path, task file path, feature name, or design rationale to the test author.
- Pass only acceptance criteria and the code surface the tests interact with.
- When structural TDD is used, prompt hygiene is mandatory in addition to physical isolation.

## Implementer Worktree Isolation

- For implementation tasks that run concurrently with any other implementation task, create or enter a task-scoped worktree using the Workspace Isolation Strategy.
- Before dispatch, hard-link seed required ignored build/cache artifact directories into the task worktree when the project is build-heavy.
- Dispatch the assigned implementer worker inside that worktree using Codex's actual worker dispatch mechanism and explicit model/skill binding.
- Do not copy plan files, review files, or `progress.md` into the worktree. The parent executor passes the full sub-plan content, prerequisite outputs, and test file paths as inline task context.
- After the implementer finishes, inspect the task worktree diff and integrate it into the coordinator workspace using native Codex result collection when available, then `wt merge` from Worktrunk (`wt`), then explicit git merge/cherry-pick/patch transfer.
- If integration conflicts or verification fails after integration, record the task as blocked and keep enough worktree state for diagnosis. Remove the worktree only after the result is integrated or intentionally abandoned.

## Implementer Dispatch

- Spawn a separate implementer sub-agent with the full task packet, test file paths, prerequisite outputs, and required skills/reference material.
- Tell the implementer explicitly that tests are immutable.
- If the implementer reports a dispute, record it and continue according to the canonical workflow.

## Progress and Artifacts

- The parent executor owns `progress.md` and should update it after each meaningful step.
- Even if a sub-agent writes files directly, the parent remains responsible for checkpointing and artifact verification.
- Record dispatch evidence in progress: planned worker, actual worker, model/effort, runtime dispatch mechanism, implementation workspace path, build/cache seeding status, integration status, and TDD isolation outcome.

## Model Assignment

- Use Codex's explicit model-selection mechanism in the worker dispatch path; do not rely on prompt wording.
- Treat plan-assigned tiers as binding.
- If the required model is unavailable, stop and ask the user how to proceed.
