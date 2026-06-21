#!/usr/bin/env python3
"""Regression tests for plan_html.render."""

import pathlib
import sys
from importlib.resources import files

from plan_html import render as tool
from plan_html.graph import parse_dependency_graph
from plan_html.markdown import parse_table_after, render_markdown

PROJECT_ROOT = pathlib.Path(__file__).resolve().parents[1]
sys.path.insert(0, str(PROJECT_ROOT / "src"))


def test_escaped_pipe_in_code_span_stays_in_table_cell():
    rendered = render_markdown(
        "| Command | Purpose |\n"
        "|---------|---------|\n"
        "| `foo\\|bar` | keep one cell |\n"
    )

    assert "<td><code>foo|bar</code></td>" in rendered
    assert "<td>keep one cell</td>" in rendered
    assert rendered.count("<td") == 2


def test_plan_table_extraction_handles_escaped_pipe():
    parsed = parse_table_after(
        "# Master Plan\n\n"
        "## Sub-Plans\n\n"
        "| # | Description |\n"
        "|---|-------------|\n"
        "| 01 | Run `foo\\|bar` safely |\n",
        "Sub-Plans",
    )

    assert parsed is not None
    header, rows = parsed
    assert header == ["#", "Description"]
    assert rows == [["01", "Run `foo|bar` safely"]]


def test_dependency_graph_parser_accepts_strict_sp_mermaid_subset():
    parsed = parse_dependency_graph(
        "# Master Plan\n\n"
        "## Dependency Graph\n\n"
        "```mermaid\n"
        "flowchart TD\n"
        "  SP01[\"01 Foundation\"]\n"
        "  SP02[\"02 First Branch\"]\n"
        "  SP03[\"03 Second Branch\"]\n"
        "  SP01 --> SP02\n"
        "  SP01 --> SP03\n"
        "```\n",
        {"01", "02", "03"},
    )

    assert parsed is not None
    assert parsed.direction == "TD"
    assert [(node.id, node.subplan) for node in parsed.nodes] == [
        ("SP01", "01"),
        ("SP02", "02"),
        ("SP03", "03"),
    ]
    assert [(edge.source, edge.target) for edge in parsed.edges] == [
        ("SP01", "SP02"),
        ("SP01", "SP03"),
    ]
    assert parsed.warnings == []


def test_dependency_graph_parser_reports_unsupported_lines():
    parsed = parse_dependency_graph(
        "# Master Plan\n\n"
        "## Dependency Graph\n\n"
        "```mermaid\n"
        "flowchart LR\n"
        "  SP01[\"01 Foundation\"]\n"
        "  SP01 -.-> SP02\n"
        "```\n",
        {"01", "02"},
    )

    assert parsed is not None
    assert parsed.warnings == ["Unsupported graph line: SP01 -.-> SP02"]


def test_render_plan_uses_rfc_plan_template_dependency_column(tmp_path):
    plan_dir = tmp_path
    (plan_dir / "00-master.md").write_text(
        "# Master Plan: Demo\n\n"
        "## RFC Baseline\n"
        "- **RFC Status**: Accepted\n\n"
        "## Explicit Deviations\n\n"
        "None\n\n"
        "## Sub-Plans\n\n"
        "| # | Sub-Plan | Depends On / Sequenced After | Model | Description |\n"
        "|---|----------|------------------------------|-------|-------------|\n"
        "| 01 | `01-first.md` | - | Mid-tier | First step |\n"
        "| 02 | `02-second.md` | 01 (logical dependency) | Cheapest | Second step |\n",
        encoding="utf-8",
    )
    (plan_dir / "01-first.md").write_text(
        "# Sub-Plan: First\n\n"
        "## Objective\n\nFirst.\n\n"
        "## Primary Files\n\n- `first.txt`\n",
        encoding="utf-8",
    )
    (plan_dir / "02-second.md").write_text(
        "# Sub-Plan: Second\n\n"
        "## Objective\n\nSecond.\n\n"
        "## Primary Files\n\n- `second.txt`\n",
        encoding="utf-8",
    )

    out_path = plan_dir / "plan.html"
    tool.render_plan(plan_dir, out_path)
    rendered = out_path.read_text(encoding="utf-8")

    assert "after 01 (logical dependency)" in rendered


