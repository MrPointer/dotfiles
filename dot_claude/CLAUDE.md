# Global Rules

## General Guidelines

- Do only what you've been asked for. If you think more could be done, consult before doing it.
- Ask for clarification if the request is ambiguous.
- Reviewer or sub-agent feedback does not override explicit user decisions. When a reviewer flags something the user already decided, note the concern in the review but do not change the plan. Present the concern to the user and let them decide whether to revisit.
- Always use the cheapest model for exploring code.
- When docs exist, always read them fully before re-exploring the code.
- `AGENTS.md` is the canonical rules file. If a tool requires `CLAUDE.md`, create a local symlink between `AGENTS.md` and `CLAUDE.md` instead of maintaining separate copies.
- When reading `AGENTS.md` or `CLAUDE.md`, their contents should appear identical because both names should resolve to the same underlying file.

## Coding Guidelines

- NEVER write example code unless expressly asked.
- Do not leave meta-comments about the edit session in code.
- Always assume the current file state is correct - never revert based on cached/stale versions.
- Clean up any temporary files, scripts, or helpers created during the task.

## Plan Execution

- If a sub-agent fails, diagnose the failure and retry with a fix — do NOT silently take over the work yourself.
- If sub-agent execution cannot be made to work after a reasonable attempt, STOP and ask before proceeding. Never fall back to a more expensive model without explicit approval.

## GitHub Access

- NEVER fetch directly from `github.com` — always use the `gh` CLI instead.
- If a `gh` command fails, diagnose and fix the failure or alert the user. Do not fall back to direct fetching.

## GitLab Access

- NEVER fetch directly from GitLab instances — always use the `glab` CLI instead.
- If a `glab` command fails, diagnose and fix the failure or alert the user. Do not fall back to direct fetching.

## Session Summaries

- Summarize work at the end of edit sessions.
- Keep summaries concise with bullet points.
- Highlight key changes and important information.
- Use a few emojis sparingly to make summaries more engaging (don't overdo it).

@RTK.md
