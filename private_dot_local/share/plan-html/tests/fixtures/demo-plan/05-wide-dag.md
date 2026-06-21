# Sub-Plan: Wider DAG Fixture

## Objective

Keep the dashboard usable when a plan has more sub-plans than the initial three-card demo.

## Context

Five nodes are enough to show whether the DAG behaves like a responsive navigation surface or a single row that only works for short plans.

## Primary Files

- `src/plan_html/render.py`
- `src/plan_html/assets/base.css`
- `tests/fixtures/demo-plan/00-master.md`
- `tests/fixtures/demo-plan/05-wide-dag.md`

## Acceptance Criteria

- [ ] The demo fixture renders five clickable DAG nodes.
- [ ] DAG nodes can wrap cleanly across rows.
- [ ] No standalone arrow separator can wrap onto a row by itself.
