# Master Plan Template

The master plan is the orchestration document. It does NOT contain implementation details — those live in sub-plans.

```markdown
# Master Plan: <Feature Name>

## Summary
<Brief description of what this feature accomplishes>

## Requirements
<Bullet list of confirmed requirements from Phase 1>

## Scope
- **In scope**: ...
- **Out of scope**: ...

## Sub-Plans

| #  | Sub-Plan                | Depends On | Model  | Description                          |
|----|-------------------------|------------|--------|--------------------------------------|
| 01 | `01-<name>.md`          | —          | Cheapest      | <What this sub-plan accomplishes>    |
| 02 | `02-<name>.md`          | 01         | Mid-tier      | <What this sub-plan accomplishes>    |
| 03 | `03-<name>.md`          | —          | Cheapest      | <What this sub-plan accomplishes>    |
...

## Execution Order
<Describe which sub-plans can run in parallel and which must be sequential>
- **Parallel group 1**: 01, 03 (no dependencies)
- **Sequential**: 02 (after 01)
...

## Execution via Worker Agents

**Worker agents are REQUIRED for plans with 2+ sub-plans.** Each sub-plan's model + skills combination maps to a worker agent definition (created during Phase 4). Worker agents are the only reliable mechanism for controlling sub-agent model selection — model requests via natural language prompts or team configuration are unreliable.

**The only exception** — skip worker agents when:
- Single sub-plan (just execute directly)
- All sub-plans are trivially small (e.g., "add one import")

**Worker Agents** (created during planning):
| Sub-Plan | Worker Agent | Model Tier | Skills |
|----------|-------------|------------|--------|
| 01       | `<tier>-<domain>-worker` | <tier> | <skills> |
| 02       | `<tier>-<domain>-worker` | <tier> | <skills> |
...

**File Ownership** (prevent conflicts during parallel execution):
| Sub-Plan | Primary Files |
|----------|---------------|
| 01       | <files this sub-plan creates/modifies> |
| 02       | <files this sub-plan creates/modifies> |
...

**Lead Agent Instructions**:
- Use this master plan as the roadmap
- Spawn each sub-plan's assigned worker agent from the table above
- Run sub-plans in the same parallel group concurrently where the framework supports it
- For sequential dependencies, wait for the prior worker to complete before spawning the next
- The lead relays information between workers when needed (workers cannot communicate directly)
- Pass the sub-plan file path and any prerequisite context when spawning
- Synthesize results when all sub-plans finish

**Coordination Points**:
<When the lead needs to relay information between sequential workers>
- After 01 completes: Pass results to worker executing 02
- If <event>: Relay to affected workers
...

## Risks & Mitigations
| Risk | Mitigation |
|------|------------|
| ...  | ...        |

## Post-Execution
If this project has component-level documentation, run the `component-docs-reviewer` agent to verify
component docs still match the actual implementation.
```
