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
  --ref=[git-ref]                   Reference the given git-ref for installation (can be any git ref - commit, branch, tag). Defaults to 'main'
  --work-env                        Treat this installation as a work environment
  --work-name                       Use the given work-name as the work environment. Defaults to 'sedg' (current workplace)
  --work-email=[email]              Use given email address as work's email address. Defaults to 'timor.gruber@solaredge.com'
  --shell=[shell]                   Install given shell if required and set it as user's default. Defaults to zsh
  --brew-shell                      Install shell using brew. By default it's installed with system's package manager
  --no-brew                         Don't install brew (Homebrew)
  --prefer-package-manager          Prefer installing tools with system's package manager rather than brew (Doesn't apply for Mac)
  --package-manager=[manager]       Package manager to use for installing prerequisites
-----------------------------------------------------"
DOTFILES_INSTALL_IMPL_USAGE
}

###
# Set default color codes for colorful prints
###
RED_COLOR="\033[0;31m"
GREEN_COLOR="\033[0;32m"
YELLOW_COLOR="\033[1;33m"
BLUE_COLOR="\033[0;34m"
NEUTRAL_COLOR="\033[0m"

function cecho {
    local string_placeholders=""
    for ((i = 1; i < $#; i++)); do
        string_placeholders+="%s"
    done

    # shellcheck disable=SC2059
    printf "${1}${string_placeholders}${NEUTRAL_COLOR}\n" "${@:2}"
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
# Retrieves the path to the given shell's user profile.
# Arguments:
#       $1 - Name of the shell to retrieve for
# Output (stdout):
#       Path to the given shell's user profile
# Returns:
#       0 on success, 1 if an unknown/unsupported shell has been specified
###
function get_shell_user_profile {
    local shell_name="${1:-}"

    case "$shell_name" in
    bash)
        echo "${HOME}/.profile"
        ;;
    zsh)
        echo "${HOME}/.zprofile"
        ;;
    *)
        return 1
        ;;
    esac
}

###
# Checks whether the current user is root.
# Returns:
#       0 if the user is root, 1 otherwise.
###
function root_user {
    local current_uid
    current_uid=$(id -u)

    ((current_uid == 0))
}

function _install_packages_with_brew {
    local packages=("$@")

    install_package_cmd=(brew install --force-bottle "${packages[@]}")

    "${install_package_cmd[@]}"
}

function _install_packages_with_package_manager {
    local packages=("$@")

    if [[ -z "$PACKAGE_MANAGER" ]]; then
        error "Package manager not set, something went wrong. Please install packages manually."
        return 1
    fi

    local install_package_cmd=()
    if [[ "$ROOT_USER" == false ]]; then
        install_package_cmd=(sudo)
    fi
    install_package_cmd+=("$PACKAGE_MANAGER" install -y "${packages[@]}")

    "${install_package_cmd[@]}"
}

###
# Install given package(s) using either system's package manager or homebrew, depending on the passed options.
# Arguments:
#       $1..$N - Variable number of packages to install
# Returns:
#       Install tool's result, zero on success.
###
function install_packages {
    if [[ "$PREFER_BREW_FOR_ALL_TOOLS" == true ]]; then
        ! _install_packages_with_brew "$@" && return 1
    else
        ! _install_packages_with_package_manager "$@" && return 1
    fi
    return 0
}

###
# Reload target shell's user profile, to activate changes.
###
function _reload_shell_user_profile {
    source "$SHELL_USER_PROFILE"
}

function _reinstall_chezmoi_as_package {
    [[ "$INSTALL_BREW" == false || "$BREW_INSTALLED_DOTFILES_MANAGER" == true ]] && return 0

    if ! command -v brew &>/dev/null; then
        warning "Brew is not available, deferring chezmoi installation as a brew package"
        return 0
    fi

    local brew_chezmoi_installed=false

    if brew list | grep -q "$DOTFILES_MANAGER"; then
        brew_chezmoi_installed=true
    fi

    if [[ "$brew_chezmoi_installed" == false ]]; then
        [ "$VERBOSE" == true ] && info "Installing $DOTFILES_MANAGER using brew"

        if ! _install_packages_with_brew "$DOTFILES_MANAGER"; then
            error "Failed installing $DOTFILES_MANAGER using brew, will keep existing binary at $DOTFILES_MANAGER_STANDALONE_BINARY_PATH"
            return 1
        fi
        success "Successfully installed $DOTFILES_MANAGER as a brew package"
    fi

    [ "$VERBOSE" == true ] && info "Removing standalone $DOTFILES_MANAGER binary"
    if ! rm -f "$DOTFILES_MANAGER_STANDALONE_BINARY_PATH"; then
        warning "Failed removing standalone chezmoi binary (downloaded at first) at $DOTFILES_MANAGER_STANDALONE_BINARY_PATH"
    else
        success "Successfully removed standalone chezmoi binary"
    fi

    return 0
}

