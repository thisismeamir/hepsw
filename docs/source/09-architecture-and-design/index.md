# Architecture and Design

This section provides an overview of the architecture and design principles behind HepSW.
HepSW is built around several key concepts that ensure transparency, reproducibility, and
adaptability in building HEP software stacks from source. We took inspiration from existing
build systems and package managers, but focused on a source-first approach that emphasizes
explicitness and clarity.

## System Overview

HepSW orchestrates the build process through several interconnected components:

```{graphviz}
digraph system_overview {
    rankdir=LR;
    node [shape=box, style="rounded,filled", fillcolor=lightblue];
    
    manifest [label="Package\nManifests", fillcolor=lightgreen];
    parser [label="Manifest\nParser"];
    resolver [label="Dependency\nResolver"];
    fetcher [label="Source\nFetcher"];
    builder [label="Build\nEngine"];
    installer [label="Install\nManager"];
    envgen [label="Environment\nGenerator"];
    
    manifest -> parser;
    parser -> resolver;
    resolver -> fetcher;
    fetcher -> builder;
    builder -> installer;
    installer -> envgen;
    
    cache [label="Source Cache", shape=cylinder, fillcolor=lightyellow];
    workspace [label="Workspace", shape=folder, fillcolor=lightyellow];
    
    fetcher -> cache [style=dashed, label="reuse"];
    builder -> workspace;
    installer -> workspace;
    envgen -> workspace;
}
```

The main components of HepSW's architecture include:

- **Manifests**: Each package is defined by a YAML document describing metadata, dependencies, build instructions, and configuration options
- **Dependency Management**: Robust system tracking relationships between packages, versions, and constraints with SAT-based conflict resolution
- **Build Engine**: Orchestrates the build process based on manifest information, retrieving sources and executing builds in controlled environments
- **Environment Management**: Tools for managing software environments including variables, paths, and configurations for seamless integration

Each of these components is discussed in detail in the following sections.

## Core Principles

HepSW's core principles are centered around the needs of physicists using HEP software and contributors to HEP projects:

```{graphviz}
digraph core_principles {
    node [shape=box, style="rounded,filled", fillcolor=lightcyan];
    
    hepsw [label="HepSW", shape=ellipse, fillcolor=lightgreen];
    
    source [label="Source First"];
    repro [label="Reproducibility"];
    adapt [label="Adaptability"];
    mod [label="Modularity"];
    
    hepsw -> source;
    hepsw -> repro;
    hepsw -> adapt;
    hepsw -> mod;
}
```

### Why These Principles Matter

**Source-first** means every build is traceable and auditable. When a build fails, you can inspect the exact commands run, the source code used, and the environment it ran in. This transparency is essential for debugging complex dependency issues common in HEP stacks. Unlike binary distributions that hide the build process, HepSW makes every step explicit and inspectable.

**Reproducibility** means the same manifest should produce identical results on different machines. HepSW achieves this through explicit dependency versions, deterministic build flags, and isolated build environments. This is critical for collaboration—if it builds on your laptop, it should build on your colleague's workstation, whether they're at a university, national lab, or working remotely.

**Adaptability** acknowledges that HEP software must work across diverse computing environments: university clusters, national labs, cloud platforms, and developer laptops. HepSW doesn't assume a specific OS distribution, filesystem layout, or centralized infrastructure like LxPlus or CVMFS. This is especially important for distributed collaborations like FCC and DUNE where team members work from institutions around the world.

**Modularity** means packages are independently buildable. You shouldn't need to build all of ROOT to get Geant4 working. Users install only what they need, and developers can test changes to individual packages without rebuilding the world. This also enables faster iteration during development and testing.

### Why Source-First?

Building software from source seems fragile and unreliable at first glance, especially when compared to binary distributions or package managers that provide pre-built packages. HepSW embraces this approach for several critical reasons:

**1. Solving Real Collaboration Problems**

