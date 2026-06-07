# Testing Anti-Patterns

Use this reference when writing or changing mocks, fakes, spies, captured call records, test helpers, or production code that exists mainly to support tests.

Tests must verify real behavior. Test doubles are tools for observation and isolation; they are not the behavior under test.

## Quick Gate

Before accepting a test that uses a test double, answer these questions:

1. What requirement-backed behavior does this test prove?
2. Is the double placed at an external boundary rather than inside the system under test?
3. Does the test depend on side effects from the thing being mocked?
4. If the test captures calls, arguments, emitted values, or sequences, does it assert the values required by the acceptance criteria?
5. Would a simpler integration-style test with real collaborators prove the behavior more clearly?

If any answer is unclear, stop and simplify the test before implementation continues.

## Anti-Patterns

| Anti-Pattern | What It Looks Like | Fix |
|--------------|--------------------|-----|
| Testing the double | Assertions only prove a mock, fake, or stub was installed | Assert the system's observable behavior with the double present |
| Unasserted capture | Calls, arguments, events, or sequences are recorded but never checked | Assert the captured values required by the acceptance criteria |
| Presence-only observation | The test only checks that something was called or that a calls list is non-empty | Assert meaningful arguments, order, count, result, or side effect |
| Over-mocking | Internal collaborators are mocked even though real code would be clearer | Mock only external boundaries or slow/unreliable dependencies |
| Mocking without understanding | A mocked method suppresses side effects the test actually needs | Understand the dependency chain, then mock at the lowest safe boundary |
| Incomplete fake data | Fake responses include only fields the immediate assertion reads | Mirror the real shape needed by downstream code and documented contracts |
| Test-only production API | Production code gains methods or flags used only by tests | Move setup/cleanup into test utilities or improve dependency injection |
| Mock complexity spiral | Mock setup dominates the test and obscures the behavior | Prefer real collaborators, narrower seams, or a higher-level test |

## Boundary Interaction Checks

Asserting boundary calls is valid only when the interaction is itself required behavior, such as retry attempts, emitted commands, audit events, or outbound requests.

When the interaction is required behavior, assert the contract:

- arguments or payload values
- call count when count is part of the requirement
- ordering when order is part of the requirement
- error handling or retry behavior when relevant

When the interaction is not required behavior, prefer asserting the user-visible output, persisted state, returned value, or externally observable side effect.

## Test-Only Production Code

Do not add production methods, flags, constructors, or state accessors solely because tests are hard to write.

Use one of these instead:

- test utilities for setup and cleanup
- dependency injection at real production seams
- existing public APIs
- narrower units with clearer boundaries

If a new production seam is genuinely useful outside tests, document the production use case in the implementation or design notes. Otherwise it belongs in test support code.

## Review Checklist

Before marking tests complete:

- The test fails for behavior, not just setup or compilation.
- The test would fail against hardcoded success that skips the required behavior.
- Every captured observation is asserted or removed.
- Mocks preserve side effects the test depends on.
- Fake data matches the real contract closely enough for this behavior.
- No production API was added only for test convenience.
