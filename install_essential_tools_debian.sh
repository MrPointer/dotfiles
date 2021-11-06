#!/usr/bin/env bash

# saner programming env: these switches turn some bugs into errors
set -o pipefail -o nounset

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

function install_python_package_managers {
    info "Installing python package managers"

    info "Installing potentially missing venv-creation package"
    if ! sudo apt install -y python3-venv; then
        error "Failed installing python-venv package, can't create venvs without it"
        return 1
    fi
    success "Successfully installed python-venv"

    info "Installing pipx"
    if ! pip3 install --user pipx; then
        error "Failed installing pipx"
        return 2
    fi
    success "Successfully installed pipx"

    return 0
}

###
# Set global variables
###
function set_globals() {
    :
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

    _set_color_defaults
}

###
# This is the script's entry point, just like in any other programming language.
###
function main() {
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

    if ! install_python_package_managers; then
        error "Failed installing python package managers"
        return 3
    fi

    success "Successfully installed essential tools"
    return 0
}

# Call main and don't do anything else
# It will pass the correct exit code to the OS
main "$@"
