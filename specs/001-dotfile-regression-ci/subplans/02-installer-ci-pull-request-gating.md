---
kind: subplan
primary_files:
  - .github/workflows/installer-ci.yml
execution:
  model_tier: mid-tier
  rationale: "Two contained edits to an existing workflow: remove one trigger and tighten one aggregate condition. Getting the result-expression logic right needs care, but no new jobs or design."
  skills:
    - configuring-github-actions
test:
  mode: not-applicable
  basis: "Barrier: trigger and job-result behavior exist only in hosted GitHub Actions runs on a pushed pull request. Packet workers cannot push, and no local Actions runner is installed. Alternative verification: parse the YAML, inspect the trigger block and the Installer CI Status condition against every result combination listed in Acceptance, and run shellcheck on any changed run script. Named gap: the hosted cancellation and trigger evidence in plan.md › Verification Strategy, collected on PR A after implement (tasks.md › Execution Handoff)."
contracts:
  produces: []
  consumes: []
---

# Run Installer CI on Pull Requests Only and Fail on Canceled Jobs

## Objective

In `.github/workflows/installer-ci.yml`, remove the `push` trigger so the
workflow runs only on `pull_request`. Change the `Installer CI Status` aggregate
job so that it fails when `build`, `test`, `privilege-test`, or `e2e-test` is
canceled or failed, and still passes when they are skipped because the pull
request has no installer changes. Remove the `changes` job lines that served
only pushes.

## Baseline Context and Constraints

- Accepted plan: `plan.md` R6 › Design Decisions › Run Installer CI on Pull
  Requests Only; › Fail Installer CI Status on Canceled Jobs; › Contracts and
  Flows › Aggregate check contracts.
- Current state: the `on:` block also triggers on pushes to `main` that touch
  installer paths. The `Installer CI Status` job passes skipped jobs and passes
  when `build`, `test`, or `e2e-test` is canceled; only a canceled
  `privilege-test` fails it.
- A job timeout is reported as canceled and must fail the aggregate.
- Keep: the `pull_request` trigger, the `changes` job and its filter, the
  existing concurrency group, the reset of workflow permissions to `{}` with
  per-job read grants, and the conditional installer jobs.
- `release.yml` is unchanged.
- Do not modify anything under `specs/001-dotfile-regression-ci/`.

## Scope and Primary Files

- `.github/workflows/installer-ci.yml`: the `on:` block, the
  `Installer CI Status` job, and in the `changes` job the checkout step guarded
  by `if: github.event_name == 'push'` and the paths-filter `base` input, which
  serve only pushes. The paths filter itself is unchanged.
- Excluded: dotfile validation jobs, the exemption output, and
  `Dotfile Validation Status`; packet 04 adds those to this same file afterward.

## Prerequisites and Integration

### Direct Predecessors

| Path | Classification | Reason |
|---|---|---|

### Transitive Context

None.

## Runtime Considerations

- None beyond a local YAML parser and shellcheck.

## Acceptance and Verification

- The parsed workflow has `pull_request` as its only trigger.
- The `changes` job has no push-only checkout step and no paths-filter `base`
  input, and its filter output is unchanged.
- Inspection of the `Installer CI Status` condition shows this outcome for each
  case:
  - `changes` failed or was canceled: fail;
  - all four installer jobs succeeded: pass;
  - all four skipped and `changes` succeeded: pass;
  - any of `build`, `test`, `privilege-test`, `e2e-test` failed: fail;
  - any of them canceled: fail.
- `python3 -c 'import yaml,sys; yaml.safe_load(open(sys.argv[1]))'
  .github/workflows/installer-ci.yml` exits 0.
- Report the alternative verification performed and restate the named hosted
  gap.

## Required Outputs and Handoff

- One packet commit.
- Handoff state for packet 04: `installer-ci.yml` has only the `pull_request`
  trigger, a `changes` job without push-only lines, and the tightened
  `Installer CI Status`; packet 04 adds jobs to this
  file on top of that state and must keep all three changes.
- This packet and packet 01 form delivery PR A.
