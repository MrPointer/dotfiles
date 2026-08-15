---
description: Produce, review, and authorize an immutable sub-plan task package.
---

## Authority and Inputs

Require Ready canonical `spec.md` and canonical `plan.md` with Design Acceptance
Accepted for its current Revision. Read `references/artifact-contracts.md`,
`references/task-planning.md`, `references/scheduling-policy.md`,
`references/testable-work.md`, `references/documentation-planning.md`,
`references/model-and-worker-selection.md`, and `references/review-lifecycle.md`.
Use `templates/tasks-template.md`, `templates/subplan-template.md`, and
`templates/review-report-template.md`; use reviewer packets
`reviewers/tasks-rfc-fidelity.md` and `reviewers/tasks-executability.md`.

## Procedure

1. Refuse change if initial local `progress.md` locks or completes the package.
   Otherwise create or reconcile `tasks.md` and all indexed sub-plans together.
   Do not create atomic tasks, checkboxes, a separate execution artifact, or
   provider-specific runtime state.
2. Apply every closed-schema and structural check in `references/task-planning.md`:
   canonical identity resolution, ordered path-only DAG, ledger projection,
   predecessor-table agreement, ownership and concurrency, four-mode testing,
   producer body/frontmatter one-to-one agreement, consumer joins, and handoff.
   Missing, unknown, ambiguous, escaping, unmatched, or contradictory data blocks
   review and authorization; do not infer or migrate it.
3. Add a final documentation packet only under
   `references/documentation-planning.md` conditions. Record the required or
   not-applicable component-documentation handoff with paths or substantive reason.
4. Dispatch the two independent reviewers using the shared binding, attribution,
   quiescence, selective-rerun, retained-pointer, and started-review-recovery rules
   in `references/review-lifecycle.md`. Append complete report rounds and preserve
   history. No blocking, pending, recovery, stale, or contradictory pointer may
   reach authorization.
5. Request the scalar Implementation Authorization only after both current reviews
   permit it. Record human rationale and decision context in the Markdown section;
   write only `pending`, `approved`, or `rejected` in frontmatter. Only `approved`
   hands off to `implement`.

Old task artifacts and old runtime state are not migrated, normalized, or inferred.
Report canonical paths, reviews, authorization, and the next permitted action.
