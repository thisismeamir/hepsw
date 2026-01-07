# Troubleshooting

This guide covers common issues you might encounter when using HepSW and their solutions.

## General Troubleshooting Steps

Before diving into specific issues, try these general debugging steps:

1. **Check HepSW version**:
   ```bash
   hepsw --version
   ```

2. **Verify workspace state**:
   ```bash
   hepsw status
   ```

3. **Update package index**:
   ```bash
   hepsw update-index
   ```

4. **Check logs**:
   ```bash
   ls -lt ~/.hepsw/logs/ | head
   cat ~/.hepsw/logs/latest.log
   ```

5. **Validate configuration**:
   ```bash
   hepsw config --check
   ```

## Installation Issues

### HepSW Command Not Found

**Problem**: Shell cannot find the `hepsw` command.

**Solution**:
```bash
# Check if HepSW is installed
which hepsw

# If not found, add to PATH (add to ~/.bashrc or ~/.zshrc)
export PATH="$HOME/.local/bin:$PATH"

# Or reinstall HepSW
pip install --user hepsw
```

### Permission Denied During Init

**Problem**: Cannot create `~/.hepsw` directory.

**Solution**:
```bash
# Check permissions
ls -ld ~/

# Create directory manually if needed
mkdir -p ~/.hepsw
chmod 755 ~/.hepsw

# Run init again
hepsw init
```

### Index Clone Failed

**Problem**: Cannot clone package index repository.

**Solution**:
```bash
# Check network connectivity
ping github.com

# Try manual clone
git clone https://github.com/thisismeamir/hepsw-package-index.git ~/.hepsw/index

# Or specify alternative index URL
hepsw config set indexURL https://alternative-url.git
hepsw init --force
```

## Package Finding and Fetching Issues

### Package Not Found

**Problem**: `hepsw find package-name` returns no results.

**Solutions**:
```bash
# Update index
hepsw update-index

# Check if searching locally vs remotely
hepsw find --remote package-name

# Search with wildcards
hepsw find "*package*"

# Check package name spelling
hepsw find --remote  # List all packages
```

### Fetch Fails with Network Error

**Problem**: Cannot download source from URL.

**Solutions**:
```bash
# Check network connectivity
wget -q --spider https://example.com/source.tar.gz

# Try with verbose output
hepsw fetch package-name --verbose

# Use alternative mirror if available
hepsw fetch package-name --mirror alternative-url

# Check for proxy settings
echo $http_proxy
echo $https_proxy
```

### Source Checksum Mismatch

**Problem**: Downloaded source checksum doesn't match manifest.

**Solutions**:
```bash
# Try re-fetching
hepsw fetch package-name --force

# Skip checksum verification (not recommended)
hepsw fetch package-name --skip-checksum

# Report to maintainers if persistent
hepsw report-issue package-name --type checksum
```

## Build Issues

### Missing Dependencies

**Problem**: Build fails due to missing dependencies.

**Solutions**:
```bash
# Check what's missing
hepsw evaluate package-name

# Fetch missing dependencies automatically
hepsw evaluate package-name --fix

# Or fetch and build dependencies manually
hepsw fetch dependency-name
hepsw build dependency-name
hepsw build package-name
```

### Compiler Not Found

**Problem**: Required compiler version not available.

**Solutions**:
```bash
# Check current compiler
gcc --version
g++ --version

# Install required compiler (Ubuntu/Debian)
sudo apt update
sudo apt install gcc-11 g++-11

# Install required compiler (macOS)
xcode-select --install
brew install gcc@11

# Use specific compiler
hepsw build package-name --env CC=gcc-11 --env CXX=g++-11
```

### CMake Version Too Old

**Problem**: System CMake version is older than required.

**Solutions**:
```bash
# Check CMake version
cmake --version

# Build newer CMake via HepSW
hepsw fetch cmake
hepsw build cmake
source ~/.hepsw/env/cmake-*.sh

# Or install via system package manager
# Ubuntu 22.04+
sudo apt install cmake

# macOS
brew install cmake
```

### Build Fails with "Command Not Found"

**Problem**: Build recipe references unavailable command.

**Solutions**:
```bash
# Check which command is missing
hepsw walk package-name --verbose

# Install missing tools
# For autotools
sudo apt install autoconf automake libtool

# For Python tools
pip install --user required-tool

# Update PATH if tool is installed
export PATH="/path/to/tool:$PATH"
```

