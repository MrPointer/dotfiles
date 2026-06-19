---
name: planning-project-features-direct
description: Create implementation plans directly from user requirements when no reviewed RFC is available or the user explicitly declines RFC-first planning. Decomposes work into self-contained sub-plans with full iterative multi-agent plan review.
---

# Planning Project Features Direct

Create actionable implementation plans for features within a single project when no reviewed RFC baseline exists or the user explicitly chose direct planning.

Direct planning owns requirement clarification, planning-focused exploration, decomposition, reviewer assignment, model selection, execution bindings, and iterative plan review. Never assume or guess when a requirement or design decision is unresolved.

This file is the direct planning spine. Load referenced files only when their trigger applies.

## Runtime Binding

Before doing any work, determine the active runtime from the system prompt and environment banner, then read exactly one adapter:

- **OpenCode runtime**: [references/runtime-opencode.md](references/runtime-opencode.md)
- **Codex runtime**: [references/runtime-codex.md](references/runtime-codex.md)
- **Claude runtime**: [references/runtime-claude.md](references/runtime-claude.md)

If the runtime signal is ambiguous, ask the user rather than guessing. Do not load or mix other runtime adapters in the same turn. If a runtime adapter conflicts with this file, this file is authoritative.

Use runtime-native names when reading or writing concrete artifacts: worker agent definitions, dispatch recipes, custom subagents, reviewer bindings, or equivalent runtime bindings.

## Core Invariants

- **No assumptions**: if something is unclear, ambiguous, or missing, ask. Do not fill gaps with reasonable defaults or best guesses.
- **Clarification continues throughout planning**: if later phases surface new ambiguity or an unresolved decision, stop and present it to the user.
- **Sub-plans are self-contained**: executing agents should not need to read the master plan, other sub-plans, or external documents to understand assigned work.
- **DAG independence is not workspace safety**: same-group sub-plans are logically independent only. Concurrent execution still requires isolated implementer worktrees and any required build/cache seeding.
- **Plans omit implementation details**: plans define what to build and why; required skills define how to code, test, lint, build, and document.
- **Convergent review is required**: direct plans go through global and project-local review until no blocking findings remain before user presentation.
- **Reviewer bindings are real artifacts**: do not replace required local reviewers with generic/global reviewers or ad hoc prompts.

## Conditional References

| Trigger | Load |
|---------|------|
| Gathering requirements, reading anchors, or exploring planning mechanics | [references/requirements-exploration.md](references/requirements-exploration.md) |
| Decomposing work into sub-plans or defining cross-boundary contracts | [references/decomposition.md](references/decomposition.md) |
| Creating plan files, assigning reviewers/models, or establishing execution bindings | [references/plan-creation-reviewers-and-bindings.md](references/plan-creation-reviewers-and-bindings.md) |
| Running initial review, presenting for approval, or handling feedback | [references/review-and-approval.md](references/review-and-approval.md) |
| Planning documentation updates | [references/documentation-sub-plan.md](references/documentation-sub-plan.md) |
| Writing plan files | [assets/master-plan-template.md](assets/master-plan-template.md), [assets/sub-plan-template.md](assets/sub-plan-template.md) |

## Workflow

### 1. Gather Requirements

Use this workflow only after `planning-project-features` routed here because no reviewed RFC is available or the user explicitly chose direct planning.

Clarify goal, scope, dependencies, constraints, and acceptance criteria until you can plan confidently. If the user asks you to use judgment, first present options and tradeoffs. If they still choose that path, proceed with best judgment and document assumptions clearly in the plan.

Read an active feature anchor when one exists, but treat it as supporting context only. It does not replace user-approved requirements. Any unresolved anchor question that affects decomposition, scope, acceptance criteria, or contracts is a planning blocker.

Use [requirements-exploration.md](references/requirements-exploration.md) for clarification questions, anchor handling, and assumption handling.

### 2. Explore For Planning Mechanics

Read existing documentation first, then inspect code only for planning mechanics that docs and requirements do not cover.

Confirm:

