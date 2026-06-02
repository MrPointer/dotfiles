---
name: authoring-rfcs
description: Use when turning a settled design direction, brainstorming outcome, context anchor, or architecture discussion into a codebase-grounded RFC. Produces self-contained RFC documents with verified current architecture, chosen design, tradeoffs, risks, stable references, and planning handoff context without creating implementation plans.
---

# RFC Authoring

Turn a validated direction into an engineering RFC that can survive review and feed planning without re-litigating architecture.

The output is a human engineering design document: self-contained, precise, grounded in the current codebase or explicit assumptions, and clear about the chosen design. It is not an anchor, implementation plan, ADR, transcript, agent handoff, or arranged notes.

This file is the RFC-authoring spine. Load referenced files only when their trigger applies.

## Relationship To Other Skills

| Skill | Owns | RFC Authoring Role |
|-------|------|--------------------|
| `brainstorming` | Exploration, framing, alternatives, user decisions, convergence | Convert the settled direction into the final RFC |
| `anchoring-context` | Durable working memory for active work | Use the anchor as input, then let the RFC become the settled design artifact |
| `applying-architecture-patterns` | Clean Architecture, Hexagonal Architecture, and DDD guidance | Load when backend architecture patterns are part of the design |
| `planning-project-features` | Decomposing an approved design into executable plans | Produce and reviewer-validate the RFC baseline planning will use |

If used after `brainstorming`, do not reopen the conversation unless reviewer feedback reveals a blocking gap. If used without `brainstorming`, confirm that the chosen direction is already settled before writing the RFC.

## Core Standard

- **Reality before proposal**: verify current architecture before proposing changes.
- **Observed vs proposed**: keep current-state facts separate from proposed design.
- **One chosen design**: commit to a direction; put alternatives in the alternatives section.
- **Architectural precision**: name actual components, boundaries, contracts, data shapes, state owners, and runtime flows when they affect the design.
- **No execution plan**: include enough design detail for planning, but no task lists, sequencing, tickets, commands, or step-by-step implementation instructions.
- **No fabrication**: mark unverified claims as assumptions or open questions.
- **Reviewer-owned quality**: architecture, risk, and clarity validation belong to reviewer subagents.
- **Pattern-aware, not pattern-driven**: use architecture patterns only when they fit the problem and current system.
- **Human RFC voice**: do not expose prompting, agent workflow, or reviewer mechanics in design prose.

## Conditional References

| Trigger | Load |
|---------|------|
| Gathering context, verifying current reality, normalizing decisions, or deciding whether architecture-pattern guidance applies | [references/input-and-reality.md](references/input-and-reality.md) |
| Drafting the RFC, choosing location/name, applying prose rules, or deciding which sections belong | [references/writing-the-rfc.md](references/writing-the-rfc.md) |
| Running reviewer subagents, handling review artifacts, classifying changes, or updating the Review Record | [references/reviewer-review.md](references/reviewer-review.md) |
| Writing the RFC body | [assets/rfc-template.md](assets/rfc-template.md) |

## Workflow

### 1. Gather Inputs

Read user-approved direction, active anchors, relevant design docs, ADRs, domain/architecture/process docs, existing RFCs, and source files that define current architecture or contracts.

Use [input-and-reality.md](references/input-and-reality.md) for companion-skill triggers, greenfield handling, and source-gathering expectations.

### 2. Verify Current Reality

Build a concise current-state model before writing the proposed design. Important `currently` claims in the RFC must be backed by source references, documentation references, or explicit user statements.

Use [input-and-reality.md](references/input-and-reality.md) for what to verify and how to separate verified facts from assumptions.

### 3. Normalize Decisions

Reduce brainstorming or discussion output into settled design inputs: chosen approach, goals, non-goals, constraints, rejected alternatives, and genuinely unresolved open questions.

If an unresolved question changes the architecture, stop and ask the user one question. Non-blocking questions can remain in the RFC as open questions.

### 4. Write The RFC

Use [assets/rfc-template.md](assets/rfc-template.md) unless the project has a stronger convention.

Write the RFC as a normal engineering document, not a transcript or agent handoff. Save it using the project's RFC convention, or `docs/rfcs/<topic>.md` when no convention exists.

Use [writing-the-rfc.md](references/writing-the-rfc.md) for human-voice rules, filename/ID rules, required conceptual sections, source references, and placeholder removal.

### 5. Run Reviewer Review

Before presenting the first complete RFC draft, run reviewer subagents once. Do not self-approve the RFC's architecture, risk profile, or clarity.

Required reviewers:

- `rfc-architect-reviewer`
- `rfc-risk-reviewer`

Use `rfc-clarity-reviewer` when the RFC is long, nuanced, assumption-heavy, produced from a long brainstorming session, or requested by the user.

Use [reviewer-review.md](references/reviewer-review.md) for reviewer prompts, output paths, cumulative artifacts, finding incorporation, change classification, targeted re-review, and Review Record updates.

### 6. Present And Stop

Present the RFC path and a brief summary of the chosen architecture. Stop unless the user explicitly asks for revisions or asks to move into planning.

If an anchor was used, update or reconcile it so settled design context points to the RFC instead of duplicating it.

## Final Rules

- Never substitute a context anchor for an RFC.
- Never let the RFC become a plan.
- Never self-approve architecture, risk profile, or clarity when reviewer agents are available.
- Never hide blocking uncertainty in vague language.
- Never claim existing architecture, contracts, storage formats, or runtime behavior without verifying them.
- Preserve approved decisions. If verified reality conflicts with one, present the conflict and ask before changing the design.
- Never refer to the human requestor as `the user` in RFC prose; recast input as requirements, constraints, decisions, assumptions, or source references.
- Never mention agents, models, prompts, subagents, reviewer mechanics, or the authoring workflow in the RFC body except for required Review Record metadata.
- Pseudo-code is allowed only when it clarifies an algorithm, contract, state transition, or data/control flow. Keep it illustrative and avoid low-level implementation recipes.
