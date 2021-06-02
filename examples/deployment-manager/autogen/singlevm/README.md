## Prerequisites

This example package generates a Deployment Manager template for a wordpress GCE
VM. Your GCP project must contain a wordpress VM image to use in the Deployment
Manager template. The following command copies the wordpress VM image from the
[click-to-deploy project](https://github.com/GoogleCloudPlatform/click-to-deploy)
to your GCP project.

```
PROJECT_ID=<YOUR_PROJECT_ID>
gcloud compute --project=$PROJECT_ID images create wordpress --source-image=wordpress-v20200629 --source-image-project=click-to-deploy-images
```

## Generating Deployment Manager Template

To use this `mpdev` package, execute the following command which will checkout this
package to a directory named `wordpress`.

```
mpdev pkg get https://github.com/GoogleCloudPlatform/marketplace-tools.git/examples/deployment-manager/autogen/singlevm/ wordpress
```

`mpdev cfg set` will be used to customize values in the 
configurations.yaml template. Particularly the
[`projectId`](../../../../docs/autogen-reference.md#cloud.deploymentmanager.autogen.ImageSpec) 
and [`image`](../../../../docs/autogen-reference.md#cloud.deploymentmanager.autogen.ImageSpec)
values must be set for your particular solution.

```
PROJECT_ID=<PROJECT_ID>
IMAGE=wordpress
```

After setting the variables above run the following commands to set the
variables in `configurations.yaml`

```bash
mpdev cfg set wordpress/ projectId $PROJECT_ID
mpdev cfg set wordpress/ image $IMAGE

# Optionally set deployment package zip output path.
# Set to GCS location (e.g. gs://<BUCKET>/<OBJECT>) if using 
# Producer Portal (Public Preview)
mpdev cfg set wordpress/ zipPath template.zip
```

Now generate a Deployment Manager template with the following command:

```
mpdev apply -f wordpress/configurations.yaml
```

The template will be zipped to `wordpress/template.zip`, the location specified
in the `DeploymentManagerTemplate` resource of `wordpress/configurations.yaml`.

To verify the template is properly configured, you can create a deployment from
the template with the following commands:

```
TMPDIR=$(mktemp -d)
unzip wordpress/template.zip -d $TMPDIR
gcloud deployment-manager deployments create wordpress --config $TMPDIR/test_config.yaml
```

## Further Customization of Deployment Manager Template

When creating your own Deployment Manager solution, you may need to
customize other fields in the autogen specification, such as `passwords`,
`deployInput`, and `postDeploy`. See the 
[autogen reference](../../../../docs/autogen-reference.md)
for explanations of the fields.

## Upload solution to Partner Portal

See instructions in this 
[guide](../../../../docs/deployment-manager-guide.md#upload-solution-to-partner-portal).

## Upload solution to Producer Portal
**Preview**: This product or feature is covered by the [Pre-GA Offerings Terms](https://cloud.devsite.corp.google.com/terms/service-terms#1) of the Google Cloud Platform Terms of Service. Pre-GA products and features may have limited support, and changes to pre-GA products and features may not be compatible with other pre-GA versions. For more information, see the [launch stage descriptions](https://cloud.devsite.corp.google.com/products#product-launch-stages).  

See instructions in this
[guide](../../../../docs/deployment-manager-guide.md#upload-solution-to-producer-portal-public-preview)
