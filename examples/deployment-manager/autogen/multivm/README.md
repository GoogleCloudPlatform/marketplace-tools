## Prerequisites

This example package generates a Deployment Manager template for a Redis 
multiple VM setup. Your GCP project must contain a Redis VM image to use in the
Deployment Manager template. The following command copies the Redis VM image
from the
[click-to-deploy project](https://github.com/GoogleCloudPlatform/click-to-deploy)
to your GCP project.

```
PROJECT_ID=<YOUR_PROJECT_ID>
gcloud compute --project=$PROJECT_ID images create redis --source-image=redis-v20200726 --source-image-project=click-to-deploy-images
```

## Generating Deployment Manager Template

To use this mpdev package, execute the following command which will checkout this
package to a directory named `redis`.

```
mpdev pkg get https://github.com/GoogleCloudPlatform/marketplace-tools.git/examples/deployment-manager/autogen/multivm/ redis
```

`mpdev cfg set` will be used to customize values in the 
configurations.yaml template. Particularly the
[`projectId`](../../../../docs/autogen-reference.md#cloud.deploymentmanager.autogen.ImageSpec)
and [`image`](../../../../docs/autogen-reference.md#cloud.deploymentmanager.autogen.ImageSpec)
values must be set for your particular solution.

```
PROJECT_ID=<PROJECT_ID>
IMAGE=redis
```

After setting the variables above run the following commands to set the
variables in `configurations.yaml`

```bash
mpdev cfg set redis/ projectId $PROJECT_ID
mpdev cfg set redis/ image $IMAGE
mpdev cfg set redis/ defaultReplicas 4
```

Now generate a Deployment Manager template with the following command:

```
mpdev apply -f redis/configurations.yaml
```

The template will be zipped to `redis/template.zip`, the location specified
in the `DeploymentManagerTemplate` resource of `redis/configurations.yaml`.

To verify the template is properly configured, you can create a deployment from
the template with the following commands:

```
TMPDIR=$(mktemp -d)
unzip redis/template.zip -d $TMPDIR
gcloud deployment-manager deployments create redis --config $TMPDIR/test_config.yaml
```

## Further Customization of Deployment Manager Template

When creating your own Deployment Manager solution, you may need to
customize other fields in the autogen specification, such as `passwords`,
`deployInput`, and `postDeploy`. See the 
[autogen reference](../../../../docs/autogen-reference.md)
for explanations of the fields.

## Upload Solution To Partner Portal

See instructions in this 
[guide](../../../../docs/deployment-manager-guide.md).

