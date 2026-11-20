# Grounding Notes: Dotfile Regression CI

## Sources

| Source | Observed At | Scope |
|---|---|---|
| `specs/001-dotfile-regression-ci/spec.md:1-84` | 2026-11-16 | Ready feature definition, required validation, exclusions, and decision context |
| `.specify/presets/timors-agentic-workflow/templates/grounding-notes-template.md:1-18` | 2026-11-16 | Required grounding-note structure |
| `.specify/presets/timors-agentic-workflow/references/planning-grounding.md:1-17` | 2026-11-16 | Evidence and decision boundaries for planning grounding |
| `installer/docs/domain.md:25-37,63-68` | 2026-11-16 | Prerequisite and optional-tool meanings and rules |
| `installer/docs/architecture.md:16-42,90-114,175-180` | 2026-11-16 | Installer orchestration and dotfiles-manager boundaries |
| `installer/docs/processes/installation.md:53-89` | 2026-11-16 | Documented installation sequence and outcomes |
| `installer/docs/processes/prerequisite-installation.md:52-69` | 2026-11-16 | Prerequisite installation flow |
| `installer/docs/processes/shell-setup.md:42-71` | 2026-11-16 | Shell installation and default-shell flow |
| `installer/docs/processes/dotfiles-setup.md:51-86,125-138` | 2026-11-16 | Chezmoi installation, initialization, clone, and apply flow |
| `installer/docs/processes/tools-installation.md:54-70,102-113` | 2026-11-16 | Optional-tool selection and non-fatal installation |
| `docs/architecture.md:3-48,50-60` | 2026-11-16 | Installer, chezmoi, and shell-runtime integration boundaries |
| `docs/dotfiles/domain.md:3-43,65-90` | 2026-11-16 | Personal/work rendering and shell-runtime concepts |
| `docs/contracts/chezmoi-data.md:1-70` | 2026-11-16 | Canonical installer-to-template data contract |
| `installer/cmd/install.go:65-134,586-655,658-762,769-801` | 2026-11-16 | Actual orchestration, data inputs, optional tools, and install flags |
| `installer/lib/dotfilesmanager/data.go:7-25` and `installer/lib/dotfilesmanager/chezmoi/data.go:12-70` | 2026-11-16 | Data types and emitted chezmoi keys |
| `installer/lib/dotfilesmanager/chezmoi/apply.go:9-48` | 2026-11-16 | Current source removal, remote initialization, branch selection, and apply |
| `installer/internal/config/tools.yaml:1-17` | 2026-11-16 | Configured optional tools |
| `dot_zshenv.tmpl:1-75`, `dot_zshrc.tmpl:1-256`, and `private_dot_work/private_sedg/zsh/dot_zshenv.tmpl:1-64` | 2026-11-16 | Personal and existing work shell-startup behavior |
| `.chezmoiexternal.toml.tmpl:1-47` | 2026-11-16 | Network-backed files used during chezmoi application |
| `run_onchange_after_install-plan-html.sh.tmpl:1-26`, `run_after_refresh-nono-profile-schema.sh.tmpl:1-56`, and `run_after_ensure-codex-hooks.sh:1-23` | 2026-11-16 | Presence guards in apply-time scripts |
| `.github/workflows/installer-ci.yml:3-45,151-177,251-289,340-434` | 2026-11-16 | Workflow triggers, change filtering, E2E checkout behavior, shell check, and aggregate status |
| `gh api repos/MrPointer/dotfiles` | 2026-11-16 | Default branch, merge methods, automatic merge setting, and current actor permissions |
| `gh api repos/MrPointer/dotfiles/rulesets` | 2026-11-16 | Active repository rulesets and identifiers |
| `gh api repos/MrPointer/dotfiles/rulesets/8895424` | 2026-11-16 | `protect-main` target, bypass, and branch rules |
| `gh api repos/MrPointer/dotfiles/rulesets/12783037` | 2026-11-16 | `require-installer-ci` target, bypass, and required status check |

## Observations

### Feature Baseline Requires Pre-Merge Validation of the Applied Revision

