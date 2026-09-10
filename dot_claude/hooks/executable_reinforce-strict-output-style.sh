#!/usr/bin/env bash
# Reinforce the Strict output style on every user turn.
#
# Claude Code reads an output style once at session start and never re-sends it,
# so a long session drifts away from the required format. This UserPromptSubmit
# hook re-injects the section contract as additional context before each prompt.
#
# It stays silent unless the Strict output style is the active default, so it
# does no harm when another style is in use. A per-project outputStyle override
# in a project settings file is not detected here; this checks the user setting.
set -euo pipefail

settings="$HOME/.claude/settings.json"

style=""
if [ -f "$settings" ]; then
  style="$(jq -r '.outputStyle // ""' "$settings" 2>/dev/null || true)"
fi

if [ "$style" != "Strict" ]; then
  exit 0
fi

reminder="Strict output-style reminder (conversational replies only, never inside files):
Use exactly these four sections, with these headings and emojis, in this order, and no others:
1. 🔍 Findings/Context
2. ⚠️ Limits/Risks
3. ✅ Done
4. ➡️ Next
Omit a section only when it is empty. Do not rename, merge, or invent sections. Keep the whole block together at the end of the reply. Escape hatch: a one or two line answer may drop all four headings together, never a subset."

jq -cn --arg ctx "$reminder" \
  '{hookSpecificOutput: {hookEventName: "UserPromptSubmit", additionalContext: $ctx}}'
