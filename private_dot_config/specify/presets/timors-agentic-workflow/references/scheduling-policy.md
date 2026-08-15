# Scheduling and Integration Order

The frontmatter predecessor DAG is the only machine-authoritative scheduling
contract. Every `after` entry blocks readiness equally regardless of its human
classification. A packet is ready only after every direct predecessor is `done`.
Do not infer readiness from Mermaid, ledger prose, table reason, consumer prose,
or an undeclared relation.

Dispatch every eligible ready packet concurrently unless a concrete binding,
harness, machine-resource, or workspace limitation prevents it. Record the
limitation, then serialize only affected otherwise-ready packets in topological
list order without changing the approved graph. A worker finish may occur in any
order, but integration is always the approved topological packet order. A verified
later sibling remains `ready` until earlier integration completes.

Continue independent ready work after another packet blocks. A wave is quiescent
only when no pending or blocked packet can safely dispatch, resume, verify, or
integrate without resolving collected binding, runtime, verification, or recovery
blockers. Then present one interaction containing all currently unresolved
blockers. After a response resume scheduling; only a later quiescent wave can
produce another batch.
