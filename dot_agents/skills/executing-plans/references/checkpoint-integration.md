# Checkpoint Integration

Use this reference when a plan has multiple sub-plans, task dependencies, task-scoped worktrees, or checkpoint commits.

## Purpose

Checkpoint integration materializes prerequisite outputs as local commits so dependent worktrees can start from real Git state. These commits are execution plumbing, not review units. The user's review surface is the final aggregate dirty diff.

## Establish The Integration Branch

Before implementation begins:

1. Record the current review branch or detached HEAD state in `progress.md`.
2. Record the execution base commit: `HEAD` before execution starts.
3. Record a local temporary integration branch name, such as `agent/<plan-id>` or the project convention.
4. Verify coordinator-owned dirty files are limited to plan files, review files, `progress.md`, anchors, or other explicit coordination artifacts.
5. If task-scoped worktrees may be created, run the dirty-state preflight in [workspace-isolation.md](workspace-isolation.md). Coordination artifacts that appear in git status must be ignored, resolved, or explicitly authorized by the user before worktree creation.
6. If implementation source, tests, config, generated code, or build-affecting files are dirty, stop and ask the user how to isolate them before checkpoint execution starts.
7. Create the integration branch from the execution base, or resume the branch recorded in progress.
8. Verify the resumed branch still descends from the recorded execution base.

Do not remove, stash, or commit coordinator-owned plan/progress artifacts as checkpoints.

## Checkpoint Task Results

After a sub-plan result is integrated and required verification passes, create a local checkpoint commit on the integration branch.

- Use the project's normal commit-signing policy.
- If signing is required and fails, stop rather than creating an unsigned checkpoint.
- Do not push checkpoint commits.
- Record the checkpoint commit in progress.
- A prerequisite task is not complete until its checkpoint exists.

Dependent tasks may start only from an integration-branch commit that contains all prerequisite outputs. If a prerequisite exists only as dirty files, checkpoint it first or serialize/block the dependent task.

## Integrate Task Worktrees

When a task ran in a task-scoped worktree:

1. Inspect the worktree result before integration.
2. Integrate it into the execution integration branch using the mechanism that matches the worktree tool.
3. Use `wt merge` for Worktrunk worktrees.
4. Use explicit git merge, cherry-pick, or patch transfer for plain `git worktree` worktrees.
5. Run the task's required verification after integration.
6. Create the checkpoint commit only after verification passes.
7. Mark the task `done` only after integration, verification, and checkpointing succeed.

Before running `wt merge` or any equivalent integration operation, verify the task worktree status contains only expected task outputs. If the tool requires committing unrelated carried files, plan/RFC/review/progress artifacts, or other unexpected untracked files, stop and ask; do not create a cleanup commit just to satisfy the tool.

If integration conflicts, checkpoint creation fails, or post-integration verification fails, mark the task `blocked: integration` or `blocked: regression`, keep enough workspace state for diagnosis, and continue with independent tasks when possible.

Remove a task worktree only after its result has been integrated and checkpointed or intentionally abandoned.

## Checkpoint Sequential Tasks

When a sequential task runs directly on the integration branch workspace, run its required verification, create the checkpoint commit, record it in progress, then mark the task `done`.

Sequential dirty output is not a completed prerequisite until it has a checkpoint commit.

## Prepare Final Review State

After all tasks are done:

1. Run the full required verification on the integration branch HEAD.
2. Record the final checkpoint range in progress: `<execution-base>..<integration-branch-head>`.
3. Convert checkpointed history into the user's review surface with mixed-reset semantics: `git reset --mixed <execution-base>`.
4. Update progress to explain that the aggregate feature changes are now dirty files.
5. Switch back to the original review branch only if it is still at the execution base and checkout is safe.
6. If switching back is unsafe, leave the dirty aggregate on the integration branch and report the branch and base clearly.

The final review state should be a normal dirty working tree that the user can inspect in their IDE. The checkpoint range remains recorded only for execution traceability.
