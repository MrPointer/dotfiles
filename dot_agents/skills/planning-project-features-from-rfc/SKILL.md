---
name: planning-project-features-from-rfc
description: Create implementation plans from a reviewed RFC. Uses the RFC as the approved design baseline, decomposes it into executable sub-plans, and runs RFC-specific plan review.
---

# Planning Project Features From RFC

Create thorough, actionable implementation plans for features within a single
project from a reviewed RFC. The RFC owns requirements, design decisions,
constraints, goals, non-goals, risks, and contracts. This skill owns the same
planning mechanics as direct feature planning: decomposition, dependency graph,
file ownership, acceptance criteria, required skills, model assignments,
execution bindings, and execution orchestration instructions.

The RFC replaces direct planning's open-ended requirement gathering and design
review. It does not replace the direct planning workflow's execution discipline.

## Runtime Binding

This skill has one canonical workflow. Runtime files only map that workflow to the active agent runtime's mechanics.

Before doing any work, determine the active runtime and read exactly one adapter:

- **OpenCode runtime** -> [references/runtime-opencode.md](references/runtime-opencode.md)
- **Codex runtime** -> [references/runtime-codex.md](references/runtime-codex.md)
- **Claude runtime** -> [references/runtime-claude.md](references/runtime-claude.md)

