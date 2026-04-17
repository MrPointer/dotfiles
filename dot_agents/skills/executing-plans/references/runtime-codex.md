# Runtime Adapter: Codex

Use this adapter only when the active runtime is Codex.

This adapter maps the canonical execution workflow in `../SKILL.md` to Codex-native mechanics. It must not redefine task ordering, dispute policy, progress rules, or the test-author/implementer separation.

## Exploration and Dispatch

- Use Codex sub-agent dispatch for exploration, test authoring, and implementation work when sub-agents materially help.
- Set the sub-agent `model` explicitly using Codex's actual dispatch mechanism when the plan assigns a tier.
- Prefer narrow prompts with only the context required for the specific role.

## Execution Bindings

- In Codex, execution bindings should use the built-in worker path for implementation and test-author work unless the project explicitly defines a more specialized worker.
- In Codex, execution bindings are usually reusable dispatch recipes rather than checked-in worker files.
- Project-local Codex skills are expected under `.agents/skills/` in the project root. Global Codex skills are expected under `~/.agents/skills/`.
- Each binding should make explicit:
  - target `model`
  - required skills to attach explicitly as Codex `skill` items from the local or global Codex skills directories
  - whether the worker should write files directly or return results for the parent to persist
- Pass model explicitly in the worker dispatch and attach required skills explicitly as `skill` items.
- Do not rely on prompt text alone for skill loading when the runtime can attach the skills directly.
- If the binding is ephemeral rather than file-backed, record its parameters in the plan metadata or execution context so retries and resumed execution reuse the same model and skills.

## Test Author Isolation

- Prefer `wt` to create an isolated workspace for the test author.
- Use `git worktree` directly only as a fallback when `wt` is unavailable or unsuitable.
- Structural TDD in Codex is allowed only when the test author can be dispatched into that isolated workspace.
- If the active Codex environment cannot run the test-author worker inside the isolated workspace, skip structural TDD and record a reason such as `runtime cannot provide isolated test-author workspace`.
- Do not reveal the plan path, task file path, feature name, or design rationale to the test author.
- Pass only acceptance criteria and the code surface the tests interact with.
- When structural TDD is used, prompt hygiene is mandatory in addition to physical isolation.

## Implementer Dispatch

- Spawn a separate implementer sub-agent with the full task, test file paths, prerequisite outputs, and required skills/reference material.
- Tell the implementer explicitly that tests are immutable.
- If the implementer reports a dispute, record it and continue according to the canonical workflow.

## Progress and Artifacts

- The parent executor owns `progress.md` and should update it after each meaningful step.
- Even if a sub-agent writes files directly, the parent remains responsible for checkpointing and artifact verification.

## Model Assignment

- Use Codex's explicit model-selection mechanism in the worker dispatch path; do not rely on prompt wording.
- Treat plan-assigned tiers as binding.
- If the required model is unavailable, stop and ask the user how to proceed.
