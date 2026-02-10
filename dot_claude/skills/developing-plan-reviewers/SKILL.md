---
name: developing-plan-reviewers
description: Create project-specific plan reviewer agents that integrate with the planning-project-features workflow. Use when a project needs domain-specialized reviewers for sub-plans (e.g., Go code reviewer, API layer reviewer, database reviewer) or when the planner warns that no suitable local reviewer exists for a sub-plan's domain.
---

# Developing Plan Reviewers

Create project-specific reviewer agents for the `planning-project-features` workflow. These agents review individual sub-plans for domain-specific correctness — catching issues that the generic global reviewers (architecture, risk) cannot.

For general agent development guidance (frontmatter fields, description examples, system prompt design), see the `developing-agents` skill.

## How Plan Reviewers Fit In

The planning workflow (Phase 5) uses two types of reviewers:

1. **Global reviewers** (`~/.claude/agents/`) — `plan-architect-reviewer` and `plan-risk-reviewer`. Project-agnostic, already exist.
2. **Local reviewers** (`.claude/agents/`) — project-specific, domain-specialized. **This is what you build.**

The planner discovers local reviewers in Phase 4 by reading their descriptions, matches each sub-plan to the most appropriate reviewer, and launches them in Phase 5. No naming convention is assumed — the description is what matters.

## Reviewer Agent Template

```markdown
---
name: plan-<domain>-reviewer
description: "Use this agent to review sub-plans that involve <domain>. Evaluates
  <what it checks> against <what standards>.

<example>
Context: A sub-plan covers <domain> implementation within a feature plan.
user: \"Review sub-plan 02-<task>.md for <domain> correctness.\"
assistant: \"I'll review the sub-plan for <domain> issues using the plan-<domain>-reviewer agent.\"
<commentary>
Sub-plan involves <domain> work. Launch the domain-specific reviewer.
</commentary>
</example>"
tools: Read, Write, Glob, Grep
skills:
  - <skill-1>
  - <skill-2>
---

You are a <domain> reviewer. Your job is to review implementation sub-plans
for <domain> correctness — ensuring the proposed approach follows established
conventions, avoids known pitfalls, and will produce correct results.

You are NOT here to praise, summarize, or restate the plan. You are here to
find what's wrong with it from a <domain> perspective.

## What You Review

You will be given a path to a specific sub-plan file (e.g.,
`.claude/plans/<feature>/02-<task>.md`). You also have access to the full
codebase to verify claims.

## How You Review

### 1. Read the Sub-Plan, Documentation, and Relevant Code

Read the sub-plan completely. Then check existing project documentation
(CLAUDE.md, AGENTS.md, architecture docs, component docs) for relevant
conventions and context. Documentation is dramatically cheaper than code
exploration. Only after reviewing docs, use Glob and Grep to examine the
files the sub-plan references and fill gaps that documentation doesn't cover.

### 2. Evaluate <Domain Area 1>

- <What to check>
- <What patterns to verify>
- <What pitfalls to look for>

### 3. Evaluate <Domain Area 2>

- <What to check>
- ...

[Add as many evaluation areas as the domain requires]

## Output Format

Write your findings to `reviews/<plan-file>.<reviewer-type>.md` inside the
plan directory. Use the exact format below.

[Insert the standard output format template from this skill]

## Rules

- **Be specific and actionable** — every finding must reference the exact
  plan section and provide a concrete recommendation.
- **Review the plan, not the code** — you evaluate whether the plan's
  strategy is sound for this domain. Code-level review happens during
  execution.
- **Don't invent requirements** — review against the sub-plan's stated
  objective and acceptance criteria.
- **Don't duplicate architecture or risk review** — focus only on your
  domain expertise.
```

### Key Choices

- **`tools: Read, Write, Glob, Grep`** — Write is only for the review output file. Match the global reviewer pattern.
- **`skills`** — the differentiator. Preloads domain knowledge so the reviewer doesn't need to discover conventions at runtime. See [Choosing Skills to Preload](#choosing-skills-to-preload).
- **No `model` field** — inherits from parent, matching the global reviewers.
- **Description** — must be clear enough for the planner to match it to sub-plans. Include what domain it covers and what it evaluates.

## Choosing Skills to Preload

The `skills` frontmatter field injects full skill content into the reviewer's context at startup. This is how domain knowledge gets into the reviewer.

**Selection criteria:**
- List only skills whose conventions the reviewer needs to evaluate sub-plans correctly
- Check both global (`~/.claude/skills/`) and local (`.claude/skills/`) skills
- Don't overload — each injected skill consumes context budget

**Examples:**
| Reviewer Domain | Skills to Preload |
|-----------------|-------------------|
| Go code | `writing-go-code`, `applying-effective-go` |
| API layer | `writing-go-code` (or API-specific skill if one exists) |
| Shell config | `configuring-zsh` |
| CI/CD | `configuring-github-actions` |
| Dotfile management | `managing-chezmoi` |

## Standard Output Format

All plan reviewers **must** follow this format. The planner depends on the structure for output normalization (Phase 5, Step 2).

```markdown
# <Domain> Review: <Sub-Plan Name>

## Verdict

<PASS | PASS WITH CONCERNS | NEEDS REVISION>

## Critical Findings
<Issues that MUST be fixed before the plan can proceed. Empty if none.>

### Finding: <short title>
- **Affects**: <plan file and section>
- **Problem**: <what's wrong from a domain perspective>
- **Recommendation**: <how to fix it>

## Concerns
<Issues that SHOULD be addressed but aren't blockers. Empty if none.>

### Concern: <short title>
- **Affects**: <plan file and section>
- **Problem**: <what's wrong>
- **Recommendation**: <how to fix it>

## Observations
<Minor notes, suggestions, or things the planner might want to consider. Empty if none.>
```

## Examples

Complete, copy-and-adapt reviewer examples:

- **[Go Codebase Reviewer](examples/plan-go-reviewer.md)** — reviews sub-plans for idiomatic Go patterns, error handling, interface design, and test strategy. Preloads `writing-go-code` and `applying-effective-go`.
- **[API Layer Reviewer](examples/plan-api-reviewer.md)** — reviews sub-plans for endpoint design, request/response contracts, error responses, and backward compatibility.

## Rules

- **Always follow the standard output format** — the planner depends on the Verdict/Critical Findings/Concerns/Observations structure
- **Write output to `reviews/<plan-file>.<reviewer-type>.md`** — never elsewhere
- **Review the plan, not the code** — evaluate whether the plan's strategy is sound for the domain; code-level review happens during execution
- **Be specific and actionable** — every finding must reference the exact plan section and provide a recommendation
- **Don't duplicate architecture or risk review** — focus only on domain expertise
- **Don't invent requirements** — review against the sub-plan's stated objective and acceptance criteria
