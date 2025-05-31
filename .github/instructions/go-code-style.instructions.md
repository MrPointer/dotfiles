---
applyTo: "**/*.go"
---

# Go Coding Style

## General Guidelines

- Use the Go standard library whenever possible. Only use third-party libraries when necessary.
- Limit line length to 120 characters.
- Write code that is easy to test:
  - Use interfaces to decouple components and improve testability.
  - Use dependency injection to pass dependencies into functions and methods.
  - Wrap even basic operations (such as OS functions and file operations) in interfaces to make them easier to mock and test.
- After each struct definition, verify interface implementation by adding:
  `var _ InterfaceName = (*StructName)(nil)`
- Provide a constructor function for each struct, named `NewStructName`.
  - Place this function immediately after the struct definition and the interface assertion line (if present).
- Format code for readability:
  - Vertically align function arguments when there are multiple arguments.
  - Insert blank lines between logical sections of code.
  - Do not separate error unwrapping from related code with a blank line; treat it as part of the same section.
- End all type and function comments with a period, following Go conventions.
- Pre-allocate collections (such as slices and maps) to their expected size when possible to reduce memory allocations and improve performance.

## Main Tech Stack

- [lipgloss]: Go library for creating visually appealing command-line applications. Used for the installer's CLI.
- [cobra]: Go library for building command-line applications. Used for the installer's CLI.
- [viper]: Go library for reading configuration files. Used for the installer's configuration management.
- [goreleaser]: Go tool for building and releasing Go applications. Used for building and releasing the installer.
- [gh-actions] (GitHub Actions): CI/CD tool for automating build and release processes. Used for building and releasing the installer.

[lipgloss]: https://github.com/charmbracelet/lipgloss
[cobra]: https://github.com/spf13/cobra
[viper]: https://github.com/spf13/viper
[goreleaser]: https://github.com/goreleaser/goreleaser
[gh-actions]: https://github.com/features/actions
