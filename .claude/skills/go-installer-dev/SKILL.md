---
name: go-installer-dev
description: Go development guide for the dotfiles-installer project. Use when writing Go code, adding features, fixing bugs, creating tests, working with package managers (brew, apt, dnf), implementing interfaces, using the commander/logger/filesystem utilities, handling privilege escalation, or understanding the codebase architecture. Covers coding patterns, testing conventions, interface design, error handling, and project structure.
---

# Go Installer Developer Guide

## Purpose

Comprehensive development guide for contributing to the dotfiles-installer Go codebase. Covers architecture, patterns, testing, and conventions.

## When to Use This Skill

Activates when working on:
- Go code in the installer directory
- Adding new features or packages
- Writing or modifying tests
- Implementing interfaces
- Working with package managers (brew, apt, dnf)
- Using commander, logger, or filesystem utilities
- Handling privilege escalation
- Understanding codebase architecture

---

## Project Structure

```
installer/
├── main.go                    # Entry point (version injection)
├── cmd/                       # CLI commands (Cobra)
│   ├── root.go               # Root command, global setup
│   ├── install.go            # Main installation workflow
│   ├── checkCompatibility.go # Compatibility check command
│   └── version.go            # Version display
├── lib/                       # Core business logic
│   ├── compatibility/        # OS/distro detection
│   ├── pkgmanager/          # Package manager interface
│   ├── brew/                # Homebrew implementation
│   ├── apt/                 # APT implementation
│   ├── dnf/                 # DNF implementation
│   ├── gpg/                 # GPG key management
│   ├── shell/               # Shell installation
│   ├── dotfilesmanager/     # Chezmoi integration
│   └── packageresolver/     # Package name resolution
├── utils/                     # Shared utilities
│   ├── logger/              # Logging with progress display
│   ├── osmanager/           # OS operations interface
│   ├── privilege/           # Sudo/doas escalation
│   ├── commander.go         # Command execution
│   ├── filesystem.go        # File operations
│   └── httpclient/          # HTTP client interface
├── cli/                       # Interactive UI components
├── internal/config/           # Embedded YAML configs
├── Taskfile.yml              # Task runner commands
└── .goreleaser.yaml          # Release configuration
```

---

## Core Patterns

### 1. Interface-Driven Design

All major components are defined as interfaces for testability:

```go
// Define interface
type PackageManager interface {
    InstallPackage(ctx context.Context, pkg RequestedPackageInfo) error
    IsPackageInstalled(ctx context.Context, name string) (bool, error)
    GetInfo(ctx context.Context) (PackageManagerInfo, error)
}

// Implement
type BrewPackageManager struct {
    logger    logger.Logger
    commander Commander
    // ...
}

func (b *BrewPackageManager) InstallPackage(ctx context.Context, pkg RequestedPackageInfo) error {
    // implementation
}
```

### 2. Dependency Injection

Pass dependencies via constructors:

```go
func NewBrewPackageManager(
    logger logger.Logger,
    commander Commander,
    osManager osmanager.OsManager,
    brewPath string,
) *BrewPackageManager {
    return &BrewPackageManager{
        logger:    logger,
        commander: commander,
        osManager: osManager,
        brewPath:  brewPath,
    }
}
```

### 3. Functional Options Pattern

For flexible command execution:

```go
result, err := commander.RunCommand(ctx, "brew", []string{"install", pkg},
    WithCaptureOutput(),
    WithEnv(env),
    WithTimeout(5 * time.Minute),
)
```

Available options: `WithEnv()`, `WithDir()`, `WithInput()`, `WithCaptureOutput()`, `WithDiscardOutput()`, `WithInteractive()`, `WithTimeout()`, `WithStdout()`, `WithStderr()`

### 4. Error Handling

Wrap errors with context:

```go
if err != nil {
    return fmt.Errorf("failed to install package %s: %w", pkg.Name, err)
}
```

### 5. Optional Types

Use `samber/mo` for safer nil handling:

```go
import "github.com/samber/mo"

type Config struct {
    Shell mo.Option[string]
}

// Usage
if shell, ok := config.Shell.Get(); ok {
    // use shell
}
```

---

## Key Interfaces

### Commander (utils/commander.go)

Execute system commands:

```go
type Commander interface {
    RunCommand(ctx context.Context, name string, args []string, opts ...Option) (Result, error)
}

// Result contains stdout, stderr, exit code, duration
```

### Logger (utils/logger/)

Logging with progress display:

