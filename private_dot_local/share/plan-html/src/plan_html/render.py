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
import os
import re
import sys

try:
    from markdown_it import MarkdownIt  # pyright: ignore[reportMissingImports]
    from mdit_py_plugins.tasklists import tasklists_plugin  # pyright: ignore[reportMissingImports]
except ImportError as exc:
    raise SystemExit(
        "render_plan_html.py requires markdown-it-py and mdit-py-plugins. "
        "Run it through `plan-html` so the pinned requirements are available."
    ) from exc

_slugs = {}


def slugify(text):
    s = re.sub(r"<[^>]+>", "", text)
    s = re.sub(r"`", "", s)
    s = re.sub(r"[^\w\s-]", "", s).strip().lower()
    s = re.sub(r"\s+", "-", s)
    if not s:
        s = "section"
    if s in _slugs:
        _slugs[s] += 1
        s = f"{s}-{_slugs[s]}"
    else:
        _slugs[s] = 0
    return s


# ----------------------------- markdown engine -----------------------------
def _make_markdown():
    markdown = MarkdownIt("commonmark", {"html": False}).enable("table")
    markdown.use(tasklists_plugin, enabled=False)

    def heading_open(tokens, idx, options, env):
        if idx + 1 < len(tokens):
            tokens[idx].attrSet("id", slugify(tokens[idx + 1].content))
        return markdown.renderer.renderToken(tokens, idx, options, env)

    def table_open(tokens, idx, options, env):
        return '<div class="tablewrap">' + markdown.renderer.renderToken(tokens, idx, options, env)

    def table_close(tokens, idx, options, env):
        return markdown.renderer.renderToken(tokens, idx, options, env) + "</div>"

    markdown.renderer.rules["heading_open"] = heading_open
    markdown.renderer.rules["table_open"] = table_open
    markdown.renderer.rules["table_close"] = table_close
    return markdown


MARKDOWN = _make_markdown()


def render_markdown(md):
    return MARKDOWN.render(md)


def _first_table(md):
    header = None
    rows = []
    current_row = None
    section = None
    in_table = False

    for token in MARKDOWN.parse(md):
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


# ----------------------------- structure helpers -----------------------------
def split_title(md):
    lines = md.split("\n")
    for idx, l in enumerate(lines):
        m = re.match(r"^#\s+(.*)$", l)
        if m:
            return m.group(1).strip(), "\n".join(lines[idx + 1 :])
    return None, md


def split_h2_sections(md):
    """[(title, body_md)] split on H2; content before first H2 gets title None."""
    sections, title, cur = [], None, []
    for l in md.split("\n"):
        m = re.match(r"^##\s+(.*)$", l)
        if m:
            if title is not None or any(s.strip() for s in cur):
                sections.append((title, "\n".join(cur)))
            title, cur = m.group(1).strip(), []
        else:
            cur.append(l)
    if title is not None or any(s.strip() for s in cur):
        sections.append((title, "\n".join(cur)))
    return sections


def section_body(md, title):
    for t, body in split_h2_sections(md):
        if t and t.lower() == title.lower():
            return body
    return ""


def parse_table_after(md, heading):
    body = section_body(md, heading)
    if not body:
        return None
    return _first_table(body)


def strip_code(s):
    return s.replace("`", "").strip()


def _table_index(header):
    idx = {}
    for k, h in enumerate(header):
        name = re.sub(r"\s+", " ", strip_code(h).lower()).strip()
        idx[name] = k
        idx.setdefault(name.split(" / ", 1)[0], k)
    return idx


# ----------------------------- tabs / cards -----------------------------
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


def build_tabs(card_id, md):
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
            inner.append(render_markdown(body))
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


# ----------------------------- theme -----------------------------
HUES = ["#6366f1", "#0d9488", "#7c3aed", "#db2777", "#ea580c"]