- files and ownership boundaries
- architectural constraints and patterns to follow
- potential conflicts and risks
- required skills for execution
- ignored build/cache directories needed by isolated TDD or parallel implementation worktrees, or that no seeding is required
- documentation gaps and specification gaps

Share findings with the user and confirm understanding before proceeding when findings affect scope, sequencing, boundaries, or risk.

Use [requirements-exploration.md](references/requirements-exploration.md) for exploration mechanics and gap handling.

### 3. Decompose The Work

Break the feature into the smallest self-contained sub-plans that form a valid dependency DAG.

Each sub-plan must include the domain context, constraints, contracts, prerequisites, acceptance criteria, primary files, required skills, reviewer assignment, and execution model needed for independent execution.

Use [decomposition.md](references/decomposition.md) for boundary selection, DAG rules, embedded context, integration contracts, data-flow integrity, skills-vs-plan boundaries, decisiveness, and documentation sub-plan triggers.

Present the proposed decomposition to the user before writing plan files when findings affect scope, sequencing, boundaries, or parallelism.

### 4. Create Plans, Reviewers, Models, And Bindings

Only after requirements, exploration, and decomposition are complete:

1. Create the plan directory in the correct standalone or epic location.
2. Write `00-master.md` and numbered sub-plan files from the templates.
3. Discover project-local reviewer bindings.
4. Assign the best local reviewer to each sub-plan, or warn the user when coverage is missing.
5. Assign execution models to every sub-plan.
6. Establish or select runtime-specific execution bindings for every required model + skill combination.
7. Create or reuse one shared test-author binding when any sub-plan has testable acceptance criteria.
8. Write lead-agent execution instructions into the master plan.

Use [plan-creation-reviewers-and-bindings.md](references/plan-creation-reviewers-and-bindings.md) for plan paths, reviewer assignment, model selection, binding rules, and master-plan execution mechanics.

### 5. Run Initial Review Loop

Run the full direct-planning review loop once before showing the plan to the user.

Default global master-plan reviewers:

- `design-reviewer`
- `plan-clarity-reviewer`
- `plan-executability-reviewer`

Then run each sub-plan's assigned project-local reviewer. Incorporate findings and re-run only affected reviewers until no blocking findings remain.

Use [review-and-approval.md](references/review-and-approval.md) for reviewer launch rules, artifact ownership, output normalization, planning-rule validation, convergence, and optional specialized reviewers.

### 6. Present For Approval

Before presenting, run the final consistency check in [review-and-approval.md](references/review-and-approval.md). Fix inconsistencies first.

Present the reviewed plan with:

- review summary
- worker-dispatch summary
- how findings were addressed
- any missing local reviewer coverage approved by the user
- any restart/reload warning for newly created persistent bindings

Remind the user that the plan intentionally omits implementation details. They review architecture and constraints now; they review actual code after execution.

When the user requests changes, incorporate them and use [review-and-approval.md](references/review-and-approval.md) to decide whether re-review is recommended or required.

### 7. Note Post-Execution Follow-Up

If the project has component documentation, record that `component-docs-reviewer` should run after all sub-plans complete to catch implementation-vs-plan drift in component docs.

## Plan Structures

Read these templates when writing plan files:

- **[Master plan template](assets/master-plan-template.md)**: orchestration document without implementation details.
- **[Sub-plan template](assets/sub-plan-template.md)**: self-contained execution unit.

## Final Rules

- Never write a plan based on incomplete information.
- Never invent requirements the user did not specify.
- Always ask for clarification when ambiguity affects scope, contracts, architecture, or acceptance criteria.
- Always decompose into sub-plans; a monolithic plan is a failure mode.
- Always embed execution-critical context into each sub-plan.
- Always list required skills in every sub-plan.
- Always run the initial review loop before presenting to the user.
- Always respect model assignments through the active runtime's real model-selection mechanism.
- Use runtime-specific execution bindings for multi-sub-plan execution.
- Write lead-agent execution instructions into every multi-sub-plan master plan.
- Save plans to the correct standalone or epic feature location; never use runtime metadata directories or random/generated filenames.
