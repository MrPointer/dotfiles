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
    case "$v_DOWNLOAD_TOOL" in
    curl)
        curl -fsSL https://raw.githubusercontent.com/MrPointer/dotfiles/"$INSTALL_REF"/install-impl.sh -o "$TMP_IMPL_INSTALL_PATH"
        IMPL_DOWNLOAD_RESULT=$?
        ;;
    wget)
        wget -q https://raw.githubusercontent.com/MrPointer/dotfiles/"$INSTALL_REF"/install-impl.sh -O "$TMP_IMPL_INSTALL_PATH"
        IMPL_DOWNLOAD_RESULT=$?
        ;;
    esac

    if [ $IMPL_DOWNLOAD_RESULT -ne 0 ]; then
        error "Failed downloading implementation script!"
        return 2
    fi

    if ! "$TMP_IMPL_INSTALL_PATH" "--package-manager" "$PKG_MANAGER" "$@"; then
        error "Real installer failed, sorry..."
        return 3
    fi

    return 0
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
    if bash_exists; then
        info "Bash exists!"
        return 0
    fi

    warning "bash does not exist, trying to install it"
    if ! install_bash_with_package_manager "$1"; then
        error "Failed installing bash using $1"
        return 5
    fi

    return 0
}

###
# Parse special arguments/options, and "swallow" the rest as they're intended for the implementation script.
###
parse_arguments() {
    while [ "$#" -gt 0 ]; do
        case $1 in
        --ref)
            [ -n "$2" ] && INSTALL_REF="${2}"
            shift 2
            ;;
        *)
            # Probably options to the real installer (implementation), simply shift past them
            shift
            ;;
        esac
    done
}

supported_system() {
    v_system="${1:?}"
    v_distro="${2:?}"
    v_pkg_manager="${3:?}"

    v_supported_distros_file="$(mktemp)" || return 1
    echo "$SUPPORTED_LINUX_DISTROS" >"$v_supported_distros_file"

    if ! grep -q "$v_distro" "$v_supported_distros_file"; then
        error "$v_distro is not yet supported, currently supported are: $SUPPORTED_LINUX_DISTROS"
        unset v_supported_distros_file
        return 3
    fi
    unset v_supported_distros_file

    unset v_system v_distro v_pkg_manager
}

_get_default_system_package_manager() {
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

_get_linux_distro_name() {
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

_get_system_type() {
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

detect_system() {
    SYSTEM_TYPE="$(_get_system_type)"

    case "$SYSTEM_TYPE" in
    linux)
        DISTRO_NAME="$(_get_linux_distro_name)"
        if [ -z "$DISTRO_NAME" ]; then
            error "Failed detecting linux distribution"
            return 2
        fi
        ;;
    mac)
        DISTRO_NAME="mac"
        ;;
    *)
        error "Unsupported system type: $SYSTEM_TYPE"
        return 3
        ;;
    esac

    PKG_MANAGER="$(_get_default_system_package_manager "$DISTRO_NAME")"
    if [ -z "$PKG_MANAGER" ]; then
        error "Failed determining package manager for distro: $DISTRO_NAME"
        return 2
    fi
    if ! hash "$PKG_MANAGER" 2>/dev/null; then
        error "Package manager '$PKG_MANAGER' couldn't be found for distro: $DISTRO_NAME, maybe you need to install it manually?"
        return 4
    fi

    printf "\n" # Print an empty line
    info "Detected system:"
    info "----------------"
    info "Type: $SYSTEM_TYPE"
    info "Distro: $DISTRO_NAME"
    info "Package manager: $PKG_MANAGER"
    printf "\n" # Print an empty line
}

set_defaults() {
    INSTALL_REF="main"
    SUPPORTED_LINUX_DISTROS="ubuntu debian"
}

main() {
    info "Installing dotfiles, but first some bootstrapping"

    set_defaults # Should never fail

    info "Detecting system"
    if ! detect_system; then
        error "Detected system is not supported, sorry"
        return 1
    fi

    if ! supported_system "$SYSTEM_TYPE" "$DISTRO_NAME" "$PKG_MANAGER"; then
        error "Detected system is not supported, sorry"
        return 1
    fi

    if ! parse_arguments "$@"; then
        error "Failed parsing arguments, aborting"
        return 2
    fi

    info "Installing bash (if required)"
    if ! install_bash "$PKG_MANAGER"; then
        error "Failed installing bash!"
        return 3
    fi

    v_DOWNLOAD_TOOL="$(get_download_tool)"
    if [ -z "$v_DOWNLOAD_TOOL" ]; then
        error "Neither 'curl' nor 'wget' are available, please install one of them manually"
        return 4
    fi

    info "Running real bootstrap installation (bash script)"
    if ! invoke_actual_installation "$@"; then
        error "Failed installing dotfiles [from bootstrap]"
        return 5
    fi

    success "Successfully completed dotfiles installation [from bootstrap]"
    return 0
}

main "$@"
