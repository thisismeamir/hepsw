# Building

Building is the core functionality of HepSW. It provides a complete workflow for analyzing, debugging, building, and installing software packages to ensure consistency, reproducibility, and reliability across different systems and environments.

## Overview of the Build Process

After fetching a package with `hepsw fetch`, the source code is available in `~/.hepsw/sources/<package-name>/<version>/src` along with its manifest and a generated `build.yml` file. The build process transforms this source code into installed, usable software.

The HepSW build workflow follows these stages:

1. **Analysis**: Understanding the package structure and requirements
2. **Evaluation**: Checking for potential issues and missing dependencies
3. **Simulation**: Walking through the build steps without execution
4. **Building**: Compiling and linking the source code
5. **Testing**: Running package tests (if available)
6. **Installation**: Placing binaries and libraries in the install directory

## Quick Start: Building a Package

For a simple, straightforward build:

```bash
# Fetch the package first
hepsw fetch root

# Build with default settings
hepsw build root
```

For a more careful approach with validation:

```bash
# Fetch the package
hepsw fetch pythia8

# Evaluate for potential issues
hepsw evaluate pythia8

# Simulate the build process
hepsw walk pythia8

# If everything looks good, build
hepsw build pythia8
```

## Pre-Build Commands

Before building, HepSW provides several commands to understand and validate the build process.

### Understanding a Package: `hepsw whatis`

The `whatis` command provides a comprehensive overview of a package, including its description, dependencies, build options, and recipe structure.

```bash
hepsw whatis <package-name>
```

Example output:

```text
Package: root
Version: 6.30.02
Description: The ROOT data analysis framework

Source:
  Type: git
  URL: https://github.com/root-project/root.git
  Tag: v6-30-02

Dependencies:
  Build:
    - cmake >=3.16
    - gcc >=9.3
    - python >=3.8 (optional, for: with-python)
    - qt5 >=5.12 (optional, for: with-gui)
    
  Runtime:
    - python >=3.8 (for: with-python)
    - qt5 >=5.12 (for: with-gui)

Build Options:
  - with-python: Enable Python bindings
  - with-gui: Enable GUI components
  - with-ssl: Enable SSL support

Recipe Steps:
  configure: 3 steps
  build: 1 step
  test: 1 step
  install: 1 step

Build Directory: ~/.hepsw/builds/root/6.30.02
Install Directory: ~/.hepsw/install/root/6.30.02
```

**Options:**
- `--version <version>` or `-v <version>`: Show information for a specific version
- `--verbose`: Show detailed recipe steps
- `--deps-tree`: Display dependency tree visually

### Evaluating a Package: `hepsw evaluate`

The `evaluate` command performs a thorough analysis of the manifest and build environment to identify potential issues before building.

```bash
hepsw evaluate <package-name> [version]
```

This command checks for:

- **Manifest Validity**: Ensures the manifest structure is correct and all required fields are present
- **Dependency Availability**: Verifies that all required dependencies are available or fetched
- **Version Conflicts**: Detects incompatible version requirements across dependencies
- **Toolchain Compatibility**: Checks that required compilers and build tools are available
- **Build Target Support**: Validates that the current system matches supported build targets
- **Environment Variables**: Ensures all required environment variables are set or can be set
- **Disk Space**: Estimates required disk space for build and installation
- **Recipe Validity**: Validates recipe step syntax and command availability

Example output:

```text
Evaluating: pythia8 version 8.310

✓ Manifest structure valid
✓ Source accessible
✓ Build toolchain available:
  - cmake 3.22.1 (required: >=3.15)
  - gcc 11.2.0 (required: >=9.0)
  
⚠ Dependencies:
  ✓ hepmc3 2.13.2 (required: >=2.13)
  ✗ lhapdf 6.5.4 (required: >=6.3) - NOT FOUND
  
✓ Build target supported: linux-x86_64
✓ Estimated build space: ~1.2 GB
✓ Recipe steps valid (4 steps total)

Issues Found: 1
- Missing dependency: lhapdf >=6.3
  
Recommendation: Run 'hepsw fetch lhapdf' before building
```

**Options:**
- `--version <version>` or `-v <version>`: Evaluate a specific version
- `--fix`: Automatically fetch missing dependencies
- `--ignore-warnings`: Only show critical errors
- `--deps-tree`: Show full dependency evaluation tree

**Exit Codes:**
- `0`: No issues found, safe to build
- `1`: Warnings found, may build but could fail
- `2`: Critical errors found, build will fail

### Simulating a Build: `hepsw walk`

