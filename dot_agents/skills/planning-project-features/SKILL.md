---
name: planning-project-features
description: Route feature planning to the correct workflow. Recommends RFC-first design when no RFC exists, sends reviewed RFCs to RFC-backed planning, and sends explicit opt-outs to direct planning.
---

# Planning Project Features

Route feature planning deterministically. This skill does not create plans itself; it chooses the correct planning workflow and then stops or delegates.

## Routing Table

| Input State | Action | Planning Workflow | Reviewers |
|-------------|--------|-------------------|-----------|
| Reviewed RFC provided | Validate the RFC review record, then use the RFC as the design baseline | `planning-project-features-from-rfc` | `plan-rfc-fidelity-reviewer`, `plan-executability-reviewer` |
| RFC provided but not reviewed | Stop and offer two choices: review the RFC first, or proceed as direct planning | `authoring-rfcs` review or `planning-project-features-direct` | Depends on chosen path |
| No RFC provided | Recommend creating an RFC first, then ask whether to create one or continue directly | `brainstorming` + `authoring-rfcs`, or `planning-project-features-direct` | Depends on chosen path |

## Workflow

### Step 1: Check For An RFC

If the user provides or references an RFC, read it and inspect its **Review Record**.

A reviewed RFC has:

- `rfc-architect-reviewer` with status `Passed` or `Passed with concerns`.
- `rfc-risk-reviewer` with status `Passed` or `Passed with concerns`.
- No `Blocking` review status that affects decomposition, scope, contracts, migration, rollout, acceptance criteria, or implementation constraints.

If these conditions hold, use `planning-project-features-from-rfc`. Do not re-gather requirements or re-litigate the RFC design.

If an RFC exists but does not meet these conditions, stop and present exactly two choices:

1. Review or revise the RFC first with `authoring-rfcs`.
2. Proceed with direct planning using `planning-project-features-direct`, treating the RFC as background context rather than an approved baseline.

### Step 2: Recommend RFC-First When Missing

If no RFC exists, recommend creating one before planning. Explain the concrete benefits briefly:

- It separates design decisions from execution decomposition.
- It gives architecture and risk reviewers the right artifact to review.
- It gives planning a stable, reviewed baseline instead of forcing it to rediscover the design.
- It reduces token waste in planning by narrowing plan review to fidelity and executability.
- It creates a stable `RFC-0001`-style reference for future plans, implementation sessions, and follow-up discussion.
- It preserves user-approved decisions so the planner does not reopen them accidentally.

Then ask the user to choose exactly one path:

1. Create or finish the RFC first with `brainstorming` and `authoring-rfcs`.
2. Continue directly with `planning-project-features-direct`.

Do not continue silently. The user chooses the route.

### Step 3: Delegate

After the route is chosen, load exactly one workflow skill:

- `planning-project-features-from-rfc` for reviewed RFC-backed planning.
- `planning-project-features-direct` for direct planning.
- `brainstorming` and then `authoring-rfcs` when the user chooses RFC-first design work.

Do not mix workflows in one planning run.

## Rules

- Do not create plans in this router skill.
- Do not make conditional reviewer choices here.
- Do not treat an unreviewed RFC as an approved design baseline.
- Do not ask the user to repeat information already present in a reviewed RFC.
- If the chosen path becomes invalid, stop and return to the routing table.
