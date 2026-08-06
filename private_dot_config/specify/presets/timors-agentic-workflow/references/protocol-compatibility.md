# Preset Protocol Compatibility

## Authoritative Contract

`preset.yml` at
`.specify/presets/timors-agentic-workflow/preset.yml` is authoritative for the
preset identity, preset and protocol versions, minimum Spec Kit requirement,
and declared command/template mappings. The current package requires Spec Kit
`>=0.12.11`; it has no upper support boundary.

- Preset ID: `timors-agentic-workflow`
- Preset version: `0.2.0`
- Execution-plan protocol version: `0.1.0`

## Shared Minimum-Version And Package-Integrity Gate

Every composed or replaced command in this preset, including commands added in
later package work, MUST complete this entire gate as its first operational
action. It is atomic: any failure stops the command before an extension hook,
prerequisite script, artifact write, progress or workspace mutation, or
dispatch. Do not substitute a global, dotfiles, or package-manager path for the
installed preset path.

1. Read and parse the installed `preset.yml`. Require schema version `1.0`, the
   identity above, and the exact `0.1.0` protocol version. Require a parseable
   `requires.speckit_version` constraint whose current declared minimum is
   `>=0.12.11`.
2. Obtain and parse the active semantic version with `specify --version`. It
   must meet the minimum and the installed manifest constraint. A later manifest
   may record an explicit known incompatibility in its constraint; honor that
   declared constraint only. Do not add a known-bad list, infer exclusions, or
   warn, prompt, or block merely because a newer version has not been tested.
3. Require every file in the complete inventory below to exist as a readable
   regular file at its installed relative path. Treat a missing, malformed, or
   incoherent package as one package-level failure; do not continue into a
   phase whose own inputs happen to be present.
4. Verify that `preset.yml` declares exactly the command/template mappings and
   strategies in the manifest mapping table below, with each declared relative
   path present in the inventory. This validates package wiring, not file
   contents.
5. Validate the structural coherence rules below. They check declared fields,
   headings, protocol values, and installed-preset paths only; they do not
   compare full prose or immutable file contents.

If the manifest, version, constraint, identity, inventory, or mapping is
missing, malformed, or incoherent, report the failed contract and stop. This
gate is structural only: do not compare checksums, hashes, immutable contents,
or exact file contents.

### Complete Required Package Inventory

| Category | Required installed relative paths |
|----------|----------------------------------|
| Root files | `LICENSE`; `README.md`; `preset.yml` |
| Commands | `commands/speckit.specify.md`; `commands/speckit.plan.md`; `commands/speckit.tasks.md`; `commands/speckit.analyze.md`; `commands/speckit.implement.md` |
| Templates | `templates/plan-template.md`; `templates/tasks-template.md`; `templates/execution-plan-template.md`; `templates/analysis-template.md`; `templates/review-report-template.md`; `templates/progress-template.md` |
| Reviewer packets | `reviewers/artifact-fidelity.md`; `reviewers/decomposition-design.md`; `reviewers/plan-clarity.md` |
| References | `references/analysis-and-approval.md`; `references/artifact-validation.md`; `references/checkpoint-integration.md`; `references/concurrency-policy.md`; `references/decomposition.md`; `references/documentation-planning.md`; `references/execution-lifecycle.md`; `references/human-review.md`; `references/model-and-worker-selection.md`; `references/planning-grounding.md`; `references/protocol-compatibility.md`; `references/testable-work.md`; `references/workspace-isolation.md` |

### Manifest Mapping Table

| Type | Name | File | Strategy |
|------|------|------|----------|
| command | `speckit.specify` | `commands/speckit.specify.md` | `wrap` |
| command | `speckit.plan` | `commands/speckit.plan.md` | `wrap` |
| template | `plan-template` | `templates/plan-template.md` | `wrap` |
| command | `speckit.tasks` | `commands/speckit.tasks.md` | `replace` |
| template | `tasks-template` | `templates/tasks-template.md` | `replace` |
| command | `speckit.analyze` | `commands/speckit.analyze.md` | `replace` |
| command | `speckit.implement` | `commands/speckit.implement.md` | `replace` |
| template | `execution-plan-template` | `templates/execution-plan-template.md` | `replace` |
| template | `analysis-template` | `templates/analysis-template.md` | `replace` |
| template | `review-report-template` | `templates/review-report-template.md` | `replace` |
| template | `progress-template` | `templates/progress-template.md` | `replace` |

### Structural Coherence Rules

- The manifest binds the five command source files to exactly
  `speckit.specify`, `speckit.plan`, `speckit.tasks`, `speckit.analyze`, and
  `speckit.implement`. Each source's first body phase is `## Package Integrity
  Gate` and consumes this installed `references/protocol-compatibility.md`
  before any command-specific body phase.
- `speckit.specify`, `speckit.plan`, and `plan-template` are the only `wrap`
  mappings. The specify and plan sources retain `{CORE_TEMPLATE}` and omit local
  `scripts` and `agent_scripts` metadata so native composition inherits those
  fields from the active lower-priority command. The `plan-template` retains
  its core composition placeholder.
- `speckit.specify` and `speckit.plan` each consume the installed shared
  `references/human-review.md` only after their respective upstream artifact
  generation and validation complete. The reference governs optional review
  mechanics; each command owns phase-specific conceptual guidance.
- `speckit.tasks`, `speckit.analyze`, `speckit.implement`, and
  `tasks-template` are `replace` mappings. The task command declares `sh`,
  `ps`, and `py` setup variants.
- Analyze and implement each consume the same installed
  `references/artifact-validation.md`. No command may substitute or define a
  divergent artifact-validation source.
- The reviewer packets declare these exact role/tier pairs:
  `artifact-fidelity` / `Mid-tier`, `decomposition-design` / `Most capable`,
  and `plan-clarity` / `Mid-tier`.
- `templates/execution-plan-template.md` declares
  `**Preset Protocol Version**: 0.1.0`, consistent with the manifest and this
  reference.
- Every required installed-preset path named by a command spine resolves to a
  required file in the complete inventory. This includes command references to
  installed policy references and templates; no command may rely on an
  undeclared installed-preset file.

## Execution-Plan Protocol Gate

`execution-plan.md` is valid only when its identity line is exactly
`# Execution Plan: <feature>` and its protocol line is exactly
`**Preset Protocol Version**: 0.1.0`. Missing, malformed, or unsupported
versions fail closed. The only recovery is regeneration through
`/speckit.tasks`; analyze and implement do not repair, infer, downgrade, or
silently accept a different protocol.

The deterministic protocol rules live only in
`.specify/presets/timors-agentic-workflow/references/artifact-validation.md`.
Later `speckit.analyze` and `speckit.implement` command contracts MUST consume
that installed source directly rather than copying or redefining its rules.

## Package Boundary

This preset defines no runtime worker, model binding, or agent-dispatch
configuration. Execution capabilities are declared only as data in validated
execution-plan artifacts.
