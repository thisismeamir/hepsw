# Contributing to HepSW

Thank you for your interest in contributing to HepSW! This guide will help you get started with contributing to the project.

## Ways to Contribute

There are several ways you can contribute to HepSW:

1. **Add Package Manifests**: Contribute new package manifests to the HepSW Package Index Repository (HPIR)
2. **Improve Documentation**: Help improve or translate documentation
3. **Report Bugs**: Submit bug reports and issues
4. **Suggest Features**: Propose new features or improvements
5. **Fix Bugs**: Submit pull requests to fix existing issues
6. **Code Contributions**: Contribute to the HepSW core codebase

## Contributing Package Manifests

The most common contribution is adding new package manifests to HPIR.

### Prerequisites

- Familiarity with YAML syntax
- Understanding of the package's build system
- Knowledge of the package's dependencies
- A GitHub account

### Process

1. **Fork the HPIR repository**:
   ```bash
   # Visit https://github.com/thisismeamir/hepsw-package-index/
   # Click "Fork" button
   ```

2. **Clone your fork**:
   ```bash
   git clone https://github.com/YOUR-USERNAME/hepsw-package-index.git
   cd hepsw-package-index
   ```

3. **Create a new branch**:
   ```bash
   git checkout -b add-package-name
   ```

4. **Create the manifest**:
   ```bash
   mkdir -p packages/package-name
   cd packages/package-name
   # Create manifest.yaml
   ```

5. **Write the manifest**:
   ```yaml
   name: package-name
   version: x.y.z
   description: A brief description of the package
   
   source:
     type: tarball
     url: https://example.com/package-name-x.y.z.tar.gz
     checksum: sha256:actual-checksum-here
   
   specifications:
     build:
       toolchain:
         - cmake >=3.15
       dependencies:
         - name: dependency-name
           version: ">=1.0.0"
   
   recipe:
     configure:
       - name: Configure
         type: cmake
         args:
           - -DCMAKE_INSTALL_PREFIX=${INSTALL_PREFIX}
     build:
       - name: Build
         type: cmake
         target: build
         parallel: ${NCORES}
     install:
       - name: Install
         type: cmake
         target: install
   ```

6. **Test your manifest locally**:
   ```bash
   # From your HepSW workspace
   hepsw fetch --third-party /path/to/your/manifest.yaml
   hepsw evaluate package-name
   hepsw walk package-name
   hepsw build package-name
   ```

7. **Commit and push**:
   ```bash
   git add packages/package-name/manifest.yaml
   git commit -m "Add manifest for package-name version x.y.z"
   git push origin add-package-name
   ```

8. **Create a Pull Request**:
    - Go to your fork on GitHub
    - Click "Pull Request"
    - Describe what package you're adding and any relevant notes
    - Submit the PR

### Manifest Guidelines

- **Completeness**: Include all required fields (name, version, description, source)
- **Checksums**: Always include SHA-256 checksums for sources
- **Dependencies**: List all build and runtime dependencies with version constraints
- **Testing**: Test the manifest on at least one platform before submitting
- **Documentation**: Add comments for non-obvious configuration options
- **Versioning**: Use semantic versioning (major.minor.patch)
- **Source URLs**: Use stable, official sources (prefer official repositories or releases)

### Checklist Before Submitting

- [ ] Manifest follows the correct YAML structure
- [ ] `hepsw evaluate` passes without critical errors
- [ ] `hepsw build` successfully builds the package
- [ ] All dependencies are specified correctly
- [ ] Checksums are verified
- [ ] Source URL is accessible and stable
- [ ] Recipe steps are well-documented
- [ ] Commit message is descriptive

## Reporting Issues

### Bug Reports

When reporting bugs, please include:

1. **Description**: Clear description of the issue
2. **Steps to Reproduce**:
   ```bash
   hepsw fetch package-name
   hepsw build package-name --with option
   # Error occurs here
   ```
3. **Expected Behavior**: What you expected to happen
4. **Actual Behavior**: What actually happened
5. **Environment**:
    - OS and version (e.g., Ubuntu 22.04, macOS 13.2)
    - HepSW version: `hepsw --version`
    - Compiler versions: `gcc --version`, `cmake --version`
6. **Build Log**: Attach or link to the relevant log file from `~/.hepsw/logs/`
7. **Manifest**: If the issue is with a specific package, include the manifest

### Feature Requests

When requesting features, please include:

1. **Use Case**: Describe the problem you're trying to solve
2. **Proposed Solution**: How you envision the feature working
3. **Alternatives**: Other solutions you've considered
4. **Examples**: Similar features in other tools (if applicable)

## Code Contributions

### Setting Up Development Environment

```bash
# Clone the repository
git clone https://github.com/thisismeamir/hepsw.git
cd hepsw

# Create a virtual environment (if Python-based)
python3 -m venv venv
source venv/bin/activate

# Install development dependencies
pip install -r requirements-dev.txt

# Run tests
pytest tests/
```

### Development Guidelines

- **Code Style**: Follow the existing code style and conventions
- **Testing**: Add tests for new features or bug fixes
- **Documentation**: Update documentation for user-facing changes
- **Commits**: Write clear, descriptive commit messages
- **Branches**: Create feature branches from `main` or `develop`

### Commit Message Format

```
type: brief description (50 chars or less)

More detailed explanation if needed. Wrap at 72 characters.
Explain what and why, not how.

Fixes #123
```

Types: `feat`, `fix`, `docs`, `test`, `refactor`, `style`, `chore`

### Pull Request Process

1. Update documentation for any user-facing changes
2. Add tests for new functionality
3. Ensure all tests pass
4. Update CHANGELOG.md if applicable
5. Request review from maintainers
6. Address review feedback
7. Once approved, a maintainer will merge your PR

## Community Guidelines

### Code of Conduct

- Be respectful and inclusive
- Welcome newcomers and help them get started
- Provide constructive feedback
- Focus on what's best for the community
- Show empathy towards other community members

### Communication Channels

- **GitHub Issues**: Bug reports, feature requests, and discussions
- **GitHub Discussions**: General questions and community discussions
- **Mailing List**: [To be announced]
- **Slack/Discord**: [To be announced]

## Recognition

Contributors are recognized in:
- `CONTRIBUTORS.md` file
- Release notes for significant contributions
- Package manifest author field for new packages

## Questions?

If you have questions about contributing:
- Check existing documentation
- Search existing issues and discussions
- Open a new discussion on GitHub
- Contact maintainers: [contact information]

## Resources

- [HepSW Documentation](https://hepsw.readthedocs.io/)
- [Package Index Repository](https://github.com/thisismeamir/hepsw-package-index/)
- [Issue Tracker](https://github.com/thisismeamir/hepsw/issues)
- [Manifest Documentation](../01-basics/basics-of-manifests.md)

Thank you for contributing to HepSW!