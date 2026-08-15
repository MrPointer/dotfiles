# Worktree Isolation, Signing Policy, and Integration

Every packet has a dedicated branch and worktree. Workers may modify and commit
only selected packet-owned implementation or documentation paths outside protected
feature directory. Only the coordinator advances the explicit active target.
Workers never alter protected planning files, rewrite prerequisite history, merge
unrelated work, push, stash, reset, revert, rewrite, or auto-commit unrelated
files. Retain blocked or unresolved branches/worktrees and evidence.

Before dispatch, resume, result acceptance, rebase, squash, amend, corrective
commit, or integration, resolve signing policy against repository-controlled
instruction files and effective Git configuration for the target repository.
Apply repository instruction scope hierarchy: nearest scoped explicit signing
directive wins; contradictory directives at equal precedence require one human
resolution. Global agent instructions and prompt text are excluded. Without an
applicable directive, normal Git configuration precedence decides whether signing
is enabled. Signed history is evidence, never policy.

Resolve exactly one mode: `required` when winning repository instruction requires
signatures; `enabled` when no instruction requires/prohibits and effective Git
configuration enables signing; `disabled` when a winning instruction prohibits
signing or no directive applies and configuration does not enable it. Absence of a
directive and enabled configuration is disabled without prompting. Record mode,
every applicable instruction path/content digest, effective signing values and
reported origins, active target observation commit, and human-resolution evidence.
Observation commit is provenance, not a policy-equality input.

Before each named operation resolve again. Equality compares mode, instruction
paths/digests, effective values/origins, and human resolution, excluding only
observation commit. An unambiguous change automatically audits and repins,
including transition to disabled; reject and recreate any candidate formed under
the old policy. Prompt only for unresolved equal-precedence contradiction. A
candidate that changes repository instruction is governed by pre-integration active
target policy; after advancement resolve, audit, and repin before later work.

Every worker/coordinator commit and operation recreating a commit applies pinned
mode. Required and enabled create and verify signatures; signing failure blocks
with no unsigned fallback. Disabled does not invoke signing. No workflow default
changes repository policy.

Before target advancement, rebase packet branch on current active HEAD, recheck
baseline, and construct a separately inspectable squash candidate. Verify candidate
tree, required signature, ownership, and output evidence before a no-new-commit
fast-forward. Worktrunk is preferred only after a non-mutating capability check
proves it can make this candidate, apply pinned mode, and avoid early target
advancement. If not, plain Git creates and verifies the candidate independently.
Resolve only conflict-local mechanics outside the protected baseline; a semantic or
scope choice blocks. Advance active branch by exactly one policy-compliant commit
per packet, record it, then make consumers ready.
