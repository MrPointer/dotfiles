# Deterministic Artifact Validation Protocol

## Authority And Scope

This is the single authoritative installed source for execution-artifact
validation in protocol `0.1.0`. `speckit.analyze` and `speckit.implement` MUST
read this file from
`.specify/presets/timors-agentic-workflow/references/artifact-validation.md`.
They MUST fail closed on every violation and MUST NOT repair, normalize, infer,
or accept unknown protocol fields. Regenerate invalid execution artifacts with
`/speckit.tasks`.

Validation checks required artifacts, headings, tables, canonical columns,
fields, identifiers, controlled vocabulary, and cross-artifact mappings.

## Artifact And Identity Gate

1. `execution-plan.md` is required.
2. Its first heading is exactly `# Execution Plan: <feature>`, where `<feature>`
   is nonempty and identifies the planned feature.
3. It contains exactly one protocol line, exactly
   `**Preset Protocol Version**: 0.1.0`.
4. A missing, duplicate, malformed, or unsupported protocol version is invalid.
5. `spec.md`, `plan.md`, and `tasks.md` are required Planning Inputs. Every row
   marked `Yes` in Planning Inputs names a present, readable artifact at its
   project-relative Path. Optional inputs may be marked `No` only when they are
   not required for this feature.

## Required Master Sections

The execution plan contains each of these `##` headings exactly once:

1. `## Summary`
2. `## Planning Inputs`
3. `## Execution Groups`
4. `## Dependency Graph`
5. `## Concurrency Policy`
6. `## Execution Order`
7. `## File Ownership`
8. `## Build/Cache Seeding`
9. `## Cross-Group Data Flow`
10. `## Execution Risks`
11. `## Project Review Roles`
12. `## Execution Group Details`
13. `## Post-Execution`

## Canonical Master Tables

Each listed section contains exactly one table with the exact canonical header
below. Header spelling, order, and column count are normative; unknown columns
are invalid in protocol `0.1.0`.

### Planning Inputs

```text
| Artifact | Path | Required |
```

`Artifact` and `Path` are nonempty. `Required` is exactly `Yes` or `No`.

### Execution Groups

```text
| ID | Group | Covers | Depends On / Sequenced After | Model | Description |
```

`ID` is an `EGNN` identifier. `Group`, `Covers`, `Model`, and `Description` are
nonempty. `Model` is exactly one of `Cheapest`, `Mid-tier`, or `Most capable`.
`Covers` is an explicit comma-separated list of canonical task IDs; it contains
no ranges, phase references, wildcards, or implicit task groups.

### File Ownership

```text
| Group | Primary Files |
```

`Group` is an existing `EGNN`. `Primary Files` is a nonempty, comma-separated
projection of that group's detail `Primary Files` field. Each project-relative
path appears in exactly one group across the master table.

### Build/Cache Seeding

```text
| Relative Path | Applies To | Purpose | Notes |
```

All cells are nonempty project-relative or declared execution values. Use one
row per seeded build or cache location.

### Cross-Group Data Flow

```text
| Flow ID | Data | Source Group | Producer Contract | Transport | Destination Group | Consumer Contract |
```

`Flow ID` is a `DFNN` identifier. `Data`, `Transport`, source, destination, and
contract cells are nonempty.

### Execution Risks

```text
| Risk | Impact | Mitigation / Acceptance |
```

All cells are nonempty.

### Project Review Roles

```text
| Role ID | Packet Path | Model Tier | Applicability Rationale |
```

Use `None` as the sole table-body value only when no project review role
applies. Otherwise, every cell is nonempty and follows the role-packet rules.

## Identifiers And Task Coverage

- A task ID is exactly `TNNN`: `T` followed by three decimal digits, beginning
  with `T001` and increasing sequentially without gaps in `tasks.md` execution
  order.
- An execution-group ID is exactly `EGNN`: `EG` followed by two decimal digits.
- A data-flow ID is exactly `DFNN`: `DF` followed by two decimal digits.
- A contract ID is exactly `CTNN`: `CT` followed by two decimal digits. Every
  defined contract ID is unique across the execution plan.
- Every task ID is owned by exactly one Execution Groups `Covers` cell.
- Every `Covers` list is explicit comma-separated task IDs only. `T001-T003`,
  `Phase 1`, `all setup tasks`, and equivalent shorthand are invalid.

## Dependency, Graph, And Order Agreement

Each Execution Groups `Depends On / Sequenced After` cell is exactly `None` or
a comma-separated sequence of these canonical predecessor records:

```text
EGNN (logical dependency)
EGNN (policy-only sequencing)
```

The group must not name itself. The canonical predecessor records in Execution
Groups and Execution Order must agree exactly on predecessor ID and type.

Dependency Graph contains one strict Mermaid block using `flowchart LR`:

- every group is declared once with a quoted node label, as `EGNN["EGNN"]`;
- every predecessor is one edge from predecessor to successor;
- every edge is plain `EGNN --> EGNN`, with exactly one edge per predecessor
  relation;