### Out of Memory During Build

**Problem**: Compilation crashes with memory errors.

**Solutions**:
```bash
# Reduce parallel jobs
hepsw build package-name --jobs 2

# Check memory usage
free -h
top

# Add swap space if needed (Linux)
sudo fallocate -l 4G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile

# Close other applications
```

### Out of Disk Space

**Problem**: Build fails due to insufficient disk space.

**Solutions**:
```bash
# Check disk usage
df -h ~/.hepsw

# Clean old builds
hepsw clean --builds
hepsw clean --old-builds

# Remove source archives after extraction
hepsw config set keepArchives false

# Check for large log files
du -sh ~/.hepsw/logs/*
hepsw clean --old-logs
```

### Build Hangs or Freezes

**Problem**: Build process stops responding.

**Solutions**:
```bash
# Kill hung build
ps aux | grep hepsw
kill -9 <process-id>

# Try with timeout
hepsw build package-name --timeout 3600

# Build with verbose output to see where it hangs
hepsw build package-name --verbose

# Check for deadlock in parallel builds
hepsw build package-name --jobs 1
```

### Test Failures

**Problem**: Package builds but tests fail.

**Solutions**:
```bash
# Skip tests and install anyway (if safe)
hepsw build package-name --skip-tests

# Run tests manually for debugging
cd ~/.hepsw/builds/package-name/*/
ctest --verbose
make test

# Check test log
cat ~/.hepsw/logs/package-name-*-build.log | grep -A 50 "test"

# Report test failures if persistent
```

## Environment Issues

### Environment Variables Not Set

**Problem**: After building, package-specific variables aren't available.

**Solutions**:
```bash
# Source the environment script
source ~/.hepsw/env/package-name-*.sh

# Check if script exists
ls ~/.hepsw/env/

# Regenerate environment
hepsw env package-name --regenerate

# Add to shell startup
echo "source ~/.hepsw/env/package-name-*.sh" >> ~/.bashrc
```

### Conflicting Package Versions

**Problem**: Multiple versions of same package cause conflicts.

**Solutions**:
```bash
# List installed versions
hepsw list package-name --all-versions

# Use specific version
source ~/.hepsw/env/package-name-1.2.3.sh

# Remove old version
hepsw remove package-name --version 1.0.0

# Use environments to isolate versions
hepsw env create my-env --with package-name@1.2.3
hepsw env activate my-env
```

### Library Not Found at Runtime

**Problem**: Binary complains about missing shared libraries.

**Solutions**:
```bash
# Check LD_LIBRARY_PATH
echo $LD_LIBRARY_PATH

# Source package environment
source ~/.hepsw/env/package-name-*.sh

# Check library location
find ~/.hepsw/install -name "libname.so"

# Add manually if needed
export LD_LIBRARY_PATH="~/.hepsw/install/package-name/*/lib:$LD_LIBRARY_PATH"

# Use ldd to debug
ldd ~/.hepsw/install/package-name/*/bin/executable
```

## Manifest Issues

### Invalid Manifest Syntax

**Problem**: Manifest has YAML syntax errors.

**Solutions**:
```bash
# Validate manifest
hepsw evaluate --manifest path/to/manifest.yaml

# Check YAML syntax
python3 -c "import yaml; yaml.safe_load(open('manifest.yaml'))"

# Common issues:
# - Incorrect indentation (use spaces, not tabs)
# - Missing colons
# - Unquoted special characters
# - Mismatched brackets
```

### Manifest Recipe Fails

**Problem**: Recipe steps don't execute as expected.

**Solutions**:
```bash
# Simulate recipe
hepsw walk package-name --verbose

# Export recipe as script for debugging
hepsw walk package-name --export-script debug.sh
bash -x debug.sh

# Check variable interpolation
hepsw walk package-name --show-env

# Test individual commands manually
cd ~/.hepsw/sources/package-name/*/src
# Run commands from recipe one by one
```

## Configuration Issues

### Configuration File Corrupted

**Problem**: `hepsw.yaml` has errors or is unreadable.

**Solutions**:
```bash
# Backup current config
cp ~/.hepsw/hepsw.yaml ~/.hepsw/hepsw.yaml.backup

# Reset to defaults
hepsw config --reset

# Or manually fix
vim ~/.hepsw/hepsw.yaml

# Validate config
hepsw config --check
```

