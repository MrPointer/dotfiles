---
name: planning-project-epics-from-rfc
description: Decompose a reviewed RFC into sequenced features for a single project. Uses the RFC as the approved design baseline and produces a persistent epic plan reviewed for RFC fidelity and feature decomposition.
---

# Planning Project Epics From RFC

Decompose a reviewed RFC into sequenced, independently plannable features. The RFC owns the design; this skill owns feature boundaries, feature sequencing, cross-feature context, and the persistent epic coordination document.

Each feature in the epic plan is described richly enough that `planning-project-features` can route and plan it separately. The epic planner does **not** invoke feature planning.

## Runtime Binding

This skill has one canonical workflow. Runtime files only map that workflow to the active agent runtime's mechanics.

Before doing any work, determine the active runtime and read exactly one adapter:

- **OpenCode runtime** -> [references/runtime-opencode.md](references/runtime-opencode.md)
- **Codex runtime** -> [references/runtime-codex.md](references/runtime-codex.md)
- **Claude runtime** -> [references/runtime-claude.md](references/runtime-claude.md)

Do not load or mix instructions from another runtime adapter in the same turn. If a runtime adapter conflicts with this file, this file is authoritative.

## Core Principles

1. **RFC Baseline**: Treat the reviewed RFC as the approved source of design decisions, constraints, goals, non-goals, risks, and planning handoff context.
2. **Epic Owns Feature Decomposition**: The epic planner creates feature boundaries, feature sequencing, status tracking, and cold-start context for later feature planning.
3. **No Cross-Feature Contracts**: Note feature dependencies but do not define interface signatures, type shapes, API contracts, or execution-level wiring. Feature planning owns those decisions.
4. **No Design Re-Litigation**: Do not reopen RFC architecture, risk, or tradeoff decisions during epic planning.
5. **No Silent Deviations**: If epic decomposition requires changing the RFC design, stop and ask whether to revise the RFC or explicitly approve an epic deviation.

## Workflow

### Phase 1: Validate RFC Baseline

Read the RFC and confirm:

- `design-reviewer` is `Passed` or `Passed with concerns`.
- No `Blocking` review status remains.
- Any `Passed with concerns` item is compatible with epic decomposition and does not require a design decision before feature boundaries can be chosen.

If validation fails, stop and return to `planning-project-epics` routing. Do not continue as RFC-backed epic planning.

Read the active feature or epic anchor if one exists. The anchor may add handoff context, but it does not override the RFC.

### Phase 2: Planning-Focused Exploration

Read existing docs first, then inspect code only for epic mechanics:

- Confirm existing systems, packages, or boundaries named by the RFC.
- Identify natural feature seams implied by existing architecture and the RFC design.
- Identify required skills later feature planners and executing agents will need.
- Identify workspace/build-cache constraints later feature planners must preserve, including ignored in-repository build/cache directories that isolated worktrees need in build-heavy ecosystems.
- Identify documentation gaps that could make future feature planning difficult.

Do not use exploration to redesign the RFC. If exploration contradicts the RFC, stop and ask whether to revise the RFC or abandon RFC-backed epic planning.

### Phase 3: Feature Decomposition

Break the RFC into features:

1. **Apply the feature-worthiness test**: Each feature must have meaningful design surface for later feature planning.
2. **Map RFC coverage**: Every RFC goal, non-goal, constraint, major design area, risk, and success criterion should be covered by one or more features or explicitly marked as cross-cutting context.
3. **Find natural boundaries**: Prefer features split by domain, architectural layer, infrastructure concern, migration phase, or independently meaningful design tradeoff.
4. **Sequence by feature dependency**: Feature dependencies are existence-level dependencies only. Do not specify contracts.
5. **Identify parallel opportunities**: Features with no dependency between them can be planned independently.
6. **Capture build-heavy workspace constraints**: If isolated feature execution will be expensive without ignored build/cache artifacts, record the cache directories as a cross-cutting concern for later feature planning. Do not define implementation steps; just preserve the constraint.
7. **Capture cross-cutting concerns**: RFC decisions that span features must be recorded so later feature plans do not diverge.

Present the decomposition to the user for review before writing the epic plan.

### Phase 4: Epic Plan Creation

Create the epic plan file using this skill's [RFC-backed epic plan template][epic-plan-template]:

```text
plans/epics/<epic-name>.md
```

Each feature description must contain enough context for a feature planner starting cold:

- What the feature accomplishes and why it is distinct.
- RFC sections and decisions the feature must preserve.
- Domain knowledge, relevant existing code, and documentation pointers.
- Dependencies on other features by existence only.
- Scope boundaries and non-goals.
- Required skills likely needed during feature planning and execution.
- Known build-heavy workspace constraints, including ignored build/cache directories isolated worktrees need available.

Set all features to `not-started`.

### Phase 5: Review Loop

Run only these reviewers before presenting to the user:

- **`plan-rfc-fidelity-reviewer`**: Checks that the epic faithfully decomposes the RFC without omitting goals, violating non-goals, or adding unapproved scope.
- **`design-reviewer`**: Checks whether feature boundaries preserve RFC design decisions, integration seams, compatibility constraints, migration/rollback strategy, and epic-specific risks introduced by the decomposition.

Pass the epic plan path, RFC path, and review output path to each reviewer. Review output lives under `plans/epics/reviews/`.

Incorporate findings and re-run only affected reviewers until no blocking findings remain. If a finding requires changing the RFC design, stop and ask whether to revise the RFC or approve an explicit deviation.

Do not re-review RFC-level risk in this workflow. RFC-level risk was already reviewed by `design-reviewer`; epic planning should preserve those risks as feature context and sequencing constraints. The design reviewer should focus risk findings on epic-specific boundary, integration, compatibility, migration, rollback, and hidden-coupling risks.

### Phase 6: User Approval And Feedback

Present the reviewed epic plan with the RFC path, review summary, feature decomposition, and any explicit approved deviations.

When the user requests changes, classify:

| Change Type | Examples | Re-review Action |
|---|---|---|
| Cosmetic / wording | Clarify a description, fix typos | None |
| Scope adjustment within a feature | Expand/narrow what a feature covers without changing RFC coverage | Re-run `plan-rfc-fidelity-reviewer` only if RFC coverage changes |
| Feature boundary change | Merge features, split a feature, move scope between features | Re-run both reviewers |
| RFC deviation | Add scope, drop RFC scope, change sequencing because RFC design does not fit | Stop for RFC revision or explicit deviation approval |

The epic plan is a living document. Feature planning or execution may reveal better feature boundaries later; update the epic deliberately and re-review affected structure when necessary.

## Rules

- Never use this workflow with an unreviewed RFC.
- Never change the RFC design silently.
- Never define cross-feature contracts.
- Never re-review RFC-level risk in this workflow.
- Never ask the user to restate information that is already in the RFC.
- Always preserve RFC goals, non-goals, constraints, and risks in the epic plan.
- Always run `plan-rfc-fidelity-reviewer` and `design-reviewer` before presenting the epic plan.

[epic-plan-template]: assets/epic-plan-template.md
