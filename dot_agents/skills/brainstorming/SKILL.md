---
name: brainstorming
description: Use when exploring a problem space before implementation or planning. Turns rough ideas into validated design specs through collaborative dialogue — challenges assumptions, explores alternatives, and produces a committed design artifact. Use before planning features, starting new work, or when requirements feel unclear.
---

# Brainstorming

Explore a problem space collaboratively. Turn rough ideas into validated design specs through genuine dialogue — not a checkbox exercise.

The output is a **design spec document**: a committed artifact that captures the problem understanding, chosen approach, and key decisions. It stands on its own — useful whether the next step is formal planning, direct implementation, or team review.

## Core Principles

1. **Diverge Before Converging**: Explore the problem space before narrowing to a solution. Challenge the framing, not just the details.
2. **One Thread at a Time**: Ask one question per message. Give each answer room to reshape your understanding before moving on.
3. **Genuine Alternatives**: When proposing approaches, each option must be a real contender with real tradeoffs — not strawmen set up to make one option look obvious.
4. **The User Decides**: Present options with your recommendation and reasoning. Never silently pick an approach.
5. **No Implementation**: This skill produces a design spec. It does not write code, create plans, or invoke other skills.

## Workflow

### Phase 1: Orientation

Before asking anything, understand where you are:

1. **Read project context**: Check the project's documentation pointers (AGENTS.md, CLAUDE.md, etc.) and read relevant docs (domain, architecture, business processes). Skim recent commits if the request relates to ongoing work.
2. **Assess scope**: If the request describes multiple independent subsystems, flag this immediately. Don't spend questions refining details of something that needs decomposition first.
3. **Identify what you already know vs. what you need to learn**: Don't ask questions the codebase already answers.

### Phase 2: Divergent Exploration

This is the critical phase that separates brainstorming from requirements gathering. The goal is to **expand the problem space** before narrowing it.

**Challenge the framing:**
- Is the stated problem the real problem, or a symptom?
- Are there upstream causes that, if addressed, would make this request unnecessary?
- What assumptions is the request built on? Are they valid?

**Map the problem space:**
- Who and what does this affect beyond the obvious?
- What adjacent systems, workflows, or user experiences does this touch?
- What are the constraints that aren't stated but exist?

**Explore the edges:**
- What happens if we do nothing?
- What's the simplest version that would still be valuable?
- What's the most ambitious version, and what would it unlock?

**Ask questions one at a time.** Prefer multiple choice when the options are clear, open-ended when the space is genuinely open. Each answer should inform your next question — don't follow a script.

**When to stop exploring**: When you can articulate the problem clearly, know the constraints, and have a mental model of the solution space. If you're asking questions whose answers won't change your understanding, move on.

### Phase 3: Approaches

Propose 2-3 different approaches with genuine tradeoffs:

1. **Lead with your recommendation** and explain why.
2. **Present alternatives honestly** — each should have a scenario where it's the better choice.
3. **Name the tradeoffs concretely** — not "more complex" but "adds a migration step and requires updating the CI pipeline."
4. **Include the constraint landscape** — what would make you change your recommendation?

Let the user choose or refine. If they push back, understand why before adjusting — their objection may reveal a constraint or preference you missed.

### Phase 4: Design

Once the approach is agreed, present the design in sections scaled to complexity:

- **Architecture**: How the pieces fit together
- **Components**: What gets built, what gets modified
- **Data flow**: How information moves through the system
- **Error handling**: What can go wrong and how it's handled
- **Constraints & boundaries**: What's in scope, what's explicitly out

Present each section and confirm before moving on. A section that's straightforward gets a few sentences. A section that's nuanced gets as much space as it needs.

**In existing codebases**: Explore current structure before proposing changes. Follow established patterns. Where existing code has problems that affect the work, include targeted improvements as part of the design — don't propose unrelated refactoring.

### Phase 5: Spec Creation

Write the validated design to a spec document:

1. **Save the spec**: `docs/specs/<topic>-design.md` (or follow project conventions if they exist for spec/design documents).
2. **Self-review** before presenting:
   - **Placeholder scan**: Any "TBD", "TODO", incomplete sections?
   - **Internal consistency**: Do sections contradict each other?
   - **Scope check**: Is this focused enough, or does it need decomposition?
   - **Ambiguity check**: Could any decision be interpreted two ways? If so, pick one and make it explicit.
   - Fix issues inline — don't flag them, just fix them.
3. **Present the spec path** and ask the user to review. Wait for approval.

If the user requests changes, make them and re-run the self-review. Only finalize once the user approves.

### Phase 6: Handoff

Present the completed spec with a brief summary:

- What problem this solves
- Which approach was chosen and why
- Key decisions made during brainstorming
- The spec file path

**Stop here.** The brainstorming skill's job is done. The user decides what happens next — planning, implementation, team review, or shelving.

## Spec Document Structure

A reference template is in this skill's `assets/` directory. Read it when creating specs:

- **[Design spec template][spec-template]**

The template is a starting point — adapt sections to fit the design's complexity. Simple designs don't need every section; complex designs may need additional ones.

## Rules

- **One question at a time** — don't overwhelm with lists of questions
- **Never skip Phase 2** — even if the request seems clear, explore before converging
- **Never chain to other skills** — this skill terminates with a spec artifact, nothing more
- **Never write code** — not even pseudo-code. Describe behavior and architecture, not implementation.
- **Respect the user's domain knowledge** — they may know things you can't derive from the codebase. When they state something confidently, don't second-guess it with "are you sure?"
- **If scope is too large, decompose first** — help the user break the problem into independent pieces before brainstorming any single piece in depth

[spec-template]: assets/spec-template.md
