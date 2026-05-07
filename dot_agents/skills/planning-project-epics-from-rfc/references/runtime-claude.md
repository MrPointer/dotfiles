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

Launch reviewer agents with `subagent_type: "general-purpose"` so they inherit their file-defined tools, including `Write`/`Edit` for review output.

Pass only:

- Epic plan file path.
- RFC path.
- Requested review output path, such as `plans/epics/reviews/<epic-name>.rfc-fidelity.md`.
- The review task.

If a reviewer writes the review file directly, use it. If it returns findings, write the review file from the response. Check whether the expected file exists after the reviewer finishes.

## Review Artifact Ownership

The epic planner owns `plans/epics/reviews/` and must verify each expected review artifact exists after reviewer completion.
