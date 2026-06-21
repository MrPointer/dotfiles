"""Render a plan directory into a single self-contained, interactive HTML review view.

Markdown remains the source of truth; this HTML is disposable output (gitignored).
Produces a "dashboard + cards" view: a summary hero (status / RFC / deviations + stat
tiles + a clickable mini-DAG of the sub-plans), each sub-plan as a colored collapsible
card with inner tabs (Objective / Context / Contracts / Acceptance checklist), and the
full master plan + reviews kept available below.

The generated HTML is fully self-contained and works offline after rendering.

Usage:
    plan-html <plan-dir> [-o <output.html>]
"""

import argparse
import glob
import html
import json
import os
import re
from importlib.resources import files

from plan_html.graph import graph_to_dict, parse_dependency_graph
from plan_html.markdown import (
    MarkdownRenderer,
    section_body,
)


# ----------------------------- structure helpers -----------------------------
def split_title(md):
    lines = md.split("\n")
    for idx, line in enumerate(lines):
        m = re.match(r"^#\s+(.*)$", line)
        if m:
            return m.group(1).strip(), "\n".join(lines[idx + 1 :])
    return None, md


def strip_code(s):
    return s.replace("`", "").strip()


def _table_index(header):
    idx = {}
    for k, h in enumerate(header):
        name = re.sub(r"\s+", " ", strip_code(h).lower()).strip()
        idx[name] = k
        idx.setdefault(name.split(" / ", 1)[0], k)
    return idx


# ----------------------------- assets / theme -----------------------------
HUES = ["#6366f1", "#0d9488", "#7c3aed", "#db2777", "#ea580c"]
CSS_ASSETS = ["assets/themes.css", "assets/base.css"]
JS_ASSETS = ["assets/app.js"]
GRAPH_JS_ASSETS = ["assets/dagre.min.js"]


def _asset_text(path):
    return files("plan_html").joinpath(*path.split("/")).read_text(encoding="utf-8")


def render_styles():
    return "\n".join(_asset_text(path) for path in CSS_ASSETS)


def render_scripts(include_graph=False):
    assets = [*GRAPH_JS_ASSETS, *JS_ASSETS] if include_graph else JS_ASSETS
    return "\n".join(_asset_text(path) for path in assets)


def _json_script(data):
    return json.dumps(data).replace("</", "<\\/")


THEME_CONTROL = """
<div class="theme-control" aria-label="Theme">
<button type="button" data-theme-choice="system" aria-pressed="true">System</button>
<button type="button" data-theme-choice="light" aria-pressed="false">Light</button>
<button type="button" data-theme-choice="dark" aria-pressed="false">Dark</button>
</div>
""".strip()


