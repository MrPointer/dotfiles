# Workspace Isolation

## Purpose

Isolation prevents concurrent groups from sharing a dirty working tree. It does not prove worker capacity, build-output safety, cache correctness, or dependency independence; those are separate ready-group gates in `execution-lifecycle.md`.

## Workspace Selection

- Resolve and record an absolute workspace path for every bound group.
- Parallel groups use separate worktrees. A dependent worktree starts only from an integration-branch checkpoint containing every prerequisite output.
- Sequential groups may use the integration workspace unless the plan requires isolation.
- A `Linear DAG` is always serialized. `Parallel allowed` may be serialized when capacity or output/cache safety cannot be proven; record the fallback and reason in progress.
- Never create a dependent workspace from prerequisite output that exists only as dirty files.

## Dirty-State Preflight

Immediately before a ready group creates, enters, or mutates a workspace, inspect tracked modifications and untracked non-ignored files in the coordinator and target workspace. Record the exact evidence in Execution Audit.

If unexpected state exists, stop before workspace mutation, edits, or dispatch. Never stash, delete, commit, copy, ignore, reset, or otherwise relocate dirty state without explicit authorization. If authorization permits proceeding, record the paths, scope, decision, and resulting target status. Unrelated state carried into a worktree still blocks dispatch.

`progress.md` is a coordinator artifact and must remain ignored. Do not copy progress, plans, reviews, or other coordinator artifacts into worker worktrees.

## Concurrent Capacity And Output Safety

Before any member of a parallel set starts, prove that all members can run concurrently and that each has:

- a distinct absolute worktree;
- an available attributable worker binding;
- nonoverlapping planned file ownership;
- safe generated-output and build directories; and
- a project-documented or verified cache strategy, or an explicit no-reuse result.

Start none until all pass. Filesystem isolation alone is not evidence that path-sensitive caches or shared generated outputs are safe.

Prefer project-documented shared caches. Seed only ignored build/cache artifacts when project instructions or targeted verification proves portability between absolute paths. Never seed source, plans, reviews, progress, or checkpoints. If reuse is unsafe or unnecessary, record that it is skipped. If required reuse cannot be proven, block.

## Workspace Result Retention

Verify dispatch actually targeted the recorded absolute path. After result return, inspect changed and untracked files against group ownership and the attributable Result. Keep a worktree when dispatch start/result is ambiguous, integration or verification fails, a checkpoint is missing, or recovery is required.

Never automatically remove a branch or worktree. Release a workspace only after its exact result is integrated and checkpointed, or after explicit abandonment records the disposition of every dirty file and checkpoint.
