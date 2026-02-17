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
| 01 | `01-<name>.md`          | —          | Haiku  | <What this sub-plan accomplishes>    |
| 02 | `02-<name>.md`          | 01         | Sonnet | <What this sub-plan accomplishes>    |
| 03 | `03-<name>.md`          | —          | Haiku  | <What this sub-plan accomplishes>    |
...

## Execution Order
<Describe which sub-plans can run in parallel and which must be sequential>
- **Parallel group 1**: 01, 03 (no dependencies)
- **Sequential**: 02 (after 01)
...

## Team Execution (Agent Teams)

**Agent Teams are REQUIRED for plans with 2+ sub-plans.** Do not use Task sub-agents — they cannot write files and consume the main context window.

**The only exception** — skip Agent Teams when:
- ❌ Single sub-plan (just execute directly)
- ❌ All sub-plans are trivially small (e.g., "add one import")

**Setup**:
1. Enable Agent Teams: `export CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1`
2. Create the team:
   ```
   Create an agent team to execute .claude/plans/<feature-name>/00-master.md
   ```

**Team Lead Instructions**:
- Use this master plan as the roadmap
- Assign sub-plans to teammates based on the dependency graph
- Each teammate should load the Required Skills listed in their assigned sub-plan before starting
- Use the recommended models from the Sub-Plans table above
- Coordinate handoffs when dependencies complete
- Synthesize results when all sub-plans finish

**Suggested Team Structure**:
```
Create a team with <N> teammates to execute .claude/plans/<feature-name>/00-master.md:
- Teammate 1: Execute 01-<name>.md using Haiku (load skills: <skill-list>)
- Teammate 2: Execute 02-<name>.md using Sonnet (load skills: <skill-list>, requires 01 complete)
- Teammate 3: Execute 03-<name>.md using Haiku (load skills: <skill-list>, can start immediately)
...
```

**File Ownership** (prevent conflicts):
| Sub-Plan | Primary Files |
|----------|---------------|
| 01       | <files this sub-plan creates/modifies> |
| 02       | <files this sub-plan creates/modifies> |
...

**Communication Points**:
<When teammates might need to coordinate>
- After 01 completes: Notify teammate 2 that dependencies are ready
- If <event>: Broadcast to all teammates about <change>
...

## Risks & Mitigations
| Risk | Mitigation |
|------|------------|
| ...  | ...        |

## Post-Execution
If this project has component-level documentation, run the `component-docs-reviewer` agent to verify
component docs still match the actual implementation.
```
