# MrPointer's dotfiles

## Motivation

Like any other dotfiles project, I'm looking to create myself
a templated solution that will help me apply it on new environments, mostly Unix ones.  
I'm using a dotfiles manager alongside some custom shell scripts to achieve this,
managing both home and office/work environments.  

## Installation

To install, simply run one of the following commands
in the environment you'd like to install in:  

| Tool | Command                                                                                        |
| ---- | ---------------------------------------------------------------------------------------------- |
| Curl | `bash -c "$(curl -fsSL https://raw.githubusercontent.com/MrPointer/dotfiles/main/install.sh)"` |
| Wget | `bash -c "$(wget -O- https://raw.githubusercontent.com/MrPointer/dotfiles/main/install.sh)"`   |

The bootstrap script will take you through a configuration process
that would query some required info, and then will install the dotfiles manager and apply it.  

### Installation Options

The following options can be passed to the command above to customize the installation:  

| Option                     | Description                                                                                    |
| -------------------------- | ---------------------------------------------------------------------------------------------- |
| `-v` or `--verbose`        | Enable verbose output                                                                          |
| `--package-manager`        | Package manager to use for installing prerequisites                                            |
| `--system-package-manager` | Treat the given package manager as a system package manager, i.e. run it as root               |
| `--work-environment`       | Treat this installation as a work environment                                                  |
| `--work-email=[email]`     | Use given email address as work's email address                                                |
| `--no-python`              | Don't install python                                                                           |
| `--no-gpg`                 | Don't install gpg                                                                              |
| `--no-brew`                | Don't install brew (Linuxbrew/Homebrew)                                                        |
| `--prefer-brew`            | Prefer installing "system" tools with brew rather than package manager (Doesn't apply for Mac) |

## Overview

### Dotfiles Manager

I'm using a dedicated dotfiles manager, [chezmoi][chezmoi-url], which provides templating abilities,
per-machine differences, and a lot more.  
Other managers might be considered in the future, especially if [chezmoi][chezmoi-url] becomes stale.  

### Installation Process

#### Bootstrap

To create a single-click installation experience, I'm bootstrapping the process
by making sure everything is available, and only then proceed with the actual installation.  

The main installation driver script is written in "Pure" shell,
guaranteed to work on almost all systems, even the strangest ones.  
It checks whether `bash` is available, trying to install it if not. 
The installation utilizes the guessed system package manager, e.g. `apt` for `Debian` systems.  
If the installation fails for some reason, then the user is prompted to manually install `bash`.  
After `bash` is properly installed, it's used to execute the actual installation script, written in `bash`.  

#### Actual Installation

The actual installation script installs the dotfiles manager in some way, preferably standalone,
creates a config file for it, prompting the user for some required info such as name and email,
and then *"applies"* the template.  
The dotfiles manager can also install some additional packages on its own, depending on the configuration,
target machine, and the manager itself.  
At last, the script tries to reinstall the dotfiles manager in a way that will get it updated,
but only if the user has configured it previously and have the correct package managers installed.

[chezmoi-url]: https://www.chezmoi.io/
