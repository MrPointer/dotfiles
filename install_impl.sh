#!/usr/bin/env bash

# saner programming env: these switches turn some bugs into errors
set -o pipefail -o nounset

function show_usage {
    cat <<DOTFILES_INSTALL_IMPL_USAGE
Usage: $PROGRAM_NAME [OPTION]... 

Implementation of the dotfiles installation, executed by the plain old shell wrapper, for compatibility reasons.

Example: $PROGRAM_NAME -h

Options:
  -h, --help                        Show this message and exit
  -v, --verbose                     Enable verbose output
  --package-manager=[manager]       Package manager to use for installing prerequisites
  --work-environment                Treat this installation as a work environment
  --work-email=[email]              Use given email address as work's email address
  --shell=[shell]                   Install given shell if required and set it as user's default. Defaults to zsh.
  --no-python                       Don't install python
  --no-gpg                          Don't install gpg
  --no-brew                         Don't install brew (Homebrew)
  --prefer-package-manager          Prefer installing tools with system's package manager rather than brew (Doesn't apply for Mac)
-----------------------------------------------------"
DOTFILES_INSTALL_IMPL_USAGE
}

function cecho {
    printf "${1}%s${NEUTRAL_COLOR}\n" "${@:2}"
}
function warning {
    cecho "$YELLOW_COLOR" "$@"
}
function error {
    cecho "$RED_COLOR" "$@" >&2
}
function info {
    cecho "$BLUE_COLOR" "$@"
}
function success {
    cecho "$GREEN_COLOR" "$@"
}

###
# Join strings, just as in Python's str.join().
# Arguments:
#       $1 - String to join with (e.g. ',')
#       $2..$N - Variable number of strings to join
# Output (stdout):
#       Single string representing the joined string
###
function join_by {
    local d=${1-} f=${2-}
    if shift 2; then printf %s "$f" "${@/#/$d}"; fi
}

###
# Checks whether the current user is root.
# Returns:
#       0 if the user is root, 1 otherwise.
###
function root_user {
    local current_uid
    current_uid=$(id -u)

    ((current_uid == 0 ))
}

function _install_packages_with_brew {
    local packages=("$@")

    install_package_cmd=(brew install "${packages[@]}")

    eval "${install_package_cmd[@]}"
}

function _install_packages_with_package_manager {
    local packages=("$@")

    if [ -z "$PACKAGE_MANAGER" ]; then
        error "Package manager not set, something went wrong. Please install packages manually."
        return 1
    fi

    local install_package_cmd=()
    if [ "$ROOT_USER" = false ]; then
        install_package_cmd=(sudo)
    fi
    install_package_cmd+=("$PACKAGE_MANAGER" install -y "${packages[@]}")

    eval "${install_package_cmd[@]}"
}

###
# Install given package(s) using either system's package manager or homebrew, depending on the passed options.
# Arguments:
#       $1..$N - Variable number of packages to install
# Returns:
#       Install tool's result, zero on success.
###
function install_packages {
    if [ "$PREFER_BREW_FOR_ALL_TOOLS" = true ]; then
        _install_packages_with_brew "$@"
    else
        _install_packages_with_package_manager "$@"
    fi
}

function _reinstall_chezmoi_as_package {
    [ "$INSTALL_BREW" = false ] && return 0

    local brew_chezmoi_installed=false

    if brew list | grep -q "$DOTFILES_MANAGER"; then
        brew_chezmoi_installed=true
    fi

    if [ "$brew_chezmoi_installed" = false ]; then
        [ "$VERBOSE" = true ] && info "Installing $DOTFILES_MANAGER using brew"

        if ! brew install "$DOTFILES_MANAGER"; then
            error "Failed installing $DOTFILES_MANAGER using brew, will keep existing binary at $CHEZMOI_BINARY_PATH"
            return 1
        fi
        [ "$VERBOSE" = true ] && success "Successfully installed $DOTFILES_MANAGER as a brew package"
    fi

    [ "$VERBOSE" = true ] && info "Removing standalone $DOTFILES_MANAGER binary"
    if ! rm "$CHEZMOI_BINARY_PATH"; then
        warning "Failed removing standalone chezmoi binary (downloaded at first) at $CHEZMOI_BINARY_PATH"
    else
        [ "$VERBOSE" = true ] && success "Successfully removed standalone chezmoi binary"
    fi

    return 0
}

###
# Finalize installation by executing post-install commands.
###
function post_install {
    [ "$VERBOSE" = true ] && info "Executing post-install commands (finalization)"

    if ! _reinstall_chezmoi_as_package; then
        error "Failed reinstalling chezmoi as an updatable package"
        return 1
    fi
    return 0
}

###
# Apply dotfiles, optionally by using a dotfiles manager.
###
function apply_dotfiles {
    [ "$VERBOSE" = true ] && info "Applying dotfiles"

    if ! eval "${APPLY_DOTFILES_CMD[@]}"; then
        return 1
    fi
    return 0
}

