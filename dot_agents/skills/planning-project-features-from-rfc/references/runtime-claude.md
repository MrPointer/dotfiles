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

Launch each reviewer agent through its own `subagent_type` (for example, `subagent_type: "plan-rfc-fidelity-reviewer"`). The named agent's frontmatter declares both its persona and its `tools` list — including `Write`/`Edit` for write-capable reviewers — so direct dispatch loads both the right system prompt and the right tools. Do not launch reviewers as `general-purpose`: that bypasses the reviewer's specialized system prompt and replaces its declared tool list with general-purpose's broader set.

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

RFC-backed plans with two or more sub-plans must include concrete lead-agent instructions, worker tables, implementer worktree isolation, and result-integration mechanics in the master plan. During execution, launch the assigned Claude worker agents rather than recreating their persona in prompt text. Do not rely on prompt wording to pick the right model, and do not let the coordinator execute a sub-plan directly when the plan assigned a worker or model tier.

## Implementer Worktree Mechanics

For sub-plans in the same parallel group, the master plan must require task-scoped implementer worktrees rather than concurrent workers in the coordinator workspace. It should reference the active execution adapter's Workspace Isolation Strategy instead of repeating the fallback chain. For build-heavy projects, it must name ignored build/cache directories that isolated worktrees need before dispatch, or explicitly state that no seeding is required. If no isolated implementer path or required seeding path can be verified, the plan must instruct the executor to serialize the group or ask the user.

The plan must keep plan files, review files, and `progress.md` coordinator-owned. Implementers receive inline task packets and prerequisite outputs, not plan paths copied into worker worktrees.

## TDD Isolation Mechanics

If any sub-plan has testable acceptance criteria, the shared test-author worker must be paired with an isolation mechanism from the active execution adapter's Workspace Isolation Strategy. For build-heavy projects, the plan must also name ignored build/cache directories the test author workspace needs before compiling or running tests. If this cannot be verified, the plan must say that structural TDD is blocked or explicitly skipped with a concrete reason; generic "runtime cannot isolate" language is not sufficient when a priority-order worktree plus worker dispatch path is available.

## Model Assignment

- Use the worker agent's explicit `model` field as the source of truth.
- Treat the plan's model tier as binding. Do not silently upgrade or downgrade.
- If the requested model cannot be used in the current Claude environment, stop and ask the user how to proceed.

## Review Artifact Ownership

The planner owns the `reviews/` directory and must verify each expected review artifact exists after reviewer completion.
