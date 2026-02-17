# Sub-Plan Template

Each sub-plan is a **self-contained execution unit**. An agent should be able to pick up a sub-plan and execute it without reading anything else.

```markdown
# Sub-Plan: <Task Name>

## Objective
<What this sub-plan accomplishes and why>

## Required Skills
<Skills the executing agent MUST load before starting>
- `skill-name` — reason it's needed

## Reviewer
<Local reviewer agent assigned during Phase 4, or "None" if no suitable reviewer was found>

## Execution Model
**Recommended**: Haiku | Sonnet | Opus
**Rationale**: <Why this model is appropriate for this sub-plan>

Examples:
- Haiku: "Standard CRUD implementation following existing patterns in the codebase"
- Sonnet: "Complex business logic with multiple edge cases and error handling scenarios"
- Opus: "Novel architectural approach requiring creative problem-solving" (rare)

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

## Implementation Steps
1. <Step with clear deliverable>
2. <Step with clear deliverable>
...

## Acceptance Criteria
- [ ] <Criterion 1>
- [ ] <Criterion 2>
...

## Post-Execution
If this project has component-level documentation, run the `component-docs-reviewer` agent to verify
component docs still match the actual implementation.
```
