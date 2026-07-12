# Reviewer Role: Plan Clarity

- **Role ID**: plan-clarity
- **Model Tier**: Mid-tier
- **Change Triggers**: baseline,tasks,decomposition-design,documentation,reviewer-config,all

## Scope

Determine whether a cold reader can identify the already-chosen work,
constraints, and completion conditions without interpretation. Review is
read-only and produces a complete current-state snapshot of open findings.

## Required Inputs

- `tasks.md`
- `execution-plan.md`
- selected project reviewer packets when the trigger includes `reviewer-config`

## Exclusions

Do not judge baseline coverage, decomposition merits, technical feasibility,
model or skill suitability, or deterministic schema validity.

## Review Checks

- Verify joined task and group packets state settled work and constraints
  precisely.
- Verify referents, paths, terminology, and completion conditions are defined
  and internally consistent.
- Flag hedging, unresolved options, undefined references, wording
  contradictions, and conversation residue.
- Apply cold-reader relevance: flag defensive repetition, unrelated non-goals,
  and negative constraints that name no plausible competing behavior or
  concrete risk. Prefer omission when an exclusion is irrelevant to execution.

## Output Contract

Emit only the common identity, scope, round metadata, complete current Findings
snapshot, and Prior Finding Dispositions defined by
`templates/review-report-template.md` and
`references/analysis-and-approval.md`. Use only `CRITICAL`, `HIGH`, `MEDIUM`,
or `LOW`, keep current findings `Open`, and preserve originating IDs in the
form `<ROLE>-RNN-FNNN`. Do not add Requirement Coverage or Execution Plan
Validation as this role's semantic conclusion.
