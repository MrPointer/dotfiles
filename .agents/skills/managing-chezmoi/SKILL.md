---
name: managing-chezmoi
description: Manage dotfiles with chezmoi. Use when adding files to chezmoi, validating source changes with chezmoi status/diff/dry-run, debugging why changes aren't appearing, working with chezmoi templates or .chezmoiignore, mapping source vs target files, resolving conflicts, or asking "how do I manage this file with chezmoi". For command uncertainties, use Context7 to fetch latest docs.
---

# Chezmoi Management

Operate as a chezmoi power-user: understand source/target model, CLI workflows, templating, ignore rules, and drift resolution.

## Quick Reference

**Commands & workflows:** See [commands-and-workflows.md](references/commands-and-workflows.md)
**Troubleshooting:** See [troubleshooting.md](references/troubleshooting.md)

## Sandbox Policy

This repo is edited from inside a nono sandbox. Treat chezmoi source files as the deliverable: do not run commands that write target files, including final `chezmoi apply`, targeted non-dry-run applies, or `--force`. Validate with `chezmoi diff`, `chezmoi apply --dry-run -v`, source review, and git diff. If the user wants changes applied to `$HOME`, give them the exact command to run outside the sandbox instead of running it here.

## Default Workflow

1. **Confirm source ownership**: before editing a file under `~`, use `chezmoi source-path <target>` or inspect whether the source uses `dot_*`, `private_dot_*`, `dot_config/`, `private_dot_config/`, or a template source
2. **Inspect state**: use `chezmoi status` and `chezmoi diff` for managed targets; use `git diff` for repo-only files that chezmoi does not manage as targets
3. **Edit source directly**: modify files in this repo (the source state)
4. **Validate safely**: run `chezmoi diff` and `chezmoi apply --dry-run -v` for managed targets; stop at source review and git diff for unmanaged repo files
5. **Verify**: `chezmoi verify` after larger refactors when it can run without writing target files

## Complexity Decision Tree

| Situation | Action |
|-----------|--------|
| Simple file change | Edit source → `chezmoi diff` → `chezmoi apply --dry-run -v` |
| New file to manage | `chezmoi add <target>` → inspect source → `chezmoi apply --dry-run -v` |
| Unexpected diffs | `chezmoi diff -v` → `chezmoi apply --dry-run -v` → troubleshoot |
| Conditional behavior | Check templates + `chezmoi data` |
| Repo-only file | Edit source → review source → `git diff` |
| Command uncertainty | Fetch docs via Context7 |

## Chezmoi Expertise Checklist

When helping with a request, consider:

- **Source vs target**: User might be editing the wrong side
- **Managed config files**: For target config under `~/.config`, look for matching `dot_config/` or `private_dot_config/` source before editing the target
- **Source file naming**: `dot_` prefixes, `private_` files, `.tmpl` suffix
- **Ignore rules**: `.chezmoiignore` patterns and conditionals
- **Templates**: `{{ ... }}` logic, `chezmoi data`
- **Safe operations**: Prefer dry-runs, source review, and diffs; source edits plus validation are the goal in this sandbox

## Context7 Usage

When uncertain about chezmoi commands/flags or version-specific behavior:

1. Call `resolve-library-id` with `libraryName="chezmoi"`
2. Call `query-docs` with a targeted question
3. Prefer newest documented behavior
