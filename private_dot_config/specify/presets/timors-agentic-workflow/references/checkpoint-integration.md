# Checkpoint Integration

## Integration State

Record the review branch, execution base, local integration branch, and `refs/agent-checkpoints/<feature>` before implementation. Resume only when the branch descends from the execution base and the progress, workspaces, Completed Artifacts, and checkpoint evidence agree.

Checkpoint commits are local execution plumbing. They follow the project's signing policy and MUST be signed. If signing or signature verification fails, stop; do not create or accept an unsigned substitute. Never push a checkpoint.

## Integrate A Group

Only a group in `ready for integration/checkpoint` may enter integration:

1. inspect the attributable Result and workspace state against planned ownership;
2. integrate only the group's expected files into the local integration branch;
3. verify produced files and shapes against every `DFNN` and producer `CTNN` contract;
4. run the exact post-integration verification;
5. create and verify a signed checkpoint;
6. move `refs/agent-checkpoints/<feature>` to the retained checkpoint;
7. persist Test Evidence, Completed Artifacts, Checkpoints, Execution Audit, and verification evidence;
8. update covered task checkboxes; then mark the group `done`.

Do not satisfy integration tooling by committing, deleting, or moving unrelated dirty files. A conflict, ownership mismatch, contract mismatch, regression, verification failure, or checkpoint failure blocks the group while preserving diagnostic state. Independent groups may continue only when the validated policy permits it.

Dependent groups start from the retained checkpoint that contains all prerequisite Completed Artifacts. A dirty result is never a completed prerequisite.

## Checkpoint States

- `pending`: no accepted signed checkpoint yet.
- `retained`: signed checkpoint and local ref are available for execution or recovery.
- `released`: execution state was explicitly accepted or abandoned and retention is no longer required.
- `missing`: progress requires a checkpoint but its commit/ref/signature evidence is unavailable; recovery blocks execution.

## Final Verification And Review State

After every group is `done`:

1. run full required verification at integration-branch tip;
2. create and verify the final signed checkpoint;
3. retain the tip at `refs/agent-checkpoints/<feature>`;
4. record `<execution-base>..<final checkpoint>` in Execution Context and Final Verification;
5. prepare the review surface with mixed-reset semantics to the execution base, leaving the aggregate implementation as a dirty diff;
6. record the aggregate branch, base, dirty paths, verification, checkpoint ref, and range.

Switch to the original review branch only when it is still at the execution base and doing so is safe. Otherwise leave the aggregate dirty diff where prepared and report its branch and base. Keep the local checkpoint ref until the user explicitly accepts or abandons the feature. Do not push it.

## Abandonment And Removal

An explicit abandonment returns the execution line to the recorded base while retaining enough evidence to account for branches, worktrees, checkpoints, artifacts, and dirty files. Ask for a disposition; never automatically reset, delete, clean, or remove them.

Preset removal during active execution is unsupported. Planned-but-unstarted removal uses reverted planning or inert history. In-flight work must finish pinned, pause with retained state, or be explicitly abandoned to base.
