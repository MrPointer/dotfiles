# Feature Definition Contract

`spec.md` supports bounded product, defect, refactor, infrastructure,
documentation, tooling, and similar work. Its closed core is metadata containing
only status `Draft` or `Ready`; Intent; Scope; semantically titled Requirements
with inline acceptance evidence; Constraints and Dependencies; Consequential
Decision Context; and Material Open Questions. Add actor, journey, entity,
terminology, failure, risk, accessibility, localization, compatibility, or
operational material only when relevant. Downstream artifacts cite section paths
and semantic titles, never serial requirement identifiers.

The complete Ready predicate is: intent and underlying need are coherent; scope
and boundaries are explicit; requirements and evidence state the desired outcome;
applicable constraints, dependencies, actors, failures, and adjacent effects are
addressed; and no unanswered decision could materially change any of those fields.
Repository feasibility and technical consequences are deliberately outside this
predicate.

`specify` does not interact: it marks Ready only when supplied intent positively
establishes every criterion without defaults, inference, or unanswered questions.
Otherwise it records material open questions and marks Draft. `clarify` may begin
from either status. When a rescan of Ready discovers a material question, it writes
that question and changes to Draft before dialogue. After every human answer it
reconciles the entire artifact and uses the same predicate. See
`references/review-lifecycle.md` for ordered invalidation and lock handling.
