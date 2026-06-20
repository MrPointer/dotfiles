# Decomposition

Use this reference when decomposing direct feature requirements into executable sub-plans or defining cross-sub-plan contracts.

## Boundary Selection

Break work along natural seams: layers, domains, files, modules, independently verifiable outcomes, or integration boundaries.

Each sub-plan should be small enough for one focused execution session. If it feels too large, split it further.

## Dependency DAG

Minimize dependencies and make each edge explicit and one-directional. The dependency graph must form a valid DAG.

Before placing independent sub-plans in the same execution group, apply [concurrency-policy.md](concurrency-policy.md). For restricted build ecosystems, keep the DAG linear even when sub-plans are logically independent.

- No sub-plan may depend on information produced by a later sub-plan.
- No sub-plan may depend on information produced by a sub-plan in the same parallel group.
- Sub-plans cannot communicate at runtime; the lead agent relays results strictly along dependency edges.
- If decomposition requires bidirectional information flow, merge or restructure the boundaries.
- Same-group placement means logical independence only; the master plan must still record that parallel execution is allowed by policy, require isolated implementer worktrees for concurrent execution, or explicitly serialize the group.
- Linear policy sequencing is allowed for operational safety. Label policy-only sequencing separately from real data-flow dependencies in the master plan.

## Embedded Context

Each sub-plan must be self-contained. Include only execution-critical context:

- domain knowledge the agent cannot derive from code, such as requirements, business rules, config formats, protocol details, or constraints
- architectural decisions the user approved or the planner verified against the codebase
- exact cross-boundary contracts that other sub-plans depend on
- prerequisites, primary files, acceptance criteria, reviewer, required skills, and execution model

Do not include internal API design, function signatures within a package, private helpers, method bodies, step-by-step coding instructions, exact test/lint/build commands, or decisions already owned by skills.

The test: if changing something would break a different sub-plan's code, it is a contract and belongs in the plan. If it only affects code within one sub-plan, it is an internal execution decision.

## Cross-Boundary Contracts

When sub-plans run in parallel, consuming agents cannot discover producer output at execution time. The plan must specify the contract: interface/type signatures, files, data shapes, commands, or artifacts.

For sequential dependencies, the later agent can read the earlier sub-plan's actual output. A pre-specified contract is still required when the earlier output constrains later work.

Contracts must satisfy these integrity rules:

- **Caller annotations**: every new public method/function introduced by a sub-plan must specify its production caller. If the caller lives in a different sub-plan, both sides reference the contract. No caller means dead code.
- **Connected data flow**: every cross-sub-plan data path must trace source -> transport mechanism -> destination, with sub-plan ownership at each hop.
- **Interface boundary checks**: if a sub-plan adds a method to a concrete type, but consumers access that type through an interface, the plan must either add the method to the interface or assign concrete-type wiring to a specific sub-plan.

No orphan public methods. No prose-only data flow such as "X stores the value on config" when the consumer needs it delivered through a channel no sub-plan owns.

## Skills Are The Agent's Authority

List the skills each sub-plan requires, but do not replicate skill content into the plan.

Skills define how to write code, test, lint, build, and document. The plan defines what to build and why. If a skill says `use task test`, the plan should not say `go test ./...`. If a skill mandates a test style, the plan should not prescribe test structure.

If the feature has an active anchor or needs cross-session handoff, list `anchoring-context` only where feature-level context may need to be read or updated. Do not use it as a substitute for execution progress.

Verification and acceptance criteria should stay behavioral: `all tests pass`, `code builds`, `lints clean`, and feature-specific outcomes. The executing agent uses loaded operational skills to determine exact commands.

## Decisiveness

Before writing a design decision, verify it against the codebase.

If a decision references an interface, method, or file, confirm it exists and state its exact shape. Do not write hedges like `if X exists` or `either A or B`. Pick the approved approach or stop for clarification.

This is especially critical for sub-plans assigned to cheaper models, which cannot resolve ambiguity on their own.

## Documentation Sub-Plan Trigger

If the feature affects documented domain concepts, architecture, or business processes, add a final documentation sub-plan. Use [documentation-sub-plan.md](documentation-sub-plan.md).

Skip documentation planning only when the feature does not affect documented concepts, flows, or architecture; the only doc impact is component-level post-execution drift review; or no project documentation exists yet.
