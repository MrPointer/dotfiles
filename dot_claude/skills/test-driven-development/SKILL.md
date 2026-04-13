---
name: test-driven-development
description: Use when implementing features or fixes. Write tests first from acceptance criteria, confirm they fail, then implement to make them pass. Prevents the common agent failure mode of writing tests that rationalize implementation rather than verify requirements.
---

# Test-Driven Development

Write tests before implementation. Tests encode what the code **should** do, derived from requirements — not what the code **happens** to do after the fact.

This matters more for agents than for humans. When an agent writes code first and tests second, it unconsciously writes tests that validate what it built rather than what was required. Inverting the order forces tests to be grounded in requirements, making them an honest verification layer.

## Core Principles

1. **Tests Come From Requirements, Not Implementation**: Every test must trace back to a stated requirement, acceptance criterion, or behavioral expectation. If you can't point to where a test's expectation comes from, it's not grounded.
2. **RED Before GREEN**: Run tests after writing them. They must fail. A test that passes immediately either tests nothing or was written with knowledge of the implementation — both are failures.
3. **Never Adjust Tests to Match Broken Code**: If a test fails after implementation, the implementation is wrong — not the test. The only reason to change a test is if the *requirement* it encodes was wrong.

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
- Test observable behavior (outputs, side effects, state changes) — not internal implementation details (private methods, internal data structures, call counts).
- Mock at external boundaries — hardware, network, system time, third-party services, databases. These should always be mocked. But don't mock internal components of the system under test — use real implementations and mock *their* external dependencies instead.
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

## Standalone Workflow

When using this skill directly (not as a test author within an executor), follow the full RED-GREEN cycle:

1. **Extract testable claims** from requirements (see above)
2. **Write tests** following the grounding rules (RED)
3. **Confirm failure** — run tests, every new test must fail. If a test passes immediately, it either tests something that already works or tests nothing meaningful.
4. **Implement** — write the implementation to make tests pass. If a test turns out to be wrong, fix the test *before* adjusting the implementation, and document why.
5. **Confirm pass** — run tests again. All new tests should pass, existing tests should still pass. If something fails, diagnose and fix the implementation directly. If you find yourself going through more than two fix attempts, stop and reassess.
6. **Refactor** — now that the tests are green, improve the implementation: clean up structure, naming, duplication, abstractions. The tests are your safety net — refactor freely as long as they keep passing. Run tests after refactoring to confirm nothing broke.

## Rules

- **Never write implementation before tests** — if you catch yourself doing it, stop, delete the implementation, write the test
- **Never adjust test expectations to match a failing implementation** — fix the code, not the test
- **Never skip the RED step** — a test that wasn't seen to fail is not trusted
- **Follow project conventions** — this skill defines the order of operations, not how to write tests. The project's testing skills and patterns govern test structure, assertions, frameworks, and file organization
