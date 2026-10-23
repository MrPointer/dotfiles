function _claude_nono() {
  local profile="$1"
  shift

  _nono_agent_run \
    "$profile" \
    claude \
    "${CLAUDE_SSH_SIGNING_PUBLIC_KEY:-$HOME/.ssh/id_ed25519.pub}" \
    "${CLAUDE_SSH_SIGNING_KEY:-$HOME/.ssh/id_ed25519}" \
    --dangerously-skip-permissions \
    "$@"
}

function _claude_acp() {
  local profile="$1"
  shift

  # npm uses the uppercase form directly; lifecycle tooling may read the lowercase form.
  local -x NPM_CONFIG_CACHE="$HOME/.claude/npm-cache"
  local -x npm_config_cache="$NPM_CONFIG_CACHE"
  mkdir -p "$NPM_CONFIG_CACHE" || return

  local NONO_AGENT_ACP=1

  _nono_agent_run \
    "$profile" \
    npx \
    "${CLAUDE_SSH_SIGNING_PUBLIC_KEY:-$HOME/.ssh/id_ed25519.pub}" \
    "${CLAUDE_SSH_SIGNING_KEY:-$HOME/.ssh/id_ed25519}" \
    -y @agentclientprotocol/claude-agent-acp@latest \
    "$@"
}

function claude() { _claude_nono my-claude-code "$@" }

function claude-chezmoi() { _claude_nono claude-code-chezmoi "$@" }

function claude-go() { _claude_nono claude-code-go "$@" }

function claude-node() { _claude_nono claude-code-node "$@" }

function claude-python() { _claude_nono claude-code-python "$@" }

function claude-rust() { _claude_nono claude-code-rust "$@" }

function claude-acp() { _claude_acp my-claude-code "$@" }

function claude-chezmoi-acp() { _claude_acp claude-code-chezmoi "$@" }
