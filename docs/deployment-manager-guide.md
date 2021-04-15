## Generating a Deployment Manager template

The `mpdev` tool can autogenerate Deployment Manager templates from an
[autogen specification](./autogen-reference.md). Follow these steps to create a
Deployment Manager template.

1. Use the `mpdev pkg get` command to retrieve a preconfigured autogen
specification. See an example [here](../examples/deployment-manager/autogen/singlevm/README.md).
1. Use `mpdev cfg set` to customize values in the retrieved autogen
specification, or edit the spec manually.
3. Execute `mpdev apply` to generate the Deployment Manager template.

## Upload solution to Partner Portal

Open [Partner Portal](https://console.cloud.google.com/partner/solutions) and 
select your solution from the list of solutions, then follow these steps:

1. Next to **Deployment Package**, click **Edit**.
1. Select **Upload a Package**, and then click **Continue**
1. Select the deployment package created by `mpdev` for **Upload a package**,
uncheck the **Metadata selection** box and then click **Continue**.
1. Click **Save** to save changes.

**WARNING:** Unchecking the **Metadata selection** box is crucial, so that
`mpdev` does not override the solution metadata configured in earlier steps of
the Partner Portal wizard.

## Upload solution to Producer Portal (Public Preview)

The script below creates and configures a GCS bucket that you can use to upload
your Deployment Manager template to Producer Portal.

```
BUCKET_NAME=<BUCKET_NAME>
gsutil mb $BUCKET_NAME
gsutil versioning set on gs://$BUCKET_NAME
gsutil iam ch "group:marketplace-ops@cloud.google.com:objectViewer" $BUCKET_NAME
```

Prior to running `mpdev apply`, ensure the `zipFilePath` specified
in your [`DeploymentManagerTemplate`](../examples/deployment-manager/autogen/singlevm/configurations.yaml)
resource is pointing to the GCS bucket created above.

After running `mpdev apply`, open
[Producer Portal](https://console.cloud.google.com/producer-portal) and select
your product, then follow these steps:

1. On the **Overview** page for your product, go to the **Deployment Package**
 section and click **EDIT**.
1. Under **Specify your GCS object location**, select the deployment package 
 object previously uploaded. 