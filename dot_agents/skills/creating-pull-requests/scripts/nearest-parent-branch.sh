#!/usr/bin/env bash
#
# Print the nearest parent branch of HEAD — the local branch the current branch
# most likely built on — for inferring a stacked PR's parent when it isn't given
# explicitly. Run it from inside the target repository.
#
# Git records nothing about what a branch was forked from, so this reconstructs
# the relationship from commit topology: it ranks every local branch by how few
# commits HEAD is ahead of that branch's divergence point (merge-base) and prints
# the closest few, nearest first, as tab-separated "<commits-ahead>\t<branch>".
# The top line is the best guess; equal top counts mean it is genuinely ambiguous.

current=$(git branch --show-current)
head_rev=$(git rev-parse HEAD)

git for-each-ref --format='%(refname:short)' refs/heads/ | while read -r b; do
  # Skip the current branch itself.
  [ "$b" = "$current" ] && continue
  # Skip strict descendants of HEAD (children further up the stack). A branch
  # sitting exactly at HEAD is kept, so the trunk still wins before the current
  # branch has any commits of its own.
  if git merge-base --is-ancestor HEAD "$b" && [ "$(git rev-parse "$b")" != "$head_rev" ]; then
    continue
  fi
  base=$(git merge-base HEAD "$b")
  printf '%s\t%s\n' "$(git rev-list --count "$base..HEAD")" "$b"
done | sort -n | head -3
