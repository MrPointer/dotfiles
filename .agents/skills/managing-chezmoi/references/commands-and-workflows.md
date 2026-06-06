# Chezmoi commands & workflows

This file is a *practical* quick reference for `chezmoi` operations in a dotfiles repo.

## Mental model

- **Source directory**: where your chezmoi-managed files live (in this repo).
- **Target directory**: where files are applied (usually `$HOME`).
- **State**: chezmoi knows differences between source and target.

## Common workflows

### 1) Edit a managed file (repo-first)

- If starting from a target path, resolve the source with `chezmoi source-path ~/.zshrc`
- Edit the source file in this repo directly (recommended for agent + human collaboration)
- `chezmoi diff` and `chezmoi apply --dry-run --verbose` (validate generated target changes without writing them)

### 2) Add a new file to be managed

- `chezmoi add ~/.config/foo/config.toml`
- Optionally verify/edit the created source:
  - `chezmoi source-path ~/.config/foo/config.toml`
- `chezmoi diff` and `chezmoi apply --dry-run --verbose`

### 3) See what would change

- `chezmoi diff`
- `chezmoi apply --dry-run --verbose`

### 4) Preview only some changes

- `chezmoi diff ~/.zshrc`
- `chezmoi apply --dry-run --verbose ~/.zshrc`
- `chezmoi apply --dry-run --verbose --include dotfiles` (if tags/filters are used)

### 5) Inspect how chezmoi maps source/target

- `chezmoi source-path ~/.zshrc`
- `chezmoi target-path <source-path>`

### 6) External templates / execution

- `chezmoi execute-template < template.tmpl`
- `chezmoi execute-template --data "{\"key\":\"value\"}" < file.tmpl` (JSON)

### 7) Repo-only files

- Some files in this repo are not managed as chezmoi targets.
- If `chezmoi source-path <target>` reports "not managed", validate with source review and `git diff` instead of trying to apply.

## Day-to-day commands

- `chezmoi status`
- `chezmoi diff`
- `chezmoi apply --dry-run --verbose`
- `chezmoi verify` (detect drift)
- `chezmoi doctor` (diagnostic)

## Init / onboarding patterns

- Local repo: `chezmoi init --source <path-to-repo>`
- Remote repo: `chezmoi init <user>/<repo>`

## Useful flags

- `--dry-run` / `-n`: don’t write
- `--verbose` / `-v`: show operations
- `--force`: overwrite conflicts; write-capable, so avoid it during sandboxed agent work
- `--debug`: debugging output

## Safety notes

- Prefer `chezmoi diff` to preview target changes from the source state.
- Prefer `chezmoi apply --dry-run -v` for generated target previews.
- Avoid executing templates that run shell commands unless you trust the source.
