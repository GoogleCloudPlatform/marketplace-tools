## Generating Deployment Manager templates using the `mpdev` tool

You can use the `mpdev` tool to create a Deployment Manager template file, also known as a _deployment package_, for VM
products. The `mpdev` tool autogenerates Deployment Manager templates from [Autogen specification files](./autogen-reference.md).

To learn how to create a deployment package and submit the package to Google, see the Google Cloud Marketplace documentation for [Creating your deployment package](https://cloud.google.com/marketplace/docs/partners/vm/create-deployment-package).

## Overview

You use the following `mpdev` commands to generate a deployment package:

1. Use the `mpdev pkg get` command to retrieve a [preconfigured Autogen
specification](../examples/deployment-manager/autogen).
1. Use `mpdev cfg set` to customize values in the retrieved Autogen
specification, or edit the specification manually. If editing manually, we
recommend **setting up auto-complete** by following these
[instructions](setup-autocomplete.md)
3. Use `mpdev apply` to generate the Deployment Manager template.

## See also 

To see examples that use `mpdev`, go to the preconfigured Autogen specifications for a [single vm](../examples/deployment-manager/autogen/singlevm/README.md) and [multi vm](../examples/deployment-manager/autogen/multivm/README.md).