def test_render_plan_includes_system_light_dark_theme_control(tmp_path):
    plan_dir = tmp_path
    (plan_dir / "00-master.md").write_text(
        "# Master Plan: Themes\n\n## Explicit Deviations\n\nNone\n\n",
        encoding="utf-8",
    )

    out_path = plan_dir / "plan.html"
    tool.render_plan(plan_dir, out_path)
    rendered = out_path.read_text(encoding="utf-8")

    assert 'class="theme-control"' in rendered
    assert 'data-theme-choice="system"' in rendered
    assert 'data-theme-choice="light"' in rendered
    assert 'data-theme-choice="dark"' in rendered
    assert "prefers-color-scheme:dark" in rendered
    assert "document.documentElement.dataset.theme" in rendered
    assert "localStorage" not in rendered


def test_render_plan_renders_same_plan_deterministically_in_one_process(tmp_path):
    plan_dir = tmp_path
    (plan_dir / "00-master.md").write_text(
        "# Master Plan: Deterministic\n\n"
        "## Explicit Deviations\n\n"
        "None\n\n"
        "## Repeated Heading\n\n"
        "Master body.\n\n"
        "## Repeated Heading\n\n"
        "Master body again.\n",
        encoding="utf-8",
    )
    (plan_dir / "01-first.md").write_text(
        "# Sub-Plan: First\n\n"
        "## Objective\n\n"
        "First.\n\n"
        "## Repeated Heading\n\n"
        "Sub-plan body.\n",
        encoding="utf-8",
    )

    one_path = plan_dir / "one.html"
    two_path = plan_dir / "two.html"
    tool.render_plan(plan_dir, one_path)
    tool.render_plan(plan_dir, two_path)
    rendered = [
        one_path.read_text(encoding="utf-8"),
        two_path.read_text(encoding="utf-8"),
    ]

    assert rendered[0] == rendered[1]
    assert 'id="repeated-heading"' in rendered[0]
    assert 'id="repeated-heading-1"' in rendered[0]


def test_fixture_demo_plan_renders_representative_review_view(tmp_path):
    fixture_dir = PROJECT_ROOT / "tests" / "fixtures" / "demo-plan"

    out_path = tmp_path / "demo-plan.html"
    tool.render_plan(fixture_dir, out_path)
    rendered = out_path.read_text(encoding="utf-8")

    assert "Master Plan: Plan HTML Demo" in rendered
    assert 'data-theme-choice="system"' in rendered
    assert 'data-theme-choice="light"' in rendered
    assert 'data-theme-choice="dark"' in rendered
    assert '<div class="num">5</div><div class="lbl">sub-plans</div>' in rendered
    assert 'class="graph-surface"' in rendered
    assert 'data-graph-controls' in rendered
    assert 'class="graph-title">Execution Order</div>' in rendered
    assert 'data-plan-graph' in rendered
    assert '"id": "SP01"' in rendered
    assert '"id": "SP01", "label": "01 Foundation", "subplan": "01", "hue": "#6366f1"' in rendered
    assert '"id": "SP02", "label": "02 Theme Control", "subplan": "02", "hue": "#0d9488"' in rendered
    assert '"source": "SP01", "target": "SP03"' in rendered
    assert '"source": "SP02", "target": "SP04"' in rendered
    assert '"source": "SP03", "target": "SP04"' in rendered
    assert '"source": "SP02", "target": "SP03"' not in rendered
    assert "mermaid.initialize" not in rendered
    assert "mermaid.min" not in rendered
    assert "var dagre=" in rendered
    assert "pointerdown" in rendered
    assert "is-panning" in rendered
    assert "cdn" not in rendered.lower()
    assert "after 01" in rendered
    assert "after 02, 03" in rendered
    assert "after 04" in rendered
    assert "<td><code>foo|bar</code></td>" in rendered
    assert "generated_preview_fixture_with_a_really_long_component_name" in rendered
    assert "Acceptance <span class='cnt'>3</span>" in rendered
    assert "clarity" in rendered
    assert "Manual Smoke Checks" in rendered


