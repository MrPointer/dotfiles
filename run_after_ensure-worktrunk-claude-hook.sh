#!/bin/bash
set -euo pipefail

# Ensures the worktrunk pre-start hook for copying Claude Code settings is registered.
# Copies .claude/settings.local.json from the primary worktree so new worktrees
# inherit accumulated tool permissions.
# Idempotent: skips if the hook already exists with the correct command.
# Migrates from the old post-create hook and ensure-claude-permissions script.

CONFIG_DIR="${XDG_CONFIG_HOME:-$HOME/.config}/worktrunk"
CONFIG_FILE="$CONFIG_DIR/config.toml"

HOOK_NAME="claude-settings"
HOOK_COMMAND='mkdir -p .claude && cp {{ primary_worktree_path }}/.claude/settings.local.json .claude/settings.local.json 2>/dev/null || true'

# Skip if worktrunk isn't installed yet
if ! command -v wt &>/dev/null; then
    exit 0
fi

# Ensure config directory exists
mkdir -p "$CONFIG_DIR"

# Migrate: remove hook from old [post-create] section if present
if [[ -f "$CONFIG_FILE" ]] && grep -q '^\[post-create\]' "$CONFIG_FILE"; then
    sed -i.bak "/^\[post-create\]/,/^\[/{/^${HOOK_NAME}\s*=/d;}" "$CONFIG_FILE"
    # Remove empty [post-create] section (header followed by another section or EOF)
    sed -i.bak '/^\[post-create\]$/{N;/^\[post-create\]\n\[/s/^\[post-create\]\n//;/^\[post-create\]\n$/d;}' "$CONFIG_FILE"
    rm -f "$CONFIG_FILE.bak"
fi

# Remove old hook command (ensure-claude-permissions) if present under [pre-start]
if [[ -f "$CONFIG_FILE" ]] && grep -q "^${HOOK_NAME}.*ensure-claude-permissions" "$CONFIG_FILE"; then
    sed -i.bak "/^${HOOK_NAME}.*ensure-claude-permissions/d" "$CONFIG_FILE"
    rm -f "$CONFIG_FILE.bak"
fi

# Skip if hook is already registered with the correct command
if [[ -f "$CONFIG_FILE" ]] && grep -q "^${HOOK_NAME}\s*=" "$CONFIG_FILE"; then
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
