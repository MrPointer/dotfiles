# Review Dispositions: Dotfile Regression CI

Human dispositions of plan review findings, one section per answered review
batch. Plan reviewers read this file. No task-planning or implementation phase
reads it.

## Batch 01

- **Plan Revision**: R3

Fixes are applied in the next plan Revision, after `spec.md` returns to Ready
(the with-tools disposition below reopened the Feature Definition).

### Post-Apply Chezmoi Checks Are Not Bound to the Checkout Source

- **Role**: plan_design
- **Round**: R01
- **Disposition**: accepted
- **Applied Change**: `plan.md` › Design Decisions › Define Validation Pass Criteria; › Contracts and Flows › Validation flow per matrix job
- **Human Reason**: Agreed as recommended: criteria 2 and 3 use the same `--source` and `--config` as the apply, and an empty managed list fails.

### GitHub API Rate Limits Are Not Addressed for Apply and Release Download

- **Role**: plan_design
- **Round**: R01
- **Disposition**: accepted
- **Applied Change**: `plan.md` › Design Decisions › Choose the Installer Binary by Pull-Request Content; › Contracts and Flows › Validation flow per matrix job; › Alternatives and Risks
- **Human Reason**: Agreed as recommended: authenticate apply, verify, and release download with the job token; list rate limiting as its own risk.

### macOS Optional-Tool Absence Rests on an Unverified Runner Inventory

- **Role**: plan_design
- **Round**: R01
- **Disposition**: accepted
- **Applied Change**: `plan.md` › Verified Current State; › Design Decisions › Ubuntu Runs in the Existing Minimal Test Image; macOS Removes Optional Tools; › Alternatives and Risks
- **Human Reason**: Accepted in a modified form after discussion. The published `macos-latest` (macOS 26 Arm64) and macOS 15 Arm64 image tool lists contain no `tools.yaml` tool and no `diff-so-fancy`, so the removal step is dropped. Each macOS job only asserts that no `tools.yaml` tool and no `diff-so-fancy` resolves on the `PATH` produced by the applied `.zshenv`, and fails naming the tool if one appears after an image change.

### Merged Installer Changes Run Ahead of the Release That Fresh Installs Use

- **Role**: plan_design
- **Round**: R01
- **Disposition**: accepted
- **Applied Change**: `plan.md` › Alternatives and Risks; › Planning Handoff
- **Human Reason**: Release automation and a versioning scheme are justified, but belong to a separate feature. This plan records the risk and the interim rule to tag a release right after any installer merge.

### Startup With Optional Tools Present Is Never Validated

- **Role**: plan_design
- **Round**: R01
- **Disposition**: accepted
- **Applied Change**: `spec.md` returned to Draft to add with-tools validation; then `plan.md` › Design Decisions (validation matrix and pass criteria)
- **Human Reason**: Needed for the reasons the finding gives. Each check splits into bare, with-tools, and per-profile variants, which requires a Feature Definition change because the spec currently excludes optional tools from validation.

### Delivery Pull-Request Boundaries Around the Release Are Unstated

- **Role**: plan_clarity
- **Round**: R01
- **Disposition**: accepted
- **Applied Change**: `plan.md` › Planning Handoff › Order; › Verification Strategy
- **Human Reason**: Agreed: PR A (installer flag and Installer CI changes) → release tag → PR B (dotfile guards and validation workflows) → ruleset changes → PR C (auto-merge workflow); verification cases name their pull request.

### Auto-Merge Window Before the New Check Is Required Is Unexplained

- **Role**: plan_clarity
- **Round**: R01
- **Disposition**: accepted
- **Applied Change**: `plan.md` › Planning Handoff › Order; › Contracts and Flows › Migration order
- **Human Reason**: Agreed: auto-merge lands only after both checks are required, which closes the window.

### Producer of Automated Schema Refresh Changes Is Not Identified

- **Role**: plan_clarity
- **Round**: R01
- **Disposition**: accepted
- **Applied Change**: `plan.md` › Verified Current State; › Design Decisions › Require Pull Requests With No Bypass
- **Human Reason**: Agreed: record that `run_after_refresh-nono-profile-schema.sh.tmpl` writes the schema into the local source and the owner commits it manually; the ruleset change alone meets the requirement.

