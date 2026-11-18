# Review: Task Plan Fidelity

- **Role**: Task plan-fidelity reviewer (`reviewers/tasks-plan-fidelity.md`)
- **Tier**: mid-tier
- **Latest Round**: R01
- **Latest Verdict**: passed-with-concerns

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
