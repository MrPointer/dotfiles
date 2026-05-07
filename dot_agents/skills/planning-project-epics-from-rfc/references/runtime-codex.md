# Runtime Adapter: Codex

Use this adapter only when the active runtime is Codex.

This adapter maps `planning-project-epics-from-rfc` to Codex-native mechanics. It must not redefine phases, routing, review policy, or epic semantics.

## Discovery

- Treat the current session's available agent list as authoritative for global Codex reviewers.
- Discover project-local reviewer agents from `.codex/agents/` when available.
- If working from a dotfile source repository, use source-directory equivalents rather than applied home-directory paths.

## Required Reviewers

RFC-backed epic planning uses exactly these reviewers before user approval:

- `plan-rfc-fidelity-reviewer`
- `plan-architect-reviewer`

Do not launch `plan-risk-reviewer`; RFC-level risk belongs to `rfc-risk-reviewer` before epic planning starts.

## Reviewer Bindings

Use Codex's custom-agent or reviewer dispatch mechanism. Do not synthesize reviewer prompts ad hoc.

Pass only:

- Epic plan file path.
- RFC path.
- Requested review output path, such as `plans/epics/reviews/<epic-name>.rfc-fidelity.md`.
- The review task.

If direct file writing is supported, ask the reviewer to write to the requested path. Otherwise have it return structured findings and write the review file from the parent agent.

## Review Artifact Ownership

The epic planner owns `plans/epics/reviews/` and must verify each expected review artifact exists after reviewer completion.
