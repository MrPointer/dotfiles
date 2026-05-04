---
name: anchoring-context
description: Maintain concise feature-level working memory across sessions, teammates, and agents. Use when work spans multiple conversations; involves meaningful decisions, rejected alternatives, unresolved questions, or teammate/subagent handoffs; requires durable context before, during, or after brainstorming, planning, debugging, or implementation; or when the user asks to resume, continue, hand off, or preserve context for a feature. Skip for quick one-off questions, trivial edits, and work whose full context fits safely in a single short interaction.
---

# Context Anchoring

Maintain a concise, durable working-memory document for active work. The anchor preserves the feature's intent, decision rationale, rejected alternatives, constraints, open questions, current state, and next step so future sessions can resume without reconstructing context from chat history.

## Core Purpose

The anchor is the shared memory for a specific piece of work. It is not a transcript, scratchpad, implementation plan, or final specification.

Capture the information a future human or agent would otherwise have to rediscover:

- What problem this work is trying to solve.
- Why this direction matters now.
- Which decisions were made and what motivated them.
- Which alternatives were rejected and why.
- Which constraints must continue shaping future choices.
- Which questions remain unresolved.
- What the current state is and where to resume.

## Relationship To Other Skills

Context anchoring is a support protocol, not a replacement for task-specific skills.

| Skill | Owns | Anchor Role |
|-------|------|-------------|
| `brainstorming` | Exploration, framing, alternatives, convergence, final design spec | Read and update the anchor as persistent working memory during long-running design conversations |
| `planning-project-features` | Breaking an approved design into executable plans | Preserve decisions, unresolved questions, and handoff context that plans should respect |
| `executing-plans` | Task execution, acceptance criteria, test status, blockers, and mechanical progress checkpoints | Track only feature-level context that future sessions need: intent, decision rationale, constraints, deviations from the design, and handoff context |
| Documentation or ADR skills | Permanent project knowledge | Graduate stable decisions after the work settles |

When combined with `brainstorming`, brainstorming drives the conversation and decision process. This skill drives only anchor management. Do not let the anchor choose the design, and do not let brainstorming silently discard durable context.

When combined with `executing-plans`, the progress file remains the source of truth for task-by-task execution state. Do not duplicate task tables, test status, per-task blockers, or completion checklists in the anchor. Use the anchor only for feature-level context that explains why the work is shaped the way it is or what future sessions must not forget.

## Anchor Location

Follow project conventions if they exist. Otherwise use:

```text
docs/context/<topic>-anchor.md
```

Use a short, stable topic name. Prefer one anchor per coherent feature or investigation. If the anchor starts covering independent subsystems, stop and ask whether to split it.

## Anchor Structure

Use [anchor-template.md](assets/anchor-template.md) as the starting point when creating a new anchor. Adapt sections to the work, but keep the document concise.

For long-running anchors, include a **Quick Summary** near the top. The quick summary must be a self-contained design snapshot, not a set of category labels. Prefer a table with `Settled Direction` and `What It Means` columns. Each row should make sense without rereading the full decision log. Avoid vague areas like `Plugins`, `Scope`, or `Authoring` unless the row itself states the decision.

## Writing Rules

- Be concise. Preserve durable context, not conversation history.
- Record intent and motivation, not only raw decisions.
- Every meaningful decision should include a reason.
- Every rejected alternative should include why it was rejected.
- Prefer short rationale over long narrative.
- Record uncertainty explicitly as an open question.
- Keep transient speculation out unless it affects a future decision.
- Do not record private chain-of-thought. Summarize user-visible rationale and durable conclusions.
- Update the anchor when reality changes. Do not preserve stale decisions as if they still apply.
- Do not update inactive feature anchors just because later unrelated work changes code that the feature once touched. Anchors are for active work continuity, not historical maintenance.
- Prefer ADRs or permanent documentation for durable architectural decisions that must remain true after the feature is complete.
- Use markdown checkboxes in **Current State**: `[x]` for settled or completed context, and `[ ]` for active resumption items.
- Store decisions by theme, not one entry per conversational turn.
- Compact overlapping decisions while preserving the **Decision**, **Reason**, **Rejected**, **Reason rejected**, and **Reconsider if** details.

## Workflow

### Start Or Resume

1. Locate the anchor before asking the user to restate context.
2. Read the anchor and any referenced spec, plan, issue, or ADR needed for the current task.
3. Summarize the current state briefly before continuing if resumption context matters.
4. If no anchor exists and the work appears non-trivial, propose or create one before context starts accumulating in chat.

### During Work

Update the anchor at natural boundaries:

- A meaningful decision is made.
- A rejected alternative becomes important to remember.
- A constraint is discovered or changed.
- An open question is created, answered, or made irrelevant.
- Implementation or investigation state changes enough to affect resumption.
- A teammate or subagent contributes context future sessions must respect.

Capture outcomes and rationale. Do not copy hidden interactions, transcripts, or exhaustive tool findings.

If an execution progress file exists, keep mechanical status there. The anchor may link to it, but should only summarize execution state when that summary affects cross-session understanding.

Compact the anchor when:

- The anchor exceeds roughly 200-300 lines.
- Several related decisions accumulate in one session.
- The quick summary no longer lets a new session resume quickly.
- The user says the anchor is hard to scan.

### End Of Session Or Handoff

Before stopping long-running work, update:

- Current state.
- Completed work.
- Open questions.
- Decisions made since the last anchor update.
- Deviations from any spec or plan.
- The next recommended step.

Run an anchor hygiene pass during handoff:

- Refresh the Quick Summary.
- Update checkbox status.
- Merge overlapping decisions.
- Move resolved open questions into Decisions.
- Remove stale "in progress" wording.
- Keep only genuinely unresolved Open Questions.

The handoff test is: a new session should be able to read the anchor and continue without anxiety or lengthy reconstruction.

### Spec And ADR Flow

During long brainstorming, the anchor is the read/write working memory. The final design spec is the polished, committed design derived from the matured anchor.

After a spec exists, keep the anchor only for active-work continuity: implementation state, deviations, new constraints, and new decisions. When the work settles, graduate durable architectural or domain decisions into ADRs or permanent documentation, then archive or stop updating the anchor unless the feature becomes active again.

Once an anchor is inactive, do not chase later code changes back into it. If later work invalidates a durable architectural decision, update the ADR or permanent documentation that owns that decision. Only reopen the feature anchor when the feature itself becomes active again.

## Quality Check

Before presenting or relying on an anchor, verify:

- It explains why key decisions were made.
- It distinguishes settled decisions from open questions.
- It names rejected alternatives that future sessions might otherwise re-propose.
- Its Quick Summary is self-contained when present.
- It has a clear next step when work is active.
- It is concise enough to load quickly in a future session.
- It does not contradict referenced specs, plans, ADRs, or current implementation state.
