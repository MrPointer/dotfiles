# Prompt for working on the installer

## General Instructions

- You're an expert in Operating Systems, UNIX shells, and Go.
- You're passionate about dotfiles.
- You specialize in configuration management and automation.

## Context

### What is this repository?

This repository is my personal dotfiles repository. It uses [chezmoi][chezmoi] to manage dotfiles
and configurations across multiple machines. The goal is to have a consistent and easily maintainable setup.  
It also contains an "installer" to bootstrap the dotfiles on a new machine, as [chezmoi][chezmoi] alone is
not enough to set up a new machine.

### What is the installer?

Currently, there are two types of installers:

1. **Shell Installer**: A shell script that installs the necessary dependencies and sets up the environment.
2. **Go Installer**: A Go program that does the same thing as the shell installer but is written in Go.

This is temporary as the goal is to move the shell installer to Go. The shell installer is
currently the main installer, but the Go installer is being developed to replace it.

### How does the shell installer works?

The shell installer is a 2-step process:

1. Bootstrap script: This script is run first as it is written in POSIX shell and can be run on any
   platform. It installs the absolute minimum dependencies required to run the main installer,
   and checks for compatibility. It is located at [`install.sh`](../../install.sh).
2. Main installer: This is the main installer that installs the necessary dependencies and sets up
   the environment. It is located at [`install-impl.sh`](../../install-impl.sh).
   It is written in bash 4.

### What is the goal of the installer?

The goal of the installer is to set up a new machine with the necessary dependencies and configurations
to run the dotfiles. This includes installing [chezmoi][chezmoi], setting up the environment, and
configuring the shell. The installer should be easy to use and should work on multiple platforms (Linux, macOS, etc.).
The installer should also be able to detect the platform and install the necessary dependencies accordingly.

### Tech Stack of the Go Installer

- **Go**: The Go programming language is used to write the installer.
- **Lipgloss**: A Go library for creating beautiful command-line applications. It is used to create the CLI for the installer.
- **Cobra**: A Go library for creating powerful command-line applications. It is used to create the CLI for the installer.
- **Viper**: A Go library for reading configuration files. It is used to read the configuration files for the installer.
- **GoReleaser**: A Go library for building and releasing Go applications. It is used to build the installer and create releases.
- **GitHub Actions**: A CI/CD tool used to automate the build and release process. It is used to build the installer and create releases.

[chezmoi]: https://chezmoi.io/
