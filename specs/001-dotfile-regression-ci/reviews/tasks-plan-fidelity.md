# Review: Task Plan Fidelity

- **Role**: Task plan-fidelity reviewer (`reviewers/tasks-plan-fidelity.md`)
- **Tier**: mid-tier
- **Latest Round**: R02
- **Latest Verdict**: passed

## Rounds

### R01

- **Trigger**: initial
- **Scope**: Complete current-state trace of the Ready Feature Definition and the accepted plan (Revision R6,
  Design Acceptance Accepted) into the task package. Covers each packet's scope, constraints, risks, accepted
  design, test-mode alignment with the plan's Verification Strategy, documentation placement, and the
  Execution Handoff delivery order and hosted evidence. Decomposition, dependency, and executability mechanics
  are left to the executability reviewer.
- **Reviewed Paths**:
  - `specs/001-dotfile-regression-ci/spec.md`
  - `specs/001-dotfile-regression-ci/plan.md`
  - `specs/001-dotfile-regression-ci/tasks.md`
  - `specs/001-dotfile-regression-ci/subplans/01-installer-prepare-only.md`
  - `specs/001-dotfile-regression-ci/subplans/02-installer-ci-pull-request-gating.md`
  - `specs/001-dotfile-regression-ci/subplans/03-dotfile-optional-tool-guards.md`
  - `specs/001-dotfile-regression-ci/subplans/04-dotfile-validation-workflow.md`
  - `specs/001-dotfile-regression-ci/subplans/05-auto-merge-workflow.md`
  - `specs/001-dotfile-regression-ci/subplans/06-repository-documentation.md`
  - Referenced by the plan's Documentation impact, read only to check its condition:
    `installer/docs/architecture.md`, `installer/docs/processes/tools-installation.md`
- **Binding Evidence**: Candidate bound by the coordinator to the task plan-fidelity role at mid-tier; role
  matches `reviewers/tasks-plan-fidelity.md`. Workspace: feature worktree on branch `001-dotfile-regression-ci`
  at HEAD `36127cc`, clean, read-only for every reviewed artifact. Invocation: the coordinator's initial
  assignment for R01 with no prior report. Started and completed 2026-11-18 in one session. Result: this round;
  the only write is this report.
- **Verdict**: passed-with-concerns

#### Current Findings

### Repository Documentation Packet Changes the Content of PR C

- **Impact**: concern
- **Location**: `tasks.md` › Overview (PR C = packets 05 and 06) and › Execution Handoff step 6;
  `subplans/06-repository-documentation.md` › Required Outputs and Handoff
- **Evidence**: `plan.md` › Planning Handoff › Order defines PR C as "the auto-merge workflow" only. Its
  Documentation impact names `docs/architecture.md` and the `docs/dotfiles/` shell-runtime docs but assigns
  them to no delivery pull request. The task package places packet 06 in PR C. As a result, the docs for the
  PR B guards land only after the hosted ruleset changes. This placement follows the preset rule that the
  documentation packet runs last, and it touches only exempt `docs/**` paths, so no gate or behavior changes.
  But it is a delivery composition the plan does not state.
- **Recommendation**: In `tasks.md` › Overview, state that `plan.md` R6 leaves the repository docs unassigned
  to a delivery pull request, and that packet 06 joins PR C as the final documentation packet. Alternatively,
  if the human prefers, deliver packet 06 as a separate docs-only pull request after PR C. Get human
  confirmation of either choice.

### Execution Handoff Abridges the Negative Ordering and Sheldon Evidence

