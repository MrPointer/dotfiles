"""Markdown rendering helpers for plan-html."""

import html
import re
from collections.abc import Sequence

try:
    from markdown_it import MarkdownIt  # pyright: ignore[reportMissingImports]
    from markdown_it.renderer import (
        RendererHTML,  # pyright: ignore[reportMissingImports]
    )
    from markdown_it.token import Token  # pyright: ignore[reportMissingImports]
    from markdown_it.utils import (  # pyright: ignore[reportMissingImports]
        EnvType,
        OptionsDict,
    )
    from mdit_py_plugins.tasklists import (
        tasklists_plugin,  # pyright: ignore[reportMissingImports]
    )
except ImportError as exc:
    raise SystemExit(
        "render_plan_html.py requires markdown-it-py and mdit-py-plugins. "
        "Run it through `plan-html` so the pinned requirements are available."
    ) from exc


TAB_BUCKET = {
    "objective": "Objective",
    "rfc context": "Context",
    "required skills": "Context",
    "execution model": "Context",
    "prerequisites": "Context",
    "context": "Context",
    "design decisions": "Context",
    "integration contracts": "Contracts",
    "primary files": "Contracts",
    "acceptance criteria": "Acceptance",
}
TAB_ORDER = ["Objective", "Context", "Contracts", "Acceptance"]


def split_h2_sections(md):
    """[(title, body_md)] split on H2; content before first H2 gets title None."""
    sections, title, cur = [], None, []
    for line in md.split("\n"):
        m = re.match(r"^##\s+(.*)$", line)
        if m:
            if title is not None or any(s.strip() for s in cur):
                sections.append((title, "\n".join(cur)))
            title, cur = m.group(1).strip(), []
        else:
            cur.append(line)
    if title is not None or any(s.strip() for s in cur):
        sections.append((title, "\n".join(cur)))
    return sections


def section_body(md, title):
    for t, body in split_h2_sections(md):
        if t and t.lower() == title.lower():
            return body
    return ""


class MarkdownRenderer:
    def __init__(self):
        self._slugs = {}
        self.markdown = self._make_markdown()

    def slugify(self, text):
        s = re.sub(r"<[^>]+>", "", text)
        s = re.sub(r"`", "", s)
        s = re.sub(r"[^\w\s-]", "", s).strip().lower()
        s = re.sub(r"\s+", "-", s)
        if not s:
            s = "section"
        if s in self._slugs:
            self._slugs[s] += 1
            s = f"{s}-{self._slugs[s]}"
        else:
            self._slugs[s] = 0
        return s

    def _make_markdown(self):
        markdown = MarkdownIt("commonmark", {"html": False}).enable("table")
        markdown.use(tasklists_plugin, enabled=False)

        def heading_open(
            renderer: RendererHTML,
            tokens: Sequence[Token],
            idx: int,
            options: OptionsDict,
            env: EnvType,
        ) -> str:
            if idx + 1 < len(tokens):
                tokens[idx].attrSet("id", self.slugify(tokens[idx + 1].content))
            return renderer.renderToken(tokens, idx, options, env)

        def table_open(
            renderer: RendererHTML,
            tokens: Sequence[Token],
            idx: int,
            options: OptionsDict,
            env: EnvType,
        ) -> str:
            return '<div class="tablewrap">' + renderer.renderToken(
                tokens, idx, options, env
            )

        def table_close(
            renderer: RendererHTML,
            tokens: Sequence[Token],
            idx: int,
            options: OptionsDict,
            env: EnvType,
        ) -> str:
            return renderer.renderToken(tokens, idx, options, env) + "</div>"

        markdown.add_render_rule("heading_open", heading_open)
        markdown.add_render_rule("table_open", table_open)
        markdown.add_render_rule("table_close", table_close)
        return markdown

    def render_markdown(self, md):
        return self.markdown.render(md)

    def first_table(self, md):
        header = None
        rows = []
        current_row = None
        section = None
        in_table = False

        for token in self.markdown.parse(md):
            if token.type == "table_open":
                in_table = True
                continue
            if not in_table:
                continue
            if token.type == "table_close":
                break
            if token.type == "thead_open":
                section = "head"
                continue
            if token.type == "tbody_open":
                section = "body"
                continue
            if token.type == "tr_open":
                current_row = []
                continue
            if token.type == "inline" and current_row is not None:
                current_row.append(token.content.strip())
                continue
            if token.type == "tr_close" and current_row is not None:
                if section == "head":
                    header = current_row
                elif section == "body":
                    rows.append(current_row)
                current_row = None

        if header is None:
            return None
        return header, rows

    def parse_table_after(self, md, heading):
        body = section_body(md, heading)
        if not body:
            return None
        return self.first_table(body)

    def build_tabs(self, card_id, md):
        buckets = {k: [] for k in TAB_ORDER}
        for title, body in split_h2_sections(md):
            if not title:
                if body.strip():
                    buckets["Objective"].append(("", body))
                continue
            bucket = TAB_BUCKET.get(title.lower(), "Context")
            buckets[bucket].append((title, body))
        active_set = False
        btns, panes = [], []
        for name in TAB_ORDER:
            chunks = buckets[name]
            if not chunks:
                continue
            inner = []
            for title, body in chunks:
                if title and title.lower() != "objective":
                    inner.append(f"<h3>{html.escape(title)}</h3>")
                inner.append(self.render_markdown(body))
            active = "" if active_set else " active"
            label = name
            if name == "Acceptance":
                cnt = len(re.findall(r"^\s*-\s*\[[ xX]\]", md, re.M))
                label = f"Acceptance <span class='cnt'>{cnt}</span>" if cnt else name
            btns.append(
                f'<button class="tab{active}" data-card="{card_id}" data-tab="{name}">{label}</button>'
            )
            panes.append(
                f'<div class="pane{active}" data-card="{card_id}" data-pane="{name}">{"".join(inner)}</div>'
            )
            active_set = True
        return f'<div class="tabs">{"".join(btns)}</div><div class="panes">{"".join(panes)}</div>'


def render_markdown(md):
    return MarkdownRenderer().render_markdown(md)


def parse_table_after(md, heading):
    return MarkdownRenderer().parse_table_after(md, heading)


def build_tabs(card_id, md):
    return MarkdownRenderer().build_tabs(card_id, md)
