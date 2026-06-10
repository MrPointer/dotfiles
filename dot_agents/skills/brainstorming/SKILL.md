---
name: brainstorming
description: Use when exploring a problem space before implementation or planning. Turns rough ideas into validated, codebase-grounded RFCs through collaborative dialogue — challenges assumptions, explores alternatives, and converges on a committed design direction before final RFC authoring.
---

# Brainstorming

Explore a problem space collaboratively. Turn rough ideas into validated RFCs through genuine dialogue — not a checkbox exercise.

The output is an **RFC document**: a committed artifact that captures the problem understanding, chosen approach, key decisions, and verified architectural shape. It stands on its own — useful whether the next step is formal planning, direct implementation, or team review.

For brainstorming that spans sessions or substantial decision context, pair this skill with `anchoring-context`. Brainstorming drives exploration and convergence; the anchor preserves concise working memory until the final RFC is produced. When an anchor is active, update it as decisions, rejected alternatives, constraints, and resolved questions settle. Do not wait until RFC creation or handoff to reconstruct the session from memory.

For final RFC authoring, pair this skill with `authoring-rfcs`. Brainstorming owns the conversation and settled direction; `authoring-rfcs` owns RFC document production and reviewer coordination.

## Core Principles

1. **Diverge Before Converging**: Explore the problem space before narrowing to a solution. Challenge the framing, not just the details.
2. **One Thread at a Time**: Ask one question per message. Give each answer room to reshape your understanding before moving on.
3. **Genuine Alternatives**: When proposing approaches, each option must be a real contender with real tradeoffs — not strawmen set up to make one option look obvious.
4. **The User Decides**: Present options with your recommendation and reasoning. Never silently pick an approach.
5. **No Implementation**: This skill produces an RFC. It does not write code or create plans. It may pair with `anchoring-context` for long-running working memory and `authoring-rfcs` for final document production, but does not chain into implementation or planning skills.

## Workflow

### Phase 1: Orientation

Before asking anything, understand where you are:

1. **Read project context**: Check the project's documentation pointers (AGENTS.md, CLAUDE.md, etc.) and read relevant docs (domain, architecture, business processes). Skim recent commits if the request relates to ongoing work.
2. **Read active anchors**: If the request resumes or continues existing work, look for a feature anchor using project conventions, then `docs/context/<topic>-anchor.md`. Read it before asking the user to restate context. If the brainstorming is likely to span sessions and no anchor exists, create one with `anchoring-context` before context accumulates only in chat. Once an anchor is active, treat it as write-through memory for the session, not as an end-of-session summary destination.
3. **Assess scope**: If the request describes multiple independent subsystems, flag this immediately. Don't spend questions refining details of something that needs decomposition first.
4. **Identify what you already know vs. what you need to learn**: Don't ask questions the codebase already answers.

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

If an anchor is active, record the chosen approach, the reason it was chosen, alternatives the user rejected, and any constraint or preference revealed by the choice before moving into detailed design.

### Phase 4: Design

Once the approach is agreed, present the design in sections scaled to complexity:

- **Architecture**: How the pieces fit together
- **Components**: What gets built, what gets modified
- **Data flow**: How information moves through the system
- **Error handling**: What can go wrong and how it's handled
- **Constraints & boundaries**: What's in scope, what's explicitly out

Present each section and confirm before moving on. A section that's straightforward gets a few sentences. A section that's nuanced gets as much space as it needs.

If an anchor is active, update it after each confirmed design section or meaningful correction. Capture what changed and why before presenting the next major section.

**In existing codebases**: Explore current structure before proposing changes. Follow established patterns. Where existing code has problems that affect the work, include targeted improvements as part of the design — don't propose unrelated refactoring.

### Phase 5: RFC Creation

Write the validated design to an RFC document with `authoring-rfcs`:

1. **Load the RFC authoring skill**: Use `authoring-rfcs` for the final RFC document structure and reviewer workflow.
2. **Pass the settled design inputs**: Include the agreed approach, constraints, rejected alternatives, unresolved non-blocking questions, and any anchor context.
3. **Verify codebase reality**: Before writing the RFC, read the relevant docs and source files needed to describe the current architecture accurately. Do not let the anchor substitute for current-state verification.
4. **Save the RFC**: follow project conventions, or use `docs/rfcs/<topic>.md` when no convention exists. Use numbered RFC IDs only when the project already uses them or the user explicitly requests them.
5. **Reconcile any anchor**: If an anchor was used, ensure the RFC reflects the settled decisions and that the anchor keeps only active open questions, handoff state, or implementation-continuity context. This is a hygiene pass, not the first time decisions should be captured. For long design sessions, reconcile the anchor into a self-contained design snapshot: checkbox current state, thematic decisions, active open questions only, and a quick summary that can be understood without reading the full decision log.
6. **Run RFC reviewer review**: Use the reviewer workflow defined by `authoring-rfcs`. Incorporate reviewer findings before presenting.
7. **Present the RFC path** and ask the user to review. Wait for approval.

If the user requests changes, make them and re-run any affected RFC reviewers as defined by `authoring-rfcs`. Only finalize once the user approves.

### Phase 6: Handoff

Present the completed RFC with a brief summary:

- What problem this solves
- Which approach was chosen and why
- Key decisions made during brainstorming
- The RFC file path

**Stop here.** The brainstorming skill's job is done. The user decides what happens next — planning, implementation, team review, or shelving.

## Rules

- **One question at a time** — don't overwhelm with lists of questions
- **Never skip Phase 2** — even if the request seems clear, explore before converging
- **Never chain to implementation or planning skills** — this skill terminates with an RFC artifact. `anchoring-context` is allowed as a companion memory protocol for long-running brainstorming. `authoring-rfcs` is allowed for final RFC production only.
- **Never write code** — not even pseudo-code. Describe behavior and architecture, not implementation.
- **Never invent architecture** — verify current-state claims against code/docs or mark them as assumptions/open questions.
- **Respect the user's domain knowledge** — they may know things you can't derive from the codebase. When they state something confidently, don't second-guess it with "are you sure?"
- **If scope is too large, decompose first** — help the user break the problem into independent pieces before brainstorming any single piece in depth
- **Keep active anchors current** — when paired with `anchoring-context`, write settled decisions, rejected alternatives, constraints, and answered questions to the anchor during the conversation. End-of-session reconciliation may clean up the anchor, but it must not be the only capture point.
