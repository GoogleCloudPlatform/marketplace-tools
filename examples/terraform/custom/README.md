This example meets the Marketplace requirements for a Terraform module

To validate your module, Marketplace will execute:

```
terraform -chdir=examples/marketplace-test init
terraform -chdir=examples/marketplace-test plan -var project=<test-project>
```

The module in examples/marketplace-test must reference your root module (i.e.
the module in the root directory)

For a Marketplace Partner to reuse this module, they must:

1.  Declare "image" as a variable in Producer Portal.
1.  Replace the default value of "image" with the image declared in Producer
    Portal
1.  Zip the module and upload to GCS

```
zip terraform.zip * -r
gsutil cp terraform.zip gs://<YOUR-BUCKET>/<FOLDER>/terraform.zip
```

Additionally, Partners should include a README.md containing instructions on how
to deploy their product.