The `walk` command simulates the build process step-by-step without actually executing commands. This is invaluable for understanding what will happen during the build and debugging recipe issues.

```bash
hepsw walk <package-name> [version]
```

This command shows:

- Each recipe step in order
- Variables and their interpolated values
- Commands that would be executed
- Environment variables that would be set
- Conditional steps and whether they would run
- Expected output directories

Example output:

```text
Walking through build process: root version 6.30.02
Source: ~/.hepsw/sources/root/6.30.02/src
Build: ~/.hepsw/builds/root/6.30.02
Install: ~/.hepsw/install/root/6.30.02

════════════════════════════════════════════════════════════
STAGE: configure
════════════════════════════════════════════════════════════

[1/3] Set build variables
  Type: set
  Variables:
    BUILD_DIR = ~/.hepsw/builds/root/6.30.02
    JOBS = 8
    TYPE = Release

[2/3] Create build directory
  Type: command
  Command: mkdir -p ~/.hepsw/builds/root/6.30.02 && cd ~/.hepsw/builds/root/6.30.02

[3/3] Configure ROOT
  Type: cmake
  Arguments:
    -DCMAKE_BUILD_TYPE=Release
    -DCMAKE_INSTALL_PREFIX=~/.hepsw/install/root/6.30.02
    -Dpython3=ON
    -Dssl=OFF
    -Dqt5web=OFF
  Note: Options with-python enabled, with-ssl disabled, with-gui disabled

════════════════════════════════════════════════════════════
STAGE: build
════════════════════════════════════════════════════════════

[1/1] Build ROOT
  Type: cmake
  Target: build
  Parallel: 8
  Command equivalent: cmake --build . --parallel 8

════════════════════════════════════════════════════════════
STAGE: test
════════════════════════════════════════════════════════════

[1/1] Run ROOT tests
  Type: command
  Command: cd ~/.hepsw/builds/root/6.30.02 && ctest --output-on-failure -j8

════════════════════════════════════════════════════════════
STAGE: install
════════════════════════════════════════════════════════════

[1/1] Install ROOT
  Type: cmake
  Target: install
  Command equivalent: cmake --build . --target install

════════════════════════════════════════════════════════════
Summary
════════════════════════════════════════════════════════════
Total steps: 6
Estimated time: 15-30 minutes (based on similar builds)
Log file: ~/.hepsw/logs/root-6.30.02-walk.log
```

**Options:**
- `--version <version>` or `-v <version>`: Walk through a specific version
- `--stage <stage>`: Only show a specific stage (configure, build, test, install)
- `--with <option>`: Enable a build option (e.g., `--with python`)
- `--without <option>`: Disable a build option
- `--show-env`: Display all environment variables at each step
- `--export-script <path>`: Export the build process as a bash script

## Building Packages

### Basic Build Command

Once you've verified the package is ready to build, use the `build` command:

```bash
hepsw build <package-name> [version]
```

This command executes the full build process:

1. Sets up the build environment
2. Executes configure steps
3. Compiles the source code
4. Runs tests (if enabled)
5. Installs to the install directory
6. Updates the workspace state in `hepsw.yaml`
7. Generates environment setup scripts

Example:

```bash
hepsw build geant4
```

The build process creates:
- Build artifacts in `~/.hepsw/builds/geant4/<version>/`
- Installed files in `~/.hepsw/install/geant4/<version>/`
- Build log in `~/.hepsw/logs/geant4-<version>-build.log`
- Environment script in `~/.hepsw/env/geant4-<version>.sh`

### Build Options and Flags

**Version Selection:**
```bash
# Build a specific version
hepsw build root --version 6.28.06

# Build the latest version
hepsw build root --latest
```

**Build Extensions:**
```bash
# Enable specific build options
hepsw build root --with python --with gui

# Disable specific options
hepsw build root --without ssl

# List available options
hepsw build root --list-options
```

**Build Configuration:**
```bash
# Set number of parallel jobs
hepsw build pythia8 --jobs 16
hepsw build pythia8 -j 16

# Use a specific build type
hepsw build root --build-type Debug
hepsw build root --build-type Release
hepsw build root --build-type RelWithDebInfo

# Custom install prefix (advanced)
hepsw build mylib --prefix /opt/custom/mylib
```

**Build Stages:**
```bash
# Run only specific stages
hepsw build root --only configure
hepsw build root --only build
hepsw build root --only test
hepsw build root --only install

# Skip specific stages
hepsw build root --skip test
hepsw build root --skip-tests  # alias
```

**Clean Builds:**
```bash
# Clean before building
hepsw build root --clean

# Rebuild from scratch (removes build directory)
hepsw build root --rebuild

# Force rebuild even if already built
hepsw build root --force
```

