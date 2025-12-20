# Distribution Packaging Guide

This document explains how HepSW is packaged for different Linux distributions.

## Automated Releases

When you push a git tag (e.g., `v0.1.0`), GitHub Actions automatically:

1. Builds binaries for:
    - Linux AMD64
    - Linux ARM64
    - macOS AMD64 (Intel)
    - macOS ARM64 (Apple Silicon)

2. Creates distribution packages:
    - `.deb` for Debian/Ubuntu
    - `.rpm` for Fedora/RHEL/CentOS
    - `PKGBUILD` for Arch Linux

3. Publishes GitHub Release with all artifacts

4. Updates Homebrew tap (macOS)

## Triggering a Release

```bash
# Tag the commit
git tag -a v0.1.0 -m "Release v0.1.0"

# Push tag
git push origin v0.1.0

# GitHub Actions does the rest
```

## Installation Methods

### Quick Install Script

```bash
curl -sSL https://raw.githubusercontent.com/thisismeamir/hepsw/main/install.sh | bash
```

### Direct Binary Download

```bash
# Linux
curl -L https://github.com/thisismeamir/hepsw/releases/latest/download/hepsw-linux-amd64.tar.gz | tar xz
sudo mv hepsw-linux-amd64 /usr/local/bin/hepsw
chmod +x /usr/local/bin/hepsw

# macOS (Homebrew)
brew tap thisismeamir/hepsw
brew install hepsw
```

### Package Managers

**Debian/Ubuntu:**
```bash
# Download .deb from releases page
wget https://github.com/thisismeamir/hepsw/releases/download/v0.1.0/hepsw_0.1.0_amd64.deb
sudo dpkg -i hepsw_0.1.0_amd64.deb
```

**Fedora/RHEL:**
```bash
# Download .rpm from releases page
wget https://github.com/thisismeamir/hepsw/releases/download/v0.1.0/hepsw-0.1.0-1.x86_64.rpm
sudo rpm -i hepsw-0.1.0-1.x86_64.rpm
```

**Arch Linux:**
```bash
# Download PKGBUILD from releases page
wget https://github.com/thisismeamir/hepsw/releases/download/v0.1.0/hepsw-0.1.0-pkgbuild.tar.gz
tar xzf hepsw-0.1.0-pkgbuild.tar.gz
cd hepsw-0.1.0
makepkg -si
```

## Repository Structure for Packages

### APT Repository (Future)

For official APT repository, we'll need:
```
deb https://packages.hepsw.org/apt stable main
```

### YUM/DNF Repository (Future)

For official YUM repository:
```
[hepsw]
name=HepSW Repository
baseurl=https://packages.hepsw.org/rpm/
enabled=1
gpgcheck=1
```

### AUR (Arch User Repository)

Submit PKGBUILD to AUR:
```bash
# Clone AUR repository
git clone ssh://aur@aur.archlinux.org/hepsw.git
cd hepsw

# Copy PKGBUILD
cp ../PKGBUILD .

# Test build
makepkg -si

# Commit and push
git add PKGBUILD .SRCINFO
git commit -m "Update to version X.Y.Z"
git push
```

## Manual Packaging

### Building .deb

```bash
VERSION=0.1.0
ARCH=amd64

# Build binary
GOOS=linux GOARCH=amd64 go build -o hepsw ./cmd/hepsw

# Create package structure
mkdir -p hepsw_${VERSION}_${ARCH}/usr/local/bin
mkdir -p hepsw_${VERSION}_${ARCH}/DEBIAN

# Copy binary
cp hepsw hepsw_${VERSION}_${ARCH}/usr/local/bin/
chmod +x hepsw_${VERSION}_${ARCH}/usr/local/bin/hepsw

# Create control file
cat > hepsw_${VERSION}_${ARCH}/DEBIAN/control << EOF
Package: hepsw
Version: $VERSION
Architecture: $ARCH
Maintainer: Your Name <email@example.com>
Description: HEP Software Build System
EOF

# Build package
dpkg-deb --build hepsw_${VERSION}_${ARCH}
```

### Building .rpm

```bash
VERSION=0.1.0
RELEASE=1

# Build binary
GOOS=linux GOARCH=amd64 go build -o hepsw ./cmd/hepsw

# Create RPM build tree
mkdir -p ~/rpmbuild/{BUILD,RPMS,SOURCES,SPECS,SRPMS}

# Create spec file
cat > ~/rpmbuild/SPECS/hepsw.spec << EOF
Name: hepsw
Version: $VERSION
Release: $RELEASE
Summary: HEP Software Build System
License: Apache-2.0

%description
HepSW build system

%install
mkdir -p %{buildroot}/usr/local/bin
cp $PWD/hepsw %{buildroot}/usr/local/bin/

%files
/usr/local/bin/hepsw
EOF

# Build
rpmbuild -bb ~/rpmbuild/SPECS/hepsw.spec
```

## Verification

All releases include SHA256 checksums:

```bash
# Verify download
sha256sum -c hepsw-linux-amd64.tar.gz.sha256
```

## GitHub Secrets Required

For full automation, set these secrets in GitHub:

- `GITHUB_TOKEN` - Automatically provided
- `HOMEBREW_TAP_TOKEN` - Personal access token for homebrew tap repo (optional)

## Notes

- Binaries are statically compiled (CGO_ENABLED=0)
- Version is embedded at build time
- All packages are unsigned (add GPG signing for production)
- Consider code signing for macOS binaries