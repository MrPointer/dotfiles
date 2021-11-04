#!/usr/bin/env sh

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
        curl -fsSL https://raw.githubusercontent.com/MrPointer/dotfiles/main/install_impl.sh -o "$TMP_IMPL_INSTALL_PATH"
        IMPL_DOWNLOAD_RESULT=$?
        ;;
    wget)
        wget -q https://raw.githubusercontent.com/MrPointer/dotfiles/main/install_impl.sh -O "$TMP_IMPL_INSTALL_PATH"
        IMPL_DOWNLOAD_RESULT=$?
        ;;
    esac

    if [ $IMPL_DOWNLOAD_RESULT -ne 0 ]; then
        echo "Failed downloading implementation script!"
        return 2
    fi

    if ! "$TMP_IMPL_INSTALL_PATH" "--package-manager" "$PKG_MANAGER" "$@"; then
        echo "Failed on actual installation of dotfiles, sorry..."
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
    printf "%s" "Checking if bash exists... "
    if bash_exists; then
        echo "OK"
        return 0
    fi

    echo "bash does not exist, trying to install it"
    if ! install_bash_with_package_manager "$1"; then
        echo "Failed installing bash using $1"
        return 5
    fi

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
        echo "Failed determining package manager for distro: $DISTRO_NAME"
        return 3
    fi
    if ! hash "$PKG_MANAGER" 2>/dev/null; then
        echo "Package manager '$PKG_MANAGER' couldn't be found for distro: $DISTRO_NAME, maybe you need to install it manually?"
        return 4
    fi

    echo "Determined system:"
    echo "Type: $SYSTEM_TYPE"
    echo "Distro: $DISTRO_NAME"
    printf "%s\n\n" "Package manager: $PKG_MANAGER"
}

main() {
    echo "Installing dotfiles, but first some bootstraping"

    if ! detect_system; then
        echo "Detected system is not supported, sorry"
        return 1
    fi

    if ! install_bash "$PKG_MANAGER"; then
        echo "Failed installing bash!"
        return 1
    fi

    DOWNLOAD_TOOL="$(get_download_tool)"
    if [ -z "$DOWNLOAD_TOOL" ]; then
        echo "Neither 'curl' nor 'wget' are available, please install one of them manually."
        return 2
    fi

    if ! invoke_actual_installation "$@"; then
        echo "Failed to install dotfiles"
        return 3
    fi

    echo "Successfully completed dotfiles installation [from bootstrap]"
    return 0
}

main "$@"