- **Impact**: concern
- **Location**: `tasks.md` › Execution Handoff › Negative and ordering evidence
- **Evidence**: `plan.md` › Verification Strategy › Ordering evidence requires two things. Dotfile validation
  must start only after `test`, `privilege-test`, and `e2e-test` succeed and must use the `build` artifact.
  Also, "a deliberately failing installer test skips dotfile validation and fails `Dotfile Validation
  Status`". The handoff says only "The first installer pull request after PR B provides the ordering
  evidence", which drops the deliberately failing installer test case. › Negative evidence also says the
  unguarded `sheldon` call "fails the aggregate". The handoff says only that it "fails the bare jobs". Both
  items are still in the cited plan scope, but an executor working from the handoff list could skip the
  negative half of the merge-gate evidence.
- **Recommendation**: In `tasks.md` › Execution Handoff › Negative and ordering evidence, restate both plan
  items:
  - The ordering evidence also needs a throwaway draft installer pull request with a deliberately failing
    installer test. Its dotfile validation must be skipped, and `Dotfile Validation Status` must fail.
  - The unguarded `sheldon` case must also fail `Dotfile Validation Status`.

Traced without deviation:

- Packet 01: the `--prepare-only` contract, the `--install-tools` interaction, the three unit-test sequences,
  and the installer process docs.
- Packet 02: the pull-request-only trigger and the canceled-job failure.
- Packet 03: use-time guards and the `$HOME`-based external-diff helper.
- Packet 04: the matrix, preparation flags, installer source selection, the Ubuntu container and macOS clean
  home, criteria 1 to 7, the exact exempt list, the caller `needs`/`if`, the aggregate conditions,
  permissions, and token use.
- Packet 05: auto-merge events, draft and same-repository conditions, and write scope.
- Delivery: the order and the outward-facing confirmations.
- Interim release rule.
- Documentation impact: `installer/docs/architecture.md` lists no install flags, so the plan's conditional
  doc update does not apply. Component documentation review is recorded as required.
- No packet introduces scope, behavior, contracts, or acceptance exceptions beyond `spec.md` and `plan.md` R6.

### R02

- **Trigger**: targeted change
- **Scope**: Complete current-state snapshot after the human accepted all R01 findings. The rerun traces the
  changed sections against `spec.md` and `plan.md` R6:
  - `tasks.md` › Overview: the placement of packet 06 in PR C.
  - `tasks.md` › Execution Handoff: the restored negative and ordering evidence.
  - Packet 01: `tools-installation.md` added to its owned docs, exact test names, and `task lint`.
  - Packet 02: the `changes` failure case, plus removal of the push-only checkout step and the paths-filter
    `base` input.
  - Packet 04: local Ubuntu with-tools runs.
  - Packet 06: a note on the PR C placement.

  Packets 03 and 05 and the frontmatter DAG are unchanged. Their R01 trace still holds.
- **Reviewed Paths**:
  - `specs/001-dotfile-regression-ci/spec.md`
  - `specs/001-dotfile-regression-ci/plan.md`
  - `specs/001-dotfile-regression-ci/tasks.md`
  - `specs/001-dotfile-regression-ci/subplans/01-installer-prepare-only.md`
  - `specs/001-dotfile-regression-ci/subplans/02-installer-ci-pull-request-gating.md`
  - `specs/001-dotfile-regression-ci/subplans/04-dotfile-validation-workflow.md`
  - `specs/001-dotfile-regression-ci/subplans/06-repository-documentation.md`
  - Carried from R01 without change: `subplans/03-dotfile-optional-tool-guards.md`,
    `subplans/05-auto-merge-workflow.md`
  - Read only to check packet 02's removed lines: `.github/workflows/installer-ci.yml`
- **Binding Evidence**: The same candidate and session as R01, kept for the task plan-fidelity role at mid-tier
  by coordinator message. Workspace: feature worktree on branch `001-dotfile-regression-ci` at HEAD `4b58aca`,
  read-only for every reviewed artifact. Invocation: the coordinator's targeted-change request for R02.
  Started and completed 2026-11-18. Result: this round; the only write is this report.
- **Verdict**: passed

#### Current Findings

None. Both R01 concerns are resolved:

- **Repository Documentation Packet Changes the Content of PR C** is resolved:
  - `tasks.md` › Overview now records that `plan.md` assigns no pull request to the repository docs, and that
    the human placed packet 06 in PR C on 2026-11-18.
  - `subplans/06-repository-documentation.md` › Baseline Context and Constraints repeats that placement.
- **Execution Handoff Abridges the Negative Ordering and Sheldon Evidence** is resolved. `tasks.md` ›
  Execution Handoff › Negative and ordering evidence now states three things:
  - An unguarded `sheldon` call fails the bare jobs and `Dotfile Validation Status`.
  - Dotfile validation starts only after `test`, `privilege-test`, and `e2e-test` succeed and uses the `build`
    artifact.
  - A deliberately failing installer test skips dotfile validation and fails `Dotfile Validation Status`.

  This matches `plan.md` › Verification Strategy.

Traced without deviation in this round:

- **Packet 01, `tools-installation.md`:** The doc now documents the `--prepare-only --install-tools`
  interaction from `plan.md` › Add a Prepare-Only Stop to the Installer. This is a proportional installer
  process doc update within the plan's "new flag and stop point" documentation impact. The exact test names
  and `task lint` change only verification mechanics.
- **Packet 02, `changes` failure case:** It follows the plan's "fails ... when any job fails".
- **Packet 02, removed lines:** The checkout step guarded by `github.event_name == 'push'` and the
  paths-filter `base: ${{ github.ref }}` input exist only for the `push` trigger that `plan.md` › Run
  Installer CI on Pull Requests Only removes. The packet keeps the installer paths filter unchanged, so full
  installer E2E stays conditional on installer changes, as `spec.md` › Preserve Conditional Full Installer E2E
  Coverage requires.
- **Packet 04, local Ubuntu with-tools runs:** They add local alternative verification of plan criterion 7. No
  scope or design changes.
- **Earlier trace:** The R01 trace of packets 03 and 05, the delivery order, the outward-facing confirmations,
  the interim release rule, and component documentation review still holds.
- **Scope:** No packet introduces scope, behavior, contracts, or acceptance exceptions beyond `spec.md` and
  `plan.md` R6.
