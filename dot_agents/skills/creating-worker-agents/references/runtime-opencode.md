# Runtime Adapter: OpenCode

Use this adapter only when OpenCode is the active runtime.

## Discovery And Placement

Search project-local `.opencode/agents/` before global `~/.config/opencode/agents/`. Use `.opencode/agents/` for project-local workers and `~/.config/opencode/agents/` for explicitly requested global workers.

In a chezmoi source repository, use the source-path equivalent when the requested artifact is global. Do not edit the applied target under `~` directly.

## Global Model Defaults

| Semantic tier | Model | Reasoning effort |
|---|---|---|
| `Cheapest` | `openai/gpt-5.5` | `low` |
| `Mid-tier` | `openai/gpt-5.5` | `medium` |
| `Most capable` | `openai/gpt-5.6-sol` | `high` |

An explicit project mapping overrides this table. If the selected model is unavailable, stop and ask rather than changing tiers.

## Worker Template

Create `.opencode/agents/{semantic-tier-slug}-{domain}-worker.md`:

```markdown
---
name: {semantic-tier-slug}-{domain}-worker
description: "Use this {domain} implementation worker for work assigned to the {semantic-tier} model tier. It is intended for {bounded-purpose}."
mode: subagent
model: {model}
reasoningEffort: {reasoning-effort}
permission:
  edit: allow
  bash: allow
  task:
    "*": deny
  skill:
    {skill-id}: allow
---

You are a {domain} implementation worker configured for the {semantic-tier} tier.

Load {exact-skill-list} immediately before working.

{shared-worker-prompt-body}
```

Replace `{shared-worker-prompt-body}` with the exact prompt body from [shared-worker-prompt.md](../assets/shared-worker-prompt.md). Adapt role wording to the domain. List every required skill under both `permission.skill` and the prompt's immediate-load instruction because OpenCode does not preload skills declaratively.

## Capability Mapping

| Planned capability | OpenCode binding |
|---|---|
| `read` | Native read access; do not disable repository reads |
| `edit` | `permission.edit: allow` |
| `shell` | `permission.bash: allow` |

Keep `permission.task."*": deny` for normal implementation workers.

## Validation

- Confirm frontmatter parses, every skill is discoverable, and the worker has `mode: subagent` plus an explicit model.
- OpenCode does not guarantee hot reload of new agent files. If the worker is absent from native dispatch after creation, tell the user to restart the session.
