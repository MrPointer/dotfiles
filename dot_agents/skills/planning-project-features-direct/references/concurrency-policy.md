# Concurrency Policy

Use this reference when deciding whether independent sub-plans may execute in parallel, when writing the master plan's execution order, or when recording build/cache requirements.

## Purpose

Logical independence is not enough to justify concurrent implementation. Some build systems are expensive in fresh worktrees, and some shared caches or output directories become contested resources when several agents compile or write binaries at the same time. Prefer predictable linear execution over theoretical parallelism when the ecosystem makes safe concurrency uncertain.

## Decide The Execution Mode

Inspect project documentation, build files, and the feature scope before assigning execution groups.

### Linear DAG Exception List

Use a Linear DAG by default only for these ecosystems:

| Ecosystem or feature scope | Execution mode |
|----------------------------|----------------|
| Rust / Cargo code implementation | Linear DAG |
| C or C++ code implementation with Make, CMake, Ninja, Bazel, Meson, Xcode, or similar build systems | Linear DAG |
| Swift / Xcode / SwiftPM code implementation | Linear DAG |
| JVM code implementation with Gradle, Maven, or SBT | Linear DAG unless project documentation explicitly defines safe parallel agent builds |
| .NET code implementation with MSBuild or `dotnet` build/test outputs | Linear DAG unless project documentation explicitly defines safe parallel agent builds |

Treat this table as the complete built-in exception list. Do not add languages to it by vibes, by generic "compiled" status, or because the project feels large.

### Default For Everything Else

For every ecosystem not listed above, keep parallel execution available when normal DAG, file ownership, and workspace-safety checks pass. Do not name unlisted ecosystems elsewhere to clarify the default; absence from this table is the signal.

Linearize an unlisted ecosystem only when project documentation, build configuration, or planning exploration reveals a concrete shared build/output/cache constraint. Record that as a `project-specific constraint`, not as a new built-in language rule.

If a feature includes code implementation in one of the Linear DAG exception ecosystems, linearize the whole feature plan by default. Do not place documentation, config, or test-only sub-plans in parallel with code sub-plans unless the user explicitly approves a project-specific mixed schedule.

Do not treat the existence of a build cache as proof that parallel execution is safe. Parallel execution is allowed only when the project documents or the planner verifies that concurrent workers use non-conflicting output paths or a safe shared cache mechanism.

## Linear DAG Rules

When the policy selects a Linear DAG:

- Keep sub-plans small and self-contained; do not collapse the feature into one monolithic plan.
- Put exactly one executable sub-plan in each execution group.
- Do not use `Parallel group` labels in the master plan.
- Do not invent data-flow dependencies. If an edge exists only to serialize execution, label it as policy-only sequencing in the master plan.
- Keep cross-sub-plan contracts only for real produced/consumed behavior, data, interfaces, files, or artifacts.
- Keep worker bindings and model assignments. Linear execution does not permit coordinator self-execution when a worker binding is assigned.

## Parallel-Allowed Rules

When the policy allows parallel execution, same-group sub-plans must still satisfy all normal DAG and execution-safety rules:

- no dependency on same-group output
- no overlapping primary file ownership
- explicit cross-boundary contracts for any producer/consumer relationship
- task-scoped implementer worktrees for concurrent implementation
- documented build/cache seeding or an explicit `None`
- serialization or user approval if isolation, cache safety, or output-path safety cannot be verified

## Required Master Plan Section

Every multi-sub-plan master plan must record the decision before `## Execution Order`:

```markdown
## Concurrency Policy

- **Decision**: Linear DAG | Parallel allowed
- **Reason**: <language/build-system/workspace rationale>
- **Linearization basis**: Rust/Cargo | C/C++ build system | Swift/Xcode/SwiftPM | JVM build system | .NET/MSBuild | project-specific constraint | None
- **Execution impact**: <one sub-plan per group, or which groups may run in parallel>
- **Override**: <explicit user-approved/project-documented exception, or None>
```

If the decision is `Linear DAG`, the `## Execution Order` section must contain only sequential groups.
