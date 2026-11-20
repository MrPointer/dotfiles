# Plan Overview: Dotfile Regression CI

- **Reflects Plan Revision**: R6
- **Authority**: Non-normative human view of `plan.md`. `plan.md` prevails on
  any conflict. No workflow phase reads this file.

## Intent

Recent dotfile changes have broken fresh installations. The full installer
end-to-end (E2E) suite runs only when installer files change, so a dotfile-only
change can reach `main` with no real check.

This feature checks relevant dotfile changes before they merge. Every change
goes through a pull request. Ubuntu and macOS are each checked with the
personal and the existing work configuration. Each combination runs twice: a
bare variant with no optional tools, and a with-tools variant with every
optional tool installed. A change that passes merges automatically. Release
automation is a separate feature and is out of scope.

## How Preparation Works

CI needs a machine with the installer's prerequisites and chezmoi config, but
with dotfiles not yet applied. The plan adds one flag, `--prepare-only`, to the
existing `install` command.

- The command runs every step up to writing the chezmoi config, skips applying
  dotfiles, and exits with success.
- Alone, it also skips optional tools. This is the bare preparation.
- With `--install-tools`, it installs the optional tools before exiting. This
  is the with-tools preparation. Tool failures stay warnings inside the
  installer, as in a full install. The validation job turns a missing tool
  into a failure.
- Without `--prepare-only`, nothing changes.

One flag on the existing command avoids a new installation mode, which the spec
excludes. One consequence is that in the with-tools variant, tools exist before
the first apply. That matches a re-apply on a machine that already has its
tools. The bare variant covers the first-apply state.

## What CI Applies

CI checks out the pull request's merge commit. This is the default for
`actions/checkout` on pull requests. After preparation, CI runs
`chezmoi init --apply` against that checkout. Chezmoi does not clone anything.

You decided this on 2026-11-17 because the merge commit is closest to what
`main` becomes. The cost is that the required check is non-strict. If `main`
moves after validation, the result can be stale. You accepted that over forcing
branch updates.

## Which Installer Binary Is Used

You decided on 2026-11-17 that the pull request's content picks the binary.

- If the pull request changes `installer/**` or `installer-ci.yml`, the
  validation jobs use the binary from the existing `build` job. The installer
  change is then validated together with the dotfiles.
- Otherwise, the jobs download the latest released installer, from the same
  release source `get.sh` uses. The download uses the job token, so it avoids
  the low unauthenticated rate limit. `build` is skipped.

An unchanged installer gains nothing from a fresh build. The latest release is
also what a fresh machine runs through `get.sh`.

This has two costs:

- Non-installer pull requests cannot be validated until a release contains
  `--prepare-only`. The plan tags that release first.
- Those pull requests validate against the latest release, not against
  unreleased installer changes on `main`. An installer pull request that also
  changes dotfiles to depend on the new installer behavior can merge before it
  is released. Fresh installs through `get.sh` then break, and later
  non-installer pull requests fail, until a release is tagged. An interim rule
  covers this (see Order of Work).

## Clean Baseline and Optional Tools

The bare baseline holds only what is needed to apply dotfiles and start the
shell.

- Ubuntu runs in the existing minimal test container. No optional tools are
  installed.
- macOS has no containers on hosted runners. The job uses an empty home and
  removes nothing. The current macOS image ships none of the optional tools, so
  the bare job only asserts that they are absent. If a later image adds one,
  the job fails and names it. You then decide the response. You made this
  decision on 2026-11-18.

Three places in the dotfiles use optional tools without checking for them:
the Sheldon call in `dot_zshrc.tmpl`, the `diff-so-fancy` pager, and the
`difft` external diff. You decided to guard them instead of making them
prerequisites. Promoting Sheldon would need a new prerequisite path that needs
separate research.

- Without Sheldon, the shell starts with the default prompt and no plugins.
- Without `diff-so-fancy`, the pager is plain `less`.
- Without `difft`, a small managed helper prints a unified diff.
- The checks run when the tool is used, not when the template renders. The
  installer installs optional tools after it applies dotfiles, so a render-time
  check would miss them until the next apply.

