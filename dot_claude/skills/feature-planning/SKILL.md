---
name: feature-planning
description: "Create implementation plans for features with rigorous requirement gathering. Use when planning new features, refactoring efforts, or any multi-step implementation. Never assumes or fills in gaps - always asks for clarification until requirements are complete."
---

# Feature Planning

Create thorough, actionable implementation plans by first ensuring all requirements are fully understood. **Never assume or guess** - ask until every gap is filled.

## Core Principles

1. **No Assumptions**: If something is unclear, ambiguous, or missing - ask. Do not fill in gaps with reasonable defaults or best guesses.
2. **Relentless Clarification**: Ask as many questions as needed. A plan built on assumptions is worse than no plan.
3. **User-Controlled Output**: Plans go to `.claude/plans/<meaningful-name>.md` unless the user specifies otherwise.
4. **Iterative Refinement**: The user reviews and approves before the plan is finalized.

## Workflow

### Phase 1: Requirement Gathering

Before writing any plan, gather complete information:

1. **Understand the goal**: What problem does this solve? What does success look like?
2. **Identify scope**: What's in scope? What's explicitly out of scope?
3. **Map dependencies**: What existing code/systems does this touch?
4. **Clarify constraints**: Performance requirements? Backward compatibility? Tech stack restrictions?
5. **Define acceptance criteria**: How will we know it's done correctly?

**Ask questions until you can answer all of the above confidently.**

If the user says "just figure it out" or "use your judgment":
1. **First, educate**: Present the options you see, explain trade-offs, and help them make an informed decision
2. **If they persist**: Respect their choice to "vibe-code" and proceed with your best judgment, but document your assumptions clearly in the plan

### Phase 2: Codebase Exploration

Once requirements are clear:

1. Search for relevant existing code, patterns, and conventions
2. Identify files that will need modification
3. Note any architectural constraints or patterns to follow
4. Flag potential conflicts or risks

Share findings with the user and confirm understanding before proceeding.

### Phase 3: Plan Creation

Only after Phases 1-2 are complete:

1. Generate a meaningful, descriptive plan filename based on the feature being planned
2. Write a structured plan with:
   - Summary of requirements (as confirmed with user)
   - Files to modify/create
   - Step-by-step implementation tasks
   - Risks and mitigation strategies
   - Open questions (if any remain)

### Phase 4: Review and Approval

Present the plan for user review. Incorporate feedback. Only mark as ready when the user explicitly approves.

## Plan File Structure

```markdown
# Plan: <Feature Name>

## Summary
<Brief description of what this plan accomplishes>

## Requirements
<Bullet list of confirmed requirements from Phase 1>

## Scope
- **In scope**: ...
- **Out of scope**: ...

## Files Affected
<List of files to modify/create with brief description of changes>

## Implementation Steps
1. <Step with clear deliverable>
2. <Step with clear deliverable>
...

## Risks & Mitigations
| Risk | Mitigation |
|------|------------|
| ... | ... |

## Acceptance Criteria
- [ ] <Criterion 1>
- [ ] <Criterion 2>
...
```

## Rules (Non-Negotiable)

- **Never write a plan based on incomplete information**
- **Never invent requirements the user didn't specify**
- **Save plans to the repository's local `.claude/plans/<meaningful-name>.md`** - not to the global `~/.claude/` directory, and never with random/generated filenames
- **Ask for clarification even if it feels repetitive** - it's better than introducing garbage
