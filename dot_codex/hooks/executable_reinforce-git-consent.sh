#!/usr/bin/env bash
# Reinforce the commit and push consent rule on every Codex user turn.
set -euo pipefail

reminder="Git consent reminder:
Commit or push only when the user asks for it as part of the current task (an ongoing piece of work, not a single agent run).
Permission lasts until that task ends and does not carry over to a later task, even in the same session. When unsure whether the task has changed, ask.
A commit request does not include pushing; a push needs its own request.
A request to create a pull request on an unpushed branch allows exactly one push to publish that branch; later pushes need a new request."

jq -cn --arg ctx "$reminder" \
  '{hookSpecificOutput: {hookEventName: "UserPromptSubmit", additionalContext: $ctx}}'