```go
type Logger interface {
    Trace(format string, args ...any)
    Debug(format string, args ...any)
    Info(format string, args ...any)
    Success(format string, args ...any)
    Warning(format string, args ...any)
    Error(format string, args ...any)
    
    StartProgress(task string)
    UpdateProgress(message string)
    FinishProgress()
    FailProgress(err error)
}
```

### FileSystem (utils/filesystem.go)

File operations:

```go
type FileSystem interface {
    PathExists(path string) (bool, error)
    ReadFile(path string) ([]byte, error)
    WriteFile(path string, data []byte, perm os.FileMode) error
    CreateDirectory(path string) error
}
```

### OsManager (utils/osmanager/)

OS operations:

```go
type OsManager interface {
    UserManager
    EnvironmentManager
    ProgramQuery
    FilePermissionManager
}

type ProgramQuery interface {
    ProgramExists(name string) bool
    GetProgramVersion(name string) (string, error)
}
```

### Privilege Escalator (utils/privilege/)

Smart privilege escalation:

```go
type Escalator interface {
    EscalateCommand(cmd string, args []string) (string, []string, error)
    IsRunningAsRoot() bool
}
```

---

## Testing Conventions

Tests are co-located with source files (`brew.go` → `brew_test.go`).

For detailed testing guidelines, see **[Test Style Reference](references/test-style.md)**, which covers:
- Test naming conventions (literate names with `Test_` prefix)
- Table-driven test patterns
- Error assertion best practices
- Unit tests vs integration tests
- Mock usage with mockery/moq

### Running Tests

```bash
task test      # Run all tests with race detection
task cov       # Generate coverage report
task bench     # Run benchmarks
```

---

## Development Commands

```bash
task build     # Build binary via goreleaser
task test      # Run tests with race detection
task fmt       # Format code (gofumpt, goimports, golines)
task lint      # Run golangci-lint and typos
task check     # Run tests + lint
task sloc      # Print lines of code stats
```

---

## Adding New Features

### Adding a New Package Manager

1. Create package in `lib/{managername}/`
2. Implement `PackageManager` interface
3. Add installer interface if installation is needed
4. Add integration tests
5. Update `packagemap.yaml` for package name mappings

### Adding a New Utility

1. Define interface in `utils/`
2. Create implementation
3. Add mock file for testing
4. Inject via constructors where needed

### Adding a New Command

1. Create file in `cmd/`
2. Define cobra.Command
3. Register in `root.go` init()
4. Add flags and run logic

---

## Configuration Files

### compatibility.yaml

Defines supported OS/distros and prerequisites:

```yaml
supported_os:
  darwin:
    name: macOS
    distros:
      - name: macOS
        architectures: [amd64, arm64]
        prerequisites: [git, curl]
```

### packagemap.yaml

Maps generic package codes to manager-specific names:

```yaml
packages:
  git:
    brew: git
    apt: git
    dnf: git
  neovim:
    brew: neovim
    apt: neovim
    dnf: neovim
```

---

## Code Style

For detailed coding conventions, see **[Code Style Reference](references/code-style.md)**, which covers:
- Interface implementation verification (`var _ Interface = (*Struct)(nil)`)
- Constructor naming and placement
- Code formatting and alignment
- Documentation conventions
- Performance best practices (pre-allocation)
- Dependency injection patterns

### Linting Rules (from .golangci.yml)

- Max cyclomatic complexity: 20
- Max function arguments: 5
- No naked returns
- Prefer named returns
- nolint directives require explanation

### Formatting

Run before committing:
```bash
task fmt
```

Uses: gofumpt, goimports, golines

---

## Quick Reference

| Task | Command |
|------|---------|
| Build | `task build` |
| Test | `task test` |
| Format | `task fmt` |
| Lint | `task lint` |
| Full check | `task check` |
| Coverage | `task cov` |

| Pattern | Usage |
|---------|-------|
| Interface | Define behavior contracts |
| DI | Pass deps via constructors |
| Functional options | Flexible command config |
| Table-driven tests | Multiple test cases |
| Error wrapping | `fmt.Errorf("...: %w", err)` |
| Optional types | `mo.Option[T]` for nullable |

## Related Files

**Skill References:**
- [Code Style Reference](references/code-style.md) - Detailed Go coding conventions
- [Test Style Reference](references/test-style.md) - Testing patterns and mock usage

**Project Configuration:**
- `Taskfile.yml` - All development tasks
- `.golangci.yml` - Linting configuration
- `.goreleaser.yaml` - Release configuration
- `internal/config/*.yaml` - Embedded configurations