- **Evidence**: The Ready feature definition requires relevant pull requests to validate personal and existing work configurations on Ubuntu and macOS, use installer preparation without applying dotfiles, apply the exact pull-request target checkout or commit, fully apply into a clean disposable home, and exercise interactive and noninteractive shell startup. Optional tools remain absent, while the existing full installer E2E suite remains conditional on installer changes (`specs/001-dotfile-regression-ci/spec.md:10-32,40-66`).
- **Bounded Inference**: Planning evidence must distinguish preparation coverage, dotfile application coverage, and shell-startup coverage because the feature defines each as a separate required outcome.
- **Relevance**: These statements are the normative baseline against which current implementation gaps are grounded.
- **Unknowns**: The feature does not choose whether the exact revision is the pull-request head SHA or the generated merge SHA.

### Current Installer Orchestration Has No Pre-Apply Stop

- **Evidence**: The install command performs system detection and compatibility/prerequisite handling, optional Homebrew setup, shell installation and default-shell configuration, GPG setup, `setupDotfilesManager`, and then optional-tool installation (`installer/cmd/install.go:65-134`). `setupDotfilesManager` calls `Install`, initializes data, and calls `Apply` as one uninterrupted function (`installer/cmd/install.go:586-618`). The registered install flags include work mode, shell and package choices, clone protocol and branch, prerequisite installation, and optional-tool installation, but no preparation or pre-apply stopping flag (`installer/cmd/install.go:769-801`).
- **Bounded Inference**: The existing orchestration cannot currently satisfy the specified “prepare, then let CI apply” boundary through a public install option.
- **Relevance**: The stopping capability must fit the actual sequence without treating optional-tool installation as a prerequisite phase; optional tools currently run after dotfiles application.
- **Unknowns**: Inspection did not establish orchestration tests that prove an early stop before `Apply`; no runtime validation was performed.

### Chezmoi Apply Clones a Remote Branch Rather Than Applying an Exact Existing Checkout

- **Evidence**: `ChezmoiManager.Apply` removes the configured source directory, builds `chezmoi init --apply`, optionally adds `--branch`, identifies the remote account, and supplies the config path (`installer/lib/dotfilesmanager/chezmoi/apply.go:9-48`). The documented state change is deletion and re-cloning of the source directory (`installer/docs/processes/dotfiles-setup.md:75-86,125-130`).
- **Bounded Inference**: Supplying a branch name selects a remote branch at clone time; it does not pin the apply to a previously checked-out commit.
- **Relevance**: CI cannot use the current apply path alone as evidence that the exact checked-out pull-request revision was applied.
- **Unknowns**: The exact head-SHA versus merge-SHA policy remains undecided.

### Personal and Work Data Already Share One Explicit Producer Contract

- **Evidence**: The installer always builds hardcoded personal identity and shell data, adds work identity and `~/.work/profile` plus `~/.work/{work_name}/profile` when `--work-env` is enabled, and adds a selected GPG key when present (`installer/cmd/install.go:621-655`). `ChezmoiManager.Initialize` writes the corresponding `data.personal`, optional `data.gpg`, and `data.system` keys (`installer/lib/dotfilesmanager/chezmoi/data.go:12-70`). The canonical contract defines presence rules and current consumers for personal and work renders (`docs/contracts/chezmoi-data.md:26-70`).
- **Bounded Inference**: Existing personal and work validation can be grounded in the same producer with different flag-selected data; this does not imply any redesign of the work configuration.
- **Relevance**: The current work configuration named in the feature is already represented by installer inputs and conditional template data.
- **Unknowns**: Runtime application of both modes on clean hosted environments was not performed during grounding.

### Optional-Tool Intent and Shell Startup Currently Conflict for Sheldon

