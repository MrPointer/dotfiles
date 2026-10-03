function _opencode_nono() {
  local profile="$1"
  shift

  local command_name=opencode
  if [[ "${1:-}" == acp ]]; then
    shift
    command_name="$HOME/.agents/shell/opencode-acp"
    local -x NONO_AGENT_ACP=1
  fi

  _nono_agent_run \
    "$profile" \
    "$HOME/.config/opencode/gitconfig" \
    "$command_name" \
    "${OPENCODE_SSH_SIGNING_PUBLIC_KEY:-$HOME/.ssh/id_ed25519.pub}" \
    "${OPENCODE_SSH_SIGNING_KEY:-$HOME/.ssh/id_ed25519}" \
    "$@"
}

function opencode() { _opencode_nono my-opencode "$@" }

function opencode-go() { _opencode_nono opencode-go "$@" }

function opencode-rust() { _opencode_nono opencode-rust "$@" }

function opencode-chezmoi() { _opencode_nono opencode-chezmoi "$@" }
