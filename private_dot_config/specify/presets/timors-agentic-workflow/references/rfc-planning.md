# RFC Planning and Acceptance

`plan.md` is the sole normative RFC. It contains verified current state and
constraints; semantically titled decisions aligned locally to Feature Definition;
architecture; applicable contracts and flows; alternatives; risks; verification;
planning handoff; and source references. Its `Revision` is `R0`, `R1`, and onward.

Design Acceptance has exactly these fields:

| Field | Contract |
|---|---|
| Status | `Pending`, `Accepted`, `Rejected`, or `Stale` |
| Applies To Revision | current `RNN` for Accepted/Rejected; prior decision revision when Stale; `None` when Pending |
| Recorded By | `Human` for Accepted/Rejected; prior value when Stale; `Pending` when Pending |
| Recorded At | ISO-8601 timestamp for Accepted/Rejected; prior value when Stale; `Pending` when Pending |
| Rationale | substantive context for Accepted/Rejected; prior value when Stale; `Pending` when Pending |

Only Accepted whose Applies To Revision equals current Revision permits task
planning. A material RFC change increments Revision and, before writing, makes
Accepted or Rejected acceptance Stale; Pending stays Pending. Editorial changes
retain Revision and acceptance. Refuse a material write after the progress lock.

The Review Record has fixed `rfc_design` and `rfc_clarity` records, each with
report path, round (`RNN` or `None`), RFC revision (`RNN` or `None`), Status
(`Pending`, `Reviewed`, `Retained`, `Recovery required`), and Verdict (`Pending`,
`Passed`, `Passed with concerns`, `Blocking`). A role is current only if Reviewed
or Retained, its round exists in its report, its revision equals current Revision,
and its verdict passes or passes with concerns. See
`references/review-lifecycle.md` for dispatch and recovery.

First complete RFC review runs both roles. Editorial, evidence repair, or a
finding-specific repair needs no rerun only when it leaves architecture, boundary,
contract, flow, state, compatibility, migration, failure behavior, risk posture,
and cold-reader meaning unchanged. Otherwise rerun affected role scope. Before a
clearly affected write make affected roles Pending with no round/revision and
provisionally mark unaffected roles Retained. After successful material write,
advance an unaffected retained pointer to the new Revision with its old round and
verdict. Ambiguous or large accumulated impact makes potentially affected roles
Pending before a human confirms scope. No acceptance remains valid until each role
is freshly Reviewed or explicitly Retained for current Revision.
