---
name: planning-project-features-from-rfc
description: Create implementation plans from a reviewed RFC. Uses the RFC as the approved design baseline, decomposes it into executable sub-plans, and runs only RFC-fidelity and executability review.
---

# Planning Project Features From RFC

Create implementation plans from a reviewed RFC. The RFC owns the design; this skill owns decomposition, execution ordering, file ownership, acceptance criteria, model and skill assignment, and execution bindings.

## Runtime Binding

Use the runtime adapters from `planning-project-features-direct/references/`. They define runtime mechanics for reviewer dispatch, execution bindings, and artifact ownership. They do not change this workflow's routing or review policy.

Before doing any work, determine the active runtime and read exactly one adapter:

- **OpenCode runtime** -> `../planning-project-features-direct/references/runtime-opencode.md`
- **Codex runtime** -> `../planning-project-features-direct/references/runtime-codex.md`
- **Claude runtime** -> `../planning-project-features-direct/references/runtime-claude.md`

## Core Principles

1. **RFC Baseline**: Treat the reviewed RFC as the approved source of design decisions, constraints, goals, non-goals, risks, and contracts.
2. **Planning Owns Mechanics**: The planner creates sub-plan boundaries, dependency graph, execution ordering, file ownership, acceptance criteria, required skills, model assignments, and execution bindings.
3. **No Design Re-Litigation**: Do not reopen RFC architecture, risk, or tradeoff decisions during planning.
4. **No Silent Deviations**: If planning requires changing the RFC design, stop and ask whether to revise the RFC or explicitly approve a plan deviation.
5. **Minimal Review**: Reviewed RFC-backed plans use only `plan-rfc-fidelity-reviewer` and `plan-executability-reviewer` before user approval.

## Workflow

### Phase 1: Validate RFC Baseline

Read the RFC and confirm:

- `rfc-architect-reviewer` is `Passed` or `Passed with concerns`.
- `rfc-risk-reviewer` is `Passed` or `Passed with concerns`.
- No `Blocking` review status remains.
- Any `Passed with concerns` item is compatible with planning and does not require a design decision before decomposition.

If validation fails, stop and return to `planning-project-features` routing. Do not continue as RFC-backed planning.

Read the active feature anchor if one exists. The anchor may add handoff context, but it does not override the RFC.

### Phase 2: Planning-Focused Exploration

Read existing docs first, then inspect code only for planning mechanics:

- Confirm file paths and ownership boundaries needed for sub-plans.
- Confirm existing interfaces, commands, schemas, or config files named by the RFC.
- Identify required skills for execution.
- Identify tests, packages, or verification scopes that acceptance criteria can reference.

Do not use exploration to redesign the RFC. If exploration contradicts the RFC, stop and ask whether to revise the RFC or abandon RFC-backed planning.

### Phase 3: Decomposition

Break the RFC into sub-plans:

1. **Identify execution seams**: Split work by files, layers, domains, and independently verifiable outcomes.
2. **Create a DAG**: The dependency graph is a planning artifact. It must be acyclic, explicit, and based on execution prerequisites created by the sub-plans.
3. **Embed RFC context**: Each sub-plan must include the RFC decisions, constraints, and contract details it needs. Executing agents should not have to read the RFC.
4. **Preserve non-goals**: Do not add work the RFC explicitly excluded.
5. **Translate risks into mechanics**: RFC risks become plan constraints, acceptance criteria, sequencing notes, or handoff constraints where relevant.

Present the decomposition to the user for review before writing plan files.

### Phase 4: Plan Creation And Execution Binding

Create the plan directory and files using the templates in this skill's `assets/` directory:

```text
<plan-directory>/
├── 00-master.md
├── 01-<first-task>.md
├── 02-<second-task>.md
├── reviews/
└── ...
```

The master plan must include the RFC path and review summary. Each sub-plan must include enough RFC-derived context to be self-contained.

Assign required skills and execution models exactly as `planning-project-features-direct` does. Establish execution bindings using the active runtime adapter.

### Phase 5: Minimal Review Loop

Run only these master-plan reviewers:

- **`plan-rfc-fidelity-reviewer`**: Checks that the plan faithfully decomposes the RFC without contradicting it, omitting required context, or adding unapproved design scope.
- **`plan-executability-reviewer`**: Checks file ownership, acceptance criteria, dependency order, verification scope, and isolated execution mechanics.

Pass the plan directory, RFC path, and review output path to each reviewer. Review output lives in the plan's `reviews/` directory.

Incorporate findings and re-run only affected reviewers until both report no blocking findings. If a finding requires changing the RFC design, stop and ask whether to revise the RFC or approve an explicit deviation.

Do not run `plan-architect-reviewer`, `plan-risk-reviewer`, or `plan-clarity-reviewer` in this workflow. Those reviewers belong to direct planning. RFC architecture, risk, and clarity were already reviewed by RFC reviewers.

### Phase 6: User Approval

Before presenting, verify:

1. Every RFC goal and success criterion is represented in the plan.
2. Every RFC non-goal remains out of scope.
3. Every RFC constraint that affects execution appears in the master plan or relevant sub-plan.
4. Every sub-plan is self-contained.
5. The execution order is a valid DAG.
6. `plan-rfc-fidelity-reviewer` and `plan-executability-reviewer` findings are resolved or explicitly documented as non-blocking.

Present the reviewed plan with the RFC path, review summary, and any explicit approved deviations. Only mark ready when the user approves.

## Plan Structures

- **[Master plan template][master-plan-template]**
- **[Sub-plan template][sub-plan-template]**

## Rules

- Never use this workflow with an unreviewed RFC.
- Never change the RFC design silently.
- Never run full direct-planning reviewers in this workflow.
- Never ask the user to restate information that is already in the RFC.
- Always embed execution-critical RFC context into each sub-plan.
- Always run `plan-rfc-fidelity-reviewer` and `plan-executability-reviewer` before presenting the plan.

[master-plan-template]: assets/master-plan-template.md
[sub-plan-template]: assets/sub-plan-template.md
