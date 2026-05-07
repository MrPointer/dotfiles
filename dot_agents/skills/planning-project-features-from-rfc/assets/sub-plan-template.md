# RFC-Backed Sub-Plan Template

Each sub-plan is a self-contained execution unit derived from the RFC. An agent should be able to execute it without reading the RFC or other sub-plans.

```markdown
# Sub-Plan: <Task Name>

## Objective
<What this sub-plan accomplishes and which RFC outcome it supports>

## RFC Context
- **RFC**: `<docs/rfcs/RFC-0001-topic.md>`
- **Relevant RFC decisions**: <decisions this sub-plan must preserve>
- **Relevant constraints / non-goals**: <constraints and exclusions this sub-plan must respect>

## Required Skills
- `skill-name` — reason it's needed

## Execution Model
**Recommended**: Cheapest | Mid-tier | Most capable
**Rationale**: <Why this model tier is appropriate>

## Prerequisites
- <Specific file, interface, state, or output expected from an earlier sub-plan>
- Or: "None — this sub-plan has no dependencies"

## Context
<Essential context embedded directly: relevant interfaces, data shapes, conventions, and RFC decisions>

## Primary Files
- `path/to/file.ext` (create | modify)

## Integration Contracts
**Produces**:
- `Type.Method()` — called by: sub-plan NN, in `consumer/location.go`

**Consumes**:
- calls: `Type.Method()` from sub-plan NN

**Interface wiring**: <Which interface changes are required, or state that consumers use the concrete type directly.>

<Omit this section entirely if this sub-plan has no cross-boundary contracts.>

## Design Decisions
<Only RFC-approved decisions or execution-local decisions the agent must follow. Do not introduce new architecture here.>

## Acceptance Criteria
- [ ] <Criterion 1>
- [ ] <Criterion 2>
```
