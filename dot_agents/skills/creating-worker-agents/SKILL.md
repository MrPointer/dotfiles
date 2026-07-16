---
name: creating-worker-agents
description: Create, repair, or update implementation worker agents and runtime bindings. Use whenever the user asks to add a project-local or global implementation worker or custom subagent, a repository lacks a worker for a model tier or technical domain, workers should follow a consistent template, or an existing worker has the wrong model, skills, permissions, or discoverability. Supports OpenCode, Claude Code, and Codex. Do not use for reviewer agents.
---

# Creating Worker Agents

Create thin, reusable implementation workers that encode runtime identity, concrete model, skills, and permissions.

## Scope

Use this skill for implementation workers only. Use the appropriate reviewer-development workflow for review agents.

## Runtime Binding

Determine the active runtime from the system prompt and environment, then read exactly one adapter:

- OpenCode: [references/runtime-opencode.md](references/runtime-opencode.md)
- Claude Code: [references/runtime-claude.md](references/runtime-claude.md)
- Codex: [references/runtime-codex.md](references/runtime-codex.md)

If the runtime is ambiguous, ask rather than guessing. The adapter owns syntax, placement, model identifiers, skill loading, permissions, and discovery behavior. This file owns the common creation workflow.

## Required Inputs

Derive these values from the user's request, existing worker conventions, and the repository's documented skill mappings:

| Input | Requirement |
|---|---|
| Semantic tier | Exactly one of `Cheapest`, `Mid-tier`, or `Most capable` |
| Domain | A short stable label derived from the work and required skills |
| Required skills | Exact discoverable skill IDs; never descriptions or guessed aliases |
| Required capabilities | Normally `read`, `edit`, and `shell`; ask before omitting one because implementation verification commonly needs all three |
| Worker purpose | The bounded implementation domain the worker should cover |

Do not create a worker if the tier, exact skills, or required capabilities are unresolved. Ask one focused question rather than guessing.

## Global Model Defaults

Use the active runtime adapter's defaults when the project has no explicit tier mapping. A project-local explicit mapping wins. Existing workers are evidence only when their tier is explicit; never infer a tier from model reputation or observed output.

Do not silently substitute a different model when a default is unavailable. Report the unavailable binding and ask the user to choose a replacement.

## Workflow

### 1. Read Project Context

Read the repository's agent instructions and relevant documentation before inspecting worker files. Use documented skill mappings and existing workers to understand the project's conventions, not to infer requirements the user did not request.

### 2. Discover Existing Bindings

Use the active runtime adapter to inspect project-local workers before global workers.

A worker is eligible only when all of these match:

- its purpose covers the implementation domain;
- its explicit semantic tier matches the requested tier;
- its concrete model source is unambiguous;
- it can load every exact required skill;
- its effective permissions cover every required capability without unsafe unrelated grants;
- it is invokable through the runtime's native worker mechanism.

If exactly one eligible worker exists, reuse it. If multiple eligible workers exist, ask the user to select one. Do not choose by filename order.

### 3. Decide Whether To Repair Or Create

Repair an existing project-local worker only when its intended domain and tier already match and the change preserves its other known consumers. Otherwise create a new project-local worker. Do not modify a global worker to satisfy one repository unless the user explicitly requests a global change.

Use `{semantic-tier-slug}-{domain}-worker` names with these semantic tier slugs:

- `cheapest`
- `mid-tier`
- `most-capable`

Keep provider and model identifiers out of the worker name. The explicit model belongs inside the runtime definition so the implementation can change without renaming the worker.

Keep the domain short and skill-oriented, such as `frontend`, `backend`, `go`, `zsh`, `docs`, or `database`.

### 4. Render The Runtime Binding

Follow the selected adapter's template so the new or repaired worker satisfies every eligibility criterion from step 2. Set an explicit model and, when the runtime supports it, reasoning effort. Append the canonical prompt body from [assets/shared-worker-prompt.md](assets/shared-worker-prompt.md) after the runtime-specific role and skill-loading instructions.

### 5. Validate Discovery

Validate the worker file's syntax and confirm its referenced skills exist. Use the runtime's read-only discovery mechanism when one exists.

If the runtime does not hot-reload worker definitions, tell the user exactly what must be restarted or reloaded. Creation is not complete until the worker is either discoverable or the restart requirement is clearly reported.

## Completion Report

Report:

- reused, repaired, or created binding identity;
- semantic tier and concrete model;
- exact skills and effective capabilities;
- definition path;
- validation performed;
- whether reload or restart is required; and
- any remaining reason the worker is not yet usable.
