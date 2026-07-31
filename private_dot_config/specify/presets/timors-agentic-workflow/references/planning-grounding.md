# Planning Grounding Policy

## Purpose

This policy composes repository grounding into `/speckit.plan` while preserving
the upstream 0.12.11 planning workflow. It does not replace the upstream
research, data-model, contracts, quickstart, technical-context, or constitution
work.

## Plan Section Contract

The composed `plan-template` contains these headings exactly once:

1. `## Current-State Grounding`
2. `## Existing Interfaces And Integration Seams`
3. `## Planning Constraints`
4. `## Documentation Impact`

Complete each section with observed, project-specific evidence. Use
project-relative paths in documentation and absolute paths only for filesystem
operations, matching the upstream command contract.

### Current-State Grounding

Record facts discovered in the repository before proposing work: existing
behavior, relevant modules, tests, configuration, generated artifacts, and
known compatibility boundaries. State the evidence for each fact. Do not place
architecture choices or their rationale in this section.

### Existing Interfaces And Integration Seams

Record interfaces that must remain compatible, including public APIs, CLI
inputs/outputs, configuration formats, persistence boundaries, generated files,
callers, consumers, and extension points. Point to `data-model.md` and
`contracts/` for interface and data detail rather than duplicating it here.

### Planning Constraints

Record constraints already imposed by the specification, constitution, existing
interfaces, repository conventions, dependencies, platform support, and
backward compatibility. A constraint is evidence, not a design choice.

### Documentation Impact

List documentation to create, update, or intentionally leave unchanged. Give a
project-relative path and rationale for each entry. Do not turn this section
into a task list.

## Artifact Ownership

Keep artifact responsibilities separate:

- `research.md` contains decisions, rationale, alternatives, and resolution of
  clarifications.
- `data-model.md` and `contracts/` contain data and interface detail.
- `execution-plan.md` contains execution decomposition and coordination.
- `plan.md` contains the upstream implementation plan plus the four grounding
  sections; it does not duplicate the other artifacts.

## Feature Reuse

When the current feature already has `plan.md` and does not have `tasks.md`,
rerun `/speckit.plan` to refresh the plan and its supporting design artifacts.
Do not skip planning merely because `plan.md` exists. A feature with only
`spec.md` follows the normal upstream planning path. In either case, preserve
the upstream setup script behavior and do not duplicate any of the four added
headings.

## Command Boundaries

The command spine stays concise. Detailed execution-plan validation is the
single policy in
`.specify/presets/timors-agentic-workflow/references/artifact-validation.md`.
Compatibility preflight policy is in
`.specify/presets/timors-agentic-workflow/references/protocol-compatibility.md`.
