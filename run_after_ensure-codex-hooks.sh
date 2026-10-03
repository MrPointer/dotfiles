#!/usr/bin/env bash
set -euo pipefail

# Register the Strict output-style reminder without owning Codex's hooks file.
# Existing hooks and unrelated configuration remain unchanged.

CONFIG_DIR="$HOME/.codex"
CONFIG_FILE="$CONFIG_DIR/hooks.json"
# Codex expands HOME when it runs the registered command.
# shellcheck disable=SC2016
HOOK_COMMAND='"$HOME/.codex/hooks/reinforce-strict-output-style.sh"'

if ! command -v codex &>/dev/null; then
  exit 0
fi

if ! command -v jq &>/dev/null; then
  echo "Cannot merge Codex hooks: jq is not installed." >&2
  exit 1
fi

mkdir -p "$CONFIG_DIR"

if [[ -f "$CONFIG_FILE" ]]; then
  if ! jq -e '
    type == "object" and
    (
      if has("hooks") then
        (.hooks | type == "object") and
        (
          if (.hooks | has("UserPromptSubmit")) then
            (.hooks.UserPromptSubmit | type == "array") and
            all(
              .hooks.UserPromptSubmit[];
              type == "object" and
              ((has("hooks") | not) or (.hooks | type == "array"))
            )
          else
            true
          end
        )
      else
        true
      end
    )
  ' "$CONFIG_FILE" >/dev/null; then
    echo "Cannot merge Codex hooks: $CONFIG_FILE has an unsupported structure." >&2
    exit 1
  fi

  matching_count="$(
    jq --arg command "$HOOK_COMMAND" '
      [
        .hooks.UserPromptSubmit[]?.hooks[]?
        | select(.type == "command" and .command == $command)
      ]
      | length
    ' "$CONFIG_FILE"
  )"
  exact_count="$(
    jq --arg command "$HOOK_COMMAND" '
      [
        .hooks.UserPromptSubmit[]?.hooks[]?
        | select(. == {type: "command", command: $command, timeout: 5})
      ]
      | length
    ' "$CONFIG_FILE"
  )"

  if [[ "$matching_count" == "1" && "$exact_count" == "1" ]]; then
    exit 0
  fi
fi

tmp="$(mktemp "${CONFIG_FILE}.tmp.XXXXXX")"
trap 'rm -f "$tmp"' EXIT

# Dollar-prefixed names in this filter are jq variables.
# shellcheck disable=SC2016
merge_filter='
  .hooks //= {} |
  .hooks.UserPromptSubmit //= [] |
  .hooks.UserPromptSubmit |= (
    map(
      if ((.hooks // []) | any(.[]; .type == "command" and .command == $command)) then
        .hooks |= map(select(.type != "command" or .command != $command)) |
        if (.hooks | length) == 0 then empty else . end
      else
        .
      end
    ) + [
      {
        hooks: [
          {
            type: "command",
            command: $command,
            timeout: 5
          }
        ]
      }
    ]
  )
'

if [[ -f "$CONFIG_FILE" ]]; then
  jq --arg command "$HOOK_COMMAND" "$merge_filter" "$CONFIG_FILE" >"$tmp"

  if chmod --reference="$CONFIG_FILE" "$tmp" 2>/dev/null; then
    :
  else
    chmod "$(stat -f '%Lp' "$CONFIG_FILE")" "$tmp"
  fi
else
  jq -n --arg command "$HOOK_COMMAND" "$merge_filter" >"$tmp"
fi

mv "$tmp" "$CONFIG_FILE"
trap - EXIT
