# Chezmoi troubleshooting

## Table of Contents

- [I edited a file but nothing changed](#i-edited-a-file-but-nothing-changed)
- [chezmoi apply keeps overwriting my manual changes](#chezmoi-apply-keeps-overwriting-my-manual-changes)
- [Template variables not available / data missing](#template-variables-not-available--data-missing)
- [Why is a file encrypted / how to edit it?](#why-is-a-file-encrypted--how-to-edit-it)
- [Conflict / merge issues](#conflict--merge-issues)
- [File exists only on some machines](#file-exists-only-on-some-machines)

---

## I edited a file but nothing changed

- You likely edited the **target** file directly.
- Use `chezmoi edit <target-path>` to edit the source, then `chezmoi apply`.

## chezmoi apply keeps overwriting my manual changes

- That's expected: source-of-truth is the chezmoi **source**.
- If the manual change is desired, either:
  - re-apply it into the source via `chezmoi edit <target-path>`, or
  - stop managing that file (remove from source)

## Template variables not available / data missing

- Check `chezmoi data` output.
- If `chezmoi` is invoked by automation (installer), confirm data file location and keys.

## Why is a file encrypted / how to edit it?

- Encrypted source files can use patterns like `*.age`, `*.gpg`, etc., depending on setup.
- Use `chezmoi edit <target-path>`; chezmoi will transparently decrypt/encrypt when configured.

## Conflict / merge issues

- Use `chezmoi diff` to understand source vs target.
- Re-run apply with `--verbose`.
- Consider `--force` only when you're sure the source is correct.

## File exists only on some machines

- Look for conditional templates and `.chezmoiignore` patterns.
- Search for template conditionals: `{{ if ... }}` and machine/user keys in `chezmoi data`.