**Dependency Handling:**
```bash
# Build with dependencies
hepsw build root --with-deps

# Build dependencies to a specific depth
hepsw build root --deps-depth 2

# Only build dependencies, not the package itself
hepsw build root --deps-only
```

**Logging and Output:**
```bash
# Verbose output
hepsw build root --verbose
hepsw build root -v

# Quiet mode (only show errors)
hepsw build root --quiet
hepsw build root -q

# Save log to custom location
hepsw build root --log-file /path/to/custom.log

# Show real-time progress
hepsw build root --progress
```

**Advanced Options:**
```bash
# Continue from a failed build
hepsw build root --resume

# Use a custom manifest
hepsw build --manifest /path/to/custom-manifest.yaml

# Set custom environment variables
hepsw build root --env CC=clang --env CXX=clang++

# Dry run (equivalent to walk but with build command)
hepsw build root --dry-run
```

### Build Process in Detail

When you run `hepsw build <package-name>`, here's what happens internally:

**1. Pre-build Phase:**
- Validates the package manifest
- Checks for required dependencies
- Verifies toolchain availability
- Sets up build and install directories
- Initializes logging

**2. Environment Setup:**
- Loads toolchain environment
- Sets required environment variables from manifest
- Configures paths (PATH, LD_LIBRARY_PATH, etc.)
- Applies user-specified environment overrides

**3. Configure Stage:**
- Executes all steps in `recipe.configure`
- Interpolates variables
- Runs CMake, configure scripts, or custom configuration
- Validates configuration success

**4. Build Stage:**
- Executes all steps in `recipe.build`
- Compiles source code with specified parallelism
- Monitors build progress and logs output
- Handles build errors and provides diagnostic information

**5. Test Stage (if not skipped):**
- Executes all steps in `recipe.test`
- Runs test suites
- Reports test results
- Optionally fails build on test failure (configurable)

**6. Install Stage:**
- Executes all steps in `recipe.install`
- Copies binaries, libraries, and headers to install directory
- Sets proper permissions
- Generates package metadata

**7. Post-build Phase:**
- Updates workspace state in `hepsw.yaml`
- Generates environment setup script
- Creates package metadata file
- Cleans up temporary files (optional)
- Reports build summary

### Build Output Structure

After a successful build, the workspace contains:

```text
~/.hepsw/
├── builds/
│   └── <package-name>/
│       └── <version>/
│           ├── CMakeCache.txt       # CMake configuration
│           ├── Makefile             # Generated makefiles
│           ├── CMakeFiles/          # CMake internals
│           └── ...                  # Other build artifacts
├── install/
│   └── <package-name>/
│       └── <version>/
│           ├── bin/                 # Executables
│           ├── lib/                 # Libraries
│           ├── include/             # Headers
│           ├── share/               # Data files
│           └── .hepsw-metadata.yaml # Package metadata
├── logs/
│   └── <package-name>-<version>-build-YYYYMMDD-HHMMSS.log
└── env/
    └── <package-name>-<version>.sh  # Environment setup script
```

The `.hepsw-metadata.yaml` file contains:

```yaml
name: root
version: 6.30.02
built_at: 2025-01-07T14:30:45Z
built_on: linux-x86_64
build_type: Release
toolchain:
  compiler: gcc 11.2.0
  cmake: 3.22.1
options:
  with-python: true
  with-gui: false
  with-ssl: false
dependencies:
  - python: 3.9.7
  - hepmc3: 2.13.2
install_prefix: /home/user/.hepsw/install/root/6.30.02
build_time: 1234  # seconds
manifest_checksum: sha256:abcdef...
```

## Build Examples

### Example 1: Simple Library Build

```bash
# Fetch and build a simple library
hepsw fetch yaml-cpp
hepsw evaluate yaml-cpp
hepsw build yaml-cpp

# Check the installation
ls ~/.hepsw/install/yaml-cpp/*/
```

### Example 2: Building with Options

```bash
# Fetch ROOT and build with Python support
hepsw fetch root
hepsw build root --with python --without gui --jobs 8

# Verify Python bindings are installed
source ~/.hepsw/env/root-6.30.02.sh
python3 -c "import ROOT; print(ROOT.__version__)"
```

### Example 3: Building with Dependencies

```bash
# Build Geant4 with all dependencies
hepsw fetch geant4
hepsw evaluate geant4 --fix  # Automatically fetch missing deps
hepsw build geant4 --with-deps --jobs 12

# This will build: clhep, expat, xerces-c, and then geant4
```

### Example 4: Rebuilding After Source Changes

