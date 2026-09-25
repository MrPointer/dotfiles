function _codex() {
  local profile="$1"
  shift

  _nono_agent_run \
    "$profile" \
    "$HOME/.codex/gitconfig" \
    codex \
    "${CODEX_SSH_SIGNING_PUBLIC_KEY:-$HOME/.ssh/id_ed25519.pub}" \
    "${CODEX_SSH_SIGNING_KEY:-$HOME/.ssh/id_ed25519}" \
    "$@"
}

function _codex_acp() {
  local profile="$1"
  shift

  local NONO_AGENT_ACP=1
  local -x INITIAL_AGENT_MODE=agent-full-access

  _nono_agent_run \
    "$profile" \
    "$HOME/.codex/gitconfig" \
    npx \
    "${CODEX_SSH_SIGNING_PUBLIC_KEY:-$HOME/.ssh/id_ed25519.pub}" \
    "${CODEX_SSH_SIGNING_KEY:-$HOME/.ssh/id_ed25519}" \
    -y @agentclientprotocol/codex-acp@latest \
    "$@"
}

function codex() { _codex my-codex "$@" }

function codex-chezmoi() { _codex codex-chezmoi "$@" }

function codex-acp() { _codex_acp my-codex "$@" }

function codex-chezmoi-acp() { _codex_acp codex-chezmoi "$@" }
