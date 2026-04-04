FROM mcr.microsoft.com/devcontainers/go:1.25-bookworm

RUN apt-get update && apt-get upgrade -y

# Use `zsh` as the default shell.
RUN chsh -s $(which zsh) vscode

# `vscode` is a built-in user in the base image.
USER vscode

# Homebrew installer must be run as non-root user, as instructed by the maintainer.
# https://docs.brew.sh/FAQ#why-does-homebrew-say-sudo-is-bad
RUN NONINTERACTIVE=1 /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)" && \
    # Add `brew` command to the PATH of zsh (default shell).
    echo 'eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv zsh)"' >> /home/vscode/.zshrc && \
    # Enable access to `brew` command and install packages.
    eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)" && \
    brew install golang-migrate mockery opentofu gh httpie pgcli go-task

# golangci-lint **should not** be installed via brew since it can cause brew to install a different version
# of Go
RUN curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.9.0