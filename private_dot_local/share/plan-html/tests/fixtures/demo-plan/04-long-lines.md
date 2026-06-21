# Sub-Plan: Long Line Resilience

## Objective

Keep cards readable when a real plan contains long unbroken inline code such as `src/plan_html/assets/generated_preview_fixture_with_a_really_long_component_name_and_no_natural_breakpoints_for_wrapping_behavior.css` in prose.

## Context

Some implementation plans include generated names, command arguments, cache keys, or file paths that are too long for normal word wrapping. Those strings should not force the surrounding card wider than the viewport.

## Integration Contracts

The renderer should also tolerate long commands in code blocks without making nearby paragraphs unreadable:

```sh
uv run plan-html tests/fixtures/demo-plan -o /tmp/plan-html-demo-with-a-deliberately-long-output-file-name-used-for-layout-regression-checks.html
```

## Primary Files

- `src/plan_html/assets/base.css`
- `tests/fixtures/demo-plan/04-long-lines.md`

## Acceptance Criteria

- [ ] Long inline code wraps inside the card instead of widening it.
- [ ] Long fenced code remains contained inside the card.
- [ ] Normal paragraphs in the same card remain readable.
