# Execution Concurrency Policy

## Decision Procedure

Record exactly one decision: `Linear DAG` or `Parallel allowed`. Logical
independence is not sufficient for parallel execution. Groups may share a
parallel execution-order record only when all of these are verified:

- no predecessor relation exists between them;
- Primary Files do not overlap;
- isolated execution workspaces are available; and
- generated output, build state, and caches are safe for concurrent access.

When the plan allows parallel groups but the active runtime cannot verify
isolation or output/cache safety, the runtime may serialize dispatch. That
operational fallback does not change the recorded graph or invent a policy-only
edge.

## Mandatory Linear DAG Cases

Use a Linear DAG for these ecosystems unless project documentation or the user
provides a concrete, recorded safe exception:

- Rust with Cargo;
- C or C++ build systems;
- Swift with Xcode or SwiftPM;
- JVM work when safe independent output/cache behavior is absent;
- .NET work when safe independent output/cache behavior is absent; or
- a concrete project-specific build, output, cache, workspace, or integration
  constraint that requires serialization.

Do not linearize other work conservatively without evidence. When Parallel
allowed is safe, represent its real fork and join rather than flattening it.

## Canonical Relations

Every predecessor record is exactly `EGNN (logical dependency)` or
`EGNN (policy-only sequencing)`. Logical dependencies carry required outputs or
ordering inherent in the work. Policy-only sequencing records serialization
required solely by the concurrency policy. Use policy-only edges for each
successive group in a Linear DAG when no logical edge already provides the
required sequence.

The Mermaid graph uses `flowchart LR`, one `EGNN["EGNN"]` declaration per
group, and one plain `EGNN --> EGNN` edge per predecessor. It contains no edge
labels, styles, or subgraphs. It is an exact, acyclic projection of all master
predecessor records.

Execution Order lists each group exactly once in a Sequential or Parallel group
record. Its immediately following `After` line exactly projects the canonical
predecessors for every group in that record. Groups sharing a Parallel record
must have the same projected predecessor set and cannot depend on each other.
Use only these record shapes:

```text
- **Sequential group N**: EGNN
  - **After**: None

- **Parallel group N**: EGNN, EGNN
  - **After**: EGNN (logical dependency), EGNN (policy-only sequencing)
```

## Build And Cache Evidence

Record each relevant ignored in-repository build or cache path with its affected
groups, purpose, and seeding or safety notes. Seeding prepares isolated
workspaces; it never makes an unsafe graph parallel. A project with no cache to
seed still records a substantive disposition permitted by the canonical table
rather than omitting the required section.
