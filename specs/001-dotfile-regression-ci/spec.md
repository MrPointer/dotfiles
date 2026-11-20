---
kind: feature-definition
status: Ready
---

# Feature Definition: Dotfile Regression CI

## Intent

Recent dotfile regressions have broken fresh installations. The current expensive installer end-to-end (E2E) coverage intentionally runs only when installer changes, so dotfile changes can reach `main` without meaningful validation. This feature makes `main` reliably consumable by validating relevant dotfile changes before they merge, rather than relying on a notification after a broken change is pushed.

## Scope

### Included

- Pull requests for every repository change.
- Pre-merge validation on Ubuntu and macOS for both personal and existing work configurations when a pull request affects generated dotfiles or their behavior.
- A preparation capability that uses the installer to establish prerequisites and chezmoi configuration, then stops before chezmoi applies dotfiles.
- CI that applies the exact pull-request target checkout or commit after preparation and validates the resulting dotfiles.
- Validation that renders templates, validates output, fully applies dotfiles into a clean disposable home, and exercises both interactive and noninteractive shell startup.
- A bare validation baseline containing only the prerequisites needed to apply dotfiles and start the shell.
- A with-tools validation variant, alongside the bare baseline, in which the installer's optional tools are present, for each operating system and configuration.
- Automatic merge without further user intervention when the applicable checks pass.

### Excluded

- Running the expensive full installer E2E suite for every dotfile edit.
- Last-passing-revision or release-promotion workflows.
- Autonomous agents that repair regressions.
- Installing optional tools in the bare validation baseline or creating broad new installation modes.
- Release automation and a release versioning scheme, which are a separate feature.
- Defining the technical setup boundary beyond the installer stopping before chezmoi applies dotfiles.
- Refactoring or redesigning the existing work configuration.
- A fixed product validation runtime target or performance-optimization work.

## Requirements

### Route Every Repository Change Through a Pull Request

Every repository change must go through a pull request. Acceptance evidence: no change, including an automatically refreshed upstream nono schema, bypasses the pull-request process.

### Validate Personal and Existing Work Configurations on Ubuntu and macOS Before Merge

Changes that affect generated dotfiles or their behavior must pass pre-merge dotfile validation for both personal and existing work configurations on Ubuntu and macOS through a pull request. Acceptance evidence: a relevant change cannot merge until its applicable validation passes for both configurations on both operating systems, and a passing change can merge automatically without further user intervention.

### Prepare Without Applying Dotfiles

The installer must be able to prepare prerequisites and chezmoi configuration while stopping before chezmoi applies dotfiles. Acceptance evidence: CI can use this stopping capability, then apply the exact target checkout or commit itself.

### Validate a Real Clean-Home Application

Dotfile validation must render templates, validate the rendered output, fully apply the dotfiles into a clean disposable home, and exercise interactive and noninteractive shell startup. Acceptance evidence: validation does not pass because application was skipped or because expected files are missing.

### Expose Undeclared Mandatory Dependencies

The bare validation baseline must contain only prerequisites needed to apply dotfiles and start the shell; optional tools must remain absent. Acceptance evidence: an optional tool is not silently added to CI to make application or shell startup pass, and any such tool found to be required becomes an explicit prerequisite through an approved change.

### Validate Dotfiles With Optional Tools Present

Pre-merge validation must also run a with-tools variant for both personal and existing work configurations on Ubuntu and macOS, in which every optional tool is installed by the installer's own optional-tool installation, with the same application, output, and shell-startup checks as the bare baseline. A failure to install any optional tool fails that variant. Acceptance evidence: a relevant change cannot merge until both the bare and the with-tools variants pass for both configurations on both operating systems; a regression that appears only when optional tools are present fails validation; and a with-tools variant cannot pass with an optional tool missing.

### Prioritize Safety Over Validation Speed

Validation must prioritize safety and regression detection over speed. Acceptance evidence: no product runtime target reduces the agreed validation coverage, and future performance optimization does not weaken that coverage.

### Preserve Conditional Full Installer E2E Coverage

The existing expensive full installer E2E coverage must remain conditional on installer changes. Acceptance evidence: a dotfile-only change does not cause that full suite to run solely because it changes dotfiles.

## Constraints and Dependencies

Validation must use installer preparation before CI applies dotfiles, because a prerequisite-free environment does not provide meaningful coverage. Its bare baseline contains only prerequisites needed to apply dotfiles and start the shell; optional tools remain absent. Its with-tools variant adds every optional tool through the installer's own optional-tool installation. Both variants must cover both personal and existing work configurations on Ubuntu and macOS. The exact checkout or commit under validation must be the one CI applies after preparation.

## Consequential Decision Context

Post-push failure notification is not sufficient because a user can leave `main` broken and later need to apply it elsewhere. Pre-merge validation with automatic merge after passing checks addresses that risk without requiring ongoing user intervention.

The feature does not use last-passing revisions, release promotion, or autonomous repair agents because the user rejected their added complexity or behavior. It also keeps the expensive full installer E2E suite conditional on installer changes rather than applying it to every dotfile edit.

The current work configuration is the validation target. A future work-configuration refactor is outside this feature.

Optional tools remain absent from the bare baseline because they are optional by intent. If validation shows one is required, it must be made an explicit prerequisite through an approved change rather than being silently included in CI.

The with-tools variant was added during planning (2026-11-18). With optional tools absent from the bare baseline and guarded in the dotfiles, the startup and tool behavior that only runs when optional tools are present would otherwise never be validated, so a regression there could reach `main`. The bare baseline still exposes undeclared mandatory dependencies; the with-tools variant covers the path most real machines run. Its tools come from the installer's own optional-tool installation so that the validated path is the one users get. A real installation treats an optional-tool failure as non-fatal, but in validation any such failure fails the variant, so with-tools coverage cannot shrink unnoticed; the accepted cost is that an upstream package outage can block merges.

Release automation and versioning are justified, because a merged installer change otherwise waits for a manual release before fresh installs receive it, but they are deferred to a separate feature to keep this feature bounded.

Unchecked direct-commit exceptions, including for schema refreshes, are excluded because a supposedly mechanical exception could admit a genuine breaking change. Requiring a pull request for every repository change avoids that risk without requiring identical checks for every change.

Safety and regression detection take priority over validation speed. This feature sets no product runtime target; performance optimization may be explored later, but not by weakening agreed coverage.

## Material Open Questions

None.
