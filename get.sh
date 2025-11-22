#!/bin/bash

# Dotfiles Installer - Get Script
# This script downloads the dotfiles installer binary and optionally runs it

set -e

# Default values
INSTALL_DIR="$HOME/.local/bin"
REPO="MrPointer/dotfiles"
BINARY_NAME="dotfiles-installer"
RUN_INSTALLER=false
INSTALLER_ARGS=""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
log() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Get system information
get_os() {
    case "$(uname -s)" in
        Darwin*) echo "darwin" ;;
        Linux*) echo "linux" ;;
        *) error "Unsupported operating system: $(uname -s)" ;;
    esac
}

get_arch() {
    case "$(uname -m)" in
        x86_64|amd64)
            if [ "$(get_os)" = "darwin" ]; then
                error "Intel Macs are no longer supported. This installer only supports Apple Silicon Macs (M1 and above)."
            fi
            echo "x86_64" ;;
        arm64|aarch64) echo "arm64" ;;
        *) error "Unsupported architecture: $(uname -m)" ;;
    esac
}

# Get latest release version from GitHub
get_latest_version() {
    if command_exists curl; then
        curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
    elif command_exists wget; then
        wget -qO- "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
    else
        error "Neither curl nor wget found. Please install one of them."
    fi
}

# Download file
download_file() {
    local url="$1"
    local output="$2"

    if command_exists curl; then
        curl -L -o "$output" "$url"
    elif command_exists wget; then
        wget -O "$output" "$url"
    else
        error "Neither curl nor wget found. Please install one of them."
    fi
}

# Show help
show_help() {
    cat << EOF
Dotfiles Installer - Get Script

USAGE:
    $0 [OPTIONS]

OPTIONS:
    -d, --dir DIR       Download directory (default: $HOME/.local/bin)
    -v, --version VER   Specific version to download (default: latest)
    -r, --run           Run the installer after download
    -h, --help          Show this help message

EXAMPLES:
    # Download latest version to default location
    $0

    # Download to custom directory
    $0 --dir /usr/local/bin

    # Download specific version
    $0 --version v1.2.3

    # Download and run immediately
    $0 --run

    # Download and run with custom options
    $0 --run -- --work-env --install-prerequisites

    # Download and run the install command non-interactively
    $0 --run -- install --non-interactive --plain

EOF
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -d|--dir)
            INSTALL_DIR="$2"
            shift 2
            ;;
        -v|--version)
            VERSION="$2"
            shift 2
            ;;
        -r|--run)
            RUN_INSTALLER=true
            shift 1
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        --)
            # Everything after -- goes to the installer
            shift
            INSTALLER_ARGS="$*"
            break
            ;;
        *)
            error "Unknown option: $1. Use --help for usage information."
            ;;
    esac
done

main() {
    log "Getting dotfiles installer..."

    # Check prerequisites
    if ! command_exists curl && ! command_exists wget; then
        error "Neither curl nor wget found. Please install one of them."
    fi

    if ! command_exists tar; then
        error "tar command not found. Please install tar."
    fi

    # Get system information
    OS=$(get_os)
    ARCH=$(get_arch)
    log "Detected system: $OS/$ARCH"

    # Get version
    if [ -z "$VERSION" ]; then
        log "Getting latest release version..."
        VERSION=$(get_latest_version)
        if [ -z "$VERSION" ]; then
            error "Failed to get latest version information"
        fi
    fi

    log "Target version: $VERSION"

    # Construct download URL
    FILENAME="${BINARY_NAME}-${VERSION#v}-${OS}-${ARCH}.tar.gz"
    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${FILENAME}"

    log "Download URL: $DOWNLOAD_URL"

    # Create install directory
    if [ ! -d "$INSTALL_DIR" ]; then
        log "Creating download directory: $INSTALL_DIR"
        mkdir -p "$INSTALL_DIR"
    fi

    # Create temporary directory
    TMP_DIR=$(mktemp -d)
    trap "rm -rf $TMP_DIR" EXIT

    log "Downloading $FILENAME..."
    download_file "$DOWNLOAD_URL" "$TMP_DIR/$FILENAME"

    if [ ! -f "$TMP_DIR/$FILENAME" ]; then
        error "Download failed: file not found"
    fi

    log "Extracting archive..."
    tar -xzf "$TMP_DIR/$FILENAME" -C "$TMP_DIR"

    # Find the binary (it should be in the extracted directory)
    BINARY_PATH="$TMP_DIR/$BINARY_NAME"
    if [ ! -f "$BINARY_PATH" ]; then
        # Sometimes it's in a subdirectory
        BINARY_PATH=$(find "$TMP_DIR" -name "$BINARY_NAME" -type f | head -n1)
        if [ -z "$BINARY_PATH" ]; then
            error "Binary not found in downloaded archive"
        fi
    fi

    log "Setting up binary at $INSTALL_DIR/$BINARY_NAME..."
    cp "$BINARY_PATH" "$INSTALL_DIR/$BINARY_NAME"
    chmod +x "$INSTALL_DIR/$BINARY_NAME"

    # Check if install directory is in PATH
    case ":$PATH:" in
        *":$INSTALL_DIR:"*)
            PATH_IN_PATH=true
            ;;
        *)
            PATH_IN_PATH=false
            ;;
    esac

    success "Download completed!"
    echo
    log "Binary installed at: $INSTALL_DIR/$BINARY_NAME"

    if [ "$PATH_IN_PATH" = true ]; then
        log "You can now run: $BINARY_NAME"
    else
        warn "Download directory is not in PATH. You can:"
        echo "  1. Run directly: $INSTALL_DIR/$BINARY_NAME"
        echo "  2. Add to PATH: export PATH=\"$INSTALL_DIR:\$PATH\""
        echo "  3. Create a symlink in /usr/local/bin (if writable)"
    fi

    echo
    log "To get started, run:"
    if [ "$PATH_IN_PATH" = true ]; then
        echo "  $BINARY_NAME --help"
        echo "  $BINARY_NAME install    # This installs your dotfiles"
    else
        echo "  $INSTALL_DIR/$BINARY_NAME --help"
        echo "  $INSTALL_DIR/$BINARY_NAME install    # This installs your dotfiles"
    fi

    # Verify installation
    log "Verifying binary..."
    if "$INSTALL_DIR/$BINARY_NAME" --version >/dev/null 2>&1; then
        success "Binary verified successfully!"
    else
        warn "Binary verification failed, but file was copied successfully"
    fi

    # Run installer if requested
    if [ "$RUN_INSTALLER" = true ]; then
        echo
        success "Running dotfiles installer (this will install your dotfiles)..."

        # Default to 'install' command if no args provided
        if [ -z "$INSTALLER_ARGS" ]; then
            INSTALLER_ARGS="install"
        fi

        # Check if we can run directly or need full path
        if [ "$PATH_IN_PATH" = true ]; then
            log "Executing: $BINARY_NAME $INSTALLER_ARGS"
            exec $BINARY_NAME $INSTALLER_ARGS
        else
            log "Executing: $INSTALL_DIR/$BINARY_NAME $INSTALLER_ARGS"
            exec "$INSTALL_DIR/$BINARY_NAME" $INSTALLER_ARGS
        fi
    fi
}

# Run main function
main "$@"