### Workspace State Inconsistent

**Problem**: Workspace state doesn't match actual files.

**Solutions**:
```bash
# Resync workspace
hepsw sync

# Rebuild state from scratch
hepsw rebuild-state

# Force clean state
hepsw init --force --clean
```

## Platform-Specific Issues

### macOS: Command Line Tools Not Found

**Problem**: Build fails on macOS with missing headers.

**Solutions**:
```bash
# Install Xcode Command Line Tools
xcode-select --install

# Reset tools if corrupted
sudo rm -rf /Library/Developer/CommandLineTools
xcode-select --install

# Check installation
xcode-select -p
```

### macOS: Library Linking Issues

**Problem**: Builds fail to find system libraries on macOS.

**Solutions**:
```bash
# Set library paths
export LIBRARY_PATH="/usr/local/lib:$LIBRARY_PATH"
export CPATH="/usr/local/include:$CPATH"

# Use Homebrew paths
export LIBRARY_PATH="$(brew --prefix)/lib:$LIBRARY_PATH"
export CPATH="$(brew --prefix)/include:$CPATH"
```

### Linux: glibc Version Mismatch

**Problem**: Binary built on newer system doesn't run on older system.

**Solutions**:
```bash
# Check glibc version
ldd --version

# Build with older toolchain
hepsw build package-name --toolchain gcc@9

# Use static linking if possible
hepsw build package-name --env LDFLAGS="-static"
```

### Windows/WSL: Path Issues

**Problem**: Windows paths cause problems in WSL.

**Solutions**:
```bash
# Use WSL paths only
cd /home/user/.hepsw

# Convert paths if needed
wslpath 'C:\Users\...'

# Avoid mounting issues
# Keep HepSW workspace in WSL filesystem, not Windows drives
```

## Performance Issues

### Slow Package Index Update

**Problem**: `hepsw update-index` takes too long.

**Solutions**:
```bash
# Use shallow clone
hepsw config set shallowClone true
hepsw update-index

# Use faster mirror
hepsw config set indexURL https://mirror-url.git

# Update only changed files
hepsw update-index --quick
```

### Slow Builds

**Problem**: Builds take longer than expected.

**Solutions**:
```bash
# Increase parallelism
hepsw build package-name --jobs $(nproc)

# Use compiler cache
sudo apt install ccache
hepsw build package-name --env CC="ccache gcc" --env CXX="ccache g++"

# Build in RAM (if enough memory)
hepsw config set buildInMemory true

# Use release build type
hepsw build package-name --build-type Release
```

## Getting Help

If you can't resolve an issue:

1. **Search existing issues**: Check the [GitHub issue tracker](https://github.com/thisismeamir/hepsw/issues)

2. **Check documentation**: Review relevant sections in the docs

3. **Collect debug information**:
   ```bash
   hepsw debug-info > debug.txt
   ```

4. **Report the issue** with:
    - HepSW version
    - Operating system and version
    - Full command that failed
    - Complete error message
    - Relevant log file content
    - Output of `hepsw debug-info`

5. **Ask the community**:
    - GitHub Discussions
    - Mailing list
    - Slack/Discord channel

## Common Error Messages

### "Manifest validation failed"
→ Check manifest syntax with `hepsw evaluate`

### "Dependency resolution failed"
→ Run `hepsw evaluate --fix` to fetch missing dependencies

### "Source verification failed"
→ Checksum mismatch; try re-fetching with `--force`

### "Build directory not clean"
→ Use `hepsw build --clean` or `--rebuild`

### "Permission denied"
→ Check file/directory permissions with `ls -la`

### "Command not found in recipe"
→ Install missing tool or check PATH

### "Incompatible toolchain version"
→ Update toolchain or use older package version

### "Workspace state inconsistent"
→ Run `hepsw sync` or `hepsw init --force`

## Preventive Measures

To avoid common issues:

1. Keep HepSW updated.
2. Regularly update index: `hepsw update-index`
3. Validate manifests before building: `hepsw evaluate`
4. Keep enough disk space (10GB+ recommended)
5. Use `hepsw walk` before complex builds
6. Source environment scripts after building
7. Document custom configurations
8. Backup your configuration: `cp ~/.hepsw/hepsw.yaml ~/backup/`
