---
name: test-driven-development
description: Use whenever implementing or fixing production code, writing or changing tests, reviewing test quality, applying TDD, or diagnosing tests that pass without proving behavior. Ensures tests are requirement-grounded and reach behavioral RED before implementation.
---

# Test-Driven Development

Write tests before implementation. Tests encode what the code **should** do, derived from requirements — not what the code **happens** to do after the fact.

This matters more for agents than for humans. When an agent writes code first and tests second, it unconsciously writes tests that validate what it built rather than what was required. Inverting the order forces tests to be grounded in requirements, making them an honest verification layer.

This file is the TDD spine. Load referenced files only when their trigger applies.

## Conditional References

| Trigger | Load |
|---------|------|
| Writing or changing mocks, fakes, spies, captured observations, test helpers, or production code mainly for tests | [references/testing-anti-patterns.md](references/testing-anti-patterns.md) |

## Core Principles

1. **Tests Come From Requirements, Not Implementation**: Every test must trace back to a stated requirement, acceptance criterion, or behavioral expectation. If you can't point to where a test's expectation comes from, it's not grounded.
2. **RED Before GREEN**: Run tests after writing them. They must fail for a behavioral reason. A test that passes immediately either tests nothing or was written with knowledge of the implementation — both are failures.
3. **Never Adjust Tests to Match Broken Code**: If a test fails after implementation, the implementation is wrong — not the test. The only reason to change a test is if the *requirement* it encodes was wrong.
4. **Compile Failure Is Not Enough**: A test that only fails because the implementation is missing, a symbol is undefined, or the package does not compile has not proven it verifies behavior. Add the minimal scaffold needed to reach an assertion failure before treating RED as complete.

## Extracting Testable Claims

Before writing any test, identify what needs to be true. Look for:

- **Acceptance criteria** — explicit success conditions (often found under headings like "Acceptance Criteria", "Success Criteria", "Definition of Done")
- **Behavioral expectations** — what the system should do in response to inputs
- **Constraints** — boundaries the implementation must respect (performance, compatibility, error handling)
- **Edge cases** — explicitly stated boundary conditions

If the input has a section with acceptance criteria, use those directly. If requirements are informal (conversation, rough description), extract testable claims yourself and confirm them with the user before writing tests.

**Each testable claim becomes one or more test cases.** Name tests after the behavior they verify, not the implementation they exercise.

## Grounding Rules

- Every assertion must trace to a specific requirement or acceptance criterion. If you can't name which requirement a test verifies, delete the test.
- Every test must contain at least one assertion that would fail against a plausible wrong implementation, not just against a missing implementation.
- Test observable behavior (outputs, side effects, state changes) — not internal implementation details (private methods, internal data structures, call counts).
- Mock at external boundaries — hardware, network, system time, third-party services, databases. These should always be mocked. But don't mock internal components of the system under test — use real implementations and mock *their* external dependencies instead.
- If a boundary mock, spy, or fake captures observations such as arguments, call records, emitted values, or sequences, assert on the captured values required by the acceptance criteria. Capturing without asserting is a no-op test.
- Don't write tests for hypothetical scenarios the requirements don't mention. Test what's specified, not what might theoretically go wrong.

## Anti-Patterns

| Anti-Pattern | What It Looks Like | Why It's Wrong |
|---|---|---|
| **Rationalized tests** | Tests written after implementation that mirror what the code does | Tests validate what was built, not what was required — bugs are invisible |
| **Assertion adjustment** | Changing expected values to match actual output after a test fails | The test was right, the implementation was wrong — you just hid the bug |
| **Mock internals** | Mocking internal components instead of their external dependencies | Tests pass in isolation but the integrated system fails — mock at the boundary, not in the middle |
| **Tautological tests** | Testing that a function returns what it returns | Proves nothing — the test will pass regardless of correctness |
| **Implementation-coupled tests** | Tests that break when you refactor internals without changing behavior | Tests should verify behavior, not implementation structure |
| **Spec-less tests** | Tests for scenarios nobody asked for, "just to be safe" | Wastes tokens, adds maintenance burden, and may encode wrong assumptions |
| **Compile-only RED** | New tests fail only because code is missing or doesn't compile | The test may pass once any implementation exists, regardless of correctness |
| **Unasserted capture** | A test records arguments, calls, or emitted values but never checks them | Observation infrastructure without assertions does not verify behavior |

## Test Quality Review

Before treating tests as a valid safety net, review them against these checks:

- The RED failure is an assertion or expectation failure tied to the required behavior, not only a compilation or setup failure.
- Each test has at least one requirement-backed assertion that would fail against a plausible wrong implementation.
- Captured or recorded observations are asserted for the values, order, or side effects required by the acceptance criteria.
- The test would fail if the implementation returned hardcoded success while skipping the required behavior.

If any check fails, fix the test or add minimal scaffold and rerun RED before implementation continues.

## Standalone Workflow

When using this skill directly (not as a test author within an executor), follow the full RED-GREEN cycle:

1. **Extract testable claims** from requirements (see above)
2. **Write tests** following the grounding rules (RED)
3. **Confirm failure** — run tests, every new test must fail for a behavioral assertion or expectation reason. If a test passes immediately, it either tests something that already works or tests nothing meaningful. If it fails only to compile, add minimal scaffold or fix the test setup and rerun until the failure proves the behavior is not implemented.
4. **Implement** — write the implementation to make tests pass. If a test turns out to be wrong, fix the test *before* adjusting the implementation, and document why.
5. **Confirm pass** — run tests again. All new tests should pass, existing tests should still pass. If something fails, diagnose and fix the implementation directly. If you find yourself going through more than two fix attempts, stop and reassess.
6. **Refactor** — now that the tests are green, improve the implementation: clean up structure, naming, duplication, abstractions. The tests are your safety net — refactor freely as long as they keep passing. Run tests after refactoring to confirm nothing broke.

## Rules

- **Never write implementation before tests** — if you catch yourself doing it, stop, delete the implementation, write the test
- **Never adjust test expectations to match a failing implementation** — fix the code, not the test
- **Never skip the RED step** — a test that wasn't seen to fail is not trusted
- **Never accept compile-only RED** — missing implementation errors are setup, not behavioral proof
- **Never leave captured observations unasserted** — if the test records values, it must check the values that matter
- **Follow project conventions** — this skill defines the order of operations, not how to write tests. The project's testing skills and patterns govern test structure, assertions, frameworks, and file organization
