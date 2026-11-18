---
kind: subplan
primary_files:
  - installer/cmd/install.go
  - installer/cmd/install_prepare_only_test.go
  - installer/docs/processes/installation.md
  - installer/docs/processes/dotfiles-setup.md
execution:
  model_tier: mid-tier
  rationale: "Bounded change to one Cobra command: one new flag, one stop point in an existing sequence, and table-driven unit tests. It needs a small seam for testing but no new package or interface design."
  skills:
    - writing-go-code
    - writing-go-tests
    - developing-cli-apps
    - testing-go-code
test:
  mode: required
  basis: "New test file installer/cmd/install_prepare_only_test.go, run from installer/ with `go test ./cmd/ -run PrepareOnly -v` and then the full `go test ./...`. It covers the three step sequences in plan.md › Verification Strategy › Installer unit tests."
contracts:
  produces: []
  consumes: []
---

# Add the `--prepare-only` Stop to the Installer

## Objective

Add the boolean flag `--prepare-only` to the installer `install` command. When it
is set, the command runs every step up to and including the chezmoi manager's
`Initialize` and does not run `Apply`. It runs `installOptionalTools` only when
`--install-tools` is also set, then exits 0. Without the flag, behavior is
unchanged. Update the installer process docs to describe the flag and the stop
point.

## Baseline Context and Constraints

- Accepted plan: `specs/001-dotfile-regression-ci/plan.md` R6 › Design Decisions ›
  Add a Prepare-Only Stop to the Installer; › Contracts and Flows ›
  `install --prepare-only` contract; › Proposed Architecture › Installer.
- Current flow: the `install` command's `Run` function in
  `installer/cmd/install.go` calls `setupDotfilesManager` and then
  `installOptionalTools`. `setupDotfilesManager` calls `dm.Install()`, then
  `initDotfilesManagerData(dm)`, then `dm.Apply()`. Flags are registered in
  `init()`. `installTools` is the existing `--install-tools` variable.
- Contract to meet:
  - `--prepare-only` alone: chezmoi installed, chezmoi config file written at the
    standard path with the same data a full install writes, no source clone, no
    dotfiles applied, no optional tools installed, exit 0.
  - `--prepare-only --install-tools`: the same, plus the existing optional-tool
    step runs before exit. Optional-tool failures stay warnings and do not change
    the exit code, as in a full install.
  - Any failing step before the optional-tool step exits non-zero, as today.
- The chezmoi data contract (`docs/contracts/chezmoi-data.md`) does not change.
- A test seam is allowed when the current code cannot be tested without one, for
  example passing the dotfiles manager into the function that owns the stop. Keep
  it inside `installer/cmd/install.go`.
- Do not modify anything under `specs/001-dotfile-regression-ci/`.

## Scope and Primary Files

- `installer/cmd/install.go`: flag variable, flag registration with help text,
  the stop after `initDotfilesManagerData`, and the optional-tool condition.
- `installer/cmd/install_prepare_only_test.go`: new unit tests.
- `installer/docs/processes/installation.md`: document `--prepare-only` in the
  flow and the steps list (step 9 and step 10 behavior when the flag is set).
- `installer/docs/processes/dotfiles-setup.md`: document that `--prepare-only`
  stops after the config is written and before "Apply Dotfiles".
- Excluded: workflows, dotfiles, `installer/lib/**`, other docs.

## Prerequisites and Integration

### Direct Predecessors

| Path | Classification | Reason |
|---|---|---|

### Transitive Context

None.

## Runtime Considerations

- Go module downloads may need network access the first time tests run.

## Acceptance and Verification

- `go test ./cmd/ -run PrepareOnly -v` (from `installer/`) passes and shows:
  - with `--prepare-only` alone, `Initialize` runs, and `Apply` and the
    optional-tool step do not;
  - with `--prepare-only --install-tools`, `Initialize` and the optional-tool
    step run, and `Apply` does not;
  - without `--prepare-only`, `Install`, `Initialize`, and `Apply` run in that
    order, unchanged.
- `go test ./...` from `installer/` passes.
- The project linter passes on the changed package, per `linting-go-code`.
- `go run . install --help` (from `installer/`) lists `--prepare-only` with help
  text that says it stops before applying dotfiles.
- Both process docs describe the flag, the stop point, and the
  `--install-tools` interaction, and no longer imply that every install applies
  dotfiles.

## Required Outputs and Handoff

- One packet commit with the flag, tests, and the two doc updates.
- Report the test names added and the test output.
- This packet and packet 02 form delivery PR A (see `tasks.md` › Execution
  Handoff).
