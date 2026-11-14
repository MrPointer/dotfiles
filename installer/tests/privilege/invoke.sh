#!/bin/bash
set -euo pipefail

work=${INSTALLER_PRIVILEGE_TEST_DIR:?}
sudo -K
if sudo -n true; then
    echo 'Expected password-required sudo with no cached authentication' >&2
    exit 2
fi
echo AUTH_READY
exec "$work/privilege.test" -test.v \
    -test.run '^Test_InstallingSystemChanges_Should_RequireAuthorization_When_RunAsNormalUser$'