###
# Finalize installation by executing post-install commands.
###
function post_install {
    if ! _reinstall_chezmoi_as_package; then
        error "Failed reinstalling chezmoi as an updatable package"
        # It's not a fatal error, we can proceed
    fi

    if [[ "$SHELL_TO_INSTALL" == "bash" ]]; then
        if ! _reload_shell_user_profile; then
            warning "Failed reloading shell profile, please attempt a manual re-login"
        fi
    else
        warning "You've installed a new shell, please re-login to apply changes"
    fi

    return 0
}

###
# Apply dotfiles, optionally by using a dotfiles manager.
###
function apply_dotfiles {
    # Always remove old dotfiles, if any, just in case
    rm -rf "$DOTFILES_CLONE_PATH" || return 1

    if [[ "$DEBUG" == true ]]; then
        APPLY_DOTFILES_CMD+=("--verbose")
    fi

    APPLY_DOTFILES_CMD+=(init --apply --ssh "$GITHUB_USERNAME")

    "${APPLY_DOTFILES_CMD[@]}"
}

###
# Prepare dotfiles environment before applying dotfiles.
# This might be a useful step for some dotfiles managers.
###
function prepare_dotfiles_environment {
    if ! mkdir -p "$ENVIRONMENT_TEMPLATE_CONFIG_DIR" &>/dev/null; then
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
        printf "\t%s\n" "full_name = \"$FULL_NAME\""
        printf "\t%s\n" "email = \"$ACTIVE_EMAIL\""
        printf "\t%s\n" "signing_key = \"$ACTIVE_GPG_SIGNING_KEY\""
        printf "\t%s\n" "work_env = $WORK_ENVIRONMENT"
    } >>"$ENVIRONMENT_TEMPLATE_FILE_PATH"

    if [[ "$WORK_ENVIRONMENT" == true ]]; then
        printf "\t%s\n" "work_name = \"$WORK_NAME\"" >>"$ENVIRONMENT_TEMPLATE_FILE_PATH"
    fi

    {
        printf "%s\n" "[data.system]"
        printf "\t%s\n" "shell = \"$SHELL_TO_INSTALL\""
        printf "\t%s\n" "user = \"$CURRENT_USER_NAME\""
    } >>"$ENVIRONMENT_TEMPLATE_FILE_PATH"

    if [[ "$WORK_ENVIRONMENT" == true ]]; then
        {
            printf "\t%s\n" "work_generic_dotfiles_dir = \"${WORK_GENERIC_DOTFILES_DIR}\""
            printf "\t%s\n" "work_specific_dotfiles_dir = \"${WORK_SPECIFIC_DOTFILES_DIR}\""
            printf "\t%s\n" "work_generic_dotfiles_profile = \"${WORK_GENERIC_DOTFILES_PROFILE}\""
            printf "\t%s\n" "work_specific_dotfiles_profile = \"${WORK_SPECIFIC_DOTFILES_PROFILE}\""
        } >>"$ENVIRONMENT_TEMPLATE_FILE_PATH"
    fi

    {
        printf "%s\n" "[data.tools_preferences]"
        printf "\t%s\n" "prefer_brew = $PREFER_BREW_FOR_ALL_TOOLS"
    } >>"$ENVIRONMENT_TEMPLATE_FILE_PATH"

    for mandatory_dir in "${MANDATORY_DIRECTORIES_BEFORE_APPLYING_DOTFILES[@]}"; do
        if ! mkdir -p "${mandatory_dir}"; then
            error "Failed creating directory ${mandatory_dir} - It must exist before applying dotfiles!"
            return 3
        fi
    done

    return 0
}

function _create_new_gpg_key {
    declare -n created_key="${1:?}"

    gpg --expert --full-gen-key || return 1
    created_key="$(gpg --list-secret-keys --keyid-format LONG | tr -s " " | awk -F"[ /]" '/^sec/ { print $3 }' | tail -n1)" || return 2
    return 0
}