def test_fixture_demo_plan_embeds_long_content_and_dag_layout_guards(tmp_path):
    fixture_dir = PROJECT_ROOT / "tests" / "fixtures" / "demo-plan"

    out_path = tmp_path / "demo-plan.html"
    tool.render_plan(fixture_dir, out_path)
    rendered = out_path.read_text(encoding="utf-8")

    assert rendered.count('class="node"') == 0
    assert 'class="graph-surface"' in rendered
    assert 'class="arrow"' not in rendered
    assert "overflow-wrap:anywhere" in rendered
    assert "graph.js" not in rendered


def test_render_plan_without_dependency_graph_uses_node_grid_fallback(tmp_path):
    plan_dir = tmp_path
    (plan_dir / "00-master.md").write_text(
        "# Master Plan: Fallback\n\n"
        "## Explicit Deviations\n\n"
        "None\n\n"
        "## Sub-Plans\n\n"
        "| # | Sub-Plan | Depends On / Sequenced After | Model | Description |\n"
        "|---|----------|------------------------------|-------|-------------|\n"
        "| 01 | `01-first.md` | - | Mid-tier | First step |\n"
        "| 02 | `02-second.md` | 01 | Cheapest | Second step |\n",
        encoding="utf-8",
    )
    (plan_dir / "01-first.md").write_text("# Sub-Plan: First\n", encoding="utf-8")
    (plan_dir / "02-second.md").write_text("# Sub-Plan: Second\n", encoding="utf-8")

    out_path = plan_dir / "plan.html"
    tool.render_plan(plan_dir, out_path)
    rendered = out_path.read_text(encoding="utf-8")

    assert rendered.count('class="node"') == 2
    assert 'class="graph-surface"' not in rendered
    assert '<script type="application/json" data-plan-graph>' not in rendered
    assert "var dagre=" not in rendered


def test_render_plan_displays_dependency_graph_warnings(tmp_path):
    plan_dir = tmp_path
    (plan_dir / "00-master.md").write_text(
        "# Master Plan: Graph Warning\n\n"
        "## Dependency Graph\n\n"
        "```mermaid\n"
        "flowchart LR\n"
        "  SP01[\"01 First\"]\n"
        "  SP01 -.-> SP02\n"
        "```\n",
        encoding="utf-8",
    )
    (plan_dir / "01-first.md").write_text("# Sub-Plan: First\n", encoding="utf-8")

    out_path = plan_dir / "plan.html"
    tool.render_plan(plan_dir, out_path)
    rendered = out_path.read_text(encoding="utf-8")

    assert 'class="graph-warning"' in rendered
    assert "Unsupported graph line: SP01 -.-&gt; SP02" in rendered


def test_main_accepts_argv_and_prints_output_path(tmp_path, capsys):
    plan_dir = tmp_path
    (plan_dir / "00-master.md").write_text("# Master Plan: CLI\n", encoding="utf-8")
    out_path = plan_dir / "cli.html"

    tool.main([str(plan_dir), "-o", str(out_path)])

    assert capsys.readouterr().out.strip() == str(out_path)
    assert out_path.exists()


def test_assets_are_packaged_with_plan_html():
    package_files = files("plan_html")

    assert package_files.joinpath("assets/base.css").is_file()
    assert package_files.joinpath("assets/themes.css").is_file()
    assert package_files.joinpath("assets/app.js").is_file()
    assert package_files.joinpath("assets/dagre.min.js").is_file()