### Exempt List Names a Root File That Does Not Exist

- **Role**: plan_clarity
- **Round**: R01
- **Disposition**: accepted
- **Applied Change**: `plan.md` › Design Decisions › Exempt Only a Fixed List of Documentation Paths
- **Human Reason**: Agreed: remove `RELEASING.md` from the list.

## Batch 02

- **Plan Revision**: R4

### Criterion 7 Names Packages, Not the Commands They Install

- **Role**: plan_design
- **Round**: R02
- **Disposition**: accepted
- **Applied Change**: `plan.md` › Design Decisions › Define Validation Pass Criteria (criterion 7)
- **Human Reason**: Accepted on the condition that the command names are known; they are taken from `installer/internal/config/packagemap.yaml` and the packages' installed commands.

### First-Start Sheldon Plugin Installation Conflicts With the Empty-Stderr Startup Check

- **Role**: plan_design
- **Round**: R02
- **Disposition**: accepted
- **Applied Change**: `plan.md` › Design Decisions › Define Validation Pass Criteria (criterion 5)
- **Human Reason**: Reasonable: a logged warm-up interactive start before the checked start.

### The Release-Download Path Is Not Proven Before the Rulesets Require It

- **Role**: plan_design
- **Round**: R02
- **Disposition**: accepted
- **Applied Change**: `plan.md` › Planning Handoff › Order; › Verification Strategy
- **Human Reason**: Agreed: a throwaway dotfile-only pull request must pass all eight jobs with the released installer before the ruleset change.

### Ordering Evidence Names a Pull Request That Has No Dotfile Validation

- **Role**: plan_clarity
- **Round**: R02
- **Disposition**: accepted
- **Applied Change**: `plan.md` › Verification Strategy › Ordering evidence
- **Human Reason**: Reasonable: the first installer pull request after PR B provides it.

### Interim Release Rule Is Missing From the Handoff and Its Trigger Is Ambiguous

- **Role**: plan_clarity
- **Round**: R02
- **Disposition**: accepted
- **Applied Change**: `plan.md` › Planning Handoff; › Alternatives and Risks
- **Human Reason**: OK: the rule moves into the handoff, triggered by any merged pull request that changes `installer/**` other than `installer/docs/**` and the installer agent guides.

### No Design Decision Names the With-Tools Requirement It Serves

- **Role**: plan_clarity
- **Round**: R02
- **Disposition**: accepted
- **Applied Change**: `plan.md` › Design Decisions (Serves lines of the prepare-only, pass-criteria, and aggregate-gate decisions)
- **Human Reason**: Citation-only fix applied directly by the planner and reported to the human.

### Trigger Evidence Would Pass Without the Change

- **Role**: plan_clarity
- **Round**: R02
- **Disposition**: accepted
- **Applied Change**: `plan.md` › Verification Strategy › Trigger evidence
- **Human Reason**: OK: evidence becomes the absence of a `push` trigger and no Installer CI run on `main` after PR A merges.

## Batch 03

- **Plan Revision**: R5

### Passing Throwaway Pull Requests Can Auto-Merge After PR C

- **Role**: plan_design
- **Round**: R03
- **Disposition**: accepted
- **Applied Change**: `plan.md` › Verification Strategy; › Planning Handoff
- **Human Reason**: Agreed: every throwaway verification pull request is opened as a draft and closed without merging, so the auto-merge workflow never acts on it regardless of timing.

### Planning Handoff Order Header Miscounts Its Pull Requests

- **Role**: plan_clarity
- **Round**: R03
- **Disposition**: accepted
- **Applied Change**: `plan.md` › Planning Handoff › Order (header wording)
- **Human Reason**: Wording-only fix applied directly by the planner and reported to the human.

### Criterion 7 Command-Name Source Is Not in Source References

- **Role**: plan_clarity
- **Round**: R03
- **Disposition**: accepted
- **Applied Change**: `plan.md` › Source References
- **Human Reason**: Citation-only fix applied directly by the planner and reported to the human.
