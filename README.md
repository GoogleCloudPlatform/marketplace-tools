# Overview

This repository contains the `mpdev` tool supporting the development of VMs deployable via
[Google Cloud Marketplace](https://console.cloud.google.com/marketplace).

# Installation

## Prerequisites

The following tools must be installed before using `mpdev`.
* [docker](https://docs.docker.com/get-docker/)
* [gsutil](https://cloud.google.com/storage/docs/gsutil_install)
* zip `sudo apt-get install zip`

## Options

### Download latest release

https://github.com/GoogleCloudPlatform/marketplace-tools/releases/latest

### Install with Homebrew

```
brew tap GoogleCloudPlatform/marketplace-tools https://github.com/GoogleCloudPlatform/marketplace-tools.git
brew install mpdev
```

### Build from source

Building from source requires the following:
* [bazel](https://docs.bazel.build/versions/master/install.html)
* [golang](https://golang.org/dl/)

```
make build
```

The mpdev binary will be created at `bazel-bin/mpdev/mpdev_/mpdev`

## Getting Started

See the [mdpev reference](./docs/mpdev-reference.md) 
and [Deployment Manager guide](./docs/deployment-manager-guide.md) documentation.