CSS = """
:root{--bg:#0f1117;--panel:#171a23;--card:#1c2030;--fg:#e7e9f0;--muted:#9aa0b4;
--line:#2a2f40;--accent:#818cf8;--ok:#34d399;--warn:#fbbf24;--code:#11141d;}
*{box-sizing:border-box}
html{scroll-behavior:smooth}
body{margin:0;background:var(--bg);color:var(--fg);
font:15px/1.65 -apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,Helvetica,Arial,sans-serif}
.wrap{max-width:1040px;margin:0 auto;padding:30px 26px 120px}
a{color:var(--accent);text-decoration:none}a:hover{text-decoration:underline}
code{background:#2a1c0f;padding:.12em .42em;border-radius:5px;font-size:.86em;
font-family:"SF Mono",ui-monospace,Menlo,Consolas,monospace;color:#ffb574;border:1px solid #4a3a1f}
pre{background:#1c150d;border:1px solid #4a3a1f;border-radius:8px;padding:13px 15px;overflow:auto}
pre code{background:none;padding:0;border:none;color:#ffc88a}
h1,h2,h3{line-height:1.3}
h3{font-size:15px;margin:20px 0 7px;color:#c3c8e0;letter-spacing:.01em}
.tablewrap{overflow:auto;margin:12px 0}
table{border-collapse:collapse;width:100%;font-size:13px}
th,td{border:1px solid var(--line);padding:7px 10px;text-align:left;vertical-align:top}
th{background:#222740;color:#cfd4ee}
tbody tr:nth-child(even){background:#1a1e2b}
hr{border:none;border-top:1px solid var(--line);margin:22px 0}
blockquote{border-left:3px solid var(--accent);margin:12px 0;padding:2px 14px;color:var(--muted);
background:#191d2a;border-radius:0 6px 6px 0}
ul,ol{padding-left:22px}li{margin:4px 0}
.contains-task-list{padding-left:0}.task-list-item{list-style:none}
.task-list-item-checkbox{margin:5px 9px 0 0;width:15px;height:15px;accent-color:var(--ok)}
/* hero */
.hero{background:linear-gradient(135deg,#1e1b4b,#312e81 55%,#1e3a5f);
border:1px solid #3a3a6a;border-radius:16px;padding:24px 26px;margin-bottom:22px;box-shadow:0 10px 40px -18px #000}
.hero h1{margin:0 0 6px;font-size:26px}
.hero .sub{color:#c7cbf0;font-size:13px;margin-bottom:18px}
.chips{display:flex;flex-wrap:wrap;gap:8px;margin-bottom:18px}
.chip{font-size:12px;font-weight:600;padding:5px 11px;border-radius:999px;border:1px solid #ffffff22;
background:#ffffff14;color:#eef}
.chip.ok{background:#064e3b;border-color:#10b98155;color:#6ee7b7}
.chip.warn{background:#4a3408;border-color:#f59e0b55;color:#fcd34d}
.chip.dev{background:#3f1d2e;border-color:#ec489955;color:#f9a8d4}
.tiles{display:grid;grid-template-columns:repeat(auto-fit,minmax(150px,1fr));gap:12px;margin-bottom:20px}
.tile{background:#ffffff10;border:1px solid #ffffff1f;border-radius:12px;padding:14px 16px}
.tile .num{font-size:24px;font-weight:800;line-height:1}
.tile .lbl{font-size:11.5px;color:#c7cbf0;margin-top:5px;text-transform:uppercase;letter-spacing:.05em}
.tile.ok .num{color:#6ee7b7}.tile.warn .num{color:#fcd34d}
/* dag */
.dag{display:flex;align-items:stretch;gap:0;flex-wrap:wrap;margin-top:4px}
.node{flex:1;min-width:150px;cursor:pointer;background:#ffffff10;border:1px solid #ffffff24;
border-top:3px solid var(--hue,#818cf8);border-radius:11px;padding:12px 13px;transition:.15s}
.node:hover{background:#ffffff1c;transform:translateY(-2px)}
.node .nh{display:flex;align-items:center;gap:8px;font-weight:700}
.node .dot{width:22px;height:22px;border-radius:50%;display:grid;place-items:center;
font-size:11px;font-weight:800;color:#0f1117;background:var(--hue,#818cf8)}
.node .nm{font-size:13.5px}
.node .nd{font-size:11.5px;color:#c2c7e6;margin:7px 0 9px;min-height:30px}
.node .badges{display:flex;gap:5px;flex-wrap:wrap}
.arrow{display:flex;align-items:center;color:#818cf8;font-size:22px;padding:0 6px}
@media(max-width:760px){.arrow{transform:rotate(90deg)}.dag{flex-direction:column}}
/* badges */
.b{font-size:11px;font-weight:700;padding:3px 9px;border-radius:6px;white-space:nowrap}
.b.model{background:#0c3b34;color:#5eead4;border:1px solid #14b8a655}
.b.dep{background:#2a2f48;color:#aab2e6;border:1px solid #ffffff1f}
.b.none{background:#3a2f10;color:#fcd34d;border:1px solid #f59e0b44}
/* cards */
.card{background:var(--card);border:1px solid var(--line);border-left:4px solid var(--hue,#6366f1);
border-radius:13px;margin:16px 0;overflow:hidden;box-shadow:0 6px 24px -18px #000}
.card>summary{list-style:none;cursor:pointer;padding:16px 20px;display:flex;align-items:center;
gap:12px;flex-wrap:wrap}
.card>summary::-webkit-details-marker{display:none}
.card>summary::before{content:"\\25B8";color:var(--hue,#818cf8);font-size:14px;transition:transform .15s}
.card[open]>summary::before{transform:rotate(90deg)}
.card .seq{width:26px;height:26px;border-radius:7px;display:grid;place-items:center;font-weight:800;
font-size:12px;color:#0f1117;background:var(--hue,#6366f1)}
.card .title{font-weight:700;font-size:16px;flex:1;min-width:180px}
.card .hb{display:flex;gap:6px;flex-wrap:wrap}
.card .inner{padding:0 22px 20px;border-top:1px solid var(--line)}
.tabs{display:flex;gap:7px;flex-wrap:wrap;margin:14px 0 4px}
.tab{font-size:12.5px;font-weight:650;padding:7px 14px;border-radius:999px;cursor:pointer;
border:1px solid var(--line);background:#11141d;color:var(--muted)}
.tab:hover{color:var(--fg);border-color:var(--accent)}
.tab.active{background:var(--hue,#6366f1);color:#0f1117;border-color:transparent}
.tab .cnt{background:#0f1117;color:#fff;border-radius:999px;padding:0 6px;margin-left:5px;font-size:10px}
.pane{display:none;animation:fade .18s ease}.pane.active{display:block}
@keyframes fade{from{opacity:0;transform:translateY(4px)}to{opacity:1}}
.toolbar{display:flex;gap:8px;margin:18px 0 6px;flex-wrap:wrap}
.toolbar button{font-size:12px;padding:6px 13px;border-radius:8px;cursor:pointer;
border:1px solid var(--line);background:var(--panel);color:var(--muted)}
.toolbar button:hover{color:var(--accent);border-color:var(--accent)}
.section-h{font-size:13px;text-transform:uppercase;letter-spacing:.08em;color:var(--muted);
margin:26px 0 2px;font-weight:700}
.card.review{border-left-color:#475569}.card.review .seq{background:#475569}
.card.master{border-left-color:#f59e0b}.card.master .seq{background:#f59e0b}
"""

