---
name: planning-project-epics
description: Route epic planning to the correct workflow. Recommends RFC-first design when no RFC exists, sends reviewed RFCs to RFC-backed epic planning, and sends explicit opt-outs to direct epic planning.
---

# Planning Project Epics

Route epic planning deterministically. This skill does not create epic plans itself; it chooses the correct epic planning workflow and then stops or delegates.

## Routing Table

| Input State | Action | Epic Planning Workflow | Reviewers |
|-------------|--------|------------------------|-----------|
| Reviewed RFC provided | Validate the RFC review record, then use the RFC as the design baseline | `planning-project-epics-from-rfc` | `plan-rfc-fidelity-reviewer`, `design-reviewer` |
| RFC provided but not reviewed | Stop and offer two choices: review the RFC first, or proceed as direct epic planning | `authoring-rfcs` review or `planning-project-epics-direct` | Depends on chosen path |
| No RFC provided | Recommend creating an RFC first, then ask whether to create one or continue directly | `brainstorming` + `authoring-rfcs`, or `planning-project-epics-direct` | Depends on chosen path |

## Workflow

### Step 1: Check For An RFC

If the user provides or references an RFC, read it and inspect its **Review Record**.

A reviewed RFC has:

- `design-reviewer` with status `Passed` or `Passed with concerns`.
- No `Blocking` review status that affects epic scope, feature boundaries, sequencing, migration, rollout, or cross-feature constraints.

If these conditions hold, use `planning-project-epics-from-rfc`. Do not re-gather requirements or re-litigate the RFC design.

If an RFC exists but does not meet these conditions, stop and present exactly two choices:

1. Review or revise the RFC first with `authoring-rfcs`.
2. Proceed with direct epic planning using `planning-project-epics-direct`, treating the RFC as background context rather than an approved baseline.

### Step 2: Recommend RFC-First When Missing

If no RFC exists, recommend creating one before epic planning. Explain the concrete benefits briefly:

- It separates design decisions from feature decomposition.
- It gives the design reviewer the right artifact to review.
- It gives epic planning a stable, reviewed baseline for feature boundaries and sequencing.
- It reduces token waste by narrowing epic review to RFC fidelity and feature-decomposition quality.
- It creates a stable RFC document reference for the epic plan and future feature plans.
- It preserves user-approved decisions so the epic planner does not reopen them accidentally.

Then ask the user to choose exactly one path:

1. Create or finish the RFC first with `brainstorming` and `authoring-rfcs`.
2. Continue directly with `planning-project-epics-direct`.

Do not continue silently. The user chooses the route.

### Step 3: Delegate

After the route is chosen, load exactly one workflow skill:

- `planning-project-epics-from-rfc` for reviewed RFC-backed epic planning.
- `planning-project-epics-direct` for direct epic planning.
- `brainstorming` and then `authoring-rfcs` when the user chooses RFC-first design work.

Do not mix workflows in one epic planning run.

## Rules

- Do not create epic plans in this router skill.
- Do not make conditional reviewer choices here.
- Do not treat an unreviewed RFC as an approved design baseline.
- Do not ask the user to repeat information already present in a reviewed RFC.
- If the chosen path becomes invalid, stop and return to the routing table.
