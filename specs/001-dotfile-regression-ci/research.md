# Research: Dotfile Regression CI

## Verified Evidence

- The install command runs chezmoi `Install`, `Initialize`, and `Apply` as one
  uninterrupted sequence inside `setupDotfilesManager`, and optional tools run
  after it. No flag stops before `Apply` (`installer/cmd/install.go:586-618,769-801`).
- `ChezmoiManager.Apply` deletes the clone directory and runs
  `chezmoi init --apply [--branch B] MrPointer --config <path>`, which clones a
  remote branch rather than applying an existing checkout
  (`installer/lib/dotfilesmanager/chezmoi/apply.go:9-48`).
- `ChezmoiManager.Initialize` writes installer data to the chezmoi config file
  (`installer/lib/dotfilesmanager/chezmoi/data.go:12-70`). The repository
  `.chezmoi.toml.tmpl` re-renders that config from the existing data during
  `chezmoi init` (`.chezmoi.toml.tmpl:1-3`).
- `installer-ci.yml` runs on every pull request; its installer build, test,
  privilege, and E2E jobs run only when `installer/**` or the workflow changes
  (`.github/workflows/installer-ci.yml:3-45,150-173`). The required aggregate
  `Installer CI Status` passes skipped jobs (`.github/workflows/installer-ci.yml:409-434`).
- The E2E job clones `GITHUB_HEAD_REF` through `--git-branch`, not the checked-out
  merge commit (`.github/workflows/installer-ci.yml:251-289`). Its shell check
  sources only `.zshenv` and passes without it (`.github/workflows/installer-ci.yml:387-407`).
- `dot_zshrc.tmpl:182-183` runs `eval "$(sheldon source)"` without a presence
  guard, while `installer/internal/config/tools.yaml:6-7` lists Sheldon as
  optional.
- Apply-time scripts for `uv`, `nono`, and `codex` exit 0 when the command is
  absent. The sedg VPN script renders only on WSL kernels
  (`private_dot_work/private_sedg/private_bin/run_once_after_configure-vpn-fix-tool.sh.tmpl:1-2`).
- Hosted state (2026-11-16): `allow_auto_merge=true`; rulesets `protect-main`
  and `require-installer-ci` target `main`, have no `pull_request` rule, and
  allow the repository admin role to bypass always. `Installer CI Status` is the
  only required check, with non-strict policy.
- Commit `71b9472` ("update nono profile schema") reached `main` directly without
  a pull request, so direct commits occur today.
- Existing pull requests merge as squash commits (`(#NN)` subjects on `main`).
- Hosted runners: `ubuntu-latest` and `macos-latest` images preinstall many
  developer tools; the E2E matrix already uses containers for other Linux
  distributions (`.github/workflows/installer-ci.yml:156-173`).

## Alternatives

### Installer Source for Preparation

- Selected by the human on 2026-11-17: use the `build` artifact when the pull
  request changes the installer; otherwise download the latest release, which
  is what `get.sh` installs on a fresh machine (`get.sh:66-71`).
- Always build from the checkout. Rejected for non-installer pull requests: no
  gain over the release.
- Always download the latest release. Rejected for installer pull requests: the
  installer change would not be validated with the dotfiles.

### How CI Applies the Exact Revision

- After preparation, run `chezmoi init --apply --source <checkout> --config
  <installer config path>` with no repository argument. Uses the existing
  checkout without a clone and still renders `.chezmoi.toml.tmpl`. Selected.
- Pass a commit to the installer's clone path. Requires new installer clone
  semantics and a network clone of a ref. Rejected.

### Ubuntu Baseline Environment

- Run in a clean `ubuntu` container on a hosted runner. Optional tools are
  absent by default. Selected.
- Run directly on the `ubuntu-latest` VM. The image preinstalls tools that could
  hide missing guards. Rejected.

### macOS Baseline Environment

