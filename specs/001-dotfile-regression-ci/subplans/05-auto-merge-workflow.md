---
kind: subplan
primary_files:
  - .github/workflows/auto-merge.yml
execution:
  model_tier: mid-tier
  rationale: "One short workflow, but it holds the only write-scoped token in pull-request CI, so its event, draft, and same-repository conditions must be exact."
  skills:
    - configuring-github-actions
test:
  mode: not-applicable
  basis: "Barrier: the workflow acts only on hosted pull-request events, workers cannot push or open pull requests, and no local Actions runner is installed. Alternative verification: parse the YAML and inspect triggers, conditions, and permissions against Acceptance. Named gap: the hosted auto-merge evidence after PR C merges, per tasks.md › Execution Handoff."
contracts:
  produces: []
  consumes: []
---

# Enable Auto-Merge From a Workflow

## Objective

Add `.github/workflows/auto-merge.yml`. On `pull_request` events `opened`,
`reopened`, and `ready_for_review`, for a non-draft pull request whose head
repository is this repository, it runs `gh pr merge --auto --squash` so GitHub
merges once all required checks pass, with no human action.

## Baseline Context and Constraints

- Accepted plan: `plan.md` R6 › Design Decisions › Enable Auto-Merge From a
  Workflow; › Alternatives and Risks › Write token in pull-request workflow.
- The workflow resets permissions to `{}` at workflow level. Its one job is the
  only job with `contents: write` and `pull-requests: write`.
- Draft pull requests are skipped; drafts are the opt-out.
- The global pull-request skill under `dot_agents/skills/creating-pull-requests/`
  is not changed.
- Do not modify anything under `specs/001-dotfile-regression-ci/`.

## Scope and Primary Files

- `.github/workflows/auto-merge.yml`: new.
- Excluded: every other workflow and skill.

## Prerequisites and Integration

### Direct Predecessors

| Path | Classification | Reason |
|---|---|---|

### Transitive Context

Delivery order puts this packet in PR C, after the ruleset changes, so auto-merge
never acts while `Dotfile Validation Status` is not yet required. That order is
enforced in `tasks.md` › Execution Handoff, not by this packet.

## Runtime Considerations

- None beyond a local YAML parser.

## Acceptance and Verification

- The YAML parses.
- Inspection shows:
  - triggers are exactly `pull_request` with types `opened`, `reopened`,
    `ready_for_review`;
  - the job condition requires a non-draft pull request whose head repository
    equals this repository;
  - the only run step is `gh pr merge --auto --squash` for the event's pull
    request, authenticated with the job token;
  - workflow permissions are `{}` and only this job has `contents: write` and
    `pull-requests: write`.
- Report the alternative verification performed and restate the named hosted
  gap.

## Required Outputs and Handoff

- One packet commit.
- This packet and packet 06 form delivery PR C.
