# Runtime Adapter: OpenCode

Use this adapter only when the active runtime is OpenCode.

This adapter maps `planning-project-epics-direct` to OpenCode-native mechanics. It must not redefine phases, routing, review policy, or epic semantics.

## Discovery

Discover reviewer agents from project-local `.opencode/agents/` first, then global `~/.config/opencode/agents/`. OpenCode also discovers compatible agents in other configured agent directories, but use OpenCode-native agents when present.

## Required Reviewers

Direct epic planning uses exactly these reviewers before user approval:

- `design-reviewer`

## Reviewer Bindings

Invoke reviewers by custom subagent name through OpenCode's normal subagent mechanism. Do not recreate reviewer personas in prompt text.

Pass only:

- Epic plan file path.
- Requested review output path, such as `plans/epics/reviews/<epic-name>.design.md`.
- The review task.

If a reviewer writes the review file directly, use it. If it returns findings, write the review file from the response. Verify the review artifact exists before continuing.

## Review Artifact Ownership

The epic planner owns `plans/epics/reviews/` and must verify each expected review artifact exists after reviewer completion.
