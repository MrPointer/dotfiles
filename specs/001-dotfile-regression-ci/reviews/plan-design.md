# Review: Plan Design

- **Role**: plan design reviewer
- **Tier**: most-capable
- **Latest Round**: R04
- **Latest Verdict**: passed

## Rounds

### R01

- **Trigger**: initial
- **Scope**: Complete current-state design review of `plan.md` Revision R3 against the Ready Feature Definition: verified current state, architecture, contracts, flows, state, failure behavior, migration order, risks, verification, and whether durable decisions live in `plan.md`.
- **Reviewed Paths**: `specs/001-dotfile-regression-ci/spec.md`, `specs/001-dotfile-regression-ci/plan.md`, `specs/001-dotfile-regression-ci/grounding-notes.md`, `specs/001-dotfile-regression-ci/research.md`; current-state checks against `installer/cmd/install.go`, `installer/lib/dotfilesmanager/chezmoi/apply.go`, `installer/internal/config/tools.yaml`, `installer/docker/ubuntu/Dockerfile`, `.github/workflows/installer-ci.yml`, `.github/workflows/release.yml`, `get.sh`, `.chezmoi.toml.tmpl`, `.chezmoiignore`, `.chezmoiexternal.toml.tmpl`, `dot_zshrc.tmpl`, `private_dot_config/sheldon/plugins.toml`, `private_dot_config/git/common.gitconfig.tmpl`.
- **Binding Evidence**: Plan design reviewer, most-capable tier, single clear role match. Invoked by the coordinator as the real R01 assignment in the feature worktree on branch `001-dotfile-regression-ci`. Read-only access to the reviewed set; the only write is this report. Repository claims were re-verified from files. Hosted repository settings and rulesets were not re-queried (no shell access in this session); hosted claims rely on the dated `gh api` evidence in `grounding-notes.md`, and no finding depends on them. Result: complete round appended.
- **Verdict**: passed-with-concerns

Verified without finding: installer sequence and missing stop point (`install.go:65-134,586-619,769-801`); `Apply` deletes the clone and runs `chezmoi init --apply` against the remote (`apply.go:9-48`); `.chezmoi.toml.tmpl` re-renders data only; unguarded `sheldon source` (`dot_zshrc.tmpl:183`); `core.pager` and `diff.external` tool references (`common.gitconfig.tmpl:19,22`); `diff-so-fancy` reaches `PATH` only through a Sheldon plugin (`plugins.toml:65-66`); `Installer CI Status` cancellation gaps (`installer-ci.yml:424-434`); E2E branch clone and weak shell check; `get.sh` latest-release download. All durable decisions in `research.md` are carried into `plan.md`.

#### Current Findings

### Post-Apply Chezmoi Checks Are Not Bound to the Checkout Source

- **Impact**: concern
- **Location**: `plan.md` › Design Decisions › Define Validation Pass Criteria (criteria 2 and 3); › Contracts and Flows › Validation flow per matrix job
- **Evidence**: Criterion 1 passes `--source <checkout> --config <path>`. Criteria 2 (`chezmoi verify --exclude scripts`) and 3 (`chezmoi managed ...`) do not name a source or config. The config that `chezmoi init` writes comes from `.chezmoi.toml.tmpl`, which holds only `data` and `[diff]`, with no `sourceDir`. `--prepare-only` creates no clone at `~/.local/share/chezmoi`. If criteria 2 and 3 run without `--source`, chezmoi reads an empty source state. `verify` then has nothing to compare, and `managed` lists no paths, so both pass trivially. The spec's acceptance evidence says validation must not pass because expected files are missing.
- **Recommendation**: In the pass criteria, require criteria 2 and 3 to use the same `--source <checkout>` and `--config <path>` as criterion 1. Add that criterion 3 fails when the managed list is empty.

### GitHub API Rate Limits Are Not Addressed for Apply and Release Download

