# Documentation Planning and Final Review

When accepted design changes documented domain concepts, architecture, or
business/operational processes and relevant project documentation exists, task
planning creates a final most-capable documentation packet. It states exact paths,
stale sections, proportional scope, and documentation skills. Skip it only when no
conceptual documentation exists, no documented concept or flow changes, or only
component-level implementation detail may drift.

Component documentation is deferred. Execution Handoff records `required` with
relevant paths when existing component documentation describes a modified
component, interface, or behavior; otherwise it records `not-applicable` with a
substantive reason. At completion, the coordinator checks both the handoff and
current repository documentation. A not-applicable handoff contradicted by current
documentation is a planning defect, not a skip.

After final code exists, bind the independent read-only most-capable
`reviewers/component-docs.md` role. Its result gates completion but creates no
human approval gate. If it finds incomplete work within accepted requirements,
design, task scope, and acceptance, coordinator or bounded delegate may repair it,
append a policy-compliant corrective commit, and rerun the review. New scope,
undefined behavior, contract change, or acceptance exception returns to Feature
Definition, RFC, or task planning ownership.
