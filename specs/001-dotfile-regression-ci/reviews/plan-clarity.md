# Review: Plan Clarity

- **Role**: plan_clarity
- **Tier**: mid-tier
- **Latest Round**: R03
- **Latest Verdict**: passed-with-concerns

## Rounds

### R01

- **Trigger**: initial
- **Scope**: Complete current-state clarity review of `plan.md` Revision R3 against the Ready Feature Definition: current-state evidence and sources, Feature Definition alignment, decision rationale, interfaces and contracts, constraints, risks, verification, planning handoff, and internal consistency of step numbers, job names, check names, file names, and the exempt-list location. Cited sources were spot-checked only where a plan claim's clarity depended on them.
- **Reviewed Paths**: `specs/001-dotfile-regression-ci/spec.md`, `specs/001-dotfile-regression-ci/plan.md`; consulted for cited evidence: `specs/001-dotfile-regression-ci/research.md`, `specs/001-dotfile-regression-ci/grounding-notes.md`, `.github/workflows/installer-ci.yml`, `private_dot_work/` file listing, repository root file listing. No prior report and no `review-dispositions.md` existed.
- **Binding Evidence**: Bound as the plan clarity reviewer at mid-tier, matching the plan_clarity role description. Invoked by the producing phase with the real R01 assignment in the feature worktree for branch `001-dotfile-regression-ci`. Read-only review of the reviewed set; the only write is this report. Started and completed in a single invocation with an attributable result.
- **Verdict**: passed-with-concerns

#### Current Findings

### Delivery Pull-Request Boundaries Around the Release Are Unstated

- **Impact**: concern
- **Location**: `plan.md` › Planning Handoff › Order (steps 1, 4, 5); › Verification Strategy › Workflow evidence and Negative evidence; › Contracts and Flows › Migration order
- **Evidence**: Step 4 tags a release containing step 1 before step 5 adds dotfile validation. A release tag must point at a commit with step 1, and steps 1–3 change `installer/**` or `installer-ci.yml`, so at least one pull request must merge before step 5. Yet Verification refers to "the pull request that delivers this feature" as if there is one. The Negative evidence case for `dot_claude/CLAUDE.md.tmpl` also needs that release to exist. A cold reader cannot tell how many pull requests to create, which steps each holds, or whether the tag is cut from `main` after a merge.
- **Recommendation**: State the pull-request slicing in Planning Handoff. For example: which steps merge to `main` before the step 4 tag, and which pull request carries step 5. Then reword Workflow evidence to name the pull request it means. Tie each Negative evidence case to a point after step 4.

### Auto-Merge Window Before the New Check Is Required Is Unexplained

- **Impact**: concern
- **Location**: `plan.md` › Planning Handoff › Order (steps 6 and 7); › Contracts and Flows › Migration order
- **Evidence**: Step 6 enables auto-merge for every non-draft same-repository pull request. Step 7, which happens later, makes `Dotfile Validation Status` required and adds the `pull_request` rule. In between, auto-merge waits only for `Installer CI Status`, so a pull request whose dotfile validation fails could merge automatically. The plan does not say whether this window is accepted, or whether step 6 must wait for step 7.
- **Recommendation**: Either state the accepted consequence of the window in Alternatives and Risks, or state in Planning Handoff and Migration order that step 6 lands only after step 7. Which one applies is a design choice; the plan needs to say it.

### Producer of Automated Schema Refresh Changes Is Not Identified

- **Impact**: concern
- **Location**: `plan.md` › Design Decisions › Require Pull Requests With No Bypass (Consequence); › Verified Current State › Hosted repository state
- **Evidence**: `spec.md` requires that "an automatically refreshed upstream nono schema" not bypass pull requests. The plan says such changes "must also arrive as pull requests", but it does not name the process that produced commit `71b9472`. It also does not say whether that process needs a change to open a pull request instead of pushing. The only refresh mechanism in the cited evidence is the apply-time script `run_after_refresh-nono-profile-schema.sh.tmpl` (grounding-notes), and the plan never ties the two together. A cold reader cannot tell whether the ruleset change alone meets the requirement, or whether the refresh path needs a task.
- **Recommendation**: In Verified Current State, name how the schema refresh reaches the repository today. Then state in the decision whether that path is unchanged (with the ruleset rejecting direct pushes) or needs a task, and add it to Planning Handoff if it does.

### Exempt List Names a Root File That Does Not Exist

- **Impact**: concern
- **Location**: `plan.md` › Design Decisions › Exempt Only a Fixed List of Documentation Paths (Repository root files)
- **Evidence**: The list includes root `RELEASING.md`, but no `RELEASING.md` exists anywhere in the repository. A cold reader cannot tell whether this entry reserves a planned file or is a mistake. Since the list lives in `installer-ci.yml`, an implementer may copy it as written or drop it.
- **Recommendation**: Remove the entry, or add a short note saying why a not-yet-existing root file is exempt.