**Determining the active runtime**: Check the system prompt and environment
banner for identifying markers (e.g., "OpenCode", "Claude Code", "Codex
CLI"). If the signal is ambiguous, ask the user rather than guessing.

Do not load or mix instructions from another runtime adapter in the same turn.
If a runtime adapter conflicts with this file, this file is authoritative.

**Terminology bridge**: This skill uses runtime-neutral terms. Claude's runtime
calls execution bindings "worker agent definitions"; Codex's runtime calls them
"dispatch recipes"; OpenCode's runtime uses custom subagent definitions.
Reviewer bindings follow the same pattern. Use whichever term is native to the
active runtime when writing or reading concrete artifacts.

## Core Principles

1. **RFC Baseline**: Treat the reviewed RFC as the approved source of design decisions, constraints, goals, non-goals, risks, and contracts.
2. **No Design Re-Litigation**: Do not reopen RFC architecture, risk, or tradeoff decisions during planning.
3. **No Silent Deviations**: If planning requires changing the RFC design, stop and ask whether to revise the RFC or explicitly approve a plan deviation.
4. **Direct-Planning Mechanics Remain**: RFC-backed planning uses the direct workflow's decomposition, model selection, execution binding, worker-dispatch, documentation sub-plan, and approval mechanics unless this file explicitly replaces a step for RFC reasons.
5. **Atomic Decomposition**: Break work into the smallest self-contained sub-plans possible. Each sub-plan should be executable in isolation.
6. **Embedded Context**: Each sub-plan includes everything an executing agent needs. The agent should not have to read the RFC, master plan, or other sub-plans to understand its assigned work.
7. **RFC-Specific Review**: Reviewed RFC-backed plans use `plan-rfc-fidelity-reviewer` and `plan-executability-reviewer`. The RFC's own architecture, risk, and clarity reviews replace direct planning's full design review loop.

## Workflow

### Phase 1: Validate RFC Baseline

Use this workflow only after `planning-project-features` has routed here because
a reviewed RFC exists and the user wants an RFC-backed plan.

Read the RFC and confirm:

- `rfc-architect-reviewer` is `Passed` or `Passed with concerns`.
- `rfc-risk-reviewer` is `Passed` or `Passed with concerns`.
- `rfc-clarity-reviewer` is `Passed`, `Passed with concerns`, or explicitly not required by the project.
- No `Blocking` review status remains.
- Any `Passed with concerns` item is compatible with planning and does not require a design decision before decomposition.

If validation fails, stop and return to `planning-project-features` routing. Do not continue as RFC-backed planning.

Look for an active feature anchor using project conventions, then `docs/context/<topic>-anchor.md`.
If one exists, read it as supporting handoff context.
The anchor can explain intent, rejected alternatives, open questions, or teammate handoff context, but it does not
override the RFC. Any unresolved anchor question that conflicts with the RFC or affects decomposition, scope,
acceptance criteria, or contracts is a planning blocker to ask about.

### Phase 2: Planning-Focused Codebase Exploration

Read existing documentation first, then inspect code only for planning mechanics.
Existing documentation is dramatically cheaper than re-exploring code from scratch.

1. **Read existing docs first**: Check AGENTS.md for documentation pointers, then read relevant docs (domain, architecture, business processes, components).
2. **Explore code only for gaps**: Search for relevant code, patterns, and conventions that documentation and the RFC do not cover.
3. Confirm file paths and ownership boundaries needed for sub-plans.
4. Confirm existing interfaces, commands, schemas, config files, or runtime bindings named by the RFC.
5. Identify tests, packages, or verification scopes that acceptance criteria can reference.
6. **Identify required skills**: Determine which skills available in the active runtime the executing agent will need to follow project conventions correctly. Check the project's `AGENTS.md` for documented skill mappings, then use the active runtime adapter for exact discovery/loading mechanics.
7. **Flag documentation gaps**: If critical areas needed for the plan are undocumented, note them. Recommend the appropriate documenting skill:
   - Missing domain knowledge -> `documenting-domain`
   - Missing architecture overview -> `documenting-architecture`
   - Missing business workflow docs -> `documenting-business-processes`
   - Missing component docs -> `documenting-components`

   Present gaps to the user. They may want to create docs before planning continues, or accept the gap and proceed.
8. **Flag RFC/specification gaps**: When implementation planning depends on behavior the RFC does not define, treat it as a blocking gap. Do not fill it with a reasonable default. Ask whether to revise the RFC or approve an explicit plan deviation.

Do not use exploration to redesign the RFC. If exploration contradicts the RFC,
stop and ask whether to revise the RFC or abandon RFC-backed planning.

Share planning-relevant findings with the user and confirm understanding before
writing plan files when the findings affect decomposition, sequencing, or scope.

### Phase 3: Decomposition

This is the most critical phase. Break the RFC into sub-plans:

1. **Identify natural boundaries**: Look for seams in the work: different layers, domains, files, modules, or independently verifiable outcomes.
2. **Minimize dependencies and enforce DAG ordering**: Each sub-plan should depend on as few other sub-plans as possible. Where dependencies exist, make them explicit and one-directional. The dependency graph must form a valid DAG. No sub-plan may depend on information produced by a later sub-plan or a sub-plan in the same parallel group. Sub-plans cannot communicate with each other at runtime; the lead agent relays results strictly along dependency edges. If decomposition requires bidirectional information flow, merge or restructure the boundaries.
3. **Embed RFC context and cross-boundary contracts**: Each sub-plan must be self-contained. What belongs in a sub-plan:
   - RFC decisions, constraints, goals, non-goals, and risks that govern this unit.
   - Domain knowledge the agent cannot derive from code: business rules, config formats, protocol details, or accepted tradeoffs.
   - Cross-boundary contracts: exact interface/type signatures, files, data shapes, or commands that other sub-plans depend on. When sub-plans run in parallel, the consuming agent cannot discover these at execution time, so the plan must specify them. For sequential dependencies, the later agent can read the earlier sub-plan's actual output; no pre-specified contract is needed unless the earlier output constrains the later work.

   Cross-boundary contracts must satisfy three additional integrity rules:

   - **Caller annotations**: Every new public method/function introduced by a sub-plan must specify its production caller. If the caller lives in a different sub-plan, both sides must reference the contract: the producer documents "called by: sub-plan N, in `Location`", the consumer documents "calls: `Method` from sub-plan M". No orphan methods.
   - **Connected data flow**: When data must flow between components owned by different sub-plans, the master plan must trace the full path: source -> transport mechanism -> destination, with sub-plan ownership at each hop. Prose descriptions like "X stores the value on config" are insufficient when the consumer needs it delivered through a channel no sub-plan was told to wire.
   - **Interface boundary checks**: If a sub-plan adds a method to a concrete type, but consumers access that type through an interface, the plan must either add the method to the interface or explicitly assign the concrete-type wiring to a specific sub-plan.

   What does NOT belong: method bodies, private helper design, step-by-step coding instructions, exact commands for testing/linting/building, or design decisions already owned by skills.
4. **Preserve RFC non-goals**: Do not add work the RFC explicitly excluded. If implementation pressure suggests a non-goal should change, stop and ask whether to revise the RFC or approve an explicit deviation.
5. **Translate RFC risks into mechanics**: RFC risks become plan constraints, acceptance criteria, sequencing notes, rollback notes, or handoff constraints where relevant.
6. **Keep sub-plans small**: A good sub-plan should be completable in a single focused session. If it feels too big, split it further.
7. **Skills are the agent's authority, not the plan's**: List the skills each sub-plan requires, but do not replicate skill content into the plan. Skills define how to write code, how to test, how to lint, and how to build. The plan defines what to build and why. If a project has skills for testing, linting, or building, list them in Required Skills and write acceptance criteria at the behavioral/verification level rather than raw command level unless the RFC or project lacks an operational skill.
8. **Sub-plans must be decisive**: Before writing a design decision, verify it against the RFC and codebase. Do not write "if X exists" or "either A or B". Pick the RFC-approved approach or stop for clarification.
9. **Plan documentation updates as a sub-plan**: If the feature affects documented domain concepts, architecture, or business processes, add a final documentation sub-plan. See [Documentation Sub-Plan](#documentation-sub-plan).

Present the decomposition to the user for review before writing the actual plan
files.

### Phase 4: Plan Creation, Model Selection, And Execution Binding

Only after Phases 1-3 are complete:

1. **Create the plan directory and files**:

   **Plan location depends on context:**
   - Standalone feature: `plans/features/<feature-name>/`
   - Feature belonging to an epic: `plans/epics/<epic-name>/<feature-name>/`

   If the user references an epic plan or provides an epic context, use the epic path. Otherwise, default to standalone.

```text
<plan-directory>/
├── 00-master.md
├── 01-<first-task>.md
├── 02-<second-task>.md
├── reviews/
└── ...
```

2. **Assign execution models**: For each sub-plan, assess complexity and recommend an execution model. This enables cost optimization by using cheaper models for straightforward work while reserving the most capable model for planning, review, and difficult execution.

**Model Selection Decision Tree** — evaluate top-down, use the first tier that fits:

1. **Most capable model** — use when the sub-plan involves:
   - Ambiguous or underspecified requirements that need interpretation
   - Multi-step reasoning across multiple systems or domains
   - Novel architectural approaches with no existing pattern to follow
   - Security-sensitive operations where mistakes are costly
   - Documentation updates (see [Documentation Sub-Plan](#documentation-sub-plan))

2. **Mid-tier model** — use when the sub-plan involves:
   - Complex business logic with multiple edge cases
   - Integration of multiple systems or packages
   - Performance-critical code requiring careful trade-offs
   - State machines or error handling with recovery logic

3. **Cheapest model** — use when the sub-plan involves:
   - Following established patterns already present in the codebase
   - CRUD operations, straightforward integrations, configuration changes
   - File moves/renames, simple data transformations
   - Test writing for existing code with clear acceptance criteria

When in doubt, prefer one tier up. Document the recommendation in each
sub-plan's `## Execution Model` field with a brief rationale.

3. **Establish execution bindings for execution**: Each sub-plan's model + skills combination needs a matching runtime-specific execution binding. The active runtime adapter defines what that binding looks like. Natural language alone is not a reliable way to control model selection or preload the right skills.

   **Search for existing execution bindings**: Check the locations and mechanisms defined by the active runtime adapter. A binding is a match if it covers the sub-plan's required model tier and skill set. A partial match can serve as a basis for an updated binding.

   **Establish missing execution bindings**: If no matching binding exists, establish one using the mechanism defined by the active runtime adapter. If the runtime uses persistent bindings, create the required artifact. If the runtime uses ephemeral bindings, record the binding parameters in the runtime-specific way so retries and resumed execution use the same model and skills.

   - **Naming and placement**: Follow the active runtime adapter's conventions so the binding is discoverable.
   - **Model control**: Set the target model using the runtime's actual model-selection mechanism.
   - **Skill preload**: Make the required skills explicit using the runtime's actual skill-loading mechanism.
   - **Identity/prompt**: Keep the binding minimal. The sub-plan provides task context; the binding provides model, skills, and runtime-native agent identity.

   **Create a test author binding**: If any sub-plan has testable acceptance criteria, create or reuse one shared project-local test-author binding at the most capable model tier. Preload the project's testing and code-writing skills, plus `test-driven-development` if available. All sub-plans with testable AC share this single binding.

   **Warn the user**: If the runtime adapter says newly established persistent bindings require discovery, reload, or session restart before they become available, note that when presenting the plan.

4. **Write lead-agent instructions into the master plan**: Multi-sub-plan plans must include explicit worker-dispatch instructions, coordination points, file ownership, cross-sub-plan data flow, and post-execution checks. The executor should not have to infer these mechanics from the planning skill.

### Phase 5: RFC-Specific Review Loop

After plan creation, run an iterative review process before the user sees the plan.
RFC-backed planning uses exactly these reviewers by default:

- **`plan-rfc-fidelity-reviewer`**: Checks that the plan faithfully decomposes the RFC without contradicting it, omitting required context, or adding unapproved design scope.
- **`plan-executability-reviewer`**: Checks file ownership, acceptance criteria, dependency order, verification scope, worker-dispatch mechanics, model enforcement, and isolated execution mechanics.

Do not run `plan-architect-reviewer`, `plan-risk-reviewer`, or
`plan-clarity-reviewer` in this workflow unless the user explicitly exits
RFC-backed planning and returns to direct planning. Those design reviewers belong
to direct planning; RFC architecture, risk, and clarity were already reviewed by
RFC reviewers.

Pass the plan directory, RFC path, and review output path to each reviewer.
Review output lives in the plan's `reviews/` directory, for example:

```text
<plan-directory>/reviews/
├── 00-master.rfc-fidelity.md
└── 00-master.executability.md
```

If a reviewer has edit permission and writes its own review artifact, let it do
so. If it returns findings instead, write the review artifact from the response.
Verify the expected review artifact exists before continuing.

Incorporate findings and re-run only affected reviewers until both report no
blocking findings. If a finding requires changing the RFC design, stop and ask
whether to revise the RFC or approve an explicit deviation.

### Phase 6: User Approval And Feedback

Before presenting, run a final consistency check:

1. **RFC coverage** — every RFC goal and success criterion is represented in the plan.
2. **Non-goal preservation** — every RFC non-goal remains out of scope.
3. **Constraint propagation** — every RFC constraint that affects execution appears in the master plan or relevant sub-plan.
4. **Model assignments** — every sub-plan's execution model matches the decision tree and any explicit project rules.
5. **Execution bindings** — every multi-sub-plan worker assignment has a runtime-specific binding and lead-agent dispatch instructions.
6. **Skill conformance** — reviewer-driven changes do not introduce patterns that contradict required skills listed in each sub-plan.
7. **Cross-sub-plan prerequisites** — every interface, type, file, or artifact referenced in a sub-plan's Prerequisites section is created by an earlier sub-plan in the dependency graph.
8. **Integration contract integrity** — every "Produces" entry in one sub-plan has a matching "Consumes" entry in another when cross-sub-plan consumption exists, and the master plan's data flow table covers every cross-boundary data path.
9. **DAG validity** — execution order forms a valid DAG: no cycles, no sub-plan depends on output from a later or same-group sub-plan, and every consumed prerequisite points to an earlier group.
10. **Anchor boundaries** — if an active feature anchor exists, the plan does not duplicate anchor content and any sub-plan using `anchoring-context` has a feature-level reason to update it.
11. **Review status** — `plan-rfc-fidelity-reviewer` and `plan-executability-reviewer` findings are resolved or explicitly documented as non-blocking.

Fix inconsistencies before proceeding.

Present the reviewed plan with the RFC path, review summary, worker-dispatch
summary, and any explicit approved deviations. Only mark ready when the user
approves.

**Remind the user**: The plan intentionally omits implementation details. The
user reviews architecture and constraints now; they review actual code after
execution. This is by design, not a gap.

#### Handling User Feedback

When the user requests changes, incorporate them and classify each change to
determine whether re-review is needed:

| Change Type | Examples | Re-review Action |
|---|---|---|
| Cosmetic / wording | Clarify a step description, rename a sub-plan, fix typos | None |
| Scoped implementation detail | Add an edge case to one sub-plan, change a file path, adjust a step | None — planner judgment is sufficient |
| Scope adjustment within a sub-plan | Add/remove acceptance criteria, change approach for one sub-plan | Re-review affected sub-plan through `plan-rfc-fidelity-reviewer` if RFC scope may be affected, and `plan-executability-reviewer` if mechanics changed |
| Structural change | New sub-plan added, dependency graph changed, boundaries shifted, sub-plans merged/split, ownership or verification scope changed across sub-plans | Re-review affected plan sections with both RFC-fidelity and executability reviewers |
| RFC design change | Goal, non-goal, contract, risk, or accepted design decision changes | Stop and ask whether to revise the RFC or approve an explicit plan deviation before re-review |

Default behavior: after incorporating feedback, state what changed and recommend
whether re-review is warranted. Do not automatically re-run reviewers unless the
workflow requires it or the user asks.

### Post-Execution: Component Documentation Review

If this project has component documentation, run the `component-docs-reviewer`
agent after all sub-plans complete to catch implementation-vs-plan drift in
component docs.

## Plan Structures

Templates for plan files are in this skill's `assets/` directory. Read them when
creating plans:

- **[Master plan template][master-plan-template]** — Orchestration document structure with RFC mapping and lead-agent execution mechanics.
- **[Sub-plan template][sub-plan-template]** — Self-contained execution unit structure with RFC context.

## Documentation Sub-Plan

When a feature affects documented domain concepts, architecture, or business
processes, add a documentation sub-plan as the final sub-plan in the execution
order. This makes doc updates part of the plan: visible, reviewable, and
deliberate.

### What Goes in the Documentation Sub-Plan

| Doc Level | Planned Upfront? | Rationale |
|---|---|---|
| Domain docs | Yes | The planner knows what domain concepts are changing |
| Architecture docs | Yes | The planner knows what structural changes are happening |
| Process docs | Yes | The planner knows what flows are being added/modified |
| Component docs | No — post-execution | Component docs describe implementation details that may drift from the plan |

### How to Write It

The documentation sub-plan follows the standard sub-plan template but its
implementation steps are doc edits, not code changes. Be specific:

- **Which existing docs to update** — file paths, sections, and stale claims to replace.
- **Which new docs to create** — file paths, existing doc patterns to follow, and what the new doc should cover.
- **Structural pattern matching** — if existing docs follow a pattern, new additions must follow it.
- **Proportionality guard** — the target document's purpose, audience, and existing abstraction level outrank feature-local emphasis. The sub-plan should say whether the feature should be mentioned as a primary concept, a small example, or only as an implementation detail. Do not let a narrow feature become the center of a broad domain, architecture, or process doc unless the RFC explicitly changes that doc's central subject.
- **Required skills** — list the relevant documenting skills.
- **Execution model** — always assign the most capable available model. Documentation requires synthesis and judgment.

### When to Skip It

Skip the documentation sub-plan when:

- The feature does not affect any documented concepts, flows, or architecture.
- The only doc impact is component-level and can be handled by post-execution component docs review.
- No project documentation exists yet; recommend creating initial docs as a separate effort.

## Rules (Non-Negotiable)

- **Never use this workflow with an unreviewed RFC.** Return to `planning-project-features` routing instead.
- **Never change the RFC design silently.** Revise the RFC or record an explicit approved deviation.
- **Never ask the user to restate information that is already in the RFC.** Use the RFC as the approved baseline.
- **Always embed execution-critical RFC context into each sub-plan.** Executing agents should not need to read the RFC.
- **Always decompose into sub-plans.** A monolithic plan is a failure mode.
- **Each sub-plan must be self-contained.** Embed context, contracts, constraints, and acceptance criteria directly.
- **Always list required skills in every sub-plan.** An executing agent without the right skills will produce subpar results or get stuck.
- **Always respect model assignments during execution.** Sub-plan model assignments are deliberate cost-optimization decisions. The assigned model must be used through the active runtime's actual model-selection mechanism.
- **Use runtime-specific execution bindings for multi-sub-plan execution.** When a plan has two or more sub-plans, dispatch each sub-plan through its assigned execution binding. If a binding fails, diagnose and retry once. If it still cannot be resolved, stop and ask the user. Never silently execute on the coordinator or a more expensive model.
- **Write lead-agent execution instructions into every multi-sub-plan master plan.** Do not rely on the executor remembering this skill.
- **Always run `plan-rfc-fidelity-reviewer` and `plan-executability-reviewer` before presenting the plan.**
- **Never run full direct-planning design reviewers in this workflow unless the user explicitly exits RFC-backed planning.** RFC architecture, risk, and clarity were already reviewed at RFC time.
- **Save plans to the correct location.** Standalone features go to `plans/features/<feature-name>/`; epic features go to `plans/epics/<epic-name>/<feature-name>/`.

[master-plan-template]: assets/master-plan-template.md
[sub-plan-template]: assets/sub-plan-template.md
