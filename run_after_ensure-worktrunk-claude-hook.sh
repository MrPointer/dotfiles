#!/bin/bash
set -euo pipefail

# Ensures the worktrunk post-create hook for Claude Code permissions is registered.
# Idempotent: skips if the hook already exists.

CONFIG_DIR="${XDG_CONFIG_HOME:-$HOME/.config}/worktrunk"
CONFIG_FILE="$CONFIG_DIR/config.toml"

HOOK_NAME="claude-settings"
HOOK_COMMAND="ensure-claude-permissions"

# Skip if worktrunk isn't installed yet
if ! command -v wt &>/dev/null; then
    exit 0
fi

# Ensure config directory exists
mkdir -p "$CONFIG_DIR"

# Skip if hook is already registered
if [[ -f "$CONFIG_FILE" ]] && grep -q "^${HOOK_NAME}\s*=" "$CONFIG_FILE"; then
    exit 0
fi

# Append the post-create hook
if [[ -f "$CONFIG_FILE" ]] && grep -q '^\[post-create\]' "$CONFIG_FILE"; then
    # Section exists, append the hook entry after it
    sed -i.bak "/^\[post-create\]/a\\
${HOOK_NAME} = '${HOOK_COMMAND}'
" "$CONFIG_FILE"
    rm -f "$CONFIG_FILE.bak"
else
    # Section doesn't exist, append both
    printf '\n[post-create]\n%s = '\''%s'\''\n' "$HOOK_NAME" "$HOOK_COMMAND" >> "$CONFIG_FILE"
fi
