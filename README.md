# HepSW

**Source-First Build System for High Energy Physics Software**

HepSW is a transparent, reproducible framework for building and managing HEP software stacks on Linux systems. Instead of distributing binaries or containers, HepSW provides explicit build instructions that reconstruct software environments deterministically from source.

## Why HepSW?

HEP software stacks are complex, with deep dependency chains and strict version requirements. Traditional approaches using binary package managers or pre-built containers hide how software is actually built, making it difficult to:

- Understand what's actually installed
- Reproduce builds across different systems
- Debug build failures or compatibility issues
- Adapt software to new environments or requirements

HepSW solves this by treating **build instructions as source code**: versioned, explicit, and reproducible.

## Features

- **Source-first**: Build everything from upstream sources
- **Transparent**: Every build step is explicit and inspectable
- **Reproducible**: Same inputs → same outputs, every time
- **Flexible**: Works across Linux distributions
- **Documented**: Each package includes build guides and compatibility notes
- **Dependency-aware**: Automatic resolution of version constraints

## Quick Start

### Installation

```bash
# From source
git clone https://github.com/thisismeamir/hepsw.git
cd hepsw
make install

# Or download binary (coming soon)
# curl -L https://github.com/thisismeamir/hepsw/releases/latest/download/hepsw-linux-amd64 -o hepsw
# chmod +x hepsw
# sudo mv hepsw /usr/local/bin/
```

### Initialize a Workspace

```bash
# Create a new workspace
hepsw init ~/hep-workspace
cd ~/hep-workspace
```

This creates:
```
~/hep-workspace/
├── toolchains/     # Compilers and build tools
├── sources/        # Source code
├── builds/         # Build directories
├── install/        # Installed software
├── env/            # Environment scripts
└── logs/           # Build logs
```

### Build a Package

```bash
# Build ROOT with dependencies
hepsw build root --with-deps

# Build specific version
hepsw build geant4 --version 11.2.0

# List available packages
hepsw list

# Get package info
hepsw info root
```

### Use Installed Software

```bash
# Source environment
source $(hepsw env path root)

# Or load multiple packages
hepsw env generate --packages root,geant4 > setup.sh
source setup.sh

# Now use the software
root -b
```

## How It Works

### Package Manifests

Each package has a manifest describing how to build it:

```yaml
name: root
version: 6.30.02
description: CERN ROOT Data Analysis Framework

source:
  type: git
  url: https://github.com/root-project/root.git
  tag: v6-30-02

dependencies:
  - name: cmake
    version: ">=3.20"
  - name: python
    version: ">=3.8"

build:
  configure: cmake -S . -B build -DCMAKE_INSTALL_PREFIX=$PREFIX
  compile: cmake --build build -j$JOBS
  install: cmake --install build

environment:
  PATH: $PREFIX/bin
  LD_LIBRARY_PATH: $PREFIX/lib
  ROOTSYS: $PREFIX
```

### Build Process

1. **Fetch**: Download sources from upstream (GitHub, tarballs, etc.)
2. **Configure**: Set up build with proper flags and paths
3. **Compile**: Build in parallel with progress tracking
4. **Install**: Install to isolated prefix
5. **Environment**: Generate shell scripts to use the software

### Package Index

