# Timor's Agentic Workflow

Preset `0.3.0` supplies one project-installed lifecycle:

`specify` → `clarify` when Draft → `plan` → `tasks` → `implement`.

`spec.md` is the Feature Definition, `plan.md` is the sole normative RFC,
`tasks.md` and `subplans/` are the reviewed execution package, and ignored local
`progress.md` is coordinator runtime evidence. Design Acceptance precedes task
creation; Implementation Authorization precedes mutation.

## Installation and Refresh

Chezmoi manages the global distribution source but does not activate a project.
From a project root:

```sh
specify init --here --integration opencode --ignore-agent-tools
specify preset add --dev "$HOME/.config/specify/presets/timors-agentic-workflow"
```

Project copies remain pinned. To distribute a source update, the operator
performs this explicit remove/add refresh:

```sh
specify preset remove timors-agentic-workflow
specify preset add --dev "$HOME/.config/specify/presets/timors-agentic-workflow"
```

After installation, manually inspect generated integration files: `specify`,
`clarify`, `plan`, `tasks`, and `implement` must reference the project-relative
installed copy; `analyze` must resolve from a lower layer. This confirmation is
operator-controlled, not a workflow command gate.

## Execution Boundaries

Task packages contain provider-neutral tiers and exact skills, not concrete
workers, model IDs, credentials, local workspace paths, or dispatch identities.
Binding is just in time and requires tier, skills, discoverability, invokability,
correct workspace targeting, and attributable results. Project-local candidates
are preferred; ambiguity requires a human choice. Runtime considerations neither
grant access nor alter eligibility.

Each packet has a dedicated branch and worktree. Workers may commit only there;
the coordinator alone advances the active branch. The approved tracked feature
directory is immutable during execution and is compared at dispatch, result
acceptance, rebase, and integration. Packet branches integrate in topological
order as one repository-policy-compliant commit per packet. Nothing is pushed.

## Replacement Cautions

This release performs no migration. It does not interpret old atomic task lists,
old execution or review artifacts, or old local runtime state. A non-executing
old feature starts again at `specify`. Before replacement, finish or explicitly
abandon an in-flight old run while its old installed copy remains available.
Premature replacement has no automatic detection, rollback, or repair: manually
restore the old preset to recover it, or preserve old state and start a fresh
lifecycle. Do not remove an active package expecting this preset to resume it.