Large HEP experiments like FCC and DUNE need consistent software environments across distributed teams working at different institutions worldwide. Not everyone has access to CERN infrastructure (LxPlus, CVMFS), and even those who do often need local development environments for testing and debugging. HepSW provides that consistency without requiring centralized infrastructure or institutional access.

**2. Addressing Binary Distribution Fragmentation**

Currently, most HEP software lacks coordinated binary distribution. When binaries exist, they're fragmented across conda-forge, CVMFS, experiment-specific repositories, and ad-hoc institutional builds. Different experiments maintain their own binary stacks with varying degrees of compatibility. HepSW doesn't replace these systems—it provides the *build recipes* that could unify them while letting users build locally when binaries aren't available or don't match their needs.

**3. Enabling Development and Testing**

HEP software stacks are complex and constantly evolving. Developers need to test their code against various dependency versions and configurations. A new update in Key4HEP should be testable against the entire stack immediately, ensuring that breaking changes are caught before release. HepSW guarantees a clean build and testing procedure out-of-the-box, essential for continuous integration and validation.

**4. Transparency and Control**

Relying on pre-built binaries can lead to compatibility issues, version mismatches, and hidden dependencies. When something breaks, debugging becomes difficult because you can't see how the software was built or what compilation flags were used. By building from source, HepSW ensures that users have full control over the build process, allowing them to adapt to changes in the software ecosystem and maintain reproducibility across different environments.

## Component Details

### Manifests: The Single Source of Truth

Each package in HepSW is defined by a YAML manifest that answers four fundamental questions:

1. **What is this package?** (name, version, description)
2. **Where does it come from?** (git repository, tarball URL, tag/commit)
3. **How do you build it?** (configure, compile, install commands)
4. **What does it need?** (dependencies with version constraints)

#### Minimal Manifest Example

```yaml
name: example-package
version: 1.0.0
description: Example HEP analysis package

source:
  type: git
  url: https://github.com/hep-org/example
  tag: v1.0.0

dependencies:
  - name: cmake
    version: ">=3.20"
  - name: root
    version: "^6.28.0"

build:
  configure: cmake -B build -DCMAKE_INSTALL_PREFIX=$PREFIX
  compile: cmake --build build -j$JOBS
  install: cmake --install build

environment:
  PATH: $PREFIX/bin
  LD_LIBRARY_PATH: $PREFIX/lib
```

#### Manifest with Build Options

HepSW supports configurable builds through user-selectable options. These options are discovered automatically by analyzing CMake projects (via Seemake) and can be specified during installation:

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
    optional: false

options:
  - name: builtin_openssl
    description: Build OpenSSL internally
    default: ON
    type: bool
  - name: pyroot
    description: Enable Python bindings
    default: ON
    type: bool
  - name: tmva
    description: Enable TMVA machine learning
    default: ON
    type: bool

build:
  configure: |
    cmake -S . -B build \
      -DCMAKE_INSTALL_PREFIX=$PREFIX \
      -Dbuiltin_openssl=$OPT_builtin_openssl \
      -Dpyroot=$OPT_pyroot \
      -Dtmva=$OPT_tmva
  compile: cmake --build build -j$JOBS
  install: cmake --install build

environment:
  PATH: $PREFIX/bin
  LD_LIBRARY_PATH: $PREFIX/lib
  PYTHONPATH: $PREFIX/lib
  ROOTSYS: $PREFIX
