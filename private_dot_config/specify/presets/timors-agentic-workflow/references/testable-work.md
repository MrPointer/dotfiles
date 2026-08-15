# Test Modes and Evidence

Every packet has one `test.mode` and a nonempty substantive `basis`; Acceptance
and Verification supplies exact commands/checks and expected outcomes without
changing mode.

| Mode | Required basis | Worker evidence |
|---|---|---|
| `required` | planned test scope and exact test location, suite, or command | implement changed behavior and tests together; run and report coverage and results |
| `existing-coverage` | exact existing test files, cases, suites, or commands and the acceptance behavior they prove | run cited evidence; if it does not prove changed behavior, add/update owned tests and report the false assumption and correction |
| `not-applicable` | substantial out-of-scope infrastructure or restructuring barrier, why it is substantial, and alternative verification | perform/report alternative verification and preserve named automated-testing gap |
| `no-testable-behavior` | why no behavior is meaningfully automatable and applicable non-test verification | perform/report named non-test verification; do not supply synthetic test evidence |

Required and existing-coverage packets reserve all test paths needed for planned
coverage or bounded correction. Existing coverage cannot cite generic suite
existence. Preference or inconvenience never justifies not-applicable. If required
or corrected existing coverage cannot produce meaningful coverage within approved
ownership and scope, block the packet and dependents; frozen mode and basis do not
change. Continuing requires explicit abandonment, task replanning, applicable
review, and renewed authorization. Not-applicable means a known gap, not absence of
testable behavior; only no-testable-behavior makes that assertion.

Progress records each packet’s planned mode/basis, commands/checks, outcomes, and
added or changed tests. It retains a not-applicable barrier and alternative through
completion and records incorrect existing-coverage assumption. Final Verification
and user completion reporting must agree with this evidence and runtime outcomes.
An approved testing gap is reportable and adds no approval gate, but insufficient
resulting verification blocks completion.
