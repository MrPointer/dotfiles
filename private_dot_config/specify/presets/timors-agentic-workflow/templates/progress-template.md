# Execution Progress: <feature>

## Execution Context

- **Feature Directory**: <project-relative path>
- **Started**: <ISO-8601 timestamp>
- **Last Updated**: <ISO-8601 timestamp>
- **Status**: in-progress | blocked | complete
- **Review Branch**: <branch or detached HEAD>
- **Execution Base**: <commit SHA>
- **Integration Branch**: <local branch>
- **Concurrency Policy**: Linear DAG | Parallel allowed
- **Final Review Mode**: aggregate dirty diff
- **Checkpoint Range**: <base>..<tip> | Pending
- **Checkpoint Ref**: refs/agent-checkpoints/<feature> | Pending

## Execution Groups

| Group | Task IDs | Planned Model | Required Skills | Required Capabilities | Resolved Worker / Model | Workspace | Tests | Implementation | Integration | Status | Notes |
|---|---|---|---|---|---|---|---|---|---|---|---|
| EGNN | TNNN | Cheapest \| Mid-tier \| Most capable | <exact skill identifiers> | <exact capabilities> | pending | pending | pending | pending | pending | pending | |

Group `Status` is exactly `pending`, `binding`, `in progress`, `blocked: <reason>`, `ready for integration/checkpoint`, or `done`. `Tests`, `Implementation`, and `Integration` are exactly `pending`, `in progress`, `done`, `skipped: <reason>`, or `blocked: <reason>`.

## Current State

<Current group, transition, or reason execution is waiting.>

## Execution Audit

| Group | Tier Evidence | Skill Evidence | Capability Evidence | Dispatch Evidence | Dirty-State Preflight | Build/Cache Reuse | Checkpoint | Verification |
|---|---|---|---|---|---|---|---|---|
| EGNN | pending | pending | pending | pending | pending | pending | pending | pending |

## Test Evidence

| Group | Tests / Exception | Verification |
|---|---|---|
| EGNN | <test evidence or exact planned exception> | pending |

## Completed Artifacts

| Group | Flow ID | Contract ID | Files / Outputs | Result Summary | Checkpoint |
|---|---|---|---|---|---|
| EGNN | DFNN or None | CTNN or None | <project-relative paths or outputs> | <attributable result summary> | <commit SHA> |

Use `None` for `Flow ID` and `Contract ID` only when the output does not cross a group boundary.

## Checkpoints

| Group | Commit | Local Ref | State |
|---|---|---|---|
| EGNN | Pending | Pending | pending |

Checkpoint `State` is exactly `pending`, `retained`, `released`, or `missing`.

## Exceptions / Blockers

<Authorization, serialized parallel fallback, contract mismatch, ambiguity, or blocker.>

## Failures

<Implementation or integration failures.>

## Regressions

<Previously passing behavior that failed during verification.>

## Final Verification

<Required full verification, final checkpoint/ref evidence, aggregate preparation, and outcome.>
