# Coordinated Analysis And Approval Protocol

## Authority And Coordinator Boundary

`/speckit.analyze` is a narrow coordinator. It owns compatibility and shared
structural preflight, role selection, dispatch, report-schema validation,
cumulative report materialization, deterministic source-preserving
aggregation, conflict exposure, verdict computation, and recording explicit
human decisions. Semantic conclusions belong only to selected reviewers in
fresh read-only contexts. The coordinator MUST NOT perform semantic review,
repair reviewer output, add a broad fallback pass, or review in its own context.

Planning artifacts and packets are read-only during analysis. Intended writes
are `specs/<feature>/reviews/<role-id>.md` and `specs/<feature>/analysis.md`.

## Preflight Order

Before extension hooks, prerequisite scripts, dispatch, or writes, run the
mandatory compatibility preflight from `protocol-compatibility.md`. After the
upstream before-hook and prerequisite script seams, but before dispatch or
report writes, read and execute the complete deterministic protocol in
`.specify/presets/timors-agentic-workflow/references/artifact-validation.md`.
That file is the sole source of structural task, execution-plan, and selected
packet invariants. Do not duplicate, infer, normalize, or repair those rules.
Any violation blocks semantic dispatch.

## Roles And Packets

The mandatory core roles are:

| Role ID | Model Tier | Change Triggers |
|---------|------------|-----------------|
| `artifact-fidelity` | Mid-tier | `baseline,tasks,decomposition-design,documentation,all` |
| `decomposition-design` | Most capable | `baseline,tasks,decomposition-design,reviewer-config,all` |
| `plan-clarity` | Mid-tier | `baseline,tasks,decomposition-design,documentation,reviewer-config,all` |

Every core or selected project packet uses the foundation packet schema:

```markdown
# Reviewer Role: <name>

- **Role ID**: <lowercase-kebab-case>
- **Model Tier**: Cheapest | Mid-tier | Most capable
- **Change Triggers**: <comma-separated subset>

## Scope
## Required Inputs
## Exclusions
## Review Checks
## Output Contract
```

The closed trigger vocabulary is `baseline`, `tasks`,
`decomposition-design`, `documentation`, `reviewer-config`, and `all`.
Project packets live at `.specify/reviewers/<role-id>.md`, cannot replace a
core packet, cannot request edits, and may require inputs only within the
project and feature directory. Their tier must match `execution-plan.md`.
Missing or malformed core or selected project packets block.

Core role questions, inputs, checks, exclusions, and role-specific output are
normative in the installed packets. Artifact fidelity owns baseline
traceability and `Requirement Coverage`; decomposition design owns technical
execution feasibility and `Execution Plan Validation`; plan clarity owns
cold-reader communication and emits only the common report.

## Reviewer Binding And Real Dispatch

The coordinator requires native `subagent-dispatch`. Each reviewer binding
must establish its planned tier, a concrete model selected through the native
runtime mechanism, read access to every required input, fresh-context
isolation, invokability, and attributable structured-result support.

Resolve each role independently. Project-local eligible bindings take
precedence over global bindings. Within that scope select, in order:

1. an explicit role-to-candidate preference supplied to `/speckit.analyze`;
2. one uniquely advertised preferred/default eligible reviewer at the tier;
3. the only eligible reviewer at the tier;
4. explicit human choice among equivalent eligible reviewers.

Use runtime metadata and configuration to verify tier source, concrete model
binding, input access, fresh-context support, and attribution support without
invocation. Never guess provider preference or silently change tier. A missing
binding or unresolved choice blocks.

There is no probe, canary, calibration task, or sacrificial call. The first
invocation is the real assigned review. Retry only when runtime evidence proves
no invocation started. A started or ambiguous invocation with no attributable
result becomes `Recovery required` and cannot be automatically redispatched.
Concrete provider/model, candidate, and dispatch IDs remain transient.

## Cumulative Role Reports

Each selected role owns one tracked cumulative report at
`specs/<feature>/reviews/<role-id>.md`. Reruns append rounds and never rewrite
earlier rounds. Identity and round records are exactly:

```markdown
# Review: <role name>

- **Role ID**: <role-id>
- **Model Tier**: <planned tier>
- **Latest Round**: RNN
- **Latest Verdict**: Passed | Passed with concerns | Blocking
- **Reviewed At**: <ISO-8601 timestamp>

## Scope

## Rounds

### RNN

- **Analysis Run**: ARNN
- **Trigger**: initial | baseline | tasks | decomposition-design | documentation | reviewer-config | all | explicit | <comma-separated union>
- **Inputs**: <project-relative artifact paths>
- **Planned Tier Verified**: Yes
- **Fresh Context Verified**: Yes
- **Attributable Result Verified**: Yes
- **Verdict**: Passed | Passed with concerns | Blocking

| ID | Category | Severity | Location | Summary | Recommendation | Resolution |
|----|----------|----------|----------|---------|----------------|------------|

#### Prior Finding Dispositions

| Prior Finding ID | Disposition | Current Finding ID | Notes |
|------------------|-------------|--------------------|-------|
```

