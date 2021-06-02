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

## Upload solution to Producer Portal
**Preview**: This product or feature is covered by the [Pre-GA Offerings Terms](https://cloud.devsite.corp.google.com/terms/service-terms#1) of the Google Cloud Platform Terms of Service. Pre-GA products and features may have limited support, and changes to pre-GA products and features may not be compatible with other pre-GA versions. For more information, see the [launch stage descriptions](https://cloud.devsite.corp.google.com/products#product-launch-stages).

The following commands create and configure a Google Cloud Storage bucket that you can use to upload
your Deployment Manager template to Producer Portal.

```
BUCKET_NAME=<BUCKET_NAME>
gsutil mb $BUCKET_NAME
gsutil versioning set on gs://$BUCKET_NAME
gsutil iam ch "group:marketplace-ops@cloud.google.com:objectViewer" $BUCKET_NAME
```

Prior to running `mpdev apply`, ensure that the `zipFilePath` that you specified
in your [`DeploymentManagerTemplate`](../examples/deployment-manager/autogen/singlevm/configurations.yaml)
resource points to the Google Cloud Storage bucket that you created.

If following an [example](../examples/deployment-manager/autogen/singlevm) in
this repo, you can use the following command to set the `zipFilePath` programmatically:
```
mpdev cfg set <DIR> zipPath gs://$BUCKET_NAME/dm-template.zip
```

After you run `mpdev apply`, upload the deployment package to Producer Portal:

1. Open [Producer Portal](https://console.cloud.google.com/producer-portal) in the Google Cloud Console.
1.  In the list of products, click the name of your product.
1. On the **Overview** page for your product, go to the **Deployment Package**
 section and click **EDIT**.
1. Under **Specify your GCS object location**, select the deployment package
 object that you uploaded to Google Cloud Storage.
[Producer Portal](https://console.cloud.google.com/producer-portal) and select
your product, then follow these steps:

1. On the **Overview** page for your product, go to the **Deployment Package**
 section and click **EDIT**.
1. Under **Specify your GCS object location**, select the deployment package
 object previously uploaded.