```bash
# Make changes to the source
cd ~/.hepsw/sources/mylib/1.0.0/src
# ... edit files ...

# Rebuild from the modified source
hepsw build mylib --rebuild
```

### Example 5: Debug Build for Development

```bash
# Build in debug mode with verbose output
hepsw build mylib --build-type Debug --verbose --skip-tests

# The debug symbols are now available for gdb
gdb ~/.hepsw/install/mylib/1.0.0/bin/myapp
```

### Example 6: Third-Party Package Build

```bash
# Fetch a third-party manifest
hepsw fetch --third-party /path/to/my-custom-lib.yaml

# Build it (must explicitly specify since it's third-party)
hepsw build my-custom-lib --third-party
```

## Troubleshooting Build Issues

### Common Build Failures

**1. Missing Dependencies:**
```text
Error: Package 'openssl' not found
```
**Solution:**
```bash
hepsw fetch openssl
hepsw build openssl
hepsw build mypackage
```

**2. Compiler Not Found:**
```text
Error: C++ compiler not found (required: gcc >=9.0)
```
**Solution:**
Install the required compiler or update your toolchain:
```bash
# On Ubuntu/Debian
sudo apt install gcc-11 g++-11

# Update alternatives
sudo update-alternatives --install /usr/bin/gcc gcc /usr/bin/gcc-11 100
```

**3. Configuration Failed:**
```text
Error in configure stage: CMake configuration failed
```
**Solution:**
Check the detailed log:
```bash
hepsw build mypackage --verbose
# or
cat ~/.hepsw/logs/mypackage-*-build.log
```

**4. Out of Disk Space:**
```text
Error: No space left on device
```
**Solution:**
Clean old builds:
```bash
hepsw clean --builds
hepsw clean --old-builds  # Keeps only the latest version
```

**5. Build Timeout:**
```text
Error: Build timed out after 3600 seconds
```
**Solution:**
Increase timeout or reduce parallelism:
```bash
hepsw build mypackage --timeout 7200 --jobs 4
```

### Debugging Build Failures

**Step 1: Check the evaluation:**
```bash
hepsw evaluate mypackage --verbose
```

**Step 2: Simulate the build:**
```bash
hepsw walk mypackage --show-env
```

**Step 3: Inspect the log:**
```bash
# View the latest log
cat $(ls -t ~/.hepsw/logs/mypackage-*-build.log | head -1)

# Search for errors
grep -i "error" ~/.hepsw/logs/mypackage-*-build.log
```

**Step 4: Try a clean rebuild:**
```bash
hepsw build mypackage --rebuild --verbose
```

**Step 5: Build manually for debugging:**
```bash
# Export the build script
hepsw walk mypackage --export-script build-debug.sh

# Edit and run manually
bash -x build-debug.sh
```

### Getting Help

If you encounter persistent build issues:

1. Check the package manifest for known issues or special requirements
2. Review the build log carefully for specific error messages
3. Try building with `--verbose` flag for more detailed output
4. Check if the package has upstream build documentation
5. Report issues to the HepSW issue tracker with:
    - Package name and version
    - Build command used
    - Full build log
    - System information (`uname -a`, compiler versions, etc.)

## Advanced Topics

### Custom Build Scripts

You can create custom build workflows by combining HepSW commands:

```bash
#!/bin/bash
# custom-root-build.sh

set -e  # Exit on error

echo "Building custom ROOT environment..."

# Fetch dependencies
hepsw fetch python cmake

# Build dependencies
hepsw build python --jobs 8
hepsw build cmake --jobs 8

# Fetch and configure ROOT
hepsw fetch root --version 6.30.02
hepsw evaluate root

# Build with custom options
hepsw build root \
  --with python \
  --with gui \
  --without ssl \
  --build-type RelWithDebInfo \
  --jobs 16 \
  --verbose

# Source the environment
source ~/.hepsw/env/root-6.30.02.sh

# Verify installation
root -b -q -e 'cout << "ROOT " << gROOT->GetVersion() << " is ready!" << endl;'

echo "Build complete!"
```

### Batch Building Multiple Packages

```bash
#!/bin/bash
# Build entire analysis stack

packages=(
  "cmake"
  "boost"
  "python"
  "hepmc3"
  "pythia8"
  "root"
  "geant4"
)

for pkg in "${packages[@]}"; do
  echo "Building $pkg..."
  hepsw fetch "$pkg" || continue
  hepsw build "$pkg" --with-deps --jobs 8 || {
    echo "Failed to build $pkg"
    exit 1
  }
done

echo "All packages built successfully!"
```

### Integration with CI/CD

