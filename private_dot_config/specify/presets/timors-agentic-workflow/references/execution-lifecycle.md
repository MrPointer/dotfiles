# Execution Lifecycle and Runtime Progress

Before progress creation, require Ready `spec.md`; current Revision Accepted
`plan.md`; a closed task package with both current non-blocking reviews;
`authorization: approved`; and fully clean named stable active branch. Pin branch
and starting commit. External checkout, reset, rebase, or unrelated commit blocks
instead of being reconciled automatically. Resolve the exact progress path, add it
to `.git/info/exclude` if absent, then create it from
`templates/progress-template.md`. Successful initial creation is the task-package
lock event. Never change `tasks.md` or a sub-plan after that event.

Progress has exactly these sections: Execution Context, Packet States, Execution
Audit, Test and Verification Evidence, Runtime Consideration Evidence, Contract
Outputs, Blockers and Recovery, Final Verification. Package status is
`in-progress`, `blocked`, `complete`, or `abandoned`; packet state is exactly
`pending`, `in-progress`, `ready`, `blocked`, or `done`. It records feature
directory, active branch/start commit, protected baseline, signing policy evidence,
canonical identity and derived worktree-root-relative locator for each packet,
worktree/branch attribution, and all durable runtime transitions.

At lock, pin the complete path set and content of every tracked file beneath the
canonical feature directory. The baseline includes every feature-local artifact,
including spec, RFC, package, sub-plans, reports, grounding, research, and tracked
additions. Progress is excluded. At dispatch, result acceptance, after rebase, and
integration, compare paths, content, tracked additions/deletions/replacements, and
file types to this exact baseline. Any difference rejects dispatch or candidate;
never accept a changed protected path.

Before dispatch validate canonical identity is a tracked regular file under feature
directory and derive the locator from worktree root plus feature location. It must
remain beneath the directory and resolve to the same selected file. Send identity
and locator, not the packet body. Permit repository browsing but no authority
expansion. For each consumed contract, resolve the exact producer/name, revalidate
its body definition and final-commit evidence, and send only that evidence plus
attributable producer/name reference. The worker reads definition in the producer
sub-plan; neither progress nor prompt reconstructs it.

Workers may make packet-scoped commits only in their worktree/branch. The
coordinator verifies complete diff, declared ownership, tests, acceptance, output
evidence, result attribution, clean returned branch, policy, and baseline before
marking `ready` or integrating. Contract Outputs has exactly one entry per globally
unique produced name with producer identity, final integrated commit, required
returned evidence, validation outcome, and consumers derived from all consumes.
Missing, extra, duplicate, or contradictory entry blocks dependent readiness.

An unstarted invocation may retry once only with evidence of no start and no
workspace mutation. Missing result/test evidence gets one request to the same
session. A known-started runtime failure or ambiguous attempt is never automatically
redispatched. Resume or replacement binding first preserves and reconciles branch,
worktree, mutations, invocation, verification, ownership, locator, and baseline
evidence; a same-session resume continues only from clean attributable evidence.
Changed authority or ambiguous mutation requires abandonment and replanning. Never
infer progress from commits.

Abandonment before completion sets status `abandoned`, archives progress to a unique
ignored timestamped path, retains branches/worktrees unless explicit disposition
allows cleanup, returns authorization to pending, and permits fresh reviewed tasks
against current state. Active-branch commits remain untouched. Completion freezes
the package permanently; later work uses a new feature package.
