# Epic Plan Template

The epic plan is a coordination document that persists across sessions. It tracks the overall effort, sequences features, and provides rich context for each feature's planning session.

```markdown
# Epic: <Name>

## Summary
<What this epic accomplishes — the end state when all features are complete>

## Requirements
<Bullet list of confirmed requirements from Phase 1, refined after Phase 4>

## Scope
- **In scope**: ...
- **Out of scope**: ...

## Cross-Cutting Concerns
<Decisions or patterns that span multiple features and must be consistent across them>
- <Concern>: <How it should be handled across features>
...

## Features

<!-- Status values: not-started | planning | planned | in-progress | done -->

### Feature 1: <Name>
**Status**: not-started
**Dependencies**: None
**Skills**: <Skills the feature planner and executing agents will need>

<Rich description — must be enough for a feature planner starting cold>

**What & Why**: <What this feature accomplishes and why it's a distinct unit of work>

**Context**:
<Domain knowledge, spec references, relevant existing code/packages, architectural context>

**Scope Boundaries**:
- Includes: ...
- Excludes (handled by other features): ...

---

### Feature 2: <Name>
**Status**: not-started
**Dependencies**: Feature 1
**Skills**: ...

...

---

<!-- Repeat for each feature -->

## Sequencing

<Visual or textual description of the dependency graph and parallel opportunities>

- **Independent** (can be planned/executed in any order): Feature 1, Feature 3
- **After Feature 1**: Feature 2
- **After Features 1, 2**: Feature 4
...

## Risks & Mitigations
| Risk | Mitigation |
|------|------------|
| ...  | ...        |
```
