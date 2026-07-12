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
scripts:
  sh: scripts/bash/setup-plan.sh --json
  ps: scripts/powershell/setup-plan.ps1 -Json
---

## Preset Compatibility Preflight

Before extension hooks, prerequisite scripts, or any write, read
`.specify/presets/timors-agentic-workflow/preset.yml` and obtain the active
Spec Kit version with `specify --version`. Parse the manifest's
`requires.speckit_version` range and fail closed unless the active version is
within `>=0.12.11,<0.13.0`. If the manifest, version, or version range cannot
be read and parsed, stop and report the compatibility failure. Do not invoke a
hook, setup script, or write a project artifact after a failed preflight.

## Preset Planning Composition

Use `.specify/presets/timors-agentic-workflow/references/planning-grounding.md`
as the policy for the four added plan sections and feature reuse. Keep the
upstream planning workflow and artifact ownership intact: decisions belong in
`research.md`; data details belong in `data-model.md` and `contracts/`; and
execution decomposition belongs in `execution-plan.md`, not `plan.md`.

{CORE_TEMPLATE}
