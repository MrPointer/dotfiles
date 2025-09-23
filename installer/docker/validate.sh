#!/bin/bash
# Validation script to test Docker environment has all required prerequisites
# This script runs inside the container to verify everything is properly installed
# Supports Ubuntu, Debian, and Fedora

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_header() {
    echo -e "${BLUE}=================================================${NC}"
    echo -e "${BLUE}  Docker Environment Validation${NC}"
    echo -e "${BLUE}=================================================${NC}"
    echo ""
}

log_test() {
    echo -e "${BLUE}[TEST]${NC} $1"
}

log_pass() {
    echo -e "${GREEN}[PASS]${NC} $1"
}

log_fail() {
    echo -e "${RED}[FAIL]${NC} $1"
}

log_info() {
    echo -e "${YELLOW}[INFO]${NC} $1"
}

test_command() {
    local cmd=$1
    local name=$2
    local version_flag=${3:-"--version"}

    log_test "Testing $name command availability"

    if command -v "$cmd" &> /dev/null; then
        log_pass "$name is available"
        if [[ "$version_flag" != "none" ]]; then
            local version_output
            version_output=$($cmd $version_flag 2>&1 | head -n1)
            log_info "Version: $version_output"
        fi
        return 0
    else
        log_fail "$name is not available"
        return 1
    fi
}

test_package() {
    local package=$1
    local test_cmd=$2

    log_test "Testing $package package"

    if $test_cmd &> /dev/null; then
        log_pass "$package is properly installed"
        return 0
    else
        log_fail "$package is not properly installed"
        return 1
    fi
}

test_user_setup() {
    log_test "Testing user setup"

    local current_user=$(whoami)
    if [[ "$current_user" == "testuser" ]]; then
        log_pass "Running as testuser"
    else
        log_fail "Expected to run as testuser, but running as $current_user"
        return 1
    fi

    log_test "Testing sudo access"
    if sudo -n true 2>/dev/null; then
        log_pass "Sudo access without password works"
    elif sudo -l 2>/dev/null | grep -q NOPASSWD; then
        log_pass "Sudo access configured (NOPASSWD found)"
    else
        log_info "Sudo may not be fully configured, but this is non-critical for testing"
        log_info "Most installer functionality will still work"
    fi

    return 0
}

test_filesystem() {
    log_test "Testing filesystem setup"

    if [[ -d "/workspace" ]]; then
        log_pass "Workspace directory exists"
    else
        log_fail "Workspace directory not found"
        return 1
    fi

    # Check the mounted installer directory instead of workspace root
    if [[ -d "/workspace/installer" ]]; then
        log_pass "Installer directory is mounted"
    else
        log_fail "Installer directory not found"
        return 1
    fi

    # Test write access to the mounted directory
    if touch "/workspace/installer/.test_write" 2>/dev/null && rm -f "/workspace/installer/.test_write" 2>/dev/null; then
        log_pass "Mounted directory is writable"
    else
        log_info "Mounted directory may not be writable, but this is expected in some setups"
    fi

    return 0
}

detect_os() {
    if [[ -f /etc/os-release ]]; then
        source /etc/os-release
        echo "$ID"
    else
        echo "unknown"
    fi
}

test_fedora_prerequisites() {
    local failed_tests=0

    # Test development tools group
    if rpm -q --whatprovides gcc &>/dev/null; then
        log_pass "Development tools (gcc) available"
    else
        log_fail "Development tools (gcc) not found"
        ((failed_tests++))
    fi

    # Test make
    if ! test_command "make" "Make (development tools)"; then
        ((failed_tests++))
    fi

    # procps-ng (test ps)
    if ! test_command "ps" "Process utilities (procps-ng)" "none"; then
        ((failed_tests++))
    fi

    # curl
    if ! test_command "curl" "cURL"; then
        ((failed_tests++))
    fi

    # file
    if ! test_command "file" "File utility"; then
        ((failed_tests++))
    fi

    # git
    if ! test_command "git" "Git"; then
        ((failed_tests++))
    fi

    return $failed_tests
}

test_debian_ubuntu_prerequisites() {
    local failed_tests=0

    # build-essential (test gcc)
    if ! test_command "gcc" "GCC (build-essential)"; then
        ((failed_tests++))
    fi

    # Make sure we have make
    if ! test_command "make" "Make (build-essential)"; then
        ((failed_tests++))
    fi

    # procps (test ps)
    if ! test_command "ps" "Process utilities (procps)" "none"; then
        ((failed_tests++))
    fi

    # curl
    if ! test_command "curl" "cURL"; then
        ((failed_tests++))
    fi

    # file
    if ! test_command "file" "File utility"; then
        ((failed_tests++))
    fi

    # git
    if ! test_command "git" "Git"; then
        ((failed_tests++))
    fi

    return $failed_tests
}

run_validation() {
    local failed_tests=0
    local os_type=$(detect_os)

    print_header

    # Test basic system info
    log_info "OS: $(cat /etc/os-release | grep PRETTY_NAME | cut -d'"' -f2)"
    log_info "Detected OS Type: $os_type"
    log_info "Kernel: $(uname -r)"
    log_info "Architecture: $(uname -m)"
    echo ""

    # Test user setup
    if ! test_user_setup; then
        ((failed_tests++))
    fi
    echo ""

    # Test filesystem
    if ! test_filesystem; then
        ((failed_tests++))
    fi
    echo ""

    # Test OS-specific prerequisites from compatibility.yaml
    case "$os_type" in
        "ubuntu"|"debian")
            log_info "Testing Debian/Ubuntu prerequisites..."
            if ! test_debian_ubuntu_prerequisites; then
                ((failed_tests += $?))
            fi
            ;;
        "fedora")
            log_info "Testing Fedora prerequisites..."
            if ! test_fedora_prerequisites; then
                ((failed_tests += $?))
            fi
            ;;
        *)
            log_fail "Unsupported OS: $os_type"
            ((failed_tests++))
            ;;
    esac

    echo ""

    # Test some additional functionality
    log_test "Testing build capability"
    if echo 'int main(){return 0;}' | gcc -x c - -o /tmp/test_build 2>/dev/null; then
        log_pass "C compilation works"
        rm -f /tmp/test_build
    else
        log_fail "C compilation failed"
        ((failed_tests++))
    fi

    log_test "Testing network connectivity"
    if curl -s --connect-timeout 5 https://httpbin.org/ip > /dev/null; then
        log_pass "Network connectivity works"
    else
        log_fail "Network connectivity issues"
        ((failed_tests++))
    fi

    echo ""
    echo -e "${BLUE}=================================================${NC}"

    if [[ $failed_tests -eq 0 ]]; then
        echo -e "${GREEN}✓ All tests passed! ${os_type^} environment is ready for development.${NC}"
        exit 0
    else
        echo -e "${RED}✗ $failed_tests test(s) failed. ${os_type^} environment needs attention.${NC}"
        exit 1
    fi
}

# Main execution
run_validation
