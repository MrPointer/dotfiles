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

## Execution Dispatch

RFC-backed plans with two or more sub-plans must include concrete lead-agent instructions and worker tables in the master plan. During execution, launch the assigned Claude worker agents rather than recreating their persona in prompt text. Do not rely on prompt wording to pick the right model, and do not let the coordinator execute a sub-plan directly when the plan assigned a worker or model tier.

## TDD Isolation Mechanics

If any sub-plan has testable acceptance criteria, the shared test-author worker must be paired with an isolation mechanism. Prefer Worktrunk when available, then Claude's native worktree mechanism, then `git worktree`. If this cannot be verified, the plan must say that structural TDD is blocked or explicitly skipped with a concrete reason; generic "runtime cannot isolate" language is not sufficient when a worktree plus worker dispatch path is available.

## Model Assignment

- Use the worker agent's explicit `model` field as the source of truth.
- Treat the plan's model tier as binding. Do not silently upgrade or downgrade.
- If the requested model cannot be used in the current Claude environment, stop and ask the user how to proceed.

## Review Artifact Ownership

The planner owns the `reviews/` directory and must verify each expected review artifact exists after reviewer completion.
