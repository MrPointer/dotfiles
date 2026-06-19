# Reviewer Review

Use this reference when running RFC reviewers, handling review artifacts, incorporating findings, classifying changes, or updating the Review Record.

## Reviewers

Before presenting the first complete RFC draft, run reviewer subagents once.

Required reviewers:

- `design-reviewer`: architectural soundness, boundary placement, contracts, data/control flow, current-architecture fit, planning readiness, technical risks, migration, compatibility, rollback gaps, operational hazards, hidden complexity, and risk mitigations

Optional reviewer:

- `rfc-clarity-reviewer`: use when the RFC is long, nuanced, assumption-heavy, produced from a long brainstorming session, or requested by the user

Ask reviewers to review the RFC as a design artifact, not an implementation plan. They should not request task sequencing, code-level implementation instructions, or decomposition details.

## Reviewer Inputs And Outputs

When launching reviewers, pass:

- RFC path
- active anchor path, if any
- relevant source references
- intended review output path

Use this default review output location unless the project has a stronger convention:

```text
docs/rfcs/reviews/<rfc-file-stem>.<reviewer>.md
```

Review output files are cumulative artifacts. When re-running a reviewer, reuse the same output path and instruct the reviewer to append a new review round. Do not let later review overwrite earlier findings. If a reviewer returns findings instead of writing the file, append those findings yourself without replacing prior rounds.

## Incorporating Findings

Incorporate reviewer findings before presenting the RFC.

Preserve explicit user decisions. If reviewer feedback conflicts with a user decision or verified codebase reality, record the concern and ask before changing the design. Non-blocking concerns may remain in `Risks And Tradeoffs` or `Open Questions` with rationale.

## Review Change Classification

After incorporating initial reviewer findings, classify RFC changes before deciding whether to run or recommend re-review. Do not restart the full review by default.

| Change Type | Examples | Re-review Action |
|-------------|----------|------------------|
| Editorial / wording | Improve phrasing, remove transcript-like language, fix typos, reorganize prose without meaning changes | None |
| Evidence / citation repair | Add source references, quote current-state details, clarify where a constraint came from | None unless evidence contradicts reviewed design |
| Finding-specific repair | Address one reviewer's finding without changing design, boundaries, contracts, flow, compatibility, or risk posture | None by default; recommend targeted re-review only when confirmation is genuinely needed |
| Cross-scope design change | Change chosen approach, component boundaries, public contracts, data/control flow, state ownership, compatibility, migration, failure behavior, or major risk mitigation | Run targeted re-review with the design reviewer |
| Large accumulated revision | Many smaller edits make the reviewed artifact materially different | Ask the user whether to spend tokens on targeted re-review |

Required reviewers do not need to re-approve every minor edit. Automatically re-run reviewers only for clear cross-scope design changes, and only affected scopes. If it is ambiguous whether a change is cross-scope, state the ambiguity and ask before launching a reviewer.

If the user requests changes after the reviewed RFC is presented, incorporate them and classify the change with the same table. Default to stating what changed and recommending whether re-review is warranted. The user can always request re-review regardless of classification.

## Review Record

Update the RFC's `Review Record` before presenting so planning can see which reviews ran, what remains open, and whether design review passed.

Use statuses such as `Passed`, `Passed with concerns`, `Blocking`, or `Not requested`. In notes, record whether a verdict came from original review, targeted re-review, or author classification that no re-review was needed after a limited edit.
