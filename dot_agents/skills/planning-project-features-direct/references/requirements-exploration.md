# Requirements And Exploration

Use this reference when gathering requirements, reading feature anchors, handling assumptions, or exploring code for direct planning.

## Requirement Gathering

Ask until you can answer these confidently:

- What problem does this solve?
- What does success look like?
- What is in scope?
- What is explicitly out of scope?
- What existing code or systems does this touch?
- What constraints apply, such as performance, compatibility, technology, security, operations, or release timing?
- What acceptance criteria prove the feature is done correctly?

When later phases reveal new ambiguity, stop and ask. Requirement gathering does not end after the first phase.

## Active Feature Anchor

Look for an active feature anchor using project conventions, then `docs/context/<topic>-anchor.md`.

If one exists, read it as supporting context before planning. It can explain intent, rejected alternatives, open questions, teammate handoff context, or recent decisions.

The anchor does not replace user-approved requirements. Any unresolved anchor question that affects decomposition, scope, acceptance criteria, or contracts is a planning blocker to ask about.

## When The User Says To Use Judgment

If the user says `just figure it out`, `use your judgment`, or similar:

1. Present the options you see.
2. Explain tradeoffs and consequences.
3. Help the user make an informed decision.
4. If they still want best-judgment planning, proceed and document assumptions clearly in the plan.

Do not silently convert ambiguity into hidden design decisions.

## Documentation-First Exploration

Once requirements are clear, read existing documentation before inspecting code.

1. Check `AGENTS.md` for documentation pointers.
2. Read relevant domain, architecture, business process, and component docs.
3. Explore code only for gaps documentation does not cover.

During code exploration, confirm only planning-relevant mechanics:

- files that will need modification
- architectural constraints and patterns to follow
- potential conflicts or risks
- required execution skills from `AGENTS.md`, project docs, and active runtime discovery
- ignored in-repository build/cache directories that isolated execution should seed, such as Rust `target/` or Bazel output trees

If no build/cache seeding is required, record that explicitly.

## Documentation Gaps

If critical areas needed for the plan are undocumented, note the gap and recommend the appropriate documenting skill:

| Gap | Recommended Skill |
|-----|-------------------|
| Missing domain knowledge | `documenting-domain` |
| Missing architecture overview | `documenting-architecture` |
| Missing business workflow docs | `documenting-business-processes` |
| Missing component docs | `documenting-components` |

Present gaps to the user. They may want to create docs before planning continues, or accept the gap and proceed.

## Specification Gaps

When the plan depends on a specification or requirements document that does not cover an implementation case, treat it as a blocking gap rather than a documentation gap.

Examples include an API spec that omits error responses, a data model that omits edge-case states, or a workflow description that omits failure paths.

Present the gap and options to the user before proceeding.

Share planning-relevant findings with the user and confirm understanding before writing plan files when findings affect decomposition, sequencing, or scope.
