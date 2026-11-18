---
kind: subplan
primary_files:
  - .github/workflows/dotfile-validation.yml
  - .github/workflows/installer-ci.yml
  - .github/scripts/validate-dotfiles.sh
execution:
  model_tier: most-capable
  rationale: "Cross-workflow design with a reusable workflow, an eight-entry matrix over a container and a hosted macOS runner, job-result conditions that decide merge gating, exact path exemptions, and seven pass criteria. Mistakes here silently let regressions merge."
  skills:
    - configuring-github-actions
    - managing-chezmoi
    - configuring-zsh
    - testing-e2e-containers
test:
  mode: not-applicable
  basis: "Barrier: the workflows run only on hosted GitHub Actions runners for a pushed pull request, workers cannot push, no local Actions runner is installed, and macOS jobs cannot run in a container. Alternative verification: parse both workflows, shellcheck the validation script, and run the script locally for all four Ubuntu cases (bare personal, bare work, with-tools personal, with-tools work) inside the installer/docker/ubuntu image with an installer built from this branch, including one run that must fail on an empty managed list. Named gap: hosted evidence for the four macOS jobs, the hosted runs of all eight jobs, the exemption cases, and ordering, collected on PR B and later per tasks.md › Execution Handoff."
contracts:
  produces: []
  consumes: []
---

# Add the Dotfile Validation Reusable Workflow and Its Merge Gate

## Objective

Add dotfile validation as a reusable workflow called from Installer CI. Eight
matrix jobs, {Ubuntu, macOS} × {personal, work} × {bare, with-tools}, prepare a
clean home with the installer's `--prepare-only` stop, apply the pull-request
merge commit from the checkout, and check pass criteria 1 to 7. Add the
exemption output and the always-running `Dotfile Validation Status` aggregate.

## Baseline Context and Constraints

- Accepted plan: `plan.md` R6 › Design Decisions › Apply the Pull-Request Merge
  Commit From the Checkout; › Choose the Installer Binary by Pull-Request
  Content; › Ubuntu Runs in the Existing Minimal Test Image; macOS Asserts the
  Bare Baseline; › Define Validation Pass Criteria; › Exempt Only a Fixed List of
  Documentation Paths; › Gate Merges on One Aggregate Check; › Run Dotfile
  Validation as a Reusable Workflow Called by Installer CI; › Contracts and
  Flows; › Proposed Architecture. Those sections define criteria 1 to 7, the
  exempt list, the command-name map, and the aggregate conditions; implement them
  exactly.
- `.github/workflows/dotfile-validation.yml`: `on: workflow_call` only. One input
  selects the installer source: the `build` job artifact, or the latest release
  downloaded for the job's platform with the job token (for example
  `gh release download`).
- `installer-ci.yml` additions:
  - the `changes` job gains an output that is true only when every changed path
    matches the exempt list at its exact location;
  - a new job calls the reusable workflow. It needs `changes`, `build`, `test`,
    `privilege-test`, and `e2e-test`, and runs only when the pull request is not
    exempt and either the installer filter is false or all four installer jobs
    succeeded;
  - a job named `Dotfile Validation Status` always runs and passes only when
    `changes` succeeded and marked the pull request exempt, or the caller job
    succeeded. It fails in every other case, including canceled or failed
    `changes`, and a non-exempt pull request whose validation failed, was
    canceled, or was skipped.
- Preparation flags: `--non-interactive --plain --install-prerequisites
  --prepare-only`, plus `--install-tools` in with-tools jobs only, plus
  `--work-env` in work jobs with default work name and email.
- Ubuntu jobs build `installer/docker/ubuntu/Dockerfile` and run as `testuser`
  inside it. macOS jobs use `macos-latest` with `HOME` set to a new empty
  directory and remove nothing.
- Every `chezmoi` command passes `--source <checkout>` and `--config <installer
  config path>` and runs with `CHEZMOI_GITHUB_ACCESS_TOKEN` set from the job
  token.
- Permissions: keep the workflow-level `{}` reset; grant only `contents: read`
  per job, and keep `pull-requests: read` on `changes`. Keep the existing
  per-pull-request concurrency group.
- Predecessor state: packet 02 left `installer-ci.yml` with only the
  `pull_request` trigger, a `changes` job without push-only lines, and an
  `Installer CI Status` that fails on canceled jobs. Keep all three. Packet 01 added `--prepare-only`. Packet 03 added the guards
  that the bare jobs need.
- The validation steps live in `.github/scripts/validate-dotfiles.sh` so the
  same checks run in the container, on macOS, and locally. The script takes the
  variant (bare or with-tools), the checkout path, and the chezmoi config path
  as inputs, and names the failing criterion and, for criterion 7, the tool.
- Do not modify anything under `specs/001-dotfile-regression-ci/`. Do not push
  or open pull requests.

## Scope and Primary Files

- `.github/workflows/dotfile-validation.yml`: new.
- `.github/workflows/installer-ci.yml`: `changes` output, caller job, and
  `Dotfile Validation Status`.
- `.github/scripts/validate-dotfiles.sh`: new; criteria 1 to 7.
- Excluded: installer code, dotfiles, `release.yml`, hosted rulesets,
  the auto-merge workflow.

## Prerequisites and Integration

### Direct Predecessors

| Path | Classification | Reason |
|---|---|---|
| subplans/01-installer-prepare-only.md | genuine implementation prerequisite | The workflow and the local alternative verification run the installer with `--prepare-only`, which packet 01 adds. |
| subplans/02-installer-ci-pull-request-gating.md | policy-only serialization | Both packets edit `.github/workflows/installer-ci.yml`; this packet builds on packet 02's committed state of that file. |
| subplans/03-dotfile-optional-tool-guards.md | genuine implementation prerequisite | The local bare run must start zsh without Sheldon and run git without difftastic and diff-so-fancy, which needs packet 03's guards. |

### Transitive Context

None.

## Runtime Considerations

- The local container run needs Docker, network access for the installer's
  prerequisites and chezmoi externals, and a GitHub token for
  `CHEZMOI_GITHUB_ACCESS_TOKEN` (for example from `gh auth token`).
- Build the installer for linux from `installer/` for the container run.

## Acceptance and Verification

- Both workflows parse as YAML.
- `shellcheck .github/scripts/validate-dotfiles.sh` passes.
- Inspection shows:
  - the matrix has exactly eight entries with the flags above;
  - the caller job's `needs` and `if` match the conditions above;
  - `Dotfile Validation Status` uses `if: always()` and passes only in the two
    cases above;
  - the exempt filter matches exactly the plan's list, with root files matched
    only at the repository root.
- Local Ubuntu run inside the `installer/docker/ubuntu` image as `testuser`,
  with an installer built from this branch: the bare personal, bare work,
  with-tools personal, and with-tools work cases each pass criteria 1 to 7. The
  with-tools runs show every listed tool present for criterion 7, and the logs
  show criteria 2 and 3 called with `--source` and `--config`.
- A local run of the criterion 3 step against an empty managed list exits
  non-zero.
- Report the alternative verification performed and restate the named hosted
  gap.

## Required Outputs and Handoff

- One packet commit.
- Handoff state: `installer-ci.yml` keeps packet 02's trigger, `changes` job
  cleanup, and `Installer CI Status`, and adds the caller job and `Dotfile Validation Status`.
- Report the exact required-check name `Dotfile Validation Status` for the
  ruleset step in `tasks.md` › Execution Handoff.
- This packet and packet 03 form delivery PR B.
