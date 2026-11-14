# GitLab

Platform-specific details for creating a PR on GitLab (`gitlab.com`). The workflow, title
format, description sections, and linking *conventions* live in `SKILL.md`; this file
supplies the concrete syntax that only applies to GitLab.

## Terminology

GitLab calls a pull request a **merge request** (MR), and the CLI is `glab mr`. Wherever
`SKILL.md` says "PR" — including "Stacked PRs" — read it as "merge request". Every
convention in `SKILL.md` applies unchanged; only the name differs.

## Emoji

GitLab renders **emoji shortcodes** — `:memo:`, `:arrows_counterclockwise:`,
`:earth_africa:` — as native inline emoji, which display more consistently than a pasted
Unicode character. So throughout the MR body, write **every** emoji as its shortcode rather
than a literal Unicode glyph. This applies to both:

- **The template's fixed section headings.** Use `## :memo: Summary` and
  `## :arrows_counterclockwise: Before & After` in place of the literal `📝` / `🔄` that
  `assets/pr-description-template.md` shows.
- **The per-decision cues in Key Decisions.** Choose the emoji exactly as `SKILL.md`
  describes, then write it as its shortcode — a globe cue becomes `:earth_africa:`, a
  cleanup `:broom:`, a database `:card_file_box:`, and so on.

## Description length limit

GitLab caps MR descriptions at roughly **1,048,576 characters (1 MiB)** — effectively no
limit for our purposes. You never need to trim for length. Keep the body concise for
readability (the anti-patterns in `SKILL.md` still apply), not to fit a budget.

## Auto-link tokens

GitLab auto-links two reference tokens inside an MR body. Both render compactly — the
visible text stays `#42` or `!15` and the title shows only on hover — so use the raw token on
every mention. No bottom-block definition is needed for either.

- `#<id>` → **issue**.
- `!<id>` → **merge request**.

```
This is the second MR in the stack on top of !15, which extracted the
validation service. It moves the webhook handler onto that service.

Because the handler now delegates to the service introduced in !15, the
inline retry logic could be deleted entirely.
```

To reference an issue or MR in a **different** project, use the fully-qualified form
`group/project#42` or `group/project!15`.

## Closing keyword

`Closes #42` closes the issue when the MR merges into the project's default branch.

## Stacked PRs

GitLab has **no diff-scoping stacked-MR support** — an MR aimed at the trunk shows every
ancestor's still-unmerged changes. Applying the Stacked PRs section in `SKILL.md`:

**Target the parent branch at creation.** Pass the parent's branch with `-b <parent-branch>`
(`--target-branch`) so the MR's diff shows only its own net change. The first MR in the stack
omits this and takes the default branch as usual.

**Find the parent's MR number** to reference it with `!<id>`. GitLab identifies an MR by its
project-scoped `iid`; list the parent's source branch and read it:

```bash
glab mr list --source-branch <parent-branch> -F json --jq '.[0].iid'
```

GitLab *also* has a native **merge request dependencies** feature (Premium and Ultimate
tiers): an MR can be marked as depending on another, and GitLab then blocks merging until the
dependency merges. This is a merge-order block, **not** diff scoping — it complements
parent-targeting, it does not replace it. When the project's tier supports it, also set the
dependent MR's "Merge request dependencies" field to the parent; there is no `glab` flag, so
set it in the MR's UI after creation.

Whether or not the native field is available, still **declare the dependency in prose** (see
the Stacked PRs section in `SKILL.md`), linking the parent with the `!<id>` token above, for
example: "Stacked on !15 — review and merge that first."

## Creating the PR

Use the GitLab CLI:

```bash
BODY="$(cat <<'EOF'
<paste the assembled body from assets/pr-description-template.md>
EOF
)"

glab mr create \
  --title "<title>" \
  --description "$BODY" \
  --yes
```

- `glab` auto-detects the project and the target branch from the git remote. It targets the
  project's default branch unless instructed otherwise or the MR is stacked; pass
  `-b <branch>` to override — a stacked MR targets its parent branch (see [Stacked PRs](#stacked-prs)).
- `--yes` skips the interactive confirmation. Do not pass `--fill` — it overwrites the title
  and description with commit info.
- Do not add assignees, reviewers, labels, milestones, or auto-merge unless asked.
- If project detection fails, pass the project explicitly with `-R <group>/<project>`.
- If the branch has not been pushed yet, push it first with `git push -u origin HEAD` (or
  add `--push` to the create command).
- Only mark the MR as a draft (`--draft`) if asked.
