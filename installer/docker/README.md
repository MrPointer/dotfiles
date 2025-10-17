# Multi-OS Docker Testing Environment

This directory contains Docker configurations for creating consistent testing environments across different operating systems with all required prerequisites pre-installed.

## Supported Operating Systems

The Docker environment supports the following operating systems based on `../internal/config/compatibility.yaml`:

- **Ubuntu Latest** - All prerequisites via `apt-get`
- **Debian Latest** - All prerequisites via `apt-get`
- **Fedora Latest** - All prerequisites via `dnf`

Each OS includes:
- Build tools (gcc, make, development packages)
- Process utilities (ps, procps)
- Network tools (curl)
- File utilities (file command)
- Version control (git)
- SSL certificates (ca-certificates)

## Quick Start

### List Available Operating Systems
```bash
task list
```

### Ubuntu Environment
```bash
# All-in-one development command
task ubuntu:dev

# Or step by step:
task ubuntu:start    # Start Ubuntu environment
task ubuntu:shell    # Enter the container shell
task ubuntu:stop     # Stop when done
```

### Debian Environment
```bash
task debian:dev      # Start Debian and enter shell
task debian:validate # Run validation tests
task debian:stop     # Stop environment
```

### Fedora Environment
```bash
task fedora:dev      # Start Fedora and enter shell
task fedora:validate # Run validation tests
task fedora:stop     # Stop environment
```

## Available Tasks

Run `task --list` to see all available tasks. Here are the main patterns:

### Per-OS Tasks
Each OS (`ubuntu`, `debian`, `fedora`) supports:
- `task <os>:start` - Build and start the environment
- `task <os>:stop` - Stop the environment
- `task <os>:shell` - Enter the container shell
- `task <os>:dev` - Start environment and enter shell (recommended)
- `task <os>:validate` - Run validation tests
- `task <os>:rebuild` - Rebuild the Docker image from scratch
- `task <os>:clean` - Clean up containers and images

### Multi-OS Tasks
- `task status` - Show status of all environments
- `task stop-all` - Stop all running environments
- `task clean-all` - Clean up all environments
- `task validate-all` - Run validation on all environments

## Directory Structure

```
docker/
├── README.md              # This file
├── Taskfile.yml          # Multi-OS task definitions
├── validate.sh           # OS-aware validation script
├── ubuntu/               # Ubuntu-specific files
│   ├── Dockerfile
│   └── docker-compose.yml
├── debian/               # Debian-specific files
│   ├── Dockerfile
│   └── docker-compose.yml
└── fedora/               # Fedora-specific files
    ├── Dockerfile
    └── docker-compose.yml
```

## Features

- **Multiple OS support** - Test across Ubuntu, Debian, and Fedora
- **OS-specific prerequisites** - Each image has the correct packages for that OS
- **Non-root user setup** - All containers run as `testuser` with sudo privileges
- **Volume mounting** - Your installer code is mounted at `/workspace/installer`
- **Persistent caches** - Separate Docker volumes for each OS build artifacts
- **Clean environments** - Fresh OS installation every time
- **Comprehensive validation** - OS-aware validation script
- **Task automation** - Streamlined workflow with OS-namespaced tasks

## Usage Examples

### Development Workflow
```bash
# Start your preferred OS environment
task ubuntu:dev         # Ubuntu development
# or
task debian:dev         # Debian development
# or
task fedora:dev         # Fedora development

# The container will start and you'll be dropped into a shell
# Your installer code is available at /workspace/installer
```

### Testing Across Multiple OS
```bash
# Validate all environments work
task validate-all

# Check status of all environments
task status

# Clean up everything when done
task clean-all
```

### Validation Results
```bash
task ubuntu:validate
# Example output:
# 🧪 Running Ubuntu validation tests...
# [PASS] Running as testuser
# [PASS] Sudo access without password works
# [PASS] GCC (build-essential) is available
# [PASS] All prerequisites installed
# ✅ All tests passed! Ubuntu environment is ready for development.
```

