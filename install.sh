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
    if command -v curl >/dev/null 2>&1; then
        echo "curl"
    elif command -v wget >/dev/null 2>&1; then
        echo "wget"
    else
        echo ""
    fi
}

invoke_actual_installation() {
    if [ "$INVOKE_LOCAL_INSTALL" = true ]; then
        # Get the path to the local implementation script
        v_script_base_dir="$(dirname "$(readlink -f "$0")")"
        v_local_install_script="${v_script_base_dir}/install-impl.sh"

        # Check if the file exists
        if [ ! -f "$v_local_install_script" ]; then
            error "Failed to find local implementation script!"
            unset v_local_install_script v_script_base_dir
            return 1
        fi
        # Check if the file is executable
        if [ ! -x "$v_local_install_script" ]; then
            error "Local implementation script is not executable!"
            unset v_local_install_script v_script_base_dir
            return 2
        fi

        TMP_IMPL_INSTALL_PATH="$v_local_install_script"
        unset v_local_install_script v_script_base_dir
    else
        # Create temporary executable file to hold the contents
        # of the downloaded implementation script
        TMP_IMPL_INSTALL_PATH="$(mktemp)"

        # Execute manually for every type of download tool to get exit code, it's impossible otherwise...
        # Shell commands executed with "-c" must be in single-quotes to catch their exit codes correctly
        IMPL_DOWNLOAD_RESULT=0
        case "$v_DOWNLOAD_TOOL" in
        curl)
            curl -fsSL -o "$TMP_IMPL_INSTALL_PATH" "https://raw.githubusercontent.com/MrPointer/dotfiles/$INSTALL_REF/install-impl.sh"
            IMPL_DOWNLOAD_RESULT=$?
            ;;
        wget)
            wget -q -O "$TMP_IMPL_INSTALL_PATH" "https://raw.githubusercontent.com/MrPointer/dotfiles/$INSTALL_REF/install-impl.sh"
            IMPL_DOWNLOAD_RESULT=$?
            ;;
        esac

        if [ $IMPL_DOWNLOAD_RESULT -ne 0 ] || [ ! -s "$TMP_IMPL_INSTALL_PATH" ]; then
            error "Failed downloading implementation script!"
            return 2
        fi

        chmod +x "$TMP_IMPL_INSTALL_PATH"
    fi

    # For macOS, find GNU getopt path
    if [ "$SYSTEM_TYPE" = "darwin" ]; then
        # Determine where GNU getopt is installed
        v_apple_silicon_path="/opt/homebrew/opt/gnu-getopt/bin"
        v_intel_path="/usr/local/opt/gnu-getopt/bin"

        if [ -d "$v_apple_silicon_path" ]; then
            # Apple Silicon Mac
            v_getopt_path="$v_apple_silicon_path"
        elif [ -d "$v_intel_path" ]; then
            # Intel Mac
            v_getopt_path="$v_intel_path"
        else
            error "GNU getopt not found on macOS, please install it manually OR open a new shell and run the script again"
            unset v_apple_silicon_path v_intel_path
            return 4
        fi

        # Execute with modified PATH environment
        info "Executing: env PATH=\"$v_getopt_path:\$PATH\" $TMP_IMPL_INSTALL_PATH --package-manager $PKG_MANAGER --system $SYSTEM_TYPE $*"
        env PATH="$v_getopt_path:$PATH" "$TMP_IMPL_INSTALL_PATH" --package-manager "$PKG_MANAGER" --system "$SYSTEM_TYPE" "$@"
        v_result=$?
    else
        # Execute normally on other systems
        info "Executing: $TMP_IMPL_INSTALL_PATH --package-manager $PKG_MANAGER --system $SYSTEM_TYPE $*"
        "$TMP_IMPL_INSTALL_PATH" --package-manager "$PKG_MANAGER" --system "$SYSTEM_TYPE" "$@"
        v_result=$?
    fi

    if [ $v_result -ne 0 ]; then
        error "Real installer failed, sorry..."
        return 3
    fi

    unset v_local_install_script
    return 0
}

install_getopt() {
    v_distro="$1"
    v_pkg_manager="$2"

    if [ "$v_distro" = "mac" ] && [ "$v_pkg_manager" = "brew" ]; then
        if brew list | grep -q gnu-getopt; then
            info "gnu-getopt already installed"
            return 0
        fi

        brew install gnu-getopt || return 1
        unset v_distro v_pkg_manager
        return 0
    fi

    unset v_distro v_pkg_manager
    return 0
}

