#!/usr/bin/env python3
"""Regression tests for plan_html.render."""

import contextlib
import io
import pathlib
import sys
import tempfile
import unittest


PROJECT_ROOT = pathlib.Path(__file__).resolve().parents[1]
sys.path.insert(0, str(PROJECT_ROOT / "src"))

from plan_html import render as tool


class MarkdownRenderingTest(unittest.TestCase):
    def test_escaped_pipe_in_code_span_stays_in_table_cell(self):
        rendered = tool.render_markdown(
            "| Command | Purpose |\n"
            "|---------|---------|\n"
            "| `foo\\|bar` | keep one cell |\n"
        )

        self.assertIn("<td><code>foo|bar</code></td>", rendered)
        self.assertIn("<td>keep one cell</td>", rendered)
        self.assertEqual(rendered.count("<td"), 2)

    def test_plan_table_extraction_handles_escaped_pipe(self):
        header, rows = tool.parse_table_after(
            "# Master Plan\n\n"
            "## Sub-Plans\n\n"
            "| # | Description |\n"
            "|---|-------------|\n"
            "| 01 | Run `foo\\|bar` safely |\n",
            "Sub-Plans",
        )

        self.assertEqual(header, ["#", "Description"])
        self.assertEqual(rows, [["01", "Run `foo|bar` safely"]])

    def test_main_uses_rfc_plan_template_dependency_column(self):
        with tempfile.TemporaryDirectory() as tmpdir:
            plan_dir = pathlib.Path(tmpdir)
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
            old_argv = sys.argv
            sys.argv = ["render_plan_html.py", str(plan_dir), "-o", str(out_path)]
            try:
                with contextlib.redirect_stdout(io.StringIO()):
                    tool.main()
            finally:
                sys.argv = old_argv

            rendered = out_path.read_text(encoding="utf-8")

        self.assertIn("after 01 (logical dependency)", rendered)


if __name__ == "__main__":
    unittest.main()