function _verify_gpg_client_installation {
    ! command -v gpg &>/dev/null && return 1
    ! command -v gpg-agent &>/dev/null && return 1

    local gpg_version
    gpg_version="$(gpg --version | head -n1 | cut -d' ' -f3)"
    if ! grep -q "^2\.[^0-1]\." <<<"$gpg_version"; then
        # gpg 2.2 or higher is NOT installed
        warning "Installed gpg version ($gpg_version) is less than 2.2"
        return 2
    fi

    return 0
}

function _install_gpg_client {
    local rc
    _verify_gpg_client_installation
    rc=$?

    ((rc == 0)) && return 0

    # Version is too low, nothing we can do for now
    ((rc == 2)) && return 1

    info "Installing gpg"
    sudo apt-get update
    if ! sudo apt-get install -y --no-install-recommends gpg gpg-agent; then
        error "Failed installing gpg tools using apt"
        return 2
    fi
}

###
# Ensures a GPG key exist in order to be able to sign git commits in the future (and maybe do other stuff).
# If a key is not already available, a new one is created instead and will be used in all managed dotfiles.
# Otherwise, the user is asked whether to reuse an existing key, and if so which one.
# The user can also decide to create a new one nevertheless.
# The script requires some interactivity.
###
function ensure_gpg_key_exist {
    info "Installing gpg client (if required)"
    if ! _install_gpg_client; then
        error "Failed installing gpg client"
        return 1
    fi
    success "Successfully installed gpg client"

    if gpg --list-secret-keys --keyid-format LONG | grep -q "sec"; then
        info "GPG keys already available"

        info "Would you like to reuse one of the available keys?"
        local answer
        select answer in "Yes" "No"; do
            case $answer in
            [Yy]*)
                local available_keys
                mapfile -t available_keys < <(gpg --list-secret-keys --keyid-format LONG | tr -s " " | awk -F"[ /]" '/^sec/ { print $3 }')

                info "Select the key to reuse:"
                local selected_key
                select selected_key in "${available_keys[@]}"; do
                    warning "Using $selected_key as the GPG key"
                    ACTIVE_GPG_SIGNING_KEY="$selected_key"
                    return 0
                done
                ;;
            [Nn]*)
                warning "Creating a new GPG key"
                break
                ;;
            esac
        done
    fi

    local new_gpg_key
    if ! _create_new_gpg_key new_gpg_key; then
        error "Failed creating a new GPG key"
        return 2
    fi
    success "Successfully created a new GPG key"

    ACTIVE_GPG_SIGNING_KEY="$new_gpg_key"
    return 0
}

###
# Install selected shell using either system's package manager or homebrew, depending on the passed options.
# If selected shell is already installed, do nothing.
# Otherwise, also configure it as user's default shell.
###
function install_shell {
    if command -v "$SHELL_TO_INSTALL" &>/dev/null; then
        return 0
    fi

    # First install our shell
    if [[ "$INSTALL_SHELL_WITH_BREW" == true ]]; then
        # User has insisted on installing it with brew, so we follow along
        ! _install_packages_with_brew "$SHELL_TO_INSTALL" && return 1
    else
        # Otherwise, we always use the system's package-manager, even if other tools are installed via brew
        ! _install_packages_with_package_manager "$SHELL_TO_INSTALL" && return 2
    fi

    # Find installed shell's location
    local shell_path
    shell_path="$(which "$SHELL_TO_INSTALL")"

    # Then configure it as user's default shell
    sudo chsh -s "$shell_path" "$CURRENT_USER_NAME"
}

###
# Install Homebrew using their official standalone script.
# The script requires some interactivity.
###
function install_brew {
    if [[ "$BREW_AVAILABLE" == true ]]; then
        return 0
    fi

    if ! bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"; then
        return 1
    fi

    # Eval brew for current session to be able to use it later, if needed
    eval "$($BREW_LOCATION_RESOLVING_CMD)"
}

