# Sub-Plan: Foundation Assets

## Objective

Separate maintained CSS and JS source from the Python renderer while keeping the generated review view self-contained.

## Context

The generated `plan.html` remains disposable output. The source markdown and package assets are the inputs that should be changed by humans.

## Primary Files

- `src/plan_html/render.py`
- `src/plan_html/assets/base.css`
- `src/plan_html/assets/app.js`

## Integration Contracts

| Command | Purpose |
|---------|---------|
| `foo\|bar` | Keeps escaped pipes inside a single table cell |
| `uv build --wheel` | Confirms package assets are included in the installable tool |

## Acceptance Criteria

- [x] CSS source assets are maintained outside the Python source.
- [x] JS source assets are maintained outside the Python source.
- [ ] Generated HTML remains a single offline file.
