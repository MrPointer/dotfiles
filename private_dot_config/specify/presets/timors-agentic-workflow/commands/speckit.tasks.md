---
description: Produce, review, and authorize an immutable sub-plan task package.
---

## Authority and Inputs

Require Ready canonical `spec.md` and canonical `plan.md` with Design Acceptance
Accepted for its current Revision. Resolve the canonical feature directory from
`.specify/feature.json` per
`.specify/presets/timors-agentic-workflow/references/artifact-contracts.md` before
reading any artifact. Read `.specify/presets/timors-agentic-workflow/references/artifact-contracts.md`,
`.specify/presets/timors-agentic-workflow/references/task-planning.md`, `.specify/presets/timors-agentic-workflow/references/scheduling-policy.md`,
`.specify/presets/timors-agentic-workflow/references/testable-work.md`, `.specify/presets/timors-agentic-workflow/references/documentation-planning.md`,
`.specify/presets/timors-agentic-workflow/references/model-and-worker-selection.md`, and `.specify/presets/timors-agentic-workflow/references/review-lifecycle.md`.
Use `.specify/presets/timors-agentic-workflow/templates/tasks-template.md`, `.specify/presets/timors-agentic-workflow/templates/subplan-template.md`, and
`.specify/presets/timors-agentic-workflow/templates/review-report-template.md`; use reviewer packets
`.specify/presets/timors-agentic-workflow/reviewers/tasks-plan-fidelity.md` and `.specify/presets/timors-agentic-workflow/reviewers/tasks-executability.md`.

## Procedure

1. Refuse change if initial local `progress.md` locks or completes the package.
   Otherwise create or reconcile `tasks.md` and all indexed sub-plans together.
   Do not create atomic tasks, checkboxes, a separate execution artifact, or
   provider-specific runtime state.
2. Apply every closed-schema and structural check in `.specify/presets/timors-agentic-workflow/references/task-planning.md`:
   canonical identity resolution, ordered path-only DAG, ledger projection,
   predecessor-table agreement, ownership and concurrency, four-mode testing,
   producer body/frontmatter one-to-one agreement, consumer joins, and handoff.
   Missing, unknown, ambiguous, escaping, unmatched, or contradictory data blocks
   review and authorization; do not infer or migrate it.
3. Add a final documentation packet only under
   `.specify/presets/timors-agentic-workflow/references/documentation-planning.md` conditions. Record the required or
   not-applicable component-documentation handoff with paths or substantive reason.
4. Dispatch the two independent reviewers using the shared binding, attribution,
   quiescence, selective-rerun, retained-pointer, and started-review-recovery rules
   in `.specify/presets/timors-agentic-workflow/references/review-lifecycle.md`. Append complete report rounds and preserve
   history. No blocking, pending, recovery, stale, or contradictory pointer may
   reach authorization.
5. Request the scalar Implementation Authorization only after both current reviews
   permit it. Record human rationale and decision context in the Markdown section;
   write only `pending`, `approved`, or `rejected` in frontmatter. Only `approved`
   hands off to `implement`.

Old task artifacts and old runtime state are not migrated, normalized, or inferred.
Report canonical paths, reviews, authorization, and the next permitted action.
