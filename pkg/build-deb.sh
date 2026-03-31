#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
PKG_DIR="$ROOT_DIR/pkg/deb"
OUT_DIR="$ROOT_DIR/pkg/out"
DOC_DIR="$PKG_DIR/usr/share/doc/gitty"

if ! command -v dpkg-deb >/dev/null 2>&1; then
  echo "[ERROR] dpkg-deb not found. Run this script on Debian/Ubuntu (or install dpkg)."
  exit 1
fi

mkdir -p "$PKG_DIR/usr/bin" "$OUT_DIR" "$DOC_DIR"

# Do not ship helper scripts in package root
rm -f "$PKG_DIR/build.sh"

if [ -f "$ROOT_DIR/dist/gitty-linux-amd64" ]; then
  cp "$ROOT_DIR/dist/gitty-linux-amd64" "$PKG_DIR/usr/bin/gitty"
else
  echo "[INFO] dist/gitty-linux-amd64 not found, building from source..."
  (cd "$ROOT_DIR/source" && GOOS=linux GOARCH=amd64 go build -ldflags='-s -w' -o "$PKG_DIR/usr/bin/gitty" .)
fi

chmod 0755 "$PKG_DIR/usr/bin/gitty"

# Debian documentation files required by lintian
cat > "$DOC_DIR/changelog" <<'EOF'
gitty (2.0.0) stable; urgency=medium

  * Add Linux package metadata and validation workflow.
  * Add AUR package descriptors and Debian package build support.

 -- Omibranch <omibranch@users.noreply.github.com>  Tue, 31 Mar 2026 18:00:00 +0000
EOF

gzip -n -f "$DOC_DIR/changelog"

cp "$ROOT_DIR/LICENSE" "$DOC_DIR/copyright"

# Ensure files in package are owned by root in resulting .deb
if [ "$(id -u)" -eq 0 ]; then
  chown -R root:root "$PKG_DIR"
fi

dpkg-deb --root-owner-group --build "$PKG_DIR" "$OUT_DIR/gitty_2.0.0_amd64.deb"

echo "[SUCCESS] Built: $OUT_DIR/gitty_2.0.0_amd64.deb"
