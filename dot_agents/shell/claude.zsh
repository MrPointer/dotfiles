function _claude_nono() {
  local profile="$1"
  shift

  _nono_agent_ensure_ssh_signing_keys \
    "${CLAUDE_SSH_SIGNING_PUBLIC_KEY:-$HOME/.ssh/id_rsa.pub}" \
    "${CLAUDE_SSH_SIGNING_KEY:-$HOME/.ssh/id_rsa}" \
    "${CLAUDE_ED25519_SSH_SIGNING_PUBLIC_KEY:-$HOME/.ssh/id_ed25519.pub}" \
    "${CLAUDE_ED25519_SSH_SIGNING_KEY:-$HOME/.ssh/id_ed25519}"

  _nono_agent_run \
    "$profile" \
    "$HOME/.claude/gitconfig" \
    claude \
    "${CLAUDE_SSH_SIGNING_PUBLIC_KEY:-$HOME/.ssh/id_rsa.pub}" \
    "${CLAUDE_SSH_SIGNING_KEY:-$HOME/.ssh/id_rsa}" \
    --dangerously-skip-permissions \
    "$@"
}

function claude() { _claude_nono my-claude-code "$@" }

function claude-go() { _claude_nono claude-code-go "$@" }

function claude-lcp-docs() { _claude_nono claude-code-lcp-docs "$@" }

function claude-python() { _claude_nono claude-code-python "$@" }

function claude-rust() { _claude_nono claude-code-rust "$@" }

function claude-typescript() { _claude_nono claude-code-typescript "$@" }
