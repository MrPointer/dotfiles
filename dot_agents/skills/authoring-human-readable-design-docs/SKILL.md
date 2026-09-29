---
name: authoring-human-readable-design-docs
description: Drive authoring of a human-audience design document — a narrative that explains a settled design for a person to understand, not for an agent to execute — one section at a time with user review between chunks, delegating the writing to the mid-tier-technical-writer agent. Also drives standalone editing passes on such a doc (clarify, simplify, restructure, re-term). Use only when the user explicitly asks to write, draft, or edit a human-readable design doc, or to turn a brainstorm, decision log, or design notes into one. This is an opt-in, manually invoked workflow — do not trigger it automatically for general writing, editing, or documentation tasks. Not for planning or agent-targeted RFCs, execution plans, or reference documentation (domain, architecture, component, business-process); those have their own paths.
---

# Authoring Human-Readable Design Docs

Drive the creation of a design document written for a person to understand a
design — the kind of narrative you hand a colleague so they grasp how something
works and why, not a spec an agent executes. The design is already settled; this
workflow turns it into readable prose and keeps it readable over later editing
passes.

Work incrementally. Author one section at a time and let the user review each
chunk before moving on. One-shot drafts of a whole doc cost more, bury mistakes,
and review worse — small chunks keep the user in control and catch drift early.

## The agent that does the writing

The `mid-tier-technical-writer` agent writes and edits the prose. The voice,
tone, and wording rules live in that agent, not here — this skill orchestrates
and defers every writing decision to it. Do not write the document prose
yourself in the coordinating session, and do not restate the agent's writing
rules here; that only creates a second, drifting copy.

Give the agent the source material, the target section, and the output file
path, and let it write the section straight into the file. Do not copy the
agent's prose back into your own response — the section already lives in the
file, so repeating it wastes tokens and adds nothing. Point the user at the file
and the section to review instead.

## The one boundary: wording versus design

Each round of user feedback is either about the writing or about the design.
Classify it before you act:

- **Wording, structure, or clarity** — rephrasing, reordering, splitting or
  merging sections, simplifying a dense passage, fixing a term. This is in
  scope. Route it to the writer agent.
- **The design itself** — "actually, let's change how X works," a new
  constraint, a reversed decision, a different approach. This is out of scope.

When feedback changes the design, stop. Do not hand it to the writer agent and
do not absorb it into the prose. The writer agent must not invent or alter
design conclusions; if it did, an unreviewed design change would end up buried
in a document instead of being decided deliberately. Tell the user plainly that
this is a design change, and suggest they take it back to the drawing board in a
separate design session, then return to authoring once the design is settled
again.

## Workflow

1. **Setup** — Confirm the section structure (take the user's, or propose one
   and get agreement) and the output file. Write to a new file and leave any
   original source untouched, so the source stays available for reference.
   Identify the source material the doc is built from.
2. **Author section by section** — Have the writer agent write one (sub)section
   at a time directly into the output file, then let the user review that
   section in the file before the next. Report only that the section landed and
   where — do not relay the written prose in your response. Defer terms or links
   that point at sections not yet written; add the backlinks once the target
   section exists, so the draft never references something that isn't there.
3. **Editing passes** — For polish on an existing doc, use the same agent for
   standalone passes: clarify wording, simplify dense or math-heavy passages,
   restructure, fix terminology, remove repetition.
4. **Wrap-up** — Check that cross-references and anchors resolve. If a pass
   surfaced a fix that also belongs in the source design material, note it for
   the user rather than silently changing only the doc.

## Coordinator behaviors

- When the user is confused by a passage, explain the concept in plain terms
  first, then let that explanation drive the rewrite. Often the plain
  explanation you just gave is the model for the fixed prose.
- Prefer incremental progress over large one-shot sections — cheaper and easier
  to review.
- User decisions override. Do not silently re-decide something the user already
  settled.
