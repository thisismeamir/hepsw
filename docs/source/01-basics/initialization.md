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
├── install/           # installed binaries and libraries
├── env/               # environment setup scripts
├── logs/              # build logs and command history
├── manifests/         # Collection of manifests that are fetched from index 
├── thirdparty/        # Unofficial/thirdparty/user-driven builds 
├── index.db           # Local copy of HepSW index database for searching. 
└── hepsw.yaml         # global configuration file
```

The `hepsw.yaml` file is the global configuration file for HepSW. It contains default settings and paths that HepSW will use when building and managing packages.
As well as user preferences about HepSW behavior.

The default configuration file looks like below:

```yaml
workspace: /home/kid-a/.hepsw
sources: /home/kid-a/.hepsw/sources
builds: /home/kid-a/.hepsw/builds
installs: /home/kid-a/.hepsw/installs
envs: /home/kid-a/.hepsw/envs
logs: /home/kid-a/.hepsw/logs
toolchains: /home/kid-a/.hepsw/toolchains
manifests: /home/kid-a/.hepsw/manifests
thirdparty: /home/kid-a/.hepsw/thirdparty
indexConfig:
  databaseURL: link to the remote database
  authToken: The token for reading the remote database
  timeout: 5s
  maxRetries: 3
  retryDelay: 1s
  cacheTTL: 1h0m0s
  enableCache: true
  lastSyncId:
    dependencies: 13
    packages: 13
    versions: 32
state:
  packages: []
  sources: []
  environments: []
userConfig:
  verbosity: ""
  parallelBuilds: 4

```

The initialization process essentially does the following:

1. Checks for/ makes ~/.hepsw directory and the subdirectories.
2. Generates the default `hepsw.yaml` configuration file if it does not already exist.
3. Creates a local database and syncs it with the remote database where we keep official software manifests.
4. Checks for the consistency between the configuration file and the state of the workspace (e.g., existing packages, sources, environments).

**Note**: To stay consistent and reliable, the state of the workspace is verified by the configuration file and not the other way around.
This means that if you clone a project in these directories manually, or made a change in the workspace without the cli tool, HepSW will delete those changes and restore the state according to the configuration file on the next command execution.


After initialization, you can start using other HepSW commands to search for packages, fetch sources, build packages, and set up environments.