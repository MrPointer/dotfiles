# Opt-in sandbox features per profile. These run when the file is sourced, not
# inside a launcher, because nono-agent-acp starts profiles by name without
# calling any launcher function below.
_nono_agent_enable codex-chezmoi docker

function _codex_nono_with_npm_cache() {
  # npm uses the uppercase form directly; lifecycle tooling may read the lowercase form.
  local -x NPM_CONFIG_CACHE="$HOME/.codex/npm-cache"
  local -x npm_config_cache="$NPM_CONFIG_CACHE"
  mkdir -p "$NPM_CONFIG_CACHE" || return

  _nono_agent_run "$@"
}

function _codex_nono() {
  local profile="$1"
  shift

  _codex_nono_with_npm_cache \
    "$profile" \
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

  _codex_nono_with_npm_cache \
    "$profile" \
    npx \
    "${CODEX_SSH_SIGNING_PUBLIC_KEY:-$HOME/.ssh/id_ed25519.pub}" \
    "${CODEX_SSH_SIGNING_KEY:-$HOME/.ssh/id_ed25519}" \
    -y @agentclientprotocol/codex-acp@latest \
    "$@"
}

function codex() { _codex_nono my-codex "$@" }

function codex-chezmoi() { _codex_nono codex-chezmoi "$@" }

function codex-acp() { _codex_acp my-codex "$@" }

function codex-chezmoi-acp() { _codex_acp codex-chezmoi "$@" }
