# Writing The RFC

Use this reference when drafting RFC prose, choosing the RFC location/name, deciding which sections belong, or polishing source references.

## Human Engineering Voice

Write in direct human engineering prose.

- Convert conversation facts into project facts, constraints, or decisions.
- Use `we` only when it describes the project/team's chosen direction or operational responsibility.
- Avoid referring to `the user`, `the agent`, `the model`, `the assistant`, `the reviewer`, or `the conversation` in the RFC body.
- Put conversation-derived information in neutral terms such as `Requirement`, `Constraint`, `Decision`, or `Assumption`.
- Use a neutral source label such as `Requirement or decision input` in the Source References table when a fact cannot be derived from code or docs.
- Do not write transcript-style rationale. Explain why the design is right, not who said what during discovery.
- Do not optimize the language for downstream agents. Planning agents benefit from the same clear human-readable design record that engineers use.

Before review, scan the RFC body and rewrite agent-centric or transcript-like phrasing, including `the user`, `the assistant`, `the agent`, `the conversation`, `explicitly accepted`, `asked for`, or `decided during brainstorming`.

## Location And Naming

Save the draft RFC using the project's convention. If none exists, use:

```text
docs/rfcs/<topic>.md
```

Use a stable kebab-case topic slug and title the document `# RFC: <Title>` unless the project uses another heading style.

Numbered RFC IDs are optional, not the default. Use IDs only when the project already has an `RFC-0001` convention, the user explicitly requests numbered RFCs, or the repository needs stable numbered cross-references across many RFCs. When IDs are appropriate, assign the next sequential ID by inspecting existing RFCs in `docs/rfcs/`, then preserve that ID across revisions.

## Conceptual Sections

Scale detail to design complexity, but keep these concepts present:

- problem, goals, non-goals, and constraints
- verified current architecture
- chosen design and decision summary
- proposed architecture: boundaries, responsibilities, contracts, data/control flow, state, lifecycle, and failure behavior
- alternatives considered and why they were rejected
- risks, tradeoffs, design success criteria, and planning handoff notes
- source references showing what was verified

Remove irrelevant optional sections rather than keeping empty headings. Do not leave placeholders.

## Design-Level Detail

The RFC must contain enough detail for planning to decompose the work without re-deciding architecture.

Pseudo-code is allowed when it clarifies an algorithm, contract, state transition, or data/control flow. Keep it illustrative and language-agnostic unless the RFC is explicitly about a language-specific contract.

Do not write practical implementation code, low-level recipes, task lists, execution sequencing, tickets, test commands, lint commands, or build commands.
