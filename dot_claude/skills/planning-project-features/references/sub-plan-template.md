# Sub-Plan Template

Each sub-plan is a **self-contained execution unit**. An agent should be able to pick up a sub-plan and execute it without reading anything else.

```markdown
# Sub-Plan: <Task Name>

## Objective
<What this sub-plan accomplishes and why>

## Required Skills
<Skills the executing agent MUST load before starting — both design skills (coding conventions, patterns) and operational skills (testing, linting, building)>
- `skill-name` — reason it's needed

## Reviewer
<Local reviewer agent assigned during Phase 4, or "None" if no suitable reviewer was found>

## Execution Model
**Recommended**: Cheapest | Mid-tier | Most capable
**Rationale**: <Why this model tier is appropriate for this sub-plan>

Examples:
- Cheapest: "Standard CRUD implementation following established patterns in the codebase"
- Mid-tier: "Complex business logic with multiple edge cases and error handling scenarios"
- Most capable: "Novel architectural approach requiring creative problem-solving" (rare)

## Prerequisites
<What must exist before this sub-plan can be executed>
- <Specific file, interface, or state expected from a prior sub-plan — include the actual signatures/shapes, not just references>
- Or: "None — this sub-plan has no dependencies"

## Context
<Essential context embedded directly — relevant interfaces, data shapes, conventions, architectural decisions the agent needs to know>

## Primary Files
<Files this sub-plan primarily creates or modifies — helps prevent conflicts in parallel execution>
- `path/to/file.ext` (create | modify)
- `path/to/other.ext` (modify)

## Integration Contracts
<Cross-sub-plan wiring this sub-plan participates in. Every new public method must have a production caller; every consumed method must be referenced.>

**Produces** (new methods/types other sub-plans consume):
- `Type.Method()` — called by: sub-plan NN, in `consumer/location.go`

**Consumes** (methods/types from other sub-plans this sub-plan calls):
- calls: `Type.Method()` from sub-plan NN

**Interface wiring**: <If a produced method must be added to an interface for consumers to reach it, state which interface and who owns the change. If consumers use the concrete type directly, state that explicitly.>

<Omit this section entirely if this sub-plan has no cross-boundary contracts.>

## Design Decisions
<Decisions the executing agent must follow — things it wouldn't arrive at independently from skills + codebase alone. Focus on the "what" and "why," not the "how.">

1. **<Decision name>**: <What to do and why — architectural choices, library selections, behavioral requirements, constraints not covered by skills>
...

Do NOT include: method body implementations, exact code patterns, step-by-step coding instructions, specific commands for testing/linting/building (skills handle these). If the agent can derive it from the loaded skills + codebase + acceptance criteria, it doesn't belong here.

## Acceptance Criteria
- [ ] <Criterion 1>
- [ ] <Criterion 2>
...

<!-- Include a ## Post-Execution section only if this sub-plan modifies interfaces or behavior
     covered by component-level documentation. When included, instruct the executing agent to run
     the `component-docs-reviewer` agent to verify component docs still match the implementation. -->
```
