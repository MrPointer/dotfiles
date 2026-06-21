# Sub-Plan: Fixture Coverage

## Objective

Provide a stable plan directory that can verify the renderer in automated tests and in manual browser smoke checks.

## Required Skills

- Test-driven development
- Managing chezmoi source files

## Primary Files

- `tests/fixtures/demo-plan/00-master.md`
- `tests/fixtures/demo-plan/01-foundation.md`
- `tests/fixtures/demo-plan/02-theme-control.md`
- `tests/fixtures/demo-plan/03-review-fixture.md`
- `tests/fixtures/demo-plan/reviews/plan-clarity-reviewer.md`

## Acceptance Criteria

- [ ] Stable fixture renders through the same CLI path as real plans.
- [ ] Manual smoke command is documented for browser checks.
- [ ] Automated tests assert representative snippets without a brittle full snapshot.
