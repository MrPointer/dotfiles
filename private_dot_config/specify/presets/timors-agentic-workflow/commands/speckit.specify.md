---
description: Create or reconcile a Feature Definition without interactive questions.
---

## Authority and Inputs

Read supplied intent, canonical feature `spec.md` when present,
`references/artifact-contracts.md`, `references/feature-definition.md`, and
`references/review-lifecycle.md`. Render `templates/spec-template.md`. `spec.md`
is the product-intent authority. Do not inspect the repository, ask a question,
choose technical architecture, or supply a default.

## Procedure

1. Construct or reconcile the complete Feature Definition from positive supplied
   information only. Preserve semantically titled requirements and their inline
   acceptance evidence. Include conditional actors, journeys, failures, risks,
   accessibility, localization, compatibility, or operations only when relevant.
2. Scan every Ready criterion in `references/feature-definition.md`. When any
   criterion is not positively established, record every material unanswered
   question in `Material Open Questions`, set `status: Draft`, and do not infer a
   resolution. Otherwise set `status: Ready` and state why it satisfies the
   predicate.
3. Before a material write, check whether local `progress.md` exists with a locked
   in-progress, blocked, or complete package. If so, refuse the change. Before a
   permitted material write, make the ordered invalidations in
   `references/review-lifecycle.md`; do not erase downstream artifacts or prior
   review history.
4. Report the canonical path, status, material questions, invalidated records, and
   only the next permitted phase: Draft hands off to `clarify`; Ready may hand off
   to `plan`.

Missing, partial, or contradictory state blocks its consumer when that consumer
reads it. Do not add a package, manifest, installed-inventory, mapping, or version
preflight.
