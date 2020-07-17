## Prerequisites

This example package generates a deployment manager template for a wordpress GCE
VM. Your GCP project must contain a wordpress VM image to use in the deployment
manager template. The following command copies the wordpress VM image from the
[click-to-deploy project](https://github.com/GoogleCloudPlatform/click-to-deploy)
to your GCP project.

```
PROJECT_ID=<YOUR_PROJECT_ID>
gcloud compute --project=$PROJECT_ID images create wordpress --source-image=wordpress-v20200629 --source-image-project=click-to-deploy-images
```

## Generating Deployment Manager Template

To use this mpdev package, execute the following command which will checkout this
package to a directory named wordpress.

```
mpdev pkg get https://github.com/marketplace-tools.git/examples/wordpress wordpress
```

`mpdev pkg` and `mpdev cfg` are wrappers around 
[`kpt pkg`](https://googlecontainertools.github.io/kpt/reference/pkg/get/]) and
[`kpt cfg`](https://googlecontainertools.github.io/kpt/reference/cfg/set/])
respectively. `mpdev cfg set` will be used to customize values in the mpdev.yaml
template. Particularly the
[`projectId`](https://github.com/GoogleCloudPlatform/marketplace-tools/tree/gibbley/autogen-docs/docs/autogen-reference.md#cloud.deploymentmanager.autogen.ImageSpec) 
and [`image`](https://github.com/GoogleCloudPlatform/marketplace-tools/tree/gibbley/autogen-docs/docs/autogen-reference.md#cloud.deploymentmanager.autogen.ImageSpec)
values must be set for your particular solution.

```
PROJECT_ID=<PROJECT_ID>
IMAGE=wordpress
```

After setting the variables above run the following commands to set the
variables in `mpdev.yaml`

```bash
mpdev cfg set wordpress projectId $PROJECT_ID
mpdev cfg set wordpress image $IMAGE
```

Now generate a deployment manager template with the following command:

```
mpdev apply -f wordpress/mpdev.yaml
```

## Upload Solution To partner portal

Open (Partner Portal)[https://pantheon.corp.google.com/partner/solutions] and 
select your solution from the list of solutions, then follow these steps:

1. Next to **Deployment Package**, click **Edit**.
1. Select to **Upload a Package**, and then click **Continue**
1. Select the deployment package created by mpdev for **Upload a package**,
uncheck the **Metadata selection** box and then click **Continue**.
1. Click **Save** to save changes.

**Warning:** Unchecking the **Metadata selection** box is crucial, so that
mpdev does not override the solution metadata configured in earlier steps of
the partner portal wizard.