Available packages are maintained in [hepsw-package-index](https://github.com/thisismeamir/hepsw-package-index):

- ROOT - Data analysis framework
- Geant4 - Simulation toolkit
- Pythia8 - Event generator
- FastJet - Jet clustering
- HepMC3 - Event record
- And more...

## Commands

```bash
hepsw init <path>              # Initialize workspace
hepsw list                     # List available packages
hepsw info <package>           # Show package details
hepsw build <package>          # Build a package
hepsw build <package> --with-deps  # Build with dependencies
hepsw validate <manifest>      # Validate manifest file
hepsw env generate <package>   # Generate environment script
hepsw env path <package>       # Get path to environment script
hepsw graph <package>          # Show dependency graph
hepsw version                  # Show version
```

### Global Flags

```bash
--workspace, -w    Path to workspace (default: $HEPSW_WORKSPACE or ./hepsw-workspace)
--verbose, -v      Enable verbose output
--quiet, -q        Suppress non-essential output
--config           Config file (default: ~/.hepsw.yaml)
```

## Configuration

Create `~/.hepsw.yaml`:

```yaml
workspace: /home/user/hep-workspace

build:
  jobs: 8
  type: Release

packages:
  root:
    version: "6.30.02"
  geant4:
    version: "11.2.0"
```

## Development

### Project Structure

```
hepsw/
├── cmd/hepsw/          # CLI entry point
├── internal/
│   ├── cli/            # Command implementations
│   ├── manifest/       # Manifest parsing
│   ├── builder/        # Build orchestration
│   ├── workspace/      # Workspace management
│   └── dependencies/   # Dependency resolution
├── pkg/types/          # Public types
├── manifests/          # Package manifests (deprecated, moved to hepsw-package-index)
└── docs/               # Documentation
```

### Building from Source

```bash
# Clone repository
git clone https://github.com/thisismeamir/hepsw.git
cd hepsw

# Install dependencies
go mod download

# Build
make build

# Run tests
make test

# Install locally
make install
```

### Adding a Package

Packages are defined in [hepsw-package-index](https://github.com/thisismeamir/hepsw-package-index):

1. Fork the repository
2. Create `packages/<package-name>/manifest.yaml`
3. Add entry to `index.yaml`
4. Submit pull request

See existing packages for examples.

## Documentation

Full documentation at [docs/](./docs/):

- [Introduction](docs/00-introduction/index.md)
- [Getting Started](docs/01-basics/index.md)
- [Layout and Workflow](docs/02-layout-and-workflow/index.md)
- [Dependencies](docs/03-dependencies/index.md)
- [Build Guides](docs/04-build-guides/index.md)
- [Environments](docs/05-environments/index.md)
- [Advanced Topics](docs/06-advanced/index.md)
- [Troubleshooting](docs/07-troubleshooting/index.md)
- [Contributing](docs/08-contribution/index.md)

## Requirements

- **OS**: Linux (any distribution)
- **Tools**: git, make, gcc/clang
- **Go**: 1.21+ (for building HepSW itself)

Individual packages have their own requirements (specified in manifests).

## Comparison

### vs Binary Package Managers (apt, dnf, conda)

| Aspect | Binary Managers | HepSW |
|--------|----------------|-------|
| Transparency | Opaque binaries | Full source visibility |
| Reproducibility | Version-dependent | Deterministic from source |
| Customization | Limited | Full control |
| Cross-distro | Distribution-specific | Works everywhere |

### vs Containers (Docker, Singularity)

| Aspect | Containers | HepSW |
|--------|-----------|-------|
| What you get | Frozen filesystem | Build instructions |
| Updates | Rebuild entire image | Rebuild specific packages |
| Debugging | Black box | Inspect any layer |
| Size | GBs | Only what you need |

### vs Other Build Tools (Spack, EasyBuild)

HepSW is similar but:
- **Simpler**: Focused on HEP, not general-purpose
- **More transparent**: YAML manifests, not complex DSLs
- **Better documented**: Each package includes usage guides
- **Opinionated**: Best practices baked in

## Contributing

Contributions welcome! See [CONTRIBUTING.md](CONTRIBUTING.md).

Areas where help is needed:
- Adding package manifests
- Improving documentation
- Testing on different Linux distributions
- Bug reports and feature requests

## Support

- **Issues**: [GitHub Issues](https://github.com/thisismeamir/hepsw/issues)
- **Discussions**: [GitHub Discussions](https://github.com/thisismeamir/hepsw/discussions)
- **Documentation**: [docs/](./docs/)

## License

Apache License 2.0 - see [LICENSE](LICENSE) for details.

## Acknowledgments

HepSW builds on decades of HEP software development practices. Special thanks to the maintainers of ROOT, Geant4, and other HEP packages for their excellent upstream work.

## Status

**Early Development** - HepSW is actively being developed. APIs and manifest formats may change. Feedback welcome!

---

**HepSW**: Build HEP software the transparent way.