# Progress File Template

The progress file is the checkpoint mechanism. It must be updated after every meaningful step so execution can resume from any point.

## Template Contents

- Plan
- Tasks
- Current State
- Execution Audit
- Test Evidence
- Completed Artifacts
- Exceptions / Blockers
- Failures
- Regressions

```markdown
# Execution Progress: <Feature/Task Name>

## Plan
- **Plan location**: <path to plan file or directory>
- **Started**: <timestamp>
- **Last updated**: <timestamp>
- **Status**: in-progress | blocked | complete
- **Review branch**: <branch or detached HEAD where execution started>
- **Execution base**: <commit SHA before execution started>
- **Integration branch**: <local agent/<plan-id> branch used for checkpoint commits>
- **Concurrency policy**: <Linear DAG | Parallel allowed; reason from master plan>
- **Execution schedule**: <serialized order or policy-allowed parallel groups>
- **Final review mode**: mixed reset to execution base, leaving aggregate changes dirty
- **Checkpoint range**: <execution-base>..<final checkpoint commit, once complete>

## Tasks

| # | Task | Tests | Implementation | Status | Notes |
|---|------|-------|----------------|--------|-------|
| 01 | <task name or file> | done | done | done | |
| 02 | <task name or file> | missing | blocked | blocked: missing test evidence | See exceptions below |
| 03 | <task name or file> | skipped | done | done | No testable AC |
| 04 | <task name or file> | explained | done | done | Existing tests cover behavior and were run |
| 05 | <task name or file> | — | — | pending | Depends on 02 |

## Current State
<What the executor is currently doing or waiting on>

## Execution Audit

| Task | Planned Worker | Actual Worker | Model / Effort | Dispatch Evidence | Implementation Workspace | Dirty-State Preflight | Build/Cache Reuse | Checkpoint Commit | Integration Status | Test Evidence | Verification |
|------|----------------|---------------|----------------|-------------------|--------------------------|-----------------------|-------------------|-------------------|--------------------|---------------|--------------|
| <task 01> | <worker from plan> | <worker actually used> | <model and effort> | <runtime command, subagent id, or reason not applicable> | <main workspace / worktree path / serialized: reason> | <clean / user-authorized dirty: paths / not applicable> | <shared cache configured / seeded: path / skipped: no safe strategy / none required / blocked: reason> | <commit SHA / pending / not applicable> | <pending / merged / blocked: reason / not applicable> | <tests added/updated / existing tests cover / skipped: reason / missing> | <focused + full command results / blocked: reason> |

## Test Evidence
<Map tasks to tests, verification results, or concrete no-test explanations>

| Task | Tests / Explanation | Verification |
|------|---------------------|--------------|
| <task 01> | tests/models/inverter_test.go | `go test ./pkg/models` passed |
| <task 02> | Existing handler tests cover the changed branch | `go test ./tests/api` passed |

## Completed Artifacts
<Files created or modified by completed tasks — needed for relaying to dependent tasks>

| Task | Files |
|------|-------|
| <task 01> | pkg/models/inverter.go, pkg/models/types.go |

## Exceptions / Blockers
<Missing test evidence, acceptance-criteria ambiguity, or implementation blockers batched for human resolution>

### Task 02: <task name>
- **Issue**: Acceptance criterion says "return error on missing resource" but does not specify status-code behavior.
- **Implementer's claim**: "Cannot choose a test expectation without deciding between 404 and 422."
- **Resolution**: <pending | resolved: description>

## Failures
<Implementation failures that are not acceptance-criteria ambiguities>

## Regressions
<Existing tests that broke during implementation>
```
