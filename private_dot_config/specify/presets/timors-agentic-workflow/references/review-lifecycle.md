# Independent Review Lifecycle

Review roles are independent. A reviewer does not modify the artifacts it
reviews; its only output is its review report, written outside the reviewed set.
Each role names a tier and a role description. Bind a candidate that meets the
requested tier and is natively
discoverable, invokable, correctly targeted so that it does not modify the
reviewed artifacts, and able to
attribute results; project-local candidates are preferred. Among these safe
candidates, match the role description against the candidate descriptions and bind
only a single clear fit. Zero clear fits or several plausible fits produce a
collected blocker for one human choice. Record the chosen binding in runtime
evidence; no tracked artifact names a candidate. Declared skills, abstract
capabilities, permission breadth, and worker names are not eligibility inputs. The
first invocation is the real review assignment; there is no probe.

Every role report is cumulative and appends a complete current-state round with
trigger, scope, reviewed paths, binding/invocation/workspace/start/result evidence,
verdict, and semantic blocking or concern findings. A targeted rerun carries
unaffected current findings into its new complete snapshot. Findings have semantic
titles, not opaque identifiers. The review pointer holds only its declared report,
round, revision when applicable, status, and verdict; detailed evidence remains in
the report.

For a plan review that may have started without attributable complete result, set
its pointer to unchanged report path, round `None`, revision `None`, status
`Recovery required`, verdict `Pending`. For a task review use unchanged report,
round `null`, status `recovery-required`, verdict `pending`. Do not retain an old
pointer value there. Recover with attributable same-session resume or a newly
selected binding only after confirming reviewed artifact unchanged, workspace
correctly targeted, and retained evidence reconciled. A recovered
result appends the next complete round and updates the pointer normally. Ambiguous
start, mutation, target, or result stops the producing phase: never self-review,
infer a verdict, request a downstream human gate, or abandon implementation work.

The producing phase never applies a finding on its own judgment. A fix limited to
wording, formatting, or source citation that leaves every decision, scope,
requirement, constraint, boundary, contract, flow, state, failure behavior, risk
posture, and acceptance evidence unchanged may be applied directly and is listed in
the phase completion report. Every other fix, including one that fills a gap no
human decided, is a proposed remediation: collect it with its finding, the
reviewer recommendation, and the affected artifact sections, and apply it only
after a human accepts it. Report each declined proposal with its human disposition
at phase completion.

A producing phase is quiescent only when no pending role can safely dispatch,
resume, or complete without a collected binding or recovery blocker. Then present
one interaction with all current blockers and proposed remediations. After a
response resume work; a later
quiescent state may present one later batch. Automatic redispatch after known start
is forbidden. Material artifact changes follow their owner’s invalidation and
fresh-review rules; a later invocation may resume recovery only against unchanged
current artifact.
