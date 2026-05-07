# Runtime Adapter: Claude

Use this adapter only when the active runtime is Claude.

This adapter maps `planning-project-features-from-rfc` to Claude-native mechanics. It must not redefine phases, routing, review policy, or plan semantics.

## Discovery

**Skill directories**:

- Global skills: `~/.claude/skills/`
- Project-local skills: `.claude/skills/`

**Agent directories**:

- Global agents: `~/.claude/agents/`
- Project-local agents: `.claude/agents/`

Project-local wins when the same name exists in both locations.

## Required Reviewers

RFC-backed feature planning uses exactly these reviewers:

- `plan-rfc-fidelity-reviewer`
- `plan-executability-reviewer`

Do not launch `plan-architect-reviewer`, `plan-risk-reviewer`, `plan-clarity-reviewer`, or project-local domain reviewers in this workflow unless the user explicitly exits RFC-backed planning and returns to direct planning.

## Reviewer Bindings

Launch reviewer agents with `subagent_type: "general-purpose"` so they inherit their file-defined tools, including `Write`/`Edit` for review output.

Pass only:

- Plan directory path.
- RFC path.
- Requested review output path, such as `reviews/00-master.rfc-fidelity.md` or `reviews/00-master.executability.md`.
- The review task.

If the reviewer writes its review artifact directly, use it. If it returns findings instead, write the review artifact from the response. Check whether the expected file exists after the reviewer finishes.

## Execution Bindings

Claude execution bindings are file-defined worker agents under the agent directories above. They are the reliable mechanism for controlling model selection and preloaded skills.

Reuse first. A worker matches when its `model` and `skills` cover the sub-plan's assigned model tier and required skills. If no binding exists, create one using the same worker rules as direct feature planning:

- Place project-specific workers in `.claude/agents/`.
- Name implementers `{model-tier}-{domain}-worker.md`.
- Name the shared test author `{model-tier}-test-author-worker.md`.
- Set `model` explicitly.
- List preloaded `skills` in frontmatter.

Newly created Claude worker agents are not discoverable to an already-running session. After creating or modifying a worker, tell the user the current session must be restarted before dispatch.

## Review Artifact Ownership

The planner owns the `reviews/` directory and must verify each expected review artifact exists after reviewer completion.
