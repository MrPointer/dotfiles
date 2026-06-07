# Workspace Isolation

Use this reference when implementation tasks run concurrently, structural TDD needs an isolated workspace, or build/cache reuse must be configured for worktrees.

## Contents

- [Purpose](#purpose)
- [Workspace Selection](#workspace-selection)
- [Dirty State Preflight](#dirty-state-preflight)
- [Isolation Priority](#isolation-priority)
- [Worker Worktree Rules](#worker-worktree-rules)
- [Build And Cache Reuse](#build-and-cache-reuse)
- [Integration And Cleanup](#integration-and-cleanup)

## Purpose

DAG independence does not make a shared dirty workspace safe. Concurrent implementers need separate worktrees, and dependent worktrees need a Git base that contains their prerequisites.

## Workspace Selection

- Sequential tasks may run in the main execution workspace unless the plan or user requires isolation.
- Concurrent implementation tasks must use task-scoped implementer worktrees.
- Structural TDD test authors must run in an isolated workspace when structural TDD is used.
- Dependent task worktrees must be created from an integration branch checkpoint that already contains all prerequisite outputs.

If prerequisite outputs exist only as local dirty files, do not create a stale worktree. Integrate and checkpoint the prerequisites first, or serialize/block the dependent task.

## Dirty State Preflight

Before creating a new isolated worktree, or running a switch command that may create one, inspect the coordinator workspace status with a command that reports tracked modifications and untracked files. If any non-ignored file is dirty, stop before creating the worktree and warn the user.

The warning must explain that worktree tools may carry untracked files into the isolated worktree, require those files to be committed before merge, or cause coordination artifacts to be committed, dropped, or lost during cleanup.

Do not commit, stash, delete, ignore, or copy dirty files to make the worktree command succeed unless the user explicitly authorizes that action.

If the dirty files are plan, RFC, review, progress, or anchor artifacts, ask whether to ignore them with project `.gitignore` or local `.git/info/exclude`, resolve them another way, or continue with the risk. If the user explicitly insists on continuing, record the dirty paths and authorization in progress, create the worktree, verify the new worktree's initial status, and block dispatch if unrelated dirty or untracked files were carried into it.

## Isolation Priority

Use the active runtime adapter's workspace isolation strategy. The default priority is:

1. Worktrunk (`wt`), when installed and suitable.
2. Plain `git worktree`.

Do not use runtime-native worktree switching for this skill unless the active runtime adapter explicitly permits it. Runtime adapters dispatch workers into the selected worktree and verify the worker actually ran there.

If no isolated implementation path can be verified for a concurrent task, serialize that task or ask the user before weakening isolation.

## Worker Worktree Rules

- Create task worktrees from the current integration branch checkpoint that includes prerequisites.
- Run the dirty-state preflight before creating the worktree. When entering an existing worktree, verify its status before dispatch.
- Do not copy plan files, review files, anchors, or `progress.md` into the worktree.
- After creation, verify the worktree starts with only expected tracked code and allowed ignored cache artifacts. If unrelated dirty or untracked files are present, stop before dispatch and ask the user how to proceed.
- Pass the sub-plan as an inline task packet, not as a plan file path.
- Pass prerequisite outputs through the executor, not by letting subagents communicate with each other.
- Verify the assigned worker can be dispatched inside the selected worktree before treating isolation as available.

## Build And Cache Reuse

Fresh worktrees do not contain ignored build artifacts. Treat build/cache reuse as a best-effort performance optimization, not as a property of workspace isolation.

Before configuring shared caches or seeding ignored artifacts, read project-local instructions such as `AGENTS.md`, runtime adapter notes, plan instructions, repository docs, or build configuration. Use the project-documented strategy when one exists.

Build caches may be sensitive to absolute paths, environment variables, config files, compiler flags, generated metadata, or dependency graph state. Filesystem-level hard-link success does not prove that a build cache is valid in a different worktree path.

Prefer project-documented shared cache mechanisms over copying artifacts. If a project documents a safe shared cache directory, tool cache, environment variable, or task-runner cache configuration, apply that mechanism and record it in progress.

Seed only ignored build/cache artifacts, and only when the project instructions or a targeted verification show that artifact copying is safe across worktree paths. Never seed source files, plan files, review files, anchors, or `progress.md`.

If no safe project-specific cache reuse strategy is documented, skip artifact seeding, record that the isolated worktree may rebuild from scratch, and proceed unless the plan or user explicitly requires cache reuse because rebuild cost is unacceptable.

If required cache reuse cannot be verified, record the failure and ask before dispatching a worker that would rebuild from an empty worktree.

## Integration And Cleanup

After a task-scoped worktree worker finishes, inspect the worktree result and integrate it according to [checkpoint-integration.md](checkpoint-integration.md).

If integration conflicts, checkpoint creation fails, or verification fails, keep enough workspace state for diagnosis. Remove the task worktree only after its result has been integrated and checkpointed or intentionally abandoned.
