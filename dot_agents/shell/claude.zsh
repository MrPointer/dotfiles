function _claude_nono() {
  local profile="$1"
  shift

  _nono_agent_run \
    "$profile" \
    "$HOME/.claude/gitconfig" \
    claude \
    "${CLAUDE_SSH_SIGNING_PUBLIC_KEY:-$HOME/.ssh/id_ed25519.pub}" \
    "${CLAUDE_SSH_SIGNING_KEY:-$HOME/.ssh/id_ed25519}" \
    --dangerously-skip-permissions \
    "$@"
}

function claude() { _claude_nono my-claude-code "$@" }

function claude-chezmoi() { _claude_nono claude-code-chezmoi "$@" }

function claude-go() { _claude_nono claude-code-go "$@" }

function claude-node() { _claude_nono claude-code-node "$@" }

function claude-python() { _claude_nono claude-code-python "$@" }

function claude-rust() { _claude_nono claude-code-rust "$@" }
