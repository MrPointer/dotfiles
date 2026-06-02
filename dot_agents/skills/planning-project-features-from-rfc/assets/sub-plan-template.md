# RFC-Backed Sub-Plan Template

Each sub-plan is a self-contained execution unit derived from the RFC. An agent
should be able to execute it without reading the RFC, master plan, or other
sub-plans.

```markdown
# Sub-Plan: <Task Name>

## Objective

<What this sub-plan accomplishes, why it matters, and which RFC outcome it supports>

## RFC Context
- **RFC**: `<docs/rfcs/topic.md>`
- **Relevant RFC decisions**: <decisions this sub-plan must preserve>
- **Relevant constraints / non-goals**: <constraints and exclusions this sub-plan must respect>

## Required Skills

<Skills the executing agent MUST load before starting — both design skills and operational skills>

- `skill-name` — reason it's needed

## Execution Model

**Recommended**: Cheapest | Mid-tier | Most capable

**Rationale**: <Why this model tier is appropriate for this sub-plan>

Examples:

- Cheapest: "Rename-only task with all paths and reference updates listed"
- Mid-tier: "Straightforward configuration change following existing patterns"
- Mid-tier: "Shell/runtime change where regressions are costly"
- Most capable: "Documentation synthesis across final implementation state"

## Prerequisites

<What must exist before this sub-plan can be executed>

- <Specific file, interface, state, or output expected from a prior sub-plan — include actual signatures/shapes when relevant>
- Or: "None — this sub-plan has no dependencies"

## Context

<Essential context embedded directly: relevant interfaces, data shapes, conventions, architectural decisions, and RFC decisions the agent needs to know>

## Primary Files

<Files this sub-plan primarily creates or modifies — helps prevent conflicts during parallel worktree integration>

- `path/to/file.ext` (create | modify | delete)

## Integration Contracts

<Cross-sub-plan wiring this sub-plan participates in. Every new public method must have a production caller; every consumed method must be referenced. Omit this section entirely if this sub-plan has no cross-boundary contracts.>

**Produces**:

- `Type.Method()` — called by: sub-plan NN, in `consumer/location.go`

**Consumes**:

- calls: `Type.Method()` from sub-plan NN

**Interface wiring**: <If a produced method must be added to an interface for consumers to reach it, state which interface and who owns the change. If consumers use the concrete type directly, state that explicitly.>

## Design Decisions

<Only RFC-approved decisions or execution-local decisions the agent must follow. Focus on the "what" and "why," not the "how." Do not introduce new architecture here.>

Do NOT include method body implementations, exact code patterns, step-by-step coding instructions, or specific commands for testing/linting/building when skills cover those operations.

## Acceptance Criteria

- [ ] <Criterion 1>
- [ ] <Criterion 2>

<!-- Include a ## Post-Execution section only if this sub-plan modifies interfaces or behavior
     covered by component-level documentation. When included, instruct the executing agent to run
     the component-docs reviewer after implementation. -->
```
