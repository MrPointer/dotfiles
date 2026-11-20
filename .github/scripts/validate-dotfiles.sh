#!/usr/bin/env bash
#
# Validate dotfiles in a home that `dotfiles-installer install --prepare-only`
# prepared: apply the checkout with chezmoi, then check pass criteria 1 to 7.
#
# Usage: validate-dotfiles.sh <bare|with-tools> <checkout-path> <chezmoi-config-path>
#
# Requires CHEZMOI_GITHUB_ACCESS_TOKEN in the environment. Runs as the user who
# owns the prepared home. Exits non-zero when any criterion fails, naming the
# criterion and, for criterion 7, each offending tool.
#
# Compatible with bash 3.2 (macOS system bash).

set -euo pipefail

usage() {
    echo "Usage: $0 <bare|with-tools> <checkout-path> <chezmoi-config-path>" >&2
    exit 2
}

[[ $# -eq 3 ]] || usage

VARIANT=$1
SOURCE_DIR=$2
CONFIG_PATH=$3

case "$VARIANT" in
    bare | with-tools) ;;
    *) usage ;;
esac

if [[ -z "${CHEZMOI_GITHUB_ACCESS_TOKEN:-}" ]]; then
    echo "ERROR: CHEZMOI_GITHUB_ACCESS_TOKEN must be set" >&2
    exit 2
fi
if [[ ! -d "$SOURCE_DIR" ]]; then
    echo "ERROR: checkout path '$SOURCE_DIR' is not a directory" >&2
    exit 2
fi
if [[ ! -f "$CONFIG_PATH" ]]; then
    echo "ERROR: chezmoi config '$CONFIG_PATH' does not exist" >&2
    exit 2
fi
if [[ -z "${HOME:-}" || ! -d "$HOME" ]]; then
    echo "ERROR: HOME must be set to the prepared home directory" >&2
    exit 2
fi

USER_NAME=${USER:-$(id -un)}

# The minimal environment of criterion 5: only these variables reach zsh.
MIN_PATH=/usr/bin:/bin:/usr/sbin:/sbin

# Optional tools of criterion 7: "<tool>:<command> [<alternate command>...]".
OPTIONAL_TOOLS=(
    "fzf:fzf"
    "bat:bat batcat"
    "sheldon:sheldon"
    "eza:eza"
    "ripgrep:rg"
    "fd:fd fdfind"
    "difftastic:difft"
    "uv:uv"
    "diff-so-fancy:diff-so-fancy"
)

WORK_DIR=$(mktemp -d)
trap 'rm -rf "$WORK_DIR"' EXIT

FAILED_CRITERIA=()

log() {
    printf '%s\n' "$*"
}

section() {
    printf '\n=== %s ===\n' "$*"
}

fail_criterion() {
    local criterion=$1
    shift
    printf 'FAIL [criterion %s]: %s\n' "$criterion" "$*"
    FAILED_CRITERIA+=("$criterion")
}

pass_criterion() {
    printf 'PASS [criterion %s]: %s\n' "$1" "$2"
}

# Locate a binary the installer may have placed outside the caller's PATH.
find_binary() {
    local name=$1
    (
        PATH="$HOME/.local/bin:/home/linuxbrew/.linuxbrew/bin:/opt/homebrew/bin:/usr/local/bin:$PATH"
        command -v "$name" 2>/dev/null
    ) || true
}

# The user's login shell when it is zsh, otherwise the zsh found by name.
find_zsh() {
    local login_shell=""
    if command -v dscl >/dev/null 2>&1; then
        login_shell=$(dscl . -read "/Users/$USER_NAME" UserShell 2>/dev/null | awk '{print $2}') || true
    elif command -v getent >/dev/null 2>&1; then
        login_shell=$(getent passwd "$USER_NAME" | cut -d: -f7) || true
    fi
    if [[ "$login_shell" == */zsh && -x "$login_shell" ]]; then
        printf '%s\n' "$login_shell"
        return
    fi
    find_binary zsh
}

CHEZMOI_BIN=$(find_binary chezmoi)
if [[ -z "$CHEZMOI_BIN" ]]; then
    echo "ERROR: chezmoi not found; the prepare step should have installed it" >&2
    exit 1
fi

# Every chezmoi command passes the same --source and --config. The command is
# logged to stderr so captured stdout stays clean; the token is inherited from
# the environment and never printed.
run_chezmoi() {
    log "+ $CHEZMOI_BIN $1 --source $SOURCE_DIR --config $CONFIG_PATH ${*:2}" >&2
    "$CHEZMOI_BIN" "$1" --source "$SOURCE_DIR" --config "$CONFIG_PATH" "${@:2}"
}

# Run zsh in the minimal environment of criterion 5, from the home directory,
# with no terminal on stdin. Usage: run_min_zsh <stdout-file> <stderr-file> <zsh-args>...
run_min_zsh() {
    local out=$1 err=$2
    shift 2
    (
        cd "$HOME" &&
            env -i \
                HOME="$HOME" \
                USER="$USER_NAME" \
                TERM=xterm-256color \
                PATH="$MIN_PATH" \
                "$ZSH_BIN" "$@" </dev/null >"$out" 2>"$err"
    )
}

