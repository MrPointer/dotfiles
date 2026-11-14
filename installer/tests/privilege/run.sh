#!/bin/bash
set -euo pipefail

if [[ $(uname -s) != Linux || $(id -u) != 0 ||
      ${INSTALLER_PRIVILEGE_TEST_DISPOSABLE:-} != 1 ]]; then
    echo 'Run as root only on disposable Ubuntu/Debian, with INSTALLER_PRIVILEGE_TEST_DISPOSABLE=1.' >&2
    exit 2
fi
if ! grep -Eq '^ID="?(ubuntu|debian)"?$' /etc/os-release; then
    echo 'The installer privilege checks require disposable Ubuntu or Debian.' >&2
    exit 2
fi
binary=${1:?usage: bash tests/privilege/run.sh /path/to/privilege.test}
scripts=$(cd "$(dirname "$0")" && pwd)
for command in sudo expect runuser useradd userdel usermod chpasswd visudo dpkg-deb flock; do
    command -v "$command" >/dev/null
done
# The checks modify /etc/shells and the package database; do not run concurrently.
exec 9>/var/lock/installer-privilege-tests.lock
flock -n 9

work=$(mktemp -d /var/tmp/installer-privilege.XXXXXX)
username="privtest$$"
shell_path="/usr/local/bin/installer-privilege-shell-$$"
sudoers="/etc/sudoers.d/installer-privilege-$$"
created_user=
package=installer-privilege-fixture
if id "$username" >/dev/null 2>&1 || [[ -e "$shell_path" ||
    -e /usr/share/installer-privilege-test || -e "$sudoers" ]] ||
    dpkg-query -W "$package" >/dev/null 2>&1; then
    rmdir "$work"
    echo 'Refusing to overwrite an existing test user, package, or fixture.' >&2
    exit 2
fi

cleanup() {
    local status=$?
    trap - EXIT
    if dpkg-query -W "$package" >/dev/null 2>&1; then
        dpkg --purge "$package" >/dev/null || status=1
    fi
    if [[ -f "$work/shells.before" ]]; then
        cp -p "$work/shells.before" /etc/shells || status=1
    fi
    rm -f "$sudoers" "$shell_path"
    if [[ -n "$created_user" ]]; then
        userdel --remove "$username" || status=1
    fi
    rm -rf "$work"
    exit "$status"
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM

chmod 755 "$work"
cp "$binary" "$work/privilege.test"
cp "$scripts/invoke.sh" "$scripts/observe-command.sh" "$scripts/test.exp" "$work/"
chmod 755 "$work/privilege.test" "$work/"*.sh
cp -p /etc/shells "$work/shells.before"
useradd --create-home --shell /bin/sh "$username"
created_user=1
printf '%s:privilege-test-only\n' "$username" | chpasswd
ln -s /bin/bash "$shell_path"
mkdir "$work/bin"
for command in apt tee usermod; do
    ln -s "$work/observe-command.sh" "$work/bin/$command"
done
cat > "$sudoers" <<EOF
Defaults:$username lecture=never, passwd_tries=1
Defaults:$username secure_path="$work/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
$username ALL=(ALL) PASSWD: ALL
EOF
chmod 440 "$sudoers"
visudo -cf /etc/sudoers

# A local package proves a root-owned effect without package download dependencies.
mkdir -p "$work/package/DEBIAN" "$work/package/usr/share/installer-privilege-test"
printf 'Package: %s\nVersion: 1.0\nArchitecture: all\nMaintainer: Test <test@example.invalid>\nDescription: Installer privilege regression fixture\n' \
    "$package" > "$work/package/DEBIAN/control"
printf 'installed by authenticated apt\n' > "$work/package/usr/share/installer-privilege-test/protected"
dpkg-deb --build --root-owner-group "$work/package" "$work/protected.deb"
printf 'installer-privilege-tests\n' > "$work/ready"

export INSTALLER_PRIVILEGE_TEST_DIR="$work"
export INSTALLER_PRIVILEGE_TEST_USER="$username"
export INSTALLER_PRIVILEGE_TEST_SHELL="$shell_path"
export TERM=xterm-256color LC_ALL=C
export PATH="$work/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
for caller in apt shell; do
    for scenario in cold reject; do
        echo "Testing installer privileges: $caller/$scenario"
        cp -p "$work/shells.before" /etc/shells
        /usr/sbin/usermod -s /bin/sh "$username"
        if dpkg-query -W "$package" >/dev/null 2>&1; then
            dpkg --purge "$package" >/dev/null
        fi
        cp -p /etc/passwd "$work/passwd.before"
        : > "$work/operations"
        chmod 666 "$work/operations"
        export INSTALLER_PRIVILEGE_TEST_CALLER="$caller"
        export INSTALLER_PRIVILEGE_TEST_SCENARIO="$scenario"
        expect "$work/test.exp"
        if [[ "$scenario" == reject ]]; then
            cmp /etc/shells "$work/shells.before"
            cmp /etc/passwd "$work/passwd.before"
            [[ ! -e /usr/share/installer-privilege-test/protected ]]
            [[ ! -s "$work/operations" ]]
        elif [[ "$caller" == apt ]]; then
            [[ $(cat /usr/share/installer-privilege-test/protected) == 'installed by authenticated apt' ]]
            [[ $(stat -c %u /usr/share/installer-privilege-test/protected) == 0 ]]
            printf '0 apt update\n0 apt install -y %s/protected.deb\n' "$work" > "$work/expected"
            cmp "$work/expected" "$work/operations"
            cmp /etc/shells "$work/shells.before"
            cmp /etc/passwd "$work/passwd.before"
        else
            grep -qx "$shell_path" /etc/shells
            [[ $(getent passwd "$username" | cut -d: -f7) == "$shell_path" ]]
            printf '0 tee -a /etc/shells\n0 usermod -s %s %s\n' "$shell_path" "$username" > "$work/expected"
            cmp "$work/expected" "$work/operations"
        fi
        echo "PASS: $caller/$scenario effects and execution identities"
    done
done
