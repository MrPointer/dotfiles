# Runtime Adapter: Claude

Use this adapter only when the active runtime is Claude.

This adapter maps the canonical planning workflow in `../SKILL.md` to Claude-native mechanics. It must not redefine phases, review policy, or plan semantics.

## Discovery

**Skill directories**:

- Global skills: `~/.claude/skills/`
- Project-local skills: `.claude/skills/` within the repository

**Agent directories** (where reviewer and worker agent definitions live):

- Global agents: `~/.claude/agents/`
- Project-local agents: `.claude/agents/` within the repository

**Precedence**: When the same name exists in both, project-local wins for loading. When searching for a reusable binding, check project-local first, then global.

Also consult `AGENTS.md` (or `CLAUDE.md`) for documented skill mappings and domain-specific conventions. Prefer existing file-defined reviewer and worker agents when they match the required model tier and skill set.

## Reviewer Bindings

**Launch mechanism**: Launch each reviewer agent through its own `subagent_type` (for example, `subagent_type: "plan-rfc-fidelity-reviewer"`). The named agent's frontmatter declares both its persona and its `tools` list — including `Write`/`Edit` for write-capable reviewers — so direct dispatch loads both the right system prompt and the right tools. Do not launch reviewers as `general-purpose`: that bypasses the reviewer's specialized system prompt and replaces its declared tool list with general-purpose's broader set, defeating the reviewer definition. Read-only agent types like `Explore` are not a substitute either; pick the actual reviewer's `subagent_type`.

**Dispatch parameters**: Pass the plan directory path and the requested review output path (e.g., `reviews/00-master.architect.md`).

**Output ownership**: If a reviewer definition can write files directly, let it write its own review artifact. If it returns findings instead, the planner writes the review artifact from the returned response. The planner checks whether the expected file exists after the reviewer finishes.

## Execution Bindings

Claude execution bindings are **file-defined worker agents** under the agent directories above. They are the only reliable mechanism for controlling a sub-agent's model and preloaded skills — requesting a model or skill via natural-language prompt or team configuration is unreliable.

**Reuse first**: Search both the global and project-local agent directories. A worker is a match if its `model` field and `skills` list cover the sub-plan's required model tier and skill set. A partial match (correct model but incomplete skills, or correct skills but different model) can serve as a basis for an updated worker — adapt rather than starting from scratch.

**Creating missing workers**:

- **Placement**: Always project-local (`.claude/agents/`), since the preloaded skills are usually project-local.
- **Naming convention**: `{model-tier}-{domain}-worker.md` — e.g., `mid-go-worker.md`, `cheap-docs-worker.md`. The model tier must match the decision-tree terminology (`cheap`, `mid`, `capable`); domain should reflect the preloaded skill set.
- **Frontmatter rules** (these are load-bearing):
  - `description` **must** be a quoted single-line YAML string with `\n` escapes for any line breaks. Multiline YAML descriptions break frontmatter parsing and the agent won't be discovered.
  - `model`: set to the target model identifier — this is what actually controls the execution model.
  - `skills`: list of skill names to preload into the agent's context at startup. The agent does not need to load skills manually.
  - Do **not** set `tools` unless you need to restrict access; omit for full tool access.
- **System prompt**: Brief role description. The sub-plan provides all task context; the agent definition provides model, skills, and identity.

**Test author worker**: If any sub-plan has testable acceptance criteria, create a single project-local test author worker. It always uses the **most capable model** available in the environment. Preload the project's testing and code-writing skills, plus `test-driven-development` if available. Use the naming convention `{model-tier}-test-author-worker.md`. All sub-plans with testable AC share this single worker.

**Session restart warning (loud)**: Newly created worker agent files are **not discoverable** to a session that was already running when they were created. After creating or modifying a worker agent, tell the user explicitly that the current session must be restarted before the worker can be dispatched. This is a frequent gotcha — do not bury it.

## Review Loop

- Run independent reviewers in parallel using multiple sub-agent dispatches in a single message.
- Re-run only affected reviewers during convergence; do not restart the full review.
- The planner is responsible for checking that each expected review artifact was actually created, regardless of whether the reviewer wrote it directly or returned findings.

## Execution Dispatch

Direct feature plans with two or more sub-plans must include concrete lead-agent instructions, worker tables, implementer worktree isolation, and result-integration mechanics in the master plan. During execution, launch the assigned Claude worker agents rather than recreating their persona in prompt text. Do not rely on prompt wording to pick the right model, and do not let the coordinator execute a sub-plan directly when the plan assigned a worker or model tier.

## Implementer Worktree Mechanics

For sub-plans in the same parallel group, the master plan must require task-scoped implementer worktrees rather than concurrent workers in the coordinator workspace. It should reference the active execution adapter's Workspace Isolation Strategy instead of repeating the fallback chain. If no isolated implementer path can be verified, the plan must instruct the executor to serialize the group or ask the user.

The plan must keep plan files, review files, and `progress.md` coordinator-owned. Implementers receive inline task packets and prerequisite outputs, not plan paths copied into worker worktrees.

## TDD Isolation Mechanics

If any sub-plan has testable acceptance criteria, the shared test-author worker must be paired with an isolation mechanism from the active execution adapter's Workspace Isolation Strategy. If this cannot be verified, the plan must say that structural TDD is blocked or explicitly skipped with a concrete reason; generic "runtime cannot isolate" language is not sufficient when a priority-order worktree plus worker dispatch path is available.

## Model Assignment

- Use the worker or reviewer definition's explicit `model` field as the source of truth.
- Treat the plan's model tier as binding. Do not silently upgrade or downgrade.
- If the assigned model cannot be used (permission mode, tool access, model unavailable), stop and ask the user how to proceed. Do not fall back to a more expensive model.
