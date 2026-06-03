function codex() {
  _nono_agent_run \
    my-codex \
    "$HOME/.config/opencode/gitconfig" \
    codex \
    "${CODEX_SSH_SIGNING_PUBLIC_KEY:-$HOME/.ssh/id_ed25519.pub}" \
    "${CODEX_SSH_SIGNING_KEY:-$HOME/.ssh/id_ed25519}" \
    "$@"
}