```

When installing, users can specify options:

```bash
hepsw build root --options pyroot=ON,tmva=OFF
```

For project-specific workflows, HepSW can suggest optimal configurations:

```bash
hepsw build root --project fcc
# Automatically enables FCC-relevant options
```

#### Manifest Storage

Manifests are stored separately from HepSW itself in the [hepsw-package-index](https://github.com/thisismeamir/hepsw-package-index) repository. This design allows:

- Community contributions without modifying HepSW core
- Version control of package definitions
- Independent evolution of packages and tool
- Easy forking for experiment-specific stacks

### Dependency Management: Handling Version Constraints

HepSW's dependency resolver is designed to handle the complex version relationships common in HEP software stacks:

```{graphviz}
digraph dependency_resolution {
    rankdir=TB;
    node [shape=box, style="rounded,filled", fillcolor=lightblue];
    
    request [label="User Request\n(root, geant4)", fillcolor=lightgreen];
    parse [label="Parse Manifests"];
    graph [label="Build Dependency\nGraph"];
    conflicts [label="Check Version\nConflicts", shape=diamond, fillcolor=lightyellow];
    sat [label="SAT Solver\nResolution"];
    multi [label="Build Multiple\nVersions"];
    order [label="Topological Sort\n(Build Order)"];
    build [label="Execute Builds", fillcolor=lightcoral];
    
    request -> parse;
    parse -> graph;
    graph -> conflicts;
    conflicts -> sat [label="conflicts"];
    conflicts -> order [label="compatible"];
    sat -> multi [label="no solution"];
    sat -> order [label="resolved"];
    multi -> order;
    order -> build;
}
```

#### Version Constraint Syntax

HepSW supports semantic versioning with standard constraint operators:

- `>=3.20` - Minimum version (inclusive)
- `~1.2` - Patch version updates allowed (1.2.0, 1.2.1, but not 1.3.0)
- `^2.0` - Minor version updates allowed (2.0.0, 2.1.0, but not 3.0.0)
- `==6.28.0` - Exact version required
- `>=3.0,<4.0` - Range specification

#### Resolution Strategy

The resolver operates in several phases:

1. **Graph Construction**: Parse all manifests for requested packages and recursively collect dependencies
2. **Constraint Collection**: Gather all version constraints from the dependency tree
3. **Conflict Detection**: Identify cases where different packages require incompatible versions
4. **Resolution**:
    - If all constraints are compatible: select highest compatible version for each package
    - If conflicts exist: invoke SAT solver to find compatible version set
    - If no solution exists: build multiple versions of conflicting packages in isolated prefixes
5. **Build Ordering**: Topological sort ensures dependencies are built before dependents

#### Handling Conflicting Dependencies

Example scenario:
- User wants ROOT 6.30 + Geant4 11.2
- ROOT 6.30 requires Python >=3.9,<3.12
- Geant4 11.2 requires Python >=3.11

**Resolution approach:**
1. SAT solver determines Python 3.11 satisfies both constraints
2. Single Python 3.11 installation is built
3. Both ROOT and Geant4 link against this shared Python

Example unresolvable conflict:
- Package A requires OpenSSL 1.1.x
- Package B requires OpenSSL 3.x

**Resolution approach:**
1. SAT solver detects incompatibility
2. HepSW builds both OpenSSL 1.1 and OpenSSL 3.x in separate prefixes
3. Package A links against OpenSSL 1.1, Package B against OpenSSL 3.x
4. User is warned about the dual installation

#### Optional Dependencies

Some package features are only enabled if certain dependencies are present:

```yaml
dependencies:
  - name: python
    version: ">=3.8"
    optional: false  # Required
  - name: cuda
    version: ">=11.0"
    optional: true   # Optional, enables GPU support
```

Users can control optional dependencies:

```bash
# Enable optional dependencies
hepsw build root --with-optional

