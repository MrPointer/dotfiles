---
name: authoring-rfcs
description: Use when turning a settled design direction, brainstorming outcome, context anchor, or architecture discussion into a codebase-grounded RFC. Produces self-contained RFC documents with verified current architecture, chosen design, tradeoffs, risks, stable IDs, and planning handoff context without creating implementation plans.
---

# RFC Authoring

Turn a validated direction into an engineering RFC that can survive review and feed planning without re-litigating the architecture.

The output is an RFC document for human engineering review: self-contained, precise, grounded in the current codebase or explicitly stated assumptions, and clear about the chosen design. It is not an anchor, implementation plan, ADR, transcript, agent handoff, or collection of well-arranged thoughts.

## Relationship To Other Skills

| Skill | Owns | RFC Authoring Role |
|-------|------|--------------------|
| `brainstorming` | Exploration, framing, alternatives, user decisions, convergence | Convert the settled direction into the final RFC |
| `anchoring-context` | Durable working memory for active work | Use the anchor as input, then let the RFC become the settled design artifact |
| `applying-architecture-patterns` | Clean Architecture, Hexagonal Architecture, and Domain-Driven Design guidance | Use as a companion reference when the RFC designs or changes backend architecture patterns |
| `planning-project-features` | Decomposing an approved design into executable plans | Produce and reviewer-validate the RFC that planning uses as its requirements and architecture baseline |

If used after `brainstorming`, do not reopen the conversation unless reviewer review reveals a blocking gap. If used without `brainstorming`, confirm that the chosen direction is already settled before writing the RFC.

## Core Standard

1. **Reality Before Proposal**: Verify the current architecture before describing the proposed one. Read existing docs first, then source files for the affected areas.
2. **Observed vs. Proposed**: Keep current-state facts separate from proposed changes. A reader should always know what exists today and what the design introduces.
3. **One Chosen Design**: The RFC commits to a direction. Alternatives belong in the alternatives section, not sprinkled through the proposed design.
4. **Architectural Precision**: Name actual components, boundaries, contracts, data shapes, state owners, and runtime flows when they affect the design.
5. **No Execution Plan**: Capture enough detail for planning to decompose the work, but do not create task lists, sequencing, tickets, test commands, or step-by-step implementation instructions.
6. **No Fabrication**: If a claim is not verified and not a user-stated requirement, mark it as an assumption or open question. Do not invent architecture to make the document feel complete.
7. **Reviewer-Owned Quality**: Architecture, risk, and clarity review belong to reviewer subagents. The RFC author drafts and revises; reviewers validate.
8. **Pattern-Aware, Not Pattern-Driven**: Use established architecture patterns when they fit the problem and existing system. Do not force Clean Architecture, Hexagonal Architecture, or DDD onto simple CRUD or codebases that do not need that structure.
9. **Human RFC Voice**: Write the RFC as a normal engineering document by and for people. Do not expose the prompting process, agent workflow, reviewer mechanics, or phrases like "the user decided" in design prose.

## Workflow

### Phase 1: Gather Inputs

Read the available design context before drafting:

- User-approved direction and constraints from the current conversation.
- Active anchor, if one exists.
- Existing RFCs, design documents, ADRs, architecture docs, domain docs, and process docs relevant to the affected area.
- Source files that define the current architecture, contracts, schemas, configuration, or runtime flow.

Load `applying-architecture-patterns` before drafting when the RFC involves any of these:

- New backend system architecture.
- Refactoring a monolith or tightly coupled backend boundaries.
- Clean Architecture, Hexagonal Architecture, ports/adapters, DDD, aggregates, repositories, domain events, bounded contexts, or microservice decomposition.
- Testability, mockability, or framework-independence as a design goal.

Use the companion skill as design vocabulary, tradeoff guidance, and a pitfall checklist. Do not copy its content into the RFC and do not apply a pattern unless it is justified by the current system and goals.

If the work is greenfield or the codebase has no relevant implementation yet, state that explicitly in the RFC and separate assumptions from verified facts.

### Phase 2: Verify Current Reality

Build a concise current-state model before writing the proposed design:

