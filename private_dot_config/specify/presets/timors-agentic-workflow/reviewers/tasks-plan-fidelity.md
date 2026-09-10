# Task Plan-Fidelity Review Packet

- **Tier**: mid-tier
- **Role**: task plan-fidelity reviewer that traces the accepted plan into the task package
- **Inputs**: Ready `spec.md`, current accepted `plan.md`, `tasks.md`, indexed sub-plans, prior report
- **Workspace**: correctly targeted read-only feature/repository view
- **Exclusions**: no mutation, no implementation dispatch, no authorization decision

Independently assess execution package scope, constraints, risks, accepted design,
and producer-owned output definitions against Feature Definition and accepted plan.
Check that no packet introduces undefined behavior, scope, contract, or acceptance
exception. Return an attributable complete current-state round with reviewed paths,
trigger, scope, semantic blocking/concern findings, and verdict. Unapproved
deviation is blocking.
