# Overview

This repository contains the `mpdev` tool supporting the development of VMs deployable via
[Google Cloud Marketplace](https://console.cloud.google.com/marketplace).

## Installation

### Prerequisites

The following tools must be installed before using `mpdev`.
* [docker](https://docs.docker.com/get-docker/)
* [gsutil](https://cloud.google.com/storage/docs/gsutil_install)

### Download latest release

https://github.com/GoogleCloudPlatform/marketplace-tools/releases/latest

### Install with Homebrew

```
brew tap GoogleCloudPlatform/mpdev https://github.com/GoogleCloudPlatform/mpdev.git
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
and [deployment manager guide](./docs/deployment-manager-guide.md) documentation.
