# Azure DevOps SSH auth uses the RSA key, so make sure the host agent holds it.
function _sedg_claude_nono() {
  _nono_agent_ensure_ssh_signing_key "$HOME/.ssh/id_rsa.pub" "$HOME/.ssh/id_rsa"
  _claude_nono "$@"
}

function claude-sedg() { _sedg_claude_nono work-sedg-claude-code "$@" }

function claude-energy-runtime() { _sedg_claude_nono work-sedg-claude-code-energy-runtime "$@" }

function claude-inverter-driver() { _sedg_claude_nono work-sedg-claude-code-inverter-driver "$@" }

function claude-lcp-docs() { _sedg_claude_nono work-sedg-claude-code-lcp-docs "$@" }

function claude-edge-metadata() { _sedg_claude_nono work-sedg-claude-code-edge-metadata "$@" }
