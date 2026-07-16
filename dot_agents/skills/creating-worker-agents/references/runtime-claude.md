# Runtime Adapter: Claude Code

Use this adapter only when Claude Code is the active runtime.

## Discovery And Placement

Search project-local `.claude/agents/` before global `~/.claude/agents/`. Use `.claude/agents/` for project-local workers and `~/.claude/agents/` for explicitly requested global workers.

## Global Model Defaults

| Semantic tier | Claude model alias |
|---|---|
| `Cheapest` | `haiku` |
| `Mid-tier` | `sonnet` |
| `Most capable` | `opus` |

An explicit project mapping overrides this table. If the selected model is unavailable, stop and ask rather than changing tiers.

## Worker Template

Create `.claude/agents/{semantic-tier-slug}-{domain}-worker.md`:

```markdown
---
name: {semantic-tier-slug}-{domain}-worker
description: "Use this {domain} implementation worker for work assigned to the {semantic-tier} model tier. It is intended for {bounded-purpose}."
model: {model-alias}
skills:
  - {skill-id}
---

You are a {domain} implementation worker configured for the {semantic-tier} tier.

{shared-worker-prompt-body}
```

Replace `{shared-worker-prompt-body}` with the exact prompt body from [shared-worker-prompt.md](../assets/shared-worker-prompt.md). Claude preloads every exact required skill through `skills`; do not depend on the worker discovering them from prose.

## Capability Mapping

Claude workers have full tool access when `tools` is omitted. Omit `tools` unless the user or repository conventions explicitly require restrictions.

Map capabilities to the smallest native tools that satisfy them:

- `read`: read and repository-search tools;
- `edit`: file-writing and editing tools;
- `shell`: shell execution.

## Validation And Discovery

- Confirm valid YAML frontmatter, an explicit `model`, and exact skill IDs.
- Newly created or modified Claude worker files are not discoverable in an already-running session. Tell the user clearly that the session must restart before dispatch.
