# Official repositories plan (AUR / pacman / apt)

## ✅ AUR (official community registry) — ready to publish

Prepared files:
- `pkg/aur/PKGBUILD`
- `pkg/aur/.SRCINFO`

Publication is maintainer-driven (push to `aur.archlinux.org`), so final step requires your AUR account/SSH key.

---

## ⚠ pacman official repos (Arch `extra`/`community`)

Direct self-publish is **not available**. Packages are accepted by Arch maintainers/TUs based on policy, demand, quality, maintenance.

Practical path:
1. Publish and maintain in AUR first.
2. Gain users/votes/issues history.
3. Request official inclusion through Arch packaging channels.

---

## ⚠ apt official repos (Debian/Ubuntu main/universe)

Direct self-publish to Debian/Ubuntu official repos is **not immediate**.
It requires Debian packaging workflow + sponsor/review + policy compliance.

Practical path:
1. Keep GitHub Releases with Linux binaries.
2. Build `.deb` with `pkg/deb/build.sh`.
3. Publish your own apt repo (or PPA) for `apt install gitty` from your source.
4. In parallel, pursue Debian/Ubuntu official inclusion process.

---

## Build `.deb`

Run on Debian/Ubuntu/Linux:
```sh
chmod +x pkg/deb/build.sh
./pkg/deb/build.sh
```

Output:
- `pkg/out/gitty_2.0.0_amd64.deb`
