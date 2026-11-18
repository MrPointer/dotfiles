# Plan: Dotfile Regression CI

- **Revision**: R6

## Feature Definition Alignment

This plan implements `spec.md` › Requirements:

- Route Every Repository Change Through a Pull Request
- Validate Personal and Existing Work Configurations on Ubuntu and macOS Before Merge
- Prepare Without Applying Dotfiles
- Validate a Real Clean-Home Application
- Expose Undeclared Mandatory Dependencies
- Validate Dotfiles With Optional Tools Present
- Prioritize Safety Over Validation Speed
- Preserve Conditional Full Installer E2E Coverage

It honors `spec.md` › Scope › Excluded and › Constraints and Dependencies. Each
Design Decision below names the requirement it serves.

## Verified Current State

Observed 2026-11-16 and 2026-11-17. Full evidence is in `grounding-notes.md` and
`research.md`.

- **Installer sequence.** The `install` command runs compatibility and
  prerequisite handling, optional Homebrew setup, shell installation, GPG setup,
  `setupDotfilesManager`, then `installOptionalTools`
  (`installer/cmd/install.go:65-134`). `setupDotfilesManager` calls the chezmoi
  manager's `Install`, then `initDotfilesManagerData` (which calls `Initialize`),
  then `Apply`, with no stop point (`installer/cmd/install.go:586-618`). No flag
  stops before `Apply` (`installer/cmd/install.go:769-801`).
- **Chezmoi paths and apply.** The standard manager writes its config to
  `<chezmoi config home>/chezmoi/chezmoi.toml` and uses `~/.local/share/chezmoi`
  as the source clone (`installer/lib/dotfilesmanager/chezmoi/manager.go:74-87`).
  `Initialize` writes installer data into that config
  (`installer/lib/dotfilesmanager/chezmoi/data.go:12-70`). `Apply` deletes the
  clone directory and runs `chezmoi init --apply [--branch B] MrPointer --config
  <path>`, which clones a remote branch
  (`installer/lib/dotfilesmanager/chezmoi/apply.go:9-48`). The repository
  `.chezmoi.toml.tmpl` re-renders the config from existing data during
  `chezmoi init` (`.chezmoi.toml.tmpl:1-3`).
- **Personal and work data.** `--work-env` adds work identity and the
  `~/.work/profile` and `~/.work/<work-name>/profile` paths; `--work-name`
  defaults to `sedg` (`installer/cmd/install.go:621-655,771-776`). The existing
  work configuration lives in `private_dot_work/` and is selected through
  `.chezmoiignore:31-37`.
- **Optional tools.** `installer/internal/config/tools.yaml` lists fzf, bat,
  sheldon, eza, ripgrep, fd, difftastic, and uv. The installer installs them
  only after dotfiles are applied.
- **Unguarded optional tools in dotfiles.**
  - `dot_zshrc.tmpl:183` runs `eval "$(sheldon source)"` without a presence
    check. Sheldon loads the plugins in `private_dot_config/sheldon/plugins.toml`.
    No other shell startup line requires those plugins; the remaining
    references are `zstyle` settings and an optional `~/.p10k.zsh` source
    (`dot_zshrc.tmpl:185-241`).
  - `private_dot_config/git/common.gitconfig.tmpl:19` sets
    `core.pager = diff-so-fancy | less --tabs=4 -RFX`. `diff-so-fancy` reaches
    `PATH` only through the Sheldon plugin `so-fancy/diff-so-fancy`.
  - `private_dot_config/git/common.gitconfig.tmpl:22` sets
    `diff.external = difft`; difftastic is an optional tool.
- **Network externals and guarded scripts.** A full apply downloads files listed
  in `.chezmoiexternal.toml.tmpl`. The `uv`, `nono`, and `codex` apply-time
  scripts exit 0 when their command is absent. The sedg VPN script renders only
  on WSL kernels.
- **Existing CI.** `.github/workflows/installer-ci.yml` runs on every pull
  request; its build, test, privilege, and E2E jobs run only when
  `installer/**` or that workflow changes (`:3-45`). The E2E job clones
  `GITHUB_HEAD_REF` through `--git-branch` instead of using the checked-out
  commit (`:251-289`), and its shell check passes without an applied
  `.zshenv` and never starts an interactive shell (`:387-407`). The aggregate
  `Installer CI Status` job passes skipped jobs, and also passes when `build`,
  `test`, or `e2e-test` is canceled; only a canceled `privilege-test` fails it
  (`:409-434`). `test`, `privilege-test`, and `e2e-test` each depend only on
  `build` and run beside each other (`:96,128,153`). Workflow permissions are
  reset to `{}` with per-job read grants (`:16-24`). `installer-ci.yml` also
  runs on pushes to `main` that touch installer paths (`:5-10`).
