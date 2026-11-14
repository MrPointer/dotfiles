# GitHub

Platform-specific details for creating a PR on GitHub. The workflow, title format,
description sections, and linking *conventions* live in `SKILL.md`; this file supplies the
concrete syntax that only applies to GitHub.

## Description length limit

GitHub caps PR bodies at **65,536 characters**. A description that follows `SKILL.md`
stays far below that, so keep the body concise for readability (the anti-patterns in
`SKILL.md` still apply), not to fit a budget.

## Auto-link tokens

GitHub uses one token for both issues and pull requests in the same repository:

- `#<number>` → **issue or pull request**. GitHub renders it compactly as `#123` and shows the
  title on hover, so it can be used on every mention without read fatigue. No bottom-block
  definition is needed for it.

To reference an issue or PR in a **different** repository, use the fully-qualified form
`owner/repo#123`.

```
This is the second PR in the stack on top of #123, which extracted the
validation service. It moves the webhook handler onto that service.

Because the handler now delegates to the service introduced in #123, the
inline retry logic could be deleted entirely.
```

## Closing keyword

`Closes #42` closes the issue when the PR merges into the repository's default branch.
GitHub does not close the issue when the PR merges into another branch, such as a stack
parent.

## Stacked PRs

Applying the Stacked PRs section in `SKILL.md`:

**Target the parent branch at creation.** Pass the parent's branch to `gh pr create` with
`--base <parent-branch>` (`-B`) so the PR's diff shows only its own net change. The first PR
in the stack omits this and takes the default branch as usual.

When the parent merges and its head branch is deleted, GitHub retargets the dependent PRs to
the parent's base automatically. If the parent's branch is kept, the dependent PR must be
retargeted by hand (out of scope, see `SKILL.md`).

**Closing keywords wait for the default branch.** GitHub ignores `Closes #42` while the PR
targets a branch other than the default branch. Still add the line as `SKILL.md` describes,
but do not promise in the description that merging this PR closes the issue.

**Find the parent's PR number** to reference it with `#<number>`:

```bash
gh pr list --head <parent-branch> --state open --json number --jq '.[0].number'
```

**Declare the dependency in prose.** GitHub has no general PR-to-PR dependency mechanism, so
follow the dependency rules in `SKILL.md` with the `#<number>` token above, for example:
"Stacked on #123 — review and merge that first."

## Creating the PR

Use the GitHub CLI:

```bash
gh pr create --title "<title>" --body "$(cat <<'EOF'
<paste the assembled body from assets/pr-description-template.md>
EOF
)"
```

- Target the repository's default branch unless instructed otherwise or the PR is stacked;
  pass `--base <branch>` to override — a stacked PR targets its parent branch (see
  [Stacked PRs](#stacked-prs)).
- Do not pass `--fill` — it replaces the title and body with commit info.
- Do not add reviewers, assignees, labels, or milestones unless asked.
- Create a draft PR (`--draft`) only when explicitly asked.
- If the branch has not been pushed yet, push it first with `git push -u origin HEAD`.
