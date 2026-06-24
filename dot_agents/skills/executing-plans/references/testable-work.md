# Testable Work

Use this reference when a task has testable acceptance criteria. It keeps test discipline inside the implementation task instead of creating a separate test-writing lane.

## Purpose

Testable behavior changes should not become "code first, maybe tests later" work. The implementer owns both tests and code, has the full task context, and uses the required testing skills to decide the right level of coverage.

This reference is intentionally light. It sets expectations and result checks; the required testing skills define the detailed mechanics.

## Default Expectation

For testable behavior changes, tell the implementer to:

- load the required testing skills, including `test-driven-development` when the sub-plan lists it
- make a test-first attempt when practical
- implement until the focused tests and required verification pass
- report tests added or updated, verification results, and any reason test-first work or new tests were not practical

Do not require a separate RED/GREEN ceremony in the plan unless a required testing skill asks for it. The worker's result only needs enough evidence for the coordinator to know that acceptance criteria were tested or consciously exempted.

## Skip Or Exception Reasons

Record a short reason when test-first work or new tests are not practical, such as:

- documentation-only, move-only, or formatting-only work
- configuration changes with no practical local test surface
- behavior already covered by existing tests that were run and reported
- area declared untestable by project docs
- acceptance criteria too ambiguous to test without user clarification

Do not spend execution time designing new testability architecture unless the sub-plan explicitly includes that work.

## Coordinator Check

Before accepting a testable implementation result, check for one of these:

- tests were added or updated and relevant verification passed
- existing tests cover the acceptance criteria and were run
- the implementer gave a concrete no-test or no-test-first reason

If none is present, send the task back once through the same implementer binding. If the result still lacks evidence or explanation, block it as missing test evidence.
