# Distro-Specific Notes

Platform quirks and package manager differences relevant to E2E container testing.

## Ubuntu / Debian (apt)

- Package manager: `apt-get` (use over `apt` in scripts - stable CLI interface)
- Zsh installs to `/usr/bin/zsh` (or `/bin/zsh` which symlinks to it)
- `DEBIAN_FRONTEND=noninteractive` is set in the Dockerfile to suppress prompts
- Debian may emit `debconf` warnings about missing terminal - these are harmless
- `zsh-common` is a dependency of `zsh` and must be removed together when cleaning

```bash
# Install
sudo apt-get install -y zsh

# Remove (for clean re-test)
sudo apt-get remove -y zsh zsh-common
```

## Fedora (dnf)

- Package manager: `dnf`
- Zsh installs to `/usr/sbin/zsh` (not `/usr/bin/zsh` like Debian)
- First `dnf` operation after container start triggers a repository metadata download (~30MB) which can be slow
- Development tools installed via `@development-tools` group (not `build-essential`)

```bash
# Install
sudo dnf install -y zsh

# Remove (for clean re-test)
sudo dnf remove -y zsh
```

## Homebrew on Linux (all distros)

- Installs to `/home/linuxbrew/.linuxbrew/`
- Bin directory: `/home/linuxbrew/.linuxbrew/bin/`
- NOT in PATH by default inside containers - always check with full path
- First brew operation downloads portable Ruby (~30MB), subsequent operations are faster
- Brew's zsh path: `/home/linuxbrew/.linuxbrew/bin/zsh` (symlink to Cellar)
- Brew lock files: if a previous `brew install` was interrupted, a lock on the Cellar may remain - start a fresh container to avoid this

```bash
# Check if brew is available (not in PATH)
test -x /home/linuxbrew/.linuxbrew/bin/brew

# Get brew prefix
/home/linuxbrew/.linuxbrew/bin/brew --prefix

# Check brew-installed zsh
ls -la /home/linuxbrew/.linuxbrew/bin/zsh
/home/linuxbrew/.linuxbrew/bin/zsh --version
```

## Container Image Names

| Distro | Image | Container (docker-compose) |
|--------|-------|---------------------------|
| Ubuntu | `ubuntu-installer-test-ubuntu:latest` | `installer-test-ubuntu-env` |
| Debian | `debian-installer-test-debian:latest` | `installer-test-debian-env` |
| Fedora | `fedora-installer-test-fedora:latest` | `installer-test-fedora-env` |
