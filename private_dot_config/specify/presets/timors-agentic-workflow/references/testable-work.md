# Testable Work Policy

## Default

Tests are the default for behavior that can regress. Derive test tasks from user
story independent tests, interface contracts, quickstart scenarios, existing
test conventions, and the observable behavior being changed. Put tests and the
implementation they prove in the same execution group so a delegated group can
reach a verified outcome without another group's uncommitted work.

When practical, order a test task before its implementation task and require
the test to demonstrate the missing or incorrect behavior before implementation.
Use exact project-relative test paths in `tasks.md`.

## Exact Test Expectation

Every execution-group detail uses exactly one of these values:

- `Tests required`
- `Existing coverage: <evidence expected>`
- `Not applicable: <reason>`

Use `Tests required` when the group creates or changes testable behavior. Use
Existing coverage only when named current tests or a concrete verification
surface already proves the result and the group does not need a new test task.
Use Not applicable only for work with no meaningful automated test surface,
with a concrete reason tied to that work. Bare `None`, generic claims such as
"docs only", and preference-based omissions are invalid.

## Behavioral Checks

Story Independent Test, group Acceptance Criteria, and group Verification must
describe observable behavior or artifacts. Acceptance and Verification are
bullets, not task checkboxes. Commands alone are insufficient unless the
expected result is also stated. A build or lint command does not replace a
behavioral test when behavior changed.
