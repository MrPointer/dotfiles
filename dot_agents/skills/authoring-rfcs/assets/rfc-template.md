# RFC Template

Use this template for RFCs that need to be precise enough for human engineering review and later planning. Adapt section names to project conventions, but preserve the distinction between verified current state, proposed design, tradeoffs, and planning handoff context.

Remove sections that truly do not apply. Do not leave placeholders in the final document.

Write the RFC as a normal engineering design document. Do not refer to "the user," "the agent," "the model," or the prompting/review workflow in body prose. Recast conversation-derived input as requirements, constraints, decisions, assumptions, or source references.

```markdown
# RFC-0000: <Title>

Status: Draft
Revision: R0
Last Updated: <YYYY-MM-DD>

## Review Record
| Reviewer | Scope | Status | Notes |
|----------|-------|--------|-------|
| `rfc-architect-reviewer` | Architecture, boundaries, contracts, current-state fit, planning readiness | Pending | <Review output path or summary> |
| `rfc-risk-reviewer` | Technical risks, migration, compatibility, rollback, hidden complexity | Pending | <Review output path or summary> |
| `rfc-clarity-reviewer` | Clarity and actionability | Not requested | <Reason if omitted, review output path if used> |

Valid review statuses: `Pending`, `Passed`, `Passed with concerns`, `Blocking`, `Not requested`.

## Summary
<One to three paragraphs describing the chosen design, the problem it solves, and the most important architectural consequence.>

## Problem
<The problem, who or what it affects, why it matters now, and why the current system does not already solve it.>

## Goals
- <Specific outcome the design must achieve>
- <Specific outcome the design must achieve>

## Non-Goals
- <Explicitly out-of-scope behavior, system, migration, or capability>
- <Explicitly out-of-scope behavior, system, migration, or capability>

## Constraints
- <Hard constraint: compatibility, persisted data, external contract, performance, security, platform, operational limit>
- <Soft constraint: convention, maintainability preference, team workflow, rollout preference>

## Current State
<Verified description of the relevant current architecture. Separate current facts from proposed changes. Name actual components, files, interfaces, configuration surfaces, data stores, and runtime flows where relevant.>

## Chosen Approach
<The committed design direction and why it is the right tradeoff for the known constraints. This section should not contain unresolved alternatives.>

## Decision Summary
| Decision | Rationale | Consequence |
|----------|-----------|-------------|
| <Settled architectural decision> | <Why this decision is appropriate for the project> | <What this means for the design or future plan> |

## Proposed Architecture
<How the system will be structured after this design. Describe boundaries, ownership, responsibilities, and how the changed pieces fit with the existing architecture.>

## Components And Responsibilities
| Component / Boundary | Responsibility | Current / New / Modified | Notes |
|----------------------|----------------|--------------------------|-------|
| <Actual component, module, service, command, config surface, or document area> | <What it owns after the design> | <Current/New/Modified> | <Important boundary, dependency, or ownership note> |

## Contracts And Interfaces
<Public APIs, CLI behavior, config keys, data schemas, events, file formats, protocol boundaries, or other contracts that the design depends on or changes. Omit if no contract changes exist.>

## Data And State
<Data ownership, lifecycle, persistence, migration concerns, cache/state invalidation, and compatibility implications. Omit if no durable or meaningful state exists.>

## Control And Data Flow
<How requests, commands, events, or data move through the system after the design. Include the main successful path and any important alternate paths.>

## Failure Modes And Recovery
| Failure Mode | Expected Behavior | Recovery / User Impact |
|--------------|-------------------|------------------------|
| <What can go wrong> | <How the system responds> | <How the system recovers or what affected actors see> |

## Security, Privacy, And Permissions
<Authentication, authorization, secret handling, data exposure, local permissions, or privacy implications. Omit if not relevant.>

## Operations And Observability
<Logging, metrics, diagnostics, rollout visibility, supportability, or operator workflows. Omit if not relevant.>

## Compatibility And Migration
<Backward compatibility, data migration, config migration, rollout constraints, or upgrade/downgrade behavior. Omit if not relevant.>

## Alternatives Considered
| Alternative | Strengths | Why Rejected |
|-------------|-----------|--------------|
| <Genuine alternative> | <Why it was plausible> | <Specific reason it was not chosen> |

## Risks And Tradeoffs
| Risk / Tradeoff | Impact | Mitigation Or Acceptance |
|-----------------|--------|--------------------------|
| <Risk or tradeoff introduced by the chosen design> | <What could happen> | <How the design handles it, or why it is acceptable> |

## Success Criteria
- <Design-level outcome that shows the chosen approach solves the problem>
- <Design-level outcome that planning can later turn into acceptance criteria>

## Open Questions
- <Only include non-blocking questions or decisions intentionally deferred beyond this design. Remove this section if none remain.>

## Planning Handoff
<Information the planning skill must preserve when decomposing the work: architectural boundaries, contracts that must stay stable, dependencies between design areas, and decisions that must not be reopened without new information. Do not include tasks or execution sequence.>

## Source References
| Source | What It Confirms |
|--------|------------------|
| `<path-or-doc>` | <Current architecture, contract, constraint, or behavior verified from this source> |
| Requirement or decision input | <Constraint, decision, or assumption that cannot be derived from code> |
```