The with-tools variant exists because the bare variant never runs the code that
only runs when tools are present. A regression there could otherwise reach
`main`. Its tools come from the installer's own installation step, so CI
validates the path users get. The cost is that an upstream package outage can
block merges.

## What Counts as a Pass

A job passes only when all of these hold in the clean home. Every chezmoi
command uses the same checkout and config as the apply, with the job token set.

1. The apply exits 0.
2. `chezmoi verify --exclude scripts` exits 0.
3. The list of managed files, directories, and symlinks is not empty, and every
   entry exists.
4. `zsh -n` accepts every applied zsh startup file.
5. Shell startup runs under a minimal environment.
   - A warm-up interactive start must exit 0. Its stderr is logged but not
     judged, so first-run work such as Sheldon cloning plugins can happen.
   - Then a noninteractive start and a second interactive start each exit 0
     with empty stderr.
6. `git diff` and the configured pager work, run inside an interactive shell.
7. Optional-tool state is checked by command name, including `diff-so-fancy`.
   In bare jobs, no tool may be present. In with-tools jobs, every tool must be
   present. A failure names the tool.

Applying nothing, an empty managed list, or a missing tool cannot pass.

## Which Changes Get Validated

A change-detection job exempts a pull request only when every changed path is
on this fixed list:

- Root files: `README.md`, `AGENTS.md`, `CLAUDE.md`
- Directories: `docs/**`, `specs/**`, `installer/docs/**`, `.specify/**`
- Installer agent guides: `installer/AGENTS.md`, `installer/CLAUDE.md`,
  `installer/docker/AGENTS.md`, `installer/docker/CLAUDE.md`

Any pull request that changes at least one other path runs the full validation.
A new or unknown path therefore fails toward validation.

You decided this on 2026-11-17. Documentation, specs, and agent guides that
chezmoi never applies cannot break a fresh installation. Running every pull
request would be stricter than needed for them.

Each entry matches only at the exact location shown. Bare file names never
match at other depths, because rendered files reuse the names. Examples are
`dot_claude/CLAUDE.md.tmpl`, `dot_codex/AGENTS.md.tmpl`, and README files under
`private_dot_config/specify/` and `private_dot_local/`. The `.specify/**` entry
is the repository-local Spec Kit copy, which chezmoi does not apply. The
rendered copy under `private_dot_config/specify/` stays validated.

The list lives in `.github/workflows/installer-ci.yml`. Changing it is
therefore a non-exempt change.

## How Dotfile Validation Joins Installer CI

You decided that dotfile validation must wait for the installer jobs when both
run, while keeping its own workflow file. GitHub Actions `needs` works only
inside one workflow run. So `dotfile-validation.yml` becomes a reusable
workflow, and `installer-ci.yml` calls it from a new job.

- The reusable workflow holds eight jobs: {Ubuntu, macOS} × {personal, work} ×
  {bare, with-tools}. The work jobs use `--work-env` with the default work name
  and email.
- On an installer pull request, the calling job waits for `test`,
  `privilege-test`, and `e2e-test`. All three must succeed first.
- `e2e-test` alone is not enough. `test` and `privilege-test` run beside it,
  not before it, so it can finish while they are still running or have failed.
- If any installer job does not succeed, dotfile validation is skipped and
  `Dotfile Validation Status` fails.
- On a pull request with no installer changes, the installer jobs are skipped
  and dotfile validation starts right after change detection.
- The installer change filter stays as it is, so a dotfile-only change still
  skips the full installer E2E suite.

The cost is that installer pull requests are slower. You accepted that.

## Gating and Merging

- `Dotfile Validation Status` is an aggregate job in `installer-ci.yml`. It sits
  outside the reusable workflow so it reports even when validation is skipped.
  It always runs.
- It passes in two cases only. Either the change-detection job succeeded and
  marked the pull request exempt, or the dotfile validation job succeeded,
  which means all eight jobs succeeded.
- It fails in every other case. That includes a failed or canceled change
  detection. For a non-exempt pull request, it also includes validation that
  fails, is canceled, or is skipped because an installer job did not succeed.
