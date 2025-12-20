# How to Use HepSW

HepSW is designed to help you build and manage the packages, dependencies, and environments. 
At this stage we would not provide a fully automated tool to do everything for you. Instead, we provide clear guidelines, steps, best practices, and reproducible build scripts that you can follow to set up things the way you want.

The reason behind this approach is that HEP software stacks are complex and diverse, and a one-size-fits-all solution may not be feasible or desirable.
The goal of HepSW is to empower you with the knowledge and tools to understand and control your software environment, rather than relying on opaque binary distributions or package managers.

## Getting Started

HepSW is has a command-line interface (CLI) toolkit written in Go to help with common tasks. You can install it via different package managers, or just use this repository as a source.

```bash
git clone https://github.com/thisismeamir/hepsw.git
```

```bash
sudo apt-get install hepsw # On Debian/Ubuntu
sudo dnf install hepsw     # On Fedora
sudo pacman -S hepsw       # On Arch Linux
```

Once installed, you can use the `hepsw` command to interact with the HepSW system. 

### Initializing a HepSW Workspace
To get started, you need to initialize a HepSW workspace. This will set up the necessary directory structure and configuration files.

```bash
hepsw init /path/to/hepsw-workspace
cd /path/to/hepsw-workspace
```

This will create a directory structure like below:

```text
/hepsw-workspace/
├── toolchains/        # compilers and core build tools
├── sources/           # cloned repositories and source tarballs
├── builds/            # out-of-source build directories
├── install/           # install prefixes
├── env/               # environment setup scripts
└── logs/              # build logs and command history
```



