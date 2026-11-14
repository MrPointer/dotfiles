# Docker E2E Testing

Docker environments for testing the installer on supported Linux distributions. Each distro has its own subdirectory with a Dockerfile, docker-compose.yml, and entrypoint.sh.

```
docker/
├── Taskfile.yml          # Orchestrates all container tasks
├── validate.sh           # Environment validation script
├── ubuntu/               # Ubuntu (apt)
├── debian/               # Debian (apt)
├── fedora/               # Fedora (dnf)
└── privilege/            # Optional local privilege test image and launcher
```

## Container Management

Run `task` commands from `installer/docker/`:

| Command | Purpose |
|---------|---------|
| `task <distro>:start` | Build image and start container |
| `task <distro>:stop` | Stop container |
| `task <distro>:shell` | Enter container shell as testuser |
| `task <distro>:dev` | Start + shell in one step |
| `task <distro>:validate` | Run environment validation |
| `task <distro>:rebuild` | Rebuild image from scratch (no cache) |
| `task <distro>:clean` | Stop container and remove image |
| `task stop-all` | Stop all containers |
| `task validate-all` | Validate all environments |
| `task privilege:test` | Run the privilege integration checks in a local container |

Where `<distro>` is `ubuntu`, `debian`, or `fedora`.

## Container Details

| Distro | Image | Container (docker-compose) |
|--------|-------|---------------------------|
| Ubuntu | `ubuntu-installer-test-ubuntu:latest` | `installer-test-ubuntu-env` |
| Debian | `debian-installer-test-debian:latest` | `installer-test-debian-env` |
| Fedora | `fedora-installer-test-fedora:latest` | `installer-test-fedora-env` |

- Docker-compose bind-mounts `installer/` at `/workspace/installer` inside containers
- Containers use `testuser` (passwordless sudo) via docker-compose, or `root` with `docker run`
- Homebrew on Linux installs to `/home/linuxbrew/.linuxbrew/` and is NOT in PATH by default

## Privilege Integration Tests

Run `task privilege:test` from `installer/docker/` for local testing on macOS.
The task runs `task test:privilege:build` first, which cross-compiles
`bin/privilege.test` for Linux. Set `PRIVILEGE_TEST_ARCH` to select the target
architecture; it defaults to the host architecture. The Docker engine and image
must support that architecture, as with the existing Docker test environments.

The task builds the optional `installer privilege test image`
(`installer-privilege-test`) from `privilege/Dockerfile`. The Ubuntu 24.04 image
contains only dependencies for repeated local runs. The test binary and shared
scripts are bind-mounted read-only. The container has no network access and runs
all checks in one disposable container, resetting state between checks.

The Docker launcher uses the Docker-independent scripts in
`../tests/privilege/`. They exercise the real APT and shell-setup caller paths
with a cold sudo cache, once with the correct password and once with a rejected
password. Expect provides the terminal and password exchange. The shared runner
creates a password-required user, verifies that no protected operation occurs
before authentication, and verifies either the protected root-owned effect after
success or no effects and no unprivileged retry after rejection.

`privilege/run.sh` runs Docker commands asynchronously so it can handle signals.
Its first INT or TERM stops the child and force-removes the active container. A
second signal aborts stalled cleanup, so removal can remain incomplete if the
Docker daemon is unresponsive.

CI runs the same checks natively on a disposable Ubuntu runner; Docker is not
part of CI. For the native runner and its safety requirements, see [Privilege
Integration Tests](../AGENTS.md#privilege-integration-tests).

For E2E testing procedures, load the `testing-e2e-containers` skill.
