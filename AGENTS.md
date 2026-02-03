# Dotfiles

Personal dotfiles managed with [chezmoi]. This repo is the **chezmoi source directory** - edit files here, not in `~`.

## Available Skills

| Skill                        | Use When                                                         |
| ---------------------------- | ---------------------------------------------------------------- |
| `writing-go-code`            | Writing/editing Go code, tests, mocks, interfaces                |
| `managing-chezmoi`           | chezmoi add/apply/diff, templates, .chezmoiignore, source naming |
| `configuring-zsh`            | .zshrc, .zshenv, plugins, PATH, completions                      |
| `configuring-github-actions` | .github/workflows, CI/CD, matrix builds                          |

## Workflow

1. **Read files first** - Understand the current state and what needs to change
2. **Load relevant skills** - Before editing, load the matching skill from the table above based on the files involved
3. **Plan the approach** - Briefly reason about what changes are needed and how they fit with existing patterns
4. **Make targeted changes** - Follow conventions from the loaded skill

## Directory Structure

```
.
├── dot_claude/                 # ~/.claude (Claude global config)
├── .github/workflows/          # CI/CD pipelines
├── installer/                  # Go installer/bootstrapper (has own CLAUDE.md)
├── dot_config/                 # ~/.config/* files (sheldon, etc.)
├── dot_zshrc, dot_zshenv, ...  # Shell config files
├── private_dot_ssh/            # ~/.ssh (private files)
└── .chezmoiignore              # Files to ignore during apply
```

For chezmoi source naming conventions (`dot_`, `private_`, `.tmpl`), load the `managing-chezmoi` skill.

## Key Conventions

1. **Zsh is the primary shell** - Most shell config is Zsh-specific
2. **Sheldon for plugin management** - Not oh-my-zsh at runtime (used only for vendored snippets)
3. **Templates for conditional config** - `{{ .chezmoi.os }}` for OS-specific logic
4. **Separate work/personal dotfiles** - Work configs loaded conditionally in work environments

## Project Motivation

The goal is an easy, portable setup process for any machine. The installer handles prerequisites (including optional [homebrew]), installs chezmoi, populates its data file with custom keys (e.g., work vs personal environment), and applies the dotfiles.

While all actions could be performed manually, automating the process saves time and reduces complexity across multiple machines.

## Installer

The installer is a Go CLI application in the `installer/` directory. It was rewritten from Bash for better maintainability and testability.

**Tech stack**: [cobra] (CLI), [huh] (interactive UI), [goreleaser] (releases)

**For development details**, see `installer/CLAUDE.md`.

[chezmoi]: https://www.chezmoi.io/
[homebrew]: https://docs.brew.sh/
[cobra]: https://github.com/spf13/cobra
[huh]: https://github.com/charmbracelet/huh
[goreleaser]: https://github.com/goreleaser/goreleaser
