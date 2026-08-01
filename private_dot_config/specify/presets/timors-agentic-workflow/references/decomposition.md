# Execution Decomposition Policy

## Artifact Ownership

Generate `tasks.md` and `execution-plan.md` as one validated pair.

- `tasks.md` exclusively owns task descriptions and checks: phase headings,
  story goals, independent tests, checkpoints, sequential `TNNN` records,
  story labels, exact paths, checkboxes, and task-local `[P]` markers.
- `execution-plan.md` exclusively owns execution groups, the global dependency
  graph and order, model tiers, group capabilities and skills, file ownership,
  build/cache seeding, data flows, integration contracts, execution risks, and
  project review roles.

Do not duplicate task descriptions in the execution plan. Do not put group IDs,
a global DAG, worker strategy, model selection, or contracts in `tasks.md`.
`[P]` says only that one task can run alongside another because it touches
different files and has no dependency on incomplete work. It does not imply
that their execution groups may run concurrently.

## Story-First Hybrid Decomposition

Start with one independently testable phase and execution group per user story.
Move work out of a story only at a concrete seam:

- isolate shared foundation that blocks multiple stories;
- isolate necessarily cross-cutting work after the stories it spans;
- split conflicting file ownership so each path has one group owner;
- split at an explicit producer/consumer contract or interface boundary;
- split when materially different required skills or capabilities apply; or
- split when each side has its own observable acceptance and verification.

Keep tests and the behavior they prove in the same group. Every task belongs to
exactly one group, and every group is delegated as one independently executable
unit with a semantic model tier. Documentation follows the separate policy in
`documentation-planning.md`.

## Group Records

Use sequential `EGNN` identifiers. `Covers` lists every owned task explicitly as
comma-separated `TNNN` values; ranges, phases, wildcards, and prose aliases are
invalid. Each group has nonempty Objective, Required Skills, Required
Capabilities, Execution Model rationale, Context, Primary Files, Acceptance
Criteria, and Verification. Only Prerequisites, Integration Contracts, and
Design Decisions may be `None`.

Required Skills is a nonempty bullet list of exact available skill IDs. Required
Capabilities is a bullet-list subset of `read`, `edit`, `shell`, `network`, and
`subagent-dispatch`. Every group requires `read`; every modifying group also
requires `edit`. Capabilities describe execution needs and do not provision a
runtime worker.

Primary Files uses project-relative paths and exactly projects to File
Ownership. Every path has one owner. Acceptance Criteria and Verification are
behavioral bullets, never task checkboxes.

## Cross-Group Contracts And Flows

Create one `DFNN` master row per cross-group data path. Give each contract a
globally unique `CTNN`. The source group records the flow once in Produces and
the destination group records it once in Consumes. Project these fields exactly:

- source and destination group ownership;
- Producer Contract and Consumer Contract IDs;
- destination group in Produces Consumers;
- explicit comma-separated Flow IDs; and
- byte-for-byte identical Produces `Shape / Artifact` and Consumes
  `Expected Shape / Artifact` strings.

Use only the canonical Produces, Consumes, and Interface Wiring tables from the
artifact-validation contract. Interface Wiring may be `None` when no wiring
applies. A group with no integration contract uses `None` for the entire field.

## Adoption

An existing upstream `tasks.md` is not adopted unless a conforming
`execution-plan.md` forms a valid pair with it. A missing, malformed, or
unsupported execution plan, or any invalid pair, requires rerunning
`/speckit.tasks`. There is no migration, repair, normalization, or inferred
protocol data.
