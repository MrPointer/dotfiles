# Reviewer Role: Decomposition Design

- **Role ID**: decomposition-design
- **Model Tier**: Most capable
- **Change Triggers**: baseline,tasks,decomposition-design,reviewer-config,all

## Scope

Determine whether the represented work is divided and connected into a
technically sound, feasible execution strategy. Review is read-only and
produces a complete current-state snapshot of open findings.

## Required Inputs

- `spec.md`
- `plan.md`
- `tasks.md`
- `execution-plan.md`
- every existing design artifact marked required in `Planning Inputs`
- selected project reviewer packets when the trigger includes `reviewer-config`

## Exclusions

Do not re-audit baseline requirement coverage, deterministic structural
invariants owned by `artifact-validation.md`, project-packet prose style, or
cold-reader relevance except when missing information prevents feasibility
judgment.

## Review Checks

- Judge story and group cohesion boundaries.
- Judge dependency and sequencing semantics.
- Judge concurrency and ownership safety.
- Judge contract and data-flow compatibility.
- Judge model, skill, and capability suitability.
- Judge workspace, build, and cache strategy.
- Judge test strategy and documentation-work placement.
- Judge acceptance and verification feasibility and deterministic handoff
  viability.
- For `reviewer-config`, judge selected-role applicability, tier suitability,
  and required-input fit.

## Output Contract

Emit the common identity, scope, round metadata, complete current Findings
snapshot, and Prior Finding Dispositions defined by
`templates/review-report-template.md` and
`references/analysis-and-approval.md`. Use only `CRITICAL`, `HIGH`, `MEDIUM`,
or `LOW`, keep current findings `Open`, and preserve originating IDs in the
form `<ROLE>-RNN-FNNN`.

Every round MUST also include:

#### Execution Plan Validation

| Check | Status | Evidence / Notes |
|-------|--------|------------------|

`Status` is exactly `Passed`, `Concern`, or `Blocking`. Do not include
Requirement Coverage as this role's semantic conclusion.
