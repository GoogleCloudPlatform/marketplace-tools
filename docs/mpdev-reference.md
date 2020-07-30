## Setup

See installation instructions [here](../README.md).

## Mpdev resources

The `mpdev` tool creates and updates artifacts defined by yaml configurations. The `mpdev` resource schema
is inspired by the schema for kubernetes resources, where a resource type is
uniquely specified by a `kind` and `apiVersion`.

Currently, the `mpdev` tool supports the following types of resources:
* [`DeploymentManagerAutogenTemplate`](https://pkg.go.dev/github.com/GoogleCloudPlatform/marketplace-tools/mpdev/internal/apply?tab=doc#DeploymentManagerAutogenTemplate)
* [`DeploymentManagerTemplate`](https://pkg.go.dev/github.com/GoogleCloudPlatform/marketplace-tools/mpdev/internal/apply?tab=doc#DeploymentManagerTemplate).

See the 
[deployment manager guide](./deployment-manager-guide.md) for how to configure
these resources.

## Commands

### Start from preconfigured mpdev template

The `pkg get` command downloads a preconfigured mpdev template. `mpdev pkg` is
a wrapper around 
[`kpt pkg`](https://googlecontainertools.github.io/kpt/reference/pkg/get]).

```bash
mpdev pkg get https://github.com/marketplace-tools.git/examples/deployment-manager/autogen/singlevm dir
```

### Customize an mpdev template

The `cfg set` command can be used to programmatically customize values in a
preconfigured mpdev template.
`mpdev cfg` is a wrapper around
[`kpt cfg`](https://googlecontainertools.github.io/kpt/reference/cfg/set]).

```bash
mpdev cfg set dir projectId YourProject
```

### Generate mpdev resources

The `apply` command creates resources from the mpdev template.

```bash
mpdev apply -f dir/configurations.yaml
```

The `--dry-run` option can be used to validate the schema of configuration files
quickly without creating the `mpdev` resources.

```bash
mpdev apply --dry-run -f dir/configurations.yaml
```
