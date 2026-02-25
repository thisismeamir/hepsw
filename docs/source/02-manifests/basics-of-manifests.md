# Basics of Manifests

HepSW manifests are a way to define and manage software.
They provide a structured format to specify information about software packages, their versions, dependencies, and configurations.
Manifests are written in YAML format, making them easy to read and write.

It is important to understand the structure and components of a HepSW manifest to effectively analyze and utilize them.

## Quick Start: Simple Manifest Example

For most packages, you only need the essential fields. Here's a minimal manifest:

```yaml
name: example-lib
version: 1.2.3
description: A simple example library

source:
  type: tarball
  url: https://example.com/example-lib-1.2.3.tar.gz

recipe:
  configure:
    - name: Configure
      command: ./configure --prefix=${INSTALL_PREFIX}
  build:
    - name: Build
      command: make -j${NCORES}
  install:
    - name: Install
      command: make install
```

This minimal manifest is sufficient for simple packages with standard build procedures and no complex dependencies.

## Structure of a Manifest

### Header

Every manifest starts with basic information about the package, including its name, version, a brief description, and more.

```yaml
name: package-name
version: x.y.z
description: A brief description of the package.
source:
  type: tarball
  url: https://example.com/source.tar.gz
  checksum: sha256:abcdef1234567890
```

#### Source Types

The source section describes where to obtain the source code, including the type, URL, and optional checksum for verification.
There are several source types supported:

- **tarball**: A compressed archive file (e.g., .tar.gz, .zip).
- **git**: A Git repository URL (can include branch or tag).
- **svn**: A Subversion repository URL.
- **local**: A local file path (not recommended for distribution).

**Note on Checksums**: Checksums are optional but highly recommended for security and reproducibility. When packages are included in the HepSW Package Index Repository (HPIR), checksums are verified against the repository's records. For local or third-party manifests, checksums ensure the downloaded source matches expectations. If omitted, HepSW will fetch the source without verification.

#### Metadata (Optional)

Apart from the core fields, additional metadata can be included, such as authors, license, homepage, and documentation links:

```yaml
metadata:
  authors:
    - Name One <email@example.com>
    - Name Two <email2@example.com>
  homepage: https://example.com
  license: MIT
  documentation: https://docs.example.com
```

#### Example: ROOT Project Header

```yaml
name: root
version: 6.30.02
description: The ROOT data analysis framework.
source:
  type: git
  url: https://github.com/root-project/root.git
  tag: v6-30-02
  checksum: sha256:abcdef1234567890
metadata:
  authors:
    - René Brun 
    - Fons Rademakers
  homepage: https://root.cern
  license: LGPL-3.0
  documentation: https://root.cern/doc
```

### Specifications

The specifications section defines the build and runtime requirements for the package.
For simple packages, this section can be minimal or omitted entirely. For complex packages with dependencies and build options, it becomes essential.

#### Minimal Specifications

```yaml
specifications:
  build:
    toolchain:
      - cmake >=3.15
      - gcc >=9.0
```

#### Full Specifications

```yaml
specifications:
  build:
    toolchain: # What tools and versions are needed to build the package
      - cmake >=3.15 
      - gcc >=9.0
    targets: # Supported build targets
      - linux-x86_64
      - linux-aarch64
      - darwin-x86_64
    extensions: # Optional build features (e.g., CMake flags, configure options)
      - with-ssl
      - with-gui
      - with-python
    dependencies: # Build-time dependencies
      - name: openssl
        version: ">=1.1.1,<4.0.0"
        isOptional: false
        forOptions:
          - with-ssl
      - name: qt5
        version: ">=5.12"
        isOptional: true
        forOptions:
          - with-gui
        withOptions:
          - with-ssl
    variables:
      parallelism: 8 # Number of parallel jobs during build
      installPrefix: /opt/hep
      
  runtime:
    dependencies: # Runtime dependencies
      - name: openssl
        version: ">=1.1.1"
        forOptions:
          - with-ssl
      - name: qt5
        version: ">=5.12"
        forOptions:
          - with-gui
          
  environment:
    build:
      variables: # Environment variables needed during build
        CXX_FLAGS: "-O3 -march=native"
        BUILD_TYPE: Release
    runtime:
      variables: # Environment variables needed at runtime
        LD_LIBRARY_PATH: "${INSTALL_PREFIX}/lib:${LD_LIBRARY_PATH}"
    self:
      variables: # Environment variables this package provides to other packages
        EXAMPLE_ROOT: "${INSTALL_PREFIX}"
        EXAMPLE_INCLUDE: "${INSTALL_PREFIX}/include"
        EXAMPLE_LIB: "${INSTALL_PREFIX}/lib"
```

