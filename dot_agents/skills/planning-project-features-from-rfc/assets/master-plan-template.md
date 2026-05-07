# RFC-Backed Master Plan Template

The master plan is the orchestration document. It decomposes an approved RFC into executable sub-plans. It does NOT redefine the RFC design.

```markdown
# Master Plan: <Feature Name>

## RFC Baseline
- **RFC**: `<docs/rfcs/RFC-0001-topic.md>`
- **RFC Status**: <Accepted | Draft approved for planning | other project status>
- **RFC Reviews**: `rfc-architect-reviewer` <status>, `rfc-risk-reviewer` <status>, `rfc-clarity-reviewer` <status or not requested>

## Summary
<Brief description of what this plan implements from the RFC>

## RFC Scope Mapping
| RFC Item | Plan Coverage |
|----------|---------------|
| <Goal / constraint / success criterion / risk> | <Sub-plan(s) or master-plan section covering it> |

## Explicit Deviations
<Approved deviations from the RFC. Use "None" when the plan does not deviate.>

## Sub-Plans
| #  | Sub-Plan | Depends On | Model | Description |
|----|----------|------------|-------|-------------|
| 01 | `01-<name>.md` | - | Cheapest | <What this sub-plan accomplishes> |
| 02 | `02-<name>.md` | 01 | Mid-tier | <What this sub-plan accomplishes> |

## Execution Order
<The execution order must form a valid DAG. Sub-plans in the same parallel group cannot depend on each other. Every dependency edge points from an earlier group to a later one.>

## Execution via Worker Agents
| Sub-Plan | Implementer Worker | Test Author Worker | Model Tier |
|----------|--------------------|--------------------|------------|
| 01 | `<tier>-<domain>-worker` | `<tier>-test-author-worker` | <tier> |

## File Ownership
| Sub-Plan | Primary Files |
|----------|---------------|
| 01 | <files this sub-plan creates/modifies> |

## Cross-Sub-Plan Data Flow
| Data | Source (Sub-Plan) | Transport | Destination (Sub-Plan) |
|------|-------------------|-----------|------------------------|
| <what> | 01 | <how it travels> | 02 |

## RFC Risk Coverage
| RFC Risk / Tradeoff | Plan Handling |
|---------------------|---------------|
| <risk> | <acceptance criteria, sequencing, rollback note, or explicit acceptance> |

## Review Summary
| Reviewer | Status | Notes |
|----------|--------|-------|
| `plan-rfc-fidelity-reviewer` | Pending | <review output path> |
| `plan-executability-reviewer` | Pending | <review output path> |

## Post-Execution
If this project has component-level documentation, run the `component-docs-reviewer` agent to verify component docs still match the actual implementation.
```