# Disable specific optional dependency
hepsw build root --without cuda
```

### Build Engine: Orchestrating the Build

The build engine executes package builds in distinct, logged phases:

```{graphviz}
digraph build_phases {
    rankdir=TB;
    node [shape=box, style="rounded,filled", fillcolor=lightblue];
    
    start [label="Start Build", shape=ellipse, fillcolor=lightgreen];
    fetch [label="Fetch Phase\nDownload Sources"];
    configure [label="Configure Phase\nPrepare Build"];
    compile [label="Compile Phase\nBuild Binaries"];
    install [label="Install Phase\nCopy to Prefix"];
    register [label="Register Phase\nRecord Metadata"];
    done [label="Build Complete", shape=ellipse, fillcolor=lightgreen];
    
    fail [label="Build Failed", shape=ellipse, fillcolor=lightcoral];
    
    start -> fetch;
    fetch -> configure [label="success"];
    fetch -> fail [label="error"];
    configure -> compile [label="success"];
    configure -> fail [label="error"];
    compile -> install [label="success"];
    compile -> fail [label="error"];
    install -> register [label="success"];
    install -> fail [label="error"];
    register -> done;
}
```

#### Build Phases in Detail

**1. Fetch Phase**
- Downloads source code to `workspace/sources/<package>/<version>/`
- Supports: git repositories (tags, branches, commits), tarballs (HTTP/HTTPS), local paths
- Uses source cache to avoid redundant downloads
- Verifies checksums if provided in manifest

**2. Configure Phase**
- Prepares the build system (CMake, Autotools, etc.)
- Sets up environment variables:
    - `$PREFIX`: Installation target directory
    - `$JOBS`: Parallel build jobs (from `-j` flag)
    - `$CMAKE_PREFIX_PATH`: Paths to dependencies
    - `$PKG_CONFIG_PATH`: For pkg-config detection
- Runs configuration commands from manifest
- Creates build directory in `workspace/builds/<package>/<version>/`

**3. Compile Phase**
- Executes compilation commands
- Supports parallel builds (defaults to number of CPU cores)
- Streams output to both terminal and log file
- Can be interrupted and resumed

**4. Install Phase**
- Installs built artifacts to isolated prefix: `workspace/install/<package>/<version>/`
- Standard directory structure:
  ```
  install/<package>/<version>/
  ├── bin/          # Executables
  ├── lib/          # Libraries
  ├── include/      # Headers
  ├── share/        # Data files
  └── etc/          # Configuration
  ```
- No system directories touched (no `/usr/local`, no sudo required)

**5. Register Phase**
- Records build metadata (timestamp, options used, dependency versions)
- Generates environment scripts
- Updates package database for dependency tracking

#### Key Design Decisions

**Isolated Prefixes**
Each package/version combination gets its own installation directory. This enables:
- Multiple versions of the same package to coexist
- Clean uninstallation (just delete the directory)
- No conflicts between packages
- Easy binary distribution (tarball the prefix)

**No Sudo Required**
Everything installs to user-writable workspace directories. This is critical for:
- Cluster environments where users lack root
- Development workflows where iteration is frequent
- Reproducibility (no system state modifications)

**Parallel Builds**
Respects `-j` flag for parallel compilation:
```bash
hepsw build root -j8  # Use 8 parallel jobs
```
Defaults to `nproc` (number of CPU cores) if not specified.

**Incremental Builds**
HepSW detects when sources haven't changed and can skip redundant work:
```bash
hepsw build root --incremental  # Reuse existing build
```

**Comprehensive Logging**
Each build phase logs to `workspace/logs/<package>-<version>-<phase>.log`:
```
logs/
├── root-6.30.02-fetch.log
├── root-6.30.02-configure.log
├── root-6.30.02-compile.log
├── root-6.30.02-install.log
└── root-6.30.02-register.log
```

This makes debugging much easier—you can pinpoint exactly which phase failed and why.

### Workspace Layout: Where Everything Lives

The workspace is the central organizational structure for all HepSW operations:

```{graphviz}
digraph workspace_layout {
    rankdir=TB;
    node [shape=folder, style=filled, fillcolor=lightyellow];
    
    workspace [label="workspace/", fillcolor=lightgreen];
    
    sources [label="sources/\nDownloaded source code"];
    builds [label="builds/\nBuild directories (temporary)"];
    install [label="install/\nInstalled software"];
    env [label="env/\nEnvironment scripts"];
    logs [label="logs/\nBuild logs"];
    cache [label="cache/\nDownloaded tarballs, git repos"];
    
    workspace -> sources;
    workspace -> builds;
    workspace -> install;
    workspace -> env;
    workspace -> logs;
    workspace -> cache;
    
    pkg_sources [label="<package>/<version>/", shape=box];
    sources -> pkg_sources;
    
    pkg_builds [label="<package>/<version>/", shape=box];
    builds -> pkg_builds;
    
    pkg_install [label="<package>/<version>/\n  bin/\n  lib/\n  include/\n  share/", shape=box];
    install -> pkg_install;
}
```

#### Directory Structure

```
workspace/
├── sources/              # Downloaded source code
│   └── <package>/
│       └── <version>/    # e.g., root/6.30.02/
├── builds/               # Build directories (can be cleaned)
│   └── <package>/
│       └── <version>/
├── install/              # Installed software (isolated by package/version)
│   └── <package>/
│       └── <version>/
│           ├── bin/      # Executables
│           ├── lib/      # Libraries
│           ├── include/  # Headers
│           └── share/    # Data files, docs
├── env/                  # Generated environment scripts
│   ├── root-6.30.02.sh
│   └── geant4-11.2.0.sh
├── logs/                 # Build logs for debugging
│   ├── root-6.30.02-configure.log
│   ├── root-6.30.02-compile.log
│   └── root-6.30.02-install.log
└── cache/                # Downloaded tarballs, git repos (shared across builds)
    ├── tarballs/
    └── git/
