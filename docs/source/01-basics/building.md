# Building

Build process is the core functionality of HepSW. It provides analysis, debugging, building, and installation process to ensure consistency, reproducibility 
and more.

## Building in HepSW

To build a software that you've fetched using `hepsw fetch` command, you must start by evaluating the build process. This is possible by either of the following:

```bash
hepsw whatis <package-name>
```
to read about what the software does, and how is it going to be built.

```bash
hepsw evaluate <package-name> <version>
```
for hepsw to evaluate and find potential bottlenecks (missing dependencies, and process evaluation), and:

```bash
hepsw walk <package-name> <version>
```

to simulate the build process.

After that if no bottleneck or fragmentation is found, it is safe to build the software using:

```bash
hepsw build <package-name> <version>
```

which then results in a folder in `~/.hepsw/builds/<package-name>/<version>` with the build process.


## Record Building (for Developers)

## Building Software Guideline

Building software from source is a fundamental skill in scientific computing and HEP software development. This section provides general guidelines applicable to most software projects.

### Prerequisites

Before building any software, ensure you have:

- **Compiler toolchain**: A C/C++ compiler (GCC, Clang, or others) appropriate for your system
- **Build tools**: Common tools like `make`, `cmake`, `autotools`, or project-specific build systems
- **Dependencies**: All required libraries and their development headers
- **System resources**: Adequate disk space and memory for compilation

### General Build Process

Most software follows a similar build pattern:

1. **Obtain the source code**: Download or clone the software repository
2. **Configure**: Set up the build environment and specify installation paths and options
3. **Compile**: Transform source code into executable binaries
4. **Test** (optional but recommended): Verify the build works correctly
5. **Install**: Place binaries and libraries in their final locations

### Common Build Systems

**CMake-based projects**:
```bash
mkdir build && cd build
cmake .. -DCMAKE_INSTALL_PREFIX=/path/to/install
make -j$(nproc)
make install
```

**Autotools-based projects**:
```bash
./configure --prefix=/path/to/install
make -j$(nproc)
make install
```

**Python projects**:
```bash
pip install .
# or
python setup.py install
```

### Best Practices

- **Out-of-source builds**: Keep build artifacts separate from source code (use a dedicated `build/` directory)
- **Parallel compilation**: Use `-j` flag with the number of CPU cores to speed up compilation
- **Debug vs Release**: Choose appropriate build type (`-DCMAKE_BUILD_TYPE=Release` or `Debug`)
- **Documentation**: Always read `README`, `INSTALL`, or `BUILD` files before starting
- **Environment variables**: Be aware of variables like `PATH`, `LD_LIBRARY_PATH`, and `PYTHONPATH`
- **Version control**: Keep track of which version/commit you're building
- **Clean builds**: When in doubt, start with a clean build directory

### Troubleshooting

Common issues when building from source:

- **Missing dependencies**: Install development packages (often ending in `-dev` or `-devel`)
- **Compiler errors**: Check compiler version compatibility
- **Linker errors**: Verify library paths and ensure all dependencies are found
- **Permission issues**: Ensure write access to installation directory or use `sudo` when necessary
- **Configuration failures**: Review configuration logs carefully for hints

### Further Reading and Resources

**Build Systems Documentation**:
- [CMake Official Documentation](https://cmake.org/documentation/) - Comprehensive guide to CMake
- [CMake Tutorial](https://cmake.org/cmake/help/latest/guide/tutorial/index.html) - Step-by-step CMake learning
- [GNU Make Manual](https://www.gnu.org/software/make/manual/) - Complete Make documentation
- [Autotools Mythbuster](https://autotools.info/) - Modern guide to Autotools
- [Meson Build System](https://mesonbuild.com/) - Fast and user-friendly build system

**Compilation and Linking**:
- [GCC Documentation](https://gcc.gnu.org/onlinedocs/) - GNU Compiler Collection manuals
- [Clang Documentation](https://clang.llvm.org/docs/) - LLVM Clang compiler documentation
- [An Introduction to GCC](https://www.network-theory.co.uk/docs/gccintro/) - Free book on GCC
- [Beginner's Guide to Linkers](https://www.lurklurk.org/linkers/linkers.html) - Understanding the linking process

**Package Management and Dependencies**:
- [Spack](https://spack.io/) - Package manager for supercomputers and HPC (widely used in HEP)
- [Conda Documentation](https://docs.conda.io/) - Cross-platform package manager
- [vcpkg](https://vcpkg.io/) - C/C++ library manager by Microsoft

**Best Practices and Guides**:
- [HSF Software Training](https://hsf-training.github.io/hsf-training-cmake-webpage/) - HEP Software Foundation CMake training
- [Modern CMake](https://cliutils.gitlab.io/modern-cmake/) - Best practices for CMake projects
- [Professional CMake](https://crascit.com/professional-cmake/) - Comprehensive CMake book
- [The Architecture of Open Source Applications](https://aosabook.org/) - Learn from real-world build systems

**Debugging Build Issues**:
- [CMake FAQ](https://gitlab.kitware.com/cmake/community/-/wikis/FAQ) - Common CMake questions
- [Stack Overflow - Build Systems](https://stackoverflow.com/questions/tagged/build) - Community Q&A
- [Compiler Explorer (Godbolt)](https://godbolt.org/) - Interactive tool to explore compilation

**HEP-Specific Resources**:
- [HSF Knowledge Base](https://hepsoftwarefoundation.org/knowledge_base.html) - HEP software development resources
- [CERN Software Development](https://ep-dep-sft.web.cern.ch/) - Software tools used at CERN
- [LCG Releases](https://lcginfo.cern.ch/) - CERN library and application releases
