---
description: Execute the implementation planning workflow using the plan template to generate design artifacts.
handoffs:
  - label: Create Tasks
    agent: speckit.tasks
    prompt: Break the plan into tasks
    send: true
  - label: Create Checklist
    agent: speckit.checklist
    prompt: Create a checklist for the following domain...
---

## Package Integrity Gate

This is the first operational action. Before extension hooks, prerequisite
scripts, artifact writes, progress or workspace mutation, or dispatch, read and
complete the shared gate in
`.specify/presets/timors-agentic-workflow/references/protocol-compatibility.md`.
Do not duplicate, weaken, or partially apply its minimum-version or complete
package-integrity checks. On any failure, report the failed required contract
and stop before continuing.

## Preset Planning Composition

Use `.specify/presets/timors-agentic-workflow/references/planning-grounding.md`
as the policy for the four added plan sections and feature reuse. Keep the
upstream planning workflow and artifact ownership intact: decisions belong in
`research.md`; data details belong in `data-model.md` and `contracts/`; and
execution decomposition belongs in `execution-plan.md`, not `plan.md`.

## Completion Report And Handoff Deferral

Execute the complete upstream workflow, including artifact generation,
validation, and mandatory post-execution hooks. Do not emit the upstream final
Completion Report or perform its handoffs until the optional human-review step
below has completed or is declined. Defer only final reporting and handoffs.

{CORE_TEMPLATE}

## Optional Planning Human Review

After the upstream planning workflow has generated and validated its artifacts,
apply `.specify/presets/timors-agentic-workflow/references/human-review.md` at
this natural completion point before the final completion reporting and
handoffs. Preserve all upstream planning behavior; this wrapper adds only the
optional review offer.

Treat the current `spec.md` requirements, scenarios, boundaries, and success
criteria as authoritative inputs. Review only technical/design decisions
introduced or derived during planning. A plan choice that follows directly from
the specification and does not introduce a meaningful planning choice must not
be presented merely to reconfirm the specification. If the plan contradicts the
specification, correct the plan rather than reopen the specification. If
planning discovers a concrete contradiction, infeasibility, or material omission
that cannot be corrected within planning ownership, identify the blocker and
route it to the stock owning workflow; do not relitigate it inside planning
review.

For a guided walkthrough, use judgment to select consequential conceptual
units. Consider, when consequential: technical approach/research decisions;
current-state/architecture/project fit; interfaces/integration/data/contracts;
constraints/compatibility/constitution; testing/validation; and documentation
impact/unresolved design choices. Keep the interaction mechanics in the shared
protocol rather than turning this into a literal artifact or section
walkthrough.
