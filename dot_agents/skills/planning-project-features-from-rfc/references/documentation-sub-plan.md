# Documentation Sub-Plan

Use this reference when a feature affects documented domain concepts, architecture, or business processes.

Add documentation updates as the final sub-plan in execution order. This makes them visible, reviewable, and deliberate.

## What To Plan Upfront

| Doc Level | Planned Upfront? | Rationale |
|-----------|------------------|-----------|
| Domain docs | Yes | The planner knows what domain concepts are changing. |
| Architecture docs | Yes | The planner knows what structural changes are happening. |
| Process docs | Yes | The planner knows what flows are being added or modified. |
| Component docs | No, post-execution | Component docs describe implementation details that may drift from the plan. |

## How To Write The Sub-Plan

Use the standard sub-plan template, but make implementation steps doc edits rather than code changes.

Be specific about:

- existing docs to update, including file paths, sections, and stale claims to replace
- new docs to create, including paths, existing doc patterns to follow, and what the new doc covers
- structural pattern matching when existing docs follow a pattern
- required documenting skills
- execution model, which should always be the most capable available model

## Proportionality Guard

The target document's purpose, audience, and existing abstraction level outrank feature-local emphasis.

The sub-plan should say whether the feature should be mentioned as a primary concept, a small example, or only an implementation detail. Do not let a narrow feature become the center of a broad domain, architecture, or process doc unless the RFC explicitly changes that document's central subject.

## When To Skip

Skip the documentation sub-plan when:

- the feature does not affect documented concepts, flows, or architecture
- the only documentation impact is component-level and can be handled by post-execution component docs review
- no project documentation exists yet; recommend creating initial docs as a separate effort
