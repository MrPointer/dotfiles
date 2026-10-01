#!/bin/bash
set -euo pipefail

# Ensures the worktrunk worktree path template is set and the pre-start hook
# for copying agent settings is registered.
# Copies missing files from the primary worktree's .claude and .agents
# directories so new worktrees inherit local agents, skills, and permissions.
# Idempotent: skips if the hook already exists with the correct command.
# Removes the legacy Claude-specific hook name from previous configurations.

CONFIG_DIR="${XDG_CONFIG_HOME:-$HOME/.config}/worktrunk"
CONFIG_FILE="$CONFIG_DIR/config.toml"

# Sibling worktrees joined with a dash; nono cannot open directories with dots.
WORKTREE_PATH_LINE='worktree-path = "{{ repo_path }}/../{{ repo }}-{{ branch | sanitize }}"'

HOOK_NAME="agents-settings"
LEGACY_HOOK_NAME="claude-settings"
HOOK_COMMAND='for dir in .claude .agents .codex .opencode; do mkdir -p "$dir"; cp -Rn "{{ primary_worktree_path }}/$dir/." "$dir/" 2>/dev/null || true; done'

# Skip if worktrunk isn't installed yet
if ! command -v wt &>/dev/null; then
    exit 0
fi

# Ensure config directory exists
mkdir -p "$CONFIG_DIR"

# Ensure the worktree path template, replacing any previous value.
if [[ -f "$CONFIG_FILE" ]] && grep -q "^worktree-path[[:space:]]*=" "$CONFIG_FILE"; then
    if ! grep -Fxq "${WORKTREE_PATH_LINE}" "$CONFIG_FILE"; then
        sed -i.bak "/^worktree-path[[:space:]]*=/c\\
${WORKTREE_PATH_LINE}
" "$CONFIG_FILE"
        rm -f "$CONFIG_FILE.bak"
    fi
else
    # Top-level keys must come before any TOML table, so prepend the line.
    { printf '%s\n' "$WORKTREE_PATH_LINE"; cat "$CONFIG_FILE" 2>/dev/null || true; } > "$CONFIG_FILE.tmp"
    mv "$CONFIG_FILE.tmp" "$CONFIG_FILE"
fi

# Remove the legacy hook name from any previous configuration.
if [[ -f "$CONFIG_FILE" ]] && grep -q "^${LEGACY_HOOK_NAME}[[:space:]]*=" "$CONFIG_FILE"; then
    sed -i.bak "/^${LEGACY_HOOK_NAME}[[:space:]]*=/d" "$CONFIG_FILE"
    rm -f "$CONFIG_FILE.bak"
fi

# Skip if hook is already registered with the correct command
if [[ -f "$CONFIG_FILE" ]] && grep -q "^${HOOK_NAME}[[:space:]]*=" "$CONFIG_FILE"; then
    if grep -Fq "${HOOK_NAME} = '${HOOK_COMMAND}'" "$CONFIG_FILE"; then
        exit 0
    fi

    # Update the previously registered command with the current copy behavior.
    sed -i.bak "/^${HOOK_NAME}[[:space:]]*=/c\\
${HOOK_NAME} = '${HOOK_COMMAND}'
" "$CONFIG_FILE"
    rm -f "$CONFIG_FILE.bak"
    exit 0
fi

# Append the pre-start hook
if [[ -f "$CONFIG_FILE" ]] && grep -q '^\[pre-start\]' "$CONFIG_FILE"; then
    # Section exists, append the hook entry after it
    sed -i.bak "/^\[pre-start\]/a\\
${HOOK_NAME} = '${HOOK_COMMAND}'
" "$CONFIG_FILE"
    rm -f "$CONFIG_FILE.bak"
else
    # Section doesn't exist, append both
    printf '\n[pre-start]\n%s = '\''%s'\''\n' "$HOOK_NAME" "$HOOK_COMMAND" >> "$CONFIG_FILE"
fi
