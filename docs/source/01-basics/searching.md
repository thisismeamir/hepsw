# Searching

You can search for specific packages in HepSW using the `hepsw search` command. This command allows you to find packages by name, description, or other metadata defined in their manifest files. 

There are two types of searches you can perform; one is in the local workspace and the other is in the upstream index.

```bash
hepsw search --local <package-name>
```

and 

```bash
hepsw search --remote <package-name>
```

if the flag is not provided, it will search in the local and remote and show results from both.

For example, to search for a package named "root" in the local workspace, you would run:

```bash
hepsw search --local root
```

Which would return something like:

```text
Found 1 package(s) in local workspace:
- root
  Version: 6.24/06
  Description: An object-oriented framework for large scale data analysis
```
To search for the same package in the upstream index, you would run:

```bash
hepsw search --remote root
```
Which would return something like:

```text
Found 3 package(s) in remote index:
- root
  Version: 6.30.02
  Description: An object-oriented framework for large scale data analysis
- rootpy
    Version: 0.9.5
    Description: Python bindings for ROOT
- root_numpy
    Version: 4.6.2
    Description: NumPy bindings for ROOT
```

if multiple versions of the same package are available, all versions will be shown in the results.

```text
- root
  Versions: 
    - 6.30.02
    - 6.24/06
  Description: An object-oriented framework for large scale data analysis
```

You can also use wildcards to search for packages. For example, to search for all packages that contain "data" in their name or description, you would run:

```bash
hepsw search --remote "*data*"
```
Which would return something like:

```text
Found 2 package(s) in remote index:
- pandas
  Version: 1.3.3
  Description: Data analysis and manipulation tool
- hdf5
    Version: 1.12.0
    Description: Hierarchical Data Format library
```

### Find Options

The `hepsw search` command supports several options to refine your search:
- `--local` or `-l`: Find in the local workspace.
- `--remote` or `-r`: Find in the upstream index.
- `--name <name>` or `-n <name>`: Find by package name.
- `--keyword <keyword>` or `-k <keyword>`: Find by keyword in the package description.
- `--version <version-constraint>` `-v <version-constraint>`: Find for a specific version of a package.
- `--depends-on <name>` or `-d <name>`: Show packages that depends on <name>.
- `--needed-for <name>` or `-a <name>`: Show packages that are needed for <name>.
- `--help`: Display help information about the `hepsw search` command.

Some examples:

```bash
hepsw search root
```
Finds for root in local and remote without any specific version in mind.

```bash
hepsw search --name pythia8 --version ">8.3"
```

Finds for pythia8 in name field specifically with version being bigger than 8.3.

```bash
hepsw search --depends-on root --version ">=6.30"
```

Finds packages that depends on root with versions bigger than 6.30.


