#!/bin/bash
set -euo pipefail

installer=$(cd "$(dirname "$0")/../.." && pwd)
image="installer-privilege-test"
child=
container=

stop_child() {
    if [[ -n "$child" ]]; then
        kill -KILL "$child" 2>/dev/null || true
        wait "$child" 2>/dev/null || true
        child=
    fi
}

cleanup() {
    local status=$?
    trap - EXIT
    # A second interrupt can abort cleanup if the Docker daemon is unresponsive.
    trap 'stop_child; exit "$status"' INT TERM
    stop_child
    if [[ -n "$container" ]]; then
        docker rm -f "$container" >/dev/null 2>&1 &
        child=$!
        wait "$child" || true
        child=
    fi
    exit "$status"
}

trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM

run_command() {
    # Bash handles signals immediately while waiting on an asynchronous child.
    # A foreground Docker command can otherwise defer the shell's signal trap.
    "$@" &
    child=$!
    wait "$child"
    child=
}

run_command docker build -t "$image" "$installer/docker/privilege"
container="installer-privilege-$$-$RANDOM"
run_command docker run --rm --init --name "$container" --network none \
    --mount "type=bind,src=$installer/bin/privilege.test,dst=/privilege.test,readonly" \
    --mount "type=bind,src=$installer/tests/privilege,dst=/checks,readonly" \
    -e INSTALLER_PRIVILEGE_TEST_DISPOSABLE=1 \
    "$image" /bin/bash /checks/run.sh /privilege.test
container=