- **Impact**: concern
- **Location**: `plan.md` › Alternatives and Risks › Network externals; › Choose the Installer Binary by Pull-Request Content; › Contracts and Flows › Validation flow per matrix job
- **Evidence**: `.chezmoiexternal.toml.tmpl:1` calls `gitHubLatestRelease "nvm-sh/nvm"`. The plan also calls the GitHub API during `chezmoi verify` re-rendering and through the latest-release lookup, which `get.sh:66-71` makes unauthenticated. The existing E2E job passes `CHEZMOI_GITHUB_ACCESS_TOKEN` for this reason (`installer-ci.yml:252-253,307-308`). Hosted runners share egress IPs, and the unauthenticated limit is low. The plan treats network failure as a rare outage that someone retries manually. With four matrix jobs on every non-exempt pull request and auto-merge, a predictable rate-limit failure blocks merges with no human watching.
- **Recommendation**: Have the plan state that the apply, verify, and release-lookup steps authenticate GitHub API calls with the job token, for example `CHEZMOI_GITHUB_ACCESS_TOKEN` and an authenticated release query. Update the Network externals risk to separate rate limiting from upstream outages.

### macOS Optional-Tool Absence Rests on an Unverified Runner Inventory

- **Impact**: concern
- **Location**: `plan.md` › Design Decisions › Ubuntu Runs in the Existing Minimal Test Image; macOS Removes Optional Tools; › Define Validation Pass Criteria (criterion 6); › Verified Current State
- **Evidence**: The macOS job uninstalls `tools.yaml` tools through Homebrew, then fails if any still resolves on `PATH`. Three things are unverified:
  - **Runner inventory.** `grounding-notes.md` (Optional-Tool Intent observation, Unknowns) says the tools actually present on hosted runners were not established.
  - **Non-Homebrew installs.** A tool installed outside Homebrew, such as `uv` from a standalone installer, cannot be removed this way. The job would then fail permanently on every pull request.
  - **Which `PATH`.** The plan does not say which `PATH` the assertion uses. The validated shell builds its own `PATH` from `.zshenv` (criterion 5), not from the runner's.

  Criterion 6 also requires `diff-so-fancy` to be absent, but `diff-so-fancy` is not in `tools.yaml`, so nothing removes it or asserts that it is gone.
- **Recommendation**:
  - Record the macOS runner's actual optional-tool inventory and install sources in Verified Current State, or name the gap as a risk.
  - State how a tool installed outside Homebrew is handled.
  - Run the absence assertion against the `PATH` that the applied `.zshenv` produces.
  - Add `diff-so-fancy` to the asserted-absent set.

### Merged Installer Changes Run Ahead of the Release That Fresh Installs Use

- **Impact**: concern
- **Location**: `plan.md` › Design Decisions › Choose the Installer Binary by Pull-Request Content; › Alternatives and Risks › Release lag
- **Evidence**: This concerns the human decision of 2026-11-17 and does not call for a redesign. An installer pull request is validated with its own `build` artifact. After it merges, `get.sh` still serves the previous release until someone tags a new one. Suppose such a pull request changes the data written for chezmoi (`docs/contracts/chezmoi-data.md`) and dotfiles in the same pull request consume the change. Then `main` passes validation but breaks fresh installs, which is the regression this feature exists to prevent (`spec.md` › Intent). Every later non-installer pull request is validated with the old release and fails until a release is cut. The Release lag risk covers only the reverse case.
- **Recommendation**: Present this to the human. If accepted, add a risk entry naming this case and an operating rule. Options: tag a release promptly after any installer change that affects dotfile consumers, or require dotfile templates to tolerate the data keys of the latest release.

### Startup With Optional Tools Present Is Never Validated

- **Impact**: concern
- **Location**: `plan.md` › Alternatives and Risks; › Guard Optional Tools in Dotfiles Instead of Promoting Them
- **Evidence**: With Sheldon guarded and absent from every validation job, no job runs the Sheldon path. That path covers `sheldon source`, `private_dot_config/sheldon/plugins.toml` (including powerlevel10k and diff-so-fancy), and the difftastic path. The usual real install adds optional tools after apply, so a broken plugin entry or a broken tool-present branch of a guard can still reach `main` undetected. The spec requires a baseline without optional tools, so this gap is a consequence of the spec, not a plan defect. However, the plan's risk list does not record it.
- **Recommendation**: Add a risk entry stating that validation covers only the startup path without optional tools, so regressions in the path with tools present remain undetected. Leave any added coverage to a later human decision.

