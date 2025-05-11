---
applyTo: "**/*.go"
---

# Go Tech Stack

## Main Tech Stack

- **Lipgloss**: A Go library for creating beautiful command-line applications. It is used to create the CLI for the installer.
- **Cobra**: A Go library for creating powerful command-line applications. It is used to create the CLI for the installer.
- **Viper**: A Go library for reading configuration files. It is used to read the configuration files for the installer.
- **GoReleaser**: A Go library for building and releasing Go applications. It is used to build the installer and create releases.
- **GitHub Actions**: A CI/CD tool used to automate the build and release process. It is used to build the installer and create releases.

## Testing Tech Stack

- **Testing**: The standard Go testing package is used for unit and integration tests.
- **Moq**: A Go library for generating mocks. It is used to generate mocks for testing.
- **Testcontainers**: A Go library for running tests in containers. It is used to run system tests in a container.