JS = """
document.addEventListener('click',e=>{
  const t=e.target.closest('.tab');
  if(t){const card=t.dataset.card,name=t.dataset.tab;
    document.querySelectorAll(`.tab[data-card="${card}"]`).forEach(b=>b.classList.toggle('active',b===t));
    document.querySelectorAll(`.pane[data-card="${card}"]`).forEach(p=>p.classList.toggle('active',p.dataset.pane===name));
    return;}
  const node=e.target.closest('.node');
  if(node){const id=node.dataset.target,el=document.getElementById(id);
    if(el){el.open=true;el.scrollIntoView({behavior:'smooth',block:'start'});}}
});
const ex=document.getElementById('expand'),co=document.getElementById('collapse');
if(ex)ex.onclick=()=>document.querySelectorAll('details.card').forEach(d=>d.open=true);
if(co)co.onclick=()=>document.querySelectorAll('details.card').forEach(d=>d.open=false);
"""


# ----------------------------- assembly -----------------------------
def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("plan_dir")
    ap.add_argument("-o", "--output", default=None)
    args = ap.parse_args()

    plan_dir = os.path.abspath(args.plan_dir)
    out_path = args.output or os.path.join(plan_dir, "plan.html")
    master_path = os.path.join(plan_dir, "00-master.md")
    if not os.path.exists(master_path):
        sys.exit(f"No 00-master.md in {plan_dir}")

    sub_paths = sorted(
        p for p in glob.glob(os.path.join(plan_dir, "[0-9][0-9]-*.md"))
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
    parsed = parse_table_after(master_md, "Sub-Plans")
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
    parsed = parse_table_after(master_md, "Review Summary")
    if parsed:
        header, rows = parsed
        idx = _table_index(header)
        for r in rows:
            rev = strip_code(r[idx.get("reviewer", 0)]) if r else ""
            stat = strip_code(r[idx.get("status", 1)]) if len(r) > 1 else ""
            reviews_meta.append((rev, stat))

    # --- sub-plan cards + dag nodes ---
    total_files = 0
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
        card_id = f"sp-{num}"

        pf = section_body(md, "Primary Files")
        nfiles = len(re.findall(r"^\s*-\s+`?[\w./-]+`?", pf, re.M))
        total_files += nfiles

        model = meta.get("model", "")
        deps = meta.get("deps", "-")
        dep_txt = "no deps" if deps in ("-", "", "–") else f"after {deps}"
        badges = (
            (f'<span class="b model">{html.escape(model)}</span>' if model else "")
            + f'<span class="b dep">{html.escape(dep_txt)}</span>'
        )

        dag_nodes.append(
            f'<div class="node" data-target="{card_id}" style="--hue:{hue}">'
            f'<div class="nh"><span class="dot">{html.escape(num)}</span>'
            f'<span class="nm">{html.escape(short.split(" ")[0] if " " in short else short)}</span></div>'
            f'<div class="nd">{html.escape(meta.get("desc","")[:90])}</div>'
            f'<div class="badges">{badges}</div></div>'
        )

        cards.append(
            f'<details class="card" id="{card_id}" open style="--hue:{hue}">'
            f'<summary><span class="seq">{html.escape(num)}</span>'
            f'<span class="title">{html.escape(short)}</span>'
            f'<span class="hb">{badges}</span></summary>'
            f'<div class="inner">{build_tabs(card_id, rest)}</div></details>'
        )

    # dag with arrows
    dag = []
    for j, node in enumerate(dag_nodes):
        if j:
            dag.append('<div class="arrow">→</div>')
        dag.append(node)
    dag_html = f'<div class="dag">{"".join(dag)}</div>' if dag else ""

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
        f'<div class="inner">{render_markdown(master_rest)}</div></details>'
    )

    review_cards = ""
    for p in review_paths:
        with open(p, encoding="utf-8") as f:
            md = f.read()
        title, rest = split_title(md)
        title = title or os.path.basename(p)
        rid = "rv-" + slugify(os.path.basename(p))
        review_cards += (
            f'<details class="card review" id="{rid}" style="--hue:#475569">'
            f'<summary><span class="seq">⚑</span>'
            f'<span class="title">{html.escape(title)}</span></summary>'
            f'<div class="inner">{render_markdown(rest)}</div></details>'
        )

    doc = f"""<!doctype html>
<html lang="en"><head><meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>{html.escape(page_title)}</title><style>{CSS}</style></head>
<body><div class="wrap">
<div class="hero">
<h1>{html.escape(page_title)}</h1>
<div class="sub">interactive review view · source of truth is the markdown in <code>{html.escape(os.path.basename(plan_dir))}/</code></div>
<div class="chips">{chips}</div>
<div class="tiles">{tiles}</div>
{dag_html}
</div>
<div class="toolbar"><button id="expand">Expand all</button><button id="collapse">Collapse all</button></div>
<div class="section-h">Sub-plans</div>
{''.join(cards)}
<div class="section-h">Reference</div>
{master_card}
{review_cards}
</div><script>{JS}</script></body></html>"""

    with open(out_path, "w", encoding="utf-8") as f:
        f.write(doc)
    print(out_path)


if __name__ == "__main__":
    main()
