#!/bin/sh
set -e

REPO="justinsautter/bitsplitter"
INSTALL_DIR="/usr/local/bin"

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
    linux)  OS="linux" ;;
    darwin) OS="darwin" ;;
    *)      echo "Unsupported OS: $OS"; exit 1 ;;
esac

# Detect architecture
ARCH=$(uname -m)
case "$ARCH" in
    x86_64|amd64)  ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *)             echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Get latest release tag
LATEST=$(curl -sI "https://github.com/${REPO}/releases/latest" | grep -i "^location:" | sed 's/.*tag\///' | tr -d '\r\n')
if [ -z "$LATEST" ]; then
    echo "Error: could not determine latest release"
    exit 1
fi

ARCHIVE="bitsplitter_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/${LATEST}/${ARCHIVE}"

echo "Downloading bitsplitter ${LATEST} for ${OS}/${ARCH}..."

TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT

curl -sL "$URL" -o "${TMPDIR}/${ARCHIVE}"
tar -xzf "${TMPDIR}/${ARCHIVE}" -C "$TMPDIR"

echo "Installing to ${INSTALL_DIR}/bitsplitter..."
sudo install -m 755 "${TMPDIR}/bitsplitter" "${INSTALL_DIR}/bitsplitter"

echo "bitsplitter installed successfully!"
bitsplitter --help 2>/dev/null || true
