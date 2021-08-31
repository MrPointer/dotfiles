#!/usr/bin/env bash

# saner programming env: these switches turn some bugs into errors
set -o pipefail -o nounset

function show_usage() {
    cat <<DOTFILES_INSTALL_IMPL_USAGE
Usage: $PROGRAM_NAME [OPTION]... 

Implementation of the dotfiles installation, executed by the plain old shell wrapper, for compatibility reasons.

Example: $PROGRAM_NAME -h

Options:
  -h, --help                Show this message and exit
  -v, --verbose             Enable verbose output
  --work-environment        Treat this installation as a work environment
  --work-email=[email]      Use given email address as work's email address
  --no-python               Don't install python
  --no-gpg                  Don't install gpg
  --no-brew                 Don't install brew (Linuxbrew/Homebrew)
  --prefer-brew             Prefer installing "system" tools with brew rather than package manager (Doesn't apply for Mac)
-----------------------------------------------------"
DOTFILES_INSTALL_IMPL_USAGE
}

function cecho() {
    printf "${1}%s${NEUTRAL_COLOR}\n" "${@:2}"
}
function warning() {
    cecho "$YELLOW_COLOR" "$@"
}
function error() {
    cecho "$RED_COLOR" "$@" >&2
}
function info() {
    cecho "$BLUE_COLOR" "$@"
}
function success() {
    cecho "$GREEN_COLOR" "$@"
}

function join_by {
    local d=${1-} f=${2-}
    if shift 2; then printf %s "$f" "${@/#/$d}"; fi
}

function _reinstall_chezmoi_as_package() {
    [ "$INSTALL_BREW" = false ] && return 0

    info "Installing chezmoi using brew"
    if ! brew install chezmoi; then
        error "Failed installing chezmoi using brew, will keep existing binary at $CHEZMOI_BINARY_PATH"
        return 1
    fi
    success "Successfully installed chezmoi as a brew package"

    info "Removing standalone chezmoi binary"
    if ! rm "$CHEZMOI_BINARY_PATH"; then
        warning "Failed removing standalone chezmoi binary (downloaded at first) at $CHEZMOI_BINARY_PATH"
    else
        success "Successfully removed standalone chezmoi binary"
    fi

    return 0
}

function post_install() {
    info "Executing post-install commands (finalization)"

    if ! _reinstall_chezmoi_as_package; then
        error "Failed reinstalling chezmoi as an updatable package"
        return 1
    fi

    return 0
}

function apply_dotfiles() {
    info "Applying dotfiles"

    if ! eval "${APPLY_DOTFILES_CMD[@]}"; then
        error "Failed applying dotfiles"
        return 1
    fi

    success "Successfully applied dotfiles"
    return 0
}

function prepare_dotfiles_environment() {
    info "Preparing dotfiles environment"

    # The first print zeroes the template file if it already has content
    if ! printf "%s\n" "[data]" >"$ENVIRONMENT_TEMPLATE_FILE_PATH"; then
        error "Failed initializing environment template file!"
        return 1
    fi

    {
        printf "%s\n" "[data.personal]"
        printf "%s\n" "email = $ACTIVE_EMAIL"

        printf "%s\n" "[data.installed]"
        printf "%s\n" "python = $INSTALL_PYTHON"
        printf "%s\n" "gpg = $INSTALL_GPG"
        printf "%s\n" "brew = $INSTALL_BREW"

        printf "%s\n" "[data.install_config]"
        printf "%s\n" "prefer_brew = $PREFER_BREW_FOR_ALL_TOOLS"
    } >>"$ENVIRONMENT_TEMPLATE_FILE_PATH"

    local text_inlined_tools
    text_inlined_tools="$(join_by , "${PACKAGE_MANAGER_INSTALLED_TOOLS[@]}")"
    printf "%s\n" "system_tools = ${text_inlined_tools}" >> "$ENVIRONMENT_TEMPLATE_FILE_PATH"
}

###
# Install chezmoi, our dotfiles manager.
# To avoid any errors and complicated checks, just install the latest binary at this stage.
###
function install_dotfiles_manager() {
    info "Installing dotfiles manager ($DOTFILES_MANAGER)"

    local installation_failed=false

    if [[ "$DOWNLOAD_TOOL" == "curl" ]]; then
        ! sh -c '$(curl -fsLS git.io/chezmoi)' && installation_failed=true
    elif [[ "$DOWNLOAD_TOOL" == "wget" ]]; then
        ! sh -c '$(wget -qO- git.io/chezmoi)' && installation_failed=true
    fi

    if [ "$installation_failed" = true ]; then
        error "Failed installing dotfiles manager ($DOTFILES_MANAGER)"
        return 1
    fi

    return 0
}

function install_dotfiles() {
    if ! install_dotfiles_manager; then
        error "Failed installing dotfiles manager"
        return 1
    fi

    if ! prepare_dotfiles_environment; then
        error "Failed preparing dotfiles environment"
        return 2
    fi

    if ! apply_dotfiles; then
        error "Failed applying dotfiles"
        return 3
    fi

    if ! post_install; then
        error "Failed executing post-install stuff (finalization)"
        return 4
    fi

    return 0
}

function get_download_tool() {
    if hash curl 2>/dev/null; then
        echo "curl"
    elif hash wget 2>/dev/null; then
        echo "wget"
    else
        echo ""
    fi
}

###
# Set global variables
###
function set_globals() {
    DOWNLOAD_TOOL="$(get_download_tool)"
    if [ -z "$DOWNLOAD_TOOL" ]; then
        error "Couldn't determine download tool, aborting"
        return 1
    fi

    if [ "$WORK_ENVIRONMENT" = true ]; then
        ACTIVE_EMAIL="$WORK_EMAIL"
    else
        ACTIVE_EMAIL="$PERSONAL_EMAIL"
    fi
}

###
# Parse arguments/options using getopt, the almighty C-based parser.
###
function parse_arguments() {
    getopt --test >/dev/null
    if (($? != 4)); then
        error "I'm sorry, 'getopt --test' failed in this environment."
        return 1
    fi

    local short_options=hv
    local long_options=help,verbose
    long_options+=,no-python,no-gpg,no-brew,prefer-brew
    long_options+=,work-environment,work-email

    # -temporarily store output to be able to check for errors
    # -activate quoting/enhanced mode (e.g. by writing out “--options”)
    # -pass arguments only via   -- "$@"   to separate them correctly
    if ! PARSED=$(
        getopt --options="$short_options" --longoptions="$long_options" \
            --name "$PROGRAM_PATH" -- "$@"
    ); then
        # getopt has complained about wrong arguments to stdout
        error "Wrong arguments to $PROGRAM_NAME" && return 2
    fi

    # read getopt’s output this way to handle the quoting right:
    eval set -- "$PARSED"

    while true; do
        case $1 in
        -h | --help)
            show_usage
            exit 0
            ;;
        -v | --verbose)
            VERBOSE=true
            shift
            ;;
        --work-environment)
            WORK_ENVIRONMENT=true
            shift
            ;;
        --work-email)
            WORK_EMAIL="$2"
            shift 2
            ;;
        --no-python)
            INSTALL_PYTHON=false
            shift
            ;;
        --no-gpg)
            INSTALL_GPG=false
            shift
            ;;
        --no-brew)
            INSTALL_BREW=false
            shift
            ;;
        --prefer-brew)
            PREFER_BREW_FOR_ALL_TOOLS=true
            shift
            ;;
        --)
            shift
            break
            ;;
        *)
            error "Programming error"
            return 3
            ;;
        esac
    done

    return 0
}

