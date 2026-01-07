# Fetching

HepSW fetches projects, packages, and libraries from source and aims to build them from scratch, by resolving the dependencies in the same manner (building from source)

Fetch process can be given in three ways:

1. Fetching a simple package which can be done using `hepsw fetch <package-name>`.
2. Fetching a set of packages which are separated by white-space `hepsw fetch <package-1> <package-2> ...`.
3. Fetching an environment which means fetching all the sources specified in that environment, `hepsw fetch -e fccsw`.

Each of these does exactly the same thing but the first does it for a single package, the second for a list of packages, and the last one for a defined set of packages.

We would talk about environment creation in later sections, for now lets focus on a single package fetching.

To fetch a package we write the following command:

```bash
hepsw fetch <package-name>
```

Which results in the following process:

1. HepSW looks in the index repository to find a package with such name.
2. If multiple packages are found, it uses the latest version.
3. It will clone/copy/download the package into `~/.hepsw/sources/<package-name>/<version>/src`
4. It copies the manifest into the package as well, so that the build process can be delayed.
5. It generates `build.yml` in `~/.hepsw/sources/<package-name>/<version>/`

Fetching is a very straight procedure, but there are several flags and options that can be added to the command:

- `--path` or `-p`: Specifies a path to be used instead of the `~/.hepsw/sources/<package-name>/<version>/src`, this is not recommended at all because it will may cause fragmentation. If you chose to do it anyway then `build.yml` still exists in the same directory (`~/.hepsw/sources/<package-name>/<version>`). With an entry `src: path/you/mentioned`
- `--deps-depth <number>` or `-d <number>`: Using this HepSW will also fetch all the dependencies up to the depth specified. This means that if `A->B->C->D` and we want to install `D` with depth set to 2 we'd get `B`, `C`, and `D`. This option allows the users to be faster in set-up, instead of fetching each dependency manually, the default is `0` but it can be set to any number using `hepsw config` which we'll get into.

## Importing Manifests Locally

Assume you wrote a manifest yourself (or want to use a third-party manifest) and want to import and test it locally on your computer. `hepsw fetch` is the gateway of importing unofficial manifests into a workspace.

To do so you must add `--third-party` or `-t` flag and pass on the path of the manifest instead of the name. This way HepSW adds the manifest to the `third-party` directory, and the `build.yml` contains the `third-party: true` tag which then makes HepSW not use it unless strictly specified for future dependency resolution.

A Third party software can be something you're working on, or a project that is not currently supported in our index repository. Make sure to request, or contribute to our index repository if it is an important software for the ecosystem.

- [HepSW Package Index Repository (HPIR)](https://github.com/thisismeamir/hepsw-package-index/)

