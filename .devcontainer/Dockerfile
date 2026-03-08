FROM mcr.microsoft.com/devcontainers/go:1.25-bookworm

RUN apt-get update && apt-get upgrade -y

ARG TARGETARCH

# Install `asdf` (multiple runtime version manager) based on the target image cpu architecture.
RUN curl -L -o asdf.tar.gz "https://github.com/asdf-vm/asdf/releases/download/v0.18.0/asdf-v0.18.0-linux-${TARGETARCH}.tar.gz" && \
    tar -xzf asdf.tar.gz && \
    rm asdf.tar.gz && \
    mv asdf /usr/local/bin/

# Use `zsh` as the default shell.
RUN chsh -s $(which zsh) vscode

# `vscode` is a built-in user in the base image.
USER vscode

# Add `ASDF_DATA_DIR/shims` to PATH so that runtimes installed by asdf (i.e ruby, python, node) can be accessed.
ENV ASDF_DATA_DIR=/home/vscode/.asdf
ENV PATH="${ASDF_DATA_DIR}/shims:${PATH}"

# Homebrew installer must be run as non-root user, as instructed by the maintainer.
# https://docs.brew.sh/FAQ#why-does-homebrew-say-sudo-is-bad
RUN NONINTERACTIVE=1 /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Add `brew` command to the PATH of zsh (default shell).
RUN echo 'eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv zsh)"' >> /home/vscode/.zshrc

# Enable access to `brew` command and install packages.
RUN eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)" &&\
    brew install golang-migrate mockery opentofu gh httpie pgcli go-task

# golangci-lint **should not** be installed via brew since it can cause brew to install a different version
# of Go
RUN curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.9.0