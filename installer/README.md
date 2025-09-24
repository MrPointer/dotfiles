# Dotfiles Installer

[![License](https://img.shields.io/github/license/MrPointer/dotfiles_installer)](https://github.com/MrPointer/dotfiles_installer)
[![Go Report Card](https://goreportcard.com/badge/github.com/MrPointer/dotfiles_installer)](https://goreportcard.com/report/github.com/MrPointer/dotfiles_installer)
[![GitHub](https://img.shields.io/github/v/release/MrPointer/dotfiles_installer?logo=github&sort=semver)](https://github.com/MrPointer/dotfiles_installer)
[![CI](https://github.com/MrPointer/dotfiles_installer/workflows/CI/badge.svg)](https://github.com/MrPointer/dotfiles_installer/actions?query=workflow%3ACI)
[![codecov](https://codecov.io/gh/MrPointer/dotfiles_installer/branch/main/graph/badge.svg)](https://codecov.io/gh/MrPointer/dotfiles_installer)

The installer of my dotfiles, used to bootstrap the system

> [!WARNING]
> **Dotfiles Installer is in early development and is not yet ready for use**

![caution](./img/caution.png)

## Project Description

## Features

- **Hierarchical Progress Display**: Shows npm-style progress indicators with spinners and timing information
- **Automatic Cursor Cleanup**: Ensures terminal cursor is always visible after program exit, even on interruption
- **Signal Handling**: Gracefully handles Ctrl+C and other termination signals with proper cleanup
- **Verbosity Control**: Multiple verbosity levels from minimal to extra-verbose output
- **Non-Interactive Mode**: Supports automated installations without user prompts

## Installation

Compiled binaries for all supported platforms can be found in the [GitHub release]. There is also a [homebrew] tap:

```shell
brew install MrPointer/tap/dotfiles_installer
```

## Quickstart

### Credits

This package was created with [copier] and the [FollowTheProcess/go_copier] project template.

[copier]: https://copier.readthedocs.io/en/stable/
[FollowTheProcess/go_copier]: https://github.com/FollowTheProcess/go_copier
[GitHub release]: https://github.com/MrPointer/dotfiles_installer/releases
[homebrew]: https://brew.sh
