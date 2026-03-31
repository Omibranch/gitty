#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/../.." && pwd)"
PKG_DIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"

if ! command -v dpkg-deb >/dev/null 2>&1; then
  echo "[ERROR] dpkg-deb not found. Run this script on Debian/Ubuntu (or install dpkg)."
  exit 1
fi

mkdir -p "$PKG_DIR/usr/bin" "$ROOT_DIR/pkg/out"

if [ -f "$ROOT_DIR/dist/gitty-linux-amd64" ]; then
  cp "$ROOT_DIR/dist/gitty-linux-amd64" "$PKG_DIR/usr/bin/gitty"
else
  echo "[INFO] dist/gitty-linux-amd64 not found, building from source..."
  (cd "$ROOT_DIR/source" && GOOS=linux GOARCH=amd64 go build -ldflags='-s -w' -o "$PKG_DIR/usr/bin/gitty" .)
fi

chmod 0755 "$PKG_DIR/usr/bin/gitty"
dpkg-deb --build "$PKG_DIR" "$ROOT_DIR/pkg/out/gitty_2.0.0_amd64.deb"

echo "[SUCCESS] Built: $ROOT_DIR/pkg/out/gitty_2.0.0_amd64.deb"
