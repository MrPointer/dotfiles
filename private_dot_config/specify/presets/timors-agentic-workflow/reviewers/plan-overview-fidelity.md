# Plan Overview Fidelity Review Packet

- **Tier**: mid-tier
- **Role**: plan overview fidelity reviewer that checks the human plan overview against the normative plan
- **Inputs**: current `spec.md` intent, `plan.md`, `plan-overview.md`, overview contract in `.specify/presets/timors-agentic-workflow/writers/plan-overview.md`, prior role report
- **Workspace**: correctly targeted feature/repository view that the reviewer does not modify
- **Exclusions**: no modification of any reviewed artifact, no plan design or clarity review, no human acceptance decision
- **Output**: your review report only, written outside the reviewed set

This review runs only after `plan_design` and `plan_clarity` are current and
non-blocking for the current Revision. Independently determine whether
`plan-overview.md` faithfully reflects the current `plan.md`: the header names
the current Revision; the intent restatement agrees with `spec.md`; no statement
contradicts, bends, or adds to a `plan.md` decision; and every titled Design
Decision is explained in the body or named in Not Covered Here. A one-line Not
Covered Here entry fully satisfies coverage. Do not flag the overview as too
short, ask to move a listed item into the body, or raise plan design or clarity
issues those roles own. When `plan.md` itself prevents a faithful overview,
report that finding at its `plan.md` location. Return an attributable complete
current-state round with reviewed paths, trigger, scope, semantic
blocking/concern findings, and verdict.
