# Master Plan: Plan HTML Demo

## RFC Baseline

- **RFC Status**: Accepted

## Explicit Deviations

None

## Sub-Plans

| # | Sub-Plan | Depends On / Sequenced After | Model | Description |
|---|----------|------------------------------|-------|-------------|
| 01 | `01-foundation.md` | - | Mid-tier | Extract source assets and preserve self-contained output |
| 02 | `02-theme-control.md` | 01 | Cheapest | Add System / Light / Dark controls with no persistence |
| 03 | `03-review-fixture.md` | 02 | Mid-tier | Add stable fixture coverage and manual smoke checks |

## Review Summary

| Reviewer | Status |
|----------|--------|
| plan-clarity-reviewer | Passed |
| plan-executability-reviewer | Passed with concerns |

## Manual Smoke Checks

- Open the generated HTML in a browser.
- Switch between System, Light, and Dark.
- Expand and collapse all cards.
- Click each DAG node and confirm it opens the matching sub-plan.
- Confirm the reference review section remains readable.