- **Evidence**: The installer domain defines optional tools as not required to apply dotfiles (`installer/docs/domain.md:25-37`), and `tools.yaml` lists Sheldon as optional (`installer/internal/config/tools.yaml:6-7`). The main interactive shell template invokes `eval "$(sheldon source)"` without a command-presence guard (`dot_zshrc.tmpl:182-183`). Later fzf behavior checks whether `fzf` exists before setup (`dot_zshrc.tmpl:197-217`); the following fzf-tab `zstyle` settings only set shell configuration (`dot_zshrc.tmpl:219-241`).
- **Bounded Inference**: In the inspected continuation of the main shell template, no later command directly requires plugins loaded by Sheldon. This is limited to static inspection of the cited startup path.
- **Relevance**: A baseline that intentionally omits optional tools can expose the unconditional Sheldon invocation during interactive startup.
- **Unknowns**: Whether Sheldon will be guarded to remain optional or promoted to a prerequisite is an unresolved user decision. The optional-tool inventory actually present on hosted environments was not established.

### Startup and Apply Coverage Extends Beyond the Unconditional Sheldon Call

- **Evidence**: `.zshenv` sets common paths, conditionally loads Homebrew, and conditionally sources readable work-profile extensions (`dot_zshenv.tmpl:1-75`). The existing work `.zshenv` extension includes path setup, directory checks, conditional framework sources, optional Cargo environment sourcing, and devbox-only environment-file discovery (`private_dot_work/private_sedg/zsh/dot_zshenv.tmpl:1-64`). The main `.zshrc` also conditionally sources work and sandbox-tool files, checks local plugin files, sources the generated NVM lazy-loader path, and guards several tool integrations (`dot_zshrc.tmpl:118-156,185-217,243-256`).
- **Bounded Inference**: Conditional sources are not automatically regressions; their conditions and the clean-home result determine whether they execute successfully.
- **Relevance**: Both interactive and noninteractive startup are needed to cover the distinct `.zshrc` and `.zshenv` paths, including the existing work configuration.
- **Unknowns**: Static inspection does not establish behavior for every hosted environment variable, filesystem condition, or external work dependency.

### A Full Apply Includes Network Externals and Guarded Scripts

- **Evidence**: `.chezmoiexternal.toml.tmpl` obtains NVM release metadata and downloads NVM, shell snippets, plugins, and other files from remote URLs (`.chezmoiexternal.toml.tmpl:1-47`). Apply-time scripts exit successfully when `uv`, `nono`, or `codex` is absent (`run_onchange_after_install-plan-html.sh.tmpl:16-18`, `run_after_refresh-nono-profile-schema.sh.tmpl:8-10`, `run_after_ensure-codex-hooks.sh:16-18`).
- **Bounded Inference**: A real clean-home apply exercises more than template rendering and file copying; it can also depend on network externals, while the cited scripts preserve absence for their primary optional commands.
- **Relevance**: Application success must not be inferred solely from rendered shell files, and guarded scripts should not be treated as evidence that their optional commands belong in the baseline.
- **Unknowns**: Network behavior and all run-script outcomes were not executed during grounding. An earlier generic installed-template resolver failed because PyYAML was unavailable; that tooling issue was unrelated to direct consumption of the installed preset template used for these notes.

### Existing CI Routes Pull Requests but Filters Expensive Installer Work

- **Evidence**: `installer-ci.yml` triggers on every pull request, while pushes to `main` are path-filtered to installer files and the workflow itself (`.github/workflows/installer-ci.yml:3-10`). Its change detector marks only those installer paths, and the build plus downstream installer jobs depend on that result (`.github/workflows/installer-ci.yml:18-45,93-153`).
- **Bounded Inference**: A dotfile-only pull request starts the workflow but skips the expensive installer build, test, privilege, and E2E chain.
- **Relevance**: This matches the requirement to preserve conditional full installer E2E coverage, but it does not provide the new dotfile validation.
- **Unknowns**: No CI run was dispatched during grounding.

### Current Pull-Request Checkout and Installer E2E Clone Are Different Revisions