- Existing components and their responsibilities.
- Current data/control flow through the affected path.
- Existing contracts, public interfaces, configuration surfaces, storage formats, or external integrations.
- Constraints created by shipped behavior, persisted data, compatibility requirements, operational practices, or team conventions.
- Current pain points that the design intentionally addresses.

Every important "currently" claim in the RFC should be backed by a source reference, documentation reference, or explicit user statement.

### Phase 3: Normalize Decisions

Before drafting, reduce the brainstorming output into settled design inputs:

- Chosen approach and why it was chosen.
- Goals and non-goals.
- Hard constraints and compatibility requirements.
- Rejected alternatives with concrete reasons.
- Open questions that remain genuinely unresolved.

If an unresolved question changes the architecture, stop and ask the user one question. Non-blocking questions can remain in the RFC as open questions.

### Phase 4: Write The RFC

Use the [RFC template][rfc-template] unless the project has a stronger convention.

Write in a direct human engineering voice:

- Convert conversation facts into project facts, constraints, or decisions. For example, write "Homebrew-only `uv` availability is acceptable for the MVP because target setups primarily use Homebrew," not "The user accepted Homebrew-only `uv` availability."
- Use "we" only when it describes the project/team's chosen direction or operational responsibility. Avoid referring to "the user," "the agent," "the model," "the assistant," "the reviewer," or "the conversation" in the RFC body.
- Put conversation-derived information in neutral terms such as "Requirement," "Constraint," "Decision," or "Assumption." Reserve a neutral source label such as "Requirement or decision input" for the **Source References** table when a fact cannot be derived from code or docs.
- Do not write transcript-style rationale. The RFC should explain why the design is right, not who said what during discovery.
- Do not optimize the language for downstream agents. Planning agents should benefit from the same clear human-readable design record that engineers use.

Before review, scan the RFC body and rewrite any agent-centric or transcript-like phrasing, including "the user," "the assistant," "the agent," "the conversation," "explicitly accepted," "asked for," or "decided during brainstorming."

Save the draft RFC using the project's convention. If none exists, use this ID and path convention:

```text
docs/rfcs/RFC-0001-<topic>.md
```

Assign the next sequential four-digit ID by inspecting existing RFCs in `docs/rfcs/`. The first RFC is `RFC-0001`. Use a stable, kebab-case topic slug after the ID, and title the document `# RFC-0001: <Title>`. Preserve the ID across revisions; update status and `Revision: R1` metadata instead of renumbering the RFC.

Scale detail to the design's complexity, but keep these sections conceptually present:

- Problem, goals, non-goals, and constraints.
- Verified current architecture.
- Chosen design and decision summary.
- Proposed architecture, including boundaries, responsibilities, contracts, data/control flow, state, lifecycle, and failure behavior.
- Alternatives considered and why they were rejected.
- Risks, tradeoffs, design success criteria, and planning handoff notes.
- Source references that show what was verified.

Do not leave placeholders. Remove irrelevant optional sections rather than keeping empty headings.

### Phase 5: Reviewer Review

Before presenting the first complete RFC draft, run reviewer subagents once. Do not self-approve the RFC's architecture, risk profile, or clarity.

Required reviewers:

- **`rfc-architect-reviewer`**: Reviews architectural soundness, boundary placement, contracts, data/control flow, fit with current architecture, and whether the RFC contains enough design detail for planning without re-deciding architecture.
- **`rfc-risk-reviewer`**: Reviews technical risks, migration and compatibility concerns, rollback gaps, operational hazards, hidden complexity, and risk mitigations.

Optional reviewer:

- **`rfc-clarity-reviewer`**: Use when the RFC is long, has nuanced decisions, contains assumptions or open questions, was produced from a long brainstorming session, or the user asks for extra clarity review.

When launching reviewers, pass the RFC path, active anchor path if any, relevant source references, and the intended review output path. Ask reviewers to review the RFC as a design artifact, not as an implementation plan. They should not request task sequencing, code-level implementation instructions, or plan decomposition details.

Review output files are cumulative artifacts. When re-running a reviewer, reuse the same review output path and instruct the reviewer to preserve existing content by appending a new review round to the file. Do not let a later review overwrite earlier findings. If a reviewer returns findings instead of writing the file, append those findings to the existing review file yourself without replacing prior rounds.

