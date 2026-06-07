# Structural TDD

Use this reference when a task has testable acceptance criteria and the executor may separate test authoring from implementation.

## Contents

- [Purpose](#purpose)
- [Required Skill Dependency](#required-skill-dependency)
- [Skip Gates](#skip-gates)
- [Isolation Gate](#isolation-gate)
- [Build And Compilation Readiness](#build-and-compilation-readiness)
- [Test Author Context](#test-author-context)
- [TDD Quality Gate](#tdd-quality-gate)
- [Returning Test Files](#returning-test-files)
- [Implementer Contract With Structural TDD](#implementer-contract-with-structural-tdd)

## Purpose

Structural TDD prevents implementers from writing tests that rationalize their own implementation. The test author receives acceptance criteria only. The implementer receives the full task plus the tests. The `test-driven-development` skill defines valid RED behavior and test-quality mechanics; this reference defines isolation, handoff, and execution coordination.

Use this structure only when the code surface and runtime can support it. If not, skip structural TDD and implement against the acceptance criteria directly.

## Required Skill Dependency

Structural TDD depends on the `test-driven-development` skill for valid RED behavior and test-quality mechanics. Load it whenever this reference is used for test authoring, RED validation, or test-quality review.

If the active runtime cannot load or pass the `test-driven-development` skill to the test author or executor, do not weaken this workflow silently. Record the failure in progress and ask whether to fix skill loading or explicitly skip structural TDD for the affected task.

## Skip Gates

Skip test authoring when the task has no testable work, such as documentation updates, pure file moves, configuration-only changes, or no acceptance criteria.

For tasks with testable acceptance criteria, check whether the code surface supports TDD:

1. Read project docs first: `AGENTS.md`, project documentation, and task-level notes.
2. If the area is declared untestable, skip immediately and record the documented reason.
3. If docs do not answer it, use the active runtime adapter's cheapest suitable exploration path to answer one yes/no question: can the touched components be tested in isolation?

The lightweight code surface check looks only for:

- interfaces or injectable dependencies for test doubles
- existing test framework and test patterns in the target area
- whether tests can target the code without restructuring it first

Do not analyze how to make the area testable. Record a short skip reason such as `declared untestable in AGENTS.md`, `no testable seams`, or `no existing test infrastructure`.

## Isolation Gate

Structural TDD requires physical isolation and prompt isolation.

- **Physical isolation**: the test author runs in a Worktrunk or plain `git worktree` workspace containing the relevant code and explicitly allowed build/cache artifacts, but not plan files, review files, anchors, or `progress.md`.
- **Prompt isolation**: the test author's prompt contains no design rationale, plan path, task path, feature name, or breadcrumb that could lead to the plan.

Use the active runtime adapter to verify that the test author can be dispatched inside the isolated workspace. Prompt hygiene alone is not enough.

If no enforceable isolation path exists, skip structural TDD and record the reason. If an isolation path exists but dispatch verification fails, record the attempted mechanism and stop to ask whether to fix isolation or explicitly skip structural TDD.

## Build And Compilation Readiness

Before spawning the test author:

1. Verify that the target package compiles at the current execution state.
2. Verify that required test infrastructure, fixtures, helpers, and generated mocks exist and are current.
3. For build-heavy projects, configure or seed build/cache reuse only according to [workspace-isolation.md](workspace-isolation.md).

If compilation or test infrastructure is blocked, resolve the blocker within the TDD framework instead of skipping TDD or inverting the order.

Acceptable scaffolding includes method stubs that panic or return zero values, adding a package to mock-generator configuration and regenerating mocks, or creating missing test helpers. These are temporary unblockers for the test author, not task implementation. Record every scaffold in progress so the implementer knows what must be replaced.

## Test Author Context

The executor passes the test author only:

- acceptance criteria, extracted and relayed as inline text
- the relevant code surface available in the isolated workspace
- the `test-driven-development` skill plus required project testing skills through the active runtime binding when available

The test author must not receive:

- design decisions or implementation strategy
- plan file paths, plan directory paths, task file paths, feature names, or anchor paths
- the complete task file
- prerequisite task rationale beyond code already present in the workspace

The test author writes tests grounded in the acceptance criteria, follows the `test-driven-development` skill, confirms RED according to that skill, and returns test file paths plus a summary of what each test verifies.

## TDD Quality Gate

Use the `test-driven-development` skill as the source of truth for RED validation and test-quality review.

Before handing tests to the implementer, verify the test author's reported RED result against that skill. If the tests fail the review, add only acceptable scaffolding from this reference, rerun RED, or send the tests back to the test author.

After the implementer reports green and before marking a task ready for integration or done, apply the same test-quality review to the test-author tests. If the tests fail the review, block the task as `blocked: weak test` and rerun test authoring.

Record the review result in progress, including the sampled files and decision.

## Returning Test Files

If the test author ran in the task's implementer worktree, leave the tests there for the implementer.

If the test author ran in a temporary isolated workspace, bring only the test files back to the main execution workspace using the active runtime adapter's mechanism, then remove the temporary workspace. Verify the files exist in the execution workspace before updating progress.

## Implementer Contract With Structural TDD

The implementer receives:

- the complete task as an inline packet
- test file paths from the test author
- prerequisite outputs relayed by the executor
- required skills through the implementer's execution binding

Tell the implementer:

- test-author tests are immutable
- make the tests pass by implementing correctly
- after tests pass, refactor implementation quality while keeping tests green
- if a test cannot be satisfied, report a dispute with file, test name, and explanation
- do not modify, skip, or delete test-author tests

The implementer returns implementation status, files created or modified, test results, and disputes.
