# Review And Approval

Use this reference when running RFC-backed plan review, checking consistency, presenting the plan, or handling user feedback.

## RFC-Backed Review Loop

Run exactly these reviewers by default:

- `plan-rfc-fidelity-reviewer`
- `plan-executability-reviewer`

Do not run `design-reviewer`, `plan-clarity-reviewer`, or project-local domain reviewers unless the user explicitly exits RFC-backed planning and returns to direct planning. RFC design and clarity were already reviewed at RFC time.

Pass each reviewer:

- plan directory path
- RFC path
- requested review output path
- review task

Review output lives in the plan's `reviews/` directory:

```text
<plan-directory>/reviews/
├── 00-master.rfc-fidelity.md
└── 00-master.executability.md
```

If a reviewer has write permission and writes its own review artifact, use it. If it returns findings instead, write the review artifact from the response. Verify the expected artifact exists before continuing.

Incorporate findings and re-run only affected reviewers until both report no blocking findings. If a finding requires changing the RFC design, stop and ask whether to revise the RFC or approve an explicit deviation.

## Final Consistency Check

Before presenting the plan, verify:

1. RFC coverage: every RFC goal and success criterion is represented.
2. Non-goal preservation: every RFC non-goal remains out of scope.
3. Constraint propagation: every execution-relevant RFC constraint appears in the master plan or relevant sub-plan.
4. Model assignments: every sub-plan model matches the decision tree and project rules.
5. Execution bindings: every multi-sub-plan worker assignment has a runtime-specific binding and lead-agent dispatch instructions.
6. Skill conformance: reviewer-driven changes do not contradict required skills listed in sub-plans.
7. Concurrency policy: the master plan records `Linear DAG` or `Parallel allowed`, exception-list ecosystems are linearized unless an explicit project-documented or user-approved override is recorded, everything else remains parallel-eligible unless a concrete shared build/output/cache constraint is documented, and policy-only sequencing is labeled rather than represented as fake data flow.
8. Cross-sub-plan prerequisites: every referenced prerequisite is created by an earlier sub-plan.
9. Integration contract integrity: every cross-sub-plan `Produces` entry has matching `Consumes` coverage, and the master data-flow table covers every cross-boundary path.
10. DAG validity: no cycles, no dependency on later or same-group output, and every consumed prerequisite points to an earlier group.
11. Parallel execution safety: when policy allows same-group execution, sub-plans have non-overlapping primary file ownership and execution instructions require task-scoped worktrees, required build/cache seeding, result integration, and serialization when isolation or output/cache safety cannot be verified.
12. Anchor boundaries: active feature anchors are not duplicated, and any sub-plan using `anchoring-context` has a feature-level reason to update it.
13. Review status: reviewer findings are resolved or explicitly documented as non-blocking.

Fix inconsistencies before presenting the plan.

## Presentation

Present the reviewed plan with:

- RFC path
- review summary
- worker-dispatch summary
- explicit approved deviations, or `None`
- any persistent-binding reload or restart warning from the runtime adapter

Only mark the plan ready when the user approves.

Remind the user that the plan intentionally omits implementation details. The user reviews architecture and constraints now; they review actual code after execution.

## Handling User Feedback

When the user requests changes, incorporate them and classify the change:

| Change Type | Examples | Re-review Action |
|-------------|----------|------------------|
| Cosmetic / wording | Clarify a step, rename a sub-plan, fix typos | None |
| Scoped implementation detail | Add an edge case to one sub-plan, change a file path, adjust a step | None; planner judgment is sufficient |
| Scope adjustment within a sub-plan | Add/remove acceptance criteria, change approach for one sub-plan | Re-review affected sub-plan through `plan-rfc-fidelity-reviewer` if RFC scope may be affected, and `plan-executability-reviewer` if mechanics changed |
| Structural change | New sub-plan, dependency graph change, boundary shift, merge/split, ownership or verification scope change | Re-review affected plan sections with both reviewers |
| RFC design change | Goal, non-goal, contract, risk, or accepted design decision changes | Stop and ask whether to revise the RFC or approve an explicit plan deviation before re-review |

Default behavior: after incorporating feedback, state what changed and recommend whether re-review is warranted. Do not automatically re-run reviewers unless the workflow requires it or the user asks.

## Post-Execution Component Documentation Review

If the project has component documentation, record that `component-docs-reviewer` should run after all sub-plans complete to catch implementation-vs-plan drift in component docs.