# ----------------------------- assembly -----------------------------
def render_plan(plan_dir, output=None):
    markdown = MarkdownRenderer()

    plan_dir = os.path.abspath(os.fspath(plan_dir))
    out_path = os.fspath(output) if output else os.path.join(plan_dir, "plan.html")
    master_path = os.path.join(plan_dir, "00-master.md")
    if not os.path.exists(master_path):
        raise SystemExit(f"No 00-master.md in {plan_dir}")

    sub_paths = sorted(
        p
        for p in glob.glob(os.path.join(plan_dir, "[0-9][0-9]-*.md"))
        if os.path.basename(p) != "00-master.md"
    )
    review_paths = sorted(glob.glob(os.path.join(plan_dir, "reviews", "*.md")))

    with open(master_path, encoding="utf-8") as f:
        master_md = f.read()
    page_title, master_rest = split_title(master_md)
    page_title = page_title or os.path.basename(plan_dir)

    # --- metadata for the hero ---
    rfc_status = ""
    m = re.search(r"^-\s+\*\*RFC Status\*\*:\s*(.+)$", master_md, re.M)
    if m:
        rfc_status = strip_code(m.group(1))
    dev_body = section_body(master_md, "Explicit Deviations")
    deviations_none = bool(re.match(r"^\s*\**none", dev_body.strip(), re.I))

    subplan_meta = {}  # number -> dict
    parsed = markdown.parse_table_after(master_md, "Sub-Plans")
    if parsed:
        header, rows = parsed
        idx = _table_index(header)
        for r in rows:
            num = strip_code(r[idx.get("#", 0)]) if r else ""

            def col(name, default=""):
                k = idx.get(name)
                return strip_code(r[k]) if k is not None and k < len(r) else default

            subplan_meta[num] = {
                "deps": col("depends on", "-"),
                "model": col("model"),
                "desc": col("description"),
            }

    reviews_meta = []
    parsed = markdown.parse_table_after(master_md, "Review Summary")
    if parsed:
        header, rows = parsed
        idx = _table_index(header)
        for r in rows:
            rev = strip_code(r[idx.get("reviewer", 0)]) if r else ""
            stat = strip_code(r[idx.get("status", 1)]) if len(r) > 1 else ""
            reviews_meta.append((rev, stat))

    known_subplans = {os.path.basename(p)[:2] for p in sub_paths}

    # --- sub-plan cards + dag nodes ---
    total_files = 0
    subplan_hues = {}
    dag_nodes, cards = [], []
    for k, p in enumerate(sub_paths):
        with open(p, encoding="utf-8") as f:
            md = f.read()
        title, rest = split_title(md)
        base = os.path.basename(p)
        num = base[:2]
        title = title or base
        short = re.sub(r"^Sub-Plan:\s*", "", title)
        meta = subplan_meta.get(num, {})
        hue = HUES[k % len(HUES)]
        subplan_hues[num] = hue
        card_id = f"sp-{num}"

        pf = section_body(md, "Primary Files")
        nfiles = len(re.findall(r"^\s*-\s+`?[\w./-]+`?", pf, re.M))
        total_files += nfiles

        model = meta.get("model", "")
        deps = meta.get("deps", "-")
        dep_txt = "no deps" if deps in ("-", "", "–") else f"after {deps}"
        badges = (
            f'<span class="b model">{html.escape(model)}</span>' if model else ""
        ) + f'<span class="b dep">{html.escape(dep_txt)}</span>'

        dag_nodes.append(
            f'<div class="node" data-target="{card_id}" style="--hue:{hue}">'
            f'<div class="nh"><span class="dot">{html.escape(num)}</span>'
            f'<span class="nm">{html.escape(short.split(" ")[0] if " " in short else short)}</span></div>'
            f'<div class="nd">{html.escape(meta.get("desc", "")[:90])}</div>'
            f'<div class="badges">{badges}</div></div>'
        )

        cards.append(
            f'<details class="card" id="{card_id}" open style="--hue:{hue}">'
            f'<summary><span class="seq">{html.escape(num)}</span>'
            f'<span class="title">{html.escape(short)}</span>'
            f'<span class="hb">{badges}</span></summary>'
            f'<div class="inner">{markdown.build_tabs(card_id, rest)}</div></details>'
        )

    plan_graph = parse_dependency_graph(master_md, known_subplans)
    if plan_graph and plan_graph.nodes:
        has_graph = True
        graph_data = graph_to_dict(plan_graph)
        for node in graph_data["nodes"]:
            node["hue"] = subplan_hues.get(node["subplan"], HUES[0])
        graph_warnings = "".join(
            f'<div class="graph-warning">{html.escape(warning)}</div>'
            for warning in plan_graph.warnings
        )
        dag_html = (
            '<div class="graph-surface" data-graph-surface>'
            '<div class="graph-toolbar" data-graph-controls>'
            '<div class="graph-title">Execution Order</div>'
            '<div class="graph-actions">'
            '<button type="button" data-graph-action="fit">Fit</button>'
            '<button type="button" data-graph-action="zoom-in">+</button>'
            '<button type="button" data-graph-action="zoom-out">-</button>'
            '</div>'
            '</div>'
            '<div class="graph-canvas" data-graph-canvas></div>'
            f'{graph_warnings}'
            f'<script type="application/json" data-plan-graph>{_json_script(graph_data)}</script>'
            '</div>'
        )
    else:
        has_graph = False
        dag_html = f'<div class="dag">{"".join(dag_nodes)}</div>' if dag_nodes else ""

    # tiles
    rev_tiles = ""
    for rev, stat in reviews_meta:
        short = rev.replace("plan-", "").replace("-reviewer", "")
        ok = "passed" in stat.lower()
        cls = "ok" if ok and "concern" not in stat.lower() else "warn"
        mark = "✓" if ok else "✗"
        rev_tiles += (
            f'<div class="tile {cls}"><div class="num">{mark}</div>'
            f'<div class="lbl">{html.escape(short)}</div></div>'
        )
    tiles = (
        f'<div class="tile"><div class="num">{len(sub_paths)}</div><div class="lbl">sub-plans</div></div>'
        f'<div class="tile"><div class="num">{total_files}</div><div class="lbl">files touched</div></div>'
        + rev_tiles
    )

    chips = ""
    if rfc_status:
        chips += f'<span class="chip">RFC: {html.escape(rfc_status)}</span>'
    chips += (
        '<span class="chip dev">deviations: none</span>'
        if deviations_none
        else '<span class="chip warn">has deviations</span>'
    )

    master_card = (
        '<details class="card master" id="master" style="--hue:#f59e0b">'
        '<summary><span class="seq">M</span><span class="title">Full master plan</span>'
        '<span class="hb"><span class="b dep">orchestration</span></span></summary>'
        f'<div class="inner">{markdown.render_markdown(master_rest)}</div></details>'
    )

    review_cards = ""
    for p in review_paths:
        with open(p, encoding="utf-8") as f:
            md = f.read()
        title, rest = split_title(md)
        title = title or os.path.basename(p)
        rid = "rv-" + markdown.slugify(os.path.basename(p))
        review_cards += (
            f'<details class="card review" id="{rid}" style="--hue:#475569">'
            f'<summary><span class="seq">⚑</span>'
            f'<span class="title">{html.escape(title)}</span></summary>'
            f'<div class="inner">{markdown.render_markdown(rest)}</div></details>'
        )

    doc = f"""<!doctype html>
<html lang="en"><head><meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>{html.escape(page_title)}</title><style>{render_styles()}</style></head>
<body><div class="wrap">
<div class="hero">
<h1>{html.escape(page_title)}</h1>
<div class="sub">interactive review view · source of truth is the markdown in <code>{html.escape(os.path.basename(plan_dir))}/</code></div>
<div class="chips">{chips}</div>
<div class="tiles">{tiles}</div>
{dag_html}
</div>
<div class="toolbar"><div class="toolbar-actions"><button id="expand">Expand all</button><button id="collapse">Collapse all</button></div>{THEME_CONTROL}</div>
<div class="section-h">Sub-plans</div>
{"".join(cards)}
<div class="section-h">Reference</div>
{master_card}
{review_cards}
</div><script>{render_scripts(has_graph)}</script></body></html>"""

    with open(out_path, "w", encoding="utf-8") as f:
        f.write(doc)
    return out_path


def main(argv=None):
    ap = argparse.ArgumentParser()
    ap.add_argument("plan_dir")
    ap.add_argument("-o", "--output", default=None)
    args = ap.parse_args(argv)

    out_path = render_plan(args.plan_dir, args.output)
    print(out_path)


if __name__ == "__main__":
    main()