```

#### Why This Layout?

**Multiple Versions Coexist**
`install/root/6.28.0` and `install/root/6.30.0` can exist side-by-side. Users choose which to activate via environment scripts.

**Clean Rebuilds**
Delete `builds/` directory without losing installed software. Source code in `sources/` can also be cleaned after successful builds.

**Disk Space Management**
- `cache/` persists across builds to avoid redownloading
- `sources/` and `builds/` are ephemeral (can be cleaned)
- `install/` contains only what you actually use

**Debugging Support**
Every build has its own log file. If compilation fails, inspect the exact error:
```bash
cat workspace/logs/root-6.30.02-compile.log
```

**Portability**
The entire workspace can be tarred up and moved to another machine (as long as the OS/architecture match):
```bash
tar czf my-hep-stack.tar.gz workspace/
```

### Environment Management: Making Software Usable

After building, software needs to be discoverable by the shell and other tools. HepSW generates environment scripts that configure the necessary variables:

```{graphviz}
digraph environment_generation {
    rankdir=LR;
    node [shape=box, style="rounded,filled", fillcolor=lightblue];
    
    installed [label="Installed\nPackages", fillcolor=lightgreen];
    metadata [label="Package\nMetadata"];
    template [label="Environment\nTemplate"];
    generate [label="Script\nGenerator"];
    script [label="Shell Script", shape=note, fillcolor=lightyellow];
    
    installed -> metadata;
    metadata -> generate;
    template -> generate;
    generate -> script;
}
```

#### Generated Environment Scripts

For a single package:

```bash
# workspace/env/root-6.30.02.sh
export PATH="/path/to/workspace/install/root/6.30.02/bin:$PATH"
export LD_LIBRARY_PATH="/path/to/workspace/install/root/6.30.02/lib:$LD_LIBRARY_PATH"
export PYTHONPATH="/path/to/workspace/install/root/6.30.02/lib:$PYTHONPATH"
export ROOTSYS="/path/to/workspace/install/root/6.30.02"
export CMAKE_PREFIX_PATH="/path/to/workspace/install/root/6.30.02:$CMAKE_PREFIX_PATH"
```

#### Usage Patterns

**Single Package Environment**
```bash
source $(hepsw env path root)
root -b  # Now works
```

**Combined Environment**
```bash
hepsw env generate --packages root,geant4,pythia8 > my-analysis.sh
source my-analysis.sh
```

**Named Environments**
```bash
# Create named environment for specific workflow
hepsw env create fcc-analysis --packages root,geant4,pythia8,fastjet

# Later, activate it
hepsw env activate fcc-analysis
```

**Project-Specific Environments**
HepSW can suggest packages for known projects:
```bash
hepsw env create --project fcc
# Suggests and installs: Key4HEP stack, FCCSW, Gaudi, etc.

