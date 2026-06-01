# Runtime Adapter: OpenCode

Use this adapter only when the active runtime is OpenCode.

This adapter maps `planning-project-features-from-rfc` to OpenCode-native mechanics. It must not redefine phases, routing, review policy, or plan semantics.

## Discovery

**Skill directories**:

- Project-local OpenCode skills: `.opencode/skills/` within the repository
- Global OpenCode skills: `~/.config/opencode/skills/`
- OpenCode also discovers compatible skills under `.claude/skills/`, `~/.claude/skills/`, `.agents/skills/`, and `~/.agents/skills/`

**Agent directories**:

- Project-local OpenCode agents: `.opencode/agents/` within the repository
- Global OpenCode agents: `~/.config/opencode/agents/`

If working from a dotfile source repository, use source-path equivalents when creating global artifacts. For example, use `private_dot_config/opencode/agents/` rather than editing `~/.config/opencode/agents/` directly.

When the same agent name exists in both locations, project-local wins. When searching for reusable reviewer or execution bindings, check project-local first, then global.

Also consult `AGENTS.md` for documented skill mappings, reviewer expectations, and project-specific conventions.

## Required Reviewers

RFC-backed feature planning uses exactly these reviewers:

- `plan-rfc-fidelity-reviewer`
- `plan-executability-reviewer`

Do not launch `plan-architect-reviewer`, `plan-risk-reviewer`, `plan-clarity-reviewer`, or project-local domain reviewers in this workflow unless the user explicitly exits RFC-backed planning and returns to direct planning.

## Reviewer Bindings

OpenCode reviewer bindings are markdown-defined custom subagents. Discover project-local first, then global. Invoke reviewers by subagent name through OpenCode's normal subagent mechanism, including `@<reviewer-name>` mentions in interactive sessions when appropriate. Do not recreate reviewer personas in prompt text.

Pass only:

- Plan directory path.
- RFC path.
- Requested review output path, such as `reviews/00-master.rfc-fidelity.md` or `reviews/00-master.executability.md`.
- The review task.

If a reviewer has `edit` permission and writes the review file directly, let it do so. If it returns findings instead, write the review file from the response. Verify the expected review artifact exists before continuing.

## Execution Bindings

OpenCode execution bindings are markdown-defined custom subagents under `.opencode/agents/` or `~/.config/opencode/agents/`.

Reuse first. A binding matches when it covers the sub-plan's model tier, permissions, and required skills. If no binding exists, create one using the same worker rules as direct feature planning:

- Name implementers `{model-tier}-{domain}-worker.md`.
- Name the shared test author `{model-tier}-test-author-worker.md`.
- Set `mode: subagent`.
- Set an explicit `model`; omitted models inherit from the caller and are not acceptable for plan-assigned model tiers.
- Grant minimum permissions.
- Allow required skills in `permission.skill` and instruct the subagent prompt to load them immediately.
- Deny or tightly scope nested `permission.task` unless the worker genuinely needs to spawn child agents.

After establishing a new persistent binding, verify that OpenCode can invoke it. If not, tell the user a reload or session restart is required.

## Execution Dispatch

RFC-backed plans with two or more sub-plans must include concrete lead-agent instructions, worker tables, implementer worktree isolation, and result-integration mechanics in the master plan. During execution, OpenCode workers are invoked through the runtime's native subagent mechanism. In interactive sessions, `@<worker-name>` mention is a valid native invocation path when it can invoke the named worker and preserve the worker's configured model and permissions.

The CLI is an acceptable fallback when the current runtime surface cannot invoke project-local custom subagents directly, or when explicit workspace routing is required:

```bash
opencode run --agent <worker-name> --dir <workspace-path> "<task prompt>"
```

Do not rely on prompt text alone to pick the right model, and do not let the coordinator execute a sub-plan directly when the plan assigned a worker or model tier.

## Implementer Worktree Mechanics

For sub-plans in the same parallel group, the master plan must require task-scoped implementer worktrees rather than concurrent workers in the coordinator workspace. It should reference the active execution adapter's Workspace Isolation Strategy instead of repeating the fallback chain. For build-heavy projects, it must name ignored build/cache directories that isolated worktrees need before dispatch, or explicitly state that no seeding is required. If no isolated implementer path or required seeding path can be verified, the plan must instruct the executor to serialize the group or ask the user.

The plan must keep plan files, review files, and `progress.md` coordinator-owned. Implementers receive inline task packets and prerequisite outputs, not plan paths copied into worker worktrees.

## TDD Isolation Mechanics

If any sub-plan has testable acceptance criteria, the test-author binding must be paired with an isolation mechanism from the active execution adapter's Workspace Isolation Strategy. Same-workspace `@<test-author-worker>` invocation is not sufficient for structural TDD unless the runtime can prove it routes that subagent into the isolated worktree. Acceptable routing includes a verified native isolated-workspace dispatch mechanism or `opencode run --agent <test-author-worker> --dir <isolated-workspace>`. For build-heavy projects, the plan must also name ignored build/cache directories the test author workspace needs before compiling or running tests. If this cannot be verified, the plan must say that structural TDD is blocked or explicitly skipped with a concrete reason; generic "runtime cannot isolate" language is not sufficient when a priority-order worktree plus either native isolated-workspace dispatch or `opencode run --dir` is available.

## Model Assignment

- Use the custom subagent's explicit `model` field as the source of truth.
- Treat the plan's model tier as binding. Do not silently upgrade or downgrade.
- Because omitted `model` values inherit from the caller in OpenCode, do not treat inherited model selection as sufficient for a binding that is supposed to encode a plan decision.
- If the requested model cannot be used in the current OpenCode environment, stop and ask the user how to proceed.

## Review Artifact Ownership

The planner owns the `reviews/` directory and must verify each expected review artifact exists after reviewer completion.
