# Optional Human Review Protocol

## Authority And Scope

This shared protocol governs the optional human-review offer at the natural
completion of `/speckit.specify` and `/speckit.plan`, after that phase has
generated and validated its artifacts. It supplies interaction mechanics only.
Each command supplies its own phase-specific, conceptual review guidance.

Do not apply this protocol to `/speckit.analyze`, task generation, or
implementation. It does not replace their existing review, approval, or
execution contracts.

Human review is optional for each phase independently. Offer a guided
walkthrough or independent review, and allow the user to decline both and
continue normal completion. There is no persisted preference, no approval
record, and no workflow gate. Do not advertise later re-entry into a completed
phase or relitigation of an earlier phase.

## Offer And Independent Review

After the owning command's normal artifact generation and validation complete,
briefly offer these choices without a roadmap:

1. **Guided walkthrough** — discuss consequential concerns selected with the
   command's phase-specific guidance.
2. **Independent review** — identify the phase's key artifacts for the user to
   read independently, encourage that final read, then continue normal
   completion reporting and handoffs.

Independent review does not start a guided walkthrough, wait for an approval,
or change the command's existing report or handoff behavior.

## Guided Walkthrough Mechanics

Select only consequential decisions owned or introduced by the current phase,
using judgment, the phase guidance, and the conversation context. Do not provide
an upfront walkthrough roadmap. Inputs and settled decisions from earlier phases
are authoritative baselines for the current phase. Cite them briefly as
constraints, rationale, or traceability context when useful, but they must not
become review units requiring renewed confirmation.

For qualifying current-phase units, review material requirements and boundaries,
model-selected decisions, meaningful assumptions, tradeoffs or rejected
alternatives, unresolved ambiguity, compatibility or operational consequences,
and risks when they are consequential. Do not limit the walkthrough to
unresolved decisions.

Use supporting artifacts such as `data-model.md`, `research.md`, or contracts
as evidence sources, not mandatory review stops. Skip mechanical projections or
repetition unless they expose a consequential decision. Do not force a literal
section-by-section or file-by-file walkthrough.

Each review unit covers one coherent concern and stays cognitively small. Group
only tightly coupled details that require shared reasoning, never merely
because they seem noncontroversial. If the user could accept one item and
reject another, split them. Avoid dense turns and fatigue. For genuine
alternatives, lead with a recommendation and its reasoning; do not invent
options where no meaningful choice exists.

Pause after each review unit. Once the user settles it, immediately update every
affected artifact owned by the current phase, restore consistency and required
validation, briefly confirm the outcome, then select the next unit using
judgment and conversation context. During an active walkthrough, allow normal
conversational redirection, repetition, or stopping without ceremonial
navigation.

## Walkthrough Completion

When the walkthrough is complete, say so. Briefly summarize the resulting
decisions and changes, point only to the key phase artifacts for independent
self-review, encourage that final read, and continue the command's normal
completion reporting and handoffs. Do not imply that self-review is mandatory.
