# Artifact Authorities and Consumer Gates

| Artifact | Authority | Required consumer state |
|---|---|---|
| `spec.md` | Feature Definition | `status: Ready` before `plan` |
| `plan.md` | Sole normative plan | current Revision has human Design Acceptance `Accepted` before `tasks` |
| `tasks.md` and indexed `subplans/` | immutable execution package | both current reviews permit and `authorization: approved` before `implement` |
| role reports | cumulative review evidence | pointer resolves to an attributable applicable non-blocking complete round |
| ignored `progress.md` | coordinator-owned runtime evidence | exact durable evidence before resume, integration, or completion |

The canonical feature directory is resolved deterministically before any artifact
read. The resolution order is a caller-supplied feature directory, then the
`feature_directory` value in `.specify/feature.json`, then the single feature
directory under `specs/` when exactly one exists. Two or more `specs/` candidates
without a `.specify/feature.json` anchor fail closed rather than resolve by guess.
The resolved directory holds the canonical `spec.md`, `plan.md`, and `tasks.md`.
Consumers read concrete fields and invariants needed for their own phase and fail
closed on missing, partial, contradictory, or malformed state. The installed
project copy is runtime policy; commands do not consult the global distribution
source or alter installed files.

Material Feature Definition writes are forbidden after the initial progress lock
or completion. Before an otherwise permitted material definition write, mark a
current Accepted or Rejected Design Acceptance Stale while preserving decision
provenance; Pending remains Pending. Return authorization to `pending` and make
existing plan and task review pointers pending with no applicable round. Before a
material accepted plan write, return authorization to `pending`, apply the same
acceptance transition, and make task review pointers pending. Consumers never
repair a partially written transition by inference.

No tracked planning artifact carries concrete provider identity, model ID,
credentials, workspace path, runtime locator, dispatch identity, or a field for
runtime compatibility negotiation. Runtime evidence belongs only in ignored
progress. Old artifacts are not converted by a consumer; see the operator
procedure in `README.md`.
