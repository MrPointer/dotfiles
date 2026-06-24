# Runtime Adapter: Codex

Use this adapter only when the active runtime is Codex.

This adapter maps the canonical planning workflow in `../SKILL.md` to Codex-native mechanics. It must not redefine phases, review policy, or plan semantics.

## Discovery

- Treat the current session's available skills list as authoritative for global Codex skills.
- Discover project-local skills and conventions from `AGENTS.md` first.
- Project-local Codex reviewer agents are expected under `.codex/agents/` in the project root. Global Codex reviewer agents, if used by the runtime, are expected under `~/.codex/agents/`.
- Project-local Codex skills are expected under `.codex/skills/` in the project root. Global Codex skills are expected under `~/.codex/skills/`.
- If working from a dotfile source repository such as a chezmoi source dir, use the source-directory equivalent of those paths rather than the applied home-directory paths.
- Discover project-local reviewer definitions from the local Codex agents directory and treat them as required local reviewer bindings rather than optional hints.
- Assume local reviewer sub-agents exist as checked-in Codex custom-agent definition files.

## Reviewer Bindings

- Codex local reviewers should be represented as file-defined custom agents in `.codex/agents/`.
- Codex custom-agent files are the source of truth for reviewer identity, model, and `developer_instructions`.
- Discover and invoke those reviewers through Codex's custom-agent mechanism.
- Use the runtime's normal custom-agent path for reviewers; do not synthesize reviewer prompts ad hoc.
- Set the reviewer model using the reviewer definition or explicit dispatch configuration, whichever the active Codex runtime actually uses.
- If the reviewer can write files in its workspace, you may ask it to write directly to the requested review path.
- If direct file writing is awkward or unreliable, have the reviewer return structured findings and write the review file from the parent agent. Preserve the standard review format.

## Local Reviewer Requirement

- The planning workflow requires project-local reviewer coverage for project-specific domains.
- Global reviewers remain part of the master-plan review loop, but they do not replace local reviewers for sub-plan review.
- If a local reviewer is missing in Codex-native runnable form, the human maintainer should add the Codex custom-agent file under `.codex/agents/`. Until then, report the gap rather than silently substituting a generic reviewer.
- If the active Codex environment cannot invoke checked-in custom reviewer agents, report that runtime incompatibility and stop. Do not recreate the reviewer in prompt text.

## Execution Bindings

- Codex execution bindings should use the built-in worker path unless the project explicitly defines a more specialized Codex worker.
- In Codex, these bindings are usually reusable dispatch recipes rather than checked-in worker files.
- For each sub-plan, create or select a binding that specifies:
  - the target `model`
  - which required skills from project-local or global skill directories must be attached explicitly as Codex `skill` items
  - which additional reference files must be attached or read
  - any runtime constraints the worker must follow
- Attach required skills explicitly as Codex `skill` items when dispatching the worker.
- If the binding is ephemeral rather than file-backed, record its parameters in the plan metadata or execution context so retries and resumed execution reuse the same model and skills.
- Keep bindings thin. The sub-plan remains the source of task truth.
- For testable implementation work, ensure the binding attaches the testing skills named by the sub-plan, whether they are project-local or global. Do not create a separate test-writing dispatch recipe.

## Review Loop

- Run independent reviewers in parallel using Codex's reviewer dispatch mechanism when practical.
- Pass only the plan path, the review output location, and the review task. Do not mix in unrelated planning rationale.
- When re-reviewing after changes, re-run only the affected reviewers, matching the canonical workflow.

## Execution Dispatch

Direct feature plans with two or more sub-plans must include concrete lead-agent instructions, worker/dispatch-recipe tables, implementer worktree isolation, and result-integration mechanics in the master plan. During execution, use Codex's actual worker dispatch mechanism with the recipe's explicit model and skills. Do not rely on prompt text alone to pick the right model, and do not let the coordinator execute a sub-plan directly when the plan assigned a worker or model tier.

## Implementer Worktree Mechanics

Apply the master plan's Concurrency Policy before assigning same-group implementation. If the policy says `Linear DAG`, the master plan must instruct the executor to run one sub-plan at a time even when sub-plans are logically independent. Only when the policy says `Parallel allowed`, same-group sub-plans must require task-scoped implementer worktrees rather than concurrent workers in the coordinator workspace. The plan should reference the active execution adapter's Workspace Isolation Strategy instead of repeating the fallback chain. It must name ignored build/cache directories that isolated worktrees need before dispatch, or explicitly state that no seeding is required. If no isolated implementer path, required seeding path, or output/cache safety can be verified, the plan must instruct the executor to serialize the group or ask the user.

The plan must keep plan files, review files, and `progress.md` coordinator-owned. Implementers receive inline task packets and prerequisite outputs, not plan paths copied into worker worktrees.

## Testable Work Mechanics

For testable behavior changes, the master plan should tell execution that the implementer owns both tests and code. The implementer should follow the sub-plan's testing skills, make a test-first attempt when practical, and report tests added or updated plus verification results. If test-first work or new tests are not practical, the implementer reports the reason instead of relying on a separate test-writing worker.

## Model Assignment

- Use Codex's explicit model-selection mechanism in the reviewer or worker dispatch path rather than prompt-only requests.
- Treat the plan's model tier as binding. Do not silently upgrade or downgrade.
- If the requested model is unavailable in the current Codex environment, stop and ask the user how to proceed.

## Review Artifact Ownership

- The planner owns the `reviews/` directory and is responsible for ensuring the expected files exist.
- After each reviewer finishes, verify whether the review artifact exists.
- If not, write it from the returned findings before continuing.
