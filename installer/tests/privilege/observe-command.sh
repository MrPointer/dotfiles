#!/bin/sh
set -eu

# Record both elevated operations and any unprivileged retry, then run the real command.
work=$(dirname "$(readlink -f "$0")")
command_name=${0##*/}
printf '%s %s %s\n' "$(id -u)" "$command_name" "$*" >> "$work/operations"
case "$command_name" in
    apt)
        # Keep apt update offline without changing the machine's configured sources.
        exec /usr/bin/apt -o Dir::Etc::sourcelist=/dev/null -o Dir::Etc::sourceparts=- "$@"
        ;;
    tee) exec /usr/bin/tee "$@" ;;
    usermod) exec /usr/sbin/usermod "$@" ;;
    *) exit 2 ;;
esac