show_output() {
    local label=$1 file=$2
    if [[ -s "$file" ]]; then
        log "--- $label ---"
        cat "$file"
        log "--- end $label ---"
    fi
}

log "Variant:        $VARIANT"
log "Source:         $SOURCE_DIR"
log "Config:         $CONFIG_PATH"
log "Home:           $HOME"
log "User:           $USER_NAME"
log "chezmoi:        $CHEZMOI_BIN ($("$CHEZMOI_BIN" --version 2>/dev/null || echo unknown))"

section "Criterion 1: chezmoi init --apply"
if run_chezmoi init --apply --no-tty; then
    pass_criterion 1 "chezmoi init --apply exited 0"
else
    fail_criterion 1 "chezmoi init --apply exited non-zero"
    log "Remaining criteria need an applied home; stopping."
    exit 1
fi

ZSH_BIN=$(find_zsh)
if [[ -z "$ZSH_BIN" ]]; then
    echo "ERROR: zsh not found; the prepare step should have installed it" >&2
    exit 1
fi
log "zsh:            $ZSH_BIN ($("$ZSH_BIN" --version 2>/dev/null || echo unknown))"

section "Criterion 2: chezmoi verify"
if run_chezmoi verify --exclude scripts; then
    pass_criterion 2 "every target matches the source state"
else
    fail_criterion 2 "chezmoi verify --exclude scripts exited non-zero"
fi

section "Criterion 3: managed targets exist"
# chezmoi does not create a file whose rendered content is empty (for example a
# modify script that passes empty stdin through), so a missing managed file is
# skipped when `chezmoi cat` renders it empty. Missing dirs and symlinks fail.
managed_file="$WORK_DIR/managed"
managed_files_file="$WORK_DIR/managed-files"
if run_chezmoi managed --include files,dirs,symlinks >"$managed_file" &&
    run_chezmoi managed --include files >"$managed_files_file"; then
    managed_count=0
    missing_count=0
    skipped_count=0
    while IFS= read -r target || [[ -n "$target" ]]; do
        [[ -n "$target" ]] || continue
        managed_count=$((managed_count + 1))
        if [[ -e "$HOME/$target" || -L "$HOME/$target" ]]; then
            continue
        fi
        if grep -Fxq -- "$target" "$managed_files_file"; then
            if ! rendered=$(run_chezmoi cat "$HOME/$target" </dev/null); then
                log "missing: ~/$target (chezmoi cat failed)"
                missing_count=$((missing_count + 1))
            elif [[ -z "${rendered//[[:space:]]/}" ]]; then
                log "skipped: ~/$target (rendered content is empty; chezmoi does not create it)"
                skipped_count=$((skipped_count + 1))
            else
                log "missing: ~/$target"
                missing_count=$((missing_count + 1))
            fi
        else
            log "missing: ~/$target"
            missing_count=$((missing_count + 1))
        fi
    done <"$managed_file"
    log "managed paths listed: $managed_count"
    log "empty-rendering files skipped: $skipped_count"
    if ((managed_count == 0)); then
        fail_criterion 3 "chezmoi managed listed no paths"
    elif ((missing_count > 0)); then
        fail_criterion 3 "$missing_count managed path(s) missing from the home"
    else
        pass_criterion 3 "all $((managed_count - skipped_count)) managed paths exist ($skipped_count empty-rendering file(s) skipped)"
    fi
else
    fail_criterion 3 "chezmoi managed exited non-zero"
fi

section "Criterion 4: zsh -n on startup files"
startup_files=("$HOME/.zshenv" "$HOME/.zshrc")
if [[ -d "$HOME/.work" ]]; then
    while IFS= read -r file; do
        startup_files+=("$file")
    done < <(find "$HOME/.work" \( -name .zshenv -o -name .zshrc \) \( -type f -o -type l \) | LC_ALL=C sort)
fi
syntax_failed=0
for file in "${startup_files[@]}"; do
    if [[ ! -f "$file" ]]; then
        log "missing startup file: $file"
        syntax_failed=1
    elif "$ZSH_BIN" -n "$file"; then
        log "ok: $file"
    else
        log "syntax error: $file"
        syntax_failed=1
    fi
done
if ((syntax_failed)); then
    fail_criterion 4 "zsh -n failed for at least one startup file"
else
    pass_criterion 4 "${#startup_files[@]} startup file(s) parse"
fi

section "Criterion 5: shell startup in a minimal environment"
log "environment: HOME USER TERM=xterm-256color PATH=$MIN_PATH"
startup_failed=0
warm_out="$WORK_DIR/warm.out"
warm_err="$WORK_DIR/warm.err"
log "+ zsh -i -c exit (warm-up; stderr not judged)"
if run_min_zsh "$warm_out" "$warm_err" -i -c exit; then
    log "warm-up exited 0"
else
    log "warm-up exited non-zero"
    startup_failed=1
fi
show_output "warm-up stdout" "$warm_out"
show_output "warm-up stderr" "$warm_err"

