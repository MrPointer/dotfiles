---
kind: subplan
primary_files:
  - dot_zshrc.tmpl
  - private_dot_config/git/common.gitconfig.tmpl
  - private_dot_local/bin/executable_git-external-diff
execution:
  model_tier: mid-tier
  rationale: "Three small runtime guards in shell and git config. They need care with zsh startup, git pager, and git external-diff calling conventions, but no cross-component design."
  skills:
    - managing-chezmoi
    - configuring-zsh
test:
  mode: not-applicable
  basis: "Barrier: the repository has no test harness for applied dotfiles. The automated coverage for these guards is the dotfile validation workflow from packet 04, which runs only on hosted runners after a push. Building a separate dotfile test framework is outside the accepted plan. Alternative verification: render both templates for personal and work data, run `zsh -n` on the rendered zshrc, run shellcheck on the helper, and run the helper and pager command by hand in a temporary git repository with and without difft and diff-so-fancy on PATH. Named gap: hosted bare and with-tools evidence for criteria 5 to 7 on PR B."
contracts:
  produces: []
  consumes: []
---

# Guard Sheldon, the Git Pager, and the Git External Diff

## Objective

Make shell startup and git output work when Sheldon, diff-so-fancy, or
difftastic is missing, and keep current behavior when they are present. Every
check happens at use time, not at template render time.

## Baseline Context and Constraints

- Accepted plan: `plan.md` R6 › Design Decisions › Guard Optional Tools in
  Dotfiles Instead of Promoting Them; › Proposed Architecture › Dotfiles.
- Current state:
  - `dot_zshrc.tmpl` runs `eval "$(sheldon source)"` with no presence check.
  - `private_dot_config/git/common.gitconfig.tmpl` sets
    `core.pager = diff-so-fancy | less --tabs=4 -RFX` and
    `diff.external = difft`.
  - `diff-so-fancy` reaches `PATH` only through the Sheldon plugin
    `so-fancy/diff-so-fancy`.
- Required behavior:
  - Sheldon: run `eval "$(sheldon source)"` only when `sheldon` resolves on
    `PATH`. Without it, the shell starts with the default prompt and no Sheldon
    plugins, and prints nothing to stderr.
  - `core.pager`: run `diff-so-fancy | less --tabs=4 -RFX` when `diff-so-fancy`
    resolves on `PATH` when the pager runs; otherwise run `less --tabs=4 -RFX`.
  - `diff.external`: point to the new helper by an absolute `$HOME`-based path.
    The helper runs `difft` with git's arguments when `difft` resolves on
    `PATH`. Otherwise it prints a unified diff of the old and new file versions,
    labeled with the repository path, and exits 0 when the files differ, so git
    does not report "external diff died".
- Render-time checks are not allowed: the installer installs optional tools
  after it applies dotfiles, so a render-time check would miss them.
- Do not run write-capable `chezmoi apply` in the user's home. Do not modify
  anything under `specs/001-dotfile-regression-ci/`.

## Scope and Primary Files

- `dot_zshrc.tmpl`: the Sheldon line only.
- `private_dot_config/git/common.gitconfig.tmpl`: the `core.pager` and
  `diff.external` values only.
- `private_dot_local/bin/executable_git-external-diff`: new managed helper,
  applied as `~/.local/bin/git-external-diff`.
- Excluded: `private_dot_config/sheldon/plugins.toml`, other startup files,
  docs, workflows.

## Prerequisites and Integration

### Direct Predecessors

| Path | Classification | Reason |
|---|---|---|

### Transitive Context

None.

## Runtime Considerations

- Rendering templates with `chezmoi execute-template` needs chezmoi locally and
  test data for personal and work configurations; use `--init` style data or a
  temporary config file, never the user's real config file.

## Acceptance and Verification

- Both templates render without error for personal data and for work data
  (`--work-env` equivalent keys).
- `zsh -n` passes on the rendered `dot_zshrc.tmpl`.
- `shellcheck` passes on the helper.
- In a temporary git repository with one modified file:
  - with `PATH` lacking `difft`, `git -c diff.external=<helper> diff` exits 0
    and prints a unified diff naming the repository path;
  - with `difft` on `PATH` (if installed locally), the same command shows
    difftastic output;
  - the rendered `core.pager` command, fed a sample diff, exits 0 and prints the
    change with and without `diff-so-fancy` on `PATH`.
- With `sheldon` absent from `PATH`, `zsh -c` sourcing the rendered zshrc
  snippet around the guard produces no stderr.
- Report the alternative verification performed and restate the named hosted
  gap.

## Required Outputs and Handoff

- One packet commit.
- Handoff state for packet 04: the guards are in place, so a bare clean-home
  apply starts zsh without Sheldon and git works without difftastic and
  diff-so-fancy. Packet 04's local container run depends on this state.
- This packet and packet 04 form delivery PR B.
