# Initialization

The first step in using HepSW is to initialize a workspace.
This workspace will contain all the necessary directories and configuration files needed to manage your HEP software stack.
To initialize a HepSW workspace, use the following command:

```bash
hepsw init
```

This command will create a directory structure like below:

```text
~/.hepsw/
├── toolchains/        # compilers and core build tools
├── sources/           # cloned repositories and source tarballs
├── builds/            # out-of-source build directories
├── install/           # install prefixes
├── env/               # environment setup scripts
├── logs/              # build logs and command history
├── index/             # package manifests and metadata cloned from upstream
├── third-party/       # Unofficial/third-party/user-driven builds 
└── hepsw.yaml         # global configuration file
```

The `hepsw.yaml` file is the global configuration file for HepSW. It contains default settings and paths that HepSW will use when building and managing packages.
As well as user preferences about HepSW behavior.

The default configuration file looks like below:

```yaml
workspace: ~/.hepsw
sourcesDir: sources
buildsDir: builds
installDir: install
envDir: env
logsDir: logs
toolchainsDir: toolchains
indexDir: index

state:
    packages: {}
    sources: {}
    environments: {}

userConfig:
    verbosity: info
    parallelBuilds: 4
```

The initialization process essentially does the following:

1. Checks for/ makes ~/.hepsw directory and the subdirectories.
2. Generates the default `hepsw.yaml` configuration file if it does not already exist.
3. Clones the package index from the upstream repository into the `index/` directory if it is not already present.
4. If the package index is already present, it will attempt to update it to the latest version.
5. Refreshes user change in index since it should not be overwritten.
6. Checks for the consistency between the configuration file and the state of the workspace (e.g., existing packages, sources, environments).

**Note**: To stay consistent and reliable, the state of the workspace is verified by the configuration file and not the other way around.
This means that if you clone a project in these directories manually, or made a change in the workspace without the cli tool, HepSW will delete those changes and restore the state according to the configuration file on the next command execution.


After initialization, you can start using other HepSW commands to search for packages, fetch sources, build packages, and set up environments.