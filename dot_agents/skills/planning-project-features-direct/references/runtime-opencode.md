# Runtime Adapter: OpenCode

Use this adapter only when the active runtime is OpenCode.

This adapter maps the canonical planning workflow in `../SKILL.md` to OpenCode-native mechanics. It must not redefine phases, review policy, or plan semantics.

## Discovery

**Skill directories**:

- Project-local OpenCode skills: `.opencode/skills/` within the repository
- Global OpenCode skills: `~/.config/opencode/skills/`
- OpenCode also discovers compatible skills under `.claude/skills/`, `~/.claude/skills/`, `.agents/skills/`, and `~/.agents/skills/`

**Agent directories** (where reviewer and worker agent definitions live):

- Project-local OpenCode agents: `.opencode/agents/` within the repository
- Global OpenCode agents: `~/.config/opencode/agents/`

**Dotfile-source note**: If the project itself is a chezmoi source directory or another dotfile source repo, use the source-path equivalents when creating global artifacts (for example `private_dot_config/opencode/agents/` rather than editing `~/.config/opencode/agents/` directly).

**Precedence**: When the same agent name exists in both places, project-local wins. When searching for reusable reviewer or execution bindings, check project-local first, then global.

Also consult `AGENTS.md` for documented skill mappings, reviewer expectations, and project-specific conventions.

## Reviewer Bindings

OpenCode reviewer bindings should be **markdown-defined custom subagents**. Store them as files in the agent directories above so they are discoverable by name and invokable through the runtime's normal subagent mechanism.

**Launch mechanism**: Use OpenCode's subagent dispatch path for reviewers. In practice, that means invoking the reviewer by its subagent name rather than recreating its persona ad hoc in prompt text. The built-in `general` and `explore` subagents are not substitutes for required project-local reviewers.

**Dispatch parameters**: Pass the plan path, the requested review output path (for example `reviews/00-master.design.md`), and the review task. Keep the prompt focused on the requested review and avoid mixing in unrelated planner rationale.

**Output ownership**: If the reviewer has `edit` permission and is designed to write files directly, let it write its own review artifact. If it returns findings instead, the planner writes the review artifact from the response. In both cases, the planner must verify that the expected file exists after the reviewer finishes.

## Local Reviewer Requirement

- The planning workflow still requires project-local reviewer coverage for project-specific domains.
- Global reviewers remain part of the master-plan loop, but they do not replace local reviewers for sub-plan review.
- If a local reviewer is missing in OpenCode-native runnable form, report the gap and stop. Do not silently substitute the built-in `general` subagent or recreate the reviewer in prompt text.
- If the active OpenCode environment cannot invoke project-local custom subagents, report that runtime incompatibility and stop rather than weakening the review loop.

## Execution Bindings

OpenCode execution bindings should also be **markdown-defined custom subagents** under `.opencode/agents/` or `~/.config/opencode/agents/`. These files are the reliable place to encode subagent identity, model selection, permissions, and binding-specific instructions.

**Reuse first**: Search both agent directories. A binding is a match if it covers the sub-plan's required model tier, has the permissions needed for the work, and its prompt is compatible with the required skills and role.

**Creating missing bindings**:

- **Placement**: Prefer project-local `.opencode/agents/` when the binding depends on project-local skills or conventions.
- **Naming convention**: `{model-tier}-{domain}-worker.md` for implementers. Keep reviewer names domain-specific and descriptive.
- **Frontmatter rules**:
  - `description` is required and should make the role discoverable to the runtime.
  - `mode: subagent` so the binding is invokable as a subagent.
  - `model`: set this explicitly. In OpenCode, subagents inherit the caller's model when `model` is omitted, which is not strong enough for plan-assigned model tiers.
  - `permission`: grant the minimum tool access needed for the role. Prefer this over the deprecated `tools` field for new bindings.
  - `permission.skill`: explicitly allow the skills the binding must load.
  - `permission.task`: usually deny or tightly scope nested subagent dispatch unless the binding genuinely needs to spawn children.
