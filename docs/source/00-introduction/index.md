# Introduction

High-Energy Physics has an enormous software stack with dependencies and modules that are interconnected. Yet most of the current documentations on these projects, are isolated and specifically interested in their own build process (which is the right thing to do). Yet for anyone engaged in use or development of these projects it is clear that a complete documentation and build process for the whole stack is necessary. HepSW is the central hub for that.

HEP SW aims to provide guidance, scripts, and environments to set-up high-energy physics stack (including stacks like key4hep) on linux distributions. This is done with some assumptions and ideological choices that in my opinion would make the experience more customizable and yet more powerful.

The goal of [HepSW](https://github.com/thisismeamir/hepsw) is to provide 

- Building guidance and scripts for everything (including base libraries) from scratch
- eep the set-up undestandable, inspectable, and reproducible
- Avoid assumptions made by binary package managers

[HepSW](https://github.com/thisismeamir/hepsw) would provide:

- A well-defined directory layout for sources, builds, installs, and environments
- Source-based environment setup scripts, installation scripts and building process guidance
- Documentation for the build processes of hep stacks

The goal is **not** to:

- Providing prebuilt binaries (until the repository becomes very mature and reliable)
- Hiding complexity behind automation
- Being tied to a specific experiment, framework, linux distro, or package-manager






