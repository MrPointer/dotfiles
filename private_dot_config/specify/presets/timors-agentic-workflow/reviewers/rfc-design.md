# RFC Design Review Packet

- **Tier**: most-capable
- **Required skills**: exact project RFC and architecture-review skill identifiers
- **Inputs**: current `spec.md`, `plan.md`, relevant grounding/research, prior role report
- **Workspace**: correctly targeted read-only feature/repository view
- **Exclusions**: no mutation, no binding choice, no human acceptance decision

Independently assess whether the current RFC is technically coherent and complete:
verified current state, architecture, boundaries, contracts, flows, state,
failure behavior, compatibility/migration, risks, verification, and Feature
Definition alignment. Check that durable decisions are in `plan.md`, not only
supporting workspaces. Return an attributable complete current-state round with
reviewed paths, trigger, scope, semantic blocking/concern findings, and verdict.
A passing verdict never implies Design Acceptance.
