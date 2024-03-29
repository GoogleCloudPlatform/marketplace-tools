This example uses the `mpdev` tool to generate a Deployment Manager template for a Wordpress Google Compute Engine virtual machine (VM).

To generate a Deployment Manager template for your VM product, 
see the Google Cloud Marketplace documentation for [Creating your deployment package](https://cloud.google.com/marketplace/docs/partners/vm/create-deployment-package).

## Prerequisites

Your Google Cloud project must contain a Wordpress VM image to use in the Deployment
Manager template. To copy the Wordpress VM image from the
[click-to-deploy project](https://github.com/GoogleCloudPlatform/click-to-deploy)
to your project, run the following command:

```
PROJECT_ID=<YOUR_PROJECT_ID>
gcloud compute --project=$PROJECT_ID images create wordpress --source-image=wordpress-v20200629 --source-image-project=click-to-deploy-images
```

## Generate a Deployment Manager template

To retrieve this example's Autogen specification, execute the following command, which checks out the
specification to a directory named `wordpress`.

```
mpdev pkg get https://github.com/GoogleCloudPlatform/marketplace-tools.git/examples/deployment-manager/autogen/singlevm/ wordpress
```

### Update the Autogen specification

You use the command `mpdev cfg set` to update the configurations.yaml file with your product's
[`projectId`](../../../../docs/autogen-reference.md#cloud.deploymentmanager.autogen.ImageSpec)
and [`image`](../../../../docs/autogen-reference.md#cloud.deploymentmanager.autogen.ImageSpec)
values.

Set the variables for the Google Cloud project and name of the VM image:

```
PROJECT_ID=<PROJECT_ID>
IMAGE=wordpress
```

Next, run the following commands to update the values in `configurations.yaml`:

```bash
mpdev cfg set wordpress/ projectId $PROJECT_ID
mpdev cfg set wordpress/ image $IMAGE
```

**Optional**: For further customizations, manually edit `configurations.yaml`.
Editing `yaml` can be error-prone. We recommend using an IDE and following these
[instructions](../../../../docs/setup-autocomplete.md) to set up auto-complete
and schema validation.

### Generate the Deployment Manager template

To generate a Deployment Manager template, run the following command with the updated
Autogen specification:

```
mpdev apply -f wordpress/configurations.yaml
```

The template is zipped to `wordpress/template.zip`, which is the location specified
in the `DeploymentManagerTemplate` resource of `wordpress/configurations.yaml`.

### Deploy the VM

To verify the template is properly configured, create a deployment from
the template:

```
TMPDIR=$(mktemp -d)
unzip wordpress/template.zip -d $TMPDIR
gcloud deployment-manager deployments create wordpress --config $TMPDIR/test_config.yaml
```

## What's next

For instructions for generating a Deployment Manager template for your
VM product, see [Creating your deployment package](https://cloud.google.com/marketplace/docs/partners/vm/create-deployment-package).