###
# Install chezmoi, our dotfiles manager.
# To avoid any errors and complicated checks, just install the latest binary at this stage.
###
function install_dotfiles_manager {
    local dotfiles_manager_bin=""
    if dotfiles_manager_bin="$(command -v "$DOTFILES_MANAGER" &>/dev/null)" && [[ -n "$dotfiles_manager_bin" ]]; then
        info "$DOTFILES_MANAGER already installed at '$dotfiles_manager_bin', skipping"
    elif [[ "$BREW_AVAILABLE" == true && -e "$DOTFILES_MANAGER_BREW_BINARY_PATH" ]]; then
        info "$DOTFILES_MANAGER already installed with brew, skipping"
        BREW_INSTALLED_DOTFILES_MANAGER=true
        dotfiles_manager_bin="$DOTFILES_MANAGER_BREW_BINARY_PATH/bin/${DOTFILES_MANAGER}"
    elif [[ -f "$DOTFILES_MANAGER_STANDALONE_BINARY_PATH" && -x "$DOTFILES_MANAGER_STANDALONE_BINARY_PATH" ]]; then
        info "$DOTFILES_MANAGER already installed at '$DOTFILES_MANAGER_STANDALONE_BINARY_PATH', skipping"
        dotfiles_manager_bin="$DOTFILES_MANAGER_STANDALONE_BINARY_PATH"
    fi

    if [[ -n "$dotfiles_manager_bin" ]]; then
        APPLY_DOTFILES_CMD=("$dotfiles_manager_bin")
        return 0
    else
        APPLY_DOTFILES_CMD=("$DOTFILES_MANAGER_STANDALONE_BINARY_PATH")
    fi

    local installation_failed=false

    if [[ "$DOWNLOAD_TOOL" == "curl" ]]; then
        ! sh -c "$(curl -fsLS get.chezmoi.io)" && installation_failed=true
    elif [[ "$DOWNLOAD_TOOL" == "wget" ]]; then
        ! sh -c "$(wget -qO- get.chezmoi.io)" && installation_failed=true
    fi

    if [[ "$installation_failed" == true ]]; then
        return 1
    fi
    return 0
}

###
# Install dotfiles. This is the main "driver" function.
###
function install_dotfiles {
    info "Installing dotfiles manager ($DOTFILES_MANAGER)"
    if ! install_dotfiles_manager; then
        error "Failed installing dotfiles manager ($DOTFILES_MANAGER)"
        return 1
    fi
    success "Successfully installed dotfiles manager, $DOTFILES_MANAGER"

    if [[ "$INSTALL_BREW" == true ]]; then
        info "Installing brew"
        if ! install_brew; then
            error "Failed installing brew"
            return 2
        fi
        success "Successfully installed brew"
    fi

    info "Installing shell"
    if ! install_shell; then
        error "Failed installing shell"
        return 2
    fi
    success "Successfully installed $SHELL_TO_INSTALL"

    info "Ensuring a GPG key exists"
    if ! ensure_gpg_key_exist; then
        error "Failed ensuring a GPG key exists"
        return 3
    fi
    success "Successfully ensured a GPG key exists"

    info "Preparing dotfiles environment"
    if ! prepare_dotfiles_environment; then
        error "Failed preparing dotfiles environment"
        return 4
    fi
    success "Successfully prepared dotfiles environment"

    info "Applying dotfiles"
    if ! apply_dotfiles; then
        error "Failed applying dotfiles"
        return 5
    fi
    success "Successfully applied dotfiles"

    info "Finalizing installation"
    if ! post_install; then
        error "Failed finalizing installation"
        return 6
    fi
    success "Successfully finalized installation"

    return 0
}

