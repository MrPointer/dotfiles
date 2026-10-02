# Human-Readable Design Doc Workflow Anchor

## Intent

Decide where human-readable design docs fit in the planning workflow. Docs that
humans and agents both use pull in opposite directions:

- Humans need the narrative: what was decided and why. Dense agent-facing RFCs
  took days to read, and small nuances were forgotten by the end.
- Agents need code-grounded facts: paths, contracts, current state.

The `authoring-human-readable-design-docs` skill produces readable docs, but
those docs have no code grounding and no review.

## Quick Summary

| Settled Direction | What It Means |
|-------------------|---------------|
| The legacy skill workflow (`brainstorming` → `authoring-rfcs` → `planning-project-features*`) is frozen; the Spec Kit preset `timors-agentic-workflow` is the target. | Do not wire the design-doc skill into the legacy brainstorming or planning skills. It stays a standalone tool. |
| The RFC stage is not wanted in any form. | Do not propose RFC-based pipelines, including "RFC for agents, design doc generated from it for humans". |
| A reviewer finding that changes design content needs human approval, even when it fills a gap nobody decided. | This is in the global rules and in Spec Kit preset `0.4.1`. Only wording, formatting, or citation fixes are applied directly. |
| For the pressured feature, the design doc wins over its older reviewed RFC. | No comparison is needed. The user handles it in that feature's own repository. |

## Current State

- [x] Analyzed the legacy skills. Planning only accepts an RFC with a passing
      `design-reviewer` Review Record. The design-doc skill has no review step,
      no save path, and is marked opt-in only.
- [x] Reviewed how Spec Kit handles design (see Spec Kit Design Model below).
- [x] Committed the global "Decision Ownership" rule (`201e889`).
- [x] Committed the Spec Kit remediation gate and version bump to `0.4.1`
      (`8ac543c`).
- [ ] User: run `chezmoi apply` outside the sandbox, and refresh the preset in
      projects that use it (remove/add, per the preset README).
- [ ] User, in the other repository: make the pressured feature's RFC match the
      design doc, then plan it with `planning-project-features-from-rfc`.
- [ ] Decide the design doc's place in Spec Kit (first open question below).
      This is the next design step.

## Spec Kit Design Model

The preset splits design by ownership, not by audience. It has no RFC.

- `spec.md` (Feature Definition): intent, scope, requirements, and the
  human-owned decisions. Code feasibility is excluded on purpose.
- `plan.md`: the only normative plan. It holds current state, decisions,
  architecture, contracts, risks, and the task-planning handoff. It is grounded
  in the code.
- `grounding-notes.md`: repository evidence only, with no recommendations.
  `research.md` is non-normative scratch space.
- The read-only reviewers `plan-design` and `plan-clarity` write only their own
  reports. A human records Design Acceptance in `plan.md` before `tasks` can
  run.
- Planning may fix only typos or metadata in `spec.md`. Any change to intent
  sends `spec.md` back to `Draft` and the `clarify` step.
- A guided human walkthrough was added in #26 and removed in #28. Since then,
  `plan.md` has no human-readable view.

## Constraints

- The Spec Kit preset must stay agent-agnostic and shareable with the team. The
  design-doc skill depends on `mid-tier-technical-writer`, which is defined per
  runtime. A Spec Kit step would need a provider-neutral binding, meaning a tier
  plus a role description, the way the reviewer packets work.
- Spec Kit's review and implementation steps have never run on a real feature.
- A design doc derived from another artifact can drop or bend nuances. Any
  derived doc needs a fidelity check against its source.

## Decisions

### Legacy Workflow

**Decision:** Freeze the legacy skills. Do not change the brainstorming routing
or the planners for design docs.

**Reason:** The user is moving to Spec Kit. The design-doc skill was added to
the legacy flow only under delivery pressure.

**Rejected:** Extending the legacy flow, which would mean:

- brainstorming routes to the design doc by default
- `planning-project-features-from-rfc` accepts design docs
- a `docs/design/` save path
- a `design-reviewer` result kept in a separate file

**Reason rejected:** It invests in a workflow that is being retired.

**Reconsider if:** The move to Spec Kit stalls.

### RFC Stage

**Decision:** Drop the RFC stage, including as an agent-only artifact behind a
human design doc.

**Reason:** It cost large amounts of tokens and time. Across several designs,
reviewer feedback was applied without asking, so the RFC carried decisions the
user never saw.

**Rejected:** "Paired" pipeline: brainstorm → reviewed RFC → design doc for
human approval → plan from the RFC.

**Reason rejected:** Same cost and the same risk of unconsulted decisions. One
earlier run worked only because the RFC already existed and was reviewed.

### Reviewer Authority

**Decision:** Agents present every finding that adds, removes, or changes design
or plan content, and apply only what the human accepts. Wording, formatting, and
citation fixes may be applied directly and are listed when reporting.

**Reason:** The old rule protected only decisions the user had already stated.
Reviewers mostly flag gaps, and gaps were being filled silently.

**Where:**

- The global rules come from `.chezmoitemplates/global-ai-core-guidelines.md.tmpl`
  and render to the Claude, Codex, and OpenCode global files.
- Spec Kit `references/review-lifecycle.md` defines "proposed remediations",
  which are shown together with blockers at review quiescence.
- `commands/speckit.plan.md` step 4 points to that rule.

### Pressured Feature

**Decision:** The design doc is authoritative. Make the RFC match it, then plan
from the RFC.

**Reason:** It is the fastest path that uses existing, working tooling.

**Recommended safeguards:**

- Commit the reconciled RFC before any review, so `git diff` shows later drift.
- Tell the session that review is report-only.

## Open Questions

- [ ] What is the design doc's role in Spec Kit? Proposal not yet accepted:
      - a non-normative human view of `plan.md`
      - generated after reviews pass and before Design Acceptance
      - checked by a fidelity review against `plan.md`
      - used by the human to decide acceptance
      Still open: where the file lives, whether it is a new command or a step
      in `/speckit.plan`, and whether `spec.md` also needs a human view.
- [ ] Should the removed walkthrough (#26) return, or does the design doc replace
      it?
- [ ] `docs/spec-kit-workflow/architecture.md` still describes the removed
      "Optional Human Review". Fix it when that question settles.

## References

- Design-doc skill: `dot_agents/skills/authoring-human-readable-design-docs/SKILL.md`
- Writer agent: `.chezmoitemplates/agent-mid-tier-technical-writer.md.tmpl`
- Legacy skills: `dot_agents/skills/brainstorming/`, `authoring-rfcs/`,
  `planning-project-features/`, `planning-project-features-from-rfc/`
- Spec Kit preset: `private_dot_config/specify/presets/timors-agentic-workflow/`
- Spec Kit architecture: `docs/spec-kit-workflow/architecture.md`