- `Installer CI Status` now also fails when `build`, `test`, `privilege-test`,
  or `e2e-test` is canceled. Today only a canceled `privilege-test` fails it.
  It still passes when those jobs are skipped because there are no installer
  changes. You decided this because no human sees a canceled run before an
  automatic merge. A job timeout counts as canceled, so it must not pass.
- `installer-ci.yml` runs on pull requests only. You decided to remove its
  `push` trigger, because every change reaches `main` through a pull request
  that already ran the same jobs. The tag-triggered release workflow is
  unchanged.
- A new `auto-merge.yml` workflow runs `gh pr merge --auto --squash` when a
  non-draft pull request from this repository is opened, reopened, or marked
  ready. This is the same action as the "Enable auto-merge" button. GitHub then
  merges once all required checks pass, with no click from you.
- Draft pull requests are the opt-out. The auto-merge job is the only one with
  write permission.
- You decided on a workflow only. The global pull-request skill is not changed,
  because it applies to every repository.

## Hosted Rule Changes

To make every change go through a pull request, the plan changes the rules on
the default branch:

- Add a `pull_request` rule to `protect-main` with zero required approvals.
- Add `Dotfile Validation Status` as a non-strict required check.
- Remove the admin role's always-bypass from both rulesets.

Consequences:

- Tool-generated changes, such as a refreshed nono schema, must also arrive as
  pull requests. No automation pushes the refresh today. The apply-time script
  writes the schema into the local source directory, and you commit it by hand.
  So the rule change alone meets the requirement and the script needs no change.
- Emergency fixes also need a pull request or a deliberate ruleset edit.

## Order of Work

The plan has six steps. Three are delivery pull requests (PR A, PR B, PR C).
Each merges to `main` before the next step.

1. **PR A**: installer `--prepare-only` with tests and docs. It also makes the
   Installer CI changes: pull-request-only trigger, and canceled-job failure.
2. **Release tag**: tag a release from `main` after PR A merges, so it contains
   `--prepare-only`.
3. **PR B**: the dotfile guards, the reusable validation workflow, its calling
   job, the exemption output, and `Dotfile Validation Status`. The new
   validation runs on PR B itself and proves the guards.
4. **Release-path check**: after PR B merges, open a throwaway draft pull
   request that changes one dotfile and no installer path. All eight jobs must
   pass with the released installer. Close it without merging.
5. **Hosted rule changes**, after step 4 passes.
6. **PR C**: the auto-merge workflow. It lands last, so auto-merge never acts
   while `Dotfile Validation Status` is not yet required.

Two more rules apply:

- **Interim release rule.** Until the separate release-automation feature
  exists, you tag a new release after any merged pull request that changes
  installer code or files. Documentation and agent-guide paths under
  `installer/` do not count. Workflow-only pull requests such as PR B do not
  trigger it.
- **Throwaway pull requests.** Every pull request opened only to produce
  evidence is a draft and is closed without merging. CI runs on drafts and the
  auto-merge workflow skips them, so such a pull request can never merge
  automatically.

Two steps change things outside the repository. Step 2 publishes a release.
Step 5 changes hosted settings. The executor confirms with you right before
each one.

## Not Covered Here

- Exact installer stop point and `--prepare-only` contract — `plan.md` ›
  Contracts and Flows
- Exact preparation flags and the macOS and Ubuntu setup details — `plan.md` ›
  Design Decisions › Ubuntu Runs in the Existing Minimal Test Image; macOS
  Asserts the Bare Baseline
- Exact command names checked for each optional tool — `plan.md` › Design
  Decisions › Define Validation Pass Criteria
- Current installer, CI, macOS image, rate-limit, and hosted repository facts —
  `plan.md` › Verified Current State
- Job graph, component ownership, and workflow permissions — `plan.md` ›
  Proposed Architecture
- Compatibility and data contract — `plan.md` › Contracts and Flows
- Alternatives considered and the full risk list — `plan.md` › Alternatives and
  Risks
- Tests and evidence for each claim — `plan.md` › Verification Strategy
- Documentation updates and skills for executors — `plan.md` › Planning Handoff
