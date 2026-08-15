# Task Package Schema and Planning Rules

`tasks.md` is immutable after execution begins. It has exactly seven human
sections: Overview, Sub-Plan Ledger, Dependency Graph, Execution Inputs and
Requirements, Review Record, Implementation Authorization, Execution Handoff. It
has no atomic task lines or checkboxes. Mermaid and the ledger are non-authoritative
human projections; frontmatter is the sole machine authority.

Its complete closed frontmatter envelope is:

```yaml
kind: task-package
subplans:
  - path: subplans/NN-semantic-name.md
    after:
      - subplans/NN-predecessor.md
reviews:
  rfc_fidelity:
    report: reviews/tasks-rfc-fidelity.md
    round: null
    status: pending
    verdict: pending
  executability:
    report: reviews/tasks-executability.md
    round: null
    status: pending
    verdict: pending
authorization: pending
```

`subplans` is nonempty and topologically ordered. Each record has exactly `path`
and `after`; all nested records and the top level reject unknown keys. Every path
is the sole globally unique canonical identity, is feature-directory-relative,
stays beneath that directory, and resolves to exactly one tracked regular sub-plan
file. Every predecessor is a distinct identity in the package. Duplicate, missing,
escaping, ambiguous, untracked, non-regular, cyclic, or mismatched identities
block structure, review, authorization, dispatch, and resume. The ordered ledger
must match paths and exact direct-predecessor sets one-to-one, with no missing,
extra, duplicate, reordered, or mismatched rows.

Each selected sub-plan has exactly these closed fields:

```yaml
kind: subplan
primary_files: [project/root/relative/path]
execution:
  model_tier: cheapest | mid-tier | most-capable
  rationale: substantive provider-neutral string
  skills: [exact-skill-id]
test:
  mode: required | existing-coverage | not-applicable | no-testable-behavior
  basis: substantive string
contracts:
  produces: [semantic-kebab-case-name]
  consumes:
    - name: semantic-kebab-case-name
      producer: subplans/NN-producer.md
```

`primary_files` and skills are nonempty; every nested record is closed. There is
no sub-plan `path`, abstract capability, runtime identity, or producer-side
consumer list. Produces names match `^[a-z0-9]+(?:-[a-z0-9]+)*$`, are globally
unique, and are strings. A consume has exactly name and producer, each resolves to
one producer identity/name pair, and the producer must be a direct predecessor.

A packet body has exactly Objective; Baseline Context and Constraints; Scope and
Primary Files; Prerequisites and Integration; Runtime Considerations; Acceptance
and Verification; Required Outputs and Handoff. Direct Predecessors is a table
with exactly Path, Classification, Reason and exactly one row per `after` entry.
Each classification is `genuine implementation prerequisite` or `policy-only
serialization`; the latter does not loosen readiness. The empty `after` list has
no rows. Transitive Context is separate non-authoritative background. Concurrent
packets may not overlap primary files. Sequential overlap is valid only when an
ordered predecessor relation and both packet handoffs state the predecessor state.

For every produced name, Required Outputs and Handoff has exactly one direct child:
`### Produced Contract: ` followed by that exact single-backticked name. Its body
is closed at the next level-three heading or section end and has exactly these
direct level-four headings, once and in order: `#### Description`, `#### Contract`,
`#### Transport`, `#### Invariants and Constraints`, `#### Required Output
Evidence`. Extra level-four peers block; deeper headings are allowed within one
child. Description defines meaning; Contract defines required shape; Transport
states delivery; Invariants lists all constraints or exactly `None`; Required
Output Evidence names worker evidence recorded against final integrated producer
commit. This is the sole semantic definition. Produces and these bodies agree
one-to-one; empty produces has no such headings. Consumer prose may explain use
but cannot repeat, narrow, widen, or override producer definition.

Every consume must name its direct genuine predecessor, resolve its body exactly,
and receive final-commit evidence at dispatch. Every successor-needed output not
fully embodied and verifiable in integrated repository state requires this join.
Derive consumer coverage by scanning all consumes. Missing, duplicate, ambiguous,
misclassified, or unjoined output blocks; routing never creates an implicit edge.

Task review pointer round is `null` or `RNN`; status is `pending`, `reviewed`,
`retained`, or `recovery-required`; verdict is `pending`, `passed`,
`passed-with-concerns`, or `blocking`; authorization is `pending`, `approved`, or
`rejected`. Authorization may be approved only when both roles are reviewed or
retained, point to existing rounds, and pass or pass with concerns. Initial review
runs both. Cosmetic or scoped implementation-detail correction needs no rerun by
default. RFC-fidelity reruns when changed scope can affect Feature Definition or
RFC scope, constraints, risks, accepted design, or producer-owned definitions.
Executability reruns when ownership, dependencies, concurrency, routing, body
agreement, producer resolution, output evidence, test mode/basis, acceptance,
verification, or handoff changes. Packet addition/removal, split/merge, boundary,
DAG, or verification-scope change reruns both. Ambiguous impact is surfaced before
review work and large accumulated revision requires human scope confirmation.
