# AUR publication (`gitty`)

This folder already contains:
- `PKGBUILD`
- `.SRCINFO`

## Publish to AUR

> Requires an AUR account and SSH key configured for `aur.archlinux.org`.

1. Clone AUR package repo:
   ```sh
   git clone ssh://aur@aur.archlinux.org/gitty.git
   cd gitty
   ```
2. Copy package files from this repository:
   - `pkg/aur/PKGBUILD`
   - `pkg/aur/.SRCINFO`
3. Commit and push:
   ```sh
   git add PKGBUILD .SRCINFO
   git commit -m "gitty 2.0.0"
   git push
   ```

After push, users can install with:
```sh
yay -S gitty
# or
paru -S gitty
```