### R02

- **Trigger**: targeted change
- **Scope**: Complete current-state clarity review of `plan.md` Revision R4 against the re-Ready Feature Definition. Covers the R4 changes (the eight-job bare/with-tools matrix, `--prepare-only --install-tools`, criterion 7, the chezmoi source/config binding, macOS assert-only, rate-limit and release risks, the three-pull-request delivery order, schema-refresh facts, and removal of `RELEASING.md`). Re-checks every R01 finding against its Batch 01 disposition, and re-checks cross-section consistency of job counts, step and pull-request references, and criterion numbers.
- **Reviewed Paths**: `specs/001-dotfile-regression-ci/spec.md`, `specs/001-dotfile-regression-ci/plan.md`, `specs/001-dotfile-regression-ci/review-dispositions.md`, `specs/001-dotfile-regression-ci/reviews/plan-clarity.md` (prior report).
- **Binding Evidence**: Bound as the plan clarity reviewer at mid-tier, matching the plan_clarity role description. Same-session continuation, invoked by the producing phase with the R02 targeted-change assignment in the feature worktree for branch `001-dotfile-regression-ci`. Read-only review of the reviewed set; the only write is this appended round. Started and completed in this invocation with an attributable result.
- **Verdict**: passed-with-concerns

R01 findings resolved by accepted Batch 01 fixes and not carried forward: Delivery Pull-Request Boundaries Around the Release Are Unstated (Planning Handoff now names PR A, the release tag, PR B, rulesets, and PR C, and verification cases name their pull request); Auto-Merge Window Before the New Check Is Required Is Unexplained (PR C lands after the ruleset step); Producer of Automated Schema Refresh Changes Is Not Identified (Verified Current State › Schema refresh path); Exempt List Names a Root File That Does Not Exist (`RELEASING.md` removed). The eight-job count and the criterion references (flow "2–7", criterion 7 in the macOS decision and the contracts) agree across sections.

#### Current Findings

### Ordering Evidence Names a Pull Request That Has No Dotfile Validation

- **Impact**: concern
- **Location**: `plan.md` › Verification Strategy › Ordering evidence; › Planning Handoff › Order (steps 1 and 3)
- **Evidence**: This follows from the accepted fix "Delivery Pull-Request Boundaries Around the Release Are Unstated". Ordering evidence is assigned to "PR A or a later installer pull request after PR B". But PR A has no dotfile validation job, reusable workflow, or `Dotfile Validation Status`; Planning Handoff adds all of these in PR B. So PR A cannot show that dotfile validation waits for the installer jobs. PR B changes `installer-ci.yml`, so it does trip the installer change filter and could supply this evidence, but the plan does not say so.
- **Recommendation**: Reassign Ordering evidence to PR B, or to an installer pull request after PR B, and remove PR A from that item.

### Interim Release Rule Is Missing From the Handoff and Its Trigger Is Ambiguous

- **Impact**: concern
- **Location**: `plan.md` › Alternatives and Risks › Installer merged ahead of its release; › Planning Handoff › Order and › Outward-facing steps
- **Evidence**: This follows from the accepted fix "Merged Installer Changes Run Ahead of the Release That Fresh Installs Use", whose disposition names Planning Handoff as an applied location. The interim rule "tag a release right after any installer pull request merges" appears only under Risks. Planning Handoff orders only the single release after PR A. "Installer pull request" is also undefined. Elsewhere, an installer change means the filter (`installer/**` or `installer-ci.yml`). PR B changes `installer-ci.yml`, so a reader cannot tell whether PR B's merge also calls for a release, or whether that release would be an outward-facing step needing human confirmation.
- **Recommendation**: State the interim rule in Planning Handoff, and define its trigger: either the installer change filter, or only `installer/**` changes. Then say whether a release is due after PR B, and confirm that every such release follows the outward-facing confirmation rule.

### No Design Decision Names the With-Tools Requirement It Serves

- **Impact**: concern
- **Location**: `plan.md` › Feature Definition Alignment; the "Serves" lines of Design Decisions › Add a Prepare-Only Stop to the Installer, › Ubuntu Runs in the Existing Minimal Test Image; macOS Asserts the Bare Baseline, › Define Validation Pass Criteria, › Gate Merges on One Aggregate Check
- **Evidence**: Feature Definition Alignment promises that each Design Decision names the requirement it serves, and lists "Validate Dotfiles With Optional Tools Present". No "Serves" line names it. The with-tools behavior lives in decisions whose "Serves" lines name other requirements: the `--install-tools` preparation, the with-tools flags, criterion 7, and the eight-job matrix. Only body text cites the new requirement, in two places. Criterion 7's bare half likewise serves "Expose Undeclared Mandatory Dependencies", which that decision's "Serves" line omits. A cold reader tracing requirements to decisions finds nothing listed for the new requirement.
- **Recommendation**: Add "Validate Dotfiles With Optional Tools Present" to the "Serves" lines of the four decisions named above. Add "Expose Undeclared Mandatory Dependencies" to Define Validation Pass Criteria.

