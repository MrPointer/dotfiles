#!/usr/bin/env python3
"""Regression tests for plan_html.render."""

import pathlib
import sys
from importlib.resources import files

from plan_html import render as tool
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
    assert "after 01" in rendered
    assert "after 02" in rendered
    assert "<td><code>foo|bar</code></td>" in rendered
    assert "Acceptance <span class='cnt'>3</span>" in rendered
    assert "clarity" in rendered
    assert "Manual Smoke Checks" in rendered


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