function brew_available {
    [[ -d /home/linuxbrew/ && -f "$DEFAULT_BREW_PATH" ]]
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
        if command -v "${download_tool}" 2>/dev/null; then
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
    if [[ -n "$INSTALL_REF" ]]; then
        APPLY_DOTFILES_CMD+=(--branch "$INSTALL_REF")
    fi

    # Can't prefer to install with brew if brew should not even be installed
    if [[ "$INSTALL_BREW" == false ]]; then
        PREFER_BREW_FOR_ALL_TOOLS=false
    else
        brew_available && BREW_AVAILABLE=true
    fi

    if ! DOWNLOAD_TOOL="$(get_download_tool)"; then
        error "Couldn't determine download tool, aborting"
        return 1
    fi

    CURRENT_USER_NAME="$(id -u -n)"

    if root_user; then
        ROOT_USER=true
    fi

    if ! SHELL_USER_PROFILE="$(get_shell_user_profile "$SHELL_TO_INSTALL")"; then
        error "Failed determining shell's user profile"
        return 2
    fi

    if [[ "$WORK_ENVIRONMENT" == true ]]; then
        ACTIVE_EMAIL="$WORK_EMAIL"
        WORK_SPECIFIC_DOTFILES_DIR="${WORK_GENERIC_DOTFILES_DIR}/${WORK_NAME}"
        WORK_SPECIFIC_DOTFILES_PROFILE="${WORK_SPECIFIC_DOTFILES_DIR}/profile"
    else
        ACTIVE_EMAIL="$PERSONAL_EMAIL"
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

    local short_options=hvd
    local long_options=help,verbose,debug
    long_options+=,ref:
    long_options+=,work-env,work-name:,work-email:
    long_options+=,shell:,brew-shell
    long_options+=,no-brew,prefer-package-manager,package-manager:

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
        -d | --debug)
            DEBUG=true
            shift
            ;;
        --ref)
            INSTALL_REF="${2:-main}"
            shift 2
            ;;
        --work-env)
            WORK_ENVIRONMENT=true
            shift
            ;;
        --work-name)
            WORK_ENVIRONMENT=true
            [ -n "$2" ] && WORK_NAME="${2}"
            shift 2
            ;;
        --work-email)
            [ -n "$2" ] && WORK_EMAIL="${2}"
            WORK_ENVIRONMENT=true
            shift 2
            ;;
        --shell)
            SHELL_TO_INSTALL="${2:-}"
            shift 2
            ;;
        --brew-shell)
            INSTALL_SHELL_WITH_BREW=true
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
        --package-manager)
            PACKAGE_MANAGER="${2:-}"
            shift 2
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

function _set_work_info_defaults {
    WORK_NAME="sedg"
    WORK_GENERIC_DOTFILES_DIR="${HOME}/.work"
    WORK_GENERIC_DOTFILES_PROFILE="${WORK_GENERIC_DOTFILES_DIR}/profile"
}

function _set_package_management_defaults {
    PACKAGE_MANAGER=""
    INSTALL_BREW=true
    PREFER_BREW_FOR_ALL_TOOLS=true
    DEFAULT_BREW_PATH="/home/linuxbrew/.linuxbrew/bin/brew"
    BREW_LOCATION_RESOLVING_CMD="$DEFAULT_BREW_PATH shellenv"
    BREW_AVAILABLE=false
    BREW_INSTALLED_DOTFILES_MANAGER=false
}

function _set_shell_defaults {
    INSTALL_SHELL_WITH_BREW=false
    SHELL_TO_INSTALL=zsh
    SHELL_USER_PROFILE=""
}

function _set_dotfiles_manager_defaults {
    DOTFILES_MANAGER=chezmoi
    DOTFILES_MANAGER_STANDALONE_BINARY_PATH="${HOME}/bin/${DOTFILES_MANAGER}"
    DOTFILES_MANAGER_BREW_BINARY_PATH="/home/linuxbrew/.linuxbrew/opt/${DOTFILES_MANAGER}"

    DOTFILES_CLONE_PATH="${HOME}/.local/share/${DOTFILES_MANAGER}"
    ENVIRONMENT_TEMPLATE_CONFIG_DIR="$HOME/.config/${DOTFILES_MANAGER}"
    ENVIRONMENT_TEMPLATE_FILE_PATH="${ENVIRONMENT_TEMPLATE_CONFIG_DIR}/${DOTFILES_MANAGER}.toml"

    MANDATORY_DIRECTORIES_BEFORE_APPLYING_DOTFILES=(
        "${HOME}/.oh-my-zsh/cache/"
    )
}

function _set_personal_info_defaults {
    GITHUB_USERNAME="MrPointer"
    FULL_NAME="Timor Gruber"
    PERSONAL_EMAIL="timor.gruber@gmail.com"
    WORK_EMAIL="timor.gruber@solaredge.com"
}

###
# Set script default values for later show_usage.
###
function set_defaults {
    VERBOSE=false
    DEBUG=false
    INSTALL_REF=main
    WORK_ENVIRONMENT=false
    ROOT_USER=false

    _set_personal_info_defaults
    _set_dotfiles_manager_defaults
    _set_shell_defaults
    _set_package_management_defaults
    _set_work_info_defaults
}

###
# This is the script's entry point, just like in any other programming language.
###
function main {
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

    info "Installing dotfiles"
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
