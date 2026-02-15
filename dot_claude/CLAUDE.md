# Global Rules

## General Guidelines

- Do only what you've been asked for. If you think more could be done, consult before doing it.
- Ask for clarification if the request is ambiguous.

## Coding Guidelines

- NEVER write example code unless expressly asked.
- Do not leave meta-comments about the edit session in code.
- Always assume the current file state is correct - never revert based on cached/stale versions.
- Clean up any temporary files, scripts, or helpers created during the task.

## Plan Execution

- When a plan specifies model assignments for sub-plans, respect them exactly.
- When a plan qualifies for Agent Teams (per the planning skill criteria), use Agent Teams (TeamCreate) — never ad-hoc Task sub-agents. Agent Team teammates have their own context windows and can write files; Task sub-agents cannot write files regardless of permission mode and consume the main context window.
- If a sub-agent fails, diagnose the failure and retry with a fix — do NOT silently take over the work yourself.
- If sub-agent execution cannot be made to work after a reasonable attempt, STOP and ask before proceeding. Never fall back to a more expensive model without explicit approval.

## Session Summaries

- Summarize work at the end of edit sessions.
- Keep summaries concise with bullet points.
- Highlight key changes and important information.
- Use a few emojis sparingly to make summaries more engaging (don't overdo it).
