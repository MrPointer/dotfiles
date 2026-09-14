# Binding, Runtime Conditions, and Result Attribution

A tracked reviewer names a provider-neutral tier and a role description. A tracked
implementation packet names a provider-neutral tier and the skills and
capabilities the group needs. Eligibility requires the requested tier, native
discoverability, invokability, correct workspace targeting, and attributable
results; a reviewer additionally must not modify the artifacts it reviews and
writes only its review report, outside the reviewed set. Prefer a
project-local candidate. For a reviewer, match the role description against the
candidate descriptions and select a single clear fit. For an implementation
packet, match the named skills and capabilities and select a single eligible
candidate. A unique fit, a sole eligible candidate, or an explicit human choice
binds; zero clear or several plausible candidates are collected blockers. Abstract
capabilities and permission breadth neither grant nor gate eligibility. Do not
provision a worker or perform a suitability probe; the first invocation is real
work. Record the concrete binding in runtime evidence, never in a tracked artifact.

Runtime Considerations are anticipated network, local tool, nested-dispatch,
cache, or similar conditions. They are planning context only: neither access grants
nor worker requirements or eligibility gates. During real work, attempt a
requirement-preserving in-scope alternative if a condition is unavailable. Return
the condition, whether/how attempted, concrete permission or environment failure,
alternative, and residual verification limitation. Block only when no alternative
can meet packet requirements or acceptance; never retroactively disqualify binding.

The coordinator’s result envelope has only: status; changed paths; commits; tests
and verification performed; produced contract-output evidence; runtime-condition
outcomes; and blockers. It cannot add or alter scope, outputs, acceptance,
verification, handoff, or semantic contract content. On conflict the selected
sub-plan prevails. An applicable empty runtime-outcome set is explicit. Reviewer
results additionally prove role, reviewed paths, round, verdict, that no reviewed
artifact was modified, and attribution; a modification of any reviewed artifact
rejects the result.
