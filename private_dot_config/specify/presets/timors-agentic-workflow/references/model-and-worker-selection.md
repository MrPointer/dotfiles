# Model And Worker Selection Policy

## Semantic Model Tiers

Assign every execution group exactly one provider-neutral tier: `Cheapest`,
`Mid-tier`, or `Most capable`. Select conservatively from the work's reasoning
needs, ambiguity, integration risk, and required verification:

- `Cheapest` for bounded, mechanical work with explicit local patterns and
  straightforward verification;
- `Mid-tier` for normal implementation requiring repository reasoning,
  contract adherence, or coordinated tests; and
- `Most capable` for unusually ambiguous, high-risk, cross-boundary, or complex
  work where a weaker tier materially threatens correctness.

Do not choose a tier from provider names, availability guesses, or cost alone.
The Execution Groups table owns the tier. Each detail's Execution Model records
substantive rationale only and must not repeat a tier token as its value.

## Runtime Boundary

Every execution group is a delegated unit, but this preset declares data only.
It does not create workers, bind provider models, dispatch agents, or configure
runtime execution. Worker selection and binding remain owned by the active
runtime at implementation time. Required Capabilities declare needs; they do
not grant capabilities.

## Skills And Capabilities

List the exact IDs of all skills needed to execute the group. The list is never
empty and contains no descriptions or aliases. List only required capabilities
from `read`, `edit`, `shell`, `network`, and `subagent-dispatch`. Include `read`
for every group and `edit` for every group that modifies files. Do not request
network or subagent dispatch without a concrete need.

## Project Review Roles

Inspect `.specify/reviewers/` for applicable project reviewer packets. Use
`None` only when no role applies. If change triggers or project context leave
applicability ambiguous, ask the user rather than silently selecting or
omitting a role.

Every selected role has a lowercase-kebab-case ID and path exactly
`.specify/reviewers/<role-id>.md`. Its plan row records the packet's semantic
Model Tier and a concrete applicability rationale. The packet must match this
schema:

```markdown
# Reviewer Role: <name>

- **Role ID**: <lowercase-kebab-case>
- **Model Tier**: Cheapest | Mid-tier | Most capable
- **Change Triggers**: <comma-separated subset>

## Scope
## Required Inputs
## Exclusions
## Review Checks
## Output Contract
```

Change Triggers contains only `baseline`, `tasks`, `decomposition-design`,
`documentation`, `reviewer-config`, or `all`. The packet's Role ID and Model
Tier must exactly match its execution-plan row.