- **Evidence**: Pull-request jobs use the checkout action without an explicit `ref`; on a pull-request event this checks out the event’s generated merge ref (`.github/workflows/installer-ci.yml:50-53,100-102,134-135,175-177`). The noninteractive E2E step instead derives `GITHUB_HEAD_REF` and passes that branch name to `--git-branch`, after which the installer removes its source and clones that remote branch (`.github/workflows/installer-ci.yml:251-289`; `installer/lib/dotfilesmanager/chezmoi/apply.go:9-29`).
- **Bounded Inference**: The normal job checkout is an event-selected merge commit, not an unpinned moving branch; the subsequent installer clone is branch-based and can move independently.
- **Relevance**: The distinction prevents the existing checkout from being misclassified and isolates the exact-revision problem to the installer-driven clone/apply path.
- **Unknowns**: The feature has not selected head SHA or generated merge SHA as the exact validation target.

### Existing Shell Evidence Can Pass Without Applied Files

- **Evidence**: The E2E artifact check reports missing chezmoi directories as expected (`.github/workflows/installer-ci.yml:340-348`). Its shell check sources only `.zshenv`; when `.zshenv` is absent it still runs a bare `zsh` command and passes (`.github/workflows/installer-ci.yml:387-407`). It does not start an interactive shell or source `.zshrc`.
- **Bounded Inference**: The current check does not prove a complete clean-home application or interactive startup.
- **Relevance**: This directly explains why existing E2E shell evidence does not satisfy the feature’s acceptance evidence.
- **Unknowns**: No current runtime result established which startup failure appears first under the intended minimal baseline.

### Workflow Token Permissions Are Already Read-Limited

- **Evidence**: The workflow declares `permissions: {}` at top level, and jobs that need repository data declare `contents: read` and, for change detection, `pull-requests: read` (`.github/workflows/installer-ci.yml:16-24,42-49,93-99,126-134`). No workflow-level write permission is declared.
- **Bounded Inference**: Repository default workflow permissions do not grant write tokens to this workflow execution because the workflow explicitly resets permissions and adds only cited per-job reads.
- **Relevance**: Hosted merge governance must be assessed separately from workflow token scope; a write-token risk is not present in the current file.
- **Unknowns**: Future automatic-merge execution ownership is not defined by the current workflow.

### Aggregate Status Accepts Skips and Most Cancellations

- **Evidence**: `Installer CI Status` always runs and fails for job failures plus cancellation of `privilege-test`, but it does not fail for cancellation of `build`, `test`, or `e2e-test`; its success message explicitly accepts skipped or canceled jobs (`.github/workflows/installer-ci.yml:409-434`).
- **Bounded Inference**: The aggregate status is designed to pass installer jobs skipped by the path filter and can also pass several cancellation states.
- **Relevance**: The hosted required check currently names this aggregate status, so its exact pass/fail semantics determine what repository rules enforce.
- **Unknowns**: No hosted check-run history was inspected to determine whether those cancellation cases have occurred.

### Hosted Rules Support Automatic Merge but Do Not Enforce Pull Requests for Every Change

- **Evidence**: `gh api repos/MrPointer/dotfiles` reported default branch `main`, `allow_auto_merge=true`, squash and rebase enabled, and `permissions.admin=true` for the current actor. `gh api repos/MrPointer/dotfiles/rulesets` reported active `protect-main` ruleset `8895424` and `require-installer-ci` ruleset `12783037`. Detail queries reported both target `~DEFAULT_BRANCH`, no excludes, and `bypass_actors=[{actor_id:5,actor_type:"RepositoryRole",bypass_mode:"always"}]`. `protect-main` contains deletion, non-fast-forward, required-linear-history, and required-signatures rules. `require-installer-ci` requires status context `Installer CI Status` with `strict_required_status_checks_policy=false` and `do_not_enforce_on_create=true`. Neither detail response contains a `pull_request` rule.
- **Bounded Inference**: Hosted settings permit automatic merge, but the active rulesets do not by themselves require every default-branch change to arrive through a pull request, and their always-bypass role leaves a bypass path.
- **Relevance**: The repository’s current governance does not yet provide acceptance evidence for “every repository change must go through a pull request.”
- **Unknowns**: Grounding did not mutate hosted settings or test bypass behavior.