verify_bash_version() {
    v_bash_version="$(bash --version | head -n 1 | grep -o 'version [0-9]\.[0-9]\.[0-9]' | cut -d ' ' -f 2)"
    if [ -z "$v_bash_version" ]; then
        error "Failed detecting bash version"
        unset v_bash_version
        return 2
    fi

    # Check if bash version is at least 4.4.0
    v_expected_bash_version="4.4.0"
    if [ "$(printf '%s\n' "$v_bash_version" "$v_expected_bash_version" | sort -V | head -n1)" = "$v_expected_bash_version" ]; then
        info "Bash exists!"
        unset v_bash_version
        return 0
    fi

    unset v_bash_version
    return 1
}

install_bash_with_package_manager() {
    case "$1" in
    apt)
        sudo apt install -y bash
        ;;
    dnf)
        sudo dnf install -y bash
        ;;
    brew)
        brew install bash
        ;;
    *) ;;

    esac
}

install_bash() {
    if command -v bash >/dev/null 2>&1; then
        verify_bash_version
        v_exit_code=$?
        if [ $v_exit_code -eq 0 ]; then
            return 0
        fi
        if [ $v_exit_code -eq 1 ]; then
            warning "Bash version is too old, trying to install a newer version using detected package manager"
        fi
        if [ $v_exit_code -eq 2 ]; then
            return 1
        fi
        unset v_exit_code
    else
        warning "bash does not exist, trying to install it"
    fi

    if ! install_bash_with_package_manager "$1"; then
        error "Failed installing bash using $1"
        return 5
    fi

    # Verify version again to ensure it's installed correctly
    verify_bash_version
    v_exit_code=$?
    if [ $v_exit_code -eq 1 ]; then
        warning "Bash version is still too old, please install a newer version manually OR open a new shell and run the script again"
        unset v_exit_code
        return 6
    fi

    unset v_exit_code
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
        --local)
            INVOKE_LOCAL_INSTALL=true
            shift
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
    echo "$SUPPORTED_DISTROS" >"$v_supported_distros_file"

    if ! grep -q "$v_distro" "$v_supported_distros_file"; then
        error "$v_distro is not yet supported, currently supported are: $SUPPORTED_DISTROS"
        unset v_supported_distros_file v_system v_distro v_pkg_manager
        return 3
    fi

    unset v_supported_distros_file v_system v_distro v_pkg_manager
    return 0
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
    v_distro=""

    if [ -f /etc/os-release ]; then
        # freedesktop.org and systemd
        . /etc/os-release
        v_distro=$NAME
    elif [ -f /etc/lsb-release ]; then
        # For some versions of Debian/Ubuntu without lsb_release command
        . /etc/lsb-release
        v_distro=$DISTRIB_ID
    elif [ -f /etc/debian_version ]; then
        # Older Debian/Ubuntu/etc.
        v_distro=Debian
    elif [ -f /etc/SuSe-release ]; then
        # Older SuSE/etc.
        v_distro=SuSE
    elif [ -f /etc/redhat-release ]; then
        # Older Red Hat, CentOS, etc.
        v_distro=RedHat
    else
        # Fall back to uname, e.g. "Linux <version>", also works for BSD, etc.
        v_distro="$(uname -s)"
    fi

    echo "$v_distro" | tr '[:upper:]' '[:lower:]'
    unset v_distro
    return 0
}

_get_system_type() {
    case "$(uname -s)" in
    Darwin)
        echo "darwin"
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
    darwin)
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
    if ! command -v "$PKG_MANAGER" >/dev/null 2>&1; then
        error "Detected '$PKG_MANAGER' as package-manager for '$DISTRO_NAME' but it's not available, maybe you need to install it manually first?"
        return 4
    fi

    printf "\n" # Print an empty line
    info "Detected system:"
    info "----------------"
    info "Type: $SYSTEM_TYPE"
    info "Distro: $DISTRO_NAME"
    info "Package manager: $PKG_MANAGER"
    info "----------------"
    printf "\n" # Print an empty line
}

set_defaults() {
    INSTALL_REF="main"
    INVOKE_LOCAL_INSTALL=false
    SUPPORTED_DISTROS="ubuntu debian mac"
}

main() {
    info "Installing dotfiles, but first some bootstrapping"

    set_defaults # Should never fail

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

    info "Installing getopt (if required)"
    if ! install_getopt "$DISTRO_NAME" "$PKG_MANAGER"; then
        error "Failed installing getopt!"
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
