# Plan Design Review Packet

- **Tier**: most-capable
- **Role**: plan design reviewer that judges the technical coherence and architecture of the plan
- **Inputs**: current `spec.md`, `plan.md`, relevant grounding/research, prior role report
- **Workspace**: correctly targeted feature/repository view that the reviewer does not modify
- **Exclusions**: no modification of any reviewed artifact, no binding choice, no human acceptance decision
- **Output**: your review report only, written outside the reviewed set

Independently assess whether the current plan is technically coherent and complete:
verified current state, architecture, boundaries, contracts, flows, state,
failure behavior, compatibility/migration, risks, verification, and Feature
Definition alignment. Check that durable decisions are in `plan.md`, not only
supporting workspaces. Return an attributable complete current-state round with
reviewed paths, trigger, scope, semantic blocking/concern findings, and verdict.
A passing verdict never implies Design Acceptance.
