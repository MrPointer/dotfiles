# Review: Task Executability

- **Role**: Task executability reviewer
- **Tier**: most-capable
- **Latest Round**: R01
- **Latest Verdict**: blocking

## Rounds

### R01

- **Trigger**: initial
- **Scope**: Complete task package for `specs/001-dotfile-regression-ci/` against accepted `plan.md` Revision R6. Checked closed frontmatter, canonical identity, and ledger agreement. Checked the predecessor DAG, direct-predecessor tables and their classifications, primary-file ownership, and concurrent and sequential overlap. Checked producer and consume contracts, model tier, and skill resolution. Checked test mode and basis against `references/testable-work.md`, plus acceptance, verification, output evidence, and handoff. Checked documentation-packet coverage per `references/documentation-planning.md`. Plan fidelity is out of scope except where a mechanic makes a packet unexecutable.
- **Reviewed Paths**: `specs/001-dotfile-regression-ci/plan.md`, `specs/001-dotfile-regression-ci/tasks.md`, `specs/001-dotfile-regression-ci/subplans/01-installer-prepare-only.md`, `specs/001-dotfile-regression-ci/subplans/02-installer-ci-pull-request-gating.md`, `specs/001-dotfile-regression-ci/subplans/03-dotfile-optional-tool-guards.md`, `specs/001-dotfile-regression-ci/subplans/04-dotfile-validation-workflow.md`, `specs/001-dotfile-regression-ci/subplans/05-auto-merge-workflow.md`, `specs/001-dotfile-regression-ci/subplans/06-repository-documentation.md`. Referenced sources read: `.github/workflows/installer-ci.yml`, `installer/cmd/install.go`, `installer/docs/processes/tools-installation.md`, and the headings of `docs/architecture.md`, `docs/dotfiles/domain.md`, and `docs/dotfiles/processes/shell-startup.md`. Skill directories: `.agents/skills/` and the user's global `~/.claude/skills/`.
- **Binding Evidence**: The coordinator invoked a Claude agent as the most-capable task executability reviewer for an initial review (round R01). Workspace: feature worktree `001-dotfile-regression-ci` at HEAD `36127cc`, which the caller reported clean and committed. The reviewer had read and search access plus one write to this report. It could not run shell or git commands, so it confirmed file existence with glob and took tracked status from the caller's statement. Review started and completed 2026-11-18. No reviewed artifact was modified.
- **Verdict**: blocking

#### Current Findings

### Ubuntu With-Tools Cases Left as a Hosted Gap Without a Barrier

- **Impact**: blocking
- **Location**: `subplans/04-dotfile-validation-workflow.md` › frontmatter `test.basis`; › Acceptance and Verification (local Ubuntu run)
- **Evidence**: The basis gives a real barrier for hosted-only behavior: macOS has no containers, and workers have no Actions runner and cannot push. The packet already runs the validation script locally inside the `installer/docker/ubuntu` image for the bare personal and bare work cases. The named gap also includes "the with-tools jobs", but no barrier is stated for running the Ubuntu with-tools cases the same way. The same container and installer build would only add `--install-tools`. These cases are the only planned local check of three things: packet 01's `--prepare-only --install-tools` path, criterion 7's presence check, and the `diff-so-fancy` pager branch. `installOptionalTools` resolves tools through the active package manager (`installer/cmd/install.go:660-718`). Whether every listed tool resolves on Ubuntu is therefore a real risk, and without a local run it would surface only after PR B is pushed. Under `references/testable-work.md`, preference or inconvenience never justifies a not-applicable gap.
- **Recommendation**: Add the Ubuntu with-tools personal and with-tools work cases to the local container run in Acceptance and Verification. Each must pass criteria 1 to 7, and criterion 7 must report every listed tool as present. Narrow the named hosted gap to the macOS jobs, the exemption cases, ordering, and hosted evidence. Alternatively, state a concrete substantial barrier that prevents the local with-tools run, and keep it in the gap.

### Prepare-Only Test Filter Can Pass Without Running Any Test

