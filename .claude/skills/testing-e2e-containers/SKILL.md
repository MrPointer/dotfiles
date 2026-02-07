---
name: testing-e2e-containers
description: Run E2E tests of the installer inside Docker containers on supported Linux distributions (Ubuntu, Debian, Fedora). Use when (1) verifying installer behavior on Linux, (2) testing new installer flags or features end-to-end, (3) running the installer in a controlled environment, (4) reproducing or debugging Linux-specific issues, or (5) validating that changes work across all supported distros before pushing to CI.
---

# E2E Container Testing

Test the installer binary inside Docker containers running the supported Linux distributions.

## Workflow

1. **Build** the installer for the container platform
2. **Start** a fresh container
3. **Run** the installer inside it
4. **Verify** results
5. **Stop** the container
6. **Repeat** for each distro (ubuntu, debian, fedora)

## Step 1: Build the Binary

Use goreleaser to build a snapshot binary. The docker-compose volumes mount `installer/` into containers at `/workspace/installer`, so the binary at `installer/bin/dotfiles-installer` becomes `/workspace/installer/bin/dotfiles-installer`.

Build for the host platform (when host matches container arch):

```bash
cd installer && task build
```

Cross-compile for Linux when running on macOS:

```bash
cd installer && GOOS=linux GOARCH=arm64 goreleaser build --single-target --skip before --snapshot --clean --output ./bin/dotfiles-installer-linux-arm64
```

Adjust `GOARCH` to match the container platform (`arm64` for Apple Silicon, `amd64` for Intel).

## Step 2: Run Tests in Containers

**Start a fresh container for each test scenario** to avoid state leakage from previous runs. Reuse the image, but always start from scratch.

```bash
cd installer/docker

# Start container
task ubuntu:start

# Run the installer
docker-compose -f ubuntu/docker-compose.yml exec --user testuser installer-test-ubuntu \
  sudo /workspace/installer/bin/dotfiles-installer install \
    --shell zsh --shell-source auto --non-interactive --install-prerequisites

# Verify results
docker-compose -f ubuntu/docker-compose.yml exec --user testuser installer-test-ubuntu \
  command -v zsh

# Tear down before next test
task ubuntu:stop
```

For standalone containers (without docker-compose), bind-mount the binary explicitly:

```bash
docker run -d --name test-ubuntu --platform linux/arm64 \
  -v "$(pwd)/../bin:/workspace/bin:ro" \
  ubuntu-installer-test-ubuntu:latest tail -f /dev/null

docker exec test-ubuntu /workspace/bin/dotfiles-installer-linux-arm64 install \
  --shell zsh --shell-source system --non-interactive

docker stop test-ubuntu && docker rm test-ubuntu
```

## Step 3: Verify Results

After running the installer, verify inside the container:

```bash
# Check PATH-visible binary
docker exec <container> command -v zsh

# Check specific paths (brew is NOT in PATH by default)
docker exec <container> ls -la /home/linuxbrew/.linuxbrew/bin/zsh  # brew
docker exec <container> ls -la /usr/bin/zsh                         # system (ubuntu/debian)
docker exec <container> ls -la /usr/sbin/zsh                        # system (fedora)

# Verify version
docker exec <container> zsh --version
```

## Testing All Distros

Run through all three distros. Always stop before starting the next test:

```bash
cd installer/docker

for distro in ubuntu debian fedora; do
  task ${distro}:start
  # ... run installer and verify ...
  task ${distro}:stop
done
```

## Gotchas

- See [distro-specific notes](references/distro-notes.md) for platform quirks (zsh paths, brew behavior, package manager differences)
- Homebrew's first run downloads portable Ruby (~30MB) and can take several minutes
- Fedora's first `dnf` triggers a repo metadata sync (~30MB) - be patient or use generous timeouts
- If brew errors with "process has already locked", a previous install was interrupted - start a fresh container
- When running long tests, prefer running in background and polling for completion over setting very long timeouts
