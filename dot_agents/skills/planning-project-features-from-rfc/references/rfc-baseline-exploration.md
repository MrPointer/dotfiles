# RFC Baseline And Exploration

Use this reference when validating an RFC baseline, reading feature anchors, or exploring the codebase for planning mechanics.

## RFC Baseline Validation

Confirm the RFC is reviewed enough for RFC-backed planning:

- `rfc-architect-reviewer` is `Passed` or `Passed with concerns`.
- `rfc-risk-reviewer` is `Passed` or `Passed with concerns`.
- `rfc-clarity-reviewer` is `Passed`, `Passed with concerns`, or explicitly not required by the project.
- No `Blocking` review status remains.
- Any `Passed with concerns` item is compatible with planning and does not require a design decision before decomposition.

If validation fails, stop and return to `planning-project-features` routing. Do not continue as RFC-backed planning.

## Feature Anchor Handling

Look for an active feature anchor using project conventions, then `docs/context/<topic>-anchor.md`.

If an anchor exists, read it as supporting handoff context. It may explain intent, rejected alternatives, open questions, teammate handoff context, or recent decisions.

The anchor does not override the RFC. Any unresolved anchor question that conflicts with the RFC or affects decomposition, scope, acceptance criteria, or contracts is a planning blocker to ask about.

## Documentation-First Exploration

Read existing documentation before inspecting code. Documentation is cheaper than re-exploring code from scratch.

1. Check `AGENTS.md` for documentation pointers.
2. Read relevant domain, architecture, business process, and component docs.
3. Explore code only for gaps the documentation and RFC do not cover.

During code exploration, confirm only planning-relevant mechanics:

- file paths and ownership boundaries needed for sub-plans
- existing interfaces, commands, schemas, config files, or runtime bindings named by the RFC
- test packages and verification scopes that acceptance criteria can reference
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

## RFC Or Specification Gaps

When implementation planning depends on behavior the RFC does not define, treat it as a blocking gap. Do not fill it with a reasonable default.

Ask whether to revise the RFC or approve an explicit plan deviation.

Do not use exploration to redesign the RFC. If exploration contradicts the RFC, stop and ask whether to revise the RFC or abandon RFC-backed planning.

Share planning-relevant findings with the user and confirm understanding before writing plan files when findings affect decomposition, sequencing, or scope.
