# Documentation Planning Policy

## Trigger

Plan a dedicated documentation execution group when implementation changes
domain terminology or rules, architecture or component boundaries, or an
end-to-end business or operational process. Also create one when `plan.md`'s
Documentation Impact names substantive documentation work that cannot remain
with one story.

Keep a small, story-specific documentation change in that story's group when it
shares the same ownership and verification seam. Do not create documentation
tasks merely to restate code or preserve planning conversation.

## Placement And Ownership

Place a dedicated documentation group after the implementation groups whose
integrated behavior it documents. Give it logical predecessors for concrete
implementation facts it consumes. It owns exact documentation paths and has
its own applicable documentation skill IDs, test expectation, behavioral
acceptance criteria, and verification.

Document only current, reader-relevant behavior. Include context that affects
the required outcome, a plausible implementation choice, an interface,
verification, or a concrete risk. Omit conversation residue, rejected
alternatives already settled elsewhere, and unrelated non-goals. State a
negative constraint only when naming the plausible competing behavior or risk
that makes the constraint relevant.

## Documentation Data Flow

When a documentation group depends on implementation results, record explicit
flows and contracts for the shapes or artifacts it consumes. Do not use a vague
"implementation complete" flow. Name the interface, behavior, schema, or
verified artifact that becomes documentation input, and keep producer and
consumer shape strings identical.

## Post-Execution

Post-Execution states how integrated behavior and documentation are checked for
agreement. Include an applicable project documentation reviewer from
`.specify/reviewers/` when its triggers match; ask the user if applicability is
ambiguous.