- **Impact**: concern
- **Location**: `subplans/01-installer-prepare-only.md` › frontmatter `test.basis`; › Acceptance and Verification (first bullet)
- **Evidence**: Verification relies on `go test ./cmd/ -run PrepareOnly -v`. Nothing requires the new test names to contain `PrepareOnly`. A worker who names tests by behavior, for example `Test_SetupDotfilesManager_StopsBeforeApply`, gets a passing run that reports "no tests to run", and the three required sequences go unproven by the named command.
- **Recommendation**: In Scope or Acceptance, require every new test function in `installer/cmd/install_prepare_only_test.go` to have a name that matches `PrepareOnly`. Also require the reported `-v` output to show at least the three sequence cases running and passing.

### Lint Check Has No Exact Command and Its Skill Is Not Declared

- **Impact**: concern
- **Location**: `subplans/01-installer-prepare-only.md` › Acceptance and Verification ("The project linter passes on the changed package, per `linting-go-code`"); › frontmatter `execution.skills`
- **Evidence**: Acceptance relies on the `linting-go-code` skill, which is missing from `execution.skills`. That list has `writing-go-code`, `writing-go-tests`, `developing-cli-apps`, and `testing-go-code`. Acceptance also gives no exact command. `installer/AGENTS.md` names `task lint` (golangci-lint and typos).
- **Recommendation**: Add `linting-go-code` to `execution.skills`, and replace the bullet with the exact command and expected outcome, for example `task lint` from `installer/` exits 0.

### Installer CI Status Acceptance Omits the Changes Job Result

- **Impact**: concern
- **Location**: `subplans/02-installer-ci-pull-request-gating.md` › Acceptance and Verification (Installer CI Status cases); › Baseline Context and Constraints (Keep list)
- **Evidence**: Today's condition also fails when `needs.changes.result == 'failure'` (`.github/workflows/installer-ci.yml:424`). The acceptance table lists outcomes only for the four installer jobs, and the Keep list does not keep the `changes` check. If `changes` fails or is canceled, all four installer jobs are skipped, which matches the "all four skipped: pass" row. A worker who rewrites the condition from the table alone could drop the existing `changes` failure check or leave a canceled `changes` passing. That would let an aggregate pass with no installer evaluation once auto-merge exists.
- **Recommendation**: Add explicit acceptance rows for `changes` failed (fail) and `changes` canceled (fail, unless the human decides otherwise). Or add "the existing `changes` failure check" to the Keep list. If the canceled-`changes` outcome is not settled in the plan, surface it as a plan question rather than leaving it to the worker.

### Push-Only Steps in the Changes Job Have No Owner After the Push Trigger Is Removed

- **Impact**: concern
- **Location**: `subplans/02-installer-ci-pull-request-gating.md` › Scope and Primary Files ("the `on:` block and the `Installer CI Status` job only"); `subplans/04-dotfile-validation-workflow.md` › Scope and Primary Files (`changes` output)
- **Evidence**: The `changes` job has a checkout step guarded by `if: github.event_name == 'push'` and a paths-filter `base: ${{ github.ref }}` that exists for push events (`.github/workflows/installer-ci.yml:28-36`). Packet 02 removes the push trigger but may not touch `changes`. Packet 04 edits `changes` only to add the exempt output. Neither packet says whether these push-only lines are removed or kept, so either worker could make an unplanned edit or leave dead configuration.
- **Recommendation**: State in packet 04's Scope (which already edits `changes`) whether the push-only checkout step and `base` input are removed or kept unchanged. Mirror that in packet 02's handoff state for packet 04.

### Optional-Tools Process Doc Becomes Stale With No Owning Packet

- **Impact**: concern
- **Location**: `tasks.md` › Execution Handoff (Component Documentation Review paths); `subplans/01-installer-prepare-only.md` › Scope and Primary Files (excludes "other docs"); `subplans/06-repository-documentation.md` › Scope (excludes installer docs)
- **Evidence**: `installer/docs/processes/tools-installation.md` is an operational process doc. Its Overview and Trigger say optional tools install "after dotfiles have been applied" once dotfiles setup completes (lines 5 and 9). With `--prepare-only --install-tools` that no longer holds. No packet owns the file. Execution Handoff only defers it to component-documentation review. `references/documentation-planning.md` expects a changed documented process to be planned with exact paths and stale sections, and reserves deferral for component-level detail.
- **Recommendation**: Add `installer/docs/processes/tools-installation.md` (Overview and Trigger) to packet 01's `primary_files` and Scope, next to the other two installer process docs, with a matching Acceptance bullet. If the human prefers deferral, record that decision and its reason in Execution Handoff.
