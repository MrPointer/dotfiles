# Preset Protocol Compatibility

## Version Contract

- Preset ID: `timors-agentic-workflow`
- Preset version: `0.1.0`
- Execution-plan protocol version: `0.1.0`
- Supported Spec Kit range: `>=0.12.11,<0.13.0`

The manifest at
`.specify/presets/timors-agentic-workflow/preset.yml` is authoritative for the
preset identity and Spec Kit compatibility range.

## Mandatory Compatibility Preflight

Every composed or replaced command in this preset, including commands added in
later package work, MUST begin with this preflight before extension hooks,
prerequisite scripts, or writes:

1. Read `.specify/presets/timors-agentic-workflow/preset.yml`.
2. Obtain the active version with `specify --version`.
3. Parse the active semantic version and `requires.speckit_version` from the
   manifest.
4. Continue only when the active version satisfies
   `>=0.12.11,<0.13.0`.

If the manifest is unavailable, YAML is malformed, the version cannot be
obtained or parsed, the constraint cannot be parsed, or the version is outside
the range, fail closed. Report the failure and stop before hooks, scripts, or
writes. Do not substitute a global, dotfiles, or package-manager path for the
installed preset path above.

## Execution-Plan Protocol Gate

`execution-plan.md` is valid only when its identity line is exactly
`# Execution Plan: <feature>` and its protocol line is exactly
`**Preset Protocol Version**: 0.1.0`. Missing, malformed, or unsupported
versions fail closed. The only recovery is regeneration through
`/speckit.tasks`; analyze and implement do not repair, infer, downgrade, or
silently accept a different protocol.

The deterministic protocol rules live only in
`.specify/presets/timors-agentic-workflow/references/artifact-validation.md`.
Later `speckit.analyze` and `speckit.implement` command contracts MUST consume
that installed source directly rather than copying or redefining its rules.

## Package Boundary

This preset defines no runtime worker, model binding, or agent-dispatch
configuration. Execution capabilities are declared only as data in validated
execution-plan artifacts.