- **Prompt/body**: Keep it thin. The sub-plan remains the source of task truth. Use the prompt to define the role and, when necessary, tell the subagent which named skills to load at the start.

**Important OpenCode difference**: OpenCode agent definitions do not have a declarative `skills` preload field. If a binding requires skills, make that explicit in two places:

1. The binding's prompt should instruct the subagent to load the named skills immediately via the `skill` tool.
2. The binding's `permission.skill` rules must allow those skill names.

Do not rely on prompt text alone to pick the right model, and do not rely on the subagent to discover required skills implicitly.

**Testing skills**: For testable implementation work, ensure the implementer binding can load the testing skills named by the sub-plan, whether they are project-local or global. Do not create a separate test-writing binding.

**Discovery warning**: OpenCode's public docs describe where agent files live, but they do not promise hot-reload semantics for newly created or renamed agent files. After establishing a new persistent binding, verify that the runtime can actually invoke it. If it cannot, tell the user that a session restart or reload is required before the new binding can be used.

## Review Loop

- Run independent reviewers in parallel using separate OpenCode subagent dispatches when practical.
- Re-run only affected reviewers during convergence; do not restart the full review.
- The planner remains responsible for checking that each expected review artifact was actually created.

## Execution Dispatch

Direct feature plans with two or more sub-plans must include concrete lead-agent instructions, worker tables, implementer worktree isolation, and result-integration mechanics in the master plan. During execution, OpenCode workers are invoked through the runtime's native subagent mechanism. In interactive sessions, `@<worker-name>` mention is a valid native invocation path when it can invoke the named worker and preserve the worker's configured model and permissions.

The CLI is an acceptable fallback when the current runtime surface cannot invoke project-local custom subagents directly, or when explicit workspace routing is required:

```bash
opencode run --agent <worker-name> --dir <workspace-path> "<task prompt>"
```

Do not rely on prompt text alone to pick the right model, and do not let the coordinator execute a sub-plan directly when the plan assigned a worker or model tier.

## Implementer Worktree Mechanics

Apply the master plan's Concurrency Policy before assigning same-group implementation. If the policy says `Linear DAG`, the master plan must instruct the executor to run one sub-plan at a time even when sub-plans are logically independent. Only when the policy says `Parallel allowed`, same-group sub-plans must require task-scoped implementer worktrees rather than concurrent workers in the coordinator workspace. The plan should reference the active execution adapter's Workspace Isolation Strategy instead of repeating the fallback chain. It must name ignored build/cache directories that isolated worktrees need before dispatch, or explicitly state that no seeding is required. If no isolated implementer path, required seeding path, or output/cache safety can be verified, the plan must instruct the executor to serialize the group or ask the user.

The plan must keep plan files, review files, and `progress.md` coordinator-owned. Implementers receive inline task packets and prerequisite outputs, not plan paths copied into worker worktrees.

## Testable Work Mechanics

For testable behavior changes, the master plan should tell execution that the implementer owns both tests and code. The implementer should follow the sub-plan's testing skills, make a test-first attempt when practical, and report tests added or updated plus verification results. If test-first work or new tests are not practical, the implementer reports the reason instead of relying on a separate test-writing worker.

## Model Assignment

- Use the custom subagent's explicit `model` field as the source of truth.
- Treat the plan's model tier as binding. Do not silently upgrade or downgrade.
- Because omitted `model` values inherit from the caller in OpenCode, do not treat inherited model selection as sufficient for a binding that is supposed to encode a plan decision.
- If the requested model cannot be used in the current OpenCode environment, stop and ask the user how to proceed.

## Review Artifact Ownership

- The planner owns the `reviews/` directory and is responsible for ensuring the expected files exist.
- After each reviewer finishes, verify whether the requested review artifact was created.
- If not, write it from the returned findings before continuing.
