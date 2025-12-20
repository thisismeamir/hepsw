#!/bin/bash
# HepSW Installation Script
# Usage: curl -sSL https://raw.githubusercontent.com/thisismeamir/hepsw/main/install.sh | bash

set -e

REPO="thisismeamir/hepsw"
INSTALL_DIR="/usr/local/bin"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

echo "Detected: $OS-$ARCH"

# Get latest release
LATEST_RELEASE=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_RELEASE" ]; then
    echo "Failed to get latest release"
    exit 1
fi

echo "Latest version: $LATEST_RELEASE"

# Download URL
BINARY_NAME="hepsw-${OS}-${ARCH}"
DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_RELEASE/${BINARY_NAME}.tar.gz"

echo "Downloading from: $DOWNLOAD_URL"

# Download and extract
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"

curl -sSL "$DOWNLOAD_URL" | tar xz

# Install
if [ -w "$INSTALL_DIR" ]; then
    mv "$BINARY_NAME" "$INSTALL_DIR/hepsw"
    chmod +x "$INSTALL_DIR/hepsw"
else
    echo "Installing to $INSTALL_DIR requires sudo..."
    sudo mv "$BINARY_NAME" "$INSTALL_DIR/hepsw"
    sudo chmod +x "$INSTALL_DIR/hepsw"
fi

# Cleanup
cd -
rm -rf "$TMP_DIR"

echo ""
echo "✓ HepSW installed successfully!"
echo ""
echo "Get started:"
echo "  hepsw init ~/hep-workspace"
echo "  hepsw build root"
echo ""
echo "Documentation: https://hepsw.readthedocs.io"