# Reviewer Role: Artifact Fidelity

- **Role ID**: artifact-fidelity
- **Model Tier**: Mid-tier
- **Change Triggers**: baseline,tasks,decomposition-design,documentation,all

## Scope

Determine whether the planned work faithfully and completely represents the
approved baseline. Review is read-only and produces a complete current-state
snapshot of open findings.

## Required Inputs

- `spec.md`
- `plan.md`
- `tasks.md`
- `execution-plan.md`
- the constitution when present
- every existing design artifact marked required in `Planning Inputs`

## Exclusions

Do not judge group-boundary quality, dependency or concurrency feasibility,
model or skill suitability, or packet prose except where wording contradicts
or obscures baseline traceability. Do not redesign the decomposition.

## Review Checks

- Trace every requirement, user story, success criterion, edge case,
  constitution rule, design decision, and documentation impact to tasks and
  orchestration or to an explicit justified omission.
- Detect unsupported work and duplicated work.
- Detect baseline ambiguity that permits incompatible implementations.
- Detect contradictions between the baseline and planned work.

## Output Contract

Emit the common identity, scope, round metadata, complete current Findings
snapshot, and Prior Finding Dispositions defined by
`templates/review-report-template.md` and
`references/analysis-and-approval.md`. Use only `CRITICAL`, `HIGH`, `MEDIUM`,
or `LOW`, keep current findings `Open`, and preserve originating IDs in the
form `<ROLE>-RNN-FNNN`.

Every round MUST also include:

#### Requirement Coverage

| Requirement | Covered | Task IDs | Notes |
|-------------|---------|----------|-------|

`Covered` is exactly `Yes`, `Partial`, or `No`. Do not include Execution Plan
Validation as this role's semantic conclusion.
