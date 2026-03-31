#!/bin/sh
set -e

REPO="Omibranch/gitty"
INSTALL_DIR="/usr/local/bin"
BIN_NAME="gitty"

# ── detect OS and arch ────────────────────────────────────────────
OS="$(uname -s)"
ARCH="$(uname -m)"

case "$OS" in
  Linux)  os="linux" ;;
  Darwin) os="darwin" ;;
  *)
    echo "[ERROR] Unsupported OS: $OS"
    exit 1
    ;;
esac

case "$ARCH" in
  x86_64|amd64)   arch="amd64" ;;
  aarch64|arm64)  arch="arm64" ;;
  *)
    echo "[ERROR] Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

ASSET="gitty-${os}-${arch}"

# ── get latest release URL ────────────────────────────────────────
echo "[INFO] Fetching latest release of gitty..."
DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${ASSET}"

# ── download ──────────────────────────────────────────────────────
TMP="$(mktemp)"
echo "[INFO] Downloading ${ASSET}..."

if command -v curl >/dev/null 2>&1; then
  curl -fsSL "$DOWNLOAD_URL" -o "$TMP"
elif command -v wget >/dev/null 2>&1; then
  wget -qO "$TMP" "$DOWNLOAD_URL"
else
  echo "[ERROR] Neither curl nor wget found. Install one and retry."
  exit 1
fi

chmod +x "$TMP"

# ── install ───────────────────────────────────────────────────────
if [ -w "$INSTALL_DIR" ]; then
  mv "$TMP" "${INSTALL_DIR}/${BIN_NAME}"
else
  echo "[INFO] Root required to install to ${INSTALL_DIR}..."
  sudo mv "$TMP" "${INSTALL_DIR}/${BIN_NAME}"
fi

echo "[SUCCESS] gitty installed to ${INSTALL_DIR}/${BIN_NAME}"
echo "[INFO] Run: gitty help"