hepsw env create --project dune
# Suggests and installs: LArSoft, ROOT, Geant4, etc.
```

#### Design Philosophy: Simple is Better

We use plain shell scripts instead of complex module systems (Environment Modules, Lmod) because:

1. **Transparency**: You can read and understand exactly what the script does
2. **Debuggability**: Easy to trace issues (`set -x`, inspect variables)
3. **Portability**: Works on any shell (bash, zsh, dash)
4. **No Dependencies**: No additional tools required
5. **Integration**: Advanced users can integrate with existing module systems if desired

If your institution uses module systems, HepSW can generate module files:
```bash
hepsw env generate --format modulefile --package root > root/6.30.02
```

## Comparison with Other Tools

HepSW is not the first tool to tackle build automation in scientific computing. Understanding how it relates to existing tools helps clarify its design choices:

### vs Spack

[Spack](https://spack.io/) is a mature, general-purpose package manager for HPC:

**Similarities:**
- Both build from source
- Both support multiple versions
- Both handle complex dependency graphs
- Both generate environment modules

**Differences:**

| Aspect | Spack | HepSW |
|--------|-------|-------|
| **Scope** | General HPC (10,000+ packages) | HEP-focused (~100 packages) |
| **Manifest Language** | Python DSL | YAML (declarative) |
| **Learning Curve** | Steep (Python API, complex syntax) | Gentle (readable YAML) |
| **Build Variants** | Comprehensive but complex | Simple, user-friendly options |
| **Documentation** | Package-centric | Build guides + usage tutorials |
| **Target Audience** | HPC system administrators | HEP physicists and developers |

**When to use Spack:**
- You need non-HEP software (compilers, MPI, etc.)
- You're a system administrator managing a cluster
- You need advanced features (compiler bootstrapping, microarchitecture optimization)

**When to use HepSW:**
- You're a HEP physicist setting up your analysis environment
- You want transparency and simplicity
- You need HEP-specific optimizations and workflows
- You want documentation that explains *how* to use the software, not just how to build it

### vs EasyBuild

[EasyBuild](https://easybuild.io/) is another HPC-focused build framework:

**Similarities:**
- Source-based builds
- Reproducibility focus
- Module file generation

**Differences:**

| Aspect | EasyBuild | HepSW |
|--------|-----------|-------|
| **Configuration** | Python-based "easyconfigs" | YAML manifests |
| **Philosophy** | System-wide installations | User-space workspaces |
| **Toolchains** | Rigid toolchain definitions | Flexible, minimal constraints |
| **Target Environment** | Clusters with shared filesystem | Local workstations + clusters |

### vs Conda/Mamba

[Conda](https://docs.conda.io/) is a popular binary package manager:

**Fundamental Difference:**
Conda distributes pre-built binaries; HepSW builds from source.

**When Conda Works Well:**
- Python-heavy workflows
- Standard packages with binaries on conda-forge
- Quick setup without compilation

**When HepSW is Better:**
- Latest versions not yet on conda-forge
- Custom build configurations needed
- Source-level debugging required
- Binary compatibility issues with your system
- You want to understand how software is built

**Complementary Use:**
Many users combine them:
```bash
conda create -n hep python=3.11 cmake numpy  # Basic tools
conda activate hep
hepsw build root geant4  # HEP-specific software from source
```

### vs Nix

[Nix](https://nixos.org/) provides purely functional package management:

**Similarities:**
- Reproducible builds
- Multiple versions coexist
- Declarative configuration

**Differences:**

| Aspect | Nix | HepSW |
|--------|-----|-------|
| **Model** | Functional, immutable | Conventional, mutable |
| **Learning Curve** | Very steep | Gentle |
| **OS Integration** | Deep (can replace OS package manager) | Shallow (workspace-based) |
| **Adoption Barrier** | High (requires buying into Nix philosophy) | Low (works like traditional tools) |

### vs Containers (Docker, Singularity)

**Fundamental Difference:**
Containers package entire environments; HepSW builds software.

**When Containers Work Well:**
- Production: reproducible deployment of complete analysis chains
- Sharing: distribute entire environment to collaborators
- Isolation: completely separate from host system

**When HepSW is Better:**
- Development: iterative builds, testing changes
- Flexibility: mix and match versions, custom builds
- Transparency: inspect and modify any component
- Size: install only what you need (not multi-GB images)

**Complementary Use:**
Build software with HepSW, then package in container for production:
```bash
# Development: use HepSW
hepsw build root geant4 my-analysis

