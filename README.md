# hepsw

**hepsw** is a reproducible, source-first collection of build scripts, documentation, and environment setups for High Energy Physics (HEP) software on **any Linux distribution**.

The goal is simple:

* build everything from source
* avoid assumptions made by binary package managers
* keep the setup understandable, inspectable, and reproducible

This repository is meant to reduce fragmentation and guesswork when working with complex HEP software ecosystems.

---

## What this repository provides

* A **well-defined directory layout** for sources, builds, installs, and environments
* **Documented build requirements** and dependencies
* Source-based **environment setup scripts**
* Guidance for reproducible builds and logging

This is not a binary distribution and not a replacement for existing tools. It is a transparent, source-first reference implementation.

---

## Supported platforms

* Any modern Linux distribution
* Designed to work without relying on distro-specific HEP packages

---

## Directory layout

```text
/opt/hep/
├── toolchains/        # compilers and core build tools
├── sources/           # cloned repositories and source tarballs
├── builds/            # out-of-source build directories
├── install/           # install prefixes
├── env/               # environment setup scripts
└── logs/              # build logs and command history
```

General rules:

* `sources/` are not modified
* `builds/` are disposable
* `install/` contains only build outputs (empty but required for build scripts)
* environments are activated explicitly

---

## Reproducibility & logging

The workflow assumes that:

* all builds are done out-of-source
* configuration and build commands are logged
* changes are documented alongside the build scripts

Reproducibility is treated as a first-class concern, not an afterthought.

---

## Scope and completeness

This repository aims to cover a **complete set of dependencies** commonly required in HEP software stacks.

If you believe something is missing or incomplete:

* **open an issue**
* describe the software and context where it is needed

Completeness improves through use and feedback.

---

## Non-goals

* Providing prebuilt binaries
* Hiding complexity behind automation
* Being tied to a specific experiment or framework

---

## Status

This repository is in its early stage. And I would gladly use some help to build a reliable source of truth for hep community.Since this is an evolving repository. The structure and principles are stable, while build scripts and documentation will continue to improve over time.

