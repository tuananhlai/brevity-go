#!/bin/bash

set -e

# Set the ZSH theme to "robbyrussell".
sed -i '/^ZSH_THEME/c\ZSH_THEME="robbyrussell"' ~/.zshrc

# Project setup
go mod tidy
cp .env.example .env