###
# Prepare dotfiles environment before applying dotfiles.
# This might be a useful step for some dotfiles managers.
###
function prepare_dotfiles_environment {
    info "Preparing dotfiles environment"

    if ! mkdir -p "$ENVIRONMENT_TEMPLATE_CONFIG_DIR" &> /dev/null; then
        error "Couldn't create environment's dotfiles config directory"
        return 1
    fi

    # The first print zeroes the template file if it already has content
    if ! printf "%s\n" "[data]" >"$ENVIRONMENT_TEMPLATE_FILE_PATH"; then
        error "Failed initializing environment template file!"
        return 2
    fi

    {
        printf "%s\n" "[data.personal]"
        printf "%s\n" "work_env = $WORK_ENVIRONMENT"
        printf "%s\n" "full_name = \"$FULL_NAME\""
        printf "%s\n" "email = \"$ACTIVE_EMAIL\""
        printf "%s\n" "signing_key = \"$ACTIVE_GPG_SIGNING_KEY\""

        printf "%s\n" "[data.system]"
        printf "%s\n" "shell = \"$SHELL_TO_INSTALL\""

        printf "%s\n" "[data.installed]"
        printf "%s\n" "python = $INSTALL_PYTHON"
        printf "%s\n" "gpg = $INSTALL_GPG"
        printf "%s\n" "brew = $INSTALL_BREW"

        printf "%s\n" "[data.install_config]"
        printf "%s\n" "prefer_brew = $PREFER_BREW_FOR_ALL_TOOLS"
    } >>"$ENVIRONMENT_TEMPLATE_FILE_PATH"

    local text_inlined_tools
    text_inlined_tools="$(join_by , "${PACKAGE_MANAGER_INSTALLED_TOOLS[@]}")"
    printf "%s\n" "system_tools = ${text_inlined_tools}" >>"$ENVIRONMENT_TEMPLATE_FILE_PATH"
}

###
# Install selected shell using either system's package manager or homebrew, depending on the passed options.
# If selected shell is already installed, do nothing.
# Otherwise, also configure it as user's default shell.
###
function install_shell {
    if hash "$SHELL_TO_INSTALL" &> /dev/null; then
        return 0
    fi

    [ "$VERBOSE" = true ] && echo "Installing shell"

    # First install the shell
    if ! install_packages "$SHELL_TO_INSTALL"; then
        return 1
    fi

    # Find installed shell's location
    local shell_path
    shell_path="$(which "$SHELL_TO_INSTALL")"

    local current_user_name
    current_user_name="$(id -u -n)"

    # Then configure it as user's default shell
    sudo chsh -s "$shell_path" "$current_user_name"
}

###
# Install git using either system's package manager or homebrew, depending on the passed options.
# If git is already installed, do nothing.
###
function install_git {
    if hash git 2>/dev/null; then
        return 0
    fi

    [ "$VERBOSE" = true ] && info "Installing git"

    install_packages git
}

###
# Install chezmoi, our dotfiles manager.
# To avoid any errors and complicated checks, just install the latest binary at this stage.
###
function install_dotfiles_manager {
    [ "$VERBOSE" = true ] && info "Installing dotfiles manager ($DOTFILES_MANAGER)"

    if hash "$DOTFILES_MANAGER" 2>/dev/null; then
        info "$DOTFILES_MANAGER already installed, skipping"
        return 0
    fi

    local installation_failed=false

    if [[ "$DOWNLOAD_TOOL" == "curl" ]]; then
        ! sh -c "$(curl -fsLS git.io/chezmoi)" && installation_failed=true
    elif [[ "$DOWNLOAD_TOOL" == "wget" ]]; then
        ! sh -c "$(wget -qO- git.io/chezmoi)" && installation_failed=true
    fi

    if [ "$installation_failed" = true ]; then
        return 1
    fi
    return 0
}

###
# Install dotfiles. This is the main "driver" function.
###
function install_dotfiles {
    if ! install_dotfiles_manager; then
        error "Failed installing dotfiles manager ($DOTFILES_MANAGER)"
        return 1
    fi

    if ! install_git; then
        error "Failed installing git"
    fi
    [ "$VERBOSE" = true ] && success "Successfully installed git"

    if ! install_shell; then
        error "Failed installing shell"
    fi
    [ "$VERBOSE" = true ] && success "Successfully installed shell"

    if ! prepare_dotfiles_environment; then
        error "Failed preparing dotfiles environment"
        return 2
    fi
    [ "$VERBOSE" = true ] && success "Successfully prepared dotfiles environment"

    if ! apply_dotfiles; then
        error "Failed applying dotfiles"
        return 3
    fi
    [ "$VERBOSE" = true ] && success "Successfully applied dotfiles"

    if ! post_install; then
        error "Failed finalizing installation"
        return 4
    fi
    [ "$VERBOSE" = true ] && success "Successfully finalized installation"

    return 0
}