### Trigger Evidence Would Pass Without the Change

- **Impact**: concern
- **Location**: `plan.md` › Verification Strategy › Trigger evidence; › Verified Current State › Existing CI
- **Evidence**: Trigger evidence checks that "a push to a branch without a pull request starts no `installer-ci.yml` run". Verified Current State says the existing push trigger fires only on pushes to `main` that touch installer paths. A push to any other branch already starts no run today, so this check cannot tell whether the `push` trigger was removed.
- **Recommendation**: Restate the evidence so it observes the change itself. For example, after PR A merges, the workflow definition on `main` has no `push` trigger, and the PR A merge commit (which touches installer paths) starts no `installer-ci.yml` run.

### R03

- **Trigger**: targeted change
- **Scope**: Complete current-state clarity review of `plan.md` Revision R5 against the Ready Feature Definition. Covers the R5 changes (criterion 5 warm-up start, criterion 7 command-name mapping, the six-step Planning Handoff with the release-path check, the interim release rule and its trigger, and the updated ordering, trigger, and release-path evidence). Re-checks every R02 finding against its Batch 02 disposition, and re-checks cross-references: step numbers, criterion numbers, and job counts.
- **Reviewed Paths**: `specs/001-dotfile-regression-ci/spec.md`, `specs/001-dotfile-regression-ci/plan.md`, `specs/001-dotfile-regression-ci/review-dispositions.md`, `specs/001-dotfile-regression-ci/reviews/plan-clarity.md` (prior report).
- **Binding Evidence**: Bound as the plan clarity reviewer at mid-tier, matching the plan_clarity role description. Same-session continuation, invoked by the producing phase with the R03 targeted-change assignment in the feature worktree for branch `001-dotfile-regression-ci`. Read-only review of the reviewed set; the only write is this appended round. Started and completed in this invocation with an attributable result.
- **Verdict**: passed-with-concerns

All four R02 findings were resolved by accepted Batch 02 fixes and are not carried forward:

- Ordering Evidence Names a Pull Request That Has No Dotfile Validation: it is now assigned to the first installer pull request after PR B.
- Interim Release Rule Is Missing From the Handoff and Its Trigger Is Ambiguous: the rule is now in Planning Handoff with a path-defined trigger, and PR B is explicitly excluded.
- No Design Decision Names the With-Tools Requirement It Serves: the Serves lines were updated as the disposition records.
- Trigger Evidence Would Pass Without the Change: the evidence now observes the absence of the `push` trigger and no run on `main`.

Cross-references agree:

- Outward-facing steps 2 and 5, release-path evidence at step 4, and PR C after step 5 all match the six-step Order.
- The flow's "criteria 2–7" and criterion 5's minimal environment, reused by criteria 6 and 7, are consistent.
- The eight-job count is consistent everywhere.

#### Current Findings

### Planning Handoff Order Header Miscounts Its Pull Requests

- **Impact**: concern
- **Location**: `plan.md` › Planning Handoff › Order (header line and step 4)
- **Evidence**: This follows from the accepted fix "The Release-Download Path Is Not Proven Before the Rulesets Require It". The Order header still reads "three pull requests, each merged to `main` before the next step". Step 4 now adds a fourth, throwaway pull request that is closed without merging. The header and step 4 now contradict each other, and a task decomposer counting pull requests from the header would miss the step 4 pull request.
- **Recommendation**: Reword the header, for example "three merged pull requests (A, B, C), one throwaway pull request at step 4, and two hosted steps", so it agrees with step 4.

### Criterion 7 Command-Name Source Is Not in Source References

- **Impact**: concern
- **Location**: `plan.md` › Design Decisions › Define Validation Pass Criteria (criterion 7); › Source References
- **Evidence**: This follows from the accepted fix "Criterion 7 Names Packages, Not the Commands They Install". Criterion 7 now takes its package-to-command mapping from `installer/internal/config/packagemap.yaml`, but Source References does not list that file. Verified Current State also never states which commands that file maps (for example `batcat` and `fdfind` under apt). A cold reader cannot find the mapping's authority from the reference list, or check whether it still matches the installer.
- **Recommendation**: Add `installer/internal/config/packagemap.yaml` to Source References. Optionally, add a one-line Verified Current State entry citing the lines that map `bat` and `fd` to their apt command names.