function _set_personal_info_defaults() {
    GITHUB_USERNAME="MrPointer"
    PERSONAL_EMAIL="timor.gruber@gmail.com"
    WORK_EMAIL="timor.gruber@solaredge.com"
}

function _set_installed_tools_defaults() {
    INSTALL_GPG=true
    INSTALL_PYTHON=true
    INSTALL_BREW=true
    PACKAGE_MANAGER_INSTALLED_TOOLS=(gpg)
    PREFER_BREW_FOR_ALL_TOOLS=false
}

function _set_dotfiles_manager_defaults() {
    DOTFILES_MANAGER=chezmoi
    CHEZMOI_BINARY_PATH="$HOME/bin/chezmoi"
    APPLY_DOTFILES_CMD=("$DOTFILES_MANAGER"
        init --apply "$GITHUB_USERNAME"
    )
    ENVIRONMENT_TEMPLATE_FILE_PATH="$HOME/.config/chezmoi/chezmoi.toml"
}

###
# Set default color codes for colorful prints.
###
function _set_color_defaults() {
    RED_COLOR="\033[0;31m"
    GREEN_COLOR="\033[0;32m"
    YELLOW_COLOR="\033[1;33m"
    BLUE_COLOR="\033[0;34m"
    NEUTRAL_COLOR="\033[0m"
}

###
# Set script default values for later show_usage.
###
function set_defaults() {
    VERBOSE=false
    WORK_ENVIRONMENT=false

    _set_color_defaults
    _set_dotfiles_manager_defaults
    _set_installed_tools_defaults
    _set_personal_info_defaults
}

###
# This is the script's entry point, just like in any other programming language.
###
function main() {
    if [[ -v ZSH_NAME && -n "$ZSH_NAME" ]]; then
        error "I'm a Bash script, please do not run me as 'zsh script_name'" \
            ", but rather execute me directly!"
        return 1
    fi

    if ! set_defaults; then
        error "Failed setting default values, aborting"
        return 1
    fi

    if ! parse_arguments "$@"; then
        error "Couldn't parse arguments, aborting"
        return 2
    fi

    if ! set_globals; then
        error "Failed setting global variables, aborting"
        return 1
    fi

    if ! install_dotfiles; then
        error "Failed installing dotfiles"
        return 3
    fi

    success "Successfully installed dotfiles!"
    return 0
}

# Call main and don't do anything else
# It will pass the correct exit code to the OS
main "$@"