# Production: containerize the workspace
FROM ubuntu:22.04
COPY workspace /opt/hep-workspace
RUN echo 'source /opt/hep-workspace/env/my-analysis.sh' >> ~/.bashrc
```

### Summary: HepSW's Unique Position

HepSW occupies a specific niche:

```{graphviz}
digraph tool_comparison {
    rankdir=TB;
    node [shape=box, style="rounded,filled"];
    
    subgraph cluster_general {
        label="General-Purpose";
        style=filled;
        fillcolor=lightgray;
        spack [label="Spack\n(HPC)", fillcolor=lightblue];
        easybuild [label="EasyBuild\n(HPC)", fillcolor=lightblue];
        nix [label="Nix\n(Universal)", fillcolor=lightblue];
    }
    
    subgraph cluster_hep {
        label="HEP-Specific";
        style=filled;
        fillcolor=lightgreen;
        hepsw [label="HepSW\n(HEP Physics)", fillcolor=yellow];
        containers [label="Containers\n(Production)", fillcolor=lightcyan];
    }
    
    subgraph cluster_binary {
        label="Binary Distribution";
        style=filled;
        fillcolor=lightyellow;
        conda [label="Conda\n(Python/Data Science)", fillcolor=lightblue];
        cvmfs [label="CVMFS\n(HEP Binary Cache)", fillcolor=lightblue];
    }
}
```

**HepSW's sweet spot:**
- HEP physicists who need source builds
- Local development environments
- Understanding and transparency
- Gentle learning curve
- Project-specific optimization (FCC, DUNE, etc.)

It's not trying to replace Spack for HPC administrators, nor Conda for Python environments, nor containers for production. It's designed specifically for HEP developers and users who want control, transparency, and simplicity.

## Developer Architecture

This section is for contributors working on HepSW itself.

### Code Organization

```
hepsw/
├── cmd/hepsw/              # CLI entry point
│   └── main.go            # Cobra command setup
├── internal/              # Internal packages (not importable)
│   ├── cli/              # Command implementations
│   │   ├── build.go
│   │   ├── env.go
│   │   ├── init.go
│   │   └── list.go
│   ├── manifest/         # Manifest parsing and validation
│   │   ├── parser.go
│   │   ├── validator.go
│   │   └── types.go
│   ├── builder/          # Build orchestration
│   │   ├── engine.go
│   │   ├── phases.go
│   │   └── logger.go
│   ├── resolver/         # Dependency resolution
│   │   ├── graph.go
│   │   ├── sat.go
│   │   └── version.go
│   ├── workspace/        # Workspace management
│   │   ├── layout.go
│   │   └── cache.go
│   └── environment/      # Environment script generation
│       ├── generator.go
│       └── templates.go
├── pkg/                  # Public packages (importable)
│   └── types/           # Shared types
└── docs/                # Documentation (Sphinx + Markdown)
```

### Key Abstractions

```{graphviz}
digraph key_abstractions {
    rankdir=TB;
    node [shape=box, style="rounded,filled", fillcolor=lightblue];
    
    manifest [label="Manifest\nYAML definition", fillcolor=lightgreen];
    package [label="Package\nParsed manifest + metadata"];
    depgraph [label="DependencyGraph\nPackages + edges"];
    buildplan [label="BuildPlan\nOrdered build sequence"];
    builder [label="Builder\nExecutes build phases"];
    workspace [label="Workspace\nManages filesystem layout"];
    
    manifest -> package [label="parse"];
    package -> depgraph [label="resolve"];
    depgraph -> buildplan [label="sort"];
    buildplan -> builder [label="execute"];
    builder -> workspace [label="uses"];
}
```

**Manifest** → Raw YAML representation
**Package** → Parsed and validated package definition
**DependencyGraph** → DAG of packages with version constraints
**BuildPlan** → Topologically sorted list of builds to execute
**Builder** → Executes build phases for a single package
**Workspace** → Manages directories, caching, paths

### Extension Points

HepSW is designed to be extensible:

**Custom Source Fetchers**
Add support for new source types:
```go
type Fetcher interface {
    Fetch(source Source, destDir string) error
    IsCached(source Source) bool
}
```

**Build System Support**
Add support for non-CMake builds:
```go
type BuildSystem interface {
    Configure(pkg Package, buildDir string) error
    Compile(pkg Package, buildDir string) error
    Install(pkg Package, buildDir, installDir string) error
}
```

**Dependency Resolvers**
Implement custom resolution strategies:
```go
type Resolver interface {
    Resolve(packages []Package) (BuildPlan, error)
}
```

### Testing Strategy

```{graphviz}
digraph testing_strategy {
    rankdir=TB;
    node [shape=box, style="rounded,filled", fillcolor=lightblue];
    
    unit [label="Unit Tests\nIndividual functions"];
    integration [label="Integration Tests\nComponent interactions"];
    e2e [label="End-to-End Tests\nComplete workflows"];
    
    unit -> integration [label="builds on"];
    integration -> e2e [label="builds on"];
    
    fixtures [label="Test Fixtures\nMinimal CMake projects", shape=cylinder, fillcolor=lightyellow];
    
    e2e -> fixtures [style=dashed];
}
```

**Unit Tests**: Test individual functions (manifest parsing, version comparison)
**Integration Tests**: Test component interactions (resolver + builder)
**End-to-End Tests**: Test complete workflows with minimal real packages

Test fixtures include simple CMake projects that mimic HEP software structure but compile in seconds.

### Performance Considerations

**Parallel Dependency Builds**
Independent packages can be built in parallel:
```bash
hepsw build root geant4 pythia8 -j4  # Build 4 packages concurrently
```

Implementation uses goroutines with dependency tracking to maximize parallelism.

**Incremental Compilation**
HepSW detects when source hasn't changed and can reuse build artifacts:
- Hash source directory
- Compare with previous build hash
- Skip configure/compile if unchanged

**Caching Strategy**
- Source cache: Persist git clones and downloaded tarballs
- Build cache: Optionally preserve build directories
- Binary cache (future): Share built packages across workspaces

## Future Enhancements

### Planned Features

**Binary Caching**
Build once, share across machines:
```bash
# Machine A
hepsw build root --cache-upload

