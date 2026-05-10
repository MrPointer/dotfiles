#!/bin/sh

tool_name=$(jq -r '.tool_name // empty' 2>/dev/null || true)

case "$tool_name" in
  Grep|Glob)
    printf '%s\n' 'Consider using `semble search` instead -- see CLAUDE.md. Grep/Glob are only for exhaustive literal matches or confirming an exact string you already know.' >&2
    ;;
esac

exit 0
