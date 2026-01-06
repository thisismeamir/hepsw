# Introduction

High-Energy Physics has an enormous software stack with dependencies and modules that are interconnected. Yet most of the current documentations on these projects, are isolated and specifically interested in their own build process (which is the right thing to do). Yet for anyone engaged in use or development of these projects it is clear that a complete documentation and build process for the whole stack is necessary. HepSW is the central hub for that.

This repository aims to fill this gap by providing guidelines, best practices, and reproducible build scripts that span the entire HEP software ecosystem. It is designed to help researchers and developers navigate the complexities of building, configuring, and maintaining HEP software stacks from source.
The reason of choosing the source-based approach is to ensure transparency, reproducibility, and adaptability across different Linux distributions and environments.

One of the key goals of HepSW is to manage the complexity of dependencies and version constraints that arise in HEP software. 
For this reason, HepSW provides structured documentation as well as manifest files that illustrates the relationships between various packages, their versions, and compatibility constraints. 

A manifest file is provided for each package, describing its role in the stack, build requirements, dependency relationships, and known compatibility constraints. 
HepSW uses the manifest files to provide a clear guide and steps to build and maintain the software stack.

A manifest file is a YAML document like below (we'll provide more details in the following sections):

```yaml
name: example-package
version: 1.2.3
description: An example package for demonstration purposes.
dependencies:
  - name: dependency-one
    version: ">=2.0.0,<3.0.0"
  - name: dependency-two
    version: "1.5.1"
build:
  system_requirements:
    - cmake
    - gcc >=9.0
  steps:
    - ./configure --prefix=/opt/hep/install/example-package/1.2.3
    - make -j8
    - make install
documentation:
  url: https://example.com/docs/example-package
  notes: |
    This package requires specific versions of its dependencies to function correctly.
```

HepSW does not act as a centralized repository of source code. Instead, it retrieves upstream sources directly—typically from official releases or version-controlled repositories such as GitHub—tracking versions explicitly and keeping local build artifacts isolated and clean. This enables the use of current software versions over time, while maintaining traceability and reproducibility across updates.

It is also very important to note that HepSW is not a binary distribution or a package manager. Instead, it is a transparent, source-first reference implementation that emphasizes clarity, reproducibility, and adaptability. Users are encouraged to inspect and understand the build processes, configurations, and dependencies involved in constructing their HEP software environments.

Since the HEP software ecosystem is constantly evolving, growing, and changing, we encourage contributions from the community to keep HepSW up-to-date and relevant. Whether it's adding new packages, updating existing ones, or improving documentation, contributions are welcome to help make HepSW a comprehensive resource for HEP software development and deployment.

HepSW provides a cli toolkit to help with common tasks. You can install it via pip, or just use this repository as a source. 
This cli tool is the main interface of the program and provides commands to interact with the manifests, build processes, and other functionalities.