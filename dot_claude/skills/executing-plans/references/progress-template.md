# Progress File Template

The progress file is the checkpoint mechanism. It must be updated after every meaningful step so execution can resume from any point.

```markdown
# Execution Progress: <Feature/Task Name>

## Plan
- **Plan location**: <path to plan file or directory>
- **Started**: <timestamp>
- **Last updated**: <timestamp>
- **Status**: in-progress | blocked | complete

## Tasks

| # | Task | Tests | Implementation | Status | Notes |
|---|------|-------|----------------|--------|-------|
| 01 | <task name or file> | done | done | done | |
| 02 | <task name or file> | done | blocked | blocked: test dispute | See disputes below |
| 03 | <task name or file> | skipped | done | done | No testable AC |
| 04 | <task name or file> | skipped | done | done | TDD skipped: no interfaces, tightly coupled to DB |
| 05 | <task name or file> | — | — | pending | Depends on 02 |

## Current State
<What the executor is currently doing or waiting on>

## Test Artifacts
<Map of tasks to their test file paths — the implementer needs these>

| Task | Test Files |
|------|-----------|
| <task 01> | tests/models/inverter_test.go |
| <task 02> | tests/api/handler_test.go, tests/api/middleware_test.go |

## Completed Artifacts
<Files created or modified by completed tasks — needed for relaying to dependent tasks>

| Task | Files |
|------|-------|
| <task 01> | pkg/models/inverter.go, pkg/models/types.go |

## Disputes
<Test disputes reported by implementers — batched for human resolution>

### Task 02: <task name>
- **Test**: test_returns_404_on_missing_resource (tests/api/handler_test.go)
- **Implementer's claim**: "Acceptance criterion says 'return error on missing resource' but doesn't specify HTTP status code. Test asserts 404 but 422 may be more appropriate for this API's conventions."
- **Resolution**: <pending | resolved: description>

## Failures
<Implementation failures that aren't disputes — the implementer couldn't make tests pass but doesn't claim the tests are wrong>

## Regressions
<Existing tests that broke during implementation>
```
