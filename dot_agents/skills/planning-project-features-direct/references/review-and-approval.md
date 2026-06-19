# Review And Approval

Use this reference when running the initial direct-planning review loop, checking consistency, presenting the plan, or handling user feedback.

## Reviewer Types

The review loop uses global and project-local reviewers.

Global reviewers are project-agnostic roles:

- `design-reviewer`: evaluates whether the plan preserves design boundaries, ownership, integration seams, compatibility constraints, migration/rollback strategy, technical risks, and hidden complexity.
- `plan-clarity-reviewer`: catches vague, ambiguous, or speculative language that would force executing agents to make planner-owned decisions.
- `plan-executability-reviewer`: checks file ownership, acceptance criteria, dependency order, verification scope, worker dispatch, and isolated execution mechanics.

Project-local reviewers are domain-specialized roles defined by the project. Discover them using the active runtime adapter and match them by description, not naming convention.

If a required reviewer binding is unavailable in the active runtime, report the gap and stop rather than recreating the reviewer ad hoc.

## Review Artifact Ownership

Review output is saved in `reviews/` within the plan directory, named `<plan-file>.<reviewer-type>.md`:

```text
<plan-directory>/reviews/
├── 00-master.design.md
├── 00-master.clarity.md
├── 00-master.executability.md
├── 01-data-model.<local-reviewer>.md
└── ...
```

Pass reviewers the plan path, requested review output path, and review task. If a reviewer writes its own artifact, use it. If it returns findings instead, write the artifact from the response. Verify the expected file exists before continuing.

The `reviews/` directory is ephemeral and normally covered by the `plans/` ignore rule, but it persists locally across sessions for reference.

## Initial Review Loop

Run this loop once before the user sees the plan.

1. Launch `design-reviewer`, `plan-clarity-reviewer`, and `plan-executability-reviewer` against the master plan. Run them in parallel when the runtime supports it.
2. Incorporate findings into the master plan and affected sub-plans.
3. After master-plan review is resolved, launch each sub-plan's assigned project-local reviewer. Sub-plan reviews can run in parallel.
4. Normalize local reviewer output to the standard format when needed: Verdict, Critical Findings, Concerns, Observations.
5. Re-run only affected reviewers until no reviewer produces blocking findings.

Do not restart the entire review loop during convergence. Only re-review plans that changed.

Reviewers evaluate architecture, contracts, constraints, risks, clarity, and acceptance criteria. If a reviewer suggests implementation specifics that belong to the executing agent, reject that suggestion.

## Planning-Rule Validation

Before incorporating a reviewer finding, verify it does not contradict this skill.

Reviewer agents may load domain-specific skills but not the planning skill itself. They may suggest changes such as model downgrades, skipping review steps, or adding implementation detail that violate planning-level constraints.

When a conflict exists, this skill's rules take precedence. Note the reviewer's rationale in the review file, but do not apply the conflicting change.

The user may request additional specialized reviewers, such as security or performance, for specific sub-plans. Add them on request; they are not part of the default flow.

## Final Consistency Check

Before presenting the plan, verify:

1. Model assignments: every sub-plan model matches the decision tree and project rules.
2. Skill conformance: reviewer-driven changes do not contradict required skills listed in sub-plans.
3. Execution bindings: every multi-sub-plan worker assignment has a runtime-specific binding and lead-agent dispatch instructions.
4. Cross-sub-plan prerequisites: every referenced prerequisite is created by an earlier sub-plan.
5. Integration contract integrity: every cross-sub-plan `Produces` entry has matching `Consumes` coverage, and the master data-flow table covers every cross-boundary path.
6. DAG validity: no cycles, no dependency on later or same-group output, and every consumed prerequisite points to an earlier group.
7. Parallel execution safety: same-group sub-plans have non-overlapping primary file ownership and execution instructions require task-scoped worktrees, required build/cache seeding, result integration, and serialization when isolation cannot be verified.
8. Anchor boundaries: active feature anchors are not duplicated, and any sub-plan using `anchoring-context` has a feature-level reason to update it.

Fix inconsistencies before presenting the plan.

## Presentation

Present the reviewed plan with:

- review summary
- worker-dispatch summary
- how findings were addressed
- missing local reviewer coverage approved by the user, if any
- any persistent-binding reload or restart warning from the runtime adapter

Only mark the plan ready when the user approves.

Remind the user that the plan intentionally omits implementation details. The user reviews architecture and constraints now; they review actual code after execution.

## Handling User Feedback

When the user requests changes, incorporate them and classify the change:

| Change Type | Examples | Re-review Action |
|-------------|----------|------------------|
| Cosmetic / wording | Clarify a step, rename a sub-plan, fix typos | None |
| Scoped implementation detail | Add an edge case to one sub-plan, change a file path, adjust a step | None; planner judgment is sufficient |
| Scope adjustment within a sub-plan | Add/remove acceptance criteria, change approach for one sub-plan | Re-review only that sub-plan with its assigned reviewer |
| Structural change | New sub-plan, dependency graph change, boundary shift, merge/split, ownership or verification scope change | Re-review affected sub-plans plus `design-reviewer` and `plan-executability-reviewer` on the master plan |

Default behavior: after incorporating feedback, state what changed and recommend whether re-review is warranted. Do not automatically re-run reviewers unless the workflow requires it or the user asks.

The user can always explicitly request a re-review regardless of change classification.

## Post-Execution Component Documentation Review

If the project has component documentation, record that `component-docs-reviewer` should run after all sub-plans complete to catch implementation-vs-plan drift in component docs.