### R02

- **Trigger**: targeted change
- **Scope**: Complete current-state design review of `plan.md` Revision R4 against the re-Ready Feature Definition. The new requirement "Validate Dotfiles With Optional Tools Present" adds an eight-job matrix, `--prepare-only --install-tools` preparation, and pass criterion 7. The release-automation exclusion was also checked. The round checks how each Batch 01 disposition was applied and carries forward any unaffected R01 finding.
- **Reviewed Paths**: `specs/001-dotfile-regression-ci/spec.md`, `specs/001-dotfile-regression-ci/plan.md`, `specs/001-dotfile-regression-ci/research.md`, `specs/001-dotfile-regression-ci/grounding-notes.md`, `specs/001-dotfile-regression-ci/review-dispositions.md`, `specs/001-dotfile-regression-ci/reviews/plan-design.md` (R01). Current-state checks against `installer/cmd/install.go` (package-manager selection and `installOptionalTools`), `installer/internal/config/tools.yaml`, `installer/internal/config/packagemap.yaml`, `private_dot_config/git/common.gitconfig.tmpl`, `private_dot_config/sheldon/plugins.toml`, `dot_zshrc.tmpl`, and `docs/dotfiles/processes/shell-startup.md`.
- **Binding Evidence**: Plan design reviewer, most-capable tier, single clear role match. The same reviewer session continued with the coordinator's R02 assignment in the feature worktree on branch `001-dotfile-regression-ci`. Access was read-only on the reviewed set; the only write is this report. Hosted settings and the published macOS runner image list were not re-queried (no shell or network access in this session). The macOS inventory claim rests on the human-accepted disposition and on `plan.md` › Verified Current State. Result: complete round appended.
- **Verdict**: passed-with-concerns

Disposition check against `review-dispositions.md` Batch 01. All five R01 findings were accepted, and each fix is present in R4:

- **Source binding.** Criteria 2–3 now share `--source`/`--config`, and an empty managed list fails.
- **Rate limits.** Apply, verify, and the release download use the job token, and rate limiting is now its own risk.
- **macOS baseline.** The approach is now assert-only on the applied interactive `PATH`, and `diff-so-fancy` is covered.
- **Installer ahead of release.** A risk entry and an interim tagging rule were added. Release automation is excluded in `spec.md`.
- **Startup with tools present.** The with-tools variant is now in `spec.md` and `plan.md`.

No R01 finding is carried forward. None of the findings below re-raises a declined finding. The first two arise from the accepted with-tools fix (from "Startup With Optional Tools Present Is Never Validated"). The third arises from the accepted delivery-order fix (from the clarity finding "Delivery Pull-Request Boundaries Around the Release Are Unstated").

Verified without finding:

- **Brew-backed tool installs.** With the default `--install-brew=true`, `createPackageManagerForSystem` returns the Homebrew manager on both platforms (`install.go:161-165`). `packagemap.yaml` maps all eight `tools.yaml` tools for brew. `installOptionalTools` silently skips tools it cannot resolve (`install.go:687-702`). Criterion 7's presence check therefore fails closed, as `spec.md` requires.
- **Durable decisions.** These are in `plan.md`. `research.md` still has stale lines: `RELEASING.md` in the exempt list, "four matrix entries", and "optional tools are also skipped". This is a supporting-document consistency matter for the clarity reviewer, not a design gap.

#### Current Findings

### Criterion 7 Names Packages, Not the Commands They Install

- **Impact**: concern
- **Location**: `plan.md` › Design Decisions › Define Validation Pass Criteria (criterion 7)
- **Evidence**: Criterion 7 checks "every `tools.yaml` tool" with `command -v`. `tools.yaml` holds package names, and two of them differ from the command they put on `PATH`:
  - `ripgrep` installs `rg`.
  - `difftastic` installs `difft`, which is the command `common.gitconfig.tmpl:22` and the planned guard look for.

  Read literally, the with-tools jobs fail on every run, because `ripgrep` and `difftastic` never resolve. The bare jobs could also pass with `rg` or `difft` present, because the absence check looks for the wrong names. That would weaken the "Expose Undeclared Mandatory Dependencies" check.
