# Dotfiles

## Overview

Dotfiles projects are curated collections of hidden configuration files (dotfiles) that users store in version control systems like GitHub. These projects often include scripts or automation tools to apply the configurations on any machine, enabling users—especially developers—to quickly set up and maintain a consistent, personalized environment across multiple systems.

## My Dotfiles

My personal dotfiles project is managed at [https://github.com/MrPointer/dotfiles](https://github.com/MrPointer/dotfiles).  
It uses [chezmoi] as the primary management tool, after installing some basic prerequisites. The installation process is handled by a dedicated installer binary written in Go.

The goal is to create an easy, portable process that can be used on any machine. [chezmoi] simplifies this by supporting templating and conditional configuration.  
The repository contains the actual dotfiles, primarily for the Zsh shell, including plugins and auto-completion for selected tools.

I separate personal and work dotfiles by storing them in different locations. In work environments, both sets are loaded, but in personal environments, only the personal dotfiles are applied, keeping the setup clean.

## Installer/Bootstrapper

### Overview

The core of this project is an installer/bootstrapper, which acts as the entry point for setup.  
This tool installs all prerequisites on a new system (supporting multiple operating systems and distributions), installs [chezmoi], populates its data file with custom keys (such as whether the environment is for work), and applies the dotfiles.

A key prerequisite is [homebrew], which itself may require additional dependencies. The installer manages these in an OS/distro-aware manner.  
Note: [homebrew] is optional and can be skipped.

### Motivation

While all installer actions could be performed manually (with [chezmoi] handling most of the configuration), automating the process saves time and reduces complexity.

### Switch to Go

The installer was originally written in Bash, with a POSIX sh entry point for compatibility. As complexity increased, maintenance and testing became difficult.

The installer has been rewritten in Go for better maintainability and testability, with unit tests for each operation.  
It is a [cobra]-based CLI application, using [huh] for interactive selections (e.g., reusing existing GPG keys).  
For interactive sessions, the installer displays a progress bar for each operation, with a hierarchical system to expand real-time details.

### Source Code Location

- The Go installer source code is in the `installer` directory at the repository root.
- `go.mod` and `go.sum` are located in the `installer` directory.
- Run Go-related tools (e.g., tests) from within the `installer` directory.

### Tech Stack

- [cobra]: Go library for building command-line applications (used for the installer's CLI).
- [viper]: Go library for configuration management.
- [goreleaser]: Go tool for building and releasing applications.
- [gh-actions] (GitHub Actions): CI/CD tool for automating build and release processes.

[chezmoi]: https://www.chezmoi.io/
[homebrew]: https://docs.brew.sh/
[huh]: https://github.com/charmbracelet/huh
[cobra]: https://github.com/spf13/cobra
[viper]: https://github.com/spf13/viper
[goreleaser]: https://github.com/goreleaser/goreleaser
[gh-actions]: https://github.com/features/actions