All three verification attestations are exactly `Yes` in a valid round and are
written only after the transient evidence is observed. Missing evidence blocks
the round. Every round is a complete current snapshot, not a delta. Its Findings
table contains only current open findings; use `None` when there are none.

Finding IDs originate as `<ROLE>-RNN-FNNN`, remain stable when retained, and
are never reused. Severity is exactly `CRITICAL`, `HIGH`, `MEDIUM`, or `LOW`;
current finding Resolution is `Open`. Prior Finding Dispositions accounts for
every finding in the preceding round with exactly `Retained`, `Resolved`, or
`Superseded`. A retained row references the same ID, a resolved row references
`None`, and a superseded row references every replacement ID. New findings
need no disposition row.

Artifact-fidelity rounds add:

```text
| Requirement | Covered | Task IDs | Notes |
```

Decomposition-design rounds add:

```text
| Check | Status | Evidence / Notes |
```

The role-specific output is additive and exhaustive. Omission of a required
table or presentation of another core role's table as the role's semantic
conclusion is invalid. Artifact-fidelity coverage values are `Yes`, `Partial`,
or `No`; decomposition validation values are `Passed`, `Concern`, or
`Blocking`.

A current `CRITICAL` or `HIGH` finding makes the role verdict `Blocking`.
Only current `MEDIUM` or `LOW` findings make it `Passed with concerns`. No
current findings makes it `Passed`.

## Change Declarations And Rerun Selection

The first run uses declaration `initial` and dispatches all three core roles
and every selected project role. Every later run requires an explicit human
declaration. Version-control inspection may suggest but cannot narrow it.
Unknown or uncertain scope becomes `all`; multiple classes use the union. A
human may broaden but not narrow the required union.

- `baseline`: behavioral or architectural changes to the specification, plan,
  design artifacts, quickstart, or constitution; rerun all core roles.
- `tasks`: non-documentation task identity, scope, story mapping, path, test
  intent, or checklist text; rerun all core roles.
- `decomposition-design`: groups, dependencies, concurrency, ownership,
  contracts, data flow, models, skills, capabilities, acceptance criteria, or
  verification; rerun all core roles.
- `documentation`: documentation-impact text or documentation-only task changes
  that do not alter behavior, architecture, execution records, dependencies,
  ownership, contracts, acceptance criteria, or verification; rerun
  artifact fidelity and plan clarity.
- `reviewer-config`: selected roles, packet content, packet tier, or required
  inputs; rerun decomposition design, plan clarity, and every added or changed
  project role.
- `all`: rerun every core and currently selected project role.

Documentation-only task changes declare `documentation`, not `tasks`. A
documentation change that also changes an execution-group record declares
`decomposition-design` too. For every declaration, add every currently selected
project role whose packet declares that class or `all`.

A newly selected role always receives an initial round. A removed role is
excluded from the next aggregate but its historical report is preserved. A
changed packet or tier requires a round under its current contract. A selected
role without a valid current-contract round blocks. Core packet changes require
a protocol release and rerun the changed core role plus the roles selected by
the feature declaration.

## Runs, Resume, And Recovery

Analysis run IDs are monotonic feature-local `ARNN`. Before dispatch, write the
aggregate with `Aggregate Verdict: Incomplete` and `Run Status: Incomplete`,
the explicit declaration, selected roles, retained rounds, required reruns, and
the existing Approval Record unchanged. A newer run cannot begin while an
incomplete run exists.

Resume the same run. Retain valid successful rounds and dispatch only unfinished
or invalid roles not marked `Recovery required`. If inputs change during the
run, the human expands its declaration and every newly affected retained role
is invalidated for this run. Failed, malformed, or interrupted work remains
Incomplete.

When an invocation may have started but has no attributable result, durably set
that role row to `Recovery required`; do not store its concrete dispatch ID.
Only an explicit human choice may clear it: recover the original invocation,
confirm from runtime evidence that it never started, or accept abandonment risk
and authorize redispatch of the real review.

A run is Complete only when every currently selected role has exactly one valid
applicable round and the aggregate references exactly those rounds.

## Aggregate Schema

`analysis.md` has exactly these required sections and canonical records:

```markdown
## Review Summary

- **Aggregate Verdict**: Incomplete | Passed | Passed with concerns | Blocking
- **Analysis Run**: ARNN
- **Run Status**: Incomplete | Complete
- **Change Declaration**: initial | baseline | tasks | decomposition-design | documentation | reviewer-config | all | <comma-separated union>
- **Selected Roles**: <comma-separated role IDs>
- **Analyzed At**: <ISO-8601 timestamp>

## Role Verdicts

| Role ID | Report | Applicable Round | Verdict | Run Status |
|---------|--------|------------------|---------|------------|

## Findings

| Aggregate ID | Source Finding IDs | Category | Severity | Location | Summary | Recommendation | Resolution |
|--------------|--------------------|----------|----------|----------|---------|----------------|------------|

## Reviewer Conflicts

| ID | Source Finding IDs | Conflict | Required Resolution | Status |
|----|--------------------|----------|---------------------|--------|

## Requirement Coverage

| Source Role / Round | Requirement | Covered | Task IDs | Notes |
|---------------------|-------------|---------|----------|-------|

## Execution Plan Validation

| Source Role / Round | Check | Status | Evidence / Notes |
|---------------------|-------|--------|------------------|

## Explicit Deviations

| ID | Aggregate / Conflict ID | Source Finding IDs | Decision | Status | Rationale |
|----|-------------------------|--------------------|----------|--------|-----------|

## Approval Record

- **Decision**: Pending | Approved | Changes requested | Approved with deviations
- **Recorded By**: Human
- **Recorded At**: <ISO-8601 timestamp or Pending>
- **Rationale**: <decision rationale or Pending>
```

`Aggregate Verdict` is `Incomplete` exactly when Run Status is `Incomplete`.
Completed runs use only `Passed`, `Passed with concerns`, or `Blocking`. Role
Run Status is exactly `Rerun`, `Retained`, `Pending`, `Failed`, or
`Recovery required`. Pending, Failed, and Recovery-required rows use
`Applicable Round: None` and the matching sentinel as Verdict. Complete runs
contain only Rerun or Retained rows with a concrete applicable `RNN` and final
role verdict.

## Deterministic Aggregation And Conflicts

The Findings and Reviewer Conflicts tables contain only current records.
Aggregate IDs are monotonic non-reused `AGG-FNNN`; conflict IDs are monotonic
non-reused `CON-FNNN`. Preserve an ID only while the exact source-ID set
represents the same defect or conflict; retire it when resolved or regrouped.

Deduplicate only findings that identify the same artifact location, defect, and
required correction. Preserve every source finding ID and use the highest
source severity. Any disagreement about severity, diagnosis, or correction is
not deduplicated; create a conflict. Do not resolve semantic disagreements.

Finding and conflict status is `Open` or `Accepted deviation`. Project the
latter only while a matching active accepted deviation exists. Any open
conflict is blocking. Aggregate severity remains `CRITICAL`, `HIGH`, `MEDIUM`,
or `LOW`. Coverage is `Yes`, `Partial`, or `No`; validation is `Passed`,
`Concern`, or `Blocking`.

Compute the verdict:

- `Incomplete` while the run is incomplete;
- `Blocking` when any selected role is Blocking or any conflict is unresolved;
- `Passed with concerns` when none block and at least one role has concerns;
- `Passed` only when every selected role Passed and no conflict is unresolved.

## Deviations And Human Approval

Deviation IDs are monotonic non-reused `DEV-FNNN`. Each row points to one
current aggregate finding or conflict and its exact current source-ID set.
Decision is exactly `Accepted` or `Rejected`; Status is exactly `Active`,
`Resolved`, or `Superseded`. Decisions and rationale are human-owned. The
coordinator may only set deterministic lifecycle status: resolved sources make
a deviation Resolved; regrouping or a changed source set makes it Superseded
and requires a new human-owned active deviation.

`Approval Record / Decision` is the sole authorization field and is exactly
`Pending`, `Approved`, `Changes requested`, or `Approved with deviations`.
Preserve the record across reruns, including incomplete runs. Do not store
hashes or fingerprints and do not automatically revoke or reset approval.

| Aggregate State | Decision That Authorizes Implementation |
|-----------------|-----------------------------------------|
| Incomplete analysis run | None |
| `Passed` | `Approved` |
| `Passed with concerns` and no blocking finding/conflict | `Approved`, or `Approved with deviations` when at least one active accepted deviation exists |
| `Blocking` with every current blocking finding/conflict covered by an active accepted deviation | `Approved with deviations` |
| `Blocking` with any uncovered blocking finding/conflict | None |

`Pending` and `Changes requested` never authorize. `Approved with deviations`
requires at least one active accepted deviation. Plain `Approved` never
authorizes Blocking. Incomplete state, uncovered blockers, and uncovered
conflicts never authorize implementation.
