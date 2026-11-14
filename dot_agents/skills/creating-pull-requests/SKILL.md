---
name: creating-pull-requests
description: "Create Pull Requests with consistent, high-quality titles and descriptions. Use when opening a PR or merge request — covers a cold read of the diff, title formatting, summary, decision log, before/after comparison, and stacked PRs. Supports GitHub and GitLab."
---

# Creating Pull Requests

Create Pull Requests with a consistent structure that is easy to review and understand.

The workflow, title format, description sections, and linking conventions in this file
apply to every hosting platform. The parts that differ per platform — the concrete
auto-link tokens, emoji syntax, description length limits, stacked-PR targeting, and the
create command — live in a per-platform reference file.

## Select the platform

Read the reference file for the platform you are targeting **before assembling and
creating the PR**. It carries the platform-specific syntax the rest of this skill defers
to:

- **GitHub** → [references/github.md](references/github.md)
- **GitLab** → [references/gitlab.md](references/gitlab.md)

If the platform is not obvious, infer it from the remote (for example, a `github.com`
remote means GitHub and a `gitlab.com` remote means GitLab) or ask.

## Gathering Context

Before writing anything, collect the raw material:

1. **Branch name** - extract the issue number if present (for example, `42` from `feature/42-local-accounts`).
2. **Commits** - `git log <base>..HEAD --oneline` to understand the full scope.
3. **Diff against base** - `git diff --no-ext-diff <base>...HEAD` (or add `--stat` first for an overview). `--no-ext-diff` keeps the output git's own raw diff regardless of any external diff tool configured on the machine. Read enough of the diff to understand what changed and why. For a stacked PR, `<base>` is the parent PR's branch, not the trunk (see [Stacked PRs](#stacked-prs)).
4. **Related context** - conversation history, plan files, or linked issues that explain motivation.

Do not rely solely on commit messages. Read the actual changes.

All PR content must describe the net change from the base branch to `HEAD`. Do not describe intermediate implementation steps, temporary refactors made during the session, or changes elsewhere in the worktree that are not part of the branch diff.

## Cold-Read the Diff

Before writing the description, get an unpoisoned read of the change. You carry the whole
session in your head, so you cannot reliably tell what a reviewer will grasp from the diff
alone and what they won't — everything looks obvious to someone who lived it. A dedicated
cold reader can tell, because it sees only the change.

1. **Stage the diff.** Create a fresh scratch folder with
   `mktemp -d "${TMPDIR:-/tmp}/<project>-<branch>.XXXXXX"` and write
   `git diff --no-ext-diff <base>...HEAD` to `pr-cold-read.diff` inside it. `<project>` is the
   last segment of the `origin` remote URL without any trailing `.git` — not the checkout
   folder's name, which in a worktree is the worktree's name. `<branch>` is the current branch
   with every `/` replaced by `-`. The random suffix keeps concurrent sessions apart, even on
   the same branch. `--no-ext-diff` bypasses any external diff tool configured on the machine
   so the reader always gets git's own raw diff text, identical across machines. This file is
   the one source of truth for both the readers and your own later judgment — do not regenerate it.
   Delete the folder once the PR description is written.
2. **Run the cheapest reader.** Spawn the `cheapest-pr-cold-reader` agent and give it **only
   the absolute path** of the staged file, with the instruction to read it. Its only tool is
   `Read`, and its instructions keep it on that one file, so the path is all it needs — do not
   paste the diff into the prompt. It is also blind by design to author context — **NEVER**
   give it the draft description, commit messages, any session context, or your own summary
   of or hints about the change; those are the author's narration, and once it sees them it
   can no longer tell you whether the *change itself* conveys the intent. It returns a
   coverage note (whether it read the whole file), a significance ranking plus, per key
   change, its inferred intent, a confidence level, and a `clear`/`unclear` verdict.
3. **Escalate the residual.** Two things send work to the `mid-tier-pr-cold-reader` agent,
   run on the same staged diff and equally blind:
   - **Unclear items** — any change the cheapest reader marked `unclear`.
   - **Size** — if the staged file is large enough to risk overflowing the cheapest reader's
     context, route it to the mid tier *from the start*. Check the file's size (`wc -l`) up
     front rather than spending a cheapest-tier read that cannot finish. Short of that, trust
     the cheapest reader's coverage note: if it could not read the whole file, treat that
     like `unclear`.

   A change the mid tier now reads as `clear` was a *capability* gap, not a genuine one —
   treat it as clear. A change it still marks `unclear` is a **real gap** — the trigger to
   add rationale (below). If even the mid tier reports it could not read the whole file, the
   diff is too large to cold-read holistically: say so and split the PR rather than describe
   it from a partial read. Do not tell either reader about the other; each must form an independent
   impression.

Then let that read drive the description:

- **Key Decisions come from the reader's significance ranking**, not your memory of what was
  hard. A fresh reader's sense of what matters is the reviewer's sense; yours is skewed toward
  whatever you wrestled with. Because the reader never saw the session, its ranking cannot
  carry interim-decision noise — that noise exists only in your head.
- **Add rationale where the read is `unclear`.** A real gap is the strongest signal a reviewer
  needs the why, and you — having lived the change — can supply it. A `clear` change is still
  surfaced if significant, but stated lean; it already carries its own why, so do not pad it.
  **NEVER FABRICATE MOTIVATION.** If you genuinely don't know why a change was made either,
  surface it plainly or flag it — do NOT invent a rationale to fill the gap. A made-up "why"
  is worse than none: it misleads the reviewer with false confidence and is exactly the kind
  of leak this whole step exists to prevent.
- **Elevate only what is in the diff.** You may promote a change the reader underweighted when
  its impact isn't visible in the diff (a one-line flag flip with outsized effect). You may
  **not** introduce anything absent from the diff — an interim decision that lives only in your
  memory is exactly the leak this step exists to prevent.
- **Correct a confident-but-wrong read** against what you actually know.

## Title

Format: `<type>: <short description>`

- **Type** - one of: `feat`, `fix`, `refactor`, `docs`, `test`, `chore`, `ci`, `perf`.
- **Description** - a short phrase (lowercase, no trailing period) that captures the change from the reviewer's perspective. Aim for under 72 characters total.

Examples:
- `feat: add support for local user accounts`
- `fix: prevent duplicate webhook deliveries on retry`
- `refactor: extract payment validation into dedicated service`

## Description Structure

Start from [assets/pr-description-template.md](assets/pr-description-template.md) and fill in the three sections. The template carries the section order, the section headings (`📝 Summary` and `🔄 Before & After` keep a fixed emoji; `## Key Decisions` stays plain and lets each decision carry its own — see below), the table shape, and a placeholder describing *what* belongs in each section — including the whole Summary, which needs no more guidance than its placeholder gives.

Some platforms cap the description length or write emoji differently (see the platform
reference). When there is a limit, keep the body concise enough to stay under it.

What the placeholders can't carry is *judgment* — when something is worth calling out and how to reason about it. That is what the rest of this section adds.

### Key Decisions

The placeholder already asks for one `###` sub-heading per significant area, and notes that a change worth no more than a bullet doesn't belong here.

Give each decision a *live* emoji cue. The section heading itself stays plain — `## Key Decisions`, no fixed emoji — because a stock emoji says nothing about *this* PR and can actively mislead: a 🔑 reads as *security keys*, which is noise on most PRs and destroys the signal on a PR that genuinely is about keys. The signal belongs one level down. Prefix each decision's `###` sub-heading with one or more emojis chosen for *that* decision — a globe 🌍 for timezone or locale work, a database 🗄️ for storage, a broom 🧹 for a cleanup — so a reviewer skimming the list gets the gist of each entry at a glance. One emoji is usually enough; add a second only when it genuinely sharpens the meaning. Skip the emoji rather than reach for a vague one — a generic ⚙️ on every heading is as meaningless as a stock 🔑.

Do not explain the order you made changes during the session unless that order is visible in the final branch state and matters to the reviewer.

The [cold read](#cold-read-the-diff) already located where rationale is needed — the changes it marked `unclear`. The cases below are the usual shapes those take:
- The PR makes internal refactors or restructuring choices that are self-contained (not driven by the issue)
- Non-obvious trade-offs or alternative approaches were considered
- The change deviates from prior patterns or conventions
- A reviewer would question the approach without explanation

Example:

```
## Key Decisions

### 🌍 Set-Timezone API
Added an endpoint so a user can set their timezone explicitly, instead of the scheduler guessing it from the browser locale.

### 🕰️🗄️ Timezone Data from a Compiled Database
Compiled the IANA timezone database into the binary rather than reading the host's zoneinfo files, so lookups stay deterministic across deployment environments.
```

### Before & After

The placeholder covers the table shape and cell style. Its purpose is to show the net effect — one row per meaningful area — so a reviewer grasps the change without reading every line of the diff. Every row compares the base branch state to the final branch state; never use intermediate states from the implementation process. Example:

```
## 🔄 Before & After

| Area | Before | After |
|------|--------|-------|
| User creation | Only OAuth providers supported | Local username/password accounts available |
| Validation | Inline in HTTP handler | Extracted to `validation.Service` |
| Retry behavior | Webhooks re-sent on any failure | Idempotency key prevents duplicate delivery |
```

### Issue Link (conditional)

If the branch name contains an issue number, add a closing line as the last *visible* line of the PR body: `Closes #42`. Both supported platforms recognize this keyword and close the issue when the PR merges. Reference-style link definitions may follow it (see [Referencing Issues & PRs](#referencing-issues--prs)).

Always use `Closes`. Omit this line entirely if there is no associated issue.

## Referencing Issues & PRs

These conventions apply on every platform. The concrete syntax — which token auto-links
an issue versus a PR, and how it renders — is in the platform reference file.

The principle behind all of it: **keep full `https://…` URLs out of the prose.**
Inline URLs bloat the source and hurt readability. So collect external URLs into a single
definitions block at the very bottom of the body — each spelled out once — and keep the
prose itself to short labels. Reference-style links are the mechanism for that.

```
See the [migration guide][migration-guide] before reviewing this change.

Closes #42

[migration-guide]: https://example.com/migration-guide
```

### Native issue and PR tokens

Each PR host auto-links its own issues and pull requests from a short token (see the
platform reference for the exact tokens). Use the native token for same-project issues and
PRs rather than a full URL or a reference-style definition.

- **Reference other PRs when they add context** — most often in [stacked PRs](#stacked-prs),
  where the reviewer needs to know which PR this one builds on.
- **Do not spray tokens.** Auto-linked references are convenient but easy to overuse; mention
  a related issue or PR only where it helps a reviewer understand the change or its
  dependencies.

### The link-definitions block

Collect **every** reference-style definition into a single block at the very bottom of the
body, below the `Closes` line. Link definitions render as nothing, so the last *visible* line
stays `Closes #42` while each URL is spelled out exactly once. One definition per label.

## Stacked PRs

A **stacked PR** is one whose branch is built directly on another PR's branch instead of on
the trunk. That creates a genuine dependency: the stacked PR builds on code that is still
under review, and it cannot merge until its parent does.

Stacking is deliberate. A feature is split so an early, self-contained part can ship first —
for reviewer load, shippability, or risk — while the parts that depend on it follow, as one
or more further PRs stacked on top.

**Determine the parent.** The parent is whatever branch the user names. If they don't say,
reconstruct it from git locally. Git records nothing about what a branch was forked from, so
this is a best-effort inference from commit topology — reliable for a real stack, but confirm
with the user when it's ambiguous. Run the bundled helper (path relative to this skill) from
inside the target repository so it reads that repo's branches:

```bash
bash scripts/nearest-parent-branch.sh
```

It ranks local branches by how few commits `HEAD` is ahead of each one's divergence point
(merge-base) and prints the closest few as `<commits-ahead>	<branch>`, nearest first. The
top line is the nearest parent. For the first PR in a stack this comes back as the trunk
(nothing sits between it and `main`) — the "not stacked" case — and that still holds before
you've committed anything, when `HEAD` is still on the trunk tip. If the top two lines share
the same count the parent is genuinely ambiguous; confirm with the user rather than pick one.

You need two things from the parent: its **branch name** (to diff against, and on some
platforms to set the target at creation) and, to declare the dependency, its **PR number** —
look that up by the parent's source branch through the platform CLI (see the platform
reference).

**Target the nearest parent, not the trunk.** The way to make a stacked PR *independently
reviewable* — its diff showing only what it adds, not everything its parents also change — is
to aim it at the branch directly below it in the stack. Given a stack A → B → C where each
branch builds on the one before:

- A is first in line and targets the trunk (`main`, or whatever the repo's default is).
- B targets A's branch.
- C targets B's branch.

Aimed this way, each PR's platform diff is exactly `<parent-branch>...HEAD` — its own net
change. Aiming a stacked PR at the trunk instead would fold every ancestor's still-unmerged
changes into its diff, which is what makes it unreviewable in isolation. The platform
reference gives the flag that sets the target.

**Describe only your own net change.** Read the parent range
(`git diff --no-ext-diff <parent-branch>...HEAD`) and let Summary, Key Decisions, and Before
& After cover only what *this* PR adds on top of its parent — never the parent's own work.

**Declare the dependency** so a reviewer knows the stack exists and in what order it merges:

- In the **dependent PR**, state the dependency at the **start of the Summary**, naming the
  parent with the platform's native PR token — for example, "Stacked on #123 — review and
  merge that first." Be explicit about merge order.
- In the **parent PR**, name the dependents too **while it is still open**, so its reviewer
  knows follow-ups build on it. Once the parent merges this note is moot.

Some platforms also have a native dependency mechanism; see the platform reference.

**Retargeting after a merge (out of scope).** When a PR targets its parent, whatever was
targeting a merged PR must be retargeted at *its* target. In A → B → C, once A merges into
the trunk, B must be retargeted from A to the trunk, while C still targets B. This skill
covers only *creating* stacked PRs with the correct target; performing that rotation as PRs
merge is left to you.

## Creating the PR

1. Assemble the body from `assets/pr-description-template.md` with the sections filled in and
   the link-definitions block at the bottom.
2. Follow the **Creating the PR** section of the platform reference for the exact create
   command and any platform flags.

## Anti-Patterns

- **Restating commits** - the commit log is one click away; the description should add context the commits lack.
- **Describing session history instead of branch diff** - PR text must reflect `base...HEAD`, not the sequence of edits made while implementing the change.
- **Claiming a stacked PR's parent changes as your own** - diff against the parent branch, not the trunk, so the description covers only this PR's net change (see [Stacked PRs](#stacked-prs)).
- **Listing every file changed** - the diff view exists for this; focus on *what* and *why*, not *where*.
- **Vague summaries** - "various improvements" or "refactor code" tells the reviewer nothing.
- **Ignoring a platform's length limit** - when the platform caps description length, trim detail before creating the PR; summarize and collapse lower-value decisions instead of overflowing.
- **Over-long descriptions** - if the summary exceeds 5 lines, it is doing too much. Push detail into Key Decisions or Before & After.
- **Spraying auto-links** - overusing native issue or PR references makes a PR unreadable for humans; follow the referencing conventions above.
