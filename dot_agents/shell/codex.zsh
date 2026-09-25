function _codex_nono_run() {
  local -x NPM_CONFIG_CACHE="$HOME/.codex/npm-cache"
  local -x npm_config_cache="$NPM_CONFIG_CACHE"
  mkdir -p "$NPM_CONFIG_CACHE" || return

  _nono_agent_run "$@"
}

function _codex() {
  local profile="$1"
  shift

  _codex_nono_run \
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

  _codex_nono_run \
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