- Hosted macOS runner; containers are not available. The published
  `macos-latest` (macOS 26 Arm64) tool list contains no `tools.yaml` tool and no
  `diff-so-fancy` (read 2026-11-18). Bare jobs only assert absence on the
  applied interactive `PATH` and fail naming any tool found. Selected by the
  human on 2026-11-18.
- Remove optional tools before validation. Dropped: nothing to remove today, and
  tools installed outside Homebrew could not be removed reliably.

### Relevance Detection

- Exempt a fixed list of non-applied documentation paths, matched at exact
  locations; every other path runs validation. Selected by the human on
  2026-11-17 after first choosing no exemptions. The list: root `README.md`,
  `RELEASING.md`, `AGENTS.md`, `CLAUDE.md`; `docs/**`, `specs/**`,
  `installer/docs/**`, `.specify/**`; `installer/AGENTS.md`,
  `installer/CLAUDE.md`, `installer/docker/AGENTS.md`,
  `installer/docker/CLAUDE.md`.
- No exemptions. Rejected as stricter than needed for files chezmoi never
  applies.
- Broad exemptions including agent tooling folders and bare file names at any
  depth. Rejected: rendered files reuse those names
  (`dot_claude/CLAUDE.md.tmpl`, `dot_codex/AGENTS.md.tmpl`,
  `private_dot_config/opencode/AGENTS.md.tmpl`).
- Validate only a list of known dotfile paths. Rejected: a new dotfile path
  could be missed.

## Working Decisions and Consequences

- A single new install flag stops after chezmoi config initialization and
  before `Apply`; optional tools are also skipped. It is not a new installation
  mode beyond that stop.
- Validation passes only when apply succeeds, `chezmoi verify` reports no
  drift, every managed target exists, rendered zsh files pass `zsh -n`, and both
  `zsh -c exit` and `zsh -i -c exit` succeed with empty stderr.
- A new required aggregate check fails on any failed or canceled validation
  job and passes only when the change is exempt or all four matrix entries pass.

## Resolved Questions (human, 2026-11-18)

- With-tools validation added to `spec.md` (bare and with-tools variants for
  each OS and profile; eight jobs). Tools come from the installer's own
  optional-tool step; a missing tool fails the job.
- With-tools preparation: `--prepare-only --install-tools` runs the optional-
  tool step before CI applies the checkout. Alternative rejected: a new
  installer mode that installs tools after CI applies, which `spec.md` excludes
  as a broad new installation mode.
- Release automation: separate feature; interim rule to tag a release after
  each installer merge.
- Delivery: PR A → release tag → PR B → rulesets → PR C.

## Resolved Questions (human, 2026-11-17)

- Exact revision: the generated pull-request merge commit, with a non-strict
  required check. Staleness after `main` moves is accepted.
- Sheldon: guard the call. Promotion was rejected because a Sheldon
  prerequisite would need a new install path that needs more research.
  Transitive check: no shell startup line needs Sheldon plugins, but git
  `core.pager` uses `diff-so-fancy`, which only the Sheldon plugin puts on
  `PATH` (`private_dot_config/git/common.gitconfig.tmpl:19`). The pager gets a
  fallback.
- difftastic: `diff.external = difft`
  (`private_dot_config/git/common.gitconfig.tmpl:22`) also gets a fallback in
  this feature.
- Auto-merge: a workflow only. The global pull-request skill applies to every
  repository and is not changed.
- Ordering: dotfile validation is a reusable workflow called by
  `installer-ci.yml`, and waits for `test`, `privilege-test`, and `e2e-test`
  when they run. Those three jobs depend only on `build`, so waiting on
  `e2e-test` alone would not cover them.
- `installer-ci.yml` runs on pull requests only.
- `Installer CI Status` fails on canceled `build`, `test`, or `e2e-test`
  (today it passes, `.github/workflows/installer-ci.yml:424-434`).
- Pull-request enforcement: required by `spec.md`; the plan removes the admin
  always-bypass.

## Open Design Questions

None.
