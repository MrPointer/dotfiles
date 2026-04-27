---
name: creating-azure-devops-pull-requests
description: "Create Azure DevOps Pull Requests with consistent, high-quality titles and descriptions. Use when opening a PR in Azure DevOps - covers JIRA-based title formatting, summary, decision log, and before/after comparison."
---

# Creating Azure DevOps Pull Requests

Create Pull Requests in Azure DevOps with a consistent structure that is easy to review and understand.

## Gathering Context

Before writing anything, collect the raw material:

1. **Branch name** - extract the JIRA ticket key if present (for example, `ABC-42` from `feature/ABC-42-local-accounts`).
2. **Commits** - `git log <base>..HEAD --oneline` to understand the full scope.
3. **Diff against base** - `git diff <base>...HEAD` (or `--stat` first for an overview). Read enough of the diff to understand what changed and why.
4. **Related context** - conversation history, plan files, linked JIRA tickets, or Azure DevOps work items that explain motivation.

Do not rely solely on commit messages. Read the actual changes.

## Title

Format: `<JIRA-ticket>: <short description>`

- **JIRA ticket** - extract the first JIRA-style key from the branch name and use it exactly as written, such as `ABC-42`.
- **Fallback** - if the branch has no ticket, use `[NO-TICKET]`. This should be rare.
- **Description** - a short sentence (lowercase, no trailing period) that captures the change from the reviewer's perspective. Aim for under 72 characters total when practical.

Examples:
- `ABC-42: add support for local user accounts`
- `PLAT-118: prevent duplicate webhook deliveries on retry`
- `[NO-TICKET]: extract payment validation into dedicated service`

## Description Structure

The PR description has three sections, plus an optional closing line at the very end. All sections use plain, direct language with no filler.

In the **generated PR description**, prefix each section heading with a matching emoji: `📝 Summary`, `🔑 Key Decisions`, `🔄 Before & After`.

### Summary (required)

3-5 lines maximum. Answer: *what changed and why does it matter?*

Write for a reviewer who has not seen the ticket. State the problem or goal, then the approach taken. Do not list files or repeat commit messages.

### Key Decisions (required)

Each significant component or area that was added, changed, or removed gets its own `###` sub-heading under this section. A short paragraph describes what changed and, when the rationale is not obvious from the ticket alone, why.

If a change is not worth more than a single bullet point, it is not worth calling out here. The diff speaks for itself.

Include rationale when:
- The PR makes internal refactors or restructuring choices that are self-contained (not driven by the ticket)
- Non-obvious trade-offs or alternative approaches were considered
- The change deviates from prior patterns or conventions
- A reviewer would question the approach without explanation

```
## 🔑 Key Decisions

### Validation Service
Added `validation.Service` to encapsulate all input validation rules, replacing the inline checks that lived in HTTP handlers. The inline approach had zero test coverage, so extracting it makes the rules unit-testable without HTTP scaffolding.

### Webhook Idempotency
Introduced an idempotency key on webhook delivery to prevent duplicates on retry. Chose optimistic locking over pessimistic locking because contention is rare and the retry cost is low.
```

### Before & After (required)

A table comparing the previous state to the new state. This helps reviewers understand the net effect without reading every line of the diff.

Each row covers one meaningful area of change. Keep cells concise - short phrases, not sentences.

```
## 🔄 Before & After

| Area | Before | After |
|------|--------|-------|
| User creation | Only OAuth providers supported | Local username/password accounts available |
| Validation | Inline in HTTP handler | Extracted to `validation.Service` |
| Retry behavior | Webhooks re-sent on any failure | Idempotency key prevents duplicate delivery |
```

### Ticket Link (conditional)

If the branch name contains a JIRA ticket key, add a closing line as the last line of the PR body using a markdown link to the ticket.

Always use `Closes`: `Closes [ABC-42](https://solaredge-prod.atlassian.net/browse/ABC-42)`

Omit this line entirely if there is no associated ticket.

## Creating the PR

Use the Azure DevOps CLI:

```bash
BODY="$(cat <<'EOF'
## 📝 Summary
<summary>

## 🔑 Key Decisions

### <area>
<details>

## 🔄 Before & After

| Area | Before | After |
|------|--------|-------|

Closes [ABC-42](https://solaredge-prod.atlassian.net/browse/ABC-42)
EOF
)"

az repos pr create \
  --detect true \
  --title "<title>" \
  --description "$BODY"
```

- Target the repository's default branch unless instructed otherwise.
- Do not add reviewers, labels, work items, or auto-complete unless asked.
- If `--detect true` is not enough in the current repo, pass `--org`, `--project`, and `--repository` explicitly.
- If the branch has not been pushed yet, push it first with `git push -u origin HEAD`.

## Anti-Patterns

- **Restating commits** - the commit log is one click away; the description should add context the commits lack.
- **Listing every file changed** - the diff view exists for this; focus on *what* and *why*, not *where*.
- **Vague summaries** - "various improvements" or "refactor code" tells the reviewer nothing.
- **Over-long descriptions** - if the summary exceeds 5 lines, it is doing too much. Push detail into Key Decisions or Before & After.
