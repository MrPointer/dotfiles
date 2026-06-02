# Plan Creation, Reviewers, And Bindings

Use this reference when creating plan files, assigning project-local reviewers, assigning execution models, establishing runtime bindings, or writing lead-agent execution instructions.

## Plan Location

Choose the plan directory from context:

- Standalone feature: `plans/features/<feature-name>/`
- Feature belonging to an epic: `plans/epics/<epic-name>/<feature-name>/`

If the user references an epic plan or provides epic context, use the epic path. Otherwise, default to standalone.

Create this structure:

```text
<plan-directory>/
├── 00-master.md
├── 01-<first-task>.md
├── 02-<second-task>.md
├── reviews/
└── ...
```

Use the templates listed in `SKILL.md`.

## Project-Local Reviewer Assignment

Use the active runtime adapter to discover project-local reviewer bindings. Every supported runtime must provide a way to represent and invoke project-local reviewer subagents, even if the storage format or dispatch mechanism differs.

Assign the most appropriate local reviewer to each sub-plan based on the reviewer's description and the sub-plan's domain. Write that reviewer into the sub-plan's `## Reviewer` field.

If no suitable local reviewer exists for a sub-plan domain, warn the user with a specific recommendation. Example: `sub-plan 03 covers database migrations, but no database reviewer exists locally`.

Ask how to proceed. Do not silently substitute a generic global reviewer for missing project-specific review coverage.

## Model Selection

Assign one execution model to each sub-plan. Evaluate top-down and use the first tier that fits.

### Most Capable Model

Use when the sub-plan involves:

- ambiguous or underspecified requirements that need interpretation
- multi-step reasoning across multiple systems or domains
- novel architectural approaches with no existing pattern to follow
- security-sensitive operations where mistakes are costly
- documentation updates

### Mid-Tier Model

Use when the sub-plan involves:

- following established patterns already present in the codebase
- CRUD operations, straightforward integrations, or configuration changes
- simple data transformations
- test writing for existing code with clear acceptance criteria
- complex business logic with multiple edge cases
- integration of multiple systems or packages
- performance-critical code requiring careful trade-offs
- state machines or error handling with recovery logic

### Cheapest Model

Use only for the smallest mechanical sub-plans where success does not require interpretation, pattern judgment, or non-trivial cross-file reasoning.

- file renames or moves with all paths and required reference updates explicitly listed
- tiny text-only mechanical edits with no behavior, schema, API, build, test, or workflow impact

Do not assign the cheapest model just because the pattern is established. Established to the planner may still be non-trivial to the cheapest model.

When in doubt, use at least the mid-tier model. The cost of a wrong model choice is rework, which is usually more expensive than the model difference. Document the recommendation in each sub-plan's `## Execution Model` field with a brief rationale.

## Execution Bindings

Each sub-plan's model + skills combination needs a matching runtime-specific execution binding when the plan has two or more sub-plans or any explicit model assignment. Natural language alone is not reliable for model selection or skill preload.

A single trivial sub-plan with no explicit model assignment may record direct execution as the binding decision.

Use the active runtime adapter for exact binding mechanics.

Search for reusable bindings first. A binding matches when it covers the required model tier and skill set. A partial match can be updated if appropriate.

If no match exists, establish one using the runtime adapter:

- follow runtime naming and placement conventions
- set the target model through the runtime's real model-selection mechanism
- make required skills explicit through the runtime's real skill-loading mechanism
- keep the binding minimal; the sub-plan provides task context, and the binding provides model, skills, permissions, and runtime-native identity

If the runtime uses ephemeral bindings, record parameters in the runtime-specific way so retries and resumed execution use the same model and skills.

## Test Author Binding

If any sub-plan has testable acceptance criteria, create or reuse one shared project-local test-author binding at the most capable model tier.

Preload the project's testing and code-writing skills, plus `test-driven-development` if available. All sub-plans with testable acceptance criteria share this binding.

If newly established persistent bindings require discovery, reload, or session restart before availability, warn the user when presenting the plan.

## Master Plan Execution Instructions

Multi-sub-plan plans must include explicit lead-agent execution mechanics. The executor should not infer them from this planning skill.

The master plan must include:

- worker-dispatch instructions and worker table
- coordination points and dependency graph
- file ownership table
- implementer worktree isolation requirements
- build/cache seeding requirements, or explicit `None`
- cross-sub-plan data flow table
- result integration and post-execution checks
- instruction that plan files, review files, and `progress.md` stay in the coordinator workspace
- instruction that workers receive inline task packets, not copied plan files

For concurrent work, require task-scoped implementer worktrees. If isolation or required cache seeding cannot be verified, instruct execution to serialize or ask the user.