HepSW can be integrated into continuous integration pipelines:

```yaml
# .gitlab-ci.yml example
build-hep-stack:
  stage: build
  script:
    - hepsw init
    - hepsw fetch --deps-depth 5 myanalysis
    - hepsw build myanalysis --with-deps --jobs 8
  artifacts:
    paths:
      - ~/.hepsw/install/
      - ~/.hepsw/env/
    expire_in: 1 week
  cache:
    paths:
      - ~/.hepsw/sources/
```

### Performance Optimization

**Maximize Build Speed:**
```bash
# Use all available cores
hepsw build mypackage --jobs $(nproc)

# Use ccache if available
hepsw build mypackage --env CC="ccache gcc" --env CXX="ccache g++"

# Skip tests during development
hepsw build mypackage --skip-tests
```

**Minimize Disk Usage:**
```bash
# Clean after installation
hepsw build mypackage --clean-after-install

# Use shared dependencies
hepsw config set shareableDeps true
```

## Build Command Reference

### `hepsw build`

Build a package from source.

**Usage:**
```bash
hepsw build [options] <package-name> [version]
```

**Common Options:**
- `--version, -v`: Specify package version
- `--with <option>`: Enable build option
- `--without <option>`: Disable build option
- `--jobs, -j`: Number of parallel jobs
- `--build-type`: Build type (Debug, Release, RelWithDebInfo)
- `--verbose`: Verbose output
- `--quiet`: Quiet mode
- `--clean`: Clean before building
- `--rebuild`: Rebuild from scratch
- `--with-deps`: Build dependencies
- `--skip-tests`: Skip test stage

**Examples:**
```bash
hepsw build root
hepsw build root --version 6.28.06 --with python --jobs 8
hepsw build pythia8 --rebuild --verbose
hepsw build geant4 --with-deps --skip-tests
```

### `hepsw evaluate`

Evaluate a package for build readiness.

**Usage:**
```bash
hepsw evaluate [options] <package-name> [version]
```

**Options:**
- `--version, -v`: Specify version
- `--fix`: Automatically fetch missing dependencies
- `--ignore-warnings`: Only show critical errors
- `--deps-tree`: Show dependency evaluation tree

### `hepsw walk`

Simulate a build process without execution.

**Usage:**
```bash
hepsw walk [options] <package-name> [version]
```

**Options:**
- `--version, -v`: Specify version
- `--stage <stage>`: Show only specific stage
- `--with <option>`: Enable build option
- `--without <option>`: Disable build option
- `--show-env`: Display environment variables
- `--export-script <path>`: Export as bash script

### `hepsw whatis`

Display package information.

**Usage:**
```bash
hepsw whatis [options] <package-name>
```

**Options:**
- `--version, -v`: Show specific version
- `--verbose`: Show detailed recipe steps
- `--deps-tree`: Display dependency tree

## Further Reading

### Build Systems and Tools

- [CMake Official Documentation](https://cmake.org/documentation/) - Comprehensive CMake guide
- [GNU Make Manual](https://www.gnu.org/software/make/manual/) - Complete Make documentation
- [Ninja Build System](https://ninja-build.org/) - Fast build system (alternative to Make)
- [Meson Build System](https://mesonbuild.com/) - Modern build system

### Compilation and Toolchains

- [GCC Documentation](https://gcc.gnu.org/onlinedocs/) - GNU Compiler Collection
- [Clang Documentation](https://clang.llvm.org/docs/) - LLVM Clang compiler
- [An Introduction to GCC](https://www.network-theory.co.uk/docs/gccintro/) - GCC guide

### HEP Software Development

- [HSF Software Training](https://hsf-training.github.io/hsf-training-cmake-webpage/) - HEP Software Foundation CMake training
- [HSF Knowledge Base](https://hepsoftwarefoundation.org/knowledge_base.html) - HEP software development resources
- [CERN Software Development](https://ep-dep-sft.web.cern.ch/) - Software tools used at CERN
- [LCG Releases](https://lcginfo.cern.ch/) - CERN library and application releases

### Package Management

- [Spack](https://spack.io/) - Package manager for HPC (widely used in HEP)
- [Conda Documentation](https://docs.conda.io/) - Cross-platform package manager
- [vcpkg](https://vcpkg.io/) - C/C++ library manager

### Best Practices

- [Modern CMake](https://cliutils.gitlab.io/modern-cmake/) - Best practices for CMake
- [Professional CMake](https://crascit.com/professional-cmake/) - Comprehensive CMake book
- [The Architecture of Open Source Applications](https://aosabook.org/) - Learn from real-world build systems