# RFC-Backed Epic Plan Template

The epic plan is a coordination document that decomposes a reviewed RFC into sequenced features. It tracks the overall effort, preserves RFC context, and gives each future feature planning session enough context to start cold.

```markdown
# Epic: <Name>

## RFC Baseline
- **RFC**: `<docs/rfcs/topic.md>`
- **RFC Status**: <Accepted | Draft approved for planning | other project status>
- **RFC Reviews**: `design-reviewer` <status>, `rfc-clarity-reviewer` <status or not requested>

## Summary
<What this epic accomplishes from the RFC — the end state when all features are complete>

## RFC Scope Mapping
| RFC Item | Epic Coverage |
|----------|---------------|
| <Goal / non-goal / constraint / success criterion / risk> | <Feature(s), cross-cutting concern, or explicit acceptance> |

## Explicit Deviations
<Approved deviations from the RFC. Use "None" when the epic does not deviate.>

## Cross-Cutting Concerns
<RFC decisions or patterns that span multiple features and must be consistent across them>
- <Concern>: <How it should be preserved across feature plans>

## Build/Workspace Constraints
<Use "None" when there are no known build-heavy isolated-worktree constraints. Otherwise, list ignored in-repository build/cache directories that future feature plans should make available in isolated worktrees before TDD or implementation dispatch.>
- <Relative build/cache dir>: <Why feature planners must preserve this constraint>

## Features

<!-- Status values: not-started | planning | planned | in-progress | done -->

### Feature 1: <Name>
**Status**: not-started
**Dependencies**: None
**RFC Coverage**: <RFC section(s), goal(s), risk(s), or success criterion covered>
**Skills**: <Skills the feature planner and executing agents will likely need>

<Rich description — must be enough for a feature planner starting cold>

**What & Why**: <What this feature accomplishes and why it is a distinct unit of work>

**Context**:
<RFC decisions, domain knowledge, relevant existing code/packages, architectural context>

**Build/Workspace Notes**:
<Known build-heavy constraints for this feature, or "None beyond epic-level constraints">

**Scope Boundaries**:
- Includes: ...
- Excludes (handled by other features or RFC non-goals): ...

---

<!-- Repeat for each feature -->

## Sequencing

<Visual or textual description of the feature dependency graph and parallel opportunities. Dependencies are existence-level only; do not define contracts.>

- **Independent**: Feature 1, Feature 3
- **After Feature 1**: Feature 2

## RFC Risk Coverage
| RFC Risk / Tradeoff | Epic Handling |
|---------------------|---------------|
| <risk> | <feature placement, sequencing constraint, cross-cutting concern, or explicit acceptance> |

## Review Summary
| Reviewer | Status | Notes |
|----------|--------|-------|
| `plan-rfc-fidelity-reviewer` | Pending | <review output path> |
| `design-reviewer` | Pending | <review output path> |
```
