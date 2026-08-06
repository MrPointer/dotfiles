---
description: Create or update the feature specification from a natural language feature description.
---

## Package Integrity Gate

This is the first operational action. Before extension hooks, prerequisite
scripts, artifact writes, progress or workspace mutation, or dispatch, read and
complete the shared gate in
`.specify/presets/timors-agentic-workflow/references/protocol-compatibility.md`.
Do not duplicate, weaken, or partially apply its minimum-version or complete
package-integrity checks. On any failure, report the failed required contract
and stop before continuing.

## Completion Report And Handoff Deferral

Execute the complete upstream workflow, including artifact generation,
validation, and mandatory post-execution hooks. Do not emit the upstream final
Completion Report or perform its handoffs until the optional human-review step
below has completed or is declined. Defer only final reporting and handoffs.

{CORE_TEMPLATE}

## Optional Specification Human Review

After the upstream specification workflow has generated and validated its
artifacts, apply
`.specify/presets/timors-agentic-workflow/references/human-review.md` at this
natural completion point before the final completion reporting and handoffs.
Preserve all upstream specification behavior; this wrapper adds only the
optional review offer.

For a guided walkthrough, use judgment to select consequential conceptual
units. Consider, when consequential: problem/outcome/scope;
scenarios/priorities; acceptance behavior/edge cases; requirements;
assumptions/dependencies/clarifications; and measurable success criteria. Keep
the interaction mechanics in the shared protocol rather than turning this into
a literal artifact or section walkthrough.
