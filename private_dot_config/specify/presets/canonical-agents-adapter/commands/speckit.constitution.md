---
description: "Present effective governance and manage only explicitly requested Spec Kit-local rules."
---

# Manage Spec Kit-local governance

This command initializes the fixed project constitution once, then manages only
Spec Kit-local governance. It never turns the fixed constitution into an
editable project-policy document.

## User input

$ARGUMENTS

## Required behavior

1. Read the installed static template at
   `.specify/presets/canonical-agents-adapter/templates/constitution-template.md`.
   If `.specify/memory/constitution.md` does not exist, copy the template to
   that path exactly. This is the only permitted write to the fixed
   constitution, and it occurs only when it is missing.
2. If `.specify/memory/constitution.md` already exists, compare it with the
   installed static template. If it differs, report the drift and stop before
   any governance write. Do not overwrite, amend, merge, repair, or regenerate
   the existing constitution. If it matches, leave it unchanged.
3. Establish the normal harness scope and read **all applicable root and
   nested `AGENTS.md` files** in that scope. Apply their normal precedence;
   never assume that the root file is the only governance source.
4. Read `.specify/memory/constitution.local.md` if it exists. This file is
   optional and is the only writable source managed by this command.
5. Read the fixed `.specify/memory/constitution.md` and present it as fixed
   infrastructure, not as a writable substitute for the local file.
6. Present the effective governance with source paths and precedence. Detect
   and report duplicate rules and contradictory rules. A local rule can add to
   or override a rule for Spec Kit only; local precedence does not alter
   non-Spec Kit project governance. Report every such conflict explicitly.
7. Treat a local-governance write as authorized only when `$ARGUMENTS`
   explicitly requests a
   concrete **Spec Kit-specific** rule to be added, changed, or removed. A
   generic request to “make a constitution” or a request about project policy
   is not sufficient authorization.
8. Without that explicit request, present the resolution and perform no local
   governance write.
   In particular, do not create `.specify/memory/constitution.local.md` merely
   because it is missing.
9. With explicit Spec Kit-specific write intent, create or amend only
   `.specify/memory/constitution.local.md`. Preserve unrelated local rules,
   ask before replacing ambiguous rules, and report the resulting local rules.
   Never write `AGENTS.md`, `.specify/memory/constitution.md`, this command,
   the static constitution template, or any other project file.

## Source combinations

Handle and report all four combinations rather than inventing missing policy:

| Applicable `AGENTS.md` | Local constitution | Required result |
|---|---|---|
| Present, including nested files | Present | Read all sources, apply normal AGENTS precedence, apply local Spec Kit precedence, and report duplicates or conflicts. |
| Present, including nested files | Missing | Apply and present AGENTS governance; report the missing optional local source; do not create it without explicit Spec Kit-specific rules. |
| Missing | Present | Report missing project governance, apply only the local Spec Kit rules, and report that they do not govern non-Spec Kit work. |
| Missing | Missing | Report both missing sources, present no invented constitution, and do not write a local file without explicit Spec Kit-specific rules. |

Nested `AGENTS.md` files are part of the first two cases whenever the normal
harness scope includes them. Do not flatten their precedence or silently omit
one.

## Safety boundary

This command is independent of non-Spec Kit work. It never copies `AGENTS.md`
contents, never changes project governance, and never silently repairs or
regenerates the fixed constitution. If a source is unavailable, contradictory,
or outside the normal harness scope, report that fact and stop at the affected
decision rather than inventing a rule.
