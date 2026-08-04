# Canonical AGENTS Adapter

`canonical-agents-adapter` is a standalone Spec Kit preset. It replaces the
constitution template and `/speckit.constitution` command with a fixed
constitution that delegates project governance to the normal harness scope and
keeps optional Spec Kit-only rules separate.

- **Preset version:** `0.1.0`
- **Spec Kit compatibility:** `>=0.12.11`
- **Fixed infrastructure:** `.specify/memory/constitution.md`
- **Writable Spec Kit-local source:** `.specify/memory/constitution.local.md`

It does not copy project rules into the preset or local file, and it does not
modify `AGENTS.md`.

## Install

Chezmoi distributes the global source package; it does not activate a preset in
a project. From the project root, install it explicitly:

```sh
specify preset add --dev "$HOME/.config/specify/presets/canonical-agents-adapter"
```

The project receives its own installed copy under
`.specify/presets/canonical-agents-adapter/`. It is independent of later
changes to the global source.

Refresh a development installation by removing and adding it again:

```sh
specify preset remove canonical-agents-adapter
specify preset add --dev "$HOME/.config/specify/presets/canonical-agents-adapter"
```

Remove it when the project no longer wants the adapter:

```sh
specify preset remove canonical-agents-adapter
```

Removal is safe for the adapter's governance boundary: it does not delete
`AGENTS.md`, `.specify/memory/constitution.local.md`, or the fixed
`.specify/memory/constitution.md`. Review any generated integration command
files and choose the project's normal cleanup or upstream replacement policy.

## Constitution lifecycle

The template is the complete, project-independent constitution artifact; it
has no project-specific placeholders. On the first `/speckit.constitution`
invocation, the command copies that template exactly to
`.specify/memory/constitution.md` when the artifact is missing. Later
invocations leave a matching artifact unchanged. If the artifact differs from
the installed template, the command reports drift and stops rather than
overwriting or repairing it.

## Governance and precedence

The fixed constitution directs every phase to read each applicable root and
nested `AGENTS.md` file under the normal harness scope, using that harness's
broad-to-narrow precedence. It then reads the optional
`.specify/memory/constitution.local.md`.

`AGENTS.md` remains project governance. The local file is only for explicitly
requested Spec Kit additions or overrides. Local precedence applies only to
Spec Kit behavior, and any duplicate or contradictory rule is reported rather
than hidden. A local rule never changes non-Spec Kit behavior.

The command presents effective governance for all source combinations:

- AGENTS present, local present: combine the sources and report precedence,
  duplicates, and conflicts.
- AGENTS present, local missing: use AGENTS governance and report the missing
  optional local source.
- AGENTS missing, local present: use only the local Spec Kit rules and report
  missing project governance.
- Both missing: report both missing sources and do not invent a constitution.

The command creates or amends `.specify/memory/constitution.local.md` only
after an explicit request for concrete Spec Kit-specific rules. A missing local
file is not created by a read or a generic constitution request. The fixed
`.specify/memory/constitution.md` is infrastructure and is never amended by
this command.

## Phase compatibility

This preset does not replace or modify the stock `specify`, `clarify`, `plan`,
`checklist`, `tasks`, `analyze`, or `implement` phases. They continue to consume
`.specify/memory/constitution.md`; the fixed constitution then requires
governance resolution from the applicable `AGENTS.md` files and optional local
source. The package harness verifies this contract text and a plan-equivalent
read, but does not invoke a model through every phase. No phase customization
was needed because the centralized constitution artifact is the existing Spec
Kit seam.

## Independence and limitations

This preset is independent of non-Spec Kit project work. It does not install
workers, add runtime integrations, change planning or execution workflows, or
reinterpret project governance. Its effective scope depends on the active
harness's definition of applicable `AGENTS.md` files and precedence. A harness
that does not expose nested files, source precedence, or file contents limits
what the command can verify; the command reports that limitation instead of
guessing.

The deterministic harness in `tests/test_harness.py` simulates the governance
resolution contract because model execution is unavailable. It uses temporary
project copies and never applies dotfiles to the real home directory.

## Validation scope

Run the package-only checks from this directory:

```sh
./tests/test_harness.py
```

The harness validates manifest mappings and text contracts, root-only and
nested AGENTS resolution, local additions, direct conflicts, missing AGENTS,
missing local, both missing, no-op and explicit local-write command safety,
fixed-constitution immutability, unchanged AGENTS bytes, and a plan-equivalent
phase consuming the effective governance. It does not invoke a model, run an
actual Spec Kit command, prove a harness's hidden AGENTS precedence, or test
runtime integration.
