# Sub-Plan: Theme Control

## Objective

Add a page-local System / Light / Dark control over token-only themes.

## Design Decisions

- System mode follows `prefers-color-scheme`.
- Light tokens are the base fallback.
- Manual choices do not use browser persistence APIs.

## Primary Files

- `src/plan_html/assets/themes.css`
- `src/plan_html/assets/app.js`

## Acceptance Criteria

- [x] System mode follows the browser color-scheme preference.
- [x] Light and Dark controls set explicit page-local overrides.
- [ ] No `localStorage`, cookies, or other persistence mechanisms are used.
