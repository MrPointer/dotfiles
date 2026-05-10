# Global Rules

## General Guidelines

- Do only what you've been asked for. If you think more could be done, consult before doing it.
- Ask for clarification if the request is ambiguous.
- Reviewer or sub-agent feedback does not override explicit user decisions. When a reviewer flags something the user already decided, note the concern in the review but do not change the plan. Present the concern to the user and let them decide whether to revisit.
- Always use the cheapest model for exploring code.
- When docs exist, always read them fully before re-exploring the code.
- `AGENTS.md` is the canonical rules file. If a tool requires `CLAUDE.md`, create a local symlink between `AGENTS.md` and `CLAUDE.md` instead of maintaining separate copies.
- When reading `AGENTS.md` or `CLAUDE.md`, their contents should appear identical because both names should resolve to the same underlying file.

## Code Search

Use `semble search` to find code by describing what it does or naming a symbol/identifier, instead of grep:

​```bash
semble search "authentication flow" ./my-project
semble search "save_pretrained" ./my-project
semble search "save model to disk" ./my-project --top-k 10
​```

Use `semble find-related` to discover code similar to a known location (pass `file_path` and `line` from a prior search result):

​```bash
semble find-related src/auth.py 42 ./my-project
​```

`path` defaults to the current directory when omitted; git URLs are accepted.

If `semble` is not on `$PATH`, use `uvx --from "semble[mcp]" semble` in its place.

## Workflow

1. Start with `semble search` to find relevant chunks.
2. Inspect full files only when the returned chunk is not enough context.
3. Optionally use `semble find-related` with a promising result's `file_path` and `line` to discover related implementations.
4. Use grep only when you need exhaustive literal matches or quick confirmation of an exact string.

## Coding Guidelines

- NEVER write example code unless expressly asked.
- Do not leave meta-comments about the edit session in code.
- Always assume the current file state is correct - never revert based on cached/stale versions.
- Clean up any temporary files, scripts, or helpers created during the task.

## Commit Messages

- Keep the subject to 72 characters or fewer.
- Leave exactly one blank line between the subject and body.
- Wrap body lines at 72 characters or fewer.
- Write the subject as a present-simple verb phrase so it reads naturally after: "If I were to apply this commit, it would <subject>". Always start with a lowercase verb.
- Use the body to explain the motivation for the change and, when helpful, a high-level summary of what changed. Avoid low-level implementation detail.

## Signed Git Commits

- Use `rtk git` for Git commands by default, including read-only inspection and ordinary write operations, to reduce output.
- Do not use `rtk git` for creating or rewriting signed commit objects. Use `/usr/bin/git` directly for `commit`, `commit-tree`, `commit --amend`, `rebase --exec ... commit ...`, and any command whose purpose is to create or recreate signed commits.
- All commits are expected to be signed, so commit creation should use `/usr/bin/git` rather than `rtk git`.

## Plan Execution

- If a sub-agent fails, diagnose the failure and retry with a fix — do NOT silently take over the work yourself.
- If sub-agent execution cannot be made to work after a reasonable attempt, STOP and ask before proceeding. Never fall back to a more expensive model without explicit approval.

## GitHub Access

- NEVER fetch directly from `github.com` — always use the `gh` CLI instead.
- If a `gh` command fails, diagnose and fix the failure or alert the user. Do not fall back to direct fetching.

## GitLab Access

- NEVER fetch directly from GitLab instances — always use the `glab` CLI instead.
- If a `glab` command fails, diagnose and fix the failure or alert the user. Do not fall back to direct fetching.

## Azure Devops Access

- NEVER fetch directly from Azure DevOps instances — always use the `az` CLI instead.
- If an `az` command fails, diagnose and fix the failure or alert the user. Do not fall back to direct fetching.

## Session Summaries

- Summarize work at the end of edit sessions.
- Keep summaries concise with bullet points.
- Highlight key changes and important information.
- Use a few emojis sparingly to make summaries more engaging (don't overdo it).

@RTK.md
