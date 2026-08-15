---
description: Execute approved sub-plans in isolated worktrees and integrate packet commits.
---

## Authority and Inputs

Require Ready canonical `spec.md`; canonical `plan.md` Design Acceptance Accepted
for its current Revision; a closed, reviewed `tasks.md`; `authorization: approved`;
and a fully clean, named, stable active branch. Read
`references/artifact-contracts.md`, `references/execution-lifecycle.md`,
`references/workspace-isolation.md`, `references/scheduling-policy.md`,
`references/model-and-worker-selection.md`, `references/testable-work.md`, and
`references/documentation-planning.md`. Use `templates/progress-template.md` and,
when required, `reviewers/component-docs.md` with
`templates/review-report-template.md`.

## Procedure

1. Perform the concrete phase checks in `references/execution-lifecycle.md`; do
   not redesign the approved package. Resolve, pin, and later refresh signing
   policy exactly under `references/workspace-isolation.md`. Add the exact local
   progress path to `.git/info/exclude` after clean preflight, then create initial
   progress from its template. Creation is the lock event and records active branch,
   starting commit, complete protected baseline, and policy evidence.
2. For each ready canonical identity, validate closed package structure, selected
   tracked regular file, predecessor status and contract evidence, ownership, test
   mode/basis, and baseline equality. Bind a worker under
   `references/model-and-worker-selection.md`. Create a dedicated branch/worktree,
   derive and record its worktree-root-relative locator, and dispatch only the
   identity, locator, invocation-specific constraints, prerequisite evidence,
   recovery context, and limited result envelope. Do not copy packet prose into
   the prompt.
3. Let `references/scheduling-policy.md` control waves. Workers may browse but not
   expand authority, modify protected feature files, rewrite history, merge
   unrelated work, or push. Validate returned changed paths, commits, verification,
   output evidence, runtime outcomes, attribution, clean branch, and baseline at
   result acceptance. Follow its narrow retry and recovery rules.
4. Before each rebase, squash, candidate construction, corrective commit, and
   integration, refresh signing-policy evidence. Rebase in approved topological
   order, recheck the baseline, construct a separately inspectable candidate,
   verify tree and applicable signature, then fast-forward exactly one packet
   commit. Use Worktrunk only after its required non-mutating capability proof;
   otherwise use the plain-Git fallback. Record final commit and contract output
   evidence before a consumer is ready.
5. After all packets are done, run aggregate verification and the required
   independent component-documentation review. Correct only bounded in-scope work
   with a policy-compliant corrective commit, then rerun affected gates. New scope,
   design, contract, or acceptance exceptions return to their owning phase. Report
   every test mode, verification outcome, runtime condition, alternative, residual
   limitation, final commit, blocker, and final status; do not push.
