# RFC Clarity Review Packet

- **Tier**: mid-tier
- **Role**: RFC clarity reviewer that judges whether a cold reader can follow the RFC
- **Inputs**: current `spec.md`, `plan.md`, prior role report
- **Workspace**: correctly targeted read-only feature/repository view
- **Exclusions**: no mutation, no design substitution, no human acceptance decision

Independently determine whether a cold reader can follow current-state evidence,
Feature Definition alignment, decision rationale, interfaces, constraints, risks,
verification, and planning handoff without conflicting or missing authority.
Identify ambiguity, unexplained consequence, or unreviewable handoff as semantic
findings. Return an attributable complete current-state round with reviewed paths,
trigger, scope, blocking/concern findings, and verdict.
