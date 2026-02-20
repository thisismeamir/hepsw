# HepSW CLI API Reference

**Version:** 1.0  
**Last Updated:** February 2026

---

## Table of Contents

1. [Introduction](#introduction)
2. [Workspace Management](#workspace-management)
3. [Package Discovery](#package-discovery)
4. [Source Management](#source-management)
5. [Build Operations](#build-operations)
6. [Environment Management](#environment-management)
7. [Configuration Management](#configuration-management)
8. [Maintenance & Utilities](#maintenance--utilities)
9. [Global Flags](#global-flags)

---

## Introduction

HepSW is a command-line tool for managing High Energy Physics (HEP) software packages. This document provides a complete reference of all CLI commands, their purposes, available flags, and usage examples.

### Command Structure

All HepSW commands follow this general structure:

```bash
hepsw [global-flags] <command> [command-flags] [arguments]
```

### Getting Help

For any command, you can use the `--help` flag to get detailed information:

```bash
hepsw --help
hepsw <command> --help
```

---

## Workspace Management

Commands for initializing and managing the HepSW workspace.

### `hepsw init`

**Purpose:** Initialize a new HepSW workspace with the required directory structure and configuration files.

**Usage:**
```bash
hepsw init [flags]
```

**What it does:**
1. Creates the `~/.hepsw/` directory structure
2. Generates the default `hepsw.yaml` configuration file
3. Clones the package index from the upstream repository
4. Sets up subdirectories: `toolchains/`, `sources/`, `builds/`, `install/`, `env/`, `logs/`, `index/`, `third-party/`

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--workspace <path>` | `-w` | Specify custom workspace location (default: `~/.hepsw`) |
| `--index-url <url>` | `-i` | Use a custom package index repository URL |
| `--skip-index` | | Skip cloning the package index |
| `--force` | `-f` | Force re-initialization (overwrites existing configuration) |
| `--verbose` | `-v` | Show detailed initialization steps |

**Examples:**

```bash
# Basic initialization
hepsw init

# Initialize with custom workspace location
hepsw init --workspace /opt/hepsw

# Initialize without fetching the package index
hepsw init --skip-index

# Force re-initialization of existing workspace
hepsw init --force
```

**Output:**
```text
Initializing HepSW workspace at ~/.hepsw/
✓ Created directory structure
✓ Generated configuration file
✓ Cloned package index (250 packages)
✓ Workspace ready
```

---

## Package Discovery

Commands for searching and exploring available packages.

### `hepsw search`

**Purpose:** Search for packages in the local workspace or upstream index by name, description, or metadata.

**Usage:**
```bash
hepsw search [flags] [query]
```

**What it does:**
- Searches both local and remote repositories by default
- Supports wildcards and pattern matching
- Can filter by version constraints, dependencies, and options

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--local` | `-l` | Search only in the local workspace |
| `--remote` | `-r` | Search only in the upstream index |
| `--name <name>` | `-n` | Search specifically by package name |
| `--keyword <keyword>` | `-k` | Search by keyword in description |
| `--version <constraint>` | `-v` | Filter by version constraint (e.g., `">=6.30"`, `">8.3"`) |
| `--depends-on <name>` | `-d` | Show packages that depend on the specified package |
| `--needed-for <name>` | `-a` | Show packages that are required by the specified package |
| `--json` | | Output results in JSON format |
| `--limit <number>` | | Limit the number of results shown |

**Examples:**

```bash
# Search for a package by name (local and remote)
hepsw search root

# Search only in local workspace
hepsw search --local root

# Search with wildcards
hepsw search --remote "*data*"

# Search by name with version constraint
hepsw search --name pythia8 --version ">8.3"

# Find packages that depend on root
hepsw search --depends-on root --version ">=6.30"

# Find dependencies of a package
hepsw search --needed-for geant4

# Get JSON output for scripting
hepsw search root --json
```

**Output Example:**
```text
Found 3 package(s):

- root
  Versions: 6.30.02, 6.28.06, 6.24/06
  Description: An object-oriented framework for large scale data analysis
  
- rootpy
  Version: 0.9.5
  Description: Python bindings for ROOT
  
- root_numpy
  Version: 4.6.2
  Description: NumPy bindings for ROOT
```

---

### `hepsw list`

**Purpose:** List all packages available in the workspace or index.

**Usage:**
```bash
hepsw list [flags]
```

**What it does:**
- Shows all available packages with their versions
- Can filter by installation status or source availability

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--local` | `-l` | List only locally available packages |
| `--remote` | `-r` | List only packages from the index |
| `--installed` | `-i` | Show only installed packages |
| `--fetched` | `-f` | Show only fetched (but not built) packages |
| `--built` | `-b` | Show only built packages |
| `--format <format>` | | Output format: `table`, `list`, `json` (default: `table`) |
| `--sort-by <field>` | | Sort by: `name`, `version`, `date` (default: `name`) |

**Examples:**

```bash
# List all packages
hepsw list

# List only installed packages
hepsw list --installed

# List fetched but not built packages
hepsw list --fetched

# List in JSON format
hepsw list --format json

# List remote packages sorted by name
hepsw list --remote --sort-by name
```

---

### `hepsw info`

**Purpose:** Display detailed information about a specific package.

**Usage:**
```bash
hepsw info <package-name> [flags]
```

**What it does:**
- Shows comprehensive package information including description, versions, dependencies, build options, and source details
- Equivalent to `hepsw whatis` but with more concise output

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--version <version>` | `-v` | Show information for a specific version |
| `--show-deps` | `-d` | Display full dependency tree |
| `--show-options` | `-o` | Show all available build options |
| `--show-recipe` | `-r` | Display recipe steps |
| `--json` | | Output in JSON format |

**Examples:**

```bash
# Get basic package information
hepsw info root

# Get information for a specific version
hepsw info root --version 6.30.02

# Show with dependency tree
hepsw info pythia8 --show-deps

# Show all build options
hepsw info root --show-options
```

**Output Example:**
```text
Package: root
Version: 6.30.02 (latest)
Description: The ROOT data analysis framework

Source:
  Type: git
  URL: https://github.com/root-project/root.git
  Tag: v6-30-02

Dependencies (Build):
  - cmake >=3.16
  - gcc >=9.3
  - python >=3.8 (optional, for: with-python)

Build Options:
  - with-python: Enable Python bindings
  - with-gui: Enable GUI components
  - with-ssl: Enable SSL support

Status: Not installed
```

---

## Source Management

Commands for fetching and managing package source code.

### `hepsw fetch`

**Purpose:** Download package source code and prepare it for building.

**Usage:**
```bash
hepsw fetch <package-name> [package-name...] [flags]
```

**What it does:**
1. Looks up the package in the index repository
2. Downloads/clones the source code to `~/.hepsw/sources/<package-name>/<version>/src`
3. Copies the manifest file
4. Generates a `build.yml` file for the build process

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--version <version>` | `-v` | Fetch a specific version (default: latest) |
| `--path <path>` | `-p` | Use custom path for source (not recommended) |
| `--deps-depth <number>` | `-d` | Also fetch dependencies up to specified depth |
| `--third-party` | `-t` | Import a third-party manifest from local path |
| `--environment <name>` | `-e` | Fetch all packages for a specific environment |
| `--force` | `-f` | Re-fetch even if already present |
| `--verify-checksum` | | Verify source checksums after download |
| `--shallow` | | Use shallow clone for git repositories (faster) |

**Examples:**

```bash
# Fetch a single package (latest version)
hepsw fetch root

# Fetch a specific version
hepsw fetch root --version 6.28.06

# Fetch multiple packages
hepsw fetch pythia8 hepmc3 boost

# Fetch with dependencies (depth 2)
hepsw fetch geant4 --deps-depth 2

# Fetch all packages for an environment
hepsw fetch --environment fccsw

# Import a third-party manifest
hepsw fetch --third-party /path/to/manifest.yaml

# Force re-fetch
hepsw fetch root --force
```

**Output:**
```text
Fetching: root@6.30.02
✓ Found in index
✓ Cloning from https://github.com/root-project/root.git
✓ Checked out tag v6-30-02
✓ Copied manifest
✓ Generated build.yml
Fetched to: ~/.hepsw/sources/root/6.30.02/
```

---

### `hepsw update`

**Purpose:** Update the package index and optionally update fetched sources.

**Usage:**
```bash
hepsw update [flags] [package-name...]
```

**What it does:**
- Updates the local package index from upstream
- Optionally updates fetched source code to newer versions

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--index-only` | `-i` | Only update the package index |
| `--sources` | `-s` | Update fetched sources to latest versions |
| `--check` | `-c` | Check for updates without applying them |
| `--all` | `-a` | Update all fetched packages |

**Examples:**

```bash
# Update package index
hepsw update

# Update index and check for source updates
hepsw update --check

# Update specific package sources
hepsw update root pythia8 --sources

# Update all fetched packages
hepsw update --all --sources
```

---

### `hepsw verify`

**Purpose:** Verify integrity of fetched sources using checksums.

**Usage:**
```bash
hepsw verify <package-name> [flags]
```

**What it does:**
- Checks source code integrity against manifest checksums
- Validates git commit hashes and tags
- Reports any corruption or mismatches

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--version <version>` | `-v` | Verify specific version |
| `--all` | `-a` | Verify all fetched packages |
| `--fix` | `-f` | Attempt to re-fetch corrupted sources |

**Examples:**

```bash
# Verify a single package
hepsw verify root

# Verify all fetched packages
hepsw verify --all

# Verify and fix if corrupted
hepsw verify root --fix
```

---

## Build Operations

Commands for analyzing, building, and testing packages.

### `hepsw whatis`

**Purpose:** Display comprehensive information about a package's structure, dependencies, and build recipe.

**Usage:**
```bash
hepsw whatis <package-name> [flags]
```

**What it does:**
- Shows detailed package overview including all metadata
- Displays dependency tree, build options, and recipe structure
- More verbose than `hepsw info`

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--version <version>` | `-v` | Show information for specific version |
| `--verbose` | | Show detailed recipe steps |
| `--deps-tree` | | Display dependency tree visually |

**Examples:**

```bash
# Get comprehensive package information
hepsw whatis root

# Show with detailed recipe
hepsw whatis root --verbose

# Show with dependency tree
hepsw whatis pythia8 --deps-tree
```

**Output Example:**
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

---

### `hepsw evaluate`

**Purpose:** Analyze a package and check if it's ready to build successfully.

**Usage:**
```bash
hepsw evaluate <package-name> [version] [flags]
```

**What it does:**
- Validates manifest structure
- Checks dependency availability
- Detects version conflicts
- Verifies toolchain compatibility
- Validates build target support
- Estimates disk space requirements
- Validates recipe syntax

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--version <version>` | `-v` | Evaluate specific version |
| `--fix` | `-f` | Automatically fetch missing dependencies |
| `--ignore-warnings` | | Only show critical errors |
| `--deps-tree` | | Show full dependency evaluation tree |

**Exit Codes:**
- `0`: No issues, safe to build
- `1`: Warnings found, may build but could fail
- `2`: Critical errors, build will fail

**Examples:**

```bash
# Evaluate a package
hepsw evaluate pythia8

# Evaluate and auto-fix issues
hepsw evaluate geant4 --fix

# Evaluate specific version
hepsw evaluate root --version 6.28.06

# Show only critical errors
hepsw evaluate pythia8 --ignore-warnings
```

**Output Example:**
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

---

### `hepsw walk`

**Purpose:** Simulate the build process step-by-step without executing commands.

**Usage:**
```bash
hepsw walk <package-name> [version] [flags]
```

**What it does:**
- Shows each recipe step in order
- Displays variable interpolation
- Shows commands that would be executed
- Shows environment variables
- Displays conditional step evaluation
- Useful for debugging recipes and understanding builds

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--version <version>` | `-v` | Walk through specific version |
| `--stage <stage>` | | Show only specific stage: `configure`, `build`, `test`, `install` |
| `--with <option>` | | Enable a build option |
| `--without <option>` | | Disable a build option |
| `--show-env` | | Display all environment variables at each step |
| `--export-script <path>` | | Export build process as bash script |

**Examples:**

```bash
# Walk through build process
hepsw walk root

# Show only configure stage
hepsw walk root --stage configure

# Walk with specific options
hepsw walk root --with python --with gui

# Show environment variables
hepsw walk pythia8 --show-env

# Export as script for manual debugging
hepsw walk root --export-script /tmp/build-root.sh
```

**Output Example:**
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

════════════════════════════════════════════════════════════
Summary
════════════════════════════════════════════════════════════
Total steps: 6
Estimated time: 15-30 minutes
```

---

### `hepsw build`

**Purpose:** Build a package from source.

**Usage:**
```bash
hepsw build <package-name> [version] [flags]
```

**What it does:**
1. Sets up the build environment
2. Executes configure steps
3. Compiles the source code
4. Runs tests (if enabled)
5. Installs to the install directory
6. Updates workspace state in `hepsw.yaml`
7. Generates environment setup scripts

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--version <version>` | `-v` | Build specific version |
| `--latest` | | Build the latest version |
| `--with <option>` | | Enable a build option (can use multiple times) |
| `--without <option>` | | Disable a build option |
| `--list-options` | | List available build options and exit |
| `--jobs <number>` | `-j` | Number of parallel jobs (default: 4) |
| `--build-type <type>` | | Build type: `Debug`, `Release`, `RelWithDebInfo`, `MinSizeRel` |
| `--prefix <path>` | | Custom installation prefix |
| `--verbose` | | Show detailed build output |
| `--quiet` | `-q` | Suppress non-error output |
| `--clean` | | Clean build directory before building |
| `--rebuild` | | Rebuild from scratch (clean + build) |
| `--with-deps` | | Also build dependencies |
| `--skip-tests` | | Skip the test stage |
| `--skip-install` | | Skip the install stage (build only) |
| `--continue-on-test-fail` | | Continue even if tests fail |
| `--timeout <seconds>` | | Build timeout in seconds (default: 3600) |
| `--keep-build-dir` | | Don't clean build directory after install |
| `--clean-after-install` | | Remove build directory after successful install |
| `--log-file <path>` | | Write build log to custom file |
| `--env <VAR=value>` | | Set environment variable for build |

**Examples:**

```bash
# Basic build
hepsw build root

# Build specific version
hepsw build root --version 6.28.06

# Build with options
hepsw build root --with python --with gui --without ssl

# Build with more parallel jobs
hepsw build pythia8 --jobs 16

# Debug build
hepsw build geant4 --build-type Debug

# Build with dependencies
hepsw build myanalysis --with-deps

# Rebuild from scratch
hepsw build root --rebuild --verbose

# Build and skip tests
hepsw build pythia8 --skip-tests --jobs 8

# Build with custom environment
hepsw build root --env CC=clang --env CXX=clang++

# Clean build with verbose output
hepsw build root --clean --verbose --jobs 12
```

**Output:**
```text
Building: root@6.30.02
✓ Manifest loaded
✓ Dependencies satisfied
✓ Build environment ready

[1/4] Configure
  Running CMake configuration...
  ✓ Configuration complete (45s)

[2/4] Build
  Compiling with 8 parallel jobs...
  ████████████████████████████████████████ 100%
  ✓ Build complete (18m 32s)

[3/4] Test
  Running tests...
  ✓ All tests passed (2m 15s)

[4/4] Install
  Installing to ~/.hepsw/install/root/6.30.02
  ✓ Installation complete (1m 05s)

✓ Build successful
Total time: 22m 37s
Log: ~/.hepsw/logs/root-6.30.02-build.log
Environment: ~/.hepsw/env/root-6.30.02.sh
```

---

### `hepsw test`

**Purpose:** Run tests for an already-built package.

**Usage:**
```bash
hepsw test <package-name> [flags]
```

**What it does:**
- Runs the test suite defined in the package manifest
- Useful for re-testing after manual changes

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--version <version>` | `-v` | Test specific version |
| `--verbose` | | Show detailed test output |
| `--parallel <number>` | `-j` | Number of parallel test jobs |
| `--filter <pattern>` | | Run only tests matching pattern |

**Examples:**

```bash
# Run all tests
hepsw test root

# Run tests with verbose output
hepsw test root --verbose

# Run tests in parallel
hepsw test pythia8 --parallel 8

# Run specific tests
hepsw test root --filter "*vector*"
```

---

### `hepsw install`

**Purpose:** Install an already-built package without rebuilding.

**Usage:**
```bash
hepsw install <package-name> [flags]
```

**What it does:**
- Runs the install step from the build directory
- Useful after building with `--skip-install`

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--version <version>` | `-v` | Install specific version |
| `--force` | `-f` | Force reinstall even if already installed |

**Examples:**

```bash
# Install a built package
hepsw install root

# Force reinstall
hepsw install root --force
```

---

## Environment Management

Commands for managing runtime environments and environment variables.

### `hepsw env`

**Purpose:** Manage and activate package environments.

**Usage:**
```bash
hepsw env [subcommand] [flags]
```

**Subcommands:**

#### `hepsw env create`

Create a new environment with a set of packages.

```bash
hepsw env create <env-name> [package...] [flags]
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--from-file <path>` | `-f` | Create from environment definition file |
| `--description <text>` | `-d` | Environment description |

**Examples:**

```bash
# Create environment with packages
hepsw env create myenv root pythia8 geant4

# Create from file
hepsw env create analysis --from-file env.yaml
```

#### `hepsw env activate`

Activate an environment (generates shell commands to source).

```bash
hepsw env activate <env-name> [flags]
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--shell <shell>` | `-s` | Target shell: `bash`, `zsh`, `fish` (auto-detected) |

**Examples:**

```bash
# Activate environment (bash/zsh)
eval $(hepsw env activate myenv)

# Activate for specific shell
eval $(hepsw env activate myenv --shell bash)
```

#### `hepsw env list`

List all available environments.

```bash
hepsw env list [flags]
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--verbose` | `-v` | Show packages in each environment |

#### `hepsw env show`

Show details of an environment.

```bash
hepsw env show <env-name> [flags]
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--packages` | `-p` | List packages in environment |
| `--variables` | `-v` | Show environment variables |

#### `hepsw env delete`

Delete an environment.

```bash
hepsw env delete <env-name> [flags]
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--force` | `-f` | Delete without confirmation |

#### `hepsw env export`

Export environment definition to a file.

```bash
hepsw env export <env-name> [flags]
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--output <path>` | `-o` | Output file path |
| `--format <format>` | | Format: `yaml`, `json` (default: `yaml`) |

**Examples:**

```bash
# Export environment
hepsw env export myenv --output myenv.yaml

# Export as JSON
hepsw env export myenv --format json --output myenv.json
```

---

### `hepsw shell`

**Purpose:** Start a new shell with an environment activated.

**Usage:**
```bash
hepsw shell <env-name> [flags]
```

**What it does:**
- Launches a new shell with the specified environment loaded
- Shows environment name in prompt

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--shell <shell>` | `-s` | Shell to use: `bash`, `zsh`, `fish` |
| `--command <cmd>` | `-c` | Run command in environment and exit |

**Examples:**

```bash
# Start interactive shell
hepsw shell myenv

# Run command in environment
hepsw shell myenv --command "root -b -q script.C"
```

---

## Configuration Management

Commands for managing HepSW configuration and preferences.

### `hepsw config`

**Purpose:** View and modify HepSW configuration.

**Usage:**
```bash
hepsw config [subcommand] [flags]
```

**Subcommands:**

#### `hepsw config show`

Display current configuration.

```bash
hepsw config show [flags]
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--all` | `-a` | Show all settings including defaults |
| `--json` | | Output as JSON |

#### `hepsw config get`

Get a specific configuration value.

```bash
hepsw config get <key>
```

**Examples:**

```bash
# Get parallel build setting
hepsw config get userConfig.parallelBuilds

# Get workspace path
hepsw config get workspace
```

#### `hepsw config set`

Set a configuration value.

```bash
hepsw config set <key> <value>
```

**Examples:**

```bash
# Set default parallel builds
hepsw config set userConfig.parallelBuilds 8

# Set verbosity level
hepsw config set userConfig.verbosity debug

# Set custom workspace
hepsw config set workspace /opt/hepsw
```

#### `hepsw config reset`

Reset configuration to defaults.

```bash
hepsw config reset [flags]
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--key <key>` | `-k` | Reset specific key only |
| `--force` | `-f` | Reset without confirmation |

**Examples:**

```bash
# Reset all configuration
hepsw config reset

# Reset specific setting
hepsw config reset --key userConfig.parallelBuilds
```

#### `hepsw config validate`

Validate configuration file.

```bash
hepsw config validate [flags]
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--fix` | `-f` | Attempt to fix issues automatically |

---

## Maintenance & Utilities

Commands for maintaining and cleaning up the workspace.

### `hepsw clean`

**Purpose:** Clean up workspace by removing build artifacts, old versions, or sources.

**Usage:**
```bash
hepsw clean [flags] [package...]
```

**What it does:**
- Removes specified build artifacts, sources, or installations
- Helps manage disk space

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--builds` | `-b` | Clean build directories |
| `--sources` | `-s` | Remove fetched sources |
| `--installs` | `-i` | Remove installed packages |
| `--logs` | `-l` | Remove log files |
| `--all` | `-a` | Clean everything (use with caution) |
| `--old-versions` | | Keep only latest version of each package |
| `--dry-run` | `-n` | Show what would be removed without removing |
| `--force` | `-f` | Skip confirmation prompts |

**Examples:**

```bash
# Clean build directory for a package
hepsw clean --builds root

# Clean old versions (keep latest)
hepsw clean --old-versions

# Clean everything for a package
hepsw clean --all root

# Dry run to see what would be removed
hepsw clean --builds --dry-run

# Clean all build directories
hepsw clean --builds --all

# Clean logs older than 30 days
hepsw clean --logs
```

---

### `hepsw status`

**Purpose:** Show the status of packages and workspace.

**Usage:**
```bash
hepsw status [package...] [flags]
```

**What it does:**
- Shows which packages are fetched, built, and installed
- Displays workspace disk usage
- Shows recent activity

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--all` | `-a` | Show status of all packages |
| `--verbose` | `-v` | Show detailed status information |
| `--format <format>` | | Output format: `table`, `list`, `json` |

**Examples:**

```bash
# Show status of specific packages
hepsw status root pythia8

# Show status of all packages
hepsw status --all

# Verbose status
hepsw status root --verbose

# JSON output
hepsw status --all --format json
```

**Output Example:**
```text
Workspace: ~/.hepsw/
Disk Usage: 4.2 GB (builds: 2.1 GB, sources: 1.8 GB, installs: 0.3 GB)

Package Status:
┌──────────┬─────────┬──────────┬────────┬────────────┐
│ Package  │ Version │ Fetched  │ Built  │ Installed  │
├──────────┼─────────┼──────────┼────────┼────────────┤
│ root     │ 6.30.02 │ ✓        │ ✓      │ ✓          │
│ pythia8  │ 8.310   │ ✓        │ ✓      │ ✗          │
│ geant4   │ 11.1.0  │ ✓        │ ✗      │ ✗          │
└──────────┴─────────┴──────────┴────────┴────────────┘
```

---

### `hepsw doctor`

**Purpose:** Diagnose issues with the workspace and system configuration.

**Usage:**
```bash
hepsw doctor [flags]
```

**What it does:**
- Checks workspace integrity
- Verifies required tools are installed
- Checks for common configuration issues
- Validates package state consistency
- Suggests fixes for detected problems

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--fix` | `-f` | Automatically fix detected issues |
| `--verbose` | `-v` | Show detailed diagnostic information |

**Examples:**

```bash
# Run diagnostics
hepsw doctor

# Run diagnostics and auto-fix
hepsw doctor --fix

# Verbose diagnostics
hepsw doctor --verbose
```

**Output Example:**
```text
Running HepSW diagnostics...

✓ Workspace structure valid
✓ Configuration file valid
✓ Package index up to date
✓ Required tools available:
  - cmake 3.22.1
  - gcc 11.2.0
  - git 2.34.1
⚠ Build directory for root@6.28.06 is incomplete
  Suggestion: Run 'hepsw clean --builds root' and rebuild

⚠ Log directory is large (2.3 GB)
  Suggestion: Run 'hepsw clean --logs'

✓ No critical issues found
2 warnings
```

---

### `hepsw export`

**Purpose:** Export workspace state or package list for reproducibility.

**Usage:**
```bash
hepsw export [flags]
```

**What it does:**
- Exports installed packages and versions
- Generates reproducible environment specifications
- Useful for CI/CD and collaboration

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--output <path>` | `-o` | Output file path (default: stdout) |
| `--format <format>` | `-f` | Format: `yaml`, `json`, `requirements` |
| `--with-versions` | | Include exact versions |
| `--installed-only` | | Export only installed packages |

**Examples:**

```bash
# Export to YAML
hepsw export --output workspace.yaml

# Export with exact versions
hepsw export --with-versions --output requirements.yaml

# Export as JSON
hepsw export --format json --output workspace.json

# Export only installed packages
hepsw export --installed-only
```

---

### `hepsw import`

**Purpose:** Import workspace state from an export file.

**Usage:**
```bash
hepsw import <file> [flags]
```

**What it does:**
- Reads exported workspace specification
- Fetches and builds packages as specified
- Recreates environments

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--fetch-only` | | Only fetch sources, don't build |
| `--build` | `-b` | Fetch and build all packages |
| `--jobs <number>` | `-j` | Parallel jobs for builds |
| `--continue-on-error` | | Continue even if some packages fail |

**Examples:**

```bash
# Import and fetch packages
hepsw import workspace.yaml

# Import and build all packages
hepsw import workspace.yaml --build --jobs 8

# Import but continue on errors
hepsw import workspace.yaml --build --continue-on-error
```

---

### `hepsw version`

**Purpose:** Display HepSW version information.

**Usage:**
```bash
hepsw version [flags]
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--short` | `-s` | Show only version number |
| `--json` | | Output as JSON |

**Examples:**

```bash
# Show full version info
hepsw version

# Show version number only
hepsw version --short
```

**Output:**
```text
HepSW version 1.0.0
Build: 2026-02-14
Commit: a1b2c3d
Go: go1.21.5
Platform: linux/amd64
```

---

## Global Flags

These flags can be used with any HepSW command.

| Flag | Short | Description |
|------|-------|-------------|
| `--help` | `-h` | Show help for the command |
| `--verbose` | `-v` | Enable verbose output |
| `--debug` | | Enable debug mode with detailed logging |
| `--quiet` | `-q` | Suppress non-error output |
| `--workspace <path>` | `-w` | Use alternate workspace directory |
| `--config <path>` | `-c` | Use alternate configuration file |
| `--color <when>` | | Colorize output: `always`, `auto`, `never` |
| `--yes` | `-y` | Automatically answer yes to prompts |
| `--no-color` | | Disable colored output |

**Examples:**

```bash
# Run command with verbose output
hepsw build root --verbose

# Use alternate workspace
hepsw --workspace /tmp/hepsw list

# Debug mode
hepsw --debug build pythia8

# Quiet mode (only errors)
hepsw --quiet fetch root

# Auto-accept all prompts
hepsw --yes clean --all
```

---

## Command Chaining and Workflows

### Common Workflows

**Initial Setup:**
```bash
# Initialize workspace
hepsw init

# Search for packages
hepsw search root

# Get package info
hepsw info root
```

**Fetch and Build:**
```bash
# Fetch package
hepsw fetch root

# Evaluate before building
hepsw evaluate root

# Simulate build
hepsw walk root

# Build package
hepsw build root --with python --jobs 8
```

**Environment Setup:**
```bash
# Create environment
hepsw env create analysis root pythia8 geant4

# Activate environment
eval $(hepsw env activate analysis)

# Or start a shell
hepsw shell analysis
```

**Maintenance:**
```bash
# Check workspace status
hepsw status --all

# Run diagnostics
hepsw doctor

# Clean old builds
hepsw clean --builds --old-versions

# Update package index
hepsw update
```

**Reproducibility:**
```bash
# Export workspace
hepsw export --with-versions --output workspace.yaml

# On another machine
hepsw init
hepsw import workspace.yaml --build
```

---

## Exit Codes

HepSW uses standard exit codes:

| Code | Meaning |
|------|---------|
| `0` | Success |
| `1` | General error or warning |
| `2` | Command-line usage error |
| `3` | Configuration error |
| `4` | Dependency error |
| `5` | Build error |
| `6` | Network error |
| `7` | File system error |

---

## Environment Variables

HepSW respects the following environment variables:

| Variable | Description |
|----------|-------------|
| `HEPSW_WORKSPACE` | Override default workspace location |
| `HEPSW_CONFIG` | Override configuration file location |
| `HEPSW_JOBS` | Default number of parallel jobs |
| `HEPSW_VERBOSITY` | Default verbosity level: `quiet`, `info`, `debug` |
| `HEPSW_NO_COLOR` | Disable colored output (any value) |
| `HEPSW_INDEX_URL` | Custom package index repository URL |

**Example:**
```bash
export HEPSW_JOBS=16
export HEPSW_VERBOSITY=debug
hepsw build root
```

---

## Notes and Best Practices

### Version Specifications

HepSW supports semantic version constraints:

- `1.2.3` - Exact version
- `>=1.2.3` - Greater than or equal to
- `>1.2.3` - Greater than
- `<2.0.0` - Less than
- `<=2.0.0` - Less than or equal to
- `>=1.2.3,<2.0.0` - Range (comma-separated)
- `~1.2.3` - Compatible with 1.2.x
- `^1.2.3` - Compatible with 1.x.x
- `latest` - Latest available version

### Build Options

Build options follow the naming convention:
- `with-<feature>` - Enable optional feature
- `without-<feature>` - Explicitly disable feature

Examples: `with-python`, `with-gui`, `with-ssl`, `without-tests`

### Workspace Management

- The workspace state is authoritative; manual changes may be reverted
- Always use HepSW commands to modify the workspace
- Regular backups of `hepsw.yaml` are recommended
- Use `hepsw export` for reproducible setups

### Performance Tips

- Use `--jobs` flag to match your CPU core count
- Use `--deps-depth` with `fetch` to prepare dependencies
- Use `--skip-tests` during development iterations
- Use `--clean-after-install` to save disk space
- Use `--shallow` cloning for faster git fetches

### Debugging Builds

1. Use `hepsw evaluate` to check readiness
2. Use `hepsw walk` to simulate the build
3. Use `--verbose` flag for detailed output
4. Check logs in `~/.hepsw/logs/`
5. Use `hepsw doctor` to diagnose issues
6. Export build script with `walk --export-script` for manual debugging

---

## Additional Resources

- **Package Index Repository:** https://github.com/thisismeamir/hepsw-package-index/
- **Issue Tracker:** Report bugs and request features
- **Documentation:** Full user guide and tutorials
- **Community:** Discussion forums and support channels

---

**End of HepSW CLI API Reference**