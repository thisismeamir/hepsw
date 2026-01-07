# HepSW

HepSW is a source-first, reproducible software framework for building, packaging, and composing High Energy Physics (HEP) software stacks on Linux systems.

It provides a single, explicit description of how HEP software is built—covering compilers, dependencies, build configurations, and environment assumptions—without relying on opaque binary package managers or pre-built images. Instead of distributing frozen states, HepSW records build intent and provenance, allowing software environments to be reconstructed deterministically across different Linux distributions, institutions, and timescales.

HepSW is designed to assist researchers in assembling and evolving their software stacks, rather than requiring them to author low-level build definitions themselves. Users may rely on curated build descriptions and workflows, while HepSW ensures that dependencies remain explicit, versioned, and internally consistent as the stack evolves. When changes occur—such as compiler upgrades or dependency updates—the system is structured to expose breakage early and transparently, rather than hiding it behind implicit state.

For each supported package, HepSW also provides structured documentation describing its role in the stack, build requirements, dependency relationships, and known compatibility constraints. This documentation is treated as part of the software definition itself, ensuring that knowledge about how and why software is built is preserved alongside the build process.

HepSW does not act as a centralized repository of source code. Instead, it retrieves upstream sources directly—typically from official releases or version-controlled repositories such as GitHub—tracking versions explicitly and keeping local build artifacts isolated and clean. This enables the use of current software versions over time, while maintaining traceability and reproducibility across updates.

This documentation serves as a practical guide to using HepSW: understanding the structure of HEP software stacks, building and updating software from source, managing dependencies and version constraints, and assembling robust environments for research and analysis without relying on fragile, ad-hoc setups.

```{toctree}
:maxdepth: 2
:numbered:

00-introduction/index
01-basics/index
02-layout-and-workflow/index
03-dependencies/index
04-build-guides/index
05-environments/index
06-advanced/index
07-troubleshooting/index
08-contribution/index
09-architecture-and-design/index
```
