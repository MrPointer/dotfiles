#!/bin/sh

get_download_tool() {
    if hash curl 2>/dev/null; then
        echo "curl"
    elif hash wget 2>/dev/null; then
        echo "wget"
    else
        echo ""
    fi
}

get_default_system_package_manager() {
    case "$1" in
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
    distro

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

    echo "$distro"
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
    "apt")
        sudo apt install -y bash
        ;;
    "dnf")
        sudo dnf install -y bash
        ;;
    *) ;;

    esac
}

install_bash() {
    case "$(get_system_type)" in
    "linux")
        DISTRO_NAME="$(get_linux_distro_name | tr '[:upper:]' '[:lower:]')"
        PKG_MANAGER="$(get_default_system_package_manager "$DISTRO_NAME")"
        if [ -z "$PKG_MANAGER" ]; then
            echo "Package manager couldn't be identified for distro: $DISTRO_NAME"
            return 3
        fi
        if ! hash "$PKG_MANAGER" 2>/dev/null; then
            echo "Package manager '$PKG_MANAGER' couldn't be found for distro: $DISTRO_NAME"
            return 4
        fi
        if ! install_bash_with_package_manager "$PKG_MANAGER"; then
            echo "Failed installing bash using $PKG_MANAGER"
            return 5
        fi
        ;;
    "mac")
        echo "Bash should already be installed on Mac..."
        return 2
        ;;
    "unsupported")
        echo "Unsupported system, sorry"
        return 1
        ;;
    *)
        echo "WTF?"
        return 2
        ;;
    esac

    return 0
}

bash_exists() {
    hash bash >/dev/null 2>&1
}

main() {
    if ! bash_exists; then
        echo "bash does not exist, trying to install it"
        if ! install_bash; then
            echo "Failed installing bash!"
            return 1
        fi
    fi

    DOWNLOAD_TOOL="$(get_download_tool)"
    if [ -z "$DOWNLOAD_TOOL" ]; then
        echo "Neither curl nor wget are available, please install one of them manually."
        return 2
    fi

    TMP_IMPL_INSTALL_PATH="$(mktemp)"

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

    if ! "$TMP_IMPL_INSTALL_PATH" "$@"; then
        echo "Failed on actual installation of dotfiles, sorry..."
        return 3
    fi

    return 0
}

main "$@"
