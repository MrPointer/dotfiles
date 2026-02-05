# Lazy NVM loading - only loads nvm when actually needed
# Saves ~500ms on every shell start by deferring nvm.sh sourcing
#
# Features:
# - Wrapper functions (nvm, node, npm, npx) that load nvm on first use
# - Git-root aware .nvmrc detection
# - Auto-switching when entering/leaving projects with .nvmrc
# - Minimal overhead: precmd detector only does file checks (no subshells)
# - chpwd hook only active while in a project

# Requires BREW_HOME to be set
[[ -z "$BREW_HOME" ]] && return
[[ ! -f "$BREW_HOME/opt/nvm/nvm.sh" ]] && return

# Export NVM_DIR for nvm to work properly
export NVM_DIR="$HOME/.nvm"

# Find .nvmrc from current dir up to git root (or current dir only if not in git repo)
_nvm_find_project_nvmrc() {
    local dir="$PWD" git_root=""

    # Find git root without subshell
    local check_dir="$dir"
    while [[ -n "$check_dir" && "$check_dir" != "/" ]]; do
        [[ -d "$check_dir/.git" ]] && { git_root="$check_dir"; break; }
        check_dir="${check_dir:h}"
    done

    # Not in git repo: only check current directory
    if [[ -z "$git_root" ]]; then
        [[ -f "$dir/.nvmrc" ]] && echo "$dir/.nvmrc"
        return
    fi

    # In git repo: check from current dir up to git root
    while [[ -n "$dir" ]]; do
        [[ -f "$dir/.nvmrc" ]] && { echo "$dir/.nvmrc"; return; }
        [[ "$dir" == "$git_root" ]] && break
        dir="${dir:h}"
    done
}

_nvm_load() {
    # Guard against re-entrancy
    [[ -n "$_NVM_LOADED" ]] && return 0
    _NVM_LOADED=1

    # Remove wrapper functions BEFORE sourcing nvm.sh
    unfunction nvm node npm npx 2>/dev/null

    # Load nvm
    source "$BREW_HOME/opt/nvm/nvm.sh"
}

# Full hook - only active while in a project
_nvm_chpwd_hook() {
    local nvmrc_path="$(_nvm_find_project_nvmrc)"

    if [[ -n "$nvmrc_path" ]]; then
        # Still in project (maybe different subdir)
        local wanted="$(cat "$nvmrc_path")"
        local current="$(nvm version)"
        local wanted_resolved="$(nvm version "$wanted" 2>/dev/null)"

        if [[ "$wanted_resolved" == "N/A" ]]; then
            nvm install
        elif [[ "$wanted_resolved" != "$current" ]]; then
            nvm use
        fi
        _NVM_PROJECT_ROOT="${nvmrc_path:h}"
    else
        # Left project - revert and deactivate
        echo "Reverting to nvm default version"
        nvm use default 2>/dev/null
        unset _NVM_PROJECT_ROOT
        add-zsh-hook -d chpwd _nvm_chpwd_hook
    fi
}

# Wrapper functions that load nvm on first use, then call the real command
nvm()  { _nvm_load; nvm "$@" }
node() { _nvm_load; node "$@" }
npm()  { _nvm_load; npm "$@" }
npx()  { _nvm_load; npx "$@" }

# Lightweight precmd detector - watches for project entry when not in project mode
_nvm_precmd_detect() {
    # Fast return if already in project mode (chpwd hook handles everything)
    [[ -n "$_NVM_PROJECT_ROOT" ]] && return

    local nvmrc_path="$(_nvm_find_project_nvmrc)"
    if [[ -n "$nvmrc_path" ]]; then
        # Entering a project - load nvm, apply version, activate chpwd hook
        _nvm_load
        _nvm_chpwd_hook
        add-zsh-hook chpwd _nvm_chpwd_hook
    fi
}

autoload -U add-zsh-hook
add-zsh-hook precmd _nvm_precmd_detect
