---
name: creating-github-pull-requests
description: "Create GitHub Pull Requests with consistent, high-quality titles and descriptions. Use when opening a PR on GitHub — covers title formatting, summary, decision log, and before/after comparison."
---

# Creating GitHub Pull Requests

Create Pull Requests on GitHub with a consistent structure that is easy to review and understand.

## Gathering Context

Before writing anything, collect the raw material:

1. **Branch name** — extract any issue number (e.g., `42` from `feature/42-local-accounts`).
2. **Commits** — `git log <base>..HEAD --oneline` to understand the full scope.
3. **Diff against base** — `git diff <base>...HEAD` (or `--stat` first for an overview). Read enough of the diff to understand what changed and why.
4. **Related context** — conversation history, plan files, or linked issues that explain motivation.

Do not rely solely on commit messages — they often lack the bigger picture. Read the actual changes.

All PR content must describe the net change from the base branch to `HEAD`. Do not describe intermediate implementation steps, temporary refactors made during the session, or changes elsewhere in the worktree that are not part of the branch diff.

## Title

Format: `<type>: <short description>`

- **Type** — one of: `feat`, `fix`, `refactor`, `docs`, `test`, `chore`, `ci`, `perf`.
- **Description** — a short sentence (lowercase, no trailing period) that captures the change from the reviewer's perspective. Aim for under 72 characters total.

Examples:
- `feat: add support for local user accounts`
- `fix: prevent duplicate webhook deliveries on retry`
- `refactor: extract payment validation into dedicated service`

## Description Structure

The PR body has three sections, plus an optional issue-closing line at the very end. All sections use plain, direct language — no filler, no marketing.

In the **generated PR description**, prefix each section heading with a matching emoji: `📝 Summary`, `🔑 Key Decisions`, `🔄 Before & After`.

### Summary (required)

3–5 lines maximum. Answer: *what changed and why does it matter?*

Write for a reviewer who has not seen the ticket. State the problem or goal, then the approach taken. Do not list files or repeat commit messages.

Describe only the final reviewer-visible outcome relative to the base branch.

### Key Decisions (required)

Each significant component or area that was added, changed, or removed gets its own `###` sub-heading under this section. A short paragraph describes what changed and, when the rationale is not obvious from the ticket alone, why.

If a change isn't worth more than a single bullet point, it's not worth calling out here — the diff speaks for itself.

Do not explain the order you made changes during the session unless that order is visible in the final branch state and matters to the reviewer.

Include rationale when:
- The PR makes internal refactors or restructuring choices that are self-contained (not driven by the ticket)
- Non-obvious trade-offs or alternative approaches were considered
- The change deviates from prior patterns or conventions
- A reviewer would question the approach without explanation

```
## Key Decisions

### Validation Service
Added `validation.Service` to encapsulate all input validation rules, replacing the inline checks that lived in HTTP handlers. The inline approach had zero test coverage — extracting it makes the rules unit-testable without HTTP scaffolding.

### Webhook Idempotency
Introduced an idempotency key on webhook delivery to prevent duplicates on retry. Chose optimistic locking over pessimistic since contention is rare and the retry cost is low.
```

### Before & After (required)

A table comparing the previous state to the new state. This helps reviewers understand the net effect without reading every line of the diff.

Each row covers one meaningful area of change. Keep cells concise — short phrases, not sentences.

Every row must compare the base branch state to the final branch state. Never use intermediate states from the implementation process as the `Before` or `After` values.

```
## Before & After

| Area | Before | After |
|------|--------|-------|
| User creation | Only OAuth providers supported | Local username/password accounts available |
| Validation | Inline in HTTP handler | Extracted to `validation.Service` |
| Retry behavior | Webhooks re-sent on any failure | Idempotency key prevents duplicate delivery |
```

### Issue Link (conditional)

If the branch name contains a GitHub issue number, add a closing keyword as the last line of the PR body. This auto-closes the issue when the PR merges.

Always use `Closes`: `Closes #42`

Omit this line entirely if there is no associated issue.

## Creating the PR

Use the GitHub CLI:

```bash
gh pr create --title "<title>" --body "$(cat <<'EOF'
<description body>
EOF
)"
```

- Target the repository's default branch unless instructed otherwise.
- Do not add reviewers, labels, or milestones unless asked.
- If the branch has not been pushed yet, push it first with `git push -u origin HEAD`.

## Anti-Patterns

- **Restating commits** — the commit log is one click away; the description should add context the commits lack.
- **Describing session history instead of branch diff** — PR text must reflect `base...HEAD`, not the sequence of edits made while implementing the change.
- **Listing every file changed** — the diff view exists for this; focus on *what* and *why*, not *where*.
- **Vague summaries** — "various improvements" or "refactor code" tells the reviewer nothing.
- **Over-long descriptions** — if the summary exceeds 5 lines, it is doing too much. Push detail into Key Decisions or Before & After.