- **Recommendation**: In criterion 7, state that each `tools.yaml` entry is checked by the command it provides, and name the mapping explicitly (at least `ripgrep`→`rg`, `difftastic`→`difft`, with the other entries checked by their own names). Also state where that mapping is maintained, so that adding a tool to `tools.yaml` updates both checks.

### First-Start Sheldon Plugin Installation Conflicts With the Empty-Stderr Startup Check

- **Impact**: concern
- **Location**: `plan.md` › Design Decisions › Define Validation Pass Criteria (criterion 5); › Alternatives and Risks › Interactive startup without a terminal
- **Evidence**:
  - In with-tools jobs, Sheldon is present, so the guarded `eval "$(sheldon source)"` (`dot_zshrc.tmpl:183`) runs.
  - On a clean home this is Sheldon's first run, so it locks and clones every plugin in `private_dot_config/sheldon/plugins.toml` (including powerlevel10k and diff-so-fancy) during the `zsh -i -c exit` of criterion 5.
  - Sheldon normally reports lock and clone progress on stderr. The plan cites no evidence of what it emits, and the startup line passes no quiet option.
  - If it writes to stderr, criterion 5's "empty stderr" fails every with-tools job on first startup. Real users also see this output, so it is not a regression the check should catch.
  - The plan's risk entry covers only tools that expect a terminal, not first-run plugin installation.
- **Recommendation**: Verify what first-run `sheldon source` writes to stderr. Then have the plan state how first-start plugin installation relates to criterion 5. For example, the empty-stderr check could apply to a startup after plugins are installed, while a failure of the first startup still fails the job. Record the choice as a risk or criterion note.

### Release-Download Validation Path Is Not Proven Before the Rulesets Require It

- **Impact**: concern
- **Location**: `plan.md` › Planning Handoff › Order (steps 3–4); › Contracts and Flows › Migration order; › Verification Strategy › Negative evidence
- **Evidence**:
  - PR B changes `.github/workflows/installer-ci.yml`, so its installer change filter is true and its own validation uses the `build` artifact.
  - The authenticated latest-release download path is first exercised by the throwaway `dot_claude/CLAUDE.md.tmpl` pull request "after PR B merges".
  - The Handoff orders the ruleset changes (step 4) after PR B merges, but not after that throwaway run.
  - If the release-download path is broken, every non-installer pull request is blocked once `Dotfile Validation Status` is required with no bypass. The only repair is a pull request that touches `installer-ci.yml`, or a deliberate ruleset edit.
  - The Migration order sentence ("has passed on its own pull request") does not cover this path.
- **Recommendation**: In Planning Handoff › Order and Migration order, require that before step 4, a non-installer pull request passes all eight matrix jobs using the released installer with `build` skipped. The existing `dot_claude/CLAUDE.md.tmpl` negative-evidence case can serve as that pull request.

### R03

- **Trigger**: targeted change
- **Scope**: Complete current-state design review of `plan.md` Revision R5. Changed areas:
  - criterion 5 warm-up interactive start;
  - criterion 7 package-to-command mapping;
  - Planning Handoff step 4 release-path check, with rulesets now step 5 and PR C step 6;
  - the interim release rule, moved into the handoff with a defined trigger;
  - updated ordering, trigger, and release-path evidence;
  - Serves citations.

  The round checks how each Batch 02 disposition was applied, and re-checks the unchanged sections for their current state.
- **Reviewed Paths**: `specs/001-dotfile-regression-ci/spec.md`, `specs/001-dotfile-regression-ci/plan.md`, `specs/001-dotfile-regression-ci/research.md`, `specs/001-dotfile-regression-ci/grounding-notes.md`, `specs/001-dotfile-regression-ci/review-dispositions.md`, `specs/001-dotfile-regression-ci/reviews/plan-design.md` (R01–R02). Current-state checks against `installer/internal/config/packagemap.yaml` and `installer/cmd/install.go` (carried from R02; unchanged).
- **Binding Evidence**: Plan design reviewer, most-capable tier, single clear role match. The same reviewer session continued with the coordinator's R03 assignment in the feature worktree on branch `001-dotfile-regression-ci`. Access was read-only on the reviewed set; the only write is this report. Hosted settings were not re-queried (no shell or network access). Result: complete round appended.
- **Verdict**: passed-with-concerns

