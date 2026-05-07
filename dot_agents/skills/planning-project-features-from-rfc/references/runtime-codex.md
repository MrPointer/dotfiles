# Runtime Adapter: Codex

Use this adapter only when the active runtime is Codex.

This adapter maps `planning-project-features-from-rfc` to Codex-native mechanics. It must not redefine phases, routing, review policy, or plan semantics.

## Discovery

- Treat the current session's available skills list as authoritative for global Codex skills.
- Discover project-local skills and conventions from `AGENTS.md` first.
- Project-local Codex reviewer agents are expected under `.codex/agents/`.
- Global Codex reviewer agents, if used by the runtime, are expected under `~/.codex/agents/`.
- If working from a dotfile source repository, use source-directory equivalents rather than applied home-directory paths.

## Required Reviewers

RFC-backed feature planning uses exactly these reviewers:

- `plan-rfc-fidelity-reviewer`
- `plan-executability-reviewer`

Do not launch `plan-architect-reviewer`, `plan-risk-reviewer`, `plan-clarity-reviewer`, or project-local domain reviewers in this workflow unless the user explicitly exits RFC-backed planning and returns to direct planning.

## Reviewer Bindings

Use Codex's custom-agent or reviewer dispatch mechanism. Do not synthesize reviewer prompts ad hoc.

Pass only:

- Plan directory path.
- RFC path.
- Requested review output path, such as `reviews/00-master.rfc-fidelity.md` or `reviews/00-master.executability.md`.
- The review task.

If direct file writing is supported, ask the reviewer to write to the requested path. Otherwise have it return structured findings and write the review file from the parent agent.

## Execution Bindings

Codex execution bindings are usually reusable dispatch recipes rather than checked-in worker files.

For each sub-plan, create or select a binding that specifies:

- Target `model`.
- Required project skills to attach explicitly as Codex `skill` items.
- Additional reference files required by the sub-plan.
- Runtime constraints the worker must follow.

If the binding is ephemeral, record its parameters in plan metadata or execution context so retries and resumed execution use the same model and skills.

## Review Artifact Ownership

The planner owns the `reviews/` directory and must verify each expected review artifact exists after reviewer completion.
