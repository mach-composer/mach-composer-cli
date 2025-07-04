#!/usr/bin/env bash

if ! (set -o pipefail 2>/dev/null); then
    echo "❌ This script requires bash. Please run it with: bash $0"
    exit 1
fi

set -euo pipefail

# Optional version from environment variable (default to latest from GitHub)
VERSION="${VERSION:-$(curl -s https://api.github.com/repos/mach-composer/mach-composer-cli/releases/latest | jq -r .tag_name)}"
VERSION_NO_V="${VERSION#v}"

# Detect OS
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "$OS" in
    linux|darwin|freebsd) ;;
    msys*|mingw*|cygwin*) OS="windows" ;;
    *)
        echo "❌ Unsupported OS: $OS"
        exit 1
        ;;
esac

# Detect architecture
ARCH="$(uname -m)"
case "$ARCH" in
    x86_64|amd64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    armv6l) ARCH="armv6" ;;
    *)
        echo "❌ Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

# Determine file extension
EXT="tar.gz"
[[ "$OS" == "windows" ]] && EXT="zip"

FILENAME="mach-composer-${VERSION_NO_V}-${OS}-${ARCH}.${EXT}"
URL="https://github.com/mach-composer/mach-composer-cli/releases/download/${VERSION}/${FILENAME}"

# Setup directories
TARGET_DIR="$HOME/.local/bin"
mkdir -p "$TARGET_DIR"

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

echo "⬇️ Downloading $FILENAME..."
curl -sL "$URL" -o "$TMP_DIR/$FILENAME"

echo "📦 Extracting $FILENAME..."
if [[ "$EXT" == "zip" ]]; then
    unzip -q "$TMP_DIR/$FILENAME" -d "$TMP_DIR"
else
    tar -xzf "$TMP_DIR/$FILENAME" -C "$TMP_DIR"
fi

# Move binary from bin/ to target dir
BIN_PATH="$TMP_DIR/bin/mach-composer"
if [[ ! -f "$BIN_PATH" ]]; then
    echo "❌ Binary not found at expected path: $BIN_PATH"
    exit 1
fi

mv "$BIN_PATH" "$TARGET_DIR/mach-composer"
chmod +x "$TARGET_DIR/mach-composer"

echo "✅ mach-composer $VERSION installed to $TARGET_DIR"

# Warn if not in PATH
if ! echo "$PATH" | grep -q "$HOME/.local/bin"; then
    echo "⚠️  $HOME/.local/bin is not in your PATH. You can add this to your shell config:"
    echo 'export PATH="$HOME/.local/bin:$PATH"'
fi