check_startup() {
    local label=$1
    shift
    local out="$WORK_DIR/startup.out" err="$WORK_DIR/startup.err"
    log "+ zsh $* ($label)"
    if run_min_zsh "$out" "$err" "$@"; then
        log "$label exited 0"
    else
        log "$label exited non-zero"
        startup_failed=1
    fi
    show_output "$label stdout" "$out"
    if [[ -s "$err" ]]; then
        log "$label wrote to stderr"
        show_output "$label stderr" "$err"
        startup_failed=1
    fi
}
check_startup "noninteractive startup" -c exit
check_startup "interactive startup" -i -c exit
if ((startup_failed)); then
    fail_criterion 5 "shell startup failed or wrote to stderr"
else
    pass_criterion 5 "noninteractive and interactive startup exit 0 with empty stderr"
fi

section "Criterion 6: git diff and core.pager"
marker_old="validate-dotfiles-original-line"
marker_new="validate-dotfiles-modified-line"
repo_dir="$WORK_DIR/repo"
mkdir -p "$repo_dir"
git_failed=0
(
    cd "$repo_dir" &&
        git init -q &&
        printf '%s\n' "$marker_old" >sample.txt &&
        git add sample.txt &&
        printf '%s\n' "$marker_new" >sample.txt
) || {
    log "could not create the temporary repository"
    git_failed=1
}

diff_out="$WORK_DIR/diff.out"
diff_err="$WORK_DIR/diff.err"
log "+ zsh -i -c 'cd <repo> && git diff'"
# shellcheck disable=SC2016 # expanded by the inner zsh
if run_min_zsh "$diff_out" "$diff_err" -i -c 'cd "$1" && git diff' zsh "$repo_dir"; then
    log "git diff exited 0"
else
    log "git diff exited non-zero"
    git_failed=1
fi
show_output "git diff stdout" "$diff_out"
show_output "git diff stderr" "$diff_err"
if ! grep -q "$marker_new" "$diff_out"; then
    log "git diff output does not show the change"
    git_failed=1
fi

sample_diff="$WORK_DIR/sample.diff"
cat >"$sample_diff" <<EOF
diff --git a/sample.txt b/sample.txt
index 0000001..0000002 100644
--- a/sample.txt
+++ b/sample.txt
@@ -1 +1 @@
-$marker_old
+$marker_new
EOF
pager_out="$WORK_DIR/pager.out"
pager_err="$WORK_DIR/pager.err"
log "+ zsh -i -c 'sh -c \"\$(git config --get core.pager)\" < sample.diff'"
# shellcheck disable=SC2016 # expanded by the inner zsh
pager_script='pager=$(git config --get core.pager) || { print -u2 "core.pager is not set"; exit 1; }
print -r -- "core.pager: $pager"
sh -c "$pager" <"$1"'
if run_min_zsh "$pager_out" "$pager_err" -i -c "$pager_script" zsh "$sample_diff"; then
    log "core.pager exited 0"
else
    log "core.pager exited non-zero"
    git_failed=1
fi
show_output "core.pager stdout" "$pager_out"
show_output "core.pager stderr" "$pager_err"
if ! sed 's/\x1b\[[0-9;]*m//g' "$pager_out" | grep -q "$marker_new"; then
    log "core.pager output does not show the change"
    git_failed=1
fi
if ((git_failed)); then
    fail_criterion 6 "git diff or core.pager failed"
else
    pass_criterion 6 "git diff and core.pager exit 0 and show the change"
fi

section "Criterion 7: optional-tool state ($VARIANT)"
tools_failed=()
for entry in "${OPTIONAL_TOOLS[@]}"; do
    tool=${entry%%:*}
    commands=${entry#*:}
    tool_out="$WORK_DIR/tool.out"
    tool_err="$WORK_DIR/tool.err"
    # shellcheck disable=SC2086 # split the command names into arguments
    # shellcheck disable=SC2016 # expanded by the inner zsh
    if run_min_zsh "$tool_out" "$tool_err" -i -c \
        'for c in "$@"; do command -v "$c" && exit 0; done; exit 1' zsh $commands; then
        log "present: $tool ($(tail -n 1 "$tool_out"))"
        [[ "$VARIANT" == bare ]] && tools_failed+=("$tool")
    else
        log "absent:  $tool (looked for: $commands)"
        [[ "$VARIANT" == with-tools ]] && tools_failed+=("$tool")
    fi
done
if ((${#tools_failed[@]} == 0)); then
    if [[ "$VARIANT" == bare ]]; then
        pass_criterion 7 "no optional tool is present"
    else
        pass_criterion 7 "every optional tool is present"
    fi
else
    for tool in "${tools_failed[@]}"; do
        if [[ "$VARIANT" == bare ]]; then
            fail_criterion 7 "optional tool present in a bare home: $tool"
        else
            fail_criterion 7 "optional tool missing in a with-tools home: $tool"
        fi
    done
fi

section "Summary"
if ((${#FAILED_CRITERIA[@]} > 0)); then
    log "FAILED criteria: $(printf '%s\n' "${FAILED_CRITERIA[@]}" | sort -u | tr '\n' ' ')"
    exit 1
fi
log "All criteria 1 to 7 passed ($VARIANT)"
