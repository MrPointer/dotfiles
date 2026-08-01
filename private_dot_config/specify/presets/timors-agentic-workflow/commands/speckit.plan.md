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

{CORE_TEMPLATE}