#### Specification Notes

1. **Conditional Dependencies**: The `forOptions` and `withOptions` fields specify conditional dependencies based on build options. A dependency with `forOptions: [with-ssl]` is only required when building with the `with-ssl` extension enabled.

2. **Environment Variables**: The `environment` section has three subsections:
    - `build`: Variables needed during the build process
    - `runtime`: Variables needed when using the package
    - `self`: Variables this package exports for dependent packages

3. **Toolchain**: Specifies required build tools and their version constraints.

4. **Targets**: Lists supported build architectures and platforms.

5. **Extensions**: Optional features that can be enabled/disabled during build.

6. **Version Constraints**: Dependencies can specify version ranges using operators like `>=`, `<`, `==`, etc.

7. **Optional Dependencies**: The `isOptional: true` field indicates a dependency that can be omitted.

### Recipe

The recipe section outlines the steps required to build and install the package.
This is similar to GitHub Actions or CI/CD pipelines where you define a series of steps to achieve a goal.

#### Simple Recipe

```yaml
recipe:
  configure:
    - name: Configure
      command: ./configure --prefix=${INSTALL_PREFIX}
  build:
    - name: Build
      command: make -j${NCORES}
  install:
    - name: Install
      command: make install
```

#### Recipe Step Types

HepSW supports multiple step types to handle different build scenarios:

##### 1. Command Step

Execute a shell command directly:

```yaml
recipe:
  build:
    - name: Compile sources
      type: command
      command: gcc -o myapp main.c -lm
```

##### 2. Script Step

Run a script file from the package or a custom location:

```yaml
recipe:
  configure:
    - name: Run configuration script
      type: script
      script: scripts/custom-configure.sh
      args:
        - --enable-optimizations
        - --prefix=${INSTALL_PREFIX}
```

##### 3. CMake Step

Run CMake with specified arguments (simplified syntax for CMake projects):

```yaml
recipe:
  configure:
    - name: Configure with CMake
      type: cmake
      args:
        - -DCMAKE_BUILD_TYPE=Release
        - -DCMAKE_INSTALL_PREFIX=${INSTALL_PREFIX}
        - -DENABLE_SSL=${OPTIONS_WITH_SSL}
  build:
    - name: Build with CMake
      type: cmake
      target: build
      parallel: ${NCORES}
```

##### 4. Set Step

Define variables to be used in subsequent steps:

```yaml
recipe:
  configure:
    - name: Define build variables
      type: set
      variables:
        BUILD_DIR: ${SOURCE_DIR}/build
        INSTALL_DIR: /opt/myapp
        COMPILER_FLAGS: "-O3 -march=native"
    - name: Configure
      command: ./configure --prefix=${INSTALL_DIR} CXXFLAGS="${COMPILER_FLAGS}"
```

##### 5. Conditional Step

Execute steps based on conditions:

```yaml
recipe:
  build:
    - name: Build with SSL support
      type: conditional
      condition: ${OPTIONS_WITH_SSL}
      steps:
        - name: Configure SSL
          command: ./configure --with-ssl
        - name: Build
          command: make ssl-modules
    - name: Build without SSL
      type: conditional
      condition: "!${OPTIONS_WITH_SSL}"
      steps:
        - name: Configure without SSL
          command: ./configure --without-ssl
```

##### 6. Parallel Step

Execute multiple steps in parallel:

```yaml
recipe:
  build:
    - name: Build modules in parallel
      type: parallel
      steps:
        - name: Build core
          command: make -C core
        - name: Build utils
          command: make -C utils
        - name: Build plugins
          command: make -C plugins
```

##### 7. Environment Step

Set environment variables for subsequent steps:

```yaml
recipe:
  configure:
    - name: Set build environment
      type: environment
      variables:
        CC: gcc-11
        CXX: g++-11
        MAKEFLAGS: "-j${NCORES}"
    - name: Configure
      command: ./configure --prefix=${INSTALL_PREFIX}
```

##### 8. AskInput Step

Prompt the user for input during the build process (useful for interactive configurations):