- styles, subgraphs, nonstandard edge forms, and edge labels are invalid.

Predecessor type is recorded only in Execution Groups and Execution Order; it
is not encoded on Mermaid edges. The Mermaid graph contains every group and
relation recorded in Execution Groups and is acyclic.

Execution Order is an ordered sequence in exactly this shape:

```text
- **Sequential group N**: EGNN
  - **After**: None

- **Parallel group N**: EGNN, EGNN
  - **After**: EGNN (logical dependency), EGNN (policy-only sequencing)
```

`N` is the group number. Every `EGNN` appears exactly once in one sequential or
parallel group. Each group is immediately followed by its `After` line. `After`
is exactly `None` or the comma-separated canonical predecessor records above,
and exactly projects the predecessor records for the listed execution groups.
The resulting order is topological for every Mermaid edge; cycles, missing
groups, duplicate groups, unrecorded edges, and order disagreements are invalid.

## Execution Group Detail Fields

`## Execution Group Details` contains one `### EGNN: <group name>` section for
master row. Each detail contains these `####` headings exactly once, in this
order:

1. `#### Objective`
2. `#### Required Skills`
3. `#### Required Capabilities`
4. `#### Execution Model`
5. `#### Prerequisites`
6. `#### Context`
7. `#### Primary Files`
8. `#### Integration Contracts`
9. `#### Design Decisions`
10. `#### Test Expectation`
11. `#### Acceptance Criteria`
12. `#### Verification`

Within an execution group detail, only `Prerequisites`, `Integration Contracts`,
and `Design Decisions` may be exactly `None`. `Test Expectation` is never bare
`None`; it may use the exact Not applicable form below. Objective, Required
Skills, Required Capabilities, Execution Model, Context, Primary Files,
Acceptance Criteria, and Verification must be substantive and nonempty.

- `Required Skills` is a nonempty bullet list. Every bullet is one exact skill
  identifier, with no prose, alias, or explanatory suffix.
- `Required Capabilities` is a nonempty bullet list whose items are exactly
  `read`, `edit`, `shell`, `network`, or `subagent-dispatch`. Every group
  includes `read`; a group that changes files includes `edit`.
- The master Execution Groups `Model` field owns the tier and is exactly
  `Cheapest`, `Mid-tier`, or `Most capable`. Detail `Execution Model` contains
  only substantive rationale for that selected tier; it must not duplicate or
  equal a tier token.
- `Primary Files` is a nonempty comma-separated list of project-relative paths
  and exactly matches that group's File Ownership projection.
- `Test Expectation` is exactly one allowed form: `Tests required`,
  `Existing coverage: <evidence expected>`, or `Not applicable: <reason>`.
  The evidence and reason placeholders are substantive and nonempty.

## Review Role Packets

Selected core and project reviewer packets use this schema:

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

For every selected Project Review Roles row:

- `Role ID` matches lowercase kebab case: `^[a-z0-9]+(?:-[a-z0-9]+)*$`.
- `Packet Path` is exactly `.specify/reviewers/<role-id>.md` for that role ID.
- The selected packet's Role ID and Model Tier match the execution-plan row.
- The selected packet contains every required schema field and section above.
- `Change Triggers` is a comma-separated subset of exactly `baseline`,
  `tasks`, `decomposition-design`, `documentation`, `reviewer-config`, and
  `all`.

`None` is allowed in Project Review Roles only when no role applies; it is not a
placeholder for a missing or malformed role packet.

## Integration Contract Tables

When a group has integration contracts, `#### Integration Contracts` contains
the following canonical records. The record labels may render as
`**Produces**`, `**Consumes**`, and `**Interface Wiring**`. Unknown columns are
invalid.

**Produces**

```text
| Contract ID | Flow IDs | Shape / Artifact | Consumers |
```

**Consumes**

```text
| Contract ID | Flow IDs | Source | Expected Shape / Artifact |
```

**Interface Wiring**

```text
| Contract ID | Interface / Boundary | Owner | Wiring |
```

`Interface Wiring` may be exactly `None` when no wiring applies. All Contract
IDs and Flow IDs use the canonical forms above. Flow IDs are explicit
comma-separated `DFNN` values; no ranges or shorthand are allowed.

For every Cross-Group Data Flow row:

1. Its Flow ID appears exactly once in the source group's Produces table and
   exactly once in the destination group's Consumes table.
2. The source and destination group ownership match the master flow row.
3. The source Produces Contract ID matches `Producer Contract`; the destination
   Consumes Contract ID matches `Consumer Contract`.
4. The destination group appears in the source Produces `Consumers` cell.
5. The Produces `Shape / Artifact` and Consumes `Expected Shape / Artifact`
   strings are identical.
6. `Data` and `Transport` in the master flow row are present and nonempty.

Duplicate flow membership, mismatched ownership, contract IDs, consumers, or
shape strings are invalid.
