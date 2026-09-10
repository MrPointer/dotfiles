---
description: Create or reconcile a Feature Definition without interactive questions.
---

## Authority and Inputs

Read supplied intent, canonical feature `spec.md` when present,
`.specify/presets/timors-agentic-workflow/references/artifact-contracts.md`, `.specify/presets/timors-agentic-workflow/references/feature-definition.md`, and
`.specify/presets/timors-agentic-workflow/references/review-lifecycle.md`. Render `.specify/presets/timors-agentic-workflow/templates/spec-template.md`. `spec.md`
is the product-intent authority. Do not inspect the repository, ask a question,
choose technical architecture, or supply a default.

## Procedure

1. Resolve the canonical feature directory. When the caller supplied a feature
   directory, use it. Otherwise derive it under `specs/`: read `feature_numbering`
   from `.specify/init-options.json`, prefix the name with the current
   `YYYYMMDD-HHMMSS` when the value is `timestamp` or with the next zero-padded
   three-digit number after the existing `specs/` entries otherwise, and name the
   directory `<prefix>-<short-name>` from a slug of the supplied intent. Do not
   write anything yet.
2. Construct or reconcile the complete Feature Definition from positive supplied
   information only. Preserve semantically titled requirements and their inline
   acceptance evidence. Include conditional actors, journeys, failures, risks,
   accessibility, localization, compatibility, or operations only when relevant.
3. Scan every Ready criterion in `.specify/presets/timors-agentic-workflow/references/feature-definition.md`. When any
   criterion is not positively established, record every material unanswered
   question in `Material Open Questions`, set `status: Draft`, and do not infer a
   resolution. Otherwise set `status: Ready` and state why it satisfies the
   predicate.
4. Before a material write, check whether local `progress.md` exists with a locked
   in-progress, blocked, or complete package. If so, refuse the change. Before a
   permitted material write, make the ordered invalidations in
   `.specify/presets/timors-agentic-workflow/references/review-lifecycle.md`; do not erase downstream artifacts or prior
   review history. Create the feature directory when absent. Seed `spec.md` from
   `.specify/presets/timors-agentic-workflow/templates/spec-template.md` when it does not already exist. Write the
   Feature Definition to `spec.md`. Persist the resolved directory to
   `.specify/feature.json` as `{"feature_directory": "<path>"}`.
5. Report the canonical path, status, material questions, invalidated records, and
   only the next permitted phase: Draft hands off to `clarify`; Ready may hand off
   to `plan`.

Missing, partial, or contradictory state blocks its consumer when that consumer
reads it. Do not add a package, manifest, installed-inventory, mapping, or version
preflight.
