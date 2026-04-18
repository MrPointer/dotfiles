# Design Spec Template

Adapt sections to fit the design's complexity. Simple designs don't need every section; complex designs may need additional ones.

```markdown
# Design Spec: <Topic>

## Problem Statement
<What problem this solves, who it affects, and why it matters>

## Constraints
- <Hard constraints: backward compatibility, performance budgets, tech stack restrictions>
- <Soft constraints: preferences, conventions, timeline considerations>

## Chosen Approach
<The selected approach and why it was chosen over alternatives>

## Alternatives Considered
| Approach | Tradeoffs | Why Not |
|----------|-----------|---------|
| <Alternative A> | <Genuine strengths and weaknesses> | <Specific reason it wasn't chosen> |
| <Alternative B> | <Genuine strengths and weaknesses> | <Specific reason it wasn't chosen> |

## Architecture
<How the pieces fit together — components, boundaries, relationships>

## Components
<What gets built, what gets modified — described at the architectural level, not implementation level>

## Data Flow
<How information moves through the system — inputs, transformations, outputs, storage>

## Error Handling
<What can go wrong and how it's handled — failure modes, recovery strategies, user-facing behavior>

## Scope
- **In scope**: <What this design covers>
- **Out of scope**: <What this design explicitly does not cover>
- **Future considerations**: <Things that were discussed but deferred>

## Key Decisions
<Decisions made during brainstorming that aren't obvious from the design itself — include rationale>

| Decision | Rationale |
|----------|-----------|
| <What was decided> | <Why — the reasoning, not just the conclusion> |

## Success Criteria
<How to know this is done and working correctly>
- <Criterion 1>
- <Criterion 2>
```
