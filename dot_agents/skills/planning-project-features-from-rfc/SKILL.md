---
name: planning-project-features-from-rfc
description: Create implementation plans from a reviewed RFC. Uses the RFC as the approved design baseline, decomposes it into executable sub-plans, and runs RFC-specific plan review.
---

# Planning Project Features From RFC

Create actionable implementation plans for features within a single project from a reviewed RFC. The RFC owns requirements, design decisions, constraints, goals, non-goals, risks, tradeoffs, and contracts. This skill turns that approved design into executable sub-plans.

This file is the RFC-backed planning spine. Load referenced files only when their trigger applies.

## Runtime Binding

Before doing any work, determine the active runtime from the system prompt and environment banner, then read exactly one adapter:

- **OpenCode runtime**: [references/runtime-opencode.md](references/runtime-opencode.md)
- **Codex runtime**: [references/runtime-codex.md](references/runtime-codex.md)
- **Claude runtime**: [references/runtime-claude.md](references/runtime-claude.md)

If the runtime signal is ambiguous, ask the user rather than guessing. Do not load or mix other runtime adapters in the same turn. If a runtime adapter conflicts with this file, this file is authoritative.

Use runtime-native names when reading or writing concrete artifacts: worker agent definitions, dispatch recipes, custom subagents, reviewer bindings, or equivalent runtime bindings.

## Core Invariants

- **RFC baseline is authoritative**: the reviewed RFC is the approved source of design decisions, constraints, goals, non-goals, risks, and contracts.
- **Do not re-litigate design**: planning does not reopen RFC architecture, risk, or tradeoff decisions.
- **No silent deviations**: if planning requires changing the RFC design, stop and ask whether to revise the RFC or approve an explicit plan deviation.
- **Direct-planning execution discipline still applies**: RFC-backed planning uses the same decomposition, model selection, execution binding, worker dispatch, documentation sub-plan, and approval mechanics unless this skill explicitly replaces them.
- **Sub-plans are self-contained**: executing agents should not need to read the RFC, master plan, or other sub-plans to understand assigned work.
- **DAG independence is not workspace safety**: same-group sub-plans are logically independent only. Concurrent execution still requires an explicit concurrency policy, isolated implementer worktrees, and verified output/cache safety.
- **RFC-backed review is narrow**: use `plan-rfc-fidelity-reviewer` and `plan-executability-reviewer`. Do not run direct-planning design reviewers unless the user explicitly exits RFC-backed planning.

## Conditional References

| Trigger | Load |
|---------|------|
| Validating the RFC, reading anchors, or exploring planning mechanics | [references/rfc-baseline-exploration.md](references/rfc-baseline-exploration.md) |
| Decomposing RFC scope into sub-plans or defining cross-boundary contracts | [references/decomposition.md](references/decomposition.md) |
| Determining whether independent sub-plans may execute in parallel | [references/concurrency-policy.md](references/concurrency-policy.md) |
| Creating plan files, assigning models, or establishing execution bindings | [references/plan-creation-and-bindings.md](references/plan-creation-and-bindings.md) |
| Running plan review, presenting for approval, or handling feedback | [references/review-and-approval.md](references/review-and-approval.md) |
| Planning documentation updates | [references/documentation-sub-plan.md](references/documentation-sub-plan.md) |
| Writing plan files | [assets/master-plan-template.md](assets/master-plan-template.md), [assets/sub-plan-template.md](assets/sub-plan-template.md) |

## Workflow

### 1. Validate RFC Baseline

Use this workflow only after `planning-project-features` routed here because a reviewed RFC exists and the user wants RFC-backed planning.

Read the RFC and verify that required RFC reviews have no blocking status. Read an active feature anchor when one exists, but treat it as supporting handoff context only. If anchor content conflicts with the RFC or exposes unresolved decomposition, scope, acceptance-criteria, or contract questions, stop and ask.

Use [rfc-baseline-exploration.md](references/rfc-baseline-exploration.md) for validation details. If validation fails, return to `planning-project-features` routing instead of continuing.

### 2. Explore For Planning Mechanics

Read existing documentation first, then inspect code only for planning mechanics the RFC or docs do not cover.

Confirm:

- file paths and ownership boundaries
- existing interfaces, commands, schemas, config files, or runtime bindings named by the RFC
- test, package, and verification scopes that acceptance criteria can reference
- required skills for execution
- language, build-system, output-path, and cache constraints that determine whether execution must be linear or may use parallel groups
- documentation gaps and RFC/specification gaps