###
# Checks which download tool is locally available from a preset list
# and outputs the first that has been found.
###
function get_download_tool {
    local optional_download_tools=(
        curl
        wget
    )

    for download_tool in "${optional_download_tools[@]}"; do
        if hash "${download_tool}" 2>/dev/null; then
            echo "${download_tool}"
            return 0
        fi
    done

    echo ""
    return 1
}

###
# Set global variables
###
function set_globals {
    if ! DOWNLOAD_TOOL="$(get_download_tool)"; then
        error "Couldn't determine download tool, aborting"
        return 1
    fi

    if root_user; then
        ROOT_USER=true
    fi

    if [ "$WORK_ENVIRONMENT" = true ]; then
        ACTIVE_EMAIL="$WORK_EMAIL"
        ACTIVE_GPG_SIGNING_KEY="$WORK_GPG_SIGNING_KEY"
    else
        ACTIVE_EMAIL="$PERSONAL_EMAIL"
        ACTIVE_GPG_SIGNING_KEY="$PERSONAL_GPG_SIGNING_KEY"
    fi
}

###
# Parse arguments/options using getopt, the almighty C-based parser.
###
function parse_arguments {
    getopt --test >/dev/null
    if (($? != 4)); then
        error "I'm sorry, 'getopt --test' failed in this environment."
        return 1
    fi

    local short_options=hv
    local long_options=help,verbose
    long_options+=,package-manager:
    long_options+=,shell:,no-python,no-gpg,no-brew,prefer-package-manager
    long_options+=,work-environment,work-email:

    # -temporarily store output to be able to check for errors
    # -activate quoting/enhanced mode (e.g. by writing out “--options”)
    # -pass arguments only via   -- "$@"   to separate them correctly
    if ! PARSED=$(
        getopt --options="$short_options" --longoptions="$long_options" \
            --name "Dotfiles installer" -- "$@"
    ); then
        # getopt has complained about wrong arguments to stdout
        error "Wrong arguments to Dotfiles installer" && return 2
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
        --package-manager)
            PACKAGE_MANAGER="${2:-}"
            shift 2
            ;;
        --work-environment)
            WORK_ENVIRONMENT=true
            shift
            ;;
        --work-email)
            WORK_EMAIL="${2:-}"
            shift 2
            ;;
        --shell)
            SHELL_TO_INSTALL="${2:-}"
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
        --prefer-package-manager)
            PREFER_BREW_FOR_ALL_TOOLS=false
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

function _set_package_management_defaults {
    PACKAGE_MANAGER=""
}

function _set_installed_tools_defaults {
    SHELL_TO_INSTALL=zsh
    INSTALL_GPG=true
    INSTALL_PYTHON=true
    INSTALL_BREW=true
    PACKAGE_MANAGER_INSTALLED_TOOLS=(gpg)
    PREFER_BREW_FOR_ALL_TOOLS=true
}

function _set_dotfiles_manager_defaults {
    DOTFILES_MANAGER=chezmoi
    CHEZMOI_BINARY_PATH="$HOME/bin/chezmoi"
    APPLY_DOTFILES_CMD=("$DOTFILES_MANAGER"
        init --apply "$GITHUB_USERNAME"
    )
    ENVIRONMENT_TEMPLATE_CONFIG_DIR="$HOME/.config/chezmoi"
    ENVIRONMENT_TEMPLATE_FILE_PATH="$ENVIRONMENT_TEMPLATE_CONFIG_DIR/chezmoi.toml"
}

function _set_personal_info_defaults {
    GITHUB_USERNAME="MrPointer"
    FULL_NAME="Timor Gruber"
    PERSONAL_EMAIL="timor.gruber@gmail.com"
    PERSONAL_GPG_SIGNING_KEY=E1B39E9320C37806
    WORK_EMAIL="timor.gruber@solaredge.com"
    WORK_GPG_SIGNING_KEY=90BBCCC1DDED66C4
}

###
# Set default color codes for colorful prints.
###
function _set_color_defaults {
    RED_COLOR="\033[0;31m"
    GREEN_COLOR="\033[0;32m"
    YELLOW_COLOR="\033[1;33m"
    BLUE_COLOR="\033[0;34m"
    NEUTRAL_COLOR="\033[0m"
}

###
# Set script default values for later show_usage.
###
function set_defaults {
    VERBOSE=false
    WORK_ENVIRONMENT=false
    ROOT_USER=false

    _set_color_defaults
    _set_personal_info_defaults
    _set_dotfiles_manager_defaults
    _set_installed_tools_defaults
    _set_package_management_defaults
}

###
# This is the script's entry point, just like in any other programming language.
###
function main {
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