```yaml
recipe:
  configure:
    - name: Get installation path
      type: askInput
      prompt: "Enter custom installation path (or press Enter for default)"
      variable: CUSTOM_INSTALL_PATH
      default: ${INSTALL_PREFIX}
    - name: Configure with custom path
      command: ./configure --prefix=${CUSTOM_INSTALL_PATH}
```

##### 9. Custom Step

Define a custom build step with user-defined behavior (requires custom handlers):

```yaml
recipe:
  build:
    - name: Custom build step
      type: custom
      handler: my_custom_builder
      params:
        optimization_level: 3
        target_arch: x86_64
```

#### Variable Interpolation

Variables can be referenced using the `${VARIABLE_NAME}` syntax. HepSW provides several built-in variables and allows you to access manifest fields directly.

**Built-in Variables:**
- `${INSTALL_PREFIX}`: Installation directory
- `${SOURCE_DIR}`: Source code directory
- `${BUILD_DIR}`: Build directory
- `${NCORES}`: Number of CPU cores for parallel builds
- `${OPTIONS_WITH_*}`: Boolean flags for enabled extensions (e.g., `${OPTIONS_WITH_SSL}`)

**Accessing Manifest Fields:**
You can access any field from the manifest using the `${manifest.path.to.field}` syntax:

```yaml
specifications:
  build:
    variables:
      parallelism: 8
      compiler: gcc

recipe:
  build:
    - name: Build
      command: make -j${manifest.specifications.build.variables.parallelism}
    - name: Report compiler
      command: echo "Using ${manifest.specifications.build.variables.compiler}"
```

**Creating Aliases with Set:**
For frequently used or complex manifest paths, use the `set` step to create shorter aliases:

```yaml
recipe:
  configure:
    - name: Create variable aliases
      type: set
      variables:
        JOBS: ${manifest.specifications.build.variables.parallelism}
        COMPILER: ${manifest.specifications.build.variables.compiler}
        INSTALL_DIR: ${manifest.specifications.build.variables.install_prefix}
    - name: Configure
      command: ./configure --prefix=${INSTALL_DIR}
    - name: Build
      command: ${COMPILER} -j${JOBS} all
```

This approach keeps your recipe steps clean and readable while still leveraging the full manifest structure.

#### Full Recipe Example

Here's a complete example combining multiple step types:

```yaml
recipe:
  configure:
    - name: Set up build environment
      type: environment
      variables:
        CC: gcc
        CXX: g++
        
    - name: Create build directory
      type: set
      variables:
        BUILD_PATH: ${SOURCE_DIR}/build
        NJOBS: ${manifest.specifications.build.variables.parallelism}
        
    - name: Check for SSL support
      type: conditional
      condition: ${OPTIONS_WITH_SSL}
      steps:
        - name: Configure with SSL
          type: cmake
          args:
            - -DCMAKE_BUILD_TYPE=Release
            - -DCMAKE_INSTALL_PREFIX=${INSTALL_PREFIX}
            - -DENABLE_SSL=ON
            
  build:
    - name: Build project
      type: cmake
      target: build
      parallel: ${NJOBS}
      
    - name: Build documentation
      type: parallel
      steps:
        - name: Build HTML docs
          command: make docs-html
        - name: Build PDF docs
          command: make docs-pdf
          
  test:
    - name: Run tests
      command: ctest --output-on-failure
      
  install:
    - name: Install binaries
      type: cmake
      target: install
```

For further details about the recipe steps and their options, please refer to the [Manifests Revisited](../03-advanced/manifests-revisited) section.

## Complete Manifest Examples

### Simple Package: A Basic Library

```yaml
name: simple-math
version: 2.1.0
description: A simple mathematics library

source:
  type: tarball
  url: https://github.com/example/simple-math/releases/download/v2.1.0/simple-math-2.1.0.tar.gz
  checksum: sha256:1234567890abcdef

specifications:
  build:
    toolchain:
      - cmake >=3.10
      - gcc >=7.0

recipe:
  configure:
    - name: Configure
      type: cmake
      args:
        - -DCMAKE_BUILD_TYPE=Release
        - -DCMAKE_INSTALL_PREFIX=${INSTALL_PREFIX}
  build:
    - name: Build
      type: cmake
      target: build
      parallel: ${NCORES}
  test:
    - name: Test
      command: ctest
  install:
    - name: Install
      type: cmake
      target: install
```

### Complex Package: ROOT with Dependencies

