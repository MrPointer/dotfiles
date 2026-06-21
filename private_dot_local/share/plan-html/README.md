# plan-html — interactive HTML review views for markdown plan directories

Renders a plan directory (a `00-master.md` + numbered sub-plans + optional `reviews/`)
into a single self-contained, dark-themed **dashboard + cards** HTML file you can open
in a browser to review a plan without scrolling a wall of markdown.

The markdown is always the source of truth. The generated `plan.html` is **disposable
output** (gitignored) — never hand-edit it; change the markdown or this script and
re-render.

## Why it exists

Reviewing long markdown plans is tiring. This produces an at-a-glance view: a hero
dashboard (status, RFC revision, deviations, stat tiles, a clickable sub-plan DAG) and
one collapsible, color-coded card per sub-plan with inner tabs (Objective / Context /
Contracts / Acceptance-as-checklist). The full master plan and reviews are kept below,
collapsed, so nothing is lost.

## Usage

```sh
plan-html <plan-dir> [-o <output.html>]
# default output: <plan-dir>/plan.html
```

- Runs from `~/.local/bin/plan-html` after chezmoi apply.
- Implementation lives under `~/.local/share/plan-html`, not in `~/.local/bin`.
- Installs as a `uv tool` from this local package; `uv` generates the executable in
  `~/.local/bin` from the package's `[project.scripts]` entry.
- Output is one self-contained `.html` (inline CSS + JS) — open it directly, works offline.

Example:

```sh
plan-html edge/plans/epics/lcpt-2751-hec-firewall/hotspot-integration
```

## Input contract (what it expects in a plan dir)

| File | Role |
|------|------|
| `00-master.md` | Required. Title (`# ...`) becomes the page title. Its tables drive the dashboard. |
| `NN-*.md` (e.g. `01-foo.md`) | Each becomes a sub-plan card. `NN` is the badge number. |
| `reviews/*.md` | Optional. Each becomes a collapsed review card under "Reference". |

The dashboard is **data-driven** from the master plan, parsed best-effort (missing
pieces just render as blanks):

- `- **RFC Status**: ...` line → RFC chip.
- `## Explicit Deviations` → "deviations: none" chip if the section starts with "None".
- `## Sub-Plans` table → per-card model + dependency badges and the DAG node descriptions.
  Recognized columns (case-insensitive): `#`, `Depends On` or
  `Depends On / Sequenced After`, `Model`, `Description`.
- `## Review Summary` table → review verdict tiles (✓ if status contains "passed").
  Recognized columns: `Reviewer`, `Status`.
- Each sub-plan's `## Primary Files` list → "files touched" tile count.

## How it's built (map for extending it)

Layout:

1. **Install hook** — `run_onchange_after_install-plan-html.sh.tmpl` installs this
   directory with `uv tool install --force ~/.local/share/plan-html` only when the
   package metadata or installed package modules change.
2. **Console script** — `pyproject.toml` declares `plan-html = "plan_html.render:main"`;
   `uv tool` generates the executable in `~/.local/bin`.
3. **Project metadata** — `pyproject.toml` declares pinned parser dependencies for
   IDE/LSP support, packaging, and runtime dependency resolution.
4. **Renderer package** — `src/plan_html/render.py` contains the Markdown renderer,
   plan-specific extraction, CSS, JS, and HTML assembly.
5. **Tests** — `tests/test_render.py` covers parser and RFC-backed plan-template regressions.

Inside `src/plan_html/render.py`, top to bottom:

1. **Markdown engine** — `markdown-it-py` plus `mdit-py-plugins` task lists,
   customized for heading ids and responsive table wrappers.
2. **Structure helpers** — `split_title`, `split_h2_sections`, `section_body`,
   `parse_table_after` (pulls a specific table out of the master using the same
   Markdown parser).
3. **Tabs/cards** — `build_tabs` groups each sub-plan's `##` sections into tabs via
   `TAB_BUCKET` (H2-title → bucket) and `TAB_ORDER`. Unknown sections fall into "Context".
4. **Theme** — `CSS` (CSS custom properties at `:root`), `HUES` (per-card accent colors),
   `JS` (tab switching, DAG-node click-to-open, expand/collapse-all).
5. **Assembly** — `main()` reads the dir, builds the hero + cards, writes the file.

Common tweaks:

- **Palette / colors** — edit the `:root` variables and `HUES` list in `CSS`. (Code blocks
  are intentionally warm orange — `#ffb574` inline, `#ffc88a` in `pre` — to stand out.)
- **Tab grouping** — edit `TAB_BUCKET` / `TAB_ORDER` to add or re-route tabs.
- **New dashboard tiles/chips** — add parsing in `main()` and a tile in the `tiles` string.
- **More markdown** — enable another `markdown-it-py` rule/plugin in `_make_markdown`.

## Limitations (by design or known)

- Requires `uv` at apply/install time so the local package can be installed as a tool.
- Uses CommonMark plus table and task-list plugins, not the full GitHub Markdown feature set.
- Table cells that need a literal pipe should use the standard escaped form, e.g. `` `a\|b` ``.
- No images/footnotes are intentionally styled in the dashboard theme yet.
- Tab grouping keys on known H2 titles; bespoke section names land in "Context".

## Reusing it elsewhere

Copy this directory anywhere that has `uv` available and install it with:

```sh
uv tool install --force /path/to/plan-html
plan-html <plan-dir>
```

Nothing here is tied to this repository except the chezmoi install hook.

If you later want it as a Claude skill, this README is the spec: wrap the script in a
skill that calls it and points the user at the produced `plan.html`.

## Resuming work in a new session

Open a session anywhere with this folder available and paste:

> I have a Python project at `~/.local/share/plan-html` that renders a markdown plan
> directory into an interactive HTML review view. Read it and its `README.md` first.
> I want to <describe change — e.g. add a light theme / a new tab / support images /
> point it at a different plan dir>. Keep generated HTML self-contained; the markdown
> stays the source of truth and `plan.html` stays disposable/gitignored.

That's enough for a fresh session to load the full design and continue.