### OS-Specific Testing
```bash
# Test Ubuntu-specific behavior
task ubuntu:start
task ubuntu:shell
# ... do Ubuntu-specific testing ...

# Switch to Fedora for RPM-based testing
task fedora:start
task fedora:shell
# ... test dnf/rpm behavior ...
```

## Troubleshooting

### Permission Issues
All containers run as `testuser` (UID 1000) by default. The mounted installer directory inherits permissions from your host filesystem.

### Rebuilding After Changes
If you modify any Dockerfile:
```bash
task <os>:rebuild    # Rebuild specific OS
# or
task clean-all       # Clean everything and rebuild as needed
```

### Container Conflicts
Each OS uses separate containers and volumes:
- Ubuntu: `installer-test-ubuntu-env`
- Debian: `installer-test-debian-env`
- Fedora: `installer-test-fedora-env`

### Cleaning Up
```bash
# Clean specific OS
task ubuntu:clean

# Clean everything
task clean-all

# Nuclear option - remove all Docker artifacts
docker system prune -a
```

## Interactive GPG Testing

⚠️ **SAFETY WARNING**: When testing locally, the scripts use isolated test environments (`/tmp/test-*-home`) to prevent any damage to your real system. Never run the installer directly on your system during development!

This directory also contains tools for testing the dotfiles installer's interactive GPG setup functionality using automated expect scripts.

### Prerequisites for Interactive Testing

**Install expect:**

macOS:
```bash
brew install expect
```

Ubuntu/Debian:
```bash
sudo apt-get install expect
```

### Interactive Testing Script

The `./test-interactive-gpg.exp` script automates GPG key setup during interactive installation using the `expect` tool.

**Basic Usage:**
```bash
# Use defaults (from installer directory)
./test-interactive-gpg.exp

# Specify custom parameters
./test-interactive-gpg.exp [installer_path] [email] [name] [passphrase]
```

**Examples:**
```bash
# Test with default values
./test-interactive-gpg.exp

# Test with custom GPG info
./test-interactive-gpg.exp \
  "./dotfiles-installer" \
  "your@email.com" \
  "Your Full Name" \
  "your-secure-passphrase"

# Test with built binary in Docker environment
./test-interactive-gpg.exp \
  "./dist/dotfiles_installer_linux_amd64_v1/dotfiles-installer" \
  "test@example.com" \
  "Test User" \
  "test-passphrase"
```

### GPG Prompts Handled

The expect script handles these GPG-specific interactive prompts:

- ✅ GPG email address input
- ✅ GPG full name input
- ✅ GPG passphrase input
- ✅ GPG "okay" confirmation (sends "O" as required by GPG)
- ✅ GPG key generation parameters (size, expiration)
- ✅ GPG key type selection
- ✅ GPG comment field

### Interactive Testing Workflow

1. **Build the installer:**
   ```bash
   cd installer
   goreleaser build --skip before --snapshot --clean
   ```
2. **Start Container** (e.g. Ubuntu)
   ```bash
   task ubuntu:dev
   ```
3. **Run interactive test:**
   ```bash
   ./test-interactive-gpg.exp
   ```

### Troubleshooting Interactive Tests

**Script hangs or times out:**
- Enable debug mode by editing the script and setting `exp_internal 1`
- Run manually and observe the exact GPG prompt text

**GPG errors in containers:**
- This is expected in containerized environments
- The script handles these gracefully for CI testing

**Custom GPG prompts not handled:**
- Add new patterns to the `expect` block in the script
- Use case-insensitive regex patterns: `(?i).*your_pattern`

## Development Tips

- Use `task status` to see which environments are running
- Each OS has its own cache volume for faster rebuilds
- The validation script automatically detects the OS and runs appropriate tests
- You can run multiple OS environments simultaneously for comparison testing
- All environments mount the same installer code, so changes are immediately available
- Interactive GPG testing works in both Docker environments and locally
- The expect script uses a 5-minute timeout to prevent hanging in CI