```yaml
name: root
version: 6.30.02
description: The ROOT data analysis framework

source:
  type: git
  url: https://github.com/root-project/root.git
  tag: v6-30-02
  checksum: sha256:abcdef1234567890

metadata:
  authors:
    - René Brun
    - Fons Rademakers
  homepage: https://root.cern
  license: LGPL-3.0
  documentation: https://root.cern/doc

specifications:
  build:
    toolchain:
      - cmake >=3.16
      - gcc >=9.3
    targets:
      - linux-x86_64
      - darwin-x86_64
    extensions:
      - with-python
      - with-gui
      - with-ssl
    dependencies:
      - name: python
        version: ">=3.8"
        forOptions:
          - with-python
      - name: qt5
        version: ">=5.12"
        forOptions:
          - with-gui
      - name: openssl
        version: ">=1.1.1"
        forOptions:
          - with-ssl
    variables:
      parallelism: 8
      build_type: Release

  runtime:
    dependencies:
      - name: python
        version: ">=3.8"
        forOptions:
          - with-python
      - name: qt5
        version: ">=5.12"
        forOptions:
          - with-gui

  environment:
    build:
      variables:
        ROOTSYS: ${SOURCE_DIR}
    runtime:
      variables:
        ROOTSYS: ${INSTALL_PREFIX}
        PATH: "${INSTALL_PREFIX}/bin:${PATH}"
        LD_LIBRARY_PATH: "${INSTALL_PREFIX}/lib:${LD_LIBRARY_PATH}"
        PYTHONPATH: "${INSTALL_PREFIX}/lib:${PYTHONPATH}"
    self:
      variables:
        ROOTSYS: ${INSTALL_PREFIX}
        ROOT_INCLUDE_PATH: ${INSTALL_PREFIX}/include

recipe:
  configure:
    - name: Set build variables
      type: set
      variables:
        BUILD_DIR: ${SOURCE_DIR}/build
        JOBS: ${manifest.specifications.build.variables.parallelism}
        TYPE: ${manifest.specifications.build.variables.build_type}

    - name: Create build directory
      command: mkdir -p ${BUILD_DIR} && cd ${BUILD_DIR}

    - name: Configure ROOT
      type: cmake
      args:
        - -DCMAKE_BUILD_TYPE=${TYPE}
        - -DCMAKE_INSTALL_PREFIX=${INSTALL_PREFIX}
        - -Dpython3=${OPTIONS_WITH_PYTHON}
        - -Dssl=${OPTIONS_WITH_SSL}
        - -Dqt5web=${OPTIONS_WITH_GUI}

  build:
    - name: Build ROOT
      type: cmake
      target: build
      parallel: ${JOBS}

  test:
    - name: Run ROOT tests
      command: cd ${BUILD_DIR} && ctest --output-on-failure -j${JOBS}

  install:
    - name: Install ROOT
      type: cmake
      target: install
```

## Validation and Checks

The user can prompt HepSW to validate a manifest file for correctness and completeness.
This analysis would reveal any missing fields, incorrect formats, dependency issues and other potential problems.

```bash
hepsw evaluate path/to/manifest.yaml
```

This command will analyze the specified manifest file and report any issues found.
The validation process checks for:
- Required fields presence
- Correct data types
- Dependency resolution
- Version constraints
- Build and runtime specifications consistency
- Recipe steps validity
- Environment variable definitions
- Overall structure adherence to the HepSW manifest schema
- Source accessibility (if checksums are provided, they are verified)
- And more

To ensure the manifest is valid and ready for use, it is recommended to run this evaluation before proceeding with building or deploying the package.

### Additional Validation Commands

The `hepsw whatis` command generates a visual representation of the manifest's structure and dependencies:

```bash
hepsw whatis path/to/manifest.yaml
```

The `hepsw walk` command simulates the build process without actually executing the steps:

```bash
hepsw walk path/to/manifest.yaml
```

This command will go through the recipe steps and display what would be executed, allowing you to verify the build process and identify any potential issues before actual execution.

## Best Practices

1. **Start Simple**: Begin with a minimal manifest and add complexity only as needed.

2. **Use Checksums**: Always include checksums for sources when distributing manifests to ensure integrity.

3. **Leverage Variables**: Use the `set` step to create readable aliases for complex manifest paths.

4. **Test Thoroughly**: Use `hepsw walk` to verify your recipe before attempting a real build.

5. **Document Extensions**: Clearly document what each extension/option does in the description or comments.

6. **Version Constraints**: Be specific with version constraints to avoid compatibility issues.

7. **Environment Variables**: Use the `self` section to properly export variables for dependent packages.