Do not use exploration to redesign the RFC. If exploration contradicts the RFC or implementation planning depends on behavior the RFC does not define, stop and ask whether to revise the RFC or approve a plan deviation.

Use [rfc-baseline-exploration.md](references/rfc-baseline-exploration.md) for exploration mechanics and gap handling.

### 3. Decompose The Work

Break the RFC into the smallest self-contained sub-plans that form a valid dependency DAG.

Apply [concurrency-policy.md](references/concurrency-policy.md) before assigning sub-plans to the same execution group. A restricted build ecosystem still gets multiple focused sub-plans, but the DAG must be linear.

Each sub-plan must include the RFC context, constraints, contracts, prerequisites, acceptance criteria, primary files, required skills, and execution model needed for independent execution. Preserve RFC non-goals and translate RFC risks into plan mechanics.

Use [decomposition.md](references/decomposition.md) for boundary selection, DAG rules, embedded context, integration contracts, data-flow integrity, and documentation sub-plan triggers.

Present the proposed decomposition to the user before writing plan files when findings affect scope, sequencing, boundaries, or parallelism.

### 4. Create Plans, Models, And Bindings

Only after baseline validation, exploration, and decomposition are complete:

1. Create the plan directory in the correct standalone or epic location.
2. Write `00-master.md` and numbered sub-plan files from the templates, including the master plan's concurrency policy.
3. Assign execution models to every sub-plan.
4. Establish or select runtime-specific execution bindings for every model + skill combination.
5. Create or reuse one shared test-author binding when any sub-plan has testable acceptance criteria.
6. Write lead-agent execution instructions into the master plan.

Use [plan-creation-and-bindings.md](references/plan-creation-and-bindings.md) for plan paths, model selection, binding rules, and master-plan execution mechanics.

### 5. Run RFC-Specific Plan Review

Run the RFC-backed review loop before showing the plan to the user:

- `plan-rfc-fidelity-reviewer`
- `plan-executability-reviewer`

Pass each reviewer the plan directory, RFC path, requested review output path, and review task. Store review artifacts in the plan's `reviews/` directory. Incorporate findings and re-run only affected reviewers until no blocking findings remain.

If a finding requires changing the RFC design, stop and ask whether to revise the RFC or approve an explicit plan deviation.

Use [review-and-approval.md](references/review-and-approval.md) for review artifact handling, reviewer exclusions, final consistency checks, and feedback classification.

### 6. Present For Approval

Before presenting, run the final consistency check in [review-and-approval.md](references/review-and-approval.md). Fix inconsistencies first.

Present the reviewed plan with:

- RFC path
- review summary
- worker-dispatch summary
- explicit approved deviations, or `None`
- any restart/reload warning for newly created persistent bindings

Remind the user that the plan intentionally omits implementation details. They review architecture and constraints now; they review actual code after execution.

When the user requests changes, incorporate them and use [review-and-approval.md](references/review-and-approval.md) to decide whether re-review is required.

### 7. Note Post-Execution Follow-Up

If the project has component documentation, record that `component-docs-reviewer` should run after all sub-plans complete to catch implementation-vs-plan drift in component docs.

## Plan Structures

Read these templates when writing plan files:

- **[Master plan template](assets/master-plan-template.md)**: orchestration document with RFC mapping and lead-agent execution mechanics.
- **[Sub-plan template](assets/sub-plan-template.md)**: self-contained execution unit with RFC context.

## Final Rules

- Never use this workflow with an unreviewed RFC.
- Never change the RFC design silently.
- Never ask the user to restate information that is already in the RFC.
- Always embed execution-critical RFC context into each sub-plan.
- Always decompose into sub-plans; a monolithic plan is a failure mode.
- Always list required skills in every sub-plan.
- Always respect model assignments through the active runtime's real model-selection mechanism.
- Always apply and record the concurrency policy before writing execution groups.
- Never parallelize implementation plans in the concurrency policy's Linear DAG exception list unless project documentation or the user explicitly approves a documented project-specific exception.
- Use runtime-specific execution bindings for multi-sub-plan execution.
- Write lead-agent execution instructions into every multi-sub-plan master plan.
- Always run `plan-rfc-fidelity-reviewer` and `plan-executability-reviewer` before presenting the plan.
- Never run full direct-planning design reviewers in this workflow unless the user explicitly exits RFC-backed planning.
- Save plans to the correct standalone or epic feature location.
