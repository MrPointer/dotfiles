# Runtime Adapter: Codex

Use this adapter only when Codex is the active runtime.

## Discovery And Placement

Search project-local `.codex/agents/*.toml` before global `~/.codex/agents/*.toml`. Use `.codex/agents/` for project-local workers and `~/.codex/agents/` for explicitly requested global workers.

Codex discovers TOML agent roles from these directories and applies each file as a role-specific configuration layer. Create a persistent custom agent rather than an ephemeral dispatch recipe.

## Global Model Defaults

| Semantic tier | Model | Reasoning effort |
|---|---|---|
| `Cheapest` | `gpt-5.5` | `low` |
| `Mid-tier` | `gpt-5.5` | `medium` |
| `Most capable` | `gpt-5.6-sol` | `high` |

An explicit project mapping overrides this table. Confirm the selected model and effort appear in the active Codex model catalog. If unavailable, stop and ask rather than changing tiers.

## Worker Template

Create `{semantic-tier-slug}-{domain}-worker.toml` in the selected agent directory:

```toml
name = "{semantic-tier-slug}-{domain}-worker"
description = "Use this {domain} implementation worker for work assigned to the {semantic-tier} model tier. It is intended for {bounded-purpose}."
model = "{model}"
model_reasoning_effort = "{reasoning-effort}"
sandbox_mode = "workspace-write"
developer_instructions = """
You are a {domain} implementation worker configured for the {semantic-tier} tier.

Use these required skills before working: {codex-skill-references}.

{shared-worker-prompt-body}
"""

[skills]
include_instructions = true

[[skills.config]]
name = "{required-skill-id}"
enabled = true
```

Add one `[[skills.config]]` table for every required skill. Use the exact skill ID as the `name` selector. In `developer_instructions`, render `{codex-skill-references}` using Codex skill references such as `$writing-go-code` so the worker applies those enabled skills rather than merely knowing they exist.

Replace `{shared-worker-prompt-body}` with the exact prompt body from [shared-worker-prompt.md](../assets/shared-worker-prompt.md).

Keep `sandbox_mode = "workspace-write"` for implementation workers so they can read, edit, run commands, build, and test within Codex's configured sandbox. Do not replace declarative fields with dispatch-time prose.

## Validation

- Parse the TOML and reject unknown or malformed fields.
- Confirm `name`, `description`, `model`, `model_reasoning_effort`, `sandbox_mode`, and `developer_instructions` are present.
- Confirm every `skills.config.name` resolves to one intended skill and is enabled.
- Confirm the selected model supports the configured reasoning effort in the active model catalog.
- Confirm Codex discovers the named role from the selected project-local or global agent directory. If the current session does not reload agent files, tell the user to restart it.