- **Real installs use the released installer.** `get.sh` downloads the latest
  GitHub release of the installer (`get.sh:66-71`), and `README.md:26-32`
  directs fresh installs through `get.sh`. Releases are built by
  `.github/workflows/release.yml` on `v*.*.*` tags.
- **Linux test image.** `installer/docker/ubuntu/Dockerfile` builds an
  `ubuntu:latest` image with only installer prerequisites (build tools, curl,
  git, ca-certificates, sudo) and a passwordless-sudo user `testuser`.
- **Hosted repository state.** `allow_auto_merge=true`; squash and rebase merges
  are enabled; `delete_branch_on_merge=true`. Rulesets `protect-main` (8895424)
  and `require-installer-ci` (12783037) target the default branch, contain no
  `pull_request` rule, and let the repository admin role bypass always.
  `Installer CI Status` is the only required check, with non-strict policy.
  Pull requests merge as squash commits. Commit `71b9472` ("update nono profile
  schema") reached `main` without a pull request.
- **Schema refresh path.** No automation pushes schema refreshes. The apply-time
  script `run_after_refresh-nono-profile-schema.sh.tmpl` writes the refreshed
  nono schema into the local chezmoi source directory when `nono` is installed,
  and the owner commits it by hand; `71b9472` is such a commit.
- **Hosted macOS image.** `macos-latest` maps to the macOS 26 Arm64 image. Its
  published tool list (`actions/runner-images`, `images/macos/macos-26-arm64-Readme.md`,
  read 2026-11-18) contains none of the `tools.yaml` tools and no
  `diff-so-fancy`; the macOS 15 Arm64 list is the same. The published list
  covers notable tools only.
- **Rate limits.** `.chezmoiexternal.toml.tmpl:1` calls `gitHubLatestRelease`,
  and `get.sh:66-71` queries the latest release unauthenticated. The existing
  E2E job passes `CHEZMOI_GITHUB_ACCESS_TOKEN` from the job token for this
  reason (`.github/workflows/installer-ci.yml:252-253`).
- **Pull-request tooling.** The GitHub path of
  `dot_agents/skills/creating-pull-requests/` does not enable auto-merge.
- **Governance.** `.specify/memory/constitution.md` is an unfilled template and
  adds no constraints.

## Design Decisions

### Add a Prepare-Only Stop to the Installer

Serves: Prepare Without Applying Dotfiles; Validate Dotfiles With Optional Tools
Present.

Add one boolean flag, `--prepare-only`, to `install`. When set, the command runs
every step up to and including `Initialize` (chezmoi installed, config file
written with installer data), skips `Apply`, and then exits successfully.

- Without `--install-tools`, it also skips `installOptionalTools`. This is the
  bare preparation.
- With `--install-tools`, it runs `installOptionalTools` after `Initialize` and
  before exiting. This is the with-tools preparation. Optional-tool failures
  stay non-fatal warnings inside the installer, as in a full install; the
  with-tools validation job enforces presence separately (see Define
  Validation Pass Criteria).

Without `--prepare-only`, behavior is unchanged.

Rationale: the feature needs exactly one stop point before chezmoi applies
dotfiles, and the with-tools variant needs the installer's own optional-tool
installation (`spec.md` › Validate Dotfiles With Optional Tools Present). Running
that step before CI applies the checkout reuses the existing installer step
instead of adding an installation mode, which `spec.md` excludes. Consequence:
in the with-tools variant, tools exist before the first apply. This matches a
re-apply on a machine that already has its tools, rather than the first apply of
a full install, where tools are installed afterward; the bare variant covers the
first-apply state.

### Apply the Pull-Request Merge Commit From the Checkout

Serves: Validate Personal and Existing Work Configurations on Ubuntu and macOS
Before Merge; Prepare Without Applying Dotfiles.

CI validates GitHub's generated pull-request merge commit, which
`actions/checkout` selects by default on `pull_request` events. After
preparation, CI runs `chezmoi init --apply --source <checkout> --config
<installer config path>` with no repository argument. Chezmoi then uses the
existing checkout without cloning and still renders `.chezmoi.toml.tmpl` from
the installer data.

Rationale (human decision, 2026-11-17): the merge commit is closest to what
`main` becomes. Consequence: the required check stays non-strict, so a result
can go stale if `main` moves before merge. The human accepted this over forcing
branch updates.

### Choose the Installer Binary by Pull-Request Content

Serves: Prepare Without Applying Dotfiles; Prioritize Safety Over Validation
Speed.

Human decision, 2026-11-17.

- When the existing installer change filter is true (the pull request changes
  `installer/**` or `.github/workflows/installer-ci.yml`), the validation jobs
  use the binary from the existing `build` job artifact. An installer change is
  then validated together with the dotfiles.
- Otherwise, the validation jobs download the latest released installer for
  their platform, from the same release source `get.sh` uses, and `build` stays
  skipped. The download is authenticated with the job token (for example
  `gh release download`), so it is not subject to the unauthenticated API rate
  limit.

Rationale: an unchanged installer gains nothing from a fresh build, and the
latest release is what a fresh machine runs through `get.sh`. Consequence:
dotfile validation of non-installer pull requests needs a release that already
contains `--prepare-only`; see Planning Handoff.

### Ubuntu Runs in the Existing Minimal Test Image; macOS Asserts the Bare Baseline

Serves: Expose Undeclared Mandatory Dependencies; Validate a Real Clean-Home
Application.

- Ubuntu: build `installer/docker/ubuntu/Dockerfile` and run preparation,
  apply, and shell checks as `testuser` inside that container. Its home starts
  empty of dotfiles, and no optional tool is installed.
- macOS: use the hosted `macos-latest` runner with `HOME` set to a new empty
  directory. The job removes nothing. Bare jobs assert the absence of optional
  tools as pass criterion 7.
- Preparation runs with `--non-interactive --plain --install-prerequisites
  --prepare-only` and default values otherwise. Bare jobs never pass
  `--install-tools`; with-tools jobs always pass it.

Rationale: a container gives Ubuntu a clean baseline. macOS has no containers on
hosted runners. The current image ships none of the optional tools (see
Verified Current State), so the job only asserts their absence and fails, naming
the tool, if a later image adds one; the human then decides how to respond
(human decision, 2026-11-18).

### Guard Optional Tools in Dotfiles Instead of Promoting Them

Serves: Expose Undeclared Mandatory Dependencies.

Human decision, 2026-11-17: guard missing tools rather than make them
prerequisites. Promoting Sheldon would require a new non-Homebrew or
Homebrew-only prerequisite path that needs separate research.

- `dot_zshrc.tmpl`: run `eval "$(sheldon source)"` only when `sheldon` resolves
  on `PATH`. Without Sheldon, the shell starts with the default prompt and no
  Sheldon plugins.
- `core.pager`: use `diff-so-fancy | less --tabs=4 -RFX` when `diff-so-fancy`
  resolves on `PATH` at pager run time; otherwise use `less --tabs=4 -RFX`.
- `diff.external`: run `difft` when it resolves on `PATH` at diff run time;
  otherwise produce a unified diff of the two file versions labeled with the
  repository path, and exit 0 when the files differ so git does not report
  "external diff died".

The checks run at use time, not at template render time, because the installer
installs optional tools after it applies dotfiles; a render-time check would
miss tools installed later until the next apply.

### Define Validation Pass Criteria

Serves: Validate a Real Clean-Home Application; Expose Undeclared Mandatory
Dependencies; Validate Dotfiles With Optional Tools Present.

A validation job passes only when all of these hold, in the prepared clean home.
Every `chezmoi` command below passes the same `--source <checkout>` and
`--config <installer config path>` as the apply, and runs with the job token
exported as `CHEZMOI_GITHUB_ACCESS_TOKEN`.

1. `chezmoi init --apply` (see above) exits 0.
2. `chezmoi verify --exclude scripts` exits 0, so every target matches the
   source state.
3. `chezmoi managed --include files,dirs,symlinks` lists at least one path, and
   every listed path exists in the home. An empty list fails.
4. `zsh -n` succeeds on each applied zsh startup file: `~/.zshenv`, `~/.zshrc`,
   and each `.zshenv` and `.zshrc` under `~/.work/`.
5. Shell startup, under a minimal environment holding only `HOME`, `USER`,
   `TERM=xterm-256color`, and `PATH=/usr/bin:/bin:/usr/sbin:/sbin`, so
   `.zshenv` must set up the rest:
   - A warm-up interactive start `zsh -i -c exit` runs first and must exit 0.
     Its stderr is written to the job log but not judged, so first-run work
     such as Sheldon cloning its plugins can happen.
   - Then noninteractive startup `zsh -c exit` and a second interactive startup
     `zsh -i -c exit` each exit 0 with empty stderr.
6. In a temporary git repository with one modified file, `git diff` and the
   configured `core.pager` command (fed a sample diff) each exit 0 and print
   the change. Both run inside `zsh -i -c` under the minimal environment of
   criterion 5, so they see the same `PATH` as an interactive user.
7. Optional-tool state, checked with `command -v` inside `zsh -i -c` under the
   minimal environment of criterion 5, against the commands the optional tools
   provide (package to command, from `installer/internal/config/packagemap.yaml`
   and each package's installed command): `fzf` → `fzf`; `bat` → `bat`, or
   `batcat` when installed through apt; `sheldon` → `sheldon`; `eza` → `eza`;
   `ripgrep` → `rg`; `fd` → `fd`, or `fdfind` when installed through apt
   (package `fd-find`); `difftastic` → `difft`; `uv` → `uv`; plus
   `diff-so-fancy` from the Sheldon plugin. A tool counts as present when any of
   its command names resolves.
   - Bare jobs: no listed tool is present. A present tool fails the job and is
     named in the output.
   - With-tools jobs: every listed tool is present. A missing tool fails the job
     and is named in the output.

### Exempt Only a Fixed List of Documentation Paths

Serves: Validate Personal and Existing Work Configurations on Ubuntu and macOS
Before Merge; Prioritize Safety Over Validation Speed.

A change-detection job marks a pull request exempt only when every changed path
matches this fixed list. Each entry matches only at the exact location shown;
bare file names never match at other depths.

- Repository root files: `README.md`, `AGENTS.md`, `CLAUDE.md`
- Directories: `docs/**`, `specs/**`, `installer/docs/**`, `.specify/**`
- Installer agent guides: `installer/AGENTS.md`, `installer/CLAUDE.md`,
  `installer/docker/AGENTS.md`, `installer/docker/CLAUDE.md`

Any pull request that changes at least one other path runs the full
validation.

Rationale (human decision, 2026-11-17): documentation, specs, and agent guides
that chezmoi never applies cannot break a fresh installation. Every other path
runs the pipeline, so a new or unknown path fails toward validation. Names are
matched by exact location because rendered files reuse them:
`dot_claude/CLAUDE.md.tmpl`, `dot_codex/AGENTS.md.tmpl`,
`private_dot_config/opencode/AGENTS.md.tmpl`, and README files under
`private_dot_config/specify/` and `private_dot_local/`. `.specify/**` is the
repository-local Spec Kit copy, which chezmoi does not apply because it is
dot-prefixed; the rendered preset copy under `private_dot_config/specify/`
stays validated. Changing the list is itself a non-exempt change, because the
list lives in `.github/workflows/installer-ci.yml`.

### Gate Merges on One Aggregate Check

Serves: Validate Personal and Existing Work Configurations on Ubuntu and macOS
Before Merge; Validate Dotfiles With Optional Tools Present.

Dotfile validation runs a matrix of eight jobs: {Ubuntu, macOS} ×
{personal, work} × {bare, with-tools} (`spec.md` › Validate Dotfiles With
Optional Tools Present). The work entries pass `--work-env` with default work
name and email. An aggregate job named `Dotfile Validation Status` always runs in
`installer-ci.yml`, outside the reusable workflow, so it reports even when
validation is skipped. It passes in two cases only:

- the change-detection job succeeded and marked the pull request exempt; or
- the dotfile validation job (and with it all eight matrix jobs) succeeded.

It fails in every other case, including when change detection fails or is
canceled, and when a non-exempt pull request's validation fails, is canceled,
or is skipped because an installer job did not succeed. The aggregate becomes a
required status check.

### Run Dotfile Validation as a Reusable Workflow Called by Installer CI

Serves: Validate Personal and Existing Work Configurations on Ubuntu and macOS
Before Merge; Preserve Conditional Full Installer E2E Coverage.

Human decision, 2026-11-17: dotfile validation must wait for the installer
jobs when both run, while keeping its own file.

- `.github/workflows/dotfile-validation.yml` is a reusable workflow
  (`on: workflow_call` only). It holds the eight matrix jobs. An input tells it
  whether to use the `build` artifact or the latest release.
- `installer-ci.yml` calls it from a new job. That job needs `changes`,
  `build`, `test`, `privilege-test`, and `e2e-test`, and runs only when the
  pull request is not exempt and either the installer change filter is false
  (installer jobs skipped) or `build`, `test`, `privilege-test`, and
  `e2e-test` all succeeded. `e2e-test` alone is not enough, because `test` and
  `privilege-test` run beside it rather than before it.
- The `changes` job gains a second output that reports whether the pull
  request is exempt from dotfile validation.
- The installer change filter, and so the full installer E2E suite, stays
  conditional on installer changes. A dotfile-only change still skips it.

Rationale: GitHub Actions `needs` works only within one workflow run; a
reusable workflow keeps a separate file while joining that run.

### Run Installer CI on Pull Requests Only

Serves: Route Every Repository Change Through a Pull Request.

Human decision, 2026-11-17: remove the `push` trigger from `installer-ci.yml`.
Every change reaches `main` through a pull request that already ran the same
jobs, so post-merge runs repeat work. The tag-triggered `release.yml` is
unchanged.

### Fail Installer CI Status on Canceled Jobs

Serves: Validate Personal and Existing Work Configurations on Ubuntu and macOS
Before Merge (automatic merge after passing checks).

Human decision, 2026-11-17. `Installer CI Status` fails when `build`, `test`,
`privilege-test`, or `e2e-test` is canceled, as well as when any job fails. It
still passes when those jobs are skipped because the pull request has no
installer changes.

Rationale: with automatic merge, no human sees a canceled run before merge. A
canceled job (including a job timeout, which GitHub reports as canceled) must
not count as a pass.

### Enable Auto-Merge From a Workflow

Serves: Validate Personal and Existing Work Configurations on Ubuntu and macOS
Before Merge (automatic merge after passing checks).

Human decision, 2026-11-17: a workflow only.

A job in a new workflow `.github/workflows/auto-merge.yml` runs on
`pull_request` events `opened`, `reopened`, and `ready_for_review`. For a
non-draft pull request whose head repository is this repository, it runs
`gh pr merge --auto --squash`, which is the same action as the "Enable
auto-merge" button. GitHub then merges once all required checks pass, with no
human action. It is the only job granted `contents: write` and
`pull-requests: write`. Draft pull requests are the opt-out.

The global pull-request skill in `dot_agents/skills/creating-pull-requests/` is
not changed, because it applies to every repository.

### Require Pull Requests With No Bypass

Serves: Route Every Repository Change Through a Pull Request; Validate Personal
and Existing Work Configurations on Ubuntu and macOS Before Merge.

Change hosted rules on the default branch:

- Add a `pull_request` rule to `protect-main` with zero required approvals.
- Add `Dotfile Validation Status` as a required status check, non-strict, in
  `require-installer-ci` next to `Installer CI Status`.
- Remove the always-bypass entry for the repository admin role from both
  rulesets.

Consequence: tool-generated changes, such as a refreshed nono schema, must also
arrive as pull requests. The refresh is committed by hand today (see Verified
Current State), so the ruleset change alone meets the requirement and the
refresh script needs no change. An emergency change also goes through a pull
request, or requires a deliberate ruleset edit.

## Proposed Architecture

- **Installer (`installer/cmd/install.go`)**: owns preparation. New state is one
  flag variable. The stop sits inside `setupDotfilesManager` after
  `initDotfilesManagerData`. When the flag is set, the install command runs
  `installOptionalTools` only if `--install-tools` is also set, then exits.
- **Dotfiles (`dot_zshrc.tmpl`, `private_dot_config/git/common.gitconfig.tmpl`,
  and one new managed helper for the external diff fallback)**: own the
  optional-tool guards. The helper is a managed executable under
  `private_dot_local/bin/` and is referenced from `diff.external` by an absolute
  `$HOME`-based path.
- **Installer CI workflow (`.github/workflows/installer-ci.yml`)**: the single
  pull-request entry point, triggered on `pull_request` only. Job graph:
  `changes` → `build` → (`test`, `privilege-test`, `e2e-test`) →
  `dotfile-validation` (calls the reusable workflow) → `Dotfile Validation
  Status`; `Installer CI Status` still aggregates the installer jobs. When the
  installer filter is false, `build` and its dependents are skipped and
  `dotfile-validation` starts right after `changes`. Permissions stay reset to
  `{}` with per-job `contents: read`; `changes` keeps `pull-requests: read`.
  The existing per-pull-request concurrency group stays.
- **Dotfile Validation reusable workflow
  (`.github/workflows/dotfile-validation.yml`)**: `workflow_call` only; eight
  matrix jobs; reads the installer from the `build` artifact or the latest
  release according to its input.
- **Auto-merge workflow (`.github/workflows/auto-merge.yml`)**: one job, write
  scope limited to that job.
- **Hosted rulesets**: enforce pull requests and the two required checks.

## Contracts and Flows

- **`install --prepare-only` contract**
  - Input: existing install flags plus `--prepare-only`.
  - Output: exit 0; chezmoi installed; config file at the standard path holds
    the same data a full install writes; no source clone created; no dotfiles
    applied. Without `--install-tools`, no optional tools installed. With
    `--install-tools`, the installer's optional-tool step has run.
  - Failure: any failing step before the optional-tool step exits non-zero as
    today. Optional-tool failures stay warnings; the with-tools job's criterion
    7 turns a missing tool into a job failure.
- **Validation flow per matrix job**: check out merge commit → obtain installer
  binary (`build` artifact or authenticated latest-release download) → run
  bare or with-tools preparation → `chezmoi init --apply --source <checkout>
  --config <path>` with `CHEZMOI_GITHUB_ACCESS_TOKEN` set → pass criteria 2–7.
  Any failure fails the job and the aggregate check.
- **Aggregate check contracts**: `Dotfile Validation Status` succeeds only
  under the conditions in "Gate Merges on One Aggregate Check".
  `Installer CI Status` succeeds only when no installer job failed or was
  canceled.
- **Run order on an installer pull request**: `build` → `test`,
  `privilege-test`, `e2e-test` → dotfile validation. If any installer job does
  not succeed, dotfile validation is skipped and `Dotfile Validation Status`
  fails.
- **Merge flow**: pull request opened → auto-merge enabled → both required
  checks report → GitHub squash-merges when both pass.
- **Data contract**: `docs/contracts/chezmoi-data.md` is unchanged; preparation
  writes the same keys.
- **Compatibility**: a full `install` without the new flag behaves as before.
  Machines without Sheldon, diff-so-fancy, or difftastic gain working startup
  and git output; machines with them keep current behavior.
- **Migration order**: see Planning Handoff › Order. A release containing
  `--prepare-only` is tagged before dotfile validation is enabled. Hosted rule
  changes happen only after the validation workflow is on `main`, has passed on
  its own pull request, and has passed the release-path check. Auto-merge lands last, so no pull request can merge
  automatically before both checks are required.

## Alternatives and Risks

Alternatives considered (details in `research.md`): head SHA instead of merge
commit; strict up-to-date policy; always build the installer (rejected for
non-installer pull requests: no gain, and the release matches real installs);
dotfile validation in a separate workflow file without a dependency on
installer jobs (rejected: no ordering across workflows); dotfile jobs inline in
`installer-ci.yml` (rejected: the human wants a separate file); clone a ref
through the installer; run Ubuntu on the hosted VM; no exemptions at all
(rejected by the human as stricter than needed for non-applied docs); a broader
exempt list including agent tooling folders (rejected because some rendered
folders share names with non-rendered ones);
promote Sheldon to a prerequisite; render-time tool checks; removing optional
tools from the macOS image before validation (dropped by the human: the current
image ships none, so only absence is asserted); CI installing optional tools
outside the installer (rejected in `spec.md`); release automation (deferred to
a separate feature in `spec.md`); enabling
auto-merge from the global pull-request skill (rejected: it applies to every
repository); a local wrapper skill (rejected: duplicates the workflow).

Risks:

- **Stale merge-commit result.** With a non-strict check, `main` can move after
  validation. Accepted by the human; the next pull request re-validates.
- **macOS image changes.** A later image may add an optional tool; the bare job
  then fails and names it, and the human decides the response. Preinstalled
  tools that are not in `tools.yaml` can still mask a missing guard; such gaps
  surface on real machines.
- **Network externals and package sources.** Apply downloads remote files, and
  with-tools preparation installs packages; an upstream outage fails
  validation and blocks merges. Safety takes priority; the run is retried
  manually.
- **GitHub API rate limits.** Apply, verify, and the release lookup call the
  GitHub API. All of them use the job token, so they do not share the low
  unauthenticated limit of hosted-runner IP addresses.
- **With-tools order differs from a first install.** With-tools jobs install
  tools before the first apply, matching a re-apply on a machine that has its
  tools; the bare jobs cover the first-apply state.
- **Bot-attributed merges.** When the workflow's token enables auto-merge, the
  resulting push to `main` does not trigger other workflows. After this
  feature no workflow triggers on pushes to `main`, so nothing is lost.
- **Release lag.** Non-installer pull requests validate with the latest
  release, not with unreleased installer changes on `main`. This matches what
  `get.sh` installs on a fresh machine.
- **Installer merged ahead of its release.** An installer pull request that
  also changes dotfiles to depend on the new installer behavior passes with its
  own build. Until a new release is tagged, fresh installs through `get.sh`
  break, and later non-installer pull requests fail validation. Mitigated by
  the interim release rule in Planning Handoff until the separate
  release-automation feature exists.
- **Release bootstrap.** Until a release contains `--prepare-only`, dotfile
  validation of non-installer pull requests cannot run. The handoff orders a
  release before enabling it.
- **Slower installer pull requests.** Dotfile validation waits for the full
  installer chain. Accepted by the human.
- **Write token in pull-request workflow.** Limited to one job, same-repository
  pull requests only, and only `gh pr merge --auto --squash`.
- **No bypass.** Removing the admin bypass blocks direct fixes. Required by
  `spec.md`; recovery is a deliberate ruleset edit.
- **Interactive startup without a terminal.** `zsh -i -c exit` runs without a
  TTY; output from tools that expect a terminal could appear on stderr and fail
  the check. This is intended to surface such output.

## Verification Strategy

- **Installer unit tests** (`installer/cmd/`, PR A): with `--prepare-only`
  alone, `Initialize` runs and `Apply` and optional tools do not; with
  `--prepare-only --install-tools`, `Initialize` and the optional-tool step run
  and `Apply` does not; without `--prepare-only` the sequence is unchanged.
- **Workflow evidence** (PR B): all eight matrix jobs and `Dotfile Validation
  Status` pass on the merge commit.
- **Negative evidence** (throwaway draft pull requests after PR B merges; see
  Planning Handoff › Throwaway pull requests):
  - An unguarded `sheldon` call fails interactive startup in the bare jobs and
    fails the aggregate.
  - A pull request that changes only exempt paths passes the aggregate with
    matrix jobs skipped.
  - A pull request that changes only `dot_claude/CLAUDE.md.tmpl` runs all eight
    matrix jobs with the released installer and `build` skipped.
  - In both of the previous cases installer E2E stays skipped.
  - A broken entry in `private_dot_config/sheldon/plugins.toml` fails the
    with-tools jobs.
- **Release-path evidence** (Planning Handoff step 4): the throwaway dotfile-only
  pull request passes all eight matrix jobs with the released installer and
  `build` skipped.
- **Empty-source evidence** (PR B): job logs show criteria 2 and 3 invoked with
  `--source <checkout>` and `--config <path>`, and the criterion 3 step fails
  when given an empty managed list.
- **Ordering evidence** (the first installer pull request after PR B merges): dotfile validation
  starts only after `test`, `privilege-test`, and `e2e-test` succeed and uses
  the `build` artifact; a deliberately failing installer test skips dotfile
  validation and fails `Dotfile Validation Status`.
- **Cancellation evidence** (PR A): canceling `e2e-test` on an installer pull request
  makes `Installer CI Status` fail.
- **Trigger evidence** (PR A): `installer-ci.yml` has no `push` trigger, and
  merging PR A starts no Installer CI run on `main`.
- **Dotfile guard evidence** (PR B): criterion 6 passes in bare jobs (difftastic
  and diff-so-fancy absent) and in with-tools jobs (both present).
- **Hosted rules evidence**: after the ruleset change, `gh api` output for both
  rulesets shows the `pull_request` rule, both required checks, and no bypass
  actors; a direct push to `main` is rejected.
- **Auto-merge evidence** (after PR C merges): a passing non-draft pull request merges with no
  manual action.

## Planning Handoff

- **Order** (three delivery pull requests, PR A, PR B, and PR C, each merged to
  `main` before the next step, plus one throwaway pull request in step 4 that is
  closed without merging):
  1. **PR A**: installer `--prepare-only` with tests and docs, plus the
     Installer CI changes (pull-request-only trigger; canceled-job failure in
     `Installer CI Status`).
  2. **Release tag**: tag a release from `main` after PR A merges, so it
     contains `--prepare-only`.
  3. **PR B**: dotfile guards (Sheldon, pager, external diff helper), the
     Dotfile Validation reusable workflow, its caller job, the exemption output,
     and `Dotfile Validation Status`. The new validation runs on PR B itself and
     proves the guards.
  4. **Release-path check**: after PR B merges, open a throwaway draft pull request
     that changes one dotfile and no installer path. All eight matrix jobs must
     pass with the released installer and `build` skipped. Close it without
     merging.
  5. **Hosted ruleset changes**, after step 4 passes.
  6. **PR C**: the auto-merge workflow. It lands after step 5, so auto-merge
     never acts while `Dotfile Validation Status` is not yet required.
- **Interim release rule** (operating rule until the separate
  release-automation feature exists): after any pull request merges that
  changes a path under `installer/**` other than `installer/docs/**`,
  `installer/AGENTS.md`, `installer/CLAUDE.md`, `installer/docker/AGENTS.md`,
  and `installer/docker/CLAUDE.md`, the owner tags a new release. Workflow-only
  pull requests such as PR B do not trigger it.
- **Throwaway pull requests**: every pull request opened only to produce
  verification evidence (step 4 and the negative-evidence cases) is opened as a
  draft and closed without merging. Pull-request CI runs on drafts, and the
  auto-merge workflow skips drafts, so such a pull request can never merge
  automatically, whenever it is opened.
- **Outward-facing steps**: step 2 publishes a release, and step 5 changes
  hosted repository settings through `gh api`. The executor confirms with the
  human immediately before each.
- **Documentation impact**: `installer/docs/processes/installation.md` and
  `installer/docs/processes/dotfiles-setup.md` (new flag and stop point);
  `installer/docs/architecture.md` if it lists install flags;
  `docs/architecture.md` (CI validation boundary); `docs/dotfiles/` shell-runtime
  docs for the Sheldon guard and git fallbacks.
- **Skills for executors**: `writing-go-code`, `writing-go-tests`,
  `developing-cli-apps`, `testing-go-code`, `configuring-github-actions`
  (PR A); `managing-chezmoi`, `configuring-zsh`, `configuring-github-actions`
  (PR B); `configuring-github-actions` (PR C).
- **Runtime considerations**: every step needs network access and `gh`
  authentication; workflow evidence comes from hosted runs.

## Source References

- `specs/001-dotfile-regression-ci/spec.md`
- `specs/001-dotfile-regression-ci/grounding-notes.md`
- `specs/001-dotfile-regression-ci/research.md`
- `installer/cmd/install.go`
- `installer/lib/dotfilesmanager/chezmoi/apply.go`, `data.go`, `manager.go`
- `installer/internal/config/tools.yaml`
- `installer/internal/config/packagemap.yaml`
- `installer/docker/ubuntu/Dockerfile`
- `get.sh`, `README.md`, `.github/workflows/release.yml`
- `dot_zshrc.tmpl`, `private_dot_config/sheldon/plugins.toml`
- `private_dot_config/git/common.gitconfig.tmpl`
- `.chezmoi.toml.tmpl`, `.chezmoiignore`, `.chezmoiexternal.toml.tmpl`
- `.github/workflows/installer-ci.yml`
- `dot_agents/skills/creating-pull-requests/`
- Hosted: `gh api repos/MrPointer/dotfiles`, rulesets 8895424 and 12783037

## Review Record

| Role | Report | Round | Revision | Status | Verdict |
|---|---|---|---|---|---|
| plan_design | reviews/plan-design.md | R04 | R6 | Reviewed | Passed |
| plan_clarity | reviews/plan-clarity.md | R03 | R6 | Retained | Passed with concerns |
| plan_overview_fidelity | reviews/plan-overview-fidelity.md | R02 | R6 | Reviewed | Passed |

## Design Acceptance

| Field | Value |
|---|---|
| Status | Accepted |
| Applies To Revision | R6 |
| Recorded By | Human |
| Recorded At | 2026-11-18T21:08:52Z |
| Rationale | The human read plan-overview.md for R6 and confirmed they are happy with the plan. |
