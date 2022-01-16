#!/usr/bin/env sh

###
# Set default color codes for colorful prints
###
RED_COLOR="\033[0;31m"
GREEN_COLOR="\033[0;32m"
YELLOW_COLOR="\033[1;33m"
BLUE_COLOR="\033[0;34m"
NEUTRAL_COLOR="\033[0m"

error() {
    printf "${RED_COLOR}%s${NEUTRAL_COLOR}\n" "$@"
}

warning() {
    printf "${YELLOW_COLOR}%s${NEUTRAL_COLOR}\n" "$@"
}

info() {
    printf "${BLUE_COLOR}%s${NEUTRAL_COLOR}\n" "$@"
}

success() {
    printf "${GREEN_COLOR}%s${NEUTRAL_COLOR}\n" "$@"
}

get_download_tool() {
    if hash curl 2>/dev/null; then
        echo "curl"
    elif hash wget 2>/dev/null; then
        echo "wget"
    else
        echo ""
    fi
}

invoke_actual_installation() {
    # Create temporary executable file to hold the contents
    # of the downloaded implementation script
    TMP_IMPL_INSTALL_PATH="$(mktemp)"
    chmod +x "$TMP_IMPL_INSTALL_PATH"

    # Execute manually for every type of download tool to get exit code, it's impossible otherwise...
    # Shell commands executed with "-c" must be in single-quotes to catch their exit codes correctly
    IMPL_DOWNLOAD_RESULT=0
    case "$DOWNLOAD_TOOL" in
    curl)
        curl -fsSL https://raw.githubusercontent.com/MrPointer/dotfiles/"$INSTALL_BRANCH"/install-impl.sh -o "$TMP_IMPL_INSTALL_PATH"
        IMPL_DOWNLOAD_RESULT=$?
        ;;
    wget)
        wget -q https://raw.githubusercontent.com/MrPointer/dotfiles/"$INSTALL_BRANCH"/install-impl.sh -O "$TMP_IMPL_INSTALL_PATH"
        IMPL_DOWNLOAD_RESULT=$?
        ;;
    esac

    if [ $IMPL_DOWNLOAD_RESULT -ne 0 ]; then
        error "Failed downloading implementation script!"
        return 2
    fi

    if ! "$TMP_IMPL_INSTALL_PATH" "--package-manager" "$PKG_MANAGER" "$@"; then
        error "Failed on actual installation of dotfiles, sorry..."
        return 3
    fi
}

get_default_system_package_manager() {
    case "$1" in
    mac | darwin)
        echo "brew"
        ;;
    ubuntu | debian | suse)
        echo "apt"
        ;;
    fedora | centos | redhat)
        echo "dnf"
        ;;
    *)
        echo ""
        ;;
    esac
}

get_linux_distro_name() {
    distro=""

    if [ -f /etc/os-release ]; then
        # freedesktop.org and systemd
        . /etc/os-release
        distro=$NAME
    elif [ -f /etc/lsb-release ]; then
        # For some versions of Debian/Ubuntu without lsb_release command
        . /etc/lsb-release
        distro=$DISTRIB_ID
    elif [ -f /etc/debian_version ]; then
        # Older Debian/Ubuntu/etc.
        distro=Debian
    elif [ -f /etc/SuSe-release ]; then
        # Older SuSE/etc.
        distro=SuSE
    elif [ -f /etc/redhat-release ]; then
        # Older Red Hat, CentOS, etc.
        distro=RedHat
    else
        # Fall back to uname, e.g. "Linux <version>", also works for BSD, etc.
        distro="$(uname -s)"
    fi

    echo "$distro" | tr '[:upper:]' '[:lower:]'
}

get_system_type() {
    case "$(uname -s)" in
    Darwin)
        echo "mac"
        ;;
    Linux)
        echo "linux"
        ;;
    *)
        echo "unsupported"
        ;;
    esac
}

install_bash_with_package_manager() {
    case "$1" in
    apt)
        sudo apt install -y bash
        ;;
    dnf)
        sudo dnf install -y bash
        ;;
    *) ;;

    esac
}

bash_exists() {
    hash bash >/dev/null 2>&1
}

install_bash() {
    info "Checking if bash exists"
    if bash_exists; then
        info "Bash exists!"
        return 0
    fi

    info "bash does not exist, trying to install it"
    if ! install_bash_with_package_manager "$1"; then
        error "Failed installing bash using $1"
        return 5
    fi

    return 0
}

###
# Parse arguments/options using getopt, the almighty C-based parser.
###
parse_arguments() {
    getopt --test >/dev/null
    if [ $? -ne 4 ]; then
        error "I'm sorry, 'getopt --test' failed in this environment."
        return 1
    fi

    long_options=branch:

    # -temporarily store output to be able to check for errors
    # -activate quoting/enhanced mode (e.g. by writing out “--options”)
    # -pass arguments only via   -- "$@"   to separate them correctly
    if ! PARSED=$(
        getopt --longoptions="$long_options" \
            --name "Dotfiles installer-bootstrapper" -- "$@"
    ); then
        # getopt has complained about wrong arguments to stdout
        error "Wrong arguments to Dotfiles installer-bootstrapper" && return 2
    fi

    # read getopt’s output this way to handle the quoting right:
    eval set -- "$PARSED"

    while true; do
        case $1 in
        --branch)
            INSTALL_BRANCH="${2:-main}"
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

    unset long_options
    return 0
}

detect_system() {
    SYSTEM_TYPE="$(get_system_type)"

    case "$SYSTEM_TYPE" in
    linux)
        DISTRO_NAME="$(get_linux_distro_name)"
        if [ -z "$DISTRO_NAME" ]; then
            echo "Failed detecting linux distribution, or distro not supported: $DISTRO_NAME"
            return 1
        fi
        ;;
    mac)
        DISTRO_NAME="mac"
        ;;
    *)
        echo "Unsupported system type: $SYSTEM_TYPE"
        return 1
        ;;
    esac

    PKG_MANAGER="$(get_default_system_package_manager "$DISTRO_NAME")"
    if [ -z "$PKG_MANAGER" ]; then
        error "Failed determining package manager for distro: $DISTRO_NAME"
        return 3
    fi
    if ! hash "$PKG_MANAGER" 2>/dev/null; then
        error "Package manager '$PKG_MANAGER' couldn't be found for distro: $DISTRO_NAME, maybe you need to install it manually?"
        return 4
    fi

    info "Determined system:"
    info "Type: $SYSTEM_TYPE"
    info "Distro: $DISTRO_NAME"
    info "Package manager: $PKG_MANAGER"
    printf "\n" # Print an empty line
}

set_defaults() {
    INSTALL_BRANCH="main"
}

main() {
    info "Installing dotfiles, but first some bootstrapping"

    set_defaults # Should never fail

    if ! detect_system; then
        error "Detected system is not supported, sorry"
        return 1
    fi

    if ! parse_arguments "$@"; then
        error "Failed parsing arguments, aborting"
        return 2
    fi

    if ! install_bash "$PKG_MANAGER"; then
        error "Failed installing bash!"
        return 3
    fi

    DOWNLOAD_TOOL="$(get_download_tool)"
    if [ -z "$DOWNLOAD_TOOL" ]; then
        error "Neither 'curl' nor 'wget' are available, please install one of them manually."
        return 4
    fi

    if ! invoke_actual_installation "$@"; then
        error "Failed to install dotfiles"
        return 5
    fi

    success "Successfully completed dotfiles installation [from bootstrap]"
    return 0
}

main "$@"
