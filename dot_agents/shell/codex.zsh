function codex() {
  _nono_agent_run \
    my-codex \
    "$HOME/.codex/gitconfig" \
    codex \
    "${CODEX_SSH_SIGNING_PUBLIC_KEY:-$HOME/.ssh/id_ed25519.pub}" \
    "${CODEX_SSH_SIGNING_KEY:-$HOME/.ssh/id_ed25519}" \
    "$@"
}

function codex-acp() {
  local NONO_AGENT_ACP=1
  local -x INITIAL_AGENT_MODE=agent-full-access

  _nono_agent_run \
    my-codex \
    "$HOME/.codex/gitconfig" \
    npx \
    "${CODEX_SSH_SIGNING_PUBLIC_KEY:-$HOME/.ssh/id_ed25519.pub}" \
    "${CODEX_SSH_SIGNING_KEY:-$HOME/.ssh/id_ed25519}" \
    -y @agentclientprotocol/codex-acp@latest \
    "$@"
}
