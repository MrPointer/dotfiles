# Runtime Adapter: OpenCode

Use this adapter only when the active runtime is OpenCode.

This adapter maps `planning-project-features-from-rfc` to OpenCode-native mechanics. It must not redefine phases, routing, review policy, or plan semantics.

## Discovery

**Skill directories**:

- Project-local OpenCode skills: `.opencode/skills/`
- Global OpenCode skills: `~/.config/opencode/skills/`
- OpenCode also discovers compatible skills under `.claude/skills/`, `~/.claude/skills/`, `.agents/skills/`, and `~/.agents/skills/`

**Agent directories**:

- Project-local OpenCode agents: `.opencode/agents/`
- Global OpenCode agents: `~/.config/opencode/agents/`

If working from a dotfile source repository, use source-path equivalents when creating global artifacts.

## Required Reviewers

RFC-backed feature planning uses exactly these reviewers:

- `plan-rfc-fidelity-reviewer`
- `plan-executability-reviewer`

Do not launch `plan-architect-reviewer`, `plan-risk-reviewer`, `plan-clarity-reviewer`, or project-local domain reviewers in this workflow unless the user explicitly exits RFC-backed planning and returns to direct planning.

## Reviewer Bindings

OpenCode reviewer bindings are markdown-defined custom subagents. Discover project-local first, then global. Invoke reviewers by subagent name through OpenCode's normal subagent mechanism; do not recreate reviewer personas in prompt text.

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

After establishing a new persistent binding, verify that OpenCode can invoke it. If not, tell the user a reload or session restart is required.

## Review Artifact Ownership

The planner owns the `reviews/` directory and must verify each expected review artifact exists after reviewer completion.
