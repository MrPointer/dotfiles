---
kind: subplan
primary_files:
  - docs/architecture.md
  - docs/dotfiles/processes/shell-startup.md
  - docs/dotfiles/domain.md
execution:
  model_tier: most-capable
  rationale: "Final documentation packet over repository-wide architecture and shell-runtime process docs; it must describe the integrated result of packets 03 to 05 accurately and in proportion."
  skills:
    - documenting-architecture
    - documenting-business-processes
    - documenting-domain
test:
  mode: no-testable-behavior
  basis: "Documentation only; no runtime behavior changes. Non-test verification: each updated statement is checked against the integrated files from packets 03, 04, and 05, and every relative link in the changed sections resolves to an existing file."
contracts:
  produces: []
  consumes: []
---

# Document the Validation Boundary and the Optional-Tool Guards

## Objective

Update repository-wide docs so they describe pre-merge dotfile validation, the
merge gate, and the runtime guards for Sheldon, the git pager, and the git
external diff.

## Baseline Context and Constraints

- Accepted plan: `plan.md` R6 › Planning Handoff › Documentation impact; the
  installer process docs are owned by packet 01, not this packet.
- Stale sections today:
  - `docs/architecture.md` › Structure and › Communication Patterns describe the
    installer, templates, and shell runtime but no CI validation boundary; the
    › Key Design Decisions table has no entry for pre-merge validation.
  - `docs/dotfiles/processes/shell-startup.md` › Flow › `.zshrc` Phase step
    "Load sheldon plugins" states an unconditional `eval "$(sheldon source)"`;
    › Failure Scenarios has no case for Sheldon missing.
  - `docs/dotfiles/domain.md` describes the common git configuration
    (`private_dot_config/git/common.gitconfig.tmpl`) without the pager and
    external-diff fallbacks.
- Keep changes proportional: add or correct only what the integrated change
  makes stale. Do not restate the plan.
- Do not modify anything under `specs/001-dotfile-regression-ci/`.

## Scope and Primary Files

- `docs/architecture.md`: the CI validation boundary (Installer CI calls the
  Dotfile Validation reusable workflow; the two required aggregate checks;
  pull-request-only changes with auto-merge) and one Key Design Decisions row.
- `docs/dotfiles/processes/shell-startup.md`: the guarded Sheldon step and a
  "Sheldon not installed" failure scenario.
- `docs/dotfiles/domain.md`: the pager and external-diff fallbacks and the
  `~/.local/bin/git-external-diff` helper.
- Excluded: installer docs, `README.md`, skills, specs.

## Prerequisites and Integration

### Direct Predecessors

| Path | Classification | Reason |
|---|---|---|
| subplans/03-dotfile-optional-tool-guards.md | genuine implementation prerequisite | The shell-startup and domain docs describe the guards and helper that packet 03 adds. |
| subplans/04-dotfile-validation-workflow.md | genuine implementation prerequisite | The architecture doc describes the validation workflow and aggregate check that packet 04 adds. |
| subplans/05-auto-merge-workflow.md | genuine implementation prerequisite | The architecture doc describes automatic merge, which packet 05 adds. |

### Transitive Context

Packet 04 builds on packets 01 and 02.

## Runtime Considerations

- None.

## Acceptance and Verification

- Each stale section listed above is updated, and each updated statement
  matches the integrated files.
- Every relative link in changed sections resolves.
- No other section changes.
- Report the non-test verification performed.

## Required Outputs and Handoff

- One packet commit.
- This packet and packet 05 form delivery PR C.
