# ccliwrapper

This is a small tool to help you generate wrapper scripts for containerized CLI tools,
i.e., CLI tools that can be started in a Docker or Podman container.

## Advantages

* One additional layer of protection compared to running everything directly on the
  host system
* This is especially useful for package managers like `npm`/`yarn` or `gem`/`bundle`
  because there are often lots of dependencies in software projects, which can run
  arbitrary code already on install
* Support for multiple CLI tools per container image via tool alias mechanism
* Make CLI tools available only once by specifying the image to use and e.g. mount
  options - no need to install anything!
* Upgrading CLI tools is then as easy as specifying a new image tag or pulling the
  latest tag again
* Standard mechanism for shell completion can be used
* Supports (rootfull) Docker and (rootless) Podman currently as engine for containers

## Limitations

* By design, applications in containers can only access those folders that are mounted
  into the container (by default only the current working directory) and ports are
  opened on a virtual network interface.

## How to use

### Option 1: Use prebuilt binary (recommended)

Binaries are automatically built for all releases and can be downloaded from the
release page at GitHub. Extract the artifact suitable for your platform and run
`./ccliwrapper --help` to get some usage hints.

Generate wrapper scripts for example config:

    ./ccliwrapper generate --config=example/ccliwrapper.yaml

The example config also functions as documentation for now.

**Note: only the x64 artifact for Linux is tested for now.**

### Option 2: Compile and run

Check out this repository and run `go build` to compile. Then see Option 1.

### Option 3: Use prebuilt docker image

There is a distroless image for `ccliwrapper` that can be easily used. Of course, this
means that you need to take care that all folders (config file, output folder) needed
are mounted.

For the example config in this repo use (with `podman`):

    podman run --rm -v $PWD/example:/example docker.io/aspettl/ccliwrapper:edge generate --config /example/ccliwrapper.yaml

Note: `ccliwrapper` tries to create folders listed in the config file for mounts (and
prints a warning if it does so). The simple reason is that a `run` command fails when
local folders do not exist. Naturally, this feature does not work properly when
`ccliwrapper` itself runs in a container - because most folders will not exist there.

## Ideas

Not yet implemented:

* Auto-pull images
* Auto-build images with additional layers, e.g. for extra tools or UID switch
* Possibility to specify additional parameters for Docker/Podman run command
* Possibility to specify environment variables that are always passed through to the container
* Native support for shell completion with caching
* Support CLI argument parsing to determine additional files or ports to attach
* Support host network mode
* Add automated tests
