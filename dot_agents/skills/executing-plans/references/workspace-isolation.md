# Workspace Isolation

Use this reference when implementation tasks run concurrently, structural TDD needs an isolated workspace, or build/cache artifacts must be seeded into worktrees.

## Purpose

DAG independence does not make a shared dirty workspace safe. Concurrent implementers need separate worktrees, and dependent worktrees need a Git base that contains their prerequisites.

## Workspace Selection

- Sequential tasks may run in the main execution workspace unless the plan or user requires isolation.
- Concurrent implementation tasks must use task-scoped implementer worktrees.
- Structural TDD test authors must run in an isolated workspace when structural TDD is used.
- Dependent task worktrees must be created from an integration branch checkpoint that already contains all prerequisite outputs.

If prerequisite outputs exist only as local dirty files, do not create a stale worktree. Integrate and checkpoint the prerequisites first, or serialize/block the dependent task.

## Isolation Priority

Use the active runtime adapter's workspace isolation strategy. The default priority is:

1. Worktrunk (`wt`), when installed and suitable.
2. Plain `git worktree`.

Do not use runtime-native worktree switching for this skill unless the active runtime adapter explicitly permits it. Runtime adapters dispatch workers into the selected worktree and verify the worker actually ran there.

If no isolated implementation path can be verified for a concurrent task, serialize that task or ask the user before weakening isolation.

## Worker Worktree Rules

- Create task worktrees from the current integration branch checkpoint that includes prerequisites.
- Do not copy plan files, review files, anchors, or `progress.md` into the worktree.
- Pass the sub-plan as an inline task packet, not as a plan file path.
- Pass prerequisite outputs through the executor, not by letting subagents communicate with each other.
- Verify the assigned worker can be dispatched inside the selected worktree before treating isolation as available.

## Build And Cache Seeding

Fresh worktrees do not contain ignored build artifacts. For build-heavy projects, identify ignored in-repository build/cache directories that materially affect compile or test performance.

Examples include Rust `target/`, Bazel output trees, or project-local build directories named by docs or config.

Seed only ignored build/cache artifacts. Never seed source files, plan files, review files, anchors, or `progress.md`.

When the source cache exists in the coordinator workspace and both paths are on the same filesystem, a hard-link-preserving copy is acceptable, such as:

```bash
cp -al <source-dir>/. <worktree-dir>/<relative-dir>/
```

Use an equivalent hard-link-preserving copy when `cp -al` is unavailable. Verify the destination contains expected files and record the seeded relative paths in progress.

If required hard-link seeding cannot be verified for a build-heavy project, record the failure and ask before dispatching a worker that would rebuild from an empty worktree.

## Integration And Cleanup

After a task-scoped worktree worker finishes, inspect the worktree result and integrate it according to [checkpoint-integration.md](checkpoint-integration.md).

If integration conflicts, checkpoint creation fails, or verification fails, keep enough workspace state for diagnosis. Remove the task worktree only after its result has been integrated and checkpointed or intentionally abandoned.
