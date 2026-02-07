# Docker E2E Testing

Docker environments for testing the installer on supported Linux distributions. Each distro has its own subdirectory with a Dockerfile, docker-compose.yml, and entrypoint.sh.

```
docker/
├── Taskfile.yml          # Orchestrates all container tasks
├── validate.sh           # Environment validation script
├── ubuntu/               # Ubuntu (apt)
├── debian/               # Debian (apt)
└── fedora/               # Fedora (dnf)
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

For E2E testing procedures, load the `testing-e2e-containers` skill.
