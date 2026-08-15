# Independent Review Lifecycle

Review roles are independent and read-only. Bind only a candidate with requested
tier, exact skills, native discoverability, invokability, correct workspace
targeting, and attributable results; project-local candidates are preferred.
Select a unique preferred candidate, a sole eligible candidate, or an explicit
human choice. Abstract capabilities and permission breadth are not eligibility
inputs. The first invocation is the real review assignment; there is no probe.

Every role report is cumulative and appends a complete current-state round with
trigger, scope, reviewed paths, binding/invocation/workspace/start/result evidence,
verdict, and semantic blocking or concern findings. A targeted rerun carries
unaffected current findings into its new complete snapshot. Findings have semantic
titles, not opaque identifiers. The review pointer holds only its declared report,
round, revision when applicable, status, and verdict; detailed evidence remains in
the report.

For an RFC review that may have started without attributable complete result, set
its pointer to unchanged report path, round `None`, revision `None`, status
`Recovery required`, verdict `Pending`. For a task review use unchanged report,
round `null`, status `recovery-required`, verdict `pending`. Do not retain an old
pointer value there. Recover with attributable same-session resume or a newly
selected binding only after confirming reviewed artifact unchanged, workspace
read-only and correctly targeted, and retained evidence reconciled. A recovered
result appends the next complete round and updates the pointer normally. Ambiguous
start, mutation, target, or result stops the producing phase: never self-review,
infer a verdict, request a downstream human gate, or abandon implementation work.

A producing phase is quiescent only when no pending role can safely dispatch,
resume, or complete without a collected binding or recovery blocker. Then present
one interaction with all current blockers. After a response resume work; a later
quiescent state may present one later batch. Automatic redispatch after known start
is forbidden. Material artifact changes follow their owner’s invalidation and
fresh-review rules; a later invocation may resume recovery only against unchanged
current artifact.
