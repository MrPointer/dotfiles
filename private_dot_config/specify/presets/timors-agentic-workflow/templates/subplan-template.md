---
kind: subplan
primary_files:
  - project/root/relative/path
execution:
  model_tier: mid-tier
  rationale: "Substantive provider-neutral tier rationale."
  skills:
    - exact-skill-id
test:
  mode: required
  basis: "Planned test scope and exact test location, suite, or command."
contracts:
  produces: []
  consumes: []
---

# [Packet Objective]

## Objective

[Bounded outcome owned by this packet.]

## Baseline Context and Constraints

[Accepted RFC context, constraints, and protected feature-directory boundary.]

## Scope and Primary Files

[Exact owned project-root-relative paths and exclusions.]

## Prerequisites and Integration

### Direct Predecessors

| Path | Classification | Reason |
|---|---|---|

<!-- Add one exact row for each `after` entry; no rows when `after` is empty. -->

### Transitive Context

[Non-authoritative inherited background, or `None`.]

## Runtime Considerations

[Anticipated conditions only; they are not access grants, eligibility rules, or gates.]

## Acceptance and Verification

[Exact commands/checks, expected outcomes, and evidence for frozen test mode/basis.]

## Required Outputs and Handoff

When produces is empty, add no produced-contract section. For every produced name,
add exactly one direct level-three heading in this form:

```markdown
### Produced Contract: `semantic-contract-name`
```

Under each such heading add exactly once and in order: `#### Description`,
`#### Contract`, `#### Transport`, `#### Invariants and Constraints`, and
`#### Required Output Evidence`. No extra direct level-four peer is allowed. The
last subsection names evidence returned by worker and recorded against final
integrated producer commit.
