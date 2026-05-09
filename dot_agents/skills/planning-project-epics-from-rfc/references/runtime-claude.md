# Runtime Adapter: Claude

Use this adapter only when the active runtime is Claude.

This adapter maps `planning-project-epics-from-rfc` to Claude-native mechanics. It must not redefine phases, routing, review policy, or epic semantics.

## Discovery

Discover reviewer agents from project-local `.claude/agents/` first, then global `~/.claude/agents/`. Project-local wins when the same name exists in both locations.

## Required Reviewers

RFC-backed epic planning uses exactly these reviewers before user approval:

- `plan-rfc-fidelity-reviewer`
- `plan-architect-reviewer`

Do not launch `plan-risk-reviewer`; RFC-level risk belongs to `rfc-risk-reviewer` before epic planning starts.

## Reviewer Bindings

Launch each reviewer agent through its own `subagent_type` (for example, `subagent_type: "plan-rfc-fidelity-reviewer"`). The named agent's frontmatter declares both its persona and its `tools` list — including `Write`/`Edit` for write-capable reviewers — so direct dispatch loads both the right system prompt and the right tools. Do not launch reviewers as `general-purpose`: that bypasses the reviewer's specialized system prompt and replaces its declared tool list with general-purpose's broader set.

Pass only:

- Epic plan file path.
- RFC path.
- Requested review output path, such as `plans/epics/reviews/<epic-name>.rfc-fidelity.md`.
- The review task.

If a reviewer writes the review file directly, use it. If it returns findings, write the review file from the response. Check whether the expected file exists after the reviewer finishes.

## Review Artifact Ownership

The epic planner owns `plans/epics/reviews/` and must verify each expected review artifact exists after reviewer completion.