Use this default review output location unless the project has a stronger convention:

```text
docs/rfcs/reviews/RFC-0001.<reviewer>.md
```

Incorporate reviewer findings into the RFC before presenting it. Preserve explicit user decisions: if reviewer feedback conflicts with a user decision or verified codebase reality, record the concern and ask the user before changing the design. Non-blocking concerns may remain in **Risks And Tradeoffs** or **Open Questions** with rationale.

#### Review Change Classification

After incorporating initial reviewer findings, classify the resulting RFC changes before deciding whether to run or recommend re-review. Do not restart the full review by default.

| Change Type | Examples | Re-review Action |
|---|---|---|
| Editorial / wording | Improve phrasing, remove transcript-like language, fix typos, reorganize prose without changing meaning | **None** |
| Evidence / citation repair | Add source references, quote verified current-state details, clarify where a constraint came from | **None** unless the added evidence contradicts the reviewed design |
| Finding-specific repair | Address one reviewer's finding without changing the chosen design, boundaries, contracts, data flow, compatibility strategy, or risk posture | **None by default**; recommend targeted re-review only if confirmation is genuinely needed |
| Cross-scope design change | Change chosen approach, component boundaries, public contracts, data/control flow, state ownership, compatibility or migration strategy, failure behavior, or major risk mitigation | Run targeted re-review of the affected scopes, commonly `rfc-architect-reviewer` and/or `rfc-risk-reviewer` |
| Large accumulated revision | Many smaller edits together make the reviewed artifact materially different, even if each edit looked local | **Ask the user** whether to spend tokens on targeted re-review; do not decide silently |

Required reviewers do not need to re-approve every minor edit. There is no automatic convergence loop after the initial review. Automatically re-run reviewers only for clear cross-scope design changes, and only for affected scopes. If it is ambiguous whether a change is cross-scope, state the ambiguity and ask the user before launching any reviewer. Large accumulated revisions always require asking before re-review.

Update the RFC's **Review Record** before presenting it so planning can see which reviews ran, what remains open, and whether architecture and risk review passed. Use statuses like `Passed`, `Passed with concerns`, `Blocking`, or `Not requested`. In the notes, record whether a verdict came from the original review, a targeted re-review, or author classification that no re-review was needed after a limited edit.

If the user requests changes after the reviewed RFC is presented, incorporate them and classify the change using the same table. Default to stating what changed and recommending whether re-review is warranted; automatically spend reviewer tokens only for clear cross-scope design changes or when the user explicitly asks for re-review. The user can always request re-review regardless of classification.

### Phase 6: Present And Stop

Present the RFC path and a brief summary of the chosen architecture. Stop there unless the user explicitly asks for revisions or asks to move into planning.

If an anchor was used, update or reconcile it so settled design context points to the RFC instead of duplicating it.

## Rules

- Never substitute a context anchor for an RFC.
- Never let the RFC become a plan; planning owns decomposition and execution details.
- Never self-approve the RFC's architecture, risk profile, or clarity when reviewer agents are available.
- Pseudo-code is allowed when it clarifies an algorithm, contract, state transition, or data/control flow. Keep it illustrative and language-agnostic unless the RFC is explicitly about a language-specific contract.
- Never write practical implementation code or low-level implementation recipes. Describe contracts, responsibilities, state, and behavior at the architectural level.
- Never hide uncertainty in vague language. Ask about blocking uncertainty or record non-blocking uncertainty as an open question.
- Never claim an existing component, interface, storage format, or runtime behavior exists without verifying it.
- Prefer specific file paths and document references over broad statements like "the existing architecture".
- Preserve approved decisions. If a verified reality conflicts with an approved decision, present the conflict and ask before changing the design.
- Never refer to the human requestor as "the user" in RFC prose. Recast their input as requirements, constraints, decisions, assumptions, or source references.
- Never mention agents, models, prompts, subagents, reviewer mechanics, or the authoring workflow in the RFC body except for the required **Review Record** metadata.

[rfc-template]: assets/rfc-template.md