# Machine B (same OS/arch)
hepsw build root --cache-fetch  # Skip compilation
```

**Cross-Compilation**
Build for different architectures:
```bash
hepsw build root --target aarch64-linux-gnu
```

**Distributed Builds**
Offload compilation to build servers:
```bash
hepsw build root --distributed
```

**Environment Snapshots**
Capture exact environment for reproducibility:
```bash
hepsw env snapshot > my-analysis.lock
# Later, reproduce exactly
hepsw env restore my-analysis.lock
```

**CI/CD Integration**
GitHub Actions and GitLab CI templates:
```yaml
- uses: hepsw/setup@v1
  with:
    packages: root geant4
```

### Research Directions

- **Automatic dependency discovery** from CMakeLists.txt (deeper Seemake integration)
- **Machine learning for build optimization** (predict optimal compiler flags)
- **Provenance tracking** (record exact commit hashes, build machine details)
- **Incremental environment updates** (change one package without rebuilding everything)

## Contributing to Architecture

The architecture is not set in stone. We welcome discussions about:

- Alternative dependency resolution strategies
- Better workspace organization schemes
- Performance optimizations
- New use cases that don't fit the current model

Please open an issue on GitHub to discuss architectural changes before implementing them.

---

**Next:** [Workspace Layout and Workflow](02-layout-and-workflow.md)