Disposition check against `review-dispositions.md` Batch 02. All three R02 plan_design findings were accepted, and each fix is present in R5:

- **Command names.** Criterion 7 now checks commands, not package names. It maps `ripgrep`→`rg` and `difftastic`→`difft`, accepts the apt names `batcat` and `fdfind`, and adds `diff-so-fancy`. The mapping is sourced from `packagemap.yaml`.
- **Sheldon first run.** Criterion 5 runs a logged warm-up interactive start that must exit 0 but whose stderr is not judged. The noninteractive and second interactive starts must have empty stderr.
- **Release path.** Handoff step 4 requires the throwaway dotfile-only pull request to pass all eight jobs with the released installer before step 5 changes the rulesets. Migration order and Release-path evidence agree with this.

The relocated interim release rule has a clear trigger consistent with the exempt list. The Trigger evidence now proves the change. No R02 finding is carried forward. The finding below arises from the accepted auto-merge ordering (from the clarity finding "Auto-Merge Window Before the New Check Is Required Is Unexplained") combined with the unchanged Negative evidence.

#### Current Findings

### Passing Throwaway Pull Requests Can Auto-Merge After PR C

- **Impact**: concern
- **Location**: `plan.md` › Verification Strategy › Negative evidence; › Planning Handoff › Order (step 6)
- **Evidence**: Negative evidence says only that its throwaway pull requests run "after PR B merges". Two of them are expected to pass:
  - the exempt-only pull request;
  - the `dot_claude/CLAUDE.md.tmpl`-only pull request.

  If either is opened after PR C (step 6) as a non-draft pull request, `auto-merge.yml` enables auto-merge. GitHub then squash-merges the throwaway change into `main` once both checks pass. Only step 4's throwaway pull request is explicitly ordered before PR C and closed without merging.
- **Recommendation**: In Negative evidence, require that these throwaway pull requests run before step 6, or be opened as drafts and closed without merging. Drafts still get both checks but are the auto-merge opt-out.

### R04

- **Trigger**: targeted change
- **Scope**: Complete current-state design review of `plan.md` Revision R6. Changed areas:
  - the new Planning Handoff › Throwaway pull requests rule, under which every verification-only pull request is a draft and is closed without merging;
  - draft marking on step 4 and on Negative evidence;
  - two direct wording and citation fixes: the handoff Order header count, and `packagemap.yaml` in Source References.

  The round checks how the Batch 03 disposition was applied. The unchanged sections keep their R03 current state, which had no other open finding.
- **Reviewed Paths**: `specs/001-dotfile-regression-ci/plan.md`, `specs/001-dotfile-regression-ci/review-dispositions.md`, `specs/001-dotfile-regression-ci/reviews/plan-design.md` (R01–R03). `spec.md`, `research.md`, and `grounding-notes.md` are unchanged since R03 per the assignment, and their R03 assessment stands.
- **Binding Evidence**: Plan design reviewer, most-capable tier, single clear role match. The same reviewer session continued with the coordinator's R04 assignment in the feature worktree on branch `001-dotfile-regression-ci`. Access was read-only on the reviewed set; the only write is this report. Hosted settings were not re-queried (no shell or network access). Result: complete round appended.
- **Verdict**: passed

Disposition check against `review-dispositions.md` Batch 03. The R03 finding "Passing Throwaway Pull Requests Can Auto-Merge After PR C" was accepted, and the fix is present in R6:

- The Planning Handoff › Throwaway pull requests rule makes every verification-only pull request a draft that is closed without merging.
- Step 4 and Negative evidence both reference this rule.
- This agrees with the auto-merge workflow's draft opt-out (Design Decisions › Enable Auto-Merge From a Workflow), so timing no longer matters.

The two direct fixes change no decision, contract, flow, or evidence. No finding is carried forward, and no new design finding was found.

#### Current Findings

None.
